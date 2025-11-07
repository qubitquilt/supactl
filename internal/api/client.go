package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Instance lifecycle action constants
const (
	instanceActionStart   = "start"
	instanceActionStop    = "stop"
	instanceActionRestart = "restart"
)

// Client is the API client for SupaControl server
type Client struct {
	ServerURL  string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(serverURL, apiKey string) *Client {
	return &Client{
		ServerURL: serverURL,
		APIKey:    apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest is a helper function to make HTTP requests
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.ServerURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// handleErrorResponse parses and returns a user-friendly error message
func (c *Client) handleErrorResponse(resp *http.Response) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response", resp.StatusCode)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if errResp.Message != "" {
		return fmt.Errorf("%s", errResp.Message)
	}
	if errResp.Error != "" {
		return fmt.Errorf("%s", errResp.Error)
	}

	return fmt.Errorf("HTTP %d: request failed", resp.StatusCode)
}

// LoginTest validates the API key and server URL
func (c *Client) LoginTest() error {
	resp, err := c.makeRequest("GET", "/api/v1/auth/me", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	if !authResp.Authenticated {
		return fmt.Errorf("authentication failed")
	}

	return nil
}

// ListInstances retrieves all instances
func (c *Client) ListInstances() ([]Instance, error) {
	resp, err := c.makeRequest("GET", "/api/v1/instances", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var listResp ListInstancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to parse instances list: %w", err)
	}

	return listResp.Instances, nil
}

// CreateInstance creates a new instance
func (c *Client) CreateInstance(name string) (*Instance, error) {
	reqBody := CreateInstanceRequest{Name: name}

	resp, err := c.makeRequest("POST", "/api/v1/instances", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("failed to parse instance response: %w", err)
	}

	return &instance, nil
}

// DeleteInstance deletes an instance
func (c *Client) DeleteInstance(name string) error {
	endpoint := fmt.Sprintf("/api/v1/instances/%s", name)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// GetInstance retrieves details about a specific instance
func (c *Client) GetInstance(name string) (*Instance, error) {
	endpoint := fmt.Sprintf("/api/v1/instances/%s", name)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("failed to parse instance response: %w", err)
	}

	return &instance, nil
}

// instanceAction performs a lifecycle action (start, stop, restart) on an instance
func (c *Client) instanceAction(name, action string) error {
	endpoint := fmt.Sprintf("/api/v1/instances/%s/%s", name, action)
	resp, err := c.makeRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// StartInstance starts a stopped instance
func (c *Client) StartInstance(name string) error {
	return c.instanceAction(name, instanceActionStart)
}

// StopInstance stops a running instance
func (c *Client) StopInstance(name string) error {
	return c.instanceAction(name, instanceActionStop)
}

// RestartInstance restarts an instance
func (c *Client) RestartInstance(name string) error {
	return c.instanceAction(name, instanceActionRestart)
}

// GetLogs retrieves logs for an instance
func (c *Client) GetLogs(name string, lines int) (string, error) {
	endpoint := fmt.Sprintf("/api/v1/instances/%s/logs?lines=%d", name, lines)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.handleErrorResponse(resp)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(bodyBytes), nil
}
