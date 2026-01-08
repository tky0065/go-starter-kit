package main

import (
	"os"
	"os/exec"
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

	cmd = exec.Command("./test-binary", testProjectName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
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
