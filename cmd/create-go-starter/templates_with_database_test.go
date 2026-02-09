package main

import (
	"strings"
	"testing"
)

// TestNewProjectTemplatesWithDatabase tests creating templates with specific database
func TestNewProjectTemplatesWithDatabase(t *testing.T) {
	tests := []struct {
		name         string
		projectName  string
		database     string
		wantDatabase string
	}{
		{
			name:         "MySQL database",
			projectName:  "test-project",
			database:     "mysql",
			wantDatabase: "mysql",
		},
		{
			name:         "PostgreSQL database (explicit)",
			projectName:  "test-project",
			database:     "postgres",
			wantDatabase: "postgres",
		},
		{
			name:         "PostgreSQL database (empty defaults to postgres)",
			projectName:  "test-project",
			database:     "",
			wantDatabase: "postgres",
		},
		{
			name:         "SQLite database",
			projectName:  "test-project",
			database:     "sqlite",
			wantDatabase: "sqlite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := NewProjectTemplatesWithDatabase(tt.projectName, tt.database)
			if templates.database != tt.wantDatabase {
				t.Errorf("NewProjectTemplatesWithDatabase(%q, %q).database = %q, want %q",
					tt.projectName, tt.database, templates.database, tt.wantDatabase)
			}
			if templates.projectName != tt.projectName {
				t.Errorf("NewProjectTemplatesWithDatabase(%q, %q).projectName = %q, want %q",
					tt.projectName, tt.database, templates.projectName, tt.projectName)
			}
		})
	}
}

// TestGoModTemplateWithDatabase tests that go.mod includes correct driver
func TestGoModTemplateWithDatabase(t *testing.T) {
	tests := []struct {
		name           string
		projectName    string
		database       string
		wantDriverPath string
	}{
		{
			name:           "MySQL driver in go.mod",
			projectName:    "test-mysql-project",
			database:       "mysql",
			wantDriverPath: "gorm.io/driver/mysql",
		},
		{
			name:           "PostgreSQL driver in go.mod (default)",
			projectName:    "test-postgres-project",
			database:       "postgres",
			wantDriverPath: "gorm.io/driver/postgres",
		},
		{
			name:           "PostgreSQL driver in go.mod (backward compat)",
			projectName:    "test-compat-project",
			database:       "",
			wantDriverPath: "gorm.io/driver/postgres",
		},
		{
			name:           "SQLite driver in go.mod",
			projectName:    "test-sqlite-project",
			database:       "sqlite",
			wantDriverPath: "gorm.io/driver/sqlite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := NewProjectTemplatesWithDatabase(tt.projectName, tt.database)
			goModContent := templates.GoModTemplate()

			// Check that it contains the right driver
			if !strings.Contains(goModContent, tt.wantDriverPath) {
				t.Errorf("GoModTemplate() for database %q does not contain %q",
					tt.database, tt.wantDriverPath)
			}

			// Check that it contains the project name
			if !strings.Contains(goModContent, "module "+tt.projectName) {
				t.Errorf("GoModTemplate() does not contain module %q", tt.projectName)
			}

			// Check that it still contains other essential dependencies
			essentialDeps := []string{
				"github.com/gofiber/fiber/v2",
				"gorm.io/gorm",
				"go.uber.org/fx",
			}
			for _, dep := range essentialDeps {
				if !strings.Contains(goModContent, dep) {
					t.Errorf("GoModTemplate() missing essential dependency %q", dep)
				}
			}
		})
	}
}

// TestBackwardCompatibility tests that NewProjectTemplates still works (defaults to postgres)
func TestBackwardCompatibility(t *testing.T) {
	projectName := "legacy-project"
	templates := NewProjectTemplates(projectName)

	// Should default to postgres
	if templates.database != "postgres" {
		t.Errorf("NewProjectTemplates(%q).database = %q, want %q",
			projectName, templates.database, "postgres")
	}

	// go.mod should contain postgres driver
	goModContent := templates.GoModTemplate()
	if !strings.Contains(goModContent, "gorm.io/driver/postgres") {
		t.Error("Backward compatibility broken: go.mod does not contain postgres driver")
	}
}
