# Contributing Guide

Thank you for your interest in contributing to Go Starter Kit!

!!! note "Translation in progress"
    This page is being translated from French. For the complete documentation, please refer to the [French version](../contributing/).

## How to Contribute

### Reporting Bugs

1. Check if the issue already exists
2. Create a new issue with:
   - Clear title
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS

### Suggesting Features

1. Open a discussion or issue
2. Describe the feature and use case
3. Provide examples if possible

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Run linter: `make lint`
6. Commit with clear message: `git commit -m "feat: add new feature"`
7. Push and create PR

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/go-starter-kit.git
cd go-starter-kit

# Add upstream remote
git remote add upstream https://github.com/tky0065/go-starter-kit.git

# Install and test
go build -o create-go-starter ./cmd/create-go-starter
go test ./...
```

## Code Style

- Follow Go conventions
- Run `go fmt ./...` before committing
- Run `go vet ./...` to check for issues
- Use meaningful variable and function names
- Add comments for exported functions

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code refactoring
- `test:` - Tests
- `chore:` - Maintenance

## Testing

```bash
# Run all tests
go test ./...

# With verbose output
go test -v ./...

# Skip E2E tests
go test -short ./...

# With race detection
go test -race ./...
```

## Documentation

When changing code, update:

- `README.md` - If user-facing changes
- `docs/*.md` - Relevant documentation
- `AGENTS.md` - If changing CLI architecture

## Questions?

- Open a [GitHub Discussion](https://github.com/tky0065/go-starter-kit/discussions)
- Check existing issues and PRs

Thank you for contributing!
