# Database Selection Guide

**Navigation:**
- 📖 [Database Selection Guide](./databases.md) ← You are here
- 🔄 [Database Migration Guide](./database-migration.md) - How to switch databases
- 📚 [Back to README](../README.md)

---

go-starter-kit supports **3 database options** to fit your project needs.

## Quick Comparison

| Database | Best For | Complexity | Production Ready | Setup Time |
|----------|----------|------------|------------------|------------|
| **PostgreSQL** | Production apps, complex queries | Medium | ✅ Yes | 2 min (Docker) |
| **MySQL** | Wide compatibility, shared hosting | Medium | ✅ Yes | 2 min (Docker) |
| **SQLite** | Prototyping, small apps, embedded | Low | ⚠️ Limited | 0 min |

## Detailed Comparison

### PostgreSQL (Default)

**Command:**
```bash
create-go-starter my-app
# OR explicitly:
create-go-starter my-app --database=postgres
```

**Strengths:**
- ✅ Advanced SQL features (JSON, arrays, full-text search)
- ✅ Excellent performance and reliability
- ✅ ACID compliant, strong data integrity
- ✅ Great for complex queries and analytics
- ✅ Active community and ecosystem

**Limitations:**
- ⚠️ Requires Docker for local development
- ⚠️ Slightly more resource-intensive than MySQL

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
- ✅ Wide compatibility and hosting support
- ✅ Excellent for read-heavy workloads
- ✅ Mature ecosystem and tooling
- ✅ Easy to find hosting providers

**Limitations:**
- ⚠️ Fewer advanced features than PostgreSQL
- ⚠️ Some variations between MySQL and MariaDB

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
- ✅ Zero configuration (no server needed)
- ✅ Perfect for rapid prototyping
- ✅ Single file database (easy backup/share)
- ✅ Great for testing and development
- ✅ Very fast for small datasets

**Limitations:**
- ⚠️ Limited concurrent writes (locks entire DB)
- ⚠️ No user/permission management
- ⚠️ Not suitable for high-traffic production
- ⚠️ Limited scalability

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
- 🎯 You're unsure (it's the default for a reason)
- 🎯 You need production-grade reliability
- 🎯 You have complex relational data

**Choose MySQL if:**
- 🎯 You're using shared hosting
- 🎯 Your team knows MySQL well
- 🎯 You have read-heavy workloads

**Choose SQLite if:**
- 🎯 You're prototyping or building an MVP
- 🎯 You want zero infrastructure setup
- 🎯 You have a small user base (<100 concurrent)

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
