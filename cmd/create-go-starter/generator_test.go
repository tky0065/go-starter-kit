package main

import (
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

	// Create the project structure first
	if err := createProjectStructure(projectPath); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate project files
	if err := generateProjectFiles(projectPath, projectName); err != nil {
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
	err := generateProjectFiles("/non/existent/path", "test-project")
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
			err := validateGoModuleName(tt.modName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGoModuleName(%s) error = %v, wantErr %v", tt.modName, err, tt.wantErr)
			}
		})
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

	// Create the complete project structure
	if err := createProjectStructure(projectPath); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate all project files
	if err := generateProjectFiles(projectPath, projectName); err != nil {
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
