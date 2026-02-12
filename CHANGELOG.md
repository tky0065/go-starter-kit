# Changelog

Toutes les modifications notables de ce projet seront documentées dans ce fichier.

Le format est basé sur [Keep a Changelog](https://keepachangelog.com/fr/1.0.0/),
et ce projet adhère au [Semantic Versioning](https://semver.org/lang/fr/).

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

- **Pluralisation**: Règles simples (ajoute 's'). Pluriels irréguliers nécessitent édition manuelle
- **Relations many-to-many**: Pas encore supportées (prévu v1.3.0)
- **Validations custom**: Doivent être ajoutées manuellement après génération
- **Swagger**: Doit être regénéré manuellement (`make swagger`) après add-model

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
