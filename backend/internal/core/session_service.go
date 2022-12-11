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

type SessionService struct {
	provider    AuthProvider
	users       UserRepository
	authService AuthService
}

type SessionServiceConfig struct {
	AuthProvider   AuthProvider
	UserRepository UserRepository
	AuthService    AuthService
}

// Starts a session given an IDToken
func (ss *SessionService) StartSession(ctx context.Context, idToken string) (Session, error) {
	decodedToken, err := ss.provider.VerifyIdToken(ctx, idToken)
	if err != nil {
		return Session{}, NewError(BadRequestError, "Invalid ID Token")
	}

	if !authenticatedWithinTime(decodedToken, time.Minute*5) {
		return Session{}, NewError(UnauthenticatedError, "Unauthenticated")
	}

	if err = ss.registerIfNewUser(ctx, decodedToken.Subject); err != nil {
		return Session{}, err
	}

	sessionExpiresIn := time.Minute * 15
	cookie, err := ss.provider.CreateSessionCookie(ctx, idToken, sessionExpiresIn)
	if err != nil {
		return Session{}, NewError(InternalError, "An error occurred whilst signing in")
	}

	return Session{
		Cookie:    cookie,
		ExpiresIn: sessionExpiresIn,
	}, nil
}

// Check that ID token was authenticated within a time duration
func authenticatedWithinTime(token *Token, duration time.Duration) bool {
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

	if !IsError(err, NotFoundError) {
		return err
	}

	idpUser, err := ss.provider.GetUserDetails(ctx, userId)
	if err != nil {
		return NewError(InternalError, "Invalid Firebase User ID")
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
func (ss *SessionService) VerifySession(ctx context.Context, sessionCookie string) (*SessionUser, error) {
	decodedToken, err := ss.provider.VerifySessionCookie(ctx, sessionCookie)
	if err != nil {
		return nil, NewError(UnauthenticatedError, "Invalid session")
	}

	userId := decodedToken.Subject

	user, err := ss.users.FindById(ctx, userId)
	if err != nil {
		return nil, NewError(InternalError, "An error occurred")
	}

	role, err := ss.authService.GetOrGiveUserRole(ctx, user.Id, "User")
	if err != nil {
		log.Printf("ERROR: FIREBASEAUTH-VERIFY - %v", err.Error())
		return nil, NewError(InternalError, "Could not verify session")
	}

	return &SessionUser{
		Role: role,
		User: user,
	}, nil
}

// Create a new Session Service given an Auth Provider (e.g. Firebase Auth)
func NewSessionService(config SessionServiceConfig) *SessionService {
	return &SessionService{
		provider:    config.AuthProvider,
		users:       config.UserRepository,
		authService: config.AuthService,
	}
}
