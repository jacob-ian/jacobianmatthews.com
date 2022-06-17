package backend

import "context"

type Session struct {
	SessionCookie string
}

type RequestUser struct {
	User  *User
	Admin bool `json:"admin"`
}

type AuthService interface {
	CreateSession(idToken string) (*Session, error)
	RevokeSession(*Session) error
	GetRequestUser() (*RequestUser, error)

	AddAdminByEmail(email string) error
	RemoveAdminByEmail(email string) error
}
