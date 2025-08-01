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
	"fmt"
	"net/url"
	"strings"

	eventv1 "github.com/fluxcd/pkg/apis/event/v1beta1"
	"github.com/hashicorp/go-retryablehttp"
)

type Grafana struct {
	URL       string
	Token     string
	ProxyURL  string
	TLSConfig *tls.Config
	Username  string
	Password  string
}

// GraphitePayload represents a Grafana API annotation in Graphite format
type GraphitePayload struct {
	When int64    `json:"when"` //optional unix timestamp (ms)
	Text string   `json:"text"`
	Tags []string `json:"tags,omitempty"`
}

// NewGrafana validates the Grafana URL and returns a Grafana object
func NewGrafana(URL string, proxyURL string, token string, tlsConfig *tls.Config, username string, password string) (*Grafana, error) {
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL %s", URL)
	}

	return &Grafana{
		URL:       URL,
		ProxyURL:  proxyURL,
		Token:     token,
		Username:  username,
		Password:  password,
		TLSConfig: tlsConfig,
	}, nil
}

// Post annotation
func (g *Grafana) Post(ctx context.Context, event eventv1.Event) error {
	// Skip Git commit status update event.
	if event.HasMetadata(eventv1.MetaCommitStatusKey, eventv1.MetaCommitStatusUpdateValue) {
		return nil
	}

	sfields := make([]string, 0, len(event.Metadata))
	// add tag to filter on grafana
	sfields = append(sfields, "flux", event.ReportingController)
	for k, v := range event.Metadata {
		key := strings.ReplaceAll(k, ":", "|")
		value := strings.ReplaceAll(v, ":", "|")
		sfields = append(sfields, fmt.Sprintf("%s: %s", key, value))
	}
	sfields = append(sfields, fmt.Sprintf("kind: %s", event.InvolvedObject.Kind))
	sfields = append(sfields, fmt.Sprintf("name: %s", event.InvolvedObject.Name))
	sfields = append(sfields, fmt.Sprintf("namespace: %s", event.InvolvedObject.Namespace))
	payload := GraphitePayload{
		When: event.Timestamp.Unix(),
		Text: fmt.Sprintf("%s/%s.%s", strings.ToLower(event.InvolvedObject.Kind), event.InvolvedObject.Name, event.InvolvedObject.Namespace),
		Tags: sfields,
	}

	opts := []postOption{
		withRequestModifier(func(req *retryablehttp.Request) {
			if (g.Username != "" && g.Password != "") && g.Token == "" {
				req.Header.Add("Authorization", "Basic "+basicAuth(g.Username, g.Password))
			}
			if g.Token != "" {
				req.Header.Add("Authorization", "Bearer "+g.Token)
			}
		}),
	}
	if g.ProxyURL != "" {
		opts = append(opts, withProxy(g.ProxyURL))
	}
	if g.TLSConfig != nil {
		opts = append(opts, withTLSConfig(g.TLSConfig))
	}

	if err := postMessage(ctx, g.URL, payload, opts...); err != nil {
		return fmt.Errorf("postMessage failed: %w", err)
	}
	return nil
}
