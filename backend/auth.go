package backend

import (
	"context"
	"time"
)

type Session struct {
	Cookie    string
	ExpiresIn time.Duration
}

type SessionUser struct {
	User  User `json:"user"`
	Admin bool `json:"admin"`
}

type AuthService interface {
	CreateSession(ctx context.Context, idToken string) (Session, error)
	VerifySession(ctx context.Context, sessionCookie string) (SessionUser, error)
	RevokeSession(ctx context.Context, sessionCookie string) error
}
