package main

import (
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

	// Check essential sections
	essentialSections := []string{
		"## Architecture",
		"## Prerequisites",
		"## Getting Started",
		"## Development",
		"## Project Structure",
	}

	for _, section := range essentialSections {
		if !strings.Contains(content, section) {
			t.Errorf("ReadmeTemplate() should contain section '%s', got:\n%s", section, content)
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
