package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/qubitquilt/supactl/internal/api"
	"github.com/qubitquilt/supactl/internal/auth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login <server_url>",
	Short: "Login to your SupaControl server",
	Long: `Login to your SupaControl server by providing your server URL and API key.

The API key can be obtained from your SupaControl dashboard.
Your credentials will be stored securely in ~/.supacontrol/config.json.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverURL := strings.TrimRight(args[0], "/")

		// Validate URL format
		if !strings.HasPrefix(serverURL, "http://") && !strings.HasPrefix(serverURL, "https://") {
			fmt.Fprintf(os.Stderr, "Error: Server URL must start with http:// or https://\n")
			os.Exit(1)
		}

		// Prompt for API key (no echo)
		var apiKey string
		prompt := &survey.Password{
			Message: "Enter your API key:",
		}
		if err := survey.AskOne(prompt, &apiKey, survey.WithValidator(survey.Required)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Test the credentials
		fmt.Println("Validating credentials...")
		client := api.NewClient(serverURL, apiKey)
		if err := client.LoginTest(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Authentication failed: %v\n", err)
			os.Exit(1)
		}

		// Save the configuration
		if err := auth.SaveConfig(serverURL, apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save credentials: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully logged in to %s\n", serverURL)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
