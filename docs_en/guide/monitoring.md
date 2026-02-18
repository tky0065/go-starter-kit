# Monitoring & Observability

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

<i class="material-icons success">new_releases</i> **New in v1.3.0** — Advanced observability is now natively integrated into generated projects via the `--observability` flag.

## Overview

Go Starter Kit offers **3 levels of observability** to monitor your projects in production:

| Level | Flag | Description |
|-------|------|-------------|
| `none` | `--observability=none` (default) | Standard behavior, no instrumentation |
| `basic` | `--observability=basic` | Advanced K8s health checks (liveness/readiness) |
| `advanced` | `--observability=advanced` | Full stack: Prometheus + Jaeger + Grafana + Health Checks |

```bash
# Generate a project with advanced observability
create-go-starter mon-app --template=full --observability=advanced
```

> <i class="material-icons warning">warning</i> **Note**: The `--observability` flag only works with the `full` template. The `minimal` and `graphql` templates are not supported.

---

## Logging with zerolog

### Configuration

The logger is configured in `pkg/logger/logger.go`:

```go
func NewLogger(config *config.Config) zerolog.Logger {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

    var logger zerolog.Logger

    if config.AppEnv == "production" {
        logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
    } else {
        logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
            With().
            Timestamp().
            Logger()
    }

    switch config.AppEnv {
    case "production":
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    case "development":
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    default:
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    }

    return logger
}
```

### Usage

Injection via fx:

```go
type UserService struct {
    logger zerolog.Logger
}

func NewUserService(logger zerolog.Logger) *UserService {
    return &UserService{logger: logger}
}
```

**Structured logging**:

```go
// Info
logger.Info().
    Str("email", user.Email).
    Uint("user_id", user.ID).
    Msg("User registered successfully")

// Error
logger.Error().
    Err(err).
    Str("operation", "create_user").
    Str("email", email).
    Msg("Failed to create user")

// Debug
logger.Debug().
    Interface("request", req).
    Msg("Received request")

// Warn
logger.Warn().
    Dur("duration", elapsed).
    Msg("Slow query detected")
```

### Enriched Logging with Tracing (advanced mode)

When `--observability=advanced` is enabled, the logger is enriched with OpenTelemetry trace identifiers:

```go
// Logs automatically include trace_id and span_id
logger.Info().
    Str("trace_id", span.SpanContext().TraceID().String()).
    Str("span_id", span.SpanContext().SpanID().String()).
    Msg("Processing request")
```

**JSON output example**:
```json
{
  "level": "info",
  "trace_id": "abc123def456...",
  "span_id": "789xyz...",
  "message": "Processing request",
  "timestamp": 1708000000
}
```

This allows correlating logs with traces in Jaeger.

### Log Levels

| Level | Usage |
|-------|-------|
| **Debug** | Detailed information for debugging |
| **Info** | Important events (user login, etc.) |
| **Warn** | Non-critical abnormal behaviors |
| **Error** | Errors requiring attention |
| **Fatal** | Critical errors (app exit) |

### Best practices

<i class="material-icons success">check_circle</i> **GOOD — Structured logging**:

```go
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Dur("duration", elapsed).
    Msg("User logged in")
```

<i class="material-icons error">cancel</i> **BAD — String formatting**:

```go
logger.Info().Msgf("User %s logged in after %v", userID, elapsed)
```

<i class="material-icons success">check_circle</i> **GOOD — No secrets**:

```go
logger.Info().Str("email", email).Msg("User login attempt")
```

<i class="material-icons error">cancel</i> **BAD — Logging secrets**:

```go
logger.Info().Str("password", password).Msg("Login")  // NEVER!
```

---

## Prometheus Metrics

### Endpoint `/metrics`

When `--observability=advanced` is enabled, a Prometheus endpoint is exposed:

```bash
GET /metrics
```

### Exposed HTTP Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total requests by method, route, and status code |
| `http_request_duration_seconds` | Histogram | HTTP latency by route (buckets p50, p90, p95, p99) |
| `http_requests_in_flight` | Gauge | Number of active in-flight requests |

### Library Used

[`fiberprometheus/v2`](https://github.com/ansrivas/fiberprometheus) v2.7.0 — native Fiber v2 integration.

### Generated Files

```
mon-app/
├── pkg/metrics/
│   └── prometheus.go                # Prometheus registry + PrometheusMetrics struct
├── internal/adapters/
│   ├── middleware/
│   │   └── metrics_middleware.go    # HTTP middleware for capturing metrics
│   └── handlers/
│       └── metrics_handler.go       # Handler GET /metrics
```

### Prometheus Configuration Example

```yaml
# prometheus.yml (automatically generated in monitoring/prometheus/)
scrape_configs:
  - job_name: 'mon-app'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Testing Metrics

```bash
# Generate traffic
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users

# View metrics
curl http://localhost:8080/metrics

# Example output
# http_requests_total{method="GET",path="/health",status="200"} 5
# http_request_duration_seconds_bucket{method="GET",path="/health",le="0.1"} 5
# http_requests_in_flight 0
```

---

## Distributed Tracing with OpenTelemetry

### Architecture

```
Application → OTLP/gRPC → Jaeger Collector → Jaeger UI
```

### Configuration

Tracing uses **OpenTelemetry** with OTLP/gRPC export to **Jaeger**:

- **Protocol**: OTLP/gRPC
- **Propagation**: W3C `traceparent` header
- **Endpoint**: Configurable via `OTEL_EXPORTER_OTLP_ENDPOINT` (default: `localhost:4317`)

### Generated Files

```
mon-app/
├── pkg/tracing/
│   └── tracer.go                    # OpenTelemetry configuration + TracerProvider
├── internal/adapters/middleware/
│   └── tracing_middleware.go        # HTTP middleware for creating spans
├── internal/infrastructure/database/
│   └── tracing.go                   # GORM instrumentation with DB spans
└── pkg/logger/
    └── logger_tracing.go            # Logger enriched with trace_id/span_id
```

### Environment Variables

```bash
# .env.example (automatically added)
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=mon-app
```

### Automatically Generated Spans

| Component | Span | Attributes |
|-----------|------|-----------|
| HTTP Middleware | `HTTP {method} {path}` | `http.method`, `http.route`, `http.status_code` |
| GORM Tracing | `db.query` | `db.statement`, `db.system=postgresql` |
| Service Layer | Custom spans | Business attributes |

### Accessing Jaeger UI

```bash
# With Docker Compose (automatically generated)
docker-compose up -d

# Jaeger UI available at
open http://localhost:16686
```

### W3C traceparent Propagation

Traces are propagated between services via the standard HTTP header:

```
traceparent: 00-<trace-id>-<span-id>-01
```

This enables distributed tracing between microservices.

---

## Advanced Health Checks

### Kubernetes-compatible Endpoints

| Endpoint | K8s Usage | Behavior |
|----------|-----------|----------|
| `GET /health/liveness` | `livenessProbe` | Returns 200 if the application is running |
| `GET /health/readiness` | `readinessProbe` | Returns 200 if the DB is reachable, 503 otherwise |
| `GET /health` | Backward compatibility | Alias for `/health/liveness` |

### Liveness Probe

```bash
GET /health/liveness
```

**Response** (200):
```json
{
  "status": "alive",
  "service": "mon-app",
  "timestamp": "2026-02-17T10:00:00Z"
}
```

### Readiness Probe

```bash
GET /health/readiness
```

**Response** (200 — DB reachable):
```json
{
  "status": "ready",
  "service": "mon-app",
  "timestamp": "2026-02-17T10:00:00Z",
  "checks": {"database": "ok"}
}
```

**Response** (503 — DB unreachable):
```json
{
  "status": "not_ready",
  "service": "mon-app",
  "timestamp": "2026-02-17T10:00:00Z",
  "checks": {"database": "error"},
  "error": "database connection failed"
}
```

### Implementation

The `HealthHandler` receives `*gorm.DB` via fx injection and checks the connection with a 2-second timeout:

```go
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
    defer cancel()

    sqlDB, err := h.db.DB()
    if err != nil {
        return c.Status(503).JSON(...)
    }

    if err := sqlDB.PingContext(ctx); err != nil {
        return c.Status(503).JSON(...)
    }

    return c.JSON(/* ready response */)
}
```

### Prometheus Metrics for Health (advanced mode)

When `--observability=advanced` is enabled, health checks expose Prometheus metrics:

| Metric | Type | Description |
|--------|------|-------------|
| `health_check_status` | Gauge | 1 = healthy, 0 = unhealthy |
| `health_check_duration_seconds` | Histogram | Health check response time |

### Kubernetes Configuration

The `deployments/kubernetes/probes.yaml` file is automatically generated:

```yaml
# deployments/kubernetes/probes.yaml
livenessProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/readiness
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 3

startupProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  failureThreshold: 30
```

---

## Grafana Dashboard

### Pre-configured Dashboard

A Grafana JSON dashboard with **7 panels** is automatically generated and provisioned:

| Panel | Type | Description |
|-------|------|-------------|
| Request Rate | Time series | Requests per second |
| Error Rate | Time series | Error percentage (4xx, 5xx) |
| Latency P95 | Time series | Latency at 95th percentile |
| Latency P99 | Time series | Latency at 99th percentile |
| Requests in Flight | Gauge | Active requests in real time |
| Status Code Distribution | Pie chart | HTTP status code distribution |
| Top Endpoints | Table | Most requested endpoints |

### Generated Files

```
mon-app/
├── monitoring/
│   ├── grafana/
│   │   ├── provisioning/
│   │   │   ├── datasources/
│   │   │   │   └── prometheus.yml     # Auto-configured Prometheus datasource
│   │   │   └── dashboards/
│   │   │       └── dashboard.yml      # Dashboard auto-provisioning
│   │   └── dashboards/
│   │       └── app-dashboard.json     # 7-panel dashboard
│   └── prometheus/
│       ├── prometheus.yml             # Scraping configuration
│       └── alert_rules.yml            # Alerting rules
```

### Accessing Grafana

```bash
# Start the full stack
docker-compose up -d

# Grafana UI
open http://localhost:3000

# Default credentials
# Username: admin
# Password: admin
```

The dashboard is automatically provisioned on Grafana startup.

### Alerting Rules

The `alert_rules.yml` file includes pre-configured alerts:

| Alert | Threshold | Description |
|-------|-----------|-------------|
| HighErrorRate | > 5% for 5 min | High error rate |
| HighLatency | p95 > 1s for 5 min | High latency |
| ServiceDown | No metrics for 1 min | Service unreachable |

---

## Docker Compose — Observability Stack

### Complete Architecture

When `--observability=advanced` is enabled, the generated `docker-compose.yml` includes the full stack:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
      - jaeger

  db:
    image: postgres:16-alpine
    # ...

  jaeger:
    image: jaegertracing/all-in-one:1.56.0
    ports:
      - "16686:16686"    # Jaeger UI
      - "4317:4317"      # OTLP gRPC

  prometheus:
    image: prom/prometheus:v2.51.0
    ports:
      - "9090:9090"      # Prometheus UI
    volumes:
      - ./monitoring/prometheus:/etc/prometheus

  grafana:
    image: grafana/grafana:10.4.0
    ports:
      - "3000:3000"      # Grafana UI
    volumes:
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning
      - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards
```

### Ports and URLs

| Service | Port | URL |
|---------|------|-----|
| Application | 8080 | `http://localhost:8080` |
| Jaeger UI | 16686 | `http://localhost:16686` |
| Prometheus UI | 9090 | `http://localhost:9090` |
| Grafana UI | 3000 | `http://localhost:3000` |

### Starting the Stack

```bash
# Start all services
docker-compose up -d

# Verify everything is working
curl http://localhost:8080/health/readiness
curl http://localhost:8080/metrics

# Open the UIs
open http://localhost:16686   # Jaeger — view traces
open http://localhost:9090    # Prometheus — explore metrics
open http://localhost:3000    # Grafana — visual dashboard
```

### Fixed Service Versions

| Service | Version | Reason |
|---------|---------|--------|
| Jaeger | 1.56.0 | Latest stable with native OTLP support |
| Prometheus | v2.51.0 | OTLP and remote write support |
| Grafana | 10.4.0 | YAML provisioning + unified alerting |

---

## Navigation

**Previous**: [Deployment](deployment.md)  
**Next**: [Best Practices](best-practices.md)  
**Index**: [Guide Index](index.md)
