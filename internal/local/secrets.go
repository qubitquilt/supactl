package local

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GeneratePassword generates a secure random password of the specified length
// Uses base64 encoding to ensure uniform distribution of random values
func GeneratePassword(length int) (string, error) {
	// Calculate number of random bytes needed
	// base64 encoding produces ~4/3 the output length, so we need ~3/4 input bytes
	numBytes := (length * 3) / 4
	if numBytes == 0 {
		numBytes = 1
	}

	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 and trim to desired length
	encoded := base64.RawURLEncoding.EncodeToString(b)
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}

// GenerateEncryptionKey generates a secure random encryption key (32 bytes for AES-256)
func GenerateEncryptionKey() (string, error) {
	return GeneratePassword(32)
}

// GenerateJWT generates a JWT token using HS256 algorithm with the golang-jwt library
func GenerateJWT(jwtSecret, role string) (string, error) {
	now := time.Now()
	exp := now.Add(10 * 365 * 24 * time.Hour) // 10 years from now

	// Create claims
	claims := jwt.MapClaims{
		"aud":  "authenticated",
		"exp":  exp.Unix(),
		"iat":  now.Unix(),
		"iss":  "supabase",
		"ref":  "localhost",
		"role": role,
		"sub":  "1234567890",
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

// Secrets holds all generated secrets for a Supabase instance
type Secrets struct {
	PostgresPassword  string
	JWTSecret         string
	DashboardPassword string
	VaultEncKey       string
	AnonKey           string
	ServiceRoleKey    string
}

// GenerateSecrets generates all required secrets for a new Supabase instance
func GenerateSecrets() (*Secrets, error) {
	postgresPassword, err := GeneratePassword(40)
	if err != nil {
		return nil, fmt.Errorf("failed to generate postgres password: %w", err)
	}

	jwtSecret, err := GeneratePassword(40)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
	}

	dashboardPassword, err := GeneratePassword(40)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dashboard password: %w", err)
	}

	vaultEncKey, err := GenerateEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate vault encryption key: %w", err)
	}

	anonKey, err := GenerateJWT(jwtSecret, "anon")
	if err != nil {
		return nil, fmt.Errorf("failed to generate anon key: %w", err)
	}

	serviceRoleKey, err := GenerateJWT(jwtSecret, "service_role")
	if err != nil {
		return nil, fmt.Errorf("failed to generate service role key: %w", err)
	}

	return &Secrets{
		PostgresPassword:  postgresPassword,
		JWTSecret:         jwtSecret,
		DashboardPassword: dashboardPassword,
		VaultEncKey:       vaultEncKey,
		AnonKey:           anonKey,
		ServiceRoleKey:    serviceRoleKey,
	}, nil
}
