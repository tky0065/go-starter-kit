---
stepsCompleted: [1, 2, 3, 4]
inputDocuments:
  - '_bmad-output/planning-artifacts/prd.md'
  - '_bmad-output/planning-artifacts/architecture.md'
---

# go-starter-kit - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for go-starter-kit, decomposing the requirements from the PRD, UX Design if it exists, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

- **FR1:** L'utilisateur peut installer l'outil via une commande `go install`.
- **FR2:** L'utilisateur peut générer un nouveau projet en fournissant un nom de projet via le CLI.
- **FR3:** Le système peut créer automatiquement une structure de dossiers respectant l'architecture hexagonale lite.
- **FR4:** Le système peut injecter dynamiquement le nom du projet dans les fichiers générés (fichiers Go, `go.mod`, Dockerfile).
- **FR5:** Le système peut initialiser automatiquement un module Go et télécharger les dépendances nécessaires.
- **FR6:** Le système peut créer un fichier `.env` par défaut à partir d'un template `.env.example`.
- **FR7:** Un visiteur peut créer un compte utilisateur (Register).
- **FR8:** Un utilisateur peut s'authentifier via des identifiants sécurisés (Login).
- **FR9:** Le système peut générer des jetons d'accès (Access Tokens) et de renouvellement (Refresh Tokens) sécurisés.
- **FR10:** Un utilisateur peut renouveler son jeton d'accès sans se reconnecter manuellement.
- **FR11:** Le système peut hacher les mots de passe de manière sécurisée avant stockage.
- **FR12:** Le système peut protéger des routes spécifiques pour n'autoriser que les utilisateurs authentifiés.
- **FR13:** Le système peut regrouper les routes API par domaine fonctionnel.
- **FR14:** Le système peut appliquer un préfixe global `/api/v1` à toutes les routes métier.
- **FR15:** Le système peut gérer automatiquement les erreurs et renvoyer des réponses JSON standardisées.
- **FR16:** Le système peut valider automatiquement les données d'entrée des requêtes HTTP selon des règles prédéfinies.
- **FR17:** Le système peut exposer une documentation interactive (Swagger UI) mise à jour automatiquement.
- **FR18:** Le système peut se connecter à une base de données PostgreSQL de manière résiliente (pool de connexions).
- **FR19:** Le système peut exécuter des migrations de base de données au démarrage ou via commande.
- **FR20:** Un développeur peut effectuer des opérations CRUD (Créer, Lire, Mettre à jour, Supprimer) sur l'entité User.
- **FR21:** Le système peut gérer l'injection de dépendances pour tous les composants majeurs (DB, Serveur, Handlers).
- **FR22:** Le système peut assurer un démarrage et un arrêt "propre" (graceful shutdown) de l'application.
- **FR23:** Le développeur peut lancer l'application en mode développement avec rechargement automatique (hot-reload).
- **FR24:** Le développeur peut exécuter l'ensemble des tests (unitaires et intégration) via une commande unique.
- **FR25:** Le système peut être exécuté dans un environnement conteneurisé (Docker).
- **FR26:** Le système peut être déployé via un pipeline de CI/CD pré-configuré (GitHub Actions).

### NonFunctional Requirements

- **NFR1 (Performance):** Les endpoints de base (Auth/Health) doivent répondre en moins de **100ms** (hors latence réseau).
- **NFR2 (Performance):** L'application conteneurisée doit être opérationnelle en moins de **2 secondes** (Cold Start).
- **NFR3 (Performance):** L'image Docker finale doit être optimisée pour peser moins de **50 Mo**.
- **NFR4 (Security):** Toutes les données sensibles (mots de passe) sont hachées avec `bcrypt` (coût par défaut >= 10).
- **NFR5 (Security):** Utilisation de l'algorithme `HS256` ou `RS256` pour JWT avec une gestion stricte de l'expiration.
- **NFR6 (Security):** Aucune clé ou secret ne doit être en dur dans le code (utilisation obligatoire de `.env`).
- **NFR7 (Security):** Protection native contre les vulnérabilités courantes (CORS, CSRF, Injection SQL via l'ORM).
- **NFR8 (Maintainability):** 100% du code généré doit respecter les standards `golangci-lint`.
- **NFR9 (Documentation):** Chaque fonction publique doit être documentée. Le README doit permettre un démarrage en 5 minutes.
- **NFR10 (Testability):** L'architecture hexagonale doit permettre de mocker 100% des dépendances externes pour les tests unitaires.
- **NFR11 (Operability):** Graceful Shutdown en moins de **5 secondes**.
- **NFR12 (Observability):** Utilisation de logs structurés (JSON) en production.
- **NFR13 (Operability):** Présence d'un endpoint `/health` pour les orchestrateurs.

### Additional Requirements

- **Architecture:** Structure hexagonale "Lite" avec renommage de `/internal/ports` en **`/internal/interfaces`**.
- **Architecture:** Utilisation de `uber-go/fx` pour l'injection de dépendances.
- **Architecture:** API groupée sous `/api/v1`.
- **Architecture:** Auth JWT stockant les Refresh Tokens en base de données PostgreSQL.
- **Architecture:** Gestion centralisée des erreurs via middleware Fiber et format JSON standardisé.
- **Architecture:** Validation déclarative via `go-playground/validator`.
- **Architecture:** Logs structurés via `rs/zerolog`.
- **Stack:** Go 1.25.5, Fiber v2.52.10, GORM v1.31.1, PostgreSQL v16.
- **DevOps:** Docker multi-stage build sur Alpine.

### FR Coverage Map

- **FR1:** Epic 1 - Installation via go install
- **FR2:** Epic 1 - Génération via CLI
- **FR3:** Epic 1 - Structure hexagonale
- **FR4:** Epic 1 - Injection dynamique du nom
- **FR5:** Epic 1 - Init go mod & deps
- **FR6:** Epic 1 - Création .env
- **FR7:** Epic 2 - Register
- **FR8:** Epic 2 - Login
- **FR9:** Epic 2 - Tokens JWT
- **FR10:** Epic 2 - Refresh Token
- **FR11:** Epic 2 - Password Hashing
- **FR12:** Epic 2 - Protected routes
- **FR13:** Epic 4 - Grouping routes
- **FR14:** Epic 4 - Préfixe /api/v1
- **FR15:** Epic 4 - Error management
- **FR16:** Epic 2 - Data validation
- **FR17:** Epic 4 - Swagger UI
- **FR18:** Epic 1 - PostgreSQL connection
- **FR19:** Epic 1 - DB Migrations
- **FR20:** Epic 3 - CRUD User
- **FR21:** Epic 1 - Dependency Injection (fx)
- **FR22:** Epic 1 - Graceful shutdown
- **FR23:** Epic 1 - Hot-reload
- **FR24:** Epic 4 - Test command
- **FR25:** Epic 1 - Dockerization
- **FR26:** Epic 4 - CI/CD GitHub Actions

## Epic List

### Epic 1: Project Initialization & Core Infrastructure
Permettre à l'utilisateur de générer un projet API Go fonctionnel, structuré et connecté à sa base de données en une seule commande.
**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR18, FR19, FR21, FR22, FR23, FR25.
**Status:** ✅ Completed (Sprint 1)

### Epic 2: Authentication & Security Foundation
Fournir un système d'authentification complet et sécurisé dès le départ.
**FRs covered:** FR7, FR8, FR9, FR10, FR11, FR12, FR16.
**Status:** ✅ Completed (Sprint 1)

### Epic 3: User Management Logic
Gérer les données et les opérations métier liées aux utilisateurs.
**FRs covered:** FR20.
**Status:** ✅ Completed (Sprint 1)

### Epic 4: Production Readiness & Developer Experience
Rendre le projet prêt pour la production et agréable à développer.
**FRs covered:** FR13, FR14, FR15, FR17, FR24, FR26.
**Status:** ✅ Completed (Sprint 1)

### Epic 5: MVP Finalization & Quality Assurance
Finaliser le MVP avec les automatisations CLI manquantes et la validation qualité.
**NFRs covered:** NFR3, NFR8, NFR9, NFR10.
**Status:** 🔄 In Progress (Sprint 2)

## Epic 1: Project Initialization & Core Infrastructure

Cette Epic pose les fondations du générateur, de la structure du projet et de l'infrastructure de base.

### Story 1.1: Installation de l'outil CLI

As a développeur,
I want installer le `go-starter-kit` via une commande standard `go install`,
So that je puisse l'utiliser globalement sur ma machine.

**Acceptance Criteria:**

**Given** Go est installé sur le système
**When** J'exécute la commande d'installation (ex: `go install .../cmd/create-go-starter@latest`)
**Then** Le binaire `create-go-starter` est disponible dans le PATH
**And** L'exécution de `create-go-starter --help` affiche les instructions d'utilisation
**And** Les sorties du CLI sont colorées (Vert pour le succès, Rouge pour l'erreur) pour une meilleure lisibilité

### Story 1.2: Génération de la structure de base (Scaffolding)

As a développeur,
I want lancer une commande pour créer l'arborescence hexagonale "Lite",
So that je puisse démarrer sur une base architecturale saine.

**Acceptance Criteria:**

**Given** Je suis dans un répertoire vide
**When** J'exécute `create-go-starter mon-projet`
**Then** Un dossier `mon-projet` est créé
**And** Il contient les répertoires : `cmd`, `internal/adapters`, `internal/domain`, `internal/interfaces`, `internal/infrastructure`, `pkg`, `deployments`
**And** Le CLI affiche une barre de progression ou des étapes claires (ex: "Creating directories...", "Copying templates...")

### Story 1.3: Injection dynamique du contexte projet

As a développeur,
I want que le CLI remplace automatiquement le nom du projet dans les fichiers générés,
So that je n'aie aucun renommage manuel à faire.

**Acceptance Criteria:**

**Given** J'exécute le scaffolding pour un projet nommé "my-api"
**When** Je vérifie les fichiers générés
**Then** Le fichier `go.mod` déclare `module my-api`
**And** Les fichiers source Go utilisent `my-api/...` pour les imports internes
**And** Le `Dockerfile` et le `Makefile` font référence à `my-api`

### Story 1.4: Initialisation du serveur Fiber, DI (fx) et DB

As a développeur,
I want que le projet inclue un serveur Fiber et une connexion PostgreSQL gérés par `fx`,
So that mon infrastructure soit prête pour les modules métier.

**Acceptance Criteria:**

**Given** Le projet est généré et PostgreSQL est disponible
**When** Je lance l'application (`go run cmd/main.go`)
**Then** `uber-go/fx` initialise le serveur Fiber et la connexion GORM/PostgreSQL
**And** Le système exécute automatiquement les migrations de base pour les entités déjà présentes
**And** Le serveur démarre sur le port 3000
**And** Un log structuré confirme la connexion réussie à la base de données

### Story 1.5: Environnement de développement (Dotenv, Makefile & Docker)

As a développeur,
I want disposer d'un fichier `.env`, d'un Makefile et d'un Dockerfile optimisé,
So that je puisse lancer et construire mon projet instantanément.

**Acceptance Criteria:**

**Given** Le projet est généré
**When** Je vérifie la racine du projet
**Then** Un fichier `.env` est créé à partir d'un template `.env.example`
**And** La commande `make dev` lance l'application avec hot-reload
**And** La commande `docker build` produit une image multi-stage légère (Alpine)
**And** Un message de succès final s'affiche en vert avec les prochaines étapes suggérées (ex: "Next steps: cd mon-projet && make dev")

## Epic 2: Authentication & Security Foundation

Cette Epic ajoute la couche de sécurité et la gestion des identités, essentielle pour toute API moderne.

### Story 2.1: Inscription des utilisateurs (Register)

As a visiteur,
I want créer un compte avec mon email et mot de passe,
So that je puisse accéder aux fonctionnalités protégées.

**Acceptance Criteria:**

**Given** Le serveur est démarré et la base de données est accessible
**When** J'envoie une requête POST `/api/v1/auth/register` avec un email valide et un mot de passe fort
**Then** Je reçois une réponse HTTP 201 Created
**And** L'utilisateur est créé en base de données avec son mot de passe haché (bcrypt)
**And** Les données sensibles (mot de passe) ne sont jamais retournées dans la réponse

### Story 2.2: Authentification (Login)

As a utilisateur,
I want me connecter avec mes identifiants,
So that je puisse obtenir des jetons d'accès sécurisés.

**Acceptance Criteria:**

**Given** J'ai un compte utilisateur existant
**When** J'envoie une requête POST `/api/v1/auth/login` avec mes identifiants corrects
**Then** Je reçois une réponse HTTP 200 OK contenant un `access_token` (JWT) et un `refresh_token`
**And** Le `access_token` contient mon ID utilisateur et une expiration courte (ex: 15min)
**And** Le `refresh_token` est stocké ou associé de manière sécurisée côté serveur

### Story 2.3: Renouvellement de session (Refresh Token)

As a utilisateur,
I want obtenir un nouveau jeton d'accès via mon Refresh Token,
So that je puisse rester connecté sans ressaisir mes identifiants.

**Acceptance Criteria:**

**Given** Mon `access_token` a expiré mais j'ai un `refresh_token` valide
**When** J'envoie une requête POST `/api/v1/auth/refresh` avec mon `refresh_token`
**Then** Je reçois une nouvelle paire de tokens valide
**And** L'ancien `refresh_token` est invalidé (rotation de tokens) pour plus de sécurité

### Story 2.4: Sécurisation des routes (Auth Middleware)

As a développeur,
I want protéger certaines routes API via un middleware,
So that seuls les utilisateurs authentifiés puissent y accéder.

**Acceptance Criteria:**

**Given** Une route API configurée comme "protégée" (ex: `/api/v1/users/me`)
**When** J'appelle cette route sans header `Authorization`
**Then** Je reçois une erreur HTTP 401 Unauthorized
**When** J'appelle cette route avec un token JWT valide
**Then** Je reçois la réponse attendue et le contexte utilisateur est injecté dans la requête

## Epic 3: User Management Logic

Cette Epic concrétise la gestion des ressources utilisateurs au niveau métier.

### Story 3.1: Gestion du Profil Utilisateur (Me)

As a utilisateur connecté,
I want consulter mon propre profil,
So that je puisse vérifier mes informations de compte.

**Acceptance Criteria:**

**Given** Je suis authentifié avec un token valide
**When** J'envoie une requête GET `/api/v1/users/me`
**Then** Je reçois une réponse HTTP 200 OK avec mes informations (ID, email, nom)
**And** Les informations retournées correspondent à l'utilisateur identifié par le token

### Story 3.2: Opérations CRUD Utilisateur

As a administrateur,
I want pouvoir lister, modifier ou supprimer des utilisateurs,
So that je puisse gérer la base d'utilisateurs.

**Acceptance Criteria:**

**Given** Une requête autorisée
**When** J'utilise les endpoints correspondants (GET /users, PUT /users/:id, DELETE /users/:id)
**Then** Le système effectue l'opération demandée en base de données
**And** Les réponses respectent le format standard de l'API (success/error)

## Epic 4: Production Readiness & Developer Experience

Cette Epic apporte la touche finale pour transformer le projet en un outil professionnel prêt pour le déploiement.

### Story 4.1: Standardisation des API (Grouping & V1)

As a développeur,
I want que toutes les routes soient préfixées par `/api/v1`,
So that je puisse versionner mon API facilement à l'avenir.

**Acceptance Criteria:**

**Given** Le serveur est démarré
**When** J'accède à une route métier (ex: Auth ou Users)
**Then** Elle doit obligatoirement être préfixée par `/api/v1`
**And** Les routes sont regroupées logiquement dans le code via les "Groups" de Fiber

### Story 4.2: Gestion centralisée des erreurs

As a développeur,
I want un mécanisme uniforme pour formater les erreurs en JSON,
So that les clients de mon API reçoivent des réponses cohérentes en cas de problème.

**Acceptance Criteria:**

**Given** Une erreur survient (ex: 404, 500, erreur de validation)
**When** Le système renvoie la réponse au client
**Then** Le corps de la réponse est un JSON structuré (ex: `{"status": "error", "message": "...", "code": "..."}`)
**And** Aucune information sensible (stack trace) n'est exposée en production

### Story 4.3: Documentation interactive (Swagger)

As a consommateur de l'API,
I want accéder à une documentation Swagger auto-générée,
So that je puisse comprendre et tester l'API sans lire le code source.

**Acceptance Criteria:**

**Given** Le serveur est démarré
**When** J'accède à l'URL `/swagger` dans mon navigateur
**Then** L'interface Swagger UI s'affiche avec la liste de tous les endpoints
**And** Je peux tester une requête (ex: Login) directement depuis l'interface

### Story 4.4: Automatisation de la Qualité (Lint & Test)

As a développeur,
I want lancer les tests et le linter via une seule commande,
So that je puisse garantir la qualité de mon code rapidement.

**Acceptance Criteria:**

**Given** Je suis à la racine du projet
**When** J'exécute `make test`
**Then** Tous les tests unitaires et d'intégration sont lancés
**When** J'exécute `make lint`
**Then** Le linter `golangci-lint` vérifie la conformité du code selon les standards définis

### Story 4.5: Intégration Continue (GitHub Actions)

As a développeur,
I want que mes tests soient exécutés automatiquement sur GitHub,
So that je puisse éviter les régressions lors du travail en équipe.

**Acceptance Criteria:**

**Given** Un dépôt GitHub configuré
**When** Je pousse du code sur la branche principale
**Then** Un workflow GitHub Actions se déclenche automatiquement
**And** Il installe Go, exécute le linting et lance les tests
**And** Le statut du commit devient rouge en cas d'échec

## Epic 5: MVP Finalization & Quality Assurance

Cette Epic finalise le MVP en ajoutant les automatisations manquantes du CLI, en améliorant la couverture de tests, et en vérifiant la conformité aux NFR (Non-Functional Requirements) définis dans le PRD.

**Objectif:** Atteindre un MVP "production-ready" avec une expérience utilisateur fluide et une qualité de code irréprochable.

**NFRs couverts:** NFR3 (Image Docker < 50 Mo), NFR8 (golangci-lint compliant), NFR9 (Documentation), NFR10 (Testabilité)

### Story 5.1: Initialisation Git automatique

As a développeur,
I want que le CLI initialise automatiquement un dépôt Git dans le projet généré,
So that je puisse commencer à versionner mon code immédiatement sans étape manuelle.

**Acceptance Criteria:**

**Given** J'exécute `create-go-starter mon-projet`
**When** La génération est terminée avec succès
**Then** Un dépôt Git est initialisé dans le dossier `mon-projet` (`git init`)
**And** Un fichier `.gitignore` approprié est présent (déjà fait)
**And** Un premier commit initial est créé avec le message "Initial commit from go-starter-kit"
**And** Si Git n'est pas installé, un avertissement s'affiche mais la génération continue

**Technical Notes:**
- Utiliser `os/exec` pour exécuter les commandes git
- Gérer gracieusement l'absence de Git sur le système
- Le commit initial doit inclure tous les fichiers générés

### Story 5.2: Installation automatique des dépendances Go

As a développeur,
I want que le CLI exécute automatiquement `go mod tidy` après la génération,
So that toutes les dépendances soient téléchargées et le projet soit immédiatement fonctionnel.

**Acceptance Criteria:**

**Given** Le projet a été généré avec succès
**When** Le CLI termine la génération des fichiers
**Then** La commande `go mod tidy` est exécutée automatiquement
**And** Les dépendances sont téléchargées dans le cache Go
**And** Un message de progression s'affiche ("📦 Installation des dépendances...")
**And** En cas d'échec réseau, un message d'erreur clair s'affiche avec les instructions pour réessayer manuellement

**Technical Notes:**
- Utiliser `os/exec` pour exécuter `go mod tidy`
- Capturer stdout/stderr pour afficher les erreurs éventuelles
- Timeout de 2 minutes maximum pour l'opération

### Story 5.3: Amélioration de la couverture de tests du CLI

As a mainteneur du projet,
I want que la couverture de tests du CLI atteigne au moins 70%,
So that le code soit fiable et les régressions détectées automatiquement.

**Acceptance Criteria:**

**Given** Le code actuel du CLI a une couverture de ~44%
**When** J'exécute `go test -cover ./cmd/create-go-starter`
**Then** La couverture affichée est >= 70%
**And** Les fonctions critiques sont testées : `validateProjectName`, `createProjectStructure`, `generateProjectFiles`, `copyEnvFile`
**And** Les cas d'erreur sont couverts (nom invalide, répertoire existant, erreurs d'écriture)

**Technical Notes:**
- Ajouter des tests pour les nouvelles fonctionnalités (git init, go mod tidy)
- Utiliser `t.TempDir()` pour les tests de système de fichiers
- Mocker les commandes externes si nécessaire

### Story 5.4: Optimisation de l'image Docker générée

As a DevOps/SRE,
I want que l'image Docker générée par le starter kit pèse moins de 50 Mo,
So that les déploiements soient rapides et les coûts de stockage minimisés.

**Acceptance Criteria:**

**Given** Un projet généré avec `create-go-starter`
**When** J'exécute `docker build -t mon-projet .`
**Then** L'image finale pèse moins de 50 Mo
**And** L'image utilise un utilisateur non-root pour la sécurité
**And** Seuls les fichiers nécessaires sont inclus (pas de sources, pas de cache)
**And** Les health checks Docker sont configurés

**Technical Notes:**
- Vérifier le Dockerfile template actuel
- Utiliser `scratch` ou `alpine` comme base
- Ajouter `USER nonroot` dans le Dockerfile
- Ajouter HEALTHCHECK instruction

### Story 5.5: Documentation des fonctions publiques du code généré

As a développeur utilisant le starter kit,
I want que toutes les fonctions publiques du code généré soient documentées,
So that je puisse comprendre rapidement le fonctionnement de chaque composant.

**Acceptance Criteria:**

**Given** Le code généré par le CLI
**When** J'examine les fichiers Go générés
**Then** Chaque fonction publique (commençant par une majuscule) a un commentaire de documentation
**And** Les commentaires suivent le format standard Go (`// FunctionName does...`)
**And** Les structures et interfaces publiques sont également documentées
**And** `go doc` affiche une documentation lisible pour chaque package

**Technical Notes:**
- Mettre à jour tous les templates dans `templates.go` et `templates_user.go`
- Vérifier avec `golint` ou `go doc`
- Priorité aux fichiers : handlers, services, repositories, middleware

### Story 5.6: Validation finale et smoke tests

As a mainteneur du projet,
I want exécuter une suite de tests de validation complète sur un projet généré,
So that je puisse confirmer que le MVP fonctionne de bout en bout.

**Acceptance Criteria:**

**Given** Le CLI avec toutes les améliorations du Sprint 2
**When** Je génère un nouveau projet de test
**Then** La génération se termine sans erreur
**And** `go build ./...` compile sans erreur
**And** `go test ./...` passe (si des tests sont inclus)
**And** `golangci-lint run` ne retourne aucune erreur
**And** Le serveur démarre et répond sur `/health`
**And** L'endpoint `/swagger` est accessible
**And** L'authentification JWT fonctionne (register/login)

**Technical Notes:**
- Créer un script de validation E2E
- Peut être un test Go avec le tag `// +build e2e`
- Documenter la procédure de validation manuelle

