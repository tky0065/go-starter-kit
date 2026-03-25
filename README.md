# create-go-starter

[![Release](https://img.shields.io/github/v/release/tky0065/go-starter-kit)](https://github.com/tky0065/go-starter-kit/releases)
[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/github/license/tky0065/go-starter-kit)](LICENSE)

Un outil CLI puissant pour générer des projets Go prêts pour la production en quelques secondes.

## Version actuelle

**v1.5.2** - Stable et prêt pour la production

 **Production ready** - Utilisé dans des projets réels
 **3 templates** - Minimal, Full (JWT), GraphQL
 **Observabilité avancée** - Prometheus, Jaeger, Grafana, Health checks K8s
 **CLI interactif** - Mode guidé, dry-run, diagnostics, alias courts
 **Bien testé** - Tests unitaires et E2E
 **Documentation complète** - Guides et exemples
 **Open source** - MIT License

## Aperçu

`create-go-starter` est un générateur de projets Go qui crée une architecture hexagonale complète avec toutes les fonctionnalités essentielles d'une application backend moderne. En une seule commande, obtenez un projet structuré avec authentification JWT, API REST, base de données, tests, et configuration Docker prête pour le déploiement.

### Fonctionnalités incluses

- **Architecture hexagonale** (Ports & Adapters) - Séparation claire des responsabilités
- **Authentification JWT** - Access tokens + Refresh tokens avec rotation sécurisée
- **API REST** avec Fiber v2 - Framework web haute performance
- **Base de données** - GORM avec PostgreSQL et migrations automatiques
- **Injection de dépendances** - uber-go/fx pour une architecture modulaire
- **Tests complets** - Tests unitaires et d'intégration
- **Documentation Swagger** - API documentée automatiquement avec OpenAPI
- **Docker** - Build multi-stage optimisé et docker-compose
- **CI/CD** - Pipeline GitHub Actions pré-configuré
- **Logging structuré** - rs/zerolog pour des logs professionnels
- **Validation** - go-playground/validator pour valider les entrées
- **Makefile** - Commandes utiles pour dev, test, build et déploiement
- **Observabilité avancée** - Prometheus `/metrics`, distributed tracing Jaeger/OpenTelemetry, health checks K8s (`/health/liveness`, `/health/readiness`), dashboard Grafana pré-configuré (`--observability=advanced`)

## Installation rapide

### Méthode 1: Installation directe (Recommandée)

Installation globale en une seule commande, sans cloner le repository:

```bash
# Version stable (recommandée)
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@v1.5.2

# Ou dernière version
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest

# Ou version de développement
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@develop
```

Le binaire sera installé dans `$GOPATH/bin` (généralement `~/go/bin`). Assurez-vous que ce répertoire est dans votre PATH.

**Note**: Si `create-go-starter` n'est pas reconnu après l'installation, ajoutez `$GOPATH/bin` à votre PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Vérifier l'installation**:

```bash
create-go-starter --help
# Devrait afficher l'aide du CLI
```

### Méthode 2: Build depuis les sources

Pour contributeurs ou personnalisation:

```bash
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit
go build -o create-go-starter ./cmd/create-go-starter
# Le binaire est maintenant disponible: ./create-go-starter
```

### Méthode 3: Build avec Makefile

```bash
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit
make build
# Le binaire est disponible: ./create-go-starter
```

Pour plus de détails, consultez le [guide d'installation complet](./docs/installation.md).

## Utilisation de base

### Créer un nouveau projet

```bash
create-go-starter mon-super-projet
```

Cette commande va:
1. Créer la structure complète du projet
2. Générer ~30+ fichiers (handlers, services, repositories, tests, etc.)
3. Configurer tous les fichiers nécessaires (.env, Dockerfile, Makefile, etc.)
4. Copier le fichier `.env.example` vers `.env`
5. Initialiser un dépôt Git avec un commit initial (si Git est disponible)

### Choisir un template

Par défaut, `create-go-starter` génère un projet **full** avec authentification JWT et gestion des utilisateurs. Vous pouvez choisir un template différent avec le flag `--template`:

```bash
create-go-starter mon-projet --template minimal    # API REST basique avec Swagger
create-go-starter mon-projet --template full       # API complète avec JWT auth (défaut)
create-go-starter mon-projet --template graphql    # API GraphQL avec gqlgen
```

**Templates disponibles**:

| Template | Description | Cas d'usage |
|----------|-------------|-------------|
| `minimal` | API REST basique avec Swagger (sans authentification) | Prototypes rapides, APIs publiques simples |
| `full` | API complète avec JWT auth, gestion utilisateurs et Swagger | Applications backend complètes (défaut) |
| `graphql` | API GraphQL avec gqlgen et GraphQL Playground | Applications nécessitant GraphQL |

Pour plus de détails sur les différences entre templates, consultez le [guide d'utilisation](./docs/usage.md#templates-disponibles).

### Choisir une base de données

`create-go-starter` supporte **PostgreSQL** (défaut), **MySQL**, et **SQLite**. Utilisez le flag `--database` pour sélectionner votre base de données:

```bash
# PostgreSQL (défaut) - Production, requêtes complexes
create-go-starter mon-app
create-go-starter mon-app --database=postgres

# MySQL - Compatibilité large, shared hosting
create-go-starter mon-app --database=mysql

# SQLite - Prototypage rapide, zéro configuration
create-go-starter mon-app --database=sqlite
```

**Comparaison rapide**:

| Database | Setup | Idéal pour |
|----------|-------|-----------|
| **PostgreSQL** | Docker | Production, requêtes complexes, fiabilité |
| **MySQL** | Docker | Shared hosting, compatibilité, read-heavy |
| **SQLite** | Aucun | Prototypage, petites apps, développement |

**Voir le [guide complet des databases](./docs/databases.md)** pour des détails sur:
- Comparaison détaillée des features
- Avantages et limitations
- Cas d'usage recommandés
- Guide de migration entre databases

### Sélectionner un framework web (Nouveau! v1.6.0)

`create-go-starter` supporte plusieurs frameworks web avec le flag `--framework` (ou `-f`):

```bash
# Fiber (défaut) — Framework performant inspiré d'Express
create-go-starter mon-app

# Avec flag explicite
create-go-starter mon-app --framework=fiber
create-go-starter mon-app -f fiber
```

| Framework | Description | Statut |
|-----------|-------------|--------|
| `fiber` | Fiber v2 - Fast HTTP framework inspired by Express **(défaut)** | ✅ Disponible |
| `gin` | Gin - High-performance HTTP web framework | 🔜 Planifié v2.0.0 |
| `echo` | Echo - Minimalist high-performance HTTP framework | 🔜 Planifié v2.0.0 |

### Activer l'observabilité avancée (Nouveau! v1.3.0)

`create-go-starter` supporte **3 niveaux d'observabilité** avec le flag `--observability`:

```bash
# Pas d'observabilité (défaut) — comportement actuel
create-go-starter mon-app

# Observabilité avancée — Prometheus, Jaeger, Grafana, Health checks K8s
create-go-starter mon-app --observability=advanced

# Combinaison complète
create-go-starter mon-app --template=full --database=postgres --observability=advanced
```

| Niveau | Description |
|--------|-------------|
| `none` | Aucune observabilité (défaut) |
| `basic` | `/health` amélioré (sans Prometheus) |
| `advanced` | Stack complète: Prometheus, Jaeger, Grafana, Health checks K8s |

**Avec `--observability=advanced`, vous obtenez:**

- **Prometheus Metrics** — Endpoint `/metrics` avec `http_requests_total`, `http_request_duration_seconds`, `http_requests_in_flight` via fiberprometheus
- **Distributed Tracing** — OpenTelemetry avec export OTLP/gRPC vers Jaeger, propagation W3C traceparent, trace_id/span_id dans les logs zerolog
- **Health Checks K8s** — `/health/liveness` (toujours 200), `/health/readiness` (vérifie DB avec timeout 2s, retourne 503 si down), `/health` alias rétrocompatible
- **Grafana Dashboard** — Dashboard JSON 7 panneaux pré-configuré, auto-provisioning datasources + dashboards
- **Docker Compose étendu** — Stack complète: PostgreSQL + Jaeger (1.56.0) + Prometheus (v2.51.0) + Grafana (10.4.0)
- **Kubernetes Probes** — `deployments/kubernetes/probes.yaml` généré automatiquement
- **Alerting Prometheus** — Règles d'alerte pré-configurées (latence, erreurs, saturation)

Pour plus de détails, consultez le [guide d'observabilité](./docs/usage.md#observabilité---observability).

### Améliorations CLI (Nouveau! v1.4.0)

La v1.4.0 apporte des améliorations majeures à l'expérience utilisateur du CLI: mode interactif guidé, prévisualisation dry-run, diagnostics environnement, et alias courts pour tous les flags.

#### Mode interactif

```bash
# Lancer le mode interactif guidé
create-go-starter --interactive
create-go-starter -i
```

Le mode interactif guide l'utilisateur étape par étape pour configurer un nouveau projet: nom du projet, choix du template, base de données, et niveau d'observabilité. Un résumé est affiché avant la génération avec confirmation.
L'expérience utilise désormais un TUI plus structuré avec écran d'accueil, sélections centrées, résumé lisible et écran de progression avec statistiques en temps réel.

#### Prévisualisation dry-run

```bash
# Voir les fichiers qui seraient générés sans les créer
create-go-starter mon-app --dry-run
create-go-starter -n -t minimal -d sqlite mon-app
```

Le mode dry-run affiche la liste complète des fichiers et répertoires qui seraient créés, avec compteur et configuration, sans écrire sur le disque.

#### Diagnostics environnement

```bash
# Vérifier que l'environnement est prêt
create-go-starter doctor
```

La commande `doctor` vérifie la version de Go (>= 1.21), Git, et Docker (binaire + daemon), et affiche un rapport clair avec le statut de chaque outil.

#### Alias courts

Tous les flags disposent maintenant d'alias courts pour une saisie plus rapide:

| Flag | Alias | Exemple |
|------|-------|---------|
| `--template` | `-t` | `-t minimal` |
| `--database` | `-d` | `-d sqlite` |
| `--observability` | `-o` | `-o advanced` |
| `--interactive` | `-i` | `-i` |
| `--dry-run` | `-n` | `-n` |
| `--help` | `-h` | `-h` |

```bash
# Exemples avec alias courts
create-go-starter -t graphql -d postgres mon-app
create-go-starter -n -t minimal -d sqlite mon-app
create-go-starter -t full -d postgres -o advanced mon-app
```

#### FAQ - Choix de base de données

**Quelle base de données dois-je choisir?**  
PostgreSQL (défaut) pour la production, MySQL pour la compatibilité large, SQLite pour le prototypage rapide.

**Puis-je changer de base de données plus tard?**  
Oui, mais nécessite de régénérer le projet. Voir le [guide de migration](./docs/database-migration.md).

**SQLite est-il adapté à la production?**  
Pour petite échelle (<100 utilisateurs concurrents) seulement. PostgreSQL/MySQL recommandés pour production.

**Ai-je besoin de Docker?**  
PostgreSQL et MySQL nécessitent Docker pour développement local. SQLite fonctionne sans Docker.

### Lancer le projet généré

#### Option 1: Configuration automatique (Recommandé) <i class="material-icons">rocket_launch</i>

```bash
cd mon-super-projet
./setup.sh
make run
```

Le script `setup.sh` automatise:
- Installation des dépendances Go
- Génération du JWT secret
- Configuration de PostgreSQL (Docker ou local)
- Vérification de l'installation

#### Option 2: Configuration manuelle

```bash
cd mon-super-projet

# Installer les dépendances et générer go.sum
go mod tidy

# Configurer le JWT secret dans .env
# JWT_SECRET=<générer avec: openssl rand -base64 32>

# Lancer la base de données (PostgreSQL)
docker run -d --name postgres \
  -e POSTGRES_DB=mon-super-projet \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Lancer l'application
make run
```

L'API sera disponible sur `http://localhost:8080`

```bash
# Tester le health check
curl http://localhost:8080/health
# {"status":"ok"}
```

### Ajouter des modèles avec Relations (Nouveau! v1.2.0)

Une fois votre projet créé, vous pouvez ajouter de nouveaux modèles avec la commande `add-model`:

```bash
# Ajouter un modèle simple
create-go-starter add-model Todo --fields "title:string,completed:bool"

# Ajouter un modèle avec relation BelongsTo (enfant → parent)
create-go-starter add-model Comment --fields "content:string" --belongs-to Todo

# Ajouter un modèle avec relation HasMany (parent → enfants)
create-go-starter add-model Category --fields "name:string:unique" --has-many Product
```

**Fonctionnalités générées automatiquement:**
- <i class="material-icons success small">check</i> Model avec tags GORM (`internal/models/`)
- <i class="material-icons success small">check</i> Repository interface et implémentation
- <i class="material-icons success small">check</i> Service avec logique métier
- <i class="material-icons success small">check</i> Handlers HTTP avec endpoints CRUD complets
- <i class="material-icons success small">check</i> Tests unitaires (service + handler)
- <i class="material-icons success small">check</i> Routes automatiquement ajoutées
- <i class="material-icons success small">check</i> Support des relations (foreign keys, preloading)

**Relations supportées:**
- `--belongs-to <Parent>` - Ajoute foreign key + champ relation (ex: `Comment` belongs to `Todo`)
- `--has-many <Child>` - Modifie le parent pour ajouter slice d'enfants (ex: `Todo` has many `Comment`)

**Exemple complet avec relations:**

```bash
# 1. Créer le projet
create-go-starter blog-api

# 2. Ajouter un modèle Category
cd blog-api
create-go-starter add-model Category --fields "name:string:unique"

# 3. Ajouter Post qui appartient à Category
create-go-starter add-model Post --fields "title:string,content:string" --belongs-to Category

# 4. Ajouter Comment qui appartient à Post
create-go-starter add-model Comment --fields "content:string,author:string" --belongs-to Post

# 5. Lancer et tester
make run
curl http://localhost:8080/api/v1/posts?include=category  # Preload relation
```

**Note:** La pluralisation utilise des règles simples. Pour les pluriels irréguliers (Person→People, Child→Children), éditez manuellement le code généré.

## Structure générée

Voici ce que `create-go-starter` génère pour vous:

```
mon-super-projet/
├── cmd/
│   └── main.go                    # Point d'entrée avec fx dependency injection
├── internal/
│   ├── models/                    # Entités de domaine partagées
│   │   └── user.go                # User, RefreshToken, AuthResponse
│   ├── domain/                    # Couche domaine (logique métier)
│   │   ├── user/                  # Domaine User
│   │   │   ├── service.go         # Logique métier (Register, Login, etc.)
│   │   │   └── module.go          # Module fx
│   │   └── errors.go              # Erreurs métier personnalisées
│   ├── adapters/                  # Adapters (HTTP, DB)
│   │   ├── handlers/              # HTTP handlers
│   │   │   ├── auth_handler.go    # Auth endpoints (register, login, refresh)
│   │   │   └── user_handler.go    # User CRUD endpoints
│   │   ├── middleware/            # Middleware Fiber
│   │   │   ├── auth_middleware.go # JWT verification
│   │   │   └── error_handler.go   # Gestion centralisée des erreurs
│   │   ├── repository/            # Implémentation des repositories
│   │   │   └── user_repository.go # GORM implementation
│   │   └── http/                  # HTTP utilities
│   │       ├── health.go          # Handler health check
│   │       └── routes.go          # Routes centralisées
│   ├── infrastructure/            # Infrastructure
│   │   ├── database/              # Configuration DB (GORM, migrations)
│   │   └── server/                # Configuration Fiber app
│   └── interfaces/                # Ports (interfaces)
│       └── user_repository.go     # Interface UserRepository
├── pkg/                           # Packages réutilisables
│   ├── auth/                      # JWT utilities
│   ├── config/                    # Chargement configuration (.env)
│   └── logger/                    # Configuration zerolog
├── .github/workflows/
│   └── ci.yml                     # Pipeline CI/CD (lint, test, build)
├── .env                           # Configuration (copié depuis .env.example)
├── .env.example                   # Template de configuration
├── .gitignore                     # Exclusions Git
├── .golangci.yml                  # Configuration du linter
├── Dockerfile                     # Build multi-stage pour production
├── Makefile                       # Commandes utiles (run, test, lint, docker, etc.)
├── setup.sh                       # Script de configuration automatique
├── go.mod                         # Module Go avec dépendances
└── README.md                      # Documentation du projet
```

Pour une explication détaillée de chaque composant, consultez le [guide d'utilisation](./docs/usage.md).

## Stack technique

Les projets générés utilisent les meilleures bibliothèques de l'écosystème Go:

| Composant | Bibliothèque | Version | Description |
|-----------|-------------|---------|-------------|
| Web Framework | [Fiber](https://gofiber.io/) | v2 | Framework HTTP rapide, inspiré d'Express |
| ORM | [GORM](https://gorm.io/) | v1 | ORM Go avec support PostgreSQL |
| Dependency Injection | [fx](https://uber-go.github.io/fx/) | latest | DI framework par Uber |
| Logging | [zerolog](https://github.com/rs/zerolog) | latest | Logger structuré haute performance |
| JWT | [golang-jwt](https://github.com/golang-jwt/jwt) | v5 | Tokens JWT pour authentification |
| Validation | [validator](https://github.com/go-playground/validator) | v10 | Validation de structs |
| Swagger | [swaggo](https://github.com/swaggo/swag) | latest | Documentation API OpenAPI |
| Crypto | golang.org/x/crypto | latest | Hashage bcrypt pour mots de passe |

## Documentation

### Guides essentiels

- **[Guide d'installation](./docs/installation.md)** - Installation détaillée avec toutes les méthodes
- **[Guide d'utilisation](./docs/usage.md)** - Utilisation du CLI et structure complète générée
- **[Guide des projets générés](./docs/generated-project-guide.md)** - Guide complet pour développer avec les projets créés (architecture, API, tests, déploiement)
- **[Changelog](./CHANGELOG.md)** - Historique des versions et nouveautés
- **[Roadmap](./ROADMAP.md)** - Prochaines fonctionnalités et vision long-terme

### Documentation avancée

- **[Architecture du CLI](./docs/cli-architecture.md)** - Documentation technique pour contributeurs
- **[Guide de contribution](./docs/contributing.md)** - Comment contribuer au projet

### Site de documentation

Le projet utilise **Material for MkDocs** pour générer le site de documentation officiel.

**Voir la documentation en ligne:** [https://tky0065.github.io/go-starter-kit/](https://tky0065.github.io/go-starter-kit/)

**Travailler avec la documentation localement:**

```bash
# Activer l'environnement Python
source venv/bin/activate

# Serveur de développement avec rechargement automatique
mkdocs serve
# Accès: http://127.0.0.1:8000/go-starter-kit/

# Construire le site
mkdocs build --clean

# Déployer sur GitHub Pages
mkdocs gh-deploy
```

**Note pour les contributeurs:** Utilisez la syntaxe HTML pour les icônes Material Design, pas la syntaxe emoji. Voir [CLAUDE.md](./CLAUDE.md) pour les détails.

## Démarrage rapide en 30 secondes

```bash
# 1. Installer l'outil
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest

# 2. Créer un projet
create-go-starter mon-projet

# 3. Configuration automatique
cd mon-projet
./setup.sh

# 4. Lancer
make run

# 5. Tester
curl http://localhost:8080/health
```

Ou configuration manuelle:
```bash
cd mon-projet
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
docker run -d --name postgres -e POSTGRES_DB=mon-projet -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16-alpine
make run
```

## Exemples d'utilisation de l'API

Une fois votre projet lancé, vous pouvez tester l'API:

```bash
# Créer un utilisateur
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securePass123"}'

# Se connecter
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securePass123"}'

# Utiliser le token retourné pour accéder aux endpoints protégés
TOKEN="<access_token_from_login_response>"
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

Pour plus d'exemples et la documentation complète de l'API, consultez le [guide des projets générés](./docs/generated-project-guide.md#api-reference).

## Commandes Makefile disponibles

Les projets générés incluent un Makefile avec des commandes utiles:

```bash
make help           # Afficher toutes les commandes disponibles
make run            # Lancer l'application
make build          # Build le binaire
make test           # Exécuter les tests
make test-coverage  # Tests avec rapport de coverage
make lint           # Linter le code (golangci-lint)
make docker-build   # Build l'image Docker
make docker-run     # Lancer le conteneur Docker
make clean          # Nettoyer les artifacts
```

## Prérequis

- **Go 1.25 ou supérieur** - [Télécharger Go](https://golang.org/dl/)
- **PostgreSQL** - Pour les projets générés (peut être lancé via Docker)
- **Git** - Pour cloner et contribuer
- **Docker** (optionnel) - Pour lancer PostgreSQL et containeriser l'application
- **golangci-lint** (optionnel) - Pour le linting

## Pourquoi create-go-starter?

### Gain de temps

Au lieu de passer des heures à configurer:
- L'architecture du projet
- L'authentification JWT
- La connexion à la base de données
- Les tests
- Le Docker et CI/CD
- La documentation Swagger

Obtenez tout cela en **une seule commande** et commencez immédiatement à développer vos fonctionnalités métier.

### Best practices intégrées

- **Architecture hexagonale** - Séparation claire entre domaine, adapters et infrastructure
- **Dependency injection** - Code testable et modulaire
- **Error handling centralisé** - Gestion cohérente des erreurs
- **Security-first** - JWT, bcrypt, validation, CORS
- **Tests** - Exemples de tests unitaires et d'intégration
- **Clean code** - Respect des conventions Go et linting strict

### Production-ready

Les projets générés sont prêts pour la production:
- Build Docker multi-stage optimisé
- CI/CD avec tests automatiques
- Logging structuré pour monitoring
- Configuration par environnement
- Health checks
- Graceful shutdown

## Contribuer

Les contributions sont les bienvenues! Consultez le [guide de contribution](./docs/contributing.md) pour commencer.

### Processus de contribution

1. Fork le projet
2. Créer une branche (`git checkout -b feature/ma-fonctionnalite`)
3. Commit les changements (`git commit -m 'feat: ajouter une fonctionnalité'`)
4. Push vers la branche (`git push origin feature/ma-fonctionnalite`)
5. Ouvrir une Pull Request

## FAQ

### Puis-je changer de base de données plus tard?
Oui! Voir le [guide de migration](./docs/database-migration.md) pour les instructions détaillées. Vous devrez régénérer le projet avec le flag `--database` et migrer vos données.

### Quelle base de données choisir pour mon projet?
- **PostgreSQL**: Pour la production, queries complexes, fiabilité
- **MySQL**: Pour shared hosting, compatibilité large
- **SQLite**: Pour prototypage, développement, petites apps

Voir le [guide des databases](./docs/databases.md) pour une comparaison complète.

### SQLite est-il adapté pour la production?
SQLite peut être utilisé en production pour des petites applications (<100 utilisateurs concurrents), mais PostgreSQL ou MySQL sont recommandés pour une croissance attendue.

### Ai-je besoin de Docker?
- **PostgreSQL**: Oui (pour développement local)
- **MySQL**: Oui (pour développement local)
- **SQLite**: Non (base de données embarquée)

### Puis-je utiliser plusieurs templates et databases?
Oui! Toutes les combinaisons sont possibles:
```bash
create-go-starter app --template=minimal --database=sqlite
create-go-starter app --template=full --database=mysql
create-go-starter app --template=graphql --database=postgres
```

## Roadmap

**Fonctionnalités complétées** (v1.0 -- v1.4):

- [x] **Templates multiples** - Trois templates disponibles (minimal, full, graphql)
- [x] **Multi-database support** - PostgreSQL, MySQL, SQLite avec migration guides
- [x] **CRUD Scaffolding** - Commande `add-model` avec relations BelongsTo/HasMany
- [x] **Observabilité avancée** - Prometheus, Jaeger/OpenTelemetry, Grafana, Health checks K8s
- [x] **CLI interactif** - Mode guidé étape par étape (`--interactive`)
- [x] **Dry-run** - Prévisualisation des fichiers sans génération (`--dry-run`)
- [x] **Diagnostics** - Commande `doctor` pour vérifier l'environnement
- [x] **Alias courts** - `-t`, `-d`, `-o`, `-i`, `-n`, `-h` pour tous les flags
- [x] **Feedback visuel** - Barre de progression et statistiques post-génération

**Fonctionnalités prévues** (v1.5+):

- [ ] Support NoSQL (MongoDB) - en attente de demande communautaire
- [ ] Choix du framework web (Gin, Echo, Chi)
- [ ] Génération de microservices
- [ ] Templates de tests E2E avancés
- [ ] Configuration Kubernetes avancée

## Licence

[MIT License](LICENSE) - Libre d'utilisation pour projets personnels et commerciaux.

## Support

- **Issues**: [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)
- **Documentation**: [docs/](./docs/)

## Remerciements

Construit avec les excellentes bibliothèques de la communauté Go. Merci aux mainteneurs de Fiber, GORM, fx, zerolog et toutes les autres dépendances.

---

**Fait avec <i class="material-icons small" style="color:#e25555">favorite</i> pour la communauté Go**

Commencez à construire votre prochaine application backend en secondes, pas en jours!
