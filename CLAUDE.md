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

- `cmd/create-go-starter/` - Main CLI application entry point
  - `main.go` - CLI implementation with flag parsing and color utilities
  - `colors_test.go` - Tests for ANSI color formatting functions
- `go.mod` - Go module definition (requires Go 1.25.0)
- `_bmad/` - BMAD workflow automation system (not part of core application)

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
