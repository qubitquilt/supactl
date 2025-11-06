package cmd

import (
	"os"
	"testing"

	"github.com/yourusername/supactl/internal/auth"
)

func TestGetAPIClientNotLoggedIn(t *testing.T) {
	// Create temporary home directory
	tempHome, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Override HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	// Ensure no config exists
	auth.ClearConfig()

	// This test verifies that getAPIClient would exit if called
	// We can't directly test the os.Exit call, but we can verify
	// the config doesn't exist
	if auth.IsLoggedIn() {
		t.Error("expected not to be logged in")
	}
}

func TestGetAPIClientLoggedIn(t *testing.T) {
	// Create temporary home directory
	tempHome, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Override HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	// Save config
	err = auth.SaveConfig("https://example.com", "test-key")
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Get client
	client := getAPIClient()

	if client == nil {
		t.Error("expected client, got nil")
	}

	if client.ServerURL != "https://example.com" {
		t.Errorf("ServerURL = %v, want %v", client.ServerURL, "https://example.com")
	}

	if client.APIKey != "test-key" {
		t.Errorf("APIKey = %v, want %v", client.APIKey, "test-key")
	}
}
