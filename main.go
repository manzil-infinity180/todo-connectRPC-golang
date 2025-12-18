package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"rahulxf.com/rpc-learning/cmd/server"
	"rahulxf.com/rpc-learning/internal/db"
	"rahulxf.com/rpc-learning/internal/gen/go/todo/v1/todov1connect"

	"rahulxf.com/rpc-learning/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongodb_url := os.Getenv("MONGODB_URL")
	mongodb, err := db.NewMongoDB(mongodb_url, "todo_db")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Disconnect()
	log.Println("Connected to MongoDB successfully")

	todoRepo := repository.NewTodoRepository(mongodb.TodoCollection())
	todoServer := server.NewTodoServer(todoRepo)
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(todoServer)
	mux.Handle(path, handler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	srv := &http.Server{
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}

	log.Println("Server starting on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
