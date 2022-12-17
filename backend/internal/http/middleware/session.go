package middleware

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type SessionMiddlewareConfig struct {
	SessionService core.SessionService
}

type SessionMiddleware struct {
	router  http.Handler
	service core.SessionService
	res     *res.ResponseWriterFactory
}

func (m *SessionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		m.router.ServeHTTP(w, r)
		return
	}

	user, err := m.service.VerifySession(r.Context(), cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})
		m.res.NewResponseWriter(w, r).HandleError(err)
		return
	}

	ctx := core.WithUserContext(r.Context(), &user)
	m.router.ServeHTTP(w, r.WithContext(ctx))
}

func (m *SessionMiddleware) Inject(handler http.Handler, writer *res.ResponseWriterFactory) http.Handler {
	m.router = handler
	m.res = writer
	return m
}

// Creates the session authentication middleware
func NewSessionMiddleware(config SessionMiddlewareConfig) *SessionMiddleware {
	return &SessionMiddleware{
		service: config.SessionService,
	}
}
