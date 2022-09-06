package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend"
	_ "github.com/lib/pq"
)

type UserService struct {
	db *sql.DB
}

func (us *UserService) FindById(ctx context.Context, id uuid.UUID) (backend.User, error) {
	return backend.User{}, backend.NewError(backend.InternalError, "Not implemented")
}

func (us *UserService) FindAll(ctx context.Context, filter backend.UserFilter) ([]backend.User, error) {
	return []backend.User{}, backend.NewError(backend.InternalError, "Not implemented")
}

func (us *UserService) Create(ctx context.Context, user backend.NewUser) (backend.User, error) {
	return backend.User{}, backend.NewError(backend.InternalError, "Not implemented")
}

func (us *UserService) Update(ctx context.Context, user backend.User) (backend.User, error) {
	return backend.User{}, backend.NewError(backend.InternalError, "Not implemented")
}

func (us *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return backend.NewError(backend.InternalError, "Not implemented")
}

var _ backend.UserService = (*UserService)(nil)

// Create a UserService
func NewUserService(ctx context.Context, db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}