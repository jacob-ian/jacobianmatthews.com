package http

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/auth"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/middleware"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/postgres"
)

type Config struct {
	Port         uint16
	Host         string
	Database     *postgres.Database
	AuthProvider core.AuthProvider
}

type Application struct {
	router         *http.ServeMux
	res            *res.ResponseWriterFactory
	server         *http.Server
	database       *postgres.Database
	sessionService core.SessionService
}

func (a *Application) Serve() error {
	log.Printf("Listening on localhost%v\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *Application) connectControllers() {
	auth.ConnectControllers(auth.AuthControllersConfig{
		BaseRoute:      "/api/auth",
		Router:         a.router,
		Res:            a.res,
		SessionService: a.sessionService,
	})
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
	sessionService := core.NewSessionService(core.CoreSessionServiceConfig{
		AuthService:    authService,
		AuthProvider:   config.AuthProvider,
		UserRepository: db.UserRepository,
	})

	mux := http.NewServeMux()
	handler := middleware.NewRequestMiddleware(mux, middleware.RequestMiddlewareConfig{
		CorsOrigin: "localhost:3001",
		Accept:     "application/json, application/grpc-web",
	})
	withAuth := middleware.NewSessionMiddleware(handler, sessionService)

	res := res.NewResponseWriterFactory(res.ResponseWriterConfig{
		Afterware: []res.Afterware{
			middleware.NewCsrfMiddleware(),
			middleware.NewSessionExpiryMiddleware(),
		},
	})

	srv := http.Server{
		Addr:    config.Host + ":" + strconv.FormatUint(uint64(config.Port), 10),
		Handler: withAuth,
	}

	app := &Application{
		database:       config.Database,
		sessionService: sessionService,
		server:         &srv,
		router:         mux,
		res:            res,
	}

	app.connectControllers()

	return app, nil
}
