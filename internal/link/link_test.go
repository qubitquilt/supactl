package link

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLink(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "supactl-link-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name        string
		projectName string
		wantErr     bool
	}{
		{
			name:        "valid project name",
			projectName: "my-project",
			wantErr:     false,
		},
		{
			name:        "project with hyphens",
			projectName: "my-test-project-v2",
			wantErr:     false,
		},
		{
			name:        "empty project name",
			projectName: "",
			wantErr:     false, // SaveLink doesn't validate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up link directory
			os.RemoveAll(linkDir)

			err := SaveLink(tt.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify link file was created
				linkPath := GetLinkPath()
				if _, err := os.Stat(linkPath); os.IsNotExist(err) {
					t.Errorf("link file was not created")
					return
				}

				// Verify content
				data, err := os.ReadFile(linkPath)
				if err != nil {
					t.Errorf("failed to read link file: %v", err)
					return
				}

				content := strings.TrimSpace(string(data))
				if content != tt.projectName {
					t.Errorf("link content = %v, want %v", content, tt.projectName)
				}
			}
		})
	}
}

func TestGetLink(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "supactl-link-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name        string
		setupFunc   func() error
		wantProject string
		wantErr     bool
	}{
		{
			name: "valid link file",
			setupFunc: func() error {
				return SaveLink("test-project")
			},
			wantProject: "test-project",
			wantErr:     false,
		},
		{
			name: "link file does not exist",
			setupFunc: func() error {
				return nil // Don't create link
			},
			wantErr: true,
		},
		{
			name: "empty link file",
			setupFunc: func() error {
				os.MkdirAll(linkDir, 0755)
				return os.WriteFile(GetLinkPath(), []byte(""), 0644)
			},
			wantErr: true,
		},
		{
			name: "link file with whitespace",
			setupFunc: func() error {
				os.MkdirAll(linkDir, 0755)
				return os.WriteFile(GetLinkPath(), []byte("  \n  "), 0644)
			},
			wantErr: true,
		},
		{
			name: "link file with trailing newline",
			setupFunc: func() error {
				return SaveLink("project-with-newline\n")
			},
			wantProject: "project-with-newline",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			os.RemoveAll(linkDir)

			// Setup
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			project, err := GetLink()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && project != tt.wantProject {
				t.Errorf("GetLink() = %v, want %v", project, tt.wantProject)
			}
		})
	}
}

func TestClearLink(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "supactl-link-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name      string
		setupFunc func() error
		wantErr   bool
	}{
		{
			name: "clear existing link",
			setupFunc: func() error {
				return SaveLink("test-project")
			},
			wantErr: false,
		},
		{
			name: "clear non-existent link",
			setupFunc: func() error {
				return nil
			},
			wantErr: false, // Should not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			os.RemoveAll(linkDir)

			// Setup
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err := ClearLink()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify link was deleted
			linkPath := GetLinkPath()
			if _, err := os.Stat(linkPath); !os.IsNotExist(err) {
				t.Errorf("link file still exists after ClearLink()")
			}
		})
	}
}

func TestIsLinked(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "supactl-link-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name      string
		setupFunc func() error
		want      bool
	}{
		{
			name: "is linked",
			setupFunc: func() error {
				return SaveLink("project")
			},
			want: true,
		},
		{
			name: "not linked",
			setupFunc: func() error {
				return nil
			},
			want: false,
		},
		{
			name: "empty link file",
			setupFunc: func() error {
				os.MkdirAll(linkDir, 0755)
				return os.WriteFile(GetLinkPath(), []byte(""), 0644)
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			os.RemoveAll(linkDir)

			// Setup
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			if got := IsLinked(); got != tt.want {
				t.Errorf("IsLinked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddToGitignore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "supactl-link-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name        string
		setupFunc   func() error
		wantPattern bool
	}{
		{
			name: "add to existing gitignore",
			setupFunc: func() error {
				return os.WriteFile(".gitignore", []byte("node_modules/\n*.log\n"), 0644)
			},
			wantPattern: true,
		},
		{
			name: "gitignore already contains pattern",
			setupFunc: func() error {
				return os.WriteFile(".gitignore", []byte(".supacontrol/\nnode_modules/\n"), 0644)
			},
			wantPattern: true,
		},
		{
			name: "no gitignore exists",
			setupFunc: func() error {
				return nil // No .gitignore
			},
			wantPattern: false, // Should not create .gitignore
		},
		{
			name: "empty gitignore",
			setupFunc: func() error {
				return os.WriteFile(".gitignore", []byte(""), 0644)
			},
			wantPattern: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			os.Remove(".gitignore")

			// Setup
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			// Call SaveLink which should trigger addToGitignore
			SaveLink("test-project")

			// Check if .gitignore was modified correctly
			if _, err := os.Stat(".gitignore"); err == nil {
				data, err := os.ReadFile(".gitignore")
				if err != nil {
					t.Fatalf("failed to read .gitignore: %v", err)
				}

				content := string(data)
				hasPattern := strings.Contains(content, ".supacontrol/")

				if hasPattern != tt.wantPattern {
					t.Errorf("gitignore contains pattern = %v, want %v", hasPattern, tt.wantPattern)
				}

				// Verify pattern only appears once
				if tt.wantPattern {
					count := strings.Count(content, ".supacontrol/")
					if count != 1 {
						t.Errorf("pattern appears %d times, want 1", count)
					}
				}
			} else if tt.wantPattern {
				t.Errorf(".gitignore was not created/modified")
			}
		})
	}
}

func TestGetLinkPath(t *testing.T) {
	want := filepath.Join(".", linkDir, linkFile)
	got := GetLinkPath()

	if got != want {
		t.Errorf("GetLinkPath() = %v, want %v", got, want)
	}
}
