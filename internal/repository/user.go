package repository

import (
	"context"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

type userRepository struct {
	clerk *user.Client
}

type UserRepository interface {
	GetUsersByID(ctx context.Context, id []string) (*clerk.UserList, error)
	GetUserByEmail(ctx context.Context, email string) (*clerk.UserList, error)
}

func NewUserRepository() UserRepository {

	key := os.Getenv("CLERK_SECRET_KEY") // Replace with your actual Clerk API key or fetch from config

	return &userRepository{
		clerk: user.NewClient(&clerk.ClientConfig{
			BackendConfig: clerk.BackendConfig{
				Key: &key,
			},
		}),
	}
}

func (r *userRepository) GetUsersByID(ctx context.Context, id []string) (*clerk.UserList, error) {
	// Implement the logic to fetch a user by ID from the database
	list, err := r.clerk.List(ctx, &user.ListParams{
		UserIDs: id,
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*clerk.UserList, error) {
	// Implement the logic to fetch a user by email from the database
	list, err := r.clerk.List(ctx, &user.ListParams{
		EmailAddresses: []string{email},
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
