package resolvers

import "github.com/chirag3003/collab-draw-backend/internal/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo *repository.Repository
}
