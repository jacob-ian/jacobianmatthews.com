package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

// Return the details for the currently signed in user
func UserInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			res.NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, ok := core.UserFromContext(r.Context())
		if !ok {
			res.NewResponseWriter(w, r).WriteError("Not signed in", http.StatusUnauthorized)
			return
		}
		res.NewResponseWriter(w, r).Write(http.StatusOK, user)
	}
}
