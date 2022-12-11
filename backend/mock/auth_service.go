package mock

import (
	"context"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type MockAuthServiceValues struct {
	GetOrGiveUserRole  MockResponse
	GetUserRole        MockResponse
	GiveUserRoleByName MockResponse
}

type MockAuthService struct {
	values MockAuthServiceValues
}

// GetOrGiveUserRole implements core.AuthService
func (as *MockAuthService) GetOrGiveUserRole(ctx context.Context, userId string, roleName string) (core.Role, error) {
	return (as.values.GetOrGiveUserRole.Value).(core.Role), as.values.GetOrGiveUserRole.Error
}

// GetUserRole implements core.AuthService
func (as *MockAuthService) GetUserRole(ctx context.Context, userId string) (core.Role, error) {
	return (as.values.GetUserRole.Value).(core.Role), as.values.GetUserRole.Error
}

// GiveUserRoleByName implements core.AuthService
func (as *MockAuthService) GiveUserRoleByName(ctx context.Context, userId string, roleName string) (core.Role, error) {
	return (as.values.GiveUserRoleByName.Value).(core.Role), as.values.GiveUserRoleByName.Error
}

var _ core.AuthService = (*MockAuthService)(nil)

func NewAuthService(values MockAuthServiceValues) *MockAuthService {
	return &MockAuthService{
		values: values,
	}
}
