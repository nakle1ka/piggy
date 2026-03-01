package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	TokenId string `json:"token_id"`
	jwt.RegisteredClaims
}

var ErrInvalidTokenPayload = errors.New("Invalid token payload")

type TokenManager interface {
	GenerateToken(userID string, exp time.Duration) (string, string, error)
	ValidateToken(token string) (*TokenClaims, error)
}

type tokenManager struct {
	secretKey string
}

func (t *tokenManager) GenerateToken(userID string, exp time.Duration) (string, string, error) {
	tokenId := uuid.New().String()

	claims := &TokenClaims{
		TokenId: tokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.secretKey))

	return signedToken, tokenId, err
}

func (t *tokenManager) ValidateToken(token string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		return claims, nil
	}

	return nil, ErrInvalidTokenPayload
}

func NewTokenManager(secretKey string) TokenManager {
	return &tokenManager{
		secretKey: secretKey,
	}
}
