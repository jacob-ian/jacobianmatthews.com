package core

import (
	"context"
	"log"
	"time"
)

type Session struct {
	Cookie    string
	ExpiresIn time.Duration
}

type SessionUser struct {
	User User `json:"user"`
	Role Role `json:"role"`
}

type SessionService interface {
	// Starts a user session given an ID Token (registers the user if new)
	StartSession(ctx context.Context, idToken string) (Session, error)
	// Verifies a session cookie and returns the session user
	VerifySession(ctx context.Context, sessionCookie string) (SessionUser, error)
}

type CoreSessionService struct {
	provider    AuthProvider
	users       UserRepository
	authService AuthService
}

type CoreSessionServiceConfig struct {
	AuthProvider   AuthProvider
	UserRepository UserRepository
	AuthService    AuthService
}

// Starts a session given an IDToken
func (ss *CoreSessionService) StartSession(ctx context.Context, idToken string) (Session, error) {
	decodedToken, err := ss.provider.VerifyIdToken(ctx, idToken)
	if err != nil {
		return Session{}, NewError(BadRequestError, "Invalid ID Token")
	}

	if err := checkRecentlyAuthenticated(decodedToken); err != nil {
		return Session{}, err
	}

	if err = ss.registerIfNewUser(ctx, decodedToken.Subject); err != nil {
		return Session{}, err
	}

	sessionExpiresIn := time.Hour * 24 * 5
	cookie, err := ss.provider.CreateSessionCookie(ctx, idToken, sessionExpiresIn)
	if err != nil {
		return Session{}, NewError(InternalError, "An error occurred whilst signing in")
	}

	return Session{
		Cookie:    cookie,
		ExpiresIn: sessionExpiresIn,
	}, nil
}

// Check that ID token was authenticated within the last 5 minutes
func checkRecentlyAuthenticated(token *Token) error {
	now := time.Now().Unix()
	timeOfAuth, ok := token.Claims["auth_time"].(int64)
	if !ok {
		log.Printf("Invalid ID Token 'auth_time'")
		return NewError(BadRequestError, "Invalid ID Token")
	}
	if now-timeOfAuth > int64(time.Minute*5) {
		return NewError(UnauthenticatedError, "Unauthenticated")
	}
	return nil
}

// Checks if given user exists and creates one if they don't
func (ss *CoreSessionService) registerIfNewUser(ctx context.Context, userId string) error {
	_, err := ss.users.FindById(ctx, userId)
	if err == nil {
		return nil
	}

	if !IsError(err, NotFoundError) {
		return err
	}

	idpUser, err := ss.provider.GetUserDetails(ctx, userId)
	if err != nil {
		return NewError(BadRequestError, "Invalid Firebase User ID")
	}

	if _, err := ss.users.Create(ctx, NewUser{
		Id:            idpUser.Id,
		Name:          idpUser.DisplayName,
		Email:         idpUser.Email,
		EmailVerified: idpUser.EmailVerified,
		ImageUrl:      idpUser.ImageUrl,
	}); err != nil {
		return err
	}

	if _, err = ss.authService.GiveUserRoleByName(ctx, userId, "User"); err != nil {
		return err
	}

	return nil
}

// Verify a session cookie and get the details of the request user
func (ss *CoreSessionService) VerifySession(ctx context.Context, sessionCookie string) (SessionUser, error) {
	decodedToken, err := ss.provider.VerifySessionCookie(ctx, sessionCookie)
	if err != nil {
		return SessionUser{}, NewError(UnauthenticatedError, "Invalid session")
	}

	userId := decodedToken.Subject

	user, err := ss.users.FindById(ctx, userId)
	if err != nil {
		return SessionUser{}, NewError(InternalError, "Could not get signed in user")
	}

	role, err := ss.authService.GetOrGiveUserRole(ctx, user.Id, "User")
	if err != nil {
		return SessionUser{}, NewError(InternalError, "Could not verify session")
	}

	return SessionUser{
		Role: role,
		User: user,
	}, nil
}

// Create a new Session Service given an Auth Provider (e.g. Firebase Auth)
func NewSessionService(config CoreSessionServiceConfig) *CoreSessionService {
	return &CoreSessionService{
		provider:    config.AuthProvider,
		users:       config.UserRepository,
		authService: config.AuthService,
	}
}
