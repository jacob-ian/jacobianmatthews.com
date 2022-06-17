package firestore

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
)

func Connect() (*firestore.Client, error) {
	projectId, exists := os.LookupEnv("FIREBASE_PROJECT_ID")

	if !exists {
		return nil, errors.New("Missing project ID environment variable")
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		return nil, err
	}
	return client, nil
}
