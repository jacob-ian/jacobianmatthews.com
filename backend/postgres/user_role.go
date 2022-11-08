package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type UserRoleRepository struct {
	db *sql.DB
}

func (urr *UserRoleRepository) FindRoleByUserId(ctx context.Context, userId string) (backend.Role, error) {
	var role backend.Role
	query := `
        SELECT role.* FROM role 
        INNER JOIN user_role
        ON role.id = user_role.role_id
        WHERE user_role.user_id = $1 AND user_role.deleted_at IS NULL
    `
	err := urr.db.QueryRowContext(ctx, query, userId).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return backend.Role{}, backend.NewError(backend.NotFoundError, "Role not found")
		}
		log.Printf("ERROR: DB_USERROLE_FINDROLEBYUSERID - %v", err.Error())
		return backend.Role{}, backend.NewError(backend.InternalError, "Could not find role")
	}
	return role, nil
}

func (urr *UserRoleRepository) FindById(ctx context.Context, id uuid.UUID) (backend.UserRole, error) {
	var userRole backend.UserRole
	query := `
        SELECT * FROM user_role
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := urr.db.QueryRowContext(ctx, query, id).Scan(&userRole)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return backend.UserRole{}, backend.NewError(backend.NotFoundError, "Role not found")
		}
		log.Printf("ERROR: DB_USERROLE_FINDBYID - %v", err.Error())
		return backend.UserRole{}, backend.NewError(backend.InternalError, "Could not find role")
	}
	return userRole, nil
}

func (urr *UserRoleRepository) Create(ctx context.Context, userId string, roleId uuid.UUID) (backend.UserRole, error) {
	var userRole backend.UserRole
	query := `
        INSERT INTO user_role (id, user_id, role_id)
        VALUES (gen_random_uuid(), $1, $2)
        RETURNING *
    `
	err := urr.db.QueryRowContext(ctx, query, userId, roleId).Scan(&userRole)
	if err != nil {
		log.Printf("ERROR: DB_USERROLE_CREATE - %v", err.Error())
		return backend.UserRole{}, backend.NewError(backend.InternalError, "Could not create user role")
	}
	return userRole, nil
}

func (urr *UserRoleRepository) UpdateByUserId(ctx context.Context, userId string, roleId string) (backend.UserRole, error) {
	var updated backend.UserRole
	query := `
        UPDATE user_role
        SET role_id = $2, updated_at = NOW()
        WHERE user_id = $1 AND deleted_at IS NULL
    `
	err := urr.db.QueryRowContext(ctx, query, userId, roleId).Scan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return backend.UserRole{}, backend.NewError(backend.NotFoundError, "User role not found")
		}
		log.Printf("ERROR: DB_USERROLE_UPDATEBYUSERID - %v", err.Error())
		return backend.UserRole{}, backend.NewError(backend.InternalError, "Could not update user role")
	}
	return updated, nil
}

func (urr *UserRoleRepository) DeleteByUserId(ctx context.Context, userId string) error {
	query := `
        UPDATE user_role
        SET deleted_at = NOW(), updated_at = NOW()
        WHERE user_id = $1 AND deleted_at IS NULL
    `
	err := urr.db.QueryRowContext(ctx, query, userId).Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return backend.NewError(backend.NotFoundError, "User role not found")
		}
		log.Printf("ERROR: DB_USERROLE_DELETEBYUSERID - %v", err.Error())
		return backend.NewError(backend.InternalError, "Could not delete user role")
	}
	return nil
}

func NewUserRoleRepository(ctx context.Context, db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{
		db: db,
	}
}
