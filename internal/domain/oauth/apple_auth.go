package oauth

import (
	"context"
)

// AppleIDToken represents the information from an Apple ID token
type AppleIDToken struct {
	Subject string // User's unique identifier
	Email   string // User's email
	Name    string // User's name
}

// AppleAuthService defines the interface for Apple authentication
type AppleAuthService interface {
	// AuthCodeWithApple authenticates a user with Apple using an authorization code
	AuthCodeWithApple(ctx context.Context, authorizationCode string) (*AppleIDToken, error)
}
