package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danilodrobac/asana-cli/internal/models"
)

func TestProjectsAPI_Create(t *testing.T) {
	var receivedMethod string
	var receivedPath string

	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "12345",
				"name": "Test Project",
			},
		})
	})
	defer server.Close()

	api := NewProjectsAPI(c)
	project, err := api.Create(&models.ProjectCreateRequest{
		Name:      "Test Project",
		Workspace: "ws123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedMethod != "POST" {
		t.Errorf("expected POST, got %s", receivedMethod)
	}
	if receivedPath != "/projects" {
		t.Errorf("expected /projects, got %s", receivedPath)
	}
	if project.GID != "12345" {
		t.Errorf("expected GID '12345', got %q", project.GID)
	}
	if project.Name != "Test Project" {
		t.Errorf("expected name 'Test Project', got %q", project.Name)
	}
}

func TestProjectsAPI_Get(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/12345" {
			t.Errorf("expected /projects/12345, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "12345",
				"name": "My Project",
			},
		})
	})
	defer server.Close()

	api := NewProjectsAPI(c)
	project, err := api.Get("12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Name != "My Project" {
		t.Errorf("expected 'My Project', got %q", project.Name)
	}
}

func TestProjectsAPI_List(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		ws := r.URL.Query().Get("workspace")
		if ws != "ws123" {
			t.Errorf("expected workspace ws123, got %s", ws)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "1", "name": "Project 1"},
				{"gid": "2", "name": "Project 2"},
			},
		})
	})
	defer server.Close()

	api := NewProjectsAPI(c)
	projects, err := api.List("ws123", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
}

func TestProjectsAPI_Delete(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/projects/12345" {
			t.Errorf("expected /projects/12345, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewProjectsAPI(c)
	err := api.Delete("12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
