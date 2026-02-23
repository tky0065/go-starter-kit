package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tky0065/go-starter-kit/cmd/create-go-starter/tui"
	"github.com/tky0065/go-starter-kit/pkg/utils"
)

// ANSI color codes
const (
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
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

// ValidDatabases contains the list of SUPPORTED database types (Story 7.1-7.3)
// Note: mongodb is in backlog (Story 7.4) and not yet implemented
var ValidDatabases = []string{DatabasePostgres, DatabaseMySQL, DatabaseSQLite}

// DefaultDatabase is the default database type when not specified
const DefaultDatabase = DatabasePostgres

// Observability level constants define the available observability modes.
// - ObservabilityNone: No observability (default) — preserves current behavior
// - ObservabilityBasic: Enhanced /health endpoint only (no Prometheus)
// - ObservabilityAdvanced: Full Prometheus metrics endpoint + middleware
const (
	ObservabilityNone     = "none"
	ObservabilityBasic    = "basic"
	ObservabilityAdvanced = "advanced"
)

// ValidObservabilityLevels contains the list of valid observability levels
var ValidObservabilityLevels = []string{ObservabilityNone, ObservabilityBasic, ObservabilityAdvanced}

// DefaultObservabilityLevel is the default observability level when not specified
const DefaultObservabilityLevel = ObservabilityNone

// Green returns the string wrapped in green ANSI code
func Green(msg string) string {
	return ColorGreen + msg + ColorReset
}

// Red returns the string wrapped in red ANSI code
func Red(msg string) string {
	return ColorRed + msg + ColorReset
}

// Yellow returns the string wrapped in yellow ANSI code
func Yellow(msg string) string {
	return ColorYellow + msg + ColorReset
}

// isTTY checks if stdout is connected to a terminal (TTY).
// Returns false in CI/CD environments, piped output, or redirected output.
func isTTY() bool {
	return tui.IsTTY()
}

// runInteractiveTUI launches the Bubble Tea interactive mode.
// Converts main package InteractiveDefaults to tui.InteractiveDefaults.
// After the TUI exits (and the alt screen buffer is restored), it prints
// the detailed success message and generation stats on the normal terminal.
func runInteractiveTUI(defaults InteractiveDefaults) error {
	// Convert main package defaults to tui package defaults
	tuiDefaults := tui.InteractiveDefaults{
		ProjectName:   defaults.ProjectName,
		Template:      defaults.Template,
		Database:      defaults.Database,
		Observability: defaults.Observability,
	}

	// Create a generator function that wraps the existing run() function
	// This allows the TUI to call the generation logic without circular imports
	generatorFunc := func(projectName, template, database, observability string, progressCallback func(current, total int)) error {
		// Call runWithCallback with quiet=true to suppress stdout output during TUI mode
		// The TUI has its own rendering via viewGenerating() and viewDone()
		return runWithCallback(projectName, template, database, observability, progressCallback, true)
	}

	result, err := tui.RunInteractiveTUI(tuiDefaults, generatorFunc)
	if err != nil {
		return err
	}

	// After TUI exits, the alternate screen buffer is restored.
	// Now print the detailed success message on the normal terminal
	// so users get the same post-generation instructions as CLI mode.
	printSuccessMessage(result.ProjectName, result.Database)

	return nil
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

// validateObservabilityLevel checks if the observability level is valid.
// Valid levels are: none, basic, advanced
func validateObservabilityLevel(level string) error {
	for _, valid := range ValidObservabilityLevels {
		if level == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid observability level '%s': valid options are: %s", level, strings.Join(ValidObservabilityLevels, ", "))
}

// validateDatabase checks if the database type is valid and supported.
// Valid databases are: postgres, mysql, sqlite
// Note: mongodb is in backlog (Story 7.4) and will return a helpful error
func validateDatabase(database string) error {
	// Check if it's a supported database
	for _, valid := range ValidDatabases {
		if database == valid {
			return nil
		}
	}

	// Special message for mongodb (in backlog)
	if database == "mongodb" {
		return fmt.Errorf("database 'mongodb' is not yet supported (planned for future release).\nSupported databases: %s", strings.Join(ValidDatabases, ", "))
	}

	// Generic error for other invalid databases
	return fmt.Errorf("invalid database '%s'.\nSupported databases: %s", database, strings.Join(ValidDatabases, ", "))
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
	// Check for subcommands before parsing flags
	if len(os.Args) > 1 && os.Args[1] == "add-model" {
		if err := runAddModel(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
			os.Exit(1)
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		// Handle doctor --help
		for _, arg := range os.Args[2:] {
			if arg == "-help" || arg == "--help" || arg == "-h" {
				fmt.Fprintf(os.Stderr, "Usage: create-go-starter doctor\n\n")
				fmt.Fprintf(os.Stderr, "Run environment diagnostics to verify your development setup.\n\n")
				fmt.Fprintf(os.Stderr, "Checks performed:\n")
				fmt.Fprintf(os.Stderr, "  Go       Verifies Go is installed and meets minimum version (1.21+)\n")
				fmt.Fprintf(os.Stderr, "  Git      Verifies Git is installed and available in PATH\n")
				fmt.Fprintf(os.Stderr, "  Docker   Verifies Docker binary is installed and daemon is running\n\n")
				fmt.Fprintf(os.Stderr, "Exit codes:\n")
				fmt.Fprintf(os.Stderr, "  0   All checks passed\n")
				fmt.Fprintf(os.Stderr, "  1   One or more checks failed\n")
				os.Exit(0)
			}
		}
		exitCode := runDoctor()
		os.Exit(exitCode)
	}

	// Parse flags manually to allow flags in any position
	// This is necessary because we want to support both:
	// - create-go-starter -database sqlite my-project
	// - create-go-starter my-project -database sqlite
	help := false
	interactive := false
	dryRun := false
	var template string = DefaultTemplate
	var database string = DefaultDatabase
	var observability string = DefaultObservabilityLevel
	var projectName string

	// Manually parse arguments to allow flags in any position
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-help" || arg == "--help" || arg == "-h" {
			help = true
		} else if arg == "-interactive" || arg == "--interactive" || arg == "-i" {
			interactive = true
		} else if arg == "-dry-run" || arg == "--dry-run" || arg == "-n" {
			dryRun = true
		} else if strings.HasPrefix(arg, "-template=") || strings.HasPrefix(arg, "--template=") || strings.HasPrefix(arg, "-t=") {
			// Handle -template=value, --template=value, -t=value syntax
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				template = parts[1]
			}
		} else if strings.HasPrefix(arg, "-database=") || strings.HasPrefix(arg, "--database=") || strings.HasPrefix(arg, "-d=") {
			// Handle -database=value, --database=value, -d=value syntax
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				database = parts[1]
			}
		} else if strings.HasPrefix(arg, "-observability=") || strings.HasPrefix(arg, "--observability=") || strings.HasPrefix(arg, "-o=") {
			// Handle -observability=value, --observability=value, -o=value syntax
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				observability = parts[1]
			}
		} else if (arg == "-template" || arg == "--template" || arg == "-t") && i+1 < len(args) {
			template = args[i+1]
			i++ // Skip next arg since we consumed it
		} else if (arg == "-database" || arg == "--database" || arg == "-d") && i+1 < len(args) {
			database = args[i+1]
			i++ // Skip next arg since we consumed it
		} else if (arg == "-observability" || arg == "--observability" || arg == "-o") && i+1 < len(args) {
			observability = args[i+1]
			i++ // Skip next arg since we consumed it
		} else if !strings.HasPrefix(arg, "-") && projectName == "" {
			// First non-flag argument is the project name
			projectName = arg
		} else if strings.HasPrefix(arg, "-") {
			// Unknown flag — print clear error message (AC: #9)
			fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("unknown flag: %s", arg)))
			os.Exit(1)
		}
	}

	// Define usage function
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: create-go-starter [options] <project-name>\n")
		fmt.Fprintf(os.Stderr, "       Flags can be placed before or after the project name\n\n")
		fmt.Fprintf(os.Stderr, "Subcommands:\n")
		fmt.Fprintf(os.Stderr, "  doctor\n")
		fmt.Fprintf(os.Stderr, "        Check your development environment (Go, Git, Docker)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -i, --interactive\n")
		fmt.Fprintf(os.Stderr, "        Launch interactive mode (guided step-by-step configuration)\n")
		fmt.Fprintf(os.Stderr, "  -n, --dry-run\n")
		fmt.Fprintf(os.Stderr, "        Preview files that would be created without generating them\n")
		fmt.Fprintf(os.Stderr, "  -d, --database string\n")
		fmt.Fprintf(os.Stderr, "        Database type to use (default \"postgres\")\n")
		fmt.Fprintf(os.Stderr, "  -t, --template string\n")
		fmt.Fprintf(os.Stderr, "        Template type to generate (default \"full\")\n")
		fmt.Fprintf(os.Stderr, "  -o, --observability string\n")
		fmt.Fprintf(os.Stderr, "        Observability level: none|basic|advanced (default \"none\")\n")
		fmt.Fprintf(os.Stderr, "  -h, --help\n")
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
		fmt.Fprintf(os.Stderr, "\nObservability:\n")
		fmt.Fprintf(os.Stderr, "  %-9s No observability (default) - current behavior preserved\n", ObservabilityNone)
		fmt.Fprintf(os.Stderr, "  %-9s Enhanced /health endpoint only (no Prometheus)\n", ObservabilityBasic)
		fmt.Fprintf(os.Stderr, "  %-9s Full Prometheus metrics endpoint + HTTP metrics middleware\n", ObservabilityAdvanced)
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter -d sqlite my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter my-project -t minimal\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter -d mysql -t minimal my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter my-project --observability=advanced\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter -d sqlite -t minimal -o none my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter --database=mysql --template=full my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter -i\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter -n my-project\n")
		fmt.Fprintf(os.Stderr, "  create-go-starter doctor\n")
	}

	if help {
		usage()
		os.Exit(0)
	}

	// Reject conflicting flags: --interactive and --dry-run cannot be used together
	if interactive && dryRun {
		fmt.Fprintln(os.Stderr, Red("--interactive and --dry-run cannot be used together"))
		os.Exit(1)
	}

	if interactive {
		defaults := InteractiveDefaults{
			ProjectName:   projectName,
			Template:      template,
			Database:      database,
			Observability: observability,
		}

		// Use Bubble Tea TUI if TTY is available and NO_COLOR is not set
		// Falls back to text mode for CI/CD, piped output, or when NO_COLOR is set
		if isTTY() && os.Getenv("NO_COLOR") == "" {
			// Modern Bubble Tea TUI experience
			if err := runInteractiveTUI(defaults); err != nil {
				fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
				os.Exit(1)
			}
		} else {
			// Fallback to text-based prompts (CI/CD compatible)
			if err := runInteractiveMode(defaults); err != nil {
				fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
				os.Exit(1)
			}
		}
		return
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
		usage() // Display usage on invalid project name
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

	// Validate observability level
	if err := validateObservabilityLevel(observability); err != nil {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
		os.Exit(1)
	}

	// Validate observability + template combination
	if observability == ObservabilityAdvanced && template != TemplateFull {
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("--observability=advanced is only supported with --template=full (got --template=%s)", template)))
		os.Exit(1)
	}

	// If dry-run flag is set, preview files without creating them (AC: #1-#6)
	if dryRun {
		if err := runDryRun(projectName, template, database, observability); err != nil {
			fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
			os.Exit(1)
		}
		return
	}

	// Run the project creation logic
	if err := run(projectName, template, database, observability); err != nil {
		// Changed to not include "Error: " prefix as Red() function will color the message itself.
		fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
		os.Exit(1)
	}
}

// run executes the main project creation logic without progress callback.
// This is the legacy function for CLI usage (non-interactive mode).
// For interactive TUI mode with progress updates, use runWithCallback() instead.
func run(projectName, template, database, observabilityLevel string) error {
	return runWithCallback(projectName, template, database, observabilityLevel, nil, false)
}

// runWithCallback executes the main project creation logic with optional progress callback.
// It validates the project name, creates the directory structure,
// generates files, and initializes git.
// The progressCallback is called for each file generated with (current, total).
// When quiet is true, all stdout output (progress messages, success message, stats) is suppressed.
// This is used when the TUI is active since it has its own rendering.
// Returns an error if any step fails (except git initialization which is non-fatal).
func runWithCallback(projectName, template, database, observabilityLevel string, progressCallback func(current, total int), quiet bool) error {
	stats := NewGenerationStats()

	// Display start message with template info
	if !quiet {
		fmt.Println(Green(fmt.Sprintf("Creating project: %s (template: %s, database: %s, observability: %s)", projectName, template, database, observabilityLevel)))
	}

	// Validate project name again to ensure safety when run() is called directly (e.g. in tests)
	if err := utils.ValidateGoModuleName(projectName); err != nil {
		return err
	}

	// Use project name as directory path (relative to current directory)
	projectPath := projectName

	// Display progress message
	if !quiet {
		fmt.Println("📁 Creating directories...")
	}

	// Create the project structure
	stats.StartStep("Creating directories")
	if err := createProjectStructure(projectPath, template); err != nil {
		return err
	}
	stats.EndStep("Creating directories")

	if !quiet {
		fmt.Println(Green("✅ Structure created"))
	}

	// Get the full file list before generating (needed for progress bar total count)
	files := getFilesForTemplate(projectPath, projectName, template, database, observabilityLevel)
	var pb *ProgressBar
	if !quiet {
		pb = NewProgressBar(len(files), 30)
	}

	// Generate project files with progress reporting
	if !quiet {
		fmt.Println("📝 Generating core files...")
	}
	stats.StartStep("Generating files")
	if err := writeFiles(files, func(current, total int) {
		// Update CLI progress bar (only when not in quiet/TUI mode)
		if pb != nil {
			pb.Update(current)
		}
		// Forward progress to TUI callback if provided (Story 10.7 AC#3)
		if progressCallback != nil {
			progressCallback(current, total)
		}
	}); err != nil {
		return err
	}
	if pb != nil {
		pb.Complete()
	}

	// Make setup.sh executable (handled previously inside template-specific functions)
	setupPath := filepath.Join(projectPath, "setup.sh")
	if _, statErr := os.Stat(setupPath); statErr == nil {
		if err := os.Chmod(setupPath, 0755); err != nil {
			return fmt.Errorf("failed to make setup.sh executable: %w", err)
		}
	}

	// Accumulate file stats (size from disk after writing)
	for _, f := range files {
		stats.AddFile(f.Path)
	}
	stats.EndStep("Generating files")

	// Display success message
	if !quiet {
		fmt.Println(Green("✅ Files generated successfully"))
	}

	// Copy .env.example to .env
	if !quiet {
		fmt.Println("🔑 Configuring environment...")
	}
	if err := copyEnvFile(projectPath); err != nil {
		return err
	}

	// Initialize Git repository (AC: 1, 2, 3, 4, 5)
	if !quiet {
		fmt.Println("🔧 Initializing Git repository...")
	}
	stats.StartStep("Git initialization")
	if err := initGitRepo(projectPath); err != nil {
		// Non-fatal: warn user but continue
		if !quiet {
			fmt.Println(Red(fmt.Sprintf("⚠️  Git warning: %v", err)))
			fmt.Println("   You can initialize the repository manually later.")
		}
	} else if !quiet && isGitAvailable() {
		fmt.Println(Green("✅ Git repository initialized with initial commit"))
	}
	stats.EndStep("Git initialization")

	// Display success message with detailed setup instructions
	if !quiet {
		printSuccessMessage(projectName, database)
	}

	// Display generation statistics (AC: #2, #3)
	if !quiet {
		stats.Display()
	}

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
	fmt.Println("    # Should return: {\"status\":\"alive\",\"service\":\"" + projectName + "\",\"timestamp\":\"...\"}") // Changed to English
	fmt.Println()
	fmt.Println("    # Advanced health checks:")
	fmt.Println("    curl http://localhost:8080/health/liveness    # Liveness probe (K8s)")
	fmt.Println("    curl http://localhost:8080/health/readiness   # Readiness probe (K8s)")
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
