/*
Copyright 2020 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package notifier

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type postOptions struct {
	proxy             string
	tlsConfig         *tls.Config
	requestModifier   func(*retryablehttp.Request)
	responseValidator func(statusCode int, body []byte) error
}

type postOption func(*postOptions)

func postMessage(ctx context.Context, address string, payload interface{}, opts ...postOption) error {
	options := &postOptions{
		// Default validateResponse function verifies that the response status code is 200, 202 or 201.
		responseValidator: func(statusCode int, body []byte) error {
			if statusCode == http.StatusOK ||
				statusCode == http.StatusAccepted ||
				statusCode == http.StatusCreated {
				return nil
			}

			return fmt.Errorf("request failed with status code %d, %s", statusCode, string(body))
		},
	}

	for _, o := range opts {
		o(options)
	}

	httpClient, err := newHTTPClient(options)
	if err != nil {
		return err
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling notification payload failed: %w", err)
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, address, data)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if options.requestModifier != nil {
		options.requestModifier(req)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := options.responseValidator(resp.StatusCode, body); err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	return nil
}

func withProxy(proxy string) postOption {
	return func(opts *postOptions) {
		opts.proxy = proxy
	}
}

func withTLSConfig(tlsConfig *tls.Config) postOption {
	return func(opts *postOptions) {
		opts.tlsConfig = tlsConfig
	}
}

func withRequestModifier(reqModifier func(*retryablehttp.Request)) postOption {
	return func(opts *postOptions) {
		opts.requestModifier = reqModifier
	}
}

func withResponseValidator(respValidator func(statusCode int, body []byte) error) postOption {
	return func(opts *postOptions) {
		opts.responseValidator = respValidator
	}
}

func newHTTPClient(opts *postOptions) (*retryablehttp.Client, error) {
	httpClient := retryablehttp.NewClient()

	transport := httpClient.HTTPClient.Transport.(*http.Transport)

	if opts.tlsConfig != nil {
		transport.TLSClientConfig = opts.tlsConfig
	}

	if opts.proxy != "" {
		proxyURL, err := url.Parse(opts.proxy)
		if err != nil {
			return nil, fmt.Errorf("unable to parse proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// Disable the timeout for the HTTP client,
	// as we set the provider timeout on the context.
	httpClient.HTTPClient.Timeout = 0

	httpClient.RetryWaitMin = 2 * time.Second
	httpClient.RetryWaitMax = 30 * time.Second
	httpClient.RetryMax = 4
	httpClient.Logger = nil

	return httpClient, nil
}
