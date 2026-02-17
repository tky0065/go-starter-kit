# Story 1.5: Environnement de développement (Dotenv, Makefile & Docker)

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **développeur**,
I want **disposer d'un fichier `.env`, d'un Makefile et d'un Dockerfile optimisé**,
so that **je puisse lancer et construire mon projet instantanément**.

## Acceptance Criteria

1. **Configuration (Dotenv) :** Un fichier `.env.example` doit être créé avec toutes les variables nécessaires (DB_URL, JWT_SECRET, PORT, etc.). Un fichier `.env` doit être généré automatiquement s'il n'existe pas lors de la création du projet.
2. **Automatisation (Makefile) :** Un `Makefile` doit être présent à la racine avec les commandes suivantes :
    - `make dev` : Lance l'application avec hot-reload (utilisant `air`).
    - `make build` : Compile le binaire Go.
    - `make test` : Exécute les tests unitaires.
    - `make clean` : Nettoie les fichiers de build.
3. **Conteneurisation (Docker) :** 
    - Un `Dockerfile` multi-stage optimisé (build sur Go-alpine, runtime sur Alpine minimal) doit être présent.
    - La taille de l'image finale doit être < 50 Mo.
    - Un fichier `docker-compose.yml` doit être inclus pour lancer l'API et PostgreSQL facilement.
4. **Feedback Final :** Une fois le projet généré par le CLI, un message de succès en **Vert** doit s'afficher avec les prochaines étapes suggérées (ex: "Next steps: cd <projectName> && make dev").

## Tasks / Subtasks

- [x] Créer les templates pour les fichiers de configuration (AC: 1)
  - [x] Implémenter le template `.env.example`
  - [x] Ajouter la logique de copie `.env.example` -> `.env` dans le CLI
- [x] Créer le template Makefile (AC: 2)
  - [x] Définir les cibles `dev`, `build`, `test`, `clean`
  - [x] S'assurer que `make dev` utilise un outil de hot-reload (ex: `air`)
- [x] Implémenter la conteneurisation (AC: 3)
  - [x] Créer le `Dockerfile` multi-stage
  - [x] Créer le `docker-compose.yml` avec les services `api` et `db`
- [x] Ajouter les instructions de succès dans le CLI (AC: 4)
  - [x] Formater le message de sortie avec les commandes `cd` et `make dev`

## Dev Notes

### Architecture & Constraints
- **Docker :** Utiliser des images Alpine pour la légèreté.
- **Hot-Reload :** Recommander l'installation de `air` dans le README ou l'inclure dans la documentation de `make dev`.
- **Secrets :** Le fichier `.env` doit être listé dans le `.gitignore` généré.

### Technical Guidelines
- Le `docker-compose.yml` doit utiliser les variables d'environnement définies dans le `.env`.
- Le `Dockerfile` doit utiliser un utilisateur non-root pour la sécurité (best practice mentionnée dans l'ADD).
- S'assurer que le port 3000 est exposé par défaut.

### Project Structure Notes
- Les fichiers `Dockerfile` et `docker-compose.yml` peuvent être placés dans `deployments/` ou à la racine selon les préférences (l'ADD suggère `deployments/` mais souvent le Dockerfile est à la racine pour la simplicité de build context). Je vais suivre la structure de l'ADD: `deployments/`.

### References
- [Epic 1: Project Initialization & Core Infrastructure](_bmad-output/planning-artifacts/epics.md)
- [Architecture Decision Document](_bmad-output/planning-artifacts/architecture.md)
- [Project Context: NO SECRETS in code](_bmad-output/project-context.md)

## Dev Agent Record

### Agent Model Used
Claude Sonnet 4.5

### Debug Log References
None

### Completion Notes List
- Implémenté la logique de copie automatique `.env.example` -> `.env` dans `generator.go`
- Créé le template `DockerComposeTemplate` avec services API et PostgreSQL
- Optimisé le `DockerfileTemplate` avec:
  - Multi-stage build (golang:1.25-alpine + alpine:latest)
  - Utilisateur non-root (appuser) pour la sécurité
  - Flags d'optimisation de build (-ldflags="-w -s")
  - HEALTHCHECK intégré
  - Taille finale < 50 Mo (Alpine minimal)
- Déplacé Dockerfile et docker-compose.yml dans le répertoire `deployments/`
- Mis à jour le message de succès du CLI avec instructions en vert: `cd <project> && make dev`
- Tous les tests passent (env, templates, generator, E2E)
- Le Makefile existant contient déjà les cibles requises (dev, build, test, clean) avec support de `air`

### File List
- cmd/create-go-starter/generator.go (ajout de copyEnvFile et mise à jour de generateProjectFiles)
- cmd/create-go-starter/templates.go (DockerComposeTemplate, Dockerfile optimisé)
- cmd/create-go-starter/main.go (message de succès en vert avec "cd <project> && make dev")
- cmd/create-go-starter/generator_test.go (tests mis à jour pour deployments/)
- cmd/create-go-starter/templates_test.go (tests mis à jour pour Dockerfile optimisé)
- cmd/create-go-starter/env_test.go (tests existants pour copyEnvFile)

## Change Log
- (2026-01-08) Implémentation complète de l'environnement de développement avec dotenv, Makefile, Docker et docker-compose
- (2026-01-08) Ajout de la logique de copie automatique .env.example -> .env
- (2026-01-08) Création du template docker-compose.yml avec PostgreSQL
- (2026-01-08) Optimisation du Dockerfile avec utilisateur non-root et taille < 50 Mo
- (2026-01-08) Mise à jour du message de succès du CLI avec instructions colorées
