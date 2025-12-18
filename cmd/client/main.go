package main

import (
	"context"
	"log"
	"net/http"
	todov1 "rahulxf.com/rpc-learning/internal/gen/go/todo/v1"
	"rahulxf.com/rpc-learning/internal/gen/go/todo/v1/todov1connect"
)

func main() {
	client := todov1connect.NewTodoServiceClient(http.DefaultClient, "http://localhost:8080")

	createRes, err := client.CreateTodo(context.Background(), &todov1.CreateTodoRequest{
		Title:       "Learn ConnectRPC",
		Description: "Master gRPC and Protocol Buffers",
	})
	if err != nil {
		log.Fatalf("CreateTodo failed: %v", err)
	}
	log.Printf("Created Todo: %s (ID: %s)\n", createRes.Todo.Title, createRes.Todo.Id)

	getRes, err := client.GetTodo(context.Background(), &todov1.GetTodoRequest{
		Id: createRes.Todo.Id,
	})
	if err != nil {
		log.Fatalf("GetTodo failed: %v", err)
	}
	log.Printf("Fetched Todo: %s - Completed: %v\n", getRes.Todo.Title, getRes.Todo.Completed)

	updateRes, err := client.UpdateTodo(context.Background(), &todov1.UpdateTodoRequest{
		Id:        createRes.Todo.Id,
		Completed: true,
	})
	if err != nil {
		log.Fatalf("UpdateTodo failed: %v", err)
	}
	log.Printf("Updated Todo: %s - Completed: %v\n", updateRes.Todo.Title, updateRes.Todo.Completed)

	deleteRes, err := client.DeleteTodo(context.Background(), &todov1.DeleteTodoRequest{
		Id: createRes.Todo.Id,
	})
	if err != nil {
		log.Fatalf("DeleteTodo failed: %v", err)
	}
	log.Printf("Deleted Todo: Success=%v\n", deleteRes.Success)

}
