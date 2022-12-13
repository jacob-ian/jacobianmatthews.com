package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type AuthControllersConfig struct {
	Router         *http.ServeMux
	BaseRoute      string
	SessionService core.SessionService
}

// Connects the authentication controllers
func ConnectControllers(config AuthControllersConfig) {
	config.Router.Handle(config.BaseRoute+"/login", LoginHandler(config.SessionService))
	config.Router.Handle(config.BaseRoute+"/logout", LogoutHandler(config.SessionService))
	config.Router.Handle(config.BaseRoute+"/me", UserInfoHandler())
	config.Router.Handle(config.BaseRoute+"/csrf", CSRFTokenHandler())
}
