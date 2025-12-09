package main

import (
	"context"
	"net/http"
	todov1 "rahulxf.com/rpc-learning/gen/todo/v1"
	todov1connect "rahulxf.com/rpc-learning/gen/todo/v1/todov1connect"
	"time"
)

type TodoServer struct {
}

func (s *TodoServer) CreateTodo(_ context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	res := &todov1.CreateTodoResponse{
		Todo: &todov1.Todo{
			Id:          "generated-id", // could be UUID
			Title:       req.Title,
			Description: req.Description,
			Completed:   false,
			CreatedAt:   time.Now().Unix(),
		},
	}
	return res, nil
}

func (s *TodoServer) GetTodo(_ context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	res := &todov1.GetTodoResponse{
		Todo: &todov1.Todo{
			Id:          "hello",
			Title:       "Hello rahulxf",
			Description: "hello rahulxf is good",
			Completed:   false,
			CreatedAt:   time.Now().Unix(),
		},
	}
	return res, nil
}

func (s *TodoServer) ListTodo(context.Context, *todov1.ListTodoRequest) (*todov1.ListTodoResponse, error) {
	res := &todov1.ListTodoResponse{
		Todos:         []*todov1.Todo{},
		NextPageToken: "",
	}

	return res, nil
}

func main() {
	todo := &TodoServer{}
	mux := http.NewServeMux()

	path, handler := todov1connect.NewTodoServiceHandler(todo)
	mux.Handle(path, handler)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}
	s.ListenAndServe()
}
