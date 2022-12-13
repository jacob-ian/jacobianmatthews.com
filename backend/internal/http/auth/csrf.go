package auth

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type csrfResponse struct {
	Message string `json:"message"`
}

// Ensures that a CSRF Token is present in the browser
func CSRFTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			res.NewResponseWriter(w, r).WriteError("Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// We don't need to do anything, the CSRF middleware will set the cookie
		res.NewResponseWriter(w, r).Write(http.StatusOK, csrfResponse{
			Message: "OK",
		})
	}
}