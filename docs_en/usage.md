# Usage Guide

This guide details the usage of `create-go-starter` and explains the structure of generated projects.

## Basic Command

The basic syntax is very simple:

```bash
create-go-starter <project-name>
```

### Example

```bash
create-go-starter my-api-backend
```

This command will create a new `my-api-backend/` directory with the entire project structure using the **full** template by default.

## Available Options

### Flags and Aliases

All flags have short aliases for faster typing.

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--help` | `-h` | | Show help |
| `--interactive` | `-i` | `false` | Launch guided interactive mode |
| `--dry-run` | `-n` | `false` | Preview files without creating them |
| `--template` | `-t` | `full` | Template: `minimal`, `full`, `graphql` |
| `--database` | `-d` | `postgres` | Database: `postgres`, `mysql`, `sqlite` |
| `--observability` | `-o` | `none` | Observability: `none`, `basic`, `advanced` |

**Sub-commands:**

| Command | Description |
|---------|-------------|
| `doctor` | Environment diagnostics (Go, Git, Docker) |
| `add-model` | Generate a CRUD model in an existing project |

**Flexible syntax:** Flags accept `=` or space as separator:
```bash
# These syntaxes are equivalent:
create-go-starter -t=minimal my-app
create-go-starter -t minimal my-app
create-go-starter --template=minimal my-app
create-go-starter --template minimal my-app
```

### Examples

```bash
# Minimal template
create-go-starter my-project --template minimal
create-go-starter my-project -t minimal

# Full template (default)
create-go-starter my-project --template full
create-go-starter my-project  # Same result

# GraphQL template
create-go-starter my-project --template graphql
create-go-starter my-project -t graphql

# Choose MySQL as database
create-go-starter my-project --database=mysql
create-go-starter my-project -d mysql

# Choose SQLite (ideal for prototyping)
create-go-starter my-project --database=sqlite
create-go-starter my-project -d sqlite

# Combine template and database
create-go-starter my-project --template=minimal --database=sqlite
create-go-starter my-project -t minimal -d sqlite

# Dry-run: preview without creating files
create-go-starter my-project --dry-run
create-go-starter my-project -n

# Dry-run with specific options
create-go-starter my-project -t minimal -d sqlite -n

# Guided interactive mode
create-go-starter --interactive
create-go-starter -i

# Full combination with observability
create-go-starter my-project -t full -d postgres -o advanced
```

> **Notes**:
> - The `--template` flag is optional. If not specified, the **full** template is used by default.
> - The `--database` flag is optional. If not specified, **PostgreSQL** is used by default.
> - The `--dry-run` flag is optional. It displays the list of files that would be generated without creating them.
> - The `--interactive` flag launches a guided assistant and ignores other configuration flags.

## Interactive Mode (--interactive) <i class="material-icons success">new_releases</i>

**New in v1.4.0!** Interactive mode guides the user step by step to configure a new project.

### Usage

```bash
create-go-starter --interactive
create-go-starter -i
```

### How It Works

Interactive mode asks for:

1. **Project name** -- With real-time validation
2. **Template** -- Choice between minimal, full (default), graphql
3. **Database** -- Choice between postgres (default), mysql, sqlite
4. **Observability** -- Choice between none (default), basic, advanced (if full template)
5. **Summary** -- Display of chosen configuration with confirmation

```
create-go-starter - Interactive Mode

Enter project name: my-app
Select template (minimal/full/graphql) [full]: full
Select database (postgres/mysql/sqlite) [postgres]: postgres
Select observability (none/basic/advanced) [none]: advanced

Configuration Summary:
  Project:       my-app
  Template:      full
  Database:      postgres
  Observability: advanced

Proceed with generation? (y/n) [y]: y
```

### Notes

- <i class="material-icons info">info</i> Interactive mode requires an interactive terminal (no stdin pipe)
- <i class="material-icons warning">warning</i> `--interactive` and `--dry-run` cannot be used together
- <i class="material-icons info">info</i> Zero external dependencies -- uses only `bufio.NewReader` from stdlib

## Dry-Run Preview (--dry-run) <i class="material-icons success">new_releases</i>

**New in v1.4.0!** Dry-run mode displays the files that would be generated without creating them.

### Usage

```bash
create-go-starter my-app --dry-run
create-go-starter -n -t minimal -d sqlite my-app
```

### How It Works

Dry-run displays:
- The configuration used (template, database, observability)
- The complete list of files that would be created
- The number of files and directories
- A warning if the target directory already exists

```
Dry-run mode: no files will be created

Configuration:
  Project:       my-app
  Template:      minimal
  Database:      sqlite
  Observability: none

Files that would be generated (23 files):
  my-app/cmd/main.go
  my-app/internal/models/user.go
  my-app/internal/domain/user/service.go
  ...

Summary: 23 files in 12 directories
```

### Compatible with All Flags

```bash
create-go-starter -n -t full -d postgres -o advanced my-app
```

## Doctor Command <i class="material-icons success">new_releases</i>

**New in v1.4.0!** The `doctor` command checks that your environment is properly configured.

### Usage

```bash
create-go-starter doctor
```

### How It Works

The command checks:

| Tool | Check | Required |
|------|-------|----------|
| **Go** | Version >= 1.21 installed | <i class="material-icons success small">check</i> Yes |
| **Git** | Binary available | <i class="material-icons success small">check</i> Recommended |
| **Docker** | Binary + running daemon | <i class="material-icons info small">info</i> Optional |

### Example Output

```
create-go-starter doctor v1.4.0

Checking environment...

  [OK] Go 1.25.5 (minimum: 1.21)
  [OK] Git 2.43.0
  [OK] Docker 24.0.7 (daemon running)

All checks passed!
```

### Exit Code

- `0` — All checks passed
- `1` — One or more checks failed

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
