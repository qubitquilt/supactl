package provider

import (
	"time"

	"github.com/qubitquilt/supactl/internal/api"
)

// RemoteProvider implements InstanceProvider for remote SupaControl server instances
type RemoteProvider struct {
	client *api.Client
}

// NewRemoteProvider creates a new remote provider
func NewRemoteProvider(serverURL, apiKey string) *RemoteProvider {
	return &RemoteProvider{
		client: api.NewClient(serverURL, apiKey),
	}
}

// mapAPIInstanceToInstance converts an API instance to a unified provider instance
func mapAPIInstanceToInstance(apiInstance *api.Instance) *Instance {
	// Parse created_at timestamp if available
	var createdAt time.Time
	if apiInstance.CreatedAt != "" {
		// Try RFC3339 format first
		t, err := time.Parse(time.RFC3339, apiInstance.CreatedAt)
		if err != nil {
			// Try alternative format
			t, err = time.Parse("2006-01-02 15:04:05", apiInstance.CreatedAt)
		}
		if err == nil {
			createdAt = t
		}
	}

	return &Instance{
		Name:        apiInstance.Name,
		Status:      apiInstance.Status,
		StudioURL:   apiInstance.StudioURL,
		APIURL:      apiInstance.APIURL,
		KongURL:     apiInstance.KongURL,
		AnonKey:     apiInstance.AnonKey,
		ServiceKey:  apiInstance.ServiceKey,
		DatabaseURL: apiInstance.DatabaseURL,
		CreatedAt:   createdAt,
	}
}

// ListInstances returns all remote instances
func (p *RemoteProvider) ListInstances() ([]Instance, error) {
	apiInstances, err := p.client.ListInstances()
	if err != nil {
		return nil, err
	}

	instances := make([]Instance, len(apiInstances))
	for i, apiInst := range apiInstances {
		instances[i] = *mapAPIInstanceToInstance(&apiInst)
	}

	return instances, nil
}

// GetInstance retrieves a specific remote instance
func (p *RemoteProvider) GetInstance(name string) (*Instance, error) {
	apiInstance, err := p.client.GetInstance(name)
	if err != nil {
		return nil, err
	}

	return mapAPIInstanceToInstance(apiInstance), nil
}

// CreateInstance creates a new remote instance
func (p *RemoteProvider) CreateInstance(name string) (*Instance, error) {
	apiInstance, err := p.client.CreateInstance(name)
	if err != nil {
		return nil, err
	}

	return mapAPIInstanceToInstance(apiInstance), nil
}

// DeleteInstance deletes a remote instance
func (p *RemoteProvider) DeleteInstance(name string) error {
	return p.client.DeleteInstance(name)
}

// StartInstance starts a remote instance
func (p *RemoteProvider) StartInstance(name string) error {
	return p.client.StartInstance(name)
}

// StopInstance stops a remote instance
func (p *RemoteProvider) StopInstance(name string) error {
	return p.client.StopInstance(name)
}

// RestartInstance restarts a remote instance
func (p *RemoteProvider) RestartInstance(name string) error {
	return p.client.RestartInstance(name)
}

// GetLogs retrieves logs for a remote instance
func (p *RemoteProvider) GetLogs(name string, lines int) (string, error) {
	return p.client.GetLogs(name, lines)
}

// ProviderType returns "remote"
func (p *RemoteProvider) ProviderType() string {
	return ProviderTypeRemote
}

// ValidateConnection validates the connection to the remote server
func (p *RemoteProvider) ValidateConnection() error {
	return p.client.LoginTest()
}

// Compile-time check to ensure RemoteProvider implements InstanceProvider
var _ InstanceProvider = (*RemoteProvider)(nil)
