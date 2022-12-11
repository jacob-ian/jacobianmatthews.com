package core

import (
	"context"
)

type AuthService interface {
	GetUserRole(ctx context.Context, userId string) (Role, error)
	GiveUserRoleByName(ctx context.Context, userId string, roleName string) (Role, error)
	GetOrGiveUserRole(ctx context.Context, userId string, roleName string) (Role, error)
}

type CoreAuthService struct {
	users     UserRepository
	roles     RoleRepository
	userRoles UserRoleRepository
}

type CoreAuthServiceConfig struct {
	UserRepository     UserRepository
	RoleRepository     RoleRepository
	UserRoleRepository UserRoleRepository
}

// Gets the user's role
func (auth *CoreAuthService) GetUserRole(ctx context.Context, userId string) (Role, error) {
	return auth.userRoles.FindRoleByUserId(ctx, userId)
}

// Gives a user a role by name
func (auth *CoreAuthService) GiveUserRoleByName(ctx context.Context, userId string, roleName string) (Role, error) {
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
func (auth *CoreAuthService) GetOrGiveUserRole(ctx context.Context, userId string, roleName string) (Role, error) {
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

func NewAuthService(config CoreAuthServiceConfig) *CoreAuthService {
	return &CoreAuthService{
		users:     config.UserRepository,
		roles:     config.RoleRepository,
		userRoles: config.UserRoleRepository,
	}
}
