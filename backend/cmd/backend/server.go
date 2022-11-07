package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jacob-ian/jacobianmatthews.com/backend/firebaseauth"
	"github.com/jacob-ian/jacobianmatthews.com/backend/http"
	"github.com/jacob-ian/jacobianmatthews.com/backend/postgres"
)

func main() {
	log.Printf("\n--------\njacobianmatthews.com/api - Copyright © Jacob Ian Matthews\n--------")
	ctx := context.Background()

	db, err := postgres.NewDatabaseClient(ctx, os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Could not create database client: %v", err.Error())
	}
	defer db.Close()

	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		log.Println("Running database migrations...")
		err := db.RunMigrations()
		if err != nil {
			log.Fatalf("Could not run database migrations: %v", err.Error())
		}
	}

	auth, err := firebaseauth.NewAuthService(ctx, firebaseauth.AuthServiceConfig{
		UserService: db.UserService,
	})
	if err != nil {
		log.Fatalf("Could not create auth service: %v", err.Error())
	}

	app, err := http.NewApplication(ctx, http.Config{
		Port:        getPort(),
		Database:    db,
		AuthService: auth,
	})
	if err != nil {
		log.Fatalf("Could not create HTTP Application: %v", err.Error())
	}

	log.Fatal(app.Serve())
	defer app.Shutdown(ctx)
}

func getPort() uint16 {
	portEnv, exists := os.LookupEnv("PORT")
	if !exists {
		return 3001
	}
	port, err := strconv.ParseUint(portEnv, 10, 16)
	if err != nil {
		log.Println("Invalid PORT variable, using default.")
		return 3001
	}
	return uint16(port)
}
