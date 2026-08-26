package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"product-inventory/internal/database/db"
)

var (
	ErrEmailTaken         = errors.New("email is already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// Service coordinates registration, login, hashing, and token operations.
type Service struct {
	queries   *db.Queries
	jwtSecret string
}

// NewService creates a new Service instance.
func NewService(queries *db.Queries, jwtSecret string) *Service {
	return &Service{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

// Register hashes password, registers a user, and returns an AuthResponse.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 1. Check if user already exists
	_, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailTaken
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// 2. Hash password using Bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Set default role if not provided
	role := req.Role
	if role == "" || (role != "admin" && role != "customer" && role != "system") {
		role = "customer"
	}

	// 4. Create user
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         role,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// 5. Generate JWT token (expires in 24 hours)
	token, err := GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

// Login validates user credentials and returns an AuthResponse with a JWT token.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 1. Fetch user by email
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 2. Compare hashed password with requested password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. Generate token
	token, err := GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
