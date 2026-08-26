package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenLifecycle(t *testing.T) {
	secret := "test-secret-key-12345"
	userID := int64(42)
	username := "testuser"
	role := "customer"
	duration := 1 * time.Hour

	// 1. Generate Token
	tokenStr, err := GenerateToken(userID, username, role, secret, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// 2. Validate Token
	claims, err := ValidateToken(tokenStr, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)

	// 3. Validate with Wrong Secret
	_, err = ValidateToken(tokenStr, "wrong-secret-key")
	assert.Error(t, err)
}

func TestRequestValidation(t *testing.T) {
	t.Run("RegisterRequest", func(t *testing.T) {
		// Valid request
		req := RegisterRequest{Username: "newuser", Password: "mypassword"}
		assert.NoError(t, req.Validate())

		// Missing username
		req = RegisterRequest{Username: "", Password: "mypassword"}
		assert.Error(t, req.Validate())

		// Password too short
		req = RegisterRequest{Username: "newuser", Password: "123"}
		assert.Error(t, req.Validate())
	})

	t.Run("LoginRequest", func(t *testing.T) {
		// Valid request
		req := LoginRequest{Username: "user", Password: "password"}
		assert.NoError(t, req.Validate())

		// Missing username
		req = LoginRequest{Username: "", Password: "password"}
		assert.Error(t, req.Validate())

		// Missing password
		req = LoginRequest{Username: "user", Password: ""}
		assert.Error(t, req.Validate())
	})
}
