package main

import "fmt"

// GoModDependencies returns the appropriate database driver dependency for go.mod
func GoModDependencies(database string) string {
	switch database {
	case "mysql":
		return "gorm.io/driver/mysql v1.5.2"
	case "sqlite":
		return "gorm.io/driver/sqlite v1.5.4"
	case "postgres", "":
		return "gorm.io/driver/postgres v1.5.4"
	default:
		return "gorm.io/driver/postgres v1.5.4"
	}
}

// DatabaseImportPath returns the import path for the database driver
func DatabaseImportPath(database string) string {
	switch database {
	case "mysql":
		return "gorm.io/driver/mysql"
	case "sqlite":
		return "gorm.io/driver/sqlite"
	case "postgres", "":
		return "gorm.io/driver/postgres"
	default:
		return "gorm.io/driver/postgres"
	}
}

// DatabaseDSN generates the DSN (Data Source Name) connection string for the specified database
func DatabaseDSN(database, projectName string) string {
	switch database {
	case "mysql":
		return "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"
	case "sqlite":
		return "./${DB_NAME}.db"
	case "postgres", "":
		return "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=${DB_PORT} sslmode=disable"
	default:
		return "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=${DB_PORT} sslmode=disable"
	}
}

// DatabaseConnectionCode returns the GORM connection code for the specified database
func DatabaseConnectionCode(database string) string {
	switch database {
	case "mysql":
		return "mysql.Open(dsn)"
	case "sqlite":
		return "sqlite.Open(dsn)"
	case "postgres", "":
		return "postgres.Open(dsn)"
	default:
		return "postgres.Open(dsn)"
	}
}

// DatabaseDockerService generates the docker-compose service configuration for the database
func DatabaseDockerService(database, projectName string) string {
	switch database {
	case "mysql":
		return fmt.Sprintf(`  db:
    image: mysql:8.0
    container_name: %s_db
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p${DB_PASSWORD}"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  mysql_data:`, projectName)
	case "postgres", "":
		return fmt.Sprintf(`  db:
    image: postgres:16-alpine
    container_name: %s_db
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:`, projectName)
	case "sqlite":
		// SQLite doesn't need a Docker service
		return "# SQLite uses a local file, no Docker service needed"
	default:
		return fmt.Sprintf(`  db:
    image: postgres:16-alpine
    container_name: %s_db
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:`, projectName)
	}
}

// DatabaseEnvVars generates the environment variable configuration for .env file
func DatabaseEnvVars(database, projectName string) string {
	switch database {
	case "mysql":
		return fmt.Sprintf(`# Database Configuration (MySQL)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=%s`, projectName)
	case "sqlite":
		return fmt.Sprintf(`# Database Configuration (SQLite - local file .db)
DB_NAME=%s.db`, projectName)
	case "postgres", "":
		return fmt.Sprintf(`# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=%s`, projectName)
	default:
		return fmt.Sprintf(`# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=%s`, projectName)
	}
}
