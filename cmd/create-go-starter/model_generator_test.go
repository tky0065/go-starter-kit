package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateModelTemplate tests model file template generation (AC: 1)
func TestGenerateModelTemplate(t *testing.T) {
	fields := []FieldDefinition{
		{Name: "Title", Type: "string", GORMTags: []string{"not_null"}, JSONName: "title"},
		{Name: "Completed", Type: "bool", GORMTags: nil, JSONName: "completed"},
	}

	content := generateModelTemplate("Todo", fields, nil)

	// Verify package
	if !strings.Contains(content, "package models") {
		t.Error("Expected 'package models'")
	}

	// Verify imports
	if !strings.Contains(content, `"time"`) {
		t.Error("Expected time import")
	}
	if !strings.Contains(content, `"gorm.io/gorm"`) {
		t.Error("Expected gorm import")
	}

	// Verify struct declaration
	if !strings.Contains(content, "type Todo struct") {
		t.Error("Expected 'type Todo struct'")
	}

	// Verify base fields
	if !strings.Contains(content, "ID") {
		t.Error("Expected ID field")
	}
	if !strings.Contains(content, "CreatedAt") {
		t.Error("Expected CreatedAt field")
	}
	if !strings.Contains(content, "UpdatedAt") {
		t.Error("Expected UpdatedAt field")
	}
	if !strings.Contains(content, "DeletedAt") {
		t.Error("Expected DeletedAt field")
	}

	// Verify user fields
	if !strings.Contains(content, "Title") {
		t.Error("Expected Title field")
	}
	if !strings.Contains(content, "Completed") {
		t.Error("Expected Completed field")
	}

	// Verify GORM tags
	if !strings.Contains(content, `gorm:"not null"`) {
		t.Error("Expected 'not null' GORM tag for Title")
	}
	if !strings.Contains(content, `gorm:"default:false"`) {
		t.Error("Expected 'default:false' GORM tag for Completed (bool)")
	}

	// Verify JSON tags
	if !strings.Contains(content, `json:"title"`) {
		t.Error("Expected json tag for title")
	}
	if !strings.Contains(content, `json:"completed"`) {
		t.Error("Expected json tag for completed")
	}
	if !strings.Contains(content, `json:"id"`) {
		t.Error("Expected json tag for id")
	}
	if !strings.Contains(content, `json:"created_at"`) {
		t.Error("Expected json tag for created_at")
	}
}

// TestGenerateRepositoryInterfaceTemplate tests interface generation (AC: 2)
func TestGenerateRepositoryInterfaceTemplate(t *testing.T) {
	content := generateRepositoryInterfaceTemplate("github.com/user/project", "Todo", nil)

	// Verify package
	if !strings.Contains(content, "package interfaces") {
		t.Error("Expected 'package interfaces'")
	}

	// Verify import
	if !strings.Contains(content, `"github.com/user/project/internal/models"`) {
		t.Error("Expected models import with correct module path")
	}

	// Verify interface declaration
	if !strings.Contains(content, "type TodoRepository interface") {
		t.Error("Expected 'type TodoRepository interface'")
	}

	// Verify CRUD methods
	methods := []string{
		"Create(ctx context.Context, todo *models.Todo) error",
		"FindByID(ctx context.Context, id uint) (*models.Todo, error)",
		"FindAll(ctx context.Context, page, limit int) ([]*models.Todo, int64, error)",
		"Update(ctx context.Context, todo *models.Todo) error",
		"Delete(ctx context.Context, id uint) error",
	}
	for _, m := range methods {
		if !strings.Contains(content, m) {
			t.Errorf("Expected method %q in interface", m)
		}
	}
}

// TestGenerateRepositoryImplTemplate tests repository implementation generation (AC: 3)
func TestGenerateRepositoryImplTemplate(t *testing.T) {
	content := generateRepositoryImplTemplate("github.com/user/project", "Todo", nil)

	// Verify package
	if !strings.Contains(content, "package repository") {
		t.Error("Expected 'package repository'")
	}

	// Verify imports
	if !strings.Contains(content, `"github.com/user/project/internal/interfaces"`) {
		t.Error("Expected interfaces import")
	}
	if !strings.Contains(content, `"github.com/user/project/internal/models"`) {
		t.Error("Expected models import")
	}

	// Verify struct and constructor
	if !strings.Contains(content, "type todoRepository struct") {
		t.Error("Expected todoRepository struct")
	}
	if !strings.Contains(content, "func NewTodoRepository(db *gorm.DB) interfaces.TodoRepository") {
		t.Error("Expected NewTodoRepository constructor")
	}

	// Verify CRUD methods
	if !strings.Contains(content, "func (r *todoRepository) Create") {
		t.Error("Expected Create method")
	}
	if !strings.Contains(content, "func (r *todoRepository) FindByID") {
		t.Error("Expected FindByID method")
	}
	if !strings.Contains(content, "func (r *todoRepository) FindAll") {
		t.Error("Expected FindAll method")
	}
	if !strings.Contains(content, "func (r *todoRepository) Update") {
		t.Error("Expected Update method")
	}
	if !strings.Contains(content, "func (r *todoRepository) Delete") {
		t.Error("Expected Delete method")
	}

	// Verify ErrRecordNotFound handling
	if !strings.Contains(content, "gorm.ErrRecordNotFound") {
		t.Error("Expected ErrRecordNotFound handling in FindByID")
	}

	// Verify pagination pattern
	if !strings.Contains(content, "Offset(offset).Limit(limit)") {
		t.Error("Expected pagination with offset/limit")
	}
}

// TestBuildGORMTag tests GORM tag building
func TestBuildGORMTag(t *testing.T) {
	tests := []struct {
		name     string
		field    FieldDefinition
		expected string
	}{
		{
			name:     "no modifiers string",
			field:    FieldDefinition{Name: "Title", Type: "string"},
			expected: `gorm:""`,
		},
		{
			name:     "not_null modifier",
			field:    FieldDefinition{Name: "Title", Type: "string", GORMTags: []string{"not_null"}},
			expected: `gorm:"not null"`,
		},
		{
			name:     "unique modifier",
			field:    FieldDefinition{Name: "Email", Type: "string", GORMTags: []string{"unique"}},
			expected: `gorm:"uniqueIndex"`,
		},
		{
			name:     "index modifier",
			field:    FieldDefinition{Name: "Category", Type: "string", GORMTags: []string{"index"}},
			expected: `gorm:"index"`,
		},
		{
			name:     "multiple modifiers",
			field:    FieldDefinition{Name: "Email", Type: "string", GORMTags: []string{"unique", "not_null"}},
			expected: `gorm:"uniqueIndex;not null"`,
		},
		{
			name:     "bool default",
			field:    FieldDefinition{Name: "Active", Type: "bool"},
			expected: `gorm:"default:false"`,
		},
		{
			name:     "bool with index",
			field:    FieldDefinition{Name: "Active", Type: "bool", GORMTags: []string{"index"}},
			expected: `gorm:"index;default:false"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildGORMTag(tt.field)
			if got != tt.expected {
				t.Errorf("buildGORMTag() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestUpdateAutoMigrate tests the AutoMigrate modification (AC: 5)
func TestUpdateAutoMigrate(t *testing.T) {
	tmpDir := t.TempDir()
	dbDir := filepath.Join(tmpDir, "internal", "infrastructure", "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package database

func setup(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		return err
	}
}
`
	if err := os.WriteFile(filepath.Join(dbDir, "database.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateAutoMigrate(tmpDir, "Todo"); err != nil {
		t.Fatalf("updateAutoMigrate() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dbDir, "database.go"))
	if err != nil {
		t.Fatalf("Failed to read database.go after update: %v", err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "&models.Todo{}") {
		t.Error("Expected &models.Todo{} in AutoMigrate")
	}

	// Verify existing models are preserved
	if !strings.Contains(contentStr, "&models.User{}") {
		t.Error("Expected &models.User{} to be preserved")
	}
	if !strings.Contains(contentStr, "&models.RefreshToken{}") {
		t.Error("Expected &models.RefreshToken{} to be preserved")
	}
}

// TestUpdateAutoMigrateIdempotent tests that running twice doesn't duplicate
func TestUpdateAutoMigrateIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbDir := filepath.Join(tmpDir, "internal", "infrastructure", "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package database

func setup(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		return err
	}
}
`
	if err := os.WriteFile(filepath.Join(dbDir, "database.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	// Should not duplicate
	if err := updateAutoMigrate(tmpDir, "Todo"); err != nil {
		t.Fatalf("updateAutoMigrate() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dbDir, "database.go"))
	if err != nil {
		t.Fatalf("Failed to read database.go after idempotent update: %v", err)
	}
	count := strings.Count(string(content), "&models.Todo{}")
	if count != 1 {
		t.Errorf("Expected exactly 1 occurrence of &models.Todo{}, got %d", count)
	}
}

// TestUpdateRepositoryModule tests the fx module modification (AC: 4)
func TestUpdateRepositoryModule(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "internal", "adapters", "repository")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package repository

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"github.com/test/myproject/internal/interfaces"
)

var Module = fx.Module("repository",
	fx.Provide(func(db *gorm.DB) interfaces.UserRepository {
		return NewUserRepository(db)
	}),
)
`
	if err := os.WriteFile(filepath.Join(repoDir, "module.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateRepositoryModule(tmpDir, "github.com/test/myproject", "Todo"); err != nil {
		t.Fatalf("updateRepositoryModule() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(repoDir, "module.go"))
	if err != nil {
		t.Fatalf("Failed to read module.go after update: %v", err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "NewTodoRepository") {
		t.Error("Expected NewTodoRepository in module")
	}
	if !strings.Contains(contentStr, "interfaces.TodoRepository") {
		t.Error("Expected interfaces.TodoRepository in module")
	}

	// Verify original content preserved
	if !strings.Contains(contentStr, "NewUserRepository") {
		t.Error("Expected NewUserRepository to be preserved")
	}
}

// TestUpdateRepositoryModuleIdempotent tests that running twice doesn't duplicate
func TestUpdateRepositoryModuleIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "internal", "adapters", "repository")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package repository

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"github.com/test/myproject/internal/interfaces"
)

var Module = fx.Module("repository",
	fx.Provide(func(db *gorm.DB) interfaces.UserRepository {
		return NewUserRepository(db)
	}),
	fx.Provide(func(db *gorm.DB) interfaces.TodoRepository {
		return NewTodoRepository(db)
	}),
)
`
	if err := os.WriteFile(filepath.Join(repoDir, "module.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	// Should not duplicate
	if err := updateRepositoryModule(tmpDir, "github.com/test/myproject", "Todo"); err != nil {
		t.Fatalf("updateRepositoryModule() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(repoDir, "module.go"))
	if err != nil {
		t.Fatalf("Failed to read module.go after idempotent update: %v", err)
	}
	count := strings.Count(string(content), "NewTodoRepository")
	if count != 1 {
		t.Errorf("Expected exactly 1 occurrence of NewTodoRepository, got %d", count)
	}
}

// TestPluralize tests the pluralization function
func TestPluralize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"todo", "todos"},
		{"product", "products"},
		{"category", "categories"},
		{"bus", "buses"},
		{"box", "boxes"},
		{"dish", "dishes"},
		{"match", "matches"},
		{"quiz", "quizes"},
		{"key", "keys"},
		{"day", "days"},
		{"", ""},
		// Irregular plurals (critical for proper REST APIs)
		{"person", "people"},
		{"child", "children"},
		{"man", "men"},
		{"woman", "women"},
		{"foot", "feet"},
		{"tooth", "teeth"},
		{"mouse", "mice"},
		// f/fe → ves
		{"knife", "knives"},
		{"wife", "wives"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := pluralize(tt.input)
			if got != tt.expected {
				t.Errorf("pluralize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestGenerateServiceTemplate tests service file template generation (Story 8.3 AC: 1)
func TestGenerateServiceTemplate(t *testing.T) {
	content := generateServiceTemplate("github.com/user/project", "Todo", nil)

	// Verify package
	if !strings.Contains(content, "package todo") {
		t.Error("Expected 'package todo'")
	}

	// Verify imports
	if !strings.Contains(content, `"github.com/user/project/internal/interfaces"`) {
		t.Error("Expected interfaces import")
	}
	if !strings.Contains(content, `"github.com/user/project/internal/models"`) {
		t.Error("Expected models import")
	}

	// Verify struct and constructor
	if !strings.Contains(content, "type Service struct") {
		t.Error("Expected Service struct")
	}
	if !strings.Contains(content, "repo interfaces.TodoRepository") {
		t.Error("Expected repository interface injection")
	}
	if !strings.Contains(content, "func NewService(repo interfaces.TodoRepository) *Service") {
		t.Error("Expected NewService constructor")
	}

	// Verify CRUD methods
	methods := []string{
		"func (s *Service) Create(",
		"func (s *Service) GetByID(",
		"func (s *Service) GetAll(",
		"func (s *Service) Update(",
		"func (s *Service) Delete(",
	}
	for _, m := range methods {
		if !strings.Contains(content, m) {
			t.Errorf("Expected method %q in service", m)
		}
	}

	// Verify pagination limits
	if !strings.Contains(content, "limit > 100") {
		t.Error("Expected pagination limit cap at 100")
	}
}

// TestGenerateServiceModuleTemplate tests the fx module generation for service
func TestGenerateServiceModuleTemplate(t *testing.T) {
	content := generateServiceModuleTemplate("Todo")

	if !strings.Contains(content, "package todo") {
		t.Error("Expected 'package todo'")
	}
	if !strings.Contains(content, `fx.Module("todo"`) {
		t.Error("Expected fx.Module(\"todo\")")
	}
	if !strings.Contains(content, "fx.Provide(NewService)") {
		t.Error("Expected fx.Provide(NewService)")
	}
}

// TestGenerateHandlerTemplate tests handler template generation (Story 8.3 AC: 2, 5)
func TestGenerateHandlerTemplate(t *testing.T) {
	fields := []FieldDefinition{
		{Name: "Title", Type: "string", GORMTags: []string{"not_null"}, JSONName: "title"},
		{Name: "Completed", Type: "bool", GORMTags: nil, JSONName: "completed"},
	}

	content := generateHandlerTemplate("github.com/user/project", "Todo", fields, nil)

	// Verify package
	if !strings.Contains(content, "package handlers") {
		t.Error("Expected 'package handlers'")
	}

	// Verify imports
	if !strings.Contains(content, `"github.com/user/project/internal/domain/todo"`) {
		t.Error("Expected domain/todo import")
	}
	if !strings.Contains(content, `"github.com/user/project/internal/models"`) {
		t.Error("Expected models import")
	}

	// Verify DTOs (AC: 5)
	if !strings.Contains(content, "type CreateTodoRequest struct") {
		t.Error("Expected CreateTodoRequest DTO")
	}
	if !strings.Contains(content, "type UpdateTodoRequest struct") {
		t.Error("Expected UpdateTodoRequest DTO")
	}
	// Verify validation tags
	if !strings.Contains(content, `validate:"required,min=1,max=255"`) {
		t.Error("Expected validation tag on string field")
	}
	if !strings.Contains(content, `validate:"omitempty`) {
		t.Error("Expected omitempty validation on update DTO fields")
	}

	// Verify handler struct
	if !strings.Contains(content, "type TodoHandler struct") {
		t.Error("Expected TodoHandler struct")
	}
	if !strings.Contains(content, "func NewTodoHandler(service *todo.Service) *TodoHandler") {
		t.Error("Expected NewTodoHandler constructor")
	}

	// Verify 5 CRUD endpoints
	endpoints := []string{
		"func (h *TodoHandler) CreateTodo(",
		"func (h *TodoHandler) GetTodo(",
		"func (h *TodoHandler) GetAllTodos(",
		"func (h *TodoHandler) UpdateTodo(",
		"func (h *TodoHandler) DeleteTodo(",
	}
	for _, e := range endpoints {
		if !strings.Contains(content, e) {
			t.Errorf("Expected endpoint %q", e)
		}
	}

	// Verify Swagger annotations
	swaggerTags := []string{
		"@Summary Create a new todo",
		"@Summary Get a todo by ID",
		"@Summary Get all todos",
		"@Summary Update a todo",
		"@Summary Delete a todo",
		"@Router /api/v1/todos",
	}
	for _, tag := range swaggerTags {
		if !strings.Contains(content, tag) {
			t.Errorf("Expected Swagger annotation %q", tag)
		}
	}

	// Verify JSON response format
	if !strings.Contains(content, `"status": "success"`) {
		t.Error("Expected standard success response")
	}
	if !strings.Contains(content, `"status":  "error"`) {
		t.Error("Expected standard error response")
	}

	// Verify body parsing and validation
	if !strings.Contains(content, "c.BodyParser(&req)") {
		t.Error("Expected body parsing")
	}
	if !strings.Contains(content, "h.validate.Struct(req)") {
		t.Error("Expected struct validation")
	}
}

// TestUpdateHandlerModule tests the handler module modification (Story 8.3 AC: 4)
func TestUpdateHandlerModule(t *testing.T) {
	tmpDir := t.TempDir()
	handlerDir := filepath.Join(tmpDir, "internal", "adapters", "handlers")
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package handlers

import (
	"go.uber.org/fx"
	"github.com/test/myproject/internal/domain/user"
)

var Module = fx.Module("handlers",
	fx.Provide(func(s *user.Service) *AuthHandler {
		return NewAuthHandler(s)
	}),
	fx.Provide(func(s *user.Service) *UserHandler {
		return NewUserHandler(s)
	}),
)
`
	if err := os.WriteFile(filepath.Join(handlerDir, "module.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateHandlerModule(tmpDir, "github.com/test/myproject", "Todo"); err != nil {
		t.Fatalf("updateHandlerModule() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(handlerDir, "module.go"))
	if err != nil {
		t.Fatalf("Failed to read module.go after update: %v", err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "NewTodoHandler") {
		t.Error("Expected NewTodoHandler in handlers module")
	}
	if !strings.Contains(contentStr, `"github.com/test/myproject/internal/domain/todo"`) {
		t.Error("Expected todo domain import in handlers module")
	}

	// Verify original content preserved
	if !strings.Contains(contentStr, "NewUserHandler") {
		t.Error("Expected NewUserHandler to be preserved")
	}
}

// TestUpdateHandlerModuleIdempotent tests idempotency
func TestUpdateHandlerModuleIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	handlerDir := filepath.Join(tmpDir, "internal", "adapters", "handlers")
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package handlers

import (
	"go.uber.org/fx"
	"github.com/test/myproject/internal/domain/user"
	"github.com/test/myproject/internal/domain/todo"
)

var Module = fx.Module("handlers",
	fx.Provide(func(s *user.Service) *UserHandler {
		return NewUserHandler(s)
	}),
	fx.Provide(func(s *todo.Service) *TodoHandler {
		return NewTodoHandler(s)
	}),
)
`
	if err := os.WriteFile(filepath.Join(handlerDir, "module.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateHandlerModule(tmpDir, "github.com/test/myproject", "Todo"); err != nil {
		t.Fatalf("updateHandlerModule() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(handlerDir, "module.go"))
	if err != nil {
		t.Fatalf("Failed to read module.go: %v", err)
	}
	count := strings.Count(string(content), "NewTodoHandler")
	if count != 1 {
		t.Errorf("Expected exactly 1 occurrence of NewTodoHandler, got %d", count)
	}
}

// TestUpdateRoutes tests route registration (Story 8.3 AC: 3)
func TestUpdateRoutes(t *testing.T) {
	tmpDir := t.TempDir()
	httpDir := filepath.Join(tmpDir, "internal", "adapters", "http")
	if err := os.MkdirAll(httpDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/test/myproject/internal/adapters/handlers"
)

func RegisterRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware fiber.Handler,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User routes (protected)
	users := v1.Group("/users", authMiddleware)
	users.Get("/me", userHandler.GetMe)
}
`
	if err := os.WriteFile(filepath.Join(httpDir, "routes.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	// Test protected routes (default)
	if err := updateRoutes(tmpDir, "Todo", false, nil); err != nil {
		t.Fatalf("updateRoutes() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(httpDir, "routes.go"))
	if err != nil {
		t.Fatalf("Failed to read routes.go: %v", err)
	}
	contentStr := string(content)

	// Verify handler parameter added
	if !strings.Contains(contentStr, "todoHandler *handlers.TodoHandler") {
		t.Error("Expected todoHandler parameter in RegisterRoutes")
	}

	// Verify route group with auth middleware
	if !strings.Contains(contentStr, `v1.Group("/todos", authMiddleware)`) {
		t.Error("Expected /todos route group with authMiddleware")
	}

	// Verify CRUD routes
	expectedRoutes := []string{
		`todos.Post("", todoHandler.CreateTodo)`,
		`todos.Get("", todoHandler.GetAllTodos)`,
		`todos.Get("/:id", todoHandler.GetTodo)`,
		`todos.Put("/:id", todoHandler.UpdateTodo)`,
		`todos.Delete("/:id", todoHandler.DeleteTodo)`,
	}
	for _, route := range expectedRoutes {
		if !strings.Contains(contentStr, route) {
			t.Errorf("Expected route %q", route)
		}
	}
}

// TestUpdateRoutesPublic tests public route generation
func TestUpdateRoutesPublic(t *testing.T) {
	tmpDir := t.TempDir()
	httpDir := filepath.Join(tmpDir, "internal", "adapters", "http")
	if err := os.MkdirAll(httpDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/test/myproject/internal/adapters/handlers"
)

func RegisterRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware fiber.Handler,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
}
`
	if err := os.WriteFile(filepath.Join(httpDir, "routes.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateRoutes(tmpDir, "Todo", true, nil); err != nil {
		t.Fatalf("updateRoutes(public) failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(httpDir, "routes.go"))
	if err != nil {
		t.Fatalf("Failed to read routes.go: %v", err)
	}
	contentStr := string(content)

	// Verify route group WITHOUT auth middleware
	if !strings.Contains(contentStr, `v1.Group("/todos")`) {
		t.Error("Expected /todos route group without authMiddleware for public routes")
	}
	if strings.Contains(contentStr, `"/todos", authMiddleware`) {
		t.Error("Public routes should NOT have authMiddleware")
	}
}

// TestUpdateMainFxModule tests main.go fx module update
func TestUpdateMainFxModule(t *testing.T) {
	tmpDir := t.TempDir()
	cmdDir := filepath.Join(tmpDir, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package main

import (
	"github.com/test/myproject/internal/domain/user"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Domain services
		user.Module,

		// Data persistence
		repository.Module,
	).Run()
}
`
	if err := os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	if err := updateMainFxModule(tmpDir, "github.com/test/myproject", "Todo"); err != nil {
		t.Fatalf("updateMainFxModule() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(cmdDir, "main.go"))
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "todo.Module") {
		t.Error("Expected todo.Module in main.go fx.New()")
	}
	if !strings.Contains(contentStr, `"github.com/test/myproject/internal/domain/todo"`) {
		t.Error("Expected todo domain import in main.go")
	}

	// Verify original content preserved
	if !strings.Contains(contentStr, "user.Module") {
		t.Error("Expected user.Module to be preserved")
	}
}

// TestGenerateModelFilesE2E tests the full generation flow (AC: 1-5)
func TestGenerateModelFilesE2E(t *testing.T) {
	tmpDir := t.TempDir()
	setupFakeProject(t, tmpDir)

	ctx := &projectContext{
		ModulePath: "github.com/test/myproject",
		Database:   "postgres",
		ProjectDir: tmpDir,
	}

	fields := []FieldDefinition{
		{Name: "Title", Type: "string", GORMTags: []string{"not_null"}, JSONName: "title"},
		{Name: "Completed", Type: "bool", GORMTags: nil, JSONName: "completed"},
	}

	if err := generateModelFiles(ctx, "Todo", fields, false, nil); err != nil {
		t.Fatalf("generateModelFiles() failed: %v", err)
	}

	// Verify model file
	modelContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "models", "todo.go"))
	if err != nil {
		t.Fatalf("Model file not created: %v", err)
	}
	if !strings.Contains(string(modelContent), "type Todo struct") {
		t.Error("Model file should contain Todo struct")
	}

	// Verify interface file
	ifaceContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "interfaces", "todo_repository.go"))
	if err != nil {
		t.Fatalf("Interface file not created: %v", err)
	}
	if !strings.Contains(string(ifaceContent), "type TodoRepository interface") {
		t.Error("Interface file should contain TodoRepository interface")
	}

	// Verify repository implementation
	repoContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "adapters", "repository", "todo_repository.go"))
	if err != nil {
		t.Fatalf("Repository file not created: %v", err)
	}
	if !strings.Contains(string(repoContent), "func NewTodoRepository") {
		t.Error("Repository file should contain NewTodoRepository")
	}

	// Verify AutoMigrate updated
	dbContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "infrastructure", "database", "database.go"))
	if err != nil {
		t.Fatalf("Failed to read database.go: %v", err)
	}
	if !strings.Contains(string(dbContent), "&models.Todo{}") {
		t.Error("database.go should contain &models.Todo{}")
	}

	// Verify repository module updated
	repoModContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "adapters", "repository", "module.go"))
	if err != nil {
		t.Fatalf("Failed to read repository module.go: %v", err)
	}
	if !strings.Contains(string(repoModContent), "NewTodoRepository") {
		t.Error("repository module.go should contain NewTodoRepository")
	}

	// Verify service file created
	serviceContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "domain", "todo", "service.go"))
	if err != nil {
		t.Fatalf("Service file not created: %v", err)
	}
	if !strings.Contains(string(serviceContent), "func NewService(") {
		t.Error("Service file should contain NewService")
	}

	// Verify service module file created
	svcModContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "domain", "todo", "module.go"))
	if err != nil {
		t.Fatalf("Service module file not created: %v", err)
	}
	if !strings.Contains(string(svcModContent), `fx.Module("todo"`) {
		t.Error("Service module should contain fx.Module(\"todo\")")
	}

	// Verify handler file created
	handlerContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "adapters", "handlers", "todo_handler.go"))
	if err != nil {
		t.Fatalf("Handler file not created: %v", err)
	}
	if !strings.Contains(string(handlerContent), "type TodoHandler struct") {
		t.Error("Handler file should contain TodoHandler struct")
	}

	// Verify handler module updated
	handlerModContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "adapters", "handlers", "module.go"))
	if err != nil {
		t.Fatalf("Failed to read handlers module.go: %v", err)
	}
	if !strings.Contains(string(handlerModContent), "NewTodoHandler") {
		t.Error("handlers module.go should contain NewTodoHandler")
	}

	// Verify routes updated
	routesContent, err := os.ReadFile(filepath.Join(tmpDir, "internal", "adapters", "http", "routes.go"))
	if err != nil {
		t.Fatalf("Failed to read routes.go: %v", err)
	}
	routesStr := string(routesContent)
	if !strings.Contains(routesStr, "todoHandler") {
		t.Error("routes.go should contain todoHandler")
	}
	if !strings.Contains(routesStr, "/todos") {
		t.Error("routes.go should contain /todos route")
	}

	// Verify main.go fx module updated
	mainContent, err := os.ReadFile(filepath.Join(tmpDir, "cmd", "main.go"))
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}
	if !strings.Contains(string(mainContent), "todo.Module") {
		t.Error("main.go should contain todo.Module")
	}
}

// TestAddModelFieldParsingSummaryWithGeneration tests that add-model generates files (AC: 1-5)
func TestAddModelFieldParsingSummaryWithGeneration(t *testing.T) {
	tmpDir := t.TempDir()
	setupFakeProject(t, tmpDir)

	cmd := exec.Command(absBinaryPath(t), "add-model", "Todo", "--fields", "title:string,completed:bool", "--yes")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Expected success for valid add-model with --yes, got error: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "generated successfully") {
		t.Errorf("Expected 'generated successfully' in output, got: %s", outputStr)
	}

	// Verify files were created
	expectedFiles := []string{
		"internal/models/todo.go",
		"internal/interfaces/todo_repository.go",
		"internal/adapters/repository/todo_repository.go",
		"internal/domain/todo/service.go",
		"internal/domain/todo/module.go",
		"internal/adapters/handlers/todo_handler.go",
	}
	for _, f := range expectedFiles {
		if _, err := os.Stat(filepath.Join(tmpDir, f)); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist after generation", f)
		}
	}
}

// TestE2EAddModelWithBuild tests that generated project compiles after add-model (AC: 4)
func TestE2EAddModelWithBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E build test in short mode")
	}

	testProjectName := "e2e-test-build-" + t.Name()
	defer os.RemoveAll(testProjectName)

	// Step 1: Generate a real project
	createCmd := exec.Command(binaryPath, testProjectName, "--database", "postgres")
	if output, err := createCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to create test project: %v\nOutput: %s", err, string(output))
	}

	// Step 2: Run add-model
	addModelCmd := exec.Command(absBinaryPath(t), "add-model", "Todo", "--fields", "title:string:not_null,completed:bool", "--yes")
	addModelCmd.Dir = testProjectName
	if output, err := addModelCmd.CombinedOutput(); err != nil {
		t.Fatalf("add-model failed: %v\nOutput: %s", err, string(output))
	}

	// Step 3: Run go mod tidy to generate go.sum
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = testProjectName
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\nOutput: %s", err, string(output))
	}

	// Step 4: Verify go build succeeds
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = testProjectName
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed after add-model: %v\nOutput: %s", err, string(output))
	}

	// Step 5: Verify generated tests compile and pass
	testCmd := exec.Command("go", "test", "./internal/domain/todo", "./internal/adapters/handlers", "-v")
	testCmd.Dir = testProjectName
	if output, err := testCmd.CombinedOutput(); err != nil {
		t.Fatalf("go test failed for generated test files: %v\nOutput: %s", err, string(output))
	}
}

// ==================== Story 8.5 Relation Tests ====================

// TestGenerateModelTemplateWithBelongsTo tests model generation with a BelongsTo relation.
func TestGenerateModelTemplateWithBelongsTo(t *testing.T) {
	fields := []FieldDefinition{
		{Name: "Content", Type: "string", GORMTags: []string{"not_null"}, JSONName: "content"},
	}
	relations := &RelationConfig{BelongsTo: "Post"}

	content := generateModelTemplate("Comment", fields, relations)

	// Verify foreign key field
	if !strings.Contains(content, "PostID") {
		t.Error("Expected PostID foreign key field")
	}
	if !strings.Contains(content, `gorm:"not null;index"`) {
		t.Error("Expected 'not null;index' GORM tag on foreign key")
	}
	if !strings.Contains(content, `json:"post_id"`) {
		t.Error("Expected json tag for post_id")
	}

	// Verify relation struct field
	if !strings.Contains(content, "Post") {
		t.Error("Expected Post relation field")
	}
	if !strings.Contains(content, `gorm:"foreignKey:PostID"`) {
		t.Error("Expected 'foreignKey:PostID' GORM tag on relation field")
	}
	if !strings.Contains(content, `json:"post,omitempty"`) {
		t.Error("Expected json tag for post with omitempty")
	}

	// Verify user fields are still present
	if !strings.Contains(content, "Content") {
		t.Error("Expected Content field")
	}

	// Verify base fields are still present
	if !strings.Contains(content, "CreatedAt") {
		t.Error("Expected CreatedAt field")
	}
	if !strings.Contains(content, "type Comment struct") {
		t.Error("Expected 'type Comment struct'")
	}
}

// TestGenerateRepositoryInterfaceTemplateWithBelongsTo tests repo interface with BelongsTo relation.
func TestGenerateRepositoryInterfaceTemplateWithBelongsTo(t *testing.T) {
	relations := &RelationConfig{BelongsTo: "Post"}

	content := generateRepositoryInterfaceTemplate("github.com/user/project", "Comment", relations)

	// Verify standard CRUD methods still present
	if !strings.Contains(content, "Create(ctx context.Context,") {
		t.Error("Expected Create method")
	}
	if !strings.Contains(content, "FindByID(ctx context.Context, id uint)") {
		t.Error("Expected FindByID method")
	}

	// Verify relation methods
	if !strings.Contains(content, "FindByIDWithRelations(ctx context.Context, id uint, includes []string)") {
		t.Error("Expected FindByIDWithRelations method")
	}
	if !strings.Contains(content, "FindByPostID(ctx context.Context, parentID uint, page, limit int)") {
		t.Error("Expected FindByPostID method")
	}
}

// TestGenerateRepositoryImplTemplateWithBelongsTo tests repo implementation with BelongsTo relation.
func TestGenerateRepositoryImplTemplateWithBelongsTo(t *testing.T) {
	relations := &RelationConfig{BelongsTo: "Post"}

	content := generateRepositoryImplTemplate("github.com/user/project", "Comment", relations)

	// Verify standard CRUD methods still present
	if !strings.Contains(content, "func (r *commentRepository) Create") {
		t.Error("Expected Create method")
	}
	if !strings.Contains(content, "func (r *commentRepository) FindByID") {
		t.Error("Expected FindByID method")
	}

	// Verify FindByIDWithRelations implementation
	if !strings.Contains(content, "func (r *commentRepository) FindByIDWithRelations") {
		t.Error("Expected FindByIDWithRelations method")
	}
	if !strings.Contains(content, "Preload") {
		t.Error("Expected Preload call in FindByIDWithRelations")
	}

	// Verify FindByPostID implementation
	if !strings.Contains(content, "func (r *commentRepository) FindByPostID") {
		t.Error("Expected FindByPostID method")
	}
	if !strings.Contains(content, "post_id = ?") {
		t.Error("Expected 'post_id = ?' WHERE clause in FindByPostID")
	}
}

// TestGenerateServiceTemplateWithBelongsTo tests service with BelongsTo relation.
func TestGenerateServiceTemplateWithBelongsTo(t *testing.T) {
	relations := &RelationConfig{BelongsTo: "Post"}

	content := generateServiceTemplate("github.com/user/project", "Comment", relations)

	// Verify standard CRUD methods
	if !strings.Contains(content, "func (s *Service) Create(") {
		t.Error("Expected Create method")
	}
	if !strings.Contains(content, "func (s *Service) GetByID(") {
		t.Error("Expected GetByID method")
	}

	// Verify relation methods
	if !strings.Contains(content, "func (s *Service) GetByIDWithRelations(") {
		t.Error("Expected GetByIDWithRelations method")
	}
	if !strings.Contains(content, "func (s *Service) GetByPostID(") {
		t.Error("Expected GetByPostID method")
	}

	// Verify pagination in GetByPostID
	if !strings.Contains(content, "s.repo.FindByPostID(") {
		t.Error("Expected FindByPostID delegation to repo")
	}
}

// TestGenerateHandlerTemplateWithBelongsTo tests handler template with BelongsTo relation.
func TestGenerateHandlerTemplateWithBelongsTo(t *testing.T) {
	fields := []FieldDefinition{
		{Name: "Content", Type: "string", GORMTags: []string{"not_null"}, JSONName: "content"},
	}
	relations := &RelationConfig{BelongsTo: "Post"}

	content := generateHandlerTemplate("github.com/user/project", "Comment", fields, relations)

	// Verify standard CRUD handlers still present
	if !strings.Contains(content, "func (h *CommentHandler) CreateComment(") {
		t.Error("Expected CreateComment handler")
	}
	if !strings.Contains(content, "func (h *CommentHandler) GetComment(") {
		t.Error("Expected GetComment handler")
	}

	// Verify GetByPost handler
	if !strings.Contains(content, "func (h *CommentHandler) GetByPost(") {
		t.Error("Expected GetByPost handler")
	}
	if !strings.Contains(content, "@Summary Get comments by post ID") {
		t.Error("Expected Swagger annotation for GetByPost")
	}
	if !strings.Contains(content, "@Router /api/v1/posts/{postId}/comments [get]") {
		t.Error("Expected nested route Swagger annotation for GetByPost")
	}

	// Verify CreateForPost handler
	if !strings.Contains(content, "func (h *CommentHandler) CreateForPost(") {
		t.Error("Expected CreateForPost handler")
	}
	if !strings.Contains(content, "@Router /api/v1/posts/{postId}/comments [post]") {
		t.Error("Expected nested route Swagger annotation for CreateForPost")
	}

	// Verify "strings" import was injected for relations
	if !strings.Contains(content, `"strings"`) {
		t.Error("Expected 'strings' import for relation handlers")
	}

	// Verify parent ID parsing
	if !strings.Contains(content, `c.Params("postId")`) {
		t.Error("Expected postId param parsing")
	}

	// Verify parent ID assignment in CreateForPost
	if !strings.Contains(content, "PostID: uint(parentID)") {
		t.Error("Expected PostID assignment in CreateForPost")
	}
}

// TestUpdateParentModelHasMany tests surgical insertion of has-many field into parent model.
func TestUpdateParentModelHasMany(t *testing.T) {
	tmpDir := t.TempDir()
	modelsDir := filepath.Join(tmpDir, "internal", "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		t.Fatal(err)
	}

	bt := "`"
	parentContent := `package models

import (
	"time"

	"gorm.io/gorm"
)

// Post represents a post entity in the system.
type Post struct {
	ID        uint           ` + bt + `gorm:"primaryKey" json:"id"` + bt + `
	Title     string         ` + bt + `gorm:"not null" json:"title"` + bt + `
	CreatedAt time.Time      ` + bt + `gorm:"autoCreateTime" json:"created_at"` + bt + `
	UpdatedAt time.Time      ` + bt + `gorm:"autoUpdateTime" json:"updated_at"` + bt + `
	DeletedAt gorm.DeletedAt ` + bt + `gorm:"index" json:"deleted_at,omitempty"` + bt + `
}
`
	if err := os.WriteFile(filepath.Join(modelsDir, "post.go"), []byte(parentContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Call updateParentModelHasMany
	if err := updateParentModelHasMany(tmpDir, "Comment", "Post"); err != nil {
		t.Fatalf("updateParentModelHasMany() failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(modelsDir, "post.go"))
	if err != nil {
		t.Fatalf("Failed to read post.go: %v", err)
	}
	contentStr := string(content)

	// Verify Comments slice field was added
	if !strings.Contains(contentStr, "Comments") {
		t.Error("Expected Comments field to be inserted")
	}
	if !strings.Contains(contentStr, "[]Comment") {
		t.Error("Expected []Comment type")
	}
	if !strings.Contains(contentStr, `gorm:"foreignKey:PostID"`) {
		t.Error("Expected gorm foreignKey:PostID tag")
	}
	if !strings.Contains(contentStr, `json:"comments,omitempty"`) {
		t.Error("Expected json comments,omitempty tag")
	}

	// Verify it was inserted before CreatedAt
	commentsIdx := strings.Index(contentStr, "Comments")
	createdAtIdx := strings.Index(contentStr, "CreatedAt")
	if commentsIdx >= createdAtIdx {
		t.Error("Expected Comments field to appear before CreatedAt")
	}

	// Verify existing fields are preserved
	if !strings.Contains(contentStr, "Title") {
		t.Error("Expected Title field to be preserved")
	}
	if !strings.Contains(contentStr, "ID") {
		t.Error("Expected ID field to be preserved")
	}
}

// TestUpdateParentModelHasManyIdempotent tests that calling twice doesn't duplicate the field.
func TestUpdateParentModelHasManyIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	modelsDir := filepath.Join(tmpDir, "internal", "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		t.Fatal(err)
	}

	bt := "`"
	parentContent := `package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           ` + bt + `gorm:"primaryKey" json:"id"` + bt + `
	Title     string         ` + bt + `gorm:"not null" json:"title"` + bt + `
	CreatedAt time.Time      ` + bt + `gorm:"autoCreateTime" json:"created_at"` + bt + `
	UpdatedAt time.Time      ` + bt + `gorm:"autoUpdateTime" json:"updated_at"` + bt + `
	DeletedAt gorm.DeletedAt ` + bt + `gorm:"index" json:"deleted_at,omitempty"` + bt + `
}
`
	if err := os.WriteFile(filepath.Join(modelsDir, "post.go"), []byte(parentContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Call twice
	if err := updateParentModelHasMany(tmpDir, "Comment", "Post"); err != nil {
		t.Fatalf("First call failed: %v", err)
	}
	if err := updateParentModelHasMany(tmpDir, "Comment", "Post"); err != nil {
		t.Fatalf("Second call failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(modelsDir, "post.go"))
	if err != nil {
		t.Fatalf("Failed to read post.go: %v", err)
	}
	count := strings.Count(string(content), "Comments")
	if count != 1 {
		t.Errorf("Expected exactly 1 occurrence of Comments, got %d", count)
	}
}

// TestUpdateRoutesWithBelongsTo tests route registration with belongs-to relation.
func TestUpdateRoutesWithBelongsTo(t *testing.T) {
	tmpDir := t.TempDir()
	httpDir := filepath.Join(tmpDir, "internal", "adapters", "http")
	if err := os.MkdirAll(httpDir, 0755); err != nil {
		t.Fatal(err)
	}

	original := `package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/test/myproject/internal/adapters/handlers"
)

func RegisterRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware fiber.Handler,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User routes (protected)
	users := v1.Group("/users", authMiddleware)
	users.Get("/me", userHandler.GetMe)
}
`
	if err := os.WriteFile(filepath.Join(httpDir, "routes.go"), []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	relations := &RelationConfig{BelongsTo: "Post"}
	if err := updateRoutes(tmpDir, "Comment", false, relations); err != nil {
		t.Fatalf("updateRoutes() with BelongsTo failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(httpDir, "routes.go"))
	if err != nil {
		t.Fatalf("Failed to read routes.go: %v", err)
	}
	contentStr := string(content)

	// Verify standard CRUD routes
	if !strings.Contains(contentStr, "commentHandler *handlers.CommentHandler") {
		t.Error("Expected commentHandler parameter in RegisterRoutes")
	}
	if !strings.Contains(contentStr, `v1.Group("/comments", authMiddleware)`) {
		t.Error("Expected /comments route group")
	}

	// Verify nested route group for BelongsTo
	if !strings.Contains(contentStr, `/posts/:postId/comments`) {
		t.Error("Expected nested route group /posts/:postId/comments")
	}
	if !strings.Contains(contentStr, "commentHandler.GetByPost") {
		t.Error("Expected GetByPost handler in nested routes")
	}
	if !strings.Contains(contentStr, "commentHandler.CreateForPost") {
		t.Error("Expected CreateForPost handler in nested routes")
	}
}

// TestGenerateServiceTestTemplateWithBelongsTo tests that generated service tests include relation tests.
func TestGenerateServiceTestTemplateWithBelongsTo(t *testing.T) {
	fields := []FieldDefinition{
		{Name: "Content", Type: "string", GORMTags: []string{"not_null"}, JSONName: "content"},
	}
	relations := &RelationConfig{BelongsTo: "Post"}

	content := generateServiceTestTemplate("github.com/user/project", "Comment", fields, relations)

	// Verify standard CRUD test functions are present
	if !strings.Contains(content, "func TestService_Create(t *testing.T)") {
		t.Error("Expected TestService_Create test function")
	}
	if !strings.Contains(content, "func TestService_GetByID(t *testing.T)") {
		t.Error("Expected TestService_GetByID test function")
	}

	// Verify extra mock methods for relations
	if !strings.Contains(content, "func (m *MockCommentRepository) FindByIDWithRelations(") {
		t.Error("Expected FindByIDWithRelations mock method")
	}
	if !strings.Contains(content, "func (m *MockCommentRepository) FindByPostID(") {
		t.Error("Expected FindByPostID mock method")
	}

	// Verify relation test functions
	if !strings.Contains(content, "func TestService_GetByIDWithRelations(t *testing.T)") {
		t.Error("Expected TestService_GetByIDWithRelations test function")
	}
	if !strings.Contains(content, "func TestService_GetByPostID(t *testing.T)") {
		t.Error("Expected TestService_GetByPostID test function")
	}

	// Verify test subtests
	if !strings.Contains(content, `t.Run("with_includes"`) {
		t.Error("Expected with_includes subtest")
	}
	if !strings.Contains(content, `t.Run("pagination_defaults"`) {
		t.Error("Expected pagination_defaults subtest")
	}
}

// TestGenerateHandlerTestTemplateWithBelongsTo tests that generated handler tests include relation tests.
func TestGenerateHandlerTestTemplateWithBelongsTo(t *testing.T) {
	fields := []FieldDefinition{
		{Name: "Content", Type: "string", GORMTags: []string{"not_null"}, JSONName: "content"},
	}
	relations := &RelationConfig{BelongsTo: "Post"}

	content := generateHandlerTestTemplate("github.com/user/project", "Comment", fields, relations)

	// Verify standard test functions
	if !strings.Contains(content, "func TestCreateComment_Success(t *testing.T)") {
		t.Error("Expected TestCreateComment_Success test function")
	}
	if !strings.Contains(content, "func TestGetComment_Found(t *testing.T)") {
		t.Error("Expected TestGetComment_Found test function")
	}

	// Verify extra mock methods
	if !strings.Contains(content, "func (m *MockCommentRepository) FindByIDWithRelations(") {
		t.Error("Expected FindByIDWithRelations mock method in handler tests")
	}
	if !strings.Contains(content, "func (m *MockCommentRepository) FindByPostID(") {
		t.Error("Expected FindByPostID mock method in handler tests")
	}

	// Verify nested route setup in setupCommentTestApp
	if !strings.Contains(content, "func setupCommentTestApp(") {
		t.Error("Expected setupCommentTestApp function")
	}
	if !strings.Contains(content, `/posts/:postId/comments`) {
		t.Error("Expected nested route in test app setup")
	}
	if !strings.Contains(content, "handler.GetByPost") {
		t.Error("Expected GetByPost handler in test app setup")
	}
	if !strings.Contains(content, "handler.CreateForPost") {
		t.Error("Expected CreateForPost handler in test app setup")
	}

	// Verify relation test functions
	if !strings.Contains(content, "func TestGetByPost_Success(t *testing.T)") {
		t.Error("Expected TestGetByPost_Success test function")
	}
	if !strings.Contains(content, "func TestGetByPost_InvalidParentID(t *testing.T)") {
		t.Error("Expected TestGetByPost_InvalidParentID test function")
	}
	if !strings.Contains(content, "func TestCreateForPost_Success(t *testing.T)") {
		t.Error("Expected TestCreateForPost_Success test function")
	}
	if !strings.Contains(content, "func TestCreateForPost_InvalidParentID(t *testing.T)") {
		t.Error("Expected TestCreateForPost_InvalidParentID test function")
	}

	// Verify nested route URLs in tests
	if !strings.Contains(content, "/api/v1/posts/1/comments") {
		t.Error("Expected nested route URL in test requests")
	}
	if !strings.Contains(content, "/api/v1/posts/abc/comments") {
		t.Error("Expected invalid parent ID URL in test requests")
	}
}
