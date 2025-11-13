package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <instance-name>",
	Short: "Delete a Supabase instance",
	Long: `Delete a Supabase instance.

WARNING: This action may be irreversible depending on your provider.
For remote instances, all data will be permanently deleted.
For local instances, only the database entry is removed (files remain).

You will be asked to confirm before the deletion proceeds.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := strings.TrimSpace(args[0])
		provider := getProvider()

		// Ask for confirmation
		var confirmed bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete '%s'?", instanceName),
			Default: false,
		}

		if err := survey.AskOne(prompt, &confirmed); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !confirmed {
			fmt.Println("Deletion cancelled.")
			return
		}

		fmt.Printf("Deleting instance '%s'...\n", instanceName)

		if err := provider.DeleteInstance(instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to delete instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted instance '%s'\n", instanceName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
