package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAttachmentsAPI_Get(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/attachments/att1" {
			t.Errorf("expected /attachments/att1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "att1",
				"name": "design.pdf",
				"host": "asana",
			},
		})
	})
	defer server.Close()

	api := NewAttachmentsAPI(c)
	att, err := api.Get("att1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.GID != "att1" {
		t.Errorf("expected GID 'att1', got %q", att.GID)
	}
	if att.Name != "design.pdf" {
		t.Errorf("expected name 'design.pdf', got %q", att.Name)
	}
}

func TestAttachmentsAPI_List(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/task1/attachments" {
			t.Errorf("expected /tasks/task1/attachments, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "att1", "name": "file1.pdf"},
				{"gid": "att2", "name": "file2.png"},
			},
		})
	})
	defer server.Close()

	api := NewAttachmentsAPI(c)
	atts, err := api.List("task1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(atts) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(atts))
	}
}

func TestAttachmentsAPI_Delete(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/attachments/att1" {
			t.Errorf("expected /attachments/att1, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewAttachmentsAPI(c)
	if err := api.Delete("att1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
