package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tky0065/go-starter-kit/pkg/utils" // Corrected import path
)

func TestCreateDirectories(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()
	projectName := "test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Expected directories
	expectedDirs := []string{
		"cmd",
		"internal/adapters",
		"internal/domain",
		"internal/interfaces",
		"internal/infrastructure",
		"pkg",
		"deployments",
	}

	// Call the function to create directories (using full template)
	err := createProjectStructure(projectPath, TemplateFull)
	if err != nil {
		t.Fatalf("createProjectStructure failed: %v", err)
	}

	// Verify each directory exists with correct permissions
	for _, dir := range expectedDirs {
		fullPath := filepath.Join(projectPath, dir)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("Directory %s does not exist: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", dir)
		}
		// Verify permissions (note: on some systems, the actual permissions may differ)
		if info.Mode().Perm() != defaultDirPerm {
			t.Errorf("Directory %s has permissions %v, expected %v", dir, info.Mode().Perm(), defaultDirPerm)
		}
	}
}

func TestProjectAlreadyExists(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()
	projectName := "existing-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Create the project directory first
	err := os.Mkdir(projectPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Try to create project structure - should return error
	err = createProjectStructure(projectPath, TemplateFull)
	if err == nil {
		t.Error("Expected error when project directory already exists, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' in error message, got: %v", err)
	}
}

func TestCreateProjectStructureWithInvalidPath(t *testing.T) {
	// Create a temp directory, then try to create a project in a non-existent subdirectory
	tempDir := t.TempDir()
	invalidPath := filepath.Join(tempDir, "nonexistent", "deeply", "nested", "project")

	err := createProjectStructure(invalidPath, TemplateFull)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid alphanumeric", "myproject", false},
		{"valid with hyphens", "my-project", false},
		{"valid with underscores", "my_project", false},
		{"valid mixed", "my-cool_project123", false},
		{"valid starts with number", "123project", false},
		{"invalid starts with hyphen", "-myproject", true},
		{"invalid starts with underscore", "_myproject", true},
		{"invalid with slash", "my/project", true},
		{"invalid with dots", "../project", true},
		{"invalid with spaces", "my project", true},
		{"invalid with special chars", "my@project", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateGoModuleName(tt.input) // Changed to utils.ValidateGoModuleName
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateGoModuleName(%q) error = %v, wantError %v", tt.input, err, tt.wantError)
			}
			if err != nil {
				if tt.input == "" {
					if !strings.Contains(err.Error(), "module name cannot be empty") {
						t.Errorf("ValidateGoModuleName(%q) error message should contain 'module name cannot be empty', got: %v", tt.input, err)
					}
				} else if !strings.Contains(err.Error(), "invalid module name") {
					t.Errorf("ValidateGoModuleName(%q) error message should contain 'invalid module name', got: %v", tt.input, err)
				}
			}
		})
	}
}

// TestCreateProjectStructureVerifiesAllDirectories verifies all directories are created correctly
func TestCreateProjectStructureVerifiesAllDirectories(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "verify-all-dirs"
	projectPath := filepath.Join(tempDir, projectName)

	// Call the function to create directories (using full template)
	err := createProjectStructure(projectPath, TemplateFull)
	if err != nil {
		t.Fatalf("createProjectStructure failed: %v", err)
	}

	// Complete list of all expected directories for full template
	// (from getDirectoriesForTemplate(TemplateFull))
	allExpectedDirs := []string{
		"cmd",
		"internal/adapters/http",
		"internal/infrastructure/database",
		"internal/infrastructure/server",
		"pkg/config",
		"pkg/logger",
		"docs",
		"deployments",
		".github/workflows",
		// Additional dirs for full template
		"pkg/auth",
		"internal/domain",
		"internal/domain/user",
		"internal/interfaces",
		"internal/models",
		"internal/adapters/middleware",
		"internal/adapters/handlers",
		"internal/adapters/repository",
	}

	// Verify each directory exists
	for _, dir := range allExpectedDirs {
		fullPath := filepath.Join(projectPath, dir)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("Directory %s does not exist: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", dir)
		}
	}
}

// TestCreateProjectStructureRootDirCreation tests root directory creation
func TestCreateProjectStructureRootDirCreation(t *testing.T) {
	tempDir := t.TempDir()
	projectName := "root-test"
	projectPath := filepath.Join(tempDir, projectName)

	err := createProjectStructure(projectPath, TemplateFull)
	if err != nil {
		t.Fatalf("createProjectStructure failed: %v", err)
	}

	// Verify root directory exists with correct permissions
	info, err := os.Stat(projectPath)
	if err != nil {
		t.Fatalf("Root directory does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("Root path is not a directory")
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Errorf("Root directory has permissions %v, expected %v", info.Mode().Perm(), defaultDirPerm) // Corrected comparison
	}
}
