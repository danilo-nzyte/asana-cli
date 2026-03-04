package cmd

import (
	"fmt"
	"os"

	selfupdate "github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

const repo = "danilodrobac/asana-cli"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update asana-cli to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "dev" {
			fmt.Fprintln(os.Stderr, "Cannot update a dev build. Install a release build or use 'go install'.")
			os.Exit(1)
		}

		latest, found, err := selfupdate.DetectLatest(cmd.Context(), selfupdate.ParseSlug(repo))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
			os.Exit(1)
		}
		if !found {
			fmt.Fprintln(os.Stderr, "No releases found.")
			os.Exit(1)
		}

		if latest.LessOrEqual(Version) {
			fmt.Printf("Already up to date (v%s).\n", Version)
			return
		}

		fmt.Printf("Updating v%s → v%s ...\n", Version, latest.Version())

		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not locate executable: %v\n", err)
			os.Exit(1)
		}

		if err := selfupdate.UpdateTo(cmd.Context(), latest.AssetURL, latest.AssetName, exe); err != nil {
			fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Updated to v%s.\n", latest.Version())
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
