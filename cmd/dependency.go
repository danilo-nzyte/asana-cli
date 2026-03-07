package cmd

import (
	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var dependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Manage task dependencies",
}

var dependencyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add dependencies to a task",
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		dependsOn, _ := cmd.Flags().GetStringSlice("depends-on")

		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}
		if len(dependsOn) == 0 {
			output.Fail("validation", "--depends-on is required (at least one GID)", client.ExitUsageError)
		}

		depsAPI := api.NewDependenciesAPI(apiClient)
		err := depsAPI.Add(task, dependsOn)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Dependencies added successfully")
	},
}

var dependencyRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a dependency from a task",
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		dependsOn, _ := cmd.Flags().GetStringSlice("depends-on")

		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}
		if len(dependsOn) == 0 {
			output.Fail("validation", "--depends-on is required", client.ExitUsageError)
		}

		depsAPI := api.NewDependenciesAPI(apiClient)
		err := depsAPI.Remove(task, dependsOn)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Dependencies removed successfully")
	},
}

var dependencyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List dependencies of a task",
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}

		depsAPI := api.NewDependenciesAPI(apiClient)
		deps, err := depsAPI.List(task)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(deps, "")
	},
}

func init() {
	dependencyAddCmd.Flags().String("task", "", "Task GID (required)")
	dependencyAddCmd.Flags().StringSlice("depends-on", nil, "GIDs of tasks this task depends on (required, repeatable)")

	dependencyRemoveCmd.Flags().String("task", "", "Task GID (required)")
	dependencyRemoveCmd.Flags().StringSlice("depends-on", nil, "GIDs of dependencies to remove (required)")

	dependencyListCmd.Flags().String("task", "", "Task GID (required)")

	dependencyCmd.AddCommand(dependencyAddCmd)
	dependencyCmd.AddCommand(dependencyRemoveCmd)
	dependencyCmd.AddCommand(dependencyListCmd)
	rootCmd.AddCommand(dependencyCmd)
}
