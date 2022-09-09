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

// Connects the authentication controllers
func (a *Application) connectAuthControllers(ctx context.Context, route string) {
	a.router.Handle(route+"/login", handleLogin(a.authService))
}

// Attempt to sign the user into the website
func handleLogin(auth backend.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		payload := &logInPayload{}
		dec := NewJsonDecoder(r.Body)
		err := dec.Decode(payload)
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
			Message:   "Success",
			ExpiresIn: int(session.ExpiresIn.Seconds()),
		})
	}
}
