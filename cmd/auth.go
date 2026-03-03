package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danilodrobac/asana-cli/internal/auth"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Asana via OAuth",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := auth.LoadConfig()
		if err != nil {
			output.Fail("auth_config", err.Error(), client.ExitAuthError)
		}

		token, err := auth.RunOAuthFlow(cfg)
		if err != nil {
			output.Fail("oauth_failed", err.Error(), client.ExitAuthError)
		}

		if err := auth.SaveToken(token); err != nil {
			output.Fail("token_save", err.Error(), client.ExitAuthError)
		}

		// Fetch user info to confirm
		c := client.New(func() (string, error) { return token.AccessToken, nil })
		body, err := c.Get("/users/me")
		if err != nil {
			output.Success(map[string]string{"status": "authenticated"}, "Logged in successfully (could not fetch user info)")
			return
		}

		var resp struct {
			Data struct {
				GID   string `json:"gid"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"data"`
		}
		json.Unmarshal(body, &resp)

		output.Success(map[string]string{
			"status": "authenticated",
			"user":   resp.Data.Name,
			"email":  resp.Data.Email,
		}, fmt.Sprintf("Logged in as %s", resp.Data.Name))
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored authentication tokens",
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.ClearToken(); err != nil {
			output.Fail("logout_failed", err.Error(), 1)
		}
		output.Success(nil, "Logged out successfully")
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	Run: func(cmd *cobra.Command, args []string) {
		// Check OAuth token
		token, err := auth.LoadToken()
		if err == nil {
			c := client.New(func() (string, error) { return token.AccessToken, nil })
			body, apiErr := c.Get("/users/me")
			if apiErr == nil {
				var resp struct {
					Data struct {
						GID   string `json:"gid"`
						Name  string `json:"name"`
						Email string `json:"email"`
					} `json:"data"`
				}
				json.Unmarshal(body, &resp)
				output.Success(map[string]interface{}{
					"method":     "oauth",
					"user":       resp.Data.Name,
					"email":      resp.Data.Email,
					"expires_at": token.ExpiresAt,
				}, fmt.Sprintf("Authenticated as %s (OAuth)", resp.Data.Name))
				return
			}
		}

		// Check PAT
		pat := os.Getenv("ASANA_ACCESS_TOKEN")
		if pat != "" {
			c := client.New(func() (string, error) { return pat, nil })
			body, apiErr := c.Get("/users/me")
			if apiErr == nil {
				var resp struct {
					Data struct {
						GID   string `json:"gid"`
						Name  string `json:"name"`
						Email string `json:"email"`
					} `json:"data"`
				}
				json.Unmarshal(body, &resp)
				output.Success(map[string]interface{}{
					"method": "pat",
					"user":   resp.Data.Name,
					"email":  resp.Data.Email,
				}, fmt.Sprintf("Authenticated as %s (PAT)", resp.Data.Name))
				return
			}
		}

		output.Fail("not_authenticated", "Not authenticated. Run 'asana-cli auth login' or set ASANA_ACCESS_TOKEN.", client.ExitAuthError)
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
