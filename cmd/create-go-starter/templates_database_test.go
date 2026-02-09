package main

import (
	"strings"
	"testing"
)

// TestGoModDependencies tests that the correct database driver is returned
func TestGoModDependencies(t *testing.T) {
	tests := []struct {
		name     string
		database string
		want     string
	}{
		{
			name:     "MySQL driver",
			database: "mysql",
			want:     "gorm.io/driver/mysql v1.5.2",
		},
		{
			name:     "PostgreSQL driver (default)",
			database: "postgres",
			want:     "gorm.io/driver/postgres v1.5.4",
		},
		{
			name:     "PostgreSQL driver (empty string defaults to postgres)",
			database: "",
			want:     "gorm.io/driver/postgres v1.5.4",
		},
		{
			name:     "SQLite driver",
			database: "sqlite",
			want:     "gorm.io/driver/sqlite v1.5.4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GoModDependencies(tt.database)
			if got != tt.want {
				t.Errorf("GoModDependencies(%q) = %q, want %q", tt.database, got, tt.want)
			}
		})
	}
}

// TestDatabaseDSN tests that the correct DSN format is generated
func TestDatabaseDSN(t *testing.T) {
	tests := []struct {
		name        string
		database    string
		projectName string
		wantSubstr  []string // Check that DSN contains these substrings
	}{
		{
			name:        "MySQL DSN",
			database:    "mysql",
			projectName: "test-project",
			wantSubstr: []string{
				"${DB_USER}:${DB_PASSWORD}",
				"@tcp(${DB_HOST}:${DB_PORT})",
				"/${DB_NAME}",
				"parseTime=True",
				"charset=utf8mb4",
			},
		},
		{
			name:        "PostgreSQL DSN (default)",
			database:    "postgres",
			projectName: "test-project",
			wantSubstr: []string{
				"host=${DB_HOST}",
				"user=${DB_USER}",
				"password=${DB_PASSWORD}",
				"dbname=${DB_NAME}",
				"port=${DB_PORT}",
				"sslmode=disable",
			},
		},
		{
			name:        "SQLite DSN",
			database:    "sqlite",
			projectName: "test-project",
			wantSubstr: []string{
				"./${DB_NAME}.db",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseDSN(tt.database, tt.projectName)
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(got, substr) {
					t.Errorf("DatabaseDSN(%q, %q) = %q, want to contain %q", tt.database, tt.projectName, got, substr)
				}
			}
		})
	}
}

// TestDatabaseConnectionCode tests that the correct GORM connection code is generated
func TestDatabaseConnectionCode(t *testing.T) {
	tests := []struct {
		name     string
		database string
		want     string
	}{
		{
			name:     "MySQL connection",
			database: "mysql",
			want:     "mysql.Open(dsn)",
		},
		{
			name:     "PostgreSQL connection (default)",
			database: "postgres",
			want:     "postgres.Open(dsn)",
		},
		{
			name:     "PostgreSQL connection (empty defaults to postgres)",
			database: "",
			want:     "postgres.Open(dsn)",
		},
		{
			name:     "SQLite connection",
			database: "sqlite",
			want:     "sqlite.Open(dsn)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseConnectionCode(tt.database)
			if got != tt.want {
				t.Errorf("DatabaseConnectionCode(%q) = %q, want %q", tt.database, got, tt.want)
			}
		})
	}
}

// TestDatabaseImportPath tests that the correct import path is returned
func TestDatabaseImportPath(t *testing.T) {
	tests := []struct {
		name     string
		database string
		want     string
	}{
		{
			name:     "MySQL import",
			database: "mysql",
			want:     "gorm.io/driver/mysql",
		},
		{
			name:     "PostgreSQL import (default)",
			database: "postgres",
			want:     "gorm.io/driver/postgres",
		},
		{
			name:     "PostgreSQL import (empty defaults to postgres)",
			database: "",
			want:     "gorm.io/driver/postgres",
		},
		{
			name:     "SQLite import",
			database: "sqlite",
			want:     "gorm.io/driver/sqlite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseImportPath(tt.database)
			if got != tt.want {
				t.Errorf("DatabaseImportPath(%q) = %q, want %q", tt.database, got, tt.want)
			}
		})
	}
}

// TestDatabaseDockerService tests docker-compose service generation
func TestDatabaseDockerService(t *testing.T) {
	tests := []struct {
		name        string
		database    string
		projectName string
		wantSubstr  []string
	}{
		{
			name:        "MySQL docker service",
			database:    "mysql",
			projectName: "test-project",
			wantSubstr: []string{
				"image: mysql:8.0",
				"MYSQL_ROOT_PASSWORD",
				"MYSQL_DATABASE",
				"MYSQL_USER",
				"MYSQL_PASSWORD",
				"3306:3306",
				"mysqladmin",
				"mysql_data",
			},
		},
		{
			name:        "PostgreSQL docker service (default)",
			database:    "postgres",
			projectName: "test-project",
			wantSubstr: []string{
				"image: postgres:16-alpine",
				"POSTGRES_DB",
				"POSTGRES_USER",
				"POSTGRES_PASSWORD",
				"5432:5432",
				"pg_isready",
				"postgres_data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseDockerService(tt.database, tt.projectName)
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(got, substr) {
					t.Errorf("DatabaseDockerService(%q, %q) missing expected substring %q", tt.database, tt.projectName, substr)
				}
			}
		})
	}
}

// TestDatabaseEnvVars tests environment variables generation
func TestDatabaseEnvVars(t *testing.T) {
	tests := []struct {
		name        string
		database    string
		projectName string
		wantSubstr  []string
	}{
		{
			name:        "MySQL env vars",
			database:    "mysql",
			projectName: "test-project",
			wantSubstr: []string{
				"DB_HOST=localhost",
				"DB_PORT=3306",
				"DB_USER=root",
				"DB_PASSWORD=",
				"DB_NAME=",
			},
		},
		{
			name:        "PostgreSQL env vars (default)",
			database:    "postgres",
			projectName: "test-project",
			wantSubstr: []string{
				"DB_HOST=localhost",
				"DB_PORT=5432",
				"DB_USER=postgres",
				"DB_PASSWORD=",
				"DB_NAME=",
			},
		},
		{
			name:        "SQLite env vars",
			database:    "sqlite",
			projectName: "test-project",
			wantSubstr: []string{
				"DB_NAME=",
				".db",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseEnvVars(tt.database, tt.projectName)
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(got, substr) {
					t.Errorf("DatabaseEnvVars(%q, %q) missing expected substring %q", tt.database, tt.projectName, substr)
				}
			}
		})
	}
}
