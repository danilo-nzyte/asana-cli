package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	c := &Client{
		httpClient:    &http.Client{Timeout: 5 * time.Second},
		tokenProvider: func() (string, error) { return "test-token", nil },
		baseURL:       server.URL,
	}
	return c, server
}

func TestGet_Success(t *testing.T) {
	var receivedAuth string
	var receivedPath string

	c, server := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedPath = r.URL.Path
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"gid": "123", "name": "Test"},
		})
	})
	defer server.Close()

	body, err := c.Get("/tasks/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAuth != "Bearer test-token" {
		t.Errorf("expected auth header 'Bearer test-token', got %q", receivedAuth)
	}
	if receivedPath != "/tasks/123" {
		t.Errorf("expected path '/tasks/123', got %q", receivedPath)
	}

	var resp struct {
		Data struct {
			GID  string `json:"gid"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Data.GID != "123" {
		t.Errorf("expected GID '123', got %q", resp.Data.GID)
	}
}

func TestGet_NotFound(t *testing.T) {
	c, server := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []map[string]string{
				{"message": "project not found"},
			},
		})
	})
	defer server.Close()

	_, err := c.Get("/projects/999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.ExitCode() != ExitNotFound {
		t.Errorf("expected exit code %d, got %d", ExitNotFound, apiErr.ExitCode())
	}
}

func TestPost_Success(t *testing.T) {
	var receivedMethod string

	c, server := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"gid": "456"},
		})
	})
	defer server.Close()

	body, err := c.Post("/tasks", map[string]string{"name": "New Task"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedMethod != "POST" {
		t.Errorf("expected POST, got %s", receivedMethod)
	}

	var resp struct {
		Data struct {
			GID string `json:"gid"`
		} `json:"data"`
	}
	json.Unmarshal(body, &resp)
	if resp.Data.GID != "456" {
		t.Errorf("expected GID '456', got %q", resp.Data.GID)
	}
}

func TestPut_Success(t *testing.T) {
	var receivedMethod string

	c, server := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"gid": "123", "name": "Updated"},
		})
	})
	defer server.Close()

	_, err := c.Put("/tasks/123", map[string]string{"name": "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedMethod != "PUT" {
		t.Errorf("expected PUT, got %s", receivedMethod)
	}
}

func TestDelete_Success(t *testing.T) {
	var receivedMethod string

	c, server := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	_, err := c.Delete("/tasks/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedMethod != "DELETE" {
		t.Errorf("expected DELETE, got %s", receivedMethod)
	}
}

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		status   int
		code     string
		exitCode int
	}{
		{401, "unauthorized", ExitAuthError},
		{403, "forbidden", ExitAuthError},
		{404, "not_found", ExitNotFound},
		{400, "bad_request", ExitValidation},
		{422, "unprocessable", ExitValidation},
		{429, "rate_limited", ExitRateLimited},
		{500, "server_error", ExitServerError},
	}

	for _, tt := range tests {
		code := ErrorCode(tt.status)
		if code != tt.code {
			t.Errorf("ErrorCode(%d) = %q, want %q", tt.status, code, tt.code)
		}
		apiErr := &APIError{StatusCode: tt.status}
		if apiErr.ExitCode() != tt.exitCode {
			t.Errorf("ExitCode for status %d = %d, want %d", tt.status, apiErr.ExitCode(), tt.exitCode)
		}
	}
}
