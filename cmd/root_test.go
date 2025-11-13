package cmd

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/qubitquilt/supactl/internal/auth"
)

func TestGetAPIClientNotLoggedIn(t *testing.T) {
	// Create temporary home directory
	tempHome, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Override HOME and USERPROFILE
	oldHome := os.Getenv("HOME")
	oldUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempHome)
	os.Setenv("USERPROFILE", tempHome)
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("USERPROFILE", oldUserProfile)
	}()

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

	// Override HOME and USERPROFILE
	oldHome := os.Getenv("HOME")
	oldUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempHome)
	os.Setenv("USERPROFILE", tempHome)
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("USERPROFILE", oldUserProfile)
		// Give time for file handles to be released on Windows
		if runtime.GOOS == "windows" {
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// Save config using legacy function
	err = auth.SaveLegacyConfig("https://example.com", "test-key")
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Give time for file to be fully written on Windows
	if runtime.GOOS == "windows" {
		time.Sleep(50 * time.Millisecond)
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
