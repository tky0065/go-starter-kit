package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestHelpFlag tests that --help flag displays usage information
func TestHelpFlag(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	// Test --help flag
	cmd = exec.Command("./test-binary", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected exit code 0 for --help, but got error: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Usage:") {
		t.Errorf("Expected usage message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "create-go-starter") {
		t.Errorf("Expected binary name in output, got: %s", outputStr)
	}
}

// TestHelpFlagShorthand tests that -h flag displays usage information
func TestHelpFlagShorthand(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	cmd = exec.Command("./test-binary", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected exit code 0 for -h, but got error: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Usage:") {
		t.Errorf("Expected usage message for -h, got: %s", outputStr)
	}
}

// TestMissingProjectName tests that missing project name returns error
func TestMissingProjectName(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	cmd = exec.Command("./test-binary")
	output, err := cmd.CombinedOutput()

	// Should exit with error code
	if err == nil {
		t.Error("Expected non-zero exit code for missing project name")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Error: Project name is required") {
		t.Errorf("Expected error message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Usage:") {
		t.Errorf("Expected usage message after error, got: %s", outputStr)
	}
}

// TestValidProjectName tests that valid project name is accepted
func TestValidProjectName(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	// Create a unique test project name
	testProjectName := "my-test-project-" + t.Name()
	defer os.RemoveAll(testProjectName) // Clean up the created directory

	cmd = exec.Command("./test-binary", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected successful execution, got error: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Creating project: "+testProjectName) {
		t.Errorf("Expected project creation message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Création des répertoires") {
		t.Errorf("Expected directory creation message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Structure terminée") {
		t.Errorf("Expected completion message, got: %s", outputStr)
	}
}

// TestInvalidProjectName tests that invalid project names are rejected
func TestInvalidProjectName(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	invalidNames := []string{
		"../evil-path",
		"my/project",
		"my project",
		"_invalid",
		// Note: "-invalid" is not tested here as it's interpreted as a flag by the flag parser
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command("./test-binary", name)
			output, err := cmd.CombinedOutput()

			// Should exit with error code
			if err == nil {
				t.Errorf("Expected non-zero exit code for invalid name %q, got success", name)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, "invalid project name") {
				t.Errorf("Expected 'invalid project name' error for %q, got: %s", name, outputStr)
			}
		})
	}
}

// TestColorOutput tests that ANSI color codes are present in output
func TestColorOutput(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	// Test success color (green)
	testProjectName := "test-app-" + t.Name()
	defer os.RemoveAll(testProjectName) // Clean up the created directory

	// Create a dummy .env.example file to prevent copyEnvFile from failing
	// The test creates the project directory, so we create the dummy file inside it after creation.
	cmd = exec.Command("./test-binary", testProjectName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected error: %v\nOutput: %s", err, string(output))
	}

	// Check for ANSI green color code
	if !strings.Contains(string(output), "\033[32m") {
		t.Error("Expected green ANSI color code in success output")
	}

	// Test error color (red)
	cmd = exec.Command("./test-binary")
	output, err = cmd.CombinedOutput()

	// We expect an error here (missing project name)
	if err == nil {
		t.Error("Expected error for missing project name")
	}

	// Check for ANSI red color code
	if !strings.Contains(string(output), "\033[31m") {
		t.Error("Expected red ANSI color code in error output")
	}
}

// TestRunFunction tests the run() function directly
func TestRunFunction(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	projectName := "test-run-project"

	// Change to temp directory so project is created there
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Run the main logic
	err = run(projectName)
	if err != nil {
		t.Errorf("run() returned error: %v", err)
	}

	// Verify project was created
	projectPath := filepath.Join(tmpDir, projectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Error("Project directory was not created")
	}

	// Verify essential files exist
	essentialFiles := []string{
		"go.mod",
		"cmd/main.go",
		"Makefile",
		".env.example",
		".env",
	}
	for _, file := range essentialFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", file)
		}
	}
}

// TestRunFunctionWithInvalidName tests that run() fails with invalid project name
func TestRunFunctionWithInvalidName(t *testing.T) {
	err := run("../invalid-path")
	if err == nil {
		t.Error("Expected error for invalid project name, got nil")
	}
	if !strings.Contains(err.Error(), "invalid project name") {
		t.Errorf("Expected 'invalid project name' error, got: %v", err)
	}
}

// TestRunFunctionWithEmptyName tests that run() fails with empty project name
func TestRunFunctionWithEmptyName(t *testing.T) {
	err := run("")
	if err == nil {
		t.Error("Expected error for empty project name, got nil")
	}
}

// TestRunFunctionWithExistingDirectory tests run() with already existing directory
func TestRunFunctionWithExistingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "existing-project"

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create existing directory
	existingDir := filepath.Join(tmpDir, projectName)
	if err := os.Mkdir(existingDir, 0755); err != nil {
		t.Fatalf("Failed to create existing directory: %v", err)
	}

	// Run should fail because directory exists
	err = run(projectName)
	if err == nil {
		t.Error("Expected error for existing directory, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

// TestPrintSuccessMessage tests that printSuccessMessage doesn't panic
func TestPrintSuccessMessage(t *testing.T) {
	// This test just verifies the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printSuccessMessage() panicked: %v", r)
		}
	}()

	printSuccessMessage("test-project")
}

// TestCoverageThreshold verifies that code coverage meets the 70% threshold (AC#1)
func TestCoverageThreshold(t *testing.T) {
	// Run tests with coverage and capture output
	cmd := exec.Command("go", "test", "-cover", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run tests with coverage: %v\nOutput: %s", err, string(output))
	}

	// Parse coverage percentage from output
	// Expected format: "coverage: XX.X% of statements"
	outputStr := string(output)
	if !strings.Contains(outputStr, "coverage:") {
		t.Fatalf("Coverage information not found in output: %s", outputStr)
	}

	// Extract coverage percentage
	lines := strings.Split(outputStr, "\n")
	var coverageStr string
	for _, line := range lines {
		if strings.Contains(line, "coverage:") && strings.Contains(line, "% of statements") {
			// Extract percentage using string manipulation
			start := strings.Index(line, "coverage: ") + len("coverage: ")
			end := strings.Index(line, "% of statements")
			if start < end && start >= 0 && end >= 0 {
				coverageStr = line[start:end]
				break
			}
		}
	}

	if coverageStr == "" {
		t.Fatalf("Could not extract coverage percentage from: %s", outputStr)
	}

	// Parse coverage as float
	var coverage float64
	if _, err := fmt.Sscanf(coverageStr, "%f", &coverage); err != nil {
		t.Fatalf("Failed to parse coverage percentage '%s': %v", coverageStr, err)
	}

	// AC#1: Coverage must be >= 70%
	const requiredCoverage = 70.0
	if coverage < requiredCoverage {
		t.Errorf("Coverage %.1f%% is below required threshold of %.1f%%", coverage, requiredCoverage)
	}

	t.Logf("✅ Coverage: %.1f%% (>= %.1f%% required)", coverage, requiredCoverage)
}
