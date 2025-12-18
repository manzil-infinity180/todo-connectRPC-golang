package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"rahulxf.com/rpc-learning/internal/models"
	"time"
)

type WorkspaceRepository struct {
	workspaceCollection *mongo.Collection
	memberCollection    *mongo.Collection
}

func NewWorkspaceRepository(workspaceCollection, memberCollection *mongo.Collection) *WorkspaceRepository {
	return &WorkspaceRepository{
		workspaceCollection: workspaceCollection,
		memberCollection:    memberCollection,
	}
}

func (r *WorkspaceRepository) Create(ctx context.Context, name, description, ownerID string) (*models.WorkspaceDocument, error) {
	ownerObjectID, err := bson.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner id: %w", err)
	}

	workspace := &models.WorkspaceDocument{
		Name:        name,
		Description: description,
		OwnerID:     ownerObjectID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	result, err := r.workspaceCollection.InsertOne(ctx, workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to insert workspace: %w", err)
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		workspace.ID = oid
	}

	// Add owner as member
	member := &models.WorkspaceMemberDocument{
		WorkspaceID: workspace.ID,
		UserID:      ownerObjectID,
		Role:        "owner",
		JoinedAt:    time.Now().Unix(),
	}

	_, err = r.memberCollection.InsertOne(ctx, member)
	if err != nil {
		return nil, fmt.Errorf("failed to add owner as member: %w", err)
	}

	return workspace, nil
}

func (r *WorkspaceRepository) AddMember(ctx context.Context, workspaceID, userID, role string) (*models.WorkspaceMemberDocument, error) {
	workspaceObjectID, err := bson.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id: %w", err)
	}
	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	var existing models.WorkspaceMemberDocument
	err = r.memberCollection.FindOne(ctx, bson.M{
		"workspace_id": workspaceObjectID,
		"user_id":      userObjectID,
	}).Decode(&existing)

	if err == nil {
		return nil, fmt.Errorf("user is already a member")
	}

	member := &models.WorkspaceMemberDocument{
		WorkspaceID: workspaceObjectID,
		UserID:      userObjectID,
		Role:        role,
		JoinedAt:    time.Now().Unix(),
	}

	result, err := r.memberCollection.InsertOne(ctx, member)
	if err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		member.ID = oid
	}

	return member, nil
}
