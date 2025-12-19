package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"rahulxf.com/rpc-learning/internal/auth"
	"rahulxf.com/rpc-learning/internal/db"
	"rahulxf.com/rpc-learning/internal/repository"
	service "rahulxf.com/rpc-learning/internal/services"

	authv1connect "rahulxf.com/rpc-learning/internal/gen/go/auth/v1/authv1connect"
	todov1connect "rahulxf.com/rpc-learning/internal/gen/go/todo/v1/todov1connect"
)

func main() {
	_ = godotenv.Load()

	mongodb, err := db.NewMongoDB(os.Getenv("MONGODB_URL"), "todo_db")
	if err != nil {
		log.Fatalf("failed to connect mongodb: %v", err)
	}
	defer mongodb.Disconnect()

	jwtManager := auth.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		1*time.Hour,
		7*24*time.Hour,
	)
	userRepo := repository.NewUserRepository(
		mongodb.UserCollection(),
	)
	workspaceRepo := repository.NewWorkspaceRepository(
		mongodb.WorkspaceCollection(),
		mongodb.WorkspaceMemberCollection(),
	)

	todoRepo := repository.NewTodoRepository(
		mongodb.TodoCollection(),
	)

	authService := service.NewAuthService(
		userRepo,
		workspaceRepo,
		jwtManager,
	)
	todoService := service.NewTodoService(
		todoRepo,
	)
	authInterceptor := auth.NewAuthInterceptor(jwtManager)
	protected := connect.WithInterceptors(authInterceptor)

	mux := http.NewServeMux()

	mux.Handle(
		authv1connect.NewAuthServiceHandler(authService),
	)
	mux.Handle(
		todov1connect.NewTodoServiceHandler(
			todoService,
			protected,
		),
	)
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(
		":8080",
		h2c.NewHandler(corsMiddleware(mux), &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization, Connect-Protocol-Version, Connect-Timeout-Ms",
		)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
