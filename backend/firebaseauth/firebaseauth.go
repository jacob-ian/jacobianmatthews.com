package firebaseauth

import (
	"context"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type AuthServiceConfig struct {
	UserService backend.UserService
}

type AuthService struct {
	client      *auth.Client
	userService backend.UserService
}

// Creates a session from a Firebase Auth ID Token
func (auth *AuthService) CreateSession(ctx context.Context, idToken string) (backend.Session, error) {
	decodedToken, err := auth.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return backend.Session{}, backend.NewError(backend.BadRequestError, "Invalid ID Token")
	}

	if !authenticatedWithin(decodedToken, time.Minute*5) {
		return backend.Session{}, backend.NewError(backend.UnauthenticatedError, "Unauthenticated")
	}

	_, err = auth.findOrCreateUser(ctx, decodedToken)
	if err != nil {
		return backend.Session{}, err
	}

	sessionExpiresIn := time.Minute * 15
	cookie, err := auth.client.SessionCookie(ctx, idToken, sessionExpiresIn)
	if err != nil {
		return backend.Session{}, backend.NewError(backend.InternalError, "An error occurred whilst signing in")
	}
	return backend.Session{
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
func (auth *AuthService) findOrCreateUser(ctx context.Context, token *auth.Token) (backend.User, error) {
	createUser := false
	user, err := auth.userService.FindById(ctx, token.UID)
	if err != nil {
		if e, ok := err.(*backend.Error); ok {
			if e.IsError(backend.NotFoundError) {
				createUser = true
			} else {
				return backend.User{}, e
			}
		}
		return backend.User{}, err
	}

	if !createUser {
		return user, nil
	}

	fUser, err := auth.client.GetUser(ctx, token.UID)
	if err != nil {
		return backend.User{}, backend.NewError(backend.InternalError, "Invalid User ID")
	}

	user, err = auth.userService.Create(ctx, backend.NewUser{
		Id:            fUser.UID,
		Name:          fUser.DisplayName,
		Email:         fUser.Email,
		EmailVerified: fUser.EmailVerified,
		ImageUrl:      fUser.PhotoURL,
	})
	if err != nil {
		return backend.User{}, err
	}

	return user, nil
}

// Verifies a session
func (auth *AuthService) VerifySession(ctx context.Context, cookie string) (*backend.SessionUser, error) {
	decodedToken, err := auth.client.VerifySessionCookieAndCheckRevoked(ctx, cookie)
	if err != nil {
		return nil, backend.NewError(backend.UnauthenticatedError, "Invalid session")
	}

	userId := decodedToken.UID

	isAdmin := false
	if decodedToken.Claims["admin"] == true {
		isAdmin = true
	}

	user, err := auth.userService.FindById(ctx, userId)
	if err != nil {
		return nil, backend.NewError(backend.InternalError, "An error occurred")
	}

	return &backend.SessionUser{
		Admin: isAdmin,
		User:  user,
	}, nil
}

func (auth *AuthService) RevokeSession(ctx context.Context, uid string) error {
	err := auth.client.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return backend.NewError(backend.InternalError, "Failed to sign out everywhere")
	}
	return nil
}

// Creates a Firebase Auth implementation of AuthService
func NewAuthService(ctx context.Context, config AuthServiceConfig) (*AuthService, error) {
	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, backend.NewError(backend.InternalError, "Could not create Firebase App")
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		return nil, backend.NewError(backend.InternalError, "Could not create Firebase Auth Client")
	}

	return &AuthService{
		client:      authClient,
		userService: config.UserService,
	}, nil
}
