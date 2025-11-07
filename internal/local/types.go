package local

// Ports represents all port configurations for a local Supabase instance
type Ports struct {
	API       int `json:"api"`
	DB        int `json:"db"`
	Shadow    int `json:"shadow"`
	Studio    int `json:"studio"`
	Inbucket  int `json:"inbucket"`
	SMTP      int `json:"smtp"`
	POP3      int `json:"pop3"`
	Pooler    int `json:"pooler"`
	Analytics int `json:"analytics"`
	KongHTTPS int `json:"kong_https"`
}

// Project represents a local Supabase project configuration
type Project struct {
	Directory string `json:"directory"`
	Ports     Ports  `json:"ports"`
}

// Database represents the local projects database structure
type Database struct {
	Projects         map[string]Project `json:"projects"`
	LastPortAssigned int                `json:"last_port_assigned"`
}

// Constants for port management
const (
	BasePort      = 54321
	PortIncrement = 1000
)
