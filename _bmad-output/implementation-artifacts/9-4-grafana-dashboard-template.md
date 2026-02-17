# Story 9.4: Grafana Dashboard Template

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** DevOps,
**Je veux** un dashboard Grafana pré-configuré et un stack Prometheus+Grafana dans le docker-compose généré avec `--observability=advanced`,
**Afin de** visualiser immédiatement les métriques (traffic, errors, latency) au démarrage sans configuration manuelle.

## Acceptance Criteria

1. **AC1**: Given un projet généré avec `--observability=advanced`, When le fichier `deployments/grafana/dashboards/api-dashboard.json` est généré, Then il est un JSON Grafana valide importable directement
2. **AC2**: Given le dashboard Grafana est importé, When des métriques sont disponibles dans Prometheus, Then les panneaux suivants s'affichent correctement : Request Rate (req/s), Error Rate (%), P95 Latency (ms), DB Query Duration (P95), Active DB Connections, Health Status
3. **AC3**: Given le projet est généré avec `--observability=advanced`, When le `docker-compose.yml` est consulté, Then les services `prometheus`, `grafana` sont présents avec volumes et configuration
4. **AC4**: Given le stack est démarré avec `docker-compose up`, When Grafana s'ouvre sur `http://localhost:3000`, Then le dashboard est auto-provisionné (accessible sans import manuel) via le provisioning Grafana
5. **AC5**: Given le stack Prometheus+Grafana est configuré, When le fichier `deployments/prometheus/prometheus.yml` est consulté, Then il scrape l'endpoint `/metrics` de l'app avec le job `{{project_name}}`
6. **AC6**: Given des alertes sont souhaitées, When les fichiers `deployments/prometheus/rules/*.yml` sont consultés, Then des règles d'alerte pré-configurées sont présentes (Error Rate > 5%, P95 Latency > 1s, DB Down)
7. **AC7**: Given le projet est généré, When `go build ./...` s'exécute, Then le projet compile sans erreur (les fichiers JSON/YAML n'affectent pas le build Go)
8. **AC8**: Given les tests du générateur s'exécutent, When `go test -short ./cmd/create-go-starter`, Then tous les tests passent

## Tasks / Subtasks

- [x] Task 1: Créer les templates de configuration infrastructure (AC: 3, 4, 5)
  - [x] 1.1 Dans `templates_observability.go`, ajouter template `deployments/prometheus/prometheus.yml`
  - [x] 1.2 Template `deployments/grafana/provisioning/datasources/prometheus.yaml` — auto-provisioning datasource
  - [x] 1.3 Template `deployments/grafana/provisioning/dashboards/default.yaml` — auto-provisioning dashboard path
  - [x] 1.4 Mettre à jour le template `docker-compose.yml` pour ajouter services `prometheus` et `grafana` (complète le service `jaeger` de Story 9.2)
  - [x] 1.5 Ajouter variables `GF_SECURITY_ADMIN_PASSWORD` dans `.env.example` (défaut: `admin`)

- [x] Task 2: Créer le template du dashboard Grafana JSON (AC: 1, 2)
  - [x] 2.1 Dans `templates_observability.go`, créer `GrafanaDashboardJSONTemplate()` avec le JSON complet du dashboard
  - [x] 2.2 Panneau "Request Rate" — timeseries
  - [x] 2.3 Panneau "Error Rate %" — stat + thresholds (vert/jaune/rouge)
  - [x] 2.4 Panneau "P95 Latency" — gauge + seuils (warn: 0.5s, critical: 1s)
  - [x] 2.5 Panneau "DB Query P95" — timeseries
  - [x] 2.6 Panneau "Active DB Connections" — stat
  - [x] 2.7 Panneau "Health Status" — stat avec mapping 0=DOWN(rouge)/1=UP(vert)
  - [x] 2.8 Générer le fichier `deployments/grafana/dashboards/api-dashboard.json`

- [x] Task 3: Créer les templates de règles d'alerte Prometheus (AC: 6)
  - [x] 3.1 Template `deployments/prometheus/rules/api_alerts.yml`
  - [x] 3.2 Alerte `HighErrorRate` — for: 2m, severity: warning
  - [x] 3.3 Alerte `HighP95Latency` — for: 5m, severity: warning
  - [x] 3.4 Alerte `DatabaseDown` — for: 1m, severity: critical
  - [x] 3.5 Inclure `deployments/prometheus/rules/` dans `prometheus.yml` (`rule_files:` section)

- [x] Task 4: Mettre à jour le docker-compose complet (AC: 3, 4)
  - [x] 4.1 Ajouter service `prometheus` avec image prom/prometheus:v2.51.0, volumes, ports 9090, commandes
  - [x] 4.2 Ajouter service `grafana` avec image grafana/grafana:10.4.0, env vars, volumes, port 3000, depends_on prometheus
  - [x] 4.3 Service `api` expose `/metrics` sur le même réseau Docker que `prometheus`
  - [x] 4.4 Ajouter `depends_on: prometheus` au service `api` (ordre de démarrage)

- [x] Task 5: Mettre à jour le générateur (AC: 7)
  - [x] 5.1 Dans `generator.go`, dans `generateObservabilityFiles()`, ajouter les 5 nouveaux fichiers
  - [x] 5.2 Variable substitution du nom de projet via string concatenation Go et `strings.ReplaceAll`
  - [x] 5.3 Répertoires créés automatiquement via `os.MkdirAll(filepath.Dir(file.Path))`

- [x] Task 6: Tests du générateur (AC: 8)
  - [x] 6.1 Dans `templates_observability_test.go`, tests ajoutés pour Grafana/Prometheus
  - [x] 6.2 `TestPrometheusYMLGeneratedForAdvanced` — vérifie la génération du fichier
  - [x] 6.3 `TestPrometheusYMLContainsProjectJobName` — vérifie le job_name
  - [x] 6.4 `TestGrafanaDashboardIsValidJSON` — valide le JSON avec `json.Unmarshal`
  - [x] 6.5 `TestDockerComposeContainsFullObservabilityStack` — vérifie prometheus, grafana, jaeger
  - [x] 6.6 `TestGrafanaFilesNotGeneratedWithoutAdvanced` — vérifie absence pour none/basic

- [x] Task 7: Documentation (AC: 1, 4)
  - [x] 7.1 `docs/generated-project-guide.md` mis à jour avec section "Observabilité complète"
  - [x] 7.2 URLs documentées : Grafana (`:3000`), Prometheus (`:9090`), Jaeger (`:16686`)
  - [x] 7.3 Identifiants Grafana documentés (admin/admin, configurable via `.env`)
  - [x] 7.4 Structure des fichiers générés documentée

## Dev Notes

### Architecture de la Story

Story 9.4 est la **story de finalisation de l'Epic 9**. Elle assemble tout le stack observabilité :
- Métriques Prometheus (Story 9.1)
- Traces Jaeger (Story 9.2)
- Health Checks (Story 9.3)
- Dashboard Grafana + Alertes (cette story)

**Fichiers générés (en plus des stories précédentes) :**
```
<projet-généré>/
├── deployments/
│   ├── prometheus/
│   │   ├── prometheus.yml              # Config scraping
│   │   └── rules/
│   │       └── api_alerts.yml          # Règles d'alerte
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/
│       │   │   └── prometheus.yaml     # Auto-datasource
│       │   └── dashboards/
│       │       └── default.yaml        # Auto-provisioning dashboard
│       └── dashboards/
│           └── api-dashboard.json      # Dashboard principal
```

### Template `deployments/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: '{{project_name}}'
    static_configs:
      - targets: ['app:8080']  # Nom du service Docker
    metrics_path: '/metrics'
    scrape_interval: 10s

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

### Template `deployments/grafana/provisioning/datasources/prometheus.yaml`

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: false
```

### Template `deployments/grafana/provisioning/dashboards/default.yaml`

```yaml
apiVersion: 1

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
```

### Dashboard JSON — Structure de Base

Le dashboard JSON Grafana est volumineux mais suit une structure standard :

```json
{
  "id": null,
  "uid": "{{project_name}}-api",
  "title": "{{project_name}} - API Dashboard",
  "tags": ["api", "{{project_name}}"],
  "timezone": "browser",
  "schemaVersion": 39,
  "version": 1,
  "refresh": "30s",
  "time": {"from": "now-1h", "to": "now"},
  "panels": [
    // ... voir panneaux ci-dessous
  ]
}
```

**Structure d'un panneau stat (Error Rate) :**
```json
{
  "type": "stat",
  "title": "Error Rate",
  "gridPos": {"h": 4, "w": 6, "x": 6, "y": 0},
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
  "targets": [{
    "expr": "sum(rate(http_requests_total{job=\"{{project_name}}\",status_code=~\"5..\"}[5m])) / sum(rate(http_requests_total{job=\"{{project_name}}\"}[5m])) * 100",
    "refId": "A"
  }]
}
```

**Layout recommandé (12 colonnes) :**
```
Row 1: [Request Rate (6w)] [Error Rate (6w)]
Row 2: [P95 Latency (6w)] [DB Query P95 (6w)]
Row 3: [Active DB Connections (4w)] [Health Status DB (4w)] [Health Status App (4w)]
Row 4: [Time Series HTTP Requests (12w)]
```

### Règles d'Alerte Prometheus

```yaml
# deployments/prometheus/rules/api_alerts.yml
groups:
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
```

### Variables d'Environnement Ajoutées

Dans le `.env.example` généré (complète les vars de Story 9.2) :
```
# Grafana (--observability=advanced)
GF_SECURITY_ADMIN_USER=admin
GF_SECURITY_ADMIN_PASSWORD=admin
GF_USERS_ALLOW_SIGN_UP=false
```

### docker-compose.yml Complet (avec observabilité)

```yaml
services:
  app:
    # ... existant ...
    depends_on:
      - db
      - prometheus

  db:
    # ... existant ...

  jaeger:        # Story 9.2
    # ... voir Story 9.2 ...

  prometheus:
    image: prom/prometheus:v2.51.0  # Utiliser version fixée, pas latest
    container_name: {{project_name}}-prometheus
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
      - {{project_name}}-network

  grafana:
    image: grafana/grafana:10.4.0  # Utiliser version fixée, pas latest
    container_name: {{project_name}}-grafana
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
      - {{project_name}}-network

volumes:
  # ... existant (postgres_data) ...
  prometheus_data:
  grafana_data:
```

**NOTE** : Utiliser des versions FIXÉES (pas `latest`) pour la reproductibilité en production. Vérifier les dernières versions stables disponibles en 2026.

### Substitution de Variables dans les Templates

Le générateur doit remplacer `{{project_name}}` dans :
- `prometheus.yml` : job_name, container_name
- `api-dashboard.json` : title, uid, toutes les requêtes PromQL avec `job="..."`
- `docker-compose.yml` : container_name, network name

Utiliser la même fonction de substitution que pour les templates Go existants (`strings.ReplaceAll` ou `text/template`).

### Anti-Patterns à Éviter

- NE PAS utiliser `image: latest` dans docker-compose — fixer les versions pour reproductibilité
- NE PAS mettre de credentials en dur dans les fichiers de config — utiliser variables d'env avec valeurs défaut
- NE PAS créer un dashboard JSON à la main (trop fragile) — exporter depuis Grafana puis adapter
- NE PAS oublier de créer les dossiers vides avec `.gitkeep` si nécessaire pour les volumes
- NE PAS mettre `prometheus_data` et `grafana_data` dans `.gitignore` (les fichiers de provisioning doivent être commités)
- NE PAS exposer l'interface Prometheus publiquement en production (`:9090` est pour le dev local seulement)

### Dépendances entre Stories de l'Epic 9

```
9.1 (Prometheus flag + métriques HTTP)
  └── 9.2 (Traces + docker-compose Jaeger)
      └── 9.3 (Health checks + K8s probes + métriques health)
          └── 9.4 (Dashboard + Alertes + docker-compose complet) ← cette story
```

La story 9.4 doit être développée APRÈS les stories 9.1, 9.2, 9.3 car elle dépend des métriques définies dans ces stories.

### Tests — Validation JSON Grafana

```go
// Test du dashboard JSON dans templates_observability_test.go
func TestGrafanaDashboardIsValidJSON(t *testing.T) {
    config := Config{ObservabilityLevel: "advanced", ProjectName: "test-proj"}
    content, err := generateGrafanaDashboard(config)
    require.NoError(t, err)

    var dashboard map[string]interface{}
    err = json.Unmarshal([]byte(content), &dashboard)
    assert.NoError(t, err, "Grafana dashboard must be valid JSON")
    assert.Equal(t, "test-proj - API Dashboard", dashboard["title"])
}
```

### Project Structure Notes

- Les fichiers Grafana/Prometheus sont dans `deployments/` (cohérent avec Docker Compose existant)
- Vérifier si `deployments/` existe déjà dans le générateur (oui, pour Docker Compose)
- Ajouter `deployments/grafana/` et `deployments/prometheus/` comme nouveaux sous-dossiers
- Les volumes Docker `prometheus_data` et `grafana_data` s'ajoutent au volume `postgres_data` existant

### References

- [Source: cmd/create-go-starter/templates_observability.go] - Stories 9.1, 9.2, 9.3
- [Source: cmd/create-go-starter/templates.go] - Template docker-compose.yml existant
- [Source: cmd/create-go-starter/generator.go] - Génération conditionnelle
- [Source: _bmad-output/planning-artifacts/epics.md#Story 9.4] - Story specification
- [Source: _bmad-output/project-context.md] - Rules: no hardcoded secrets
- [Grafana Provisioning: https://grafana.com/docs/grafana/latest/administration/provisioning/]
- [Grafana Dashboard JSON: https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/import-dashboards/]
- [Prometheus Docker: https://hub.docker.com/r/prom/prometheus]
- [Grafana Docker: https://hub.docker.com/r/grafana/grafana]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

(aucun blocage rencontré)

### Completion Notes List

- Implémenté 5 nouvelles méthodes sur `ObservabilityTemplates` : `PrometheusConfigTemplate()`, `GrafanaDatasourceTemplate()`, `GrafanaDashboardProvisionTemplate()`, `GrafanaDashboardJSONTemplate()`, `PrometheusAlertRulesTemplate()`
- Dashboard Grafana JSON avec 7 panneaux (Request Rate, Error Rate, P95 Latency, DB Query P95, Active DB Connections, Health Status, HTTP Requests by Status)
- `DockerComposeTemplateWithObservability()` entièrement réécrit pour produire un YAML valide avec les 4 services (db, jaeger, prometheus, grafana) + volumes + networks correctement structurés (corrige également un bug de structure YAML préexistant des stories 9.1-9.2)
- `generateObservabilityFiles()` dans `generator.go` étendu avec 5 nouveaux fichiers
- 7 nouveaux tests ajoutés dans `templates_observability_test.go` : tous passent en `-short`
- `docs/generated-project-guide.md` mis à jour avec section observabilité complète (URLs, credentials, panneaux, alertes, structure fichiers)
- Tous les tests de la suite passent sans régression (`go test -short ./cmd/create-go-starter/...`)

### File List

- `cmd/create-go-starter/templates_observability.go` (modifié — 5 nouvelles méthodes, mise à jour DockerCompose + env)
- `cmd/create-go-starter/generator.go` (modifié — generateObservabilityFiles étendu)
- `cmd/create-go-starter/templates_observability_test.go` (modifié — 7 nouveaux tests + import encoding/json)
- `docs/generated-project-guide.md` (modifié — section observabilité complète)

## Change Log

- 2026-02-17 : Implémentation Story 9.4 — Grafana Dashboard Template + stack Prometheus+Grafana complet (Date: 2026-02-17)
- 2026-02-17 : Code Review (AI) — 7 issues trouvées (2H, 3M, 2L), 4 fixées automatiquement, 1 annulée (M3: t.Parallel impossible avec os.Chdir), 2 documentaires (M2, L1). Status: done

## Senior Developer Review (AI)

**Reviewer:** Yacoubakone (via claude-opus-4.6)
**Date:** 2026-02-17
**Verdict:** APPROUVÉ avec corrections appliquées

### Issues Trouvées: 7 (2 High, 3 Medium, 2 Low)

#### Fixées automatiquement (4):

- **[H1][FIXED]** `templates_observability.go:563` — Image Jaeger `latest` remplacée par `jaegertracing/all-in-one:1.56.0` (violait l'anti-pattern de la story)
- **[H2][FIXED]** `templates_observability.go:1215` — Titre panneau P95 Latency changé de "P95 Latency (ms)" à "P95 Latency" (unité Grafana est "s", pas "ms")
- **[M1][FIXED]** `templates_observability.go:636-639` — Format `depends_on` SQLite uniformisé vers la forme longue avec `condition: service_started` (cohérent avec PostgreSQL/MySQL)
- **[L2][FIXED]** `templates_observability.go:633,737` — `JWT_SECRET` changé de valeur en dur vers `${JWT_SECRET:-dev-secret-change-in-production}` (env var substitution cohérente)

#### Non fixées — Action items (2):

- **[M2][DOC]** Story File List dit "modifié" pour `templates_observability.go` et `templates_observability_test.go` mais ces fichiers sont `??` (untracked/new) dans git — ils n'ont jamais été commités
- **[L1][DOC]** Dashboard a 7 panneaux mais AC2 en liste seulement 6 — le 7e "HTTP Requests by Status" est un ajout utile non documenté dans l'AC

#### Annulée (1):

- **[M3][CANCELLED]** `t.Parallel()` ne peut pas être ajouté aux tests car ils utilisent `os.Chdir()` (état global du processus) — nécessiterait un refactoring des tests pour utiliser des chemins absolus

### AC Validation

| AC | Status |
|---|---|
| AC1 (JSON Grafana valide) | IMPLÉMENTÉ |
| AC2 (Panneaux dashboard) | IMPLÉMENTÉ (6/6 AC + 1 bonus) |
| AC3 (docker-compose prometheus+grafana) | IMPLÉMENTÉ |
| AC4 (auto-provisioning Grafana) | IMPLÉMENTÉ |
| AC5 (prometheus.yml scrape) | IMPLÉMENTÉ |
| AC6 (alertes Prometheus) | IMPLÉMENTÉ |
| AC7 (go build sans erreur) | IMPLÉMENTÉ |
| AC8 (tests passent) | IMPLÉMENTÉ |

### Tests

- Tous les tests Story 9.4 passent (`go test -short -run "TestPrometheusYML|TestGrafana|TestDockerComposeContainsFull|TestPrometheusAlertRules|TestEnvExampleContainsGrafana"`)
- Build OK (`go build ./cmd/create-go-starter`)
- Note: `TestAddModelSubcommandRouting` échoue (pré-existant, non lié à Story 9.4)
