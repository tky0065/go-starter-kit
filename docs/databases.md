# Database Selection Guide

**Navigation:**

- :material-book-open: [Database Selection Guide](./databases.md) ← You are here
- :material-sync: [Database Migration Guide](./database-migration.md) - How to switch databases
- :material-arrow-left: [Back to README](../README.md)

---

go-starter-kit supports **3 database options** to fit your project needs.

## Quick Comparison

| Database | Best For | Complexity | Production Ready | Setup Time |
|----------|----------|------------|------------------|------------|
| **PostgreSQL** | Production apps, complex queries | Medium | :material-check-circle:{ .success } Yes | 2 min (Docker) |
| **MySQL** | Wide compatibility, shared hosting | Medium | :material-check-circle:{ .success } Yes | 2 min (Docker) |
| **SQLite** | Prototyping, small apps, embedded | Low | :material-alert-circle:{ .warning } Limited | 0 min |

## Detailed Comparison

### PostgreSQL (Default)

**Command:**
```bash
create-go-starter my-app
# OR explicitly:
create-go-starter my-app --database=postgres
```

**Strengths:**

- :material-check: Advanced SQL features (JSON, arrays, full-text search)
- :material-check: Excellent performance and reliability
- :material-check: ACID compliant, strong data integrity
- :material-check: Great for complex queries and analytics
- :material-check: Active community and ecosystem

**Limitations:**

- :material-alert: Requires Docker for local development
- :material-alert: Slightly more resource-intensive than MySQL

**When to use:**
- Production applications with complex data
- Apps requiring advanced SQL features
- Projects needing strong data integrity

**Docker Setup:**
```yaml
# Automatically included in docker-compose.yml
docker-compose up -d
```

**DSN Format:**
```
user:password@tcp(host:5432)/dbname?sslmode=disable
```

---

### MySQL/MariaDB

**Command:**
```bash
create-go-starter my-app --database=mysql
```

**Strengths:**

- :material-check: Wide compatibility and hosting support
- :material-check: Excellent for read-heavy workloads
- :material-check: Mature ecosystem and tooling
- :material-check: Easy to find hosting providers

**Limitations:**

- :material-alert: Fewer advanced features than PostgreSQL
- :material-alert: Some variations between MySQL and MariaDB

**When to use:**
- Shared hosting environments
- Read-heavy applications
- Teams familiar with MySQL
- Need for wide hosting compatibility

**Docker Setup:**
```yaml
# Automatically included in docker-compose.yml
docker-compose up -d
```

**DSN Format:**
```
user:password@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

---

### SQLite

**Command:**
```bash
create-go-starter my-app --database=sqlite
```

**Strengths:**

- :material-check: Zero configuration (no server needed)
- :material-check: Perfect for rapid prototyping
- :material-check: Single file database (easy backup/share)
- :material-check: Great for testing and development
- :material-check: Very fast for small datasets

**Limitations:**

- :material-alert: Limited concurrent writes (locks entire DB)
- :material-alert: No user/permission management
- :material-alert: Not suitable for high-traffic production
- :material-alert: Limited scalability

**When to use:**
- Rapid prototyping and MVPs
- Desktop applications
- Embedded systems
- Development and testing
- Small-scale production (<100 concurrent users)

**No Docker Needed:**
```bash
# Just run your app, SQLite file auto-created
go run ./cmd/main.go
# Creates: ./my_database.db
```

**DSN Format:**
```
./database.db
```

---

## Decision Matrix

**Choose PostgreSQL if:**

- :material-target: You're unsure (it's the default for a reason)
- :material-target: You need production-grade reliability
- :material-target: You have complex relational data

**Choose MySQL if:**

- :material-target: You're using shared hosting
- :material-target: Your team knows MySQL well
- :material-target: You have read-heavy workloads

**Choose SQLite if:**

- :material-target: You're prototyping or building an MVP
- :material-target: You want zero infrastructure setup
- :material-target: You have a small user base (<100 concurrent)

---

## Configuration Examples

### PostgreSQL

**.env.example:**
```bash
# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=myapp
DB_SSLMODE=disable
```

### MySQL

**.env.example:**
```bash
# Database Configuration (MySQL)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=myapp
```

### SQLite

**.env.example:**
```bash
# Database Configuration (SQLite - embedded)
DB_NAME=myapp.db
```

---

## Migration Guide

See [database-migration.md](./database-migration.md) for detailed migration instructions between databases.

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
- Best for: Single-user, low-concurrency scenarios
- Write performance: Limited (database-level locking)
- Read performance: Excellent for small datasets
- Connection pooling: Not applicable (file-based)

---

## Frequently Asked Questions

### Can I switch databases later?
Yes, but it requires regenerating the project with the new database flag and migrating data. See [database-migration.md](./database-migration.md) for details.

### Which database should I use for my SaaS?
For a production SaaS application, we recommend **PostgreSQL** for its reliability, ACID compliance, and advanced features. MySQL is also a solid choice if you're more familiar with it.

### Can I use SQLite in production?
SQLite can be used for small-scale production (<100 concurrent users), but we recommend PostgreSQL or MySQL for applications expecting growth.

### Do I need Docker?
- **PostgreSQL**: Yes (for local development)
- **MySQL**: Yes (for local development)
- **SQLite**: No (embedded database)

### What about MongoDB or NoSQL?
NoSQL support (MongoDB) was considered but deferred to future releases. Current focus is on SQL databases with GORM support.

---

## Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [GORM Documentation](https://gorm.io/docs/)
- [Database Migration Guide](./database-migration.md)

---

**Last Updated:** 2026-02-09
**Related Stories:** Epic 7 - Multi-Database Support (Stories 7.1-7.5)
