package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	// Call the function to create directories
	err := createProjectStructure(projectPath)
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
	err = createProjectStructure(projectPath)
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

	err := createProjectStructure(invalidPath)
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
			err := validateProjectName(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("validateProjectName(%q) error = %v, wantError %v", tt.input, err, tt.wantError)
			}
			if err != nil && !strings.Contains(err.Error(), "invalid project name") {
				t.Errorf("validateProjectName(%q) error message should contain 'invalid project name', got: %v", tt.input, err)
			}
		})
	}
}
