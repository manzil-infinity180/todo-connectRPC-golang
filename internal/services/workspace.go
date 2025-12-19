package service

import (
	"context"
	"fmt"
	"rahulxf.com/rpc-learning/internal/models"

	"connectrpc.com/connect"
	"rahulxf.com/rpc-learning/internal/auth"
	"rahulxf.com/rpc-learning/internal/repository"

	workspacev1 "rahulxf.com/rpc-learning/internal/gen/go/workspace/v1"
	"rahulxf.com/rpc-learning/internal/gen/go/workspace/v1/workspacev1connect"
)

type WorkspaceService struct {
	repo *repository.WorkspaceRepository
}

func NewWorkspaceService(repo *repository.WorkspaceRepository) workspacev1connect.WorkspaceServiceHandler {
	return &WorkspaceService{repo: repo}
}

func userIDFromCtx(ctx context.Context) (string, error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || id == "" {
		return "", fmt.Errorf("unauthenticated")
	}
	return id, nil
}

func (s *WorkspaceService) CreateWorkspace(
	ctx context.Context,
	req *connect.Request[workspacev1.CreateWorkspaceRequest],
) (*connect.Response[workspacev1.CreateWorkspaceResponse], error) {

	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	ws, err := s.repo.Create(ctx, req.Msg.Name, req.Msg.Description, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&workspacev1.CreateWorkspaceResponse{
		Workspace: mapWorkspace(ws),
	}), nil
}

func (s *WorkspaceService) GetWorkspace(
	ctx context.Context,
	req *connect.Request[workspacev1.GetWorkspaceRequest],
) (*connect.Response[workspacev1.GetWorkspaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *WorkspaceService) ListWorkspaces(
	ctx context.Context,
	req *connect.Request[workspacev1.ListWorkspacesRequest],
) (*connect.Response[workspacev1.ListWorkspacesResponse], error) {

	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	workspaces, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &workspacev1.ListWorkspacesResponse{}
	for _, ws := range workspaces {
		res.Workspaces = append(res.Workspaces, mapWorkspace(ws))
	}

	return connect.NewResponse(res), nil
}

func (s *WorkspaceService) UpdateWorkspace(
	ctx context.Context,
	req *connect.Request[workspacev1.UpdateWorkspaceRequest],
) (*connect.Response[workspacev1.UpdateWorkspaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *WorkspaceService) DeleteWorkspace(
	ctx context.Context,
	req *connect.Request[workspacev1.DeleteWorkspaceRequest],
) (*connect.Response[workspacev1.DeleteWorkspaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *WorkspaceService) AddMember(
	ctx context.Context,
	req *connect.Request[workspacev1.AddMemberRequest],
) (*connect.Response[workspacev1.AddMemberResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *WorkspaceService) RemoveMember(
	ctx context.Context,
	req *connect.Request[workspacev1.RemoveMemberRequest],
) (*connect.Response[workspacev1.RemoveMemberResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *WorkspaceService) UpdateMemberRole(
	ctx context.Context,
	req *connect.Request[workspacev1.UpdateMemberRoleRequest],
) (*connect.Response[workspacev1.UpdateMemberRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *WorkspaceService) ListMembers(
	ctx context.Context,
	req *connect.Request[workspacev1.ListMembersRequest],
) (*connect.Response[workspacev1.ListMembersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func mapWorkspace(ws *models.WorkspaceDocument) *workspacev1.Workspace {
	return &workspacev1.Workspace{
		Id:          ws.ID.Hex(),
		Name:        ws.Name,
		Description: ws.Description,
		OwnerId:     ws.OwnerID.Hex(),
		CreatedAt:   ws.CreatedAt,
		UpdatedAt:   ws.UpdatedAt,
	}
}
