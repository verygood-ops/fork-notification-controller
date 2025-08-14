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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fluxcd/pkg/apis/event/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestOpsgenie_Post(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var payload OpsgenieAlert
		err = json.Unmarshal(b, &payload)
		require.NoError(t, err)

	}))
	defer ts.Close()

	tests := []struct {
		name  string
		event func() v1beta1.Event
	}{
		{
			name:  "test event",
			event: testEvent,
		},
		{
			name: "test event with empty metadata",
			event: func() v1beta1.Event {
				events := testEvent()
				events.Metadata = nil
				return events
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			opsgenie, err := NewOpsgenie(ts.URL, "", nil, "token")
			require.NoError(t, err)

			err = opsgenie.Post(context.TODO(), tt.event())
			require.NoError(t, err)
		})
	}
}
