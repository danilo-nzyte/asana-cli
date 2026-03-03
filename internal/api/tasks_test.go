package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danilodrobac/asana-cli/internal/models"
)

func TestTasksAPI_Create(t *testing.T) {
	var receivedBody map[string]interface{}

	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "99999",
				"name": "New Task",
			},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	task, err := api.Create(&models.TaskCreateRequest{
		Name:     "New Task",
		Projects: []string{"proj123"},
		DueOn:    "2025-06-15",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.GID != "99999" {
		t.Errorf("expected GID '99999', got %q", task.GID)
	}
	if task.Name != "New Task" {
		t.Errorf("expected name 'New Task', got %q", task.Name)
	}
}

func TestTasksAPI_Get(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/99999" {
			t.Errorf("expected /tasks/99999, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":       "99999",
				"name":      "My Task",
				"completed": false,
			},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	task, err := api.Get("99999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Name != "My Task" {
		t.Errorf("expected 'My Task', got %q", task.Name)
	}
	if task.Completed {
		t.Error("expected completed=false")
	}
}

func TestTasksAPI_List(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		proj := r.URL.Query().Get("project")
		if proj != "proj123" {
			t.Errorf("expected project=proj123, got %s", proj)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "1", "name": "Task 1"},
				{"gid": "2", "name": "Task 2"},
				{"gid": "3", "name": "Task 3"},
			},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	tasks, err := api.List("proj123", nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestTasksAPI_Search(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces/ws123/tasks/search" {
			t.Errorf("expected /workspaces/ws123/tasks/search, got %s", r.URL.Path)
		}
		text := r.URL.Query().Get("text")
		if text != "bug fix" {
			t.Errorf("expected text='bug fix', got %q", text)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "5", "name": "Fix login bug"},
			},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	tasks, err := api.Search("ws123", "bug fix", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Name != "Fix login bug" {
		t.Errorf("expected 'Fix login bug', got %q", tasks[0].Name)
	}
}

func TestTasksAPI_Delete(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	if err := api.Delete("99999"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
