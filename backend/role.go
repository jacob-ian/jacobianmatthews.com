package backend

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Id        uuid.UUID `field:"id"`
	Name      string    `field:"name"`
	CreatedAt time.Time `field:"created_at"`
	DeletedAt time.Time `field:"deleted_at"`
}

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (Role, error)
	FindAll(ctx context.Context) ([]Role, error)
	Create(ctx context.Context, name string) (Role, error)
	Delete(ctx context.Context, name string) error
}
