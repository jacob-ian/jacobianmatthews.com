package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func (ur *UserRepository) FindById(ctx context.Context, id string) (core.User, error) {
	var user core.User
	statement := `
        SELECT * 
        FROM users
        WHERE 
            id = $1 
            AND deleted_at IS NULL;
    `
	err := ur.db.QueryRowContext(ctx, statement, id).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.User{}, core.NewError(core.NotFoundError, "User not found")
		}
		log.Printf("ERROR: DB_USERS_FINDBYID - %v", err.Error())
		return core.User{}, core.NewError(core.InternalError, "An error occurred")
	}
	return user, nil
}

func (us *UserRepository) FindAll(ctx context.Context, filter core.UserFilter) ([]core.User, error) {
	statement := `
        SELECT * 
        FROM users 
        WHERE deleted_at IS NULL
        ORDER BY name ASC;
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
		return []core.User{}, core.NewError(core.InternalError, "Could not get users")
	}
	defer rows.Close()

	var users []core.User
	for rows.Next() {
		var user core.User
		err := rows.Scan(&user)
		if err != nil {
			log.Printf("ERROR: DB_USERS_FINDALL - %v", err.Error())
			return []core.User{}, core.NewError(core.InternalError, "Could not get users")
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Printf("ERROR: DB_USERS_FINDALL - %v", err.Error())
		return []core.User{}, core.NewError(core.InternalError, "Could not get users")
	}
	return users, nil
}

func (us *UserRepository) Create(ctx context.Context, user core.NewUser) (core.User, error) {
	var newUser core.User
	statement := `
        INSERT INTO users(id, name, email, email_verified, image_url) VALUES
            ($1, $2, $3, $4, $5)
        RETURNING *;
    `
	err := us.db.QueryRowContext(ctx, statement, user.Id, user.Name, user.Email, user.EmailVerified, user.ImageUrl).Scan(&newUser)
	if err != nil {
		log.Printf("ERROR: DB_USERS_CREATE - %v", err.Error())
		return core.User{}, core.NewError(core.InternalError, "Could not create user")
	}

	return newUser, nil
}

func (us *UserRepository) Update(ctx context.Context, user core.User) (core.User, error) {
	var updated core.User
	statement := `
        UPDATE users
        SET 
            updated_at = NOW(), 
            name = $2, 
            email = $3, 
            email_verified = $4, 
            image_url = $5
        WHERE 
            id = $1 
            AND deleted_at IS NULL
        RETURNING *;
    `
	err := us.db.QueryRowContext(ctx, statement, user.Id, user.Name, user.Email, user.EmailVerified, user.ImageUrl).Scan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.User{}, core.NewError(core.NotFoundError, "User not found")
		}
		log.Printf("ERROR: DB_USERS_UPDATE - %v", err.Error())
		return core.User{}, core.NewError(core.InternalError, "Could not update user")
	}
	return updated, nil
}

func (us *UserRepository) Delete(ctx context.Context, id string) error {
	statement := `
        UPDATE users
        SET 
            deleted_at = NOW(), 
            updated_at = NOW()
        WHERE 
            id = $1 
            AND deleted_at IS NULL;
    `
	err := us.db.QueryRowContext(ctx, statement, id).Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.NewError(core.NotFoundError, "User not found")
		}
		log.Printf("ERROR: DB_USERS_DELETE - %v", err.Error())
		return core.NewError(core.InternalError, "Could not delete user")
	}
	return nil
}

// Create a postgres implemented UserRepository
func NewUserRepository(ctx context.Context, db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
