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

//
// General steps on how to hook up Flux notifications to a Webex space:
// From the Webex App UI:
// - create a Webex space where you want notifications to be sent
// - add the bot email address to the Webex space (see next section)
//
// Register to https://developer.webex.com/, after signing in:
// - create a bot for forwarding FluxCD notifications to a Webex Space (User profile icon|MyWebexApps|Create a New App|Create a Bot)
// - make a note of the bot email address, this email needs to be added to the Webex space
// - generate a bot access token, this is the ID to use in the webex provider manifest token field
// - find the room ID associated to the webex space using https://developer.webex.com/docs/api/v1/rooms/list-rooms
// - this is the ID to use in the webex provider manifest channel field
//

// Webex holds the hook URL
type Webex struct {
	// mandatory: this should be set to the universal webex API server https://webexapis.com/v1/messages
	URL string
	// mandatory: webex room ID, specifies on which webex space notifications must be sent
	RoomId string
	// mandatory: webex bot access token, this access token must be generated after creating a webex bot
	Token string

	// optional: use a proxy as needed
	ProxyURL  string
	TLSConfig *tls.Config
}

// WebexPayload holds the message text
type WebexPayload struct {
	RoomId   string `json:"roomId,omitempty"`
	Markdown string `json:"markdown,omitempty"`
}

// NewWebex validates the Webex URL and returns a Webex object
func NewWebex(hookURL, proxyURL string, tlsConfig *tls.Config, channel string, token string) (*Webex, error) {

	_, err := url.ParseRequestURI(hookURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Webex hook URL %s: '%w'", hookURL, err)
	}

	return &Webex{
		URL:       hookURL,
		ProxyURL:  proxyURL,
		RoomId:    channel,
		Token:     token,
		TLSConfig: tlsConfig,
	}, nil
}

func (s *Webex) CreateMarkdown(event *eventv1.Event) string {
	var b strings.Builder
	emoji := "✅"
	if event.Severity == eventv1.EventSeverityError {
		emoji = "💣"
	}
	fmt.Fprintf(&b, "%s **%s/%s.%s**\n", emoji, strings.ToLower(event.InvolvedObject.Kind), event.InvolvedObject.Name, event.InvolvedObject.Namespace)
	fmt.Fprintf(&b, "%s\n", event.Message)

	if len(event.Metadata) > 0 {
		for k, v := range event.Metadata {
			fmt.Fprintf(&b, ">**%s**: %s\n", k, v)
		}
	}
	return b.String()
}

// Post Webex message
func (s *Webex) Post(ctx context.Context, event eventv1.Event) error {
	// Skip Git commit status update event.
	if event.HasMetadata(eventv1.MetaCommitStatusKey, eventv1.MetaCommitStatusUpdateValue) {
		return nil
	}

	payload := WebexPayload{
		RoomId:   s.RoomId,
		Markdown: s.CreateMarkdown(&event),
	}

	opts := []postOption{
		withRequestModifier(func(req *retryablehttp.Request) {
			req.Header.Add("Authorization", "Bearer "+s.Token)
		}),
	}
	if s.ProxyURL != "" {
		opts = append(opts, withProxy(s.ProxyURL))
	}
	if s.TLSConfig != nil {
		opts = append(opts, withTLSConfig(s.TLSConfig))
	}

	if err := postMessage(ctx, s.URL, payload, opts...); err != nil {
		return fmt.Errorf("postMessage failed: %w", err)
	}

	return nil
}
