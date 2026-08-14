package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestListOwnerAgentRoutesRequestAndTypedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/agent-routes" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/agent-routes")
		}
		if r.URL.RawQuery != "" {
			t.Errorf("unexpected query = %q", r.URL.RawQuery)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer owner-token" {
			t.Errorf("Authorization = %q, want bearer owner token", got)
		}
		if r.Body != nil && r.ContentLength > 0 {
			t.Error("unexpected request body")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"id": "route-capability-2",
					"agent_id": "agent-2",
					"label": "Second agent",
					"expires_at": null,
					"revoked_at": null,
					"route_version": 1,
					"status": "active",
					"created_at": "2025-01-02T00:00:00Z",
					"updated_at": "2025-01-02T00:00:00Z",
					"path": "/a/route-capability-2"
				},
				{
					"id": "route-capability-1",
					"agent_id": "agent-1",
					"label": null,
					"expires_at": "2025-02-01T00:00:00Z",
					"revoked_at": null,
					"route_version": 1,
					"status": "expired",
					"created_at": "2025-01-01T00:00:00Z",
					"updated_at": "2025-01-01T00:00:00Z",
					"path": "/a/route-capability-1",
					"url": "https://view.example.test/a/route-capability-1"
				}
			]
		}`))
	}))
	defer server.Close()

	generated, err := api.NewClientWithResponses(
		server.URL,
		api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer owner-token")
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("NewClientWithResponses() error = %v", err)
	}

	resp, err := generated.ListOwnerAgentRoutesWithResponse(context.Background())
	if err != nil {
		t.Fatalf("ListOwnerAgentRoutesWithResponse() error = %v", err)
	}
	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil || resp.JSON200.Data == nil {
		t.Fatalf("unexpected typed response: status=%d body=%+v", resp.StatusCode(), resp.JSON200)
	}
	routes := *resp.JSON200.Data
	if len(routes) != 2 {
		t.Fatalf("route count = %d, want 2", len(routes))
	}
	if routes[0].Id != "route-capability-2" || routes[0].AgentId != "agent-2" || routes[0].Status != api.Active {
		t.Errorf("unexpected first route: %+v", routes[0])
	}
	if routes[1].Id != "route-capability-1" || routes[1].AgentId != "agent-1" || routes[1].Status != api.Expired {
		t.Errorf("unexpected second route: %+v", routes[1])
	}
	if routes[1].ExpiresAt == nil || routes[1].Url == nil {
		t.Errorf("nullable fields were not decoded: %+v", routes[1])
	}
}

type recordingLogger struct {
	messages []string
}

func (l *recordingLogger) Logf(format string, args ...any) {
	l.messages = append(l.messages, fmt.Sprintf(format, args...))
}

type recordingTracer struct {
	request *http.Request
}

func (t *recordingTracer) TraceRequest(_ context.Context, req *http.Request) {
	t.request = req
}

func (t *recordingTracer) TraceResponse(_ context.Context, _ *http.Response) {}

func TestAgentRouteObservabilityRedactsCapabilitiesAndTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/agent-routes/route-capability-1" {
			t.Errorf("actual path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer owner-token" {
			t.Errorf("actual Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer server.Close()

	logger := &recordingLogger{}
	tracer := &recordingTracer{}
	c, err := NewClient(
		WithAPIKey("owner-token"),
		WithBaseURL(server.URL),
		WithLogger(logger),
		WithTracer(tracer),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	resp, err := c.GetAgentRoute(context.Background(), "route-capability-1")
	if err != nil {
		t.Fatalf("GetAgentRoute() error = %v", err)
	}
	resp.Body.Close()

	logs := strings.Join(logger.messages, "\n")
	if strings.Contains(logs, "route-capability-1") || strings.Contains(logs, "owner-token") {
		t.Fatalf("logs leaked capability or token: %q", logs)
	}
	if !strings.Contains(logs, "/v1/agent-routes/REDACTED") {
		t.Fatalf("logs do not contain redacted route path: %q", logs)
	}
	if tracer.request == nil {
		t.Fatal("tracer did not receive a request")
	}
	if got := tracer.request.URL.Path; got != "/v1/agent-routes/REDACTED" {
		t.Fatalf("traced path = %q, want redacted path", got)
	}
	if got := tracer.request.Header.Get("Authorization"); got != "REDACTED" {
		t.Fatalf("traced Authorization = %q, want REDACTED", got)
	}
}

func TestAgentRouteManagementRequests(t *testing.T) {
	activeRoute := `{
		"id": "route-capability-1",
		"agent_id": "agent-1",
		"label": "Production website",
		"expires_at": null,
		"revoked_at": null,
		"route_version": 1,
		"status": "active",
		"created_at": "2025-01-01T00:00:00Z",
		"updated_at": "2025-01-01T00:00:00Z",
		"path": "/a/route-capability-1",
		"url": "https://view.example.test/a/route-capability-1"
	}`
	revokedRoute := `{
		"id": "route-capability-1",
		"agent_id": "agent-1",
		"label": "Production website",
		"expires_at": null,
		"revoked_at": "2025-01-01T01:00:00Z",
		"route_version": 2,
		"status": "revoked",
		"created_at": "2025-01-01T00:00:00Z",
		"updated_at": "2025-01-01T01:00:00Z",
		"path": "/a/route-capability-1",
		"url": "https://view.example.test/a/route-capability-1"
	}`

	tests := []struct {
		name       string
		method     string
		path       string
		response   string
		invoke     func(context.Context, *Client) (*http.Response, error)
		assertBody func(*testing.T, *http.Response)
	}{
		{
			name:     "list",
			method:   http.MethodGet,
			path:     "/v1/agents/agent-1/routes",
			response: `{"data":[` + activeRoute + `]}`,
			invoke: func(ctx context.Context, c *Client) (*http.Response, error) {
				return c.ListAgentRoutes(ctx, "agent-1")
			},
			assertBody: func(t *testing.T, resp *http.Response) {
				var result api.ListAgentRoutesApiResponseBody
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Fatalf("decode list response: %v", err)
				}
				if result.Data == nil || len(*result.Data) != 1 || (*result.Data)[0].Id != "route-capability-1" {
					t.Fatalf("unexpected list response: %+v", result)
				}
			},
		},
		{
			name:     "view",
			method:   http.MethodGet,
			path:     "/v1/agent-routes/route-capability-1",
			response: `{"data":` + activeRoute + `}`,
			invoke: func(ctx context.Context, c *Client) (*http.Response, error) {
				return c.GetAgentRoute(ctx, "route-capability-1")
			},
			assertBody: func(t *testing.T, resp *http.Response) {
				var result api.GetAgentRouteApiResponseBody
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Fatalf("decode view response: %v", err)
				}
				if result.Data.Id != "route-capability-1" || result.Data.Status != api.Active {
					t.Fatalf("unexpected view response: %+v", result.Data)
				}
			},
		},
		{
			name:     "revoke",
			method:   http.MethodPost,
			path:     "/v1/agent-routes/route-capability-1/revoke",
			response: `{"data":` + revokedRoute + `}`,
			invoke: func(ctx context.Context, c *Client) (*http.Response, error) {
				return c.RevokeAgentRoute(ctx, "route-capability-1")
			},
			assertBody: func(t *testing.T, resp *http.Response) {
				var result api.RevokeAgentRouteApiResponseBody
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Fatalf("decode revoke response: %v", err)
				}
				if result.Data.Status != api.Revoked || result.Data.RevokedAt == nil || result.Data.RouteVersion != 2 {
					t.Fatalf("unexpected revoke response: %+v", result.Data)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("method = %q, want %q", r.Method, tt.method)
				}
				if r.URL.Path != tt.path {
					t.Errorf("path = %q, want %q", r.URL.Path, tt.path)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer owner-token" {
					t.Errorf("Authorization = %q, want bearer owner token", got)
				}
				if r.Body != nil && r.ContentLength > 0 {
					t.Errorf("unexpected request body for %s", tt.name)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			c, err := NewClient(WithAPIKey("owner-token"), WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			resp, err := tt.invoke(context.Background(), c)
			if err != nil {
				t.Fatalf("%s request error = %v", tt.name, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
			}
			tt.assertBody(t, resp)
		})
	}
}
