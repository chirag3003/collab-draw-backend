package resolvers

//go:generate go run github.com/99designs/gqlgen generate

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/chirag3003/collab-draw-backend/graph/model"
	"github.com/chirag3003/collab-draw-backend/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ProjectSubscriber struct {
	sockedID string
	channel  chan *model.ProjectSubscription
}

type Resolver struct {
	Repo               *repository.Repository
	projectSubscribers map[string][]ProjectSubscriber
	subscribersMutex   sync.RWMutex
}

func NewResolver(repo *repository.Repository) *Resolver {
	return &Resolver{
		Repo:               repo,
		projectSubscribers: make(map[string][]ProjectSubscriber),
	}
}

func generateRandom8DigitString() string {
	// The range for an 8-digit number is [10000000, 99999999].
	// We generate a number in the range [0, 89999999] and add 10000000 to it.
	minN := 10000000
	maxNum := 90000000 // maxNum is exclusive in rand.Intn, so 99999999 is maxNum - 1
	randomNumber := rand.IntN(maxNum) + minN

	// Convert the integer to a string
	return fmt.Sprintf("%d", randomNumber)
}

// Subscribe adds a subscriber for a specific project
func (r *Resolver) subscribeToProject(projectID string, ch chan *model.ProjectSubscription) string {
	r.subscribersMutex.Lock()
	defer r.subscribersMutex.Unlock()
	subscriber := ProjectSubscriber{
		channel:  ch,
		sockedID: generateRandom8DigitString(),
	}
	r.projectSubscribers[projectID] = append(r.projectSubscribers[projectID], subscriber)
	return subscriber.sockedID
}

// Unsubscribe removes a subscriber for a specific project
func (r *Resolver) unsubscribeFromProject(projectID string, socketID string) {
	r.subscribersMutex.Lock()
	defer r.subscribersMutex.Unlock()

	subscribers := r.projectSubscribers[projectID]
	for i, subscriber := range subscribers {
		if subscriber.sockedID == socketID {
			r.projectSubscribers[projectID] = append(subscribers[:i], subscribers[i+1:]...)
			close(subscriber.channel)
			break
		}
	}

	// Clean up empty subscriber lists
	if len(r.projectSubscribers[projectID]) == 0 {
		delete(r.projectSubscribers, projectID)
	}
}

// Broadcast sends a project update to all subscribers
func (r *Resolver) broadcastProjectUpdate(projectID string, project *model.ProjectSubscription, fromID string) {
	r.subscribersMutex.RLock()
	defer r.subscribersMutex.RUnlock()

	if subscribers, ok := r.projectSubscribers[projectID]; ok {
		for _, subscriber := range subscribers {
			if subscriber.sockedID == fromID {
				continue
			}
			project.SocketID = subscriber.sockedID
			select {
			case subscriber.channel <- project:
			default:
				// Channel is full or closed, skip
			}
		}
	}
}
