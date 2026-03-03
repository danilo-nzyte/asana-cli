package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danilodrobac/asana-cli/internal/models"
)

func TestCustomFieldsAPI_Create(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "cf1",
				"name": "Priority",
				"type": "enum",
			},
		})
	})
	defer server.Close()

	api := NewCustomFieldsAPI(c)
	cf, err := api.Create(&models.CustomFieldCreateRequest{
		Name:      "Priority",
		Workspace: "ws123",
		Type:      "enum",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cf.GID != "cf1" {
		t.Errorf("expected GID 'cf1', got %q", cf.GID)
	}
}

func TestCustomFieldsAPI_List(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces/ws123/custom_fields" {
			t.Errorf("expected /workspaces/ws123/custom_fields, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "cf1", "name": "Priority", "type": "enum"},
				{"gid": "cf2", "name": "Points", "type": "number"},
			},
		})
	})
	defer server.Close()

	api := NewCustomFieldsAPI(c)
	cfs, err := api.List("ws123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfs) != 2 {
		t.Fatalf("expected 2 custom fields, got %d", len(cfs))
	}
}

func TestCustomFieldsAPI_Get(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/custom_fields/cf1" {
			t.Errorf("expected /custom_fields/cf1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "cf1",
				"name": "Priority",
				"type": "enum",
			},
		})
	})
	defer server.Close()

	api := NewCustomFieldsAPI(c)
	cf, err := api.Get("cf1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cf.Name != "Priority" {
		t.Errorf("expected 'Priority', got %q", cf.Name)
	}
}
