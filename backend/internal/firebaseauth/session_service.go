package firebaseauth

import (
	"context"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type SessionServiceConfig struct {
	UserRepository core.UserRepository
}

type SessionService struct {
	client *auth.Client
	users  core.UserRepository
}

// Creates a session from a Firebase Auth ID Token
func (ss *SessionService) CreateSession(ctx context.Context, idToken string) (core.Session, error) {
	decodedToken, err := ss.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return core.Session{}, core.NewError(core.BadRequestError, "Invalid ID Token")
	}

	if !authenticatedWithin(decodedToken, time.Minute*5) {
		return core.Session{}, core.NewError(core.UnauthenticatedError, "Unauthenticated")
	}

	_, err = ss.findOrCreateUser(ctx, decodedToken)
	if err != nil {
		return core.Session{}, err
	}

	sessionExpiresIn := time.Minute * 15
	cookie, err := ss.client.SessionCookie(ctx, idToken, sessionExpiresIn)
	if err != nil {
		return core.Session{}, core.NewError(core.InternalError, "An error occurred whilst signing in")
	}
	return core.Session{
		Cookie:    cookie,
		ExpiresIn: sessionExpiresIn,
	}, nil
}

// Check that ID token was authenticated within a time duration
func authenticatedWithin(token *auth.Token, duration time.Duration) bool {
	now := time.Now().Unix()
	timeOfAuth := token.Claims["auth_time"].(int64)
	return now-timeOfAuth < int64(duration)
}

// Finds the associated user or creates one
func (ss *SessionService) findOrCreateUser(ctx context.Context, token *auth.Token) (core.User, error) {
	createUser := false
	user, err := ss.users.FindById(ctx, token.UID)
	if err != nil {
		if e, ok := err.(*core.Error); ok {
			if e.IsError(core.NotFoundError) {
				createUser = true
			} else {
				return core.User{}, e
			}
		}
		return core.User{}, err
	}

	if !createUser {
		return user, nil
	}

	fUser, err := ss.client.GetUser(ctx, token.UID)
	if err != nil {
		return core.User{}, core.NewError(core.InternalError, "Invalid User ID")
	}

	user, err = ss.users.Create(ctx, core.NewUser{
		Id:            fUser.UID,
		Name:          fUser.DisplayName,
		Email:         fUser.Email,
		EmailVerified: fUser.EmailVerified,
		ImageUrl:      fUser.PhotoURL,
	})
	if err != nil {
		return core.User{}, err
	}

	return user, nil
}

// Verifies a session
func (ss *SessionService) VerifySession(ctx context.Context, cookie string) (*core.SessionUser, error) {
	decodedToken, err := ss.client.VerifySessionCookieAndCheckRevoked(ctx, cookie)
	if err != nil {
		return nil, core.NewError(core.UnauthenticatedError, "Invalid session")
	}

	userId := decodedToken.UID

	isAdmin := false
	if decodedToken.Claims["admin"] == true {
		isAdmin = true
	}

	user, err := ss.users.FindById(ctx, userId)
	if err != nil {
		return nil, core.NewError(core.InternalError, "An error occurred")
	}

	return &core.SessionUser{
		Admin: isAdmin,
		User:  user,
	}, nil
}

func (auth *SessionService) RevokeSession(ctx context.Context, uid string) error {
	err := auth.client.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return core.NewError(core.InternalError, "Failed to sign out everywhere")
	}
	return nil
}

// Creates a Firebase Auth implementation of AuthService
func NewSessionService(ctx context.Context, config SessionServiceConfig) (*SessionService, error) {
	firebaseApp, err := firebase.NewApp(ctx, &firebase.Config{})
	if err != nil {
		return nil, core.NewError(core.InternalError, "Could not create Firebase App")
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		return nil, core.NewError(core.InternalError, "Could not create Firebase Auth Client")
	}

	return &SessionService{
		client: authClient,
		users:  config.UserRepository,
	}, nil
}
