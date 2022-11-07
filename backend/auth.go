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

type SessionService interface {
	CreateSession(ctx context.Context, idToken string) (Session, error)
	VerifySession(ctx context.Context, sessionCookie string) (*SessionUser, error)
	RevokeSession(ctx context.Context, uid string) error
}

type AuthService interface {
	GetUserRole(ctx context.Context, userId string) (Role, error)
	GrantUserRole(ctx context.Context, userId string, roleName string) (UserRole, error)
}
