package cmd

import (
	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var attachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage Asana attachments",
}

var attachmentUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file attachment to a task",
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		file, _ := cmd.Flags().GetString("file")

		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}
		if file == "" {
			output.Fail("validation", "--file is required", client.ExitUsageError)
		}

		attachAPI := api.NewAttachmentsAPI(apiClient)
		attachment, err := attachAPI.Upload(task, file)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(attachment, "Attachment uploaded successfully")
	},
}

var attachmentGetCmd = &cobra.Command{
	Use:   "get [GID]",
	Short: "Get an attachment by GID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		attachAPI := api.NewAttachmentsAPI(apiClient)
		attachment, err := attachAPI.Get(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(attachment, "")
	},
}

var attachmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List attachments for a task",
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		if task == "" {
			output.Fail("validation", "--task is required", client.ExitUsageError)
		}

		attachAPI := api.NewAttachmentsAPI(apiClient)
		attachments, err := attachAPI.List(task)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(attachments, "")
	},
}

var attachmentDeleteCmd = &cobra.Command{
	Use:   "delete [GID]",
	Short: "Delete an attachment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		attachAPI := api.NewAttachmentsAPI(apiClient)
		err := attachAPI.Delete(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Attachment deleted successfully")
	},
}

func init() {
	attachmentUploadCmd.Flags().String("task", "", "Task GID (required)")
	attachmentUploadCmd.Flags().String("file", "", "File path to upload (required)")

	attachmentListCmd.Flags().String("task", "", "Task GID (required)")

	attachmentCmd.AddCommand(attachmentUploadCmd)
	attachmentCmd.AddCommand(attachmentGetCmd)
	attachmentCmd.AddCommand(attachmentListCmd)
	attachmentCmd.AddCommand(attachmentDeleteCmd)
	rootCmd.AddCommand(attachmentCmd)
}
