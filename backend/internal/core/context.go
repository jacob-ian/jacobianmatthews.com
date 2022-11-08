package core

import (
	"context"
)

var UserContextKey string = "user"

// Adds the session user to the request context
func WithUserContext(ctx context.Context, user *SessionUser) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// Gets the session user from the request context
func UserFromContext(ctx context.Context) (*SessionUser, bool) {
	user, ok := ctx.Value(UserContextKey).(*SessionUser)
	return user, ok
}
