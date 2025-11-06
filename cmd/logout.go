package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/supactl/internal/auth"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your SupaControl server",
	Long: `Logout from your SupaControl server by removing your stored credentials.

This will delete the configuration file at ~/.supacontrol/config.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !auth.IsLoggedIn() {
			fmt.Println("You are not currently logged in.")
			return
		}

		if err := auth.ClearConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to logout: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Successfully logged out.")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
