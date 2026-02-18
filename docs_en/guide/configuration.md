# Configuration

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Environment Variables

The `.env` file contains all the configuration:

```bash
# Application
APP_NAME=mon-projet
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mon-projet
DB_SSLMODE=disable

# JWT
JWT_SECRET=          # MUST BE FILLED IN!
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

### Generating a Secure JWT_SECRET

**CRITICAL**: Always generate a strong secret:

```bash
openssl rand -base64 32
```

Example output:
```
XqR7nP3vY2kL9wH4sT6mU8jC1bN5aD0f
```

Add it to `.env`:
```
JWT_SECRET=XqR7nP3vY2kL9wH4sT6mU8jC1bN5aD0f
```

### Per-Environment Configuration

#### Development

```bash
APP_ENV=development
DB_HOST=localhost
DB_SSLMODE=disable
```

#### Staging

```bash
APP_ENV=staging
DB_HOST=staging-db.example.com
DB_SSLMODE=require
JWT_SECRET=<secret-depuis-vault>
```

#### Production

```bash
APP_ENV=production
DB_HOST=prod-db.example.com
DB_SSLMODE=require
DB_PASSWORD=<secret-depuis-secrets-manager>
JWT_SECRET=<secret-depuis-secrets-manager>
```

**Best practice**: Use secrets managers:

- **AWS**: Secrets Manager, Parameter Store
- **GCP**: Secret Manager
- **Kubernetes**: Secrets
- **HashiCorp**: Vault

### PostgreSQL Configuration

#### Option 1: Local PostgreSQL

**macOS (Homebrew)**:

```bash
brew install postgresql@16
brew services start postgresql@16
createdb mon-projet
```

**Linux (apt)**:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo -u postgres createdb mon-projet
```

#### Option 2: Docker

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_DB=mon-projet \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

#### Option 3: Docker Compose

If a `docker-compose.yml` is generated:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: mon-projet
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Start it:

```bash
docker-compose up -d postgres
```

#### Verifying the Connection

```bash
# With psql
psql -h localhost -U postgres -d mon-projet

# Or test from the app
make run
# Check the logs: "Database connected successfully"
```

---



---

## Navigation

**Previous**: [Project Structure](structure.md)  
**Next**: [Development](development.md)  
**Index**: [Guide Index](index.md)
