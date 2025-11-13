package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "Name:\t%s\n", instance.Name)
		fmt.Fprintf(w, "Status:\t%s\n", instance.Status)
		fmt.Fprintf(w, "Studio URL:\t%s\n", instance.StudioURL)

		if instance.APIURL != "" {
			fmt.Fprintf(w, "API URL:\t%s\n", instance.APIURL)
		}
		if instance.KongURL != "" {
			fmt.Fprintf(w, "Kong URL:\t%s\n", instance.KongURL)
		}
		if instance.DatabaseURL != "" {
			fmt.Fprintf(w, "Database URL:\t%s\n", instance.DatabaseURL)
		}
		if instance.Directory != "" {
			fmt.Fprintf(w, "Directory:\t%s\n", instance.Directory)
		}
		if instance.DBPort != 0 {
			fmt.Fprintf(w, "DB Port:\t%d\n", instance.DBPort)
		}
		if instance.AnonKey != "" {
			fmt.Fprintf(w, "Anon Key:\t%s\n", instance.AnonKey)
		}
		if instance.ServiceKey != "" {
			fmt.Fprintf(w, "Service Key:\t%s\n", instance.ServiceKey)
		}
		if !instance.CreatedAt.IsZero() {
			fmt.Fprintf(w, "Created:\t%s\n", instance.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
