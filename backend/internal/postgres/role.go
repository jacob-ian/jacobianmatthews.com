package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type RoleRepository struct {
	db *sql.DB
}

func (rr *RoleRepository) FindByName(ctx context.Context, name string) (core.Role, error) {
	var role core.Role
	query := `
        SELECT * FROM roles
        WHERE name = $1 AND deleted_at IS NULL
    `
	err := rr.db.QueryRowContext(ctx, query, name).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Role{}, core.NewError(core.NotFoundError, "Role not found")
		}
		log.Printf("ERROR: DB_ROLE_FINDBYNAME - %v", err.Error())
		return core.Role{}, core.NewError(core.InternalError, "Could not get role")
	}
	return role, nil
}

func (rr *RoleRepository) FindAll(ctx context.Context) ([]core.Role, error) {
	query := `SELECT * FROM roles WHERE deleted_at IS NULL`
	rows, err := rr.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: DB_ROLE_FINDALL - %v", err.Error())
		return []core.Role{}, core.NewError(core.InternalError, "Could not get roles")
	}
	var roles []core.Role
	for rows.Next() {
		var role core.Role
		err := rows.Scan(&role)
		if err != nil {
			log.Printf("ERROR: DB_ROLE_FINDALL - %v", err.Error())
			return []core.Role{}, core.NewError(core.InternalError, "Could not get roles")
		}
		roles = append(roles, role)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("ERROR: DB_ROLE_FINDALL - %v", err.Error())
		return []core.Role{}, core.NewError(core.InternalError, "Could not get roles")
	}
	return roles, nil
}

func (rr *RoleRepository) Create(ctx context.Context, name string) (core.Role, error) {
	var role core.Role
	query := `
        INSERT INTO roles (id, name)
        VALUES (gen_random_uuid(), $1)
        RETURNING *;
    `
	err := rr.db.QueryRowContext(ctx, query, name).Scan(&role)
	if err != nil {
		log.Printf("ERROR: DB_ROLE_CREATE - %v", err.Error())
		return core.Role{}, core.NewError(core.InternalError, "Could not create role")
	}
	return role, nil
}

func (rr *RoleRepository) Delete(ctx context.Context, name string) error {
	query := `
        UPDATE roles
        SET deleted_at = NOW(), updated_at = NOW()
        WHERE name = $1 AND deleted_at IS NULL
    `
	err := rr.db.QueryRowContext(ctx, query, name).Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.NewError(core.NotFoundError, "Role not found")
		}
		log.Printf("ERROR: DB_ROLE_DELETE - %v", err.Error())
		return core.NewError(core.InternalError, "Could not delete role")
	}
	return nil
}

// Create a postgres implementation of RoleRepository
func NewRoleRepository(ctx context.Context, db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}
