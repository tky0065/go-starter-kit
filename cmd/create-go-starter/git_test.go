package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestIsGitAvailable tests the git availability detection function
func TestIsGitAvailable(t *testing.T) {
	// This test assumes git is installed on the test machine
	// If git is not installed, this test will verify the function returns false
	available := isGitAvailable()

	// Check if git is actually available using exec.LookPath
	_, err := exec.LookPath("git")
	expectedAvailable := err == nil

	if available != expectedAvailable {
		t.Errorf("isGitAvailable() = %v, want %v", available, expectedAvailable)
	}
}

// TestInitGitRepoSuccess tests successful git repository initialization
func TestInitGitRepoSuccess(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed, skipping test")
	}

	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a test file so we have something to commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Initialize git repo
	err := initGitRepo(tmpDir)
	if err != nil {
		t.Errorf("initGitRepo() error = %v, want nil", err)
	}

	// Verify .git directory exists
	gitDir := filepath.Join(tmpDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error(".git directory was not created")
	}

	// Verify initial commit was made
	cmd := exec.Command("git", "log", "--oneline", "-1")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Failed to get git log: %v", err)
	}

	if !strings.Contains(string(output), "Initial commit from go-starter-kit") {
		t.Errorf("Expected initial commit message, got: %s", string(output))
	}
}

// TestInitGitRepoWithFiles tests that all files are included in initial commit
func TestInitGitRepoWithFiles(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed, skipping test")
	}

	tmpDir := t.TempDir()

	// Create multiple test files
	files := []string{"file1.txt", "file2.go", "README.md"}
	for _, f := range files {
		filePath := filepath.Join(tmpDir, f)
		if err := os.WriteFile(filePath, []byte("content for "+f), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", f, err)
		}
	}

	// Create a subdirectory with a file
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	subFile := filepath.Join(subDir, "subfile.txt")
	if err := os.WriteFile(subFile, []byte("subdir content"), 0644); err != nil {
		t.Fatalf("Failed to create subdir file: %v", err)
	}

	// Initialize git repo
	if err := initGitRepo(tmpDir); err != nil {
		t.Fatalf("initGitRepo() error = %v", err)
	}

	// Verify all files are tracked
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to list git files: %v", err)
	}

	trackedFiles := string(output)
	for _, f := range files {
		if !strings.Contains(trackedFiles, f) {
			t.Errorf("File %s was not tracked by git", f)
		}
	}
	if !strings.Contains(trackedFiles, "subdir/subfile.txt") {
		t.Error("Subdirectory file was not tracked by git")
	}
}

// TestInitGitRepoInvalidPath tests behavior with an invalid path
func TestInitGitRepoInvalidPath(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed, skipping test")
	}

	err := initGitRepo("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

// TestInitGitRepoEmptyDirectory tests git init on empty directory
func TestInitGitRepoEmptyDirectory(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed, skipping test")
	}

	tmpDir := t.TempDir()

	// Initialize git repo on empty directory
	// This should still work - git allows empty initial commits with --allow-empty
	// But our implementation stages files, so with no files, git add . is a no-op
	// The commit should still succeed with the message

	err := initGitRepo(tmpDir)
	if err != nil {
		t.Errorf("initGitRepo() on empty directory error = %v", err)
	}

	// Verify .git directory exists
	gitDir := filepath.Join(tmpDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error(".git directory was not created for empty directory")
	}
}

// TestIsGitAvailableReturnsBool tests that isGitAvailable returns a boolean
func TestIsGitAvailableReturnsBool(t *testing.T) {
	result := isGitAvailable()
	// The function should return either true or false, not panic
	if result != true && result != false {
		t.Error("isGitAvailable should return a boolean")
	}
}

// TestInitGitRepoCommitMessage tests the exact commit message format
func TestInitGitRepoCommitMessage(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed, skipping test")
	}

	tmpDir := t.TempDir()

	// Create a file to commit
	testFile := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test Project"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := initGitRepo(tmpDir); err != nil {
		t.Fatalf("initGitRepo() error = %v", err)
	}

	// Get the exact commit message
	cmd := exec.Command("git", "log", "--format=%s", "-1")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit message: %v", err)
	}

	expectedMessage := "Initial commit from go-starter-kit"
	actualMessage := strings.TrimSpace(string(output))
	if actualMessage != expectedMessage {
		t.Errorf("Commit message = %q, want %q", actualMessage, expectedMessage)
	}
}

// TestE2EGitIntegration is an end-to-end test that verifies the full CLI
// creates a project with git initialization (AC: 1, 2, 3, 5)
func TestE2EGitIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed, skipping E2E test")
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	// Create a unique project name
	projectName := "test-git-e2e-" + t.Name()
	defer os.RemoveAll(projectName)

	// Run the CLI
	cmd = exec.Command("./test-binary", projectName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI execution failed: %v\nOutput: %s", err, string(output))
	}

	// Verify output mentions git initialization
	outputStr := string(output)
	if !strings.Contains(outputStr, "Initializing Git repository") {
		t.Error("Expected git initialization message in output")
	}
	if !strings.Contains(outputStr, "Git repository initialized") {
		t.Error("Expected git success message in output")
	}

	// AC1: Verify .git directory exists
	gitDir := filepath.Join(projectName, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error("AC1 FAILED: .git directory was not created")
	}

	// AC2 & AC3: Verify initial commit exists with correct message
	cmd = exec.Command("git", "log", "--format=%s", "-1")
	cmd.Dir = projectName
	commitOutput, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get git log: %v", err)
	}

	expectedMsg := "Initial commit from go-starter-kit"
	if !strings.Contains(string(commitOutput), expectedMsg) {
		t.Errorf("AC2 FAILED: Expected commit message %q, got %q", expectedMsg, string(commitOutput))
	}

	// AC3: Verify files are tracked
	cmd = exec.Command("git", "ls-files")
	cmd.Dir = projectName
	filesOutput, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to list git files: %v", err)
	}

	// Check that key generated files are tracked
	trackedFiles := string(filesOutput)
	expectedFiles := []string{"go.mod", "Makefile", "README.md", ".env.example"}
	for _, f := range expectedFiles {
		if !strings.Contains(trackedFiles, f) {
			t.Errorf("AC3 FAILED: Expected file %q to be tracked in git", f)
		}
	}

	// AC5: Verify git init happened after file generation (implicit - if files are tracked, it worked)
	// This is verified by the fact that all expected files are in the initial commit
}
