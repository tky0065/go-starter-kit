# Part 1: Installation and Configuration

<i class="material-icons info small">circle</i> **Part 1/4** - Estimated time: 15 minutes

[<i class="material-icons">arrow_back</i> Back to index](index.md)

---

## Goal

Create a complete REST API for a blog with:

- **Articles (Posts)** with author, title, content, tags
- **Comments** on articles
- **JWT Authentication** (already included in create-go-starter)
- **Comprehensive tests**
- **Docker deployment**

At the end of this tutorial, you will have a production-ready Blog API with all best practices.

## Prerequisites

### Required Software

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **PostgreSQL** or **Docker** - For the database
- **curl** or **Postman** - To test the API
- Code editor (VS Code, GoLand, etc.)

### Recommended Knowledge

- Go basics (structs, interfaces, error handling)
- REST API concepts
- Familiarity with SQL/PostgreSQL (basic)

No need to be an expert! This tutorial explains each step in detail.

---

## Step 1: CLI Installation

### Global installation (recommended)

The simplest method to install `create-go-starter`:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

This command downloads, compiles and installs the CLI globally.

### Verification

```bash
create-go-starter --help
```

You should see the help output displayed.

**Note**: If the command is not found, add `$GOPATH/bin` to your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Step 2: Project Generation

### Create the project

```bash
create-go-starter blog-api
```

This command generates **~45 files** with all the necessary architecture.

### Generated structure

```bash
cd blog-api
tree -L 3
```

**Result**:
```
blog-api/
├── cmd/
│   └── main.go                       # Entry point with fx DI
├── internal/
│   ├── models/
│   │   └── user.go                   # Entities: User, RefreshToken, AuthResponse
│   ├── domain/
│   │   ├── user/                     # User domain (pre-generated)
│   │   │   ├── service.go
│   │   │   └── module.go
│   │   └── errors.go
│   ├── adapters/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go
│   │   │   └── error_handler.go
│   │   └── repository/
│   │       └── user_repository.go
│   ├── infrastructure/
│   │   ├── database/
│   │   └── server/
│   └── interfaces/                   # Ports (interfaces)
│       └── user_repository.go
├── pkg/
│   ├── auth/                         # JWT utilities
│   ├── config/                       # Configuration
│   └── logger/                       # Zerolog logger
├── docs/
│   ├── README.md
│   └── quick-start.md
├── .env                              # Configuration (auto-copied)
├── .env.example
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

<i class="material-icons success">check_circle</i> **Checkpoint 1**: The project is generated successfully.

---

## Step 3: Initial Configuration

### 3.1 Install dependencies

```bash
cd blog-api
go mod tidy
```

This command downloads all dependencies (Fiber, GORM, fx, etc.).

### 3.2 Configure PostgreSQL

You have 2 options:

#### Option A: Docker (recommended)

```bash
docker run -d \
  --name blog-postgres \
  -e POSTGRES_DB=blog_api \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

#### Option B: Local PostgreSQL

If PostgreSQL is installed locally:

```bash
createdb blog_api
```

### 3.3 Configure environment variables

Generate a secure JWT secret:

```bash
JWT_SECRET=$(openssl rand -base64 32)
echo "Generated JWT_SECRET: $JWT_SECRET"
```

Edit the `.env` file:

```bash
nano .env
```

Contents of `.env`:

```env
# Application
APP_NAME=blog-api
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=blog_api
DB_SSLMODE=disable

# JWT
JWT_SECRET=<paste_generated_secret_here>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

**Important**: Replace `<paste_generated_secret_here>` with the generated JWT_SECRET.

---

## Step 4: Test the Base Project

### 4.1 Launch the application

```bash
make run
```

You should see:

```
2024/01/10 10:00:00 INF Starting blog-api server on :8080
```

### 4.2 Test the health check

In another terminal:

```bash
curl http://localhost:8080/health
```

**Expected response**:
```json
{"status":"ok"}
```

### 4.3 Test the default authentication

#### Create a user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Response**:
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "user": {
    "id": 1,
    "email": "admin@blog.com",
    "created_at": "2024-01-10T10:05:00Z"
  }
}
```

#### Log in

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Same response** with access_token and refresh_token.

#### Test a protected route

```bash
# Replace <ACCESS_TOKEN> with the received token
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response**:
```json
[
  {
    "id": 1,
    "email": "admin@blog.com",
    "created_at": "2024-01-10T10:05:00Z"
  }
]
```

<i class="material-icons success">check_circle</i> **Checkpoint 2**: The base project works perfectly with User and Auth.

---

## Part 1 Summary

<i class="material-icons success">check</i> Installation of the `create-go-starter` CLI
<i class="material-icons success">check</i> Generation of a complete project
<i class="material-icons success">check</i> PostgreSQL and JWT configuration
<i class="material-icons success">check</i> Authentication testing

You now have a working project with JWT authentication and user management.

---

## Navigation

[<i class="material-icons">arrow_back</i> Back to index](index.md) | [Part 2: Create Your First Domain <i class="material-icons">arrow_forward</i>](02-first-domain.md)
