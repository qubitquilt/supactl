package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/qubitquilt/supactl/internal/api"
	"github.com/qubitquilt/supactl/internal/auth"
)

var (
	version = "1.0.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "supactl",
	Short: "SupaControl CLI - Manage your Supabase instances",
	Long: `supactl is a command-line interface for managing self-hosted Supabase
instances via a central SupaControl server.

Use supactl to create, list, delete, and manage Supabase instances from your terminal.`,
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

// getAPIClient creates and returns an API client, or exits if not logged in
func getAPIClient() *api.Client {
	config, err := auth.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: You are not logged in. Please run 'supactl login <server_url>' first.\n")
		os.Exit(1)
	}

	return api.NewClient(config.ServerURL, config.APIKey)
}
