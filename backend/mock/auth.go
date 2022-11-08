package mock

import (
	"context"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type MockVerifySessionOutput struct {
	Value *core.SessionUser
	Error error
}

type MockCreateSessionOutput struct {
	Value core.Session
	Error error
}

type AuthServiceOutput struct {
	VerifySession MockVerifySessionOutput
	CreateSession MockCreateSessionOutput
	RevokeSession error
}

type AuthService struct {
	output AuthServiceOutput
}

func (a *AuthService) CreateSession(ctx context.Context, idToken string) (core.Session, error) {
	return a.output.CreateSession.Value, a.output.CreateSession.Error
}
func (a *AuthService) VerifySession(ctx context.Context, sessionCookie string) (*core.SessionUser, error) {
	return a.output.VerifySession.Value, a.output.VerifySession.Error
}
func (a *AuthService) RevokeSession(ctx context.Context, sessionCookie string) error {
	return a.output.RevokeSession
}

func NewAuthService(o AuthServiceOutput) *AuthService {
	return &AuthService{output: o}
}
