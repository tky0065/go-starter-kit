package main

// GraphQLGoModTemplate returns the go.mod file content for GraphQL template.
// Includes gqlgen and gofiber/adaptor for net/http handler integration.
func (t *ProjectTemplates) GraphQLGoModTemplate() string {
	return `module ` + t.projectName + `

go 1.25.5

require (
	github.com/99designs/gqlgen v0.17.73
	github.com/go-playground/validator/v10 v10.30.1
	github.com/gofiber/adaptor/v2 v2.2.1
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.33.0
	github.com/vektah/gqlparser/v2 v2.5.27
	go.uber.org/fx v1.24.0
	golang.org/x/crypto v0.31.0
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.31.1
)
`
}

// GraphQLMainGoTemplate returns the cmd/main.go file content for GraphQL template.
func (t *ProjectTemplates) GraphQLMainGoTemplate() string {
	return `package main

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/fx"

	"` + t.projectName + `/internal/infrastructure/database"
	"` + t.projectName + `/internal/infrastructure/server"
	"` + t.projectName + `/pkg/logger"
)

// @title ` + t.projectName + ` GraphQL API
// @version 1.0
// @description A GraphQL API built with Go, gqlgen, Fiber, and GORM
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	// Load environment variables from .env file
	// This is primarily for local development; in production, use system environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or couldn't be loaded")
	}

	fx.New(
		// Core infrastructure
		logger.Module,
		database.Module,

		// HTTP server with GraphQL (must be last as it depends on other modules)
		server.Module,
	).Run()
}
`
}

// GqlGenYmlTemplate returns the gqlgen.yml configuration file content.
func (t *ProjectTemplates) GqlGenYmlTemplate() string {
	return `# gqlgen configuration file
# See https://gqlgen.com/config/ for documentation

# Schema files to load
schema:
  - graph/*.graphqls

# Where to generate the server code
exec:
  filename: graph/generated/generated.go
  package: generated

# Where to generate the models
model:
  filename: graph/model/models_gen.go
  package: model

# Where to put the resolver implementations
resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

# Enable autobind to automatically bind Go types to GraphQL types
autobind:
  - "` + t.projectName + `/graph/model"
  - "` + t.projectName + `/internal/models"

# Model mappings - map GraphQL types to existing Go types
models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  User:
    model: ` + t.projectName + `/internal/models.User
`
}

// GraphQLSchemaTemplate returns the GraphQL schema file content.
func (t *ProjectTemplates) GraphQLSchemaTemplate() string {
	return `# GraphQL Schema for ` + t.projectName + `
# This schema defines the types, queries, and mutations for the API

# Scalar types for custom data
scalar Time

# User type - represents a user in the system
type User {
  id: ID!
  email: String!
  createdAt: Time!
  updatedAt: Time!
}

# Input type for creating a new user
input NewUser {
  email: String!
  password: String!
}

# Input type for updating a user
input UpdateUser {
  email: String
}

# Pagination info for list queries
type PageInfo {
  page: Int!
  limit: Int!
  total: Int!
  hasNextPage: Boolean!
}

# Paginated users response
type UsersConnection {
  users: [User!]!
  pageInfo: PageInfo!
}

# Root Query type - all read operations
type Query {
  # Get a user by ID
  user(id: ID!): User
  
  # Get all users with pagination
  users(page: Int = 1, limit: Int = 10): UsersConnection!
  
  # Health check
  health: String!
}

# Root Mutation type - all write operations
type Mutation {
  # Create a new user
  createUser(input: NewUser!): User!
  
  # Update an existing user
  updateUser(id: ID!, input: UpdateUser!): User!
  
  # Delete a user (soft delete)
  deleteUser(id: ID!): Boolean!
}
`
}

// GraphQLResolverTemplate returns the resolver.go file content.
func (t *ProjectTemplates) GraphQLResolverTemplate() string {
	return `package graph

import (
	"` + t.projectName + `/internal/interfaces"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver holds the dependencies for GraphQL resolvers.
// It follows the hexagonal architecture pattern by depending on interfaces.
type Resolver struct {
	UserRepo interfaces.UserRepository
}

// NewResolver creates a new Resolver with the required dependencies.
func NewResolver(userRepo interfaces.UserRepository) *Resolver {
	return &Resolver{
		UserRepo: userRepo,
	}
}
`
}

// GraphQLSchemaResolversTemplate returns the schema.resolvers.go file content.
func (t *ProjectTemplates) GraphQLSchemaResolversTemplate() string {
	return `package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.73

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"` + t.projectName + `/graph/generated"
	"` + t.projectName + `/graph/model"
	"` + t.projectName + `/internal/models"
)

var emailRegex = regexp.MustCompile(` + "`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`" + `)

// CreateUser creates a new user in the database.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*models.User, error) {
	// Validate email format
	if !emailRegex.MatchString(input.Email) {
		log.Warn().Str("email", input.Email).Msg("Invalid email format")
		return nil, errors.New("invalid email format")
	}

	// Normalize email
	normalizedEmail := strings.ToLower(strings.TrimSpace(input.Email))

	// Validate password length
	if len(input.Password) < 8 {
		log.Warn().Msg("Password too short")
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if email already exists
	existingUser, err := r.UserRepo.GetUserByEmail(ctx, normalizedEmail)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check existing email")
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUser != nil {
		log.Warn().Str("email", normalizedEmail).Msg("Email already exists")
		return nil, errors.New("email already exists")
	}

	// Hash the password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:        normalizedEmail,
		PasswordHash: string(hashedBytes),
	}

	if err := r.UserRepo.CreateUser(ctx, user); err != nil {
		log.Error().Err(err).Str("email", normalizedEmail).Msg("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Info().Uint("user_id", user.ID).Str("email", normalizedEmail).Msg("User created successfully")
	return user, nil
}

// UpdateUser updates an existing user.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUser) (*models.User, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := r.UserRepo.FindByID(ctx, uint(userID))
	if err != nil {
		log.Error().Err(err).Uint("user_id", uint(userID)).Msg("Failed to find user")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		log.Warn().Uint("user_id", uint(userID)).Msg("User not found")
		return nil, errors.New("user not found")
	}

	// Validate and update email if provided
	if input.Email != nil {
		if !emailRegex.MatchString(*input.Email) {
			return nil, errors.New("invalid email format")
		}
		normalizedEmail := strings.ToLower(strings.TrimSpace(*input.Email))
		
		// Check if new email already exists (for different user)
		existingUser, err := r.UserRepo.GetUserByEmail(ctx, normalizedEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
		
		user.Email = normalizedEmail
	}

	if err := r.UserRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Uint("user_id", user.ID).Msg("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	log.Info().Uint("user_id", user.ID).Msg("User updated successfully")
	return user, nil
}

// DeleteUser soft-deletes a user.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	if err := r.UserRepo.Delete(ctx, uint(userID)); err != nil {
		log.Error().Err(err).Uint("user_id", uint(userID)).Msg("Failed to delete user")
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	log.Info().Uint("user_id", uint(userID)).Msg("User deleted successfully")
	return true, nil
}

// User retrieves a single user by ID.
func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := r.UserRepo.FindByID(ctx, uint(userID))
	if err != nil {
		log.Error().Err(err).Uint("user_id", uint(userID)).Msg("Failed to find user")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		log.Debug().Uint("user_id", uint(userID)).Msg("User not found")
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Users retrieves all users with pagination.
func (r *queryResolver) Users(ctx context.Context, page *int, limit *int) (*model.UsersConnection, error) {
	p := 1
	l := 10

	if page != nil && *page > 0 {
		p = *page
	}
	if limit != nil {
		if *limit > 100 {
			return nil, errors.New("limit cannot exceed 100")
		}
		if *limit > 0 {
			l = *limit
		}
	}

	users, total, err := r.UserRepo.FindAll(ctx, p, l)
	if err != nil {
		log.Error().Err(err).Int("page", p).Int("limit", l).Msg("Failed to find users")
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	hasNextPage := int64(p*l) < total

	log.Debug().Int("page", p).Int("limit", l).Int64("total", total).Msg("Users query executed")

	return &model.UsersConnection{
		Users: users,
		PageInfo: &model.PageInfo{
			Page:        p,
			Limit:       l,
			Total:       int(total),
			HasNextPage: hasNextPage,
		},
	}, nil
}

// Health returns a health check status.
func (r *queryResolver) Health(ctx context.Context) (string, error) {
	return "ok", nil
}

// Mutation returns the mutation resolver.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns the query resolver.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
`
}

// GraphQLModelTemplate returns the graph/model/models.go file content.
func (t *ProjectTemplates) GraphQLModelTemplate() string {
	return `package model

import "` + t.projectName + `/internal/models"

// NewUser is the input type for creating a new user.
type NewUser struct {
	Email    string ` + "`json:\"email\"`" + `
	Password string ` + "`json:\"password\"`" + `
}

// UpdateUser is the input type for updating a user.
type UpdateUser struct {
	Email *string ` + "`json:\"email,omitempty\"`" + `
}

// PageInfo contains pagination information.
type PageInfo struct {
	Page        int  ` + "`json:\"page\"`" + `
	Limit       int  ` + "`json:\"limit\"`" + `
	Total       int  ` + "`json:\"total\"`" + `
	HasNextPage bool ` + "`json:\"hasNextPage\"`" + `
}

// UsersConnection is the paginated response for users query.
type UsersConnection struct {
	Users    []*models.User ` + "`json:\"users\"`" + `
	PageInfo *PageInfo      ` + "`json:\"pageInfo\"`" + `
}
`
}

// GraphQLGeneratedTemplate returns a placeholder for graph/generated/generated.go.
// This file will be overwritten when running 'go generate ./...' or 'go run github.com/99designs/gqlgen generate'
func (t *ProjectTemplates) GraphQLGeneratedTemplate() string {
	return `// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.
// This is a placeholder file. Run 'go generate ./...' to generate the actual code.

package generated

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"` + t.projectName + `/internal/models"
	"` + t.projectName + `/graph/model"
)

// Config holds the configuration for the GraphQL server
type Config struct {
	Resolvers  ResolverRoot
	Directives DirectiveRoot
	Complexity ComplexityRoot
}

// ResolverRoot is the root resolver interface
type ResolverRoot interface {
	Mutation() MutationResolver
	Query() QueryResolver
}

// DirectiveRoot holds directive implementations
type DirectiveRoot struct{}

// ComplexityRoot holds complexity functions
type ComplexityRoot struct{}

// MutationResolver is the interface for mutation operations
type MutationResolver interface {
	CreateUser(ctx context.Context, input model.NewUser) (*models.User, error)
	UpdateUser(ctx context.Context, id string, input model.UpdateUser) (*models.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
}

// QueryResolver is the interface for query operations
type QueryResolver interface {
	User(ctx context.Context, id string) (*models.User, error)
	Users(ctx context.Context, page *int, limit *int) (*model.UsersConnection, error)
	Health(ctx context.Context) (string, error)
}

// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.
// This is a placeholder - run 'go generate ./...' to generate the actual implementation.
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	panic("Run 'go generate ./...' or 'go run github.com/99designs/gqlgen generate' to generate this file")
}
`
}

// GraphQLServerTemplate returns the internal/infrastructure/server/server.go for GraphQL template.
func (t *ProjectTemplates) GraphQLServerTemplate() string {
	return `// Package server provides HTTP server configuration and lifecycle management.
// It creates and configures a Fiber application with GraphQL support using
// gqlgen and gofiber/adaptor for net/http handler integration.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"` + t.projectName + `/graph"
	"` + t.projectName + `/graph/generated"
	"` + t.projectName + `/internal/interfaces"
	"` + t.projectName + `/pkg/config"
)

// Module provides the Fiber server dependency via fx with automatic lifecycle management.
var Module = fx.Module("server",
	fx.Provide(NewServer),
	fx.Invoke(registerHooks),
)

// NewServer creates and configures a new Fiber application with GraphQL support.
func NewServer(log zerolog.Logger, db *gorm.DB, userRepo interfaces.UserRepository) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "` + t.projectName + `",
		// Increase buffer sizes for GraphQL queries
		ReadBufferSize:  32768, // 32KB
		WriteBufferSize: 32768,
		// Request timeout
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Add recovery middleware to prevent crashes
	app.Use(recover.New())

	// Add CORS middleware for frontend integration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Add rate limiting to prevent abuse
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			log.Warn().Str("ip", c.IP()).Msg("Rate limit exceeded")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	// Add request logging middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))

	// Ignore common browser requests
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	// Create GraphQL resolver with dependencies
	resolver := graph.NewResolver(userRepo)

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// GraphQL Playground at root
	app.Get("/", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(w, r)
	}))

	// GraphQL query endpoint
	app.All("/query", adaptor.HTTPHandler(srv))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Info().Msg("Fiber server initialized with GraphQL support")

	return app
}

// registerHooks registers fx lifecycle hooks for server startup and graceful shutdown.
func registerHooks(lifecycle fx.Lifecycle, app *fiber.App, log zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := config.GetEnv("APP_PORT", "8080")
			log.Info().Str("port", port).Msg("Starting Fiber server with GraphQL")
			log.Info().Str("playground", "http://localhost:"+port+"/").Msg("GraphQL Playground available")
			log.Info().Str("endpoint", "http://localhost:"+port+"/query").Msg("GraphQL endpoint")

			// Start server in background goroutine
			go func() {
				if err := app.Listen(":" + port); err != nil {
					log.Error().Err(err).Msg("Server stopped unexpectedly")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Shutting down Fiber server gracefully")
			return app.ShutdownWithContext(ctx)
		},
	})
}
`
}

// GraphQLDatabaseTemplate returns the internal/infrastructure/database/database.go for GraphQL template.
func (t *ProjectTemplates) GraphQLDatabaseTemplate() string {
	return `// Package database provides PostgreSQL database connectivity and management.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"` + t.projectName + `/internal/models"
	"` + t.projectName + `/pkg/config"
)

// Module provides the database dependency via fx with automatic lifecycle management.
var Module = fx.Module("database",
	fx.Provide(NewDatabase),
	fx.Provide(NewUserRepository),
	fx.Invoke(registerHooks),
)

// NewDatabase creates a new GORM database connection configured from environment variables.
func NewDatabase(logger zerolog.Logger) (*gorm.DB, error) {
	// Build DSN from environment variables
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "5432"),
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", "postgres"),
		config.GetEnv("DB_NAME", "` + t.projectName + `"),
		config.GetEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info().Msg("Successfully connected to database")

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// AutoMigrate database schemas
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	logger.Info().Msg("Database migrations completed successfully")

	return db, nil
}

// registerHooks registers fx lifecycle hooks for graceful database shutdown.
func registerHooks(lifecycle fx.Lifecycle, db *gorm.DB, logger zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Closing database connection")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
}
`
}

// GraphQLUserRepositoryTemplate returns the user repository for GraphQL template.
func (t *ProjectTemplates) GraphQLUserRepositoryTemplate() string {
	return `package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"` + t.projectName + `/internal/interfaces"
	"` + t.projectName + `/internal/models"
)

// UserRepository implements user data persistence using GORM.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user record into the database.
func (r *UserRepository) CreateUser(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// GetUserByEmail retrieves a user by their email address.
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

// FindByID retrieves a user by their unique identifier.
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

// FindAll retrieves all users with pagination support.
func (r *UserRepository) FindAll(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

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

// Update persists changes to an existing user record.
func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

// Delete performs a soft delete on the user.
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// SaveRefreshToken stores a new refresh token (not used in GraphQL template).
func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	return nil // Not implemented for GraphQL template
}

// GetRefreshToken retrieves a refresh token (not used in GraphQL template).
func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return nil, nil // Not implemented for GraphQL template
}

// RevokeRefreshToken marks a refresh token as revoked (not used in GraphQL template).
func (r *UserRepository) RevokeRefreshToken(ctx context.Context, tokenID uint) error {
	return nil // Not implemented for GraphQL template
}

// RotateRefreshToken performs token rotation (not used in GraphQL template).
func (r *UserRepository) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *models.RefreshToken) error {
	return nil // Not implemented for GraphQL template
}
`
}

// GraphQLModelsUserTemplate returns the internal/models/user.go for GraphQL template.
func (t *ProjectTemplates) GraphQLModelsUserTemplate() string {
	return `// Package models defines the domain entities used throughout the application.
package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the domain entity for a user account.
type User struct {
	ID           uint           ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
	Email        string         ` + "`gorm:\"uniqueIndex;not null\" json:\"email\"`" + `
	PasswordHash string         ` + "`gorm:\"not null\" json:\"-\"`" + `
	CreatedAt    time.Time      ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
	UpdatedAt    time.Time      ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
	DeletedAt    gorm.DeletedAt ` + "`gorm:\"index\" json:\"deleted_at,omitempty\"`" + `
}

// RefreshToken is not used in GraphQL template but kept for interface compatibility.
type RefreshToken struct {
	ID        uint      ` + "`gorm:\"primaryKey\"`" + `
	UserID    uint      ` + "`gorm:\"not null;index\"`" + `
	Token     string    ` + "`gorm:\"uniqueIndex;not null\"`" + `
	ExpiresAt time.Time ` + "`gorm:\"not null\"`" + `
	Revoked   bool      ` + "`gorm:\"not null;default:false\"`" + `
	CreatedAt time.Time ` + "`gorm:\"autoCreateTime\"`" + `
	UpdatedAt time.Time ` + "`gorm:\"autoUpdateTime\"`" + `
}
`
}

// GraphQLInterfacesTemplate returns the internal/interfaces/user_repository.go for GraphQL template.
func (t *ProjectTemplates) GraphQLInterfacesTemplate() string {
	return `// Package interfaces defines the ports (abstractions) for the hexagonal architecture.
package interfaces

import (
	"context"
	"time"

	"` + t.projectName + `/internal/models"
)

// UserRepository defines the interface for user data persistence operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindAll(ctx context.Context, page, limit int) ([]*models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	// Token methods kept for interface compatibility
	SaveRefreshToken(ctx context.Context, UserID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
	RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *models.RefreshToken) error
}
`
}

// GraphQLEnvTemplate returns the .env.example file content for GraphQL template.
func (t *ProjectTemplates) GraphQLEnvTemplate() string {
	return `# Application Configuration
APP_NAME=` + t.projectName + `
APP_ENV=development
APP_PORT=8080

# CORS Configuration (comma-separated list of allowed origins)
CORS_ORIGINS=http://localhost:3000,http://localhost:5173

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=` + t.projectName + `
DB_SSLMODE=disable
`
}

// GraphQLReadmeTemplate returns the README.md file content for GraphQL template.
func (t *ProjectTemplates) GraphQLReadmeTemplate() string {
	return `# ` + t.projectName + `

Application backend Go avec API GraphQL, générée avec create-go-starter.

## Fonctionnalités

- **API GraphQL** avec gqlgen - Génération de code type-safe
- **GraphQL Playground** - Interface interactive à la racine
- **Base de données** - GORM avec PostgreSQL et migrations automatiques
- **Injection de dépendances** - uber-go/fx pour architecture modulaire
- **Docker** - Build multi-stage optimisé
- **Logging structuré** - rs/zerolog pour logs professionnels
- **Architecture hexagonale** - Séparation claire des responsabilités

## Prérequis

- **Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **PostgreSQL** - Base de données (peut être lancée via Docker)
- **Docker** (optionnel) - Pour containerisation

## Installation rapide

### 1. Installer les dépendances

` + "```bash" + `
go mod tidy
` + "```" + `

### 2. Générer le code GraphQL

` + "```bash" + `
go generate ./...
` + "```" + `

### 3. Lancer PostgreSQL

` + "```bash" + `
docker run -d \
  --name postgres \
  -e POSTGRES_DB=` + t.projectName + ` \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
` + "```" + `

### 4. Lancer l'application

` + "```bash" + `
make run
` + "```" + `

## Endpoints

- **GraphQL Playground**: ` + "`http://localhost:8080/`" + `
- **GraphQL Endpoint**: ` + "`http://localhost:8080/query`" + `
- **Health Check**: ` + "`http://localhost:8080/health`" + `

## Exemples de requêtes GraphQL

### Créer un utilisateur

` + "```graphql" + `
mutation {
  createUser(input: { email: "test@example.com", password: "password123" }) {
    id
    email
    createdAt
  }
}
` + "```" + `

### Lister les utilisateurs

` + "```graphql" + `
query {
  users(page: 1, limit: 10) {
    users {
      id
      email
      createdAt
    }
    pageInfo {
      page
      limit
      total
      hasNextPage
    }
  }
}
` + "```" + `

### Obtenir un utilisateur par ID

` + "```graphql" + `
query {
  user(id: "1") {
    id
    email
    createdAt
    updatedAt
  }
}
` + "```" + `

### Mettre à jour un utilisateur

` + "```graphql" + `
mutation {
  updateUser(id: "1", input: { email: "newemail@example.com" }) {
    id
    email
  }
}
` + "```" + `

### Supprimer un utilisateur

` + "```graphql" + `
mutation {
  deleteUser(id: "1")
}
` + "```" + `

## Architecture

` + "```" + `
` + t.projectName + `/
├── cmd/                      # Point d'entrée
│   └── main.go               # Bootstrap avec fx
├── graph/                    # GraphQL
│   ├── generated/            # Code généré par gqlgen
│   ├── model/                # Modèles GraphQL
│   ├── resolver.go           # Resolver principal
│   ├── schema.graphqls       # Schéma GraphQL
│   └── schema.resolvers.go   # Implémentation des resolvers
├── internal/
│   ├── infrastructure/       # Infrastructure
│   │   ├── database/         # Configuration DB + Repository
│   │   └── server/           # Configuration Fiber + GraphQL
│   ├── interfaces/           # Ports (interfaces)
│   └── models/               # Entités domaine
├── pkg/                      # Packages réutilisables
│   ├── config/               # Configuration
│   └── logger/               # Logger
├── gqlgen.yml                # Configuration gqlgen
├── .env                      # Configuration
├── Dockerfile                # Build Docker
└── Makefile                  # Commandes
` + "```" + `

## Développement

### Commandes Make

| Commande | Description |
|----------|-------------|
| ` + "`make help`" + ` | Afficher l'aide |
| ` + "`make run`" + ` | Lancer l'application |
| ` + "`make build`" + ` | Compiler le binaire |
| ` + "`make test`" + ` | Tests avec race detector |
| ` + "`make generate`" + ` | Générer code GraphQL |
| ` + "`make docker-build`" + ` | Build image Docker |

### Régénérer le code GraphQL

Après modification du schéma (` + "`graph/schema.graphqls`" + `):

` + "```bash" + `
go generate ./...
# ou
go run github.com/99designs/gqlgen generate
` + "```" + `

## Stack technique

| Composant | Bibliothèque | Description |
|-----------|-------------|-------------|
| GraphQL | [gqlgen](https://gqlgen.com/) | Génération GraphQL type-safe |
| Web Framework | [Fiber](https://gofiber.io/) v2 | Framework HTTP rapide |
| Adaptor | [gofiber/adaptor](https://github.com/gofiber/adaptor) | Bridge net/http vers Fiber |
| ORM | [GORM](https://gorm.io/) | ORM avec PostgreSQL |
| DI | [fx](https://uber-go.github.io/fx/) | Dependency injection |
| Logging | [zerolog](https://github.com/rs/zerolog) | Logger structuré |

## Licence

MIT

---

**Généré avec [create-go-starter](https://github.com/tky0065/go-starter-kit)** 🚀
`
}

// GraphQLMakefileTemplate returns the Makefile for GraphQL template.
func (t *ProjectTemplates) GraphQLMakefileTemplate() string {
	return `.PHONY: help build run test clean generate lint docker-build docker-run

# Binary name
BINARY_NAME=` + t.projectName + `

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd
	@echo "Build complete: $(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@go run ./cmd

generate: ## Generate GraphQL code
	@echo "Generating GraphQL code..."
	@go generate ./...
	@echo "Generation complete"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

test: ## Run tests with race detection
	@echo "Running tests..."
	@go test -v -race ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

docker-build: ## Build docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-run: ## Run docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME):latest
`
}

// GraphQLSetupScriptTemplate returns the setup.sh for GraphQL template.
func (t *ProjectTemplates) GraphQLSetupScriptTemplate() string {
	return `#!/bin/bash

# setup.sh - Automated setup script for ` + t.projectName + ` (GraphQL template)
# This script configures your development environment

set -e  # Exit on error

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Helper functions
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_step() {
    echo -e "\n${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Welcome message
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Configuration automatique de ` + t.projectName + ` (GraphQL)${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}\n"

# ============================================================================
# STEP 1: Check Prerequisites
# ============================================================================
print_step "Étape 1/5: Vérification des prérequis"

MISSING_DEPS=0

# Check Go
if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_success "Go est installé: $GO_VERSION"
else
    print_error "Go n'est pas installé. Installez Go 1.25+ depuis https://golang.org/dl/"
    MISSING_DEPS=1
fi

# Check Docker (optional but recommended)
if command_exists docker; then
    print_success "Docker est installé"
    DOCKER_AVAILABLE=1
else
    print_info "Docker n'est pas installé (optionnel). PostgreSQL devra être installé localement."
    DOCKER_AVAILABLE=0
fi

if [ $MISSING_DEPS -eq 1 ]; then
    print_error "Des dépendances obligatoires sont manquantes. Installez-les et relancez ce script."
    exit 1
fi

# ============================================================================
# STEP 2: Install Go Dependencies
# ============================================================================
print_step "Étape 2/5: Installation des dépendances Go"

print_info "Exécution de 'go mod tidy'..."
if go mod tidy; then
    print_success "Dépendances Go installées avec succès"
else
    print_error "Échec de l'installation des dépendances Go"
    exit 1
fi

# ============================================================================
# STEP 3: Generate GraphQL Code
# ============================================================================
print_step "Étape 3/5: Génération du code GraphQL"

print_info "Exécution de 'go generate ./...'..."
if go generate ./... 2>/dev/null; then
    print_success "Code GraphQL généré avec succès"
else
    print_info "Génération ignorée (exécutez 'go generate ./...' manuellement après setup)"
fi

# ============================================================================
# STEP 4: Configure Environment
# ============================================================================
print_step "Étape 4/5: Configuration de l'environnement"

if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        print_success "Fichier .env créé depuis .env.example"
    else
        print_error ".env.example introuvable"
        exit 1
    fi
else
    print_info "Fichier .env existe déjà"
fi

# ============================================================================
# STEP 5: PostgreSQL Setup
# ============================================================================
print_step "Étape 5/5: Configuration de PostgreSQL"

if [ $DOCKER_AVAILABLE -eq 1 ]; then
    echo -n "Voulez-vous démarrer PostgreSQL avec Docker? (Y/n): "
    read -r USE_DOCKER
    if [[ ! $USE_DOCKER =~ ^[Nn]$ ]]; then
        if docker ps -a --format '{{.Names}}' | grep -q "^postgres$"; then
            print_info "Conteneur PostgreSQL existe déjà"
            if docker ps --format '{{.Names}}' | grep -q "^postgres$"; then
                print_success "PostgreSQL est déjà en cours d'exécution"
            else
                docker start postgres
                print_success "PostgreSQL démarré"
            fi
        else
            print_info "Création du conteneur PostgreSQL..."
            docker run -d \
                --name postgres \
                -e POSTGRES_DB=` + t.projectName + ` \
                -e POSTGRES_PASSWORD=postgres \
                -p 5432:5432 \
                postgres:16-alpine
            sleep 5
            print_success "PostgreSQL démarré avec Docker"
        fi
    fi
fi

# ============================================================================
# Summary
# ============================================================================
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✅ Configuration terminée avec succès!${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}\n"

print_info "Prochaines étapes:"
echo "  1. Générer le code GraphQL:  go generate ./..."
echo "  2. Lancer l'application:     make run"
echo "  3. GraphQL Playground:       http://localhost:8080/"
echo "  4. Endpoint GraphQL:         http://localhost:8080/query"
echo ""
print_success "Bon développement! 🚀"
`
}

// GraphQLDockerComposeTemplate returns the docker-compose.yml for GraphQL template.
func (t *ProjectTemplates) GraphQLDockerComposeTemplate() string {
	return `version: '3.8'

services:
  # PostgreSQL Database
  db:
    image: postgres:16-alpine
    container_name: ` + t.projectName + `_db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: ` + t.projectName + `
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - ` + t.projectName + `_network

  # Application API
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ` + t.projectName + `_api
    environment:
      APP_NAME: ` + t.projectName + `
      APP_ENV: development
      APP_PORT: 8080
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: ` + t.projectName + `
      DB_SSLMODE: disable
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
    networks:
      - ` + t.projectName + `_network
    command: /app/` + t.projectName + `

volumes:
  postgres_data:

networks:
  ` + t.projectName + `_network:
    driver: bridge
`
}

// GraphQLGenerateGoTemplate returns the generate.go file for go generate support.
func (t *ProjectTemplates) GraphQLGenerateGoTemplate() string {
	return `//go:generate go run github.com/99designs/gqlgen generate

package graph
`
}

// GraphQLResolverTestTemplate returns the schema.resolvers_test.go file for testing resolvers.
func (t *ProjectTemplates) GraphQLResolverTestTemplate() string {
	return `package graph

import (
	"context"
	"testing"

	"` + t.projectName + `/graph/model"
	"` + t.projectName + `/internal/models"
)

// MockUserRepository is a mock implementation of UserRepository for testing.
type MockUserRepository struct {
	users map[uint]*models.User
	nextID uint
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*models.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, int64(len(users)), nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	return nil
}

func (m *MockUserRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return nil, nil
}

func (m *MockUserRepository) RevokeRefreshToken(ctx context.Context, tokenID uint) error {
	return nil
}

func (m *MockUserRepository) RotateRefreshToken(ctx context.Context, oldTokenID uint, newToken *models.RefreshToken) error {
	return nil
}

func TestCreateUser(t *testing.T) {
	repo := NewMockUserRepository()
	resolver := NewResolver(repo)
	mutationResolver := &mutationResolver{resolver}

	ctx := context.Background()

	t.Run("Valid user creation", func(t *testing.T) {
		input := model.NewUser{
			Email:    "test@example.com",
			Password: "password123",
		}

		user, err := mutationResolver.CreateUser(ctx, input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if user.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
		}
		if user.ID == 0 {
			t.Error("Expected user ID to be set")
		}
	})

	t.Run("Invalid email format", func(t *testing.T) {
		input := model.NewUser{
			Email:    "invalid-email",
			Password: "password123",
		}

		_, err := mutationResolver.CreateUser(ctx, input)
		if err == nil {
			t.Error("Expected error for invalid email format")
		}
	})

	t.Run("Password too short", func(t *testing.T) {
		input := model.NewUser{
			Email:    "test2@example.com",
			Password: "short",
		}

		_, err := mutationResolver.CreateUser(ctx, input)
		if err == nil {
			t.Error("Expected error for short password")
		}
	})

	t.Run("Duplicate email", func(t *testing.T) {
		input := model.NewUser{
			Email:    "duplicate@example.com",
			Password: "password123",
		}

		// Create first user
		_, err := mutationResolver.CreateUser(ctx, input)
		if err != nil {
			t.Fatalf("Expected no error on first creation, got %v", err)
		}

		// Try to create duplicate
		_, err = mutationResolver.CreateUser(ctx, input)
		if err == nil {
			t.Error("Expected error for duplicate email")
		}
	})
}

func TestUserQuery(t *testing.T) {
	repo := NewMockUserRepository()
	resolver := NewResolver(repo)
	queryResolver := &queryResolver{resolver}

	ctx := context.Background()

	// Create a test user
	testUser := &models.User{Email: "query@example.com", PasswordHash: "hashed"}
	repo.CreateUser(ctx, testUser)

	t.Run("Find existing user", func(t *testing.T) {
		user, err := queryResolver.User(ctx, "1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if user.Email != "query@example.com" {
			t.Errorf("Expected email 'query@example.com', got '%s'", user.Email)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		_, err := queryResolver.User(ctx, "999")
		if err == nil {
			t.Error("Expected error for non-existent user")
		}
	})

	t.Run("Invalid user ID", func(t *testing.T) {
		_, err := queryResolver.User(ctx, "invalid")
		if err == nil {
			t.Error("Expected error for invalid user ID")
		}
	})
}

func TestUsersQuery(t *testing.T) {
	repo := NewMockUserRepository()
	resolver := NewResolver(repo)
	queryResolver := &queryResolver{resolver}

	ctx := context.Background()

	// Create test users
	for i := 1; i <= 5; i++ {
		user := &models.User{Email: fmt.Sprintf("user%d@example.com", i), PasswordHash: "hashed"}
		repo.CreateUser(ctx, user)
	}

	t.Run("Default pagination", func(t *testing.T) {
		result, err := queryResolver.Users(ctx, nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(result.Users) != 5 {
			t.Errorf("Expected 5 users, got %d", len(result.Users))
		}
	})

	t.Run("Limit exceeds maximum", func(t *testing.T) {
		limit := 101
		_, err := queryResolver.Users(ctx, nil, &limit)
		if err == nil {
			t.Error("Expected error when limit exceeds 100")
		}
	})
}
`
}

// GraphQLDocsReadmeTemplate returns the docs/README.md for GraphQL template.
func (t *ProjectTemplates) GraphQLDocsReadmeTemplate() string {
	return `# Documentation ` + t.projectName + `

Documentation pour le projet ` + t.projectName + ` (template GraphQL).

## Table des matières

1. [Démarrage rapide](./quick-start.md)

## Aide rapide

- **Lancer le projet**: ` + "`make run`" + `
- **GraphQL Playground**: ` + "`http://localhost:8080/`" + `
- **Endpoint GraphQL**: ` + "`http://localhost:8080/query`" + `
- **Health Check**: ` + "`http://localhost:8080/health`" + `

## Ressources

- [gqlgen Documentation](https://gqlgen.com/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
`
}

// GraphQLQuickStartTemplate returns the docs/quick-start.md for GraphQL template.
func (t *ProjectTemplates) GraphQLQuickStartTemplate() string {
	return `# Démarrage rapide

Guide pour lancer ` + t.projectName + ` (GraphQL) en 5 minutes.

## Prérequis

- Go 1.25+
- PostgreSQL (ou Docker)

## Installation

### 1. Installer les dépendances

` + "```bash" + `
go mod tidy
` + "```" + `

### 2. Générer le code GraphQL

` + "```bash" + `
go generate ./...
` + "```" + `

### 3. Configurer la base de données

**Docker (Recommandé)**

` + "```bash" + `
docker run -d \
  --name postgres \
  -e POSTGRES_DB=` + t.projectName + ` \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
` + "```" + `

### 4. Lancer l'application

` + "```bash" + `
make run
` + "```" + `

## Tester l'API

Ouvrez le GraphQL Playground: ` + "`http://localhost:8080/`" + `

### Créer un utilisateur

` + "```graphql" + `
mutation {
  createUser(input: { email: "test@example.com", password: "password123" }) {
    id
    email
    createdAt
  }
}
` + "```" + `

### Lister les utilisateurs

` + "```graphql" + `
query {
  users {
    users {
      id
      email
    }
    pageInfo {
      total
    }
  }
}
` + "```" + `

## Développement

### Modifier le schéma GraphQL

1. Éditez ` + "`graph/schema.graphqls`" + `
2. Régénérez le code: ` + "`go generate ./...`" + `
3. Implémentez les nouveaux resolvers dans ` + "`graph/schema.resolvers.go`" + `

Bon développement! 🚀
`
}
