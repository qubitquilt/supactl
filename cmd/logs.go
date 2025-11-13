package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	logLines int
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs <instance-name>",
	Short: "View logs for a Supabase instance",
	Long: `View logs for a Supabase instance.

This command retrieves and displays the recent logs from the instance containers.
Works with both remote and local instances based on your current context.
Use the --lines flag to control how many lines to display.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := strings.TrimSpace(args[0])
		provider := getProvider()

		fmt.Printf("Fetching logs for instance '%s'...\n\n", instanceName)

		logs, err := provider.GetLogs(instanceName, logLines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to fetch logs: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(logs)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().IntVarP(&logLines, "lines", "n", 100, "Number of lines to show")
}
