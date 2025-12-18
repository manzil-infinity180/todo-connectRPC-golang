package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB(uri, dbname string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	return &MongoDB{
		Client: client,
		DB:     client.Database(dbname),
	}, nil
}

func (m *MongoDB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) TodoCollection() *mongo.Collection {
	return m.DB.Collection("todos")
}

func (m *MongoDB) UserCollection() *mongo.Collection {
	return m.DB.Collection("users")
}

func (m *MongoDB) WorkspaceCollection() *mongo.Collection {
	return m.DB.Collection("workspaces")
}

func (m *MongoDB) WorkspaceMemberCollection() *mongo.Collection {
	return m.DB.Collection("workspace_members")
}
