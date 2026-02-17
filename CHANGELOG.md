# Changelog

Toutes les modifications notables de ce projet seront documentées dans ce fichier.

Le format est basé sur [Keep a Changelog](https://keepachangelog.com/fr/1.0.0/),
et ce projet adhère au [Semantic Versioning](https://semver.org/lang/fr/).

## [1.3.0] - 2026-02-17

### <i class="material-icons success">new_releases</i> Observabilité Avancée (Epic 9)

Stack d'observabilité complète pour les projets générés en production: métriques Prometheus, distributed tracing OpenTelemetry/Jaeger, health checks Kubernetes, et dashboard Grafana pré-configuré.

### <i class="material-icons success">star</i> Fonctionnalités

#### Story 9.1 — Endpoint Prometheus Metrics

- **Flag `--observability`** avec 3 niveaux: `none` (défaut), `basic`, `advanced`
- **Endpoint `/metrics`** compatible Prometheus via `fiberprometheus/v2 v2.7.0`
- **Métriques HTTP exposées**:
  - `http_requests_total` (Counter) — requêtes totales par méthode, route, status
  - `http_request_duration_seconds` (Histogram) — latence HTTP par route
  - `http_requests_in_flight` (Gauge) — requêtes actives en cours
- **Middleware metrics** — capture automatique des métriques sur toutes les routes
- **Validation CLI** — `--observability=advanced` requiert `--template=full`

#### Story 9.2 — Distributed Tracing (OpenTelemetry + Jaeger)

- **OpenTelemetry SDK** avec exporter OTLP/gRPC vers Jaeger
- **Propagation W3C traceparent** — correlation entre services
- **Middleware tracing** — span automatique par requête HTTP avec attributs (method, route, status)
- **GORM tracing** — spans automatiques pour les requêtes SQL (query, table, rows_affected)
- **Logger enrichi** — `trace_id` et `span_id` injectés dans chaque log zerolog
- **Service Jaeger** — ajouté au Docker Compose (jaeger:1.56.0) avec UI sur port 16686
- **Variable d'environnement** — `OTEL_EXPORTER_OTLP_ENDPOINT` dans `.env.example`

#### Story 9.3 — Health Checks Avancés (Kubernetes-ready)

- **Endpoints dédiés**:
  - `GET /health/liveness` — toujours 200 si l'application tourne (K8s liveness probe)
  - `GET /health/readiness` — vérifie la connexion DB avec timeout 2s, retourne 503 si down (K8s readiness probe)
  - `GET /health` — alias vers liveness (rétrocompatibilité)
- **HealthHandler** avec injection `*gorm.DB` via fx
- **Kubernetes probes** — `deployments/kubernetes/probes.yaml` généré automatiquement
- **Métriques health** — `health_check_status` exposé sur `/metrics` quand `--observability=advanced`

#### Story 9.4 — Dashboard Grafana pré-configuré

- **Dashboard JSON 7 panneaux** — Request Rate, Error Rate, Latency (p50/p95/p99), Active Requests, Health Status, DB Latency, Top Endpoints
- **Auto-provisioning Grafana** — datasources + dashboards YAML configurés automatiquement
- **Configuration Prometheus** — `prometheus.yml` avec scrape config et alert rules
- **Règles d'alerte** — latence élevée (>500ms), taux d'erreur (>5%), saturation
- **Docker Compose étendu** — stack complète:
  - PostgreSQL (base de données)
  - Jaeger 1.56.0 (tracing UI: port 16686)
  - Prometheus v2.51.0 (métriques: port 9090)
  - Grafana 10.4.0 (dashboards: port 3000)

### <i class="material-icons">build</i> Architecture

#### Nouveau fichier CLI

- `cmd/create-go-starter/templates_observability.go` — **~1369 lignes**, 15+ fonctions de templates:
  1. `PrometheusTemplate()` — Registry et métriques Prometheus
  2. `MetricsMiddlewareTemplate()` — Middleware capture HTTP metrics
  3. `MetricsHandlerTemplate()` — Handler GET /metrics
  4. `TracerTemplate()` — Initialisation OpenTelemetry SDK
  5. `TracingMiddlewareTemplate()` — Middleware spans HTTP
  6. `GORMTracingTemplate()` — Plugin GORM pour spans SQL
  7. `LoggerWithTracingTemplate()` — zerolog enrichi avec trace_id/span_id
  8. `AdvancedHealthHandlerTemplate()` — Liveness + Readiness endpoints
  9. `KubernetesProbesTemplate()` — probes.yaml pour K8s
  10. `GrafanaDashboardJSONTemplate()` — Dashboard 7 panneaux
  11. `PrometheusConfigTemplate()` — prometheus.yml
  12. `PrometheusAlertRulesTemplate()` — alert.rules.yml
  13. `GrafanaDatasourceTemplate()` — Provisioning datasources
  14. `GrafanaDashboardProvisioningTemplate()` — Provisioning dashboards
  15. `ObservabilityDockerComposeTemplate()` — Docker Compose étendu

#### Fichiers modifiés

- `main.go` — ajout flag `--observability`, validation (`advanced` requiert `full`)
- `generator.go` — nouvelle fonction `generateObservabilityFiles()` orchestrant la génération
- `templates.go` — `AdvancedHealthHandlerTemplate()`, Dockerfile mis à jour
- `templates_user.go` — `HandlerModuleTemplate()` et `RoutesTemplate()` adaptés pour health checks avancés

### <i class="material-icons success">check_circle</i> Tests

- **Nouveaux tests** dans `templates_observability_test.go`
- **Tests mis à jour**: main_test.go, generator_test.go, smoke_test.go, template_minimal_test.go
- **Tests E2E** mis à jour: e2e_mysql_test.go, e2e_sqlite_test.go, database_integration_test.go
- **Tous les tests passent** — `go test ./...` (36.672s, full suite)

### <i class="material-icons">description</i> Documentation

- Guide monitoring complet dans `docs/guide/monitoring.md`
- Section observabilité dans `docs/usage.md`
- Architecture templates_observability.go dans `docs/cli-architecture.md`
- Exemples et configuration dans `docs/generated-project-guide.md`

### <i class="material-icons info">info</i> Améliorations

- **Production-ready** — Stack observabilité complète en une commande
- **Zero configuration** — Docker Compose avec tous les services pré-configurés
- **Kubernetes-native** — Probes et métriques compatibles K8s
- **Rétrocompatibilité** — `/health` préservé comme alias, `--observability=none` par défaut

### <i class="material-icons warning">warning</i> Limitations connues

- **Template full uniquement** — `--observability=advanced` requiert `--template=full`
- **Templates minimal et graphql** — Observabilité non supportée (sanity check CLI)
- **Jaeger uniquement** — Export OTLP vers d'autres backends (Zipkin, Tempo) nécessite configuration manuelle
- **Dashboard Grafana** — Pré-configuré pour Prometheus, personnalisation manuelle pour métriques custom

### <i class="material-icons">link</i> Liens

- **Documentation**: https://tky0065.github.io/go-starter-kit/guide/monitoring/
- **Release**: [v1.3.0](https://github.com/tky0065/go-starter-kit/releases/tag/v1.3.0)

## [1.2.0] - 2026-02-12

### <i class="material-icons success">new_releases</i> Générateur CRUD Scaffolding

Nouvelle commande `add-model` pour générer automatiquement des modèles CRUD complets dans les projets existants.

### <i class="material-icons success">star</i> Fonctionnalités

#### Commande add-model

- **Génération automatique de modèles CRUD** via `create-go-starter add-model <ModelName> --fields "field:type,..."`
- **8 fichiers générés automatiquement** par modèle:
  - `internal/models/<model>.go` - Model avec tags GORM
  - `internal/interfaces/<model>_repository.go` - Interface repository
  - `internal/adapters/repository/<model>_repository.go` - Implémentation GORM
  - `internal/domain/<model>/service.go` - Logique métier CRUD
  - `internal/domain/<model>/module.go` - Module fx pour injection
  - `internal/adapters/handlers/<model>_handler.go` - Handlers HTTP REST
  - `internal/domain/<model>/service_test.go` - Tests unitaires service
  - `internal/adapters/handlers/<model>_handler_test.go` - Tests handlers
- **Mise à jour automatique de 3 fichiers existants**:
  - `internal/infrastructure/database/database.go` - Migration auto du modèle
  - `internal/adapters/http/routes.go` - Routes CRUD REST
  - `cmd/main.go` - Enregistrement du module fx

#### Types de champs supportés

- **Types primitifs**: `string`, `int`, `uint`, `float64`, `bool`, `time`
- **Modificateurs GORM**: `unique`, `not_null`, `index`
- **Syntaxe**: `field:type:modifier1:modifier2`

**Exemples**:
```bash
create-go-starter add-model Todo --fields "title:string,completed:bool"
create-go-starter add-model Product --fields "name:string:unique:not_null,price:float64,stock:int:index"
```

#### Relations entre modèles

- **BelongsTo (N:1)**: Ajoute foreign key et relation au modèle enfant
  ```bash
  create-go-starter add-model Comment --fields "content:string" --belongs-to Post
  ```
  - Génère: `PostID uint` + `Post Post` dans Comment
  - Endpoints: `GET /posts/:postId/comments`, `POST /posts/:postId/comments`
  - Preloading: `GET /comments/:id?include=post`

- **HasMany (1:N)**: Ajoute slice d'enfants au modèle parent
  ```bash
  create-go-starter add-model Category --fields "name:string" --has-many Product
  ```
  - Génère: `Products []Product` dans Category
  - Preloading: `GET /categories/:id?include=products`

- **Support relations imbriquées**: Category → Post → Comment (3+ niveaux)

#### Flags disponibles

- `--fields` (requis): Définition des champs
- `--belongs-to <Model>`: Ajoute relation BelongsTo
- `--has-many <Model>`: Ajoute relation HasMany
- `--public`: Routes sans authentification (pas de middleware JWT)
- `--yes`, `-y`: Skip confirmation prompt
- `--help`, `-h`: Afficher l'aide

#### Endpoints générés automatiquement

Pour chaque modèle (exemple: Todo), les endpoints suivants sont générés:

- `GET /api/v1/todos` - Liste tous les todos (pagination support)
- `GET /api/v1/todos/:id` - Récupère un todo par ID
- `POST /api/v1/todos` - Crée un nouveau todo
- `PUT /api/v1/todos/:id` - Met à jour un todo
- `DELETE /api/v1/todos/:id` - Supprime un todo (soft delete)

**Endpoints relationnels** (si `--belongs-to Parent`):
- `GET /api/v1/parents/:parentId/todos` - Liste des enfants d'un parent
- `POST /api/v1/parents/:parentId/todos` - Crée enfant lié au parent

**Query preloading**: `?include=relation1,relation2` pour charger les relations

#### Tests automatiques

Chaque modèle généré inclut:
- **Tests service**: Create, GetByID, Update, Delete, List
- **Tests handler**: HTTP endpoints avec fiber test utilities
- **Couverture**: ~80%+ du code généré

### <i class="material-icons">build</i> Architecture

#### Nouveaux fichiers CLI

- `cmd/create-go-starter/add_model.go` - CLI add-model orchestration
- `cmd/create-go-starter/model.go` - Structures de données (Field, Model, Relation)
- `cmd/create-go-starter/model_generator.go` - 13 générateurs de templates:
  1. `GenerateModel()` - Entity avec GORM tags
  2. `GenerateRepositoryInterface()` - Port d'accès données
  3. `GenerateRepositoryImpl()` - Implémentation GORM
  4. `GenerateService()` - Logique métier CRUD
  5. `GenerateServiceModule()` - Module fx
  6. `GenerateHandler()` - HTTP handlers REST
  7. `GenerateServiceTests()` - Tests unitaires service
  8. `GenerateHandlerTests()` - Tests HTTP handlers
  9. `UpdateDatabaseMigrations()` - AutoMigrate dans database.go
  10. `UpdateRoutes()` - Routes dans routes.go
  11. `UpdateMainModule()` - Module fx dans main.go
  12. `UpdateParentModel()` - Ajoute HasMany au parent
  13. `AddBelongsToImport()` - Import parent dans model enfant

#### Patterns

- **Pluralisation intelligente**: Todo→Todos, Category→Categories
- **Nommage cohérent**: PascalCase (types), camelCase (variables), snake_case (DB)
- **Gestion des relations**: Détection automatique des modèles existants
- **Validation stricte**: Vérifie que les modèles parents existent avant relations

### <i class="material-icons success">check_circle</i> Tests

- **Couverture**: 85%+ (ajout de 15+ tests pour add-model)
- **Tests unitaires**: Field parsing, Model validation, Relations
- **Tests d'intégration**: Génération complète de modèles avec relations
- **Tests E2E**: Smoke tests avec génération Todo + Comment + Category→Product

### <i class="material-icons">description</i> Documentation

- Guide complet add-model dans `docs/usage.md`
- Architecture add-model dans `docs/cli-architecture.md`
- Exemples avancés dans `docs/generated-project-guide.md`
- Tutoriel complet Blog System (Category→Post→Comment)
- Conventions et limitations documentées

### <i class="material-icons info">info</i> Améliorations

- **Expérience développeur**: Génération CRUD en <2 secondes
- **Productivité**: 90%+ de code boilerplate automatisé
- **Maintenabilité**: Code généré suit les best practices Go
- **Testabilité**: Tests générés prêts à être étendus

### <i class="material-icons warning">warning</i> Limitations connues

- **Template minimal**: add-model non supporté (pas de structure internal/models/)
- **Template graphql**: add-model non supporté (architecture GraphQL différente - prévu v1.3.0)
- **Template full**: ✅ Entièrement supporté avec add-model
- **Pluralisation**: Règles simples (ajoute 's'). Pluriels irréguliers nécessitent édition manuelle
- **Relations many-to-many**: Pas encore supportées (prévu v1.3.0)
- **Validations custom**: Doivent être ajoutées manuellement après génération
- **Swagger**: Doit être regénéré manuellement (`make swagger`) après add-model
- **Flags automation**: `--dir` et `--yes` pour création projet prévus v1.2.1

### <i class="material-icons">link</i> Liens

- **Documentation**: https://tky0065.github.io/go-starter-kit/usage/#ajouter-des-modeles-add-model
- **Release**: [v1.2.0](https://github.com/tky0065/go-starter-kit/releases/tag/v1.2.0)
- **Changelog détaillé**: [Epic 8 Stories](https://github.com/tky0065/go-starter-kit/milestone/8?closed=1)

## [1.0.0] - 2026-01-15

### <i class="material-icons">rocket_launch</i> Première version stable

Premier release officiel de `create-go-starter`, un générateur CLI pour créer des projets Go prêts pour la production avec architecture hexagonale.

### <i class="material-icons success">star</i> Fonctionnalités

#### Templates de Projet
- **3 templates au choix** via le flag `--template`:
  - `minimal` - API REST basique avec Swagger (sans authentification) - ~20 fichiers
  - `full` - API complète avec JWT auth et gestion utilisateurs (défaut) - ~35 fichiers
  - `graphql` - API GraphQL avec gqlgen et GraphQL Playground - ~23 fichiers

#### Architecture & Stack Technique
- **Architecture hexagonale** (Ports & Adapters) pour séparation claire des responsabilités
- **JWT Authentication**:
  - Access tokens + Refresh tokens avec rotation sécurisée
  - Middleware de sécurisation des routes
  - Gestion de session avec renouvellement automatique
- **User CRUD**:
  - Opérations complètes (Create, Read, Update, Delete)
  - Gestion du profil utilisateur
  - Hachage sécurisé des mots de passe (bcrypt)
- **API REST** avec Fiber v2 - Framework web haute performance
- **Base de données** PostgreSQL avec GORM et migrations automatiques
- **Injection de dépendances** avec uber-go/fx
- **Logging structuré** avec rs/zerolog
- **Validation** avec go-playground/validator

#### Documentation & API
- **Swagger/OpenAPI** - Documentation auto-générée avec swaggo/swag
- **Standardisation des API** - Format de réponse uniforme
- **Gestion centralisée des erreurs** - Codes d'erreur standardisés

#### DevOps & Qualité
- **Docker**:
  - Build multi-stage optimisé
  - docker-compose pré-configuré pour dev
  - Image de production légère basée sur Alpine
- **CI/CD**:
  - Pipeline GitHub Actions pré-configuré
  - Lint automatique avec golangci-lint
  - Tests automatisés
- **Tests**:
  - Tests unitaires pour handlers, services, repositories
  - Tests d'intégration
  - Couverture de tests du CLI
  - 8 tests de résolveurs GraphQL (template graphql)
- **Makefile** avec commandes utiles (dev, test, build, docker)

#### Automatisation
- **Initialisation Git automatique** avec commit initial
- **Installation automatique des dépendances** Go (`go mod tidy`)
- **Script setup.sh** pour configuration automatique du projet
- **Documentation inline** avec GoDoc pour toutes les fonctions publiques

### <i class="material-icons success">bar_chart</i> Qualité

- <i class="material-icons success">check</i> Tests unitaires et d'intégration complets
- <i class="material-icons success">check</i> Lint avec golangci-lint
- <i class="material-icons success">check</i> Documentation complète
- <i class="material-icons success">check</i> Exemples et guides d'utilisation

### <i class="material-icons">menu_book</i> Documentation

- Guide d'installation complet
- Guide d'utilisation avec exemples
- Documentation de l'architecture CLI
- Guide du projet généré
- Quick start guide
- GitHub Pages: https://tky0065.github.io/go-starter-kit/

### <i class="material-icons">build</i> Configuration Requise

- Go 1.23+ (recommandé: 1.25.5)
- PostgreSQL 12+ (ou Docker pour exécution via conteneur)
- Git (optionnel, pour initialisation automatique)

### <i class="material-icons">install_desktop</i> Installation

```bash
# Installation globale depuis GitHub
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@v1.0.0

# Ou version latest
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

### <i class="material-icons">rocket_launch</i> Utilisation

```bash
# Template par défaut (full)
create-go-starter mon-projet

# Template minimal
create-go-starter --template=minimal mon-projet

# Template GraphQL
create-go-starter --template=graphql mon-projet
```

### <i class="material-icons warning">bug_report</i> Corrections de Bugs

- Fix: Ajout des imports manquants (`fmt`, `time`) dans le template de tests GraphQL
- Fix: Gestion correcte du flag `--template` (nécessite `--template=value` ou position avant le nom)

### <i class="material-icons">lock</i> Sécurité

- Tokens JWT sécurisés avec expiration configurable
- Refresh tokens avec rotation automatique
- Hachage bcrypt pour les mots de passe
- Validation stricte des entrées utilisateur
- Configuration des secrets via variables d'environnement

### <i class="material-icons">handyman</i> Développement

Le projet a été développé avec les meilleures pratiques de développement logiciel:

- Architecture hexagonale pour maintenabilité
- Tests automatisés pour fiabilité
- CI/CD pour déploiement continu
- Documentation complète pour faciliter l'utilisation
- Code propre et bien structuré

### 🙏 Remerciements

Merci à tous les contributeurs et aux projets open-source utilisés :
- [Fiber](https://github.com/gofiber/fiber) - Framework web
- [fx](https://github.com/uber-go/fx) - Injection de dépendances
- [GORM](https://gorm.io/) - ORM
- [zerolog](https://github.com/rs/zerolog) - Logging
- [swaggo](https://github.com/swaggo/swag) - Swagger
- [gqlgen](https://github.com/99designs/gqlgen) - GraphQL

---

## Format du Versioning

- **MAJOR** (X.0.0): Changements incompatibles avec les versions précédentes
- **MINOR** (1.X.0): Ajout de fonctionnalités rétro-compatibles
- **PATCH** (1.0.X): Corrections de bugs rétro-compatibles

[1.0.0]: https://github.com/tky0065/go-starter-kit/releases/tag/v1.0.0
