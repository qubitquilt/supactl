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
	Use:   "delete <project-name>",
	Short: "Delete a Supabase instance",
	Long: `Delete a Supabase instance from your SupaControl server.

WARNING: This action is irreversible. All data associated with the instance
will be permanently deleted.

You will be asked to confirm before the deletion proceeds.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.TrimSpace(args[0])
		client := getAPIClient()

		// Ask for confirmation
		var confirmed bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete '%s'? This action is irreversible.", projectName),
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

		fmt.Printf("Deleting instance '%s'...\n", projectName)

		if err := client.DeleteInstance(projectName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to delete instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted instance '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
