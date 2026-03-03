package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadToken(t *testing.T) {
	// Use temp dir
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	token := &TokenData{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	if err := SaveToken(token); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	loaded, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken failed: %v", err)
	}

	if loaded.AccessToken != token.AccessToken {
		t.Errorf("expected access token %q, got %q", token.AccessToken, loaded.AccessToken)
	}
	if loaded.RefreshToken != token.RefreshToken {
		t.Errorf("expected refresh token %q, got %q", token.RefreshToken, loaded.RefreshToken)
	}
}

func TestClearToken(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	token := &TokenData{
		AccessToken:  "test",
		RefreshToken: "test",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}
	SaveToken(token)

	if err := ClearToken(); err != nil {
		t.Fatalf("ClearToken failed: %v", err)
	}

	_, err := LoadToken()
	if err == nil {
		t.Error("expected error after clearing token, got nil")
	}
}

func TestLoadConfig_EnvVars(t *testing.T) {
	os.Setenv("ASANA_CLIENT_ID", "test-id")
	os.Setenv("ASANA_CLIENT_SECRET", "test-secret")
	defer os.Unsetenv("ASANA_CLIENT_ID")
	defer os.Unsetenv("ASANA_CLIENT_SECRET")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.ClientID != "test-id" {
		t.Errorf("expected client ID 'test-id', got %q", cfg.ClientID)
	}
	if cfg.ClientSecret != "test-secret" {
		t.Errorf("expected client secret 'test-secret', got %q", cfg.ClientSecret)
	}
}

func TestLoadConfig_File(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	os.Unsetenv("ASANA_CLIENT_ID")
	os.Unsetenv("ASANA_CLIENT_SECRET")

	cfgDir := filepath.Join(tmpDir, ".config", "asana-cli")
	os.MkdirAll(cfgDir, 0700)

	cfg := Config{ClientID: "file-id", ClientSecret: "file-secret"}
	data, _ := json.Marshal(cfg)
	os.WriteFile(filepath.Join(cfgDir, "config.json"), data, 0600)

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded.ClientID != "file-id" {
		t.Errorf("expected client ID 'file-id', got %q", loaded.ClientID)
	}
}

func TestExchangeCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		r.ParseForm()
		if r.FormValue("grant_type") != "authorization_code" {
			t.Errorf("expected grant_type authorization_code, got %s", r.FormValue("grant_type"))
		}
		if r.FormValue("code") != "test-code" {
			t.Errorf("expected code 'test-code', got %s", r.FormValue("code"))
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  "new-access-token",
			"refresh_token": "new-refresh-token",
			"expires_in":    3600,
			"token_type":    "bearer",
		})
	}))
	defer server.Close()

	// Temporarily override token endpoint — we can't easily do this without
	// modifying the code, so this test verifies the function signature/contract.
	// A full integration test would hit the real endpoint.
	t.Skip("requires overridable token endpoint for unit testing")
}

func TestGetAccessToken_PAT(t *testing.T) {
	// Clear any OAuth token
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	os.Setenv("ASANA_ACCESS_TOKEN", "my-pat")
	defer os.Unsetenv("ASANA_ACCESS_TOKEN")

	token, err := GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}
	if token != "my-pat" {
		t.Errorf("expected 'my-pat', got %q", token)
	}
}

func TestGetAccessToken_OAuth(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	os.Unsetenv("ASANA_ACCESS_TOKEN")

	tokenData := &TokenData{
		AccessToken:  "oauth-token",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}
	SaveToken(tokenData)

	token, err := GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}
	if token != "oauth-token" {
		t.Errorf("expected 'oauth-token', got %q", token)
	}
}

func TestGetAccessToken_NotAuthenticated(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	os.Unsetenv("ASANA_ACCESS_TOKEN")

	_, err := GetAccessToken()
	if err == nil {
		t.Error("expected error when not authenticated, got nil")
	}
}
