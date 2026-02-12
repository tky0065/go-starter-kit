# Monitoring & Logging

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Logging avec zerolog

#### Configuration

Le logger est configuré dans `pkg/logger/logger.go`:

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

    // Set level based on env
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

#### Utilisation

Injection via fx:

```go
type UserService struct {
    logger zerolog.Logger
}

func NewUserService(logger zerolog.Logger) *UserService {
    return &UserService{logger: logger}
}
```

**Logging structuré**:

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

// Fatal (exits)
logger.Fatal().
    Err(err).
    Msg("Cannot connect to database")
```

#### Niveaux de log

| Niveau | Usage |
|--------|-------|
| **Debug** | Informations détaillées pour debugging |
| **Info** | Événements importants (user login, etc.) |
| **Warn** | Comportements anormaux non-critiques |
| **Error** | Erreurs nécessitant attention |
| **Fatal** | Erreurs critiques (app exit) |

#### Best practices

**:material-check-circle: BON - Structured logging**:

```go
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Dur("duration", elapsed).
    Msg("User logged in")
```

**❌ MAUVAIS - String formatting**:

```go
logger.Info().Msgf("User %s logged in after %v", userID, elapsed)
```

**:material-check-circle: BON - Pas de secrets**:

```go
logger.Info().Str("email", email).Msg("User login attempt")
```

**❌ MAUVAIS - Logging secrets**:

```go
logger.Info().Str("password", password).Msg("Login")  // NEVER!
```

### Monitoring (recommandations)

Pour la production, intégrer:

#### 1. Prometheus + Grafana

**Metrics à collecter**:
- Request rate (requests/sec)
- Response time (p50, p95, p99)
- Error rate (%)
- Database connection pool
- CPU/Memory usage

**Implémenter avec fiber/prometheus**:

```go
import "github.com/gofiber/adaptor/v2"
import "github.com/prometheus/client_golang/prometheus/promhttp"

app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
```

#### 2. Jaeger / OpenTelemetry

**Distributed tracing** pour suivre les requêtes à travers les services.

#### 3. Sentry

**Error tracking** en temps réel:

```go
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.ClientOptions{
    Dsn: os.Getenv("SENTRY_DSN"),
})

// Capture errors
sentry.CaptureException(err)
```

#### 4. APM (Application Performance Monitoring)

- **New Relic**: APM complet
- **Datadog**: Monitoring + logs
- **Elastic APM**: Open source

### Health checks

L'endpoint `/health` est crucial pour:

- Load balancers
- Kubernetes probes
- Monitoring tools

**Amélioré**:

```go
type HealthResponse struct {
    Status   string            `json:"status"`
    Version  string            `json:"version"`
    Services map[string]string `json:"services"`
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
    // Check database
    dbStatus := "ok"
    if err := h.db.Exec("SELECT 1").Error; err != nil {
        dbStatus = "error"
    }

    response := HealthResponse{
        Status:  "ok",
        Version: "1.0.0",
        Services: map[string]string{
            "database": dbStatus,
        },
    }

    if dbStatus != "ok" {
        return c.Status(fiber.StatusServiceUnavailable).JSON(response)
    }

    return c.JSON(response)
}
```

---


---

## Navigation

**Previous**: [Déploiement](deployment.md)  
**Next**: [Bonnes pratiques](best-practices.md)  
**Index**: [Guide Index](index.md)
