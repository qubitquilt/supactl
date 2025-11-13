package provider

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/qubitquilt/supactl/internal/local"
)

// LocalProvider implements InstanceProvider for local Docker-based instances
type LocalProvider struct {
	db *local.Database
}

// NewLocalProvider creates a new local provider
func NewLocalProvider() (*LocalProvider, error) {
	db, err := local.LoadDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to load local database: %w", err)
	}

	return &LocalProvider{db: db}, nil
}

// reloadDatabase reloads the database from disk (for operations that might have changed it)
func (p *LocalProvider) reloadDatabase() error {
	db, err := local.LoadDatabase()
	if err != nil {
		return fmt.Errorf("failed to reload local database: %w", err)
	}
	p.db = db
	return nil
}

// mapProjectToInstance converts a local project to a unified instance
func mapProjectToInstance(name string, project *local.Project) *Instance {
	// Determine status by checking if containers are running
	status := "stopped"
	if isProjectRunning(name, project.Directory) {
		status = "running"
	}

	hostIP := getHostIP()

	return &Instance{
		Name:      name,
		Status:    status,
		StudioURL: fmt.Sprintf("http://%s:%d", hostIP, project.Ports.Studio),
		APIURL:    fmt.Sprintf("http://%s:%d/rest/v1/", hostIP, project.Ports.API),
		Directory: project.Directory,
		DBPort:    project.Ports.DB,
		CreatedAt: time.Time{}, // Local instances don't track creation time
	}
}

// getHostIP attempts to get the host IP, defaulting to localhost
func getHostIP() string {
	output, err := exec.Command("hostname", "-I").Output()
	if err == nil && len(output) > 0 {
		// Take first IP
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			return fields[0]
		}
	}
	return "localhost"
}

// isProjectRunning checks if a project's containers are running
func isProjectRunning(projectID, directory string) bool {
	dockerDir := filepath.Join(directory, "supabase", "docker")
	cmd := exec.Command("docker", "compose", "-p", projectID, "ps", "-q")
	cmd.Dir = dockerDir

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// If there are container IDs in the output, project is running
	return len(strings.TrimSpace(string(output))) > 0
}

// ListInstances returns all local instances
func (p *LocalProvider) ListInstances() ([]Instance, error) {
	if err := p.reloadDatabase(); err != nil {
		return nil, err
	}

	instances := make([]Instance, 0, len(p.db.Projects))
	for name, project := range p.db.Projects {
		inst := mapProjectToInstance(name, &project)
		instances = append(instances, *inst)
	}

	return instances, nil
}

// GetInstance retrieves a specific local instance
func (p *LocalProvider) GetInstance(name string) (*Instance, error) {
	if err := p.reloadDatabase(); err != nil {
		return nil, err
	}

	project, err := p.db.GetProject(name)
	if err != nil {
		return nil, err
	}

	return mapProjectToInstance(name, project), nil
}

// CreateInstance creates a new local instance
// Note: This is a simplified version. For full setup, use the existing local add command
func (p *LocalProvider) CreateInstance(name string) (*Instance, error) {
	return nil, fmt.Errorf("creating local instances requires additional parameters. Use 'supactl local add' command instead")
}

// DeleteInstance removes a local instance from the database
// Note: This does not delete the files on disk
func (p *LocalProvider) DeleteInstance(name string) error {
	if err := p.reloadDatabase(); err != nil {
		return err
	}

	// Check if project exists
	if !p.db.ProjectExists(name) {
		return fmt.Errorf("project '%s' not found", name)
	}

	// Remove from database
	if err := p.db.RemoveProject(name); err != nil {
		return err
	}

	// Save database
	if err := local.SaveDatabase(p.db); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	return nil
}

// StartInstance starts a local instance
func (p *LocalProvider) StartInstance(name string) error {
	if err := p.reloadDatabase(); err != nil {
		return err
	}

	project, err := p.db.GetProject(name)
	if err != nil {
		return err
	}

	// Check if already running
	if isProjectRunning(name, project.Directory) {
		return fmt.Errorf("instance '%s' is already running", name)
	}

	return local.DockerComposeUp(name, project.Directory)
}

// StopInstance stops a local instance
func (p *LocalProvider) StopInstance(name string) error {
	if err := p.reloadDatabase(); err != nil {
		return err
	}

	project, err := p.db.GetProject(name)
	if err != nil {
		return err
	}

	// Check if running
	if !isProjectRunning(name, project.Directory) {
		return fmt.Errorf("instance '%s' is not running", name)
	}

	return local.DockerComposeDown(name, project.Directory)
}

// RestartInstance restarts a local instance
func (p *LocalProvider) RestartInstance(name string) error {
	if err := p.reloadDatabase(); err != nil {
		return err
	}

	project, err := p.db.GetProject(name)
	if err != nil {
		return err
	}

	dockerDir := filepath.Join(project.Directory, "supabase", "docker")

	// Run docker compose restart
	cmd := exec.Command("docker", "compose", "-p", name, "restart")
	cmd.Dir = dockerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart instance: %w", err)
	}

	return nil
}

// GetLogs retrieves logs for a local instance
func (p *LocalProvider) GetLogs(name string, lines int) (string, error) {
	if err := p.reloadDatabase(); err != nil {
		return "", err
	}

	project, err := p.db.GetProject(name)
	if err != nil {
		return "", err
	}

	dockerDir := filepath.Join(project.Directory, "supabase", "docker")

	// Run docker compose logs
	cmd := exec.Command("docker", "compose", "-p", name, "logs", "--tail", fmt.Sprintf("%d", lines))
	cmd.Dir = dockerDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}

	return string(output), nil
}

// ProviderType returns "local"
func (p *LocalProvider) ProviderType() string {
	return ProviderTypeLocal
}

// Compile-time check to ensure LocalProvider implements InstanceProvider
var _ InstanceProvider = (*LocalProvider)(nil)
