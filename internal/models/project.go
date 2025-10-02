package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Project struct {
	ID          bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name        string         `bson:"name" json:"name"`
	Description string         `bson:"description" json:"description"`
	Owner       bson.ObjectID  `bson:"owner" json:"owner"`
	Workspace   *bson.ObjectID `bson:"workspace,omitempty" json:"workspace,omitempty"`
	Personal    bool           `bson:"personal" json:"personal"`
	AppState    string         `bson:"app_state" json:"appState"`
	Elements    string         `bson:"elements" json:"elements"`
	CreatedAt   string         `bson:"created_at" json:"createdAt"`
	UpdatedAt   int            `bson:"updated_at" json:"updatedAt"`
}
