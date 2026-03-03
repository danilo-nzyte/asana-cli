package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestDependenciesAPI_Add(t *testing.T) {
	var receivedBody map[string]interface{}

	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tasks/task1/addDependencies" {
			t.Errorf("expected /tasks/task1/addDependencies, got %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewDependenciesAPI(c)
	err := api.Add("task1", []string{"dep1", "dep2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDependenciesAPI_List(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/task1/dependencies" {
			t.Errorf("expected /tasks/task1/dependencies, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "dep1", "name": "Dependency 1"},
			},
		})
	})
	defer server.Close()

	api := NewDependenciesAPI(c)
	deps, err := api.List("task1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	if deps[0].Name != "Dependency 1" {
		t.Errorf("expected 'Dependency 1', got %q", deps[0].Name)
	}
}

func TestDependenciesAPI_Remove(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tasks/task1/removeDependencies" {
			t.Errorf("expected /tasks/task1/removeDependencies, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewDependenciesAPI(c)
	err := api.Remove("task1", []string{"dep1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
