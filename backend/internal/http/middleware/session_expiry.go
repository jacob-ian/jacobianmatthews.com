package middleware

import (
	"net/http"
	"time"
)

type SessionExpiryMiddleware struct {
}

// Update the expiry of the session cookie if it exists
func (a SessionExpiryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	// If cookie is set to be deleted, don't reset the expiry
	if cookie.MaxAge >= 0 {
		return nil
	}

	expiresIn := time.Minute * 15
	http.SetCookie(w, &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		MaxAge:   int(expiresIn.Seconds()),
		SameSite: cookie.SameSite,
		HttpOnly: cookie.HttpOnly,
		Secure:   cookie.Secure,
	})
	return nil
}

// Creates middleware that extends session expiry by another 15 minutes
func NewSessionExpiryMiddleware() *SessionExpiryMiddleware {
	return &SessionExpiryMiddleware{}
}
