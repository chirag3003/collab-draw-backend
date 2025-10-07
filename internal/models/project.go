package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Project struct {
	ID          bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name        string         `bson:"name" json:"name"`
	Description string         `bson:"description" json:"description"`
	Owner       string         `bson:"owner" json:"owner"`
	Members     []string       `bson:"members" json:"members"`
	Workspace   *bson.ObjectID `bson:"workspace,omitempty" json:"workspace,omitempty"`
	Personal    bool           `bson:"personal" json:"personal"`
	Elements    string         `bson:"elements" json:"elements"`
	CreatedAt   string         `bson:"created_at" json:"createdAt"`
	UpdatedAt   string         `bson:"updated_at" json:"updatedAt"`
}
