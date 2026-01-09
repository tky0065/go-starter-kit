package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCopyEnvFile verifies that .env.example is correctly copied to .env
func TestCopyEnvFile(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.Mkdir(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}

	// Create .env.example file
	envExamplePath := filepath.Join(projectPath, ".env.example")
	envExampleContent := "APP_NAME=test\nAPP_PORT=3000\n"
	if err := os.WriteFile(envExamplePath, []byte(envExampleContent), 0644); err != nil {
		t.Fatalf("Failed to create .env.example: %v", err)
	}

	// Test copying .env.example to .env
	if err := copyEnvFile(projectPath); err != nil {
		t.Fatalf("copyEnvFile() failed: %v", err)
	}

	// Verify .env file was created
	envPath := filepath.Join(projectPath, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Fatal(".env file was not created")
	}

	// Verify .env content matches .env.example
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("Failed to read .env file: %v", err)
	}

	if string(envContent) != envExampleContent {
		t.Errorf("Expected .env content to match .env.example\nGot: %s\nWant: %s", string(envContent), envExampleContent)
	}
}

// TestCopyEnvFileSkipsIfExists verifies that existing .env files are not overwritten
func TestCopyEnvFileSkipsIfExists(t *testing.T) {
	tmpDir := t.TempDir()

	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.Mkdir(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}

	// Create .env.example
	envExamplePath := filepath.Join(projectPath, ".env.example")
	envExampleContent := "APP_NAME=example\n"
	if err := os.WriteFile(envExamplePath, []byte(envExampleContent), 0644); err != nil {
		t.Fatalf("Failed to create .env.example: %v", err)
	}

	// Create existing .env with different content
	envPath := filepath.Join(projectPath, ".env")
	existingContent := "APP_NAME=existing\nAPP_SECRET=secret123\n"
	if err := os.WriteFile(envPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to create existing .env: %v", err)
	}

	// Attempt to copy - should skip
	if err := copyEnvFile(projectPath); err != nil {
		t.Fatalf("copyEnvFile() failed: %v", err)
	}

	// Verify .env was NOT overwritten
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("Failed to read .env file: %v", err)
	}

	if string(envContent) != existingContent {
		t.Errorf("Existing .env file was overwritten\nGot: %s\nWant: %s", string(envContent), existingContent)
	}
}

// TestCopyEnvFileErrorsIfNoExample verifies appropriate error when .env.example is missing
func TestCopyEnvFileErrorsIfNoExample(t *testing.T) {
	tmpDir := t.TempDir()

	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.Mkdir(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}

	// Do NOT create .env.example

	// Attempt to copy - should fail
	err := copyEnvFile(projectPath)
	if err == nil {
		t.Fatal("Expected error when .env.example is missing, but got nil")
	}
}

// TestEnvTemplateContainsRequiredVariables verifies .env.example has all required variables
func TestEnvTemplateContainsRequiredVariables(t *testing.T) {
	templates := NewProjectTemplates("test-project")
	envContent := templates.EnvTemplate()

	requiredVars := []string{
		"APP_NAME=",
		"APP_PORT=",
		"DB_HOST=",
		"DB_PORT=",
		"DB_USER=",
		"DB_PASSWORD=",
		"DB_NAME=",
		"JWT_SECRET=",
	}

	for _, reqVar := range requiredVars {
		if !contains(envContent, reqVar) {
			t.Errorf(".env.example template missing required variable: %s", reqVar)
		}
	}
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
