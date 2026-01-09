package main

// UserEntityTemplate returns the internal/domain/user/entity.go file content
func (t *ProjectTemplates) UserEntityTemplate() string {
	return `package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents the domain entity for a user.
type User struct {
	ID           uint           ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
	Email        string         ` + "`gorm:\"uniqueIndex;not null\" json:\"email\"`" + `
	PasswordHash string         ` + "`gorm:\"not null\" json:\"-\"`" + `
	CreatedAt    time.Time      ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
	UpdatedAt    time.Time      ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
	DeletedAt    gorm.DeletedAt ` + "`gorm:\"index\" json:\"deleted_at,omitempty\"`" + `
}
`
}

// UserRefreshTokenTemplate returns the internal/domain/user/refresh_token.go file content
func (t *ProjectTemplates) UserRefreshTokenTemplate() string {
	return `package user

import "time"

// RefreshToken represents a refresh token for session management
type RefreshToken struct {
	ID        uint      ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
	UserID    uint      ` + "`gorm:\"not null;index\" json:\"user_id\"`" + `
	Token     string    ` + "`gorm:\"uniqueIndex;not null\" json:\"token\"`" + `
	ExpiresAt time.Time ` + "`gorm:\"not null\" json:\"expires_at\"`" + `
	Revoked   bool      ` + "`gorm:\"not null;default:false\" json:\"revoked\"`" + `
	CreatedAt time.Time ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
}

// IsExpired returns true if the token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsRevoked returns true if the token has been revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.Revoked
}
`
}

// UserInterfacesTemplate returns the internal/interfaces/services.go file content
func (t *ProjectTemplates) UserInterfacesTemplate() string {
	return `package interfaces

import (
	"context"

	"` + t.projectName + `/internal/domain/user"
)

// AuthService defines the interface for authentication operations.
// Implemented by internal/domain/user/Service.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*user.User, error)
	Authenticate(ctx context.Context, email, password string) (*user.AuthResponse, error)
	RefreshToken(ctx context.Context, oldToken string) (*user.AuthResponse, error)
}

// UserService defines the business logic operations for user management.
// Implemented by internal/domain/user/Service.
type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*user.User, error)
	GetAll(ctx context.Context, page, limit int) ([]*user.User, int64, error)
	UpdateUser(ctx context.Context, userID uint, email string) (*user.User, error)
	DeleteUser(ctx context.Context, userID uint) error
}

// TokenService defines the interface for token generation.
// Implemented by pkg/auth/JWTService.
type TokenService interface {
	GenerateTokens(userID uint) (accessToken string, refreshToken string, expiresIn int64, err error)
}
`
}

// UserRepositoryInterfaceTemplate returns the internal/interfaces/user_repository.go file content
func (t *ProjectTemplates) UserRepositoryInterfaceTemplate() string {
	return `package interfaces

import (
	"context"
	"time"

	"` + t.projectName + `/internal/domain/user"
)

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *user.User) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id uint) (*user.User, error)
	FindAll(ctx context.Context, page, limit int) ([]*user.User, int64, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id uint) error
	SaveRefreshToken(ctx context.Context, UserID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*user.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
	RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *user.RefreshToken) error
}
`
}

// UserRepositoryTemplate returns the internal/adapters/repository/user_repository.go file content
func (t *ProjectTemplates) UserRepositoryTemplate() string {
	return `package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"` + t.projectName + `/internal/domain"
	"` + t.projectName + `/internal/domain/user"
)

// UserRepository implements user data persistence using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// GetUserByEmail retrieves a user by email, returns nil if not found
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByID retrieves a user by ID, returns nil if not found
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindAll retrieves all users from the database (excluding soft-deleted)
func (r *UserRepository) FindAll(ctx context.Context, page, limit int) ([]*user.User, int64, error) {
	var users []*user.User
	var total int64
	
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

// Delete performs a soft delete on the user (sets deleted_at)
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, id).Error
}

// SaveRefreshToken saves a refresh token for the given user
func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	refreshToken := &user.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	return r.db.WithContext(ctx).Create(refreshToken).Error
}

// GetRefreshToken retrieves a refresh token by token string
func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (*user.RefreshToken, error) {
	var rt user.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (r *UserRepository) RevokeRefreshToken(ctx context.Context, tokenID uint) error {
	return r.db.WithContext(ctx).Model(&user.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked", true).Error
}

// RotateRefreshToken performs atomic token rotation: revocation of old and creation of new
func (r *UserRepository) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *user.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Revoke old token with optimistic locking check
		result := tx.Model(&user.RefreshToken{}).
			Where("id = ? AND revoked = ?", oldTokenID, false).
			Update("revoked", true)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return domain.ErrRefreshTokenRevoked
		}

		// 2. Create new token
		if err := tx.Create(newToken).Error; err != nil {
			return err
		}

		return nil
	})
}
`
}

// DomainErrorsTemplate returns the internal/domain/errors.go file content
func (t *ProjectTemplates) DomainErrorsTemplate() string {
	return `package domain

import "errors"

// Domain-wide errors
var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrRefreshTokenExpired    = errors.New("refresh token expired")
	ErrRefreshTokenRevoked    = errors.New("refresh token has been revoked")
)
`
}

// UserServiceTemplate returns the internal/domain/user/service.go file content
func (t *ProjectTemplates) UserServiceTemplate() string {
	return `package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"` + t.projectName + `/internal/domain"
)

// UserRepository defines the persistence interface required by the service
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	FindAll(ctx context.Context, page, limit int) ([]*User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
	RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *RefreshToken) error
}

// TokenService defines the token generation interface required by the service
type TokenService interface {
	GenerateTokens(userID uint) (accessToken string, refreshToken string, expiresIn int64, err error)
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
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrEmailAlreadyRegistered
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

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
	AccessToken  string ` + "`json:\"access_token\"`" + `
	RefreshToken string ` + "`json:\"refresh_token\"`" + `
	ExpiresIn    int64  ` + "`json:\"expires_in\"`" + `
}

// Authenticate validates user credentials and returns JWT tokens
func (s *Service) Authenticate(ctx context.Context, email, password string) (*AuthResponse, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if u == nil {
		return nil, domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, refreshToken, expiresIn, err := s.tokenService.GenerateTokens(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

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
	rt, err := s.repo.GetRefreshToken(ctx, oldToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if rt == nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	if rt.IsExpired() {
		return nil, domain.ErrRefreshTokenExpired
	}

	if rt.IsRevoked() {
		fmt.Printf("SECURITY ALERT: Attempt to use revoked refresh token ID: %d UserID: %d\n", rt.ID, rt.UserID)
		return nil, domain.ErrRefreshTokenRevoked
	}

	accessToken, refreshToken, expiresIn, err := s.tokenService.GenerateTokens(rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	newRefreshToken := &RefreshToken{
		UserID:    rt.UserID,
		Token:     refreshToken,
		ExpiresAt: refreshExpiresAt,
		Revoked:   false,
	}

	err = s.repo.RotateRefreshToken(ctx, rt.ID, newRefreshToken)
	if err != nil {
		if err == domain.ErrRefreshTokenRevoked {
			fmt.Printf("SECURITY ALERT: Race condition on refresh token rotation ID: %d\n", rt.ID)
			return nil, domain.ErrRefreshTokenRevoked
		}
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// GetProfile retrieves a user's profile by their ID
func (s *Service) GetProfile(ctx context.Context, userID uint) (*User, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if u == nil {
		return nil, domain.ErrUserNotFound
	}

	return u, nil
}

// GetAll retrieves all users from the database
func (s *Service) GetAll(ctx context.Context, page, limit int) ([]*User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	users, total, err := s.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, total, nil
}

// UpdateUser updates a user's email address
func (s *Service) UpdateUser(ctx context.Context, userID uint, email string) (*User, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if u == nil {
		return nil, domain.ErrUserNotFound
	}

	if email != u.Email {
		existing, err := s.repo.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existing != nil {
			return nil, domain.ErrEmailAlreadyRegistered
		}
	}

	u.Email = email

	err = s.repo.Update(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return u, nil
}

// DeleteUser performs a soft delete on a user
func (s *Service) DeleteUser(ctx context.Context, userID uint) error {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if u == nil {
		return domain.ErrUserNotFound
	}

	err = s.repo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
`
}

// UserHandlerTemplate returns the internal/adapters/handlers/user_handler.go file content
func (t *ProjectTemplates) UserHandlerTemplate() string {
	return `package handlers

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"` + t.projectName + `/internal/domain"
	"` + t.projectName + `/internal/interfaces"
	"` + t.projectName + `/pkg/auth"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service  interfaces.UserService
	validate *validator.Validate
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service interfaces.UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

// ProfileResponse represents the user profile response
type ProfileResponse struct {
	ID        uint   ` + "`json:\"id\"`" + `
	Email     string ` + "`json:\"email\"`" + `
	CreatedAt string ` + "`json:\"created_at\"`" + `
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with data"
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"error":  "Unable to extract user information",
		})
	}

	u, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"error":  "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to retrieve user profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": ProfileResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		},
		"meta": fiber.Map{},
	})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with data"
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [get]
// @Security BearerAuth
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	users, total, err := h.service.GetAll(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to retrieve users",
		})
	}

	userResponses := make([]ProfileResponse, len(users))
	for i, u := range users {
		userResponses[i] = ProfileResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   userResponses,
		"meta":   fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email string ` + "`json:\"email\" validate:\"required,email\"`" + `
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update user request"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with data"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid user ID",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Validation failed: " + err.Error(),
		})
	}

	u, err := h.service.UpdateUser(c.Context(), uint(userID), req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"error":  "User not found",
			})
		}
		if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "Email already in use",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": ProfileResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		},
		"meta": fiber.Map{},
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Soft delete a user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid user ID",
		})
	}

	err = h.service.DeleteUser(c.Context(), uint(userID))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"error":  "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
		"meta":    fiber.Map{},
	})
}
`
}
