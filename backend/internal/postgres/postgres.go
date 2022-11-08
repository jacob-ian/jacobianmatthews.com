package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	_ "github.com/lib/pq"
)

//go:embed migrations
var embedFs embed.FS

// A PostgreSQL database implementation
type Database struct {
	Db                 *sql.DB
	UserRepository     *UserRepository
	RoleRepository     *RoleRepository
	UserRoleRepository *UserRoleRepository
}

// Close the Database connection
func (db *Database) Close() error {
	return db.Db.Close()
}

func (db *Database) RunMigrations() error {
	driver, err := pgmigrate.WithInstance(db.Db, &pgmigrate.Config{})
	if err != nil {
		return core.NewError(core.InternalError, err.Error())
	}

	source, err := iofs.New(embedFs, "postgres/migrations")
	if err != nil {
		return core.NewError(core.InternalError, err.Error())
	}

	m, err := migrate.NewWithInstance("embed", source, "postgres", driver)
	if err != nil {
		return core.NewError(core.InternalError, err.Error())
	}

	err = m.Up()
	if err != nil {
		return core.NewError(core.InternalError, err.Error())
	}

	return nil
}

// Create a new PostgreSQL database client
func NewDatabaseClient(ctx context.Context, connStr string) (*Database, error) {
	if connStr == "" {
		return nil, core.NewError(core.InternalError, "Missing database connection string")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, core.NewError(core.InternalError, fmt.Sprintf("Could not connect to PostgreSQL: %v", err.Error()))
	}

	log.Println("Connected to PostgreSQL")

	return &Database{
		Db:                 db,
		UserRepository:     NewUserRepository(ctx, db),
		RoleRepository:     NewRoleRepository(ctx, db),
		UserRoleRepository: NewUserRoleRepository(ctx, db),
	}, nil
}
