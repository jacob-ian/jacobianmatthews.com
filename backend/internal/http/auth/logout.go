package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type logOutResponse struct {
	Message string `json:"message"`
}

// Revokes the user's session (signs the user out)
func LogoutHandler(sessionService core.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			res.NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, ok := core.UserFromContext(r.Context())
		if !ok {
			res.NewResponseWriter(w, r).WriteError("Not signed in", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: 0,
		})

		res.NewResponseWriter(w, r).Write(http.StatusOK, logOutResponse{
			Message: "Successfully signed out",
		})
	}
}