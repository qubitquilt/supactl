package cmd

import (
	"fmt"
	"os"

	"github.com/qubitquilt/supactl/internal/link"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of linked project",
	Long: `Show detailed status and information about the linked project.

This command requires the current directory to be linked to a project.
Run 'supactl link' first if you haven't already.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the linked project name
		projectName, err := link.GetLink()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := getAPIClient()

		// Fetch instance details
		instance, err := client.GetInstance(projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get instance details: %v\n", err)
			os.Exit(1)
		}

		// Display instance information
		fmt.Printf("Project: %s\n\n", instance.Name)
		fmt.Printf("  Status:     %s\n", instance.Status)
		fmt.Printf("  Studio URL: %s\n", instance.StudioURL)

		if instance.APIURL != "" {
			fmt.Printf("  API URL:    %s\n", instance.APIURL)
		}
		if instance.KongURL != "" {
			fmt.Printf("  Kong URL:   %s\n", instance.KongURL)
		}
		if instance.AnonKey != "" {
			fmt.Printf("  Anon Key:   %s\n", instance.AnonKey)
		}
		if instance.CreatedAt != "" {
			fmt.Printf("  Created:    %s\n", instance.CreatedAt)
		}

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
