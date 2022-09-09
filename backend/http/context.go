package http

import (
	"context"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

var userKey string = "user"

// Adds the session user to the request context
func WithUserContext(ctx context.Context, user *backend.SessionUser) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// Gets the session user from the request context
func UserFromContext(ctx context.Context) (*backend.SessionUser, bool) {
	user, ok := ctx.Value(userKey).(*backend.SessionUser)
	return user, ok
}
