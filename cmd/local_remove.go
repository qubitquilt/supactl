package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/qubitquilt/supactl/internal/local"
	"github.com/spf13/cobra"
)

var localRemoveCmd = &cobra.Command{
	Use:   "remove <project-id>",
	Short: "Remove a local Supabase instance from configuration",
	Long: `Remove a local Supabase instance from the configuration database.

This command will:
  1. Stop the instance if it's running
  2. Remove the project from the local database

Note: This does NOT delete the project directory or Docker images.
To completely remove all data, you'll need to manually delete the directory
and run 'docker system prune' to clean up unused Docker resources.

Example:
  supactl local remove my-project`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]

		// Check Docker requirements
		if err := checkDockerRequirements(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Load database
		db, err := getLocalDatabase()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Get project to ensure it exists
		project, err := db.GetProject(projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintf(os.Stderr, "\nAvailable projects:\n")

			if len(db.Projects) == 0 {
				fmt.Fprintf(os.Stderr, "  (none)\n")
			} else {
				for id := range db.Projects {
					fmt.Fprintf(os.Stderr, "  - %s\n", id)
				}
			}
			os.Exit(1)
		}

		// Confirm removal
		var confirm bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to remove project '%s'?", projectID),
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !confirm {
			fmt.Println("Removal cancelled.")
			return
		}

		// Stop the instance first
		fmt.Printf("Stopping Supabase instance '%s'...\n", projectID)
		if err := local.DockerComposeDown(projectID, project.Directory); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to stop instance: %v\n", err)
			fmt.Fprintf(os.Stderr, "Continuing with removal...\n\n")
		}

		// Remove from database
		if err := db.RemoveProject(projectID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Save database
		if err := local.SaveDatabase(db); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save database: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nProject '%s' has been removed from the configuration.\n", projectID)
		fmt.Printf("\nNote: The project directory has NOT been deleted:\n")
		fmt.Printf("  %s\n", project.Directory)
		fmt.Printf("\nTo completely remove all files, run:\n")
		fmt.Printf("  rm -rf %s\n", project.Directory)
		fmt.Printf("\nTo clean up Docker resources, run:\n")
		fmt.Printf("  docker system prune\n")
	},
}

func init() {
	localCmd.AddCommand(localRemoveCmd)
}
