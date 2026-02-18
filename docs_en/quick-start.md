# Quick Start (5 minutes)

<i class="material-icons success">bolt</i> Generate and launch a production-ready Go project in 5 minutes.

---

## 1. Install the CLI (30 seconds)

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

**Verification**:
```bash
create-go-starter doctor
```

---

## 2. Create the project (10 seconds)

```bash
create-go-starter my-api
cd my-api
```

<i class="material-icons success small">check</i> **Generated**: ~45 files with complete hexagonal architecture

> **Tip**: Use `create-go-starter -n my-api` to preview the files before creating them, or `create-go-starter -i` for the guided interactive mode.

---

## 3. Automatic setup (2 minutes)

```bash
./setup.sh
```

**This script does**:

- <i class="material-icons success small">check</i> Installs Go dependencies
- <i class="material-icons success small">check</i> Generates a JWT secret
- <i class="material-icons success small">check</i> Configures PostgreSQL (Docker)
- <i class="material-icons success small">check</i> Runs migrations

---

## 4. Launch the application (5 seconds)

```bash
make run
```

**Console output**:
```
INFO  Server starting on :8080
INFO  Database connected successfully
INFO  Migrations applied: 2
```

---

## 5. Test the API (1 minute)

### Health check

```bash
curl http://localhost:8080/health
```

**Response**:
```json
{"status":"healthy","timestamp":"2026-02-12T10:30:00Z"}
```

### Create a user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "SecurePass123!"
  }'
```

**Response**:
```json
{
  "user": {
    "id": 1,
    "email": "demo@example.com"
  },
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc..."
}
```

### Access your profile (protected)

```bash
TOKEN="<your_access_token>"

curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

**Response**:
```json
{
  "id": 1,
  "email": "demo@example.com",
  "created_at": "2026-02-12T10:32:00Z"
}
```

---

## <i class="material-icons success">celebration</i> Congratulations!

You now have a complete REST API with:

- <i class="material-icons success small">check</i> **Hexagonal Architecture** (Ports & Adapters)
- <i class="material-icons success small">check</i> **JWT Authentication** (access + refresh tokens)
- <i class="material-icons success small">check</i> **Database** (PostgreSQL + GORM)
- <i class="material-icons success small">check</i> **Validation** (go-playground/validator)
- <i class="material-icons success small">check</i> **Logging** (zerolog)
- <i class="material-icons success small">check</i> **Dependency Injection** (uber-go/fx)
- <i class="material-icons success small">check</i> **Tests** (structure ready)
- <i class="material-icons success small">check</i> **Docker** (Dockerfile + compose)

---

## Next steps

### <i class="material-icons">menu_book</i> Learn

- **[Full tutorial](tutorial/index.md)** - Build a Blog API (1h30)
- **[Complete guide](guide/index.md)** - Architecture and patterns
- **[Architecture](guide/architecture.md)** - Understanding hexagonal architecture

### <i class="material-icons">code</i> Develop

- **[Add a domain](guide/development.md#add-model)** - Use `add-model`
- **[API Reference](guide/api-reference.md)** - All endpoints
- **[Tests](guide/testing.md)** - Testing strategies

### <i class="material-icons">rocket_launch</i> Deploy

- **[Docker](guide/deployment.md#docker)** - Containerization
- **[Kubernetes](guide/deployment.md#kubernetes)** - Orchestration
- **[CI/CD](guide/deployment.md#cicd)** - Automation

---

## Need help?

- <i class="material-icons info">help</i> **[FAQ](reference/faq.md)** - Frequently asked questions
- <i class="material-icons">bug_report</i> **[Issues](https://github.com/tky0065/go-starter-kit/issues)** - Report a bug
- <i class="material-icons">forum</i> **[Discussions](https://github.com/tky0065/go-starter-kit/discussions)** - Ask a question
