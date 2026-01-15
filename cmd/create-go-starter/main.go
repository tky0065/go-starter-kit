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
		"internal/models",
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

	// Run the project creation logic
	if err := run(projectName); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}
}

// run executes the main project creation logic.
// It validates the project name, creates the directory structure,
// generates files, and initializes git.
// Returns an error if any step fails (except git initialization which is non-fatal).
func run(projectName string) error {
	// Validate project name
	if err := validateProjectName(projectName); err != nil {
		return err
	}

	// Display start message
	fmt.Println(Green(fmt.Sprintf("Creating project: %s", projectName)))

	// Use project name as directory path (relative to current directory)
	projectPath := projectName

	// Display progress message
	fmt.Println("📁 Création des répertoires...")

	// Create the project structure
	if err := createProjectStructure(projectPath); err != nil {
		return err
	}

	fmt.Println(Green("✅ Structure terminée"))

	// Generate project files with dynamic context injection
	fmt.Println("📝 Génération des fichiers de base...")

	if err := generateProjectFiles(projectPath, projectName); err != nil {
		return err
	}

	// Display success message
	fmt.Println(Green("✅ Fichiers générés avec succès"))

	// Copy .env.example to .env
	fmt.Println("🔑 Configuration de l'environnement...")
	if err := copyEnvFile(projectPath); err != nil {
		return err
	}

	// Initialize Git repository (AC: 1, 2, 3, 4, 5)
	fmt.Println("🔧 Initialisation du dépôt Git...")
	if err := initGitRepo(projectPath); err != nil {
		// Non-fatal: warn user but continue
		fmt.Println(Red(fmt.Sprintf("⚠️  Avertissement Git: %v", err)))
		fmt.Println("   Vous pouvez initialiser le dépôt manuellement plus tard.")
	} else if isGitAvailable() {
		fmt.Println(Green("✅ Dépôt Git initialisé avec un commit initial"))
	}

	// Display success message with detailed setup instructions
	printSuccessMessage(projectName)

	return nil
}

// printSuccessMessage displays the final success message and setup instructions
func printSuccessMessage(projectName string) {
	fmt.Printf("\n%s\n", Green("════════════════════════════════════════════════════════════════"))
	fmt.Printf("%s\n", Green("🎉 Projet '"+projectName+"' créé avec succès!"))
	fmt.Printf("%s\n\n", Green("════════════════════════════════════════════════════════════════"))

	fmt.Println("📋 Prochaines étapes - Configuration initiale:")
	fmt.Println()

	fmt.Println(Green("OPTION 1: Configuration automatique (Recommandé) 🚀"))
	fmt.Println("  cd " + projectName)
	fmt.Println("  ./setup.sh")
	fmt.Println()

	fmt.Println(Green("OPTION 2: Configuration manuelle"))
	fmt.Println()
	fmt.Println("1️⃣  Naviguer vers le projet:")
	fmt.Println("    cd " + projectName)
	fmt.Println()

	fmt.Println("2️⃣  Configurer PostgreSQL (choisir une option):")
	fmt.Println()
	fmt.Println("    Option A - Docker (Recommandé):")
	fmt.Println("    docker run -d --name postgres \\")
	fmt.Println("      -e POSTGRES_DB=" + projectName + " \\")
	fmt.Println("      -e POSTGRES_PASSWORD=postgres \\")
	fmt.Println("      -p 5432:5432 \\")
	fmt.Println("      postgres:16-alpine")
	fmt.Println()
	fmt.Println("    Option B - PostgreSQL local:")
	fmt.Println("    # macOS: brew install postgresql && brew services start postgresql")
	fmt.Println("    # Linux: sudo apt install postgresql && sudo systemctl start postgresql")
	fmt.Println("    createdb " + projectName)
	fmt.Println()

	fmt.Println("3️⃣  Générer le JWT secret (OBLIGATOIRE):")
	fmt.Println("    openssl rand -base64 32")
	fmt.Println()
	fmt.Println("    Puis éditer .env et ajouter:")
	fmt.Println("    JWT_SECRET=<le_secret_généré>")
	fmt.Println()

	fmt.Println("4️⃣  Lancer l'application:")
	fmt.Println("    make run")
	fmt.Println()

	fmt.Println("5️⃣  Vérifier l'installation:")
	fmt.Println("    curl http://localhost:8080/health")
	fmt.Println("    # Devrait retourner: {\"status\":\"ok\"}")
	fmt.Println()

	fmt.Println(Green("📚 Documentation complète:"))
	fmt.Println("   - Guide rapide: " + projectName + "/docs/quick-start.md")
	fmt.Println("   - README:       " + projectName + "/README.md")
	fmt.Println()

	fmt.Println(Green("⚠️  IMPORTANT:"))
	fmt.Println("   • PostgreSQL DOIT être démarré avant de lancer l'application")
	fmt.Println("   • JWT_SECRET DOIT être configuré dans .env")
	fmt.Println("   • Le fichier .env a été créé automatiquement depuis .env.example")
	fmt.Println()

	fmt.Println(Green("✨ Bon développement avec " + projectName + "!"))
}
