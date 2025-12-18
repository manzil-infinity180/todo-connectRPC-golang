package server

import (
	"context"
	"fmt"
	todov1 "rahulxf.com/rpc-learning/internal/gen/go/todo/v1"

	"rahulxf.com/rpc-learning/internal/repository"
)

type TodoServer struct {
	repo *repository.TodoRepository
}

func NewTodoServer(repo *repository.TodoRepository) *TodoServer {
	return &TodoServer{repo: repo}
}

func (s *TodoServer) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	todo, err := s.repo.Create(ctx, req.Title, req.Description)
	if err != nil {
		return nil, err
	}

	return &todov1.CreateTodoResponse{
		Todo: &todov1.Todo{
			Id:          todo.ID.Hex(),
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}

func (s *TodoServer) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("id is required")
	}

	todo, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &todov1.GetTodoResponse{
		Todo: &todov1.Todo{
			Id:          todo.ID.Hex(),
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}

func (s *TodoServer) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("id is required")
	}

	err := s.repo.Update(ctx, req.Id, req.Completed)
	if err != nil {
		return nil, err
	}

	// Fetch updated todo
	todo, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &todov1.UpdateTodoResponse{
		Todo: &todov1.Todo{
			Id:          todo.ID.Hex(),
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}

func (s *TodoServer) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("id is required")
	}

	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &todov1.DeleteTodoResponse{
		Success: true,
	}, nil
}

func (s *TodoServer) ListTodo(context.Context, *todov1.ListTodoRequest) (*todov1.ListTodoResponse, error) {
	res := &todov1.ListTodoResponse{
		Todos:         []*todov1.Todo{},
		NextPageToken: "",
	}

	return res, nil
}
