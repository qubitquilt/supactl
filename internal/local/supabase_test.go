package local

import (
	"testing"
)

func TestValidateProjectID(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		wantErr   bool
	}{
		{"valid lowercase", "myproject", false},
		{"valid with numbers", "project123", false},
		{"valid with hyphen", "my-project", false},
		{"valid with underscore", "my_project", false},
		{"valid complex", "proj-123_test", false},
		{"valid single char", "a", false},
		{"valid single number", "1", false},
		{"valid ends with hyphen", "project-", false}, // supascale.sh allows this
		{"valid ends with underscore", "project_", false}, // supascale.sh allows this
		{"invalid uppercase", "MyProject", true},
		{"invalid starts with hyphen", "-project", true},
		{"invalid starts with underscore", "_project", true},
		{"invalid with dot", "my.project", true},
		{"invalid with space", "my project", true},
		{"invalid with special char", "my@project", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectID(tt.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectID(%q) error = %v, wantErr %v", tt.projectID, err, tt.wantErr)
			}
		})
	}
}
