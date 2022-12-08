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
	err := ur.db.QueryRowContext(ctx, statement, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.EmailVerified,
		&user.ImageUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.User{}, core.NewError(core.NotFoundError, "User not found")
		}
		log.Printf("ERROR: DB_USERS_FINDBYID - %v", err.Error())
		return core.User{}, core.NewError(core.InternalError, "Could not get user")
	}
	return user, nil
}

func (us *UserRepository) FindAll(ctx context.Context) ([]core.User, error) {
	statement := `
        SELECT * 
        FROM users 
        WHERE deleted_at IS NULL
        ORDER BY name ASC;
    `
	rows, err := us.db.QueryContext(ctx, statement)
	if err != nil {
		log.Printf("ERROR: DB_USERS_FINDALL - %v", err.Error())
		return []core.User{}, core.NewError(core.InternalError, "Could not get users")
	}
	defer rows.Close()

	var users []core.User
	for rows.Next() {
		var user core.User
		if err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.EmailVerified,
			&user.ImageUrl,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		); err != nil {
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
	err := us.db.QueryRowContext(ctx, statement, user.Id, user.Name, user.Email, user.EmailVerified, user.ImageUrl).Scan(
		&newUser.Id,
		&newUser.Name,
		&newUser.Email,
		&newUser.EmailVerified,
		&newUser.ImageUrl,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
		&newUser.DeletedAt,
	)
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
	if err := us.db.QueryRowContext(
		ctx,
		statement,
		user.Id,
		user.Name,
		user.Email,
		user.EmailVerified,
		user.ImageUrl,
	).Scan(
		&updated.Id,
		&updated.Name,
		&updated.Email,
		&updated.EmailVerified,
		&updated.ImageUrl,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	); err != nil {
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
	if err := us.db.QueryRowContext(ctx, statement, id).Err(); err != nil {
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
