package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type UserRoleRepository struct {
	db *sql.DB
}

func (urr *UserRoleRepository) FindRoleByUserId(ctx context.Context, userId string) (core.Role, error) {
	var role core.Role
	query := `
        SELECT roles.* 
        FROM roles
        INNER JOIN user_role
        ON roles.id = user_role.role_id
        WHERE user_role.user_id = $1 
            AND user_role.deleted_at IS NULL;
    `
	if err := urr.db.QueryRowContext(ctx, query, userId).Scan(
		&role.Id,
		&role.Name,
		&role.CreatedAt,
		&role.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Role{}, core.NewError(core.NotFoundError, "Role not found")
		}
		log.Printf("ERROR: DB_USERROLE_FINDROLEBYUSERID - %v", err.Error())
		return core.Role{}, core.NewError(core.InternalError, "Could not find role")
	}
	return role, nil
}

func (urr *UserRoleRepository) FindById(ctx context.Context, id uuid.UUID) (core.UserRole, error) {
	var userRole core.UserRole
	query := `
        SELECT * 
        FROM user_role
        WHERE id = $1 
            AND deleted_at IS NULL;
    `
	if err := urr.db.QueryRowContext(ctx, query, id).Scan(
		&userRole.Id,
		&userRole.UserId,
		&userRole.RoleId,
		&userRole.CreatedAt,
		&userRole.UpdatedAt,
		&userRole.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.UserRole{}, core.NewError(core.NotFoundError, "User role not found")
		}
		log.Printf("ERROR: DB_USERROLE_FINDBYID - %v", err.Error())
		return core.UserRole{}, core.NewError(core.InternalError, "Could not find user role")
	}
	return userRole, nil
}

func (urr *UserRoleRepository) Create(ctx context.Context, userId string, roleId uuid.UUID) (core.UserRole, error) {
	var userRole core.UserRole
	query := `
        INSERT INTO user_role (id, user_id, role_id) VALUES 
            (gen_random_uuid(), $1, $2)
        RETURNING *;
    `
	if err := urr.db.QueryRowContext(
		ctx,
		query,
		userId,
		roleId,
	).Scan(
		&userRole.Id,
		&userRole.UserId,
		&userRole.RoleId,
		&userRole.CreatedAt,
		&userRole.UpdatedAt,
		&userRole.DeletedAt,
	); err != nil {
		log.Printf("ERROR: DB_USERROLE_CREATE - %v", err.Error())
		return core.UserRole{}, core.NewError(core.InternalError, "Could not create user role")
	}
	return userRole, nil
}

func (urr *UserRoleRepository) UpdateByUserId(ctx context.Context, userId string, roleId uuid.UUID) (core.UserRole, error) {
	var updated core.UserRole
	query := `
        UPDATE user_role
        SET 
            role_id = $2, 
            updated_at = NOW()
        WHERE   
            user_id = $1 
            AND deleted_at IS NULL
        RETURNING *;
    `
	if err := urr.db.QueryRowContext(
		ctx,
		query,
		userId,
		roleId,
	).Scan(
		&updated.Id,
		&updated.UserId,
		&updated.RoleId,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.UserRole{}, core.NewError(core.NotFoundError, "User role not found")
		}
		log.Printf("ERROR: DB_USERROLE_UPDATEBYUSERID - %v", err.Error())
		return core.UserRole{}, core.NewError(core.InternalError, "Could not update user role")
	}
	return updated, nil
}

func (urr *UserRoleRepository) DeleteByUserId(ctx context.Context, userId string) error {
	query := `
        UPDATE user_role
        SET 
            deleted_at = NOW(), 
            updated_at = NOW()
        WHERE   
            user_id = $1 
            AND deleted_at IS NULL
    `
	if err := urr.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.NewError(core.NotFoundError, "User role not found")
		}
		log.Printf("ERROR: DB_USERROLE_DELETEBYUSERID - %v", err.Error())
		return core.NewError(core.InternalError, "Could not delete user role")
	}
	return nil
}

func NewUserRoleRepository(ctx context.Context, db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{
		db: db,
	}
}
