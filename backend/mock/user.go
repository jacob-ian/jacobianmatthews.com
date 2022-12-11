package mock

import (
	"context"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

type MockUserRepositoryValues struct {
	Create   MockResponse
	Delete   MockResponse
	FindAll  MockResponse
	FindById MockResponse
	Update   MockResponse
}

type MockUserRepository struct {
	values MockUserRepositoryValues
}

// Create implements core.UserRepository
func (ur *MockUserRepository) Create(ctx context.Context, user core.NewUser) (core.User, error) {
	return (ur.values.Create.Value).(core.User), ur.values.Create.Error
}

// Delete implements core.UserRepository
func (ur *MockUserRepository) Delete(ctx context.Context, id string) error {
	return ur.values.Delete.Error
}

// FindAll implements core.UserRepository
func (ur *MockUserRepository) FindAll(ctx context.Context) ([]core.User, error) {
	return (ur.values.FindAll.Value).([]core.User), ur.values.FindAll.Error
}

// FindById implements core.UserRepository
func (ur *MockUserRepository) FindById(ctx context.Context, id string) (core.User, error) {
	return (ur.values.FindById.Value).(core.User), ur.values.FindById.Error
}

// Update implements core.UserRepository
func (ur *MockUserRepository) Update(ctx context.Context, user core.User) (core.User, error) {
	return (ur.values.Update.Value).(core.User), ur.values.Update.Error
}

var _ (core.UserRepository) = (*MockUserRepository)(nil)

func NewUserRepository(values MockUserRepositoryValues) *MockUserRepository {
	return &MockUserRepository{
		values: values,
	}
}
