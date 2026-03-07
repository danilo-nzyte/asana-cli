package cmd

import (
	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var sectionCmd = &cobra.Command{
	Use:   "section",
	Short: "Manage Asana sections",
}

var sectionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new section in a project",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		project, _ := cmd.Flags().GetString("project")

		if name == "" {
			output.Fail("validation", "--name is required", client.ExitUsageError)
		}
		if project == "" {
			output.Fail("validation", "--project is required", client.ExitUsageError)
		}

		sectionsAPI := api.NewSectionsAPI(apiClient)
		section, err := sectionsAPI.Create(project, &models.SectionCreateRequest{Name: name})
		if err != nil {
			handleAPIError(err)
		}
		output.Success(section, "Section created successfully")
	},
}

var sectionGetCmd = &cobra.Command{
	Use:   "get [GID]",
	Short: "Get a section by GID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sectionsAPI := api.NewSectionsAPI(apiClient)
		section, err := sectionsAPI.Get(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(section, "")
	},
}

var sectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sections in a project",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			output.Fail("validation", "--project is required", client.ExitUsageError)
		}

		sectionsAPI := api.NewSectionsAPI(apiClient)
		sections, err := sectionsAPI.List(project)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(sections, "")
	},
}

var sectionUpdateCmd = &cobra.Command{
	Use:   "update [GID]",
	Short: "Update a section",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := &models.SectionUpdateRequest{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			req.Name = &v
		}

		sectionsAPI := api.NewSectionsAPI(apiClient)
		section, err := sectionsAPI.Update(args[0], req)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(section, "Section updated successfully")
	},
}

var sectionDeleteCmd = &cobra.Command{
	Use:   "delete [GID]",
	Short: "Delete a section",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sectionsAPI := api.NewSectionsAPI(apiClient)
		err := sectionsAPI.Delete(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Section deleted successfully")
	},
}

var sectionAddTaskCmd = &cobra.Command{
	Use:   "add-task",
	Short: "Add a task to a section",
	Run: func(cmd *cobra.Command, args []string) {
		section, _ := cmd.Flags().GetString("section")
		task, _ := cmd.Flags().GetString("task")

		if section == "" {
			output.Fail("validation", "--section is required", client.ExitUsageError)
		}
		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}

		sectionsAPI := api.NewSectionsAPI(apiClient)
		err := sectionsAPI.AddTask(section, task)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Task added to section successfully")
	},
}

func init() {
	sectionCreateCmd.Flags().String("name", "", "Section name (required)")
	sectionCreateCmd.Flags().String("project", "", "Project GID (required)")

	sectionListCmd.Flags().String("project", "", "Project GID (required)")

	sectionUpdateCmd.Flags().String("name", "", "New section name")

	sectionAddTaskCmd.Flags().String("section", "", "Section GID (required)")
	sectionAddTaskCmd.Flags().String("task", "", "Task GID (required)")

	sectionCmd.AddCommand(sectionCreateCmd)
	sectionCmd.AddCommand(sectionGetCmd)
	sectionCmd.AddCommand(sectionListCmd)
	sectionCmd.AddCommand(sectionUpdateCmd)
	sectionCmd.AddCommand(sectionDeleteCmd)
	sectionCmd.AddCommand(sectionAddTaskCmd)
	rootCmd.AddCommand(sectionCmd)
}
