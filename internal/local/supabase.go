package local

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

const (
	supabaseRepoURL = "https://github.com/supabase/supabase"
)

// ValidateProjectID validates the project ID format
func ValidateProjectID(projectID string) error {
	// Project ID must start with a letter or number, and contain only lowercase letters, numbers, hyphens, and underscores
	pattern := `^[a-z0-9][a-z0-9_-]*$`
	matched, err := regexp.MatchString(pattern, projectID)
	if err != nil {
		return fmt.Errorf("failed to validate project ID: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid project ID '%s':\n"+
			"  - Must start with a letter or number\n"+
			"  - Can contain only lowercase letters, numbers, hyphens, and underscores\n"+
			"  - No dots, spaces, or special characters allowed", projectID)
	}

	return nil
}

// CloneSupabaseRepo clones the Supabase repository into the specified directory
func CloneSupabaseRepo(directory string) error {
	// Check if directory already exists
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		return fmt.Errorf("directory already exists: %s", directory)
	}

	// Create the directory
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Clone the repository with depth 1 (shallow clone)
	fmt.Printf("Cloning Supabase repository into %s...\n", directory)
	cmd := exec.Command("git", "clone", "--depth", "1", supabaseRepoURL, filepath.Join(directory, "supabase"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Clean up on failure
		os.RemoveAll(directory)
		return fmt.Errorf("failed to clone Supabase repository: %w", err)
	}

	// Verify the docker directory exists
	dockerDir := filepath.Join(directory, "supabase", "docker")
	if _, err := os.Stat(dockerDir); os.IsNotExist(err) {
		os.RemoveAll(directory)
		return fmt.Errorf("docker directory not found in cloned repository")
	}

	return nil
}

// SetupEnvFile copies .env.example to .env and updates it with secrets
func SetupEnvFile(directory string, secrets *Secrets, ports *Ports) error {
	dockerDir := filepath.Join(directory, "supabase", "docker")
	envExamplePath := filepath.Join(dockerDir, ".env.example")
	envPath := filepath.Join(dockerDir, ".env")

	// Check if .env.example exists
	if _, err := os.Stat(envExamplePath); os.IsNotExist(err) {
		return fmt.Errorf(".env.example not found: %s", envExamplePath)
	}

	// Copy .env.example to .env
	fmt.Println("Creating .env file from .env.example...")
	content, err := os.ReadFile(envExamplePath)
	if err != nil {
		return fmt.Errorf("failed to read .env.example: %w", err)
	}

	if err := os.WriteFile(envPath, content, 0600); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	// Update .env file with secrets and ports
	fmt.Println("Updating .env file with generated secrets...")
	if err := UpdateEnvFile(envPath, secrets, ports); err != nil {
		return err
	}

	return nil
}

// SetupConfigurationFiles updates docker-compose.yml and config.toml
func SetupConfigurationFiles(directory, projectID string, ports *Ports) error {
	// Update docker-compose.yml
	composePath := filepath.Join(directory, "supabase", "docker", "docker-compose.yml")
	if _, err := os.Stat(composePath); err == nil {
		fmt.Println("Updating docker-compose.yml...")
		if err := UpdateDockerComposeFile(composePath, projectID, ports); err != nil {
			return err
		}
	} else {
		fmt.Printf("Warning: docker-compose.yml not found at %s\n", composePath)
	}

	// Update config.toml (if it exists)
	configPath := filepath.Join(directory, "supabase", "supabase", "config.toml")
	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("Updating config.toml...")
		if err := UpdateConfigToml(configPath, projectID, ports); err != nil {
			return err
		}
	} else {
		fmt.Println("Note: config.toml not found (this is normal for newer Supabase versions)")
	}

	return nil
}

// SetupProject orchestrates the full project setup process
func SetupProject(projectID, directory string, db *Database) (*Secrets, error) {
	// Validate project ID
	if err := ValidateProjectID(projectID); err != nil {
		return nil, err
	}

	// Check if project already exists in database
	if db.ProjectExists(projectID) {
		return nil, fmt.Errorf("project '%s' already exists", projectID)
	}

	// Clone Supabase repository
	if err := CloneSupabaseRepo(directory); err != nil {
		return nil, err
	}

	// Generate secrets
	fmt.Println("Generating secrets...")
	secrets, err := GenerateSecrets()
	if err != nil {
		os.RemoveAll(directory)
		return nil, err
	}

	// Add project to database (this allocates ports)
	project, err := db.AddProject(projectID, directory)
	if err != nil {
		os.RemoveAll(directory)
		return nil, err
	}

	// Setup .env file
	if err := SetupEnvFile(directory, secrets, &project.Ports); err != nil {
		os.RemoveAll(directory)
		db.RemoveProject(projectID)
		return nil, err
	}

	// Setup configuration files
	if err := SetupConfigurationFiles(directory, projectID, &project.Ports); err != nil {
		os.RemoveAll(directory)
		db.RemoveProject(projectID)
		return nil, err
	}

	// Save database
	if err := SaveDatabase(db); err != nil {
		os.RemoveAll(directory)
		db.RemoveProject(projectID)
		return nil, fmt.Errorf("failed to save database: %w", err)
	}

	return secrets, nil
}
