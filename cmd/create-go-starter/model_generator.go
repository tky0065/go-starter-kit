package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// pluralize returns a simple lowercase English plural form of a word.
// Handles irregular plurals, common cases: y→ies, s/x/z/ch/sh→es, default→s.
func pluralize(s string) string {
	if s == "" {
		return ""
	}
	lower := strings.ToLower(s)

	// Handle irregular plurals (most common cases)
	irregulars := map[string]string{
		"person":     "people",
		"child":      "children",
		"man":        "men",
		"woman":      "women",
		"foot":       "feet",
		"tooth":      "teeth",
		"mouse":      "mice",
		"goose":      "geese",
		"ox":         "oxen",
		"datum":      "data",
		"medium":     "media",
		"criterion":  "criteria",
		"phenomenon": "phenomena",
	}

	if plural, ok := irregulars[lower]; ok {
		return plural
	}

	// Handle consonant + y → ies
	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		prev := lower[len(lower)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return lower[:len(lower)-1] + "ies"
		}
	}

	// Handle s, x, z, ch, sh → es
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") ||
		strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return lower + "es"
	}

	// Handle f/fe → ves (knife→knives, wife→wives)
	if strings.HasSuffix(lower, "f") {
		return lower[:len(lower)-1] + "ves"
	}
	if strings.HasSuffix(lower, "fe") {
		return lower[:len(lower)-2] + "ves"
	}

	// Default: add s
	return lower + "s"
}

// pluralizePascal returns PascalCase plural form (e.g., "Todo" → "Todos").
func pluralizePascal(s string) string {
	p := pluralize(s)
	if p == "" {
		return ""
	}
	return strings.ToUpper(p[:1]) + p[1:]
}

// generateModelFiles creates all the files for a new model and modifies existing files.
// It generates: model, interface, repository, service, service module, handler,
// service tests, handler tests, and modifies: database.go (AutoMigrate),
// repository/module.go, handlers/module.go, routes.go.
// When relations are specified, it also generates relation fields and nested endpoints.
func generateModelFiles(ctx *projectContext, modelName string, fields []FieldDefinition, isPublic bool, relations *RelationConfig) error {
	lowerName := strings.ToLower(modelName)

	// Files to create
	newFiles := []FileGenerator{
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "models", lowerName+".go"),
			Content: generateModelTemplate(modelName, fields, relations),
		},
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "interfaces", lowerName+"_repository.go"),
			Content: generateRepositoryInterfaceTemplate(ctx.ModulePath, modelName, relations),
		},
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "adapters", "repository", lowerName+"_repository.go"),
			Content: generateRepositoryImplTemplate(ctx.ModulePath, modelName, relations),
		},
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "domain", lowerName, "service.go"),
			Content: generateServiceTemplate(ctx.ModulePath, modelName, relations),
		},
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "domain", lowerName, "module.go"),
			Content: generateServiceModuleTemplate(modelName),
		},
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "adapters", "handlers", lowerName+"_handler.go"),
			Content: generateHandlerTemplate(ctx.ModulePath, modelName, fields, relations),
		},
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "domain", lowerName, "service_test.go"),
			Content: generateServiceTestTemplate(ctx.ModulePath, modelName, fields, relations),
		},
		{
			Path:    filepath.Join(ctx.ProjectDir, "internal", "adapters", "handlers", lowerName+"_handler_test.go"),
			Content: generateHandlerTestTemplate(ctx.ModulePath, modelName, fields, relations),
		},
	}

	// Create new files
	for _, file := range newFiles {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}
		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", file.Path, err)
		}
	}

	// Modify existing files
	if err := updateAutoMigrate(ctx.ProjectDir, modelName); err != nil {
		return fmt.Errorf("failed to update AutoMigrate: %w", err)
	}

	if err := updateRepositoryModule(ctx.ProjectDir, ctx.ModulePath, modelName); err != nil {
		return fmt.Errorf("failed to update repository module: %w", err)
	}

	if err := updateHandlerModule(ctx.ProjectDir, ctx.ModulePath, modelName); err != nil {
		return fmt.Errorf("failed to update handler module: %w", err)
	}

	if err := updateRoutes(ctx.ProjectDir, modelName, isPublic, relations); err != nil {
		return fmt.Errorf("failed to update routes: %w", err)
	}

	if err := updateMainFxModule(ctx.ProjectDir, ctx.ModulePath, modelName); err != nil {
		return fmt.Errorf("failed to update main.go fx module: %w", err)
	}

	// If has-many is specified, surgically update the parent model to add the slice field
	if relations != nil && relations.HasMany != "" {
		if err := updateParentModelHasMany(ctx.ProjectDir, modelName, relations.HasMany); err != nil {
			return fmt.Errorf("failed to update parent model with has-many: %w", err)
		}
	}

	return nil
}

// generateModelTemplate generates the content for internal/models/<name>.go
// When relations.BelongsTo is set, adds foreign key and relation fields.
func generateModelTemplate(modelName string, fields []FieldDefinition, relations *RelationConfig) string {
	bt := "`"
	lowerName := strings.ToLower(modelName)

	// Check if we need time import (for user time.Time fields; base fields always need it)
	var fieldLines []string

	// Base GORM field: ID
	fieldLines = append(fieldLines, fmt.Sprintf("\tID        uint           %sgorm:\"primaryKey\" json:\"id\"%s", bt, bt))

	// User-defined fields
	for _, f := range fields {
		gormTag := buildGORMTag(f)
		jsonTag := fmt.Sprintf("json:\"%s\"", f.JSONName)
		fieldLines = append(fieldLines, fmt.Sprintf("\t%-9s %-14s %s%s %s%s", f.Name, f.Type, bt, gormTag, jsonTag, bt))
	}

	// BelongsTo relation fields: foreign key + relation struct
	if relations != nil && relations.BelongsTo != "" {
		parent := relations.BelongsTo
		parentLower := strings.ToLower(parent)
		parentSnake := toSnakeCase(parent)
		// Foreign key field: <Parent>ID uint `gorm:"not null;index" json:"<parent>_id"`
		fieldLines = append(fieldLines, fmt.Sprintf("\t%-9s %-14s %sgorm:\"not null;index\" json:\"%s_id\"%s",
			parent+"ID", "uint", bt, parentSnake, bt))
		// Relation field: <Parent> <Parent> `gorm:"foreignKey:<Parent>ID" json:"<parent>,omitempty"`
		fieldLines = append(fieldLines, fmt.Sprintf("\t%-9s %-14s %sgorm:\"foreignKey:%sID\" json:\"%s,omitempty\"%s",
			parent, parent, bt, parent, parentLower, bt))
	}

	// Timestamp fields
	fieldLines = append(fieldLines, fmt.Sprintf("\tCreatedAt time.Time      %sgorm:\"autoCreateTime\" json:\"created_at\"%s", bt, bt))
	fieldLines = append(fieldLines, fmt.Sprintf("\tUpdatedAt time.Time      %sgorm:\"autoUpdateTime\" json:\"updated_at\"%s", bt, bt))
	fieldLines = append(fieldLines, fmt.Sprintf("\tDeletedAt gorm.DeletedAt %sgorm:\"index\" json:\"deleted_at,omitempty\"%s", bt, bt))

	return fmt.Sprintf(`package models

import (
	"time"

	"gorm.io/gorm"
)

// %s represents a %s entity in the system.
type %s struct {
%s
}
`, modelName, lowerName, modelName, strings.Join(fieldLines, "\n"))
}

// buildGORMTag creates the GORM tag string for a field.
func buildGORMTag(f FieldDefinition) string {
	var parts []string

	for _, mod := range f.GORMTags {
		switch mod {
		case "unique":
			parts = append(parts, "uniqueIndex")
		case "not_null":
			parts = append(parts, "not null")
		case "index":
			parts = append(parts, "index")
		}
	}

	// Add type-specific defaults
	switch f.Type {
	case "bool":
		parts = append(parts, "default:false")
	}

	if len(parts) == 0 {
		return "gorm:\"\""
	}
	return fmt.Sprintf("gorm:\"%s\"", strings.Join(parts, ";"))
}

// generateRepositoryInterfaceTemplate generates the content for internal/interfaces/<name>_repository.go
// When relations are specified, adds FindByIDWithRelations and FindByParentID methods.
func generateRepositoryInterfaceTemplate(modulePath, modelName string, relations *RelationConfig) string {
	lowerName := strings.ToLower(modelName)

	// Build extra methods for relations
	var extraMethods string
	if relations != nil && relations.BelongsTo != "" {
		extraMethods += fmt.Sprintf("\tFindByIDWithRelations(ctx context.Context, id uint, includes []string) (*models.%s, error)\n", modelName)
		extraMethods += fmt.Sprintf("\tFindBy%sID(ctx context.Context, parentID uint, page, limit int) ([]*models.%s, int64, error)\n",
			relations.BelongsTo, modelName)
	}

	return fmt.Sprintf(`package interfaces

import (
	"context"

	"%s/internal/models"
)

// %sRepository defines the contract for %s persistence operations.
type %sRepository interface {
	Create(ctx context.Context, %s *models.%s) error
	FindByID(ctx context.Context, id uint) (*models.%s, error)
	FindAll(ctx context.Context, page, limit int) ([]*models.%s, int64, error)
	Update(ctx context.Context, %s *models.%s) error
	Delete(ctx context.Context, id uint) error
%s}
`, modulePath, modelName, lowerName, modelName, lowerName, modelName, modelName, modelName, lowerName, modelName, extraMethods)
}

// generateRepositoryImplTemplate generates the content for internal/adapters/repository/<name>_repository.go
// When relations.BelongsTo is set, adds FindByIDWithRelations and FindByParentID implementations.
func generateRepositoryImplTemplate(modulePath, modelName string, relations *RelationConfig) string {
	lowerName := strings.ToLower(modelName)
	structName := lowerName + "Repository"

	// Build extra methods for relations
	var extraMethods string
	if relations != nil && relations.BelongsTo != "" {
		extraMethods = fmt.Sprintf(`
func (r *%s) FindByIDWithRelations(ctx context.Context, id uint, includes []string) (*models.%s, error) {
	var %s models.%s
	query := r.db.WithContext(ctx)

	for _, include := range includes {
		query = query.Preload(include)
	}

	err := query.First(&%s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &%s, err
}

func (r *%s) FindBy%sID(ctx context.Context, parentID uint, page, limit int) ([]*models.%s, int64, error) {
	var items []*models.%s
	var total int64

	r.db.WithContext(ctx).Model(&models.%s{}).Where("%s_id = ?", parentID).Count(&total)

	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Where("%s_id = ?", parentID).Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}
`,
			structName, modelName,
			lowerName, modelName,
			lowerName,
			lowerName,
			structName, relations.BelongsTo, modelName,
			modelName,
			modelName, toSnakeCase(relations.BelongsTo),
			toSnakeCase(relations.BelongsTo),
		)
	}

	return fmt.Sprintf(`package repository

import (
	"context"
	"errors"

	"%s/internal/interfaces"
	"%s/internal/models"
	"gorm.io/gorm"
)

type %s struct {
	db *gorm.DB
}

// New%sRepository creates a new GORM-backed %sRepository.
func New%sRepository(db *gorm.DB) interfaces.%sRepository {
	return &%s{db: db}
}

func (r *%s) Create(ctx context.Context, %s *models.%s) error {
	return r.db.WithContext(ctx).Create(%s).Error
}

func (r *%s) FindByID(ctx context.Context, id uint) (*models.%s, error) {
	var %s models.%s
	err := r.db.WithContext(ctx).First(&%s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &%s, err
}

func (r *%s) FindAll(ctx context.Context, page, limit int) ([]*models.%s, int64, error) {
	var items []*models.%s
	var total int64

	r.db.WithContext(ctx).Model(&models.%s{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *%s) Update(ctx context.Context, %s *models.%s) error {
	return r.db.WithContext(ctx).Save(%s).Error
}

func (r *%s) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.%s{}, id).Error
}
%s`,
		modulePath, modulePath,
		structName,
		modelName, modelName, modelName, modelName, structName,
		structName, lowerName, modelName, lowerName,
		structName, modelName, lowerName, modelName, lowerName, lowerName,
		structName, modelName, modelName, modelName,
		structName, lowerName, modelName, lowerName,
		structName, modelName,
		extraMethods,
	)
}

// updateAutoMigrate adds the new model to the AutoMigrate call in database.go.
// This function performs surgical string manipulation to preserve existing code:
// 1. Finds the db.AutoMigrate(...) call using string search
// 2. Uses parenthesis depth counting to find the closing ) of the call
// 3. Inserts the new model reference before the closing )
// 4. Checks for idempotency - skips if model already present
//
// Example transformation:
//
//	Before: db.AutoMigrate(&models.User{}, &models.RefreshToken{})
//	After:  db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Todo{})
func updateAutoMigrate(projectDir, modelName string) error {
	dbFile := filepath.Join(projectDir, "internal", "infrastructure", "database", "database.go")

	content, err := os.ReadFile(dbFile)
	if err != nil {
		return fmt.Errorf("failed to read database.go: %w", err)
	}

	contentStr := string(content)

	// Find the AutoMigrate call and add the new model
	newModel := fmt.Sprintf("&models.%s{}", modelName)

	// Check if model is already in AutoMigrate (idempotency)
	if strings.Contains(contentStr, newModel) {
		return nil // Already present
	}

	// Find the AutoMigrate call - look for the closing ); pattern
	// Pattern: db.AutoMigrate(...existing models...)
	autoMigrateIdx := strings.Index(contentStr, "db.AutoMigrate(")
	if autoMigrateIdx == -1 {
		return fmt.Errorf("AutoMigrate call not found in database.go - ensure project structure matches go-starter-kit template")
	}

	// Find the closing parenthesis of AutoMigrate using depth counting
	// This handles nested parentheses correctly (e.g., function calls within AutoMigrate)
	depth := 0
	closeIdx := -1
	for i := autoMigrateIdx; i < len(contentStr); i++ {
		if contentStr[i] == '(' {
			depth++
		} else if contentStr[i] == ')' {
			depth--
			if depth == 0 {
				closeIdx = i
				break
			}
		}
	}

	if closeIdx == -1 {
		return fmt.Errorf("could not find closing parenthesis for AutoMigrate")
	}

	// Insert the new model before the closing parenthesis
	newContent := contentStr[:closeIdx] + ", " + newModel + contentStr[closeIdx:]

	return os.WriteFile(dbFile, []byte(newContent), 0644)
}

// updateRepositoryModule adds the new repository provider to repository/module.go.
// This function performs surgical fx module manipulation:
// 1. Validates that the interfaces import is present
// 2. Finds the last fx.Provide block within the fx.Module("repository")
// 3. Inserts a new fx.Provide after the last one with proper formatting
// 4. Checks for idempotency - skips if constructor already present
//
// Example transformation:
//
//	Before:
//	  var Module = fx.Module("repository",
//	      fx.Provide(func(db *gorm.DB) interfaces.UserRepository {
//	          return NewUserRepository(db)
//	      }),
//	  )
//	After:
//	  var Module = fx.Module("repository",
//	      fx.Provide(func(db *gorm.DB) interfaces.UserRepository {
//	          return NewUserRepository(db)
//	      }),
//	      fx.Provide(func(db *gorm.DB) interfaces.TodoRepository {
//	          return NewTodoRepository(db)
//	      }),
//	  )
func updateRepositoryModule(projectDir, modulePath, modelName string) error {
	moduleFile := filepath.Join(projectDir, "internal", "adapters", "repository", "module.go")

	content, err := os.ReadFile(moduleFile)
	if err != nil {
		return fmt.Errorf("failed to read repository module.go: %w", err)
	}

	contentStr := string(content)

	// Check if already added
	constructorName := fmt.Sprintf("New%sRepository", modelName)
	if strings.Contains(contentStr, constructorName) {
		return nil // Already present
	}

	// Add the interfaces import if not present
	interfacesImport := fmt.Sprintf("\"%s/internal/interfaces\"", modulePath)
	if !strings.Contains(contentStr, interfacesImport) {
		// This shouldn't happen in a properly generated project, but handle it
		return fmt.Errorf("interfaces import not found in repository module.go")
	}

	// Find the closing of the fx.Module and add the new provide before it
	// Pattern: fx.Module("repository",\n\tfx.Provide(...),\n)
	// We need to add another fx.Provide before the closing )
	moduleIdx := strings.Index(contentStr, "fx.Module(\"repository\"")
	if moduleIdx == -1 {
		return fmt.Errorf("fx.Module(\"repository\") not found in module.go")
	}

	// Find the last closing parenthesis of fx.Module
	// The module ends with a closing ) on its own line
	// Look for the pattern: "\n)\n" after the fx.Module start
	searchFrom := moduleIdx
	lastProvideEnd := -1

	// Find all fx.Provide blocks within the module
	for i := searchFrom; i < len(contentStr); i++ {
		remaining := contentStr[i:]
		if strings.HasPrefix(remaining, "fx.Provide(") {
			// Find the end of this Provide block
			depth := 0
			for j := i; j < len(contentStr); j++ {
				if contentStr[j] == '(' {
					depth++
				} else if contentStr[j] == ')' {
					depth--
					if depth == 0 {
						lastProvideEnd = j + 1
						i = j // Skip past this Provide
						break
					}
				}
			}
		}
	}

	if lastProvideEnd == -1 {
		return fmt.Errorf("no fx.Provide found in repository module.go")
	}

	// Insert new provide after the last one
	newProvide := fmt.Sprintf(`,
	fx.Provide(func(db *gorm.DB) interfaces.%sRepository {
		return New%sRepository(db)
	})`, modelName, modelName)

	newContent := contentStr[:lastProvideEnd] + newProvide + contentStr[lastProvideEnd:]

	return os.WriteFile(moduleFile, []byte(newContent), 0644)
}

// generateServiceTemplate generates the content for internal/domain/<name>/service.go.
// The service depends on the repository interface (not the concrete implementation).
// When relations.BelongsTo is set, adds GetByIDWithRelations and GetByParentID methods.
func generateServiceTemplate(modulePath, modelName string, relations *RelationConfig) string {
	lowerName := strings.ToLower(modelName)

	// Build extra methods for relations
	var extraMethods string
	if relations != nil && relations.BelongsTo != "" {
		extraMethods = fmt.Sprintf(`
// GetByIDWithRelations retrieves a %s by ID with optional relation preloading.
func (s *Service) GetByIDWithRelations(ctx context.Context, id uint, includes []string) (*models.%s, error) {
	return s.repo.FindByIDWithRelations(ctx, id, includes)
}

// GetBy%sID retrieves all %ss for a given parent %s.
func (s *Service) GetBy%sID(ctx context.Context, parentID uint, page, limit int) ([]*models.%s, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.FindBy%sID(ctx, parentID, page, limit)
}
`,
			lowerName, modelName,
			relations.BelongsTo, lowerName, strings.ToLower(relations.BelongsTo),
			relations.BelongsTo, modelName,
			relations.BelongsTo,
		)
	}

	return fmt.Sprintf(`package %s

import (
	"context"

	"%s/internal/interfaces"
	"%s/internal/models"
)

// Service handles %s business logic.
type Service struct {
	repo interfaces.%sRepository
}

// NewService creates a new %s service.
func NewService(repo interfaces.%sRepository) *Service {
	return &Service{repo: repo}
}

// Create creates a new %s.
func (s *Service) Create(ctx context.Context, %s *models.%s) error {
	return s.repo.Create(ctx, %s)
}

// GetByID retrieves a %s by ID.
func (s *Service) GetByID(ctx context.Context, id uint) (*models.%s, error) {
	return s.repo.FindByID(ctx, id)
}

// GetAll retrieves all %ss with pagination.
func (s *Service) GetAll(ctx context.Context, page, limit int) ([]*models.%s, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.FindAll(ctx, page, limit)
}

// Update updates an existing %s.
func (s *Service) Update(ctx context.Context, %s *models.%s) error {
	return s.repo.Update(ctx, %s)
}

// Delete deletes a %s by ID.
func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
%s`,
		lowerName,
		modulePath, modulePath,
		lowerName, modelName,
		lowerName, modelName,
		lowerName, lowerName, modelName, lowerName,
		lowerName, modelName,
		lowerName, modelName,
		lowerName, lowerName, modelName, lowerName,
		lowerName,
		extraMethods,
	)
}

// generateServiceModuleTemplate generates the content for internal/domain/<name>/module.go.
func generateServiceModuleTemplate(modelName string) string {
	lowerName := strings.ToLower(modelName)
	return fmt.Sprintf(`package %s

import (
	"go.uber.org/fx"
)

// Module provides %s domain services via fx dependency injection.
var Module = fx.Module("%s",
	fx.Provide(NewService),
)
`, lowerName, lowerName, lowerName)
}

// getSemanticValidation returns appropriate validation tags based on field name and type.
// Detects common patterns like email, url, password, phone, etc.
func getSemanticValidation(fieldName, fieldType string, isRequired bool) string {
	lowerName := strings.ToLower(fieldName)

	// For non-string types, use basic validation
	if fieldType != "string" {
		// Bool fields should not use "required" since false is a valid value
		if fieldType == "bool" {
			return "omitempty"
		}
		if isRequired {
			return "required"
		}
		return "omitempty"
	}

	// String type - detect semantic meaning
	var validations []string
	if isRequired {
		validations = append(validations, "required")
	} else {
		validations = append(validations, "omitempty")
	}

	// Email fields
	if strings.Contains(lowerName, "email") {
		validations = append(validations, "email")
		return strings.Join(validations, ",")
	}

	// URL fields
	if strings.Contains(lowerName, "url") || strings.Contains(lowerName, "website") || strings.Contains(lowerName, "link") {
		validations = append(validations, "url")
		return strings.Join(validations, ",")
	}

	// Password fields - enforce minimum security
	if strings.Contains(lowerName, "password") || strings.Contains(lowerName, "passwd") {
		validations = append(validations, "min=8", "max=72") // bcrypt max is 72
		return strings.Join(validations, ",")
	}

	// Phone fields
	if strings.Contains(lowerName, "phone") || strings.Contains(lowerName, "mobile") || strings.Contains(lowerName, "telephone") {
		validations = append(validations, "min=10", "max=20") // International phone formats
		return strings.Join(validations, ",")
	}

	// UUID fields
	if strings.Contains(lowerName, "uuid") {
		validations = append(validations, "uuid")
		return strings.Join(validations, ",")
	}

	// Default string validation
	validations = append(validations, "min=1", "max=255")
	return strings.Join(validations, ",")
}

// generateHandlerTemplate generates the content for internal/adapters/handlers/<name>_handler.go.
// It includes Create/Update DTOs with validation tags and all 5 CRUD endpoints with Swagger annotations.
// When relations.BelongsTo is set, adds GetByParent and CreateForParent handlers with ?include= support.
func generateHandlerTemplate(modulePath, modelName string, fields []FieldDefinition, relations *RelationConfig) string {
	lowerName := strings.ToLower(modelName)
	pluralName := pluralizePascal(modelName)
	pluralLower := strings.ToLower(pluralName)
	bt := "`"

	// Build Create DTO fields (all required fields)
	var createFields []string
	for _, f := range fields {
		validateTag := getSemanticValidation(f.Name, f.Type, true)
		createFields = append(createFields,
			fmt.Sprintf("\t%-9s %-14s %sjson:\"%s\" validate:\"%s\"%s", f.Name, f.Type, bt, f.JSONName, validateTag, bt))
	}

	// Build Update DTO fields (all optional with pointers for strings)
	var updateFields []string
	for _, f := range fields {
		fieldType := f.Type
		validateTag := getSemanticValidation(f.Name, f.Type, false)
		if f.Type == "string" {
			fieldType = "*string"
		} else if f.Type == "bool" {
			fieldType = "*bool"
		} else if f.Type == "int" || f.Type == "uint" || f.Type == "float64" {
			fieldType = "*" + f.Type
		} else if f.Type == "time.Time" {
			fieldType = "*time.Time"
		}
		updateFields = append(updateFields,
			fmt.Sprintf("\t%-9s %-14s %sjson:\"%s\" validate:\"%s\"%s", f.Name, fieldType, bt, f.JSONName, validateTag, bt))
	}

	// Build model field assignments for Create
	var createAssignments []string
	for _, f := range fields {
		createAssignments = append(createAssignments, fmt.Sprintf("\t\t%s: req.%s,", f.Name, f.Name))
	}

	// Build model field updates for Update (only non-nil pointer fields)
	var updateAssignments []string
	for _, f := range fields {
		if f.Type == "string" || f.Type == "bool" || f.Type == "int" || f.Type == "uint" || f.Type == "float64" || f.Type == "time.Time" {
			updateAssignments = append(updateAssignments, fmt.Sprintf(`	if req.%s != nil {
		item.%s = *req.%s
	}`, f.Name, f.Name, f.Name))
		}
	}

	result := fmt.Sprintf(`package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"%s/internal/domain/%s"
	"%s/internal/models"
)

// Create%sRequest represents the request body for creating a %s.
type Create%sRequest struct {
%s
}

// Update%sRequest represents the request body for updating a %s.
type Update%sRequest struct {
%s
}

// %sHandler handles HTTP requests for %s operations.
type %sHandler struct {
	service  *%s.Service
	validate *validator.Validate
}

// New%sHandler creates a new %sHandler.
func New%sHandler(service *%s.Service) *%sHandler {
	return &%sHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Create%s creates a new %s.
// @Summary Create a new %s
// @Tags %s
// @Accept json
// @Produce json
// @Param %s body Create%sRequest true "%s data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/%s [post]
func (h *%sHandler) Create%s(c *fiber.Ctx) error {
	var req Create%sRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	item := &models.%s{
%s
	}
	if err := h.service.Create(c.Context(), item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create %s",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   item,
	})
}

// Get%s retrieves a %s by ID.
// @Summary Get a %s by ID
// @Tags %s
// @Produce json
// @Param id path int true "%s ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/%s/{id} [get]
func (h *%sHandler) Get%s(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid ID",
		})
	}
	item, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get %s",
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "%s not found",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   item,
	})
}

// GetAll%s retrieves all %s with pagination.
// @Summary Get all %s
// @Tags %s
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/%s [get]
func (h *%sHandler) GetAll%s(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	items, total, err := h.service.GetAll(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get %s",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   items,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// Update%s updates an existing %s.
// @Summary Update a %s
// @Tags %s
// @Accept json
// @Produce json
// @Param id path int true "%s ID"
// @Param %s body Update%sRequest true "Updated %s data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/%s/{id} [put]
func (h *%sHandler) Update%s(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid ID",
		})
	}

	item, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get %s",
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "%s not found",
		})
	}

	var req Update%sRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

%s

	if err := h.service.Update(c.Context(), item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update %s",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   item,
	})
}

// Delete%s deletes a %s.
// @Summary Delete a %s
// @Tags %s
// @Param id path int true "%s ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/%s/{id} [delete]
func (h *%sHandler) Delete%s(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid ID",
		})
	}

	if err := h.service.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete %s",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "%s deleted successfully",
	})
}
`,
		// imports
		modulePath, lowerName, modulePath,
		// CreateRequest
		modelName, lowerName, modelName, strings.Join(createFields, "\n"),
		// UpdateRequest
		modelName, lowerName, modelName, strings.Join(updateFields, "\n"),
		// Handler struct
		modelName, lowerName, modelName, lowerName,
		// NewHandler
		modelName, modelName, modelName, lowerName, modelName, modelName,
		// Create handler
		modelName, lowerName, lowerName, pluralLower,
		lowerName, modelName, modelName,
		pluralLower,
		modelName, modelName, modelName,
		modelName, strings.Join(createAssignments, "\n"),
		lowerName,
		// Get handler
		modelName, lowerName, lowerName, pluralLower, modelName,
		pluralLower,
		modelName, modelName,
		lowerName,
		modelName,
		// GetAll handler
		pluralName, pluralLower, pluralLower, pluralLower,
		pluralLower,
		modelName, pluralName,
		pluralLower,
		// Update handler
		modelName, lowerName, lowerName, pluralLower,
		modelName,
		lowerName, modelName, lowerName,
		pluralLower,
		modelName, modelName,
		lowerName,
		modelName,
		modelName,
		strings.Join(updateAssignments, "\n"),
		lowerName,
		// Delete handler
		modelName, lowerName, lowerName, pluralLower, modelName,
		pluralLower,
		modelName, modelName,
		lowerName,
		modelName,
	)

	// Add relation-specific handlers if BelongsTo is set
	if relations != nil && relations.BelongsTo != "" {
		parentName := relations.BelongsTo
		parentLower := strings.ToLower(parentName)
		parentPlural := pluralize(parentLower)

		relationHandlers := fmt.Sprintf(`

// GetBy%s retrieves all %s for a given %s.
// @Summary Get %s by %s ID
// @Tags %s
// @Produce json
// @Param %sId path int true "%s ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/%s/{%sId}/%s [get]
func (h *%sHandler) GetBy%s(c *fiber.Ctx) error {
	parentID, err := strconv.ParseUint(c.Params("%sId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid %s ID",
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	items, total, err := h.service.GetBy%sID(c.Context(), uint(parentID), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get %s",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   items,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// CreateFor%s creates a new %s under a specific %s.
// @Summary Create a %s for a %s
// @Tags %s
// @Accept json
// @Produce json
// @Param %sId path int true "%s ID"
// @Param %s body Create%sRequest true "%s data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/%s/{%sId}/%s [post]
func (h *%sHandler) CreateFor%s(c *fiber.Ctx) error {
	parentID, err := strconv.ParseUint(c.Params("%sId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid %s ID",
		})
	}

	var req Create%sRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	item := &models.%s{
%s
		%sID: uint(parentID),
	}
	if err := h.service.Create(c.Context(), item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create %s",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   item,
	})
}
`,
			// GetByParent
			parentName, pluralLower, parentLower,
			pluralLower, parentLower,
			pluralLower,
			parentLower, parentName,
			parentPlural, parentLower, pluralLower,
			modelName, parentName,
			parentLower,
			parentLower,
			parentName,
			pluralLower,
			// CreateForParent
			parentName, lowerName, parentLower,
			lowerName, parentLower,
			pluralLower,
			parentLower, parentName,
			lowerName, modelName, modelName,
			parentPlural, parentLower, pluralLower,
			modelName, parentName,
			parentLower,
			parentLower,
			modelName,
			modelName,
			strings.Join(createAssignments, "\n"),
			parentName,
			lowerName,
		)
		result += relationHandlers
	}

	// Add ?include= support on GetByID - modify the import to include "strings" if needed
	if relations != nil && relations.BelongsTo != "" {
		// We need "strings" import - inject it into the handler file
		result = strings.Replace(result, "\"strconv\"", "\"strconv\"\n\t\"strings\"", 1)
	}

	return result
}

// updateHandlerModule adds a new handler to the handlers fx.Module in module.go.
func updateHandlerModule(projectDir, modulePath, modelName string) error {
	moduleFile := filepath.Join(projectDir, "internal", "adapters", "handlers", "module.go")

	content, err := os.ReadFile(moduleFile)
	if err != nil {
		return fmt.Errorf("failed to read handlers module.go: %w", err)
	}

	contentStr := string(content)
	lowerName := strings.ToLower(modelName)

	// Check if already added (idempotency)
	constructorName := fmt.Sprintf("New%sHandler", modelName)
	if strings.Contains(contentStr, constructorName) {
		return nil
	}

	// Add the domain import
	domainImport := fmt.Sprintf("\"%s/internal/domain/%s\"", modulePath, lowerName)
	if !strings.Contains(contentStr, domainImport) {
		// Find the import block and add the new import
		importEnd := strings.Index(contentStr, ")")
		if importEnd == -1 {
			return fmt.Errorf("import block not found in handlers module.go")
		}
		contentStr = contentStr[:importEnd] + "\t" + domainImport + "\n" + contentStr[importEnd:]
	}

	// Find the last fx.Provide block within the handlers module
	moduleIdx := strings.Index(contentStr, "fx.Module(\"handlers\"")
	if moduleIdx == -1 {
		return fmt.Errorf("fx.Module(\"handlers\") not found in handlers module.go")
	}

	lastProvideEnd := -1
	for i := moduleIdx; i < len(contentStr); i++ {
		remaining := contentStr[i:]
		if strings.HasPrefix(remaining, "fx.Provide(") {
			depth := 0
			for j := i; j < len(contentStr); j++ {
				if contentStr[j] == '(' {
					depth++
				} else if contentStr[j] == ')' {
					depth--
					if depth == 0 {
						lastProvideEnd = j + 1
						i = j
						break
					}
				}
			}
		}
	}

	if lastProvideEnd == -1 {
		return fmt.Errorf("no fx.Provide found in handlers module.go")
	}

	newProvide := fmt.Sprintf(`,
	fx.Provide(func(s *%s.Service) *%sHandler {
		return New%sHandler(s)
	})`, lowerName, modelName, modelName)

	contentStr = contentStr[:lastProvideEnd] + newProvide + contentStr[lastProvideEnd:]

	return os.WriteFile(moduleFile, []byte(contentStr), 0644)
}

// updateRoutes adds new CRUD routes for the model to routes.go.
// If isPublic is true, routes are not protected by authMiddleware.
func updateRoutes(projectDir, modelName string, isPublic bool, relations *RelationConfig) error {
	routesFile := filepath.Join(projectDir, "internal", "adapters", "http", "routes.go")

	content, err := os.ReadFile(routesFile)
	if err != nil {
		return fmt.Errorf("failed to read routes.go: %w", err)
	}

	contentStr := string(content)
	lowerName := strings.ToLower(modelName)
	pluralLower := pluralize(lowerName)
	pluralName := pluralizePascal(modelName)

	// Check if already added (idempotency)
	handlerParam := fmt.Sprintf("%sHandler", lowerName)
	if strings.Contains(contentStr, handlerParam) {
		return nil
	}

	// 1. Add the handler parameter to RegisterRoutes signature
	// Find the closing ) of the function signature
	sigStart := strings.Index(contentStr, "func RegisterRoutes(")
	if sigStart == -1 {
		return fmt.Errorf("RegisterRoutes function not found in routes.go")
	}

	// Find the closing ) of the parameter list
	depth := 0
	sigEnd := -1
	for i := sigStart; i < len(contentStr); i++ {
		if contentStr[i] == '(' {
			depth++
		} else if contentStr[i] == ')' {
			depth--
			if depth == 0 {
				sigEnd = i
				break
			}
		}
	}
	if sigEnd == -1 {
		return fmt.Errorf("could not find closing ) of RegisterRoutes signature")
	}

	// Add new handler parameter before closing )
	// The Go template uses trailing commas, so we just add a new line with the param
	newParam := fmt.Sprintf("\n\t%sHandler *handlers.%sHandler,", lowerName, modelName)
	contentStr = contentStr[:sigEnd] + newParam + "\n" + contentStr[sigEnd:]

	// Recalculate positions after insertion
	// 2. Add route group before the closing } of the function
	// Find the last } in the file (closing of RegisterRoutes function body)
	lastBrace := strings.LastIndex(contentStr, "}")
	if lastBrace == -1 {
		return fmt.Errorf("could not find closing } of RegisterRoutes")
	}

	var routeGroup string
	if isPublic {
		routeGroup = fmt.Sprintf(`
	// %s routes (public)
	%s := v1.Group("/%s")
	%s.Post("", %sHandler.Create%s)
	%s.Get("", %sHandler.GetAll%s)
	%s.Get("/:id", %sHandler.Get%s)
	%s.Put("/:id", %sHandler.Update%s)
	%s.Delete("/:id", %sHandler.Delete%s)
`,
			modelName,
			pluralLower, pluralLower,
			pluralLower, lowerName, modelName,
			pluralLower, lowerName, pluralName,
			pluralLower, lowerName, modelName,
			pluralLower, lowerName, modelName,
			pluralLower, lowerName, modelName,
		)
	} else {
		routeGroup = fmt.Sprintf(`
	// %s routes (protected)
	%s := v1.Group("/%s", authMiddleware)
	%s.Post("", %sHandler.Create%s)
	%s.Get("", %sHandler.GetAll%s)
	%s.Get("/:id", %sHandler.Get%s)
	%s.Put("/:id", %sHandler.Update%s)
	%s.Delete("/:id", %sHandler.Delete%s)
`,
			modelName,
			pluralLower, pluralLower,
			pluralLower, lowerName, modelName,
			pluralLower, lowerName, pluralName,
			pluralLower, lowerName, modelName,
			pluralLower, lowerName, modelName,
			pluralLower, lowerName, modelName,
		)
	}

	contentStr = contentStr[:lastBrace] + routeGroup + contentStr[lastBrace:]

	// Add nested routes for belongs-to relations
	if relations != nil && relations.BelongsTo != "" {
		parentLower := strings.ToLower(relations.BelongsTo)
		parentPlural := pluralize(parentLower)

		// Find the closing } again (position shifted after previous insertion)
		lastBrace = strings.LastIndex(contentStr, "}")
		nestedRoutes := fmt.Sprintf(`
	// %s nested routes under %s
	%sNested := v1.Group("/%s/:%sId/%s")
	%sNested.Get("", %sHandler.GetBy%s)
	%sNested.Post("", %sHandler.CreateFor%s)
`,
			modelName, relations.BelongsTo,
			pluralLower, parentPlural, parentLower, pluralLower,
			pluralLower, lowerName, relations.BelongsTo,
			pluralLower, lowerName, relations.BelongsTo,
		)
		contentStr = contentStr[:lastBrace] + nestedRoutes + contentStr[lastBrace:]
	}

	return os.WriteFile(routesFile, []byte(contentStr), 0644)
}

// updateParentModelHasMany surgically inserts a has-many slice field into an existing parent model file.
// It finds the struct definition for the parent model, locates the CreatedAt field as an anchor,
// and inserts a slice field referencing the child model before it.
func updateParentModelHasMany(projectDir, childModelName, parentModelName string) error {
	parentLower := strings.ToLower(parentModelName)
	modelFile := filepath.Join(projectDir, "internal", "models", parentLower+".go")

	content, err := os.ReadFile(modelFile)
	if err != nil {
		return fmt.Errorf("failed to read parent model file %s: %w", modelFile, err)
	}

	contentStr := string(content)

	// Check if already added (idempotency)
	childPlural := pluralizePascal(childModelName)
	if strings.Contains(contentStr, childPlural+" ") {
		return nil
	}

	bt := "`"
	// Find anchor point to insert before it (try multiple anchors for robustness)
	anchors := []string{"CreatedAt", "UpdatedAt", "DeletedAt"}
	var anchorIdx = -1

	for _, anchor := range anchors {
		anchorIdx = strings.Index(contentStr, anchor)
		if anchorIdx != -1 {
			break
		}
	}

	if anchorIdx == -1 {
		return fmt.Errorf("could not find timestamp fields (CreatedAt/UpdatedAt/DeletedAt) in parent model %s - unable to determine insertion point", parentModelName)
	}

	// Find the start of the line containing the anchor
	lineStart := anchorIdx
	for lineStart > 0 && contentStr[lineStart-1] != '\n' {
		lineStart--
	}

	// Build the has-many field line with proper GORM tags
	hasMany := fmt.Sprintf("\t%s []%s %sgorm:\"foreignKey:%sID\" json:\"%s,omitempty\"%s\n",
		childPlural, childModelName, bt, parentModelName, strings.ToLower(childPlural), bt)

	contentStr = contentStr[:lineStart] + hasMany + contentStr[lineStart:]

	return os.WriteFile(modelFile, []byte(contentStr), 0644)
}

// updateMainFxModule adds the new domain module import and fx.Module entry to main.go.
func updateMainFxModule(projectDir, modulePath, modelName string) error {
	mainFile := filepath.Join(projectDir, "cmd", "main.go")

	content, err := os.ReadFile(mainFile)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	contentStr := string(content)
	lowerName := strings.ToLower(modelName)

	// Check if already added (idempotency)
	moduleRef := fmt.Sprintf("%s.Module", lowerName)
	if strings.Contains(contentStr, moduleRef) {
		return nil
	}

	// 1. Add the domain import
	domainImport := fmt.Sprintf("\"%s/internal/domain/%s\"", modulePath, lowerName)
	if !strings.Contains(contentStr, domainImport) {
		// Find the import block - look for the last import before the closing )
		// We insert after the last domain import or after the existing imports
		importClose := -1
		importOpen := strings.Index(contentStr, "import (")
		if importOpen == -1 {
			return fmt.Errorf("import block not found in main.go")
		}
		for i := importOpen; i < len(contentStr); i++ {
			if contentStr[i] == ')' {
				importClose = i
				break
			}
		}
		if importClose == -1 {
			return fmt.Errorf("could not find closing ) of import block in main.go")
		}
		contentStr = contentStr[:importClose] + "\t" + domainImport + "\n" + contentStr[importClose:]
	}

	// 2. Add the fx.Module entry after "// Domain services" section
	// Look for the user.Module line and add after it
	userModuleIdx := strings.Index(contentStr, "user.Module")
	if userModuleIdx == -1 {
		// Try to find repository.Module as fallback insertion point
		userModuleIdx = strings.Index(contentStr, "repository.Module")
		if userModuleIdx == -1 {
			return fmt.Errorf("could not find insertion point for fx.Module in main.go")
		}
	}

	// Find the end of the line containing the anchor
	lineEnd := strings.Index(contentStr[userModuleIdx:], ",")
	if lineEnd == -1 {
		return fmt.Errorf("could not find end of fx.Module line in main.go")
	}
	insertIdx := userModuleIdx + lineEnd + 1

	newModule := fmt.Sprintf("\n\t\t%s.Module,", lowerName)
	contentStr = contentStr[:insertIdx] + newModule + contentStr[insertIdx:]

	return os.WriteFile(mainFile, []byte(contentStr), 0644)
}

// generateServiceTestTemplate generates the content for internal/domain/<name>/service_test.go.
// It creates comprehensive unit tests with a mock repository for all CRUD operations:
// Create (success, error), GetByID (found, not found), GetAll (empty, with data),
// Update (success, error), Delete (success, error).
func generateServiceTestTemplate(modulePath, modelName string, fields []FieldDefinition, relations *RelationConfig) string {
	lowerName := strings.ToLower(modelName)

	// Build mock method implementations
	// The mock implements the repository interface methods
	mockMethods := generateMockMethods(modelName)

	// Build sample field assignments for test data
	sampleFields := generateSampleFields(fields)

	// Add extra mock methods for relations
	var extraMockMethods string
	var extraTests string
	if relations != nil && relations.BelongsTo != "" {
		extraMockMethods = fmt.Sprintf(`
func (m *Mock%sRepository) FindByIDWithRelations(ctx context.Context, id uint, includes []string) (*models.%s, error) {
	args := m.Called(ctx, id, includes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.%s), args.Error(1)
}

func (m *Mock%sRepository) FindBy%sID(ctx context.Context, parentID uint, page, limit int) ([]*models.%s, int64, error) {
	args := m.Called(ctx, parentID, page, limit)
	return args.Get(0).([]*models.%s), args.Get(1).(int64), args.Error(2)
}`,
			modelName, modelName, modelName,
			modelName, relations.BelongsTo, modelName, modelName,
		)

		extraTests = fmt.Sprintf(`

func TestService_GetByIDWithRelations(t *testing.T) {
	t.Run("with_includes", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		expected := &models.%s{%s}
		includes := []string{"%s"}
		repo.On("FindByIDWithRelations", mock.Anything, uint(1), includes).Return(expected, nil).Once()

		result, err := svc.GetByIDWithRelations(context.Background(), 1, includes)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		includes := []string{"%s"}
		repo.On("FindByIDWithRelations", mock.Anything, uint(999), includes).Return(nil, nil).Once()

		result, err := svc.GetByIDWithRelations(context.Background(), 999, includes)
		assert.NoError(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestService_GetBy%sID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		expected := []*models.%s{{%s}}
		repo.On("FindBy%sID", mock.Anything, uint(1), 1, 10).Return(expected, int64(1), nil).Once()

		items, total, err := svc.GetBy%sID(context.Background(), 1, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, int64(1), total)
		repo.AssertExpectations(t)
	})

	t.Run("empty", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		repo.On("FindBy%sID", mock.Anything, uint(999), 1, 10).Return([]*models.%s{}, int64(0), nil).Once()

		items, total, err := svc.GetBy%sID(context.Background(), 999, 1, 10)
		assert.NoError(t, err)
		assert.Empty(t, items)
		assert.Equal(t, int64(0), total)
		repo.AssertExpectations(t)
	})

	t.Run("pagination_defaults", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		repo.On("FindBy%sID", mock.Anything, uint(1), 1, 10).Return([]*models.%s{}, int64(0), nil).Once()

		_, _, err := svc.GetBy%sID(context.Background(), 1, 0, 0)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
`,
			// GetByIDWithRelations with_includes
			modelName, lowerName, modelName, sampleFields, relations.BelongsTo,
			// GetByIDWithRelations not_found
			modelName, lowerName, relations.BelongsTo,
			// GetByParentID success
			relations.BelongsTo,
			modelName, lowerName, modelName, sampleFields,
			relations.BelongsTo,
			relations.BelongsTo,
			// GetByParentID empty
			modelName, lowerName,
			relations.BelongsTo, modelName,
			relations.BelongsTo,
			// GetByParentID pagination_defaults
			modelName, lowerName,
			relations.BelongsTo, modelName,
			relations.BelongsTo,
		)
	}

	result := fmt.Sprintf(`package %s_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"%s/internal/domain/%s"
	"%s/internal/models"
)

// Mock%sRepository implements interfaces.%sRepository for testing.
type Mock%sRepository struct {
	mock.Mock
}

%s
%s

func TestService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		item := &models.%s{%s}
		repo.On("Create", mock.Anything, item).Return(nil).Once()

		err := svc.Create(context.Background(), item)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		item := &models.%s{%s}
		repo.On("Create", mock.Anything, item).Return(errors.New("db error")).Once()

		err := svc.Create(context.Background(), item)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_GetByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		expected := &models.%s{%s}
		repo.On("FindByID", mock.Anything, uint(1)).Return(expected, nil).Once()

		result, err := svc.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		repo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil).Once()

		result, err := svc.GetByID(context.Background(), 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestService_GetAll(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		repo.On("FindAll", mock.Anything, 1, 10).Return([]*models.%s{}, int64(0), nil).Once()

		items, total, err := svc.GetAll(context.Background(), 1, 10)
		assert.NoError(t, err)
		assert.Empty(t, items)
		assert.Equal(t, int64(0), total)
		repo.AssertExpectations(t)
	})

	t.Run("with_data", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		expected := []*models.%s{
			{%s},
			{%s},
		}
		repo.On("FindAll", mock.Anything, 1, 10).Return(expected, int64(2), nil).Once()

		items, total, err := svc.GetAll(context.Background(), 1, 10)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, int64(2), total)
		repo.AssertExpectations(t)
	})

	t.Run("pagination_defaults", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		repo.On("FindAll", mock.Anything, 1, 10).Return([]*models.%s{}, int64(0), nil).Once()

		_, _, err := svc.GetAll(context.Background(), 0, 0)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		item := &models.%s{%s}
		repo.On("Update", mock.Anything, item).Return(nil).Once()

		err := svc.Update(context.Background(), item)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		item := &models.%s{%s}
		repo.On("Update", mock.Anything, item).Return(errors.New("db error")).Once()

		err := svc.Update(context.Background(), item)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		repo.On("Delete", mock.Anything, uint(1)).Return(nil).Once()

		err := svc.Delete(context.Background(), 1)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(Mock%sRepository)
		svc := %s.NewService(repo)

		repo.On("Delete", mock.Anything, uint(1)).Return(errors.New("db error")).Once()

		err := svc.Delete(context.Background(), 1)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
%s`,
		// package
		lowerName,
		// imports
		modulePath, lowerName, modulePath,
		// Mock struct
		modelName, modelName, modelName,
		// Mock methods
		mockMethods,
		// Extra mock methods for relations
		extraMockMethods,
		// Create success
		modelName, lowerName, modelName, sampleFields,
		// Create error
		modelName, lowerName, modelName, sampleFields,
		// GetByID found
		modelName, lowerName, modelName, sampleFields,
		// GetByID not found
		modelName, lowerName,
		// GetAll empty
		modelName, lowerName, modelName,
		// GetAll with data
		modelName, lowerName, modelName, sampleFields, sampleFields,
		// GetAll pagination
		modelName, lowerName, modelName,
		// Update success
		modelName, lowerName, modelName, sampleFields,
		// Update error
		modelName, lowerName, modelName, sampleFields,
		// Delete success
		modelName, lowerName,
		// Delete error
		modelName, lowerName,
		// Extra relation tests
		extraTests,
	)
	return result
}

// generateMockMethods generates mock repository method implementations for testing.
func generateMockMethods(modelName string) string {
	lowerName := strings.ToLower(modelName)
	return fmt.Sprintf(`func (m *Mock%sRepository) Create(ctx context.Context, %s *models.%s) error {
	args := m.Called(ctx, %s)
	return args.Error(0)
}

func (m *Mock%sRepository) FindByID(ctx context.Context, id uint) (*models.%s, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.%s), args.Error(1)
}

func (m *Mock%sRepository) FindAll(ctx context.Context, page, limit int) ([]*models.%s, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*models.%s), args.Get(1).(int64), args.Error(2)
}

func (m *Mock%sRepository) Update(ctx context.Context, %s *models.%s) error {
	args := m.Called(ctx, %s)
	return args.Error(0)
}

func (m *Mock%sRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}`,
		modelName, lowerName, modelName, lowerName,
		modelName, modelName, modelName,
		modelName, modelName, modelName,
		modelName, lowerName, modelName, lowerName,
		modelName,
	)
}

// generateSampleFields returns sample field value assignments for test data creation.
func generateSampleFields(fields []FieldDefinition) string {
	var parts []string
	for _, f := range fields {
		var val string
		switch f.Type {
		case "string":
			val = fmt.Sprintf(`%s: "test"`, f.Name)
		case "int":
			val = fmt.Sprintf(`%s: 1`, f.Name)
		case "uint":
			val = fmt.Sprintf(`%s: 1`, f.Name)
		case "float64":
			val = fmt.Sprintf(`%s: 1.0`, f.Name)
		case "bool":
			val = fmt.Sprintf(`%s: false`, f.Name)
		case "time.Time":
			// Skip time fields in test samples for simplicity
			continue
		default:
			continue
		}
		parts = append(parts, val)
	}
	return strings.Join(parts, ", ")
}

// generateHandlerTestTemplate generates the content for internal/adapters/handlers/<name>_handler_test.go.
// It creates comprehensive HTTP handler tests for all CRUD endpoints with success and error cases.
// Uses integration-style testing with mock repository rather than mock service.
func generateHandlerTestTemplate(modulePath, modelName string, fields []FieldDefinition, relations *RelationConfig) string {
	lowerName := strings.ToLower(modelName)
	pluralName := pluralizePascal(modelName)
	pluralLower := strings.ToLower(pluralName)

	// Build JSON body for create request
	createJSONFields := generateCreateJSONFields(fields)
	invalidJSONFields := generateInvalidJSONFields(fields)

	// Build mock methods for repository interface
	repoMockMethods := generateMockMethods(modelName)

	// Add extra mock methods and route setup for relations
	var extraMockMethods string
	var extraRouteSetup string
	var extraTests string
	if relations != nil && relations.BelongsTo != "" {
		parentName := relations.BelongsTo
		parentLower := strings.ToLower(parentName)
		parentPlural := pluralize(parentLower)

		extraMockMethods = fmt.Sprintf(`
func (m *Mock%sRepository) FindByIDWithRelations(ctx context.Context, id uint, includes []string) (*models.%s, error) {
	args := m.Called(ctx, id, includes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.%s), args.Error(1)
}

func (m *Mock%sRepository) FindBy%sID(ctx context.Context, parentID uint, page, limit int) ([]*models.%s, int64, error) {
	args := m.Called(ctx, parentID, page, limit)
	return args.Get(0).([]*models.%s), args.Get(1).(int64), args.Error(2)
}`,
			modelName, modelName, modelName,
			modelName, parentName, modelName, modelName,
		)

		extraRouteSetup = fmt.Sprintf(`
	// Nested routes for %s relation
	%sGroup := api.Group("/%s/:%sId/%s")
	%sGroup.Get("", handler.GetBy%s)
	%sGroup.Post("", handler.CreateFor%s)`,
			parentLower,
			parentLower, parentPlural, parentLower, pluralLower,
			parentLower, parentName,
			parentLower, parentName,
		)

		extraTests = fmt.Sprintf(`

func TestGetBy%s_Success(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	expected := []*models.%s{{%s}}
	mockRepo.On("FindBy%sID", mock.Anything, uint(1), 1, 10).Return(expected, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/%s/1/%s", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "success", result["status"])
	mockRepo.AssertExpectations(t)
}

func TestGetBy%s_InvalidParentID(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/%s/abc/%s", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateFor%s_Success(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.%s")).Return(nil).Once()

	body := %s
	req := httptest.NewRequest(http.MethodPost, "/api/v1/%s/1/%s", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "success", result["status"])
	mockRepo.AssertExpectations(t)
}

func TestCreateFor%s_InvalidParentID(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	body := %s
	req := httptest.NewRequest(http.MethodPost, "/api/v1/%s/abc/%s", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
`,
			// GetByParent success
			parentName,
			modelName, modelName,
			modelName, generateSampleFields(fields),
			parentName,
			parentPlural, pluralLower,
			// GetByParent invalid parent ID
			parentName,
			modelName, modelName,
			parentPlural, pluralLower,
			// CreateForParent success
			parentName,
			modelName, modelName,
			modelName,
			createJSONFields,
			parentPlural, pluralLower,
			// CreateForParent invalid parent ID
			parentName,
			modelName, modelName,
			createJSONFields,
			parentPlural, pluralLower,
		)
	}

	result := fmt.Sprintf(`package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"%s/internal/adapters/handlers"
	"%s/internal/domain/%s"
	"%s/internal/models"
)

// Mock%sRepository mocks the repository for handler testing.
type Mock%sRepository struct {
	mock.Mock
}

%s
%s

func setup%sTestApp(mockRepo *Mock%sRepository) *fiber.App {
	app := fiber.New()
	svc := %s.NewService(mockRepo)
	handler := handlers.New%sHandler(svc)
	api := app.Group("/api/v1")
	%s := api.Group("/%s")
	%s.Post("", handler.Create%s)
	%s.Get("", handler.GetAll%s)
	%s.Get("/:id", handler.Get%s)
	%s.Put("/:id", handler.Update%s)
	%s.Delete("/:id", handler.Delete%s)
%s
	return app
}

func TestCreate%s_Success(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.%s")).Return(nil).Once()

	body := %s
	req := httptest.NewRequest(http.MethodPost, "/api/v1/%s", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "success", result["status"])
	mockRepo.AssertExpectations(t)
}

func TestCreate%s_InvalidBody(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/%s", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreate%s_ValidationError(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	body := %s
	req := httptest.NewRequest(http.MethodPost, "/api/v1/%s", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGet%s_Found(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	expected := &models.%s{%s}
	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(expected, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/%s/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "success", result["status"])
	mockRepo.AssertExpectations(t)
}

func TestGet%s_NotFound(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/%s/999", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGet%s_InvalidID(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/%s/abc", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetAll%s_EmptyList(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	mockRepo.On("FindAll", mock.Anything, 1, 10).Return([]*models.%s{}, int64(0), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/%s", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGetAll%s_WithPagination(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	items := []*models.%s{{%s}}
	mockRepo.On("FindAll", mock.Anything, 2, 5).Return(items, int64(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/%s?page=2&limit=5", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	meta := result["meta"].(map[string]interface{})
	assert.Equal(t, float64(2), meta["page"])
	assert.Equal(t, float64(5), meta["limit"])
	mockRepo.AssertExpectations(t)
}

func TestUpdate%s_Success(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	existing := &models.%s{%s}
	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil).Once()
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.%s")).Return(nil).Once()

	body := %s
	req := httptest.NewRequest(http.MethodPut, "/api/v1/%s/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdate%s_NotFound(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil).Once()

	body := %s
	req := httptest.NewRequest(http.MethodPut, "/api/v1/%s/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdate%s_ValidationError(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	existing := &models.%s{%s}
	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/api/v1/%s/1", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDelete%s_Success(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/%s/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestDelete%s_NotFound(t *testing.T) {
	mockRepo := new(Mock%sRepository)
	app := setup%sTestApp(mockRepo)

	mockRepo.On("Delete", mock.Anything, uint(999)).Return(errors.New("not found")).Once()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/%s/999", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}
%s`,
		// imports
		modulePath, modulePath, lowerName, modulePath,
		// Mock struct
		modelName, modelName,
		// Repository mock methods
		repoMockMethods,
		// Extra mock methods for relations
		extraMockMethods,
		// setup func
		modelName, modelName, lowerName, modelName,
		pluralLower, pluralLower,
		pluralLower, modelName,
		pluralLower, pluralName,
		pluralLower, modelName,
		pluralLower, modelName,
		pluralLower, modelName,
		// Extra route setup for relations
		extraRouteSetup,
		// Create success
		modelName, modelName, modelName, modelName,
		createJSONFields, pluralLower,
		// Create invalid body
		modelName, modelName, modelName, pluralLower,
		// Create validation error
		modelName, modelName, modelName, invalidJSONFields, pluralLower,
		// Get found
		modelName, modelName, modelName,
		modelName, generateSampleFields(fields),
		pluralLower,
		// Get not found
		modelName, modelName, modelName, pluralLower,
		// Get invalid ID
		modelName, modelName, modelName, pluralLower,
		// GetAll empty
		pluralName, modelName, modelName, modelName, pluralLower,
		// GetAll with pagination
		pluralName, modelName, modelName,
		modelName, generateSampleFields(fields),
		pluralLower,
		// Update success
		modelName, modelName, modelName,
		modelName, generateSampleFields(fields),
		modelName,
		createJSONFields, pluralLower,
		// Update not found
		modelName, modelName, modelName,
		createJSONFields, pluralLower,
		// Update validation error
		modelName, modelName, modelName,
		modelName, generateSampleFields(fields),
		pluralLower,
		// Delete success
		modelName, modelName, modelName, pluralLower,
		// Delete not found
		modelName, modelName, modelName, pluralLower,
		// Extra relation tests
		extraTests,
	)
	return result
}

// generateCreateJSONFields returns a JSON string literal for test request bodies.
func generateCreateJSONFields(fields []FieldDefinition) string {
	var parts []string
	for _, f := range fields {
		switch f.Type {
		case "string":
			parts = append(parts, fmt.Sprintf(`"%s":"test value"`, f.JSONName))
		case "int", "uint":
			parts = append(parts, fmt.Sprintf(`"%s":1`, f.JSONName))
		case "float64":
			parts = append(parts, fmt.Sprintf(`"%s":1.5`, f.JSONName))
		case "bool":
			parts = append(parts, fmt.Sprintf(`"%s":false`, f.JSONName))
		}
	}
	return "`{" + strings.Join(parts, ",") + "}`"
}

// generateInvalidJSONFields returns a JSON string literal with empty/invalid values for validation testing.
func generateInvalidJSONFields(fields []FieldDefinition) string {
	var parts []string
	for _, f := range fields {
		switch f.Type {
		case "string":
			parts = append(parts, fmt.Sprintf(`"%s":""`, f.JSONName))
		case "int", "uint":
			parts = append(parts, fmt.Sprintf(`"%s":0`, f.JSONName))
		case "float64":
			parts = append(parts, fmt.Sprintf(`"%s":0`, f.JSONName))
		case "bool":
			parts = append(parts, fmt.Sprintf(`"%s":false`, f.JSONName))
		}
	}
	return "`{" + strings.Join(parts, ",") + "}`"
}
