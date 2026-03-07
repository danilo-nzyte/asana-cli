package cmd

import (
	"encoding/json"

	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var customFieldCmd = &cobra.Command{
	Use:   "custom-field",
	Short: "Manage Asana custom fields",
}

var customFieldCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new custom field",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		cfType, _ := cmd.Flags().GetString("type")
		enumOptionsStr, _ := cmd.Flags().GetString("enum-options")

		if name == "" {
			output.Fail("validation", "--name is required", client.ExitUsageError)
		}
		if cfType == "" {
			output.Fail("validation", "--type is required (text, number, enum)", client.ExitUsageError)
		}
		if workspaceID == "" {
			output.Fail("validation", "--workspace or ASANA_WORKSPACE_ID is required", client.ExitUsageError)
		}

		req := &models.CustomFieldCreateRequest{
			Name:      name,
			Workspace: workspaceID,
			Type:      cfType,
		}

		if enumOptionsStr != "" {
			var opts []models.EnumOption
			if err := json.Unmarshal([]byte(enumOptionsStr), &opts); err != nil {
				output.Fail("validation", "invalid --enum-options JSON: "+err.Error(), client.ExitUsageError)
			}
			req.EnumOptions = opts
		}

		cfAPI := api.NewCustomFieldsAPI(apiClient)
		cf, err := cfAPI.Create(req)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(cf, "Custom field created successfully")
	},
}

var customFieldGetCmd = &cobra.Command{
	Use:   "get [GID]",
	Short: "Get a custom field by GID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfAPI := api.NewCustomFieldsAPI(apiClient)
		cf, err := cfAPI.Get(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(cf, "")
	},
}

var customFieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom fields in workspace",
	Run: func(cmd *cobra.Command, args []string) {
		if workspaceID == "" {
			output.Fail("validation", "--workspace or ASANA_WORKSPACE_ID is required", client.ExitUsageError)
		}

		cfAPI := api.NewCustomFieldsAPI(apiClient)
		cfs, err := cfAPI.List(workspaceID)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(cfs, "")
	},
}

var customFieldUpdateCmd = &cobra.Command{
	Use:   "update [GID]",
	Short: "Update a custom field",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := &models.CustomFieldUpdateRequest{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			req.Name = &v
		}

		cfAPI := api.NewCustomFieldsAPI(apiClient)
		cf, err := cfAPI.Update(args[0], req)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(cf, "Custom field updated successfully")
	},
}

var customFieldDeleteCmd = &cobra.Command{
	Use:   "delete [GID]",
	Short: "Delete a custom field",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfAPI := api.NewCustomFieldsAPI(apiClient)
		err := cfAPI.Delete(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Custom field deleted successfully")
	},
}

func init() {
	customFieldCreateCmd.Flags().String("name", "", "Custom field name (required)")
	customFieldCreateCmd.Flags().String("type", "", "Custom field type: text, number, enum (required)")
	customFieldCreateCmd.Flags().String("enum-options", "", "Enum options as JSON array")

	customFieldUpdateCmd.Flags().String("name", "", "New custom field name")

	customFieldCmd.AddCommand(customFieldCreateCmd)
	customFieldCmd.AddCommand(customFieldGetCmd)
	customFieldCmd.AddCommand(customFieldListCmd)
	customFieldCmd.AddCommand(customFieldUpdateCmd)
	customFieldCmd.AddCommand(customFieldDeleteCmd)
	rootCmd.AddCommand(customFieldCmd)
}
