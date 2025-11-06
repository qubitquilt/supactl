package cmd

import (
	"regexp"
	"testing"
)

func TestProjectNameValidation(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantValid   bool
	}{
		{
			name:        "valid lowercase",
			projectName: "myproject",
			wantValid:   true,
		},
		{
			name:        "valid with hyphens",
			projectName: "my-project",
			wantValid:   true,
		},
		{
			name:        "valid with numbers",
			projectName: "project123",
			wantValid:   true,
		},
		{
			name:        "valid complex name",
			projectName: "my-project-v2",
			wantValid:   true,
		},
		{
			name:        "single character",
			projectName: "a",
			wantValid:   true,
		},
		{
			name:        "single number",
			projectName: "1",
			wantValid:   true,
		},
		{
			name:        "invalid uppercase",
			projectName: "MyProject",
			wantValid:   false,
		},
		{
			name:        "invalid starts with hyphen",
			projectName: "-myproject",
			wantValid:   false,
		},
		{
			name:        "invalid ends with hyphen",
			projectName: "myproject-",
			wantValid:   false,
		},
		{
			name:        "invalid special characters",
			projectName: "my_project",
			wantValid:   false,
		},
		{
			name:        "invalid spaces",
			projectName: "my project",
			wantValid:   false,
		},
		{
			name:        "invalid empty",
			projectName: "",
			wantValid:   false,
		},
		{
			name:        "invalid dots",
			projectName: "my.project",
			wantValid:   false,
		},
		{
			name:        "valid multiple hyphens",
			projectName: "my-super-long-project-name",
			wantValid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test against the regex from create.go
			regex := regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)
			valid := regex.MatchString(tt.projectName)

			if valid != tt.wantValid {
				t.Errorf("projectName %q validation = %v, want %v", tt.projectName, valid, tt.wantValid)
			}
		})
	}
}
