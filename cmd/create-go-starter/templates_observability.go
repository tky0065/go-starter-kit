package main

import "strings"

// ObservabilityTemplates holds all templates for observability file generation.
// These templates are generated only when --observability=advanced is specified.
type ObservabilityTemplates struct {
	projectName string
}

// NewObservabilityTemplates creates a new ObservabilityTemplates instance.
func NewObservabilityTemplates(projectName string) *ObservabilityTemplates {
	return &ObservabilityTemplates{projectName: projectName}
}

// PrometheusTemplate returns the template for pkg/metrics/prometheus.go.
// It sets up a fiberprometheus registry with HTTP metrics middleware and
// health check metrics (HealthCheckDuration, HealthCheckStatus).
func (t *ObservabilityTemplates) PrometheusTemplate() string {
	return `package metrics

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics encapsulates Prometheus configuration for Fiber.
// It includes HTTP request metrics via fiberprometheus and custom health check metrics.
type PrometheusMetrics struct {
	prometheus          *fiberprometheus.FiberPrometheus
	HealthCheckDuration *prometheus.HistogramVec
	HealthCheckStatus   *prometheus.GaugeVec
}

// NewPrometheusMetrics creates a new PrometheusMetrics instance.
// The appName is used as a label prefix for all HTTP metrics.
// Health check metrics are registered globally in the default Prometheus registry.
func NewPrometheusMetrics(appName string) *PrometheusMetrics {
	prom := fiberprometheus.New(appName)

	healthDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "health_check_duration_seconds",
		Help:    "Duration of health check in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"check"})

	healthStatus := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "health_check_status",
		Help: "Health check status: 1 if healthy, 0 if unhealthy",
	}, []string{"check"})

	prometheus.MustRegister(healthDuration, healthStatus)

	return &PrometheusMetrics{
		prometheus:          prom,
		HealthCheckDuration: healthDuration,
		HealthCheckStatus:   healthStatus,
	}
}

// RegisterAt registers the /metrics endpoint on the Fiber app.
func (p *PrometheusMetrics) RegisterAt(app *fiber.App, path string) {
	p.prometheus.RegisterAt(app, path)
}

// Middleware returns the Fiber handler that captures HTTP metrics.
// Tracks: request count (by method, path, status), latency histogram.
func (p *PrometheusMetrics) Middleware() fiber.Handler {
	return p.prometheus.Middleware
}
`
}

// MetricsMiddlewareTemplate returns the template for
// internal/adapters/middleware/metrics_middleware.go.
// This file re-exports the Middleware helper for use in fx modules.
func (t *ObservabilityTemplates) MetricsMiddlewareTemplate() string {
	return `package middleware

import (
	"` + t.projectName + `/pkg/metrics"

	"github.com/gofiber/fiber/v2"
)

// MetricsMiddleware wraps PrometheusMetrics to expose the Fiber middleware.
type MetricsMiddleware struct {
	prom *metrics.PrometheusMetrics
}

// NewMetricsMiddleware creates a MetricsMiddleware instance.
func NewMetricsMiddleware(prom *metrics.PrometheusMetrics) *MetricsMiddleware {
	return &MetricsMiddleware{prom: prom}
}

// Middleware returns the Fiber handler for HTTP metrics collection.
func (m *MetricsMiddleware) Middleware() fiber.Handler {
	return m.prom.Middleware()
}
`
}

// MetricsHandlerTemplate returns the template for
// internal/adapters/handlers/metrics_handler.go.
// This handler exposes the Prometheus /metrics endpoint.
func (t *ObservabilityTemplates) MetricsHandlerTemplate() string {
	return `package handlers

import (
	"` + t.projectName + `/pkg/metrics"

	"github.com/gofiber/fiber/v2"
)

// MetricsHandler handles the /metrics endpoint for Prometheus scraping.
type MetricsHandler struct {
	prom *metrics.PrometheusMetrics
}

// NewMetricsHandler creates a MetricsHandler instance.
func NewMetricsHandler(prom *metrics.PrometheusMetrics) *MetricsHandler {
	return &MetricsHandler{prom: prom}
}

// Register registers the /metrics route on the provided Fiber app.
// @Summary     Prometheus metrics
// @Description Returns Prometheus-formatted metrics for scraping
// @Tags        observability
// @Produce     plain
// @Router      /metrics [get]
func (h *MetricsHandler) Register(app *fiber.App) {
	h.prom.RegisterAt(app, "/metrics")
}
`
}

// ObservabilityGoModExtension returns the additional go.mod requires
// needed when observability=advanced is enabled.
func ObservabilityGoModExtension() string {
	return "\tgithub.com/ansrivas/fiberprometheus/v2 v2.7.0\n"
}

// GoModTemplateWithObservability returns a go.mod that includes fiberprometheus and OpenTelemetry dependencies.
// This is used instead of GoModTemplate() when observability=advanced.
func (t *ObservabilityTemplates) GoModTemplateWithObservability(database string) string {
	databaseDriver := GoModDependencies(database)

	return `module ` + t.projectName + `

go 1.25.5

require (
	github.com/ansrivas/fiberprometheus/v2 v2.7.0
	github.com/go-playground/validator/v10 v10.30.1
	github.com/gofiber/contrib/jwt v1.1.2
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/joho/godotenv v1.5.1
	github.com/prometheus/client_golang v1.22.0
	github.com/rs/zerolog v1.33.0
	github.com/swaggo/fiber-swagger v1.3.0
	github.com/swaggo/swag v1.16.4
	go.opentelemetry.io/otel v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.30.0
	go.opentelemetry.io/otel/sdk v1.30.0
	go.opentelemetry.io/otel/trace v1.30.0
	go.uber.org/fx v1.24.0
	golang.org/x/crypto v0.32.0
	google.golang.org/grpc v1.68.0
	` + databaseDriver + `
	gorm.io/gorm v1.31.1
)
`
}

// TracerTemplate returns the template for pkg/tracing/tracer.go.
// It initialises the OpenTelemetry TracerProvider with OTLP/gRPC exporter and registers
// shutdown via the fx lifecycle hook.
func (t *ObservabilityTemplates) TracerTemplate() string {
	return `// Package tracing provides OpenTelemetry distributed tracing setup.
// It initialises the TracerProvider with an OTLP/gRPC exporter pointing to Jaeger (or any
// OTLP-compatible backend) and registers graceful shutdown through the fx lifecycle.
package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/fx"
)

// Params holds the fx dependencies for TracerProvider initialisation.
type Params struct {
	fx.In
	LC fx.Lifecycle
}

// NewTracerProvider initialises and registers the global OpenTelemetry TracerProvider.
// It reads configuration from environment variables:
//   - OTEL_EXPORTER_OTLP_ENDPOINT: gRPC endpoint (default: localhost:4317)
//   - APP_NAME: service name used as the OTel resource attribute (default: ` + t.projectName + `)
//
// The provider exports spans via OTLP/gRPC (insecure) to the configured endpoint.
// Shutdown is registered on the fx lifecycle OnStop hook to flush pending spans.
func NewTracerProvider(params Params) *sdktrace.TracerProvider {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}
	serviceName := os.Getenv("APP_NAME")
	if serviceName == "" {
		serviceName = "` + t.projectName + `"
	}

	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		// Log the error and fall back to a no-op provider rather than crashing.
		// This allows the application to start even if the tracing backend is unavailable.
		otel.SetTracerProvider(sdktrace.NewTracerProvider())
		return sdktrace.NewTracerProvider()
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
		// NOTE: AlwaysSample is suitable for development.
		// For production, consider: sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Register as the global TracerProvider so otel.Tracer() calls use it.
	otel.SetTracerProvider(tp)
	// Register W3C TraceContext + Baggage propagators globally.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	params.LC.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	})

	return tp
}
`
}

// TracingMiddlewareTemplate returns the template for
// internal/adapters/middleware/tracing_middleware.go.
// The middleware extracts W3C traceparent headers from incoming requests or creates a new
// root span, propagates the context through the Fiber handler chain, and records HTTP
// attributes on each span.
func (t *ObservabilityTemplates) TracingMiddlewareTemplate() string {
	return `package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// TracingMiddleware returns a Fiber handler that instruments each request with an
// OpenTelemetry span.  It extracts an existing W3C traceparent from incoming headers
// (parent-based propagation) or starts a new root span when none is present.
// The resulting context is stored via c.SetUserContext so downstream handlers and
// services can create child spans or enrich logs with trace/span IDs.
func TracingMiddleware() fiber.Handler {
	tracer := otel.Tracer("fiber")
	// IMPORTANT: The global propagator is set during TracerProvider initialisation
	// (tracing.NewTracerProvider).  This middleware MUST be registered after the
	// TracerProvider is initialised, otherwise GetTextMapPropagator returns a no-op.
	// In the fx dependency graph, the server depends on *sdktrace.TracerProvider,
	// which guarantees the correct initialisation order.
	propagator := otel.GetTextMapPropagator()

	return func(c *fiber.Ctx) error {
		// Extract trace context from incoming HTTP headers (W3C traceparent / baggage).
		ctx := propagator.Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))

		// Start a new span for this request.
		spanName := c.Method() + " " + c.Path()
		ctx, span := tracer.Start(ctx, spanName)
		defer span.End()

		// Propagate the context so handlers can access the active span.
		c.SetUserContext(ctx)

		// Annotate the span with standard HTTP attributes.
		span.SetAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.url", c.OriginalURL()),
			attribute.String("http.scheme", c.Protocol()),
			attribute.String("net.host.name", c.Hostname()),
		)

		err := c.Next()

		// Record the HTTP status code after the handler chain completes.
		span.SetAttributes(attribute.Int("http.status_code", c.Response().StatusCode()))
		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}
`
}

// GORMTracingTemplate returns the template for pkg/tracing/db_tracing.go.
// It registers OpenTelemetry callback hooks on a GORM *gorm.DB instance to create child
// spans for every database operation (query, create, update, delete).  The function is
// invoked via fx.Invoke so it runs automatically after the database is initialised.
// The database parameter determines the correct db.system attribute value.
func (t *ObservabilityTemplates) GORMTracingTemplate(database string) string {
	// Map database type to OTel semantic convention db.system value
	dbSystem := "postgresql"
	switch database {
	case "mysql":
		dbSystem = "mysql"
	case "sqlite":
		dbSystem = "sqlite"
	}

	return `package tracing

import (
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// RegisterGORMCallbacks attaches OpenTelemetry tracing callbacks to the provided GORM DB.
// It is called automatically via fx.Invoke after the database module is initialised.
// Each operation creates a child span with db.system, db.operation, and db.sql.table attributes.
// SQL parameter values are intentionally omitted to avoid leaking sensitive data.
func RegisterGORMCallbacks(db *gorm.DB) {
	tracer := otel.Tracer("gorm")
	dbSystemAttr := attribute.String("db.system", "` + dbSystem + `")

	// --- Query callbacks (SELECT) ---
	db.Callback().Query().Before("gorm:query").Register("otel:before_query",
		func(db *gorm.DB) {
			if db.Statement == nil || db.Statement.Context == nil {
				return
			}
			ctx, _ := tracer.Start(db.Statement.Context, "db.query",
				trace.WithAttributes(
					dbSystemAttr,
					attribute.String("db.operation", "query"),
					attribute.String("db.sql.table", db.Statement.Table),
				),
			)
			db.Statement.Context = ctx
		},
	)
	db.Callback().Query().After("gorm:query").Register("otel:after_query",
		afterCallback,
	)

	// --- Create callbacks (INSERT) ---
	db.Callback().Create().Before("gorm:create").Register("otel:before_create",
		func(db *gorm.DB) {
			if db.Statement == nil || db.Statement.Context == nil {
				return
			}
			ctx, _ := tracer.Start(db.Statement.Context, "db.create",
				trace.WithAttributes(
					dbSystemAttr,
					attribute.String("db.operation", "insert"),
					attribute.String("db.sql.table", db.Statement.Table),
				),
			)
			db.Statement.Context = ctx
		},
	)
	db.Callback().Create().After("gorm:create").Register("otel:after_create",
		afterCallback,
	)

	// --- Update callbacks (UPDATE) ---
	db.Callback().Update().Before("gorm:update").Register("otel:before_update",
		func(db *gorm.DB) {
			if db.Statement == nil || db.Statement.Context == nil {
				return
			}
			ctx, _ := tracer.Start(db.Statement.Context, "db.update",
				trace.WithAttributes(
					dbSystemAttr,
					attribute.String("db.operation", "update"),
					attribute.String("db.sql.table", db.Statement.Table),
				),
			)
			db.Statement.Context = ctx
		},
	)
	db.Callback().Update().After("gorm:update").Register("otel:after_update",
		afterCallback,
	)

	// --- Delete callbacks (DELETE) ---
	db.Callback().Delete().Before("gorm:delete").Register("otel:before_delete",
		func(db *gorm.DB) {
			if db.Statement == nil || db.Statement.Context == nil {
				return
			}
			ctx, _ := tracer.Start(db.Statement.Context, "db.delete",
				trace.WithAttributes(
					dbSystemAttr,
					attribute.String("db.operation", "delete"),
					attribute.String("db.sql.table", db.Statement.Table),
				),
			)
			db.Statement.Context = ctx
		},
	)
	db.Callback().Delete().After("gorm:delete").Register("otel:after_delete",
		afterCallback,
	)
}

// afterCallback ends the span and records any error (except ErrRecordNotFound which is normal).
func afterCallback(db *gorm.DB) {
	if db.Statement == nil || db.Statement.Context == nil {
		return
	}
	span := trace.SpanFromContext(db.Statement.Context)
	if !span.IsRecording() {
		return
	}
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		span.RecordError(db.Error)
		span.SetStatus(codes.Error, db.Error.Error())
	}
	span.End()
}
`
}

// LoggerWithTracingTemplate returns pkg/logger/logger.go extended with the WithTraceContext
// helper.  This replaces the standard LoggerTemplate() when observability=advanced so that
// handlers can enrich structured log entries with trace_id and span_id for log-trace
// correlation in Jaeger / Grafana.
func (t *ObservabilityTemplates) LoggerWithTracingTemplate() string {
	return `// Package logger provides structured logging utilities using zerolog.
// It configures the logger based on the application environment, using JSON format
// in production for log aggregation systems and console format in development
// for human readability.  The package also provides WithTraceContext to correlate
// log entries with OpenTelemetry distributed traces.
package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Module provides the logger dependency via fx for application-wide logging.
var Module = fx.Module("logger",
	fx.Provide(NewLogger),
)

// NewLogger creates a new zerolog logger instance configured for the current environment.
// In production (APP_ENV=production), it outputs JSON format for log aggregation.
// In other environments, it uses a human-readable console format with colours.
func NewLogger() zerolog.Logger {
	env := os.Getenv("APP_ENV")

	var logger zerolog.Logger
	if env == "production" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	return logger
}

// WithTraceContext enriches the given logger with trace_id and span_id fields extracted
// from the active OpenTelemetry span in ctx.  If no recording span is present the
// original logger is returned unchanged.
//
// Usage in a Fiber handler:
//
//	log := logger.WithTraceContext(c.UserContext(), h.log)
//	log.Info().Msg("handling request")
func WithTraceContext(ctx context.Context, log zerolog.Logger) zerolog.Logger {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return log
	}
	sc := span.SpanContext()
	return log.With().
		Str("trace_id", sc.TraceID().String()).
		Str("span_id", sc.SpanID().String()).
		Logger()
}
`
}

// EnvTemplateWithObservability returns a .env.example that includes the standard
// application variables plus OpenTelemetry configuration for advanced observability.
// The database parameter controls which DB env vars are included.
func (t *ObservabilityTemplates) EnvTemplateWithObservability(database string) string {
	return `# Application Configuration
APP_NAME=` + t.projectName + `
APP_ENV=development
APP_PORT=8080

` + DatabaseEnvVars(database, t.projectName) + `

# JWT Configuration
# IMPORTANT: Generate a secure random secret for production!
# Example: openssl rand -base64 32
JWT_SECRET=
JWT_EXPIRY=24h

# OpenTelemetry / Distributed Tracing (--observability=advanced)
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true

# Grafana (--observability=advanced)
GF_SECURITY_ADMIN_USER=admin
GF_SECURITY_ADMIN_PASSWORD=admin
GF_USERS_ALLOW_SIGN_UP=false
`
}

// DockerComposeTemplateWithObservability returns a docker-compose.yml with the complete
// observability stack: database, Jaeger (tracing), Prometheus (metrics), and Grafana
// (dashboards + alerting). The database parameter selects between postgres, mysql, sqlite.
func (t *ObservabilityTemplates) DockerComposeTemplateWithObservability(database string) string {
	pn := t.projectName

	jaegerService := `
  # Jaeger — distributed tracing backend (OpenTelemetry OTLP)
  jaeger:
    image: jaegertracing/all-in-one:1.56.0
    container_name: ` + pn + `-jaeger
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    ports:
      - "16686:16686"   # Jaeger UI
      - "4317:4317"     # OTLP gRPC
      - "4318:4318"     # OTLP HTTP
    networks:
      - ` + pn + `_network
`

	prometheusService := `
  # Prometheus — metrics collection & storage
  prometheus:
    image: prom/prometheus:v2.51.0
    container_name: ` + pn + `-prometheus
    volumes:
      - ./deployments/prometheus:/etc/prometheus
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=15d'
      - '--web.enable-lifecycle'
    ports:
      - "9090:9090"
    networks:
      - ` + pn + `_network
`

	grafanaService := `
  # Grafana — metrics visualization & dashboards
  grafana:
    image: grafana/grafana:10.4.0
    container_name: ` + pn + `-grafana
    environment:
      - GF_SECURITY_ADMIN_USER=${GF_SECURITY_ADMIN_USER:-admin}
      - GF_SECURITY_ADMIN_PASSWORD=${GF_SECURITY_ADMIN_PASSWORD:-admin}
      - GF_USERS_ALLOW_SIGN_UP=false
    volumes:
      - ./deployments/grafana/provisioning:/etc/grafana/provisioning
      - ./deployments/grafana/dashboards:/var/lib/grafana/dashboards
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"
    depends_on:
      - prometheus
    networks:
      - ` + pn + `_network
`

	otelEnvVars := `      OTEL_EXPORTER_OTLP_ENDPOINT: jaeger:4317
      OTEL_EXPORTER_OTLP_INSECURE: "true"
`

	// SQLite — no separate database service
	if database == "sqlite" {
		return `services:
  # Application API (SQLite - no separate DB service needed)
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ` + pn + `_api
    environment:
      APP_NAME: ` + pn + `
      APP_ENV: development
      APP_PORT: 8080
      DB_NAME: ` + pn + `.db
      JWT_SECRET: ${JWT_SECRET:-dev-secret-change-in-production}
      JWT_EXPIRY: 24h
      ` + otelEnvVars + `    ports:
      - "8080:8080"
    depends_on:
      jaeger:
        condition: service_started
      prometheus:
        condition: service_started
    networks:
      - ` + pn + `_network
    volumes:
      - .:/app
      - sqlite_data:/app/data
    command: /app/` + pn + `
` + jaegerService + prometheusService + grafanaService + `
volumes:
  sqlite_data:
  prometheus_data:
  grafana_data:

networks:
  ` + pn + `_network:
    driver: bridge
`
	}

	// MySQL / PostgreSQL — inline service definition (avoids volumes: section collision)
	dbComment := "# PostgreSQL Database"
	dbService := `  db:
    image: postgres:16-alpine
    container_name: ` + pn + `_db
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - ` + pn + `_network
`
	dbVolume := "  postgres_data:"
	dbHost := "db"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := pn

	if database == "mysql" {
		dbComment = "# MySQL Database"
		dbService = `  db:
    image: mysql:8.0
    container_name: ` + pn + `_db
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p${DB_PASSWORD}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - ` + pn + `_network
`
		dbVolume = "  mysql_data:"
		dbPort = "3306"
		dbUser = "root"
		dbPassword = "secret"
	}

	return `services:
  ` + dbComment + `
` + dbService + jaegerService + prometheusService + grafanaService + `
  # Application API
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ` + pn + `_api
    environment:
      APP_NAME: ` + pn + `
      APP_ENV: development
      APP_PORT: 8080
      DB_HOST: ` + dbHost + `
      DB_PORT: ` + dbPort + `
      DB_USER: ` + dbUser + `
      DB_PASSWORD: ` + dbPassword + `
      DB_NAME: ` + dbName + `
      DB_SSLMODE: disable
      JWT_SECRET: ${JWT_SECRET:-dev-secret-change-in-production}
      JWT_EXPIRY: 24h
      ` + otelEnvVars + `    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
      jaeger:
        condition: service_started
      prometheus:
        condition: service_started
    networks:
      - ` + pn + `_network
    volumes:
      - .:/app
    command: /app/` + pn + `

volumes:
` + dbVolume + `
  prometheus_data:
  grafana_data:

networks:
  ` + pn + `_network:
    driver: bridge
`
}

// MainGoTemplateWithObservability returns the cmd/main.go with Prometheus metrics and
// OpenTelemetry distributed tracing wired via fx.
// This is used instead of UpdatedMainGoTemplate() when observability=advanced.
func (t *ObservabilityTemplates) MainGoTemplateWithObservability() string {
	return `package main

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/fx"

	"` + t.projectName + `/internal/adapters/handlers"
	"` + t.projectName + `/internal/adapters/repository"
	"` + t.projectName + `/internal/domain/user"
	"` + t.projectName + `/internal/infrastructure/database"
	"` + t.projectName + `/internal/infrastructure/server"
	"` + t.projectName + `/pkg/auth"
	"` + t.projectName + `/pkg/logger"
	"` + t.projectName + `/pkg/metrics"
	"` + t.projectName + `/pkg/tracing"
)

// @title ` + t.projectName + ` API
// @version 1.0
// @description A Go starter kit with authentication, user management, and CRUD operations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables from .env file
	// This is primarily for local development; in production, use system environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or couldn't be loaded")
	}

	fx.New(
		// Core infrastructure
		logger.Module,
		database.Module,

		// Observability — Prometheus metrics
		fx.Provide(func() *metrics.PrometheusMetrics {
			return metrics.NewPrometheusMetrics("` + t.projectName + `")
		}),

		// Observability — Distributed tracing (OpenTelemetry + Jaeger)
		// TracerProvider must be provided before ServerModule so the global tracer is set.
		fx.Provide(tracing.NewTracerProvider),
		fx.Invoke(tracing.RegisterGORMCallbacks),

		// Authentication & authorization
		auth.Module,

		// Domain services
		user.Module,

		// Data persistence
		repository.Module,

		// HTTP handlers
		handlers.Module,

		// HTTP server (must be last as it depends on handlers)
		server.Module,
	).Run()
}
`
}

// ServerTemplateWithObservability returns the server.go with Prometheus metrics and
// OpenTelemetry tracing middleware registered.
// This is used instead of ServerTemplate() when observability=advanced.
func (t *ObservabilityTemplates) ServerTemplateWithObservability() string {
	return `// Package server provides HTTP server configuration and lifecycle management.
// It creates and configures a Fiber application with middleware, error handling,
// and graceful shutdown support through fx lifecycle hooks. This package is part
// of the infrastructure layer and coordinates all HTTP-related concerns.
package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"` + t.projectName + `/pkg/config"
	"` + t.projectName + `/pkg/metrics"
	httpRoutes "` + t.projectName + `/internal/adapters/http"
	"` + t.projectName + `/internal/adapters/middleware"

	// Swagger docs - generated by swag init
	_ "` + t.projectName + `/docs"
)

// Module provides the Fiber server dependency via fx with automatic lifecycle management.
var Module = fx.Module("server",
	fx.Provide(NewServer),
	fx.Invoke(registerHooks),
	fx.Invoke(httpRoutes.RegisterRoutes),
)

// NewServer creates and configures a new Fiber application with centralized error handling.
// It sets up the application name, error handler, Prometheus metrics middleware, distributed
// tracing middleware, and common routes like favicon handling.
// The tp parameter ensures the TracerProvider is initialised before the server starts.
func NewServer(logger zerolog.Logger, db *gorm.DB, prom *metrics.PrometheusMetrics, tp *sdktrace.TracerProvider) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "` + t.projectName + `",
		ErrorHandler: middleware.ErrorHandler,
		// Increase buffer sizes to prevent "Request Header Fields Too Large" errors
		ReadBufferSize:  16384, // 16KB (default is 4KB)
		WriteBufferSize: 16384,
	})

	// Register Prometheus metrics endpoint and middleware (must be first)
	prom.RegisterAt(app, "/metrics")
	app.Use(prom.Middleware())

	// Register distributed tracing middleware (W3C traceparent propagation)
	app.Use(middleware.TracingMiddleware())

	// Ignore common browser requests (favicon, apple-touch-icon)
	// These would otherwise pollute error logs
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Get("/apple-touch-icon*.png", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	logger.Info().Msg("Fiber server initialized with Prometheus metrics and distributed tracing")

	return app
}

// registerHooks registers fx lifecycle hooks for server startup and graceful shutdown.
// It starts the server in a background goroutine on startup and properly shuts it down
// when the application receives a termination signal.
func registerHooks(lifecycle fx.Lifecycle, app *fiber.App, logger zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := config.GetEnv("APP_PORT", "8080")
			logger.Info().Str("port", port).Msg("Starting Fiber server")

			// Start server in background goroutine
			go func() {
				if err := app.Listen(":" + port); err != nil {
					logger.Error().Err(err).Msg("Server stopped unexpectedly")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Shutting down Fiber server gracefully")
			return app.ShutdownWithContext(ctx)
		},
	})
}

`
}

// HealthHandlerWithMetricsTemplate returns internal/adapters/handlers/health_handler.go
// with Prometheus metrics instrumentation for liveness and readiness probes.
// Used when --observability=advanced is specified.
func (t *ObservabilityTemplates) HealthHandlerWithMetricsTemplate() string {
	return `package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"` + t.projectName + `/pkg/metrics"
)

// HealthHandler handles Kubernetes-compatible health check endpoints with Prometheus metrics.
// It instruments readiness checks with health_check_duration_seconds and health_check_status.
type HealthHandler struct {
	db   *gorm.DB
	prom *metrics.PrometheusMetrics
}

// NewHealthHandler creates a new HealthHandler with GORM and Prometheus metrics injection.
// It is registered in the fx dependency injection container via handlers.Module.
func NewHealthHandler(db *gorm.DB, prom *metrics.PrometheusMetrics) *HealthHandler {
	return &HealthHandler{db: db, prom: prom}
}

// HealthResponse represents the health check response body.
// It follows Kubernetes probe conventions with status, service, timestamp, and optional checks.
type HealthResponse struct {
	Status    string            ` + "`json:\"status\"`" + `
	Service   string            ` + "`json:\"service\"`" + `
	Timestamp string            ` + "`json:\"timestamp\"`" + `
	Checks    map[string]string ` + "`json:\"checks,omitempty\"`" + `
	Error     string            ` + "`json:\"error,omitempty\"`" + `
}

// Liveness reports whether the application process is alive.
// Always returns 200 as long as the app is running — used by K8s livenessProbe.
// @Summary     Liveness probe
// @Description Returns 200 if the application is running (K8s livenessProbe)
// @Tags        health
// @Produce     json
// @Success     200 {object} HealthResponse
// @Router      /health/liveness [get]
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(HealthResponse{
		Status:    "alive",
		Service:   "` + t.projectName + `",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness reports whether the application is ready to receive traffic.
// Verifies database connectivity with a 2s timeout and records Prometheus health metrics.
// Returns 503 if any dependency is unavailable.
// @Summary     Readiness probe
// @Description Returns 200 if the application and its dependencies are ready (K8s readinessProbe)
// @Tags        health
// @Produce     json
// @Success     200 {object} HealthResponse
// @Failure     503 {object} HealthResponse
// @Router      /health/readiness [get]
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	start := time.Now()

	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	checks := make(map[string]string)

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		checks["database"] = "error"
		duration := time.Since(start).Seconds()
		if h.prom != nil {
			h.prom.HealthCheckDuration.WithLabelValues("database").Observe(duration)
			h.prom.HealthCheckStatus.WithLabelValues("database").Set(0)
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(HealthResponse{
			Status:    "not_ready",
			Service:   "` + t.projectName + `",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Checks:    checks,
			Error:     "database connection failed",
		})
	}

	checks["database"] = "ok"
	duration := time.Since(start).Seconds()
	if h.prom != nil {
		h.prom.HealthCheckDuration.WithLabelValues("database").Observe(duration)
		h.prom.HealthCheckStatus.WithLabelValues("database").Set(1)
	}

	return c.Status(fiber.StatusOK).JSON(HealthResponse{
		Status:    "ready",
		Service:   "` + t.projectName + `",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}
`
}

// KubernetesProbesTemplate returns the deployments/kubernetes/probes.yaml content.
// This file is generated for all full-template projects as a K8s deployment reference.
func (t *ObservabilityTemplates) KubernetesProbesTemplate() string {
	return `# Kubernetes Probe Configuration Reference for ` + t.projectName + `
# Copy these probe configurations into your Deployment manifest.
# Docs: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/

# livenessProbe: restarts the Pod if the application is deadlocked or unresponsive.
# /health/liveness always returns 200 as long as the process is running.
livenessProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  initialDelaySeconds: 10  # Wait 10s before first check (app startup time)
  periodSeconds: 15         # Check every 15s
  timeoutSeconds: 5         # Fail if no response within 5s
  failureThreshold: 3       # Restart after 3 consecutive failures

# readinessProbe: removes the Pod from the load balancer if not ready to serve traffic.
# /health/readiness checks database connectivity — returns 503 if DB is unreachable.
readinessProbe:
  httpGet:
    path: /health/readiness
    port: 8080
  initialDelaySeconds: 5   # Wait 5s before first check
  periodSeconds: 10         # Check every 10s
  timeoutSeconds: 5         # Fail if no response within 5s
  failureThreshold: 3       # Remove from load balancer after 3 consecutive failures

# startupProbe: gives the app time to start before liveness checks begin.
# Useful for apps with slow startup (migrations, etc.).
startupProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  failureThreshold: 30      # Allow up to 30 * 5s = 150s for startup
  periodSeconds: 5
  timeoutSeconds: 5
`
}

// PrometheusConfigTemplate returns the content for deployments/prometheus/prometheus.yml.
// This configures Prometheus to scrape the application's /metrics endpoint and itself.
// rule_files points to the alert rules directory.
func (t *ObservabilityTemplates) PrometheusConfigTemplate() string {
	return `global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: '` + t.projectName + `'
    static_configs:
      - targets: ['api:8080']
    metrics_path: '/metrics'
    scrape_interval: 10s

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
`
}

// GrafanaDatasourceTemplate returns the content for
// deployments/grafana/provisioning/datasources/prometheus.yaml.
// Grafana auto-provisions this Prometheus datasource on startup.
func (t *ObservabilityTemplates) GrafanaDatasourceTemplate() string {
	return `apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: false
`
}

// GrafanaDashboardProvisionTemplate returns the content for
// deployments/grafana/provisioning/dashboards/default.yaml.
// This tells Grafana to load dashboards from /var/lib/grafana/dashboards automatically.
func (t *ObservabilityTemplates) GrafanaDashboardProvisionTemplate() string {
	return `apiVersion: 1

providers:
  - name: 'default'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 30
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards
`
}

// GrafanaDashboardJSONTemplate returns the content for
// deployments/grafana/dashboards/api-dashboard.json.
// The dashboard contains 7 panels: Request Rate, Error Rate, P95 Latency,
// DB Query P95, Active DB Connections, Health Status DB, and a full-width
// HTTP requests time series.
func (t *ObservabilityTemplates) GrafanaDashboardJSONTemplate() string {
	pn := t.projectName
	tmpl := `{
  "id": null,
  "uid": "__PN__-api",
  "title": "__PN__ - API Dashboard",
  "tags": ["api", "__PN__"],
  "timezone": "browser",
  "schemaVersion": 39,
  "version": 1,
  "refresh": "30s",
  "time": {"from": "now-1h", "to": "now"},
  "panels": [
    {
      "id": 1,
      "type": "timeseries",
      "title": "Request Rate (req/s)",
      "gridPos": {"h": 8, "w": 6, "x": 0, "y": 0},
      "datasource": {"type": "prometheus"},
      "fieldConfig": {
        "defaults": {"unit": "reqps", "color": {"mode": "palette-classic"}}
      },
      "targets": [
        {
          "expr": "sum(rate(http_requests_total{job=\"__PN__\"}[5m]))",
          "refId": "A",
          "legendFormat": "req/s"
        }
      ]
    },
    {
      "id": 2,
      "type": "stat",
      "title": "Error Rate (%)",
      "gridPos": {"h": 8, "w": 6, "x": 6, "y": 0},
      "datasource": {"type": "prometheus"},
      "fieldConfig": {
        "defaults": {
          "unit": "percent",
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "green", "value": null},
              {"color": "yellow", "value": 1},
              {"color": "red", "value": 5}
            ]
          }
        }
      },
      "options": {"reduceOptions": {"calcs": ["lastNotNull"]}},
      "targets": [
        {
          "expr": "sum(rate(http_requests_total{job=\"__PN__\",status_code=~\"5..\"}[5m])) / sum(rate(http_requests_total{job=\"__PN__\"}[5m])) * 100",
          "refId": "A"
        }
      ]
    },
    {
      "id": 3,
      "type": "gauge",
      "title": "P95 Latency",
      "gridPos": {"h": 8, "w": 6, "x": 0, "y": 8},
      "datasource": {"type": "prometheus"},
      "fieldConfig": {
        "defaults": {
          "unit": "s",
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "green", "value": null},
              {"color": "yellow", "value": 0.5},
              {"color": "red", "value": 1}
            ]
          }
        }
      },
      "options": {"reduceOptions": {"calcs": ["lastNotNull"]}},
      "targets": [
        {
          "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{job=\"__PN__\"}[5m])) by (le))",
          "refId": "A",
          "legendFormat": "P95"
        }
      ]
    },
    {
      "id": 4,
      "type": "timeseries",
      "title": "DB Query Duration P95 (s)",
      "gridPos": {"h": 8, "w": 6, "x": 6, "y": 8},
      "datasource": {"type": "prometheus"},
      "fieldConfig": {
        "defaults": {"unit": "s", "color": {"mode": "palette-classic"}}
      },
      "targets": [
        {
          "expr": "histogram_quantile(0.95, sum(rate(db_query_duration_seconds_bucket[5m])) by (le))",
          "refId": "A",
          "legendFormat": "DB P95"
        }
      ]
    },
    {
      "id": 5,
      "type": "stat",
      "title": "Active DB Connections",
      "gridPos": {"h": 8, "w": 4, "x": 0, "y": 16},
      "datasource": {"type": "prometheus"},
      "fieldConfig": {
        "defaults": {"unit": "short", "color": {"mode": "thresholds"},
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "green", "value": null},
              {"color": "yellow", "value": 20},
              {"color": "red", "value": 50}
            ]
          }
        }
      },
      "options": {"reduceOptions": {"calcs": ["lastNotNull"]}},
      "targets": [
        {"expr": "db_connections_open", "refId": "A"}
      ]
    },
    {
      "id": 6,
      "type": "stat",
      "title": "Health Status (Database)",
      "gridPos": {"h": 8, "w": 4, "x": 4, "y": 16},
      "datasource": {"type": "prometheus"},
      "fieldConfig": {
        "defaults": {
          "unit": "short",
          "mappings": [
            {"type": "value", "options": {"0": {"text": "DOWN", "color": "red"}, "1": {"text": "UP", "color": "green"}}}
          ],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "red", "value": null},
              {"color": "green", "value": 1}
            ]
          }
        }
      },
      "options": {"reduceOptions": {"calcs": ["lastNotNull"]}},
      "targets": [
        {"expr": "health_check_status{check=\"database\"}", "refId": "A"}
      ]
    },
    {
      "id": 7,
      "type": "timeseries",
      "title": "HTTP Requests by Status",
      "gridPos": {"h": 8, "w": 12, "x": 0, "y": 24},
      "datasource": {"type": "prometheus"},
      "fieldConfig": {
        "defaults": {"unit": "reqps", "color": {"mode": "palette-classic"}}
      },
      "targets": [
        {
          "expr": "sum by (status_code) (rate(http_requests_total{job=\"__PN__\"}[5m]))",
          "refId": "A",
          "legendFormat": "{{status_code}}"
        }
      ]
    }
  ]
}`
	return strings.ReplaceAll(tmpl, "__PN__", pn)
}

// PrometheusAlertRulesTemplate returns the content for
// deployments/prometheus/rules/api_alerts.yml.
// Defines three pre-configured alerts: HighErrorRate, HighP95Latency, DatabaseDown.
func (t *ObservabilityTemplates) PrometheusAlertRulesTemplate() string {
	return `groups:
  - name: api_alerts
    interval: 30s
    rules:
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{status_code=~"5.."}[5m])) /
          sum(rate(http_requests_total[5m])) > 0.05
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High error rate on {{ $labels.job }}"
          description: "Error rate is {{ $value | humanizePercentage }} (> 5%)"

      - alert: HighP95Latency
        expr: |
          histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1.0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High P95 latency"
          description: "P95 latency is {{ $value | humanizeDuration }} (> 1s)"

      - alert: DatabaseDown
        expr: health_check_status{check="database"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database is down"
          description: "The database health check is failing for > 1 minute"
`
}
