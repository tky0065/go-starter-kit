package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tky0065/go-starter-kit/pkg/utils" // Added for shared validation
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework); err != nil {
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
	err := generateProjectFiles("/non/existent/path", "test-project", DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework)
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
	err = generateProjectFiles(projectPath, "", DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework)
	if err == nil {
		t.Error("generateProjectFiles() should return error for empty module name")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Error message should mention 'empty', got: %v", err)
	}

	// Test with invalid module name
	err = generateProjectFiles(projectPath, "-invalid", DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework)
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
	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework); err != nil {
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
		"internal/adapters/handlers/health_handler.go",
		"internal/adapters/http/routes.go",
		"deployments/kubernetes/probes.yaml",
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
	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework); err != nil {
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

	if err := generateProjectFiles(projectPath, projectName, DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework); err != nil {
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
	if err := generateProjectFiles(projectPath, projectName, TemplateGraphQL, DefaultDatabase, DefaultObservabilityLevel, DefaultFramework); err != nil {
		t.Fatalf("generateProjectFiles(graphql, DefaultObservabilityLevel, DefaultFramework) error = %v", err)
	}

	// List of all expected GraphQL template files
	expectedFiles := []string{
		"go.mod",
		"cmd/main.go",
		"gqlgen.yml",
		"graph/schema.graphqls",
		"graph/email_helpers.go",
		"graph/resolver.go",
		"graph/schema.resolvers.go",
		"graph/generate.go",
		"graph/model/models.go",
		"graph/model/models_gen.go",
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

	generatedPath := filepath.Join(projectPath, "graph", "generated", "generated.go")
	content, err = os.ReadFile(generatedPath)
	if err != nil {
		t.Errorf("Failed to read graph/generated/generated.go: %v", err)
	}
	if strings.Contains(string(content), "panic(\"Run 'go generate ./...'") {
		t.Error("graph/generated/generated.go should contain runnable gqlgen output, not the placeholder panic")
	}

	resolverPath := filepath.Join(projectPath, "graph", "schema.resolvers.go")
	content, err = os.ReadFile(resolverPath)
	if err != nil {
		t.Errorf("Failed to read graph/schema.resolvers.go: %v", err)
	}
	if !strings.Contains(string(content), "func (r *Resolver) User() generated.UserResolver") {
		t.Error("graph/schema.resolvers.go should expose the generated.UserResolver implementation")
	}
	if !strings.Contains(string(content), "return strconv.FormatUint(uint64(obj.ID), 10), nil") {
		t.Error("graph/schema.resolvers.go should implement the User.id resolver without panic")
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

func getFreePort(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to allocate free port: %v", err)
	}
	defer listener.Close()

	return fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
}

func waitForHTTPReady(t *testing.T, client *http.Client, url string, done <-chan error, logs *bytes.Buffer) {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			t.Fatalf("Generated GraphQL server exited before becoming ready: %v\nOutput:\n%s", err, logs.String())
		default:
		}

		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}

		time.Sleep(250 * time.Millisecond)
	}

	t.Fatalf("Timed out waiting for generated GraphQL server readiness at %s\nOutput:\n%s", url, logs.String())
}

func postGraphQL(t *testing.T, client *http.Client, url, query string) map[string]any {
	t.Helper()

	payload, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatalf("Failed to marshal GraphQL payload: %v", err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("GraphQL request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read GraphQL response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected GraphQL status %d: %s", resp.StatusCode, string(body))
	}

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("Failed to decode GraphQL response: %v\nBody:\n%s", err, string(body))
	}

	if errs, ok := decoded["errors"]; ok {
		t.Fatalf("GraphQL response returned errors: %v", errs)
	}

	return decoded
}

// TestE2EGraphQLProjectBuilds is an end-to-end test that verifies
// a generated GraphQL project can resolve dependencies, start, and answer HTTP requests.
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
	if err := generateProjectFiles(projectPath, projectName, TemplateGraphQL, DatabaseSQLite, DefaultObservabilityLevel, DefaultFramework); err != nil {
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

	t.Run("ServerSmokeTest", func(t *testing.T) {
		port := getFreePort(t)
		baseURL := "http://127.0.0.1:" + port
		client := &http.Client{Timeout: 5 * time.Second}

		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "go", "run", "./cmd")
		cmd.Dir = projectPath
		cmd.Env = append(os.Environ(),
			"APP_PORT="+port,
			"DB_NAME=test-graphql",
		)

		var logs bytes.Buffer
		cmd.Stdout = &logs
		cmd.Stderr = &logs

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start generated GraphQL project: %v", err)
		}

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		defer func() {
			cancel()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Logf("Generated GraphQL server did not exit within timeout")
			}
		}()

		waitForHTTPReady(t, client, baseURL+"/health", done, &logs)

		resp, err := client.Get(baseURL + "/")
		if err != nil {
			t.Fatalf("Failed to query GraphQL playground: %v", err)
		}
		playgroundBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read GraphQL playground response: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected GraphQL playground status %d: %s", resp.StatusCode, string(playgroundBody))
		}
		if !strings.Contains(string(playgroundBody), "GraphQL Playground") {
			t.Fatalf("GraphQL playground response should contain 'GraphQL Playground', got:\n%s", string(playgroundBody))
		}

		resp, err = client.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Failed to query health endpoint: %v", err)
		}
		healthBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read health response: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected health status %d: %s", resp.StatusCode, string(healthBody))
		}
		if !strings.Contains(string(healthBody), `"status":"ok"`) {
			t.Fatalf("Health response should contain status ok, got: %s", string(healthBody))
		}

		healthResp := postGraphQL(t, client, baseURL+"/query", `query { health }`)
		data, ok := healthResp["data"].(map[string]any)
		if !ok || data["health"] != "ok" {
			t.Fatalf("Unexpected GraphQL health response: %#v", healthResp)
		}

		createResp := postGraphQL(t, client, baseURL+"/query", `mutation { createUser(input: { email: "smoke@example.com", password: "password123" }) { id email } }`)
		createData, ok := createResp["data"].(map[string]any)
		if !ok {
			t.Fatalf("Unexpected GraphQL createUser response: %#v", createResp)
		}
		user, ok := createData["createUser"].(map[string]any)
		if !ok {
			t.Fatalf("GraphQL createUser response missing createUser payload: %#v", createResp)
		}
		if user["email"] != "smoke@example.com" {
			t.Fatalf("Expected created user email smoke@example.com, got %#v", user["email"])
		}
		if user["id"] == "" {
			t.Fatalf("Expected created user ID to be populated, got %#v", user["id"])
		}

		usersResp := postGraphQL(t, client, baseURL+"/query", `query { users(page: 1, limit: 10) { users { id email } pageInfo { total hasNextPage } } }`)
		usersData, ok := usersResp["data"].(map[string]any)
		if !ok {
			t.Fatalf("Unexpected GraphQL users response: %#v", usersResp)
		}
		usersConnection, ok := usersData["users"].(map[string]any)
		if !ok {
			t.Fatalf("GraphQL users response missing connection payload: %#v", usersResp)
		}
		users, ok := usersConnection["users"].([]any)
		if !ok || len(users) == 0 {
			t.Fatalf("Expected at least one user in GraphQL users query, got %#v", usersConnection["users"])
		}
		firstUser, ok := users[0].(map[string]any)
		if !ok {
			t.Fatalf("Unexpected user payload: %#v", users[0])
		}
		if firstUser["id"] == "" || firstUser["email"] != "smoke@example.com" {
			t.Fatalf("Unexpected first user payload: %#v", firstUser)
		}
		pageInfo, ok := usersConnection["pageInfo"].(map[string]any)
		if !ok || pageInfo["total"] == nil {
			t.Fatalf("Unexpected pageInfo payload: %#v", usersConnection["pageInfo"])
		}
	})
}
