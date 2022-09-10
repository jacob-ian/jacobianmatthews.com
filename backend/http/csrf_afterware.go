package http

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type CsrfAfterware struct {
}

func (m *CsrfAfterware) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return backend.NewError(http.StatusInternalServerError, "CSRF Error")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "csrfToken",
		Value:    fmt.Sprintf("%x", token),
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func NewCsrfAfterware() *CsrfAfterware {
	return &CsrfAfterware{}
}
