package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2ESQLiteProjectGeneration tests the complete project generation workflow for SQLite
// This test validates all acceptance criteria for Story 7.3 (SQLite Support)
func TestE2ESQLiteProjectGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	tmpDir := t.TempDir()
	projectName := "test-sqlite-project"
	projectPath := filepath.Join(tmpDir, projectName)

	// Generate project with SQLite
	t.Run("AC1_GenerateWithSQLite", func(t *testing.T) {
		if err := createProjectStructure(projectPath, TemplateFull); err != nil {
			t.Fatalf("Failed to create project structure: %v", err)
		}

		if err := generateProjectFiles(projectPath, projectName, DefaultTemplate, "sqlite"); err != nil {
			t.Fatalf("Failed to generate project files: %v", err)
		}

		t.Log("✅ AC#1: Project generated with --database=sqlite")
	})

	// AC1: Verify go.mod contains SQLite driver
	t.Run("AC1_VerifyGoMod", func(t *testing.T) {
		goModPath := filepath.Join(projectPath, "go.mod")
		content, err := os.ReadFile(goModPath)
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}

		goModContent := string(content)
		if !strings.Contains(goModContent, "gorm.io/driver/sqlite") {
			t.Error("go.mod should contain gorm.io/driver/sqlite driver")
		}
		if !strings.Contains(goModContent, "v1.5.4") {
			t.Error("go.mod should contain SQLite driver version v1.5.4")
		}

		t.Log("✅ AC#1: go.mod contains gorm.io/driver/sqlite v1.5.4")
	})

	// AC2: Verify DSN points to local .db file
	t.Run("AC2_VerifyDSN", func(t *testing.T) {
		// Check .env.example for DB_NAME configuration
		envPath := filepath.Join(projectPath, ".env.example")
		content, err := os.ReadFile(envPath)
		if err != nil {
			t.Fatalf("Failed to read .env.example: %v", err)
		}

		envContent := string(content)
		if !strings.Contains(envContent, "DB_NAME=") {
			t.Error(".env.example should contain DB_NAME configuration")
		}
		if !strings.Contains(envContent, ".db") {
			t.Error(".env.example should reference .db file extension")
		}

		// Verify DSN template format
		if !strings.Contains(envContent, projectName+".db") {
			t.Error(".env.example should contain project name with .db extension")
		}

		// Should NOT contain host/port/user/password for SQLite
		if strings.Contains(envContent, "DB_HOST") ||
			strings.Contains(envContent, "DB_PORT") ||
			strings.Contains(envContent, "DB_USER") ||
			strings.Contains(envContent, "DB_PASSWORD") {
			t.Error("SQLite .env should NOT contain DB_HOST/PORT/USER/PASSWORD")
		}

		t.Log("✅ AC#2: DSN points to local .db file (simplified configuration)")
	})

	// AC3: Verify docker-compose.yml has NO database service
	t.Run("AC3_VerifyDockerCompose", func(t *testing.T) {
		dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
		content, err := os.ReadFile(dockerComposePath)
		if err != nil {
			t.Fatalf("Failed to read docker-compose.yml: %v", err)
		}

		dockerComposeContent := string(content)

		// Should NOT contain database service indicators
		if strings.Contains(dockerComposeContent, "postgres:") ||
			strings.Contains(dockerComposeContent, "mysql:") {
			t.Error("docker-compose.yml should NOT contain database service for SQLite")
		}

		// Should NOT contain database service configuration
		if strings.Contains(dockerComposeContent, "POSTGRES_") ||
			strings.Contains(dockerComposeContent, "MYSQL_") {
			t.Error("docker-compose.yml should NOT contain database environment variables")
		}

		// Should contain sqlite_data volume
		if !strings.Contains(dockerComposeContent, "sqlite_data") {
			t.Error("docker-compose.yml should contain sqlite_data volume")
		}

		// Should contain comment about SQLite
		if !strings.Contains(dockerComposeContent, "SQLite") {
			t.Error("docker-compose.yml should contain comment explaining SQLite usage")
		}

		t.Log("✅ AC#3: docker-compose.yml has NO database service (SQLite is embedded)")
	})

	// AC4: Verify .gitignore excludes .db files
	t.Run("AC4_VerifyGitignore", func(t *testing.T) {
		gitignorePath := filepath.Join(projectPath, ".gitignore")
		content, err := os.ReadFile(gitignorePath)
		if err != nil {
			t.Fatalf("Failed to read .gitignore: %v", err)
		}

		gitignoreContent := string(content)
		if !strings.Contains(gitignoreContent, "*.db") {
			t.Error(".gitignore should contain *.db pattern")
		}
		if !strings.Contains(gitignoreContent, "*.db-shm") {
			t.Error(".gitignore should contain *.db-shm pattern")
		}
		if !strings.Contains(gitignoreContent, "*.db-wal") {
			t.Error(".gitignore should contain *.db-wal pattern")
		}

		t.Log("✅ AC#4: .gitignore excludes SQLite database files")
	})

	// AC4: Verify project compiles
	t.Run("AC4_Compilation", func(t *testing.T) {
		// Run go mod tidy
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectPath
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			t.Fatalf("go mod tidy failed: %v\nOutput: %s", err, string(output))
		}

		// Build the project
		buildCmd := exec.Command("go", "build", "./...")
		buildCmd.Dir = projectPath
		if output, err := buildCmd.CombinedOutput(); err != nil {
			t.Fatalf("Project compilation failed: %v\nOutput: %s", err, string(output))
		}

		t.Log("✅ AC#4: Project with SQLite compiles successfully")
	})

	// AC4: Verify database connection code
	t.Run("AC4_VerifyConnectionCode", func(t *testing.T) {
		dbFilePath := filepath.Join(projectPath, "internal", "infrastructure", "database", "database.go")
		content, err := os.ReadFile(dbFilePath)
		if err != nil {
			t.Fatalf("Failed to read database.go: %v", err)
		}

		dbContent := string(content)
		if !strings.Contains(dbContent, "gorm.io/driver/sqlite") {
			t.Error("database.go should import sqlite driver")
		}
		if !strings.Contains(dbContent, "sqlite.Open") {
			t.Error("database.go should use sqlite.Open() for connection")
		}

		t.Log("✅ AC#4: Database connection code uses sqlite.Open()")
	})

	t.Log("🎉 All SQLite E2E tests passed!")
}

// TestE2ESQLiteVsPostgresComparison validates that SQLite and PostgreSQL generate different configurations
func TestE2ESQLiteVsPostgresComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E comparison test in short mode")
	}

	tmpDir := t.TempDir()

	// Generate SQLite project
	sqliteProjectPath := filepath.Join(tmpDir, "sqlite-project")
	if err := createProjectStructure(sqliteProjectPath, TemplateFull); err != nil {
		t.Fatalf("Failed to create SQLite project: %v", err)
	}
	if err := generateProjectFiles(sqliteProjectPath, "sqlite-project", DefaultTemplate, "sqlite"); err != nil {
		t.Fatalf("Failed to generate SQLite project: %v", err)
	}

	// Generate PostgreSQL project
	postgresProjectPath := filepath.Join(tmpDir, "postgres-project")
	if err := createProjectStructure(postgresProjectPath, TemplateFull); err != nil {
		t.Fatalf("Failed to create PostgreSQL project: %v", err)
	}
	if err := generateProjectFiles(postgresProjectPath, "postgres-project", DefaultTemplate, "postgres"); err != nil {
		t.Fatalf("Failed to generate PostgreSQL project: %v", err)
	}

	// Compare docker-compose files
	t.Run("CompareDockerCompose", func(t *testing.T) {
		sqliteCompose, err := os.ReadFile(filepath.Join(sqliteProjectPath, "docker-compose.yml"))
		if err != nil {
			t.Fatalf("Failed to read SQLite docker-compose.yml: %v", err)
		}

		postgresCompose, err := os.ReadFile(filepath.Join(postgresProjectPath, "docker-compose.yml"))
		if err != nil {
			t.Fatalf("Failed to read PostgreSQL docker-compose.yml: %v", err)
		}

		sqliteContent := string(sqliteCompose)
		postgresContent := string(postgresCompose)

		// SQLite should NOT have db service
		if strings.Contains(sqliteContent, "postgres:16-alpine") {
			t.Error("SQLite docker-compose should NOT contain postgres image")
		}

		// PostgreSQL SHOULD have db service
		if !strings.Contains(postgresContent, "postgres:16-alpine") {
			t.Error("PostgreSQL docker-compose should contain postgres image")
		}

		t.Log("✅ Docker Compose configurations differ correctly between SQLite and PostgreSQL")
	})

	// Compare .env files
	t.Run("CompareEnvFiles", func(t *testing.T) {
		sqliteEnv, err := os.ReadFile(filepath.Join(sqliteProjectPath, ".env.example"))
		if err != nil {
			t.Fatalf("Failed to read SQLite .env.example: %v", err)
		}

		postgresEnv, err := os.ReadFile(filepath.Join(postgresProjectPath, ".env.example"))
		if err != nil {
			t.Fatalf("Failed to read PostgreSQL .env.example: %v", err)
		}

		sqliteContent := string(sqliteEnv)
		postgresContent := string(postgresEnv)

		// SQLite should have minimal config
		if strings.Contains(sqliteContent, "DB_HOST") {
			t.Error("SQLite .env should NOT contain DB_HOST")
		}

		// PostgreSQL should have full config
		if !strings.Contains(postgresContent, "DB_HOST") {
			t.Error("PostgreSQL .env should contain DB_HOST")
		}
		if !strings.Contains(postgresContent, "DB_PORT") {
			t.Error("PostgreSQL .env should contain DB_PORT")
		}

		t.Log("✅ Environment configurations differ correctly between SQLite and PostgreSQL")
	})

	t.Log("🎉 SQLite vs PostgreSQL comparison tests passed!")
}
