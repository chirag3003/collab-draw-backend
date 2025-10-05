package resolvers

//go:generate go run github.com/99designs/gqlgen generate

import (
	"sync"

	"github.com/chirag3003/collab-draw-backend/graph/model"
	"github.com/chirag3003/collab-draw-backend/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo               *repository.Repository
	projectSubscribers map[string][]chan *model.Project
	subscribersMutex   sync.RWMutex
}

func NewResolver(repo *repository.Repository) *Resolver {
	return &Resolver{
		Repo:               repo,
		projectSubscribers: make(map[string][]chan *model.Project),
	}
}

// Subscribe adds a subscriber for a specific project
func (r *Resolver) subscribeToProject(projectID string, ch chan *model.Project) {
	r.subscribersMutex.Lock()
	defer r.subscribersMutex.Unlock()
	r.projectSubscribers[projectID] = append(r.projectSubscribers[projectID], ch)
}

// Unsubscribe removes a subscriber for a specific project
func (r *Resolver) unsubscribeFromProject(projectID string, ch chan *model.Project) {
	r.subscribersMutex.Lock()
	defer r.subscribersMutex.Unlock()

	subscribers := r.projectSubscribers[projectID]
	for i, subscriber := range subscribers {
		if subscriber == ch {
			r.projectSubscribers[projectID] = append(subscribers[:i], subscribers[i+1:]...)
			close(ch)
			break
		}
	}

	// Clean up empty subscriber lists
	if len(r.projectSubscribers[projectID]) == 0 {
		delete(r.projectSubscribers, projectID)
	}
}

// Broadcast sends a project update to all subscribers
func (r *Resolver) broadcastProjectUpdate(projectID string, project *model.Project) {
	r.subscribersMutex.RLock()
	defer r.subscribersMutex.RUnlock()

	if subscribers, ok := r.projectSubscribers[projectID]; ok {
		for _, ch := range subscribers {
			select {
			case ch <- project:
			default:
				// Channel is full or closed, skip
			}
		}
	}
}
