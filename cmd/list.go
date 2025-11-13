package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Supabase instances",
	Long: `List all Supabase instances managed by your current context.

Displays a table with instance name, status, and Studio URL.
Works with both remote and local instances based on your current context.`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := getProvider()

		instances, err := provider.ListInstances()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to list instances: %v\n", err)
			os.Exit(1)
		}

		if len(instances) == 0 {
			fmt.Println("No instances found.")
			fmt.Printf("Create your first instance with:\n")
			fmt.Printf("  - Remote: supactl create <instance-name> (after 'supactl login')\n")
			fmt.Printf("  - Local:  supactl local add <instance-name>\n")
			return
		}

		// Create a tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "INSTANCE NAME\tSTATUS\tSTUDIO URL")
		fmt.Fprintln(w, "-------------\t------\t----------")

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
	rootCmd.AddCommand(listCmd)
}
