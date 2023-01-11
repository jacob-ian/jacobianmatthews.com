package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type AuthControllerServices struct {
	SessionService core.SessionService
}

type AuthControllersConfig struct {
	Router    *http.ServeMux
	Res       *res.ResponseWriterFactory
	BaseRoute string
	Services  AuthControllerServices
}

// Connects the authentication controllers
func ConnectControllers(config AuthControllersConfig) {
	res, route, router, services :=
		config.Res,
		config.BaseRoute,
		config.Router,
		config.Services

	sessionService := services.SessionService
	router.Handle(route+"/login", NewLoginHandler(res, sessionService))
	router.Handle(route+"/logout", NewLogoutHandler(res, sessionService))
	router.Handle(route+"/me", NewUserInfoHandler(res))
	router.Handle(route+"/csrf", NewCSRFTokenHandler(res))
}
