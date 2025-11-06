package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var projectNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <project-name>",
	Short: "Create a new Supabase instance",
	Long: `Create a new Supabase instance on your SupaControl server.

The project name must be lowercase, alphanumeric, and may contain hyphens.
It must start and end with an alphanumeric character.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.TrimSpace(args[0])

		// Validate project name
		if !projectNameRegex.MatchString(projectName) {
			fmt.Fprintf(os.Stderr, "Error: Project name '%s' is invalid.\n", projectName)
			fmt.Fprintf(os.Stderr, "Name must be lowercase, alphanumeric, and may contain hyphens.\n")
			fmt.Fprintf(os.Stderr, "It must start and end with an alphanumeric character.\n")
			os.Exit(1)
		}

		client := getAPIClient()

		fmt.Printf("Creating instance '%s'...\n", projectName)

		instance, err := client.CreateInstance(projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nSuccessfully created instance '%s'\n\n", instance.Name)
		fmt.Printf("  Status:     %s\n", instance.Status)
		fmt.Printf("  Studio URL: %s\n", instance.StudioURL)
		if instance.APIURL != "" {
			fmt.Printf("  API URL:    %s\n", instance.APIURL)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
