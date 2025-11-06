package api

// Instance represents a Supabase instance managed by SupaControl
type Instance struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	StudioURL    string `json:"studio_url"`
	APIURL       string `json:"api_url"`
	KongURL      string `json:"kong_url"`
	AnonKey      string `json:"anon_key,omitempty"`
	ServiceKey   string `json:"service_key,omitempty"`
	DatabaseURL  string `json:"database_url,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// ListInstancesResponse represents the response from the list instances endpoint
type ListInstancesResponse struct {
	Instances []Instance `json:"instances"`
}

// CreateInstanceRequest represents a request to create a new instance
type CreateInstanceRequest struct {
	Name string `json:"name"`
}

// AuthResponse represents the response from the auth/me endpoint
type AuthResponse struct {
	User struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"user"`
	Authenticated bool `json:"authenticated"`
}
