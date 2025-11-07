package cmd

import (
	"fmt"
	"os"

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

// getLocalDatabase loads the local projects database and exits on error
func getLocalDatabase() *local.Database {
	db, err := local.LoadDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading local database: %v\n", err)
		os.Exit(1)
	}
	return db
}

// checkDockerRequirements ensures Docker and Docker Compose are available
func checkDockerRequirements() {
	if err := local.CheckDockerAvailable(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Docker is required but not available.\n")
		fmt.Fprintf(os.Stderr, "Please install Docker and ensure it's running.\n")
		fmt.Fprintf(os.Stderr, "Visit https://docs.docker.com/get-docker/ for installation instructions.\n")
		os.Exit(1)
	}

	if err := local.CheckDockerComposeAvailable(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Docker Compose is required but not available.\n")
		fmt.Fprintf(os.Stderr, "Please install Docker Compose.\n")
		fmt.Fprintf(os.Stderr, "Visit https://docs.docker.com/compose/install/ for installation instructions.\n")
		os.Exit(1)
	}
}
