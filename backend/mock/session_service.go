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

type SessionServiceOutput struct {
	VerifySession MockVerifySessionOutput
	CreateSession MockCreateSessionOutput
	RevokeSession error
}

type SessionService struct {
	output SessionServiceOutput
}

func (a *SessionService) CreateSession(ctx context.Context, idToken string) (core.Session, error) {
	return a.output.CreateSession.Value, a.output.CreateSession.Error
}
func (a *SessionService) VerifySession(ctx context.Context, sessionCookie string) (*core.SessionUser, error) {
	return a.output.VerifySession.Value, a.output.VerifySession.Error
}
func (a *SessionService) RevokeSession(ctx context.Context, sessionCookie string) error {
	return a.output.RevokeSession
}

func NewSessionService(o SessionServiceOutput) *SessionService {
	return &SessionService{output: o}
}
