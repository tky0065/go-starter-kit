package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tky0065/go-starter-kit/pkg/utils"
)

// ANSI color codes
const (
	ColorGreen = "\033[32m"
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

// Directory permissions for created folders
const defaultDirPerm os.FileMode = 0755

// Template constants define the available project generation templates.
// - TemplateMinimal: Basic REST API with Swagger (no authentication)
// - TemplateFull: Complete hexagonal architecture with JWT auth, user management (default)
// - TemplateGraphQL: GraphQL API with gqlgen and GraphQL Playground (not yet implemented)
const (
	TemplateMinimal = "minimal"
	TemplateFull    = "full"
	TemplateGraphQL = "graphql"
)

// Template descriptions (in English for consistency with code)
const (
	TemplateMinimalDesc = "Basic REST API with Swagger (no authentication)"
	TemplateFullDesc    = "Complete API with JWT auth, user management, and Swagger (default)"
	TemplateGraphQLDesc = "GraphQL API with gqlgen and GraphQL Playground"
)

// ValidTemplates contains the list of valid template types
var ValidTemplates = []string{TemplateMinimal, TemplateFull, TemplateGraphQL}

// DefaultTemplate is the default template type when not specified
const DefaultTemplate = TemplateFull

// Database type constants
const (
	DatabasePostgres = "postgres"
	DatabaseMySQL    = "mysql"
	DatabaseSQLite   = "sqlite"
	DatabaseMongoDB  = "mongodb"
)

// ValidDatabases contains the list of supported database types
var ValidDatabases = []string{DatabasePostgres, DatabaseMySQL, DatabaseSQLite, DatabaseMongoDB}

// DefaultDatabase is the default database type when not specified
const DefaultDatabase = DatabasePostgres

// Green returns the string wrapped in green ANSI code
func Green(msg string) string {
	return ColorGreen + msg + ColorReset
}

// Red returns the string wrapped in red ANSI code
func Red(msg string) string {
	return ColorRed + msg + ColorReset
}

// validateTemplate checks if the template type is valid.
// Valid templates are: minimal, full, graphql
func validateTemplate(template string) error {
	for _, valid := range ValidTemplates {
		if template == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid template '%s': valid options are: %s", template, strings.Join(ValidTemplates, ", "))
}

// validateDatabase checks if the database type is valid.
// Valid databases are: postgres, mysql, sqlite, mongodb
func validateDatabase(database string) error {
	for _, valid := range ValidDatabases {
		if database == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid database '%s': valid options are: %s", database, strings.Join(ValidDatabases, ", "))
}

// createProjectStructure creates the hexagonal architecture directory structure.
// It returns an error if the directory already exists or if creation fails.
// The template parameter determines which directories to create.
func createProjectStructure(projectPath, template string) error {
	// Check if project directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory %s already exists. Please choose a different name or remove the existing directory", projectPath)
	}

	// Get directories for the specified template
	directories := getDirectoriesForTemplate(template)

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
	// Parse flags manually to allow flags in any position
	// This is necessary because we want to support both:
	// - create-go-starter -database sqlite my-project
	// - create-go-starter my-project -database sqlite
	help := false
	var template string = DefaultTemplate
	var database string = DefaultDatabase
	var projectName string

	// Manually parse arguments to allow flags in any position
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-help" || arg == "--help" || arg == "-h" {
			help = true
		} else if strings.HasPrefix(arg, "-template=") || strings.HasPrefix(arg, "--template=") {
			// Handle -template=value syntax
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				template = parts[1]
			}
		} else if strings.HasPrefix(arg, "-database=") || strings.HasPrefix(arg, "--database=") {
			// Handle -database=value syntax
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				database = parts[1]
			}
		} else if (arg == "-template" || arg == "--template") && i+1 < len(args) {
			template = args[i+1]
			i++ // Skip next arg since we consumed it
		} else if (arg == "-database" || arg == "--database") && i+1 < len(args) {
			database = args[i+1]
			i++ // Skip next arg since we consumed it
		} else if !strings.HasPrefix(arg, "-") && projectName == "" {
			// First non-flag argument is the project name
			projectName = arg
		}
	}

	// Define usage function
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: create-go-starter [options] <project-name>\n")
		fmt.Fprintf(os.Stderr, "       Flags can be placed before or after the project name\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -database string\n")
		fmt.Fprintf(os.Stderr, "        Database type to use (default \"postgres\")\n")
		fmt.Fprintf(os.Stderr, "  -template string\n")
		fmt.Fprintf(os.Stderr, "        Template type to generate (default \"full\")\n")
		fmt.Fprintf(os.Stderr, "  -h, -help\n")
		fmt.Fprintf(os.Stderr, "        Show help message\n")
		fmt.Fprintf(os.Stderr, "\nTemplates:\n")
		fmt.Fprintf(os.Stderr, "  %-9s %s\n", TemplateMinimal, TemplateMinimalDesc)
		fmt.Fprintf(os.Stderr, "  %-9s %s\n", TemplateFull, TemplateFullDesc)
		fmt.Fprintf(os.Stderr, "  %-9s %s\n", TemplateGraphQL, TemplateGraphQLDesc)
		fmt.Fprintf(os.Stderr, "\nDatabases:\n")
		fmt.Fprintf(os.Stderr, "  %-9s PostgreSQL (default) - Production-ready, advanced features\n", DatabasePostgres)
		fmt.Fprintf(os.Stderr, "  %-9s MySQL/MariaDB - Wide compatibility, shared hosting\n", DatabaseMySQL)
		fmt.Fprintf(os.Stderr, "  %-9s SQLite - Quick prototyping, embedded apps\n", DatabaseSQLite)
		fmt.Fprintf(os.Stderr, "  %-9s MongoDB - NoSQL, document-oriented\n", DatabaseMongoDB)
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter -database sqlite my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter my-project -template minimal\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter -database mysql -template minimal my-project\n")
	}

	if help {
		usage()
		os.Exit(0)
	}

	if projectName == "" {
		// Changed to not include "Error: " here, as Red() function will add color to the message itself.
		fmt.Fprintln(os.Stderr, Red("Project name is required"))
		usage()
		os.Exit(1)
	}

	// Validate project name using the shared utility
	if err := utils.ValidateGoModuleName(projectName); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
		flag.Usage() // Display usage on invalid project name
		os.Exit(1)
	}

	// Validate template
	if err := validateTemplate(template); err != nil {
		// Changed to not include "Error: " prefix as Red() function will color the message itself.
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
		os.Exit(1)
	}

	// Validate database
	if err := validateDatabase(database); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
		os.Exit(1)
	}

	// Run the project creation logic
	if err := run(projectName, template, database); err != nil {
		// Changed to not include "Error: " prefix as Red() function will color the message itself.
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
		os.Exit(1)
	}
}

// run executes the main project creation logic.
// It validates the project name, creates the directory structure,
// generates files, and initializes git.
// Returns an error if any step fails (except git initialization which is non-fatal).
func run(projectName, template, database string) error {
	// Display start message with template info
	fmt.Println(Green(fmt.Sprintf("Creating project: %s (template: %s, database: %s)", projectName, template, database)))

	// Validate project name again to ensure safety when run() is called directly (e.g. in tests)
	if err := utils.ValidateGoModuleName(projectName); err != nil {
		return err
	}

	// Use project name as directory path (relative to current directory)
	projectPath := projectName

	// Display progress message
	fmt.Println("📁 Creating directories...") // Changed to English

	// Create the project structure
	if err := createProjectStructure(projectPath, template); err != nil {
		return err
	}

	fmt.Println(Green("✅ Structure created")) // Changed to English

	// Generate project files with dynamic context injection
	fmt.Println("📝 Generating core files...") // Changed to English

	if err := generateProjectFiles(projectPath, projectName, template, database); err != nil {
		return err
	}

	// Display success message
	fmt.Println(Green("✅ Files generated successfully")) // Changed to English

	// Copy .env.example to .env
	fmt.Println("🔑 Configuring environment...") // Changed to English
	if err := copyEnvFile(projectPath); err != nil {
		return err
	}

	// Initialize Git repository (AC: 1, 2, 3, 4, 5)
	fmt.Println("🔧 Initializing Git repository...") // Changed to English
	if err := initGitRepo(projectPath); err != nil {
		// Non-fatal: warn user but continue
		fmt.Println(Red(fmt.Sprintf("⚠️  Git warning: %v", err)))           // Changed to English
		fmt.Println("   You can initialize the repository manually later.") // Changed to English
	} else if isGitAvailable() {
		fmt.Println(Green("✅ Git repository initialized with initial commit")) // Changed to English
	}

	// Display success message with detailed setup instructions
	printSuccessMessage(projectName, database)

	return nil
}

// getDatabaseDisplayName returns the human-readable name for a database type
func getDatabaseDisplayName(database string) string {
	switch database {
	case "mysql":
		return "MySQL"
	case "sqlite":
		return "SQLite"
	case "mongodb":
		return "MongoDB"
	case "postgres", "":
		return "PostgreSQL"
	default:
		return "PostgreSQL"
	}
}

// getDockerCommand returns the Docker setup command for the specified database
func getDockerCommand(projectName, database string) string {
	switch database {
	case "mysql":
		return `docker run -d --name mysql \
      -e MYSQL_ROOT_PASSWORD=secret \
      -e MYSQL_DATABASE=` + projectName + ` \
      -p 3306:3306 \
      mysql:8.0`
	case "sqlite":
		return "# SQLite uses a local file, no Docker needed!\n    # The database will be created automatically at: ./" + projectName + ".db"
	case "mongodb":
		return `docker run -d --name mongodb \
      -e MONGO_INITDB_ROOT_USERNAME=admin \
      -e MONGO_INITDB_ROOT_PASSWORD=admin \
      -p 27017:27017 \
      mongo:latest`
	case "postgres", "":
		return `docker run -d --name postgres \
      -e POSTGRES_DB=` + projectName + ` \
      -e POSTGRES_PASSWORD=postgres \
      -p 5432:5432 \
      postgres:16-alpine`
	default:
		return `docker run -d --name postgres \
      -e POSTGRES_DB=` + projectName + ` \
      -e POSTGRES_PASSWORD=postgres \
      -p 5432:5432 \
      postgres:16-alpine`
	}
}

// getLocalSetupCommand returns the local setup instructions for the specified database
func getLocalSetupCommand(projectName, database string) string {
	switch database {
	case "mysql":
		return `# macOS: brew install mysql && brew services start mysql
    # Linux: sudo apt install mysql-server && sudo systemctl start mysql
    # Windows: Download from https://dev.mysql.com/downloads/mysql/
    mysql -u root -p"secret" -e "CREATE DATABASE ` + projectName + `;"`
	case "sqlite":
		return "# SQLite needs no separate setup - uses local file automatically"
	case "mongodb":
		return `# macOS: brew install mongodb-community && brew services start mongodb-community
    # Linux: sudo apt install mongodb && sudo systemctl start mongodb
    # Use MongoDB Compass for GUI: https://www.mongodb.com/products/compass`
	case "postgres", "":
		return `# macOS: brew install postgresql && brew services start postgresql
    # Linux: sudo apt install postgresql && sudo systemctl start postgresql
    createdb ` + projectName
	default:
		return `# macOS: brew install postgresql && brew services start postgresql
    # Linux: sudo apt install postgresql && sudo systemctl start postgresql
    createdb ` + projectName
	}
}

// printSuccessMessage displays the final success message and setup instructions
func printSuccessMessage(projectName, database string) {
	fmt.Printf("\n%s\n", Green("════════════════════════════════════════════════════════════════"))
	fmt.Printf("%s\n", Green(fmt.Sprintf("🎉 Project '%s' created successfully!", projectName))) // Changed to English
	fmt.Printf("%s\n\n", Green("════════════════════════════════════════════════════════════════"))

	fmt.Println("📋 Next steps - Initial setup:") // Changed to English
	fmt.Println()

	fmt.Println(Green("OPTION 1: Automatic setup (Recommended) 🚀")) // Changed to English
	fmt.Println("  cd " + projectName)
	fmt.Println("  ./setup.sh")
	fmt.Println()

	fmt.Println(Green("OPTION 2: Manual setup")) // Changed to English
	fmt.Println()
	fmt.Println("1️⃣  Navigate to your project:") // Changed to English
	fmt.Println("    cd " + projectName)
	fmt.Println()

	// Database-specific setup instructions
	fmt.Printf("2️⃣  Configure %s (choose one option):\n", getDatabaseDisplayName(database)) // Changed to English
	fmt.Println()
	fmt.Println("    Option A - Docker (Recommended):") // Changed to English
	dockerCmd := getDockerCommand(projectName, database)
	fmt.Println(dockerCmd)
	fmt.Println()
	fmt.Println("    Option B - Local setup:") // Changed to English
	localCmd := getLocalSetupCommand(projectName, database)
	fmt.Println(localCmd)
	fmt.Println()

	fmt.Println("3️⃣  Generate JWT secret (REQUIRED):") // Changed to English
	fmt.Println("    openssl rand -base64 32")
	fmt.Println()
	fmt.Println("    Then edit .env and add:")       // Changed to English
	fmt.Println("    JWT_SECRET=<generated_secret>") // Changed to English
	fmt.Println()

	fmt.Println("4️⃣  Start the application:") // Changed to English
	fmt.Println("    make run")
	fmt.Println()

	fmt.Println("5️⃣  Verify installation:") // Changed to English
	fmt.Println("    curl http://localhost:8080/health")
	fmt.Println("    # Should return: {\"status\":\"ok\"}") // Changed to English
	fmt.Println()

	fmt.Println(Green("📚 Full documentation:"))                                    // Changed to English
	fmt.Println("   - Quick Start Guide: " + projectName + "/docs/quick-start.md") // Changed to English
	fmt.Println("   - README:            " + projectName + "/README.md")           // Changed to English
	fmt.Println()

	fmt.Println(Green("⚠️  IMPORTANT:"))                                                                       // Changed to English
	fmt.Printf("   • %s MUST be started before launching the application\n", getDatabaseDisplayName(database)) // Changed to English
	fmt.Println("   • JWT_SECRET MUST be configured in .env")                                                  // Changed to English
	fmt.Println("   • The .env file was automatically created from .env.example")                              // Changed to English
	fmt.Println()

	fmt.Println(Green("✨ Happy developing with " + projectName + "!")) // Changed to English
}
