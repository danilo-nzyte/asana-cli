package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danilodrobac/asana-cli/internal/models"
)

func TestCommentsAPI_Create(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tasks/task123/stories" {
			t.Errorf("expected /tasks/task123/stories, got %s", r.URL.Path)
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "story1",
				"text": "Hello",
				"type": "comment",
			},
		})
	})
	defer server.Close()

	api := NewCommentsAPI(c)
	comment, err := api.Create("task123", &models.CommentCreateRequest{Text: "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.GID != "story1" {
		t.Errorf("expected GID 'story1', got %q", comment.GID)
	}
	if comment.Text != "Hello" {
		t.Errorf("expected text 'Hello', got %q", comment.Text)
	}
}

func TestCommentsAPI_List(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/task123/stories" {
			t.Errorf("expected /tasks/task123/stories, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "s1", "text": "Comment 1", "type": "comment"},
				{"gid": "s2", "text": "Comment 2", "type": "comment"},
			},
		})
	})
	defer server.Close()

	api := NewCommentsAPI(c)
	comments, err := api.List("task123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}

func TestCommentsAPI_Delete(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/stories/story1" {
			t.Errorf("expected /stories/story1, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewCommentsAPI(c)
	if err := api.Delete("story1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
