package http

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
	"github.com/jacob-ian/jacobianmatthews.com/backend/postgres"
)

type Config struct {
	Port        uint16
	Host        string
	Database    *postgres.Database
	AuthService backend.AuthService
}

type Application struct {
	router      *http.ServeMux
	server      *http.Server
	database    *postgres.Database
	authService backend.AuthService
}

func (a *Application) Serve() error {
	log.Printf("Listening on localhost%v\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *Application) connectControllers(ctx context.Context) {
	a.connectAuthControllers(ctx, "/api/auth")
}

func (a *Application) Shutdown(ctx context.Context) error {
	log.Println("Shutting down...")
	return a.server.Shutdown(ctx)
}

// Creates a new HTTP Applicaton
func NewApplication(ctx context.Context, config Config) (*Application, error) {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    config.Host + ":" + strconv.FormatUint(uint64(config.Port), 10),
		Handler: mux,
	}
	app := &Application{
		database:    config.Database,
		authService: config.AuthService,
		server:      &srv,
		router:      mux,
	}
	app.connectControllers(ctx)
	return app, nil
}
