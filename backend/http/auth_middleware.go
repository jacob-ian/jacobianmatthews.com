package http

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type AuthMiddleware struct {
	handler http.Handler
	service backend.AuthService
}

func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		m.handler.ServeHTTP(w, r)
		return
	}

	user, err := m.service.VerifySession(r.Context(), cookie.Value)
	if err != nil {
		NewResponseWriter(w, r).HandleError(err)
		return
	}

	ctx := WithUserContext(r.Context(), user)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}

// Creates the authentication middleware
func NewAuthMiddleware(h http.Handler, a backend.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		handler: h,
		service: a,
	}
}
