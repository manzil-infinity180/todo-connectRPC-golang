package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TodoDocument struct {
	ID          bson.ObjectID  `bson:"_id,omitempty"`
	Title       string         `bson:"title"`
	Description string         `bson:"description"`
	Completed   bool           `bson:"completed"`
	CreatedAt   int64          `bson:"created_at"`
	CreatedBy   bson.ObjectID  `bson:"created_by"`
	UpdatedAt   int64          `bson:"updated_at"`
	WorkspaceID bson.ObjectID  `bson:"workspace_id"`
	AssignedTo  *bson.ObjectID `bson:"assigned_to,omitempty"`
}

type UserDocument struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	Email        string        `bson:"email"`
	PasswordHash string        `bson:"password_hash"`
	Name         string        `bson:"name"`
	CreatedAt    int64         `bson:"created_at"`
	UpdatedAt    int64         `bson:"updated_at"`
}

type WorkspaceDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	OwnerID     bson.ObjectID `bson:"owner_id"`
	CreatedAt   int64         `bson:"created_at"`
	UpdatedAt   int64         `bson:"updated_at"`
}

type WorkspaceMemberDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	WorkspaceID bson.ObjectID `bson:"workspace_id"`
	UserID      bson.ObjectID `bson:"user_id"`
	Role        string        `bson:"role"`
	JoinedAt    int64         `bson:"joined_at"`
}
