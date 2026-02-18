# Story 10.1: Interactive Mode

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want lancer `create-go-starter --interactive`,
so that je sois guidé étape par étape dans la configuration du projet via un menu interactif.

## Acceptance Criteria

1. **Given** j'exécute `create-go-starter --interactive`, **When** le mode interactif démarre, **Then** le CLI me pose des questions séquentielles: nom du projet, template (minimal/full/graphql), database (postgres/mysql/sqlite), observability (none/basic/advanced).
2. **Given** le mode interactif est actif, **When** une option est présentée, **Then** chaque option affiche une description claire de son utilité.
3. **Given** le mode interactif est actif, **When** je réponds aux questions, **Then** les valeurs saisies sont validées (validateTemplate, validateDatabase, validateObservabilityLevel) avec un message d'erreur clair si invalide, et redemande la valeur.
4. **Given** toutes les réponses sont fournies, **When** je confirme, **Then** `run(projectName, template, database, observability)` est appelé avec les valeurs collectées.
5. **Given** le flag `--interactive` ou `-interactive` est passé seul, **When** le CLI démarre, **Then** le mode interactif s'active sans nécessiter de nom de projet en argument.
6. **Given** le mode interactif est actif, **When** l'utilisateur appuie sur Ctrl+C, **Then** le CLI termine proprement avec le message "Cancelled" sans erreur fatale.

## Tasks / Subtasks

- [x] Task 1 : Ajouter le flag `--interactive` dans le parsing d'arguments de `main()` (AC: #1, #5)
  - [x] Détecter `-interactive` et `--interactive` dans la boucle `for i := 0; i < len(args); i++`
  - [x] Mettre à jour la fonction `usage()` pour documenter le flag
  - [x] Mettre à jour les exemples dans `usage()`
- [x] Task 2 : Implémenter `runInteractiveMode() error` dans un nouveau fichier `interactive.go` (AC: #1, #2, #3, #4, #6)
  - [x] Créer `cmd/create-go-starter/interactive.go`
  - [x] Implémenter la lecture de l'input via `bufio.NewReader(os.Stdin)` (stdlib only, no deps)
  - [x] Implémenter `promptStringFromReader(r *bufio.Reader, prompt, defaultVal string) (string, error)` helper
  - [x] Implémenter `promptSelectFromReader(r *bufio.Reader, prompt string, options []string, descriptions []string, defaultIdx int) (string, error)` helper
  - [x] Collecter: projectName, template, database, observability dans l'ordre
  - [x] Afficher un résumé avant confirmation finale avec `[y/N]`
  - [x] Appeler `run(projectName, template, database, observability)` avec les valeurs collectées
- [x] Task 3 : Intégrer dans `main()` (AC: #4, #5)
  - [x] Si `--interactive` détecté, appeler `runInteractiveMode()` et retourner
  - [x] S'assurer que `--interactive` peut coexister avec d'autres flags (prendre les valeurs comme defaults)
- [x] Task 4 : Tests (AC: #1-#6)
  - [x] Créer `cmd/create-go-starter/interactive_test.go`
  - [x] Tester `promptStringFromReader` avec mock de stdin via `strings.NewReader`
  - [x] Tester `promptSelectFromReader` avec différentes sélections
  - [x] Tester le cas Ctrl+C / EOF pour termination propre
  - [x] Vérifier que les validations sont appelées avec les mauvaises valeurs

## Dev Notes

### Architecture et Patterns

**Fichier cible principal :** `cmd/create-go-starter/main.go` (parsing) + nouveau `cmd/create-go-starter/interactive.go`

**Pattern de parsing actuel (main.go:214-251) :** La boucle de parsing est manuelle. Ajouter la détection de `--interactive` comme les autres flags :
```go
} else if arg == "-interactive" || arg == "--interactive" {
    interactive = true
}
```

**Pattern de sous-commande :** Les sous-commandes comme `add-model` sont détectées **avant** le parsing des flags (main.go:195-201). Le flag `--interactive` est différent — c'est un flag, pas une sous-commande, donc il va dans la boucle de parsing.

**Design de `runInteractiveMode()`** — Utiliser UNIQUEMENT la stdlib (pas de dépendances externes) :
```go
// cmd/create-go-starter/interactive.go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func runInteractiveMode() error {
    reader := bufio.NewReader(os.Stdin)
    // ... collect projectName, template, database, observability
    // ... validate each input
    // ... display summary + confirmation
    // ... call run(...)
}
```

**Pas de dépendances externes :** Le `go.mod` n'a aucune dépendance externe. Implémenter avec `bufio.NewReader` pour la lecture, `fmt.Print` pour les prompts. Ne PAS ajouter `survey`, `bubbletea`, ou d'autres libs.

**Validation existante :** Réutiliser `validateTemplate()`, `validateDatabase()`, `validateObservabilityLevel()` (main.go:95-133) pour valider chaque input.

**Constants réutilisables** (main.go:28-76) :
- `ValidTemplates = []string{TemplateMinimal, TemplateFull, TemplateGraphQL}`
- `ValidDatabases = []string{DatabasePostgres, DatabaseMySQL, DatabaseSQLite}`
- `ValidObservabilityLevels = []string{ObservabilityNone, ObservabilityBasic, ObservabilityAdvanced}`
- `DefaultTemplate = "full"`, `DefaultDatabase = "postgres"`, `DefaultObservabilityLevel = "none"`
- `TemplateMinimalDesc`, `TemplateFullDesc`, `TemplateGraphQLDesc`

**Descriptions à afficher dans le prompt select :**
- templates : utiliser les constantes `TemplateMinimalDesc`, `TemplateFullDesc`, `TemplateGraphQLDesc`
- databases : "PostgreSQL - Production-ready, advanced features", "MySQL/MariaDB - Wide compatibility, shared hosting", "SQLite - Quick prototyping, embedded apps"
- observability : "None - No observability (default)", "Basic - Enhanced /health endpoint only", "Advanced - Full Prometheus metrics + middleware"

**Colors :** `Green()`, `Red()`, `Yellow()` sont disponibles dans main.go:78-91. Les utiliser pour l'affichage.

**EOF/Ctrl+C handling :** `bufio.Reader.ReadString('\n')` retourne `io.EOF` ou une erreur sur Ctrl+C. Retourner une erreur propre : `fmt.Errorf("interrupted")` et afficher "Cancelled." avant de retourner.

### Structure du prompt interactif (UX)

```
══════════════════════════════════════════
  create-go-starter — Interactive Mode
══════════════════════════════════════════

? Project name: [my-project]
  → Enter your project name (used as directory and Go module name)

? Template: (use number to select)
  1) full     - Complete API with JWT auth, user management, and Swagger (default)
  2) minimal  - Basic REST API with Swagger (no authentication)
  3) graphql  - GraphQL API with gqlgen and GraphQL Playground
  Select [1]:

? Database: (use number to select)
  1) postgres - PostgreSQL (default) - Production-ready, advanced features
  2) mysql    - MySQL/MariaDB - Wide compatibility, shared hosting
  3) sqlite   - SQLite - Quick prototyping, embedded apps
  Select [1]:

? Observability: (use number to select)
  1) none     - No observability (default)
  2) basic    - Enhanced /health endpoint only (no Prometheus)
  3) advanced - Full Prometheus metrics endpoint + HTTP metrics middleware
  Select [1]:

══════════════════════════════════════════
Summary:
  Project:       my-project
  Template:      full
  Database:      postgres
  Observability: none
══════════════════════════════════════════

? Confirm? [Y/n]:
```

### Project Structure Notes

- **Nouveau fichier :** `cmd/create-go-starter/interactive.go` (co-localisé avec main.go)
- **Nouveau fichier test :** `cmd/create-go-starter/interactive_test.go`
- **Fichier modifié :** `cmd/create-go-starter/main.go` (ajout parsing `--interactive`, appel `runInteractiveMode()`)
- **Pas de nouveaux répertoires** requis

### Testing Approach

Tester via **redirection stdin** en utilisant `strings.NewReader` :
```go
// Dans les tests, utiliser des pipes ou strings.Reader pour simuler stdin
func TestPromptString(t *testing.T) {
    // Pas de test directement sur os.Stdin — utiliser des fonctions avec io.Reader paramètre
}
```

**Recommandation :** Décomposer `runInteractiveMode` en helpers acceptant un `io.Reader` pour la testabilité :
```go
func runInteractiveModeWithReader(r io.Reader) error
func promptStringFromReader(r *bufio.Reader, prompt, defaultVal string) (string, error)
```

### References

- Parsing flags actuel : [Source: cmd/create-go-starter/main.go:213-251]
- Constants templates/databases/observability : [Source: cmd/create-go-starter/main.go:28-76]
- Fonctions de validation : [Source: cmd/create-go-starter/main.go:94-133]
- Fonctions couleur : [Source: cmd/create-go-starter/main.go:78-91]
- Fonction `run()` principale : [Source: cmd/create-go-starter/main.go:343-395]
- Pattern sous-commande `add-model` : [Source: cmd/create-go-starter/main.go:195-201]
- Architecture hexagonale CLI : [Source: _bmad-output/planning-artifacts/architecture.md]
- Epic 10 description : [Source: _bmad-output/planning-artifacts/epics.md#Epic-10]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

Aucun debug nécessaire — implémentation directe selon les spécifications de la story.

### Completion Notes List

- ✅ **Task 1** : Flag `--interactive` / `-interactive` ajouté dans le parsing manuel de `main.go`. La fonction `usage()` et les exemples ont été mis à jour.
- ✅ **Task 2** : `interactive.go` créé avec `runInteractiveMode()` (wrapper os.Stdin) et `runInteractiveModeWithReader(io.Reader)` (testable). Helpers `promptStringFromReader`, `promptSelectFromReader`, `promptConfirmFromReader`, `collectProjectName` implémentés via stdlib uniquement (bufio, fmt, strings, strconv) — aucune dépendance externe ajoutée.
- ✅ **Task 3** : Dans `main()`, si `interactive == true`, appel à `runInteractiveMode()` puis `return`. Le mode interactif s'active sans projet requis en argument.
- ✅ **Task 4** : `interactive_test.go` créé avec 8 tests couvrant: promptString (4 cas), promptSelect (8 cas), EOF/Ctrl+C (2 tests), cancellation (1 test), validation re-prompt (1 test), succès complet (1 test integration). Le test d'intégration (success) est skippé en `-short` pour éviter les interférences via `os.Chdir`.
- **Décision architecturale** : Utilisation de `io.Reader` comme paramètre plutôt que `os.Stdin` directement — pattern recommandé par la story pour la testabilité.
- **Point notable** : `ValidateGoModuleName` accepte les noms commençant par un chiffre (regex `^[a-zA-Z0-9]...`). Les tests de validation testent plutôt les sélections invalides dans les menus.

### File List

- `cmd/create-go-starter/interactive.go` (nouveau)
- `cmd/create-go-starter/interactive_test.go` (nouveau)
- `cmd/create-go-starter/main.go` (modifié : flag `--interactive`, `usage()`, intégration dans `main()`)
- `_bmad-output/implementation-artifacts/sprint-status.yaml` (modifié : statut epic-10 et stories)

## Change Log

- 2026-02-17 : Story 10.1 implémentée — Mode interactif CLI avec `--interactive` flag, prompts guidés (projet, template, database, observability), validation, résumé et confirmation. Tests unitaires et d'intégration ajoutés.
- 2026-02-17 : Code Review (AI) — 8 issues trouvées (2 HIGH, 4 MEDIUM, 2 LOW). Corrections appliquées :
  - [H1] `InteractiveDefaults` struct ajoutée pour passer les flags CLI comme defaults au mode interactif
  - [H2] Test `TestRunInteractiveModeInvalidProjectName` corrigé (utilise `@invalid!` au lieu de `123invalid`)
  - [M1] Binaire `create-go-starter` retiré du tracking git (`git rm --cached`)
  - [M2] File List complétée avec `sprint-status.yaml`
  - [M3] Validation `observability=advanced + template!=full` ajoutée dans le mode interactif
  - [M4] Format du prompt de confirmation corrigé (`promptConfirmFromReader` avec `defaultYes` param)
  - Tests ajoutés : `TestRunInteractiveModeWithDefaults`, `TestRunInteractiveModeObservabilityValidation`, `TestPromptConfirmFromReader` (6 cas)
  - Fix cosmétique : suppression du double `(default)` dans l'affichage des menus
