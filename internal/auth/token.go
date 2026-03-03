package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	tokenEndpoint = "https://app.asana.com/-/oauth_token"
)

// TokenData holds stored OAuth tokens.
type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Config holds OAuth client credentials.
type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// configDir returns the configuration directory path.
func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "asana-cli")
}

// TokenPath returns the path to the token file.
func TokenPath() string {
	return filepath.Join(configDir(), "token.json")
}

// ConfigPath returns the path to the config file.
func ConfigPath() string {
	return filepath.Join(configDir(), "config.json")
}

// LoadConfig loads OAuth client credentials from config file or env vars.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		ClientID:     os.Getenv("ASANA_CLIENT_ID"),
		ClientSecret: os.Getenv("ASANA_CLIENT_SECRET"),
	}

	// Try config file if env vars not set
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		data, err := os.ReadFile(ConfigPath())
		if err == nil {
			var fileCfg Config
			if json.Unmarshal(data, &fileCfg) == nil {
				if cfg.ClientID == "" {
					cfg.ClientID = fileCfg.ClientID
				}
				if cfg.ClientSecret == "" {
					cfg.ClientSecret = fileCfg.ClientSecret
				}
			}
		}
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("OAuth credentials not configured. Set ASANA_CLIENT_ID and ASANA_CLIENT_SECRET env vars, or create %s", ConfigPath())
	}

	return cfg, nil
}

// SaveToken writes token data to disk.
func SaveToken(token *TokenData) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	return os.WriteFile(TokenPath(), data, 0600)
}

// LoadToken reads token data from disk.
func LoadToken() (*TokenData, error) {
	data, err := os.ReadFile(TokenPath())
	if err != nil {
		return nil, fmt.Errorf("no stored token (run 'asana-cli auth login'): %w", err)
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("corrupt token file: %w", err)
	}

	return &token, nil
}

// ClearToken removes the stored token.
func ClearToken() error {
	err := os.Remove(TokenPath())
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// ExchangeCode exchanges an authorization code for tokens.
func ExchangeCode(cfg *Config, code string, redirectURI string) (*TokenData, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"redirect_uri":  {redirectURI},
		"code":          {code},
	}

	resp, err := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}, nil
}

// RefreshAccessToken uses the refresh token to get a new access token.
func RefreshAccessToken(cfg *Config, token *TokenData) (*TokenData, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"refresh_token": {token.RefreshToken},
	}

	resp, err := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token refresh failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	refreshToken := tokenResp.RefreshToken
	if refreshToken == "" {
		refreshToken = token.RefreshToken
	}

	return &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}, nil
}

// GetAccessToken returns a valid access token, refreshing if needed.
// Falls back to ASANA_ACCESS_TOKEN env var.
func GetAccessToken() (string, error) {
	// Try OAuth token first
	token, err := LoadToken()
	if err == nil {
		// Check if token needs refresh
		if time.Now().After(token.ExpiresAt.Add(-60 * time.Second)) {
			cfg, cfgErr := LoadConfig()
			if cfgErr != nil {
				return "", fmt.Errorf("token expired and cannot refresh: %w", cfgErr)
			}
			newToken, refErr := RefreshAccessToken(cfg, token)
			if refErr != nil {
				return "", fmt.Errorf("failed to refresh token: %w", refErr)
			}
			if err := SaveToken(newToken); err != nil {
				return "", fmt.Errorf("failed to save refreshed token: %w", err)
			}
			return newToken.AccessToken, nil
		}
		return token.AccessToken, nil
	}

	// Fall back to PAT
	pat := os.Getenv("ASANA_ACCESS_TOKEN")
	if pat != "" {
		return pat, nil
	}

	return "", fmt.Errorf("not authenticated. Run 'asana-cli auth login' or set ASANA_ACCESS_TOKEN")
}
