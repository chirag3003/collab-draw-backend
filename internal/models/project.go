package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Project struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	OwnerID     bson.ObjectID `bson:"owner_id" json:"ownerId"`
	Workspace   bson.ObjectID `bson:"workspace" json:"workspace"`
	CreatedAt   int64         `bson:"created_at" json:"createdAt"`
}
