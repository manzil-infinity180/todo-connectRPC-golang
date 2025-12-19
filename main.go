package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"rahulxf.com/rpc-learning/internal/auth"
	"rahulxf.com/rpc-learning/internal/db"
	"rahulxf.com/rpc-learning/internal/gen/go/auth/v1/authv1connect"
	service "rahulxf.com/rpc-learning/internal/services"
	"time"

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

	jwtManager := auth.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		1*time.Hour,   // access token
		168*time.Hour, // refresh token (7 days)
	)

	userRepo := repository.NewUserRepository(mongodb.UserCollection())
	workspaceRepo := repository.NewWorkspaceRepository(
		mongodb.WorkspaceCollection(),
		mongodb.WorkspaceMemberCollection(),
	)
	repository.NewTodoRepository(mongodb.TodoCollection())

	authService := service.NewAuthService(userRepo, workspaceRepo, jwtManager)
	//workspaceService := service.NewWorkspaceService(workspaceRepo)
	// todoService := service.NewTodoService(todoRepo, workspaceRepo)

	//authInterceptor := auth.NewAuthInterceptor(jwtManager)
	//interceptors := connect.WithInterceptors(authInterceptor)

	mux := http.NewServeMux()
	mux.Handle(authv1connect.NewAuthServiceHandler(authService))

	handler := corsMiddleware(mux)

	//p := new(http.Protocols)
	//p.SetHTTP1(true)
	//p.SetUnencryptedHTTP2(true)
	//
	//srv := &http.Server{
	//	Addr:      "localhost:8080",
	//	Handler:   mux,
	//	Protocols: p,
	//}
	//
	//log.Println("Server starting on http://localhost:8080")
	//if err := srv.ListenAndServe(); err != nil {
	//	log.Fatalf("Server failed: %v", err)
	//}

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(
		":8080",
		h2c.NewHandler(handler, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
