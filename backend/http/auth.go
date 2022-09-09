package http

import (
	"context"
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
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
	a.router.Handle(route+"/login", handleLogin(a.authService))
	a.router.Handle(route+"/logout", handleLogout(a.authService))
	a.router.Handle(route+"/me", handleMe(a.authService))
}

// Attempt to sign the user into the website
func handleLogin(auth backend.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		payload := &logInPayload{}
		err := NewJsonDecoder(r.Body).Decode(payload)
		if err != nil {
			NewResponseWriter(w, r).HandleError(err)
			return
		}

		csrfCookie, err := r.Cookie("csrfToken")
		if err != nil || csrfCookie.Value != payload.CsrfToken {
			NewResponseWriter(w, r).WriteError("Invalid CSRF", http.StatusBadRequest)
			return
		}

		session, err := auth.CreateSession(r.Context(), payload.IdToken)
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

		NewResponseWriter(w, r).Write(http.StatusOK, logInResponse{
			Message:   "Signed in",
			ExpiresIn: int(session.ExpiresIn.Seconds()),
		})
	}
}

// Revokes the user's session (signs the user out)
func handleLogout(auth backend.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cookie, err := r.Cookie("session")
		if err != nil {
			NewResponseWriter(w, r).WriteError("Not signed in", http.StatusUnauthorized)
			return
		}
		err = auth.RevokeSession(r.Context(), cookie.Value)
		if err != nil {
			NewResponseWriter(w, r).HandleError(err)
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
func handleMe(auth backend.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, ok := UserFromContext(r.Context())
		if !ok {
			NewResponseWriter(w, r).WriteError("Not signed in", http.StatusUnauthorized)
			return
		}

		NewResponseWriter(w, r).Write(http.StatusOK, user)
	}
}
