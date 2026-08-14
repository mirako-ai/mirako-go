package api

import (
	"strings"
	"testing"
)

func TestCreateAgentRouteSchemaValidation(t *testing.T) {
	spec, err := GetSwagger()
	if err != nil {
		t.Fatalf("GetSwagger() error = %v", err)
	}

	schemaRef, ok := spec.Components.Schemas["CreateAgentRouteInput"]
	if !ok || schemaRef.Value == nil {
		t.Fatal("CreateAgentRouteInput schema is missing")
	}
	schema := schemaRef.Value

	valid := []struct {
		name  string
		value any
	}{
		{name: "permanent route", value: map[string]any{}},
		{name: "label", value: map[string]any{"label": "Production website"}},
		{name: "temporary route", value: map[string]any{"validity_seconds": float64(86400)}},
	}
	for _, tt := range valid {
		t.Run("accepts "+tt.name, func(t *testing.T) {
			if err := schema.VisitJSON(tt.value); err != nil {
				t.Fatalf("VisitJSON() error = %v", err)
			}
		})
	}

	invalid := []struct {
		name  string
		value any
	}{
		{name: "zero validity", value: map[string]any{"validity_seconds": float64(0)}},
		{name: "negative validity", value: map[string]any{"validity_seconds": float64(-1)}},
		{name: "label over 100 characters", value: map[string]any{"label": strings.Repeat("a", 101)}},
		{name: "client supplied route id", value: map[string]any{"route_id": "not-allowed"}},
		{name: "client supplied response id", value: map[string]any{"id": "not-allowed"}},
		{name: "client supplied agent id", value: map[string]any{"agent_id": "not-allowed"}},
		{name: "client supplied user id", value: map[string]any{"user_id": "not-allowed"}},
		{name: "client supplied expiration", value: map[string]any{"expires_at": "2025-01-01T00:00:00Z"}},
		{name: "client supplied revocation", value: map[string]any{"revoked_at": "2025-01-01T00:00:00Z"}},
		{name: "client supplied route version", value: map[string]any{"route_version": float64(2)}},
		{name: "client supplied status", value: map[string]any{"status": "active"}},
		{name: "client supplied creation timestamp", value: map[string]any{"created_at": "2025-01-01T00:00:00Z"}},
		{name: "client supplied update timestamp", value: map[string]any{"updated_at": "2025-01-01T00:00:00Z"}},
	}
	for _, tt := range invalid {
		t.Run("rejects "+tt.name, func(t *testing.T) {
			if err := schema.VisitJSON(tt.value); err == nil {
				t.Fatal("VisitJSON() error = nil, want validation error")
			}
		})
	}
}
