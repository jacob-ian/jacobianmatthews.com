package mock

import (
	"context"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type MockSessionServiceValues struct {
	StartSession  MockResponse
	VerifySession MockResponse
}

type MockSessionService struct {
	values MockSessionServiceValues
}

// StartSession implements core.SessionService
func (ss *MockSessionService) StartSession(ctx context.Context, idToken string) (core.Session, error) {
	return (ss.values.StartSession.Value).(core.Session), ss.values.StartSession.Error
}

// VerifySession implements core.SessionService
func (ss *MockSessionService) VerifySession(ctx context.Context, sessionCookie string) (core.SessionUser, error) {
	return (ss.values.VerifySession.Value).(core.SessionUser), ss.values.VerifySession.Error
}

var _ core.SessionService = (*MockSessionService)(nil)

func NewSessionService(values MockSessionServiceValues) *MockSessionService {
	return &MockSessionService{
		values: values,
	}
}
