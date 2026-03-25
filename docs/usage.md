# Guide d'utilisation

Ce guide détaille l'utilisation de `create-go-starter` et explique en profondeur la structure des projets générés.

## Commande de base

La syntaxe de base est très simple:

```bash
create-go-starter <nom-du-projet>
```

### Exemple

```bash
create-go-starter mon-api-backend
```

Cette commande va créer un nouveau répertoire `mon-api-backend/` avec toute la structure du projet en utilisant le template **full** par défaut.

## Templates disponibles

`create-go-starter` propose **trois templates** pour répondre à différents besoins de projets. Choisissez le template avec le flag `--template`:

```bash
create-go-starter mon-projet --template minimal    # API REST basique
create-go-starter mon-projet --template full       # API complète avec auth (défaut)
create-go-starter mon-projet --template graphql    # API GraphQL
```

### Vue d'ensemble des templates

| Template | Description | Cas d'usage |
|----------|-------------|-------------|
| `minimal` | API REST basique avec Swagger (sans authentification) | Prototypes rapides, APIs publiques simples, microservices sans auth |
| `full` | API complète avec JWT auth, gestion utilisateurs et Swagger (**défaut**) | Applications backend complètes, APIs nécessitant authentification |
| `graphql` | API GraphQL avec gqlgen et GraphQL Playground | Applications nécessitant GraphQL, clients frontend modernes |

### Comparaison détaillée des fonctionnalités

| Fonctionnalité | minimal | full | graphql |
|----------------|---------|------|---------|
| **API REST** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons error">cancel</i> |
| **API GraphQL** | <i class="material-icons error">cancel</i> | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> |
| **Authentification JWT** | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons error">cancel</i> |
| **Gestion utilisateurs** | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Documentation Swagger** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons error">cancel</i> |
| **GraphQL Playground** | <i class="material-icons error">cancel</i> | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> |
| **Base de données (GORM)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **PostgreSQL** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Dependency Injection (fx)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Logging structuré (zerolog)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Architecture hexagonale** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Tests unitaires** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Docker** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **CI/CD (GitHub Actions)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |

### Différences structurelles majeures

#### Template `minimal`

**Caractéristiques**:
- API REST simple avec endpoints CRUD de base
- Pas d'authentification ni d'autorisation
- Swagger pour documentation API
- Parfait pour commencer rapidement sans complexité

**Structure spécifique**:
- `internal/adapters/handlers/user_handler.go` - Handlers CRUD simples sans auth
- `internal/domain/user/service.go` - Logique métier basique
- Documentation Swagger automatique via annotations

**Endpoints générés**:
```
GET    /health                  # Health check (alias liveness — rétrocompatibilité)
GET    /health/liveness         # Liveness probe K8s (toujours 200 si app tourne)
GET    /health/readiness        # Readiness probe K8s (200 si DB ok, 503 si DB down)
GET    /api/v1/users            # Liste tous les utilisateurs
GET    /api/v1/users/:id        # Récupère un utilisateur
POST   /api/v1/users            # Crée un utilisateur
PUT    /api/v1/users/:id        # Met à jour un utilisateur
DELETE /api/v1/users/:id        # Supprime un utilisateur
GET    /swagger/*               # Documentation Swagger UI
```

**Cas d'usage recommandés**:
- Prototypes et POCs rapides
- APIs publiques sans données sensibles
- Microservices internes sans besoins d'authentification
- Apprentissage de l'architecture hexagonale

---

#### Template `full` (défaut)

**Caractéristiques**:
- API REST complète avec authentification JWT
- Système d'auth avec access tokens + refresh tokens
- Gestion complète des utilisateurs (CRUD + auth)
- Swagger avec authentification Bearer token
- Production-ready avec sécurité intégrée

**Structure spécifique**:
- `internal/adapters/handlers/auth_handler.go` - Endpoints register, login, refresh
- `internal/adapters/handlers/user_handler.go` - CRUD protégé par JWT
- `internal/adapters/middleware/auth_middleware.go` - Vérification JWT
- `pkg/auth/` - Génération et validation des tokens JWT
- `internal/models/user.go` - User + RefreshToken avec bcrypt

**Endpoints générés**:
```
GET    /health                      # Health check (alias liveness — rétrocompatibilité)
GET    /health/liveness             # Liveness probe K8s
GET    /health/readiness            # Readiness probe K8s (vérifie DB)
POST   /api/v1/auth/register        # Inscription utilisateur
POST   /api/v1/auth/login           # Connexion (retourne access + refresh tokens)
POST   /api/v1/auth/refresh         # Rafraîchir l'access token
GET    /api/v1/users                # Liste utilisateurs (<i class="material-icons warning small">lock</i> JWT requis)
GET    /api/v1/users/:id            # Récupère utilisateur (<i class="material-icons warning small">lock</i> JWT requis)
PUT    /api/v1/users/:id            # Met à jour utilisateur (<i class="material-icons warning small">lock</i> JWT requis)
DELETE /api/v1/users/:id            # Supprime utilisateur (<i class="material-icons warning small">lock</i> JWT requis)
GET    /swagger/*                   # Documentation Swagger UI
```

**Cas d'usage recommandés**:
- Applications backend complètes
- APIs nécessitant authentification et autorisation
- SaaS et applications multi-utilisateurs
- APIs exposées publiquement avec données sensibles

---

#### Template `graphql`

**Caractéristiques**:
- API GraphQL complète avec gqlgen
- GraphQL Playground pour explorer l'API interactivement
- Schéma GraphQL typé avec resolvers
- Gestion des utilisateurs avec mutations et queries
- Architecture hexagonale adaptée à GraphQL

**Structure spécifique**:
- `graph/schema.graphqls` - Schéma GraphQL (types, queries, mutations)
- `graph/resolver.go` - Resolver principal
- `graph/schema.resolvers.go` - Implémentation des resolvers
- `gqlgen.yml` - Configuration gqlgen
- `internal/infrastructure/server/server.go` - Serveur GraphQL avec Playground

**Schéma GraphQL généré**:
```graphql
type User {
  id: ID!
  email: String!
  createdAt: Time!
  updatedAt: Time!
}

type Query {
  users(limit: Int, offset: Int): [User!]!
  user(id: ID!): User
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): Boolean!
}
```

**Endpoints générés**:
```
GET    /health                  # Health check
POST   /graphql                 # Endpoint GraphQL
GET    /                        # GraphQL Playground UI
```

**Cas d'usage recommandés**:
- Applications frontend modernes (React, Vue, Angular)
- APIs nécessitant des queries flexibles
- Applications mobile avec besoins de données spécifiques
- Projets privilégiant GraphQL à REST

---

### Comment choisir le bon template?

**Choisissez `minimal` si**:
- Vous voulez un prototype rapide
- Votre API est publique sans données sensibles
- Vous n'avez pas besoin d'authentification
- Vous voulez apprendre l'architecture hexagonale simplement

**Choisissez `full` si**:
- Vous construisez une application backend complète
- Vous avez besoin d'authentification JWT
- Vous voulez des utilisateurs avec login/register
- Vous préférez REST à GraphQL
- Vous voulez un projet production-ready immédiatement

**Choisissez `graphql` si**:
- Vous construisez une API GraphQL
- Votre frontend utilise Apollo Client, Relay ou urql
- Vous préférez un schéma typé fort
- Vous voulez GraphQL Playground pour l'exploration
- Vos clients ont des besoins de données variables



## Bases de données supportées

Go Starter Kit supporte **3 types de bases de données** pour s'adapter à vos besoins:

### Vue d'ensemble

| Base de données | Idéal pour | Configuration | Commande |
|-----------------|------------|---------------|----------|
| **PostgreSQL** | Production, queries complexes | Docker | `create-go-starter mon-app` |
| **MySQL** | Compatibilité large, shared hosting | Docker | `create-go-starter mon-app --database=mysql` |
| **SQLite** | Prototypage, petites apps | Aucune | `create-go-starter mon-app --database=sqlite` |

### PostgreSQL (défaut)

```bash
create-go-starter mon-app
# OU explicitement:
create-go-starter mon-app --database=postgres
```

**Avantages**:
- <i class="material-icons success small">check</i> Features SQL avancées (JSON, arrays, full-text search)
- <i class="material-icons success small">check</i> Excellente performance et fiabilité
- <i class="material-icons success small">check</i> ACID compliant, forte intégrité des données
- <i class="material-icons success small">check</i> Idéal pour production

**Quand l'utiliser**: Applications production, données complexes, besoin de fiabilité

### MySQL/MariaDB

```bash
create-go-starter mon-app --database=mysql
```

**Avantages**:
- <i class="material-icons success small">check</i> Compatibilité large avec hébergeurs
- <i class="material-icons success small">check</i> Excellent pour workloads read-heavy
- <i class="material-icons success small">check</i> Écosystème mature et bien documenté

**Quand l'utiliser**: Shared hosting, équipes MySQL, compatibilité large

### SQLite

```bash
create-go-starter mon-app --database=sqlite
```

**Avantages**:
- <i class="material-icons success small">check</i> Zéro configuration (pas de serveur)
- <i class="material-icons success small">check</i> Parfait pour prototypage rapide
- <i class="material-icons success small">check</i> Base de données dans un seul fichier
- <i class="material-icons success small">check</i> Très rapide pour petits datasets

**Limitations**:
- <i class="material-icons warning small">warning</i> Écritures concurrentes limitées
- <i class="material-icons warning small">warning</i> Pas adapté pour production à grande échelle

**Quand l'utiliser**: Prototypage, MVPs, développement, apps embarquées, petite production (<100 utilisateurs)

### Guide complet

Pour une comparaison détaillée, exemples de configuration, et guides de migration, consultez:

- **[Guide de sélection des databases](databases.md)** - Comparaison complète et aide au choix
- **[Guide de migration](database-migration.md)** - Migration entre databases

## Framework web (--framework)

<i class="material-icons success">new_releases</i> **Nouveau dans v1.6.0!** Go Starter Kit supporte plusieurs frameworks web via le flag `--framework` (ou `-f`).

### Frameworks disponibles

| Framework | Description | Statut |
|-----------|-------------|--------|
| `fiber` | Fiber v2 - Fast HTTP framework inspired by Express **(défaut)** | <i class="material-icons success small">circle</i> Disponible |
| `gin` | Gin - High-performance HTTP web framework | <i class="material-icons warning small">circle</i> Planifié v2.0.0 |
| `echo` | Echo - Minimalist high-performance HTTP framework | <i class="material-icons warning small">circle</i> Planifié v2.0.0 |

```bash
# Fiber (défaut) — pas besoin de spécifier
create-go-starter mon-app

# Avec flag explicite
create-go-starter mon-app --framework=fiber
create-go-starter mon-app -f fiber
```

> <i class="material-icons warning">warning</i> Gin et Echo retournent une erreur "not yet supported (planned for v2.0.0)" — ils seront disponibles dans les stories 11.2 et 11.3.

## Observabilité (--observability)

<i class="material-icons success">new_releases</i> **Nouveau dans v1.3.0!** Go Starter Kit supporte **3 niveaux d'observabilité** pour monitorer vos projets générés en production.

### Niveaux disponibles

| Niveau | Description | Ce qui est généré |
|--------|-------------|-------------------|
| `none` | Aucune observabilité (défaut) | Comportement standard préservé |
| `basic` | Health checks K8s avancés | `/health/liveness`, `/health/readiness` |
| `advanced` | Stack complète d'observabilité | Prometheus + Jaeger + Grafana + Health Checks K8s |

> <i class="material-icons warning">warning</i> Le flag `--observability` ne fonctionne qu'avec `--template=full`.

### Mode `advanced` — Stack d'observabilité complète

```bash
# Générer un projet avec observabilité avancée
create-go-starter mon-app --observability=advanced

# Combiné avec database et template
create-go-starter mon-app --template=full --database=postgres --observability=advanced
```

**Fichiers générés** (en plus des fichiers standard) :

```
mon-app/
├── pkg/
│   ├── metrics/
│   │   └── prometheus.go                # Registry Prometheus + métriques HTTP
│   └── tracing/
│       └── tracer.go                    # Configuration OpenTelemetry + TracerProvider
├── internal/adapters/
│   ├── middleware/
│   │   ├── metrics_middleware.go        # Middleware Prometheus pour métriques HTTP
│   │   └── tracing_middleware.go        # Middleware OpenTelemetry pour spans HTTP
│   ├── handlers/
│   │   └── metrics_handler.go           # Handler GET /metrics
│   └── http/
│       └── health.go                    # Health checks avancés (liveness/readiness)
├── internal/infrastructure/database/
│   └── tracing.go                       # Instrumentation GORM avec spans DB
├── pkg/logger/
│   └── logger_tracing.go               # Logger enrichi avec trace_id/span_id
├── monitoring/
│   ├── grafana/
│   │   ├── provisioning/
│   │   │   ├── datasources/
│   │   │   │   └── prometheus.yml       # Datasource Prometheus auto-configurée
│   │   │   └── dashboards/
│   │   │       └── dashboard.yml        # Auto-provisioning des dashboards
│   │   └── dashboards/
│   │       └── app-dashboard.json       # Dashboard Grafana 7 panneaux
│   └── prometheus/
│       ├── prometheus.yml               # Configuration scraping Prometheus
│       └── alert_rules.yml              # Règles d'alerting
├── deployments/
│   └── kubernetes/
│       └── probes.yaml                  # Configuration probes K8s
└── docker-compose.yml                   # Stack complète (DB + Jaeger + Prometheus + Grafana)
```

#### Prometheus Metrics

**Métriques exposées sur `GET /metrics`** :

| Métrique | Type | Description |
|----------|------|-------------|
| `http_requests_total` | Counter | Requêtes totales par méthode, route et code status |
| `http_request_duration_seconds` | Histogram | Latence HTTP par route (p50, p90, p95, p99) |
| `http_requests_in_flight` | Gauge | Nombre de requêtes actives en cours |

**Librairie** : [`fiberprometheus/v2`](https://github.com/ansrivas/fiberprometheus) v2.7.0.

```bash
# Tester les métriques
curl http://localhost:8080/metrics
```

#### Distributed Tracing (OpenTelemetry + Jaeger)

**Architecture** : Application → OTLP/gRPC → Jaeger Collector → Jaeger UI

**Fonctionnalités** :
- <i class="material-icons success small">check</i> Spans automatiques pour chaque requête HTTP
- <i class="material-icons success small">check</i> Spans automatiques pour chaque query GORM
- <i class="material-icons success small">check</i> Propagation W3C `traceparent` entre services
- <i class="material-icons success small">check</i> Logs zerolog enrichis avec `trace_id` et `span_id`

**Variables d'environnement** :
```bash
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=mon-app
```

**Accéder à Jaeger UI** :
```bash
open http://localhost:16686
```

#### Health Checks Kubernetes-compatible

| Endpoint | Probe K8s | Comportement |
|----------|-----------|--------------|
| `GET /health/liveness` | livenessProbe | 200 si l'app tourne |
| `GET /health/readiness` | readinessProbe | 200 si DB ok, 503 si DB down |
| `GET /health` | — | Alias rétrocompatible vers liveness |

**Readiness check** : Vérifie la connexion DB avec un timeout de 2 secondes. Retourne 503 avec détails si la DB est inaccessible.

**Fichier Kubernetes** : `deployments/kubernetes/probes.yaml` est automatiquement généré avec les configurations recommandées pour `livenessProbe`, `readinessProbe` et `startupProbe`.

#### Grafana Dashboard

**Dashboard 7 panneaux** pré-configuré et auto-provisionné :

| Panneau | Type | Description |
|---------|------|-------------|
| Request Rate | Time series | Requêtes/seconde |
| Error Rate | Time series | % erreurs (4xx, 5xx) |
| Latency P95 | Time series | 95ème percentile |
| Latency P99 | Time series | 99ème percentile |
| Requests in Flight | Gauge | Requêtes actives |
| Status Codes | Pie chart | Répartition des codes |
| Top Endpoints | Table | Endpoints les plus utilisés |

**Règles d'alerting** incluses :
- High Error Rate (> 5% pendant 5 min)
- High Latency (p95 > 1s pendant 5 min)
- Service Down (absence de métriques pendant 1 min)

**Accéder à Grafana** :
```bash
open http://localhost:3000
# Credentials: admin / admin
```

#### Docker Compose — Stack complète

Le `docker-compose.yml` généré inclut tous les services :

| Service | Image | Port | URL |
|---------|-------|------|-----|
| app | Build local | 8080 | `http://localhost:8080` |
| db | postgres:16-alpine | 5432 | — |
| jaeger | jaegertracing/all-in-one:1.56.0 | 16686 | `http://localhost:16686` |
| prometheus | prom/prometheus:v2.51.0 | 9090 | `http://localhost:9090` |
| grafana | grafana/grafana:10.4.0 | 3000 | `http://localhost:3000` |

```bash
# Démarrer toute la stack
docker-compose up -d

# Vérifier
curl http://localhost:8080/health/readiness
curl http://localhost:8080/metrics
```

### Mode `basic` — Health checks améliorés

```bash
create-go-starter mon-app --observability=basic
```

Génère les endpoints `/health/liveness` et `/health/readiness` sans Prometheus, Jaeger ni Grafana.

### Mode `none` (défaut)

```bash
create-go-starter mon-app
# équivalent à:
create-go-starter mon-app --observability=none
```

Comportement standard préservé — aucune régression.

---

## Options disponibles

### Flags et alias

Tous les flags disposent d'alias courts pour une saisie plus rapide.

| Flag | Alias | Défaut | Description |
|------|-------|--------|-------------|
| `--help` | `-h` | | Afficher l'aide |
| `--interactive` | `-i` | `false` | Lancer le mode interactif guidé |
| `--dry-run` | `-n` | `false` | Prévisualiser les fichiers sans les créer |
| `--template` | `-t` | `full` | Template: `minimal`, `full`, `graphql` |
| `--database` | `-d` | `postgres` | Base de données: `postgres`, `mysql`, `sqlite` |
| `--observability` | `-o` | `none` | Observabilité: `none`, `basic`, `advanced` |

**Sous-commandes:**

| Commande | Description |
|----------|-------------|
| `doctor` | Diagnostics environnement (Go, Git, Docker) |
| `add-model` | Générer un modèle CRUD dans un projet existant |

**Syntaxe flexible:** Les flags acceptent `=` ou espace comme séparateur:
```bash
# Ces syntaxes sont équivalentes:
create-go-starter -t=minimal mon-app
create-go-starter -t minimal mon-app
create-go-starter --template=minimal mon-app
create-go-starter --template minimal mon-app
```

### Exemples

```bash
# Template minimal
create-go-starter mon-projet --template minimal
create-go-starter mon-projet -t minimal

# Template full (défaut)
create-go-starter mon-projet --template full
create-go-starter mon-projet  # Même résultat

# Template GraphQL
create-go-starter mon-projet --template graphql
create-go-starter mon-projet -t graphql

# Choisir MySQL comme base de données
create-go-starter mon-projet --database=mysql
create-go-starter mon-projet -d mysql

# Choisir SQLite (idéal pour prototypage)
create-go-starter mon-projet --database=sqlite
create-go-starter mon-projet -d sqlite

# Combiner template et database
create-go-starter mon-projet --template=minimal --database=sqlite
create-go-starter mon-projet -t minimal -d sqlite

# Dry-run : prévisualiser sans créer de fichiers
create-go-starter mon-projet --dry-run
create-go-starter mon-projet -n

# Dry-run avec options spécifiques
create-go-starter mon-projet -t minimal -d sqlite -n

# Mode interactif guidé
create-go-starter --interactive
create-go-starter -i

# Combinaison complète avec observabilité
create-go-starter mon-projet -t full -d postgres -o advanced
```

> **Notes**:
> - Le flag `--template` est optionnel. Si non spécifié, le template **full** est utilisé par défaut.
> - Le flag `--database` est optionnel. Si non spécifié, **PostgreSQL** est utilisé par défaut.
> - Le flag `--dry-run` est optionnel. Il affiche la liste des fichiers qui seraient générés sans les créer.
> - Le flag `--interactive` lance un assistant guidé et ignore les autres flags de configuration.

## Mode Interactif (--interactive) <i class="material-icons success">new_releases</i>

**Nouveau dans v1.4.0!** Le mode interactif guide l'utilisateur étape par étape pour configurer un nouveau projet avec une interface terminal plus soignée.

### Utilisation

```bash
create-go-starter --interactive
create-go-starter -i
```

### Fonctionnement

Le mode interactif affiche un écran d'accueil, puis guide successivement:

1. **Nom du projet** -- Avec validation en temps réel
2. **Template** -- Choix entre minimal, full (défaut), graphql
3. **Base de données** -- Choix entre postgres (défaut), mysql, sqlite
4. **Observabilité** -- Choix entre none (défaut), basic, advanced (si template full)
5. **Résumé** -- Affichage de la configuration choisie avec confirmation

Exemple de parcours:

```text
Welcome screen
  Create New Project
  Help
  Exit

Project Name
  mon-app

Select a Template
  full

Select a Database
  postgres

Select Observability
  advanced

Configuration Summary
  Project Name    mon-app
  Template        full
  Database        postgres
  Observability   advanced
```

### Notes

- <i class="material-icons info">info</i> Le mode interactif nécessite un terminal interactif (pas de pipe stdin)
- <i class="material-icons warning">warning</i> `--interactive` et `--dry-run` ne peuvent pas être utilisés ensemble
- <i class="material-icons info">info</i> Le flux interactif propose une hiérarchie visuelle plus claire: écran d'accueil, résumé centré et progression avec statistiques

## Prévisualisation Dry-Run (--dry-run) <i class="material-icons success">new_releases</i>

**Nouveau dans v1.4.0!** Le mode dry-run affiche les fichiers qui seraient générés sans les créer.

### Utilisation

```bash
create-go-starter mon-app --dry-run
create-go-starter -n -t minimal -d sqlite mon-app
```

### Fonctionnement

Le dry-run affiche:
- La configuration utilisée (template, database, observabilité)
- La liste complète des fichiers qui seraient créés
- Le nombre de fichiers et répertoires
- Un avertissement si le répertoire cible existe déjà

```
Dry-run mode: no files will be created

Configuration:
  Project:       mon-app
  Template:      minimal
  Database:      sqlite
  Observability: none

Files that would be generated (23 files):
  mon-app/cmd/main.go
  mon-app/internal/models/user.go
  mon-app/internal/domain/user/service.go
  ...

Summary: 23 files in 12 directories
```

### Compatible avec tous les flags

```bash
create-go-starter -n -t full -d postgres -o advanced mon-app
```

## Commande Doctor <i class="material-icons success">new_releases</i>

**Nouveau dans v1.4.0!** La commande `doctor` vérifie que votre environnement est correctement configuré.

### Utilisation

```bash
create-go-starter doctor
```

### Fonctionnement

La commande vérifie:

| Outil | Vérification | Requis |
|-------|-------------|--------|
| **Go** | Version >= 1.21 installée | <i class="material-icons success small">check</i> Oui |
| **Git** | Binaire disponible | <i class="material-icons success small">check</i> Recommandé |
| **Docker** | Binaire + daemon actif | <i class="material-icons info small">info</i> Optionnel |

### Exemple de sortie

```
create-go-starter doctor v1.6.0

Checking environment...

  [OK] Go 1.25.5 (minimum: 1.21)
  [OK] Git 2.43.0
  [OK] Docker 24.0.7 (daemon running)

All checks passed!
```

### Code de sortie

- `0` -- Tout est OK
- `1` -- Un ou plusieurs problèmes détectés

Utile pour les scripts CI/CD:
```bash
create-go-starter doctor && echo "Ready!" || echo "Fix issues first"
```

## Ajouter des modèles (add-model) <i class="material-icons success">new_releases</i>

**Nouveau dans v1.2.0!** Une fois votre projet créé, vous pouvez générer automatiquement de nouveaux modèles CRUD complets avec la commande `add-model`.

### Syntaxe de base

```bash
create-go-starter add-model <ModelName> --fields "field1:type1,field2:type2,..."
```

### Exemple simple

```bash
cd mon-projet  # Naviguer dans un projet go-starter-kit existant

# Ajouter un modèle Todo
create-go-starter add-model Todo --fields "title:string,completed:bool"
```

**Cette commande génère automatiquement:**
- <i class="material-icons success small">circle</i> `internal/models/todo.go` - Model avec tags GORM
- <i class="material-icons success small">circle</i> `internal/interfaces/todo_repository.go` - Interface repository
- <i class="material-icons success small">circle</i> `internal/adapters/repository/todo_repository.go` - Implémentation GORM
- <i class="material-icons success small">circle</i> `internal/domain/todo/service.go` - Logique métier
- <i class="material-icons success small">circle</i> `internal/domain/todo/module.go` - Module fx
- <i class="material-icons success small">circle</i> `internal/adapters/handlers/todo_handler.go` - Handlers HTTP
- <i class="material-icons success small">circle</i> `internal/domain/todo/service_test.go` - Tests service
- <i class="material-icons success small">circle</i> `internal/adapters/handlers/todo_handler_test.go` - Tests handler

**Et modifie automatiquement:**
- <i class="material-icons warning small">circle</i> `internal/infrastructure/database/database.go` - Ajoute AutoMigrate
- <i class="material-icons warning small">circle</i> `internal/adapters/http/routes.go` - Ajoute routes CRUD
- <i class="material-icons warning small">circle</i> `cmd/main.go` - Ajoute module fx

### Types de champs supportés

| Type | Go Type | Exemple | Description |
|------|---------|---------|-------------|
| `string` | `string` | `"title:string"` | Chaîne de caractères |
| `int` | `int` | `"age:int"` | Entier signé |
| `uint` | `uint` | `"count:uint"` | Entier non-signé |
| `float64` | `float64` | `"price:float64"` | Nombre décimal |
| `bool` | `bool` | `"active:bool"` | Booléen (true/false) |
| `time` | `time.Time` | `"birthdate:time"` | Date et heure |

### Modificateurs GORM

Ajoutez des modificateurs après le type pour personnaliser les contraintes de base de données:

| Modificateur | GORM Tag | Description |
|--------------|----------|-------------|
| `unique` | `gorm:"unique"` | Contrainte d'unicité |
| `not_null` | `gorm:"not null"` | Ne peut pas être NULL |
| `index` | `gorm:"index"` | Créer un index DB |

**Syntaxe:**
```
field:type:modifier1:modifier2:...
```

**Exemples:**
```bash
# Email unique et obligatoire
create-go-starter add-model User --fields "email:string:unique:not_null"

# Product avec contraintes
create-go-starter add-model Product --fields "name:string:unique:not_null,price:float64,stock:int:index"
```

### Relations entre modèles

`add-model` supporte les relations **BelongsTo** (enfant → parent) et **HasMany** (parent → enfants).

#### BelongsTo (enfant appartient à un parent)

Ajoute une **foreign key** et un champ de relation au modèle enfant.

```bash
# Le modèle parent DOIT exister d'abord
create-go-starter add-model Todo --fields "title:string,completed:bool"

# Créer Comment qui appartient à Todo
create-go-starter add-model Comment --fields "content:string" --belongs-to Todo
```

**Résultat dans `internal/models/comment.go`:**
```go
type Comment struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Content   string         `gorm:"not null" json:"content"`
    TodoID    uint           `gorm:"not null;index" json:"todo_id"`      // Foreign key
    Todo      Todo           `gorm:"foreignKey:TodoID" json:"todo,omitempty"` // Relation
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Endpoints générés:**
- `GET /api/v1/comments/:id?include=todo` - Récupère comment avec preload du todo
- `GET /api/v1/todos/:todoId/comments` - Liste tous les comments d'un todo
- `POST /api/v1/todos/:todoId/comments` - Créer un comment pour un todo

#### HasMany (parent a plusieurs enfants)

Modifie le **modèle parent** pour ajouter un slice d'enfants.

```bash
# Créer Category avec relation has-many vers Product
create-go-starter add-model Category --fields "name:string:unique" --has-many Product
```

**Résultat dans `internal/models/category.go`:**
```go
type Category struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"unique;not null" json:"name"`
    Products  []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"` // Slice ajouté
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Note:** Le modèle `Product` DOIT déjà exister avant d'utiliser `--has-many`.

#### Exemple complet avec relations

```bash
# 1. Créer le modèle parent
create-go-starter add-model Category --fields "name:string:unique"

# 2. Créer le modèle enfant avec belongs-to
create-go-starter add-model Product --fields "name:string,price:float64" --belongs-to Category

# 3. Optionnel: Ajouter has-many au parent après coup
# (pas nécessaire si fait pendant création de Product)

# 4. Utiliser le preloading dans les requêtes
# GET /api/v1/products/1?include=category
# GET /api/v1/categories/1?include=products
```

### Flags disponibles

| Flag | Description | Requis |
|------|-------------|--------|
| `--fields` | Définition des champs (`"name:type:mod,..."`) | Oui <i class="material-icons success small">check</i> |
| `--belongs-to <Model>` | Ajoute relation BelongsTo (foreign key) | Non |
| `--has-many <Model>` | Ajoute relation HasMany (slice d'enfants) | Non |
| `--public` | Routes publiques (sans auth middleware) | Non |
| `--yes`, `-y` | Skip confirmation prompt | Non |
| `--help`, `-h` | Afficher l'aide | Non |

### Notes importantes

<i class="material-icons warning">warning</i> **Pluralisation:** La commande utilise des règles simples de pluralisation (ajoute un 's'). Pour les pluriels irréguliers (Person→People, Child→Children), éditez manuellement le code généré.

<i class="material-icons warning">warning</i> **Relations many-to-many:** Pas encore supportées. Utilisez `--belongs-to` pour créer des relations via tables de jointure manuelles.

<i class="material-icons info">info</i> **Swagger:** Après génération, exécutez `make swagger` ou `swag init -g cmd/api/main.go -o docs/swagger` pour mettre à jour la documentation API.

### Aide en ligne

```bash
create-go-starter add-model --help
```

### Exemple Complet: Blog System

Voici un exemple complet de création d'un système de blog avec 3 modèles imbriqués: **Category** → **Post** → **Comment**.

#### 1. Créer le projet initial

```bash
create-go-starter blog-api
cd blog-api
./setup.sh  # Configurer la base de données
```

#### 2. Ajouter le modèle Category (racine)

```bash
create-go-starter add-model Category --fields "name:string:unique:not_null,description:string"
```

**Fichiers générés:**
- `internal/models/category.go` - Entity
- `internal/interfaces/category_repository.go` - Port
- `internal/adapters/repository/category_repository.go` - GORM impl
- `internal/domain/category/service.go` - Business logic
- `internal/domain/category/module.go` - fx module
- `internal/adapters/handlers/category_handler.go` - HTTP handlers
- Tests: `service_test.go`, `handler_test.go`

**Endpoints disponibles:**
- `POST /api/v1/categories` - Créer catégorie
- `GET /api/v1/categories` - Liste catégories
- `GET /api/v1/categories/:id` - Récupérer catégorie
- `PUT /api/v1/categories/:id` - Modifier catégorie
- `DELETE /api/v1/categories/:id` - Supprimer catégorie

#### 3. Ajouter le modèle Post (enfant de Category)

```bash
create-go-starter add-model Post --fields "title:string:not_null,content:string,published:bool" --belongs-to Category
```

**Ce qui est ajouté automatiquement:**
- Champs dans `internal/models/post.go`:
  ```go
  CategoryID uint     `gorm:"not null;index" json:"category_id"`
  Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
  ```
- Routes imbriquées dans `internal/adapters/http/routes.go`:
  ```go
  categories.Get("/:categoryId/posts", postHandler.GetByParent)
  categories.Post("/:categoryId/posts", postHandler.CreateForParent)
  ```

**Endpoints disponibles:**
- CRUD standard: `POST/GET/PUT/DELETE /api/v1/posts`
- **Relation parent**: `GET /api/v1/categories/:categoryId/posts` - Posts d'une catégorie
- **Relation parent**: `POST /api/v1/categories/:categoryId/posts` - Créer post dans catégorie
- **Preloading**: `GET /api/v1/posts/:id?include=category` - Post avec sa catégorie

#### 4. Ajouter le modèle Comment (enfant de Post)

```bash
create-go-starter add-model Comment --fields "author:string:not_null,content:string:not_null" --belongs-to Post
```

**Ce qui est ajouté automatiquement:**
- Champs dans `internal/models/comment.go`:
  ```go
  PostID  uint `gorm:"not null;index" json:"post_id"`
  Post    Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
  ```
- Routes imbriquées:
  ```go
  posts.Get("/:postId/comments", commentHandler.GetByParent)
  posts.Post("/:postId/comments", commentHandler.CreateForParent)
  ```

**Endpoints disponibles:**
- CRUD standard: `POST/GET/PUT/DELETE /api/v1/comments`
- **Relation parent**: `GET /api/v1/posts/:postId/comments` - Comments d'un post
- **Relation parent**: `POST /api/v1/posts/:postId/comments` - Créer comment sur post
- **Preloading**: `GET /api/v1/comments/:id?include=post` - Comment avec son post

#### 5. Ajouter HasMany au parent (optionnel)

Si vous voulez preload les enfants depuis le parent:

```bash
# Ajouter Posts à Category
create-go-starter add-model Category --has-many Post

# Ajouter Comments à Post
create-go-starter add-model Post --has-many Comment
```

**Ce qui est ajouté:**
- Dans `internal/models/category.go`:
  ```go
  Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
  ```
- Dans `internal/models/post.go`:
  ```go
  Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
  ```

**Nouveaux endpoints de preloading:**
- `GET /api/v1/categories/:id?include=posts` - Catégorie avec tous ses posts
- `GET /api/v1/posts/:id?include=comments` - Post avec tous ses comments
- `GET /api/v1/posts/:id?include=category,comments` - Post avec catégorie ET comments

#### 6. Tester le système complet

```bash
# Rebuild le projet
go mod tidy
go build ./...

# Générer Swagger (optionnel)
make swagger

# Lancer tests
make test

# Démarrer le serveur
make run
```

#### 7. Exemples d'utilisation API

**Créer une catégorie:**
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Technology", "description": "Tech articles"}'
```

**Créer un post dans cette catégorie:**
```bash
curl -X POST http://localhost:8080/api/v1/categories/1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Go",
    "content": "Go is a statically typed language...",
    "published": true
  }'
```

**Créer un commentaire sur ce post:**
```bash
curl -X POST http://localhost:8080/api/v1/posts/1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "author": "John Doe",
    "content": "Great article!"
  }'
```

**Récupérer un post avec sa catégorie et ses comments:**
```bash
curl http://localhost:8080/api/v1/posts/1?include=category,comments
```

**Réponse attendue:**
```json
{
  "data": {
    "id": 1,
    "title": "Introduction to Go",
    "content": "Go is a statically typed language...",
    "published": true,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Technology",
      "description": "Tech articles"
    },
    "comments": [
      {
        "id": 1,
        "author": "John Doe",
        "content": "Great article!",
        "post_id": 1
      }
    ],
    "created_at": "2026-02-12T10:00:00Z",
    "updated_at": "2026-02-12T10:00:00Z"
  }
}
```

### Conventions et Limitations

#### Pluralisation

La commande `add-model` utilise des règles de pluralisation simples:

**Règles appliquées:**
- Ajoute 's': `Todo` → `todos`, `Product` → `products`
- Remplace 'y' par 'ies': `Category` → `categories`, `Company` → `companies`
- Ajoute 'es' pour s/x/z: `Class` → `classes`, `Box` → `boxes`, `Quiz` → `quizzes`

**Pluriels irréguliers non supportés:**
- `Person` → `People` (sera `persons`)
- `Child` → `Children` (sera `childs`)
- `Mouse` → `Mice` (sera `mouses`)

<i class="material-icons info">info</i> **Solution:** Éditez manuellement les fichiers générés pour corriger les pluriels irréguliers:
- Renommer fichiers: `internal/domain/persons/` → `internal/domain/people/`
- Mettre à jour imports et références dans le code
- Mettre à jour routes: `/api/v1/persons` → `/api/v1/people`

#### Relations supportées

| Relation | Status | Description |
|----------|--------|-------------|
| **BelongsTo** (N:1) | <i class="material-icons success">check_circle</i> Supporté | Enfant appartient à un parent |
| **HasMany** (1:N) | <i class="material-icons success">check_circle</i> Supporté | Parent a plusieurs enfants |
| **Many-to-Many** | <i class="material-icons warning">error</i> Pas encore | Prévue pour une version future |

**Workaround Many-to-Many:**

Pour créer une relation many-to-many entre `User` et `Role`:

```bash
# 1. Créer la table de jointure
create-go-starter add-model UserRole --fields "user_id:uint:index,role_id:uint:index" --belongs-to User --belongs-to Role

# 2. Éditer manuellement pour ajouter contrainte unique composite
# Dans internal/models/user_role.go:
# type UserRole struct {
#     UserID uint `gorm:"uniqueIndex:user_role_unique;not null"`
#     RoleID uint `gorm:"uniqueIndex:user_role_unique;not null"`
#     ...
# }
```

#### Champs réservés

Ces champs sont **automatiquement ajoutés** et ne doivent PAS être spécifiés dans `--fields`:

| Champ | Type | Description |
|-------|------|-------------|
| `ID` | `uint` | Clé primaire |
| `CreatedAt` | `time.Time` | Date de création (auto) |
| `UpdatedAt` | `time.Time` | Date de modification (auto) |
| `DeletedAt` | `gorm.DeletedAt` | Soft delete (nullable) |

<i class="material-icons warning">warning</i> Si vous spécifiez ces champs, ils seront **ignorés** par le générateur.

#### Validations custom

Les validations de base sont générées automatiquement, mais les validations complexes doivent être ajoutées manuellement.

**Validations automatiques:**
- `not_null` → Tag GORM: `gorm:"not null"`
- `unique` → Tag GORM: `gorm:"unique"`

**Validations à ajouter manuellement:**

```go
// Dans internal/models/user.go
type User struct {
    Email string `gorm:"unique;not null" json:"email" validate:"required,email"` // Ajouter validate tag
    Age   int    `gorm:"not null" json:"age" validate:"gte=0,lte=150"`          // Validation range
}
```

```go
// Dans internal/domain/user/service.go
func (s *Service) Create(ctx context.Context, user *models.User) error {
    // Ajouter validation custom
    if user.Age < 18 {
        return domain.ErrValidation("user must be at least 18 years old")
    }
    return s.repo.Create(ctx, user)
}
```

#### Swagger documentation

Après avoir généré un nouveau modèle, **vous devez regénérer la documentation Swagger**:

```bash
make swagger
# Ou
swag init -g cmd/api/main.go -o docs/swagger
```

<i class="material-icons warning">warning</i> Sans regénération, les nouveaux endpoints n'apparaîtront pas dans Swagger UI.

**Ajouter annotations Swagger** (optionnel):

```go
// Dans internal/adapters/handlers/todo_handler.go

// Create godoc
// @Summary      Create todo
// @Description  Create a new todo item
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo body models.Todo true "Todo object"
// @Success      201  {object}  models.Todo
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/v1/todos [post]
func (h *Handler) Create(c *fiber.Ctx) error {
    // ...
}
```

#### Routes et middleware

**Par défaut**, toutes les routes générées sont **protégées par JWT** (authentification requise).

Pour créer des **routes publiques** (sans authentification):

```bash
create-go-starter add-model Article --fields "title:string,content:string" --public
```

Cela génère les routes **sans** middleware `auth.RequireAuth`:

```go
// routes.go - routes publiques
api.Get("/articles", articleHandler.List)
api.Get("/articles/:id", articleHandler.GetByID)
api.Post("/articles", articleHandler.Create)        // Public!
api.Put("/articles/:id", articleHandler.Update)     // Public!
api.Delete("/articles/:id", articleHandler.Delete)  // Public!
```

<i class="material-icons warning">warning</i> **Utilisez --public avec précaution** pour éviter les failles de sécurité.

## Conventions de nommage

Le nom du projet doit respecter certaines règles:

### Caractères autorisés

- **Lettres**: a-z, A-Z
- **Chiffres**: 0-9
- **Tirets**: - (hyphen)
- **Underscores**: _ (trait de soulignement)

### Restrictions

- Pas d'espaces
- Pas de caractères spéciaux (/, \, @, #, etc.)
- Pas de points (.)
- Doit commencer par une lettre ou un chiffre (pas de tiret au début)

### Exemples valides

```bash
create-go-starter mon-projet           # ✅ Valide
create-go-starter my-awesome-api       # ✅ Valide
create-go-starter user_service         # ✅ Valide
create-go-starter app2024              # ✅ Valide
create-go-starter MonProjet            # ✅ Valide
```

### Exemples invalides

```bash
create-go-starter mon projet           # contient un espace
create-go-starter mon-projet!          # caractère spécial
create-go-starter my.project           # contient un point
create-go-starter -mon-projet          # commence par un tiret
create-go-starter mon/projet           # contient un slash
```

## Structure générée

Voici la structure complète créée par `create-go-starter`:

```
mon-projet/
├── cmd/
│   └── main.go                              # Point d'entrée de l'application
│
├── internal/
│   ├── models/
│   │   └── user.go                          # Entités partagées: User, RefreshToken, AuthResponse
│   │
│   ├── adapters/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go              # Endpoints auth (register, login, refresh)
│   │   │   ├── auth_handler_test.go         # Tests handlers auth
│   │   │   ├── user_handler.go              # Endpoints CRUD users
│   │   │   └── user_handler_test.go         # Tests handlers users
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go           # Middleware JWT authentication
│   │   │   ├── auth_middleware_test.go      # Tests middleware auth
│   │   │   ├── error_handler.go             # Middleware gestion centralisée erreurs
│   │   │   └── error_handler_test.go        # Tests error handler
│   │   ├── repository/
│   │   │   ├── user_repository.go           # Implémentation GORM du repository
│   │   │   └── user_repository_test.go      # Tests repository
│   │   └── http/
│   │       ├── health.go                    # Handler health check
│   │       ├── health_test.go               # Tests health check
│   │       └── routes.go                    # Routes centralisées de l'API
│   │
│   ├── domain/
│   │   ├── errors.go                        # Erreurs métier personnalisées
│   │   ├── errors_test.go                   # Tests erreurs
│   │   └── user/
│   │       ├── service.go                   # Logique métier (Register, Login, etc.)
│   │       ├── service_test.go              # Tests service
│   │       └── module.go                    # Module fx pour dependency injection
│   │
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── database.go                  # Configuration GORM et connexion DB
│   │   │   ├── database_test.go             # Tests database
│   │   │   └── module.go                    # Module fx pour database
│   │   └── server/
│   │       ├── server.go                    # Configuration Fiber app
│   │       ├── server_test.go               # Tests server
│   │       └── module.go                    # Module fx pour server
│   │
│   └── interfaces/
│       ├── auth_service.go                  # Interface AuthService (port)
│       ├── user_service.go                  # Interface UserService (port)
│       └── user_repository.go               # Interface UserRepository (port)
│
├── pkg/
│   ├── auth/
│   │   ├── jwt.go                           # Génération et parsing JWT tokens
│   │   ├── jwt_test.go                      # Tests JWT
│   │   ├── middleware.go                    # Middleware JWT pour Fiber
│   │   ├── middleware_test.go               # Tests middleware
│   │   └── module.go                        # Module fx pour auth
│   ├── config/
│   │   ├── env.go                           # Chargement variables d'environnement
│   │   ├── env_test.go                      # Tests config
│   │   └── module.go                        # Module fx pour config
│   └── logger/
│       ├── logger.go                        # Configuration zerolog
│       ├── logger_test.go                   # Tests logger
│       └── module.go                        # Module fx pour logger
│
├── .github/
│   └── workflows/
│       └── ci.yml                           # Pipeline CI/CD GitHub Actions
│
├── .env                                     # Variables d'environnement (créé automatiquement)
├── .env.example                             # Template de configuration
├── .gitignore                               # Exclusions Git
├── .golangci.yml                            # Configuration golangci-lint
├── Dockerfile                               # Build Docker multi-stage
├── Makefile                                 # Commandes utiles (run, test, lint, etc.)
├── setup.sh                                 # Script de configuration automatique
├── go.mod                                   # Module Go et dépendances
└── README.md                                # Documentation du projet
```

**Total**: ~45+ fichiers créés automatiquement!

## Explication détaillée de chaque répertoire

### `/cmd`

**Rôle**: Point d'entrée de l'application.

**Contenu**:
- `main.go`: Bootstrap de l'application avec uber-go/fx
  - Initialise tous les modules (database, server, logger, etc.)
  - Configure le lifecycle (OnStart, OnStop)
  - Lance l'application avec `fx.New().Run()`

**Pattern**: Un seul fichier minimal qui orchestre les dépendances.

### `/internal/models`

**Rôle**: Entités de domaine partagées (domain entities) utilisées à travers toute l'application.

**Contenu**:
- `user.go`: Définit toutes les entités liées aux utilisateurs
  - `User`: Entité principale avec tags GORM (ID, Email, PasswordHash, timestamps)
  - `RefreshToken`: Entité pour la gestion des refresh tokens JWT
  - `AuthResponse`: DTO pour les réponses d'authentification

**Pourquoi un package séparé?**
- **Évite les dépendances circulaires**: Avant, `interfaces` → `domain/user` → `interfaces` créait un cycle
- **Maintenant**: `interfaces` → `models` ← `domain/user` (pas de cycle!)
- **Centralisation**: Les entités sont définies en un seul endroit
- **Réutilisabilité**: Tous les layers (domain, interfaces, adapters) peuvent importer models sans conflit

**Import**:
```go
import "mon-projet/internal/models"

user := &models.User{
    Email: "user@example.com",
    PasswordHash: hashedPassword,
}
```

### `/internal/domain`

**Rôle**: Couche métier (logique business), cœur de l'architecture hexagonale.

**Contenu**:
- `errors.go`: Définition des erreurs métier (DomainError, NotFoundError, ValidationError, etc.)
- `user/`: Domaine User
  - `service.go`: Logique métier (Register, Login, GetUserByID, UpdateUser, DeleteUser)
  - `module.go`: Module fx qui fournit le service

**Principe**: Le domaine ne doit **jamais** importer d'autres packages (adapters, infrastructure). Les dépendances sont inversées via interfaces (ports). Les entités sont maintenant dans `internal/models` pour éviter les cycles de dépendances.

### `/internal/adapters`

**Rôle**: Adaptateurs qui connectent le domaine au monde extérieur (HTTP, DB).

#### `/internal/adapters/handlers`

**Rôle**: Handlers HTTP qui exposent l'API REST.

**Contenu**:
- `auth_handler.go`:
  - `Register`: POST /api/v1/auth/register
  - `Login`: POST /api/v1/auth/login
  - `RefreshToken`: POST /api/v1/auth/refresh
- `user_handler.go`:
  - `List`: GET /api/v1/users
  - `GetByID`: GET /api/v1/users/:id
  - `Update`: PUT /api/v1/users/:id
  - `Delete`: DELETE /api/v1/users/:id

**Pattern**: Handlers parsent les requêtes, valident, appellent les services du domaine, retournent les réponses.

#### `/internal/adapters/middleware`

**Rôle**: Middleware Fiber pour cross-cutting concerns.

**Contenu**:
- `auth_middleware.go`: Vérifie le JWT token dans les requêtes
- `error_handler.go`: Gestion centralisée des erreurs (convertit DomainError en réponses HTTP)

#### `/internal/adapters/repository`

**Rôle**: Implémentation du pattern Repository avec GORM.

**Contenu**:
- `user_repository.go`: Implémentation de l'interface UserRepository
  - Create, FindByID, FindByEmail, Update, Delete
  - Utilise GORM pour les opérations DB

**Pattern**: Repository isole le domaine de la couche de persistance.

#### `/internal/adapters/http`

**Rôle**: Routes HTTP et handlers utilitaires.

**Contenu**:
- `health.go`: Endpoint GET /health pour monitoring
- `routes.go`: Configuration centralisée de toutes les routes de l'API

**Avantages de la centralisation des routes**:
- Vue d'ensemble de toutes les routes en un seul fichier
- Facilite la documentation et le versioning de l'API
- Séparation claire entre définition des routes et logique des handlers

### `/internal/infrastructure`

**Rôle**: Configuration de l'infrastructure (DB, serveur).

#### `/internal/infrastructure/database`

**Rôle**: Configuration et connexion à la base de données.

**Contenu**:
- `database.go`:
  - Connexion PostgreSQL via GORM
  - Configuration du pool de connexions
  - AutoMigrate des entités (`models.User`, `models.RefreshToken`)
  - Gestion du lifecycle (fermeture connexion)

#### `/internal/infrastructure/server`

**Rôle**: Configuration du serveur HTTP Fiber.

**Contenu**:
- `server.go`:
  - Configuration Fiber app
  - Middleware error handler
  - Lifecycle du serveur (start, graceful shutdown)
  - Note: Les routes sont enregistrées via `server.Module` qui invoque `httpRoutes.RegisterRoutes()` avec `fx.Invoke`

### `/internal/interfaces`

**Rôle**: Définition des ports (interfaces) pour l'architecture hexagonale.

**Contenu**:
- `user_repository.go`: Interface pour la persistance utilisateur

**Principe**:
- Les adapters dépendent de ces interfaces, pas des implémentations concrètes
- Les interfaces référencent `models.*` pour les types de retour/paramètres
- Exemple: `CreateUser(ctx context.Context, user *models.User) error`

### `/pkg`

**Rôle**: Packages réutilisables, peuvent être importés par d'autres projets.

#### `/pkg/auth`

**Rôle**: Utilitaires JWT.

**Contenu**:
- `jwt.go`: Génération et validation de tokens JWT
  - GenerateAccessToken (15min)
  - GenerateRefreshToken (7 jours)
  - ParseToken
- `middleware.go`: Middleware Fiber pour JWT

#### `/pkg/config`

**Rôle**: Chargement de la configuration.

**Contenu**:
- `env.go`: Charge les variables .env (godotenv) et les expose via struct Config

#### `/pkg/logger`

**Rôle**: Configuration du logger.

**Contenu**:
- `logger.go`: Configure zerolog (niveau, format, output)

### `/.github/workflows`

**Rôle**: CI/CD avec GitHub Actions.

**Contenu**:
- `ci.yml`: Pipeline qui exécute:
  1. golangci-lint (quality checks)
  2. Tests avec PostgreSQL (service container)
  3. Build verification

## Fichiers de configuration détaillés

### `.env` et `.env.example`

Le fichier `.env.example` est un template, et `.env` est copié automatiquement lors de la génération.

**Variables**:

```bash
# Application
APP_NAME=mon-projet                # Nom de l'app (utilisé dans logs)
APP_ENV=development                # Environnement (development, staging, production)
APP_PORT=8080                      # Port HTTP

# Database PostgreSQL
DB_HOST=localhost                  # Hôte de la DB
DB_PORT=5432                       # Port PostgreSQL
DB_USER=postgres                   # Utilisateur DB
DB_PASSWORD=postgres               # Mot de passe DB
DB_NAME=mon-projet                 # Nom de la base de données
DB_SSLMODE=disable                 # SSL mode (require pour production)

# JWT Authentication
JWT_SECRET=                        # SECRET CRITIQUE - À générer!
JWT_EXPIRY=15m                     # Durée des access tokens (15 minutes)
REFRESH_TOKEN_EXPIRY=168h          # Durée des refresh tokens (7 jours)
```

**Important**: Générez un JWT_SECRET sécurisé:

```bash
openssl rand -base64 32
```

Puis ajoutez-le dans `.env`:

```bash
JWT_SECRET=votre_secret_genere_ici
```

### `go.mod`

Définit le module Go et les dépendances:

```go
module mon-projet

go 1.25

require (
    github.com/gofiber/fiber/v2 v2.x.x
    github.com/golang-jwt/jwt/v5 v5.x.x
    gorm.io/gorm v1.x.x
    gorm.io/driver/postgres v1.x.x
    go.uber.org/fx v1.x.x
    github.com/rs/zerolog v1.x.x
    github.com/go-playground/validator/v10 v10.x.x
    golang.org/x/crypto v0.x.x
    github.com/joho/godotenv v1.x.x
    // ... autres dépendances
)
```

### `Makefile`

Commandes utiles pour le développement:

```makefile
.PHONY: help run build test test-coverage lint clean docker-build docker-run

help:           # Afficher toutes les commandes
run:            # Lancer l'application (go run)
build:          # Compiler le binaire
test:           # Exécuter les tests avec race detector
test-coverage:  # Tests avec rapport HTML de coverage
lint:           # Linter avec golangci-lint
clean:          # Nettoyer les artifacts
docker-build:   # Build l'image Docker
docker-run:     # Lancer le conteneur Docker
```

### `.golangci.yml`

Configuration du linter avec règles strictes:

- errcheck, gosimple, govet, ineffassign, staticcheck
- gofmt, goimports
- misspell, revive
- Exclusions pour tests et code généré

### `Dockerfile`

Build multi-stage optimisé:

**Stage 1 (builder)**:
- Image golang:1.25-alpine
- Build du binaire statique

**Stage 2 (runtime)**:
- Image alpine:latest (légère)
- Copie du binaire seulement
- Exécution en tant que non-root
- EXPOSE 8080

**Taille finale**: ~15-20MB (vs ~1GB avec image golang complète)

## Workflow après génération

Une fois le projet créé, vous avez deux options pour configurer votre projet:

### Option A: Configuration automatique avec setup.sh (Recommandé) <i class="material-icons success">rocket_launch</i>

Le CLI génère automatiquement un script `setup.sh` qui automatise toute la configuration initiale.

**Fonctionnalités du script**:
- <i class="material-icons success small">check_circle</i> Vérification des prérequis (Go, OpenSSL, Docker)
- <i class="material-icons success small">check_circle</i> Installation des dépendances Go (`go mod tidy`)
- <i class="material-icons success small">check_circle</i> Génération automatique du JWT secret
- <i class="material-icons success small">check_circle</i> Configuration de PostgreSQL (Docker ou local)
- <i class="material-icons success small">check_circle</i> Exécution des tests
- <i class="material-icons success small">check_circle</i> Vérification de l'installation

**Utilisation**:

```bash
cd mon-projet
./setup.sh
```

Le script est **interactif** et vous guidera à travers les choix:
- Docker ou PostgreSQL local
- Régénération du JWT secret si déjà configuré
- Validation à chaque étape

**Après l'exécution du script**:

```bash
make run
```

C'est tout! Votre application est prête.

---

### Option B: Configuration manuelle

Si vous préférez configurer manuellement ou si le script setup.sh échoue, suivez ces étapes:

### Étape 1: Naviguer dans le projet

```bash
cd mon-projet
```

### Étape 2: Configurer les variables d'environnement

```bash
# Générer un JWT secret
openssl rand -base64 32

# Éditer .env et ajouter le secret
nano .env  # ou vim, code, etc.
```

Ajoutez:
```
JWT_SECRET=le_secret_genere
```

### Étape 3: Installer les dépendances

```bash
go mod tidy
```

### Étape 4: Lancer PostgreSQL

**Option A: Docker (recommandé)**

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_DB=mon-projet \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

**Option B: Installation locale**

```bash
# macOS
brew install postgresql
brew services start postgresql
createdb mon-projet

# Linux
sudo apt install postgresql
sudo systemctl start postgresql
sudo -u postgres createdb mon-projet
```

### Étape 5: Lancer l'application

```bash
make run
```

Ou directement:

```bash
go run cmd/main.go
```

Vous devriez voir:

```
{"level":"info","time":"...","message":"Starting server on :8080"}
{"level":"info","time":"...","message":"Database connected successfully"}
```

### Étape 6: Tester l'API

```bash
# Health check
curl http://localhost:8080/health
# {"status":"ok"}

# Register un utilisateur
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## Commandes Make disponibles

Le Makefile généré inclut ces commandes:

| Commande | Description | Usage |
|----------|-------------|-------|
| `make help` | Afficher l'aide | Voir toutes les commandes |
| `make run` | Lancer l'app | Développement local |
| `make build` | Compiler le binaire | Créer l'exécutable |
| `make test` | Tests avec race detector | Vérifier le code |
| `make test-coverage` | Tests + rapport HTML | Voir le coverage |
| `make lint` | golangci-lint | Vérifier la qualité |
| `make clean` | Nettoyer | Supprimer artifacts |
| `make docker-build` | Build image Docker | Containerisation |
| `make docker-run` | Lancer conteneur | Test Docker |

**Exemples**:

```bash
# Développement quotidien
make run          # Lance l'app avec hot-reload (si air installé)

# Avant commit
make test         # Vérifie que tous les tests passent
make lint         # Vérifie la qualité du code

# Build pour production
make build        # Crée le binaire
./mon-projet      # Exécute le binaire

# Docker
make docker-build # Build l'image
make docker-run   # Teste le conteneur
```

## Prochaines étapes

Maintenant que vous comprenez la structure, consultez:

1. **[Guide des projets générés](./generated-project-guide.md)** - Guide complet pour:
   - Comprendre l'architecture hexagonale
   - Développer de nouvelles fonctionnalités
   - Utiliser l'API (endpoints, authentification)
   - Écrire des tests
   - Déployer en production

2. **[Architecture du CLI](./cli-architecture.md)** - Si vous voulez:
   - Comprendre comment fonctionne le générateur
   - Contribuer au projet
   - Étendre les templates

3. **Commencer à développer**:
   ```bash
   # Lancer l'app
   make run

   # Dans un autre terminal, tester l'API
   curl http://localhost:8080/health

   # Lire le code généré
   cat internal/domain/user/service.go
   cat internal/adapters/handlers/auth_handler.go
   ```

## Conseils

- **Lisez le code généré**: Chaque fichier est un exemple de best practice Go
- **Modifiez selon vos besoins**: La structure est un point de départ, adaptez-la
- **Suivez les patterns**: Repository, dependency injection, error handling centralisé
- **Testez régulièrement**: `make test` avant chaque commit
- **Utilisez le linter**: `make lint` pour maintenir la qualité

Bon développement! <i class="material-icons success">rocket_launch</i>
