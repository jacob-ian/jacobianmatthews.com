package firebaseauth

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type FirebaseAuthProvider struct {
	client *auth.Client
}

func (fb *FirebaseAuthProvider) VerifyIdToken(ctx context.Context, token string) (*core.Token, error) {
	decodedToken, err := fb.client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, core.NewError(core.BadRequestError, "Invalid Token")
	}
	return &core.Token{
		Subject: decodedToken.Subject,
		Claims:  decodedToken.Claims,
	}, nil
}

func (fb *FirebaseAuthProvider) VerifySessionCookie(ctx context.Context, token string) (*core.Token, error) {
	decodedToken, err := fb.client.VerifySessionCookie(ctx, token)
	if err != nil {
		return nil, core.NewError(core.BadRequestError, "Invalid Session")
	}
	return &core.Token{
		Subject: decodedToken.Subject,
		Claims:  decodedToken.Claims,
	}, nil
}

func (fb *FirebaseAuthProvider) CreateSessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	cookie, err := fb.client.SessionCookie(ctx, idToken, expiresIn)
	if err != nil {
		return "", core.NewError(core.InternalError, "Could not create session cookie")
	}
	return cookie, nil
}

func (fb *FirebaseAuthProvider) GetUserDetails(ctx context.Context, userId string) (core.AuthProviderUser, error) {
	user, err := fb.client.GetUser(ctx, userId)
	if err != nil {
		return core.AuthProviderUser{}, core.NewError(core.BadRequestError, "Invalid Firebase User ID")
	}
	return core.AuthProviderUser{
		Id:            user.UID,
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		ImageUrl:      user.PhotoURL,
	}, nil
}

// Creates a Firebase Auth implementation of AuthProvider
func NewAuthProvider(ctx context.Context) (*FirebaseAuthProvider, error) {
	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, core.NewError(core.InternalError, fmt.Sprintf("Could not create Firebase App: %v", err.Error()))
	}

	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		return nil, core.NewError(core.InternalError, fmt.Sprintf("Could not create Firebase Auth: %v", err.Error()))
	}

	return &FirebaseAuthProvider{
		client: client,
	}, nil
}
