package auth

import (
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	a := NewAuthenticator()

	token := a.CreateToken([]string{"read", "write"}, 1*time.Hour)

	if token == "" {
		t.Fatal("CreateToken returned empty token")
	}

	if !a.ValidateToken(token) {
		t.Fatal("ValidateToken should return true for created token")
	}
}

func TestTokenExpiration(t *testing.T) {
	a := NewAuthenticator()

	token := a.CreateToken([]string{"read"}, -1*time.Second) // Already expired

	if a.ValidateToken(token) {
		t.Fatal("ValidateToken should return false for expired token")
	}
}

func TestRevokeToken(t *testing.T) {
	a := NewAuthenticator()

	token := a.CreateToken([]string{"read"}, 1*time.Hour)

	if !a.ValidateToken(token) {
		t.Fatal("Token should be valid before revocation")
	}

	a.RevokeToken(token)

	if a.ValidateToken(token) {
		t.Fatal("Token should be invalid after revocation")
	}
}

func TestHasScope(t *testing.T) {
	a := NewAuthenticator()

	token := a.CreateToken([]string{"read", "write"}, 1*time.Hour)

	if !a.HasScope(token, "read") {
		t.Fatal("Token should have read scope")
	}

	if !a.HasScope(token, "write") {
		t.Fatal("Token should have write scope")
	}

	if a.HasScope(token, "admin") {
		t.Fatal("Token should not have admin scope")
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		header   string
		expected string
	}{
		{"Bearer abc123", "abc123"},
		{"Bearer ", ""},
		{"Basic abc123", ""},
		{"", ""},
		{"abc123", ""},
	}

	for _, test := range tests {
		token := ExtractToken(test.header)
		if token != test.expected {
			t.Fatalf("ExtractToken(%q) = %q, want %q", test.header, token, test.expected)
		}
	}
}
