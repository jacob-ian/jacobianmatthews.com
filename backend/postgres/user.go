package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
	_ "github.com/lib/pq"
)

type UserService struct {
	db *sql.DB
}

func (us *UserService) FindById(ctx context.Context, id string) (backend.User, error) {
	var user backend.User
	statement := `
        SELECT * FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := us.db.QueryRowContext(ctx, statement, id).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return backend.User{}, backend.NewError(backend.NotFoundError, "User not found")
		}
		log.Printf("ERROR: DB_USERS_FINDBYID - %v", err.Error())
		return backend.User{}, backend.NewError(backend.InternalError, "An error occurred")
	}
	return user, nil
}

func (us *UserService) FindAll(ctx context.Context, filter backend.UserFilter) ([]backend.User, error) {
	statement := `
        SELECT * FROM users 
        ORDER BY name ASC
    `
	conditions := make([]string, 2)
	if filter.Name != "" {
		conditions = append(conditions, "name = $1")
	}
	if filter.Email != "" {
		conditions = append(conditions, "email = $2")
	}

	clen := len(conditions)
	if clen != 0 {
		statement = statement + "\nWHERE "
	}

	for i := range conditions {
		c := conditions[i]
		statement = statement + c
		if clen > 1 && i != clen-1 {
			statement = statement + ", "
		} else {
			statement = statement + " "
		}
	}
	rows, err := us.db.QueryContext(ctx, statement, filter.Name, filter.Email)
	if err != nil {
		log.Printf("ERROR: DB_USERS_FINDALL - %v", err.Error())
		return []backend.User{}, backend.NewError(backend.InternalError, "Could not get users")
	}
	defer rows.Close()

	var users []backend.User
	for rows.Next() {
		var user backend.User
		err := rows.Scan(&user)
		if err != nil {
			log.Printf("ERROR: DB_USERS_FINDALL - %v", err.Error())
			return []backend.User{}, backend.NewError(backend.InternalError, "Could not get users")
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Printf("ERROR: DB_USERS_FINDALL - %v", err.Error())
		return []backend.User{}, backend.NewError(backend.InternalError, "Could not get users")
	}
	return users, nil
}

func (us *UserService) Create(ctx context.Context, user backend.NewUser) (backend.User, error) {
	var newUser backend.User
	statement := `
        INSERT INTO users (id, name, email, email_verified, image_url)
        VALUES($1, $2, $3, $4, $5)
        RETURNING *;
    `
	err := us.db.QueryRowContext(ctx, statement, user.Id, user.Name, user.Email, user.EmailVerified, user.ImageUrl).Scan(&newUser)
	if err != nil {
		log.Printf("ERROR: DB_USERS_CREATE - %v", err.Error())
		return backend.User{}, backend.NewError(backend.InternalError, "Could not create user")
	}

	return newUser, nil
}

func (us *UserService) Update(ctx context.Context, user backend.User) (backend.User, error) {
	var updated backend.User
	statement := `
        UPDATE users
        SET updated_at = NOW(), name = $2, email = $3, email_verified = $4, image_url = $5
        WHERE id = $1 AND deleted_at IS NULL
        RETURNING *;
    `
	err := us.db.QueryRowContext(ctx, statement, user.Id, user.Name, user.Email, user.EmailVerified, user.ImageUrl).Scan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return backend.User{}, backend.NewError(backend.NotFoundError, "User not found")
		}
		log.Printf("ERROR: DB_USERS_UPDATE - %v", err.Error())
		return backend.User{}, backend.NewError(backend.InternalError, "Could not update user")
	}
	return updated, nil
}

func (us *UserService) Delete(ctx context.Context, id string) error {
	statement := `
        UPDATE users
        SET deleted_at = NOW(), updated_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := us.db.QueryRowContext(ctx, statement, id).Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return backend.NewError(backend.NotFoundError, "User not found")
		}
		log.Printf("ERROR: DB_USERS_DELETE - %v", err.Error())
		return backend.NewError(backend.InternalError, "Could not delete user")
	}
	return nil
}

var _ backend.UserService = (*UserService)(nil)

// Create a UserService
func NewUserService(ctx context.Context, db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}
