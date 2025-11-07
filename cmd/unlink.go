package cmd

import (
	"fmt"
	"os"

	"github.com/qubitquilt/supactl/internal/link"
	"github.com/spf13/cobra"
)

// unlinkCmd represents the unlink command
var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink current directory from remote instance",
	Long: `Unlink the current directory from its linked remote Supabase instance.

This will remove the .supacontrol/project file from the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !link.IsLinked() {
			fmt.Println("This directory is not linked to any project.")
			return
		}

		projectName, _ := link.GetLink()

		if err := link.ClearLink(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to unlink: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully unlinked from '%s'\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}
