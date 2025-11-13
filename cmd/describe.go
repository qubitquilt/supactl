package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command (kubectl-style)
var describeCmd = &cobra.Command{
	Use:   "describe instance <instance-name>",
	Short: "Show detailed instance information (kubectl-style)",
	Long: `Show detailed information about a specific instance (kubectl-style alternative to 'status').

This command provides a kubectl-style interface for viewing instance details.
Works with both remote and local instances based on your current context.

Examples:
  supactl describe instance my-project`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		instanceName := strings.TrimSpace(args[1])

		if resourceType != "instance" {
			fmt.Fprintf(os.Stderr, "Error: Unknown resource type '%s'. Only 'instance' is supported.\n", resourceType)
			os.Exit(1)
		}

		provider := getProvider()

		// Fetch instance details
		instance, err := provider.GetInstance(instanceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get instance: %v\n", err)
			os.Exit(1)
		}

		// Display instance information in kubectl describe style
		fmt.Printf("Name:\t\t%s\n", instance.Name)
		fmt.Printf("Status:\t\t%s\n", instance.Status)
		fmt.Printf("Studio URL:\t%s\n", instance.StudioURL)

		if instance.APIURL != "" {
			fmt.Printf("API URL:\t%s\n", instance.APIURL)
		}
		if instance.KongURL != "" {
			fmt.Printf("Kong URL:\t%s\n", instance.KongURL)
		}
		if instance.DatabaseURL != "" {
			fmt.Printf("Database URL:\t%s\n", instance.DatabaseURL)
		}
		if instance.Directory != "" {
			fmt.Printf("Directory:\t%s\n", instance.Directory)
		}
		if instance.DBPort != 0 {
			fmt.Printf("DB Port:\t%d\n", instance.DBPort)
		}
		if instance.AnonKey != "" {
			fmt.Printf("Anon Key:\t%s\n", instance.AnonKey)
		}
		if instance.ServiceKey != "" {
			fmt.Printf("Service Key:\t%s\n", instance.ServiceKey)
		}
		if !instance.CreatedAt.IsZero() {
			fmt.Printf("Created:\t%s\n", instance.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
