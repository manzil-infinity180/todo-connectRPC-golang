package auth

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey       contextKey = "user_id"
	WorkspaceIDsKey contextKey = "workspace_ids"
)

type AuthInterceptor struct {
	jwtManager *JWTManager
}

func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
	}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if i.isPublicEndpoint(req.Spec().Procedure) {
			return next(ctx, req)
		}
		token, err := i.extractToken(req.Header())
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing or invalid token: %w", err))
		}
		claims, err := i.jwtManager.ValidateToken(token)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid token: %w", err))
		}

		ctx = context.WithValue(ctx, UserIDKey, claims.UserId)
		ctx = context.WithValue(ctx, WorkspaceIDsKey, claims.WorkspaceIds)

		return next(ctx, req)
	}
}

func (i *AuthInterceptor) isPublicEndpoint(procedure string) bool {
	publicEndpoints := []string{
		"/auth.v1.AuthService/Register",
		"/auth.v1.AuthService/Login",
		"/auth.v1.AuthService/RefreshToken",
	}
	for _, endpoint := range publicEndpoints {
		if procedure == endpoint {
			return true
		}
	}
	return false
}

func (i *AuthInterceptor) extractToken(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}
	cookie := header.Get("Cookie")
	if cookie != "" {
		cookies := strings.Split(cookie, "; ")
		for _, c := range cookies {
			parts := strings.SplitN(c, "=", 2)
			if len(parts) == 2 && parts[0] == "access_token" {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("token not found")
}
