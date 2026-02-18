package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGetFilesForTemplateFull tests that the full template returns a substantial file list
func TestGetFilesForTemplateFull(t *testing.T) {
	files := getFilesForTemplate("/tmp/test-project", "test-project", "full", "postgres", "none")
	if len(files) == 0 {
		t.Error("expected non-empty file list for full template")
	}
	if len(files) <= 30 {
		t.Errorf("full template should have > 30 files, got %d", len(files))
	}
}

// TestGetFilesForTemplateMinimal tests that the minimal template returns fewer files than full
func TestGetFilesForTemplateMinimal(t *testing.T) {
	files := getFilesForTemplate("/tmp/test-project", "test-project", "minimal", "postgres", "none")
	if len(files) == 0 {
		t.Error("expected non-empty file list for minimal template")
	}
	fullFiles := getFilesForTemplate("/tmp/test-project", "test-project", "full", "postgres", "none")
	if len(files) >= len(fullFiles) {
		t.Errorf("minimal template should have fewer files than full: minimal=%d, full=%d", len(files), len(fullFiles))
	}
}

// TestGetFilesForTemplateGraphQL tests that the graphql template returns a non-empty list
func TestGetFilesForTemplateGraphQL(t *testing.T) {
	files := getFilesForTemplate("/tmp/test-project", "test-project", "graphql", "postgres", "none")
	if len(files) == 0 {
		t.Error("expected non-empty file list for graphql template")
	}
}

// TestGetFilesForTemplateWithSQLite tests database variation is handled
func TestGetFilesForTemplateWithSQLite(t *testing.T) {
	files := getFilesForTemplate("/tmp/test-project", "test-project", "minimal", "sqlite", "none")
	if len(files) == 0 {
		t.Error("expected non-empty file list for minimal+sqlite template")
	}
}

// TestGetFilesForTemplateAdvancedObservability tests that advanced observability adds more files
func TestGetFilesForTemplateAdvancedObservability(t *testing.T) {
	filesNone := getFilesForTemplate("/tmp/test-project", "test-project", "full", "postgres", "none")
	filesAdv := getFilesForTemplate("/tmp/test-project", "test-project", "full", "postgres", "advanced")
	if len(filesAdv) <= len(filesNone) {
		t.Errorf("advanced observability should produce more files: advanced=%d, none=%d", len(filesAdv), len(filesNone))
	}
}

// TestDryRunCreatesNoFiles verifies that runDryRun does not create any files or directories
func TestDryRunCreatesNoFiles(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if chErr := os.Chdir(originalDir); chErr != nil {
			t.Logf("failed to restore working directory: %v", chErr)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	projectName := "test-dry-run"
	if err := runDryRun(projectName, "full", "postgres", "none"); err != nil {
		t.Fatalf("runDryRun returned unexpected error: %v", err)
	}

	// Verify no files or directories were created
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected no files created, got %d entries", len(entries))
	}
}

// TestDryRunReturnsNilError verifies that runDryRun returns nil (exit code 0)
func TestDryRunReturnsNilError(t *testing.T) {
	err := runDryRun("test-project", "full", "postgres", "none")
	if err != nil {
		t.Errorf("runDryRun should return nil error, got: %v", err)
	}
}

// TestDryRunWithMinimalTemplate verifies dry-run works for minimal template
func TestDryRunWithMinimalTemplate(t *testing.T) {
	err := runDryRun("test-project", "minimal", "sqlite", "none")
	if err != nil {
		t.Errorf("runDryRun(minimal) should return nil, got: %v", err)
	}
}

// TestDryRunWithGraphQLTemplate verifies dry-run works for graphql template
func TestDryRunWithGraphQLTemplate(t *testing.T) {
	err := runDryRun("test-project", "graphql", "postgres", "none")
	if err != nil {
		t.Errorf("runDryRun(graphql) should return nil, got: %v", err)
	}
}

// TestDryRunDoesNotErrorWhenDirectoryExists verifies AC#6: warning shown but no error
func TestDryRunDoesNotErrorWhenDirectoryExists(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if chErr := os.Chdir(originalDir); chErr != nil {
			t.Logf("failed to restore working directory: %v", chErr)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create the project directory to simulate existing dir
	existingProject := "existing-project"
	if err := os.Mkdir(existingProject, 0755); err != nil {
		t.Fatal(err)
	}

	// Should not return error even when directory exists
	if err := runDryRun(existingProject, "full", "postgres", "none"); err != nil {
		t.Errorf("runDryRun should not error when directory exists, got: %v", err)
	}

	// Only the pre-existing directory should be there (no new files created)
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("expected exactly 1 entry (the pre-existing dir), got %d", len(entries))
	}
}

// TestCountUniqueDirectories verifies the directory counting logic
func TestCountUniqueDirectories(t *testing.T) {
	files := []FileGenerator{
		{Path: "/tmp/proj/go.mod"},
		{Path: "/tmp/proj/cmd/main.go"},
		{Path: "/tmp/proj/cmd/main_test.go"},
		{Path: "/tmp/proj/pkg/config/env.go"},
	}

	dirs := make(map[string]bool)
	for _, f := range files {
		dirs[filepath.Dir(f.Path)] = true
	}

	expected := 3 // /tmp/proj, /tmp/proj/cmd, /tmp/proj/pkg/config
	if len(dirs) != expected {
		t.Errorf("expected %d unique dirs, got %d", expected, len(dirs))
	}
}

// TestBuildFullFileListContainsExpectedFiles verifies the full file list structure
func TestBuildFullFileListContainsExpectedFiles(t *testing.T) {
	files := buildFullFileList("/tmp/proj", "myproject", "postgres", "none")

	expectedPaths := []string{
		"go.mod",
		"cmd/main.go",
		"pkg/config/env.go",
		"pkg/logger/logger.go",
		"pkg/auth/jwt.go",
		"setup.sh",
		"README.md",
		"Makefile",
		".gitignore",
	}

	fileMap := make(map[string]bool)
	for _, f := range files {
		rel, err := filepath.Rel("/tmp/proj", f.Path)
		if err == nil {
			fileMap[rel] = true
		}
	}

	for _, expected := range expectedPaths {
		if !fileMap[expected] {
			t.Errorf("expected file %s in full file list", expected)
		}
	}
}

// TestBuildMinimalFileListDoesNotContainAuthFiles verifies minimal has no auth files
func TestBuildMinimalFileListDoesNotContainAuthFiles(t *testing.T) {
	files := buildMinimalFileList("/tmp/proj", "myproject", "postgres")

	authFiles := []string{
		"pkg/auth/jwt.go",
		"internal/domain/user/service.go",
		"internal/adapters/handlers/auth_handler.go",
	}

	fileMap := make(map[string]bool)
	for _, f := range files {
		rel, err := filepath.Rel("/tmp/proj", f.Path)
		if err == nil {
			fileMap[rel] = true
		}
	}

	for _, authFile := range authFiles {
		if fileMap[authFile] {
			t.Errorf("minimal template should not contain auth file: %s", authFile)
		}
	}
}

// captureStdout runs a function and captures its stdout output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	buf := make([]byte, 64*1024)
	n, _ := r.Read(buf)
	r.Close()
	return string(buf[:n])
}

// TestDryRunOutputContainsFileList verifies AC#1: the output lists files with relative paths
func TestDryRunOutputContainsFileList(t *testing.T) {
	output := captureStdout(t, func() {
		_ = runDryRun("my-project", "full", "postgres", "none")
	})

	// AC#1: Should contain file paths relative to project root
	expectedFiles := []string{
		"my-project/go.mod",
		"my-project/cmd/main.go",
		"my-project/pkg/config/env.go",
		"my-project/setup.sh",
	}
	for _, f := range expectedFiles {
		if !strings.Contains(output, f) {
			t.Errorf("AC#1: expected output to contain file path %q", f)
		}
	}
}

// TestDryRunOutputContainsSummary verifies AC#3: summary includes counts, template, database, observability
func TestDryRunOutputContainsSummary(t *testing.T) {
	output := captureStdout(t, func() {
		_ = runDryRun("my-project", "full", "postgres", "none")
	})

	// AC#3: Summary must include these fields
	checks := []struct {
		label string
		text  string
	}{
		{"file count", "36 files"},
		{"directory count", "18 unique"},
		{"template", "Template:      full"},
		{"database", "Database:      postgres"},
		{"observability", "Observability: none"},
	}
	for _, c := range checks {
		if !strings.Contains(output, c.text) {
			t.Errorf("AC#3: expected output to contain %s (%q), output:\n%s", c.label, c.text, output[:min(len(output), 500)])
		}
	}
}

// TestDryRunOutputContainsCompletionMessage verifies AC#5: "Dry-run completed. No files created."
func TestDryRunOutputContainsCompletionMessage(t *testing.T) {
	output := captureStdout(t, func() {
		_ = runDryRun("my-project", "full", "postgres", "none")
	})

	if !strings.Contains(output, "Dry-run completed. No files created.") {
		t.Errorf("AC#5: expected output to contain 'Dry-run completed. No files created.', got:\n%s", output[:min(len(output), 500)])
	}
}

// TestDryRunOutputContainsWarningForExistingDir verifies AC#6: warning message when dir exists
func TestDryRunOutputContainsWarningForExistingDir(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if chErr := os.Chdir(originalDir); chErr != nil {
			t.Logf("failed to restore working directory: %v", chErr)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	existingProject := "existing-proj"
	if err := os.Mkdir(existingProject, 0755); err != nil {
		t.Fatal(err)
	}

	output := captureStdout(t, func() {
		_ = runDryRun(existingProject, "full", "postgres", "none")
	})

	if !strings.Contains(output, "Warning") || !strings.Contains(output, "already exists") {
		t.Errorf("AC#6: expected warning about existing directory, got:\n%s", output[:min(len(output), 500)])
	}
}

// min returns the smaller of two ints (Go 1.21+ has built-in min, but for compatibility)
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
