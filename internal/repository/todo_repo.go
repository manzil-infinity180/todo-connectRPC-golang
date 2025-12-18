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

func (r *TodoRepository) Create(ctx context.Context, workspaceID, title, description, createdBy string, assignedTo *string) (*models.TodoDocument, error) {
	workspaceObjectID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id: %w", err)
	}
	createdByObjectID, err := bson.ObjectIDFromHex(createdBy)
	if err != nil {
		return nil, fmt.Errorf("invalid created_by id: %w", err)
	}

	todo := &models.TodoDocument{
		WorkspaceID: workspaceObjectID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedBy:   createdByObjectID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if assignedTo != nil && *assignedTo != "" {
		assignedToObjectID, err := bson.ObjectIDFromHex(*assignedTo)
		if err != nil {
			return nil, fmt.Errorf("invalid assigned_to id: %w", err)
		}
		todo.AssignedTo = &assignedToObjectID
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

func (r *TodoRepository) GetByID(ctx context.Context, id, workspaceID string) (*models.TodoDocument, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}
	workspaceObjectID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id: %w", err)
	}
	var todo models.TodoDocument

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "workspace_id": workspaceObjectID}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("todo not found")

		}
		return nil, fmt.Errorf("failed")
	}
	return &todo, nil
}

func (r *TodoRepository) Update(ctx context.Context, id, workspaceID string, completed *bool, title, description, assignedTo *string) (*models.TodoDocument, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	workspaceObjectID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now().Unix(),
		},
	}

	if completed != nil {
		update["$set"].(bson.M)["completed"] = *completed
	}

	if title != nil {
		update["$set"].(bson.M)["title"] = *title
	}

	if description != nil {
		update["$set"].(bson.M)["description"] = *description
	}

	if assignedTo != nil {
		if *assignedTo == "" {
			update["$unset"] = bson.M{"assigned_to": ""}
		} else {
			assignedToObjectID, err := bson.ObjectIDFromHex(*assignedTo)
			if err != nil {
				return nil, fmt.Errorf("invalid assigned_to id: %w", err)
			}
			update["$set"].(bson.M)["assigned_to"] = assignedToObjectID
		}
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":          objectID,
		"workspace_id": workspaceObjectID,
	}, update)

	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("todo not found")
	}

	return r.GetByID(ctx, id, workspaceID)
}

func (r *TodoRepository) Delete(ctx context.Context, id, workspaceID string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	workspaceObjectID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf("invalid workspace id: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{
		"_id":          objectID,
		"workspace_id": workspaceObjectID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}
