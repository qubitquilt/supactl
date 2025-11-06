package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory for testing and returns cleanup function
func TempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "supactl-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup
}

// MockServer creates a mock HTTP server for testing API calls
type MockServer struct {
	Server   *httptest.Server
	Handlers map[string]http.HandlerFunc
}

// NewMockServer creates a new mock server
func NewMockServer() *MockServer {
	ms := &MockServer{
		Handlers: make(map[string]http.HandlerFunc),
	}

	ms.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		if handler, ok := ms.Handlers[key]; ok {
			handler(w, r)
			return
		}

		// Default 404 response
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Not Found",
			"message": "Endpoint not found",
			"code":    404,
		})
	}))

	return ms
}

// Close closes the mock server
func (ms *MockServer) Close() {
	ms.Server.Close()
}

// URL returns the mock server URL
func (ms *MockServer) URL() string {
	return ms.Server.URL
}

// On registers a handler for a specific HTTP method and path
func (ms *MockServer) On(method, path string, handler http.HandlerFunc) {
	ms.Handlers[method+" "+path] = handler
}

// RespondJSON is a helper to respond with JSON
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// RespondError is a helper to respond with an error
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, map[string]interface{}{
		"error":   http.StatusText(statusCode),
		"message": message,
		"code":    statusCode,
	})
}

// CreateTestFile creates a test file with content
func CreateTestFile(t *testing.T, dir, filename, content string) string {
	t.Helper()
	filePath := filepath.Join(dir, filename)

	// Create parent directory if needed
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	return filePath
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads a file and returns its content
func ReadFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	return string(content)
}

// SetEnvVar sets an environment variable and returns cleanup function
func SetEnvVar(t *testing.T, key, value string) func() {
	t.Helper()
	oldValue, existed := os.LookupEnv(key)

	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}

	return func() {
		if existed {
			os.Setenv(key, oldValue)
		} else {
			os.Unsetenv(key)
		}
	}
}
