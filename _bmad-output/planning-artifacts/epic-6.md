# Epic 6: Templates Multiples

## Objectif

Permettre aux utilisateurs de choisir entre différents types de projets Go lors de la génération:
- **minimal**: API basique sans authentification, avec Swagger
- **full**: Structure complète actuelle (JWT, CRUD users, Swagger)
- **graphql**: API GraphQL avec gqlgen et GraphQL Playground

## Stories

### Story 6.1: Flag CLI pour sélection de template

**En tant que** développeur,
**Je veux** pouvoir spécifier le type de template via un flag CLI,
**Afin de** générer le type de projet adapté à mes besoins.

**Acceptance Criteria:**

- **Given** l'utilisateur exécute `create-go-starter mon-projet --template=minimal`
- **When** le CLI parse les arguments
- **Then** le template minimal est sélectionné
- **And** si aucun flag n'est fourni, le template "full" est utilisé par défaut

**Tâches techniques:**
1. Ajouter flag `--template` avec valeurs: minimal, full, graphql
2. Valider la valeur du flag (erreur si invalide)
3. Passer le type de template au générateur
4. Afficher le template sélectionné dans la sortie

---

### Story 6.2: Template Minimal

**En tant que** développeur,
**Je veux** générer un projet Go minimal avec API REST et Swagger,
**Afin de** démarrer rapidement sans la complexité de l'authentification.

**Acceptance Criteria:**

- **Given** l'utilisateur exécute `create-go-starter mon-projet --template=minimal`
- **When** le projet est généré
- **Then** la structure inclut: Fiber, GORM, Swagger, Health check, Logger
- **And** PAS d'authentification JWT ni de gestion utilisateurs
- **And** le projet compile et démarre sans erreur

**Fichiers générés (minimal):**
- cmd/main.go
- internal/infrastructure/database/database.go
- internal/infrastructure/server/server.go
- internal/adapters/http/health.go
- internal/adapters/http/routes.go
- pkg/config/env.go
- pkg/logger/logger.go
- .env.example, Dockerfile, Makefile, go.mod, README.md
- Documentation Swagger

---

### Story 6.3: Refactoring du template Full

**En tant que** développeur,
**Je veux** que le template "full" soit la référence actuelle,
**Afin de** maintenir la compatibilité avec les projets existants.

**Acceptance Criteria:**

- **Given** l'utilisateur exécute `create-go-starter mon-projet` (sans flag)
- **When** le projet est généré
- **Then** le comportement est identique à la version actuelle
- **And** le template "full" est sélectionné par défaut

**Tâches techniques:**
1. Refactorer les templates existants en "full"
2. Organiser le code pour supporter plusieurs templates
3. S'assurer que les tests existants passent toujours

---

### Story 6.4: Template GraphQL

**En tant que** développeur,
**Je veux** générer un projet Go avec API GraphQL,
**Afin de** créer des APIs GraphQL modernes.

**Acceptance Criteria:**

- **Given** l'utilisateur exécute `create-go-starter mon-projet --template=graphql`
- **When** le projet est généré
- **Then** la structure inclut: gqlgen, GraphQL Playground, GORM, Logger
- **And** un schéma GraphQL de base avec type User et mutations
- **And** le projet compile et le playground est accessible

**Fichiers générés (graphql):**
- cmd/main.go
- graph/schema.graphqls
- graph/resolver.go
- graph/model/models_gen.go
- internal/infrastructure/database/database.go
- internal/infrastructure/server/server.go (avec GraphQL handler)
- pkg/config/env.go
- pkg/logger/logger.go
- gqlgen.yml
- .env.example, Dockerfile, Makefile, go.mod, README.md

---

### Story 6.5: Documentation et aide CLI

**En tant que** développeur,
**Je veux** voir les templates disponibles via l'aide CLI,
**Afin de** comprendre les options de génération.

**Acceptance Criteria:**

- **Given** l'utilisateur exécute `create-go-starter --help`
- **When** l'aide s'affiche
- **Then** la liste des templates disponibles est affichée avec description
- **And** le template par défaut est indiqué

**Exemple de sortie:**
```
Usage: create-go-starter [options] <project-name>

Options:
  --template string   Template type (default "full")
                      - minimal: Basic REST API with Swagger (no auth)
                      - full: Complete API with JWT auth, users, Swagger
                      - graphql: GraphQL API with gqlgen and Playground
  --help              Show this help message
```

---

## Dépendances

- Epic 5 (terminé): Base stable du CLI

## Risques

1. **Complexité du code**: Gérer plusieurs templates peut compliquer le code
   - Mitigation: Utiliser une architecture modulaire avec interfaces

2. **Maintenance**: Chaque template doit être maintenu séparément
   - Mitigation: Partager le maximum de code commun entre templates

## Critères de succès

- [x] Les 3 templates génèrent des projets fonctionnels
- [x] Les tests couvrent tous les templates
- [x] La documentation est mise à jour
- [x] Le flag --help affiche correctement les options
