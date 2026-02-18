# Database Selection Guide

**Navigation:**

- <i class="material-icons">menu_book</i> [Database Selection Guide](./databases.md) ← You are here
- <i class="material-icons">sync</i> [Database Migration Guide](./database-migration.md) - How to switch databases
- <i class="material-icons">arrow_back</i> [Back to README](../README.md)

---

go-starter-kit supports **3 database options** to fit your project's needs.

## Quick Comparison

| Database | Best for | Complexity | Production ready | Setup time |
|----------|----------|------------|------------------|------------|
| **PostgreSQL** | Production apps, complex queries | Medium | <i class="material-icons success">check_circle</i> Yes | 2 min (Docker) |
| **MySQL** | Broad compatibility, shared hosting | Medium | <i class="material-icons success">check_circle</i> Yes | 2 min (Docker) |
| **SQLite** | Prototyping, small apps, embedded | Low | <i class="material-icons warning">error</i> Limited | 0 min |

## Detailed Comparison

### PostgreSQL (Default)

**Command:**
```bash
create-go-starter mon-app
# OR explicitly:
create-go-starter mon-app --database=postgres
```

**Strengths:**

- <i class="material-icons success">check</i> Advanced SQL features (JSON, arrays, full-text search)
- <i class="material-icons success">check</i> Excellent performance and reliability
- <i class="material-icons success">check</i> ACID compliant, strong data integrity
- <i class="material-icons success">check</i> Ideal for complex queries and analytics
- <i class="material-icons success">check</i> Active community and rich ecosystem

**Limitations:**

- <i class="material-icons warning">warning</i> Requires Docker for local development
- <i class="material-icons warning">warning</i> Slightly more resource-intensive than MySQL

**When to use:**
- Production applications with complex data
- Applications requiring advanced SQL features
- Projects requiring strong data integrity

**Docker configuration:**
```yaml
# Automatically included in docker-compose.yml
docker-compose up -d
```

**DSN format:**
```
user:password@tcp(host:5432)/dbname?sslmode=disable
```

---

### MySQL/MariaDB

**Command:**
```bash
create-go-starter mon-app --database=mysql
```

**Strengths:**

- <i class="material-icons success">check</i> Broad compatibility and hosting support
- <i class="material-icons success">check</i> Excellent for read-heavy workloads
- <i class="material-icons success">check</i> Mature ecosystem and tooling
- <i class="material-icons success">check</i> Easy to find hosting providers

**Limitations:**

- <i class="material-icons warning">warning</i> Fewer advanced features than PostgreSQL
- <i class="material-icons warning">warning</i> Some variations between MySQL and MariaDB

**When to use:**
- Shared hosting environments
- Read-heavy applications
- Teams familiar with MySQL
- Need for broad hosting compatibility

**Docker configuration:**
```yaml
# Automatically included in docker-compose.yml
docker-compose up -d
```

**DSN format:**
```
user:password@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

---

### SQLite

**Command:**
```bash
create-go-starter mon-app --database=sqlite
```

**Strengths:**

- <i class="material-icons success">check</i> Zero configuration (no server needed)
- <i class="material-icons success">check</i> Perfect for rapid prototyping
- <i class="material-icons success">check</i> Single-file database (easy backup/sharing)
- <i class="material-icons success">check</i> Ideal for testing and development
- <i class="material-icons success">check</i> Very fast for small datasets

**Limitations:**

- <i class="material-icons warning">warning</i> Limited concurrent writes (locks entire DB)
- <i class="material-icons warning">warning</i> No user/permission management
- <i class="material-icons warning">warning</i> Not suitable for high-traffic production
- <i class="material-icons warning">warning</i> Limited scalability

**When to use:**
- Rapid prototyping and MVPs
- Desktop applications
- Embedded systems
- Development and testing
- Small-scale production (<100 concurrent users)

**No Docker needed:**
```bash
# Simply run your app, the SQLite file is created automatically
go run ./cmd/main.go
# Creates: ./my_database.db
```

**DSN format:**
```
./database.db
```

---

## Decision Matrix

**Choose PostgreSQL if:**

- <i class="material-icons">center_focus_strong</i> You're unsure (it's the default for good reason)
- <i class="material-icons">center_focus_strong</i> You need production-grade reliability
- <i class="material-icons">center_focus_strong</i> You have complex relational data

**Choose MySQL if:**

- <i class="material-icons">center_focus_strong</i> You're using shared hosting
- <i class="material-icons">center_focus_strong</i> Your team is well-versed in MySQL
- <i class="material-icons">center_focus_strong</i> You have read-heavy workloads

**Choose SQLite if:**

- <i class="material-icons">center_focus_strong</i> You're prototyping or building an MVP
- <i class="material-icons">center_focus_strong</i> You want zero infrastructure
- <i class="material-icons">center_focus_strong</i> You have a small user base (<100 concurrent)

---

## Configuration Examples

### PostgreSQL

**.env.example:**
```bash
# Configuration base de données (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=monapp
DB_SSLMODE=disable
```

### MySQL

**.env.example:**
```bash
# Configuration base de données (MySQL)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=monapp
```

### SQLite

**.env.example:**
```bash
# Configuration base de données (SQLite - embarqué)
DB_NAME=monapp.db
```

---

## Migration Guide

See [database-migration.md](./database-migration.md) for detailed instructions on migrating between databases.

---

## Performance Considerations

### PostgreSQL
- Best for: Complex queries, ACID transactions, concurrent writes
- Write performance: Excellent (MVCC architecture)
- Read performance: Excellent with proper indexing
- Connection pooling: Recommended for production

### MySQL
- Best for: Read-heavy workloads, simple queries
- Write performance: Good (row-level locking)
- Read performance: Excellent (query cache)
- Connection pooling: Recommended for production

### SQLite
- Best for: Single-user scenarios, low concurrency
- Write performance: Limited (database-level locking)
- Read performance: Excellent for small datasets
- Connection pooling: Not applicable (file-based)

---

## Frequently Asked Questions

### Can I change the database later?
Yes, but it requires regenerating the project with the new database flag and migrating the data. See [database-migration.md](./database-migration.md) for details.

### Which database should I use for my SaaS?
For a production SaaS application, we recommend **PostgreSQL** for its reliability, ACID compliance, and advanced features. MySQL is also a good choice if you're more familiar with it.

### Can I use SQLite in production?
SQLite can be used for small-scale production (<100 concurrent users), but we recommend PostgreSQL or MySQL for applications that plan to grow.

### Do I need Docker?
- **PostgreSQL**: Yes (for local development)
- **MySQL**: Yes (for local development)
- **SQLite**: No (embedded database)

### What about MongoDB or NoSQL?
NoSQL support (MongoDB) has been considered but deferred to future versions. The current focus is on SQL databases with GORM support.

---

## Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [GORM Documentation](https://gorm.io/docs/)
- [Database Migration Guide](./database-migration.md)

---

**Last updated:** 2026-02-09
