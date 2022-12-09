package http

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type AuthMiddleware struct {
	handler http.Handler
	service *core.SessionService
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

	ctx := core.WithUserContext(r.Context(), user)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}

// Creates the authentication middleware
func NewAuthMiddleware(h http.Handler, s *core.SessionService) *AuthMiddleware {
	return &AuthMiddleware{
		handler: h,
		service: s,
	}
}
