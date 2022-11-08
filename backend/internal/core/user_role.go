package core

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	Id        uuid.UUID `field:"id"`
	UserId    string    `field:"user_id"`
	RoleId    string    `field:"role_id"`
	CreatedAt time.Time `field:"created_at"`
	UpdatedAt time.Time `field:"updated_at"`
	DeletedAt time.Time `field:"deleted_at"`
}

type UserRoleRepository interface {
	FindRoleByUserId(ctx context.Context, userId string) (Role, error)
	FindById(ctx context.Context, id uuid.UUID) (UserRole, error)
	Create(ctx context.Context, userId string, roleId uuid.UUID) (UserRole, error)
	UpdateByUserId(ctx context.Context, userId string, roleId uuid.UUID) (UserRole, error)
	DeleteByUserId(ctx context.Context, userId string) error
}
