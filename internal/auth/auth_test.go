package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenLifecycle(t *testing.T) {
	secret := "test-secret-key-12345"
	userID := int64(42)
	email := "testuser@example.com"
	role := "customer"
	duration := 1 * time.Hour

	// 1. Generate Token
	tokenStr, err := GenerateToken(userID, email, role, secret, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// 2. Validate Token
	claims, err := ValidateToken(tokenStr, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)

	// 3. Validate with Wrong Secret
	_, err = ValidateToken(tokenStr, "wrong-secret-key")
	assert.Error(t, err)
}

func TestRequestValidation(t *testing.T) {
	t.Run("RegisterRequest", func(t *testing.T) {
		// Valid request
		req := RegisterRequest{Email: "newuser@example.com", Password: "mypassword"}
		assert.NoError(t, req.Validate())

		// Missing email
		req = RegisterRequest{Email: "", Password: "mypassword"}
		assert.Error(t, req.Validate())

		// Invalid email format
		req = RegisterRequest{Email: "invalidemail", Password: "mypassword"}
		assert.Error(t, req.Validate())

		// Password too short
		req = RegisterRequest{Email: "newuser@example.com", Password: "123"}
		assert.Error(t, req.Validate())
	})

	t.Run("LoginRequest", func(t *testing.T) {
		// Valid request
		req := LoginRequest{Email: "user@example.com", Password: "password"}
		assert.NoError(t, req.Validate())

		// Missing email
		req = LoginRequest{Email: "", Password: "password"}
		assert.Error(t, req.Validate())

		// Invalid email format
		req = LoginRequest{Email: "invalid", Password: "password"}
		assert.Error(t, req.Validate())

		// Missing password
		req = LoginRequest{Email: "user@example.com", Password: ""}
		assert.Error(t, req.Validate())
	})
}
