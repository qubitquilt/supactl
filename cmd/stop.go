package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop <instance-name>",
	Short: "Stop a running Supabase instance",
	Long: `Stop a running Supabase instance.

This command works with both remote and local instances based on your current context.
The instance data will be preserved and can be started again later.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := strings.TrimSpace(args[0])
		provider := getProvider()

		fmt.Printf("Stopping instance '%s'...\n", instanceName)

		if err := provider.StopInstance(instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to stop instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully stopped instance '%s'\n", instanceName)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
