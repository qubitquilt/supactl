package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/qubitquilt/supactl/internal/testutil"
)

func TestNewClient(t *testing.T) {
	serverURL := "https://example.com"
	apiKey := "test-key"

	client := NewClient(serverURL, apiKey)

	if client.ServerURL != serverURL {
		t.Errorf("ServerURL = %v, want %v", client.ServerURL, serverURL)
	}

	if client.APIKey != apiKey {
		t.Errorf("APIKey = %v, want %v", client.APIKey, apiKey)
	}

	if client.HTTPClient == nil {
		t.Error("HTTPClient is nil")
	}
}

func TestLoginTest(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   interface{}
		wantErr    bool
	}{
		{
			name:       "successful authentication",
			statusCode: http.StatusOK,
			response: AuthResponse{
				User: struct {
					ID       string `json:"id"`
					Email    string `json:"email"`
					Username string `json:"username"`
				}{
					ID:       "123",
					Email:    "test@example.com",
					Username: "testuser",
				},
				Authenticated: true,
			},
			wantErr: false,
		},
		{
			name:       "authentication failed",
			statusCode: http.StatusOK,
			response: AuthResponse{
				Authenticated: false,
			},
			wantErr: true,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response: ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid API key",
				Code:    401,
			},
			wantErr: true,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			response: ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Something went wrong",
				Code:    500,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.NewMockServer()
			defer server.Close()

			server.On("GET", "/api/v1/auth/me", func(w http.ResponseWriter, r *http.Request) {
				// Verify Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-key" {
					t.Errorf("Authorization header = %v, want %v", authHeader, "Bearer test-key")
				}

				testutil.RespondJSON(w, tt.statusCode, tt.response)
			})

			client := NewClient(server.URL(), "test-key")
			err := client.LoginTest()

			if (err != nil) != tt.wantErr {
				t.Errorf("LoginTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListInstances(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   interface{}
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "list with instances",
			statusCode: http.StatusOK,
			response: ListInstancesResponse{
				Instances: []Instance{
					{
						Name:      "project-1",
						Status:    "running",
						StudioURL: "https://studio.example.com",
					},
					{
						Name:      "project-2",
						Status:    "stopped",
						StudioURL: "https://studio2.example.com",
					},
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:       "empty list",
			statusCode: http.StatusOK,
			response: ListInstancesResponse{
				Instances: []Instance{},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response: ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid credentials",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.NewMockServer()
			defer server.Close()

			server.On("GET", "/api/v1/instances", func(w http.ResponseWriter, r *http.Request) {
				testutil.RespondJSON(w, tt.statusCode, tt.response)
			})

			client := NewClient(server.URL(), "test-key")
			instances, err := client.ListInstances()

			if (err != nil) != tt.wantErr {
				t.Errorf("ListInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(instances) != tt.wantCount {
				t.Errorf("ListInstances() returned %d instances, want %d", len(instances), tt.wantCount)
			}
		})
	}
}

func TestCreateInstance(t *testing.T) {
	tests := []struct {
		name         string
		projectName  string
		statusCode   int
		response     interface{}
		wantErr      bool
		wantInstance *Instance
	}{
		{
			name:        "successful creation",
			projectName: "new-project",
			statusCode:  http.StatusCreated,
			response: Instance{
				Name:      "new-project",
				Status:    "creating",
				StudioURL: "https://studio.example.com",
				APIURL:    "https://api.example.com",
			},
			wantErr: false,
			wantInstance: &Instance{
				Name:      "new-project",
				Status:    "creating",
				StudioURL: "https://studio.example.com",
				APIURL:    "https://api.example.com",
			},
		},
		{
			name:        "invalid project name",
			projectName: "Invalid-Name",
			statusCode:  http.StatusBadRequest,
			response: ErrorResponse{
				Error:   "Bad Request",
				Message: "Project name must be lowercase",
			},
			wantErr: true,
		},
		{
			name:        "project already exists",
			projectName: "existing-project",
			statusCode:  http.StatusConflict,
			response: ErrorResponse{
				Error:   "Conflict",
				Message: "Project already exists",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.NewMockServer()
			defer server.Close()

			server.On("POST", "/api/v1/instances", func(w http.ResponseWriter, r *http.Request) {
				// Verify request body
				var reqBody CreateInstanceRequest
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				if reqBody.Name != tt.projectName {
					t.Errorf("request name = %v, want %v", reqBody.Name, tt.projectName)
				}

				testutil.RespondJSON(w, tt.statusCode, tt.response)
			})

			client := NewClient(server.URL(), "test-key")
			instance, err := client.CreateInstance(tt.projectName)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if instance.Name != tt.wantInstance.Name {
					t.Errorf("instance.Name = %v, want %v", instance.Name, tt.wantInstance.Name)
				}
				if instance.Status != tt.wantInstance.Status {
					t.Errorf("instance.Status = %v, want %v", instance.Status, tt.wantInstance.Status)
				}
			}
		})
	}
}

func TestDeleteInstance(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		statusCode  int
		response    interface{}
		wantErr     bool
	}{
		{
			name:        "successful deletion",
			projectName: "project-to-delete",
			statusCode:  http.StatusNoContent,
			response:    nil,
			wantErr:     false,
		},
		{
			name:        "project not found",
			projectName: "non-existent",
			statusCode:  http.StatusNotFound,
			response: ErrorResponse{
				Error:   "Not Found",
				Message: "Project not found",
			},
			wantErr: true,
		},
		{
			name:        "deletion forbidden",
			projectName: "protected-project",
			statusCode:  http.StatusForbidden,
			response: ErrorResponse{
				Error:   "Forbidden",
				Message: "Cannot delete this project",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.NewMockServer()
			defer server.Close()

			server.On("DELETE", "/api/v1/instances/"+tt.projectName, func(w http.ResponseWriter, r *http.Request) {
				if tt.response == nil {
					w.WriteHeader(tt.statusCode)
				} else {
					testutil.RespondJSON(w, tt.statusCode, tt.response)
				}
			})

			client := NewClient(server.URL(), "test-key")
			err := client.DeleteInstance(tt.projectName)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetInstance(t *testing.T) {
	tests := []struct {
		name         string
		projectName  string
		statusCode   int
		response     interface{}
		wantErr      bool
		wantInstance *Instance
	}{
		{
			name:        "get existing instance",
			projectName: "my-project",
			statusCode:  http.StatusOK,
			response: Instance{
				Name:      "my-project",
				Status:    "running",
				StudioURL: "https://studio.example.com",
				APIURL:    "https://api.example.com",
				KongURL:   "https://kong.example.com",
				AnonKey:   "anon-key-123",
				CreatedAt: "2025-01-01T00:00:00Z",
			},
			wantErr: false,
			wantInstance: &Instance{
				Name:      "my-project",
				Status:    "running",
				StudioURL: "https://studio.example.com",
				APIURL:    "https://api.example.com",
				KongURL:   "https://kong.example.com",
				AnonKey:   "anon-key-123",
				CreatedAt: "2025-01-01T00:00:00Z",
			},
		},
		{
			name:        "instance not found",
			projectName: "non-existent",
			statusCode:  http.StatusNotFound,
			response: ErrorResponse{
				Error:   "Not Found",
				Message: "Instance not found",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.NewMockServer()
			defer server.Close()

			server.On("GET", "/api/v1/instances/"+tt.projectName, func(w http.ResponseWriter, r *http.Request) {
				testutil.RespondJSON(w, tt.statusCode, tt.response)
			})

			client := NewClient(server.URL(), "test-key")
			instance, err := client.GetInstance(tt.projectName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if instance.Name != tt.wantInstance.Name {
					t.Errorf("instance.Name = %v, want %v", instance.Name, tt.wantInstance.Name)
				}
				if instance.Status != tt.wantInstance.Status {
					t.Errorf("instance.Status = %v, want %v", instance.Status, tt.wantInstance.Status)
				}
				if instance.AnonKey != tt.wantInstance.AnonKey {
					t.Errorf("instance.AnonKey = %v, want %v", instance.AnonKey, tt.wantInstance.AnonKey)
				}
			}
		})
	}
}

func TestHandleErrorResponse(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		wantErrMsg   string
	}{
		{
			name:         "error with message",
			statusCode:   http.StatusBadRequest,
			responseBody: `{"message": "Invalid request", "error": "Bad Request"}`,
			wantErrMsg:   "Invalid request",
		},
		{
			name:         "error with only error field",
			statusCode:   http.StatusUnauthorized,
			responseBody: `{"error": "Unauthorized"}`,
			wantErrMsg:   "Unauthorized",
		},
		{
			name:         "invalid JSON",
			statusCode:   http.StatusInternalServerError,
			responseBody: `invalid json`,
			wantErrMsg:   "HTTP 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.NewMockServer()
			defer server.Close()

			server.On("GET", "/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			})

			client := NewClient(server.URL(), "test-key")
			resp, err := client.makeRequest("GET", "/test", nil)
			if err != nil {
				t.Fatalf("makeRequest failed: %v", err)
			}

			err = client.handleErrorResponse(resp)
			if err == nil {
				t.Error("expected error, got nil")
				return
			}

			if !contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("error message = %v, want to contain %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
