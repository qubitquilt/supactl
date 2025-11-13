package cmd

import (
	"fmt"
	"os"

	"github.com/qubitquilt/supactl/internal/api"
	"github.com/qubitquilt/supactl/internal/auth"
	"github.com/qubitquilt/supactl/internal/provider"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "supactl",
	Short: "SupaControl CLI - Manage your Supabase instances",
	Long: `supactl is a command-line interface for managing self-hosted Supabase
instances in two modes:

1. Remote Mode: Manage instances via a SupaControl server
2. Local Mode: Manage local Docker-based instances

Use contexts to switch between different providers (local or remote servers).`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("supactl version {{.Version}}\n")
}

// getProvider creates and returns the appropriate provider based on the current context
func getProvider() provider.InstanceProvider {
	config, err := auth.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	ctx, err := config.GetCurrentContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'supactl config use-context <name>' to set a context.\n")
		os.Exit(1)
	}

	switch ctx.Provider {
	case provider.ProviderTypeRemote:
		if ctx.ServerURL == "" || ctx.APIKey == "" {
			fmt.Fprintf(os.Stderr, "Error: Current context '%s' is a remote context but is missing credentials.\n", config.CurrentContext)
			fmt.Fprintf(os.Stderr, "Run 'supactl login <server_url>' or 'supactl config set-context %s --server=<url> --api-key=<key>'\n", config.CurrentContext)
			os.Exit(1)
		}
		return provider.NewRemoteProvider(ctx.ServerURL, ctx.APIKey)

	case provider.ProviderTypeLocal:
		localProvider, err := provider.NewLocalProvider()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to initialize local provider: %v\n", err)
			os.Exit(1)
		}
		return localProvider

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown provider type '%s' in context '%s'\n", ctx.Provider, config.CurrentContext)
		os.Exit(1)
		return nil
	}
}

// getAPIClient creates and returns an API client, or exits if not logged in
// Deprecated: Use getProvider() instead for context-aware provider selection
func getAPIClient() *api.Client {
	config, err := auth.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: You are not logged in. Please run 'supactl login <server_url>' first.\n")
		os.Exit(1)
	}

	ctx, err := config.GetCurrentContext()
	if err != nil || ctx.Provider != provider.ProviderTypeRemote {
		fmt.Fprintf(os.Stderr, "Error: Current context is not a remote context. Please run 'supactl login <server_url>' first.\n")
		os.Exit(1)
	}

	return api.NewClient(ctx.ServerURL, ctx.APIKey)
}
