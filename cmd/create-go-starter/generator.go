package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tky0065/go-starter-kit/pkg/utils"
)

// getDirectoriesForTemplate returns the list of directories to create based on the template type.
// Minimal template has fewer directories (no auth, user, handlers, repository).
// Full template has all directories for complete hexagonal architecture.
func getDirectoriesForTemplate(template string) []string {
	// Common directories for all templates
	commonDirs := []string{
		"cmd",
		"internal/adapters/http",
		"internal/infrastructure/database",
		"internal/infrastructure/server",
		"pkg/config",
		"pkg/logger",
		"docs",
		"deployments",
		".github/workflows",
	}

	switch template {
	case TemplateMinimal:
		// Minimal template: only basic infrastructure, no auth
		return commonDirs
	case TemplateGraphQL:
		// GraphQL template: includes graph directories for gqlgen
		graphqlDirs := append(commonDirs,
			"internal/interfaces",
			"internal/models",
			"graph",
			"graph/model",
			"graph/generated",
		)
		return graphqlDirs
	case TemplateFull:
		// Full template: includes auth, user management, handlers, repository
		fullDirs := append(commonDirs,
			"pkg/auth",
			"internal/domain",
			"internal/domain/user",
			"internal/interfaces",
			"internal/models",
			"internal/adapters/middleware",
			"internal/adapters/handlers",
			"internal/adapters/repository",
		)
		return fullDirs
	default:
		// Default to full template directories (defensive programming)
		fullDirs := append(commonDirs,
			"pkg/auth",
			"internal/domain",
			"internal/domain/user",
			"internal/interfaces",
			"internal/models",
			"internal/adapters/middleware",
			"internal/adapters/handlers",
			"internal/adapters/repository",
		)
		return fullDirs
	}
}

// FileGenerator represents a file to be generated
type FileGenerator struct {
	Path    string
	Content string
}

// generateProjectFiles creates all the initial project files with templates.
// The template parameter specifies the type of project to generate (minimal, full, graphql).
// For this story (6.1), only the "full" template is implemented. Other templates will be implemented in future stories.
// This switch statement clarifies intent and returns an explicit error for unimplemented templates.
func generateProjectFiles(projectPath, projectName, template string) error {
	// Validate that the project directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// Validate the module name
	if err := utils.ValidateGoModuleName(projectName); err != nil { // Updated call
		return err
	}

	switch template {
	case "full":
		return generateFullTemplateFiles(projectPath, projectName)
	case "minimal":
		return generateMinimalTemplateFiles(projectPath, projectName)
	case "graphql":
		return generateGraphQLTemplateFiles(projectPath, projectName)
	default:
		// This case should ideally not be reached if validateTemplate is called beforehand.
		return fmt.Errorf("unsupported template '%s'", template)
	}
}

// generateFullTemplateFiles generates all files for the "full" template.
// This function was extracted from the original generateProjectFiles to improve modularity.
func generateFullTemplateFiles(projectPath, projectName string) error {
	// Create templates instance
	templates := NewProjectTemplates(projectName)

	// Define all files to generate
	files := []FileGenerator{
		{
			Path:    filepath.Join(projectPath, "go.mod"),
			Content: templates.GoModTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "cmd", "main.go"),
			Content: templates.UpdatedMainGoTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "config", "env.go"),
			Content: templates.ConfigTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "logger", "logger.go"),
			Content: templates.LoggerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "auth", "jwt.go"),
			Content: templates.JWTAuthTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "auth", "middleware.go"),
			Content: templates.JWTMiddlewareTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "auth", "module.go"),
			Content: templates.AuthModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "errors.go"),
			Content: templates.DomainErrorsTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "models", "user.go"),
			Content: templates.ModelsUserTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "user", "service.go"),
			Content: templates.UserServiceTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "domain", "user", "module.go"),
			Content: templates.UserModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "interfaces", "services.go"),
			Content: templates.UserInterfacesTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "interfaces", "user_repository.go"),
			Content: templates.UserRepositoryInterfaceTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "middleware", "error_handler.go"),
			Content: templates.ErrorHandlerMiddlewareTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "repository", "user_repository.go"),
			Content: templates.UserRepositoryTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "repository", "module.go"),
			Content: templates.RepositoryModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "auth_handler.go"),
			Content: templates.AuthHandlerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "user_handler.go"),
			Content: templates.UserHandlerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "module.go"),
			Content: templates.HandlerModuleTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "http", "health.go"),
			Content: templates.HealthHandlerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "http", "routes.go"),
			Content: templates.RoutesTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "database", "database.go"),
			Content: templates.DatabaseTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "server", "server.go"),
			Content: templates.ServerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "Dockerfile"),
			Content: templates.DockerfileTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docker-compose.yml"),
			Content: templates.DockerComposeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "Makefile"),
			Content: templates.MakefileTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".env.example"),
			Content: templates.EnvTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".gitignore"),
			Content: templates.GitignoreTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".golangci.yml"),
			Content: templates.GolangCILintTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".github", "workflows", "ci.yml"),
			Content: templates.GitHubActionsWorkflowTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "README.md"),
			Content: templates.ReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "README.md"),
			Content: templates.DocsReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "docs.go"),
			Content: templates.SwaggerDocsTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "quick-start.md"),
			Content: templates.QuickStartTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "setup.sh"),
			Content: templates.SetupScriptTemplate(),
		},
	}

	// Write all files
	for _, file := range files {
		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}

		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}
	}

	// Make setup.sh executable
	setupPath := filepath.Join(projectPath, "setup.sh")
	if err := os.Chmod(setupPath, 0755); err != nil {
		return fmt.Errorf("failed to make setup.sh executable: %w", err)
	}

	return nil
}

// generateMinimalTemplateFiles generates all files for the "minimal" template.
// This template includes basic infrastructure without authentication.
func generateMinimalTemplateFiles(projectPath, projectName string) error {
	// Create templates instance
	templates := NewProjectTemplates(projectName)

	// Define all files to generate for minimal template
	// Note: No auth-related files (pkg/auth, internal/domain/user, handlers, repository)
	files := []FileGenerator{
		{
			Path:    filepath.Join(projectPath, "go.mod"),
			Content: templates.MinimalGoModTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "cmd", "main.go"),
			Content: templates.MinimalMainGoTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "config", "env.go"),
			Content: templates.ConfigTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "logger", "logger.go"),
			Content: templates.LoggerTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "http", "health.go"),
			Content: templates.HealthHandlerTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "http", "routes.go"),
			Content: templates.MinimalRoutesTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "database", "database.go"),
			Content: templates.MinimalDatabaseTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "server", "server.go"),
			Content: templates.MinimalServerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "Dockerfile"),
			Content: templates.DockerfileTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, "docker-compose.yml"),
			Content: templates.MinimalDockerComposeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "Makefile"),
			Content: templates.MakefileTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, ".env.example"),
			Content: templates.MinimalEnvTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".gitignore"),
			Content: templates.GitignoreTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, ".golangci.yml"),
			Content: templates.GolangCILintTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, ".github", "workflows", "ci.yml"),
			Content: templates.GitHubActionsWorkflowTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, "README.md"),
			Content: templates.MinimalReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "README.md"),
			Content: templates.MinimalDocsReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "docs.go"),
			Content: templates.SwaggerDocsTemplate(), // Same as full template
		},
		{
			Path:    filepath.Join(projectPath, "docs", "quick-start.md"),
			Content: templates.MinimalQuickStartTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "setup.sh"),
			Content: templates.MinimalSetupScriptTemplate(),
		},
	}

	// Write all files
	for _, file := range files {
		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}

		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}
	}

	// Make setup.sh executable
	setupPath := filepath.Join(projectPath, "setup.sh")
	if err := os.Chmod(setupPath, 0755); err != nil {
		return fmt.Errorf("failed to make setup.sh executable: %w", err)
	}

	return nil
}

// generateGraphQLTemplateFiles generates all files for the "graphql" template.
// This template includes GraphQL support with gqlgen, gofiber/adaptor, and GORM.
func generateGraphQLTemplateFiles(projectPath, projectName string) error {
	// Create templates instance
	templates := NewProjectTemplates(projectName)

	// Define all files to generate for GraphQL template
	files := []FileGenerator{
		// Core Go files
		{
			Path:    filepath.Join(projectPath, "go.mod"),
			Content: templates.GraphQLGoModTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "cmd", "main.go"),
			Content: templates.GraphQLMainGoTemplate(),
		},
		// GraphQL files
		{
			Path:    filepath.Join(projectPath, "gqlgen.yml"),
			Content: templates.GqlGenYmlTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "graph", "schema.graphqls"),
			Content: templates.GraphQLSchemaTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "graph", "resolver.go"),
			Content: templates.GraphQLResolverTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "graph", "schema.resolvers.go"),
			Content: templates.GraphQLSchemaResolversTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "graph", "generate.go"),
			Content: templates.GraphQLGenerateGoTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "graph", "schema.resolvers_test.go"),
			Content: templates.GraphQLResolverTestTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "graph", "model", "models.go"),
			Content: templates.GraphQLModelTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "graph", "generated", "generated.go"),
			Content: templates.GraphQLGeneratedTemplate(),
		},
		// Infrastructure
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "server", "server.go"),
			Content: templates.GraphQLServerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "database", "database.go"),
			Content: templates.GraphQLDatabaseTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "database", "user_repository.go"),
			Content: templates.GraphQLUserRepositoryTemplate(),
		},
		// Domain
		{
			Path:    filepath.Join(projectPath, "internal", "interfaces", "user_repository.go"),
			Content: templates.GraphQLInterfacesTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "models", "user.go"),
			Content: templates.GraphQLModelsUserTemplate(),
		},
		// Packages
		{
			Path:    filepath.Join(projectPath, "pkg", "config", "env.go"),
			Content: templates.ConfigTemplate(), // Reuse from base templates
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "logger", "logger.go"),
			Content: templates.LoggerTemplate(), // Reuse from base templates
		},
		// Configuration files
		{
			Path:    filepath.Join(projectPath, ".env.example"),
			Content: templates.GraphQLEnvTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".gitignore"),
			Content: templates.GitignoreTemplate(), // Reuse from base templates
		},
		{
			Path:    filepath.Join(projectPath, ".golangci.yml"),
			Content: templates.GolangCILintTemplate(), // Reuse from base templates
		},
		{
			Path:    filepath.Join(projectPath, ".github", "workflows", "ci.yml"),
			Content: templates.GitHubActionsWorkflowTemplate(), // Reuse from base templates
		},
		// Build files
		{
			Path:    filepath.Join(projectPath, "Dockerfile"),
			Content: templates.DockerfileTemplate(), // Reuse from base templates
		},
		{
			Path:    filepath.Join(projectPath, "docker-compose.yml"),
			Content: templates.GraphQLDockerComposeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "Makefile"),
			Content: templates.GraphQLMakefileTemplate(),
		},
		// Documentation
		{
			Path:    filepath.Join(projectPath, "README.md"),
			Content: templates.GraphQLReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "README.md"),
			Content: templates.GraphQLDocsReadmeTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docs", "quick-start.md"),
			Content: templates.GraphQLQuickStartTemplate(),
		},
		// Setup script
		{
			Path:    filepath.Join(projectPath, "setup.sh"),
			Content: templates.GraphQLSetupScriptTemplate(),
		},
	}

	// Write all files
	for _, file := range files {
		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}

		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}
	}

	// Make setup.sh executable
	setupPath := filepath.Join(projectPath, "setup.sh")
	if err := os.Chmod(setupPath, 0755); err != nil {
		return fmt.Errorf("failed to make setup.sh executable: %w", err)
	}

	return nil
}
