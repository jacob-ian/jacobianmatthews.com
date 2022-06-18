package backend

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id" firestore:"id"`
	Name      string    `json:"name" firestore:"name"`
	Email     string    `json:"email" firestore:"email"`
	ImageUrl  string    `json:"imageUrl" firestore:"imageUrl"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type UserService interface {
	FindAll(ctx context.Context, filter GetUserFilter) ([]*User, error)
	FindById(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, user NewUser) (*User, error)
	Update(ctx context.Context, user User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"imageUrl"`
}

type GetUserFilter struct {
	Id     *uuid.UUID `json:"id"`
	Name   *string    `json:"name"`
	Email  *string    `json:"email"`
	Limit  *int       `json:"limit"`
	Offset *int       `json:"offset"`
}
