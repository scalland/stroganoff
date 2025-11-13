package auth

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/gocr/internal/config"
)

// Token represents an authentication token
type Token struct {
	Value     string
	ExpiresAt time.Time
	Scopes    []string
}

// Authenticator handles API authentication
type Authenticator struct {
	tokens map[string]*Token
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator() *Authenticator {
	return &Authenticator{
		tokens: make(map[string]*Token),
	}
}

// ValidateToken validates an authentication token
func (a *Authenticator) ValidateToken(token string) bool {
	cfg := config.GetInstance().GetAPI()

	if !cfg.AuthEnabled {
		return true
	}

	if token == "" {
		return false
	}

	// Check if token exists and is not expired
	if t, ok := a.tokens[token]; ok {
		if time.Now().Before(t.ExpiresAt) {
			return true
		}
		delete(a.tokens, token)
	}

	return false
}

// CreateToken creates a new authentication token
func (a *Authenticator) CreateToken(scopes []string, duration time.Duration) string {
	token := generateToken()
	a.tokens[token] = &Token{
		Value:     token,
		ExpiresAt: time.Now().Add(duration),
		Scopes:    scopes,
	}
	return token
}

// RevokeToken revokes a token
func (a *Authenticator) RevokeToken(token string) {
	delete(a.tokens, token)
}

// HasScope checks if a token has a specific scope
func (a *Authenticator) HasScope(token, scope string) bool {
	if t, ok := a.tokens[token]; ok {
		for _, s := range t.Scopes {
			if s == scope {
				return true
			}
		}
	}
	return false
}

// generateToken generates a random token
func generateToken() string {
	timestamp := time.Now().UnixNano()
	random := fmt.Sprintf("%d-%s", timestamp, time.Now().String())
	hash := sha256.Sum256([]byte(random))
	return fmt.Sprintf("%x", hash)
}

// ExtractToken extracts token from authorization header
func ExtractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
