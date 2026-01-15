package main

import (
	"github.com/tky0065/go-starter-kit/pkg/utils" // Added for shared validation
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateProjectFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	projectName := "test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the project structure first (using full template)
	if err := createProjectStructure(projectPath, TemplateFull); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate project files
	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate); err != nil {
		t.Fatalf("generateProjectFiles() error = %v", err)
	}

	// Test that go.mod exists and contains project name
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Errorf("Failed to read go.mod: %v", err)
	}
	if !strings.Contains(string(content), "module "+projectName) {
		t.Errorf("go.mod should contain 'module %s', got:\n%s", projectName, string(content))
	}

	// Test that main.go exists and contains project name
	mainGoPath := filepath.Join(projectPath, "cmd", "main.go")
	content, err = os.ReadFile(mainGoPath)
	if err != nil {
		t.Errorf("Failed to read cmd/main.go: %v", err)
	}
	if !strings.Contains(string(content), projectName) {
		t.Errorf("cmd/main.go should contain project name '%s', got:\n%s", projectName, string(content))
	}

	// Test that Dockerfile exists and contains project name
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	content, err = os.ReadFile(dockerfilePath)
	if err != nil {
		t.Errorf("Failed to read Dockerfile: %v", err)
	}
	if !strings.Contains(string(content), projectName) {
		t.Errorf("Dockerfile should contain '%s', got:\n%s", projectName, string(content))
	}

	// Test that Makefile exists and contains project name
	makefilePath := filepath.Join(projectPath, "Makefile")
	content, err = os.ReadFile(makefilePath)
	if err != nil {
		t.Errorf("Failed to read Makefile: %v", err)
	}
	if !strings.Contains(string(content), "BINARY_NAME="+projectName) {
		t.Errorf("Makefile should contain 'BINARY_NAME=%s', got:\n%s", projectName, string(content))
	}

	// Test that .env.example exists
	envPath := filepath.Join(projectPath, ".env.example")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Error(".env.example should exist")
	}

	// Test that .gitignore exists
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Error(".gitignore should exist")
	}

	// Test that README.md exists
	readmePath := filepath.Join(projectPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Error("README.md should exist")
	}
}

func TestGenerateProjectFilesWithInvalidPath(t *testing.T) {
	// Test with non-existent directory
	err := generateProjectFiles("/non/existent/path", "test-project", DefaultTemplate)
	if err == nil {
		t.Error("generateProjectFiles() should return error for non-existent path")
	}
}

func TestValidateGoModuleName(t *testing.T) {
	tests := []struct {
		name    string
		modName string
		wantErr bool
	}{
		{
			name:    "valid simple name",
			modName: "myproject",
			wantErr: false,
		},
		{
			name:    "valid with hyphens",
			modName: "my-awesome-project",
			wantErr: false,
		},
		{
			name:    "valid with underscores",
			modName: "my_cool_app",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			modName: "myapp2024",
			wantErr: false,
		},
		{
			name:    "invalid with spaces",
			modName: "my project",
			wantErr: true,
		},
		{
			name:    "invalid starting with hyphen",
			modName: "-myproject",
			wantErr: true,
		},
		{
			name:    "invalid starting with underscore",
			modName: "_myproject",
			wantErr: true,
		},
		{
			name:    "invalid with special chars",
			modName: "my@project",
			wantErr: true,
		},
		{
			name:    "empty name",
			modName: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateGoModuleName(tt.modName) // Changed to utils.ValidateGoModuleName
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGoModuleName(%s) error = %v, wantErr %v", tt.modName, err, tt.wantErr)
			}
		})
	}
}

// TestGenerateProjectFilesWithInvalidModuleName tests that invalid module names are rejected
func TestGenerateProjectFilesWithInvalidModuleName(t *testing.T) {
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "test-project")

	// Create the project directory first
	if err := os.Mkdir(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Test with empty module name
	var err error // Declared err here
	err = generateProjectFiles(projectPath, "", DefaultTemplate)
	if err == nil {
		t.Error("generateProjectFiles() should return error for empty module name")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Error message should mention 'empty', got: %v", err)
	}

	// Test with invalid module name
	err = generateProjectFiles(projectPath, "-invalid", DefaultTemplate)
	if err == nil {
		t.Error("generateProjectFiles() should return error for invalid module name")
	}
	if !strings.Contains(err.Error(), "invalid module name") {
		t.Errorf("Error message should mention 'invalid module name', got: %v", err)
	}
}

// TestGenerateProjectFilesCreatesAllRequiredFiles tests that all expected files are created
func TestGenerateProjectFilesCreatesAllRequiredFiles(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "complete-test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the project structure first (using full template)
	if err := createProjectStructure(projectPath, TemplateFull); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate project files
	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate); err != nil {
		t.Fatalf("generateProjectFiles() error = %v", err)
	}

	// List of all expected files
	expectedFiles := []string{
		"go.mod",
		"cmd/main.go",
		"pkg/config/env.go",
		"pkg/logger/logger.go",
		"pkg/auth/jwt.go",
		"pkg/auth/middleware.go",
		"pkg/auth/module.go",
		"internal/domain/errors.go",
		"internal/models/user.go",
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
		"internal/adapters/http/health.go",
		"internal/adapters/http/routes.go",
		"internal/infrastructure/database/database.go",
		"internal/infrastructure/server/server.go",
		"Dockerfile",
		"docker-compose.yml",
		"Makefile",
		".env.example",
		".gitignore",
		".golangci.yml",
		".github/workflows/ci.yml",
		"README.md",
		"docs/README.md",
		"docs/docs.go",
		"docs/quick-start.md",
		"setup.sh",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", file)
		}
	}

	// Verify setup.sh is executable
	setupPath := filepath.Join(projectPath, "setup.sh")
	info, err := os.Stat(setupPath)
	if err != nil {
		t.Fatalf("Failed to stat setup.sh: %v", err)
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Error("setup.sh should be executable")
	}
}

// TestE2EGeneratedProjectBuilds is an end-to-end test that verifies
// a generated project can actually be built successfully
func TestE2EGeneratedProjectBuilds(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	projectName := "e2e-test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the complete project structure (using full template)
	if err := createProjectStructure(projectPath, TemplateFull); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate all project files
	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate); err != nil {
		t.Fatalf("Failed to generate project files: %v", err)
	}

	// Try to build the generated project
	t.Run("BuildGeneratedProject", func(t *testing.T) {
		// First, tidy dependencies to generate go.sum
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectPath
		tidyOutput, err := tidyCmd.CombinedOutput()
		if err != nil {
			t.Errorf("go mod tidy failed: %v\nOutput:\n%s", err, string(tidyOutput))
			return
		}

		// Then build the project
		cmd := exec.Command("go", "build", "-o", filepath.Join(projectPath, "test-binary"), "./cmd")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Generated project failed to build: %v\nOutput:\n%s", err, string(output))
			return
		}

		// Verify binary was created
		binaryPath := filepath.Join(projectPath, "test-binary")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			t.Error("Binary was not created after successful build")
		}
	})

	// Try to run go mod verify
	t.Run("GoModVerify", func(t *testing.T) {
		cmd := exec.Command("go", "mod", "verify")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("go mod verify failed: %v\nOutput:\n%s", err, string(output))
		}
	})

	// Verify go.mod is valid
	t.Run("GoModValidation", func(t *testing.T) {
		cmd := exec.Command("go", "list", "-m")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("go list -m failed: %v\nOutput:\n%s", err, string(output))
		}

		// Check that the module name is correct
		if !strings.Contains(string(output), projectName) {
			t.Errorf("Module name should be '%s', got: %s", projectName, string(output))
		}
	})
}

// TestGoModTidyWorkflow specifically tests go mod tidy command execution (AC#4)
func TestGoModTidyWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	tmpDir := t.TempDir()
	projectName := "test-mod-tidy-workflow"
	projectPath := filepath.Join(tmpDir, projectName)

	// Create project structure and files (using full template)
	if err := createProjectStructure(projectPath, TemplateFull); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate); err != nil {
		t.Fatalf("Failed to generate project files: %v", err)
	}

	// Verify go.mod exists
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Fatal("go.mod was not created")
	}

	// Test go mod tidy execution
	t.Run("GoModTidyExecution", func(t *testing.T) {
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("go mod tidy failed: %v\nOutput:\n%s", err, string(output))
			return
		}

		// Verify go.sum is created after tidy
		goSumPath := filepath.Join(projectPath, "go.sum")
		if _, err := os.Stat(goSumPath); os.IsNotExist(err) {
			t.Error("go.sum was not created after go mod tidy")
		}

		t.Logf("✅ go mod tidy executed successfully")
	})

	// Test that generated dependencies are valid
	t.Run("DependenciesValid", func(t *testing.T) {
		// Verify that we can download dependencies
		cmd := exec.Command("go", "mod", "download")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("go mod download failed: %v\nOutput:\n%s", err, string(output))
			return
		}

		t.Logf("✅ Dependencies downloaded successfully")
	})
}

// TestGenerateGraphQLTemplateFiles tests the GraphQL template generation
func TestGenerateGraphQLTemplateFiles(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "graphql-test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the project structure first (using graphql template)
	if err := createProjectStructure(projectPath, TemplateGraphQL); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate GraphQL project files
	if err := generateProjectFiles(projectPath, projectName, TemplateGraphQL); err != nil {
		t.Fatalf("generateProjectFiles(graphql) error = %v", err)
	}

	// List of all expected GraphQL template files
	expectedFiles := []string{
		"go.mod",
		"cmd/main.go",
		"gqlgen.yml",
		"graph/schema.graphqls",
		"graph/resolver.go",
		"graph/schema.resolvers.go",
		"graph/generate.go",
		"graph/model/models.go",
		"graph/generated/generated.go",
		"internal/infrastructure/server/server.go",
		"internal/infrastructure/database/database.go",
		"internal/infrastructure/database/user_repository.go",
		"internal/interfaces/user_repository.go",
		"internal/models/user.go",
		"pkg/config/env.go",
		"pkg/logger/logger.go",
		".env.example",
		".gitignore",
		".golangci.yml",
		".github/workflows/ci.yml",
		"Dockerfile",
		"docker-compose.yml",
		"Makefile",
		"README.md",
		"docs/README.md",
		"docs/quick-start.md",
		"setup.sh",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected GraphQL file %s does not exist", file)
		}
	}

	// Test that go.mod contains gqlgen dependencies
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Errorf("Failed to read go.mod: %v", err)
	}
	if !strings.Contains(string(content), "github.com/99designs/gqlgen") {
		t.Error("go.mod should contain gqlgen dependency")
	}
	if !strings.Contains(string(content), "github.com/gofiber/adaptor") {
		t.Error("go.mod should contain gofiber/adaptor dependency")
	}

	// Test that gqlgen.yml contains correct configuration
	gqlgenPath := filepath.Join(projectPath, "gqlgen.yml")
	content, err = os.ReadFile(gqlgenPath)
	if err != nil {
		t.Errorf("Failed to read gqlgen.yml: %v", err)
	}
	if !strings.Contains(string(content), "graph/generated/generated.go") {
		t.Error("gqlgen.yml should configure generated output path")
	}

	// Test that schema.graphqls contains expected types
	schemaPath := filepath.Join(projectPath, "graph", "schema.graphqls")
	content, err = os.ReadFile(schemaPath)
	if err != nil {
		t.Errorf("Failed to read schema.graphqls: %v", err)
	}
	if !strings.Contains(string(content), "type User") {
		t.Error("schema.graphqls should contain User type")
	}
	if !strings.Contains(string(content), "type Query") {
		t.Error("schema.graphqls should contain Query type")
	}

	// Verify setup.sh is executable
	setupPath := filepath.Join(projectPath, "setup.sh")
	info, err := os.Stat(setupPath)
	if err != nil {
		t.Fatalf("Failed to stat setup.sh: %v", err)
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Error("setup.sh should be executable")
	}
}

// TestGetDirectoriesForGraphQLTemplate tests that correct directories are created for GraphQL template
func TestGetDirectoriesForGraphQLTemplate(t *testing.T) {
	dirs := getDirectoriesForTemplate(TemplateGraphQL)

	// Check for GraphQL-specific directories
	expectedDirs := []string{
		"graph",
		"graph/model",
		"graph/generated",
		"internal/interfaces",
		"internal/models",
	}

	for _, expected := range expectedDirs {
		found := false
		for _, dir := range dirs {
			if dir == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected directory %s not found in GraphQL template directories", expected)
		}
	}

	// Check that auth-related directories are NOT included (GraphQL template doesn't need JWT auth)
	authDirs := []string{
		"pkg/auth",
		"internal/domain/user",
		"internal/adapters/handlers",
	}
	for _, authDir := range authDirs {
		for _, dir := range dirs {
			if dir == authDir {
				t.Errorf("GraphQL template should not include %s directory", authDir)
			}
		}
	}
}

// TestE2EGraphQLProjectBuilds is an end-to-end test that verifies
// a generated GraphQL project can pass go mod tidy
func TestE2EGraphQLProjectBuilds(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	projectName := "e2e-graphql-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the complete project structure (using graphql template)
	if err := createProjectStructure(projectPath, TemplateGraphQL); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate all project files
	if err := generateProjectFiles(projectPath, projectName, TemplateGraphQL); err != nil {
		t.Fatalf("Failed to generate project files: %v", err)
	}

	// Run go mod tidy to verify all dependencies are valid
	t.Run("GoModTidy", func(t *testing.T) {
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("go mod tidy failed for GraphQL project: %v\nOutput:\n%s", err, string(output))
			return
		}
		t.Logf("✅ go mod tidy executed successfully for GraphQL project")
	})

	// Verify go.mod is valid
	t.Run("GoModValidation", func(t *testing.T) {
		cmd := exec.Command("go", "list", "-m")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("go list -m failed: %v\nOutput:\n%s", err, string(output))
		}

		// Check that the module name is correct
		if !strings.Contains(string(output), projectName) {
			t.Errorf("Module name should be '%s', got: %s", projectName, string(output))
		}
	})
}
