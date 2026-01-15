# Story 5.1: Initialisation Git automatique

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want que le CLI initialise automatiquement un dépôt Git dans le projet généré,
so that je puisse commencer à versionner mon code immédiatement sans étape manuelle.

## Acceptance Criteria

1. Un dépôt Git est initialisé dans le dossier du projet généré (`git init`).
2. Un premier commit initial est créé avec le message "Initial commit from go-starter-kit".
3. Tous les fichiers générés sont inclus dans le commit initial (`git add .`).
4. Si Git n'est pas installé sur le système, un avertissement informatif s'affiche mais la génération du projet continue normalement (pas d'erreur fatale).
5. L'initialisation se produit après la création réussie de tous les fichiers et de la structure.

## Tasks / Subtasks

- [x] Implémenter la détection de Git et l'initialisation (AC: 1, 4)
  - [x] Créer une fonction `initGitRepo(projectPath string)` dans `cmd/create-go-starter/main.go`
  - [x] Utiliser `os/exec` pour vérifier si `git` est disponible
- [x] Réaliser le premier commit (AC: 2, 3)
  - [x] Exécuter `git init` dans le répertoire cible
  - [x] Exécuter `git add .`
  - [x] Exécuter `git commit -m "Initial commit from go-starter-kit"`
- [x] Intégrer dans le workflow de `main.go` (AC: 5)
  - [x] Appeler `initGitRepo` après la génération des fichiers et de l'environnement
  - [x] Ajouter des messages de progression colorés dans le terminal
- [x] Tester le comportement (AC: 1, 2, 3, 4)
  - [x] Vérifier qu'un dossier `.git` est présent dans un nouveau projet
  - [x] Vérifier l'historique git (`git log`)
  - [x] Simuler l'absence de git pour vérifier la gestion d'erreur gracieuse

## Dev Notes

- **Architecture :** Utilisation du package standard `os/exec` pour les commandes shell.
- **Conventions :** Les messages de sortie doivent utiliser les fonctions `Green()` et `Red()` définies dans `main.go`.
- **Handoff :** L'exécution des commandes doit se faire avec `cmd.Dir = projectPath`.
- **Sécurité :** Ne pas inclure de secrets dans le commit initial (le `.gitignore` est déjà généré par le `generator.go`).

### Project Structure Notes

- **Composants à modifier :** 
  - `cmd/create-go-starter/main.go` : Pour ajouter la logique d'orchestration et l'appel aux commandes git.
- **Impact :** Ajoute une dépendance logicielle externe (git) qui doit être gérée de manière optionnelle.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 5.1: Initialisation Git automatique]
- [Source: cmd/create-go-starter/main.go#main] flow de génération actuel.

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash / Claude Sonnet 4

### Debug Log References

- Aucune erreur de debug rencontrée pendant l'implémentation.

### Completion Notes List

- ✅ **AC1**: Fonction `initGitRepo()` créée dans `git.go` - exécute `git init` dans le répertoire projet
- ✅ **AC2**: Commit initial créé avec le message exact "Initial commit from go-starter-kit"
- ✅ **AC3**: Utilisation de `git add .` pour inclure tous les fichiers générés dans le commit initial
- ✅ **AC4**: Fonction `isGitAvailable()` détecte la présence de Git - affiche un avertissement informatif si Git n'est pas installé et continue sans erreur fatale
- ✅ **AC5**: Appel à `initGitRepo()` intégré dans `main.go` après `copyEnvFile()`, garantissant que l'initialisation Git se produit après la création de tous les fichiers
- ✅ Tests unitaires complets pour `isGitAvailable()` et `initGitRepo()`
- ✅ Test E2E `TestE2EGitIntegration` vérifiant le flux complet avec validation des ACs
- ✅ Messages de progression colorés ajoutés dans le terminal (🔧 et ✅)

### File List

- `cmd/create-go-starter/git.go` (nouveau)
- `cmd/create-go-starter/git_test.go` (nouveau)
- `cmd/create-go-starter/main.go` (modifié)

## Change Log

- **2026-01-14**: Implémentation complète de l'initialisation Git automatique avec détection de disponibilité et gestion gracieuse d'erreur (Story 5.1)
