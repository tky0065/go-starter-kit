package main

// ModelsUserTemplate returns the internal/models/user.go file content with domain entities
func (t *ProjectTemplates) ModelsUserTemplate() string {
	return `package models

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

// AuthResponse represents the authentication response with tokens
type AuthResponse struct {
	AccessToken  string ` + "`json:\"access_token\"`" + `
	RefreshToken string ` + "`json:\"refresh_token\"`" + `
	ExpiresIn    int64  ` + "`json:\"expires_in\"`" + `
}
`
}

// UserEntityTemplate is deprecated - models are now in the models package
func (t *ProjectTemplates) UserEntityTemplate() string {
	return ``
}

// UserRefreshTokenTemplate is deprecated - models are now in the models package
func (t *ProjectTemplates) UserRefreshTokenTemplate() string {
	return ``
}

// UserInterfacesTemplate returns the internal/interfaces/services.go file content
func (t *ProjectTemplates) UserInterfacesTemplate() string {
	return `package interfaces

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

	"` + t.projectName + `/internal/models"
)

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindAll(ctx context.Context, page, limit int) ([]*models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	SaveRefreshToken(ctx context.Context, UserID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
	RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *models.RefreshToken) error
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
	"` + t.projectName + `/internal/models"
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
func (r *UserRepository) CreateUser(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// GetUserByEmail retrieves a user by email, returns nil if not found
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
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
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
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
func (r *UserRepository) FindAll(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// Use the same query base for both Count and Find to ensure consistency
	query := r.db.WithContext(ctx).Model(&models.User{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Updates(u).Error
}

// Delete performs a soft delete on the user (sets deleted_at)
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// SaveRefreshToken saves a refresh token for the given user
func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	return r.db.WithContext(ctx).Create(refreshToken).Error
}

// GetRefreshToken retrieves a refresh token by token string
func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
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
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked", true).Error
}

// RotateRefreshToken performs atomic token rotation: revocation of old and creation of new
func (r *UserRepository) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *models.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Revoke old token with optimistic locking check
		result := tx.Model(&models.RefreshToken{}).
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

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// AppError represents a structured application error with HTTP status and details.
type AppError struct {
	Code    string ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
	Status  int    ` + "`json:\"-\"`" + ` // HTTP Status, not serialized in JSON
	Details any    ` + "`json:\"details,omitempty\"`" + `
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// NewNotFoundError creates a 404 error.
func NewNotFoundError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusNotFound,
		Details: nil,
	}
}

// NewBadRequestError creates a 400 error with optional validation details.
func NewBadRequestError(msg string, code string, details any) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusBadRequest,
		Details: details,
	}
}

// NewInternalError creates a 500 error.
func NewInternalError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusInternalServerError,
		Details: nil,
	}
}

// NewUnauthorizedError creates a 401 error.
func NewUnauthorizedError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusUnauthorized,
		Details: nil,
	}
}

// NewForbiddenError creates a 403 error.
func NewForbiddenError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusForbidden,
		Details: nil,
	}
}

// NewConflictError creates a 409 error.
func NewConflictError(msg string, code string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  fiber.StatusConflict,
		Details: nil,
	}
}

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

// ErrorHandlerMiddlewareTemplate returns the internal/adapters/middleware/error_handler.go file content
func (t *ProjectTemplates) ErrorHandlerMiddlewareTemplate() string {
	return `package middleware

import (
	"errors"
	"os"

	"` + t.projectName + `/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ErrorHandler is a centralized error handler for Fiber that formats all errors
// into a consistent JSON structure following the API standardization requirements.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := fiber.StatusInternalServerError
	resp := fiber.Map{
		"status":  "error",
		"code":    "INTERNAL_SERVER_ERROR",
		"message": "Internal server error",
		"details": nil,
	}

	// Flag to check if we should mask the error message (Production)
	isProd := os.Getenv("APP_ENV") == "production"

	// 1. Handle Domain standard errors (map standard errors to AppErrors)
	if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
		err = domain.NewConflictError("Email already registered", "EMAIL_ALREADY_REGISTERED")
	} else if errors.Is(err, domain.ErrInvalidCredentials) {
		err = domain.NewUnauthorizedError("Invalid email or password", "INVALID_CREDENTIALS")
	} else if errors.Is(err, domain.ErrUserNotFound) {
		err = domain.NewNotFoundError("User not found", "USER_NOT_FOUND")
	} else if errors.Is(err, domain.ErrInvalidRefreshToken) || errors.Is(err, domain.ErrRefreshTokenExpired) || errors.Is(err, domain.ErrRefreshTokenRevoked) {
		err = domain.NewUnauthorizedError(err.Error(), "AUTH_TOKEN_ERROR")
	}

	// 2. Handle Fiber Errors (including 404, 405, etc.)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		resp["message"] = fiberErr.Message
		resp["code"] = mapHTTPStatusToCode(code)
	}

	// 3. Handle Domain AppErrors (business logic errors)
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		code = appErr.Status
		resp["message"] = appErr.Message
		resp["code"] = appErr.Code
		resp["details"] = appErr.Details
	}

	// AC3: Mask internal error messages in production
	if code == fiber.StatusInternalServerError && isProd {
		resp["message"] = "Internal server error"
	}

	// Logging with context
	log.Error().
		Err(err).
		Int("status", code).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Msg("API Error")

	return c.Status(code).JSON(resp)
}

// mapHTTPStatusToCode maps HTTP status codes to readable error codes.
func mapHTTPStatusToCode(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "BAD_REQUEST"
	case fiber.StatusUnauthorized:
		return "UNAUTHORIZED"
	case fiber.StatusForbidden:
		return "FORBIDDEN"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	case fiber.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case fiber.StatusConflict:
		return "CONFLICT"
	case fiber.StatusUnprocessableEntity:
		return "UNPROCESSABLE_ENTITY"
	case fiber.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return "HTTP_ERROR"
	}
}
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
	"` + t.projectName + `/internal/interfaces"
	"` + t.projectName + `/internal/models"
)

// Service handles user business logic
type Service struct {
	repo         interfaces.UserRepository
	tokenService interfaces.TokenService
}

// NewService creates a new user service
func NewService(repo interfaces.UserRepository) *Service {
	return &Service{repo: repo}
}

// NewServiceWithJWT creates a new user service with JWT support
func NewServiceWithJWT(repo interfaces.UserRepository, tokenService interfaces.TokenService) *Service {
	return &Service{
		repo:         repo,
		tokenService: tokenService,
	}
}

// Register creates a new user account with the given email and password
func (s *Service) Register(ctx context.Context, email, password string) (*models.User, error) {
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

	newUser := &models.User{
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// Authenticate validates user credentials and returns JWT tokens
func (s *Service) Authenticate(ctx context.Context, email, password string) (*models.AuthResponse, error) {
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

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RefreshToken validates an existing refresh token and generates new tokens
func (s *Service) RefreshToken(ctx context.Context, oldToken string) (*models.AuthResponse, error) {
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
	newRefreshToken := &models.RefreshToken{
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

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// GetProfile retrieves a user's profile by their ID
func (s *Service) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
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
func (s *Service) GetAll(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
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
func (s *Service) UpdateUser(ctx context.Context, userID uint, email string) (*models.User, error) {
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
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"` + t.projectName + `/internal/domain"
	"` + t.projectName + `/internal/domain/user"
	"` + t.projectName + `/pkg/auth"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service  *user.Service
	validate *validator.Validate
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service *user.Service) *UserHandler {
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
// @Router /users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return domain.NewUnauthorizedError("Unable to extract user information", "UNAUTHORIZED")
	}

	u, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return err // Handled by middleware
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
// @Description Get a list of all users with pagination. Maximum limit is 100 users per page.
// @Tags users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Users per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with data"
// @Failure 500 {object} map[string]string
// @Router /users [get]
// @Security BearerAuth
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	users, total, err := h.service.GetAll(c.Context(), page, limit)
	if err != nil {
		return err // Handled by middleware
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
// @Router /users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil || userID <= 0 {
		return domain.NewBadRequestError("Invalid user ID", "INVALID_ID", nil)
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return domain.NewBadRequestError("Invalid request body", "INVALID_JSON", nil)
	}

	if err := h.validate.Struct(&req); err != nil {
		return domain.NewBadRequestError("Validation failed: "+err.Error(), "VALIDATION_FAILED", nil)
	}

	u, err := h.service.UpdateUser(c.Context(), uint(userID), req.Email)
	if err != nil {
		return err // Handled by middleware
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
// @Router /users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil || userID <= 0 {
		return domain.NewBadRequestError("Invalid user ID", "INVALID_ID", nil)
	}

	err = h.service.DeleteUser(c.Context(), uint(userID))
	if err != nil {
		return err // Handled by middleware
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
		"meta":    fiber.Map{},
	})
}
`
}

// HandlerModuleTemplate returns the internal/adapters/handlers/module.go file content
func (t *ProjectTemplates) HandlerModuleTemplate() string {
	return `package handlers

import (
	"go.uber.org/fx"

	"` + t.projectName + `/internal/domain/user"
)

var Module = fx.Module("handlers",
	fx.Provide(func(s *user.Service) *AuthHandler {
		return NewAuthHandler(s)
	}),
	fx.Provide(func(s *user.Service) *UserHandler {
		return NewUserHandler(s)
	}),
)
`
}

// AuthHandlerTemplate returns the internal/adapters/handlers/auth_handler.go file content
func (t *ProjectTemplates) AuthHandlerTemplate() string {
	return `package handlers

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"` + t.projectName + `/internal/domain"
	"` + t.projectName + `/internal/domain/user"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	service  *user.Service
	validate *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(service *user.Service) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Email    string ` + "`json:\"email\" validate:\"required,email,max=255\"`" + `
	Password string ` + "`json:\"password\" validate:\"required,min=8,max=72\"`" + `
}

// RegisterResponse represents the user registration response
type RegisterResponse struct {
	ID        uint   ` + "`json:\"id\"`" + `
	Email     string ` + "`json:\"email\"`" + `
	CreatedAt string ` + "`json:\"created_at\"`" + `
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} map[string]interface{} "Standard JSON Envelope with user data"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 409 {object} map[string]string "Email already registered"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return domain.NewBadRequestError("Invalid request body", "INVALID_JSON", nil)
	}

	if err := h.validate.Struct(&req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			switch field {
			case "Email":
				validationErrors["email"] = "Email must be valid and max 255 characters"
			case "Password":
				validationErrors["password"] = "Password must be between 8 and 72 characters"
			default:
				validationErrors[field] = err.Error()
			}
		}
		return domain.NewBadRequestError("Validation failed", "VALIDATION_FAILED", validationErrors)
	}

	user, err := h.service.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
			return domain.NewConflictError("Email already registered", "EMAIL_ALREADY_REGISTERED")
		}
		return err // Handled by middleware
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": RegisterResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		"meta": fiber.Map{},
	})
}

// LoginRequest represents the authentication request
type LoginRequest struct {
	Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password string ` + "`json:\"password\" validate:\"required\"`" + `
}

// Login godoc
// @Summary Authenticate user
// @Description Login with email and password to receive JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with tokens"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return domain.NewBadRequestError("Invalid request body", "INVALID_JSON", nil)
	}

	if err := h.validate.Struct(&req); err != nil {
		return domain.NewBadRequestError("Validation failed: email and password required", "VALIDATION_FAILED", nil)
	}

	authResp, err := h.service.Authenticate(c.Context(), req.Email, req.Password)
	if err != nil {
		return err // Handled by middleware
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   authResp,
		"meta":   fiber.Map{},
	})
}

// RefreshRequest represents the token refresh request
type RefreshRequest struct {
	RefreshToken string ` + "`json:\"refresh_token\" validate:\"required\"`" + `
}

// Refresh godoc
// @Summary Refresh access token
// @Description Use refresh token to obtain new access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} map[string]interface{} "Standard JSON Envelope with new tokens"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return domain.NewBadRequestError("Invalid request body", "INVALID_JSON", nil)
	}

	if err := h.validate.Struct(&req); err != nil {
		return domain.NewBadRequestError("Refresh token is required", "VALIDATION_FAILED", nil)
	}

	authResp, err := h.service.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return err // Handled by middleware
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   authResp,
		"meta":   fiber.Map{},
	})
}
`
}

// JWTAuthTemplate returns the pkg/auth/jwt.go file content
func (t *ProjectTemplates) JWTAuthTemplate() string {
	return `package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"` + t.projectName + `/pkg/config"
)

var (
	// ErrInvalidToken is returned when the JWT token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrMissingUserID is returned when user ID is missing from token claims
	ErrMissingUserID = errors.New("missing user ID in token")
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	secretKey string
	expiresIn time.Duration
}

// NewJWTService creates a new JWT service instance
func NewJWTService() *JWTService {
	secret := config.GetEnv("JWT_SECRET", "")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	expiryStr := config.GetEnv("JWT_EXPIRY", "24h")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid JWT_EXPIRY format: %v", err))
	}

	return &JWTService{
		secretKey: secret,
		expiresIn: expiry,
	}
}

// GenerateTokens creates a new JWT access token and refresh token for the given user ID
func (s *JWTService) GenerateTokens(userID uint) (accessToken string, refreshToken string, expiresIn int64, err error) {
	// Create access token claims
	now := time.Now()
	expiresAt := now.Add(s.expiresIn)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
	}

	// Generate access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token (longer expiry, same structure)
	refreshExpiresAt := now.Add(7 * 24 * time.Hour) // 7 days
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     refreshExpiresAt.Unix(),
		"iat":     now.Unix(),
		"type":    "refresh",
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessToken, refreshToken, int64(s.expiresIn.Seconds()), nil
}

// GetUserID extracts the user ID from the JWT token stored in the Fiber context
func GetUserID(c *fiber.Ctx) (uint, error) {
	// Get user from JWT middleware (stored by gofiber/contrib/jwt)
	user := c.Locals("user")
	if user == nil {
		return 0, ErrInvalidToken
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrMissingUserID
	}

	return uint(userIDFloat), nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
`
}

// JWTMiddlewareTemplate returns the pkg/auth/middleware.go file content
func (t *ProjectTemplates) JWTMiddlewareTemplate() string {
	return `package auth

import (
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"` + t.projectName + `/pkg/config"
)

// NewJWTMiddleware creates a new JWT authentication middleware
// It supports both "Bearer <token>" and raw "<token>" formats for Swagger UI compatibility
func NewJWTMiddleware() fiber.Handler {
	secret := config.GetEnv("JWT_SECRET", "")
	if secret == "" {
		panic("JWT_SECRET environment variable is required for middleware")
	}

	// Create the JWT middleware
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(secret),
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"code":    "UNAUTHORIZED",
				"message": "Missing or invalid authentication token",
			})
		},
	})

	// Return a wrapper that normalizes the Authorization header
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		
		// If token is provided without "Bearer " prefix, add it
		// This makes Swagger UI work without typing "Bearer "
		if auth != "" && !strings.HasPrefix(auth, "Bearer ") {
			c.Request().Header.Set("Authorization", "Bearer "+auth)
		}
		
		return jwtMiddleware(c)
	}
}
`
}

// UserModuleTemplate returns the internal/domain/user/module.go file content
func (t *ProjectTemplates) UserModuleTemplate() string {
	return `package user

import (
	"go.uber.org/fx"
)

// Module provides user domain services via fx
var Module = fx.Module("user",
	// Provide service with JWT support - TokenService is injected by fx from auth.Module
	fx.Provide(NewServiceWithJWT),
)
`
}

// RepositoryModuleTemplate returns the internal/adapters/repository/module.go file content
func (t *ProjectTemplates) RepositoryModuleTemplate() string {
	return `package repository

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"` + t.projectName + `/internal/interfaces"
)

// Module provides repository implementations via fx
var Module = fx.Module("repository",
	fx.Provide(func(db *gorm.DB) interfaces.UserRepository {
		return NewUserRepository(db)
	}),
)
`
}

// AuthModuleTemplate returns the pkg/auth/module.go file content
func (t *ProjectTemplates) AuthModuleTemplate() string {
	return `package auth

import (
	"go.uber.org/fx"
	"` + t.projectName + `/internal/interfaces"
)

// Module provides auth services via fx
var Module = fx.Module("auth",
	fx.Provide(func() interfaces.TokenService {
		return NewJWTService()
	}),
	fx.Provide(NewJWTMiddleware),
)
`
}

// RoutesTemplate returns the internal/adapters/http/routes.go file content
func (t *ProjectTemplates) RoutesTemplate() string {
	return `package http

import (
	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"

	"` + t.projectName + `/internal/adapters/handlers"
)

// RegisterRoutes configures all application routes
func RegisterRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware fiber.Handler,
) {
	// Health & Swagger
	RegisterHealthRoutes(app)
	app.Get("/swagger/*", swagger.WrapHandler)

	// API v1
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)

	// User routes (protected)
	users := v1.Group("/users", authMiddleware)
	users.Get("/me", userHandler.GetMe)
	users.Get("", userHandler.GetAllUsers)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
`
}
