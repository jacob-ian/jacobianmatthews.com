package middleware

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

// The request context key for retreiving the request's csrf token
const CSRFContextKey string = "csrfToken"

// Add the CSRF token to context
func WithCSRFContext(ctx context.Context, csrfToken string) context.Context {
	return context.WithValue(ctx, CSRFContextKey, csrfToken)
}

// Gets the request CSRF Token from context
func CSRFTokenFromContext(ctx context.Context) (string, bool) {
	csrfToken, ok := ctx.Value(CSRFContextKey).(string)
	return csrfToken, ok
}

type CsrfMiddleware struct {
	handler http.Handler
	writer  *res.ResponseWriterFactory
}

func (m *CsrfMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := r.Context()

	cookie, err := r.Cookie("csrfToken")
	if err == nil {
		context = WithCSRFContext(context, cookie.Value)
	}

	if !errors.Is(err, http.ErrNoCookie) {
		m.writer.NewResponseWriter(w, r).WriteError("Bad CSRF", http.StatusBadRequest)
		return
	}

	token := make([]byte, 16)
	if _, err = rand.Read(token); err != nil {
		m.writer.NewResponseWriter(w, r).WriteError("CSRF Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "csrfToken",
		Value:    fmt.Sprintf("%x", token),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	m.handler.ServeHTTP(w, r.WithContext(context))
}

func (m *CsrfMiddleware) Inject(h http.Handler, w *res.ResponseWriterFactory) http.Handler {
	m.handler = h
	m.writer = w
	return m
}

// Creates middleware that sets a CSRF token cookie to each response
func NewCsrfMiddleware() *CsrfMiddleware {
	return &CsrfMiddleware{}
}
