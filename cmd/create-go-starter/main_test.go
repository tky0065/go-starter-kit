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
	err = run(projectName, DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel)
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
	err := run("../invalid-path", DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel)
	if err == nil {
		t.Error("Expected error for invalid project name, got nil")
	}
	if !strings.Contains(err.Error(), "invalid module name") && !strings.Contains(err.Error(), "invalid project name") {
		t.Errorf("Expected 'invalid module name' error, got: %v", err)
	}
}

// TestRunFunctionWithEmptyName tests that run() fails with empty project name
func TestRunFunctionWithEmptyName(t *testing.T) {
	err := run("", DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel)
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
	err = run(projectName, DefaultTemplate, DefaultDatabase, DefaultObservabilityLevel)
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

	printSuccessMessage("test-project", "postgres")
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
			wantNoErr:      true, // GraphQL template is now implemented
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

// TestTemplateFlagBothSyntaxes tests that both --template=minimal and --template minimal work
// This is a regression test for the bug where --template minimal (with space) didn't work
func TestTemplateFlagBothSyntaxes(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantTemplate   string
		cleanupProject string
	}{
		{
			name:           "template flag with equals",
			args:           []string{"--template=minimal", "test-equals-syntax"},
			wantTemplate:   "minimal",
			cleanupProject: "test-equals-syntax",
		},
		{
			name:           "template flag with space",
			args:           []string{"--template", "minimal", "test-space-syntax"},
			wantTemplate:   "minimal",
			cleanupProject: "test-space-syntax",
		},
		{
			name:           "full template with space",
			args:           []string{"--template", "full", "test-full-space"},
			wantTemplate:   "full",
			cleanupProject: "test-full-space",
		},
		{
			name:           "graphql template with space",
			args:           []string{"--template", "graphql", "test-graphql-space"},
			wantTemplate:   "graphql",
			cleanupProject: "test-graphql-space",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.cleanupProject)

			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Errorf("Expected no error, got: %v\nOutput: %s", err, string(output))
				return
			}

			outputStr := string(output)
			expectedMsg := fmt.Sprintf("template: %s", tt.wantTemplate)
			if !strings.Contains(outputStr, expectedMsg) {
				t.Errorf("Expected output to contain %q, got: %s", expectedMsg, outputStr)
			}

			// Verify that the correct template was actually generated
			// For minimal template, pkg/auth should NOT exist
			if tt.wantTemplate == "minimal" {
				authPath := filepath.Join(tt.cleanupProject, "pkg", "auth")
				if _, err := os.Stat(authPath); !os.IsNotExist(err) {
					t.Errorf("Minimal template should not have pkg/auth directory, but it exists")
				}

				// Verify that pkg/config and pkg/logger DO exist
				configPath := filepath.Join(tt.cleanupProject, "pkg", "config")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("Minimal template should have pkg/config directory, but it doesn't exist")
				}

				loggerPath := filepath.Join(tt.cleanupProject, "pkg", "logger")
				if _, err := os.Stat(loggerPath); os.IsNotExist(err) {
					t.Errorf("Minimal template should have pkg/logger directory, but it doesn't exist")
				}
			}

			// For full template, pkg/auth SHOULD exist
			if tt.wantTemplate == "full" {
				authPath := filepath.Join(tt.cleanupProject, "pkg", "auth")
				if _, err := os.Stat(authPath); os.IsNotExist(err) {
					t.Errorf("Full template should have pkg/auth directory, but it doesn't exist")
				}
			}
		})
	}
}

// ============================================================================
// Story 10.5: Short Flag Aliases Tests (AC: #1-#9)
// ============================================================================

// TestShortFlagInteractive tests that -i is equivalent to --interactive (AC: #1)
func TestShortFlagInteractive(t *testing.T) {
	cmd := exec.Command(binaryPath, "-i")
	// Pipe empty stdin to make interactive mode exit quickly (EOF triggers exit)
	cmd.Stdin = strings.NewReader("")
	output, _ := cmd.CombinedOutput()
	// -i should launch interactive mode - verify the flag is recognized (no "unknown flag" error)
	outputStr := string(output)
	if strings.Contains(outputStr, "unknown flag: -i") {
		t.Errorf("Expected -i to be recognized as --interactive, got: %s", outputStr)
	}
	// Should not show "Project name is required" (which would mean interactive flag was not set)
	if strings.Contains(outputStr, "Project name is required") {
		t.Errorf("Expected -i to trigger interactive mode, not require project name: %s", outputStr)
	}
	// Verify the output mentions interactive mode or prompts the user
	// (The mode should output something indicating interactive behavior)
	if strings.Contains(outputStr, "Usage:") {
		t.Errorf("Expected -i to enter interactive mode, not show usage: %s", outputStr)
	}
}

// TestShortFlagTemplate tests that -t minimal is equivalent to --template=minimal (AC: #2)
func TestShortFlagTemplate(t *testing.T) {
	testProjectName := "test-short-template"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-t", "minimal", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -t minimal, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "template: minimal") {
		t.Errorf("Expected 'template: minimal' in output with -t flag, got: %s", outputStr)
	}
}

// TestShortFlagTemplateWithEquals tests that -t=minimal works (AC: #2)
func TestShortFlagTemplateWithEquals(t *testing.T) {
	testProjectName := "test-short-template-eq"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-t=minimal", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -t=minimal, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "template: minimal") {
		t.Errorf("Expected 'template: minimal' in output with -t= flag, got: %s", outputStr)
	}
}

// TestShortFlagDatabase tests that -d sqlite is equivalent to --database=sqlite (AC: #3)
func TestShortFlagDatabase(t *testing.T) {
	testProjectName := "test-short-database"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-d", "sqlite", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -d sqlite, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "database: sqlite") {
		t.Errorf("Expected 'database: sqlite' in output with -d flag, got: %s", outputStr)
	}
}

// TestShortFlagDatabaseWithEquals tests that -d=sqlite works (AC: #3)
func TestShortFlagDatabaseWithEquals(t *testing.T) {
	testProjectName := "test-short-database-eq"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-d=sqlite", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -d=sqlite, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "database: sqlite") {
		t.Errorf("Expected 'database: sqlite' in output with -d= flag, got: %s", outputStr)
	}
}

// TestShortFlagObservability tests that -o basic is equivalent to --observability=basic (AC: #4)
func TestShortFlagObservability(t *testing.T) {
	testProjectName := "test-short-obs"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-o", "basic", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -o basic, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "observability: basic") {
		t.Errorf("Expected 'observability: basic' in output with -o flag, got: %s", outputStr)
	}
}

// TestShortFlagObservabilityWithEquals tests that -o=basic works (AC: #4)
func TestShortFlagObservabilityWithEquals(t *testing.T) {
	testProjectName := "test-short-obs-eq"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-o=basic", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -o=basic, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "observability: basic") {
		t.Errorf("Expected 'observability: basic' in output with -o= flag, got: %s", outputStr)
	}
}

// TestShortFlagDryRun tests that -n is equivalent to --dry-run (AC: #5)
func TestShortFlagDryRun(t *testing.T) {
	testProjectName := "test-short-dryrun"

	cmd := exec.Command(binaryPath, "-n", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -n (dry-run), got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "DRY RUN") && !strings.Contains(outputStr, "dry-run") && !strings.Contains(outputStr, "dry run") {
		t.Errorf("Expected dry-run output with -n flag, got: %s", outputStr)
	}
	// Verify no directory was created
	if _, err := os.Stat(testProjectName); !os.IsNotExist(err) {
		os.RemoveAll(testProjectName)
		t.Errorf("Dry-run should not create directory, but %s was created", testProjectName)
	}
}

// TestShortFlagHelpStillWorks tests that -h still works (AC: #6)
func TestShortFlagHelpStillWorks(t *testing.T) {
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

// TestHelpShowsShortAliases tests that --help displays short flag aliases (AC: #7)
func TestHelpShowsShortAliases(t *testing.T) {
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected exit code 0 for --help, but got error: %v", err)
	}

	outputStr := string(output)

	// Verify short aliases are shown in help
	if !strings.Contains(outputStr, "-d,") {
		t.Errorf("Expected -d, alias for --database in help, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "-t,") {
		t.Errorf("Expected -t, alias for --template in help, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "-o,") {
		t.Errorf("Expected -o, alias for --observability in help, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "-i,") {
		t.Errorf("Expected -i, alias for --interactive in help, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "-n,") {
		t.Errorf("Expected -n, alias for --dry-run in help, got: %s", outputStr)
	}
}

// TestCombinedShortFlags tests that multiple short flags work together (AC: #8)
func TestCombinedShortFlags(t *testing.T) {
	testProjectName := "test-combined-short"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-d", "sqlite", "-t", "minimal", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with combined short flags, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "database: sqlite") {
		t.Errorf("Expected 'database: sqlite' in combined short flags output, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "template: minimal") {
		t.Errorf("Expected 'template: minimal' in combined short flags output, got: %s", outputStr)
	}
}

// TestMixedShortAndLongFlags tests that short and long flags can be combined (AC: #8)
func TestMixedShortAndLongFlags(t *testing.T) {
	testProjectName := "test-mixed-flags"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "-d", "sqlite", "--template=full", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with mixed flags, got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "database: sqlite") {
		t.Errorf("Expected 'database: sqlite' in mixed flags output, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "template: full") {
		t.Errorf("Expected 'template: full' in mixed flags output, got: %s", outputStr)
	}
}

// TestUnknownFlagError tests that unknown flags show a clear error message (AC: #9)
func TestUnknownFlagError(t *testing.T) {
	cmd := exec.Command(binaryPath, "-x", "my-project")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected non-zero exit code for unknown flag -x")
		os.RemoveAll("my-project")
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "unknown flag: -x") {
		t.Errorf("Expected 'unknown flag: -x' in error output, got: %s", outputStr)
	}
}

// TestUnknownLongFlagError tests that unknown long flags also show a clear error (AC: #9)
func TestUnknownLongFlagError(t *testing.T) {
	cmd := exec.Command(binaryPath, "--unknown-flag", "my-project")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected non-zero exit code for unknown flag --unknown-flag")
		os.RemoveAll("my-project")
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "unknown flag: --unknown-flag") {
		t.Errorf("Expected 'unknown flag: --unknown-flag' in error output, got: %s", outputStr)
	}
}

// TestUnknownFlagWithValueDoesNotConsume tests that unknown flag -x followed by a value
// does not silently consume the value as projectName (AC: #9)
func TestUnknownFlagWithValueDoesNotConsume(t *testing.T) {
	cmd := exec.Command(binaryPath, "-x", "sqlite", "my-project")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected non-zero exit code for unknown flag -x followed by value")
		os.RemoveAll("my-project")
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "unknown flag: -x") {
		t.Errorf("Expected 'unknown flag: -x' in error output, got: %s", outputStr)
	}
}

// TestLongFlagsBackwardCompatibility tests that long flags still work after adding short aliases (AC: #8)
func TestLongFlagsBackwardCompatibility(t *testing.T) {
	testProjectName := "test-backward-compat"
	defer os.RemoveAll(testProjectName)

	cmd := exec.Command(binaryPath, "--database=postgres", "--template=full", testProjectName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with long flags (backward compat), got error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "database: postgres") {
		t.Errorf("Expected 'database: postgres' with long flags, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "template: full") {
		t.Errorf("Expected 'template: full' with long flags, got: %s", outputStr)
	}
}
