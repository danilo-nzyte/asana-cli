package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danilodrobac/asana-cli/internal/models"
)

func TestPortfoliosAPI_Create(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gid":  "port1",
				"name": "My Portfolio",
			},
		})
	})
	defer server.Close()

	api := NewPortfoliosAPI(c)
	portfolio, err := api.Create(&models.PortfolioCreateRequest{
		Name:      "My Portfolio",
		Workspace: "ws123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if portfolio.GID != "port1" {
		t.Errorf("expected GID 'port1', got %q", portfolio.GID)
	}
}

func TestPortfoliosAPI_AddItem(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/portfolios/port1/addItem" {
			t.Errorf("expected /portfolios/port1/addItem, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewPortfoliosAPI(c)
	err := api.AddItem("port1", "proj123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPortfoliosAPI_RemoveItem(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/portfolios/port1/removeItem" {
			t.Errorf("expected /portfolios/port1/removeItem, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	})
	defer server.Close()

	api := NewPortfoliosAPI(c)
	err := api.RemoveItem("port1", "proj123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPortfoliosAPI_List(t *testing.T) {
	c, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"gid": "p1", "name": "Portfolio 1"},
			},
		})
	})
	defer server.Close()

	api := NewPortfoliosAPI(c)
	portfolios, err := api.List("ws123", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(portfolios) != 1 {
		t.Fatalf("expected 1 portfolio, got %d", len(portfolios))
	}
}
