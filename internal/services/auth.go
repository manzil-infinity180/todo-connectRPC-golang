package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"rahulxf.com/rpc-learning/internal/auth"
	authv1 "rahulxf.com/rpc-learning/internal/gen/go/auth/v1"
	"rahulxf.com/rpc-learning/internal/gen/go/auth/v1/authv1connect"
	"rahulxf.com/rpc-learning/internal/repository"
)

type AuthService struct {
	userRepo      *repository.UserRepository
	workspaceRepo *repository.WorkspaceRepository
	jwtManager    *auth.JWTManager
}

func NewAuthService(
	userRepo *repository.UserRepository,
	workspaceRepo *repository.WorkspaceRepository,
	jwtManager *auth.JWTManager,
) authv1connect.AuthServiceHandler {
	return &AuthService{
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
		jwtManager:    jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *connect.Request[authv1.RegisterRequest]) (*connect.Response[authv1.RegisterResponse], error) {
	user, err := s.userRepo.Create(ctx, req.Msg.Email, req.Msg.Password, req.Msg.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeAlreadyExists, err)
	}

	// Create default workspace for user
	_, err = s.workspaceRepo.Create(ctx, fmt.Sprintf("%s's Workspace", user.Name), "Default workspace", user.ID.Hex())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create default workspace: %w", err))
	}

	return connect.NewResponse(&authv1.RegisterResponse{
		UserId:  user.ID.Hex(),
		Message: "User registered successfully",
	}), nil
}

func (s *AuthService) Login(ctx context.Context, req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.LoginResponse], error) {
	user, err := s.userRepo.ComparePassword(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	// Get user's workspaces
	workspaces, err := s.workspaceRepo.ListByUserID(ctx, user.ID.Hex())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	workspaceIDs := make([]string, len(workspaces))
	for i, ws := range workspaces {
		workspaceIDs[i] = ws.ID.Hex()
	}

	// Generate tokens
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(user.ID.Hex(), workspaceIDs)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&authv1.LoginResponse{
		UserId:      user.ID.Hex(),
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	})

	// Set cookies
	resp.Header().Set("Set-Cookie", fmt.Sprintf(
		"access_token=%s; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=%d",
		accessToken,
		3600, // 1 hour
	))

	resp.Header().Add("Set-Cookie", fmt.Sprintf(
		"Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=%d",
		604800, // 7 days
	))

	return resp, nil
}

func (s *AuthService) Logout(ctx context.Context, req *connect.Request[authv1.LogoutRequest]) (*connect.Response[authv1.LogoutResponse], error) {
	resp := connect.NewResponse(&authv1.LogoutResponse{
		Success: true,
	})

	// Clear cookies
	resp.Header().Set("Set-Cookie", "access_token=; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=0")
	resp.Header().Add("Set-Cookie", "refresh_token=; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=0")

	return resp, nil
}
