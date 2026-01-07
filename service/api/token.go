package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// Secret key for signing JWT tokens (in production, this should be in environment variables)
var jwtKey = []byte("your-very-secret-key-change-this-in-production")

// tokenClaims represents the JWT payload
type tokenClaims struct {
	Sub string `json:"sub"` // Subject (user ID)
	Exp int64  `json:"exp"` // Expiration time
}

// sign creates an HMAC signature for the given data
func sign(data string) []byte {
	h := hmac.New(sha256.New, jwtKey)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// createToken generates a new JWT token for the given user ID
func createToken(userID string) (string, error) {
	// JWT header
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	// JWT payload
	claims := tokenClaims{
		Sub: userID,
		Exp: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Create signature
	sig := sign(header + "." + payload)
	signature := base64.RawURLEncoding.EncodeToString(sig)

	// Combine all parts
	token := header + "." + payload + "." + signature
	return token, nil
}

// ParseToken validates and parses a JWT token, returning the user ID
func ParseToken(tokenString string) (string, error) {
	// Split token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", ErrUnauthorized
	}

	header, payload, sigB64 := parts[0], parts[1], parts[2]

	// Decode and verify signature
	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return "", ErrUnauthorized
	}

	expected := sign(header + "." + payload)
	if !hmac.Equal(sig, expected) {
		return "", ErrUnauthorized
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", ErrUnauthorized
	}

	// Parse claims
	var claims tokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", ErrUnauthorized
	}

	// Check expiration and validity
	if claims.Exp < time.Now().Unix() {
		return "", ErrUnauthorized
	}

	if claims.Sub == "" {
		return "", ErrUnauthorized
	}

	return claims.Sub, nil
}
