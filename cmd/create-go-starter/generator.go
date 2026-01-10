package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Valid Go module name pattern
var validGoModuleNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// validateGoModuleName validates that a module name is valid for Go modules.
// Valid names must:
// - Start with a letter or number
// - Contain only letters, numbers, hyphens, or underscores
// - Not be empty
func validateGoModuleName(name string) error {
	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if !validGoModuleNamePattern.MatchString(name) {
		return fmt.Errorf("invalid module name '%s': must start with a letter or number and contain only letters, numbers, hyphens, or underscores", name)
	}

	return nil
}

// FileGenerator represents a file to be generated
type FileGenerator struct {
	Path    string
	Content string
}

// generateProjectFiles creates all the initial project files with templates
func generateProjectFiles(projectPath, projectName string) error {
	// Validate that the project directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// Validate the module name
	if err := validateGoModuleName(projectName); err != nil {
		return err
	}

	// Create templates instance
	templates := NewProjectTemplates(projectName)

	// Define all files to generate
	files := []FileGenerator{
		{
			Path:    filepath.Join(projectPath, "go.mod"),
			Content: templates.GoModTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "cmd", "main.go"),
			Content: templates.UpdatedMainGoTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "config", "env.go"),
			Content: templates.ConfigTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "logger", "logger.go"),
			Content: templates.LoggerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "auth", "jwt.go"),
			Content: templates.JWTAuthTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "auth", "middleware.go"),
			Content: templates.JWTMiddlewareTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "auth", "module.go"),
			Content: templates.AuthModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "errors.go"),
			Content: templates.DomainErrorsTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "user", "entity.go"),
			Content: templates.UserEntityTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "user", "refresh_token.go"),
			Content: templates.UserRefreshTokenTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "user", "service.go"),
			Content: templates.UserServiceTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "user", "module.go"),
			Content: templates.UserModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "interfaces", "services.go"),
			Content: templates.UserInterfacesTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "interfaces", "user_repository.go"),
			Content: templates.UserRepositoryInterfaceTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "middleware", "error_handler.go"),
			Content: templates.ErrorHandlerMiddlewareTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "repository", "user_repository.go"),
			Content: templates.UserRepositoryTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "repository", "module.go"),
			Content: templates.RepositoryModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "auth_handler.go"),
			Content: templates.AuthHandlerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "user_handler.go"),
			Content: templates.UserHandlerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "module.go"),
			Content: templates.HandlerModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "http", "health.go"),
			Content: templates.HealthHandlerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "database", "database.go"),
			Content: templates.DatabaseTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "server", "server.go"),
			Content: templates.ServerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "Dockerfile"),
			Content: templates.DockerfileTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docker-compose.yml"),
			Content: templates.DockerComposeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "Makefile"),
			Content: templates.MakefileTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".env.example"),
			Content: templates.EnvTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".gitignore"),
			Content: templates.GitignoreTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".golangci.yml"),
			Content: templates.GolangCILintTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".github", "workflows", "ci.yml"),
			Content: templates.GitHubActionsWorkflowTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "README.md"),
			Content: templates.ReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "README.md"),
			Content: templates.DocsReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "quick-start.md"),
			Content: templates.QuickStartTemplate(),
		},
	}

	// Write all files
	for _, file := range files {
		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}

		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}
	}

	return nil
}
