package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

func TestAttachUserContext(t *testing.T) {
	user := &core.SessionUser{
		Role: core.Role{
			Id:        uuid.Must(uuid.NewRandom()),
			Name:      "Admin",
			CreatedAt: time.Now(),
		},
		User: core.User{
			Id:        "id",
			Name:      "lolname",
			Email:     "lol@email",
			ImageUrl:  "img",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	ctx := context.Background()
	got := core.WithUserContext(ctx, user)

	ctxUser := got.Value(core.UserContextKey)
	if ctxUser != user {
		t.Errorf("Unexpected user in context: got %v want %v", ctxUser, user)
	}
}

func TestUserFromContextExists(t *testing.T) {
	user := &core.SessionUser{
		Role: core.Role{
			Id:        uuid.Must(uuid.NewRandom()),
			Name:      "Admin",
			CreatedAt: time.Now(),
		},
		User: core.User{
			Id:        "id",
			Name:      "lolname",
			Email:     "lol@email",
			ImageUrl:  "img",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	ctx := context.WithValue(context.Background(), core.UserContextKey, user)
	got, gotOk := core.UserFromContext(ctx)

	if !gotOk {
		t.Errorf("Unexpected user in context: got %v want %v", got, user)
	}
}

func TestUserFromContextNotExist(t *testing.T) {
	ctx := context.Background()
	got, gotOk := core.UserFromContext(ctx)
	if gotOk {
		t.Errorf("Unexpected user in context: got %v want %v", got, nil)
	}
}
