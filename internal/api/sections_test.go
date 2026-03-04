package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danilodrobac/asana-cli/internal/models"
)

func TestSectionsCreate(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/projects/proj123/sections" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "sec456",
				"name": "Test Section",
			},
		})
	})
	defer server.Close()

	api := NewSectionsAPI(c)
	section, err := api.Create("proj123", &models.SectionCreateRequest{Name: "Test Section"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if section.GID != "sec456" {
		t.Errorf("expected GID sec456, got %s", section.GID)
	}
	if section.Name != "Test Section" {
		t.Errorf("expected name 'Test Section', got %s", section.Name)
	}
}

func TestSectionsList(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/projects/proj123/sections" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "sec1", "name": "Section 1"},
				{"gid": "sec2", "name": "Section 2"},
			},
		})
	})
	defer server.Close()

	api := NewSectionsAPI(c)
	sections, err := api.List("proj123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sections) != 2 {
		t.Errorf("expected 2 sections, got %d", len(sections))
	}
}

func TestSectionsAddTask(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/sections/sec123/addTask" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{},
		})
	})
	defer server.Close()

	api := NewSectionsAPI(c)
	err := api.AddTask("sec123", "task456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
