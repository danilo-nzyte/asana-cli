package cmd

import (
	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Asana projects",
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		team, _ := cmd.Flags().GetString("team")
		notes, _ := cmd.Flags().GetString("notes")

		if name == "" {
			output.Fail("validation", "--name is required", client.ExitUsageError)
		}

		req := &models.ProjectCreateRequest{
			Name:      name,
			Workspace: workspaceID,
			Team:      team,
			Notes:     notes,
		}

		projectsAPI := api.NewProjectsAPI(apiClient)
		project, err := projectsAPI.Create(req)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(project, "Project created successfully")
	},
}

var projectGetCmd = &cobra.Command{
	Use:   "get [GID]",
	Short: "Get a project by GID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectsAPI := api.NewProjectsAPI(apiClient)
		project, err := projectsAPI.Get(args[0])
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(project, "")
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Run: func(cmd *cobra.Command, args []string) {
		team, _ := cmd.Flags().GetString("team")
		archivedFlag, _ := cmd.Flags().GetBool("archived")
		archivedSet := cmd.Flags().Changed("archived")

		var archived *bool
		if archivedSet {
			archived = &archivedFlag
		}

		projectsAPI := api.NewProjectsAPI(apiClient)
		projects, err := projectsAPI.List(workspaceID, team, archived)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(projects, "")
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update [GID]",
	Short: "Update a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := &models.ProjectUpdateRequest{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			req.Name = &v
		}
		if cmd.Flags().Changed("notes") {
			v, _ := cmd.Flags().GetString("notes")
			req.Notes = &v
		}
		if cmd.Flags().Changed("archived") {
			v, _ := cmd.Flags().GetBool("archived")
			req.Archived = &v
		}

		projectsAPI := api.NewProjectsAPI(apiClient)
		project, err := projectsAPI.Update(args[0], req)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(project, "Project updated successfully")
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [GID]",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectsAPI := api.NewProjectsAPI(apiClient)
		err := projectsAPI.Delete(args[0])
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(nil, "Project deleted successfully")
	},
}

func init() {
	projectCreateCmd.Flags().String("name", "", "Project name (required)")
	projectCreateCmd.Flags().String("team", "", "Team GID")
	projectCreateCmd.Flags().String("notes", "", "Project description")

	projectListCmd.Flags().String("team", "", "Filter by team GID")
	projectListCmd.Flags().Bool("archived", false, "Filter by archived status")

	projectUpdateCmd.Flags().String("name", "", "New project name")
	projectUpdateCmd.Flags().String("notes", "", "New project description")
	projectUpdateCmd.Flags().Bool("archived", false, "Set archived status")

	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	rootCmd.AddCommand(projectCmd)
}
