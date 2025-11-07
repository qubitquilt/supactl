package local

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	databaseFile = ".supascale_database.json"
)

// GetDatabasePath returns the full path to the local projects database file
func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, databaseFile), nil
}

// LoadDatabase loads the local projects database from disk
func LoadDatabase() (*Database, error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return an empty database
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return &Database{
			Projects:         make(map[string]Project),
			LastPortAssigned: BasePort,
		}, nil
	}

	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database file: %w", err)
	}

	var db Database
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("failed to parse database file: %w", err)
	}

	// Ensure projects map is initialized
	if db.Projects == nil {
		db.Projects = make(map[string]Project)
	}

	return &db, nil
}

// SaveDatabase saves the local projects database to disk
func SaveDatabase(db *Database) error {
	dbPath, err := GetDatabasePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	// Write with 0600 permissions (read/write for user only)
	if err := os.WriteFile(dbPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write database file: %w", err)
	}

	return nil
}

// ProjectExists checks if a project ID already exists in the database
func (db *Database) ProjectExists(projectID string) bool {
	_, exists := db.Projects[projectID]
	return exists
}

// GetProject retrieves a project by ID
func (db *Database) GetProject(projectID string) (*Project, error) {
	project, exists := db.Projects[projectID]
	if !exists {
		return nil, fmt.Errorf("project '%s' not found", projectID)
	}
	return &project, nil
}

// AddProject adds a new project to the database and allocates ports
func (db *Database) AddProject(projectID, directory string) (*Project, error) {
	if db.ProjectExists(projectID) {
		return nil, fmt.Errorf("project '%s' already exists", projectID)
	}

	// Allocate ports based on last assigned port
	basePort := db.LastPortAssigned
	ports := Ports{
		API:       basePort,
		DB:        basePort + 1,
		Shadow:    basePort - 1,
		Studio:    basePort + 2,
		Inbucket:  basePort + 3,
		SMTP:      basePort + 4,
		POP3:      basePort + 5,
		Pooler:    basePort + 8,
		Analytics: basePort + 6,
		KongHTTPS: basePort + 443,
	}

	project := Project{
		Directory: directory,
		Ports:     ports,
	}

	db.Projects[projectID] = project
	db.LastPortAssigned = basePort + PortIncrement

	return &project, nil
}

// RemoveProject removes a project from the database
func (db *Database) RemoveProject(projectID string) error {
	if !db.ProjectExists(projectID) {
		return fmt.Errorf("project '%s' not found", projectID)
	}
	delete(db.Projects, projectID)
	return nil
}
