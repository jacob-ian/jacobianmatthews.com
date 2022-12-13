package middleware

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type CsrfMiddleware struct {
}

func (m *CsrfMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return core.NewError(http.StatusInternalServerError, "CSRF Error")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "csrfToken",
		Value:    fmt.Sprintf("%x", token),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})
	return nil
}

// Creates middleware that sets a CSRF token cookie to each response
func NewCsrfMiddleware() *CsrfMiddleware {
	return &CsrfMiddleware{}
}
