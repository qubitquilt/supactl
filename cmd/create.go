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
	Use:   "create <instance-name>",
	Short: "Create a new Supabase instance",
	Long: `Create a new Supabase instance.

This command works with remote contexts only. For local instances, use 'supactl local add'.
The instance name must be lowercase, alphanumeric, and may contain hyphens.
It must start and end with an alphanumeric character.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := strings.TrimSpace(args[0])

		// Validate project name
		if !projectNameRegex.MatchString(instanceName) {
			fmt.Fprintf(os.Stderr, "Error: Instance name '%s' is invalid.\n", instanceName)
			fmt.Fprintf(os.Stderr, "Name must be lowercase, alphanumeric, and may contain hyphens.\n")
			fmt.Fprintf(os.Stderr, "It must start and end with an alphanumeric character.\n")
			os.Exit(1)
		}

		provider := getProvider()

		fmt.Printf("Creating instance '%s'...\n", instanceName)

		instance, err := provider.CreateInstance(instanceName)
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
