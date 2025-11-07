package cmd

import (
	"fmt"
	"os"

	"github.com/qubitquilt/supactl/internal/local"
	"github.com/spf13/cobra"
)

var localStopCmd = &cobra.Command{
	Use:   "stop <project-id>",
	Short: "Stop a local Supabase instance",
	Long: `Stop a local Supabase instance and remove all containers.

This command will stop all running Supabase services for the specified project
and clean up Docker resources (containers, volumes, networks).

Example:
  supactl local stop my-project`,
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

		// Get project
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

		// Stop the instance
		fmt.Printf("Stopping Supabase instance '%s'...\n", projectID)
		fmt.Printf("Directory: %s/supabase/docker\n\n", project.Directory)

		if err := local.DockerComposeDown(projectID, project.Directory); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nSupabase instance '%s' has been stopped.\n", projectID)
	},
}

func init() {
	localCmd.AddCommand(localStopCmd)
}
