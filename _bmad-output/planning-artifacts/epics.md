---
stepsCompleted: ['step-01-validate-prerequisites', 'step-02-design-epics']
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
  - _bmad-output/planning-artifacts/architecture.md
---

# go-starter-kit - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for go-starter-kit, decomposing the requirements from the PRD, UX Design if it exists, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

**FR1:** L'utilisateur peut installer l'outil via une commande `go install`.

**FR2:** L'utilisateur peut générer un nouveau projet en fournissant un nom de projet via le CLI.

**FR3:** Le système peut créer automatiquement une structure de dossiers respectant l'architecture hexagonale lite.

**FR4:** Le système peut injecter dynamiquement le nom du projet dans les fichiers générés (fichiers Go, `go.mod`, Dockerfile).

**FR5:** Le système peut initialiser automatiquement un module Go et télécharger les dépendances nécessaires.

**FR6:** Le système peut créer un fichier `.env` par défaut à partir d'un template `.env.example`.

**FR7:** Un visiteur peut créer un compte utilisateur (Register).

**FR8:** Un utilisateur peut s'authentifier via des identifiants sécurisés (Login).

**FR9:** Le système peut générer des jetons d'accès (Access Tokens) et de renouvellement (Refresh Tokens) sécurisés.

**FR10:** Un utilisateur peut renouveler son jeton d'accès sans se reconnecter manuellement.

**FR11:** Le système peut hacher les mots de passe de manière sécurisée avant stockage.

**FR12:** Le système peut protéger des routes spécifiques pour n'autoriser que les utilisateurs authentifiés.

**FR13:** Le système peut regrouper les routes API par domaine fonctionnel.

**FR14:** Le système peut appliquer un préfixe global `/api/v1` à toutes les routes métier.

**FR15:** Le système peut gérer automatiquement les erreurs et renvoyer des réponses JSON standardisées.

**FR16:** Le système peut valider automatiquement les données d'entrée des requêtes HTTP selon des règles prédéfinies.

**FR17:** Le système peut exposer une documentation interactive (Swagger UI) mise à jour automatiquement.

**FR18:** Le système peut se connecter à une base de données PostgreSQL de manière résiliente (pool de connexions).

**FR19:** Le système peut exécuter des migrations de base de données au démarrage ou via commande.

**FR20:** Un développeur peut effectuer des opérations CRUD (Créer, Lire, Mettre à jour, Supprimer) sur l'entité User.

**FR21:** Le système peut gérer l'injection de dépendances pour tous les composants majeurs (DB, Serveur, Handlers).

**FR22:** Le système peut assurer un démarrage et un arrêt "propre" (graceful shutdown) de l'application.

**FR23:** Le développeur peut lancer l'application en mode développement avec rechargement automatique (hot-reload).

**FR24:** Le développeur peut exécuter l'ensemble des tests (unitaires et intégration) via une commande unique.

**FR25:** Le système peut être exécuté dans un environnement conteneurisé (Docker).

**FR26:** Le système peut être déployé via un pipeline de CI/CD pré-configuré (GitHub Actions).

### NonFunctional Requirements

**NFR1 - Performance:** Les endpoints de base (Auth/Health) doivent répondre en moins de 100ms (hors latence réseau).

**NFR2 - Performance:** L'application conteneurisée doit être opérationnelle en moins de 2 secondes (cold start).

**NFR3 - Performance:** L'image Docker finale doit peser moins de 50 Mo (utilisation de multi-stage builds).

**NFR4 - Security:** Toutes les données sensibles (mots de passe) sont hachées avec bcrypt (coût par défaut >= 10).

**NFR5 - Security:** Utilisation de l'algorithme JWT HS256 ou RS256 avec une gestion stricte de l'expiration.

**NFR6 - Security:** Aucune clé ou secret ne doit être en dur dans le code (utilisation obligatoire de .env ou variables d'environnement).

**NFR7 - Security:** Protection native contre les vulnérabilités courantes OWASP (CORS, CSRF, Injection SQL via l'ORM).

**NFR8 - Maintainability:** 100% du code généré doit respecter les standards golangci-lint.

**NFR9 - Maintainability:** Chaque fonction publique doit être documentée. Le README doit permettre un démarrage en 5 minutes.

**NFR10 - Maintainability:** L'architecture hexagonale doit permettre de mocker 100% des dépendances externes pour les tests unitaires.

**NFR11 - Operability:** Le système doit intercepter les signaux d'arrêt et fermer les connexions proprement en moins de 5 secondes (graceful shutdown).

**NFR12 - Operability:** Utilisation de logs structurés (JSON) en production pour faciliter l'indexation.

**NFR13 - Operability:** Présence d'un endpoint /health pour les orchestrateurs (Kubernetes/Docker).

### Additional Requirements

**Starter Template (CRITICAL for Epic 1, Story 1):**
- Custom opinionated stack: Go 1.25.5 + Fiber v2.52.10 + GORM v1.31.1 + uber-go/fx v1.24.1
- golangci-lint v2.7.2 pour garantir la qualité du code

**Infrastructure & DevOps:**
- Docker multi-stage build optimisé sur Alpine (< 50Mo)
- Docker Compose pour l'environnement de développement local
- CI/CD via GitHub Actions (linting golangci-lint + tests automatisés)
- Hot-reload via air pour l'environnement de développement
- Makefile pour automatisation des commandes (dev, test, build, lint)

**Architecture Technique:**
- Architecture hexagonale "Lite" avec renommage de /internal/ports en /internal/interfaces
- Injection de dépendances via uber-go/fx pour orchestrer le cycle de vie des composants
- Validation déclarative via tags de structure avec go-playground/validator/v10
- Logs structurés (JSON) via rs/zerolog
- Documentation Swagger auto-générée via swaggo/swag accessible sur /swagger

**Sécurité & Authentification:**
- JWT via golang-jwt/jwt/v5
- Access Tokens courts et Refresh Tokens persistés en base de données PostgreSQL
- Password hashing via bcrypt (cost >= 10)
- Middleware d'authentification centralisé utilisant les groupes de Fiber

**Data & Persistence:**
- PostgreSQL (v16+) via Docker
- ORM GORM (v1.31.1)
- GORM AutoMigrate pour une expérience "Zero-to-Hero"
- Soft delete via gorm.DeletedAt par défaut pour toutes les entités principales

**Patterns de Cohérence (CRITICAL for all agents):**
- Nommage Database: snake_case au pluriel pour tables (users, refresh_tokens), snake_case pour colonnes (created_at, password_hash)
- Nommage API: snake_case pour endpoints préfixés /api/v1, noms au pluriel (/api/v1/users)
- Nommage JSON: snake_case pour tous les champs (user_id, access_token)
- Nommage Code Go: camelCase pour variables/fonctions, PascalCase pour Structs/Interfaces, acronymes en majuscules (UserID, APIKey, JSONResponse)
- Format de réponse API standard avec enveloppe (status, data, meta)
- Format d'erreur API unifié (status, message, code, details)
- Tests colocalisés avec le code source (user_service.go et user_service_test.go)

**Structure de Projet:**
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

### FR Coverage Map

**Note:** Les FR1-FR26 ont été entièrement implémentées dans les Epics 1-6 (v1.0.0 - COMPLÉTÉS).
Les epics suivants (Epic 7-12) correspondent aux nouvelles fonctionnalités de la ROADMAP (v1.1.0+).

**Epics Complétés (v1.0.0):**
- FR1-FR6 → Epic 1: Installation & Scaffolding CLI ✅
- FR7-FR12 → Epic 2: Authentification JWT ✅
- FR20 → Epic 3: User Management & CRUD ✅
- FR13-FR17 → Epic 4: API Infrastructure & Documentation ✅
- FR18-FR19, FR21-FR26 → Epic 5: DevOps & Automation ✅
- Templates multiples → Epic 6: Multiple Project Templates ✅

**Nouveaux Epics (v1.1.0+):** Basés sur ROADMAP - nouvelles fonctionnalités à implémenter.

## Epic List

### Epic 7: Multi-Database Support (v1.1.0)
**Objectif:** Permettre aux utilisateurs de choisir leur base de données préférée (MySQL, SQLite, MongoDB) lors de la génération du projet.

**FRs couverts:** Nouvelles fonctionnalités (ROADMAP)
- Support MySQL/MariaDB avec driver et templates
- Support SQLite pour prototypage rapide
- Support MongoDB (NoSQL)
- Tests E2E pour chaque DB

**Priorité:** Haute | **Estimation:** 3-4 semaines

---

### Epic 8: CRUD Scaffolding Generator (v1.2.0)
**Objectif:** Générer automatiquement du code CRUD pour de nouveaux modèles dans un projet existant via une commande CLI.

**FRs couverts:** Nouvelles fonctionnalités (ROADMAP)
- Commande `create-go-starter add-model <name>` avec parsing interactif/YAML
- Génération complète: model, repository, service, handler, tests
- Mise à jour automatique des routes et Swagger
- Support des relations (one-to-many, many-to-many)

**Priorité:** Haute | **Estimation:** 4-5 semaines

---

### Epic 9: Advanced Observability (v1.3.0)
**Objectif:** Ajouter des outils de monitoring et d'observabilité avancés pour les projets en production.

**FRs couverts:** Nouvelles fonctionnalités (ROADMAP)
- Prometheus metrics endpoint `/metrics`
- Distributed tracing avec OpenTelemetry (Jaeger/Zipkin)
- Health checks avancés (`/health/liveness`, `/health/readiness`)
- Dashboard Grafana pré-configuré avec visualisations

**Priorité:** Moyenne | **Estimation:** 3-4 semaines

---

### Epic 10: CLI Enhancements & Developer Experience
**Objectif:** Améliorer l'expérience développeur avec mode interactif, dry-run et outils de diagnostics.

**FRs couverts:** Améliorations backlog (ROADMAP)
- Mode interactif `--interactive` pour sélection guidée
- Flag `--dry-run` pour preview sans génération
- Commande `doctor` pour diagnostics automatiques
- Progress bar, colored diff et statistiques post-génération

**Priorité:** Moyenne | **Estimation:** 2-3 semaines

---

### Epic 11: Multi-Framework Support (v2.0.0)
**Objectif:** Supporter d'autres frameworks web Go populaires (Gin, Echo) en plus de Fiber.

**FRs couverts:** Nouvelles fonctionnalités (ROADMAP)
- Support Gin framework avec templates adaptés
- Support Echo framework avec templates adaptés
- Abstraction framework-agnostic pour business logic
- Migration guide entre frameworks

**Priorité:** Basse | **Estimation:** 5-6 semaines

---

### Epic 12: Plugin System & Extensibility (v2.x)
**Objectif:** Permettre à la communauté de créer et installer des plugins pour étendre les fonctionnalités.

**FRs couverts:** Vision long-terme (ROADMAP)
- Architecture de plugins modulaire
- Commande `create-go-starter plugin install <name>`
- Registry de plugins communautaires
- Plugins OAuth2, payment providers (Stripe, PayPal), cloud services (AWS, GCP)

**Priorité:** Vision Long-Terme | **Estimation:** 8-10 semaines

---

## Epic 7: Multi-Database Support (v1.1.0)

Permettre aux utilisateurs de choisir leur base de données préférée (MySQL, SQLite, MongoDB) lors de la génération du projet via un flag `--database`.

### Story 7.1: Database Selection Flag

As a développeur,
I want pouvoir spécifier le type de base de données via un flag `--database`,
So that je peux générer un projet avec la base de données de mon choix.

**Acceptance Criteria:**

**Given** l'utilisateur exécute `create-go-starter mon-projet --database=mysql`
**When** le CLI parse les arguments
**Then** MySQL est sélectionné comme base de données
**And** si aucun flag n'est fourni, PostgreSQL est utilisé par défaut

---

### Story 7.2: MySQL/MariaDB Support

As a développeur,
I want générer un projet avec support MySQL/MariaDB,
So that je peux utiliser MySQL comme base de données.

**Acceptance Criteria:**

**Given** l'utilisateur exécute `create-go-starter mon-projet --database=mysql`
**When** le projet est généré
**Then** le driver MySQL est configuré dans go.mod
**And** les templates utilisent la syntaxe MySQL pour les migrations
**And** le docker-compose.yml contient un service MySQL
**And** le projet compile et se connecte à MySQL sans erreur

---

### Story 7.3: SQLite Support

As a développeur,
I want générer un projet avec support SQLite,
So that je peux prototyper rapidement sans serveur de base de données externe.

**Acceptance Criteria:**

**Given** l'utilisateur exécute `create-go-starter mon-projet --database=sqlite`
**When** le projet est généré
**Then** le driver SQLite est configuré
**And** la base de données est un fichier local (.db)
**And** pas de service Docker pour la DB
**And** le projet démarre et fonctionne avec SQLite

---

### Story 7.4: MongoDB Support (Optional)

As a développeur,
I want générer un projet avec support MongoDB,
So that je peux utiliser une base de données NoSQL.

**Acceptance Criteria:**

**Given** l'utilisateur exécute `create-go-starter mon-projet --database=mongodb`
**When** le projet est généré
**Then** le driver mongo-go-driver est configuré
**And** l'architecture est adaptée pour NoSQL (pas de GORM)
**And** le docker-compose.yml contient MongoDB
**And** le projet se connecte à MongoDB

---

### Story 7.5: Database Tests & Documentation

As a développeur,
I want que tous les types de bases de données soient testés et documentés,
So that je puisse choisir en toute confiance.

**Acceptance Criteria:**

**Given** les 4 types de DB sont implémentés
**When** les tests E2E s'exécutent
**Then** tous les tests passent pour chaque type de DB
**And** la documentation explique les différences et cas d'usage
**And** le README est mis à jour avec exemples pour chaque DB

---

## Epic 8: CRUD Scaffolding Generator (v1.2.0)

Générer automatiquement du code CRUD pour de nouveaux modèles dans un projet existant via une commande CLI `add-model`.

### Story 8.1: Add-Model CLI Command

As a développeur,
I want exécuter `create-go-starter add-model <name>` dans un projet existant,
So that le code CRUD pour ce modèle soit généré automatiquement.

**Acceptance Criteria:**

**Given** je suis dans un projet go-starter-kit existant
**When** j'exécute `create-go-starter add-model Todo --fields "title:string,completed:bool"`
**Then** la commande parse les fields
**And** affiche un résumé de ce qui sera généré
**And** demande confirmation avant génération

---

### Story 8.2: Model & Repository Generation

As a développeur,
I want que le modèle et le repository soient générés automatiquement,
So that je n'aie pas à écrire le code boilerplate.

**Acceptance Criteria:**

**Given** la commande `add-model Todo` est validée
**When** la génération s'exécute
**Then** `internal/models/todo.go` est créé avec la struct et tags GORM
**And** `internal/adapters/repository/todo_repository.go` est créé avec CRUD
**And** `internal/interfaces/todo_repository.go` interface est créée
**And** le code compile sans erreur

---

### Story 8.3: Service & Handler Generation

As a développeur,
I want que le service et les handlers soient générés,
So that l'API REST soit complète pour ce modèle.

**Acceptance Criteria:**

**Given** le modèle et repository sont générés
**When** la génération continue
**Then** `internal/domain/todo/service.go` est créé avec logique métier
**And** `internal/adapters/handlers/todo_handler.go` est créé avec endpoints CRUD
**And** les routes sont automatiquement ajoutées dans `routes.go`
**And** tous les endpoints sont accessibles via API

---

### Story 8.4: Tests & Swagger Auto-Update

As a développeur,
I want que les tests et Swagger soient mis à jour automatiquement,
So that le nouveau modèle soit documenté et testé.

**Acceptance Criteria:**

**Given** le code CRUD est généré
**When** la génération se termine
**Then** les tests unitaires sont générés pour service et handler
**And** les annotations Swagger sont ajoutées
**And** la documentation Swagger affiche les nouveaux endpoints
**And** les tests passent

---

### Story 8.5: Relations Support

As a développeur,
I want définir des relations entre modèles (one-to-many, many-to-many),
So that je puisse modéliser des données complexes.

**Acceptance Criteria:**

**Given** j'exécute `add-model Comment --fields "content:string" --belongs-to Todo`
**When** la génération s'exécute
**Then** la relation GORM est correctement configurée
**And** les migrations incluent les foreign keys
**And** les endpoints permettent de récupérer les relations

---

## Epic 9: Advanced Observability (v1.3.0)

Ajouter des outils de monitoring et d'observabilité avancés (Prometheus, OpenTelemetry, Grafana) pour les projets en production.

### Story 9.1: Prometheus Metrics Endpoint

As a DevOps,
I want un endpoint `/metrics` exposant des métriques Prometheus,
So that je puisse monitorer l'application.

**Acceptance Criteria:**

**Given** le projet est généré avec flag `--observability=advanced`
**When** l'application démarre
**Then** l'endpoint `/metrics` est disponible
**And** les métriques HTTP (latence, status codes, throughput) sont exposées
**And** les métriques DB (connections, query time) sont exposées

---

### Story 9.2: Distributed Tracing

As a développeur,
I want activer le distributed tracing avec OpenTelemetry,
So that je puisse tracer les requêtes à travers les services.

**Acceptance Criteria:**

**Given** OpenTelemetry est configuré
**When** une requête traverse l'API
**Then** un trace ID est généré et propagé
**And** les spans sont exportés vers Jaeger/Zipkin
**And** les logs incluent le trace ID pour corrélation

---

### Story 9.3: Advanced Health Checks

As a DevOps,
I want des health checks avancés (`/health/liveness`, `/health/readiness`),
So that Kubernetes puisse gérer l'application correctement.

**Acceptance Criteria:**

**Given** le projet est généré
**When** les endpoints de health sont appelés
**Then** `/health/liveness` vérifie que l'app est vivante
**And** `/health/readiness` vérifie DB et dépendances externes
**And** les réponses sont au format attendu par K8s

---

### Story 9.4: Grafana Dashboard Template

As a DevOps,
I want un dashboard Grafana pré-configuré,
So that je puisse visualiser les métriques immédiatement.

**Acceptance Criteria:**

**Given** le projet est généré avec observability
**When** je déploie le stack (app + Prometheus + Grafana)
**Then** un dashboard JSON est fourni
**And** le dashboard affiche traffic, errors, latency
**And** des alertes sont pré-configurées

---

## Epic 10: CLI Enhancements & Developer Experience

Améliorer l'expérience développeur avec mode interactif, dry-run, diagnostics et feedback visuel.

### Story 10.1: Interactive Mode

As a développeur,
I want lancer `create-go-starter --interactive`,
So that je sois guidé étape par étape dans la configuration.

**Acceptance Criteria:**

**Given** j'exécute `create-go-starter --interactive`
**When** le mode interactif démarre
**Then** le CLI me demande: nom du projet, template, database, observability
**And** chaque option affiche une description claire
**And** je peux naviguer avec les flèches et valider

---

### Story 10.2: Dry-Run Preview

As a développeur,
I want utiliser `--dry-run` pour prévisualiser sans générer,
So that je puisse voir ce qui sera créé.

**Acceptance Criteria:**

**Given** j'exécute `create-go-starter mon-projet --dry-run`
**When** la preview s'affiche
**Then** la liste complète des fichiers à générer est affichée
**And** un diff coloré montre le contenu
**And** aucun fichier n'est réellement créé

---

### Story 10.3: Doctor Command

As a développeur,
I want exécuter `create-go-starter doctor`,
So que le CLI vérifie mon environnement et détecte les problèmes.

**Acceptance Criteria:**

**Given** j'exécute `create-go-starter doctor`
**When** le diagnostic s'exécute
**Then** Go version, Git, Docker sont vérifiés
**And** les problèmes détectés sont affichés avec solutions
**And** un rapport de santé est généré

---

### Story 10.4: Visual Feedback & Progress

As a développeur,
I want voir une progress bar et des statistiques pendant la génération,
So que je sache ce qui se passe.

**Acceptance Criteria:**

**Given** la génération démarre
**When** les fichiers sont créés
**Then** une progress bar affiche l'avancement
**And** les statistiques finales incluent: nombre de fichiers, taille totale, durée
**And** le output est coloré et lisible

---

### Story 10.5: Aliases pour les Options

As a développeur,
I want utiliser des raccourcis pour les options fréquentes,
So que je gagne du temps en ligne de commande.

**Acceptance Criteria:**

**Given** j'utilise le CLI
**When** je tape des commandes
**Then** les aliases courts fonctionnent (-d pour --database, -f pour --framework)
**And** l'aide affiche les aliases disponibles
**And** les aliases sont documentés

---

### Story 10.6: Migration vers Bubble Tea TUI Framework

As a développeur CLI,
I want migrer l'interface utilisateur vers Bubble Tea (https://github.com/charmbracelet/bubbletea),
So that le CLI offre une expérience moderne, professionnelle et intuitive avec des composants interactifs riches.

**Acceptance Criteria:**

**Given** l'utilisateur lance une commande CLI interactive
**When** Bubble Tea est initialisé
**Then** l'interface utilise le framework Bubble Tea avec le pattern Elm Architecture
**And** les composants de base (input, select, confirm) utilisent Bubble Tea
**And** le rendu est fluide et responsive
**And** les raccourcis clavier sont intuitifs (↑↓ navigation, Enter validation, Ctrl+C exit)
**And** l'interface s'adapte à la taille du terminal
**And** le code est structuré avec des Models, Messages et Updates clairs

**Technical Requirements:**
- Bibliothèque: `github.com/charmbracelet/bubbletea` (latest stable)
- Composants Charm: `github.com/charmbracelet/bubbles` pour input, list, spinner, progress
- Lipgloss: `github.com/charmbracelet/lipgloss` pour le styling
- Architecture: Elm Architecture (Model-Update-View)
- Compatibilité: Maintenir la compatibilité avec les flags existants (mode non-interactif)

---

### Story 10.7: Interface Interactive Avancée avec Bubble Tea

As a développeur,
I want une interface interactive avancée avec menus, formulaires multi-étapes et feedback visuel riche,
So that l'expérience CLI soit fluide, moderne et agréable à utiliser.

**Acceptance Criteria:**

**Given** le mode interactif est lancé avec `--interactive`
**When** l'utilisateur navigue dans l'interface
**Then** un menu principal stylé présente les options (template, database, framework, observability)
**And** les formulaires multi-étapes guident l'utilisateur progressivement
**And** une progress bar animée affiche l'avancement de la génération
**And** des spinners indiquent les opérations en cours
**And** les confirmations utilisent des composants visuels clairs (✓/✗)
**And** les erreurs sont affichées avec des messages colorés et des icônes
**And** le diff de preview utilise la coloration syntaxique
**And** l'interface est cohérente avec les couleurs du branding go-starter-kit

**Technical Requirements:**
- Composants avancés: multi-select, paginated lists, viewport scrolling
- Animations: smooth transitions, loading indicators
- Thème personnalisé: palette de couleurs cohérente (vert/bleu/gris)
- Help screens: raccourcis clavier contextuels (? pour aide)
- État persistant: navigation back/forward dans les formulaires
- Tests: tests unitaires pour les Models et Updates
- Documentation: guide d'utilisation de l'interface interactive

---

## Epic 11: Multi-Framework Support (v2.0.0)

Supporter d'autres frameworks web Go (Gin, Echo) en plus de Fiber pour donner plus de choix aux utilisateurs.

### Story 11.1: Framework Selection Flag

As a développeur,
I want spécifier `--framework=gin` ou `--framework=echo`,
So that je puisse choisir mon framework préféré.

**Acceptance Criteria:**

**Given** j'exécute `create-go-starter mon-projet --framework=gin`
**When** le CLI parse les arguments
**Then** Gin est sélectionné comme framework
**And** Fiber reste le défaut si non spécifié

---

### Story 11.2: Gin Framework Templates

As a développeur,
I want générer un projet avec Gin,
So que je puisse utiliser Gin au lieu de Fiber.

**Acceptance Criteria:**

**Given** `--framework=gin` est spécifié
**When** le projet est généré
**Then** les handlers utilisent la syntaxe Gin
**And** le middleware est adapté pour Gin
**And** le projet compile et démarre avec Gin

---

### Story 11.3: Echo Framework Templates

As a développeur,
I want générer un projet avec Echo,
So que je puisse utiliser Echo au lieu de Fiber.

**Acceptance Criteria:**

**Given** `--framework=echo` est spécifié
**When** le projet est généré
**Then** les handlers utilisent Echo
**And** le middleware est adapté
**And** le projet compile et démarre avec Echo

---

### Story 11.4: Framework-Agnostic Abstraction

As a développeur,
I want que la logique métier soit indépendante du framework,
So que je puisse migrer facilement entre frameworks.

**Acceptance Criteria:**

**Given** un projet est généré
**When** j'examine le code
**Then** la couche domain est 100% indépendante du framework
**And** seuls les adapters dépendent du framework
**And** un guide de migration est fourni

---

## Epic 12: Plugin System & Extensibility (v2.x)

Permettre à la communauté de créer et installer des plugins pour étendre les fonctionnalités (OAuth, Stripe, AWS, etc.).

### Story 12.1: Plugin Architecture

As a contributeur,
I want comprendre l'architecture de plugins,
So que je puisse créer mes propres plugins.

**Acceptance Criteria:**

**Given** la documentation plugin est publiée
**When** je la lis
**Then** l'architecture modulaire est expliquée
**And** un plugin exemple est fourni
**And** les hooks disponibles sont documentés

---

### Story 12.2: Plugin Install Command

As a développeur,
I want exécuter `create-go-starter plugin install stripe`,
So que le plugin Stripe soit ajouté à mon projet.

**Acceptance Criteria:**

**Given** un plugin existe dans le registry
**When** j'exécute `plugin install stripe`
**Then** le plugin est téléchargé et installé
**And** les dépendances sont ajoutées à go.mod
**And** les fichiers de configuration sont générés

---

### Story 12.3: Community Plugin Registry

As a contributeur,
I want publier mon plugin dans le registry communautaire,
So que d'autres puissent l'utiliser.

**Acceptance Criteria:**

**Given** j'ai créé un plugin
**When** je soumets au registry
**Then** le plugin est validé et publié
**And** il apparaît dans `plugin list`
**And** les utilisateurs peuvent l'installer

---

### Story 12.4: Core Plugins (OAuth, Payment, Cloud)

As a développeur,
I want des plugins officiels pour OAuth2, Stripe et AWS,
So que je puisse intégrer rapidement ces services.

**Acceptance Criteria:**

**Given** les plugins officiels existent
**When** j'installe un plugin
**Then** OAuth2 (Google, GitHub) est configuré automatiquement
**And** Stripe payment est intégré avec endpoints
**And** AWS S3 est configuré pour file upload