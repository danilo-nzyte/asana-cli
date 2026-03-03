package cmd

import (
	"encoding/json"

	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage Asana tasks",
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		project, _ := cmd.Flags().GetString("project")
		assignee, _ := cmd.Flags().GetString("assignee")
		dueOn, _ := cmd.Flags().GetString("due-on")
		notes, _ := cmd.Flags().GetString("notes")
		customFieldsStr, _ := cmd.Flags().GetString("custom-fields")

		if name == "" {
			output.Fail("validation", "--name is required", client.ExitUsageError)
		}

		req := &models.TaskCreateRequest{
			Name:      name,
			Assignee:  assignee,
			Notes:     notes,
			DueOn:     dueOn,
			Workspace: workspaceID,
		}

		if project != "" {
			req.Projects = []string{project}
		}

		if customFieldsStr != "" {
			var cf map[string]interface{}
			if err := json.Unmarshal([]byte(customFieldsStr), &cf); err != nil {
				output.Fail("validation", "invalid --custom-fields JSON: "+err.Error(), client.ExitUsageError)
			}
			req.CustomFields = cf
		}

		tasksAPI := api.NewTasksAPI(apiClient)
		task, err := tasksAPI.Create(req)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(task, "Task created successfully")
	},
}

var taskGetCmd = &cobra.Command{
	Use:   "get [GID]",
	Short: "Get a task by GID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tasksAPI := api.NewTasksAPI(apiClient)
		task, err := tasksAPI.Get(args[0])
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(task, "")
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		assignee, _ := cmd.Flags().GetString("assignee")
		completedFlag, _ := cmd.Flags().GetBool("completed")
		completedSet := cmd.Flags().Changed("completed")

		var completed *bool
		if completedSet {
			completed = &completedFlag
		}

		tasksAPI := api.NewTasksAPI(apiClient)
		tasks, err := tasksAPI.List(project, completed, assignee)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(tasks, "")
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update [GID]",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := &models.TaskUpdateRequest{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			req.Name = &v
		}
		if cmd.Flags().Changed("notes") {
			v, _ := cmd.Flags().GetString("notes")
			req.Notes = &v
		}
		if cmd.Flags().Changed("completed") {
			v, _ := cmd.Flags().GetBool("completed")
			req.Completed = &v
		}
		if cmd.Flags().Changed("due-on") {
			v, _ := cmd.Flags().GetString("due-on")
			req.DueOn = &v
		}
		if cmd.Flags().Changed("assignee") {
			v, _ := cmd.Flags().GetString("assignee")
			req.Assignee = &v
		}

		tasksAPI := api.NewTasksAPI(apiClient)
		task, err := tasksAPI.Update(args[0], req)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(task, "Task updated successfully")
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete [GID]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tasksAPI := api.NewTasksAPI(apiClient)
		err := tasksAPI.Delete(args[0])
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(nil, "Task deleted successfully")
	},
}

var taskSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search tasks in workspace",
	Run: func(cmd *cobra.Command, args []string) {
		query, _ := cmd.Flags().GetString("query")
		project, _ := cmd.Flags().GetString("project")
		assignee, _ := cmd.Flags().GetString("assignee")

		if workspaceID == "" {
			output.Fail("validation", "--workspace or ASANA_WORKSPACE_ID is required for search", client.ExitUsageError)
		}

		tasksAPI := api.NewTasksAPI(apiClient)
		tasks, err := tasksAPI.Search(workspaceID, query, project, assignee)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
			}
			output.Fail("unknown", err.Error(), 1)
		}
		output.Success(tasks, "")
	},
}

func init() {
	taskCreateCmd.Flags().String("name", "", "Task name (required)")
	taskCreateCmd.Flags().String("project", "", "Project GID")
	taskCreateCmd.Flags().String("assignee", "", "Assignee (GID or email)")
	taskCreateCmd.Flags().String("due-on", "", "Due date (YYYY-MM-DD)")
	taskCreateCmd.Flags().String("notes", "", "Task description")
	taskCreateCmd.Flags().String("custom-fields", "", "Custom fields as JSON object")

	taskListCmd.Flags().String("project", "", "Project GID (required for list)")
	taskListCmd.Flags().String("assignee", "", "Filter by assignee")
	taskListCmd.Flags().Bool("completed", false, "Include completed tasks")

	taskUpdateCmd.Flags().String("name", "", "New task name")
	taskUpdateCmd.Flags().String("notes", "", "New task description")
	taskUpdateCmd.Flags().Bool("completed", false, "Set completed status")
	taskUpdateCmd.Flags().String("due-on", "", "New due date (YYYY-MM-DD)")
	taskUpdateCmd.Flags().String("assignee", "", "New assignee")

	taskSearchCmd.Flags().String("query", "", "Search query text")
	taskSearchCmd.Flags().String("project", "", "Filter by project GID")
	taskSearchCmd.Flags().String("assignee", "", "Filter by assignee")

	taskCmd.AddCommand(taskCreateCmd)
	taskCmd.AddCommand(taskGetCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskUpdateCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskSearchCmd)
	rootCmd.AddCommand(taskCmd)
}
