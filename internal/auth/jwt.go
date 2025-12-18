package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserId       string   `json:"user_id"`
	WorkspaceIds []string `json:"workspace_id"`
	jwt.RegisteredClaims
}
type JWTManager struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (m *JWTManager) GenerateAccessToken(userId string, workspaceIds []string) (string, int64, error) {
	expireAt := time.Now().Add(m.accessTokenDuration)
	claims := &Claims{
		UserId:       userId,
		WorkspaceIds: workspaceIds,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expireAt.Unix(), nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
