package local

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDatabase_NewFile(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome) // For Windows

	db, err := LoadDatabase()
	if err != nil {
		t.Fatalf("LoadDatabase failed: %v", err)
	}

	if db.Projects == nil {
		t.Error("Projects map should be initialized")
	}

	if db.LastPortAssigned != BasePort {
		t.Errorf("LastPortAssigned should be %d, got %d", BasePort, db.LastPortAssigned)
	}
}

func TestSaveAndLoadDatabase(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	// Create a database
	db := &Database{
		Projects: map[string]Project{
			"test-project": {
				Directory: "/home/user/test-project",
				Ports: Ports{
					API:    54321,
					DB:     54322,
					Studio: 54323,
				},
			},
		},
		LastPortAssigned: 55321,
	}

	// Save it
	if err := SaveDatabase(db); err != nil {
		t.Fatalf("SaveDatabase failed: %v", err)
	}

	// Load it back
	loaded, err := LoadDatabase()
	if err != nil {
		t.Fatalf("LoadDatabase failed: %v", err)
	}

	// Verify
	if len(loaded.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(loaded.Projects))
	}

	project, exists := loaded.Projects["test-project"]
	if !exists {
		t.Fatal("test-project should exist")
	}

	if project.Directory != "/home/user/test-project" {
		t.Errorf("Directory mismatch: got %s", project.Directory)
	}

	if project.Ports.API != 54321 {
		t.Errorf("API port mismatch: got %d", project.Ports.API)
	}

	if loaded.LastPortAssigned != 55321 {
		t.Errorf("LastPortAssigned mismatch: got %d", loaded.LastPortAssigned)
	}
}

func TestDatabaseFilePermissions(t *testing.T) {
	// Skip on Windows as permissions work differently
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	// Create temporary home directory
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	db := &Database{
		Projects:         make(map[string]Project),
		LastPortAssigned: BasePort,
	}

	if err := SaveDatabase(db); err != nil {
		t.Fatalf("SaveDatabase failed: %v", err)
	}

	dbPath, err := GetDatabasePath()
	if err != nil {
		t.Fatalf("GetDatabasePath failed: %v", err)
	}

	info, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	mode := info.Mode()
	if mode.Perm() != 0600 {
		t.Errorf("Expected permissions 0600, got %o", mode.Perm())
	}
}

func TestProjectExists(t *testing.T) {
	db := &Database{
		Projects: map[string]Project{
			"existing": {Directory: "/home/user/existing"},
		},
		LastPortAssigned: BasePort,
	}

	if !db.ProjectExists("existing") {
		t.Error("existing project should exist")
	}

	if db.ProjectExists("non-existing") {
		t.Error("non-existing project should not exist")
	}
}

func TestAddProject(t *testing.T) {
	db := &Database{
		Projects:         make(map[string]Project),
		LastPortAssigned: BasePort,
	}

	project, err := db.AddProject("test-proj", "/home/user/test-proj")
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	if project.Directory != "/home/user/test-proj" {
		t.Errorf("Directory mismatch: got %s", project.Directory)
	}

	// Verify ports are assigned correctly
	if project.Ports.API != BasePort {
		t.Errorf("API port should be %d, got %d", BasePort, project.Ports.API)
	}

	if project.Ports.DB != BasePort+1 {
		t.Errorf("DB port should be %d, got %d", BasePort+1, project.Ports.DB)
	}

	// Verify last port assigned was incremented
	expectedNext := BasePort + PortIncrement
	if db.LastPortAssigned != expectedNext {
		t.Errorf("LastPortAssigned should be %d, got %d", expectedNext, db.LastPortAssigned)
	}

	// Verify project is in database
	if !db.ProjectExists("test-proj") {
		t.Error("Project should exist after adding")
	}
}

func TestAddProject_Duplicate(t *testing.T) {
	db := &Database{
		Projects: map[string]Project{
			"existing": {Directory: "/home/user/existing"},
		},
		LastPortAssigned: BasePort,
	}

	_, err := db.AddProject("existing", "/home/user/existing")
	if err == nil {
		t.Error("AddProject should fail for duplicate project ID")
	}
}

func TestRemoveProject(t *testing.T) {
	db := &Database{
		Projects: map[string]Project{
			"test-proj": {Directory: "/home/user/test-proj"},
		},
		LastPortAssigned: BasePort,
	}

	err := db.RemoveProject("test-proj")
	if err != nil {
		t.Fatalf("RemoveProject failed: %v", err)
	}

	if db.ProjectExists("test-proj") {
		t.Error("Project should not exist after removal")
	}
}

func TestRemoveProject_NotFound(t *testing.T) {
	db := &Database{
		Projects:         make(map[string]Project),
		LastPortAssigned: BasePort,
	}

	err := db.RemoveProject("non-existing")
	if err == nil {
		t.Error("RemoveProject should fail for non-existing project")
	}
}

func TestGetProject(t *testing.T) {
	db := &Database{
		Projects: map[string]Project{
			"test-proj": {
				Directory: "/home/user/test-proj",
				Ports:     Ports{API: 54321},
			},
		},
		LastPortAssigned: BasePort,
	}

	project, err := db.GetProject("test-proj")
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}

	if project.Directory != "/home/user/test-proj" {
		t.Errorf("Directory mismatch: got %s", project.Directory)
	}

	if project.Ports.API != 54321 {
		t.Errorf("API port mismatch: got %d", project.Ports.API)
	}
}

func TestGetProject_NotFound(t *testing.T) {
	db := &Database{
		Projects:         make(map[string]Project),
		LastPortAssigned: BasePort,
	}

	_, err := db.GetProject("non-existing")
	if err == nil {
		t.Error("GetProject should fail for non-existing project")
	}
}

func TestGetDatabasePath(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	path, err := GetDatabasePath()
	if err != nil {
		t.Fatalf("GetDatabasePath failed: %v", err)
	}

	expected := filepath.Join(tmpHome, databaseFile)
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}
