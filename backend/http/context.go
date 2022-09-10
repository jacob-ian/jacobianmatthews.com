package http

import (
	"context"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

var UserContextKey string = "user"

// Adds the session user to the request context
func WithUserContext(ctx context.Context, user *backend.SessionUser) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// Gets the session user from the request context
func UserFromContext(ctx context.Context) (*backend.SessionUser, bool) {
	user, ok := ctx.Value(UserContextKey).(*backend.SessionUser)
	return user, ok
}
