package firebaseauth

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type SessionServiceConfig struct {
	AuthService    *core.AuthService
	UserRepository core.UserRepository
}

type SessionService struct {
	firebase *auth.Client
	auth     *core.AuthService
	users    core.UserRepository
}

// Creates a session from a Firebase Auth ID Token
func (ss *SessionService) CreateSession(ctx context.Context, idToken string) (core.Session, error) {
	decodedToken, err := ss.firebase.VerifyIDToken(ctx, idToken)
	if err != nil {
		return core.Session{}, core.NewError(core.BadRequestError, "Invalid ID Token")
	}

	if !authenticatedWithinTime(decodedToken, time.Minute*5) {
		return core.Session{}, core.NewError(core.UnauthenticatedError, "Unauthenticated")
	}

	if err = ss.registerIfNewUser(ctx, decodedToken.UID); err != nil {
		return core.Session{}, err
	}

	sessionExpiresIn := time.Minute * 15
	cookie, err := ss.firebase.SessionCookie(ctx, idToken, sessionExpiresIn)
	if err != nil {
		return core.Session{}, core.NewError(core.InternalError, "An error occurred whilst signing in")
	}

	return core.Session{
		Cookie:    cookie,
		ExpiresIn: sessionExpiresIn,
	}, nil
}

// Check that ID token was authenticated within a time duration
func authenticatedWithinTime(token *auth.Token, duration time.Duration) bool {
	now := time.Now().Unix()
	timeOfAuth := token.Claims["auth_time"].(int64)
	return now-timeOfAuth < int64(duration)
}

// Checks if given user exists and creates one if they don't
func (ss *SessionService) registerIfNewUser(ctx context.Context, userId string) error {
	_, err := ss.users.FindById(ctx, userId)
	if err == nil {
		return nil
	}

	if !core.IsError(err, core.NotFoundError) {
		return err
	}

	firebaseUser, err := ss.firebase.GetUser(ctx, userId)
	if err != nil {
		return core.NewError(core.InternalError, "Invalid Firebase User ID")
	}

	if _, err := ss.users.Create(ctx, core.NewUser{
		Id:            firebaseUser.UID,
		Name:          firebaseUser.DisplayName,
		Email:         firebaseUser.Email,
		EmailVerified: firebaseUser.EmailVerified,
		ImageUrl:      firebaseUser.PhotoURL,
	}); err != nil {
		return err
	}

	if _, err = ss.auth.GiveUserRoleByName(ctx, userId, "User"); err != nil {
		return err
	}

	return nil
}

// Verifies a session
func (ss *SessionService) VerifySession(ctx context.Context, cookie string) (*core.SessionUser, error) {
	decodedToken, err := ss.firebase.VerifySessionCookie(ctx, cookie)
	if err != nil {
		return nil, core.NewError(core.UnauthenticatedError, "Invalid session")
	}

	userId := decodedToken.UID

	user, err := ss.users.FindById(ctx, userId)
	if err != nil {
		return nil, core.NewError(core.InternalError, "An error occurred")
	}

	role, err := ss.auth.GetOrGiveUserRole(ctx, user.Id, "User")
	if err != nil {
		log.Printf("ERROR: FIREBASEAUTH-VERIFY - %v", err.Error())
		return nil, core.NewError(core.InternalError, "Could not verify session")
	}

	return &core.SessionUser{
		Role: role,
		User: user,
	}, nil
}

func (ss *SessionService) RevokeSession(ctx context.Context, uid string) error {
	err := ss.firebase.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return core.NewError(core.InternalError, "Failed to sign out everywhere")
	}
	return nil
}

// Creates a Firebase Auth implementation of SessionService
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
		firebase: authClient,
		users:    config.UserRepository,
		auth:     config.AuthService,
	}, nil
}
