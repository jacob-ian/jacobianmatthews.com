package http

import (
	"context"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/postgres"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	Port         uint16
	Host         string
	Database     *postgres.Database
	AuthProvider core.AuthProvider
}

type Application struct {
	router         *http.ServeMux
	server         *http.Server
	database       *postgres.Database
	sessionService *core.SessionService
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
	db := config.Database
	authService := core.NewAuthService(core.CoreAuthServiceConfig{
		UserRepository:     db.UserRepository,
		UserRoleRepository: db.UserRoleRepository,
		RoleRepository:     db.RoleRepository,
	})
	sessionService := core.NewSessionService(core.SessionServiceConfig{
		AuthService:    authService,
		AuthProvider:   config.AuthProvider,
		UserRepository: db.UserRepository,
	})

	mux := http.NewServeMux()
	handler := NewGlobalMiddleware(mux, GlobalMiddlewareConfig{
		CorsOrigin: "localhost:3001",
		Accept:     "application/json, application/grpc-web",
	})
	withAuth := NewAuthMiddleware(handler, sessionService)

	srv := http.Server{
		Addr:    config.Host + ":" + strconv.FormatUint(uint64(config.Port), 10),
		Handler: withAuth,
	}

	app := &Application{
		database:       config.Database,
		sessionService: sessionService,
		server:         &srv,
		router:         mux,
	}
	app.connectControllers(ctx)

	return app, nil
}
