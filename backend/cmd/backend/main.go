package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/jacob-ian/jacobianmatthews.com/backend/firebaseauth"
	"github.com/jacob-ian/jacobianmatthews.com/backend/postgres"
)

var Port int

func main() {
	ctx := context.Background()

	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	database, err := postgres.NewDatabaseClient(ctx, dbConnStr)
	if err != nil {
		log.Fatalf("Could not create database client: %v", err.Error())
	}

	defer database.Close()

	authService, err := newFirebaseAuthService(ctx, database)
	if err != nil {
		log.Fatalf("Could not create auth service: %v", err.Error())
	}
}

func newFirebaseAuthService(ctx context.Context, database *postgres.Database) (*firebaseauth.AuthService, error) {
	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Could not connect to Firebase Admin: %v", err.Error())
	}

	authService, err := firebaseauth.NewAuthService(ctx, firebaseApp, database)
	if err != nil {
		log.Fatalf("Could not create Auth Service: %v", err.Error())
	}
	return authService, nil
}

func getPort() int {
	portEnv, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatalln("Missing PORT environment variable")
	}
	port, err := strconv.ParseUint(portEnv, 10, 16)
	if err != nil {
		log.Fatalln("Invalid PORT environment variable")
	}
	return int(port)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%v %v %v %v", r.Method, r.URL.Path, w.Header().Values("*"), r.UserAgent())
	})
}
