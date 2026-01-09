package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRepository defines the persistence interface required by the service
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	SaveRefreshToken(ctx context.Context, UserID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
	RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *RefreshToken) error
}

// TokenService defines the token generation interface required by the service
type TokenService interface {
	GenerateTokens(UserID uint) (accessToken string, refreshToken string, expiresIn int64, err error)
}

// Service handles user business logic
type Service struct {
	repo         UserRepository
	tokenService TokenService
}

// NewService creates a new user service
func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

// NewServiceWithJWT creates a new user service with JWT support
func NewServiceWithJWT(repo UserRepository, tokenService TokenService) *Service {
	return &Service{
		repo:         repo,
		tokenService: tokenService,
	}
}

// Register creates a new user account with the given email and password
func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {
	// Check if user exists
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	// Hash password with explicit cost validation
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Validate hash was generated
	if len(hashedBytes) == 0 {
		return nil, fmt.Errorf("password hash generation produced empty result")
	}

	newUser := &User{
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// AuthResponse represents the authentication response with tokens
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Authenticate validates user credentials and returns JWT tokens
func (s *Service) Authenticate(ctx context.Context, email, password string) (*AuthResponse, error) {
	// Get user by email
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user exists
	if u == nil {
		return nil, ErrInvalidCredentials
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, expiresIn, err := s.tokenService.GenerateTokens(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Calculate refresh token expiration (7 days)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Save refresh token
	err = s.repo.SaveRefreshToken(ctx, u.ID, refreshToken, refreshExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RefreshToken validates an existing refresh token and generates new tokens
func (s *Service) RefreshToken(ctx context.Context, oldToken string) (*AuthResponse, error) {
	// Retrieve the refresh token from database
	rt, err := s.repo.GetRefreshToken(ctx, oldToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Validate token exists
	if rt == nil {
		return nil, ErrInvalidRefreshToken
	}

	// Validate token is not expired
	if rt.IsExpired() {
		return nil, ErrRefreshTokenExpired
	}

	// Validate token is not revoked
	if rt.IsRevoked() {
		// LOG SECURITY ALERT: Revoked token usage attempt (Potential theft)
		fmt.Printf("SECURITY ALERT: Attempt to use revoked refresh token ID: %d UserID: %d\n", rt.ID, rt.UserID)
		return nil, ErrRefreshTokenRevoked
	}

	// Generate new tokens
	accessToken, refreshToken, expiresIn, err := s.tokenService.GenerateTokens(rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Prepare new refresh token object
	// Calculate new refresh token expiration (7 days)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	newRefreshToken := &RefreshToken{
		UserID:    rt.UserID,
		Token:     refreshToken,
		ExpiresAt: refreshExpiresAt,
		Revoked:   false,
	}

	// Perform atomic rotation (Revoke Old + Save New)
	err = s.repo.RotateRefreshToken(ctx, rt.ID, newRefreshToken)
	if err != nil {
		if err == ErrRefreshTokenRevoked {
			// Race condition detected
			fmt.Printf("SECURITY ALERT: Race condition on refresh token rotation ID: %d\n", rt.ID)
			return nil, ErrRefreshTokenRevoked
		}
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
