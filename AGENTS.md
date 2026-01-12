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
go test ./...                                    # Run all tests
go test -v ./...                                 # Verbose output
go test -run TestHelpFlag ./cmd/create-go-starter  # Single test by name
go test -run TestInvalid ./cmd/create-go-starter   # Tests matching pattern
go test -short ./...                             # Skip E2E tests
```

### Lint and Format
```bash
go fmt ./...           # Format code
go vet ./...           # Vet for common mistakes
make lint              # Run golangci-lint
```

### Run
```bash
go run ./cmd/create-go-starter <project-name>
make run PROJECT_NAME=my-app
```

## Code Style Guidelines

### File Structure
```
cmd/create-go-starter/
  main.go           # CLI entry point, flag parsing, color utilities
  generator.go      # File generation orchestrator
  templates.go      # Core infrastructure templates
  templates_user.go # User domain templates
  *_test.go         # Tests co-located with source
```

### Imports Organization
Standard Go convention - stdlib only (no external dependencies):
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
1. **Wrap errors with context** using `%w`:
```go
return fmt.Errorf("failed to create project directory: %w", err)
```

2. **Early return pattern** - validate and exit early on errors

3. **Validation functions return errors** - not booleans

### Testing Patterns

1. **Table-driven tests** (preferred):
```go
tests := []struct {
    name      string
    input     string
    wantError bool
}{...}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {...})
}
```

2. **Integration tests** - build binary and execute with `exec.Command`

3. **E2E tests** - skip in short mode with `t.Skip()`

4. **Use `t.TempDir()`** for automatic cleanup

### Template Pattern
Templates use string concatenation (not `text/template`):
```go
func (t *ProjectTemplates) GoModTemplate() string {
    return `module ` + t.ProjectName + `

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
- **No external deps in CLI**: The CLI tool itself uses only stdlib
- **Generated projects use**: fx (DI), fiber (HTTP), gorm (ORM), zap (logging)

## Documentation Requirement

**CRITICAL**: Always update documentation when changing code.

Files to update:
- `README.md` - Project overview
- `docs/usage.md` - Usage guide
- `docs/generated-project-guide.md` - Generated project guide
- `docs/cli-architecture.md` - CLI architecture
- `CLAUDE.md` / `GEMINI.md` / `AGENT.md` - AI context files


## Quick Reference

| Task | Command |
|------|---------|
| Build | `make build` |
| Test all | `go test ./...` |
| Test single | `go test -run TestName ./cmd/create-go-starter` |
| Test verbose | `go test -v ./...` |
| Format | `go fmt ./...` |
| Lint | `make lint` |
| Run | `go run ./cmd/create-go-starter <name>` |
