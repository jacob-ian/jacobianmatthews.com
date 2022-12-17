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
)

type Services struct {
	AuthService    core.AuthService
	SessionService core.SessionService
}

type Repositories struct {
	UserRepository     core.UserRepository
	UserRoleRepository core.UserRoleRepository
	RoleRepository     core.RoleRepository
}

type Config struct {
	Port         uint16
	Host         string
	Repositories Repositories
	Services     Services
}

type Middleware interface {
	Inject(handler http.Handler, writer *res.ResponseWriterFactory) http.Handler
}

type Application struct {
	router   *http.ServeMux
	res      *res.ResponseWriterFactory
	server   *http.Server
	db       Repositories
	services Services
}

func (a *Application) Serve() error {
	log.Printf("Listening on localhost%v\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *Application) connectMiddleware(middleware []Middleware) {
	var handler http.Handler
	handler = a.router
	for i := range middleware {
		handler = middleware[i].Inject(handler, a.res)
	}
}

func (a *Application) connectControllers() {
	auth.ConnectControllers(auth.AuthControllersConfig{
		BaseRoute:      "/api/auth",
		Router:         a.router,
		Res:            a.res,
		SessionService: a.services.SessionService,
	})
}

func (a *Application) Shutdown(ctx context.Context) error {
	log.Println("Shutting down...")
	return a.server.Shutdown(ctx)
}

// Creates a new HTTP Applicaton
func NewApplication(config Config) *Application {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    config.Host + ":" + strconv.FormatUint(uint64(config.Port), 10),
		Handler: mux,
	}

	writer := res.NewResponseWriterFactory(res.ResponseWriterConfig{
		Afterware: []res.Afterware{
			middleware.NewCsrfMiddleware(),
			middleware.NewSessionExpiryMiddleware(),
		},
	})

	app := &Application{
		server:   &srv,
		router:   mux,
		res:      writer,
		db:       config.Repositories,
		services: config.Services,
	}

	app.connectMiddleware([]Middleware{
		middleware.NewRequestMiddleware(middleware.RequestMiddlewareConfig{
			CorsOrigin: "localhost:3001",
			Accept:     "application/json, application/grpc-web",
		}),
		middleware.NewSessionMiddleware(middleware.SessionMiddlewareConfig{
			SessionService: app.services.SessionService,
		}),
	})

	app.connectControllers()

	return app
}
