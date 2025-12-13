package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TodoDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Completed   bool          `bson:"completed"`
	CreatedAt   int64         `bson:"created_at"`
	UpdatedAt   int64         `bson:"updated_at"`
}
