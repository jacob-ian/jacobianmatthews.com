package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/firebaseauth"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/postgres"
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

	firebaseAuthProvider, err := firebaseauth.NewAuthProvider(ctx)
	if err != nil {
		log.Fatalf("Could not create auth provider: %v", err.Error())
	}

	authService := core.NewAuthService(core.CoreAuthServiceConfig{
		UserRepository:     db.UserRepository,
		UserRoleRepository: db.UserRoleRepository,
		RoleRepository:     db.RoleRepository,
	})

	app := http.NewApplication(http.Config{
		Port: getPort(),
		Services: http.Services{
			AuthService: authService,
			SessionService: core.NewSessionService(core.CoreSessionServiceConfig{
				UserRepository: db.UserRepository,
				AuthProvider:   firebaseAuthProvider,
				AuthService:    authService,
			}),
		},
		Repositories: http.Repositories{
			UserRepository:     db.UserRepository,
			UserRoleRepository: db.UserRoleRepository,
			RoleRepository:     db.RoleRepository,
		},
	})

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
