package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qubitquilt/supactl/internal/local"
	"github.com/spf13/cobra"
)

var localAddCmd = &cobra.Command{
	Use:   "add <project-id>",
	Short: "Add a new local Supabase instance",
	Long: `Add a new local Supabase instance by cloning the Supabase repository,
generating secure credentials, and configuring Docker Compose.

The project ID must:
  - Start with a letter or number
  - Contain only lowercase letters, numbers, hyphens, and underscores
  - No dots, spaces, or special characters allowed

This command will:
  1. Create a new directory for the project
  2. Clone the Supabase repository
  3. Generate secure passwords and JWT tokens
  4. Configure .env file with generated secrets
  5. Update docker-compose.yml with unique ports
  6. Save project configuration to the local database

Example:
  supactl local add my-project`,
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

		// Determine project directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
			os.Exit(1)
		}
		directory := filepath.Join(homeDir, projectID)

		// Setup the project
		fmt.Printf("Creating local Supabase instance '%s'...\n", projectID)
		fmt.Printf("Directory: %s\n\n", directory)

		secrets, err := local.SetupProject(projectID, directory, db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Get the project details for port information
		project, err := db.GetProject(projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Print success message
		fmt.Println()
		fmt.Println("----------------------------------------------------------------------")
		fmt.Println("SUCCESS: PROJECT CREATED AND CONFIGURED")
		fmt.Println("----------------------------------------------------------------------")
		fmt.Printf("Project '%s' has been successfully created and configured.\n", projectID)
		fmt.Printf("Generated secrets have been saved to:\n")
		fmt.Printf("  %s/supabase/docker/.env\n", directory)
		fmt.Println()
		fmt.Println("Generated credentials:")
		fmt.Printf("  DASHBOARD_USERNAME: supabase\n")
		fmt.Printf("  DASHBOARD_PASSWORD: %s\n", secrets.DashboardPassword)
		fmt.Printf("  POSTGRES_PASSWORD:  %s\n", secrets.PostgresPassword)
		fmt.Printf("  VAULT_ENC_KEY:      %s\n", secrets.VaultEncKey)
		fmt.Printf("  JWT_SECRET:         %s\n", secrets.JWTSecret)
		fmt.Println()
		fmt.Println("Generated JWT keys:")
		fmt.Printf("  ANON_KEY:           %s\n", secrets.AnonKey)
		fmt.Printf("  SERVICE_ROLE_KEY:   %s\n", secrets.ServiceRoleKey)
		fmt.Println()
		fmt.Println("Assigned ports:")
		fmt.Printf("  API Port:      %d\n", project.Ports.API)
		fmt.Printf("  DB Port:       %d\n", project.Ports.DB)
		fmt.Printf("  Studio Port:   %d\n", project.Ports.Studio)
		fmt.Printf("  Inbucket Port: %d\n", project.Ports.Inbucket)
		fmt.Println("----------------------------------------------------------------------")
		fmt.Println()
		fmt.Println("Configuration complete! Start your instance with:")
		fmt.Printf("  supactl local start %s\n", projectID)
	},
}

func init() {
	localCmd.AddCommand(localAddCmd)
}
