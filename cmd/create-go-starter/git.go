package main

import (
	"fmt"
	"os/exec"
)

// isGitAvailable checks if git is installed and available in the system PATH.
// Returns true if git is available, false otherwise.
func isGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// initGitRepo initializes a git repository in the given project path,
// stages all files, and creates an initial commit.
// If git is not available, it prints a warning and returns nil (non-fatal).
// Returns an error only if git is available but the operation fails.
func initGitRepo(projectPath string) error {
	if !isGitAvailable() {
		fmt.Println("⚠️  Git n'est pas installé. Initialisation Git ignorée.")
		fmt.Println("   Vous pouvez initialiser le dépôt manuellement plus tard avec:")
		fmt.Println("   cd " + projectPath + " && git init && git add . && git commit -m \"Initial commit\"")
		return nil
	}

	// git init
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w (output: %s)", err, string(output))
	}

	// git add .
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stage files: %w (output: %s)", err, string(output))
	}

	// git commit -m "Initial commit from go-starter-kit"
	// Use --allow-empty in case there are no files (edge case)
	cmd = exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit from go-starter-kit")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create initial commit: %w (output: %s)", err, string(output))
	}

	return nil
}
