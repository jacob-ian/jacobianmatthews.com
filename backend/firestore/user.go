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

// Create implements backend.UserService
func (userService *UserService) Create(user backend.NewUser) (*backend.User, error) {
	newId, err := uuid.NewRandom()
	if err != nil {
		return nil, backend.InternalError
	}
	ctx := context.Background()
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

	var newUser backend.User
	dataToErr := data.DataTo(&newUser)
	if dataToErr != nil {
		log.Println("DataTo error")
		return nil, backend.InternalError
	}
	return &newUser, nil
}

// Delete implements backend.UserService
func (*UserService) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

// GetById implements backend.UserService
func (*UserService) GetById(id uuid.UUID) (*backend.User, error) {
	panic("unimplemented")
}

// Update implements backend.UserService
func (*UserService) Update(user backend.User) (*backend.User, error) {
	panic("unimplemented")
}
