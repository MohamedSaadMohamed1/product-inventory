package auth

import (
	"errors"
	"strings"
)

// RegisterRequest defines the request payload for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // Optional, defaults to 'customer'
}

// Validate validates the RegisterRequest.
func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return errors.New("invalid email address format")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

// LoginRequest defines the request payload for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the LoginRequest.
func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return errors.New("invalid email address format")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// AuthResponse defines the payload returned on successful register/login.
type AuthResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
