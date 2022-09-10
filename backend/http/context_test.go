package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend"
	"github.com/jacob-ian/jacobianmatthews.com/backend/http"
)

func TestAttachUserContext(t *testing.T) {
	user := &backend.SessionUser{
		Admin: true,
		User: backend.User{
			Id:        uuid.UUID{},
			Name:      "lolname",
			Email:     "lol@email",
			ImageUrl:  "img",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}

	ctx := context.Background()
	got := http.WithUserContext(ctx, user)

	ctxUser := got.Value(http.UserContextKey)
	if ctxUser != user {
		t.Errorf("Unexpected user in context: got %v want %v", ctxUser, user)
	}
}

func TestUserFromContextExists(t *testing.T) {
	user := &backend.SessionUser{
		Admin: true,
		User: backend.User{
			Id:        uuid.UUID{},
			Name:      "lolname",
			Email:     "lol@email",
			ImageUrl:  "img",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}

	ctx := context.WithValue(context.Background(), http.UserContextKey, user)
	got, gotOk := http.UserFromContext(ctx)

	if !gotOk {
		t.Errorf("Unexpected user in context: got %v want %v", got, user)
	}
}

func TestUserFromContextNotExist(t *testing.T) {
	ctx := context.Background()
	got, gotOk := http.UserFromContext(ctx)
	if gotOk {
		t.Errorf("Unexpected user in context: got %v want %v", got, nil)
	}
}
