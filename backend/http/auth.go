package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type success struct {
	Message string `json:"message"`
}

// Connects the authentication controllers
func (a *Application) connectAuthControllers(ctx context.Context, route string) {
	a.router.Handle(route+"/login", handleLogin(a.authService))
}

// Attempt to sign the user into the website
func handleLogin(auth backend.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Actual logic
		json.NewEncoder(w).Encode(success{Message: "Success"})
	}
}
