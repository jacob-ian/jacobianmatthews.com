package backend

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ImageUrl  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type UserService interface {
	FindAll(ctx context.Context, filter UserFilter) ([]User, error)
	FindById(ctx context.Context, id uuid.UUID) (User, error)
	Create(ctx context.Context, user NewUser) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"imageUrl"`
}

type UserFilter struct {
	Id     *uuid.UUID `json:"id"`
	Name   *string    `json:"name"`
	Email  *string    `json:"email"`
	Limit  *int       `json:"limit"`
	Offset *int       `json:"offset"`
}
