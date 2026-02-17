# Story 9.2: Distributed Tracing

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** développeur,
**Je veux** activer le distributed tracing avec OpenTelemetry dans les projets générés avec `--observability=advanced`,
**Afin de** tracer les requêtes à travers les services et corréler les logs avec les traces pour diagnostiquer les performances.

## Acceptance Criteria

1. **AC1**: Given un projet généré avec `--observability=advanced`, When une requête traverse l'API, Then un `trace_id` unique est généré et propagé via les headers HTTP (W3C `traceparent`)
2. **AC2**: Given OpenTelemetry est configuré, When les spans sont créés, Then ils sont exportés vers Jaeger (défaut) via OTLP/gRPC sur `localhost:4317`
3. **AC3**: Given une requête est tracée, When les logs zerolog sont émis dans le handler, Then ils incluent les champs `trace_id` et `span_id` pour la corrélation
4. **AC4**: Given le middleware Fiber de tracing est actif, When une requête arrive, Then le span parent est extrait des headers entrants (propagation W3C) ou un nouveau span root est créé
5. **AC5**: Given une requête implique une requête DB via GORM, When le tracing est actif, Then les callbacks GORM créent des spans enfants avec `db.system=postgresql`, `db.statement` (sans valeurs sensibles)
6. **AC6**: Given `--observability=advanced`, When le `docker-compose.yml` est généré, Then un service `jaeger` est inclus (image `jaegertracing/all-in-one:latest`) avec ports 16686 (UI) et 4317 (OTLP)
7. **AC7**: Given le projet est généré, When `go build ./...` s'exécute, Then le projet compile sans erreur
8. **AC8**: Given les tests du générateur s'exécutent, When `go test -short ./cmd/create-go-starter`, Then tous les tests passent

## Tasks / Subtasks

- [x] Task 1: Ajouter les templates OpenTelemetry dans `templates_observability.go` (AC: 1, 2, 4)
  - [x] 1.1 Template `pkg/tracing/tracer.go` — initialisation du TracerProvider OpenTelemetry avec exporter OTLP/gRPC
  - [x] 1.2 Template `internal/adapters/middleware/tracing_middleware.go` — middleware Fiber extrayant/créant spans avec propagation W3C
  - [x] 1.3 Template configuration tracer dans `pkg/tracing/tracer.go` : service name depuis `APP_NAME` env var, endpoint depuis `OTEL_EXPORTER_OTLP_ENDPOINT` (défaut: `localhost:4317`)
  - [x] 1.4 Template `shutdown` function dans `tracer.go` pour flush des spans (intégré au lifecycle fx)

- [x] Task 2: Intégrer le tracing avec zerolog (AC: 3)
  - [x] 2.1 Template de hook zerolog dans `pkg/logger/logger.go` (extension) — enrichissement des logs avec `trace_id` et `span_id` depuis le contexte
  - [x] 2.2 Template d'extraction des IDs depuis `otel.Span.SpanContext()` : `TraceID().String()` et `SpanID().String()`
  - [x] 2.3 S'assurer que le contexte Fiber est propagé : `c.UserContext()` → zerolog log avec trace context

- [x] Task 3: Intégrer le tracing GORM (AC: 5)
  - [x] 3.1 Dans `templates_observability.go`, ajouter template d'extension GORM avec callbacks OpenTelemetry
  - [x] 3.2 Callback `BeforeQuery` : créer span enfant `db.query` avec attributs `db.system`, `db.name`, `db.operation`
  - [x] 3.3 Callback `AfterQuery` : terminer le span avec statut (erreur ou succès)
  - [x] 3.4 Masquer les valeurs des paramètres SQL dans `db.statement` (sécurité)

- [x] Task 4: Mise à jour du docker-compose généré (AC: 6)
  - [x] 4.1 Dans `templates_observability.go`, template extension du `docker-compose.yml` pour ajouter Jaeger
  - [x] 4.2 Service Jaeger : `image: jaegertracing/all-in-one:latest`, ports `16686:16686` (UI), `4317:4317` (OTLP gRPC), `4318:4318` (OTLP HTTP)
  - [x] 4.3 Ajouter variable `OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4317` dans le service app du docker-compose
  - [x] 4.4 Ajouter `OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317` dans le `.env.example` généré

- [x] Task 5: Intégration fx lifecycle (AC: 2, 7)
  - [x] 5.1 Enregistrer le TracerProvider dans fx avec `fx.Provide(tracing.NewTracerProvider)`
  - [x] 5.2 Enregistrer le shutdown dans `fx.Hook` — appel de `tp.Shutdown(ctx)` au `OnStop`
  - [x] 5.3 Enregistrer le middleware tracing dans le setup du serveur Fiber
  - [x] 5.4 Vérifier que l'ordre fx est correct : TracerProvider avant ServerSetup

- [x] Task 6: Mise à jour go.mod template (AC: 7)
  - [x] 6.1 Ajouter dans le template go.mod pour `--observability=advanced` :
    - `go.opentelemetry.io/otel v1.30.0` (version stable 2024/2025)
    - `go.opentelemetry.io/otel/trace v1.30.0`
    - `go.opentelemetry.io/otel/sdk v1.30.0`
    - `go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.30.0`
    - `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.30.0`
    - Approche manuelle (middleware custom) plutôt qu'otelfiber pour éviter les conflits de version
  - [x] 6.2 Vérifier la compatibilité des versions entre elles

- [x] Task 7: Tests du générateur (AC: 8)
  - [x] 7.1 Dans `templates_observability_test.go`, ajouter tests pour le tracing
  - [x] 7.2 Test: `tracer.go` généré contient `NewTracerProvider` et `otlptracegrpc`
  - [x] 7.3 Test: `tracing_middleware.go` généré contient la propagation W3C
  - [x] 7.4 Test: `docker-compose.yml` généré avec `--observability=advanced` inclut service Jaeger
  - [x] 7.5 Test: `.env.example` généré inclut `OTEL_EXPORTER_OTLP_ENDPOINT`
  - [x] 7.6 Test E2E: projet généré compile avec `go build ./...` (couvert par les tests existants `go test -short`)

## Dev Notes

### Architecture du Tracing

Cette story **étend** `templates_observability.go` créé en Story 9.1. Le tracing est une couche orthogonale aux métriques, mais utilise le même flag `--observability=advanced`.

**Nouveaux fichiers générés dans le projet cible :**
```
<projet-généré>/
├── pkg/tracing/
│   └── tracer.go           # TracerProvider + shutdown
└── internal/adapters/
    └── middleware/
        └── tracing_middleware.go  # Fiber middleware W3C propagation
```

**Fichiers existants étendus dans le projet généré :**
```
pkg/logger/logger.go          # Hook zerolog pour trace_id/span_id
internal/infrastructure/database.go  # GORM callbacks pour DB spans
docker-compose.yml            # Service Jaeger ajouté
.env.example                  # OTEL vars ajoutées
cmd/<project>/main.go         # TracerProvider enregistré via fx
```

### Template `pkg/tracing/tracer.go`

```go
package tracing

import (
    "context"
    "os"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
    "go.uber.org/fx"
)

// Params pour fx DI
type Params struct {
    fx.In
    LC fx.Lifecycle
}

// NewTracerProvider initialise et enregistre le TracerProvider OTLP.
func NewTracerProvider(params Params) *sdktrace.TracerProvider {
    endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
    if endpoint == "" {
        endpoint = "localhost:4317"
    }
    serviceName := os.Getenv("APP_NAME")
    if serviceName == "" {
        serviceName = "{{project_name}}" // Injecté par le générateur
    }

    exporter, err := otlptracegrpc.New(
        context.Background(),
        otlptracegrpc.WithEndpoint(endpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        otel.SetTracerProvider(sdktrace.NewTracerProvider())
        return sdktrace.NewTracerProvider()
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(serviceName),
        )),
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
    )
    otel.SetTracerProvider(tp)

    params.LC.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return tp.Shutdown(ctx)
        },
    })
    return tp
}
```

### Template Middleware Fiber

```go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gofiber/fiber/otelfiber"
)

// TracingMiddleware retourne le middleware OpenTelemetry pour Fiber.
func TracingMiddleware() fiber.Handler {
    return otelfiber.Middleware(
        otelfiber.WithSpanNameFormatter(func(ctx *fiber.Ctx) string {
            return ctx.Method() + " " + ctx.Route().Path
        }),
    )
}
```

**Alternative manuelle si otelfiber non disponible :**
```go
func TracingMiddleware() fiber.Handler {
    tracer := otel.Tracer("fiber")
    propagator := otel.GetTextMapPropagator()

    return func(c *fiber.Ctx) error {
        ctx := propagator.Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))
        ctx, span := tracer.Start(ctx, c.Method()+" "+c.Path())
        defer span.End()

        c.SetUserContext(ctx)
        span.SetAttributes(
            attribute.String("http.method", c.Method()),
            attribute.String("http.url", c.OriginalURL()),
        )

        err := c.Next()

        span.SetAttributes(attribute.Int("http.status_code", c.Response().StatusCode()))
        if err != nil {
            span.RecordError(err)
        }
        return err
    }
}
```

### Corrélation Logs ↔ Traces

```go
// Extension de pkg/logger/logger.go
import (
    "go.opentelemetry.io/otel/trace"
    "github.com/rs/zerolog"
)

// WithTraceContext enrichit le logger avec trace_id et span_id.
func WithTraceContext(ctx context.Context, logger zerolog.Logger) zerolog.Logger {
    span := trace.SpanFromContext(ctx)
    if !span.IsRecording() {
        return logger
    }
    sc := span.SpanContext()
    return logger.With().
        Str("trace_id", sc.TraceID().String()).
        Str("span_id", sc.SpanID().String()).
        Logger()
}
```

**Usage dans les handlers :**
```go
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
    log := logger.WithTraceContext(c.UserContext(), h.log)
    log.Info().Msg("Getting user profile")
    // ...
}
```

### GORM Spans

```go
// callbacks pour database.go
func registerTracingCallbacks(db *gorm.DB) {
    tracer := otel.Tracer("gorm")

    db.Callback().Query().Before("gorm:query").Register("otel:before_query",
        func(db *gorm.DB) {
            ctx, span := tracer.Start(db.Statement.Context, "db.query")
            span.SetAttributes(
                attribute.String("db.system", "postgresql"),
                attribute.String("db.operation", "query"),
                attribute.String("db.sql.table", db.Statement.Table),
            )
            db.Statement.Context = ctx
        })

    db.Callback().Query().After("gorm:query").Register("otel:after_query",
        func(db *gorm.DB) {
            span := trace.SpanFromContext(db.Statement.Context)
            if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
                span.RecordError(db.Error)
                span.SetStatus(codes.Error, db.Error.Error())
            }
            span.End()
        })
}
```

### Variables d'Environnement Ajoutées

Dans le `.env.example` généré :
```
# OpenTelemetry (--observability=advanced)
APP_NAME={{project_name}}
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true
```

### docker-compose.yml — Service Jaeger

```yaml
# Ajouté au docker-compose.yml généré
  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: {{project_name}}-jaeger
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    ports:
      - "16686:16686"   # Jaeger UI
      - "4317:4317"     # OTLP gRPC
      - "4318:4318"     # OTLP HTTP
    networks:
      - {{project_name}}-network
```

### Versions OpenTelemetry

**CRITIQUE** : Les versions OpenTelemetry doivent être cohérentes entre elles (sinon build errors). Utiliser la même version de base pour tous les packages `go.opentelemetry.io/otel/*` :
- `go.opentelemetry.io/otel` v1.24.0
- `go.opentelemetry.io/otel/trace` v1.24.0 (même version)
- `go.opentelemetry.io/otel/sdk` v1.24.0 (même version)
- Vérifier les versions 2026 disponibles avant d'implémenter

### Dépendance avec Story 9.1

- Story 9.2 **étend** `templates_observability.go` créé en 9.1
- Le flag `--observability=advanced` est déjà défini (pas de modification de main.go requise)
- Les callbacks GORM de 9.1 (métriques) et 9.2 (traces) coexistent dans le même `database.go` étendu
- Le middleware Fiber de 9.1 (Prometheus) et 9.2 (tracing) sont deux middlewares distincts, appliqués dans cet ordre : `app.Use(prometheusMiddleware)` puis `app.Use(tracingMiddleware)`

### Anti-Patterns à Éviter

- NE PAS utiliser `trace.SpanFromContext` sans vérifier `IsRecording()` — peut retourner un NoopSpan
- NE PAS inclure les valeurs des paramètres SQL dans les spans (données sensibles)
- NE PAS utiliser `AlwaysSample()` en production réelle (suggérer `ParentBased(TraceIDRatioBased(0.1))` dans les commentaires du code)
- NE PAS oublier `defer span.End()` — les spans non terminés ne sont pas exportés
- NE PAS utiliser `otlptracegrpc.WithTLSClientConfig(nil)` — utiliser `WithInsecure()` pour le dev local
- NE PAS bloquer l'application si Jaeger n'est pas disponible — le SDK OTel ignore les erreurs d'export silencieusement

### Project Structure Notes

- Cette story ne crée pas de nouveau flag, elle étend les templates du flag existant `--observability=advanced`
- `templates_observability.go` dans le générateur CLI est enrichi (PAS remplacé)
- Les tests de la story 9.1 doivent continuer à passer après cette story

### References

- [Source: cmd/create-go-starter/templates_observability.go] - Fichier créé en Story 9.1
- [Source: cmd/create-go-starter/generator.go] - Pattern génération conditionnelle
- [Source: _bmad-output/planning-artifacts/architecture.md#Infrastructure & Deployment] - Stack zerolog, fx lifecycle
- [Source: _bmad-output/planning-artifacts/epics.md#Story 9.2] - Story specification
- [Source: _bmad-output/project-context.md] - Rules fx, no globals, graceful shutdown
- [OpenTelemetry Go: https://opentelemetry.io/docs/languages/go/] - Getting started
- [otelfiber: https://github.com/open-telemetry/opentelemetry-go-contrib] - Fiber instrumentation
- [Jaeger OTLP: https://www.jaegertracing.io/docs/latest/apis/] - OTLP endpoints

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

- Implémenté le distributed tracing OpenTelemetry complet pour les projets générés avec `--observability=advanced`
- Ajouté `TracerTemplate()` : `pkg/tracing/tracer.go` avec TracerProvider OTLP/gRPC, propagation W3C et lifecycle fx
- Ajouté `TracingMiddlewareTemplate()` : middleware Fiber manuel pour extraction/création de spans W3C
- Ajouté `GORMTracingTemplate()` : `pkg/tracing/db_tracing.go` avec callbacks BeforeQuery/AfterQuery (sans valeurs SQL sensibles)
- Ajouté `LoggerWithTracingTemplate()` : `pkg/logger/logger.go` étendu avec `WithTraceContext` pour corrélation logs/traces
- Ajouté `EnvTemplateWithObservability()` : `.env.example` avec `OTEL_EXPORTER_OTLP_ENDPOINT` et vars OTel
- Ajouté `DockerComposeTemplateWithObservability()` : `docker-compose.yml` avec service Jaeger all-in-one (ports 16686, 4317, 4318)
- Mis à jour `GoModTemplateWithObservability()` : dépendances OTel v1.30.0 cohérentes (otel, otel/trace, otel/sdk, otlptracegrpc, google.golang.org/grpc)
- Mis à jour `MainGoTemplateWithObservability()` : `fx.Provide(tracing.NewTracerProvider)` + `fx.Invoke(tracing.RegisterGORMCallbacks)`
- Mis à jour `ServerTemplateWithObservability()` : paramètre `tp *sdktrace.TracerProvider` + `app.Use(middleware.TracingMiddleware())`
- Mis à jour `generator.go` : utilisation des templates enrichis pour logger, .env.example, docker-compose, et génération des fichiers tracing
- Ajouté 6 nouveaux tests dans `templates_observability_test.go` couvrant AC 1-7 — tous passent ✅
- Décision architecture : middleware manuel plutôt qu'otelfiber pour éviter les conflits de version Go modules

### Code Review Fixes (revue adversariale — claude-opus-4.6)

- **HIGH #1**: File List corrigée — 8 fichiers manquants ajoutés (main.go, main_test.go, generator_test.go, database_integration_test.go, e2e_mysql_test.go, e2e_sqlite_test.go, smoke_test.go, template_minimal_test.go)
- **HIGH #2**: Erreur ignorée dans TracerTemplate() corrigée — `exporter, _ :=` remplacé par gestion d'erreur avec fallback no-op provider
- **MEDIUM #3**: `pkg/tracing` et `pkg/metrics` ajoutés dans `getDirectoriesForTemplate()` pour cohérence avec le pattern de création de répertoires
- **MEDIUM #4**: Callbacks GORM étendus — ajout de Create, Update, Delete en plus de Query pour tracer INSERT/UPDATE/DELETE
- **MEDIUM #5**: `db.system` rendu dynamique — accepte le paramètre `database` et mappe vers postgresql/mysql/sqlite
- **MEDIUM #6**: `version: '3.8'` déprécié retiré du Docker Compose template
- **LOW #7**: Commentaire d'ordre d'initialisation ajouté pour la dépendance implicite propagator/TracerProvider
- **LOW #8**: Tests ajoutés pour db_tracing.go — vérifie RegisterGORMCallbacks, callbacks CRUD, et db.system dynamique (MySQL)
- **LOW #9**: Références AC dans les commentaires de test corrigées — préfixées avec le numéro de story (9.1/9.2)

### File List

- cmd/create-go-starter/templates_observability.go (modified — ajout templates tracing)
- cmd/create-go-starter/templates_observability_test.go (modified — ajout tests tracing)
- cmd/create-go-starter/generator.go (modified — observabilityLevel param, generateObservabilityFiles)
- cmd/create-go-starter/main.go (modified — observability constants, validation, CLI flag, run() 4 params)
- cmd/create-go-starter/main_test.go (modified — run() signature 3→4 params)
- cmd/create-go-starter/generator_test.go (modified — generateProjectFiles() signature 4→5 params)
- cmd/create-go-starter/database_integration_test.go (modified — generateProjectFiles() signature update)
- cmd/create-go-starter/e2e_mysql_test.go (modified — generateProjectFiles() signature update)
- cmd/create-go-starter/e2e_sqlite_test.go (modified — generateProjectFiles() signature update)
- cmd/create-go-starter/smoke_test.go (modified — run() signature update)
- cmd/create-go-starter/template_minimal_test.go (modified — run() signature update)
