# CLI Architecture

Technical documentation for contributors explaining how `create-go-starter` works internally.

!!! note "Translation in progress"
    This page is being translated from French. For the complete documentation, please refer to the [French version](../cli-architecture/).

## Overview

The CLI is built with **zero external dependencies** - it uses only the Go standard library.

## File Structure

```
cmd/create-go-starter/
├── main.go              # CLI entry point, flag parsing, color utilities
├── generator.go         # File generation orchestrator, validation
├── templates.go         # Core infrastructure templates (server, db, config)
├── templates_user.go    # User domain templates (handlers, services, auth)
└── *_test.go            # Tests co-located with source files
```

## Components

### main.go

- CLI entry point
- Flag parsing (`-h`, `--help`)
- Color utilities for terminal output
- Main function orchestration

### generator.go

- Project directory creation
- File generation orchestration
- Project name validation
- Error handling

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
2. Add file generation call in `generator.go`
3. Add tests for the new template

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestValidProjectName ./cmd/create-go-starter

# Verbose output
go test -v -run TestGoModTemplate ./cmd/create-go-starter
```

## Next Steps

- [Contributing Guide](contributing.md) - How to contribute
- [Usage Guide](usage.md) - Using the CLI
