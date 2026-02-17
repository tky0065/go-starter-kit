# Story 9.3: Advanced Health Checks

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** DevOps,
**Je veux** des health checks avancés (`/health/liveness` et `/health/readiness`) dans les projets générés,
**Afin que** Kubernetes puisse gérer correctement le cycle de vie de l'application (restart sur liveness failure, no-traffic sur readiness failure).

## Acceptance Criteria

1. **AC1**: Given le projet est généré (avec ou sans `--observability`), When `GET /health/liveness` est appelé, Then retourne HTTP 200 avec `{"status":"alive","timestamp":"..."}` si l'application tourne
2. **AC2**: Given le projet est généré, When `GET /health/readiness` est appelé et la DB est accessible, Then retourne HTTP 200 avec `{"status":"ready","checks":{"database":"ok"},"timestamp":"..."}`
3. **AC3**: Given la connexion DB est défaillante, When `GET /health/readiness` est appelé, Then retourne HTTP 503 avec `{"status":"not_ready","checks":{"database":"error"},"error":"database connection failed"}` (message générique par sécurité — ne pas exposer les détails internes)
4. **AC4**: Given le projet est généré avec `--observability=advanced`, When `/health/readiness` est appelé, Then les métriques de health check sont exposées dans `/metrics` (`health_check_duration_seconds`, `health_check_status`)
5. **AC5**: Given le projet généré utilise les health checks avancés, When l'ancien endpoint `GET /health` est appelé, Then il retourne également HTTP 200 (rétrocompatibilité — redirige ou alias vers `/health/liveness`)
6. **AC6**: Given les health checks sont configurés, When le projet est généré, Then un fichier de probe K8s est inclus dans `deployments/kubernetes/probes.yaml` avec livenessProbe et readinessProbe
7. **AC7**: Given le projet est généré, When `go build ./...` s'exécute, Then le projet compile sans erreur
8. **AC8**: Given les tests s'exécutent, When `go test -short ./cmd/create-go-starter`, Then tous les tests passent

## Tasks / Subtasks

- [x] Task 1: Créer le template du handler health avancé (AC: 1, 2, 3, 5)
  - [x] 1.1 Dans `templates.go` (méthode `AdvancedHealthHandlerTemplate()`), créer template `internal/adapters/handlers/health_handler.go`
  - [x] 1.2 Handler `Liveness()` : retourne `{"status":"alive","service":"{{project_name}}","timestamp":"<ISO8601>"}` — toujours 200 si app tourne
  - [x] 1.3 Handler `Readiness()` : appelle `db.DB().PingContext(ctx)` avec timeout de 2s, retourne 200 ou 503 selon résultat
  - [x] 1.4 Struct `HealthResponse` avec champs JSON snake_case : `status`, `service`, `timestamp`, `checks`, `error`
  - [x] 1.5 Compatibilité : `app.Get("/health", healthHandler.Liveness)` dans routes.go — alias de liveness

- [x] Task 2: Mettre à jour les routes du serveur généré (AC: 1, 2, 5)
  - [x] 2.1 `RoutesTemplate()` dans `templates_user.go` : 3 routes health
    - `app.Get("/health", healthHandler.Liveness)` — rétrocompatibilité
    - `app.Get("/health/liveness", healthHandler.Liveness)`
    - `app.Get("/health/readiness", healthHandler.Readiness)`
  - [x] 2.2 Ces routes ne passent PAS par le middleware JWT (enregistrées avant les routes /api/v1)
  - [x] 2.3 Ces routes NE sont PAS derrière le préfixe `/api/v1`

- [x] Task 3: Intégrer `HealthHandler` dans le module fx (AC: 7)
  - [x] 3.1 `NewHealthHandler(db *gorm.DB) *HealthHandler` — injecté automatiquement par fx
  - [x] 3.2 `fx.Provide(NewHealthHandler)` ajouté dans `HandlerModuleTemplate()`
  - [x] 3.3 `*handlers.HealthHandler` injecté dans `RegisterRoutes` via fx

- [x] Task 4: Métriques de health check avec Prometheus (AC: 4)
  - [x] 4.1 `PrometheusTemplate()` mis à jour avec `HealthCheckDuration` (HistogramVec) et `HealthCheckStatus` (GaugeVec)
  - [x] 4.2 `HealthHandlerWithMetricsTemplate()` dans `templates_observability.go` — enregistre métriques dans `Readiness()`
  - [x] 4.3 Version avancée générée uniquement si `--observability=advanced`

- [x] Task 5: Fichier de configuration Kubernetes (AC: 6)
  - [x] 5.1 `KubernetesProbesTemplate()` dans `templates_observability.go` — génère `deployments/kubernetes/probes.yaml`
  - [x] 5.2 Contenu : livenessProbe, readinessProbe, startupProbe avec commentaires explicatifs
  - [x] 5.3 Valeurs recommandées : initialDelaySeconds, periodSeconds, timeoutSeconds, failureThreshold

- [x] Task 6: Tests du générateur (AC: 8)
  - [x] 6.1 Dans `templates_observability_test.go`, 7 nouveaux tests pour health checks
  - [x] 6.2 Test: `health_handler.go` contient `Liveness` et `Readiness`
  - [x] 6.3 Test: routes `/health/liveness` et `/health/readiness` présentes dans routes.go généré
  - [x] 6.4 Test: `deployments/kubernetes/probes.yaml` est généré
  - [x] 6.5 Test: `HealthResponse` a les bons champs JSON snake_case
  - [x] 6.6 Test: la route `/health` (rétrocompatibilité) est toujours présente

- [x] Task 7: Documentation (AC: 1, 2, 3)
  - [x] 7.1 Annotations Swagger `@Summary`, `@Router`, `@Produce`, `@Success`, `@Failure` sur Liveness et Readiness
  - [x] 7.2 `docs/generated-project-guide.md` mis à jour avec section health checks K8s et exemples
  - [x] 7.3 `docs/usage.md` mis à jour avec les 3 endpoints `/health`, `/health/liveness`, `/health/readiness`

## Dev Notes

### Architecture des Health Checks

**IMPORTANT** : Les health checks avancés sont générés **pour TOUS les projets** (pas seulement `--observability=advanced`). C'est une bonne pratique universelle. Seule l'intégration Prometheus (AC4) est conditionnelle à `--observability=advanced`.

**Cette story modifie :**
1. Le template `internal/infrastructure/server.go` — nouvelles routes
2. Le template `internal/adapters/handlers/health_handler.go` — remplace ou étend le handler existant
3. Le générateur (`generator.go`) — générer `probes.yaml` et mettre à jour le health handler

### Template `internal/adapters/handlers/health_handler.go`

```go
package handlers

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

// HealthHandler gère les endpoints de health check.
type HealthHandler struct {
    db *gorm.DB
}

// NewHealthHandler crée un nouveau HealthHandler.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
    return &HealthHandler{db: db}
}

// HealthResponse représente la réponse d'un health check.
type HealthResponse struct {
    Status    string            `json:"status"`
    Service   string            `json:"service"`
    Timestamp string            `json:"timestamp"`
    Checks    map[string]string `json:"checks,omitempty"`
    Error     string            `json:"error,omitempty"`
}

// Liveness vérifie que l'application est vivante.
// @Summary Liveness probe
// @Description Vérifie que l'application tourne (utilisé par K8s livenessProbe)
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health/liveness [get]
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
    return c.Status(fiber.StatusOK).JSON(HealthResponse{
        Status:    "alive",
        Service:   "{{project_name}}", // Injecté par le générateur
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    })
}

// Readiness vérifie que l'application est prête à recevoir du trafic.
// @Summary Readiness probe
// @Description Vérifie la DB et les dépendances (utilisé par K8s readinessProbe)
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health/readiness [get]
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
    defer cancel()

    checks := make(map[string]string)
    var checkError string

    // Vérifier la connexion DB
    sqlDB, err := h.db.DB()
    if err != nil || sqlDB.PingContext(ctx) != nil {
        checks["database"] = "error"
        checkError = "database connection failed"
        return c.Status(fiber.StatusServiceUnavailable).JSON(HealthResponse{
            Status:    "not_ready",
            Service:   "{{project_name}}",
            Timestamp: time.Now().UTC().Format(time.RFC3339),
            Checks:    checks,
            Error:     checkError,
        })
    }
    checks["database"] = "ok"

    return c.Status(fiber.StatusOK).JSON(HealthResponse{
        Status:    "ready",
        Service:   "{{project_name}}",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Checks:    checks,
    })
}
```

### Routes dans server.go

```go
// Remplacer/étendre la configuration des routes dans server.go généré

// Health routes (pas derrière /api/v1 ni le middleware auth)
healthHandler := handlers.NewHealthHandler(db) // Via fx
app.Get("/health", healthHandler.Liveness)           // Rétrocompatibilité
app.Get("/health/liveness", healthHandler.Liveness)
app.Get("/health/readiness", healthHandler.Readiness)
```

**Remarque** : Vérifier le template server.go existant (`templates.go`) pour trouver l'endroit exact où insérer ces routes, afin de ne pas briser les routes existantes.

### Fichier K8s `deployments/kubernetes/probes.yaml`

```yaml
# Référence de configuration K8s pour les probes
# Copier dans votre Deployment K8s

# livenessProbe: redémarre le pod si l'app est bloquée
livenessProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 15
  failureThreshold: 3
  timeoutSeconds: 5

# readinessProbe: retire le pod du load balancer si pas prêt
readinessProbe:
  httpGet:
    path: /health/readiness
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 3
  timeoutSeconds: 5

# startupProbe: laisse du temps au démarrage (désactive liveness pendant ce temps)
startupProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  failureThreshold: 30
  periodSeconds: 5
  timeoutSeconds: 5
```

### Format de Réponse K8s Compatible

Kubernetes attend uniquement un HTTP status code 200 pour "ready/alive". Le format JSON est pour les humains et les outils de monitoring. Nginx/Kong utilisent aussi ce format.

**Règle critique** : `/health/liveness` doit retourner 200 MÊME si la DB est down. La liveness sert uniquement à détecter les deadlocks/panics, pas les dépendances externes.

### Intégration Prometheus (si --observability=advanced)

```go
// Dans health_handler.go, conditionnel (build tags ou runtime check)
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
    start := time.Now()
    // ... check DB ...
    duration := time.Since(start).Seconds()

    // Métriques (seulement si Prometheus activé)
    if h.metrics != nil {
        h.metrics.HealthCheckDuration.WithLabelValues("database").Observe(duration)
        if checks["database"] == "ok" {
            h.metrics.HealthCheckStatus.WithLabelValues("database").Set(1)
        } else {
            h.metrics.HealthCheckStatus.WithLabelValues("database").Set(0)
        }
    }
    // ...
}
```

**Pattern recommandé** : Utiliser une interface `MetricsProvider` optionnelle dans `HealthHandler` plutôt que des build tags, pour éviter la complexité.

### Handler Existant vs Nouveau

**Vérifier d'abord** le template actuel dans `templates.go` pour voir si un handler `/health` existe déjà.

Si oui : remplacer entièrement par `HealthHandler`
Si non : créer `health_handler.go` et l'enregistrer via fx

Ne pas avoir deux handlers gérant `/health` en même temps.

### Convention de Réponse Cohérente avec le Reste de l'API

Le format health check est **volontairement différent** du format API standard (`{status, data, meta}`) car il suit les conventions K8s/infrastructure. C'est une exception documentée et acceptée.

### Anti-Patterns à Éviter

- NE PAS faire de requêtes DB complexes dans `/health/liveness` — juste vérifier si l'app tourne
- NE PAS utiliser un timeout trop long dans `/health/readiness` — max 2-3 secondes
- NE PAS protéger `/health/*` avec le middleware JWT (K8s doit pouvoir appeler sans token)
- NE PAS casser l'endpoint `/health` existant (AC5 — rétrocompatibilité)
- NE PAS retourner 200 si la DB est vraiment down pour `/health/readiness` (K8s retirera le pod du load balancer)
- NE PAS logger chaque health check (trop verbeux en production — K8s appelle toutes les 10-15 secondes)

### Project Structure Notes

- Les health checks avancés sont générés POUR TOUS LES PROJETS (pas conditionnel à `--observability`)
- L'intégration Prometheus dans les health checks est conditionnelle à `--observability=advanced`
- Le fichier `deployments/kubernetes/probes.yaml` est généré inconditionnellement
- Vérifier que `deployments/` existe déjà dans le générateur (il y a déjà Docker Compose)

### References

- [Source: cmd/create-go-starter/templates.go] - Template server.go existant avec route /health
- [Source: cmd/create-go-starter/generator.go] - Pattern génération conditionnelle
- [Source: cmd/create-go-starter/templates_observability.go] - Créé en Story 9.1 (Prometheus metrics)
- [Source: _bmad-output/planning-artifacts/architecture.md#API & Communication Patterns] - Routing Fiber
- [Source: _bmad-output/planning-artifacts/epics.md#Story 9.3] - Story specification
- [Source: _bmad-output/project-context.md] - Rules: JWT routes, no globals, fx DI
- [K8s Probes: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/]
- [GORM DB.DB(): https://gorm.io/docs/connecting_to_the_database.html#Connection-Pool]

## Change Log

- 2026-02-17: **Code Review** — 5 constats identifiés et corrigés : (1) `printSuccessMessage()` affichait l'ancien format `{"status":"ok"}` → corrigé avec nouveau format + endpoints avancés, (2) `ReadmeTemplate()`, `QuickStartTemplate()`, `DocsReadmeTemplate()` ne documentaient pas `/health/liveness` et `/health/readiness` → corrigés, (3) AC3 spécifiait `"connection refused"` mais code hardcode `"database connection failed"` par sécurité → AC corrigé pour refléter le choix, (4) Dockerfile HEALTHCHECK utilisait `/health` au lieu de `/health/liveness` → corrigé, (5) File List ne mentionnait pas `main.go` → corrigé. Build/vet/tests passent après corrections.
- 2026-02-17: Implémentation Story 9.3 — Health checks avancés K8s générés pour template full. Remplacement du simple handler `/health` par `/health/liveness` + `/health/readiness` avec injection fx + `*gorm.DB`. Ajout `deployments/kubernetes/probes.yaml`. Intégration Prometheus pour `--observability=advanced`. 7 nouveaux tests. Documentation mise à jour.

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

Aucun blocage majeur. Décisions techniques prises :
1. `AdvancedHealthHandlerTemplate()` ajoutée sur `ProjectTemplates` dans `templates.go` (pas `ObservabilityTemplates`) car elle s'applique à tous les projets full template.
2. Ancien `internal/adapters/http/health.go` supprimé du template full — remplacé par `internal/adapters/handlers/health_handler.go` avec injection fx.
3. `RoutesTemplate()` mis à jour : plus de `RegisterHealthRoutes(app)`, injection du `*handlers.HealthHandler` via fx.
4. `PrometheusMetrics` enrichi avec `HealthCheckDuration` et `HealthCheckStatus` pour observability=advanced.
5. `KubernetesProbesTemplate()` sur `ObservabilityTemplates` — réutilisée pour tous les projets full.

### Completion Notes List

Toutes les ACs satisfaites :
- ✅ AC1: `/health/liveness` retourne 200 avec `{"status":"alive","service":"...","timestamp":"..."}`
- ✅ AC2: `/health/readiness` retourne 200 avec `{"checks":{"database":"ok"}}` si DB accessible
- ✅ AC3: `/health/readiness` retourne 503 avec `{"checks":{"database":"error"}}` si DB down
- ✅ AC4: `health_check_duration_seconds` Histogram et `health_check_status` Gauge exposés dans `/metrics` pour `--observability=advanced`
- ✅ AC5: `GET /health` alias de liveness — rétrocompatibilité maintenue
- ✅ AC6: `deployments/kubernetes/probes.yaml` généré pour tous les projets full
- ✅ AC7: Le projet généré compile (`go build ./...` — vérifié par tests unitaires)
- ✅ AC8: `go test -short ./cmd/create-go-starter` — tous les tests passent (4.3s)

### File List

Fichiers modifiés/créés dans le générateur :
- `cmd/create-go-starter/templates.go` — ajout `AdvancedHealthHandlerTemplate()` sur `ProjectTemplates`; mise à jour `ReadmeTemplate()`, `QuickStartTemplate()`, `DocsReadmeTemplate()`, `DockerfileTemplate()` pour les nouveaux endpoints health
- `cmd/create-go-starter/main.go` — mise à jour `printSuccessMessage()` avec le nouveau format de réponse health et les endpoints avancés
- `cmd/create-go-starter/templates_user.go` — `HandlerModuleTemplate()` et `RoutesTemplate()` mis à jour
- `cmd/create-go-starter/templates_observability.go` — `PrometheusTemplate()` enrichi; ajout `HealthHandlerWithMetricsTemplate()` et `KubernetesProbesTemplate()`; `GoModTemplateWithObservability()` avec `prometheus/client_golang v1.22.0`
- `cmd/create-go-starter/generator.go` — remplacement `health.go` → `health_handler.go`; ajout `probes.yaml`
- `cmd/create-go-starter/generator_test.go` — fichiers attendus mis à jour
- `cmd/create-go-starter/templates_observability_test.go` — 7 nouveaux tests health checks

Documentation mise à jour :
- `docs/usage.md` — endpoints health mis à jour
- `docs/generated-project-guide.md` — section health checks K8s et K8s probes config

Fichiers générés dans les projets créés :
- `internal/adapters/handlers/health_handler.go` (nouveau — remplace `internal/adapters/http/health.go`)
- `internal/adapters/http/routes.go` (modifié — 3 routes health, injection HealthHandler)
- `internal/adapters/handlers/module.go` (modifié — `fx.Provide(NewHealthHandler)`)
- `deployments/kubernetes/probes.yaml` (nouveau)
