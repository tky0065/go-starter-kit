package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestNewProjectTemplates(t *testing.T) {
	projectName := "my-awesome-project"
	templates := NewProjectTemplates(projectName)

	if templates.projectName != projectName {
		t.Errorf("Expected project name %s, got %s", projectName, templates.projectName)
	}
}

func TestGoModTemplate(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantModule  string
		wantGoVer   string
	}{
		{
			name:        "simple project name",
			projectName: "my-project",
			wantModule:  "module my-project",
			wantGoVer:   "go 1.25.5",
		},
		{
			name:        "project with hyphens",
			projectName: "my-awesome-api",
			wantModule:  "module my-awesome-api",
			wantGoVer:   "go 1.25.5",
		},
		{
			name:        "project with underscores",
			projectName: "my_cool_app",
			wantModule:  "module my_cool_app",
			wantGoVer:   "go 1.25.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := NewProjectTemplates(tt.projectName)
			content := templates.GoModTemplate()

			if !strings.Contains(content, tt.wantModule) {
				t.Errorf("GoModTemplate() should contain '%s', got:\n%s", tt.wantModule, content)
			}

			if !strings.Contains(content, tt.wantGoVer) {
				t.Errorf("GoModTemplate() should contain '%s', got:\n%s", tt.wantGoVer, content)
			}

			// Check required dependencies
			requiredDeps := []string{
				"github.com/gofiber/fiber/v2 v2.52.10",
				"gorm.io/gorm v1.31.1",
				"go.uber.org/fx v1.24.0",
			}

			for _, dep := range requiredDeps {
				if !strings.Contains(content, dep) {
					t.Errorf("GoModTemplate() should contain dependency '%s', got:\n%s", dep, content)
				}
			}
		})
	}
}

func TestMainGoTemplate(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
	}{
		{
			name:        "basic project",
			projectName: "test-project",
		},
		{
			name:        "project with special chars",
			projectName: "my-api_v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := NewProjectTemplates(tt.projectName)
			content := templates.MainGoTemplate()

			// Check project name is included
			if !strings.Contains(content, tt.projectName) {
				t.Errorf("MainGoTemplate() should contain project name '%s', got:\n%s", tt.projectName, content)
			}

			// Check it's a valid Go main package
			if !strings.Contains(content, "package main") {
				t.Error("MainGoTemplate() should contain 'package main'")
			}

			// Check it has a main function
			if !strings.Contains(content, "func main()") {
				t.Error("MainGoTemplate() should contain 'func main()'")
			}
		})
	}
}

func TestDockerfileTemplate(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantBinary  string
	}{
		{
			name:        "simple binary name",
			projectName: "my-app",
			wantBinary:  "my-app",
		},
		{
			name:        "binary with hyphens",
			projectName: "my-awesome-app",
			wantBinary:  "my-awesome-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := NewProjectTemplates(tt.projectName)
			content := templates.DockerfileTemplate()

			// Check binary name appears in build command
			buildCmd := "go build -o " + tt.wantBinary
			if !strings.Contains(content, buildCmd) {
				t.Errorf("DockerfileTemplate() should contain build command '%s', got:\n%s", buildCmd, content)
			}

			// Check binary name in CMD
			runCmd := `CMD ["./` + tt.wantBinary + `"]`
			if !strings.Contains(content, runCmd) {
				t.Errorf("DockerfileTemplate() should contain run command '%s', got:\n%s", runCmd, content)
			}

			// Check multi-stage build
			if !strings.Contains(content, "FROM golang:1.25-alpine AS builder") {
				t.Error("DockerfileTemplate() should use multi-stage build")
			}

			// Check port exposure
			if !strings.Contains(content, "EXPOSE 8080") {
				t.Error("DockerfileTemplate() should expose port 8080")
			}

			// Verify no go.sum reference (not generated initially)
			if strings.Contains(content, "go.sum") {
				t.Error("DockerfileTemplate() should not reference go.sum as it's not generated")
			}

			// Check build path is ./cmd not ./cmd/main.go
			if !strings.Contains(content, "./cmd") {
				t.Error("DockerfileTemplate() should build from ./cmd")
			}
		})
	}
}

func TestMakefileTemplate(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
	}{
		{
			name:        "basic makefile",
			projectName: "my-project",
		},
		{
			name:        "makefile with special chars",
			projectName: "my-api_service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := NewProjectTemplates(tt.projectName)
			content := templates.MakefileTemplate()

			// Check BINARY_NAME variable
			binaryVar := "BINARY_NAME=" + tt.projectName
			if !strings.Contains(content, binaryVar) {
				t.Errorf("MakefileTemplate() should contain '%s', got:\n%s", binaryVar, content)
			}

			// Check essential targets
			essentialTargets := []string{
				".PHONY:",
				"help:",
				"build:",
				"run:",
				"test:",
				"clean:",
			}

			for _, target := range essentialTargets {
				if !strings.Contains(content, target) {
					t.Errorf("MakefileTemplate() should contain target '%s', got:\n%s", target, content)
				}
			}
		})
	}
}

func TestEnvTemplate(t *testing.T) {
	projectName := "test-app"
	templates := NewProjectTemplates(projectName)
	content := templates.EnvTemplate()

	// Check APP_NAME uses project name
	appName := "APP_NAME=" + projectName
	if !strings.Contains(content, appName) {
		t.Errorf("EnvTemplate() should contain '%s', got:\n%s", appName, content)
	}

	// Check DB_NAME uses project name
	dbName := "DB_NAME=" + projectName
	if !strings.Contains(content, dbName) {
		t.Errorf("EnvTemplate() should contain '%s', got:\n%s", dbName, content)
	}

	// Check essential env vars
	essentialVars := []string{
		"APP_ENV=",
		"APP_PORT=",
		"DB_HOST=",
		"JWT_SECRET=",
	}

	for _, envVar := range essentialVars {
		if !strings.Contains(content, envVar) {
			t.Errorf("EnvTemplate() should contain '%s', got:\n%s", envVar, content)
		}
	}
}

func TestGitignoreTemplate(t *testing.T) {
	projectName := "my-app"
	templates := NewProjectTemplates(projectName)
	content := templates.GitignoreTemplate()

	// Check project binary is ignored
	if !strings.Contains(content, projectName) {
		t.Errorf("GitignoreTemplate() should ignore binary '%s', got:\n%s", projectName, content)
	}

	// Check essential ignores
	essentialIgnores := []string{
		"*.exe",
		".env",
		".vscode/",
		".idea/",
		".DS_Store",
		"vendor/",
	}

	for _, ignore := range essentialIgnores {
		if !strings.Contains(content, ignore) {
			t.Errorf("GitignoreTemplate() should contain '%s', got:\n%s", ignore, content)
		}
	}
}

func TestReadmeTemplate(t *testing.T) {
	projectName := "awesome-project"
	templates := NewProjectTemplates(projectName)
	content := templates.ReadmeTemplate()

	// Check title uses project name
	title := "# " + projectName
	if !strings.Contains(content, title) {
		t.Errorf("ReadmeTemplate() should contain title '%s', got:\n%s", title, content)
	}

	// Check project structure includes project name
	if !strings.Contains(content, projectName+"/") {
		t.Errorf("ReadmeTemplate() should contain project structure with '%s/', got:\n%s", projectName, content)
	}

	// Check essential sections (README is in French)
	essentialSections := []string{
		"## Architecture",
		"## Prérequis",
		"## Installation rapide",
		"## Développement",
		"## API Endpoints",
	}

	for _, section := range essentialSections {
		if !strings.Contains(content, section) {
			t.Errorf("ReadmeTemplate() should contain section '%s', got:\n%s", section, content)
		}
	}
}

func TestDockerComposeTemplate(t *testing.T) {
	projectName := "my-api-project"
	templates := NewProjectTemplates(projectName)
	content := templates.DockerComposeTemplate()

	// Check version
	if !strings.Contains(content, "version:") {
		t.Error("DockerComposeTemplate() should contain version declaration")
	}

	// Check services
	requiredServices := []string{"db:", "api:"}
	for _, service := range requiredServices {
		if !strings.Contains(content, service) {
			t.Errorf("DockerComposeTemplate() should contain service '%s'", service)
		}
	}

	// Check PostgreSQL configuration
	if !strings.Contains(content, "postgres") {
		t.Error("DockerComposeTemplate() should use postgres image")
	}

	// Check project name in container names
	expectedDBContainer := projectName + "_db"
	if !strings.Contains(content, expectedDBContainer) {
		t.Errorf("DockerComposeTemplate() should contain db container name '%s'", expectedDBContainer)
	}

	expectedAPIContainer := projectName + "_api"
	if !strings.Contains(content, expectedAPIContainer) {
		t.Errorf("DockerComposeTemplate() should contain api container name '%s'", expectedAPIContainer)
	}

	// Check database name uses project name
	expectedDBName := "POSTGRES_DB: " + projectName
	if !strings.Contains(content, expectedDBName) {
		t.Errorf("DockerComposeTemplate() should set POSTGRES_DB to '%s'", projectName)
	}

	// Check healthcheck for database
	if !strings.Contains(content, "healthcheck:") {
		t.Error("DockerComposeTemplate() should include healthcheck for database")
	}

	// Check depends_on with health condition
	if !strings.Contains(content, "depends_on:") {
		t.Error("DockerComposeTemplate() should have depends_on for api service")
	}

	if !strings.Contains(content, "condition: service_healthy") {
		t.Error("DockerComposeTemplate() should wait for db to be healthy")
	}

	// Check ports exposed
	if !strings.Contains(content, "5432:5432") {
		t.Error("DockerComposeTemplate() should expose PostgreSQL port 5432")
	}

	if !strings.Contains(content, "8080:8080") {
		t.Error("DockerComposeTemplate() should expose API port 8080")
	}

	// Check volumes
	if !strings.Contains(content, "postgres_data:") {
		t.Error("DockerComposeTemplate() should define postgres_data volume")
	}

	// Check network
	expectedNetwork := projectName + "_network"
	if !strings.Contains(content, expectedNetwork) {
		t.Errorf("DockerComposeTemplate() should define network '%s'", expectedNetwork)
	}

	// Check environment variables for API
	envVars := []string{"APP_NAME:", "DB_HOST:", "DB_USER:", "JWT_SECRET:"}
	for _, envVar := range envVars {
		if !strings.Contains(content, envVar) {
			t.Errorf("DockerComposeTemplate() should contain environment variable '%s'", envVar)
		}
	}
}

// Story 1.4: Tests for Fiber, fx, GORM, and zerolog templates

func TestConfigTemplate(t *testing.T) {
	projectName := "test-app"
	templates := NewProjectTemplates(projectName)
	content := templates.ConfigTemplate()

	// Check package declaration
	if !strings.Contains(content, "package config") {
		t.Error("ConfigTemplate() should contain 'package config'")
	}

	// Check GetEnv function
	if !strings.Contains(content, "func GetEnv") {
		t.Error("ConfigTemplate() should have GetEnv function")
	}

	// Check os import
	if !strings.Contains(content, "import \"os\"") {
		t.Error("ConfigTemplate() should import os")
	}

	// Check function signature
	if !strings.Contains(content, "GetEnv(key, defaultValue string) string") {
		t.Error("ConfigTemplate() GetEnv should have correct signature")
	}
}

func TestLoggerTemplate(t *testing.T) {
	projectName := "test-app"
	templates := NewProjectTemplates(projectName)
	content := templates.LoggerTemplate()

	// Check package declaration
	if !strings.Contains(content, "package logger") {
		t.Error("LoggerTemplate() should contain 'package logger'")
	}

	// Check zerolog import
	if !strings.Contains(content, "github.com/rs/zerolog") {
		t.Error("LoggerTemplate() should import zerolog")
	}

	// Check fx integration
	if !strings.Contains(content, "go.uber.org/fx") {
		t.Error("LoggerTemplate() should import fx")
	}

	// Check Module export
	if !strings.Contains(content, "var Module = fx.Module") {
		t.Error("LoggerTemplate() should export fx Module")
	}

	// Check NewLogger function
	if !strings.Contains(content, "func NewLogger") {
		t.Error("LoggerTemplate() should have NewLogger function")
	}
}

func TestDatabaseTemplate(t *testing.T) {
	projectName := "test-app"
	templates := NewProjectTemplates(projectName)
	content := templates.DatabaseTemplate()

	// Check package declaration
	if !strings.Contains(content, "package database") {
		t.Error("DatabaseTemplate() should contain 'package database'")
	}

	// Check GORM imports
	if !strings.Contains(content, "gorm.io/gorm") {
		t.Error("DatabaseTemplate() should import gorm")
	}

	if !strings.Contains(content, "gorm.io/driver/postgres") {
		t.Error("DatabaseTemplate() should import postgres driver")
	}

	// Check config package import
	if !strings.Contains(content, "pkg/config") {
		t.Error("DatabaseTemplate() should import config package")
	}

	// Check uses config.GetEnv instead of local getEnv
	if !strings.Contains(content, "config.GetEnv") {
		t.Error("DatabaseTemplate() should use config.GetEnv")
	}

	// Check connection pool configuration
	if !strings.Contains(content, "SetMaxOpenConns") {
		t.Error("DatabaseTemplate() should configure connection pool with SetMaxOpenConns")
	}

	// Check fx integration
	if !strings.Contains(content, "go.uber.org/fx") {
		t.Error("DatabaseTemplate() should import fx")
	}

	// Check Module export
	if !strings.Contains(content, "var Module = fx.Module") {
		t.Error("DatabaseTemplate() should export fx Module")
	}

	// Check NewDatabase function
	if !strings.Contains(content, "func NewDatabase") {
		t.Error("DatabaseTemplate() should have NewDatabase function")
	}

	// Check AutoMigrate mention
	if !strings.Contains(content, "AutoMigrate") {
		t.Error("DatabaseTemplate() should include AutoMigrate")
	}

	// Check graceful shutdown
	if !strings.Contains(content, "OnStop") {
		t.Error("DatabaseTemplate() should implement OnStop hook for graceful shutdown")
	}
}

func TestServerTemplate(t *testing.T) {
	projectName := "test-app"
	templates := NewProjectTemplates(projectName)
	content := templates.ServerTemplate()

	// Check package declaration
	if !strings.Contains(content, "package server") {
		t.Error("ServerTemplate() should contain 'package server'")
	}

	// Check Fiber import
	if !strings.Contains(content, "github.com/gofiber/fiber/v2") {
		t.Error("ServerTemplate() should import fiber")
	}

	// Check fx integration
	if !strings.Contains(content, "go.uber.org/fx") {
		t.Error("ServerTemplate() should import fx")
	}

	// Check Module export
	if !strings.Contains(content, "var Module = fx.Module") {
		t.Error("ServerTemplate() should export fx Module")
	}

	// Check NewServer function
	if !strings.Contains(content, "func NewServer") {
		t.Error("ServerTemplate() should have NewServer function")
	}

	// Check OnStart hook
	if !strings.Contains(content, "OnStart") {
		t.Error("ServerTemplate() should implement OnStart hook")
	}

	// Check OnStop hook for graceful shutdown
	if !strings.Contains(content, "OnStop") {
		t.Error("ServerTemplate() should implement OnStop hook for graceful shutdown")
	}

	// Check Shutdown method
	if !strings.Contains(content, "Shutdown") {
		t.Error("ServerTemplate() should call Shutdown for graceful shutdown")
	}

	// Check ShutdownWithContext (respects context timeout)
	if !strings.Contains(content, "ShutdownWithContext") {
		t.Error("ServerTemplate() should use ShutdownWithContext to respect context timeout")
	}

	// Check config package import
	if !strings.Contains(content, "pkg/config") {
		t.Error("ServerTemplate() should import config package")
	}

	// Check uses config.GetEnv
	if !strings.Contains(content, "config.GetEnv") {
		t.Error("ServerTemplate() should use config.GetEnv")
	}
}

func TestHealthHandlerTemplate(t *testing.T) {
	projectName := "test-app"
	templates := NewProjectTemplates(projectName)
	content := templates.HealthHandlerTemplate()

	// Check package declaration
	if !strings.Contains(content, "package http") {
		t.Error("HealthHandlerTemplate() should contain 'package http'")
	}

	// Check Fiber import
	if !strings.Contains(content, "github.com/gofiber/fiber/v2") {
		t.Error("HealthHandlerTemplate() should import fiber")
	}

	// Check health route /health
	if !strings.Contains(content, "/health") {
		t.Error("HealthHandlerTemplate() should register /health route")
	}

	// Check JSON response with "status": "ok"
	if !strings.Contains(content, `"status"`) && !strings.Contains(content, `"ok"`) {
		t.Error("HealthHandlerTemplate() should return JSON with status ok")
	}

	// Check RegisterRoutes function
	if !strings.Contains(content, "func RegisterRoutes") && !strings.Contains(content, "func RegisterHealthRoutes") {
		t.Error("HealthHandlerTemplate() should have a route registration function")
	}
}

func TestUpdatedMainGoTemplate(t *testing.T) {
	projectName := "test-app"
	templates := NewProjectTemplates(projectName)
	content := templates.UpdatedMainGoTemplate()

	// Check package main
	if !strings.Contains(content, "package main") {
		t.Error("UpdatedMainGoTemplate() should contain 'package main'")
	}

	// Check fx.New usage
	if !strings.Contains(content, "fx.New") {
		t.Error("UpdatedMainGoTemplate() should use fx.New")
	}

	// Check all module imports
	requiredModules := []string{
		"logger.Module",
		"database.Module",
		"auth.Module",
		"user.Module",
		"repository.Module",
		"handlers.Module",
		"server.Module",
	}

	for _, mod := range requiredModules {
		if !strings.Contains(content, mod) {
			t.Errorf("UpdatedMainGoTemplate() should import %s", mod)
		}
	}

	// Check Run() call
	if !strings.Contains(content, ".Run()") {
		t.Error("UpdatedMainGoTemplate() should call Run() on fx.App")
	}

	// Check imports for infrastructure packages
	if !strings.Contains(content, "pkg/logger") {
		t.Error("UpdatedMainGoTemplate() should import logger package")
	}

	if !strings.Contains(content, "internal/infrastructure/database") {
		t.Error("UpdatedMainGoTemplate() should import database package")
	}

	if !strings.Contains(content, "internal/infrastructure/server") {
		t.Error("UpdatedMainGoTemplate() should import server package")
	}
}

// Test User Templates
func TestUserEntityTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.ModelsUserTemplate() // User entity is now in ModelsUserTemplate

	requiredContent := []string{
		"package models",
		"type User struct",
		"ID           uint",
		"Email        string",
		"PasswordHash string",
		"CreatedAt    time.Time",
		"UpdatedAt    time.Time",
		"DeletedAt    gorm.DeletedAt",
		`json:"-"`, // Password should not be in JSON
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("ModelsUserTemplate() should contain '%s'", required)
		}
	}
}

func TestUserRefreshTokenTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.ModelsUserTemplate() // RefreshToken is now in ModelsUserTemplate

	requiredContent := []string{
		"package models",
		"type RefreshToken struct",
		"UserID    uint",
		"Token     string",
		"ExpiresAt time.Time",
		"Revoked   bool",
		"IsExpired() bool",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("ModelsUserTemplate() should contain '%s'", required)
		}
	}
}

func TestUserServiceTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.UserServiceTemplate()

	requiredContent := []string{
		"package user",
		"type Service struct",
		"func NewService(",
		"func NewServiceWithJWT(",
		"func (s *Service) Register(",
		"func (s *Service) Authenticate(",
		"func (s *Service) RefreshToken(",
		"func (s *Service) GetProfile(",
		"func (s *Service) GetAll(",
		"func (s *Service) UpdateUser(",
		"func (s *Service) DeleteUser(",
		"bcrypt.GenerateFromPassword",
		"bcrypt.CompareHashAndPassword",
		projectName + "/internal/domain",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("UserServiceTemplate() should contain '%s'", required)
		}
	}
}

func TestUserRepositoryTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.UserRepositoryTemplate()

	requiredContent := []string{
		"package repository",
		"type UserRepository struct",
		"func NewUserRepository(",
		"func (r *UserRepository) CreateUser(",
		"func (r *UserRepository) GetUserByEmail(",
		"func (r *UserRepository) FindByID(",
		"func (r *UserRepository) FindAll(",
		"func (r *UserRepository) Update(",
		"func (r *UserRepository) Delete(",
		"func (r *UserRepository) SaveRefreshToken(",
		"func (r *UserRepository) GetRefreshToken(",
		"func (r *UserRepository) RevokeRefreshToken(",
		"func (r *UserRepository) RotateRefreshToken(",
		"WithContext(ctx)",
		projectName + "/internal/models",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("UserRepositoryTemplate() should contain '%s'", required)
		}
	}
}

func TestAuthHandlerTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.AuthHandlerTemplate()

	requiredContent := []string{
		"package handlers",
		"type AuthHandler struct",
		"func NewAuthHandler(",
		"func (h *AuthHandler) Register(",
		"func (h *AuthHandler) Login(",
		"func (h *AuthHandler) Refresh(",
		"type RegisterRequest struct",
		"type LoginRequest struct",
		"type RefreshRequest struct",
		`validate:"required,email,max=255"`,
		`validate:"required,min=8,max=72"`,
		"fiber.StatusCreated",
		"validator.New()",
		projectName + "/internal/domain/user",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("AuthHandlerTemplate() should contain '%s'", required)
		}
	}
}

func TestUserHandlerTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.UserHandlerTemplate()

	requiredContent := []string{
		"package handlers",
		"type UserHandler struct",
		"func NewUserHandler(",
		"func (h *UserHandler) GetMe(",
		"func (h *UserHandler) GetAllUsers(",
		"func (h *UserHandler) UpdateUser(",
		"func (h *UserHandler) DeleteUser(",
		"auth.GetUserID",
		"validator.New()",
		projectName + "/pkg/auth",
		// Response format validation
		`"page":  page`,
		`"limit": limit`,
		`"total": total`,
		`"message": "User deleted successfully"`,
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("UserHandlerTemplate() should contain '%s'", required)
		}
	}
}

func TestJWTAuthTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.JWTAuthTemplate()

	requiredContent := []string{
		"package auth",
		"type JWTService struct",
		"func NewJWTService()",
		"func (s *JWTService) GenerateTokens(",
		"func GetUserID(c *fiber.Ctx)",
		"func (s *JWTService) ValidateToken(",
		"jwt.NewWithClaims",
		"jwt.SigningMethodHS256",
		"ErrInvalidToken",
		"ErrMissingUserID",
		projectName + "/pkg/config",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("JWTAuthTemplate() should contain '%s'", required)
		}
	}
}

func TestJWTMiddlewareTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.JWTMiddlewareTemplate()

	requiredContent := []string{
		"package auth",
		"func NewJWTMiddleware()",
		"jwtware.New(",
		"SigningKey:",
		"JWTAlg: jwtware.HS256",
		"ErrorHandler:",
		"fiber.StatusUnauthorized",
		projectName + "/pkg/config",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("JWTMiddlewareTemplate() should contain '%s'", required)
		}
	}
}

func TestUserInterfacesTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.UserInterfacesTemplate()

	requiredContent := []string{
		"package interfaces",
		"type TokenService interface",
		"GenerateTokens(userID uint)",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("UserInterfacesTemplate() should contain '%s'", required)
		}
	}
}

func TestUserRepositoryInterfaceTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.UserRepositoryInterfaceTemplate()

	requiredContent := []string{
		"package interfaces",
		"type UserRepository interface",
		"CreateUser(ctx context.Context,",
		"GetUserByEmail(ctx context.Context,",
		"FindByID(ctx context.Context,",
		"FindAll(ctx context.Context,",
		"Update(ctx context.Context,",
		"Delete(ctx context.Context,",
		"SaveRefreshToken(ctx context.Context,",
		"GetRefreshToken(ctx context.Context,",
		"RevokeRefreshToken(ctx context.Context,",
		"RotateRefreshToken(ctx context.Context,",
		projectName + "/internal/models",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("UserRepositoryInterfaceTemplate() should contain '%s'", required)
		}
	}
}

func TestUserModuleTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.UserModuleTemplate()

	requiredContent := []string{
		"package user",
		"var Module = fx.Module(",
		"fx.Provide(NewServiceWithJWT)",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("UserModuleTemplate() should contain '%s'", required)
		}
	}
}

func TestRepositoryModuleTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.RepositoryModuleTemplate()

	requiredContent := []string{
		"package repository",
		"var Module = fx.Module(",
		"fx.Provide(",
		"NewUserRepository(",
		projectName + "/internal/interfaces",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("RepositoryModuleTemplate() should contain '%s'", required)
		}
	}
}

func TestAuthModuleTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.AuthModuleTemplate()

	requiredContent := []string{
		"package auth",
		"var Module = fx.Module(",
		"fx.Provide(",
		"NewJWTService()",
		"NewJWTMiddleware",
		projectName + "/internal/interfaces",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("AuthModuleTemplate() should contain '%s'", required)
		}
	}
}

func TestHandlerModuleTemplate(t *testing.T) {
	projectName := "my-auth-app"
	templates := NewProjectTemplates(projectName)
	content := templates.HandlerModuleTemplate()

	requiredContent := []string{
		"package handlers",
		"var Module = fx.Module(",
		"fx.Provide(func(s *user.Service) *AuthHandler",
		"fx.Provide(func(s *user.Service) *UserHandler",
		projectName + "/internal/domain/user",
	}

	for _, required := range requiredContent {
		if !strings.Contains(content, required) {
			t.Errorf("HandlerModuleTemplate() should contain '%s'", required)
		}
	}
}

// TestE2EDockerImageSize tests that generated Docker image is under 50MB (AC #1)
// This is an end-to-end test that requires Docker to be installed and running
func TestE2EDockerImageSize(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E Docker test in short mode")
	}

	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping Docker image size test")
	}

	// Create temporary project
	tmpDir := t.TempDir()
	projectName := "test-docker-size"
	projectPath := filepath.Join(tmpDir, projectName)

	// Use createProjectStructure function
	if err := createProjectStructure(projectPath); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate files using templates
	templates := NewProjectTemplates(projectName)

	// Create Dockerfile
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(templates.DockerfileTemplate()), 0644); err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	// Create go.mod
	gomodPath := filepath.Join(projectPath, "go.mod")
	if err := os.WriteFile(gomodPath, []byte(templates.GoModTemplate()), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create cmd directory and main.go
	cmdDir := filepath.Join(projectPath, "cmd")
	if err := os.MkdirAll(cmdDir, defaultDirPerm); err != nil {
		t.Fatalf("Failed to create cmd directory: %v", err)
	}

	mainGoPath := filepath.Join(cmdDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(templates.MainGoTemplate()), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to tidy go module: %v", err)
	}

	// Build Docker image
	imageName := "test-docker-optimization:latest"
	cmd = exec.Command("docker", "build", "-t", imageName, ".")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build Docker image: %v\nOutput: %s", err, string(output))
	}

	// Cleanup image after test
	defer func() {
		exec.Command("docker", "rmi", imageName).Run()
	}()

	// Check image size using docker images --format
	cmd = exec.Command("docker", "images", "--format", "table {{.Repository}}:{{.Tag}}\t{{.Size}}", imageName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get image size: %v", err)
	}

	// Parse size (format is like "54.9MB" or "14.9MB")
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		t.Fatalf("Unexpected docker images output: %s", string(output))
	}

	sizeStr := strings.Fields(lines[1])
	if len(sizeStr) < 2 {
		t.Fatalf("Could not parse image size from: %s", lines[1])
	}

	size := sizeStr[len(sizeStr)-1] // Last field should be size
	t.Logf("Docker image size: %s", size)

	// Verify it's under 50MB (accepting MB or GB units)
	if strings.Contains(size, "GB") {
		t.Errorf("Docker image size %s exceeds 50MB requirement (size in GB)", size)
	} else if strings.Contains(size, "MB") {
		// Extract numeric part
		sizeNum := strings.TrimSuffix(size, "MB")
		if sizeValue, err := strconv.ParseFloat(sizeNum, 64); err == nil {
			if sizeValue > 50.0 {
				t.Errorf("Docker image size %.1fMB exceeds 50MB requirement", sizeValue)
			} else {
				t.Logf("✅ Docker image size %.1fMB is under 50MB requirement", sizeValue)
			}
		}
	}
}
