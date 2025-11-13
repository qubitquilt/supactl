package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// getCmd represents the get command (kubectl-style)
var getCmd = &cobra.Command{
	Use:   "get instances",
	Short: "Display instances (kubectl-style)",
	Long: `Display instances in the current context (kubectl-style alternative to 'list').

This command provides a kubectl-style interface for listing instances.
Works with both remote and local instances based on your current context.

Examples:
  supactl get instances`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "instances" {
			fmt.Fprintf(os.Stderr, "Error: Unknown resource type '%s'. Only 'instances' is supported.\n", args[0])
			os.Exit(1)
		}

		provider := getProvider()

		instances, err := provider.ListInstances()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to list instances: %v\n", err)
			os.Exit(1)
		}

		if len(instances) == 0 {
			fmt.Println("No instances found.")
			return
		}

		// Create a tabwriter for formatted output (kubectl-style)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATUS\tSTUDIO-URL")

		for _, instance := range instances {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				instance.Name,
				instance.Status,
				instance.StudioURL,
			)
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
