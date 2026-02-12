package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// detectProjectContext verifies we're in a go-starter-kit project and detects the database type.
type projectContext struct {
	ModulePath string // e.g., "github.com/user/myproject"
	Database   string // "postgres", "mysql", or "sqlite"
	ProjectDir string // absolute path to project root
}

// RelationConfig holds the relation flags for model generation.
type RelationConfig struct {
	BelongsTo string // Parent model name for belongs-to relation (e.g., "Todo")
	HasMany   string // Child model name for has-many relation (e.g., "Comment")
}

// validateParentModel checks that the specified parent/related model file exists.
func validateParentModel(projectPath, modelName string) error {
	modelFile := filepath.Join(projectPath, "internal", "models", strings.ToLower(modelName)+".go")
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		return fmt.Errorf("related model '%s' not found at %s. Create it first with: add-model %s --fields ...",
			modelName, modelFile, modelName)
	}
	return nil
}

// detectProject checks if the current directory is a valid go-starter-kit project.
func detectProject(dir string) (*projectContext, error) {
	// Check go.mod exists
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("go.mod not found: not a Go project (run this command from your project root)")
	}

	// Check internal/models/ directory exists
	modelsDir := filepath.Join(dir, "internal", "models")
	if _, err := os.Stat(modelsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("internal/models/ not found: not a go-starter-kit project")
	}

	// Read go.mod for module path and validation
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}

	contentStr := string(content)

	// Validate this is a go-starter-kit project by checking for key dependencies
	requiredDeps := []string{
		"github.com/gofiber/fiber",
		"go.uber.org/fx",
		"gorm.io/gorm",
	}
	for _, dep := range requiredDeps {
		if !strings.Contains(contentStr, dep) {
			return nil, fmt.Errorf("not a go-starter-kit project: missing required dependency %q in go.mod\n(go-starter-kit projects use fiber, fx, and gorm)", dep)
		}
	}

	modulePath := extractModulePath(contentStr)
	if modulePath == "" {
		return nil, fmt.Errorf("could not extract module path from go.mod")
	}

	database, err := detectDatabase(contentStr)
	if err != nil {
		return nil, err
	}

	return &projectContext{
		ModulePath: modulePath,
		Database:   database,
		ProjectDir: dir,
	}, nil
}

// extractModulePath extracts the module path from go.mod content.
func extractModulePath(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}

// detectDatabase detects the database driver from go.mod content.
func detectDatabase(content string) (string, error) {
	switch {
	case strings.Contains(content, "gorm.io/driver/postgres"):
		return "postgres", nil
	case strings.Contains(content, "gorm.io/driver/mysql"):
		return "mysql", nil
	case strings.Contains(content, "gorm.io/driver/sqlite"):
		return "sqlite", nil
	default:
		return "", fmt.Errorf("database driver not found in go.mod: expected gorm.io/driver/(postgres|mysql|sqlite)")
	}
}

// runAddModel handles the `add-model` subcommand.
func runAddModel(args []string) error {
	// Parse add-model specific flags
	var fieldsStr string
	var showHelp bool
	var skipConfirm bool
	var isPublic bool
	var modelName string
	var belongsTo string
	var hasMany string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-help" || arg == "--help" || arg == "-h" {
			showHelp = true
		} else if arg == "-yes" || arg == "--yes" || arg == "-y" {
			skipConfirm = true
		} else if arg == "-public" || arg == "--public" {
			isPublic = true
		} else if strings.HasPrefix(arg, "-fields=") || strings.HasPrefix(arg, "--fields=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				fieldsStr = parts[1]
			}
		} else if (arg == "-fields" || arg == "--fields") && i+1 < len(args) {
			fieldsStr = args[i+1]
			i++
		} else if strings.HasPrefix(arg, "-belongs-to=") || strings.HasPrefix(arg, "--belongs-to=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				belongsTo = parts[1]
			}
		} else if (arg == "-belongs-to" || arg == "--belongs-to") && i+1 < len(args) {
			belongsTo = args[i+1]
			i++
		} else if strings.HasPrefix(arg, "-has-many=") || strings.HasPrefix(arg, "--has-many=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				hasMany = parts[1]
			}
		} else if (arg == "-has-many" || arg == "--has-many") && i+1 < len(args) {
			hasMany = args[i+1]
			i++
		} else if !strings.HasPrefix(arg, "-") && modelName == "" {
			modelName = arg
		}
	}

	if showHelp {
		printAddModelHelp()
		return nil
	}

	if modelName == "" {
		return fmt.Errorf("model name is required\n\nUsage: create-go-starter add-model <ModelName> --fields \"name:type,...\"\nRun 'create-go-starter add-model --help' for more information")
	}

	// Validate model name
	if err := validateModelName(modelName); err != nil {
		return err
	}

	if fieldsStr == "" {
		return fmt.Errorf("--fields flag is required\n\nUsage: create-go-starter add-model %s --fields \"name:type,...\"", modelName)
	}

	// Parse fields
	fields, err := parseFields(fieldsStr)
	if err != nil {
		return fmt.Errorf("invalid fields: %w", err)
	}

	// Validate relation model names
	if belongsTo != "" {
		if err := validateModelName(belongsTo); err != nil {
			return fmt.Errorf("invalid --belongs-to model name: %w", err)
		}
	}
	if hasMany != "" {
		if err := validateModelName(hasMany); err != nil {
			return fmt.Errorf("invalid --has-many model name: %w", err)
		}
	}

	// Detect project context
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	ctx, err := detectProject(dir)
	if err != nil {
		return err
	}

	// Validate parent model exists for belongs-to
	if belongsTo != "" {
		if err := validateParentModel(ctx.ProjectDir, belongsTo); err != nil {
			return err
		}
	}

	// Validate child model exists for has-many (the parent being created will reference it)
	if hasMany != "" {
		if err := validateParentModel(ctx.ProjectDir, hasMany); err != nil {
			return err
		}
	}

	// Check if model already exists
	modelFileName := strings.ToLower(modelName) + ".go"
	modelFilePath := filepath.Join(ctx.ProjectDir, "internal", "models", modelFileName)
	if _, err := os.Stat(modelFilePath); err == nil {
		return fmt.Errorf("model file already exists: %s", modelFilePath)
	}

	// Build relation config
	relations := &RelationConfig{
		BelongsTo: belongsTo,
		HasMany:   hasMany,
	}

	// Display summary
	printModelSummary(modelName, fields, ctx, relations)

	// Ask for confirmation unless --yes
	if !skipConfirm {
		fmt.Print("\nProceed? (y/n) ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	// Generate all model files
	fmt.Println("\n  Generating model files...")
	if err := generateModelFiles(ctx, modelName, fields, isPublic, relations); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	lowerName := strings.ToLower(modelName)
	fmt.Println(Green(fmt.Sprintf("  Model %s generated successfully!", modelName)))
	fmt.Printf("\n   Files created:\n")
	fmt.Printf("   + internal/models/%s.go\n", lowerName)
	fmt.Printf("   + internal/interfaces/%s_repository.go\n", lowerName)
	fmt.Printf("   + internal/adapters/repository/%s_repository.go\n", lowerName)
	fmt.Printf("   + internal/domain/%s/service.go\n", lowerName)
	fmt.Printf("   + internal/domain/%s/module.go\n", lowerName)
	fmt.Printf("   + internal/adapters/handlers/%s_handler.go\n", lowerName)
	fmt.Printf("   + internal/domain/%s/service_test.go\n", lowerName)
	fmt.Printf("   + internal/adapters/handlers/%s_handler_test.go\n", lowerName)
	fmt.Printf("\n   Files modified:\n")
	fmt.Printf("   ~ internal/infrastructure/database/database.go (AutoMigrate)\n")
	fmt.Printf("   ~ internal/adapters/repository/module.go (fx.Provide)\n")
	fmt.Printf("   ~ internal/adapters/handlers/module.go (fx.Provide)\n")
	fmt.Printf("   ~ internal/adapters/http/routes.go (CRUD routes)\n")
	fmt.Printf("   ~ cmd/main.go (fx.Module)\n")
	if hasMany != "" {
		fmt.Printf("   ~ internal/models/%s.go (added has-many relation)\n", strings.ToLower(hasMany))
	}

	// Remind user to update Swagger docs
	fmt.Println(Yellow("\n  Next steps:"))
	fmt.Println("   1. Run tests: go test ./...")
	fmt.Println("   2. Update Swagger docs: swag init -g cmd/api/main.go -o docs/swagger")
	fmt.Println("   3. Or use Makefile: make swagger")

	return nil
}

// printModelSummary displays what will be generated.
func printModelSummary(modelName string, fields []FieldDefinition, ctx *projectContext, relations *RelationConfig) {
	fmt.Println(Green(fmt.Sprintf("\n  Model: %s", modelName)))
	fmt.Printf("   Database: %s\n", ctx.Database)
	fmt.Printf("   Module:   %s\n\n", ctx.ModulePath)

	fmt.Println("   Fields:")
	for _, f := range fields {
		tags := ""
		if len(f.GORMTags) > 0 {
			tags = " [" + strings.Join(f.GORMTags, ", ") + "]"
		}
		fmt.Printf("   - %s (%s)%s\n", f.Name, f.Type, tags)
	}

	if relations.BelongsTo != "" {
		fmt.Printf("\n   Relations:\n")
		fmt.Printf("   - BelongsTo: %s (foreign key: %sID)\n", relations.BelongsTo, relations.BelongsTo)
	}
	if relations.HasMany != "" {
		fmt.Printf("\n   Relations:\n")
		fmt.Printf("   - HasMany: %s (adds %s slice to %s)\n", relations.HasMany, pluralize(relations.HasMany), modelName)
	}

	lower := strings.ToLower(modelName)
	fmt.Println("\n   Files to create:")
	fmt.Printf("   + internal/models/%s.go\n", lower)
	fmt.Printf("   + internal/interfaces/%s_repository.go\n", lower)
	fmt.Printf("   + internal/adapters/repository/%s_repository.go\n", lower)
	fmt.Printf("   + internal/domain/%s/service.go\n", lower)
	fmt.Printf("   + internal/domain/%s/module.go\n", lower)
	fmt.Printf("   + internal/adapters/handlers/%s_handler.go\n", lower)
	fmt.Printf("   + internal/domain/%s/service_test.go\n", lower)
	fmt.Printf("   + internal/adapters/handlers/%s_handler_test.go\n", lower)
	fmt.Println("\n   Files to modify:")
	fmt.Printf("   ~ internal/infrastructure/database/database.go (AutoMigrate)\n")
	fmt.Printf("   ~ internal/adapters/repository/module.go (fx.Provide)\n")
	fmt.Printf("   ~ internal/adapters/handlers/module.go (fx.Provide)\n")
	fmt.Printf("   ~ internal/adapters/http/routes.go (CRUD routes)\n")
	fmt.Printf("   ~ cmd/main.go (fx.Module)\n")
	if relations.HasMany != "" {
		fmt.Printf("   ~ internal/models/%s.go (add has-many relation)\n", strings.ToLower(relations.HasMany))
	}
}

// printAddModelHelp displays the help message for the add-model subcommand.
func printAddModelHelp() {
	fmt.Fprintf(os.Stderr, "Usage: create-go-starter add-model <ModelName> [options]\n\n")
	fmt.Fprintf(os.Stderr, "Add a new model to an existing go-starter-kit project.\n\n")
	fmt.Fprintf(os.Stderr, "Arguments:\n")
	fmt.Fprintf(os.Stderr, "  <ModelName>    Name of the model in PascalCase (e.g., Todo, BlogPost)\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  --fields string         Field definitions (required)\n")
	fmt.Fprintf(os.Stderr, "  --belongs-to <Model>    Add belongs-to relation (adds foreign key to this model)\n")
	fmt.Fprintf(os.Stderr, "  --has-many <Model>      Add has-many relation (adds slice to this model's parent)\n")
	fmt.Fprintf(os.Stderr, "  --public                Generate public routes (no auth middleware)\n")
	fmt.Fprintf(os.Stderr, "  --yes, -y               Skip confirmation prompt\n")
	fmt.Fprintf(os.Stderr, "  --help, -h              Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Field Syntax:\n")
	fmt.Fprintf(os.Stderr, "  name:type[:modifier1:modifier2...]\n\n")
	fmt.Fprintf(os.Stderr, "Supported Types:\n")
	fmt.Fprintf(os.Stderr, "  string     Go string type\n")
	fmt.Fprintf(os.Stderr, "  int        Go int type\n")
	fmt.Fprintf(os.Stderr, "  uint       Go uint type\n")
	fmt.Fprintf(os.Stderr, "  float64    Go float64 type\n")
	fmt.Fprintf(os.Stderr, "  bool       Go bool type\n")
	fmt.Fprintf(os.Stderr, "  time       Go time.Time type\n\n")
	fmt.Fprintf(os.Stderr, "GORM Modifiers:\n")
	fmt.Fprintf(os.Stderr, "  unique     Add unique constraint\n")
	fmt.Fprintf(os.Stderr, "  not_null   Add NOT NULL constraint\n")
	fmt.Fprintf(os.Stderr, "  index      Add database index\n\n")
	fmt.Fprintf(os.Stderr, "Relations:\n")
	fmt.Fprintf(os.Stderr, "  --belongs-to    Adds a foreign key field and relation to the parent model.\n")
	fmt.Fprintf(os.Stderr, "                  The parent model must already exist in internal/models/.\n")
	fmt.Fprintf(os.Stderr, "  --has-many      Modifies the specified parent model to add a slice of the\n")
	fmt.Fprintf(os.Stderr, "                  new model. The parent model must already exist.\n\n")
	fmt.Fprintf(os.Stderr, "Notes:\n")
	fmt.Fprintf(os.Stderr, "  - Pluralization uses simple rules (adds 's'). For irregular plurals\n")
	fmt.Fprintf(os.Stderr, "    (e.g., Person->People, Child->Children), manually edit the generated code.\n")
	fmt.Fprintf(os.Stderr, "  - Many-to-many relations are not yet supported (use --belongs-to for now).\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  create-go-starter add-model Todo --fields \"title:string,completed:bool\"\n")
	fmt.Fprintf(os.Stderr, "  create-go-starter add-model Product --fields \"name:string:unique:not_null,price:float64,stock:int\"\n")
	fmt.Fprintf(os.Stderr, "  create-go-starter add-model Comment --fields \"content:string\" --belongs-to Todo\n")
	fmt.Fprintf(os.Stderr, "  create-go-starter add-model Category --fields \"name:string:unique\" --has-many Product\n")
}
