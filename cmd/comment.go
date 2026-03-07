package cmd

import (
	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage Asana comments (stories)",
}

var commentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Add a comment to a task",
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		text, _ := cmd.Flags().GetString("text")

		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}
		if text == "" {
			output.Fail("validation", "--text is required", client.ExitUsageError)
		}

		commentsAPI := api.NewCommentsAPI(apiClient)
		comment, err := commentsAPI.Create(task, &models.CommentCreateRequest{Text: text})
		if err != nil {
			handleAPIError(err)
		}
		output.Success(comment, "Comment created successfully")
	},
}

var commentGetCmd = &cobra.Command{
	Use:   "get [GID]",
	Short: "Get a comment by GID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commentsAPI := api.NewCommentsAPI(apiClient)
		comment, err := commentsAPI.Get(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(comment, "")
	},
}

var commentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List comments for a task",
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}

		commentsAPI := api.NewCommentsAPI(apiClient)
		comments, err := commentsAPI.List(task)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(comments, "")
	},
}

var commentUpdateCmd = &cobra.Command{
	Use:   "update [GID]",
	Short: "Update a comment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := &models.CommentUpdateRequest{}
		if cmd.Flags().Changed("text") {
			v, _ := cmd.Flags().GetString("text")
			req.Text = &v
		}

		commentsAPI := api.NewCommentsAPI(apiClient)
		comment, err := commentsAPI.Update(args[0], req)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(comment, "Comment updated successfully")
	},
}

var commentDeleteCmd = &cobra.Command{
	Use:   "delete [GID]",
	Short: "Delete a comment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commentsAPI := api.NewCommentsAPI(apiClient)
		err := commentsAPI.Delete(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Comment deleted successfully")
	},
}

func init() {
	commentCreateCmd.Flags().String("task", "", "Task GID (required)")
	commentCreateCmd.Flags().String("text", "", "Comment text (required)")

	commentListCmd.Flags().String("task", "", "Task GID (required)")

	commentUpdateCmd.Flags().String("text", "", "New comment text")

	commentCmd.AddCommand(commentCreateCmd)
	commentCmd.AddCommand(commentGetCmd)
	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentUpdateCmd)
	commentCmd.AddCommand(commentDeleteCmd)
	rootCmd.AddCommand(commentCmd)
}
