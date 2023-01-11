package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

// Return the details for the currently signed in user
func NewUserInfoHandler(W *res.ResponseWriterFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			W.NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, ok := core.UserFromContext(r.Context())
		if !ok {
			W.NewResponseWriter(w, r).WriteError("Not signed in", http.StatusUnauthorized)
			return
		}
		W.NewResponseWriter(w, r).Write(http.StatusOK, user)
	}
}
