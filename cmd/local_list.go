package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var localListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all local Supabase instances",
	Long: `List all configured local Supabase instances with their port assignments
and directory locations.

Example:
  supactl local list`,
	Run: func(cmd *cobra.Command, args []string) {
		db := getLocalDatabase()

		if len(db.Projects) == 0 {
			fmt.Println("No local projects configured yet.")
			fmt.Println()
			fmt.Println("Create a new project with:")
			fmt.Println("  supactl local add <project-id>")
			return
		}

		fmt.Println("Configured Local Supabase Projects:")
		fmt.Println("====================================")
		fmt.Println()

		for projectID, project := range db.Projects {
			fmt.Printf("Project ID: %s\n", projectID)
			fmt.Printf("  Directory:     %s\n", project.Directory)
			fmt.Printf("  API Port:      %d\n", project.Ports.API)
			fmt.Printf("  DB Port:       %d\n", project.Ports.DB)
			fmt.Printf("  Studio Port:   %d\n", project.Ports.Studio)
			fmt.Printf("  Inbucket Port: %d\n", project.Ports.Inbucket)
			fmt.Printf("  Analytics Port: %d\n", project.Ports.Analytics)
			fmt.Println()
		}
	},
}

func init() {
	localCmd.AddCommand(localListCmd)
}
