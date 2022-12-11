package core

import (
	"context"
	"time"
)

// A decoded auth token
type Token struct {
	Subject string
	Claims  map[string]any
}

// The user details provided by the Auth Provider
type AuthProviderUser struct {
	Id            string
	DisplayName   string
	Email         string
	EmailVerified bool
	ImageUrl      string
}

// An authentication provider
type AuthProvider interface {
	// Verifies an ID token
	VerifyIdToken(ctx context.Context, token string) (*Token, error)
	// Verifies a session cookie
	VerifySessionCookie(ctx context.Context, sessionCookie string) (*Token, error)
	// Creates a session cookie (token)
	CreateSessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error)
	// Gets the user's details from the Identity Provider (Google, Apple)
	GetUserDetails(ctx context.Context, userId string) (AuthProviderUser, error)
}
