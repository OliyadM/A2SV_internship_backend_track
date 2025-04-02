package infrastructure

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	passService := NewPasswordService()

	t.Run("Success", func(t *testing.T) {
		password := "mypassword123"
		hashed, err := passService.HashPassword(password)
		assert.NoError(t, err, "HashPassword should not return an error")
		assert.NotEmpty(t, hashed, "Hashed password should not be empty")
		assert.NotEqual(t, password, hashed, "Hashed password should differ from plain text")

		// Verify the hash is valid by comparing it with bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		assert.NoError(t, err, "Generated hash should match the original password")
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		hashed, err := passService.HashPassword("")
		assert.NoError(t, err, "HashPassword should handle empty password without error")
		assert.NotEmpty(t, hashed, "Hashed password should still be generated")

		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(""))
		assert.NoError(t, err, "Hash of empty password should be valid")
	})
}

func TestComparePassword(t *testing.T) {
	passService := NewPasswordService()

	t.Run("Success", func(t *testing.T) {
		password := "mypassword123"
		hashed, err := passService.HashPassword(password)
		assert.NoError(t, err)

		err = passService.ComparePassword(hashed, password)
		assert.NoError(t, err, "ComparePassword should succeed with correct password")
	})

	t.Run("WrongPassword", func(t *testing.T) {
		password := "mypassword123"
		hashed, err := passService.HashPassword(password)
		assert.NoError(t, err)

		err = passService.ComparePassword(hashed, "wrongpassword")
		assert.Error(t, err, "ComparePassword should fail with incorrect password")
		assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err, "Error should be mismatch error")
	})

	t.Run("InvalidHash", func(t *testing.T) {
		err := passService.ComparePassword("invalid-hash", "mypassword123")
		assert.Error(t, err, "ComparePassword should fail with invalid hash")
		assert.NotEqual(t, bcrypt.ErrMismatchedHashAndPassword, err, "Error should indicate invalid hash format")
	})

	t.Run("EmptyHash", func(t *testing.T) {
		err := passService.ComparePassword("", "mypassword123")
		assert.Error(t, err, "ComparePassword should fail with empty hash")
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		hashed, err := passService.HashPassword("")
		assert.NoError(t, err)

		err = passService.ComparePassword(hashed, "")
		assert.NoError(t, err, "ComparePassword should succeed with empty password if hashed correctly")
	})
}
