package backend

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageUrl string    `json:"imageUrl"`
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"imageUrl"`
}

type UserService interface {
	GetById(id uuid.UUID) (*User, error)
	Create(user NewUser) (*User, error)
	Update(user User) (*User, error)
	Delete(id uuid.UUID) error
}
