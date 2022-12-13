package middleware

import (
	"net/http"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type SessionMiddleware struct {
	handler  http.Handler
	sessions core.SessionService
}

func (m *SessionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		m.handler.ServeHTTP(w, r)
		return
	}

	user, err := m.sessions.VerifySession(r.Context(), cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})
		res.NewResponseWriter(w, r).HandleError(err)
		return
	}

	ctx := core.WithUserContext(r.Context(), &user)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}

// Creates the session authentication middleware
func NewSessionMiddleware(h http.Handler, s core.SessionService) *SessionMiddleware {
	return &SessionMiddleware{
		handler:  h,
		sessions: s,
	}
}
