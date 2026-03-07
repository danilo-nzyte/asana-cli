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
	tasks, err := api.List("proj123", nil, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestTasksAPI_List_OptFields(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		optFields := r.URL.Query().Get("opt_fields")
		if optFields != "name,due_on" {
			t.Errorf("expected opt_fields='name,due_on', got %q", optFields)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	_, err := api.List("proj123", nil, "", "name,due_on")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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
	tasks, err := api.Search("ws123", "bug fix", "", "", "")
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

func TestTasksAPI_Search_OptFields(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		optFields := r.URL.Query().Get("opt_fields")
		if optFields != "name,notes" {
			t.Errorf("expected opt_fields='name,notes', got %q", optFields)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	_, err := api.Search("ws123", "", "", "", "name,notes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

func TestTasksAPI_MyTasks(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces/ws123/tasks/search" {
			t.Errorf("expected /workspaces/ws123/tasks/search, got %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("assignee.any") != "user123" {
			t.Errorf("expected assignee.any=user123, got %q", q.Get("assignee.any"))
		}
		if q.Get("is_subtask") != "false" {
			t.Errorf("expected is_subtask=false, got %q", q.Get("is_subtask"))
		}
		if q.Get("completed") != "false" {
			t.Errorf("expected completed=false, got %q", q.Get("completed"))
		}
		if q.Get("sort_by") != "due_on" {
			t.Errorf("expected sort_by=due_on, got %q", q.Get("sort_by"))
		}
		if q.Get("sort_ascending") != "true" {
			t.Errorf("expected sort_ascending=true, got %q", q.Get("sort_ascending"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "1", "name": "Task 1", "due_on": "2026-03-01"},
				{"gid": "2", "name": "Task 2", "due_on": "2026-03-05"},
			},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	tasks, err := api.MyTasks("ws123", "user123", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestTasksAPI_MyTasks_WithProject(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("projects.any") != "proj456" {
			t.Errorf("expected projects.any=proj456, got %q", q.Get("projects.any"))
		}
		if q.Get("assignee.any") != "user123" {
			t.Errorf("expected assignee.any=user123, got %q", q.Get("assignee.any"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	_, err := api.MyTasks("ws123", "user123", "proj456", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTasksAPI_MyTasks_OptFields(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		optFields := r.URL.Query().Get("opt_fields")
		if optFields != "name,due_on,notes" {
			t.Errorf("expected opt_fields='name,due_on,notes', got %q", optFields)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{},
		})
	})
	defer server.Close()

	api := NewTasksAPI(c)
	_, err := api.MyTasks("ws123", "user123", "", "name,due_on,notes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
