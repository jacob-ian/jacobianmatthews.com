package http

import (
	"context"
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/json"
)

type logInPayload struct {
	IdToken   string `json:"idToken" required:"true"`
	CsrfToken string `json:"csrfToken" required:"true"`
}

type logInResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expiresIn"`
}

type logOutResponse struct {
	Message string `json:"message"`
}

// Connects the authentication controllers
func (a *Application) connectAuthControllers(ctx context.Context, route string) {
	a.router.Handle(route+"/login", handleLogin(a.sessionService))
	a.router.Handle(route+"/logout", handleLogout(a.sessionService))
	a.router.Handle(route+"/me", handleMe())
}

// Attempt to sign the user into the website
func handleLogin(sessionService *core.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		payload := &logInPayload{}
		err := json.NewJsonDecoder(r.Body).Decode(payload)
		if err != nil {
			NewResponseWriter(w, r).HandleError(err)
			return
		}

		csrfCookie, err := r.Cookie("csrfToken")
		if err != nil || csrfCookie.Value != payload.CsrfToken {
			NewResponseWriter(w, r).WriteError("Invalid CSRF", http.StatusBadRequest)
			return
		}

		session, err := sessionService.StartSession(r.Context(), payload.IdToken)
		if err != nil {
			NewResponseWriter(w, r).HandleError(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    session.Cookie,
			MaxAge:   int(session.ExpiresIn.Seconds()),
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
		})

		NewResponseWriter(w, r).Write(http.StatusCreated, logInResponse{
			Message:   "Signed in",
			ExpiresIn: int(session.ExpiresIn.Seconds()),
		})
	}
}

// Revokes the user's session (signs the user out)
func handleLogout(sessionService *core.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, ok := core.UserFromContext(r.Context())
		if !ok {
			NewResponseWriter(w, r).WriteError("Not signed in", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: 0,
		})

		NewResponseWriter(w, r).Write(http.StatusOK, logOutResponse{
			Message: "Successfully signed out",
		})
	}
}

// Return the details for the currently signed in user
func handleMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, ok := core.UserFromContext(r.Context())
		if !ok {
			NewResponseWriter(w, r).WriteError("Not signed in", http.StatusUnauthorized)
			return
		}
		NewResponseWriter(w, r).Write(http.StatusOK, user)
	}
}
