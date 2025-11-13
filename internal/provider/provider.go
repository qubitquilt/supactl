package provider

import "time"

// Instance represents a unified Supabase instance across both remote and local providers.
// This abstraction allows the CLI to work with instances regardless of their backend.
type Instance struct {
	// Core fields (common to both remote and local)
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StudioURL string    `json:"studio_url"`
	APIURL    string    `json:"api_url"`
	CreatedAt time.Time `json:"created_at,omitempty"`

	// Remote-specific fields (optional, populated only for remote instances)
	KongURL     string `json:"kong_url,omitempty"`
	AnonKey     string `json:"anon_key,omitempty"`
	ServiceKey  string `json:"service_key,omitempty"`
	DatabaseURL string `json:"database_url,omitempty"`

	// Local-specific fields (optional, populated only for local instances)
	Directory string `json:"directory,omitempty"`
	DBPort    int    `json:"db_port,omitempty"`
}

// InstanceProvider defines the abstract contract for managing Supabase instances.
// This interface is implemented by RemoteProvider (SupaControl API) and LocalProvider (Docker).
type InstanceProvider interface {
	// ListInstances returns all instances managed by this provider
	ListInstances() ([]Instance, error)

	// GetInstance retrieves detailed information about a specific instance
	GetInstance(name string) (*Instance, error)

	// CreateInstance creates a new instance with the given name
	CreateInstance(name string) (*Instance, error)

	// DeleteInstance permanently deletes an instance
	DeleteInstance(name string) error

	// StartInstance starts a stopped instance
	StartInstance(name string) error

	// StopInstance stops a running instance
	StopInstance(name string) error

	// RestartInstance restarts an instance (stop + start)
	RestartInstance(name string) error

	// GetLogs retrieves the most recent logs for an instance
	GetLogs(name string, lines int) (string, error)

	// ProviderType returns the type of this provider ("remote" or "local")
	ProviderType() string
}

// ProviderType constants
const (
	ProviderTypeRemote = "remote"
	ProviderTypeLocal  = "local"
)
