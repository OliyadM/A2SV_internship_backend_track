package infrastructure

import (
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret"
	jwtService := NewJWTService(secret)

	t.Run("Success", func(t *testing.T) {
		tokenString, err := jwtService.GenerateToken("1", "testuser", "user")
		assert.NoError(t, err, "GenerateToken should not return an error")
		assert.NotEmpty(t, tokenString, "Token string should not be empty")

		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		assert.NoError(t, err, "Token should parse without error")
		assert.True(t, token.Valid, "Token should be valid")

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok, "Claims should be of type MapClaims")
		assert.Equal(t, "1", claims["sub"], "User ID should match")
		assert.Equal(t, "testuser", claims["name"], "Username should match")
		assert.Equal(t, "user", claims["role"], "Role should match")

		exp, ok := claims["exp"].(float64)
		assert.True(t, ok, "Expiration should be a number")
		assert.InDelta(t, time.Now().Add(24*time.Hour).Unix(), int64(exp), 2, "Expiration should be ~24 hours from now")
	})
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	jwtService := NewJWTService(secret)

	t.Run("Success", func(t *testing.T) {
		
		tokenString, err := jwtService.GenerateToken("1", "testuser", "user")
		assert.NoError(t, err)

		
		token, err := jwtService.ValidateToken(tokenString)
		assert.NoError(t, err, "ValidateToken should not return an error")
		assert.True(t, token.Valid, "Token should be valid")

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "1", claims["sub"])
		assert.Equal(t, "testuser", claims["name"])
		assert.Equal(t, "user", claims["role"])
	})

	t.Run("InvalidToken", func(t *testing.T) {
		
		_, err := jwtService.ValidateToken("invalid.token.here")
		assert.Error(t, err, "ValidateToken should return an error for invalid token")
	})

	t.Run("WrongSecret", func(t *testing.T) {
		
		wrongService := NewJWTService("wrong-secret")
		tokenString, err := wrongService.GenerateToken("1", "testuser", "user")
		assert.NoError(t, err)

		
		_, err = jwtService.ValidateToken(tokenString)
		assert.Error(t, err, "ValidateToken should fail with wrong secret")
		assert.Contains(t, err.Error(), "signature is invalid", "Error should indicate invalid signature")
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		// Create a custom token with an expired time
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  "1",
			"name": "testuser",
			"role": "user",
			"exp":  time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		})
		tokenString, err := token.SignedString([]byte(secret))
		assert.NoError(t, err)

		
		_, err = jwtService.ValidateToken(tokenString)
		assert.Error(t, err, "ValidateToken should fail for expired token")
		assert.Contains(t, err.Error(), "token is expired", "Error should indicate expiration")
	})
}
