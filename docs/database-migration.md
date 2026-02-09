# Database Migration Guide

**Navigation:**
- 📖 [Database Selection Guide](./databases.md) - Which database to choose
- 🔄 [Database Migration Guide](./database-migration.md) ← You are here
- 📚 [Back to README](../README.md)

---

This guide helps you migrate between different database systems in go-starter-kit projects.

## Prerequisites

⚠️ **Important:** Database migration is not automated. You'll need to:
1. Export data from source database
2. Regenerate project with target database
3. Import data to target database
4. Test thoroughly

---

## Migration Paths

### PostgreSQL ↔ MySQL

**Difficulty:** 🟢 Easy (both SQL, GORM compatible)

**Steps:**

1. **Export data from source database**
   ```bash
   # PostgreSQL export
   pg_dump -U postgres -h localhost -d myapp -F c -b -v -f backup.dump

   # MySQL export
   mysqldump -u root -p myapp > backup.sql
   ```

2. **Regenerate project with target database**
   ```bash
   # In your project directory
   create-go-starter myapp --database=mysql  # or postgres
   ```

3. **Convert SQL syntax if needed**
   - Most differences are handled by GORM
   - Common changes:
     - Serial → AUTO_INCREMENT (handled by GORM)
     - BOOLEAN → TINYINT(1) (handled by GORM)
     - TEXT vs LONGTEXT (usually compatible)

4. **Import data**
   ```bash
   # PostgreSQL import
   pg_restore -U postgres -h localhost -d myapp backup.dump

   # MySQL import
   mysql -u root -p myapp < backup.sql
   ```

5. **Test migrations and queries**
   ```bash
   go run ./cmd/main.go
   # Verify all endpoints work correctly
   ```

**Common Issues:**
- Serial vs AUTO_INCREMENT (handled by GORM)
- Some data types differ (TEXT vs LONGTEXT)
- Function syntax may differ (DATE_ADD vs INTERVAL)

---

### SQL → SQLite (Downgrade)

**Difficulty:** 🟡 Medium (feature reduction)

**When to do this:**
- Moving from production to local development
- Creating a portable demo version
- Simplifying infrastructure for small projects

**Steps:**

1. **Export data as SQL or CSV**
   ```bash
   # PostgreSQL CSV export
   psql -U postgres -d myapp -c "COPY users TO '/tmp/users.csv' WITH CSV HEADER;"

   # MySQL CSV export
   mysql -u root -p myapp -e "SELECT * FROM users INTO OUTFILE '/tmp/users.csv'
     FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
   ```

2. **Regenerate project**
   ```bash
   create-go-starter myapp --database=sqlite
   ```

3. **Import limited dataset**
   ```bash
   # SQLite doesn't handle large datasets well
   # Import only essential data
   sqlite3 myapp.db < import.sql
   ```

4. **Remove advanced SQL features**
   - No stored procedures
   - Limited transaction support
   - No concurrent writes

**Limitations:**
- ⚠️ No concurrent writes (database-level locking)
- ⚠️ Limited data types
- ⚠️ No stored procedures
- ⚠️ No user/permission management
- ⚠️ Not suitable for production at scale

**Recommended Use Cases:**
- Local development environments
- Portable demos
- Testing and CI/CD pipelines
- Small-scale applications (<100 users)

---

### SQLite → SQL (Upgrade)

**Difficulty:** 🟢 Easy (adding features)

**When to do this:**
- Scaling from prototype to production
- Need for concurrent write access
- Require advanced SQL features

**Steps:**

1. **Export SQLite data**
   ```bash
   # Dump to SQL
   sqlite3 myapp.db .dump > backup.sql

   # Or export tables individually
   sqlite3 myapp.db -header -csv "SELECT * FROM users;" > users.csv
   ```

2. **Regenerate project**
   ```bash
   create-go-starter myapp --database=postgres  # or mysql
   ```

3. **Convert schema if needed**
   - GORM migrations should handle most differences
   - Review and adjust any custom SQL

4. **Import data**
   ```bash
   # PostgreSQL
   psql -U postgres -d myapp < converted_backup.sql

   # MySQL
   mysql -u root -p myapp < converted_backup.sql
   ```

5. **Test thoroughly**
   - Verify data integrity
   - Test all API endpoints
   - Run full test suite

---

## Data Export/Import Examples

### Export from PostgreSQL

```bash
# Full database dump (compressed)
pg_dump -U postgres -h localhost -d myapp -F c -b -v -f myapp_backup.dump

# SQL format (human-readable)
pg_dump -U postgres -h localhost -d myapp -f myapp_backup.sql

# Specific tables only
pg_dump -U postgres -h localhost -d myapp -t users -t posts -f partial_backup.sql

# CSV export for a specific table
psql -U postgres -d myapp -c "COPY users TO '/tmp/users.csv' WITH CSV HEADER;"
```

### Export from MySQL

```bash
# Full database dump
mysqldump -u root -p myapp > myapp_backup.sql

# Specific tables only
mysqldump -u root -p myapp users posts > partial_backup.sql

# CSV export
mysql -u root -p myapp -e "SELECT * FROM users INTO OUTFILE '/tmp/users.csv'
  FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
```

### Export from SQLite

```bash
# Full database dump
sqlite3 myapp.db .dump > myapp_backup.sql

# Specific table
sqlite3 myapp.db "SELECT * FROM users;" > users.txt

# CSV export
sqlite3 myapp.db -header -csv "SELECT * FROM users;" > users.csv
```

### Import to PostgreSQL

```bash
# From PostgreSQL dump
pg_restore -U postgres -h localhost -d myapp myapp_backup.dump

# From SQL file
psql -U postgres -d myapp < myapp_backup.sql

# From CSV
psql -U postgres -d myapp -c "\COPY users FROM '/tmp/users.csv' WITH CSV HEADER;"
```

### Import to MySQL

```bash
# From SQL file
mysql -u root -p myapp < myapp_backup.sql

# From CSV
mysql -u root -p myapp -e "LOAD DATA LOCAL INFILE '/tmp/users.csv'
  INTO TABLE users
  FIELDS TERMINATED BY ','
  ENCLOSED BY '\"'
  LINES TERMINATED BY '\n'
  IGNORE 1 ROWS;"
```

### Import to SQLite

```bash
# From SQL file
sqlite3 myapp.db < myapp_backup.sql

# From CSV
sqlite3 myapp.db <<EOF
.mode csv
.import users.csv users
EOF
```

---

## Testing After Migration

**Checklist:**
- [ ] All migrations run successfully
- [ ] Data integrity verified (counts, relationships)
- [ ] All API endpoints work correctly
- [ ] Authentication/authorization works
- [ ] All tests pass
- [ ] Performance is acceptable
- [ ] No data loss detected

**Validation Commands:**
```bash
# Run all tests
go test ./...

# Check data integrity
go run ./cmd/main.go
# Test all API endpoints manually or with automated tests

# Verify record counts
# Compare source and target database counts
```

---

## Rollback Plan

Always have a rollback plan before migrating:

1. **Backup original database**
   ```bash
   # Keep original database backup safe
   cp myapp_backup.dump myapp_backup_SAFE.dump
   ```

2. **Test migration in staging first**
   - Never test migrations directly in production
   - Use a staging environment identical to production

3. **Have downtime window planned**
   - Schedule migration during low-traffic periods
   - Communicate downtime to users

4. **Document rollback steps**
   ```bash
   # Rollback script example
   #!/bin/bash
   echo "Rolling back to original database..."
   pg_restore -U postgres -h localhost -d myapp myapp_backup_SAFE.dump
   echo "Rollback complete"
   ```

---

## Migration Checklist

**Before Migration:**
- [ ] Full backup of source database
- [ ] Backup verified (test restore)
- [ ] Staging environment ready
- [ ] Migration plan documented
- [ ] Rollback plan prepared
- [ ] Downtime scheduled
- [ ] Team notified

**During Migration:**
- [ ] Stop application (prevent data changes)
- [ ] Export data from source
- [ ] Regenerate project with new database
- [ ] Import data to target
- [ ] Run GORM migrations
- [ ] Verify data integrity

**After Migration:**
- [ ] All tests pass
- [ ] API endpoints verified
- [ ] Performance acceptable
- [ ] Monitoring in place
- [ ] Old database backed up
- [ ] Documentation updated

---

## Common Migration Scenarios

### Scenario 1: Prototype to Production
**Path:** SQLite → PostgreSQL

Use Case: You built an MVP with SQLite, now scaling to production.

**Steps:**
1. Export SQLite data
2. Regenerate with PostgreSQL
3. Set up PostgreSQL in Docker/cloud
4. Import data
5. Test thoroughly

---

### Scenario 2: Shared Hosting to Cloud
**Path:** MySQL → PostgreSQL

Use Case: Moving from shared hosting to cloud infrastructure.

**Steps:**
1. Export MySQL data
2. Regenerate with PostgreSQL
3. Set up managed PostgreSQL (AWS RDS, GCP Cloud SQL)
4. Import data
5. Update connection strings

---

### Scenario 3: Simplifying Development
**Path:** PostgreSQL → SQLite

Use Case: Want faster local development without Docker.

**Steps:**
1. Export minimal test data
2. Regenerate with SQLite
3. Import test data
4. Keep production on PostgreSQL

---

## Need Help?

- Check [Database Guide](./databases.md) for detailed DB information
- Open an issue on [GitHub](https://github.com/yourusername/go-starter-kit/issues)
- See examples in `/examples` directory (if available)

---

## Additional Resources

- [PostgreSQL Migration Tools](https://www.postgresql.org/docs/current/migration.html)
- [MySQL Migration Guide](https://dev.mysql.com/doc/refman/8.0/en/migration.html)
- [SQLite Import/Export](https://www.sqlite.org/cli.html#csv_import)
- [GORM Migration](https://gorm.io/docs/migration.html)

---

**Last Updated:** 2026-02-09
**Related Stories:** Epic 7 - Multi-Database Support (Story 7.5)
