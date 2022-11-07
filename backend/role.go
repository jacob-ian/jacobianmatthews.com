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

type UserRole struct {
	Id     uuid.UUID `field:"id"`
	UserId string    `field:"user_id"`
	RoleId string    `field:"role_id"`
}

type UserRoleRepository interface {
	FindByUserId(ctx context.Context, userId string) (UserRole, error)
	FindById(ctx context.Context, id uuid.UUID) (UserRole, error)
	Create(ctx context.Context, userId string, roleId uuid.UUID) (UserRole, error)
	Update(ctx context.Context, update UserRole) (UserRole, error)
	DeleteByUserId(ctx context.Context, userId string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
