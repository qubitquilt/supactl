package cmd

import (
	"fmt"

	"github.com/qubitquilt/supactl/internal/local"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Manage local Supabase instances directly",
	Long: `Manage local Supabase instances directly on your machine using Docker.

The local commands allow you to create and manage multiple Supabase instances
on a single machine without requiring a SupaControl server. Each instance runs
in Docker with its own isolated environment and port configuration.

This is ideal for:
  - Local development and testing
  - Single-machine deployments
  - Learning and experimenting with Supabase

Examples:
  supactl local add my-project       # Create a new local instance
  supactl local list                 # List all local instances
  supactl local start my-project     # Start an instance
  supactl local stop my-project      # Stop an instance
  supactl local remove my-project    # Remove an instance`,
}

func init() {
	rootCmd.AddCommand(localCmd)
}

// getLocalDatabase loads the local projects database
func getLocalDatabase() (*local.Database, error) {
	db, err := local.LoadDatabase()
	if err != nil {
		return nil, fmt.Errorf("loading local database: %w", err)
	}
	return db, nil
}

// checkDockerRequirements ensures Docker and Docker Compose are available
func checkDockerRequirements() error {
	if err := local.CheckDockerAvailable(); err != nil {
		return fmt.Errorf("Docker is required but not available.\nPlease install Docker and ensure it's running.\nVisit https://docs.docker.com/get-docker/ for installation instructions")
	}

	if err := local.CheckDockerComposeAvailable(); err != nil {
		return fmt.Errorf("Docker Compose is required but not available.\nPlease install Docker Compose.\nVisit https://docs.docker.com/compose/install/ for installation instructions")
	}

	return nil
}
