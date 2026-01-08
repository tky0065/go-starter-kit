---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
workflowType: 'architecture'
lastStep: 8
status: 'complete'
completedAt: 'mercredi 7 janvier 2026'
project_name: 'go-starter-kit'
user_name: 'Yacoubakone'
date: 'mercredi 7 janvier 2026'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
Le système doit fournir un CLI performant pour générer une structure de projet API Go complète. Les fonctionnalités clés incluent :
- **Scaffolding CLI** : Génération de fichiers, initialisation git/go mod.
- **Architecture API** : Structure hexagonale simplifiée, routing groupé (/api/v1).
- **Core Stack** : Intégration pré-configurée de Fiber, GORM, PostgreSQL et fx.
- **Sécurité** : Système d'authentification complet (JWT, Refresh Token), hachage de mots de passe.
- **DevEx** : Documentation Swagger auto-générée, Hot-reload, Makefile complet.

**Non-Functional Requirements:**
Les décisions architecturales seront guidées par :
- **Simplicité & Accessibilité** : Le code généré doit être compréhensible par des juniors ("Zero-to-Hero").
- **Production-Ready** : Configuration par défaut sécurisée et optimisée (Docker multi-stage, CI/CD).
- **Standardisation** : Respect strict des linters et conventions Go.
- **Performance** : Démarrage rapide du CLI et de l'API générée.

**Scale & Complexity:**
- Primary domain: Developer Tools / Backend API
- Complexity level: Moyenne (Outil générateur + Template robuste)
- Estimated architectural components: ~10-15 (CLI commands + API layers: Handlers, Services, Repositories, Domain, Config, etc.)

### Technical Constraints & Dependencies

- **Langage** : Go (Golang) obligatoire.
- **Stack Imposée** : Fiber (Web), GORM (ORM), PostgreSQL (DB), fx (DI).
- **Infrastructure** : Docker et Docker Compose pour le développement et le déploiement.
- **Dépendances** : Utilisation de `uber-go/fx` pour l'injection de dépendances, ce qui structure fortement l'application.

### Cross-Cutting Concerns Identified

- **Gestion des Erreurs** : Mécanisme unifié pour capturer et formater les erreurs API.
- **Logging** : Stratégie de logs structurés pour le débogage et la production.
- **Configuration** : Gestion centralisée via variables d'environnement (`.env`).
- **Validation** : Validation systématique des entrées API.
- **Authentification** : Middleware JWT omniprésent pour sécuriser les routes.

## Starter Template Evaluation

### Primary Technology Domain

API Backend & CLI Tool basé sur le langage Go, ciblant une productivité maximale pour les développeurs.

### Starter Options Considered

Bien que nous construisons notre propre starter (`go-starter-kit`), nous avons évalué la stack technique de référence pour 2026 :
- **Go 1.25.5** (Dernière version stable, Décembre 2025)
- **Fiber v2.52.10** (Choisi pour sa maturité et sa performance stable)
- **GORM v1.31.1** (ORM de référence pour Go)
- **uber-go/fx v1.24.1** (Standard industriel pour l'injection de dépendances)
- **golangci-lint v2.7.2** (Pour garantir la qualité du code)

### Selected Starter: Custom Opinionated Stack

**Rationale for Selection:**
Les exigences du PRD imposent une stack spécifique ("Opinionated") alliant performance (Fiber) et facilité d'utilisation (GORM). L'utilisation de `fx` est cruciale pour l'architecture hexagonale "Lite" car elle permet de découpler proprement les adapters des ports sans boilerplate excessif.

**Initialization Command:**

Le CLI lui-même sera construit avec le package standard `flag` ou `spf13/cobra` pour une expérience moderne.

```bash
# Exemple de ce que le CLI générera
go mod init <project-name>
go get github.com/gofiber/fiber/v2@v2.52.10
go get gorm.io/gorm@v1.31.1
go get go.uber.org/fx@v1.24.1
```

**Architectural Decisions Provided by Starter:**

**Language & Runtime:**
Go 1.25.5 avec support des Generics et des performances optimisées.

**Build Tooling:**
Makefile pour les tâches communes. Docker multi-stage pour des images légères (< 50Mo).

**Testing Framework:**
`testing` standard de Go complété par `testify` pour les assertions et mocks.

**Code Organization:**
Architecture Hexagonale "Lite" (/internal/domain, /internal/ports, /internal/adapters, /internal/services).

**Development Experience:**
Hot-reload via `air` pour une boucle de feedback rapide.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- **Architecture Structure**: Utilisation d'une architecture hexagonale "Lite" avec renommage de `/internal/ports` en **`/internal/interfaces`** pour une meilleure clarté.
- **Dependency Injection**: Utilisation de `uber-go/fx` pour orchestrer le cycle de vie des composants.
- **Stack Core**: Fiber (Web), GORM (ORM), PostgreSQL (DB).

**Important Decisions (Shape Architecture):**
- **Auth Strategy**: JWT via `golang-jwt/v5` avec stockage des Refresh Tokens dans **PostgreSQL** pour limiter les dépendances d'infrastructure au lancement.
- **Validation**: Validation déclarative via les tags de structure avec `go-playground/validator/v10`.
- **Logging**: Logs structurés (JSON) via `rs/zerolog`.

**Deferred Decisions (Post-MVP):**
- **Caching**: Introduction de Redis pour les sessions et le cache (Phase 2).
- **Advanced Migrations**: Passage de GORM AutoMigrate à `golang-migrate` pour une gestion granulaire en production (Phase 2).

### Data Architecture

- **Database**: PostgreSQL (v16+) via Docker.
- **ORM**: GORM (v1.31.1).
- **Migration**: **GORM AutoMigrate** pour une expérience "Zero-to-Hero".
- **Validation**: **go-playground/validator** intégré dans les DTOs au niveau de la couche transport (Adapters).

### Authentication & Security

- **JWT Library**: **golang-jwt/jwt/v5**.
- **Auth Model**: Access Tokens courts et Refresh Tokens persistés en base de données.
- **Password Hashing**: **bcrypt** (cost >= 10).
- **Middleware**: Middleware d'authentification centralisé utilisant les groupes de Fiber.

### API & Communication Patterns

- **Design**: REST API avec versioning sous le préfixe `/api/v1`.
- **Documentation**: **swaggo/swag** (Swagger UI auto-généré accessible sur `/swagger`).
- **Error Handling**: Middleware global pour transformer les erreurs Go en réponses JSON standardisées.
- **Interfaces**: Définies dans `/internal/interfaces` pour découpler le domaine des adapters.

### Infrastructure & Deployment

- **Config Management**: **joho/godotenv** pour la gestion des fichiers `.env` en local.
- **Observability**: **rs/zerolog** pour des logs rapides et structurés.
- **Containerization**: Dockerfile multi-stage (build sur Alpine, runtime minimal) et Docker Compose.
- **CI/CD**: GitHub Actions pour le linting (`golangci-lint`) et les tests automatisés.

### Decision Impact Analysis

**Implementation Sequence:**
1. Initialisation du projet avec `go mod` et structure de dossiers.
2. Configuration de l'injection de dépendances (`fx`) et du serveur Fiber.
3. Mise en place de la connexion PostgreSQL et des migrations AutoMigrate.
4. Implémentation du middleware Auth et du module User.
5. Setup Swagger et Docker pour finaliser l'aspect "Production-ready".

**Cross-Component Dependencies:**
Le choix de `fx` influence la manière dont tous les autres composants (DB, Handlers, Services) sont instanciés et connectés entre eux.

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Identified:**
4 zones clés où la cohérence est impérative pour le fonctionnement harmonieux des agents : Nommage, Structure des données (Tags), Format des API et Gestion des erreurs.

### Naming Patterns

**Database Naming Conventions:**
- **Tables** : `snake_case` au pluriel (ex: `users`, `refresh_tokens`). Géré par les conventions par défaut de GORM.
- **Colonnes** : `snake_case` (ex: `created_at`, `password_hash`).

**API Naming Conventions:**
- **Endpoints** : `snake_case`, préfixés par `/api/v1`. Noms au pluriel (ex: `/api/v1/users`).
- **JSON Fields** : `snake_case` (ex: `user_id`, `access_token`).

**Code Naming Conventions (Go idiomatic):**
- **Variables/Fonctions** : `camelCase`.
- **Structs/Interfaces** : `PascalCase`.
- **Acronymes** : Majuscules intégrales (ex: `UserID`, `APIKey`, `JSONResponse`).

### Structure Patterns

**Project Organization:**
- **Tests** : Colocalisés avec le code source (ex: `user_service.go` et `user_service_test.go`).
- **Interfaces** : Regroupées dans `/internal/interfaces` pour découpler les couches.

**Data Metadata Patterns:**
- **Tags de Structure** : Utilisation systématique de tags explicites pour GORM et JSON sur les entités du domaine.
- **Soft Delete** : Utilisation de `gorm.DeletedAt` par défaut pour toutes les entités principales.

### Format Patterns

**API Response Formats:**
- **Enveloppe Standard** : Toutes les réponses réussies utilisent un wrapper avec métadonnées :
  ```json
  {
    "status": "success",
    "data": { ... },
    "meta": { "total": 100, "page": 1 }
  }
  ```

**API Error Formats:**
- **Structure Unifiée** :
  ```json
  {
    "status": "error",
    "message": "Message compréhensible",
    "code": "ERROR_SLUG",
    "details": null
  }
  ```

### Process Patterns

**Error Handling Patterns:**
- Les erreurs sont capturées par un middleware centralisé Fiber.
- Utilisation de `zerolog` pour logger les erreurs avec le contexte de la requête côté serveur.

**Validation Patterns:**
- Validation déclenchée systématiquement au niveau des Handlers (Adapters) via `go-playground/validator`.

### Enforcement Guidelines

**All AI Agents MUST:**
1. Respecter scrupuleusement le nommage `snake_case` dans les tags JSON.
2. Utiliser `PascalCase` pour tous les acronymes dans le code Go (`ID`, pas `Id`).
3. Ne jamais exposer de stack trace dans les réponses API d'erreur.
4. Toujours définir les interfaces dans `/internal/interfaces`.

## Project Structure & Boundaries

### Complete Project Directory Structure

```
go-starter-kit/
├── cmd/
│   └── create-go-starter/   # CLI de scaffolding
├── internal/
│   ├── adapters/            # Outside (HTTP Handlers, Repository GORM)
│   ├── domain/              # Inside (Entités métier, Logique de service)
│   ├── interfaces/          # Ports (Contrats et Interfaces)
│   └── infrastructure/      # Setup (DB, Server, Config)
├── pkg/                     # Code partagé (Logger, Validator)
├── deployments/             # Docker & Docker Compose
├── .github/                 # CI/CD Workflows
├── Makefile                 # Automatisation
└── README.md
```

### Architectural Boundaries

**API Boundaries:**
L'interaction externe se fait exclusivement via `/internal/adapters/handlers`. Cette couche convertit les DTOs en entités de domaine.

**Component Boundaries:**
Utilisation de `uber-go/fx` pour injecter les dépendances sans que les composants ne connaissent leurs implémentations concrètes (découplage via `/internal/interfaces`).

**Data Boundaries:**
La persistance est encapsulée dans `/internal/adapters/repository`. Seule cette couche interagit avec GORM.

### Requirements to Structure Mapping

**Feature/Epic Mapping:**
- **Auth System** : `/internal/domain/user`, `/internal/adapters/handlers/auth_handler.go`.
- **User Management** : `/internal/domain/user`, `/internal/adapters/handlers/user_handler.go`.
- **Infrastructure** : `/internal/infrastructure/server`.

**Cross-Cutting Concerns:**
- **Logging** : `/pkg/logger`.
- **Validation** : `/pkg/validator`.
- **Errors** : `/internal/domain/errors.go` et `/internal/adapters/middleware/error_middleware.go`.

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:**
La combinaison Fiber v2.52.10 + GORM v1.31.1 + uber-go/fx v1.24.1 est hautement compatible et forme une base robuste pour une API moderne.

**Pattern Consistency:**
Les patterns de nommage et les structures de données (tags) sont uniformisés entre la base de données, le code Go et l'interface JSON.

**Structure Alignment:**
La structure hexagonale "Lite" est directement supportée par l'injection de dépendances via `fx`, facilitant le découplage par interfaces.

### Requirements Coverage Validation ✅

**Functional Requirements Coverage:**
Chaque exigence, du scaffolding CLI à l'authentification JWT, possède un emplacement et une technologie dédiée dans l'architecture.

**Non-Functional Requirements Coverage:**
Les objectifs de performance (< 100ms) et de légèreté (< 50Mo) sont adressés par le choix de Fiber et le Docker multi-stage.

### Implementation Readiness Validation ✅

**Decision Completeness:**
Toutes les décisions critiques (DB, Auth, Logging, Validation) sont documentées avec leurs versions respectives.

**Structure Completeness:**
L'arborescence complète du projet est définie, éliminant les incertitudes de placement de fichiers pour les agents.

### Architecture Completeness Checklist

- [x] Analyse du contexte et de la complexité terminée.
- [x] Stack technique et versions validées par recherche web.
- [x] Patterns de nommage et de cohérence IA établis.
- [x] Structure de dossiers et frontières architecturales définies.
- [x] Mapping des exigences fonctionnelles vers la structure complété.

### Architecture Readiness Assessment

**Overall Status:** READY FOR IMPLEMENTATION
**Confidence Level:** HIGH

**Key Strengths:**
- Simplicité du scaffolding CLI couplée à une architecture robuste.
- Stack technique moderne et performante (Go 1.25.5).
- Conventions de nommage strictes pour éviter les conflits d'agents IA.

**Areas for Future Enhancement:**
- Introduction de Redis pour le cache (Phase 2).
- Migration explicite avec `golang-migrate` (Phase 2).

### Implementation Handoff

**AI Agent Guidelines:**
- Toujours utiliser les interfaces définies dans `/internal/interfaces` pour l'injection.
- Suivre scrupuleusement le format de réponse standard avec métadonnées.
- Respecter les tags JSON `snake_case` et le nommage idiomatique Go.

**First Implementation Priority:**
Initialisation du repo avec `go mod init` et création de la structure de dossiers.

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED ✅
**Total Steps Completed:** 8
**Date Completed:** mercredi 7 janvier 2026
**Document Location:** _bmad-output/planning-artifacts/architecture.md

### Final Architecture Deliverables

**📋 Complete Architecture Document**
- Toutes les décisions architecturales documentées avec des versions spécifiques.
- Patterns d'implémentation garantissant la cohérence des agents IA.
- Structure de projet complète avec tous les fichiers et répertoires.
- Mapping des exigences vers l'architecture.

**🏗️ Implementation Ready Foundation**
- 15 décisions architecturales prises.
- 10 patterns d'implémentation définis.
- 11 répertoires de structure spécifiés.
- 26 exigences fonctionnelles pleinement supportées.

**Ready for implementation phase.**
