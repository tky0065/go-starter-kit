# Database Migration Guide

**Navigation:**

- <i class="material-icons">menu_book</i> [Database Selection Guide](./databases.md) - Which database to choose
- <i class="material-icons">sync</i> [Database Migration Guide](./database-migration.md) ← You are here
- <i class="material-icons">arrow_back</i> [Back to README](../README.md)

---

This guide helps you migrate between different database systems in go-starter-kit projects.

## Prerequisites

!!! warning "Important"
    Database migration is not automated. You will need to:
1. Export data from the source database
2. Regenerate the project with the target database
3. Import data into the target database
4. Test thoroughly

---

## Migration Paths

### PostgreSQL ↔ MySQL

**Difficulty:** <i class="material-icons success small">circle</i> Easy (both SQL, GORM compatible)

**Steps:**

1. **Export data from the source database**
   ```bash
   # Export PostgreSQL
   pg_dump -U postgres -h localhost -d monapp -F c -b -v -f backup.dump

   # Export MySQL
   mysqldump -u root -p monapp > backup.sql
   ```

2. **Regenerate the project with the target database**
   ```bash
   # Postgres → MySQL
   create-go-starter monapp --database=mysql
   
   # MySQL → Postgres
   create-go-starter monapp --database=postgres
   ```

3. **Convert SQL syntax if necessary**
   - Most differences are handled by GORM
   - Common changes:
     - Serial → AUTO_INCREMENT (handled by GORM)
     - BOOLEAN → TINYINT(1) (handled by GORM)
     - TEXT vs LONGTEXT (generally compatible)
   
   **For MySQL → PostgreSQL:**
   ```bash
   # Convert MySQL dump to PostgreSQL format
   # Option 1: Use pgloader (recommended)
   pgloader mysql://user:pass@localhost/monapp postgresql://postgres:pass@localhost/monapp
   
   # Option 2: Manual conversion
   # Replace AUTO_INCREMENT with SERIAL
   # Replace `backticks` with "quotes"
   # Adjust data types
   ```
   
   **For PostgreSQL → MySQL:**
   ```bash
   # Export as SQL then convert
   pg_dump -U postgres monapp -f dump.sql
   # Edit dump.sql: replace SERIAL with AUTO_INCREMENT, etc.
   mysql -u root -p monapp < dump_converted.sql
   ```

4. **Import data**
   ```bash
   # Import PostgreSQL
   pg_restore -U postgres -h localhost -d monapp backup.dump

   # Import MySQL
   mysql -u root -p monapp < backup.sql
   ```

5. **Test migrations and queries**
   ```bash
   go run ./cmd/main.go
   # Verify that all endpoints work correctly
   ```

**Common issues:**
- Serial vs AUTO_INCREMENT (handled by GORM)
- Some data types differ (TEXT vs LONGTEXT)
- Function syntax may differ (DATE_ADD vs INTERVAL)
- Identifiers: MySQL uses `backticks`, PostgreSQL uses "double quotes"

**Recommended tools:**
- **pgloader**: Automatic MySQL→PostgreSQL migration
- **mysql2postgres**: Schema conversion script
- **GORM AutoMigrate**: Automatically recreate the schema (custom data loss)

---

### SQL → SQLite (Downgrade)

**Difficulty:** <i class="material-icons warning small">circle</i> Medium (feature reduction)

**When to do this:**
- Moving from production to local development
- Creating a portable demo version
- Simplifying infrastructure for small projects

**Steps:**

1. **Export data as SQL or CSV**
   ```bash
   # Export PostgreSQL as CSV
   psql -U postgres -d monapp -c "COPY users TO '/tmp/users.csv' WITH CSV HEADER;"

   # Export MySQL as CSV
   mysql -u root -p monapp -e "SELECT * FROM users INTO OUTFILE '/tmp/users.csv'
     FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
   ```

2. **Regenerate the project**
   ```bash
   create-go-starter monapp --database=sqlite
   ```

3. **Import a limited dataset**
   ```bash
   # SQLite doesn't handle large datasets well
   # Import only essential data
   sqlite3 monapp.db < import.sql
   ```

4. **Remove advanced SQL features**
   - No stored procedures
   - Limited transaction support
   - No concurrent writes

**Limitations:**

- <i class="material-icons warning">warning</i> No concurrent writes (database-level locking)
- <i class="material-icons warning">warning</i> Limited data types
- <i class="material-icons warning">warning</i> No stored procedures
- <i class="material-icons warning">warning</i> No user/permission management
- <i class="material-icons warning">warning</i> Not suitable for large-scale production

**Recommended use cases:**
- Local development environments
- Portable demos
- Testing and CI/CD pipelines
- Small-scale applications (<100 users)

---

### SQLite → SQL (Upgrade)

**Difficulty:** <i class="material-icons success small">circle</i> Easy (adding features)

**When to do this:**
- Moving from prototype to production
- Need for concurrent write access
- Need for advanced SQL features

**Steps:**

1. **Export SQLite data**
   ```bash
   # Dump as SQL
   sqlite3 monapp.db .dump > backup.sql

   # Or export tables individually
   sqlite3 monapp.db -header -csv "SELECT * FROM users;" > users.csv
   ```

2. **Regenerate the project**
   ```bash
   create-go-starter monapp --database=postgres  # or mysql
   ```

3. **Convert the schema if necessary**
   - GORM migrations should handle most differences
   - Review and adjust any custom SQL

4. **Import data**
   ```bash
   # PostgreSQL
   psql -U postgres -d monapp < converted_backup.sql

   # MySQL
   mysql -u root -p monapp < converted_backup.sql
   ```

5. **Test thoroughly**
   - Verify data integrity
   - Test all API endpoints
   - Run the full test suite

---

## Data Export/Import Examples

### Export from PostgreSQL

```bash
# Full database dump (compressed)
pg_dump -U postgres -h localhost -d monapp -F c -b -v -f monapp_backup.dump

# SQL format (human-readable)
pg_dump -U postgres -h localhost -d monapp -f monapp_backup.sql

# Specific tables only
pg_dump -U postgres -h localhost -d monapp -t users -t posts -f partial_backup.sql

# CSV export for a specific table
psql -U postgres -d monapp -c "COPY users TO '/tmp/users.csv' WITH CSV HEADER;"
```

### Export from MySQL

```bash
# Full database dump
mysqldump -u root -p monapp > monapp_backup.sql

# Specific tables only
mysqldump -u root -p monapp users posts > partial_backup.sql

# CSV export
mysql -u root -p monapp -e "SELECT * FROM users INTO OUTFILE '/tmp/users.csv'
  FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
```

### Export from SQLite

```bash
# Full database dump
sqlite3 monapp.db .dump > monapp_backup.sql

# Specific table
sqlite3 monapp.db "SELECT * FROM users;" > users.txt

# CSV export
sqlite3 monapp.db -header -csv "SELECT * FROM users;" > users.csv
```

### Import into PostgreSQL

```bash
# From a PostgreSQL dump
pg_restore -U postgres -h localhost -d monapp monapp_backup.dump

# From a SQL file
psql -U postgres -d monapp < monapp_backup.sql

# From CSV
psql -U postgres -d monapp -c "\COPY users FROM '/tmp/users.csv' WITH CSV HEADER;"
```

### Import into MySQL

```bash
# From a SQL file
mysql -u root -p monapp < monapp_backup.sql

# From CSV
mysql -u root -p monapp -e "LOAD DATA LOCAL INFILE '/tmp/users.csv'
  INTO TABLE users
  FIELDS TERMINATED BY ','
  ENCLOSED BY '\"'
  LINES TERMINATED BY '\n'
  IGNORE 1 ROWS;"
```

### Import into SQLite

```bash
# From a SQL file
sqlite3 monapp.db < monapp_backup.sql

# From CSV
sqlite3 monapp.db <<EOF
.mode csv
.import users.csv users
EOF
```

---

## Post-Migration Testing

**Checklist:**
- [ ] All migrations run successfully
- [ ] Data integrity verified (counts, relationships)
- [ ] All API endpoints work correctly
- [ ] Authentication/authorization works
- [ ] All tests pass
- [ ] Performance is acceptable
- [ ] No data loss detected

**Validation commands:**
```bash
# Run all tests
go test ./...

# Verify data integrity
go run ./cmd/main.go
# Test all API endpoints manually or with automated tests

# Verify record counts
# Compare counts between source and target databases
```

---

## Rollback Plan

Always have a rollback plan before migrating:

1. **Back up the original database**
   ```bash
   # Keep the original database backup safe
   cp monapp_backup.dump monapp_backup_SAFE.dump
   ```

2. **Test the migration in staging first**
   - Never test migrations directly in production
   - Use a staging environment identical to production

3. **Schedule a maintenance window**
   - Plan the migration during low-traffic periods
   - Communicate the maintenance to users

4. **Document rollback steps**
   ```bash
   # Example rollback script
   #!/bin/bash
   echo "Retour à la base de données d'origine..."
   pg_restore -U postgres -h localhost -d monapp monapp_backup_SAFE.dump
   echo "Rollback terminé"
   ```

---

## Migration Checklist

**Before migration:**
- [ ] Complete backup of the source database
- [ ] Backup verified (restore test)
- [ ] Staging environment ready
- [ ] Migration plan documented
- [ ] Rollback plan prepared
- [ ] Maintenance scheduled
- [ ] Team notified

**During migration:**
- [ ] Stop the application (prevent data changes)
- [ ] Export data from the source
- [ ] Regenerate the project with the new database
- [ ] Import data into the target
- [ ] Run GORM migrations
- [ ] Verify data integrity

**After migration:**
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

Use case: You built an MVP with SQLite, now moving to production.

**Steps:**
1. Export SQLite data
2. Regenerate with PostgreSQL
3. Set up PostgreSQL in Docker/cloud
4. Import data
5. Test thoroughly

---

### Scenario 2: Shared Hosting to Cloud
**Path:** MySQL → PostgreSQL

Use case: Moving from shared hosting to cloud infrastructure.

**Steps:**
1. Export MySQL data
2. Regenerate with PostgreSQL
3. Set up managed PostgreSQL (AWS RDS, GCP Cloud SQL)
4. Import data
5. Update connection strings

---

### Scenario 3: Simplify Development
**Path:** PostgreSQL → SQLite

Use case: Faster local development without Docker.

**Steps:**
1. Export a minimal test dataset
2. Regenerate with SQLite
3. Import test data
4. Keep production on PostgreSQL

---

## Need help?

- Check the [Database Guide](./databases.md) for detailed information
- Open an issue on [GitHub](https://github.com/tky0065/go-starter-kit/issues)
- See examples in the `/examples` directory (if available)

---

## Additional Resources

- [PostgreSQL Migration Tools](https://www.postgresql.org/docs/current/migration.html)
- [MySQL Migration Guide](https://dev.mysql.com/doc/refman/8.0/en/migration.html)
- [SQLite Import/Export](https://www.sqlite.org/cli.html#csv_import)
- [GORM Migration](https://gorm.io/docs/migration.html)

---

**Last updated:** 2026-02-09
