package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mirako-ai/mirako-go/api"
)

func TestCreateAgentRouteRequest(t *testing.T) {
	label := "Production website"
	validitySeconds := int64(86400)

	tests := []struct {
		name string
		body api.CreateAgentRouteJSONRequestBody
		want map[string]any
	}{
		{
			name: "permanent route",
			body: api.CreateAgentRouteJSONRequestBody{},
			want: map[string]any{},
		},
		{
			name: "labeled temporary route",
			body: api.CreateAgentRouteJSONRequestBody{
				Label:           &label,
				ValiditySeconds: &validitySeconds,
			},
			want: map[string]any{
				"label":            "Production website",
				"validity_seconds": float64(86400),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotBody map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
				}
				if r.URL.Path != "/v1/agents/agent-1/routes" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/v1/agents/agent-1/routes")
				}
				if got := r.Header.Get("Authorization"); got != "Bearer owner-token" {
					t.Errorf("Authorization = %q, want bearer owner token", got)
				}
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", got)
				}
				if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
					t.Errorf("decode request body: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{
					"data": {
						"id": "route-capability-1",
						"agent_id": "agent-1",
						"label": null,
						"expires_at": null,
						"revoked_at": null,
						"route_version": 1,
						"status": "active",
						"created_at": "2025-01-01T00:00:00Z",
						"updated_at": "2025-01-01T00:00:00Z",
						"path": "/a/route-capability-1",
						"url": "https://view.example.test/a/route-capability-1"
					}
				}`))
			}))
			defer server.Close()

			c, err := NewClient(
				WithAPIKey("owner-token"),
				WithBaseURL(server.URL),
			)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			resp, err := c.CreateAgentRoute(context.Background(), "agent-1", tt.body)
			if err != nil {
				t.Fatalf("CreateAgentRoute() error = %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
			}
			if !reflect.DeepEqual(gotBody, tt.want) {
				t.Errorf("body = %#v, want %#v", gotBody, tt.want)
			}

			var result api.CreateAgentRouteApiResponseBody
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			route := result.Data
			if route.Id != "route-capability-1" || route.AgentId != "agent-1" || route.Path != "/a/route-capability-1" {
				t.Errorf("unexpected typed response identity: %+v", route)
			}
			if route.Url == nil || *route.Url != "https://view.example.test/a/route-capability-1" {
				t.Errorf("typed response URL = %v", route.Url)
			}
			if route.Status != api.Active || route.RouteVersion != 1 {
				t.Errorf("typed response lifecycle = status %q, version %d", route.Status, route.RouteVersion)
			}
			if route.ExpiresAt != nil || route.RevokedAt != nil || route.Label != nil {
				t.Errorf("permanent route nullable fields = label %v, expires %v, revoked %v", route.Label, route.ExpiresAt, route.RevokedAt)
			}
			if route.CreatedAt.IsZero() || route.UpdatedAt.IsZero() {
				t.Errorf("typed response timestamps = created %v, updated %v", route.CreatedAt, route.UpdatedAt)
			}
		})
	}
}
