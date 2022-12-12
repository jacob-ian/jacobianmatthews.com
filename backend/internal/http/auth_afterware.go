package http

import (
	"net/http"
	"time"
)

type AuthAfterware struct {
}

// Update the expiry of the session cookie if it exists
func (a AuthAfterware) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
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

// Creates auth afterware to update session expiry
func NewAuthAfterware() AuthAfterware {
	return AuthAfterware{}
}
