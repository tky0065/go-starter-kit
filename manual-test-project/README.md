# manual-test-project

A Go application scaffolded with create-go-starter.

## Architecture

This project follows Hexagonal Architecture (Ports and Adapters) pattern:

- **cmd/**: Application entry points
- **internal/adapters/**: External adapters (HTTP handlers, etc.)
- **internal/domain/**: Core business logic
- **internal/interfaces/**: Port definitions
- **internal/infrastructure/**: Infrastructure concerns (DB, config, etc.)
- **pkg/**: Public libraries

## Prerequisites

- Go 1.25 or later
- Docker (optional, for containerized deployment)
- PostgreSQL (or use Docker Compose)

## Getting Started

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Copy environment file:
   ```bash
   cp .env.example .env
   ```

3. Run the application:
   ```bash
   make run
   ```

## Development

### Available Commands

```bash
make help          # Show all available commands
make build         # Build the binary
make run           # Run the application
make test          # Run tests
make clean         # Clean build artifacts
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

### Running Tests

```bash
go test ./...
```

## Project Structure

```
manual-test-project/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── adapters/            # HTTP, gRPC adapters
│   ├── domain/              # Business logic
│   ├── interfaces/          # Port definitions
│   └── infrastructure/      # DB, config, logging
├── pkg/                     # Public libraries
├── deployments/             # Docker, K8s configs
├── .env.example             # Environment template
├── Dockerfile               # Docker build config
├── Makefile                 # Build automation
├── go.mod                   # Go modules
└── README.md                # This file
```

## License

MIT
