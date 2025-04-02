package infrastructure

import (
	"task_manager/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServiceImpl struct {
	secret string
}

func NewJWTService(secret string) domain.JWTService {
	return &JWTServiceImpl{secret: secret}
}

func (s *JWTServiceImpl) GenerateToken(userID, username, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"name": username,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.secret))
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
}
