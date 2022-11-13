package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	connStr, ok := os.LookupEnv("TEST_DB_CONNECTION_STRING")
	if !ok {
		log.Panicf("Missing TEST_FB_CONNECTION_STRING")
	}
	testDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panicf("Could not connect to testing DB: %v", err.Error())
	}
	defer testDB.Close()

	os.Exit(m.Run())
}
