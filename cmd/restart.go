package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart <instance-name>",
	Short: "Restart a Supabase instance",
	Long: `Restart a Supabase instance.

This command works with both remote and local instances based on your current context.
Useful for applying configuration changes or recovering from issues.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := strings.TrimSpace(args[0])
		provider := getProvider()

		fmt.Printf("Restarting instance '%s'...\n", instanceName)

		if err := provider.RestartInstance(instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to restart instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully restarted instance '%s'\n", instanceName)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
