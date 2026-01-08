# Story 1.5: Environnement de développement (Dotenv, Makefile & Docker)

Status: ready-for-dev

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

- [ ] Créer les templates pour les fichiers de configuration (AC: 1)
  - [ ] Implémenter le template `.env.example`
  - [ ] Ajouter la logique de copie `.env.example` -> `.env` dans le CLI
- [ ] Créer le template Makefile (AC: 2)
  - [ ] Définir les cibles `dev`, `build`, `test`, `clean`
  - [ ] S'assurer que `make dev` utilise un outil de hot-reload (ex: `air`)
- [ ] Implémenter la conteneurisation (AC: 3)
  - [ ] Créer le `Dockerfile` multi-stage
  - [ ] Créer le `docker-compose.yml` avec les services `api` et `db`
- [ ] Ajouter les instructions de succès dans le CLI (AC: 4)
  - [ ] Formater le message de sortie avec les commandes `cd` et `make dev`

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
Gemini 2.0 Flash

### Debug Log References
None

### Completion Notes List
- Story for development environment automation and containerization created.
- Integrated best practices for Docker (non-root, multi-stage, light images).
- Aligned with previous stories for CLI feedback.

### File List
- .env.example (Template)
- Makefile (Template)
- deployments/Dockerfile (Template)
- deployments/docker-compose.yml (Template)
- cmd/create-go-starter/main.go (Modification pour le message de succès)
