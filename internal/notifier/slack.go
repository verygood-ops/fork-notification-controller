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
	"net/url"
	"strings"

	eventv1 "github.com/fluxcd/pkg/apis/event/v1beta1"
	"github.com/hashicorp/go-retryablehttp"
)

// Slack holds the hook URL
type Slack struct {
	URL       string
	ProxyURL  string
	Token     string
	Username  string
	Channel   string
	TLSConfig *tls.Config
}

// SlackPayload holds the channel and attachments
type SlackPayload struct {
	Channel     string            `json:"channel"`
	Username    string            `json:"username"`
	IconUrl     string            `json:"icon_url"`
	IconEmoji   string            `json:"icon_emoji"`
	Text        string            `json:"text,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment holds the markdown message body
type SlackAttachment struct {
	Color      string       `json:"color"`
	AuthorName string       `json:"author_name"`
	Text       string       `json:"text"`
	MrkdwnIn   []string     `json:"mrkdwn_in"`
	Fields     []SlackField `json:"fields"`
}

type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// NewSlack validates the Slack URL and returns a Slack object
func NewSlack(hookURL string, proxyURL string, token string, tlsConfig *tls.Config, username string, channel string) (*Slack, error) {
	_, err := url.ParseRequestURI(hookURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Slack hook URL %s: '%w'", hookURL, err)
	}

	return &Slack{
		Channel:   channel,
		Username:  username,
		URL:       hookURL,
		ProxyURL:  proxyURL,
		Token:     token,
		TLSConfig: tlsConfig,
	}, nil
}

// Post Slack message
func (s *Slack) Post(ctx context.Context, event eventv1.Event) error {
	// Skip Git commit status update event.
	if event.HasMetadata(eventv1.MetaCommitStatusKey, eventv1.MetaCommitStatusUpdateValue) {
		return nil
	}

	payload := SlackPayload{
		Username: s.Username,
	}

	if s.Channel != "" {
		payload.Channel = s.Channel
	}

	if payload.Username == "" {
		payload.Username = event.ReportingController
	}

	color := "good"
	if event.Severity == eventv1.EventSeverityError {
		color = "danger"
	}

	sfields := make([]SlackField, 0, len(event.Metadata))
	for k, v := range event.Metadata {
		sfields = append(sfields, SlackField{k, v, false})
	}

	a := SlackAttachment{
		Color:      color,
		AuthorName: fmt.Sprintf("%s/%s.%s", strings.ToLower(event.InvolvedObject.Kind), event.InvolvedObject.Name, event.InvolvedObject.Namespace),
		Text:       event.Message,
		MrkdwnIn:   []string{"text"},
		Fields:     sfields,
	}

	payload.Attachments = []SlackAttachment{a}

	opts := []postOption{
		withRequestModifier(func(request *retryablehttp.Request) {
			if s.Token != "" {
				request.Header.Add("Authorization", "Bearer "+s.Token)
			}
		}),
	}
	if s.ProxyURL != "" {
		opts = append(opts, withProxy(s.ProxyURL))
	}
	if s.TLSConfig != nil {
		opts = append(opts, withTLSConfig(s.TLSConfig))
	}
	if s.URL == "https://slack.com/api/chat.postMessage" {
		opts = append(opts, withResponseValidator(validateSlackResponse))
	}

	if err := postMessage(ctx, s.URL, payload, opts...); err != nil {
		return fmt.Errorf("postMessage failed: %w", err)
	}

	return nil
}

// validateSlackResponse validates that a chat.postMessage API response is successful.
// chat.postMessage API always returns 200 OK.
// See https://api.slack.com/methods/chat.postMessage.
//
// On the other hand, incoming webhooks return more expressive HTTP status codes.
// See https://api.slack.com/messaging/webhooks#handling_errors.
func validateSlackResponse(_ int, body []byte) error {
	type slackResponse struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
	}

	slackResp := slackResponse{}
	if err := json.Unmarshal(body, &slackResp); err != nil {
		return fmt.Errorf("unable to unmarshal response body: %w", err)
	}

	if slackResp.Ok {
		return nil
	}
	return fmt.Errorf("Slack responded with error: %s", slackResp.Error)
}
