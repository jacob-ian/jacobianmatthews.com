package backend

import (
	"context"
	"errors"
)

var InvalidToken = errors.New("Invalid ID Token")
var SignInFail = errors.New("Could not sign in")
var InvalidSession = errors.New("Invalid Session")

type SessionCookie string

type SessionUser struct {
	User  *User `json:"user"`
	Admin bool  `json:"admin"`
}

type AuthService interface {
	CreateSession(ctx context.Context, idToken string) (SessionCookie, error)
	VerifySession(ctx context.Context, session SessionCookie) (SessionUser, error)
	RevokeSession(ctx context.Context, session SessionCookie) error
}
