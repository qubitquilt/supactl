package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDir  = ".supacontrol"
	configFile = "config.json"
)

// ContextConfig represents the configuration for a single context
type ContextConfig struct {
	Provider  string `json:"provider"`            // "local" or "remote"
	ServerURL string `json:"server_url,omitempty"` // Only for remote contexts
	APIKey    string `json:"api_key,omitempty"`    // Only for remote contexts
}

// Config represents the complete configuration with multiple contexts
type Config struct {
	CurrentContext string                    `json:"current-context"`
	Contexts       map[string]*ContextConfig `json:"contexts"`
}

// LegacyConfig represents the old configuration format (for migration)
type LegacyConfig struct {
	ServerURL string `json:"server_url"`
	APIKey    string `json:"api_key"`
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, configDir, configFile), nil
}

// LoadConfig loads the configuration from disk, handling both new and legacy formats
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config with local context
			return &Config{
				CurrentContext: "local",
				Contexts: map[string]*ContextConfig{
					"local": {
						Provider: "local",
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Try to unmarshal as new format first
	var config Config
	if err := json.Unmarshal(data, &config); err == nil && config.Contexts != nil {
		// Ensure contexts map is initialized
		if config.Contexts == nil {
			config.Contexts = make(map[string]*ContextConfig)
		}
		// Ensure local context exists
		if _, exists := config.Contexts["local"]; !exists {
			config.Contexts["local"] = &ContextConfig{Provider: "local"}
		}
		// Set default context if not set
		if config.CurrentContext == "" {
			config.CurrentContext = "local"
		}
		return &config, nil
	}

	// Try to unmarshal as legacy format for migration
	var legacyConfig LegacyConfig
	if err := json.Unmarshal(data, &legacyConfig); err == nil && legacyConfig.ServerURL != "" {
		// Migrate legacy config to new format
		config := &Config{
			CurrentContext: "default",
			Contexts: map[string]*ContextConfig{
				"local": {
					Provider: "local",
				},
				"default": {
					Provider:  "remote",
					ServerURL: legacyConfig.ServerURL,
					APIKey:    legacyConfig.APIKey,
				},
			},
		}
		// Auto-save migrated config
		if err := SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save migrated config: %v\n", err)
		}
		return config, nil
	}

	return nil, fmt.Errorf("failed to parse config file: invalid format")
}

// SaveConfig saves the configuration to disk with secure permissions
func SaveConfig(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDirPath := filepath.Join(homeDir, configDir)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDirPath, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write file with 0600 permissions (read/write for user only)
	if err := os.WriteFile(configPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetCurrentContext returns the current context configuration
func (c *Config) GetCurrentContext() (*ContextConfig, error) {
	if c.CurrentContext == "" {
		return nil, fmt.Errorf("no current context set")
	}

	ctx, exists := c.Contexts[c.CurrentContext]
	if !exists {
		return nil, fmt.Errorf("current context '%s' not found", c.CurrentContext)
	}

	return ctx, nil
}

// SetCurrentContext changes the current context
func (c *Config) SetCurrentContext(name string) error {
	if _, exists := c.Contexts[name]; !exists {
		return fmt.Errorf("context '%s' does not exist", name)
	}
	c.CurrentContext = name
	return nil
}

// AddContext adds or updates a context
func (c *Config) AddContext(name string, ctx *ContextConfig) {
	if c.Contexts == nil {
		c.Contexts = make(map[string]*ContextConfig)
	}
	c.Contexts[name] = ctx
}

// RemoveContext removes a context (cannot remove 'local' or current context)
func (c *Config) RemoveContext(name string) error {
	if name == "local" {
		return fmt.Errorf("cannot remove 'local' context")
	}
	if name == c.CurrentContext {
		return fmt.Errorf("cannot remove current context '%s'. Switch to another context first", name)
	}
	if _, exists := c.Contexts[name]; !exists {
		return fmt.Errorf("context '%s' does not exist", name)
	}
	delete(c.Contexts, name)
	return nil
}

// ListContexts returns all context names
func (c *Config) ListContexts() []string {
	names := make([]string, 0, len(c.Contexts))
	for name := range c.Contexts {
		names = append(names, name)
	}
	return names
}

// ClearConfig removes the configuration file
func ClearConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.Remove(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already cleared
		}
		return fmt.Errorf("failed to remove config file: %w", err)
	}

	return nil
}

// IsLoggedIn checks if the current context is a valid remote context with credentials
// This is used for backward compatibility with existing commands
func IsLoggedIn() bool {
	config, err := LoadConfig()
	if err != nil {
		return false
	}

	ctx, err := config.GetCurrentContext()
	if err != nil {
		return false
	}

	return ctx.Provider == "remote" && ctx.ServerURL != "" && ctx.APIKey != ""
}

// Legacy functions for backward compatibility

// SaveLegacyConfig saves a legacy-style config as a remote context named "default"
// Deprecated: Use SaveConfig with context-aware Config instead
func SaveLegacyConfig(serverURL, apiKey string) error {
	config, err := LoadConfig()
	if err != nil {
		// If config doesn't exist, create new one
		config = &Config{
			CurrentContext: "default",
			Contexts:       make(map[string]*ContextConfig),
		}
	}

	// Add/update default context
	config.Contexts["default"] = &ContextConfig{
		Provider:  "remote",
		ServerURL: serverURL,
		APIKey:    apiKey,
	}

	// Ensure local context exists
	if _, exists := config.Contexts["local"]; !exists {
		config.Contexts["local"] = &ContextConfig{Provider: "local"}
	}

	// Set current context to default
	config.CurrentContext = "default"

	return SaveConfig(config)
}
