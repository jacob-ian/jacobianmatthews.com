package mock

import (
	"context"
	"time"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type MockAuthProviderValues struct {
	CreateSessionCookie MockResponse
	GetUserDetails      MockResponse
	VerifyIdToken       MockResponse
	VerifySessionCookie MockResponse
}

type MockAuthProvider struct {
	CreateSessionCookieRes MockResponse
	GetUserDetailsRes      MockResponse
	VerifyIdTokenRes       MockResponse
	VerifySessionCookieRes MockResponse
}

// CreateSessionCookie implements core.AuthProvider
func (ap *MockAuthProvider) CreateSessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	return (ap.CreateSessionCookieRes.Value).(string), ap.CreateSessionCookieRes.Error
}

// GetUserDetails implements core.AuthProvider
func (ap *MockAuthProvider) GetUserDetails(ctx context.Context, userId string) (core.AuthProviderUser, error) {
	return (ap.GetUserDetailsRes.Value).(core.AuthProviderUser), ap.GetUserDetailsRes.Error
}

// VerifyIdToken implements core.AuthProvider
func (ap *MockAuthProvider) VerifyIdToken(ctx context.Context, token string) (*core.Token, error) {
	return (ap.VerifyIdTokenRes.Value).(*core.Token), ap.VerifyIdTokenRes.Error
}

// VerifySessionCookie implements core.AuthProvider
func (ap *MockAuthProvider) VerifySessionCookie(ctx context.Context, sessionCookie string) (*core.Token, error) {
	return (ap.VerifySessionCookieRes.Value).(*core.Token), ap.VerifySessionCookieRes.Error
}

var _ core.AuthProvider = (*MockAuthProvider)(nil)

func NewAuthProvider(config MockAuthProviderValues) *MockAuthProvider {
	return &MockAuthProvider{
		CreateSessionCookieRes: config.CreateSessionCookie,
		GetUserDetailsRes:      config.GetUserDetails,
		VerifyIdTokenRes:       config.VerifyIdToken,
		VerifySessionCookieRes: config.VerifySessionCookie,
	}
}
