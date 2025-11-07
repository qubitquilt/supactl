package local

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GeneratePassword generates a secure random password of the specified length
func GeneratePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

// GenerateEncryptionKey generates a secure random encryption key (32 bytes for AES-256)
func GenerateEncryptionKey() (string, error) {
	return GeneratePassword(32)
}

// base64URLEncode encodes data to base64 URL format (no padding)
func base64URLEncode(data []byte) string {
	encoded := base64.RawURLEncoding.EncodeToString(data)
	return strings.TrimRight(encoded, "=")
}

// JWTHeader represents the JWT header
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// JWTPayload represents the JWT payload
type JWTPayload struct {
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Iss string `json:"iss"`
	Ref string `json:"ref"`
	Role string `json:"role"`
	Sub string `json:"sub"`
}

// GenerateJWT generates a JWT token using HS256 algorithm
func GenerateJWT(jwtSecret, role string) (string, error) {
	// Create header
	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}
	headerEncoded := base64URLEncode(headerJSON)

	// Create payload
	now := time.Now().Unix()
	exp := now + (315360000) // 10 years from now
	payload := JWTPayload{
		Aud:  "authenticated",
		Exp:  exp,
		Iat:  now,
		Iss:  "supabase",
		Ref:  "localhost",
		Role: role,
		Sub:  "1234567890",
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	payloadEncoded := base64URLEncode(payloadJSON)

	// Create signature
	signatureInput := headerEncoded + "." + payloadEncoded
	h := hmac.New(sha256.New, []byte(jwtSecret))
	h.Write([]byte(signatureInput))
	signature := base64URLEncode(h.Sum(nil))

	// Combine all parts
	token := signatureInput + "." + signature
	return token, nil
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
