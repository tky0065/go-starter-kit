package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	//"regexp" // Removed as regex is not used directly anymore
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary once for all tests
	const binaryName = "test-create-go-starter"
	cmd := exec.Command("go", "build", "-o", binaryName, ".")
	cmd.Stderr = os.Stderr // Ensure build errors are visible
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build test binary: %v\n", err)
		os.Exit(1)
	}
	binaryPath = "./" + binaryName

	// Run all tests
	code := m.Run()

	// Clean up the binary
	os.Remove(binaryPath)
	os.Exit(code)
}

// TestHelpFlag tests that --help flag displays usage information
func TestHelpFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "--help")
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
	cmd := exec.Command(binaryPath, "-h")
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
	cmd := exec.Command(binaryPath)
	output, err := cmd.CombinedOutput()

	// Should exit with error code
	if err == nil {
		t.Error("Expected non-zero exit code for missing project name")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Project name is required") { // Removed "Error: " prefix
		t.Errorf("Expected error message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Usage:") {
		t.Errorf("Expected usage message after error, got: %s", outputStr)
	}
}

// TestValidProjectName tests that valid project name is accepted
func TestValidProjectName(t *testing.T) {
	// Create a unique test project name
	testProjectName := "my-test-project-" + t.Name()
	defer os.RemoveAll(testProjectName) // Clean up the created directory

	cmd := exec.Command(binaryPath, testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected successful execution, got error: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Creating project: "+testProjectName) {
		t.Errorf("Expected project creation message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Creating directories...") { // Updated to English
		t.Errorf("Expected directory creation message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Structure created") { // Updated to English
		t.Errorf("Expected completion message, got: %s", outputStr)
	}
}

// TestInvalidProjectName tests that invalid project names are rejected
func TestInvalidProjectName(t *testing.T) {
	invalidNames := []string{
		"../evil-path",
		"my/project",
		"my project",
		"_invalid",
		// Note: "-invalid" is not tested here as it's interpreted as a flag by the flag parser
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, name)
			output, err := cmd.CombinedOutput()

			// Should exit with error code
			if err == nil {
				t.Errorf("Expected non-zero exit code for invalid name %q, got success", name)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, "invalid module name") && !strings.Contains(outputStr, "invalid project name") {
				t.Errorf("Expected 'invalid module name' error for %q, got: %s", name, outputStr)
			}
		})
	}
}

// TestColorOutput tests that ANSI color codes are present in output
func TestColorOutput(t *testing.T) {
	// Test success color (green)
	testProjectName := "test-app-" + t.Name()
	defer os.RemoveAll(testProjectName) // Clean up the created directory

	// Create a dummy .env.example file to prevent copyEnvFile from failing
	// The test creates the project directory, so we create the dummy file inside it after creation.
	cmd := exec.Command(binaryPath, testProjectName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected error: %v\nOutput: %s", err, string(output))
	}

	// Check for ANSI green color code
	if !strings.Contains(string(output), "\033[32m") {
		t.Error("Expected green ANSI color code in success output")
	}

	// Test error color (red)
	cmd = exec.Command(binaryPath)
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
	err = run(projectName, DefaultTemplate)
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
	err := run("../invalid-path", DefaultTemplate)
	if err == nil {
		t.Error("Expected error for invalid project name, got nil")
	}
	if !strings.Contains(err.Error(), "invalid module name") && !strings.Contains(err.Error(), "invalid project name") {
		t.Errorf("Expected 'invalid module name' error, got: %v", err)
	}
}

// TestRunFunctionWithEmptyName tests that run() fails with empty project name
func TestRunFunctionWithEmptyName(t *testing.T) {
	err := run("", DefaultTemplate)
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
	err = run(projectName, DefaultTemplate)
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

// TestValidateTemplateValid tests validateTemplate with valid templates (AC: 1, 2)
func TestValidateTemplateValid(t *testing.T) {
	tests := []struct {
		name     string
		template string
	}{
		{"minimal template", "minimal"},
		{"full template", "full"},
		{"graphql template", "graphql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplate(tt.template)
			if err != nil {
				t.Errorf("validateTemplate(%q) returned error: %v, want nil", tt.template, err)
			}
		})
	}
}

// TestValidateTemplateInvalid tests validateTemplate with invalid templates (AC: 3)
func TestValidateTemplateInvalid(t *testing.T) {
	tests := []struct {
		name     string
		template string
	}{
		{"empty template", ""},
		{"unknown template", "unknown"},
		{"invalid case", "FULL"},
		{"invalid spacing", " full"},
		{"invalid template", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplate(tt.template)
			if err == nil {
				t.Errorf("validateTemplate(%q) returned nil, want error", tt.template)
			}
			if !strings.Contains(err.Error(), "invalid template") {
				t.Errorf("validateTemplate(%q) error = %v, want error containing 'invalid template'", tt.template, err)
			}
			if !strings.Contains(err.Error(), "minimal, full, graphql") {
				t.Errorf("validateTemplate(%q) error = %v, want error listing valid options", tt.template, err)
			}
		})
	}
}

// TestTemplateDefaultValue tests that default template is "full" (AC: 2)
func TestTemplateDefaultValue(t *testing.T) {
	if DefaultTemplate != "full" {
		t.Errorf("DefaultTemplate = %q, want %q", DefaultTemplate, "full")
	}
}

// TestValidTemplatesContains tests that ValidTemplates contains expected values
func TestValidTemplatesContains(t *testing.T) {
	expected := []string{"minimal", "full", "graphql"}
	if len(ValidTemplates) != len(expected) {
		t.Errorf("ValidTemplates has %d elements, want %d", len(ValidTemplates), len(expected))
	}
	for _, exp := range expected {
		found := false
		for _, valid := range ValidTemplates {
			if valid == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ValidTemplates does not contain %q", exp)
		}
	}
}

// TestTemplateFlagParsing tests that --template flag is parsed correctly (AC: 1)
func TestTemplateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantInOutput   string
		wantNoErr      bool
		cleanupProject string
	}{
		{
			name:           "minimal template flag",
			args:           []string{"--template=minimal", "test-proj-minimal"},
			wantInOutput:   "template: minimal",
			wantNoErr:      true, // Minimal template is now implemented (Story 6.2)
			cleanupProject: "test-proj-minimal",
		},
		{
			name:           "full template flag",
			args:           []string{"--template=full", "test-proj-full"},
			wantInOutput:   "template: full",
			wantNoErr:      true,
			cleanupProject: "test-proj-full",
		},
		{
			name:           "graphql template flag",
			args:           []string{"--template=graphql", "test-proj-graphql"},
			wantInOutput:   "template: graphql",
			wantNoErr:      false, // Should return error because graphql is not implemented yet
			cleanupProject: "test-proj-graphql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.cleanupProject) // Moved inside t.Run

			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.wantNoErr { // Expect no error
				if err != nil {
					t.Errorf("Expected no error, got: %v\nOutput: %s", err, string(output))
				}
				if !strings.Contains(string(output), tt.wantInOutput) {
					t.Errorf("Output should contain %q, got: %s", tt.wantInOutput, string(output))
				}
			} else { // Expect error
				if err == nil {
					t.Errorf("Expected error for template '%s', got success", tt.args[0])
				}
				if !strings.Contains(string(output), "not yet implemented") { // Specific error message for unimplemented templates
					t.Errorf("Expected 'not yet implemented' error for %q, got: %s", tt.args[0], string(output))
				}
			}
		})
	}
}

// TestTemplateFlagDefault tests that default template is used when not specified (AC: 2)
func TestTemplateFlagDefault(t *testing.T) {
	testProjectName := "test-default-template"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected successful execution, got error: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "template: full") {
		t.Errorf("Expected default template 'full' in output, got: %s", string(output))
	}
}

// TestInvalidTemplateFlagError tests that invalid template shows error (AC: 3)
func TestInvalidTemplateFlagError(t *testing.T) {
	cmd := exec.Command(binaryPath, "--template=invalid", "test-proj-invalid")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for invalid template")
		os.RemoveAll("test-proj-invalid")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "invalid template") {
		t.Errorf("Expected 'invalid template' in error, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "minimal, full, graphql") {
		t.Errorf("Expected valid options in error, got: %s", outputStr)
	}
}

// TestHelpShowsTemplateFlag tests that --help shows template documentation (AC: 4)
func TestHelpShowsTemplateFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected exit code 0 for --help, but got error: %v", err)
	}

	outputStr := string(output)

	// Check main flag line (relaxed check)
	if !strings.Contains(outputStr, "-template string") || !strings.Contains(outputStr, "Template type to generate") {
		t.Errorf("Expected template flag definition in help, got: %s", outputStr)
	}

	// Check templates section header
	if !strings.Contains(outputStr, "Templates:") {
		t.Errorf("Expected 'Templates:' section in help, got: %s", outputStr)
	}

	// Check each template option description
	if !strings.Contains(outputStr, "  minimal   Basic REST API with Swagger (no authentication)") {
		t.Errorf("Expected 'minimal' template description in help, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "  full      Complete API with JWT auth, user management, and Swagger (default)") {
		t.Errorf("Expected 'full' template description in help, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "  graphql   GraphQL API with gqlgen and GraphQL Playground") {
		t.Errorf("Expected 'graphql' template description in help, got: %s", outputStr)
	}
}
