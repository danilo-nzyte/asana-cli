package cmd

import (
	"github.com/danilodrobac/asana-cli/internal/api"
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/models"
	"github.com/danilodrobac/asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "Manage Asana portfolios",
}

var portfolioCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new portfolio",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")

		if name == "" {
			output.Fail("validation", "--name is required", client.ExitUsageError)
		}
		if workspaceID == "" {
			output.Fail("validation", "--workspace or ASANA_WORKSPACE_ID is required", client.ExitUsageError)
		}

		req := &models.PortfolioCreateRequest{
			Name:      name,
			Workspace: workspaceID,
			Color:     color,
		}

		portfoliosAPI := api.NewPortfoliosAPI(apiClient)
		portfolio, err := portfoliosAPI.Create(req)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(portfolio, "Portfolio created successfully")
	},
}

var portfolioGetCmd = &cobra.Command{
	Use:   "get [GID]",
	Short: "Get a portfolio by GID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portfoliosAPI := api.NewPortfoliosAPI(apiClient)
		portfolio, err := portfoliosAPI.Get(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(portfolio, "")
	},
}

var portfolioListCmd = &cobra.Command{
	Use:   "list",
	Short: "List portfolios",
	Run: func(cmd *cobra.Command, args []string) {
		owner, _ := cmd.Flags().GetString("owner")

		portfoliosAPI := api.NewPortfoliosAPI(apiClient)
		portfolios, err := portfoliosAPI.List(workspaceID, owner)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(portfolios, "")
	},
}

var portfolioUpdateCmd = &cobra.Command{
	Use:   "update [GID]",
	Short: "Update a portfolio",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := &models.PortfolioUpdateRequest{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			req.Name = &v
		}
		if cmd.Flags().Changed("color") {
			v, _ := cmd.Flags().GetString("color")
			req.Color = &v
		}

		portfoliosAPI := api.NewPortfoliosAPI(apiClient)
		portfolio, err := portfoliosAPI.Update(args[0], req)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(portfolio, "Portfolio updated successfully")
	},
}

var portfolioDeleteCmd = &cobra.Command{
	Use:   "delete [GID]",
	Short: "Delete a portfolio",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portfoliosAPI := api.NewPortfoliosAPI(apiClient)
		err := portfoliosAPI.Delete(args[0])
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Portfolio deleted successfully")
	},
}

var portfolioAddItemCmd = &cobra.Command{
	Use:   "add-item [GID]",
	Short: "Add a project to a portfolio",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		item, _ := cmd.Flags().GetString("item")
		if item == "" {
			output.Fail("validation", "--item is required", client.ExitUsageError)
		}

		portfoliosAPI := api.NewPortfoliosAPI(apiClient)
		err := portfoliosAPI.AddItem(args[0], item)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Item added to portfolio")
	},
}

var portfolioRemoveItemCmd = &cobra.Command{
	Use:   "remove-item [GID]",
	Short: "Remove a project from a portfolio",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		item, _ := cmd.Flags().GetString("item")
		if item == "" {
			output.Fail("validation", "--item is required", client.ExitUsageError)
		}

		portfoliosAPI := api.NewPortfoliosAPI(apiClient)
		err := portfoliosAPI.RemoveItem(args[0], item)
		if err != nil {
			handleAPIError(err)
		}
		output.Success(nil, "Item removed from portfolio")
	},
}

func init() {
	portfolioCreateCmd.Flags().String("name", "", "Portfolio name (required)")
	portfolioCreateCmd.Flags().String("color", "", "Portfolio color")

	portfolioListCmd.Flags().String("owner", "", "Filter by owner GID")

	portfolioUpdateCmd.Flags().String("name", "", "New portfolio name")
	portfolioUpdateCmd.Flags().String("color", "", "New portfolio color")

	portfolioAddItemCmd.Flags().String("item", "", "Project GID to add (required)")
	portfolioRemoveItemCmd.Flags().String("item", "", "Project GID to remove (required)")

	portfolioCmd.AddCommand(portfolioCreateCmd)
	portfolioCmd.AddCommand(portfolioGetCmd)
	portfolioCmd.AddCommand(portfolioListCmd)
	portfolioCmd.AddCommand(portfolioUpdateCmd)
	portfolioCmd.AddCommand(portfolioDeleteCmd)
	portfolioCmd.AddCommand(portfolioAddItemCmd)
	portfolioCmd.AddCommand(portfolioRemoveItemCmd)
	rootCmd.AddCommand(portfolioCmd)
}
