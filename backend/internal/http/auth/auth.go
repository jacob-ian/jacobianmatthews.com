package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type AuthControllersConfig struct {
	Router         *http.ServeMux
	Res            *res.ResponseWriterFactory
	BaseRoute      string
	SessionService core.SessionService
}

// Connects the authentication controllers
func ConnectControllers(config AuthControllersConfig) {
	res, route, router, sessionService :=
		config.Res,
		config.BaseRoute,
		config.Router,
		config.SessionService

	router.Handle(route+"/login", LoginHandler(res, sessionService))
	router.Handle(route+"/logout", LogoutHandler(res, sessionService))
	router.Handle(route+"/me", UserInfoHandler(res))
	router.Handle(route+"/csrf", CSRFTokenHandler(res))
}
