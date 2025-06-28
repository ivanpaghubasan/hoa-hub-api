package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/constants"
)

type IJWTManager interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ParseToken(tokenStr string) (*Claims, error)
}

type JWTManager struct {
	secretKey string
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string) IJWTManager {
	return &JWTManager{secretKey: secret}
}

func (j *JWTManager) GenerateToken(userID uuid.UUID) (string, error) {
	claims := &Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, constants.ErrInvalidToken
	}
	return token.Claims.(*Claims), nil
}
