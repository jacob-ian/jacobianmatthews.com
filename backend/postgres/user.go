package postgres

import (
	"context"
	"database/sql"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
	_ "github.com/lib/pq"
)

type UserService struct {
	db *sql.DB
}

func (us *UserService) FindById(ctx context.Context, id string) (backend.User, error) {
	var user backend.User
	err := us.db.QueryRow("SELECT * from users WHERE id = $1", id).Scan(&user)
	if err != nil {
		return backend.User{}, backend.NewError(backend.InternalError, "An error occurred")
	}
	if user == (backend.User{}) {
		return backend.User{}, backend.NewError(backend.NotFoundError, "User cannot be found")
	}
	return user, nil
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

func (us *UserService) Delete(ctx context.Context, id string) error {
	return backend.NewError(backend.InternalError, "Not implemented")
}

var _ backend.UserService = (*UserService)(nil)

// Create a UserService
func NewUserService(ctx context.Context, db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}
