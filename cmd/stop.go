package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop <project-name>",
	Short: "Stop a running Supabase instance",
	Long: `Stop a running Supabase instance on your SupaControl server.

This command will stop all containers associated with the instance.
The instance data will be preserved and can be started again later.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.TrimSpace(args[0])
		client := getAPIClient()

		fmt.Printf("Stopping instance '%s'...\n", projectName)

		if err := client.StopInstance(projectName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to stop instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully stopped instance '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
