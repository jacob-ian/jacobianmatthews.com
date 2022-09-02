package backend

import (
	"context"
)

type SessionCookie string

type SessionUser struct {
	User  User `json:"user"`
	Admin bool `json:"admin"`
}

type AuthService interface {
	CreateSession(ctx context.Context, idToken string) (SessionCookie, error)
	VerifySession(ctx context.Context, session SessionCookie) (SessionUser, error)
	RevokeSession(ctx context.Context, session SessionCookie) error
}
