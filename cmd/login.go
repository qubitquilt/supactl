package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/qubitquilt/supactl/internal/api"
	"github.com/qubitquilt/supactl/internal/auth"
	"github.com/qubitquilt/supactl/internal/provider"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login <server_url>",
	Short: "Login to your SupaControl server",
	Long: `Login to your SupaControl server by providing your server URL and API key.

The API key can be obtained from your SupaControl dashboard.
Your credentials will be stored securely in ~/.supacontrol/config.json as the 'default' context.`,
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

		// Save the configuration using new context management functions
		config, err := auth.LoadConfig()
		if err != nil {
			// If config doesn't exist, create new one
			config = &auth.Config{
				CurrentContext: "default",
				Contexts:       make(map[string]*auth.ContextConfig),
			}
		}

		// Add/update default context
		config.AddContext("default", &auth.ContextConfig{
			Provider:  provider.ProviderTypeRemote,
			ServerURL: serverURL,
			APIKey:    apiKey,
		})

		// Ensure local context exists
		if _, exists := config.Contexts["local"]; !exists {
			config.Contexts["local"] = &auth.ContextConfig{Provider: provider.ProviderTypeLocal}
		}

		// Set current context to default
		if err := config.SetCurrentContext("default"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to set current context: %v\n", err)
			os.Exit(1)
		}

		if err := auth.SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save credentials: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully logged in to %s\n", serverURL)
		fmt.Printf("Context 'default' is now active.\n")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
