package link

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	linkDir  = ".supacontrol"
	linkFile = "project"
)

// GetLinkPath returns the path to the local project link file
func GetLinkPath() string {
	return filepath.Join(".", linkDir, linkFile)
}

// SaveLink saves the project name to the local link file
func SaveLink(projectName string) error {
	linkDirPath := filepath.Join(".", linkDir)

	// Create link directory if it doesn't exist
	if err := os.MkdirAll(linkDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create link directory: %w", err)
	}

	linkPath := GetLinkPath()
	if err := os.WriteFile(linkPath, []byte(projectName), 0644); err != nil {
		return fmt.Errorf("failed to write link file: %w", err)
	}

	// Try to add .supacontrol/ to .gitignore if it exists
	if err := addToGitignore(); err != nil {
		// Non-fatal error, just continue
		fmt.Printf("Warning: could not update .gitignore: %v\n", err)
	}

	return nil
}

// GetLink reads the project name from the local link file
func GetLink() (string, error) {
	linkPath := GetLinkPath()

	data, err := os.ReadFile(linkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no project linked. Run 'supactl link' to get started")
		}
		return "", fmt.Errorf("failed to read link file: %w", err)
	}

	projectName := strings.TrimSpace(string(data))
	if projectName == "" {
		return "", fmt.Errorf("link file is empty. Run 'supactl link' to link a project")
	}

	return projectName, nil
}

// ClearLink removes the local link file
func ClearLink() error {
	linkPath := GetLinkPath()

	if err := os.Remove(linkPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already cleared
		}
		return fmt.Errorf("failed to remove link file: %w", err)
	}

	return nil
}

// IsLinked checks if a project is currently linked
func IsLinked() bool {
	_, err := GetLink()
	return err == nil
}

// addToGitignore adds .supacontrol/ to .gitignore if not already present
func addToGitignore() error {
	gitignorePath := ".gitignore"
	pattern := ".supacontrol/"

	// Check if .gitignore exists
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// .gitignore doesn't exist, skip
		return nil
	}

	// Read existing .gitignore
	data, err := os.ReadFile(gitignorePath)
	if err != nil {
		return err
	}

	content := string(data)

	// Check if pattern already exists
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == pattern {
			// Pattern already exists
			return nil
		}
	}

	// Add pattern to .gitignore
	if !strings.HasSuffix(content, "\n") && content != "" {
		content += "\n"
	}
	content += pattern + "\n"

	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}
