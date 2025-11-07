package local

import (
	"strings"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword(40)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}

	if len(password) != 40 {
		t.Errorf("Expected password length 40, got %d", len(password))
	}

	// Verify it contains only base64 URL-safe characters (alphanumeric, -, _)
	for _, c := range password {
		isValid := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
		if !isValid {
			t.Errorf("Password contains invalid character for base64 URL encoding: %c", c)
		}
	}
}

func TestGeneratePassword_Uniqueness(t *testing.T) {
	// Generate multiple passwords and ensure they're unique
	passwords := make(map[string]bool)
	for i := 0; i < 100; i++ {
		password, err := GeneratePassword(40)
		if err != nil {
			t.Fatalf("GeneratePassword failed: %v", err)
		}
		if passwords[password] {
			t.Error("Generated duplicate password")
		}
		passwords[password] = true
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey failed: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected key length 32 (for AES-256), got %d", len(key))
	}
}

func TestGenerateJWT(t *testing.T) {
	jwtSecret := "test-secret-key-for-jwt-generation"
	role := "anon"

	token, err := GenerateJWT(jwtSecret, role)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	// Verify JWT structure (should have 3 parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("JWT should have 3 parts, got %d", len(parts))
	}

	// Verify each part is not empty
	for i, part := range parts {
		if part == "" {
			t.Errorf("JWT part %d is empty", i)
		}
	}
}

func TestGenerateJWT_DifferentRoles(t *testing.T) {
	jwtSecret := "test-secret-key"

	anonToken, err := GenerateJWT(jwtSecret, "anon")
	if err != nil {
		t.Fatalf("GenerateJWT failed for anon: %v", err)
	}

	serviceToken, err := GenerateJWT(jwtSecret, "service_role")
	if err != nil {
		t.Fatalf("GenerateJWT failed for service_role: %v", err)
	}

	// Tokens should be different for different roles
	if anonToken == serviceToken {
		t.Error("Tokens for different roles should be different")
	}
}

func TestGenerateSecrets(t *testing.T) {
	secrets, err := GenerateSecrets()
	if err != nil {
		t.Fatalf("GenerateSecrets failed: %v", err)
	}

	// Verify all secrets are generated
	if secrets.PostgresPassword == "" {
		t.Error("PostgresPassword should not be empty")
	}
	if secrets.JWTSecret == "" {
		t.Error("JWTSecret should not be empty")
	}
	if secrets.DashboardPassword == "" {
		t.Error("DashboardPassword should not be empty")
	}
	if secrets.VaultEncKey == "" {
		t.Error("VaultEncKey should not be empty")
	}
	if secrets.AnonKey == "" {
		t.Error("AnonKey should not be empty")
	}
	if secrets.ServiceRoleKey == "" {
		t.Error("ServiceRoleKey should not be empty")
	}

	// Verify password lengths
	if len(secrets.PostgresPassword) != 40 {
		t.Errorf("PostgresPassword length should be 40, got %d", len(secrets.PostgresPassword))
	}
	if len(secrets.JWTSecret) != 40 {
		t.Errorf("JWTSecret length should be 40, got %d", len(secrets.JWTSecret))
	}
	if len(secrets.DashboardPassword) != 40 {
		t.Errorf("DashboardPassword length should be 40, got %d", len(secrets.DashboardPassword))
	}
	if len(secrets.VaultEncKey) != 32 {
		t.Errorf("VaultEncKey length should be 32, got %d", len(secrets.VaultEncKey))
	}

	// Verify JWT keys have correct structure
	if len(strings.Split(secrets.AnonKey, ".")) != 3 {
		t.Error("AnonKey should be a valid JWT")
	}
	if len(strings.Split(secrets.ServiceRoleKey, ".")) != 3 {
		t.Error("ServiceRoleKey should be a valid JWT")
	}

	// Verify all secrets are unique
	values := []string{
		secrets.PostgresPassword,
		secrets.JWTSecret,
		secrets.DashboardPassword,
		secrets.VaultEncKey,
	}
	uniqueMap := make(map[string]bool)
	for _, val := range values {
		if uniqueMap[val] {
			t.Error("Duplicate secret value found")
		}
		uniqueMap[val] = true
	}
}
