package service

import (
	"context"

	"connectrpc.com/connect"
	"rahulxf.com/rpc-learning/internal/auth"
	"rahulxf.com/rpc-learning/internal/models"
	"rahulxf.com/rpc-learning/internal/repository"

	todov1 "rahulxf.com/rpc-learning/internal/gen/go/todo/v1"
	"rahulxf.com/rpc-learning/internal/gen/go/todo/v1/todov1connect"
)

type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) todov1connect.TodoServiceHandler {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(
	ctx context.Context,
	req *connect.Request[todov1.CreateTodoRequest],
) (*connect.Response[todov1.CreateTodoResponse], error) {

	userID := ctx.Value(auth.UserIDKey).(string)

	todo, err := s.repo.Create(
		ctx,
		req.Msg.WorkspaceId,
		req.Msg.Title,
		req.Msg.Description,
		userID,
		req.Msg.AssignedTo,
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.CreateTodoResponse{
		Todo: mapTodo(todo),
	}), nil
}

func (s *TodoService) GetTodo(
	ctx context.Context,
	req *connect.Request[todov1.GetTodoRequest],
) (*connect.Response[todov1.GetTodoResponse], error) {

	todo, err := s.repo.GetByID(ctx, req.Msg.Id, req.Msg.WorkspaceId)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&todov1.GetTodoResponse{
		Todo: mapTodo(todo),
	}), nil
}

func (s *TodoService) UpdateTodo(
	ctx context.Context,
	req *connect.Request[todov1.UpdateTodoRequest],
) (*connect.Response[todov1.UpdateTodoResponse], error) {

	todo, err := s.repo.Update(
		ctx,
		req.Msg.Id,
		req.Msg.WorkspaceId,
		req.Msg.Completed,
		req.Msg.Title,
		req.Msg.Description,
		req.Msg.AssignedTo,
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.UpdateTodoResponse{
		Todo: mapTodo(todo),
	}), nil
}

func (s *TodoService) DeleteTodo(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {

	if err := s.repo.Delete(ctx, req.Msg.Id, req.Msg.WorkspaceId); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&todov1.DeleteTodoResponse{Success: true}), nil
}

func (s *TodoService) ListTodo(
	ctx context.Context,
	req *connect.Request[todov1.ListTodoRequest],
) (*connect.Response[todov1.ListTodoResponse], error) {
	return connect.NewResponse(&todov1.ListTodoResponse{}), nil
}

func mapTodo(t *models.TodoDocument) *todov1.Todo {
	var assignedTo string
	if t.AssignedTo != nil {
		assignedTo = t.AssignedTo.Hex()
	}

	return &todov1.Todo{
		Id:          t.ID.Hex(),
		WorkspaceId: t.WorkspaceID.Hex(),
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedBy:   t.CreatedBy.Hex(),
		AssignedTo:  assignedTo,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
