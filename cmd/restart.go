package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart <project-name>",
	Short: "Restart a Supabase instance",
	Long: `Restart a Supabase instance on your SupaControl server.

This command will stop and then start all containers associated with the instance.
Useful for applying configuration changes or recovering from issues.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.TrimSpace(args[0])
		client := getAPIClient()

		fmt.Printf("Restarting instance '%s'...\n", projectName)

		if err := client.RestartInstance(projectName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to restart instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully restarted instance '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
