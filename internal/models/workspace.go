package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Workspace struct {
	ID          bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name        string          `bson:"name" json:"name"`
	Description string          `bson:"description" json:"description"`
	Owner       bson.ObjectID   `bson:"owner_id" json:"ownerId"`
	Members     []bson.ObjectID `bson:"members" json:"members"`
	CreatedAt   string          `bson:"created_at" json:"createdAt"`
}
