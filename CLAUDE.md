# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `go-starter-kit`, a Go CLI tool generator that scaffolds new Go projects. The tool is called `create-go-starter` and helps bootstrap Go project structures.

## Commands

### Build and Install
```bash
# Build the binary
go build -o create-go-starter ./cmd/create-go-starter

# Install the binary to GOBIN (typically ~/go/bin/)
go install ./cmd/create-go-starter

# Run directly without installing
go run ./cmd/create-go-starter <project-name>
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./cmd/create-go-starter

# Run tests with verbose output
go test -v ./...

# Run a specific test function
go test -run TestColors ./cmd/create-go-starter
```

### Development
```bash
# Format code
go fmt ./...

# Vet code for common mistakes
go vet ./...

# Run the tool
./create-go-starter <project-name>

# Or with go run
go run ./cmd/create-go-starter <project-name>
```

## Project Structure

### CLI Tool Structure (this repository)
- `cmd/create-go-starter/` - Main CLI application entry point
  - `main.go` - CLI implementation with flag parsing and color utilities
  - `generator.go` - File generation orchestrator
  - `templates.go` - Core templates (config, server, domain)
  - `templates_user.go` - User domain specific templates
  - `colors_test.go` - Tests for ANSI color formatting functions
- `go.mod` - Go module definition (requires Go 1.25.0)
- `_bmad/` - BMAD workflow automation system (not part of core application)

### Generated Project Structure (projects created by the CLI)
When you run `create-go-starter my-project`, the following structure is generated:
- `internal/models/` - Shared domain entities (User, RefreshToken, AuthResponse). Prevents circular dependencies.
- `internal/domain/` - Business logic (services, not entities)
- `internal/interfaces/` - Port definitions (Hexagonal Architecture)
- `internal/adapters/` - HTTP handlers, middleware, repository implementations
- `internal/infrastructure/` - Database, server configuration
- `pkg/` - Shared libraries (config, logger, auth)

## Architecture Notes

### CLI Tool Design
The tool uses a simple command-line interface with:
- Flag parsing using the standard `flag` package
- ANSI color support via helper functions `Green()` and `Red()`
- Basic scaffolding placeholder ready for expansion

### Code Organization
- All CLI logic currently lives in a single `main.go` file
- Color utilities are embedded in the main package and tested separately
- The tool is designed to be distributed as a single binary via `go install`

### Testing Approach
- Unit tests are co-located with source code (e.g., `colors_test.go` alongside `main.go`)
- Tests use standard Go testing patterns with table-driven tests where appropriate

## Documentation Maintenance

**CRITICAL REQUIREMENT**: Documentation must ALWAYS be updated when making changes to the codebase.

### When to Update Documentation

Update documentation immediately after ANY of the following changes:

1. **Adding new features**:
   - New templates added → Update `docs/cli-architecture.md`
   - New generated files → Update `docs/usage.md` and `docs/generated-project-guide.md`
   - New domain entities → Update all architecture diagrams and examples

2. **Modifying architecture**:
   - Package structure changes → Update ALL documentation files
   - Dependency changes → Update architecture diagrams
   - Pattern changes → Update examples in `docs/generated-project-guide.md`

3. **Changing templates**:
   - Template modifications → Update corresponding documentation sections
   - New imports or dependencies → Update code examples throughout docs

4. **Bug fixes that affect structure**:
   - If the fix changes how generated projects work → Update docs immediately

### Documentation Files to Check

After ANY code change, review and update these files as needed:

- `README.md` - Main project overview and quick start
- `docs/usage.md` - Usage guide and generated structure
- `docs/generated-project-guide.md` - Complete guide for generated projects
- `docs/cli-architecture.md` - CLI tool architecture documentation
- `CLAUDE.md` - This file (AI context)
- `GEMINI.md` - Gemini AI context

### Documentation Update Process

1. **Make code changes** (templates, generator, etc.)
2. **Test the changes** (generate a test project, verify it works)
3. **Update ALL affected documentation** (use search to find all references)
4. **Verify documentation** (read through to ensure consistency)
5. **Commit code AND documentation together** in the same commit or sequential commits

### Example Workflow

```bash
# 1. Make code changes
vim cmd/create-go-starter/templates_user.go

# 2. Test changes
go build -o create-go-starter ./cmd/create-go-starter
./create-go-starter test-project
cd test-project && go build ./...

# 3. Update documentation
vim docs/cli-architecture.md
vim docs/generated-project-guide.md
vim docs/usage.md

# 4. Commit together
git add cmd/ docs/ README.md CLAUDE.md
git commit -m "feat: add new feature

- Implement feature X in templates
- Update all documentation to reflect changes
- Add examples in generated-project-guide.md"

git push origin main
```

### Why This Matters

- **Users rely on accurate docs**: Out-of-date documentation causes confusion and wastes time
- **Future contributors**: Need current docs to understand the system
- **AI assistants**: Use these docs as context for helping with the project
- **Consistency**: Prevents drift between code and documentation

**Remember**: Code without updated documentation is incomplete work. Always update docs as part of your changes.
