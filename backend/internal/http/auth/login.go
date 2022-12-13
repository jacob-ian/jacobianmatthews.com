package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
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

// Attempt to sign the user into the website
func LoginHandler(W *res.ResponseWriterFactory, sessionService core.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			W.NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		payload := logInPayload{}
		err := json.NewJsonDecoder(r.Body).Decode(&payload)
		if err != nil {
			W.NewResponseWriter(w, r).HandleError(err)
			return
		}

		csrfCookie, err := r.Cookie("csrfToken")
		if err != nil || csrfCookie.Value != payload.CsrfToken {
			W.NewResponseWriter(w, r).WriteError("Invalid CSRF", http.StatusBadRequest)
			return
		}

		session, err := sessionService.StartSession(r.Context(), payload.IdToken)
		if err != nil {
			W.NewResponseWriter(w, r).HandleError(err)
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

		W.NewResponseWriter(w, r).Write(http.StatusCreated, logInResponse{
			Message:   "Signed in",
			ExpiresIn: int(session.ExpiresIn.Seconds()),
		})
	}
}
