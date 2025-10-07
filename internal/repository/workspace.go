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
	projects  *mongo.Collection
}

type WorkspaceRepository interface {
	CreateWorkspace(context context.Context, data *models.Workspace) error
	GetAllWorkspaces(context context.Context) ([]*models.Workspace, error)
	GetWorkspaceByID(context context.Context, id string, userID string) (*models.Workspace, error)
	GetWorkspacesByUser(context context.Context, userID string) (*[]models.Workspace, error)
	GetSharedWorkspaces(context context.Context, userID string) (*[]models.Workspace, error)
	UpdateWorkspaceMetadata(context context.Context, id string, name string, description string, userID string) error
	DeleteWorkspace(context context.Context, id string, userID string) error
	AddMemberToWorkspace(context context.Context, workspaceID string, userID string) error
	RemoveMemberFromWorkspace(context context.Context, workspaceID string, userID string) error
}

func NewWorkspaceRepository() WorkspaceRepository {
	return &workspaceRepository{
		workspace: db.GetCollection(config.WORKSPACE),
		projects:  db.GetCollection(config.PROJECT),
	}
}

func (r *workspaceRepository) CreateWorkspace(context context.Context, data *models.Workspace) error {
	data.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := r.workspace.InsertOne(context, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *workspaceRepository) GetAllWorkspaces(context context.Context) ([]*models.Workspace, error) {
	var workspaces []*models.Workspace
	cursor, err := r.workspace.Find(context, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *workspaceRepository) GetWorkspaceByID(context context.Context, id string, userID string) (*models.Workspace, error) {
	var workspace models.Workspace
	ID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.workspace.FindOne(context, bson.M{"_id": ID,
		"$or": bson.A{
			bson.M{"owner_id": userID},
			bson.M{"members": userID},
		},
	}).Decode(&workspace)
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
	cursor, err := r.workspace.Find(context, bson.M{"owner_id": userID})
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
	cursor, err := r.workspace.Find(context, bson.M{"members": bson.M{"$in": []string{userID}}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &workspaces); err != nil {
		return nil, err
	}
	return &workspaces, nil
}

func (r *workspaceRepository) UpdateWorkspaceMetadata(context context.Context, id string, name string, description string, userID string) error {
	ID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.workspace.UpdateOne(context, bson.M{"_id": ID, "owner_id": userID}, bson.M{
		"$set": bson.M{
			"name":        name,
			"description": description,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *workspaceRepository) DeleteWorkspace(context context.Context, id string, userID string) error {
	ID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.projects.DeleteMany(context, bson.M{"workspace": ID, "owner_id": userID})
	if err != nil {
		return err
	}
	_, err = r.workspace.DeleteOne(context, bson.M{"_id": ID, "owner_id": userID})
	if err != nil {
		return err
	}
	return nil
}

func (r *workspaceRepository) AddMemberToWorkspace(context context.Context, workspaceID string, userID string) error {
	ID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return err
	}
	_, err = r.workspace.UpdateOne(context, bson.M{"_id": ID}, bson.M{
		"$addToSet": bson.M{
			"members": userID,
		},
	})
	if err != nil {
		return err
	}
	_, err = r.projects.UpdateMany(context, bson.M{"workspace": ID}, bson.M{
		"$addToSet": bson.M{
			"members": userID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *workspaceRepository) RemoveMemberFromWorkspace(context context.Context, workspaceID string, userID string) error {
	ID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return err
	}
	_, err = r.workspace.UpdateOne(context, bson.M{"_id": ID}, bson.M{
		"$pull": bson.M{
			"members": userID,
		},
	})
	if err != nil {
		return err
	}
	_, err = r.projects.UpdateMany(context, bson.M{"workspace": ID}, bson.M{
		"$pull": bson.M{
			"members": userID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
