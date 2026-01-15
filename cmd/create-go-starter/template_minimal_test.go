package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMinimalGoModTemplate tests that MinimalGoModTemplate generates correct content (AC: 1)
func TestMinimalGoModTemplate(t *testing.T) {
	projectName := "minimal-test-project"
	templates := NewProjectTemplates(projectName)
	content := templates.MinimalGoModTemplate()

	// Check module declaration with project name
	if !strings.Contains(content, "module "+projectName) {
		t.Errorf("MinimalGoModTemplate() should contain 'module %s', got:\n%s", projectName, content)
	}

	// Check Go version
	if !strings.Contains(content, "go 1.25.5") {
		t.Errorf("MinimalGoModTemplate() should contain 'go 1.25.5', got:\n%s", content)
	}

	// Check required dependencies for minimal template
	requiredDeps := []string{
		"github.com/gofiber/fiber/v2",
		"github.com/rs/zerolog",
		"go.uber.org/fx",
		"gorm.io/gorm",
		"gorm.io/driver/postgres",
		"github.com/swaggo/swag",
		"github.com/swaggo/fiber-swagger",
	}

	for _, dep := range requiredDeps {
		if !strings.Contains(content, dep) {
			t.Errorf("MinimalGoModTemplate() should contain '%s' dependency", dep)
		}
	}

	// Check that JWT-related dependencies are NOT present (AC: 2)
	jwtDeps := []string{
		"github.com/gofiber/contrib/jwt",
		"github.com/golang-jwt/jwt",
	}

	for _, dep := range jwtDeps {
		if strings.Contains(content, dep) {
			t.Errorf("MinimalGoModTemplate() should NOT contain '%s' (no auth in minimal)", dep)
		}
	}
}

// TestMinimalMainGoTemplate tests that minimal main.go is correct (AC: 1)
func TestMinimalMainGoTemplate(t *testing.T) {
	projectName := "minimal-test-project"
	templates := NewProjectTemplates(projectName)
	content := templates.MinimalMainGoTemplate()

	// Check package main
	if !strings.Contains(content, "package main") {
		t.Error("MinimalMainGoTemplate() should contain 'package main'")
	}

	// Check fx.New usage
	if !strings.Contains(content, "fx.New") {
		t.Error("MinimalMainGoTemplate() should use fx.New")
	}

	// Check required modules for minimal template
	requiredModules := []string{
		"logger.Module",
		"database.Module",
		"server.Module",
	}

	for _, mod := range requiredModules {
		if !strings.Contains(content, mod) {
			t.Errorf("MinimalMainGoTemplate() should import %s", mod)
		}
	}

	// Check that auth-related modules are NOT present (AC: 2)
	authModules := []string{
		"auth.Module",
		"user.Module",
		"handlers.Module",
		"repository.Module",
	}

	for _, mod := range authModules {
		if strings.Contains(content, mod) {
			t.Errorf("MinimalMainGoTemplate() should NOT import %s (no auth in minimal)", mod)
		}
	}

	// Check Swagger annotations are present
	if !strings.Contains(content, "@title") {
		t.Error("MinimalMainGoTemplate() should have Swagger @title annotation")
	}

	// Check Run() call
	if !strings.Contains(content, ".Run()") {
		t.Error("MinimalMainGoTemplate() should call Run() on fx.App")
	}
}

// TestMinimalRoutesTemplate tests that minimal routes are correct (AC: 1, 4)
func TestMinimalRoutesTemplate(t *testing.T) {
	projectName := "minimal-test-project"
	templates := NewProjectTemplates(projectName)
	content := templates.MinimalRoutesTemplate()

	// Check package declaration
	if !strings.Contains(content, "package http") {
		t.Error("MinimalRoutesTemplate() should contain 'package http'")
	}

	// Check health route registration is present (AC: 4)
	if !strings.Contains(content, "RegisterHealthRoutes") {
		t.Error("MinimalRoutesTemplate() should call RegisterHealthRoutes for /health endpoint")
	}

	// Check swagger route is present (AC: 4)
	if !strings.Contains(content, "/swagger/*") {
		t.Error("MinimalRoutesTemplate() should register /swagger/* route")
	}

	// Check that auth routes are NOT present (AC: 2)
	authRoutes := []string{
		"/auth/register",
		"/auth/login",
		"/auth/refresh",
		"/users",
	}

	for _, route := range authRoutes {
		if strings.Contains(content, route) {
			t.Errorf("MinimalRoutesTemplate() should NOT contain '%s' (no auth in minimal)", route)
		}
	}

	// Check that auth middleware is NOT used
	if strings.Contains(content, "authMiddleware") {
		t.Error("MinimalRoutesTemplate() should NOT use authMiddleware")
	}
}

// TestMinimalServerTemplate tests that minimal server is correct (AC: 1)
func TestMinimalServerTemplate(t *testing.T) {
	projectName := "minimal-test-project"
	templates := NewProjectTemplates(projectName)
	content := templates.MinimalServerTemplate()

	// Check package declaration
	if !strings.Contains(content, "package server") {
		t.Error("MinimalServerTemplate() should contain 'package server'")
	}

	// Check Fiber import
	if !strings.Contains(content, "github.com/gofiber/fiber/v2") {
		t.Error("MinimalServerTemplate() should import fiber")
	}

	// Check fx integration
	if !strings.Contains(content, "go.uber.org/fx") {
		t.Error("MinimalServerTemplate() should import fx")
	}

	// Check Module export
	if !strings.Contains(content, "var Module = fx.Module") {
		t.Error("MinimalServerTemplate() should export fx Module")
	}

	// Check that auth middleware import is NOT present (AC: 2)
	if strings.Contains(content, "internal/adapters/middleware") {
		t.Error("MinimalServerTemplate() should NOT import middleware package (no auth)")
	}

	// Check OnStart and OnStop hooks
	if !strings.Contains(content, "OnStart") {
		t.Error("MinimalServerTemplate() should implement OnStart hook")
	}

	if !strings.Contains(content, "OnStop") {
		t.Error("MinimalServerTemplate() should implement OnStop hook for graceful shutdown")
	}
}

// TestMinimalEnvTemplate tests that minimal .env is correct (AC: 1, 2)
func TestMinimalEnvTemplate(t *testing.T) {
	projectName := "minimal-test-project"
	templates := NewProjectTemplates(projectName)
	content := templates.MinimalEnvTemplate()

	// Check APP_NAME uses project name
	if !strings.Contains(content, "APP_NAME="+projectName) {
		t.Errorf("MinimalEnvTemplate() should contain 'APP_NAME=%s'", projectName)
	}

	// Check DB_NAME uses project name
	if !strings.Contains(content, "DB_NAME="+projectName) {
		t.Errorf("MinimalEnvTemplate() should contain 'DB_NAME=%s'", projectName)
	}

	// Check essential env vars
	essentialVars := []string{
		"APP_ENV=",
		"APP_PORT=",
		"DB_HOST=",
		"DB_PORT=",
		"DB_USER=",
		"DB_PASSWORD=",
	}

	for _, envVar := range essentialVars {
		if !strings.Contains(content, envVar) {
			t.Errorf("MinimalEnvTemplate() should contain '%s'", envVar)
		}
	}

	// Check that JWT-related env vars are NOT present (AC: 2)
	if strings.Contains(content, "JWT_SECRET") {
		t.Error("MinimalEnvTemplate() should NOT contain 'JWT_SECRET' (no auth in minimal)")
	}

	if strings.Contains(content, "JWT_EXPIRY") {
		t.Error("MinimalEnvTemplate() should NOT contain 'JWT_EXPIRY' (no auth in minimal)")
	}
}

// TestMinimalReadmeTemplate tests that minimal README is correct (AC: 1)
func TestMinimalReadmeTemplate(t *testing.T) {
	projectName := "minimal-test-project"
	templates := NewProjectTemplates(projectName)
	content := templates.MinimalReadmeTemplate()

	// Check title uses project name
	if !strings.Contains(content, "# "+projectName) {
		t.Errorf("MinimalReadmeTemplate() should contain '# %s'", projectName)
	}

	// Check essential sections
	essentialSections := []string{
		"## Architecture",
		"## Prérequis",
		"GET /health",
		"GET /swagger",
	}

	for _, section := range essentialSections {
		if !strings.Contains(content, section) {
			t.Errorf("MinimalReadmeTemplate() should contain '%s'", section)
		}
	}

	// Check that auth-related content is NOT present (AC: 2)
	authContent := []string{
		"JWT",
		"Authentification",
		"/auth/register",
		"/auth/login",
	}

	for _, auth := range authContent {
		if strings.Contains(content, auth) {
			t.Errorf("MinimalReadmeTemplate() should NOT contain '%s' (no auth in minimal)", auth)
		}
	}
}

// TestGenerateMinimalProjectFiles tests that minimal template generates correct files (AC: 1, 2, 3)
func TestGenerateMinimalProjectFiles(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "minimal-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the project structure first (using minimal template directories)
	if err := createProjectStructure(projectPath, TemplateMinimal); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate project files with minimal template
	if err := generateProjectFiles(projectPath, projectName, TemplateMinimal); err != nil {
		t.Fatalf("generateProjectFiles() with minimal template error = %v", err)
	}

	// Expected files for minimal template
	expectedFiles := []string{
		"go.mod",
		"cmd/main.go",
		"pkg/config/env.go",
		"pkg/logger/logger.go",
		"internal/adapters/http/health.go",
		"internal/adapters/http/routes.go",
		"internal/infrastructure/database/database.go",
		"internal/infrastructure/server/server.go",
		"Dockerfile",
		"Makefile",
		".env.example",
		".gitignore",
		".golangci.yml",
		"README.md",
		"docs/README.md",
		"docs/docs.go",
		"setup.sh",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist for minimal template", file)
		}
	}

	// Files that should NOT exist in minimal template (AC: 2)
	excludedFiles := []string{
		"pkg/auth/jwt.go",
		"pkg/auth/middleware.go",
		"pkg/auth/module.go",
		"internal/models/user.go",
		"internal/domain/errors.go",
		"internal/domain/user/service.go",
		"internal/domain/user/module.go",
		"internal/interfaces/services.go",
		"internal/interfaces/user_repository.go",
		"internal/adapters/middleware/error_handler.go",
		"internal/adapters/repository/user_repository.go",
		"internal/adapters/repository/module.go",
		"internal/adapters/handlers/auth_handler.go",
		"internal/adapters/handlers/user_handler.go",
		"internal/adapters/handlers/module.go",
	}

	for _, file := range excludedFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); err == nil {
			t.Errorf("File %s should NOT exist in minimal template (no auth)", file)
		}
	}
}

// TestE2EMinimalProjectBuilds tests that generated minimal project compiles (AC: 3)
func TestE2EMinimalProjectBuilds(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tempDir := t.TempDir()
	projectName := "e2e-minimal-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the project structure (using minimal template directories)
	if err := createProjectStructure(projectPath, TemplateMinimal); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate project files with minimal template
	if err := generateProjectFiles(projectPath, projectName, TemplateMinimal); err != nil {
		t.Fatalf("Failed to generate project files with minimal template: %v", err)
	}

	// Run go mod tidy
	t.Run("GoModTidy", func(t *testing.T) {
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("go mod tidy failed: %v\nOutput:\n%s", err, string(output))
		}
	})

	// Build the project (AC: 3)
	t.Run("BuildMinimalProject", func(t *testing.T) {
		cmd := exec.Command("go", "build", "-o", filepath.Join(projectPath, "test-binary"), "./cmd")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Minimal project failed to build: %v\nOutput:\n%s", err, string(output))
		}

		// Verify binary was created
		binaryPath := filepath.Join(projectPath, "test-binary")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			t.Error("Binary was not created after successful build")
		}
	})

	// Verify go.mod content
	t.Run("GoModValidation", func(t *testing.T) {
		goModPath := filepath.Join(projectPath, "go.mod")
		content, err := os.ReadFile(goModPath)
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}

		// Verify no JWT dependencies
		if strings.Contains(string(content), "golang-jwt/jwt") {
			t.Error("Minimal template go.mod should NOT contain JWT dependencies")
		}

		if strings.Contains(string(content), "gofiber/contrib/jwt") {
			t.Error("Minimal template go.mod should NOT contain JWT middleware dependency")
		}
	})
}

// TestMinimalTemplateNoAuthFiles verifies no auth-related files exist (AC: 2)
func TestMinimalTemplateNoAuthFiles(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "no-auth-check"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the project structure (using minimal template directories)
	if err := createProjectStructure(projectPath, TemplateMinimal); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate project files with minimal template
	if err := generateProjectFiles(projectPath, projectName, TemplateMinimal); err != nil {
		t.Fatalf("generateProjectFiles() with minimal template error = %v", err)
	}

	// Read main.go and verify no auth imports
	mainGoPath := filepath.Join(projectPath, "cmd", "main.go")
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	authImports := []string{
		"pkg/auth",
		"internal/domain/user",
		"internal/adapters/handlers",
		"internal/adapters/repository",
	}

	for _, imp := range authImports {
		if strings.Contains(string(content), imp) {
			t.Errorf("Minimal main.go should NOT import '%s'", imp)
		}
	}

	// Read routes.go and verify no auth routes
	routesPath := filepath.Join(projectPath, "internal", "adapters", "http", "routes.go")
	content, err = os.ReadFile(routesPath)
	if err != nil {
		t.Fatalf("Failed to read routes.go: %v", err)
	}

	if strings.Contains(string(content), "authMiddleware") {
		t.Error("Minimal routes.go should NOT reference authMiddleware")
	}

	if strings.Contains(string(content), "authHandler") {
		t.Error("Minimal routes.go should NOT reference authHandler")
	}
}

// TestGetDirectoriesForMinimalTemplate tests that minimal template has correct directories
func TestGetDirectoriesForMinimalTemplate(t *testing.T) {
	dirs := getDirectoriesForTemplate(TemplateMinimal)

	// Required directories for minimal
	requiredDirs := []string{
		"cmd",
		"internal/adapters/http",
		"internal/infrastructure/database",
		"internal/infrastructure/server",
		"pkg/config",
		"pkg/logger",
		"docs",
	}

	for _, dir := range requiredDirs {
		found := false
		for _, d := range dirs {
			if d == dir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("getDirectoriesForTemplate(minimal) should include '%s'", dir)
		}
	}

	// Directories that should NOT be included for minimal
	excludedDirs := []string{
		"pkg/auth",
		"internal/domain/user",
		"internal/adapters/handlers",
		"internal/adapters/repository",
		"internal/adapters/middleware",
		"internal/interfaces",
		"internal/models",
	}

	for _, dir := range excludedDirs {
		for _, d := range dirs {
			if d == dir {
				t.Errorf("getDirectoriesForTemplate(minimal) should NOT include '%s'", dir)
			}
		}
	}
}
