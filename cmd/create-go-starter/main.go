package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// ANSI color codes
const (
	ColorGreen = "\033[32m"
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

// Directory permissions for created folders
const defaultDirPerm os.FileMode = 0755

// Valid project name pattern (alphanumeric, hyphens, underscores)
var validProjectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// Green returns the string wrapped in green ANSI code
func Green(msg string) string {
	return ColorGreen + msg + ColorReset
}

// Red returns the string wrapped in red ANSI code
func Red(msg string) string {
	return ColorRed + msg + ColorReset
}

// validateProjectName checks if the project name contains only valid characters.
// Valid names must start with alphanumeric and contain only alphanumeric, hyphens, or underscores.
func validateProjectName(name string) error {
	if !validProjectNamePattern.MatchString(name) {
		return fmt.Errorf("invalid project name '%s': must start with a letter or number and contain only letters, numbers, hyphens, or underscores", name)
	}
	return nil
}

// createProjectStructure creates the hexagonal architecture directory structure.
// It returns an error if the directory already exists or if creation fails.
func createProjectStructure(projectPath string) error {
	// Check if project directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory %s already exists. Please choose a different name or remove the existing directory", projectPath)
	}

	// Define the directory structure (Hexagonal Architecture Lite)
	directories := []string{
		"cmd",
		"internal/adapters/http",
		"internal/adapters/middleware",
		"internal/domain",
		"internal/interfaces",
		"internal/infrastructure/database",
		"internal/infrastructure/server",
		"pkg/config",
		"pkg/logger",
		"deployments",
	}

	// Create the project root directory
	if err := os.Mkdir(projectPath, defaultDirPerm); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create each subdirectory
	for _, dir := range directories {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, defaultDirPerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// copyEnvFile copies the generated .env.example to .env if .env doesn't already exist.
func copyEnvFile(projectPath string) error {
	envExamplePath := filepath.Join(projectPath, ".env.example")
	envPath := filepath.Join(projectPath, ".env")

	// Check if .env.example exists
	if _, err := os.Stat(envExamplePath); os.IsNotExist(err) {
		return fmt.Errorf(".env.example not found: %w", err)
	}

	// Check if .env already exists
	if _, err := os.Stat(envPath); err == nil {
		// .env already exists, skip copying
		return nil
	}

	// Read .env.example
	content, err := os.ReadFile(envExamplePath)
	if err != nil {
		return fmt.Errorf("failed to read .env.example: %w", err)
	}

	// Write to .env
	if err := os.WriteFile(envPath, content, 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	return nil
}

func main() {
	// Parse flags
	help := flag.Bool("help", false, "Show help message")
	flag.BoolVar(help, "h", false, "Show help message (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: create-go-starter [options] <project-name>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, Red("Error: Project name is required"))
		flag.Usage()
		os.Exit(1)
	}

	projectName := args[0]

	// Validate project name
	if err := validateProjectName(projectName); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}

	// Display start message
	fmt.Println(Green(fmt.Sprintf("Creating project: %s", projectName)))

	// Use project name as directory path (relative to current directory)
	projectPath := projectName

	// Display progress message
	fmt.Println("📁 Création des répertoires...")

	// Create the project structure
	if err := createProjectStructure(projectPath); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}

	fmt.Println(Green("✅ Structure terminée"))

	// Generate project files with dynamic context injection
	fmt.Println("📝 Génération des fichiers de base...")

	if err := generateProjectFiles(projectPath, projectName); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}

	// Display success message
	fmt.Println(Green("✅ Fichiers générés avec succès"))

	// Copy .env.example to .env
	fmt.Println("🔑 Configuration de l'environnement...")
	if err := copyEnvFile(projectPath); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}

	// Display success message
	fmt.Printf("\n🎉 Projet '%s' créé avec succès!\n", Green(projectName))
	fmt.Printf("\nProchaines étapes:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  go mod download\n")
	fmt.Printf("  make run\n")
}
