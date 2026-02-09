package main

import (
	"strings"
	"testing"
)

// TestDockerComposeTemplateMySQL tests docker-compose generation for MySQL
func TestDockerComposeTemplateMySQL(t *testing.T) {
	projectName := "test-mysql-project"
	templates := NewProjectTemplatesWithDatabase(projectName, "mysql")
	content := templates.DockerComposeTemplate()

	// Should contain MySQL image
	if !strings.Contains(content, "mysql:8.0") {
		t.Error("DockerComposeTemplate() for MySQL should contain mysql:8.0 image")
	}

	// Should contain MySQL environment variables
	mysqlEnvVars := []string{
		"MYSQL_ROOT_PASSWORD",
		"MYSQL_DATABASE",
		"MYSQL_USER",
		"MYSQL_PASSWORD",
	}
	for _, envVar := range mysqlEnvVars {
		if !strings.Contains(content, envVar) {
			t.Errorf("DockerComposeTemplate() for MySQL should contain %s", envVar)
		}
	}

	// Should have MySQL port 3306
	if !strings.Contains(content, "3306") {
		t.Error("DockerComposeTemplate() for MySQL should use port 3306")
	}

	// Should have MySQL healthcheck
	if !strings.Contains(content, "mysqladmin") {
		t.Error("DockerComposeTemplate() for MySQL should have mysqladmin healthcheck")
	}

	// Should have mysql_data volume
	if !strings.Contains(content, "mysql_data") {
		t.Error("DockerComposeTemplate() for MySQL should have mysql_data volume")
	}

	// API service should have correct DB_PORT
	if !strings.Contains(content, "DB_PORT: 3306") {
		t.Error("DockerComposeTemplate() for MySQL should set DB_PORT to 3306 in api service")
	}

	// API service should have correct DB_USER
	if !strings.Contains(content, "DB_USER: root") {
		t.Error("DockerComposeTemplate() for MySQL should set DB_USER to root in api service")
	}
}

// TestDockerComposeTemplatePostgres tests docker-compose generation for PostgreSQL (default)
func TestDockerComposeTemplatePostgres(t *testing.T) {
	projectName := "test-postgres-project"
	templates := NewProjectTemplatesWithDatabase(projectName, "postgres")
	content := templates.DockerComposeTemplate()

	// Should contain PostgreSQL image
	if !strings.Contains(content, "postgres:16-alpine") {
		t.Error("DockerComposeTemplate() for PostgreSQL should contain postgres:16-alpine image")
	}

	// Should contain PostgreSQL environment variables
	pgEnvVars := []string{
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
	}
	for _, envVar := range pgEnvVars {
		if !strings.Contains(content, envVar) {
			t.Errorf("DockerComposeTemplate() for PostgreSQL should contain %s", envVar)
		}
	}

	// Should have PostgreSQL port 5432
	if !strings.Contains(content, "5432") {
		t.Error("DockerComposeTemplate() for PostgreSQL should use port 5432")
	}

	// Should have PostgreSQL healthcheck
	if !strings.Contains(content, "pg_isready") {
		t.Error("DockerComposeTemplate() for PostgreSQL should have pg_isready healthcheck")
	}

	// Should have postgres_data volume
	if !strings.Contains(content, "postgres_data") {
		t.Error("DockerComposeTemplate() for PostgreSQL should have postgres_data volume")
	}
}

// TestDockerComposeTemplateSQLite tests docker-compose generation for SQLite
func TestDockerComposeTemplateSQLite(t *testing.T) {
	projectName := "test-sqlite-project"
	templates := NewProjectTemplatesWithDatabase(projectName, "sqlite")
	content := templates.DockerComposeTemplate()

	// Should NOT contain database service
	if strings.Contains(content, "db:") {
		t.Error("DockerComposeTemplate() for SQLite should not contain db service")
	}

	// Should NOT contain postgres or mysql images
	if strings.Contains(content, "postgres:") || strings.Contains(content, "mysql:") {
		t.Error("DockerComposeTemplate() for SQLite should not contain database images")
	}

	// Should contain api service only
	if !strings.Contains(content, "api:") {
		t.Error("DockerComposeTemplate() for SQLite should contain api service")
	}

	// Should have sqlite_data volume
	if !strings.Contains(content, "sqlite_data") {
		t.Error("DockerComposeTemplate() for SQLite should have sqlite_data volume")
	}

	// Should have .db file reference
	if !strings.Contains(content, ".db") {
		t.Error("DockerComposeTemplate() for SQLite should reference .db file")
	}

	// Should NOT have DB_HOST (SQLite uses local file)
	if strings.Contains(content, "DB_HOST:") {
		t.Error("DockerComposeTemplate() for SQLite should not have DB_HOST")
	}
}

// TestDockerComposeTemplateBackwardCompat tests backward compatibility (defaults to postgres)
func TestDockerComposeTemplateBackwardCompat(t *testing.T) {
	projectName := "test-legacy-project"
	templates := NewProjectTemplates(projectName) // Old function, should default to postgres
	content := templates.DockerComposeTemplate()

	// Should contain PostgreSQL (default)
	if !strings.Contains(content, "postgres:16-alpine") {
		t.Error("DockerComposeTemplate() with backward compat should default to PostgreSQL")
	}

	// Should have all required services and components
	requiredComponents := []string{
		"services:",
		"db:",
		"api:",
		"volumes:",
		"networks:",
	}
	for _, component := range requiredComponents {
		if !strings.Contains(content, component) {
			t.Errorf("DockerComposeTemplate() should contain %s", component)
		}
	}
}
