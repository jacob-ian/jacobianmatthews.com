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
        WHERE user_role.user_id = $1 AND user_role.deleted_at IS NOT NULL
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
        WHERE id = $1 AND deleted_at IS NOT NULL
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
