package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <project-name>",
	Short: "Start a Supabase instance",
	Long: `Start a stopped Supabase instance on your SupaControl server.

This command will start all containers associated with the instance.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.TrimSpace(args[0])
		client := getAPIClient()

		fmt.Printf("Starting instance '%s'...\n", projectName)

		if err := client.StartInstance(projectName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to start instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully started instance '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
