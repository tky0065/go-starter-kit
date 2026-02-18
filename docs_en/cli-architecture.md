# CLI Architecture

Technical documentation for contributors explaining how `create-go-starter` works internally.

## Overview

The CLI is built with **zero external dependencies** - it uses only the Go standard library.

## File Structure

```
cmd/create-go-starter/
├── main.go              # CLI entry point, custom flag parsing, color utilities
├── generator.go         # File generation orchestrator (build + write), validation
├── templates.go         # Core infrastructure templates (server, db, config)
├── templates_user.go    # User domain templates (handlers, services, auth)
├── interactive.go       # Interactive mode (--interactive / -i)
├── dryrun.go            # Dry-run preview (--dry-run / -n)
├── doctor.go            # Doctor command (environment diagnostics)
├── progress.go          # Progress bar for file generation
├── stats.go             # Generation statistics display
├── version.go           # Version constant
└── *_test.go            # Tests co-located with source files
```

## Components

### main.go

- CLI entry point
- Custom flag parsing (replaced stdlib `flag` package for alias support)
- Short aliases: `-t`, `-d`, `-o`, `-i`, `-n`, `-h`
- Color utilities for terminal output
- Sub-command dispatch (`doctor`, `add-model`)
- Main function orchestration

### generator.go

- Project directory creation
- Two-phase generation: `build*FileList()` then `writeFiles()`
- Project name validation
- Error handling
- Dry-run support (build without write)

### interactive.go

- Step-by-step guided wizard using `bufio.NewReader`
- Project name input with validation
- Template, database, observability selection
- Configuration summary with confirmation
- Zero external dependencies

### dryrun.go

- Displays files that would be generated without creating them
- Shows configuration summary
- File and directory count
- Existing directory warnings

### doctor.go

- Environment diagnostics command
- Go version check (>= 1.21)
- Git availability check
- Docker binary + daemon status check
- Colored pass/fail output

### progress.go / stats.go

- Visual progress bar during file generation
- Auto-disables on non-TTY or NO_COLOR environments
- Generation statistics (files created, directories, duration)

### templates.go / templates_user.go

- Template definitions using string concatenation
- Project name interpolation
- All generated file contents

## Template Pattern

Templates use string concatenation (not `text/template`):

```go
func (t *ProjectTemplates) GoModTemplate() string {
    return `module ` + t.projectName + `

go 1.25.5
`
}
```

## Adding New Templates

1. Add template method to `ProjectTemplates`
2. Add file entry in the `build*FileList()` function in `generator.go`
3. Add tests for the new template

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestValidProjectName ./cmd/create-go-starter

# Verbose output
go test -v -run TestGoModTemplate ./cmd/create-go-starter

# Skip E2E tests
go test -short ./...
```

## Next Steps

- [Contributing Guide](contributing.md) - How to contribute
- [Usage Guide](usage.md) - Using the CLI
