package repository

import (
	"context"
	"errors"
	"time"

	"github.com/chirag3003/collab-draw-backend/internal/config"
	"github.com/chirag3003/collab-draw-backend/internal/db"
	"github.com/chirag3003/collab-draw-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type workspaceRepository struct {
	workspace *mongo.Collection
}

type WorkspaceRepository interface {
	CreateWorkspace(context context.Context, data *models.Workspace) error
	GetWorkspaceByID(context context.Context, id string) (*models.Workspace, error)
	GetWorkspacesByUser(context context.Context, userID string) (*[]models.Workspace, error)
	GetSharedWorkspaces(context context.Context, userID string) (*[]models.Workspace, error)
}

func NewWorkspaceRepository() WorkspaceRepository {
	return &workspaceRepository{
		workspace: db.GetCollection(config.WORKSPACE),
	}
}

func (r *workspaceRepository) CreateWorkspace(context context.Context, data *models.Workspace) error {
	data.CreatedAt = time.Now().String()
	_, err := r.workspace.InsertOne(context, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *workspaceRepository) GetWorkspaceByID(context context.Context, id string) (*models.Workspace, error) {
	var workspace models.Workspace
	ID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.workspace.FindOne(context, bson.M{"_id": ID}).Decode(&workspace)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) GetWorkspacesByUser(context context.Context, userID string) (*[]models.Workspace, error) {
	var workspaces []models.Workspace
	ID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.workspace.Find(context, bson.M{"members": ID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &workspaces); err != nil {
		return nil, err
	}
	return &workspaces, nil
}

func (r *workspaceRepository) GetSharedWorkspaces(context context.Context, userID string) (*[]models.Workspace, error) {
	var workspaces []models.Workspace
	ID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.workspace.Find(context, bson.M{"members": bson.M{"$in": []bson.ObjectID{ID}}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &workspaces); err != nil {
		return nil, err
	}
	return &workspaces, nil
}
