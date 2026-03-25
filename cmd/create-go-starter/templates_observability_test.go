package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestValidateObservabilityValid tests validateObservabilityLevel with valid values (AC: 1)
func TestValidateObservabilityValid(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"none level", ObservabilityNone},
		{"basic level", ObservabilityBasic},
		{"advanced level", ObservabilityAdvanced},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateObservabilityLevel(tt.level)
			if err != nil {
				t.Errorf("validateObservabilityLevel(%q) returned error: %v, want nil", tt.level, err)
			}
		})
	}
}

// TestValidateObservabilityInvalid tests validateObservabilityLevel with invalid values (AC: 1)
func TestValidateObservabilityInvalid(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"empty level", ""},
		{"unknown level", "monitoring"},
		{"invalid case", "ADVANCED"},
		{"invalid spacing", " advanced"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateObservabilityLevel(tt.level)
			if err == nil {
				t.Errorf("validateObservabilityLevel(%q) returned nil, want error", tt.level)
			}
			if !strings.Contains(err.Error(), "invalid observability") {
				t.Errorf("validateObservabilityLevel(%q) error = %v, want error containing 'invalid observability'", tt.level, err)
			}
		})
	}
}

// TestObservabilityDefaultValue tests that default observability level is "none" (9.2-AC: default flag value)
func TestObservabilityDefaultValue(t *testing.T) {
	if DefaultObservabilityLevel != ObservabilityNone {
		t.Errorf("DefaultObservabilityLevel = %q, want %q", DefaultObservabilityLevel, ObservabilityNone)
	}
}

// TestObservabilityFlagParsing tests that --observability flag is parsed via CLI (AC: 1)
func TestObservabilityFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantInOutput   string
		wantNoErr      bool
		cleanupProject string
	}{
		{
			name:           "observability none flag",
			args:           []string{"--observability=none", "test-obs-none"},
			wantInOutput:   "observability: none",
			wantNoErr:      true,
			cleanupProject: "test-obs-none",
		},
		{
			name:           "observability basic flag",
			args:           []string{"--observability=basic", "test-obs-basic"},
			wantInOutput:   "observability: basic",
			wantNoErr:      true,
			cleanupProject: "test-obs-basic",
		},
		{
			name:           "observability advanced flag",
			args:           []string{"--observability=advanced", "test-obs-advanced"},
			wantInOutput:   "observability: advanced",
			wantNoErr:      true,
			cleanupProject: "test-obs-advanced",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.cleanupProject)

			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.wantNoErr {
				if err != nil {
					t.Errorf("Expected no error, got: %v\nOutput: %s", err, string(output))
				}
				if !strings.Contains(string(output), tt.wantInOutput) {
					t.Errorf("Output should contain %q, got: %s", tt.wantInOutput, string(output))
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for observability '%s', got success", tt.args[0])
				}
			}
		})
	}
}

// TestInvalidObservabilityFlagError tests that invalid observability shows error (AC: 1)
func TestInvalidObservabilityFlagError(t *testing.T) {
	cmd := exec.Command(binaryPath, "--observability=invalid", "test-obs-invalid")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for invalid observability level")
		os.RemoveAll("test-obs-invalid")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "invalid observability") {
		t.Errorf("Expected 'invalid observability' in error, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "none, basic, advanced") {
		t.Errorf("Expected valid options in error, got: %s", outputStr)
	}
}

// TestObservabilityNoneGeneratesNoObservabilityFiles tests that --observability=none (9.2-AC: default behaviour)
// produces exactly the same project as without the flag (zero regression)
func TestObservabilityNoneGeneratesNoObservabilityFiles(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-obs-none-gen"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() with observability=none returned error: %v", err)
	}

	// Observability files should NOT exist for none mode
	observabilityFiles := []string{
		filepath.Join(projectName, "pkg", "metrics", "prometheus.go"),
		filepath.Join(projectName, "internal", "adapters", "middleware", "metrics_middleware.go"),
		filepath.Join(projectName, "internal", "adapters", "handlers", "metrics_handler.go"),
	}

	for _, file := range observabilityFiles {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			t.Errorf("Observability file %s should NOT exist for mode 'none', but it does", file)
		}
	}
}

// TestObservabilityAdvancedGeneratesFiles tests that --observability=advanced (AC: 2, 3, 4)
// generates the expected Prometheus/metrics files
func TestObservabilityAdvancedGeneratesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-obs-advanced-gen"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() with observability=advanced returned error: %v", err)
	}

	// Observability files SHOULD exist for advanced mode
	observabilityFiles := []string{
		filepath.Join(projectName, "pkg", "metrics", "prometheus.go"),
		filepath.Join(projectName, "internal", "adapters", "middleware", "metrics_middleware.go"),
		filepath.Join(projectName, "internal", "adapters", "handlers", "metrics_handler.go"),
	}

	for _, file := range observabilityFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Observability file %s should exist for mode 'advanced', but it doesn't", file)
		}
	}
}

// TestPrometheusGoContainsRegistry tests that prometheus.go contains key content (9.1-AC: 4)
func TestPrometheusGoContainsRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-prom-content"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	prometheusFile := filepath.Join(projectName, "pkg", "metrics", "prometheus.go")
	content, err := os.ReadFile(prometheusFile)
	if err != nil {
		t.Fatalf("Failed to read prometheus.go: %v", err)
	}

	contentStr := string(content)

	// Should contain Prometheus setup
	if !strings.Contains(contentStr, "NewPrometheusMetrics") {
		t.Errorf("prometheus.go should contain 'NewPrometheusMetrics', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "fiberprometheus") {
		t.Errorf("prometheus.go should contain 'fiberprometheus', got: %s", contentStr)
	}
}

// TestMetricsMiddlewareContainsHistogram tests that metrics_middleware.go contains latency histogram (9.1-AC: 5)
func TestMetricsMiddlewareContainsHistogram(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-metrics-middleware"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	middlewareFile := filepath.Join(projectName, "internal", "adapters", "middleware", "metrics_middleware.go")
	content, err := os.ReadFile(middlewareFile)
	if err != nil {
		t.Fatalf("Failed to read metrics_middleware.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Middleware") {
		t.Errorf("metrics_middleware.go should contain 'Middleware' function, got: %s", contentStr)
	}
}

// TestHelpShowsObservabilityFlag tests that --help shows observability flag documentation (AC: 1)
func TestHelpShowsObservabilityFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected exit code 0 for --help, but got error: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "-observability") {
		t.Errorf("Expected '--observability' flag in help, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "none") || !strings.Contains(outputStr, "advanced") {
		t.Errorf("Expected observability options in help, got: %s", outputStr)
	}
}

// TestAdvancedGoModIncludesFiberprometheus tests that go.mod includes fiberprometheus (9.2-AC: 7)
func TestAdvancedGoModIncludesFiberprometheus(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-gomod-obs"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	goModFile := filepath.Join(projectName, "go.mod")
	content, err := os.ReadFile(goModFile)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "fiberprometheus") {
		t.Errorf("go.mod should contain 'fiberprometheus' for advanced observability, got:\n%s", contentStr)
	}
}

// TestNoneGoModExcludesFiberprometheus tests that go.mod does NOT include fiberprometheus (9.2-AC: no-regression)
func TestNoneGoModExcludesFiberprometheus(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-gomod-none"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	goModFile := filepath.Join(projectName, "go.mod")
	content, err := os.ReadFile(goModFile)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "fiberprometheus") {
		t.Errorf("go.mod should NOT contain 'fiberprometheus' for none observability, got:\n%s", contentStr)
	}
}

// TestAdvancedMainGoIncludesMetrics tests that main.go includes metrics wiring (9.1-AC: 2)
func TestAdvancedMainGoIncludesMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-main-obs"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	mainGoFile := filepath.Join(projectName, "cmd", "main.go")
	content, err := os.ReadFile(mainGoFile)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "pkg/metrics") {
		t.Errorf("main.go should import 'pkg/metrics' for advanced observability, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "NewPrometheusMetrics") {
		t.Errorf("main.go should contain 'NewPrometheusMetrics' fx.Provide for advanced observability, got:\n%s", contentStr)
	}
}

// TestAdvancedServerGoIncludesMetrics tests that server.go includes metrics wiring (9.1-AC: 2, 3)
func TestAdvancedServerGoIncludesMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-server-obs"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	serverFile := filepath.Join(projectName, "internal", "infrastructure", "server", "server.go")
	content, err := os.ReadFile(serverFile)
	if err != nil {
		t.Fatalf("Failed to read server.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "pkg/metrics") {
		t.Errorf("server.go should import 'pkg/metrics' for advanced observability, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "prom.RegisterAt") {
		t.Errorf("server.go should contain 'prom.RegisterAt' for /metrics endpoint, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "prom.Middleware()") {
		t.Errorf("server.go should contain 'prom.Middleware()' for HTTP metrics collection, got:\n%s", contentStr)
	}
}

// TestAdvancedRejectedWithMinimalTemplate tests that --observability=advanced + --template=minimal is rejected
func TestAdvancedRejectedWithMinimalTemplate(t *testing.T) {
	cmd := exec.Command(binaryPath, "--observability=advanced", "--template=minimal", "test-obs-minimal-reject")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for --observability=advanced with --template=minimal")
		os.RemoveAll("test-obs-minimal-reject")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "only supported with --template=full") {
		t.Errorf("Expected error about template restriction, got: %s", outputStr)
	}
}

// TestAdvancedRejectedWithGraphQLTemplate tests that --observability=advanced + --template=graphql is rejected
func TestAdvancedRejectedWithGraphQLTemplate(t *testing.T) {
	cmd := exec.Command(binaryPath, "--observability=advanced", "--template=graphql", "test-obs-graphql-reject")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for --observability=advanced with --template=graphql")
		os.RemoveAll("test-obs-graphql-reject")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "only supported with --template=full") {
		t.Errorf("Expected error about template restriction, got: %s", outputStr)
	}
}

// TestTracerGoContainsTracerProvider tests that tracer.go is generated with NewTracerProvider (AC: 1, 2, 4)
func TestTracerGoContainsTracerProvider(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-tracer-content"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	tracerFile := filepath.Join(projectName, "pkg", "tracing", "tracer.go")
	content, err := os.ReadFile(tracerFile)
	if err != nil {
		t.Fatalf("Failed to read tracer.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "NewTracerProvider") {
		t.Errorf("tracer.go should contain 'NewTracerProvider', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "otlptracegrpc") {
		t.Errorf("tracer.go should contain 'otlptracegrpc', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "OTEL_EXPORTER_OTLP_ENDPOINT") {
		t.Errorf("tracer.go should contain 'OTEL_EXPORTER_OTLP_ENDPOINT', got: %s", contentStr)
	}
}

// TestTracingMiddlewareContainsW3CPropagation tests that tracing_middleware.go contains W3C propagation (AC: 4)
func TestTracingMiddlewareContainsW3CPropagation(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-tracing-mw"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	middlewareFile := filepath.Join(projectName, "internal", "adapters", "middleware", "tracing_middleware.go")
	content, err := os.ReadFile(middlewareFile)
	if err != nil {
		t.Fatalf("Failed to read tracing_middleware.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "TracingMiddleware") {
		t.Errorf("tracing_middleware.go should contain 'TracingMiddleware', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "propagation") {
		t.Errorf("tracing_middleware.go should contain W3C propagation code, got: %s", contentStr)
	}
}

// TestDockerComposeContainsJaeger tests that docker-compose.yml with advanced observability includes Jaeger (AC: 6)
func TestDockerComposeContainsJaeger(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-docker-jaeger"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	dockerFile := filepath.Join(projectName, "docker-compose.yml")
	content, err := os.ReadFile(dockerFile)
	if err != nil {
		t.Fatalf("Failed to read docker-compose.yml: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "jaeger") {
		t.Errorf("docker-compose.yml should contain 'jaeger' service for advanced observability, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "jaegertracing/all-in-one") {
		t.Errorf("docker-compose.yml should contain Jaeger image, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "16686") {
		t.Errorf("docker-compose.yml should expose Jaeger UI port 16686, got: %s", contentStr)
	}
}

// TestEnvExampleContainsOTELEndpoint tests that .env.example contains OTEL vars (9.2-AC: 6)
func TestEnvExampleContainsOTELEndpoint(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-env-otel"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	envFile := filepath.Join(projectName, ".env.example")
	content, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("Failed to read .env.example: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "OTEL_EXPORTER_OTLP_ENDPOINT") {
		t.Errorf(".env.example should contain 'OTEL_EXPORTER_OTLP_ENDPOINT', got: %s", contentStr)
	}
}

// TestDBTracingGoContainsCallbacks tests that db_tracing.go is generated with GORM callbacks (AC: 5)
func TestDBTracingGoContainsCallbacks(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-db-tracing"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	dbTracingFile := filepath.Join(projectName, "pkg", "tracing", "db_tracing.go")
	content, err := os.ReadFile(dbTracingFile)
	if err != nil {
		t.Fatalf("Failed to read db_tracing.go: %v", err)
	}

	contentStr := string(content)
	// Verify RegisterGORMCallbacks function exists
	if !strings.Contains(contentStr, "RegisterGORMCallbacks") {
		t.Errorf("db_tracing.go should contain 'RegisterGORMCallbacks', got: %s", contentStr)
	}
	// Verify db.system attribute matches default database (postgresql)
	if !strings.Contains(contentStr, "postgresql") {
		t.Errorf("db_tracing.go should contain db.system 'postgresql' for default database, got: %s", contentStr)
	}
	// Verify all CRUD callbacks are registered
	for _, op := range []string{"otel:before_query", "otel:before_create", "otel:before_update", "otel:before_delete"} {
		if !strings.Contains(contentStr, op) {
			t.Errorf("db_tracing.go should contain callback %q, got: %s", op, contentStr)
		}
	}
}

// TestDBTracingGoMySQLSystem tests that db_tracing.go uses db.system=mysql for MySQL projects (AC: 5)
func TestDBTracingGoMySQLSystem(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-db-tracing-mysql"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, "mysql", ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	dbTracingFile := filepath.Join(projectName, "pkg", "tracing", "db_tracing.go")
	content, err := os.ReadFile(dbTracingFile)
	if err != nil {
		t.Fatalf("Failed to read db_tracing.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, `"mysql"`) {
		t.Errorf("db_tracing.go should contain db.system 'mysql' for MySQL database, got: %s", contentStr)
	}
}

// TestGoModContainsOpenTelemetry tests that go.mod includes OpenTelemetry packages (AC: 7)
func TestGoModContainsOpenTelemetry(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-gomod-otel"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	goModFile := filepath.Join(projectName, "go.mod")
	content, err := os.ReadFile(goModFile)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "go.opentelemetry.io/otel") {
		t.Errorf("go.mod should contain 'go.opentelemetry.io/otel', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "otlptracegrpc") {
		t.Errorf("go.mod should contain 'otlptracegrpc', got: %s", contentStr)
	}
}

// TestHealthHandlerGeneratedForAllProjects tests that health_handler.go is generated for non-advanced mode (9.3-AC: 1)
func TestHealthHandlerGeneratedForAllProjects(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-health-none"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	healthFile := filepath.Join(projectName, "internal", "adapters", "handlers", "health_handler.go")
	if _, err := os.Stat(healthFile); os.IsNotExist(err) {
		t.Errorf("health_handler.go should be generated for all projects, but it doesn't exist at %s", healthFile)
	}
}

// TestHealthHandlerContainsLivenessAndReadiness tests that health_handler.go has required functions (9.3-AC: 1, 2)
func TestHealthHandlerContainsLivenessAndReadiness(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-health-funcs"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	healthFile := filepath.Join(projectName, "internal", "adapters", "handlers", "health_handler.go")
	content, err := os.ReadFile(healthFile)
	if err != nil {
		t.Fatalf("Failed to read health_handler.go: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "func (h *HealthHandler) Liveness") {
		t.Errorf("health_handler.go should contain 'Liveness' method, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "func (h *HealthHandler) Readiness") {
		t.Errorf("health_handler.go should contain 'Readiness' method, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "NewHealthHandler") {
		t.Errorf("health_handler.go should contain 'NewHealthHandler' constructor, got: %s", contentStr)
	}
}

// TestHealthResponseHasJSONSnakeCaseFields tests HealthResponse has snake_case JSON fields (9.3-AC: 2)
func TestHealthResponseHasJSONSnakeCaseFields(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-health-response"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	healthFile := filepath.Join(projectName, "internal", "adapters", "handlers", "health_handler.go")
	content, err := os.ReadFile(healthFile)
	if err != nil {
		t.Fatalf("Failed to read health_handler.go: %v", err)
	}

	contentStr := string(content)

	// Verify JSON snake_case field tags (AC: HealthResponse struct)
	for _, field := range []string{`json:"status"`, `json:"service"`, `json:"timestamp"`, `json:"checks,omitempty"`, `json:"error,omitempty"`} {
		if !strings.Contains(contentStr, field) {
			t.Errorf("health_handler.go HealthResponse should have %q JSON tag, got: %s", field, contentStr)
		}
	}
}

// TestHealthRoutesInServerRoutes tests that /health, /health/liveness, /health/readiness are registered (9.3-AC: 1, 5)
func TestHealthRoutesInServerRoutes(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-health-routes"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	routesFile := filepath.Join(projectName, "internal", "adapters", "http", "routes.go")
	content, err := os.ReadFile(routesFile)
	if err != nil {
		t.Fatalf("Failed to read routes.go: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, `"/health/liveness"`) {
		t.Errorf("routes.go should contain '/health/liveness' route, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, `"/health/readiness"`) {
		t.Errorf("routes.go should contain '/health/readiness' route, got: %s", contentStr)
	}
	// Backward compatibility: /health should still be registered (AC5)
	if !strings.Contains(contentStr, `"/health"`) {
		t.Errorf("routes.go should contain '/health' backward-compat route, got: %s", contentStr)
	}
}

// TestKubernetesProbesYamlGenerated tests that probes.yaml is generated for all full projects (9.3-AC: 6)
func TestKubernetesProbesYamlGenerated(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-probes-yaml"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	probesFile := filepath.Join(projectName, "deployments", "kubernetes", "probes.yaml")
	if _, err := os.Stat(probesFile); os.IsNotExist(err) {
		t.Errorf("deployments/kubernetes/probes.yaml should be generated, but it doesn't exist")
	}
}

// TestKubernetesProbesYamlContainsProbes tests probes.yaml has liveness/readiness/startup probes (9.3-AC: 6)
func TestKubernetesProbesYamlContainsProbes(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-probes-content"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityNone, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	probesFile := filepath.Join(projectName, "deployments", "kubernetes", "probes.yaml")
	content, err := os.ReadFile(probesFile)
	if err != nil {
		t.Fatalf("Failed to read probes.yaml: %v", err)
	}

	contentStr := string(content)

	for _, probe := range []string{"livenessProbe", "readinessProbe", "startupProbe"} {
		if !strings.Contains(contentStr, probe) {
			t.Errorf("probes.yaml should contain '%s', got: %s", probe, contentStr)
		}
	}
	if !strings.Contains(contentStr, "/health/liveness") {
		t.Errorf("probes.yaml should reference '/health/liveness' path, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "/health/readiness") {
		t.Errorf("probes.yaml should reference '/health/readiness' path, got: %s", contentStr)
	}
}

// TestAdvancedHealthHandlerHasPrometheusMetrics tests health_handler.go with metrics for advanced mode (9.3-AC: 4)
func TestAdvancedHealthHandlerHasPrometheusMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-health-metrics"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	healthFile := filepath.Join(projectName, "internal", "adapters", "handlers", "health_handler.go")
	content, err := os.ReadFile(healthFile)
	if err != nil {
		t.Fatalf("Failed to read health_handler.go: %v", err)
	}

	contentStr := string(content)

	// Advanced version should import and use metrics
	if !strings.Contains(contentStr, "pkg/metrics") {
		t.Errorf("health_handler.go (advanced) should import 'pkg/metrics', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "HealthCheckDuration") {
		t.Errorf("health_handler.go (advanced) should reference 'HealthCheckDuration', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "HealthCheckStatus") {
		t.Errorf("health_handler.go (advanced) should reference 'HealthCheckStatus', got: %s", contentStr)
	}
}

// TestLoggerContainsWithTraceContext tests that logger.go contains WithTraceContext when observability=advanced (AC: 3)
func TestLoggerContainsWithTraceContext(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-logger-trace"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	err = run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	loggerFile := filepath.Join(projectName, "pkg", "logger", "logger.go")
	content, err := os.ReadFile(loggerFile)
	if err != nil {
		t.Fatalf("Failed to read logger.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "WithTraceContext") {
		t.Errorf("logger.go should contain 'WithTraceContext' for advanced observability, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "trace_id") {
		t.Errorf("logger.go should contain 'trace_id' field, got: %s", contentStr)
	}
}

// ── Story 9.4 Tests ─────────────────────────────────────────────────────────

// TestPrometheusYMLGeneratedForAdvanced tests that deployments/prometheus/prometheus.yml
// is generated when --observability=advanced (9.4-AC: 5)
func TestPrometheusYMLGeneratedForAdvanced(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-prom-yml"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	promFile := filepath.Join(projectName, "deployments", "prometheus", "prometheus.yml")
	if _, err := os.Stat(promFile); os.IsNotExist(err) {
		t.Errorf("deployments/prometheus/prometheus.yml should exist for advanced observability, but it doesn't")
	}
}

// TestPrometheusYMLContainsProjectJobName tests that prometheus.yml contains the
// project name in the job_name (9.4-AC: 5)
func TestPrometheusYMLContainsProjectJobName(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "my-api-project"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	promFile := filepath.Join(projectName, "deployments", "prometheus", "prometheus.yml")
	content, err := os.ReadFile(promFile)
	if err != nil {
		t.Fatalf("Failed to read prometheus.yml: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, projectName) {
		t.Errorf("prometheus.yml should contain project name %q as job_name, got:\n%s", projectName, contentStr)
	}
	if !strings.Contains(contentStr, "job_name") {
		t.Errorf("prometheus.yml should contain 'job_name', got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "/metrics") {
		t.Errorf("prometheus.yml should contain '/metrics' metrics_path, got:\n%s", contentStr)
	}
}

// TestGrafanaDashboardIsValidJSON tests that api-dashboard.json is a valid importable
// Grafana JSON (9.4-AC: 1)
func TestGrafanaDashboardIsValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-grafana-json"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	dashFile := filepath.Join(projectName, "deployments", "grafana", "dashboards", "api-dashboard.json")
	content, err := os.ReadFile(dashFile)
	if err != nil {
		t.Fatalf("Failed to read api-dashboard.json: %v", err)
	}

	var dashboard map[string]interface{}
	if err := json.Unmarshal(content, &dashboard); err != nil {
		t.Errorf("api-dashboard.json is not valid JSON: %v\nContent:\n%s", err, string(content))
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, projectName) {
		t.Errorf("api-dashboard.json should contain project name %q in title/uid, got:\n%s", projectName, contentStr)
	}
	if !strings.Contains(contentStr, "panels") {
		t.Errorf("api-dashboard.json should contain 'panels' array, got:\n%s", contentStr)
	}
}

// TestDockerComposeContainsFullObservabilityStack tests that docker-compose.yml with
// --observability=advanced includes prometheus, grafana, and jaeger services (9.4-AC: 3)
func TestDockerComposeContainsFullObservabilityStack(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-dc-full-stack"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	dockerFile := filepath.Join(projectName, "docker-compose.yml")
	content, err := os.ReadFile(dockerFile)
	if err != nil {
		t.Fatalf("Failed to read docker-compose.yml: %v", err)
	}

	contentStr := string(content)

	for _, service := range []string{"prometheus", "grafana", "jaeger"} {
		if !strings.Contains(contentStr, service+":") {
			t.Errorf("docker-compose.yml should contain '%s:' service, got:\n%s", service, contentStr)
		}
	}
	if !strings.Contains(contentStr, "prom/prometheus") {
		t.Errorf("docker-compose.yml should contain Prometheus image, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "grafana/grafana") {
		t.Errorf("docker-compose.yml should contain Grafana image, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "9090") {
		t.Errorf("docker-compose.yml should expose Prometheus port 9090, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "3000") {
		t.Errorf("docker-compose.yml should expose Grafana port 3000, got:\n%s", contentStr)
	}
}

// TestGrafanaFilesNotGeneratedWithoutAdvanced tests that Grafana files are NOT
// generated when observability is none or basic (9.4-AC: 8 — no-regression)
func TestGrafanaFilesNotGeneratedWithoutAdvanced(t *testing.T) {
	for _, level := range []string{ObservabilityNone, ObservabilityBasic} {
		t.Run("observability="+level, func(t *testing.T) {
			tmpDir := t.TempDir()
			projectName := "test-no-grafana-" + level

			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get working directory: %v", err)
			}
			defer os.Chdir(originalWd)
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}

			if err := run(projectName, DefaultTemplate, DefaultDatabase, level, DefaultFramework); err != nil {
				t.Fatalf("run() returned error: %v", err)
			}

			grafanaFiles := []string{
				filepath.Join(projectName, "deployments", "grafana", "dashboards", "api-dashboard.json"),
				filepath.Join(projectName, "deployments", "prometheus", "prometheus.yml"),
				filepath.Join(projectName, "deployments", "prometheus", "rules", "api_alerts.yml"),
			}

			for _, f := range grafanaFiles {
				if _, err := os.Stat(f); !os.IsNotExist(err) {
					t.Errorf("Grafana/Prometheus file %s should NOT exist for observability=%s", f, level)
				}
			}
		})
	}
}

// TestPrometheusAlertRulesGenerated tests that alert rules file is created (9.4-AC: 6)
func TestPrometheusAlertRulesGenerated(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-alert-rules"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	rulesFile := filepath.Join(projectName, "deployments", "prometheus", "rules", "api_alerts.yml")
	content, err := os.ReadFile(rulesFile)
	if err != nil {
		t.Fatalf("Failed to read api_alerts.yml: %v", err)
	}

	contentStr := string(content)
	for _, alert := range []string{"HighErrorRate", "HighP95Latency", "DatabaseDown"} {
		if !strings.Contains(contentStr, alert) {
			t.Errorf("api_alerts.yml should contain alert '%s', got:\n%s", alert, contentStr)
		}
	}
}

// TestGrafanaProvisioningFilesGenerated tests that datasource and dashboard
// provisioning files are created for auto-provisioning (9.4-AC: 4)
func TestGrafanaProvisioningFilesGenerated(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-grafana-prov"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	provisioningFiles := []string{
		filepath.Join(projectName, "deployments", "grafana", "provisioning", "datasources", "prometheus.yaml"),
		filepath.Join(projectName, "deployments", "grafana", "provisioning", "dashboards", "default.yaml"),
	}

	for _, f := range provisioningFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("Grafana provisioning file %s should exist for advanced observability", f)
		}
	}
}

// TestEnvExampleContainsGrafanaVars tests that .env.example contains GF_ vars (9.4-AC)
func TestEnvExampleContainsGrafanaVars(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-env-grafana"

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := run(projectName, DefaultTemplate, DefaultDatabase, ObservabilityAdvanced, DefaultFramework); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	envFile := filepath.Join(projectName, ".env.example")
	content, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("Failed to read .env.example: %v", err)
	}

	contentStr := string(content)
	for _, gfVar := range []string{"GF_SECURITY_ADMIN_USER", "GF_SECURITY_ADMIN_PASSWORD"} {
		if !strings.Contains(contentStr, gfVar) {
			t.Errorf(".env.example should contain '%s' for advanced observability, got:\n%s", gfVar, contentStr)
		}
	}
}
