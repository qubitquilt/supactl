package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/yourusername/supactl/internal/link"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link current directory to a remote instance",
	Long: `Link the current directory to a remote Supabase instance.

This command will present a list of your available instances and allow you
to select one to link to the current directory. A .supacontrol/project file
will be created to store the link.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getAPIClient()

		// Check if already linked
		if link.IsLinked() {
			existingProject, _ := link.GetLink()
			fmt.Printf("This directory is already linked to '%s'\n", existingProject)
			fmt.Println("Run 'supactl unlink' first to unlink it.")
			return
		}

		// Get list of instances
		fmt.Println("Fetching your instances...")
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

		// Build options for the selector
		options := make([]string, len(instances))
		for i, instance := range instances {
			options[i] = fmt.Sprintf("%s (%s)", instance.Name, instance.Status)
		}

		// Present interactive selector
		var selectedIndex int
		prompt := &survey.Select{
			Message: "Select a project to link:",
			Options: options,
		}

		if err := survey.AskOne(prompt, &selectedIndex); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		selectedProject := instances[selectedIndex].Name

		// Save the link
		if err := link.SaveLink(selectedProject); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save link: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nSuccessfully linked to '%s'\n", selectedProject)
		fmt.Println("Run 'supactl status' to see project details.")
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
