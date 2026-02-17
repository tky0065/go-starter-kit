# Story 1.4: Initialisation du serveur Fiber, DI (fx) et DB

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **développeur**,
I want **que le projet inclue un serveur Fiber et une connexion PostgreSQL gérés par `fx`**,
so that **mon infrastructure soit prête pour les modules métier**.

## Acceptance Criteria

1. **Injection de Dépendances (fx) :** Le cycle de vie de l'application (Démarrage/Arrêt) doit être orchestré par `uber-go/fx`.
2. **Serveur Fiber :** Un serveur Fiber v2 doit être initialisé et démarrer sur le port défini par les variables d'environnement (défaut: 3000).
3. **Connexion Base de Données :** Une connexion GORM vers PostgreSQL doit être établie au démarrage.
4. **Migrations Automatiques :** Le système doit exécuter `AutoMigrate` pour les entités de base au démarrage.
5. **Logs de Démarrage :** Des logs structurés (zerolog) doivent confirmer la connexion réussie à la DB et le démarrage du serveur.
6. **Graceful Shutdown :** L'application doit fermer proprement le serveur Fiber et la connexion DB lors de la réception d'un signal d'arrêt (SIGINT/SIGTERM).

## Tasks / Subtasks

- [ ] Configurer le câblage `fx` dans `internal/infrastructure/server` (AC: 1)
  - [ ] Créer le module `fx` pour Fiber
  - [ ] Implémenter le hook de démarrage (`OnStart`) et d'arrêt (`OnStop`)
- [ ] Initialiser la connexion PostgreSQL avec GORM (AC: 3, 4)
  - [ ] Créer le module `fx` pour la DB dans `internal/infrastructure/database` (ou similaire)
  - [ ] Configurer le pool de connexions
  - [ ] Ajouter l'appel à `AutoMigrate`
- [ ] Intégrer le logger `zerolog` (AC: 5)
  - [ ] Créer un module `fx` pour le logger dans `pkg/logger`
- [ ] Assembler le tout dans le `main.go` généré (AC: 1, 2)
  - [ ] Utiliser `fx.New(...)` pour lancer l'application
- [ ] Implémenter le endpoint de santé `/health` (AC: 2)
  - [ ] Ajouter une route simple retournant `{"status": "ok"}`

## Dev Notes

### Architecture & Constraints
- **Pattern :** Architecture Hexagonale Lite. Les fichiers de setup d'infrastructure vont dans `internal/infrastructure`.
- **DI :** AUCUNE instanciation manuelle dans le `main.go`. Tout doit passer par `fx.Provide`.
- **Database :** Utiliser GORM avec le driver `postgres`.

### Technical Guidelines
- **Versions :**
  - Fiber v2.52.10
  - GORM v1.31.1
  - fx v1.24.0
- **Naming :** Utiliser des constantes pour les clés de configuration (PORT, DB_URL).
- **Graceful Shutdown :** Fiber possède une méthode `Shutdown()` qui doit être appelée dans le hook `OnStop` de fx.

### Project Structure Notes
- Le serveur doit être configuré pour accepter le contexte Go pour le graceful shutdown.
- Les logs doivent être au format JSON en production (configuré via zerolog).

### References
- [Epic 1: Project Initialization & Core Infrastructure](_bmad-output/planning-artifacts/epics.md)
- [Architecture Decision Document](_bmad-output/planning-artifacts/architecture.md)
- [Project Context: Graceful shutdown via fx](_bmad-output/project-context.md)

## Dev Agent Record

### Agent Model Used
Gemini 2.0 Flash

### Debug Log References
None

### Completion Notes List
- Detailed story created for infrastructure initialization.
- Focused on the wiring of fx, Fiber, and GORM.
- Includes requirements for health checks and logging.

### File List
- internal/infrastructure/server/server.go
- internal/infrastructure/database/database.go
- pkg/logger/logger.go
- cmd/main.go (Generated version)
