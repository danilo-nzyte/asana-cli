package cmd

import (
	"os"

	"github.com/danilodrobac/asana-cli/internal/auth"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	apiClient   *client.Client
	workspaceID string
)

var rootCmd = &cobra.Command{
	Use:   "asana-cli",
	Short: "CLI for managing Asana resources",
	Long:  "A command-line interface for managing Asana projects, tasks, portfolios, custom fields, attachments, comments, and dependencies.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip client init for auth commands
		if cmd.Parent() != nil && cmd.Parent().Use == "auth" {
			return
		}
		if cmd.Use == "auth" {
			return
		}

		apiClient = client.New(auth.GetAccessToken)

		if workspaceID == "" {
			workspaceID = os.Getenv("ASANA_WORKSPACE_ID")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&workspaceID, "workspace", "", "Workspace GID (default: ASANA_WORKSPACE_ID env var)")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(10)
	}
}
