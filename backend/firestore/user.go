package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ backend.UserService = (*UserService)(nil)

type UserService struct {
	collection *firestore.CollectionRef
}

func NewUserService(client *firestore.Client) *UserService {
	return &UserService{collection: client.Collection("users")}
}

// Finds users given a filter
func (us *UserService) FindAll(ctx context.Context, filter backend.GetUserFilter) ([]*backend.User, error) {
	query := us.collection.Where("deletedAt", "==", nil) // don't get soft deleted records
	if filter.Email != nil {
		query.Where("email", "==", *filter.Email)
	}

	if filter.Id != nil {
		query.Where("id", "==", *filter.Id)
	}

	if filter.Name != nil {
		query.Where("name", "==", *filter.Name)
	}

	if filter.Limit != nil {
		query.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		query.Offset(*filter.Offset)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, backend.InternalError
	}
	var data []*backend.User = []*backend.User{}
	for _, doc := range docs {
		var user *backend.User
		error := doc.DataTo(&user)
		if error != nil {
			return nil, backend.InternalError
		}
		data = append(data, user)
	}
	return data, nil
}

// Finds a user by their ID
func (*UserService) FindById(ctx context.Context, id uuid.UUID) (*backend.User, error) {
	panic("unimplemented")
}

// Creates a user
func (userService *UserService) Create(ctx context.Context, user backend.NewUser) (*backend.User, error) {
	newId, err := uuid.NewRandom()
	if err != nil {
		return nil, backend.InternalError
	}
	_, createErr := userService.collection.Doc(newId.String()).Create(ctx, user)
	if createErr != nil {
		log.Println("Could not create user")
		return nil, backend.InternalError
	}
	data, getErr := userService.collection.Doc(newId.String()).Get(ctx)
	if getErr != nil {
		if status.Code(getErr) == codes.NotFound {
			return nil, backend.InternalError
		}
		return nil, backend.InternalError
	}
	if !data.Exists() {
		return nil, backend.InternalError
	}

	var newUser *backend.User
	dataToErr := data.DataTo(newUser)
	if dataToErr != nil {
		log.Println("DataTo error")
		return nil, backend.InternalError
	}
	return newUser, nil
}

// Updates a user
func (*UserService) Update(ctx context.Context, user backend.User) (*backend.User, error) {
	panic("unimplemented")
}

// (Soft) Deletes a user
func (*UserService) Delete(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}
