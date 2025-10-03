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

type projectRepository struct {
	project *mongo.Collection
}

type ProjectRepository interface {
	NewProject(context context.Context, data *models.Project) error
	UpdateProject(context context.Context, id string, appState string, elements string) error
	GetAll(context context.Context) ([]*models.Project, error)
	GetProjectByID(context context.Context, id string) (*models.Project, error)
	GetProjectsByUserID(context context.Context, userID string) ([]*models.Project, error)
	GetPersonalProjects(context context.Context, userID string) ([]*models.Project, error)
	GetProjectsByWorkspaceID(context context.Context, workspaceID string) ([]*models.Project, error)
	DeleteProject(context context.Context, id string) (bool, error)
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{
		project: db.GetCollection(config.PROJECT),
	}
}

func (r *projectRepository) NewProject(context context.Context, data *models.Project) error {
	data.CreatedAt = time.Now().String()
	_, err := r.project.InsertOne(context, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *projectRepository) UpdateProject(context context.Context, id string, appState string, elements string) error {
	ID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"appState":  appState,
			"elements":  elements,
			"updatedAt": time.Now().Unix(),
		},
	}
	_, err = r.project.UpdateOne(context, bson.M{"_id": ID}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *projectRepository) GetAll(context context.Context) ([]*models.Project, error) {
	var projects []*models.Project
	cursor, err := r.project.Find(context, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetProjectByID(context context.Context, id string) (*models.Project, error) {
	var project models.Project
	ID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.project.FindOne(context, bson.M{"_id": ID}).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetProjectsByUserID(context context.Context, userID string) ([]*models.Project, error) {
	ID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var projects []*models.Project
	cursor, err := r.project.Find(context, bson.M{"owner": ID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetPersonalProjects(context context.Context, userID string) ([]*models.Project, error) {
	ID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var projects []*models.Project
	cursor, err := r.project.Find(context, bson.M{"personal": true, "owner": ID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetProjectsByWorkspaceID(context context.Context, workspaceID string) ([]*models.Project, error) {
	ID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, err
	}
	var projects []*models.Project
	cursor, err := r.project.Find(context, bson.M{"workspace": ID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) DeleteProject(context context.Context, id string) (bool, error) {
	ID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	result, err := r.project.DeleteOne(context, bson.M{"_id": ID})
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 0 {
		return false, nil
	}
	return true, nil
}
