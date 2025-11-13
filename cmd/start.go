package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <instance-name>",
	Short: "Start a Supabase instance",
	Long: `Start a stopped Supabase instance.

This command works with both remote and local instances based on your current context.
Use 'supactl config use-context <name>' to switch between contexts.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := strings.TrimSpace(args[0])
		provider := getProvider()

		fmt.Printf("Starting instance '%s'...\n", instanceName)

		if err := provider.StartInstance(instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to start instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully started instance '%s'\n", instanceName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
