package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenCreatorDependencies struct {
	SecretKey string
	TokenTTL  time.Duration
}

type JWTTokenManager struct {
	secretKey string
	tokenTTL  time.Duration
}

type claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func NewJWTTokenManager(deps JWTTokenCreatorDependencies) *JWTTokenManager {
	return &JWTTokenManager{
		secretKey: deps.SecretKey,
		tokenTTL:  deps.TokenTTL,
	}
}

func (creator *JWTTokenManager) Create(ctx context.Context, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(creator.tokenTTL).Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte(creator.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (creator *JWTTokenManager) ExtractUserID(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(creator.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	if claims.UserID == "" {
		return "", fmt.Errorf("invalid token claims, user id not found")
	}

	return claims.UserID, nil
}
