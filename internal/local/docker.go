package local

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DockerComposeUp starts the Docker Compose services for a project
func DockerComposeUp(projectID, directory string) error {
	dockerDir := filepath.Join(directory, "supabase", "docker")

	// Check if directory exists
	if _, err := os.Stat(dockerDir); os.IsNotExist(err) {
		return fmt.Errorf("docker directory not found: %s", dockerDir)
	}

	// Check if .env file exists
	envFile := filepath.Join(dockerDir, ".env")
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return fmt.Errorf(".env file not found: %s", envFile)
	}

	// Run docker compose up -d
	cmd := exec.Command("docker", "compose", "-p", projectID, "up", "-d")
	cmd.Dir = dockerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose up failed: %w", err)
	}

	return nil
}

// DockerComposeDown stops and removes the Docker Compose services for a project
func DockerComposeDown(projectID, directory string) error {
	dockerDir := filepath.Join(directory, "supabase", "docker")

	// Check if directory exists
	if _, err := os.Stat(dockerDir); os.IsNotExist(err) {
		return fmt.Errorf("docker directory not found: %s", dockerDir)
	}

	// Run docker compose down -v --remove-orphans
	cmd := exec.Command("docker", "compose", "-p", projectID, "down", "-v", "--remove-orphans")
	cmd.Dir = dockerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose down failed: %w", err)
	}

	return nil
}

// CheckDockerAvailable checks if Docker is available on the system
func CheckDockerAvailable() error {
	cmd := exec.Command("docker", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker is not available or not running: %w", err)
	}
	return nil
}

// CheckDockerComposeAvailable checks if Docker Compose is available
func CheckDockerComposeAvailable() error {
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose is not available: %w", err)
	}
	return nil
}
