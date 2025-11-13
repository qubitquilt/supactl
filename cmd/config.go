package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/qubitquilt/supactl/internal/auth"
	"github.com/qubitquilt/supactl/internal/provider"
	"github.com/spf13/cobra"
)

// configCmd represents the config command (kubectl-style context management)
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify supactl configuration (kubectl-style)",
	Long: `Modify supactl configuration files using subcommands.

This command provides kubectl-style context management for switching between
different providers (local or remote servers).

Available Commands:
  get-contexts      List all contexts
  use-context       Set the current context
  current-context   Display the current context
  set-context       Create or update a context
  delete-context    Delete a context`,
}

// configGetContextsCmd lists all contexts
var configGetContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "List all contexts",
	Long:  `List all available contexts in the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := auth.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load config: %v\n", err)
			os.Exit(1)
		}

		if len(config.Contexts) == 0 {
			fmt.Println("No contexts found.")
			return
		}

		fmt.Printf("CURRENT   NAME       PROVIDER    SERVER\n")

		// Sort context names for consistent output
		names := make([]string, 0, len(config.Contexts))
		for name := range config.Contexts {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			ctx := config.Contexts[name]
			current := " "
			if name == config.CurrentContext {
				current = "*"
			}

			serverURL := "-"
			if ctx.ServerURL != "" {
				serverURL = ctx.ServerURL
			}

			fmt.Printf("%-9s %-10s %-11s %s\n", current, name, ctx.Provider, serverURL)
		}
	},
}

// configUseContextCmd switches the current context
var configUseContextCmd = &cobra.Command{
	Use:   "use-context <context-name>",
	Short: "Set the current context",
	Long: `Set the current context to use for all commands.

Example:
  supactl config use-context local
  supactl config use-context production`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]

		config, err := auth.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load config: %v\n", err)
			os.Exit(1)
		}

		if err := config.SetCurrentContext(contextName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := auth.SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Switched to context '%s'\n", contextName)
	},
}

// configCurrentContextCmd shows the current context
var configCurrentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Display the current context",
	Long:  `Display the name of the current context.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := auth.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(config.CurrentContext)
	},
}

// configSetContextCmd creates or updates a context
var (
	setContextProvider string
	setContextServer   string
	setContextAPIKey   string
)

var configSetContextCmd = &cobra.Command{
	Use:   "set-context <context-name>",
	Short: "Create or update a context",
	Long: `Create or update a context with specified provider and credentials.

Examples:
  # Set local context
  supactl config set-context local --provider=local

  # Set remote context
  supactl config set-context prod --provider=remote --server=https://api.example.com --api-key=sk_...`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]

		// Validate provider
		if setContextProvider != provider.ProviderTypeLocal && setContextProvider != provider.ProviderTypeRemote {
			fmt.Fprintf(os.Stderr, "Error: Provider must be '%s' or '%s'\n", provider.ProviderTypeLocal, provider.ProviderTypeRemote)
			os.Exit(1)
		}

		// Validate remote context has required fields
		if setContextProvider == provider.ProviderTypeRemote && (setContextServer == "" || setContextAPIKey == "") {
			fmt.Fprintf(os.Stderr, "Error: Remote contexts require --server and --api-key flags\n")
			os.Exit(1)
		}

		config, err := auth.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Create or update context
		ctx := &auth.ContextConfig{
			Provider: setContextProvider,
		}

		if setContextProvider == provider.ProviderTypeRemote {
			ctx.ServerURL = setContextServer
			ctx.APIKey = setContextAPIKey
		}

		config.AddContext(contextName, ctx)

		if err := auth.SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Context '%s' created/updated\n", contextName)
	},
}

// configDeleteContextCmd deletes a context
var configDeleteContextCmd = &cobra.Command{
	Use:   "delete-context <context-name>",
	Short: "Delete a context",
	Long: `Delete a context from the configuration.

Cannot delete the 'local' context or the current context.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]

		config, err := auth.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load config: %v\n", err)
			os.Exit(1)
		}

		if err := config.RemoveContext(contextName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := auth.SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Context '%s' deleted\n", contextName)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configGetContextsCmd)
	configCmd.AddCommand(configUseContextCmd)
	configCmd.AddCommand(configCurrentContextCmd)
	configCmd.AddCommand(configSetContextCmd)
	configCmd.AddCommand(configDeleteContextCmd)

	// Flags for set-context command
	configSetContextCmd.Flags().StringVar(&setContextProvider, "provider", "", "Provider type (local or remote)")
	configSetContextCmd.Flags().StringVar(&setContextServer, "server", "", "Server URL (for remote provider)")
	configSetContextCmd.Flags().StringVar(&setContextAPIKey, "api-key", "", "API key (for remote provider)")
	configSetContextCmd.MarkFlagRequired("provider")
}
