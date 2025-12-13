package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"rahulxf.com/rpc-learning/internal/models"
	"time"
)

type TodoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(collection *mongo.Collection) *TodoRepository {
	return &TodoRepository{collection: collection}
}

func (r *TodoRepository) Create(ctx context.Context, title, description string) (*models.TodoDocument, error) {
	todo := &models.TodoDocument{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	result, err := r.collection.InsertOne(ctx, todo)

	if err != nil {
		return nil, fmt.Errorf("failed to insert todo: %w", err)
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		todo.ID = oid
	} else {
		return nil, fmt.Errorf("unexpected InsertedID type: %T", result.InsertedID)
	}

	return todo, nil
}

func (r *TodoRepository) GetByID(ctx context.Context, id string) (*models.TodoDocument, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var todo models.TodoDocument

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("todo not found")

		}
		return nil, fmt.Errorf("failed")
	}
	return &todo, nil
}

//func (r *TodoRepository) List(ctx context.Context, pageSize int32, pageToken string) ([]*models.TodoDocument), string {
//
//}

func (r *TodoRepository) Update(ctx context.Context, id string, completed bool) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"completed":  completed,
			"updated_at": time.Now().Unix(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("todo not found")
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}
