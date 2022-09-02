package firebaseauth

import (
	"context"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend"
	"github.com/jacob-ian/jacobianmatthews.com/backend/postgres"
)

type AuthService struct {
	client   *auth.Client
	database *postgres.Database
}

func (auth *AuthService) CreateSession(ctx context.Context, idToken string) (backend.SessionCookie, error) {
	decodedToken, err := auth.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", backend.NewError(backend.BadRequestError, "Invalid ID Token")
	}
	if time.Now().Unix()-decodedToken.Claims["auth_time"].(int64) > 5*60 {
		return "", backend.NewError(backend.BadRequestError, "Invalid ID Token")
	}

	sessionExpiresIn := time.Minute * 15
	cookie, err := auth.client.SessionCookie(ctx, idToken, sessionExpiresIn)
	if err != nil {
		return "", backend.NewError(backend.InternalError, "An error occurred whilst signing in")
	}
	return backend.SessionCookie(cookie), nil
}

func (auth *AuthService) VerifySession(ctx context.Context, session backend.SessionCookie) (backend.SessionUser, error) {
	decodedToken, err := auth.client.VerifySessionCookieAndCheckRevoked(ctx, string(session))
	if err != nil {
		return backend.SessionUser{}, backend.NewError(backend.UnauthenticatedError, "Invalid session")
	}

	userId, err := uuid.Parse(decodedToken.UID)
	if err != nil {
		return backend.SessionUser{}, backend.NewError(backend.BadRequestError, "Invalid ID Token")
	}

	isAdmin := false
	if decodedToken.Claims["admin"] == true {
		isAdmin = true
	}

	user, err := auth.database.UserService.FindById(ctx, userId)
	if err != nil {
		return backend.SessionUser{}, backend.NewError(backend.InternalError, "An error occurred")
	}

	return backend.SessionUser{
		Admin: isAdmin,
		User:  user,
	}, nil
}

func (auth *AuthService) RevokeSession(ctx context.Context, session backend.SessionCookie) error {
	return backend.NewError(backend.InternalError, "Not implemented")
}

var _ backend.AuthService = (*AuthService)(nil)

// Creates a Firebase Auth implementation of AuthService
func NewAuthService(ctx context.Context, database *postgres.Database) (*AuthService, error) {
	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, backend.NewError(backend.InternalError, "Could not create Firebase App")
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		return nil, backend.NewError(backend.InternalError, "Could not create Firebase Auth Client")
	}

	return &AuthService{
		client:   authClient,
		database: database,
	}, nil
}
