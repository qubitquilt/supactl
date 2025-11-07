package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/qubitquilt/supactl/internal/local"
	"github.com/spf13/cobra"
)

var localStartCmd = &cobra.Command{
	Use:   "start <project-id>",
	Short: "Start a local Supabase instance",
	Long: `Start a local Supabase instance using Docker Compose.

This command will start all Supabase services (PostgreSQL, Kong, GoTrue, etc.)
in Docker containers with the configured ports.

Example:
  supactl local start my-project`,
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

		// Start the instance
		fmt.Printf("Starting Supabase instance '%s'...\n", projectID)
		fmt.Printf("Directory: %s/supabase/docker\n\n", project.Directory)

		if err := local.DockerComposeUp(projectID, project.Directory); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Get host IP for display
		hostIP := "localhost"
		if output, err := exec.Command("hostname", "-I").Output(); err == nil {
			if len(output) > 0 {
				// Take first IP
				for i, b := range output {
					if b == ' ' || b == '\n' {
						hostIP = string(output[:i])
						break
					}
				}
			}
		}

		fmt.Println()
		fmt.Printf("Supabase is now running for project '%s':\n", projectID)
		fmt.Printf("  Studio URL: http://%s:%d\n", hostIP, project.Ports.Studio)
		fmt.Printf("  API URL:    http://%s:%d/rest/v1/\n", hostIP, project.Ports.API)
		fmt.Printf("  DB URL:     postgresql://postgres:postgres@%s:%d/postgres\n", hostIP, project.Ports.DB)
	},
}

func init() {
	localCmd.AddCommand(localStartCmd)
}
