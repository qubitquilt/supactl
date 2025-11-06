package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	// Create a temporary home directory
	tempHome, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Override home directory for testing
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name      string
		serverURL string
		apiKey    string
		wantErr   bool
	}{
		{
			name:      "valid config",
			serverURL: "https://test.example.com",
			apiKey:    "test-api-key-123",
			wantErr:   false,
		},
		{
			name:      "empty server URL",
			serverURL: "",
			apiKey:    "test-key",
			wantErr:   false, // SaveConfig doesn't validate, just saves
		},
		{
			name:      "special characters in key",
			serverURL: "https://test.com",
			apiKey:    "key-with-special-chars!@#$%",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing config
			configPath, _ := GetConfigPath()
			os.RemoveAll(filepath.Dir(configPath))

			err := SaveConfig(tt.serverURL, tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file exists
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("config file was not created")
					return
				}

				// Verify file permissions
				info, err := os.Stat(configPath)
				if err != nil {
					t.Errorf("failed to stat config file: %v", err)
					return
				}

				if info.Mode().Perm() != 0600 {
					t.Errorf("config file permissions = %v, want 0600", info.Mode().Perm())
				}

				// Verify content
				data, err := os.ReadFile(configPath)
				if err != nil {
					t.Errorf("failed to read config file: %v", err)
					return
				}

				var config Config
				if err := json.Unmarshal(data, &config); err != nil {
					t.Errorf("failed to parse config: %v", err)
					return
				}

				if config.ServerURL != tt.serverURL {
					t.Errorf("ServerURL = %v, want %v", config.ServerURL, tt.serverURL)
				}
				if config.APIKey != tt.apiKey {
					t.Errorf("APIKey = %v, want %v", config.APIKey, tt.apiKey)
				}
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name      string
		setupFunc func() error
		wantURL   string
		wantKey   string
		wantErr   bool
	}{
		{
			name: "valid config",
			setupFunc: func() error {
				return SaveConfig("https://example.com", "my-api-key")
			},
			wantURL: "https://example.com",
			wantKey: "my-api-key",
			wantErr: false,
		},
		{
			name: "config does not exist",
			setupFunc: func() error {
				return nil // Don't create config
			},
			wantErr: true,
		},
		{
			name: "invalid JSON",
			setupFunc: func() error {
				configPath, _ := GetConfigPath()
				os.MkdirAll(filepath.Dir(configPath), 0700)
				return os.WriteFile(configPath, []byte("invalid json"), 0600)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			configPath, _ := GetConfigPath()
			os.RemoveAll(filepath.Dir(configPath))

			// Setup
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			config, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if config.ServerURL != tt.wantURL {
					t.Errorf("ServerURL = %v, want %v", config.ServerURL, tt.wantURL)
				}
				if config.APIKey != tt.wantKey {
					t.Errorf("APIKey = %v, want %v", config.APIKey, tt.wantKey)
				}
			}
		})
	}
}

func TestClearConfig(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name      string
		setupFunc func() error
		wantErr   bool
	}{
		{
			name: "clear existing config",
			setupFunc: func() error {
				return SaveConfig("https://example.com", "key")
			},
			wantErr: false,
		},
		{
			name: "clear non-existent config",
			setupFunc: func() error {
				return nil // No config
			},
			wantErr: false, // Should not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			configPath, _ := GetConfigPath()
			os.RemoveAll(filepath.Dir(configPath))

			// Setup
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err := ClearConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify file was deleted
			if _, err := os.Stat(configPath); !os.IsNotExist(err) {
				t.Errorf("config file still exists after ClearConfig()")
			}
		})
	}
}

func TestIsLoggedIn(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name      string
		setupFunc func() error
		want      bool
	}{
		{
			name: "logged in with valid config",
			setupFunc: func() error {
				return SaveConfig("https://example.com", "key")
			},
			want: true,
		},
		{
			name: "not logged in",
			setupFunc: func() error {
				return nil
			},
			want: false,
		},
		{
			name: "empty server URL",
			setupFunc: func() error {
				return SaveConfig("", "key")
			},
			want: false,
		},
		{
			name: "empty API key",
			setupFunc: func() error {
				return SaveConfig("https://example.com", "")
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			configPath, _ := GetConfigPath()
			os.RemoveAll(filepath.Dir(configPath))

			// Setup
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			if got := IsLoggedIn(); got != tt.want {
				t.Errorf("IsLoggedIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	// Save old HOME
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name    string
		homeDir string
		wantErr bool
	}{
		{
			name:    "valid home directory",
			homeDir: "/home/testuser",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.homeDir)

			path, err := GetConfigPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedPath := filepath.Join(tt.homeDir, configDir, configFile)
				if path != expectedPath {
					t.Errorf("GetConfigPath() = %v, want %v", path, expectedPath)
				}
			}
		})
	}
}
