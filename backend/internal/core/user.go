package core

import (
	"context"
	"time"
)

type User struct {
	Id            string    `json:"id" field:"id"`
	Name          string    `json:"name" field:"name"`
	Email         string    `json:"email" field:"email"`
	EmailVerified bool      `json:"emailVerified" field:"email_verified"`
	ImageUrl      string    `json:"imageUrl" field:"image_url"`
	CreatedAt     time.Time `json:"createdAt" field:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" field:"updated_at"`
	DeletedAt     time.Time `json:"deletedAt" field:"deleted_at"`
}

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	// Finds a user by ID. Throws a NotFound error if cannot find the User by ID
	FindById(ctx context.Context, id string) (User, error)
	Create(ctx context.Context, user NewUser) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
}

type NewUser struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	ImageUrl      string `json:"imageUrl"`
}
