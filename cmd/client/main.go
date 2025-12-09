package main

import (
	"context"
	"log"
	"net/http"
	todov1 "rahulxf.com/rpc-learning/gen/todo/v1"
	"rahulxf.com/rpc-learning/gen/todo/v1/todov1connect"
)

func main() {
	client := todov1connect.NewTodoServiceClient(http.DefaultClient, "http://localhost:8080")
	//res, err := client.GetTodo(context.Background(), &todov1.GetTodoRequest{
	//	Id: "5",
	//})
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println(res.Todo.Title)

	res, err := client.CreateTodo(context.Background(), &todov1.CreateTodoRequest{
		Title:       "hello ji rahul here",
		Description: "hello rahuxf your friend here",
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Todo.Title)
}
