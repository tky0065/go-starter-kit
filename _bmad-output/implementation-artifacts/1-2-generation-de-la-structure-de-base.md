# Story 1.2: Génération de la structure de base (Scaffolding)

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **développeur**,
I want **lancer une commande pour créer l'arborescence hexagonale "Lite"**,
so that **je puisse démarrer sur une base architecturale saine**.

## Acceptance Criteria

1. **Dossier de projet :** Un dossier `mon-projet` (nom fourni en argument) doit être créé dans le répertoire courant.
2. **Structure Hexagonale Lite :** Le projet généré doit contenir exactement les répertoires suivants :
    - `cmd/`
    - `internal/adapters/`
    - `internal/domain/`
    - `internal/interfaces/`
    - `internal/infrastructure/`
    - `pkg/`
    - `deployments/`
3. **Feedback Utilisateur :** Le CLI doit afficher des messages d'étape clairs (ex: "📁 Création des répertoires...", "✅ Structure terminée").
4. **Gestion des Erreurs :**
    - Si le dossier cible existe déjà, afficher un message d'erreur en **Rouge** et quitter avec `os.Exit(1)`.
    - Si l'utilisateur ne fournit pas de nom de projet, afficher l'aide et quitter.

## Tasks / Subtasks

- [x] Gérer l'argument du nom de projet (AC: 1, 4)
  - [x] Extraire le premier argument non-flag comme `projectName`
  - [x] Vérifier si `projectName` est vide et afficher l'usage si nécessaire
- [x] Vérifier l'existence du répertoire cible (AC: 4)
  - [x] Utiliser `os.Stat` pour vérifier si le dossier existe déjà
  - [x] Afficher l'erreur en Rouge si nécessaire
- [x] Implémenter la création de l'arborescence (AC: 2, 3)
  - [x] Créer une liste des répertoires requis
  - [x] Boucler sur la liste et créer chaque répertoire avec `os.MkdirAll`
  - [x] Afficher un message de progression pour chaque étape majeure
- [x] Intégrer les couleurs ANSI (AC: 3, 4)
  - [x] Réutiliser les utilitaires de couleurs créés dans la Story 1.1

## Dev Notes

### Architecture & Constraints
- **Pattern :** Le générateur lui-même est un outil CLI simple. Le projet *généré* suit l'architecture Hexagonale Lite.
- **Paths :** Tous les répertoires doivent être créés relativement au nouveau dossier projet.
- **Naming :** `/internal/interfaces` est utilisé à la place de `/internal/ports` (Décision architecturale).

### Technical Guidelines
- Utiliser le package `os` pour la manipulation des fichiers/répertoires.
- Réutiliser la logique de couleurs de `cmd/create-go-starter/main.go`.
- S'assurer que les permissions des dossiers sont correctes (ex: `0755`).

### Project Structure Notes
- Le code source du générateur reste dans `cmd/create-go-starter/main.go`.
- Cette story se concentre uniquement sur la création des **dossiers**. La création des fichiers (`go.mod`, `main.go` du projet, etc.) fera l'objet de la Story 1.3.

### References
- [Epic 1: Project Initialization & Core Infrastructure](_bmad-output/planning-artifacts/epics.md)
- [Architecture Decision Document](_bmad-output/planning-artifacts/architecture.md#complete-project-directory-structure)
- [Project Context: interfaces MUST be defined in /internal/interfaces](_bmad-output/project-context.md)

## Dev Agent Record

### Agent Model Used
Claude Sonnet 4.5 (implementation)

### Debug Log References
None

### Implementation Plan
- Followed TDD (red-green-refactor) cycle:
  1. RED: Created failing tests for `createProjectStructure` function
  2. GREEN: Implemented minimal code to make tests pass
  3. REFACTOR: Improved code structure and updated existing tests
- Implemented directory creation with proper error handling
- Integrated with existing color utilities from Story 1.1
- Added comprehensive test coverage for happy path and error scenarios

### Completion Notes List
- ✅ Created `createProjectStructure()` function that generates hexagonal architecture directories
- ✅ Implemented error handling for existing directories with red error messages
- ✅ Added user feedback messages ("📁 Création des répertoires...", "✅ Structure terminée")
- ✅ Created comprehensive unit tests in `scaffold_test.go`:
  - TestCreateDirectories: validates all directories are created
  - TestProjectAlreadyExists: validates error when directory exists
  - TestCreateProjectStructureWithInvalidPath: validates error handling for invalid paths
- ✅ Updated existing integration tests to handle new output messages
- ✅ All tests pass (9/9)
- ✅ Linter passes with no warnings
- ✅ Manual integration testing confirms proper functionality

### File List
- cmd/create-go-starter/main.go (Modified - added createProjectStructure, validateProjectName functions and integration)
- cmd/create-go-starter/scaffold_test.go (Created - comprehensive unit tests for scaffolding and validation)
- cmd/create-go-starter/main_test.go (Modified - updated tests for new output messages and added invalid name tests)
- cmd/create-go-starter/colors_test.go (Unchanged - from Story 1.1, included for completeness)

## Senior Developer Review (AI)

### Review Date
2026-01-07

### Reviewer
Claude Sonnet 4.5 (Code Review Agent)

### Review Outcome
✅ **APPROVE** - All issues fixed automatically

### Review Summary
Performed adversarial code review and identified 12 issues (0 CRITICAL, 7 MEDIUM, 5 LOW). All issues were automatically corrected:

**Major Improvements Made:**
1. ✅ Added project name validation with regex pattern (security improvement)
2. ✅ Improved error messages with actionable suggestions
3. ✅ Added comprehensive test coverage for edge cases (13 total tests)
4. ✅ Fixed permission verification in tests
5. ✅ Improved test robustness (removed emoji dependencies, better invalid path handling)
6. ✅ Added documentation (godoc comments)
7. ✅ Extracted magic number to named constant (defaultDirPerm)
8. ✅ Corrected misleading code comments
9. ✅ Updated File List to include all modified files

### Action Items
All action items were resolved during the review. No outstanding issues.

### Files Reviewed
- cmd/create-go-starter/main.go ✅
- cmd/create-go-starter/scaffold_test.go ✅
- cmd/create-go-starter/main_test.go ✅
- cmd/create-go-starter/colors_test.go ✅

### Test Results After Review
- **Tests:** 13/13 passing ✅
- **Linter:** 0 warnings ✅
- **Coverage:** All ACs validated ✅

## Change Log
- 2026-01-07: Story implementation completed - Directory scaffolding with Hexagonal Architecture Lite structure, comprehensive tests, error handling, and user feedback messages
- 2026-01-07: Code review completed - Added project name validation, improved error messages, enhanced test coverage (12 issues fixed)
