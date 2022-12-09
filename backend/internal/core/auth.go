package core

import (
	"context"
	"time"
)

type Session struct {
	Cookie    string
	ExpiresIn time.Duration
}

type SessionUser struct {
	User User `json:"user"`
	Role Role `json:"role"`
}

type SessionService interface {
	CreateSession(ctx context.Context, idToken string) (Session, error)
	VerifySession(ctx context.Context, sessionCookie string) (*SessionUser, error)
	RevokeSession(ctx context.Context, uid string) error
}

type AuthService struct {
	users     UserRepository
	roles     RoleRepository
	userRoles UserRoleRepository
}

type AuthServiceConfig struct {
	UserRepository     UserRepository
	RoleRepository     RoleRepository
	UserRoleRepository UserRoleRepository
}

// Gets the user's role
func (auth *AuthService) GetUserRole(ctx context.Context, userId string) (Role, error) {
	return auth.userRoles.FindRoleByUserId(ctx, userId)
}

// Gives a user a role by name
func (auth *AuthService) GiveUserRoleByName(ctx context.Context, userId string, roleName string) (Role, error) {
	role, err := auth.roles.FindByName(ctx, roleName)
	if err != nil {
		return Role{}, err
	}

	_, err = auth.userRoles.Create(ctx, userId, role.Id)
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

// Gets the user's role. If they don't have one, it gives the user the role provided by name.
func (auth *AuthService) GetOrGiveUserRole(ctx context.Context, userId string, roleName string) (Role, error) {
	role, err := auth.GetUserRole(ctx, userId)
	if err == nil {
		return role, nil
	}
	if !IsError(err, NotFoundError) {
		return Role{}, err
	}
	role, err = auth.GiveUserRoleByName(ctx, userId, roleName)
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

func NewAuthService(config AuthServiceConfig) *AuthService {
	return &AuthService{
		users:     config.UserRepository,
		roles:     config.RoleRepository,
		userRoles: config.UserRoleRepository,
	}
}
