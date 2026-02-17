package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2EMySQLProjectGeneration tests end-to-end MySQL project generation
func TestE2EMySQLProjectGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Create temp directory for test project
	tmpDir := t.TempDir()
	projectName := "test-mysql-e2e"
	projectPath := filepath.Join(tmpDir, projectName)

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Generate full project with MySQL database
	database := "mysql"
	template := "full"

	// Create all necessary directories
	dirs := getDirectoriesForTemplate(template)
	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Generate project files
	if err := generateProjectFiles(projectPath, projectName, template, database, DefaultObservabilityLevel); err != nil {
		t.Fatalf("Failed to generate project files: %v", err)
	}

	t.Run("verify_go_mod_contains_mysql_driver", func(t *testing.T) {
		goModPath := filepath.Join(projectPath, "go.mod")
		content, err := os.ReadFile(goModPath)
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}

		goModContent := string(content)
		if !strings.Contains(goModContent, "gorm.io/driver/mysql") {
			t.Error("go.mod should contain gorm.io/driver/mysql")
		}

		// Should contain version number
		if !strings.Contains(goModContent, "v1.5.2") {
			t.Error("go.mod should contain MySQL driver version v1.5.2")
		}

		// Should NOT contain postgres driver anymore (we're generating MySQL project)
		// Actually, templates generate common GORM, so this check is optional
		// But at minimum MySQL driver should be present
	})

	t.Run("verify_docker_compose_contains_mysql", func(t *testing.T) {
		dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
		content, err := os.ReadFile(dockerComposePath)
		if err != nil {
			t.Fatalf("Failed to read docker-compose.yml: %v", err)
		}

		dockerComposeContent := string(content)

		// Check MySQL image
		if !strings.Contains(dockerComposeContent, "mysql:8.0") {
			t.Error("docker-compose.yml should contain mysql:8.0 image")
		}

		// Check MySQL environment variables
		requiredEnvVars := []string{
			"MYSQL_ROOT_PASSWORD",
			"MYSQL_DATABASE",
			"MYSQL_USER",
			"MYSQL_PASSWORD",
		}
		for _, envVar := range requiredEnvVars {
			if !strings.Contains(dockerComposeContent, envVar) {
				t.Errorf("docker-compose.yml should contain %s", envVar)
			}
		}

		// Check MySQL port
		if !strings.Contains(dockerComposeContent, "3306") {
			t.Error("docker-compose.yml should use port 3306 for MySQL")
		}

		// Check MySQL healthcheck
		if !strings.Contains(dockerComposeContent, "mysqladmin") {
			t.Error("docker-compose.yml should have mysqladmin healthcheck")
		}

		// Check MySQL volume
		if !strings.Contains(dockerComposeContent, "mysql_data") {
			t.Error("docker-compose.yml should have mysql_data volume")
		}
	})

	t.Run("verify_env_file_contains_mysql_config", func(t *testing.T) {
		envPath := filepath.Join(projectPath, ".env.example")
		content, err := os.ReadFile(envPath)
		if err != nil {
			t.Fatalf("Failed to read .env.example: %v", err)
		}

		envContent := string(content)

		// Check MySQL-specific configuration
		if !strings.Contains(envContent, "DB_PORT=3306") {
			t.Error(".env.example should have DB_PORT=3306 for MySQL")
		}

		if !strings.Contains(envContent, "DB_USER=root") {
			t.Error(".env.example should have DB_USER=root for MySQL")
		}

		// Should NOT have PostgreSQL configuration
		if strings.Contains(envContent, "DB_PORT=5432") {
			t.Error(".env.example should NOT have DB_PORT=5432 (PostgreSQL) for MySQL project")
		}

		if strings.Contains(envContent, "DB_USER=postgres") {
			t.Error(".env.example should NOT have DB_USER=postgres for MySQL project")
		}

		// Should indicate MySQL in comment
		if !strings.Contains(envContent, "MySQL") {
			t.Error(".env.example should mention MySQL in database configuration comment")
		}
	})

	t.Run("verify_database_go_uses_mysql_driver", func(t *testing.T) {
		databaseGoPath := filepath.Join(projectPath, "internal", "infrastructure", "database", "database.go")
		content, err := os.ReadFile(databaseGoPath)
		if err != nil {
			t.Fatalf("Failed to read database.go: %v", err)
		}

		databaseGoContent := string(content)

		// Check for MySQL import
		if !strings.Contains(databaseGoContent, "gorm.io/driver/mysql") {
			t.Error("database.go should import gorm.io/driver/mysql")
		}

		// Check for mysql.Open() usage
		if !strings.Contains(databaseGoContent, "mysql.Open(dsn)") {
			t.Error("database.go should use mysql.Open(dsn)")
		}

		// Check for parseTime parameter in DSN
		if !strings.Contains(databaseGoContent, "parseTime=True") {
			t.Error("MySQL DSN should include parseTime=True parameter")
		}

		// Check for charset=utf8mb4
		if !strings.Contains(databaseGoContent, "charset=utf8mb4") {
			t.Error("MySQL DSN should include charset=utf8mb4 parameter")
		}

		// Should NOT have postgres driver
		if strings.Contains(databaseGoContent, "gorm.io/driver/postgres") {
			t.Error("database.go should NOT have PostgreSQL driver for MySQL project")
		}

		if strings.Contains(databaseGoContent, "postgres.Open(dsn)") {
			t.Error("database.go should NOT use postgres.Open() for MySQL project")
		}
	})

	t.Run("verify_project_compiles", func(t *testing.T) {
		// First run go mod tidy to download dependencies
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectPath
		tidyOutput, tidyErr := tidyCmd.CombinedOutput()
		if tidyErr != nil {
			t.Fatalf("go mod tidy failed.\nError: %v\nOutput:\n%s", tidyErr, string(tidyOutput))
		}

		// Try to build the project
		cmd := exec.Command("go", "build", "./...")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Project should compile successfully.\nError: %v\nOutput:\n%s", err, string(output))
		}

		t.Logf("Project compiled successfully")
	})

	t.Run("verify_go_mod_tidy_works", func(t *testing.T) {
		// Ensure dependencies are valid
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("go mod tidy should succeed.\nError: %v\nOutput:\n%s", err, string(output))
		}

		t.Logf("go mod tidy completed successfully")
	})
}

// TestE2EMySQLVsPostgresComparison compares MySQL and PostgreSQL generated projects
func TestE2EMySQLVsPostgresComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E comparison test in short mode")
	}

	// This test generates both MySQL and PostgreSQL projects and compares them
	tmpDir := t.TempDir()
	template := "minimal" // Use minimal for faster testing

	// Generate MySQL project
	mysqlProjectName := "test-mysql-compare"
	mysqlPath := filepath.Join(tmpDir, mysqlProjectName)
	if err := os.MkdirAll(mysqlPath, 0755); err != nil {
		t.Fatalf("Failed to create MySQL project directory: %v", err)
	}
	dirs := getDirectoriesForTemplate(template)
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(mysqlPath, dir), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
	}
	if err := generateProjectFiles(mysqlPath, mysqlProjectName, template, "mysql", DefaultObservabilityLevel); err != nil {
		t.Fatalf("Failed to generate MySQL project: %v", err)
	}

	// Generate PostgreSQL project
	pgProjectName := "test-postgres-compare"
	pgPath := filepath.Join(tmpDir, pgProjectName)
	if err := os.MkdirAll(pgPath, 0755); err != nil {
		t.Fatalf("Failed to create PostgreSQL project directory: %v", err)
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(pgPath, dir), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
	}
	if err := generateProjectFiles(pgPath, pgProjectName, template, "postgres", DefaultObservabilityLevel); err != nil {
		t.Fatalf("Failed to generate PostgreSQL project: %v", err)
	}

	// Compare go.mod files
	t.Run("compare_go_mod_drivers", func(t *testing.T) {
		mysqlGoMod, err := os.ReadFile(filepath.Join(mysqlPath, "go.mod"))
		if err != nil {
			t.Fatalf("Failed to read MySQL go.mod: %v", err)
		}

		pgGoMod, err := os.ReadFile(filepath.Join(pgPath, "go.mod"))
		if err != nil {
			t.Fatalf("Failed to read PostgreSQL go.mod: %v", err)
		}

		// MySQL should have mysql driver
		if !strings.Contains(string(mysqlGoMod), "gorm.io/driver/mysql") {
			t.Error("MySQL project should have MySQL driver in go.mod")
		}

		// PostgreSQL should have postgres driver
		if !strings.Contains(string(pgGoMod), "gorm.io/driver/postgres") {
			t.Error("PostgreSQL project should have PostgreSQL driver in go.mod")
		}

		// They should have different drivers
		if strings.Contains(string(mysqlGoMod), "gorm.io/driver/postgres") {
			t.Error("MySQL project should NOT have PostgreSQL driver")
		}
		if strings.Contains(string(pgGoMod), "gorm.io/driver/mysql") {
			t.Error("PostgreSQL project should NOT have MySQL driver")
		}
	})

	// Both should compile
	t.Run("both_projects_compile", func(t *testing.T) {
		// Tidy and build MySQL project
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = mysqlPath
		if tidyOutput, tidyErr := tidyCmd.CombinedOutput(); tidyErr != nil {
			t.Errorf("MySQL go mod tidy failed.\nError: %v\nOutput:\n%s", tidyErr, string(tidyOutput))
		}

		cmd := exec.Command("go", "build", "./...")
		cmd.Dir = mysqlPath
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("MySQL project should compile.\nError: %v\nOutput:\n%s", err, string(output))
		}

		// Tidy and build PostgreSQL project
		tidyCmd = exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = pgPath
		if tidyOutput, tidyErr := tidyCmd.CombinedOutput(); tidyErr != nil {
			t.Errorf("PostgreSQL go mod tidy failed.\nError: %v\nOutput:\n%s", tidyErr, string(tidyOutput))
		}

		cmd = exec.Command("go", "build", "./...")
		cmd.Dir = pgPath
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("PostgreSQL project should compile.\nError: %v\nOutput:\n%s", err, string(output))
		}
	})
}
