package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	selfupdate "github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

const repo = "danilo-nzyte/asana-cli"

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

		// Re-download skill files for the new version
		tag := "v" + latest.Version()
		updateSkills(tag)
	},
}

func updateSkills(tag string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not determine home directory for skill update: %v\n", err)
		return
	}

	skills := []struct {
		file string
		dir  string
	}{
		{"SKILL.md", filepath.Join(home, ".claude", "skills", "asana")},
		{"WORK-QUEUE.md", filepath.Join(home, ".claude", "skills", "work-queue")},
	}

	for _, s := range skills {
		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/skill/%s", repo, tag, s.file)
		dest := filepath.Join(s.dir, "SKILL.md")

		if err := downloadFile(url, s.dir, dest); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not update skill %s: %v\n", s.file, err)
			continue
		}
		fmt.Printf("Skill updated: %s\n", dest)
	}
}

func downloadFile(url, dir, dest string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	return os.WriteFile(dest, data, 0644)
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
