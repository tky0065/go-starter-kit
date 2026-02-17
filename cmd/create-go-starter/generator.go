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
			"pkg/tracing",
			"pkg/metrics",
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
			"pkg/tracing",
			"pkg/metrics",
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
// The database parameter specifies the database type (postgres, mysql, sqlite, mongodb).
// The observabilityLevel parameter specifies the observability mode (none, basic, advanced).
// Database-specific templates are used to generate the correct driver, DSN, and configurations.
// This switch statement clarifies intent and returns an explicit error for unimplemented templates.
func generateProjectFiles(projectPath, projectName, template, database, observabilityLevel string) error {
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
		return generateFullTemplateFiles(projectPath, projectName, database, observabilityLevel)
	case "minimal":
		return generateMinimalTemplateFiles(projectPath, projectName, database)
	case "graphql":
		return generateGraphQLTemplateFiles(projectPath, projectName, database)
	default:
		// This case should ideally not be reached if validateTemplate is called beforehand.
		return fmt.Errorf("unsupported template '%s'", template)
	}
}

// generateFullTemplateFiles generates all files for the "full" template.
// This function was extracted from the original generateProjectFiles to improve modularity.
// The observabilityLevel parameter controls which observability files are generated.
func generateFullTemplateFiles(projectPath, projectName, database, observabilityLevel string) error {
	// Create templates instance with database support
	templates := NewProjectTemplatesWithDatabase(projectName, database)

	// Determine file contents based on observability level.
	// When observability=advanced, key files are replaced with enriched versions that include
	// Prometheus metrics, OpenTelemetry tracing, and Jaeger docker-compose configuration.
	goModContent := templates.GoModTemplate()
	mainGoContent := templates.UpdatedMainGoTemplate()
	serverContent := templates.ServerTemplate()
	loggerContent := templates.LoggerTemplate()
	envContent := templates.EnvTemplate()
	dockerComposeContent := templates.DockerComposeTemplate()
	// Health handler: base version (no metrics) for all projects; overridden below for advanced.
	healthHandlerContent := templates.AdvancedHealthHandlerTemplate()
	obsTemplates := NewObservabilityTemplates(projectName)

	if observabilityLevel == ObservabilityAdvanced {
		goModContent = obsTemplates.GoModTemplateWithObservability(database)
		mainGoContent = obsTemplates.MainGoTemplateWithObservability()
		serverContent = obsTemplates.ServerTemplateWithObservability()
		loggerContent = obsTemplates.LoggerWithTracingTemplate()
		envContent = obsTemplates.EnvTemplateWithObservability(database)
		dockerComposeContent = obsTemplates.DockerComposeTemplateWithObservability(database)
		// Advanced health handler: includes Prometheus metrics instrumentation (AC4)
		healthHandlerContent = obsTemplates.HealthHandlerWithMetricsTemplate()
	}

	// Define all files to generate
	files := []FileGenerator{
		{
			Path:    filepath.Join(projectPath, "go.mod"),
			Content: goModContent,
		},
		{
			Path:    filepath.Join(projectPath, "cmd", "main.go"),
			Content: mainGoContent,
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "config", "env.go"),
			Content: templates.ConfigTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "logger", "logger.go"),
			Content: loggerContent,
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
			// Advanced health handler with Liveness/Readiness K8s probes (Story 9.3)
			// Replaces the simple health.go — provides /health, /health/liveness, /health/readiness
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "health_handler.go"),
			Content: healthHandlerContent,
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "http", "routes.go"),
			Content: templates.RoutesTemplate(),
		},
		{
			// K8s probe configuration reference (Story 9.3, AC6) — generated for all full projects
			Path:    filepath.Join(projectPath, "deployments", "kubernetes", "probes.yaml"),
			Content: obsTemplates.KubernetesProbesTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "database", "database.go"),
			Content: templates.DatabaseTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "infrastructure", "server", "server.go"),
			Content: serverContent,
		},
		{
			Path:    filepath.Join(projectPath, "Dockerfile"),
			Content: templates.DockerfileTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "docker-compose.yml"),
			Content: dockerComposeContent,
		},
		{
			Path:    filepath.Join(projectPath, "Makefile"),
			Content: templates.MakefileTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, ".env.example"),
			Content: envContent,
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

	// Generate observability files conditionally (AC: 2, 3, 4, 7)
	if observabilityLevel == ObservabilityAdvanced {
		if err := generateObservabilityFiles(projectPath, projectName, database); err != nil {
			return fmt.Errorf("generating observability files: %w", err)
		}
	}

	return nil
}

// generateObservabilityFiles generates all observability files for the advanced mode.
// This includes Prometheus metrics files (Story 9.1), OpenTelemetry tracing files (Story 9.2),
// and Grafana/Prometheus infrastructure files (Story 9.4):
//   - pkg/metrics/prometheus.go
//   - internal/adapters/middleware/metrics_middleware.go
//   - internal/adapters/handlers/metrics_handler.go
//   - pkg/tracing/tracer.go
//   - pkg/tracing/db_tracing.go
//   - internal/adapters/middleware/tracing_middleware.go
//   - deployments/prometheus/prometheus.yml
//   - deployments/prometheus/rules/api_alerts.yml
//   - deployments/grafana/provisioning/datasources/prometheus.yaml
//   - deployments/grafana/provisioning/dashboards/default.yaml
//   - deployments/grafana/dashboards/api-dashboard.json
func generateObservabilityFiles(projectPath, projectName, database string) error {
	obsTemplates := NewObservabilityTemplates(projectName)

	files := []FileGenerator{
		// Prometheus metrics (Story 9.1)
		{
			Path:    filepath.Join(projectPath, "pkg", "metrics", "prometheus.go"),
			Content: obsTemplates.PrometheusTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "middleware", "metrics_middleware.go"),
			Content: obsTemplates.MetricsMiddlewareTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "handlers", "metrics_handler.go"),
			Content: obsTemplates.MetricsHandlerTemplate(),
		},
		// OpenTelemetry distributed tracing (Story 9.2)
		{
			Path:    filepath.Join(projectPath, "pkg", "tracing", "tracer.go"),
			Content: obsTemplates.TracerTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "pkg", "tracing", "db_tracing.go"),
			Content: obsTemplates.GORMTracingTemplate(database),
		},
		{
			Path:    filepath.Join(projectPath, "internal", "adapters", "middleware", "tracing_middleware.go"),
			Content: obsTemplates.TracingMiddlewareTemplate(),
		},
		// Prometheus + Grafana infrastructure (Story 9.4)
		{
			Path:    filepath.Join(projectPath, "deployments", "prometheus", "prometheus.yml"),
			Content: obsTemplates.PrometheusConfigTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "deployments", "prometheus", "rules", "api_alerts.yml"),
			Content: obsTemplates.PrometheusAlertRulesTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "deployments", "grafana", "provisioning", "datasources", "prometheus.yaml"),
			Content: obsTemplates.GrafanaDatasourceTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "deployments", "grafana", "provisioning", "dashboards", "default.yaml"),
			Content: obsTemplates.GrafanaDashboardProvisionTemplate(),
		},
		{
			Path:    filepath.Join(projectPath, "deployments", "grafana", "dashboards", "api-dashboard.json"),
			Content: obsTemplates.GrafanaDashboardJSONTemplate(),
		},
	}

	for _, file := range files {
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}
		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}
	}

	return nil
}

// generateMinimalTemplateFiles generates all files for the "minimal" template.
// This template includes basic infrastructure without authentication.
func generateMinimalTemplateFiles(projectPath, projectName, database string) error {
	// Create templates instance with database support
	templates := NewProjectTemplatesWithDatabase(projectName, database)

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
func generateGraphQLTemplateFiles(projectPath, projectName, database string) error {
	// Create templates instance with database support
	templates := NewProjectTemplatesWithDatabase(projectName, database)

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
