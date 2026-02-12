package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// absBinaryPath returns the absolute path to the test binary.
// This is needed when running commands with cmd.Dir set to a different directory.
func absBinaryPath(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(binaryPath)
	if err != nil {
		t.Fatalf("Failed to get absolute binary path: %v", err)
	}
	return abs
}

// TestAddModelHelp tests that add-model --help displays usage information (AC: 4)
func TestAddModelHelp(t *testing.T) {
	cmd := exec.Command(binaryPath, "add-model", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected exit code 0 for add-model --help, got error: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	expectedStrings := []string{
		"add-model",
		"ModelName",
		"--fields",
		"Supported Types:",
		"string",
		"int",
		"bool",
		"time",
		"GORM Modifiers:",
		"unique",
		"not_null",
		"index",
		"Examples:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("add-model --help output should contain %q, got: %s", expected, outputStr)
		}
	}
}

// TestAddModelHelpShortFlag tests that add-model -h displays usage information (AC: 4)
func TestAddModelHelpShortFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "add-model", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected exit code 0 for add-model -h, got error: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "add-model") {
		t.Errorf("add-model -h should show help, got: %s", string(output))
	}
}

// TestAddModelMissingModelName tests error when model name is missing (AC: 3)
func TestAddModelMissingModelName(t *testing.T) {
	cmd := exec.Command(binaryPath, "add-model")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for missing model name")
	}

	if !strings.Contains(string(output), "model name is required") {
		t.Errorf("Expected 'model name is required' error, got: %s", string(output))
	}
}

// TestAddModelMissingFields tests error when --fields is missing (AC: 3)
func TestAddModelMissingFields(t *testing.T) {
	cmd := exec.Command(binaryPath, "add-model", "Todo")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for missing --fields flag")
	}

	if !strings.Contains(string(output), "--fields flag is required") {
		t.Errorf("Expected '--fields flag is required' error, got: %s", string(output))
	}
}

// TestAddModelInvalidModelName tests error for invalid model names (AC: 3)
func TestAddModelInvalidModelName(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
	}{
		{"lowercase", "todo"},
		{"with hyphen", "Blog-Post"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, "add-model", tt.modelName, "--fields", "title:string")
			output, err := cmd.CombinedOutput()

			if err == nil {
				t.Errorf("Expected error for invalid model name %q", tt.modelName)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, "model name") {
				t.Errorf("Expected model name validation error for %q, got: %s", tt.modelName, outputStr)
			}
		})
	}
}

// TestAddModelNotInProject tests error when not in a go-starter-kit project (AC: 3)
func TestAddModelNotInProject(t *testing.T) {
	// Create an empty temp directory (not a go-starter-kit project)
	tmpDir := t.TempDir()

	cmd := exec.Command(absBinaryPath(t), "add-model", "Todo", "--fields", "title:string")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error when not in a go-starter-kit project")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "go.mod not found") && !strings.Contains(outputStr, "not a") {
		t.Errorf("Expected project detection error, got: %s", outputStr)
	}
}

// TestAddModelExistingModel tests error when model already exists (AC: 3)
func TestAddModelExistingModel(t *testing.T) {
	// Create a fake go-starter-kit project
	tmpDir := t.TempDir()
	setupFakeProject(t, tmpDir)

	// Create an existing model file
	modelPath := filepath.Join(tmpDir, "internal", "models", "todo.go")
	if err := os.WriteFile(modelPath, []byte("package models\n"), 0644); err != nil {
		t.Fatalf("Failed to create existing model file: %v", err)
	}

	cmd := exec.Command(absBinaryPath(t), "add-model", "Todo", "--fields", "title:string")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for existing model")
	}

	if !strings.Contains(string(output), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %s", string(output))
	}
}

// TestAddModelInvalidFields tests error for invalid field definitions (AC: 3)
func TestAddModelInvalidFields(t *testing.T) {
	// Create a fake go-starter-kit project
	tmpDir := t.TempDir()
	setupFakeProject(t, tmpDir)

	cmd := exec.Command(absBinaryPath(t), "add-model", "Todo", "--fields", "title:invalidtype")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for invalid field type")
	}

	if !strings.Contains(string(output), "unsupported type") {
		t.Errorf("Expected 'unsupported type' error, got: %s", string(output))
	}
}

// TestAddModelSubcommandRouting tests that add-model routes correctly and doesn't interfere with create (AC: 1)
func TestAddModelSubcommandRouting(t *testing.T) {
	// Regular project creation should still work
	testProjectName := "test-routing-project"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, testProjectName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Regular project creation should still work, got error: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "Creating project:") {
		t.Errorf("Expected project creation message, got: %s", string(output))
	}
}

// TestAddModelFieldParsingSummary tests that valid fields show a summary (AC: 1)
func TestAddModelFieldParsingSummary(t *testing.T) {
	tmpDir := t.TempDir()
	setupFakeProject(t, tmpDir)

	// Use --yes to skip confirmation
	cmd := exec.Command(absBinaryPath(t), "add-model", "Todo", "--fields", "title:string,completed:bool", "--yes")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success for valid add-model with --yes, got error: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Todo") {
		t.Errorf("Expected model name 'Todo' in summary, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "postgres") {
		t.Errorf("Expected database type in summary, got: %s", outputStr)
	}
}

// TestDetectDatabase tests database detection from go.mod content (AC: 2)
func TestDetectDatabase(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantDB  string
		wantErr bool
	}{
		{
			name:    "postgres",
			content: "module myproject\n\nrequire gorm.io/driver/postgres v1.5.0",
			wantDB:  "postgres",
		},
		{
			name:    "mysql",
			content: "module myproject\n\nrequire gorm.io/driver/mysql v1.5.0",
			wantDB:  "mysql",
		},
		{
			name:    "sqlite",
			content: "module myproject\n\nrequire gorm.io/driver/sqlite v1.5.0",
			wantDB:  "sqlite",
		},
		{
			name:    "no driver",
			content: "module myproject\n\nrequire github.com/some/lib v1.0.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := detectDatabase(tt.content)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("detectDatabase() returned error: %v", err)
			}
			if db != tt.wantDB {
				t.Errorf("detectDatabase() = %q, want %q", db, tt.wantDB)
			}
		})
	}
}

// TestExtractModulePath tests module path extraction from go.mod content
func TestExtractModulePath(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "standard module",
			content: "module github.com/user/project\n\ngo 1.25",
			want:    "github.com/user/project",
		},
		{
			name:    "simple module",
			content: "module myproject\n\ngo 1.25",
			want:    "myproject",
		},
		{
			name:    "empty content",
			content: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractModulePath(tt.content)
			if got != tt.want {
				t.Errorf("extractModulePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

// setupFakeProject creates a minimal go-starter-kit project structure for testing.
func setupFakeProject(t *testing.T, dir string) {
	t.Helper()

	// Create directory structure
	dirs := []string{
		"cmd",
		"internal/models",
		"internal/interfaces",
		"internal/adapters/repository",
		"internal/adapters/handlers",
		"internal/adapters/http",
		"internal/infrastructure/database",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(dir, d), 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", d, err)
		}
	}

	// Create go.mod with postgres driver and required go-starter-kit dependencies
	goModContent := `module github.com/test/myproject

go 1.25

require (
	github.com/gofiber/fiber/v2 v2.52.5
	go.uber.org/fx v1.23.0
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create database.go with AutoMigrate
	databaseContent := `package database

import (
	"gorm.io/gorm"
	"github.com/test/myproject/internal/models"
)

func NewDatabase() (*gorm.DB, error) {
	// ... connection setup ...
	if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		return nil, err
	}
	return db, nil
}
`
	if err := os.WriteFile(filepath.Join(dir, "internal", "infrastructure", "database", "database.go"), []byte(databaseContent), 0644); err != nil {
		t.Fatalf("Failed to create database.go: %v", err)
	}

	// Create repository/module.go with fx.Module
	repoModuleContent := `package repository

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"github.com/test/myproject/internal/interfaces"
)

var Module = fx.Module("repository",
	fx.Provide(func(db *gorm.DB) interfaces.UserRepository {
		return NewUserRepository(db)
	}),
)
`
	if err := os.WriteFile(filepath.Join(dir, "internal", "adapters", "repository", "module.go"), []byte(repoModuleContent), 0644); err != nil {
		t.Fatalf("Failed to create repository module.go: %v", err)
	}

	// Create handlers/module.go with fx.Module
	handlerModuleContent := `package handlers

import (
	"go.uber.org/fx"
	"github.com/test/myproject/internal/domain/user"
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
	if err := os.WriteFile(filepath.Join(dir, "internal", "adapters", "handlers", "module.go"), []byte(handlerModuleContent), 0644); err != nil {
		t.Fatalf("Failed to create handlers module.go: %v", err)
	}

	// Create routes.go
	routesContent := `package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/test/myproject/internal/adapters/handlers"
)

func RegisterRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware fiber.Handler,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)

	// User routes (protected)
	users := v1.Group("/users", authMiddleware)
	users.Get("/me", userHandler.GetMe)
}
`
	if err := os.WriteFile(filepath.Join(dir, "internal", "adapters", "http", "routes.go"), []byte(routesContent), 0644); err != nil {
		t.Fatalf("Failed to create routes.go: %v", err)
	}

	// Create cmd/main.go with fx.New
	mainContent := `package main

import (
	"github.com/test/myproject/internal/domain/user"
	"github.com/test/myproject/internal/adapters/repository"
	"github.com/test/myproject/internal/adapters/handlers"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Domain services
		user.Module,

		// Data persistence
		repository.Module,

		// HTTP handlers
		handlers.Module,
	).Run()
}
`
	if err := os.WriteFile(filepath.Join(dir, "cmd", "main.go"), []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}
}

// TestE2EAddModelWithRealProject tests add-model on an actual generated project (AC: 1, 2, 3, 4)
func TestE2EAddModelWithRealProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	// Create a real project first
	testProjectName := "e2e-test-add-model-" + t.Name()
	defer os.RemoveAll(testProjectName)

	// Step 1: Generate a real go-starter-kit project
	createCmd := exec.Command(binaryPath, testProjectName, "--database", "postgres")
	createOutput, err := createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create test project: %v\nOutput: %s", err, string(createOutput))
	}

	// Verify project was created
	if _, err := os.Stat(testProjectName); os.IsNotExist(err) {
		t.Fatalf("Test project directory was not created")
	}

	// Step 2: Run add-model command in the generated project
	addModelCmd := exec.Command(absBinaryPath(t), "add-model", "Product", "--fields", "name:string:not_null,price:float64,stock:int", "--yes")
	addModelCmd.Dir = testProjectName
	addModelOutput, err := addModelCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("add-model failed on real project: %v\nOutput: %s", err, string(addModelOutput))
	}

	// Step 3: Verify the output contains expected information
	outputStr := string(addModelOutput)
	expectedContents := []string{
		"Product",                // Model name
		"postgres",               // Database type detected
		"Name (string)",          // Field parsed correctly (PascalCase)
		"Price (float64)",        // Field parsed correctly (PascalCase)
		"Stock (int)",            // Field parsed correctly (PascalCase)
		"[not_null]",             // Modifier shown
		"generated successfully", // Success message (Story 8.2)
	}

	for _, expected := range expectedContents {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("E2E add-model output should contain %q, got: %s", expected, outputStr)
		}
	}

	// Step 4: Verify generated files exist
	generatedFiles := []string{
		filepath.Join(testProjectName, "internal", "models", "product.go"),
		filepath.Join(testProjectName, "internal", "interfaces", "product_repository.go"),
		filepath.Join(testProjectName, "internal", "adapters", "repository", "product_repository.go"),
		filepath.Join(testProjectName, "internal", "domain", "product", "service.go"),
		filepath.Join(testProjectName, "internal", "domain", "product", "module.go"),
		filepath.Join(testProjectName, "internal", "adapters", "handlers", "product_handler.go"),
	}
	for _, f := range generatedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("Expected generated file %s to exist", f)
		}
	}

	// Step 5: Verify AutoMigrate was updated
	dbContent, err := os.ReadFile(filepath.Join(testProjectName, "internal", "infrastructure", "database", "database.go"))
	if err != nil {
		t.Fatalf("Failed to read database.go: %v", err)
	}
	if !strings.Contains(string(dbContent), "&models.Product{}") {
		t.Error("Expected database.go to contain &models.Product{} in AutoMigrate")
	}

	// Step 6: Verify repository module was updated
	repoModContent, err := os.ReadFile(filepath.Join(testProjectName, "internal", "adapters", "repository", "module.go"))
	if err != nil {
		t.Fatalf("Failed to read repository module.go: %v", err)
	}
	if !strings.Contains(string(repoModContent), "NewProductRepository") {
		t.Error("Expected repository module.go to contain NewProductRepository")
	}

	// Step 7: Verify handler module was updated
	handlerModContent, err := os.ReadFile(filepath.Join(testProjectName, "internal", "adapters", "handlers", "module.go"))
	if err != nil {
		t.Fatalf("Failed to read handlers module.go: %v", err)
	}
	if !strings.Contains(string(handlerModContent), "NewProductHandler") {
		t.Error("Expected handlers module.go to contain NewProductHandler")
	}

	// Step 8: Verify routes were updated
	routesContent, err := os.ReadFile(filepath.Join(testProjectName, "internal", "adapters", "http", "routes.go"))
	if err != nil {
		t.Fatalf("Failed to read routes.go: %v", err)
	}
	routesStr := string(routesContent)
	if !strings.Contains(routesStr, "productHandler") {
		t.Error("Expected routes.go to contain productHandler parameter")
	}
	if !strings.Contains(routesStr, "/products") {
		t.Error("Expected routes.go to contain /products route group")
	}

	// Step 9: Verify main.go fx module was updated
	mainContent, err := os.ReadFile(filepath.Join(testProjectName, "cmd", "main.go"))
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}
	if !strings.Contains(string(mainContent), "product.Module") {
		t.Error("Expected main.go to contain product.Module")
	}

	// Step 10: Verify go.mod detection worked correctly (should detect postgres)
	if !strings.Contains(outputStr, "postgres") {
		t.Errorf("Expected database detection to find 'postgres', got: %s", outputStr)
	}
}
