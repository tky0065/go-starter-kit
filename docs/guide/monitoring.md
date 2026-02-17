# Monitoring & Observabilité

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

<i class="material-icons success">new_releases</i> **Nouveau dans v1.3.0** — L'observabilité avancée est désormais intégrée nativement dans les projets générés via le flag `--observability`.

## Vue d'ensemble

Go Starter Kit propose **3 niveaux d'observabilité** pour monitorer vos projets en production :

| Niveau | Flag | Description |
|--------|------|-------------|
| `none` | `--observability=none` (défaut) | Comportement standard, aucune instrumentation |
| `basic` | `--observability=basic` | Health checks avancés K8s (liveness/readiness) |
| `advanced` | `--observability=advanced` | Stack complète : Prometheus + Jaeger + Grafana + Health Checks |

```bash
# Générer un projet avec observabilité avancée
create-go-starter mon-app --template=full --observability=advanced
```

> <i class="material-icons warning">warning</i> **Note** : Le flag `--observability` ne fonctionne qu'avec le template `full`. Les templates `minimal` et `graphql` ne sont pas supportés.

---

## Logging avec zerolog

### Configuration

Le logger est configuré dans `pkg/logger/logger.go` :

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

### Utilisation

Injection via fx :

```go
type UserService struct {
    logger zerolog.Logger
}

func NewUserService(logger zerolog.Logger) *UserService {
    return &UserService{logger: logger}
}
```

**Logging structuré** :

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

### Logging enrichi avec tracing (mode advanced)

Quand `--observability=advanced` est activé, le logger est enrichi avec les identifiants de trace OpenTelemetry :

```go
// Logs incluent automatiquement trace_id et span_id
logger.Info().
    Str("trace_id", span.SpanContext().TraceID().String()).
    Str("span_id", span.SpanContext().SpanID().String()).
    Msg("Processing request")
```

**Exemple de sortie JSON** :
```json
{
  "level": "info",
  "trace_id": "abc123def456...",
  "span_id": "789xyz...",
  "message": "Processing request",
  "timestamp": 1708000000
}
```

Cela permet de corréler les logs avec les traces dans Jaeger.

### Niveaux de log

| Niveau | Usage |
|--------|-------|
| **Debug** | Informations détaillées pour debugging |
| **Info** | Événements importants (user login, etc.) |
| **Warn** | Comportements anormaux non-critiques |
| **Error** | Erreurs nécessitant attention |
| **Fatal** | Erreurs critiques (app exit) |

### Best practices

<i class="material-icons success">check_circle</i> **BON — Structured logging** :

```go
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Dur("duration", elapsed).
    Msg("User logged in")
```

<i class="material-icons error">cancel</i> **MAUVAIS — String formatting** :

```go
logger.Info().Msgf("User %s logged in after %v", userID, elapsed)
```

<i class="material-icons success">check_circle</i> **BON — Pas de secrets** :

```go
logger.Info().Str("email", email).Msg("User login attempt")
```

<i class="material-icons error">cancel</i> **MAUVAIS — Logging secrets** :

```go
logger.Info().Str("password", password).Msg("Login")  // NEVER!
```

---

## Prometheus Metrics

### Endpoint `/metrics`

Quand `--observability=advanced` est activé, un endpoint Prometheus est exposé :

```bash
GET /metrics
```

### Métriques HTTP exposées

| Métrique | Type | Description |
|----------|------|-------------|
| `http_requests_total` | Counter | Requêtes totales par méthode, route et code status |
| `http_request_duration_seconds` | Histogram | Latence HTTP par route (buckets p50, p90, p95, p99) |
| `http_requests_in_flight` | Gauge | Nombre de requêtes actives en cours |

### Librairie utilisée

[`fiberprometheus/v2`](https://github.com/ansrivas/fiberprometheus) v2.7.0 — intégration native Fiber v2.

### Fichiers générés

```
mon-app/
├── pkg/metrics/
│   └── prometheus.go                # Registry Prometheus + PrometheusMetrics struct
├── internal/adapters/
│   ├── middleware/
│   │   └── metrics_middleware.go    # Middleware HTTP pour capturer les métriques
│   └── handlers/
│       └── metrics_handler.go       # Handler GET /metrics
```

### Exemple de configuration Prometheus

```yaml
# prometheus.yml (généré automatiquement dans monitoring/prometheus/)
scrape_configs:
  - job_name: 'mon-app'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Tester les métriques

```bash
# Générer du trafic
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users

# Voir les métriques
curl http://localhost:8080/metrics

# Exemple de sortie
# http_requests_total{method="GET",path="/health",status="200"} 5
# http_request_duration_seconds_bucket{method="GET",path="/health",le="0.1"} 5
# http_requests_in_flight 0
```

---

## Distributed Tracing avec OpenTelemetry

### Architecture

```
Application → OTLP/gRPC → Jaeger Collector → Jaeger UI
```

### Configuration

Le tracing utilise **OpenTelemetry** avec export OTLP/gRPC vers **Jaeger** :

- **Protocole** : OTLP/gRPC
- **Propagation** : W3C `traceparent` header
- **Endpoint** : Configurable via `OTEL_EXPORTER_OTLP_ENDPOINT` (défaut : `localhost:4317`)

### Fichiers générés

```
mon-app/
├── pkg/tracing/
│   └── tracer.go                    # Configuration OpenTelemetry + TracerProvider
├── internal/adapters/middleware/
│   └── tracing_middleware.go        # Middleware HTTP pour créer des spans
├── internal/infrastructure/database/
│   └── tracing.go                   # Instrumentation GORM avec spans DB
└── pkg/logger/
    └── logger_tracing.go            # Logger enrichi avec trace_id/span_id
```

### Variables d'environnement

```bash
# .env.example (ajouté automatiquement)
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=mon-app
```

### Spans générés automatiquement

| Composant | Span | Attributs |
|-----------|------|-----------|
| HTTP Middleware | `HTTP {method} {path}` | `http.method`, `http.route`, `http.status_code` |
| GORM Tracing | `db.query` | `db.statement`, `db.system=postgresql` |
| Service Layer | Custom spans | Attributs métier |

### Accéder à Jaeger UI

```bash
# Avec Docker Compose (généré automatiquement)
docker-compose up -d

# Jaeger UI disponible sur
open http://localhost:16686
```

### Propagation W3C traceparent

Les traces sont propagées entre services via le header HTTP standard :

```
traceparent: 00-<trace-id>-<span-id>-01
```

Cela permet le tracing distribué entre microservices.

---

## Health Checks Avancés

### Endpoints Kubernetes-compatible

| Endpoint | Usage K8s | Comportement |
|----------|-----------|--------------|
| `GET /health/liveness` | `livenessProbe` | Retourne 200 si l'application tourne |
| `GET /health/readiness` | `readinessProbe` | Retourne 200 si la DB est accessible, 503 sinon |
| `GET /health` | Rétrocompatibilité | Alias vers `/health/liveness` |

### Liveness Probe

```bash
GET /health/liveness
```

**Réponse** (200) :
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

**Réponse** (200 — DB accessible) :
```json
{
  "status": "ready",
  "service": "mon-app",
  "timestamp": "2026-02-17T10:00:00Z",
  "checks": {"database": "ok"}
}
```

**Réponse** (503 — DB inaccessible) :
```json
{
  "status": "not_ready",
  "service": "mon-app",
  "timestamp": "2026-02-17T10:00:00Z",
  "checks": {"database": "error"},
  "error": "database connection failed"
}
```

### Implémentation

Le `HealthHandler` reçoit `*gorm.DB` via injection fx et vérifie la connexion avec un timeout de 2 secondes :

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

### Métriques Prometheus pour Health (mode advanced)

Quand `--observability=advanced` est activé, les health checks exposent des métriques Prometheus :

| Métrique | Type | Description |
|----------|------|-------------|
| `health_check_status` | Gauge | 1 = healthy, 0 = unhealthy |
| `health_check_duration_seconds` | Histogram | Temps de réponse des checks |

### Configuration Kubernetes

Le fichier `deployments/kubernetes/probes.yaml` est automatiquement généré :

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

### Dashboard pré-configuré

Un dashboard Grafana JSON avec **7 panneaux** est automatiquement généré et provisionné :

| Panneau | Type | Description |
|---------|------|-------------|
| Request Rate | Time series | Requêtes par seconde |
| Error Rate | Time series | Pourcentage d'erreurs (4xx, 5xx) |
| Latency P95 | Time series | Latence au 95ème percentile |
| Latency P99 | Time series | Latence au 99ème percentile |
| Requests in Flight | Gauge | Requêtes actives en temps réel |
| Status Code Distribution | Pie chart | Répartition des codes HTTP |
| Top Endpoints | Table | Endpoints les plus sollicités |

### Fichiers générés

```
mon-app/
├── monitoring/
│   ├── grafana/
│   │   ├── provisioning/
│   │   │   ├── datasources/
│   │   │   │   └── prometheus.yml     # Datasource Prometheus auto-configurée
│   │   │   └── dashboards/
│   │   │       └── dashboard.yml      # Auto-provisioning des dashboards
│   │   └── dashboards/
│   │       └── app-dashboard.json     # Dashboard 7 panneaux
│   └── prometheus/
│       ├── prometheus.yml             # Configuration scraping
│       └── alert_rules.yml            # Règles d'alerting
```

### Accéder à Grafana

```bash
# Démarrer la stack complète
docker-compose up -d

# Grafana UI
open http://localhost:3000

# Credentials par défaut
# Username: admin
# Password: admin
```

Le dashboard est automatiquement provisionné au démarrage de Grafana.

### Règles d'alerting

Le fichier `alert_rules.yml` inclut des alertes pré-configurées :

| Alerte | Seuil | Description |
|--------|-------|-------------|
| HighErrorRate | > 5% pendant 5 min | Taux d'erreurs élevé |
| HighLatency | p95 > 1s pendant 5 min | Latence élevée |
| ServiceDown | Absence de métriques pendant 1 min | Service injoignable |

---

## Docker Compose — Stack Observabilité

### Architecture complète

Quand `--observability=advanced` est activé, le `docker-compose.yml` généré inclut la stack complète :

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

### Ports et URLs

| Service | Port | URL |
|---------|------|-----|
| Application | 8080 | `http://localhost:8080` |
| Jaeger UI | 16686 | `http://localhost:16686` |
| Prometheus UI | 9090 | `http://localhost:9090` |
| Grafana UI | 3000 | `http://localhost:3000` |

### Démarrer la stack

```bash
# Démarrer tous les services
docker-compose up -d

# Vérifier que tout fonctionne
curl http://localhost:8080/health/readiness
curl http://localhost:8080/metrics

# Ouvrir les UIs
open http://localhost:16686   # Jaeger — voir les traces
open http://localhost:9090    # Prometheus — explorer les métriques
open http://localhost:3000    # Grafana — dashboard visuel
```

### Versions fixes des services

| Service | Version | Raison |
|---------|---------|--------|
| Jaeger | 1.56.0 | Dernière stable avec support OTLP natif |
| Prometheus | v2.51.0 | Support OTLP et remote write |
| Grafana | 10.4.0 | Provisioning YAML + alerting unifié |

---

## Navigation

**Previous**: [Déploiement](deployment.md)  
**Next**: [Bonnes pratiques](best-practices.md)  
**Index**: [Guide Index](index.md)
