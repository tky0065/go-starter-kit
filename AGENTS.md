# AGENTS.md

Instructions for AI coding agents working in this repository.

## Project Overview

`go-starter-kit` is a Go CLI tool (`create-go-starter`) that scaffolds new Go projects with hexagonal architecture. The CLI has **zero external dependencies** - it uses only the Go standard library.

## Build/Lint/Test Commands

### Build
```bash
go build -o create-go-starter ./cmd/create-go-starter  # Build binary
make build                                              # Creates bin/create-go-starter
go install ./cmd/create-go-starter                      # Install globally
```

### Test
```bash
go test ./...                                           # Run all tests
go test -v ./...                                        # Verbose output
go test -short ./...                                    # Skip E2E tests (faster)
go test -race ./...                                     # Race condition detection

# Single test by exact name
go test -run TestHelpFlag ./cmd/create-go-starter
go test -run TestGoModTemplate ./cmd/create-go-starter

# Tests matching pattern (regex)
go test -run TestInvalid ./cmd/create-go-starter        # All tests with "Invalid"
go test -run "Test.*Template" ./cmd/create-go-starter   # All template tests

# Single test with verbose output
go test -v -run TestValidProjectName ./cmd/create-go-starter

# Run tests in specific package only
go test ./cmd/create-go-starter
```

### Lint and Format
```bash
go fmt ./...           # Format code (always run before commit)
go vet ./...           # Check for common mistakes
make lint              # Run golangci-lint (if installed)
```

### Run
```bash
go run ./cmd/create-go-starter <project-name>
make run PROJECT_NAME=my-app
./bin/create-go-starter my-app                          # After make build
```

## Code Style Guidelines

### File Structure
```
cmd/create-go-starter/
  main.go              # CLI entry point, flag parsing, color utilities
  generator.go         # File generation orchestrator, validation
  templates.go         # Core infrastructure templates (server, db, config)
  templates_user.go    # User domain templates (handlers, services, auth)
  *_test.go            # Tests co-located with source files
```

### Imports Organization
Standard Go convention - **stdlib only** (no external dependencies in CLI):
```go
import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
)
```

### Naming Conventions
| Type | Convention | Example |
|------|------------|---------|
| Exported functions | PascalCase | `Green()`, `NewProjectTemplates()` |
| Unexported functions | camelCase | `validateProjectName()`, `copyEnvFile()` |
| Constants | PascalCase | `ColorGreen`, `ColorReset` |
| Types/Structs | PascalCase | `ProjectTemplates`, `FileGenerator` |
| Variables | camelCase | `projectName`, `projectPath` |
| Files | snake_case | `templates_user.go`, `colors_test.go` |
| Test functions | `Test` prefix | `TestValidProjectName`, `TestE2EGeneratedProjectBuilds` |

### Error Handling
```go
// 1. Wrap errors with context using %w
return fmt.Errorf("failed to create project directory: %w", err)

// 2. Early return pattern - validate and exit early
if err != nil {
    return err
}

// 3. Validation functions return errors, not booleans
func validateProjectName(name string) error {
    if !validProjectNamePattern.MatchString(name) {
        return fmt.Errorf("invalid project name '%s': ...", name)
    }
    return nil
}
```

### Testing Patterns
```go
// 1. Table-driven tests (preferred)
tests := []struct {
    name      string
    input     string
    wantError bool
}{
    {"valid name", "my-project", false},
    {"invalid chars", "my@project", true},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        err := validateProjectName(tt.input)
        if (err != nil) != tt.wantError {
            t.Errorf("got error %v, wantError %v", err, tt.wantError)
        }
    })
}

// 2. Use t.TempDir() for automatic cleanup
tmpDir := t.TempDir()

// 3. E2E tests - skip in short mode
if testing.Short() {
    t.Skip("skipping E2E test in short mode")
}
```

### Template Pattern
Templates use string concatenation (not `text/template`):
```go
func (t *ProjectTemplates) GoModTemplate() string {
    return `module ` + t.projectName + `

go 1.25.5
`
}
```

### File Permissions
- Directories: `0755` (use `defaultDirPerm` constant)
- Regular files: `0644`
- Executable scripts: `0755`

## Architecture Notes

- **Hexagonal Architecture**: Generated projects use ports & adapters pattern
- **No external deps in CLI**: The CLI tool itself uses only Go stdlib
- **Generated projects use**: fx (DI), fiber (HTTP), gorm (ORM), zerolog (logging)
- **Swagger integration**: Auto-generated via swaggo/swag annotations

## Generated Project Stack

| Component | Library | Purpose |
|-----------|---------|---------|
| HTTP | gofiber/fiber v2 | High-performance web framework |
| DI | uber-go/fx | Dependency injection |
| ORM | gorm.io/gorm | Database ORM with PostgreSQL |
| Logging | rs/zerolog | Structured JSON logging |
| Auth | golang-jwt/jwt | JWT authentication |
| Swagger | swaggo/swag | API documentation |

## Documentation Requirement

**CRITICAL**: Always update documentation when changing code.

Files to update after code changes:
- `README.md` - Project overview
- `docs/usage.md` - Usage guide  
- `docs/generated-project-guide.md` - Generated project guide
- `docs/cli-architecture.md` - CLI architecture
- `CLAUDE.md` / `AGENTS.md` - AI context files

## Quick Reference

| Task | Command |
|------|---------|
| Build | `make build` |
| Test all | `go test ./...` |
| Test single | `go test -run TestName ./cmd/create-go-starter` |
| Test verbose | `go test -v -run TestName ./cmd/create-go-starter` |
| Skip E2E | `go test -short ./...` |
| Format | `go fmt ./...` |
| Vet | `go vet ./...` |
| Lint | `make lint` |
| Run | `go run ./cmd/create-go-starter <name>` |
| Install | `go install ./cmd/create-go-starter` |
