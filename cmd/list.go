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
	Long: `List all Supabase instances managed by your SupaControl server.

Displays a table with project name, status, and Studio URL for each instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getAPIClient()

		instances, err := client.ListInstances()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to list instances: %v\n", err)
			os.Exit(1)
		}

		if len(instances) == 0 {
			fmt.Println("No instances found.")
			fmt.Println("Create your first instance with: supactl create <project-name>")
			return
		}

		// Create a tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "PROJECT NAME\tSTATUS\tSTUDIO URL")
		fmt.Fprintln(w, "------------\t------\t----------")

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
