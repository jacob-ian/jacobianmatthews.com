package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
	_ "github.com/lib/pq"
)

// A PostgreSQL database implementation
type Database struct {
	db          *sql.DB
	UserService *UserService
}

// Close the Database connection
func (db *Database) Close() error {
	return db.db.Close()
}

// Create a new PostgreSQL database client
func NewDatabaseClient(ctx context.Context, connStr string) (*Database, error) {
	if connStr == "" {
		return nil, backend.NewError(backend.InternalError, "Missing database connection string")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, backend.NewError(backend.InternalError, fmt.Sprintf("Could not connect to PostgreSQL: %v", err.Error()))
	}

	log.Println("Connected to PostgreSQL")

	return &Database{
		db:          db,
		UserService: NewUserService(ctx, db),
	}, nil
}
