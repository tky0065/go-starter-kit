# Usage Guide

This guide details the usage of `create-go-starter` and explains in depth the structure of generated projects.

!!! note "Translation in progress"
    This page is being translated from French. For the complete documentation, please refer to the [French version](../usage/).

## Basic Command

The basic syntax is very simple:

```bash
create-go-starter <project-name>
```

### Example

```bash
create-go-starter my-api-backend
```

This command will create a new `my-api-backend/` directory with the entire project structure.

## Available Options

```bash
create-go-starter --help    # Display help
create-go-starter -h        # Alias for --help
```

## Naming Conventions

The project name must follow certain rules:

### Allowed Characters

- **Letters**: a-z, A-Z
- **Numbers**: 0-9
- **Hyphens**: -
- **Underscores**: _

### Restrictions

- No spaces
- No special characters (/, \, @, #, etc.)
- No dots (.)
- Must start with a letter or number

### Valid Examples

```bash
create-go-starter my-project           # OK
create-go-starter my-awesome-api       # OK
create-go-starter user_service         # OK
```

### Invalid Examples

```bash
create-go-starter my project           # Space not allowed
create-go-starter my-project!          # Special character
create-go-starter -my-project          # Starts with hyphen
```

## Generated Structure

Here's the complete structure created by `create-go-starter`:

```
my-project/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── models/                    # Shared domain entities
│   ├── domain/                    # Business logic layer
│   │   └── user/
│   │       ├── service.go
│   │       └── module.go
│   ├── adapters/                  # HTTP handlers, repositories
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── repository/
│   ├── infrastructure/            # DB, server configuration
│   │   ├── database/
│   │   └── server/
│   └── interfaces/                # Ports (interfaces)
├── pkg/                           # Reusable packages
│   ├── auth/
│   ├── config/
│   └── logger/
├── .github/workflows/ci.yml
├── .env
├── Dockerfile
├── Makefile
└── go.mod
```

## Workflow After Generation

### Option A: Automatic Setup (Recommended)

```bash
cd my-project
./setup.sh
make run
```

### Option B: Manual Setup

```bash
cd my-project

# Generate JWT secret
openssl rand -base64 32
# Add to .env: JWT_SECRET=<generated_secret>

# Install dependencies
go mod tidy

# Start PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_DB=my-project \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Run the application
make run
```

## Available Make Commands

| Command | Description |
|---------|-------------|
| `make help` | Display help |
| `make run` | Run the app |
| `make build` | Compile binary |
| `make test` | Run tests |
| `make lint` | Run linter |
| `make docker-build` | Build Docker image |

## Next Steps

- [Generated Project Guide](generated-project-guide.md) - Complete guide for development
- [CLI Architecture](cli-architecture.md) - Understand internals
