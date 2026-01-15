# Story 6.5: Documentation et aide CLI

Status: done

## Story

**En tant que** développeur,
**Je veux** trouver une documentation claire sur les options de templates disponibles (minimal, full, graphql),
**Afin de** choisir le point de départ le plus adapté à mon nouveau projet.

## Acceptance Criteria

1. **AC1**: Given je consulte le `README.md`, When je cherche la section "Utilisation", Then je trouve l'explication du flag `--template` et la liste des valeurs possibles.
2. **AC2**: Given je consulte `docs/usage.md`, When je lis la section "Options disponibles", Then elle est à jour et liste les templates avec une brève description de ce que chacun contient.
3. **AC3**: Given j'exécute `create-go-starter --help`, When la commande s'exécute, Then la sortie correspond à la documentation.
4. **AC4**: Given la Roadmap du `README.md`, Then le point "Templates multiples" est marqué comme terminé.

## Tasks / Subtasks

- [x] Task 1: Mettre à jour README.md (AC: 1, 4)
  - [x] 1.1 Mettre à jour la section "Utilisation de base" avec `--template`
  - [x] 1.2 Mettre à jour la Roadmap en marquant "Templates multiples" comme fait
  - [x] 1.3 Ajouter une brève mention des 3 types de templates

- [x] Task 2: Mettre à jour docs/usage.md (AC: 2)
  - [x] 2.1 Ajouter une section "Templates disponibles"
  - [x] 2.2 Créer un tableau comparatif des fonctionnalités (Auth, Swagger, GraphQL, etc.) pour chaque template
  - [x] 2.3 Expliquer les différences structurelles majeures

- [x] Task 3: Vérification de cohérence (AC: 3)
  - [x] 3.1 Vérifier que les descriptions dans `main.go` correspondent à la documentation
  - [x] 3.2 Vérifier que l'aide CLI affiche correctement les options valides

## Dev Notes

### Stratégie de Documentation

Utiliser un tableau comparatif dans `docs/usage.md` pour clarifier les différences en un coup d'œil :

| Feature | Minimal | Full (Défaut) | GraphQL |
|---------|---------|---------------|---------|
| REST API | ✅ | ✅ | ❌ |
| GraphQL API | ❌ | ❌ | ✅ |
| JWT Auth | ❌ | ✅ | ❌ |
| Swagger | ✅ | ✅ | ❌ |
| DB (GORM) | ✅ | ✅ | ✅ |

### Cohérence

- Les termes (`minimal`, `full`, `graphql`) doivent être identiques partout (code, doc, help).
- Ne pas documenter de fonctionnalités futures ou hypothétiques pour éviter la confusion.

### References

- Story 6.1 (Flag CLI)
- Story 6.2 (Template Minimal)
- Story 6.3 (Refactor Full)
- Story 6.4 (Template GraphQL)

## Dev Agent Record

### Agent Model Used

Claude 3.7 Sonnet (dev agent)

### Completion Notes List

**Story 6.5: Documentation et aide CLI - Completed Successfully**

**Objective**: Provide clear documentation on available template options (minimal, full, graphql) so developers can choose the right starting point for their projects.

**Tasks Completed**:

1. **README.md Updates** (AC1, AC4):
   - ✅ Added comprehensive "Choisir un template" section with flag examples
   - ✅ Created comparison table of all 3 templates with descriptions and use cases
   - ✅ Updated Roadmap to mark "Templates multiples" as completed [x]
   - ✅ Added link to detailed usage guide for template differences

2. **docs/usage.md Updates** (AC2):
   - ✅ Added complete "Templates disponibles" section with:
     - Overview table with use cases
     - Detailed feature comparison table (14 features × 3 templates)
     - In-depth structural differences for each template
     - Generated endpoints documentation for each template
     - "How to choose the right template" decision guide
   - ✅ Updated "Options disponibles" section with --template flag examples
   - ✅ Documented default behavior (full template when no flag)

3. **Consistency Verification** (AC3):
   - ✅ Verified CLI --help output matches documentation exactly:
     - Template names: minimal, full, graphql
     - Descriptions match between main.go and docs
     - Default template (full) clearly indicated
   - ✅ Confirmed terminology consistency across all files

**Verification Results**:

```bash
# AC3 Verification: CLI help output
$ go run ./cmd/create-go-starter --help

Usage: create-go-starter [options] <project-name>

Options:
  -h	Show help message (shorthand)
  -help
    	Show help message
  -template string
    	Template type to generate (default "full")

Templates:
  minimal   Basic REST API with Swagger (no authentication)
  full      Complete API with JWT auth, user management, and Swagger (default)
  graphql   GraphQL API with gqlgen and GraphQL Playground
```

✅ Output matches documentation exactly (AC3 verified)

**Definition of Done Checklist**:

- [x] **AC1**: README.md "Utilisation" section explains --template flag and lists all values
- [x] **AC2**: docs/usage.md "Options disponibles" is up-to-date with template descriptions
- [x] **AC3**: CLI --help output corresponds to documentation
- [x] **AC4**: README.md Roadmap marks "Templates multiples" as completed
- [x] All tasks and subtasks completed (9/9)
- [x] Documentation is comprehensive and user-friendly
- [x] Terminology is consistent (minimal, full, graphql) everywhere
- [x] No future/hypothetical features documented (only what exists)

**Quality Notes**:

- **Comprehensive comparison**: Created detailed 14-feature comparison table helping users make informed decisions
- **Structured decision guide**: "Comment choisir le bon template?" section guides users based on project needs
- **Endpoint documentation**: Each template's generated endpoints are fully documented
- **Consistency**: All terminology matches between code (main.go constants), CLI help, README, and docs/usage.md
- **French documentation**: All user-facing docs in French as per project convention
- **No breaking changes**: Documentation-only updates, no code changes required

**References**:
- Story 6.1: Flag CLI implementation (--template flag)
- Story 6.2: Template Minimal implementation
- Story 6.3: Refactor Full template
- Story 6.4: Template GraphQL implementation

### File List

**Files Modified**:

1. `README.md`
   - Added "Choisir un template" section with examples
   - Added template comparison table
   - Updated Roadmap to mark "Templates multiples" as [x] completed
   - Removed duplicate items from Roadmap (GraphQL support, api-only)

2. `docs/usage.md`
   - Added comprehensive "Templates disponibles" section (180+ lines)
   - Updated "Options disponibles" with --template flag documentation
   - Added detailed feature comparison table
   - Added structural differences for each template
   - Added endpoint documentation for each template
   - Added decision guide for choosing templates

3. `_bmad-output/implementation-artifacts/6-5-documentation-et-aide-cli.md`
   - Changed status from `ready-for-dev` to `done`
   - Marked all tasks and subtasks as completed
   - Added this Dev Agent Record

**No code changes required** - This was a documentation-only story ensuring users understand the three available templates.

### Change Log

**2026-01-15 - Story 6.5 Documentation et aide CLI**

**Documentation Updates**:

1. **README.md**:
   - Added section "Choisir un template" after "Créer un nouveau projet"
   - Template comparison table (3 templates × 3 columns)
   - Roadmap: Moved "Templates multiples" to "Fonctionnalités complétées" section
   - Cross-reference link to docs/usage.md#templates-disponibles

2. **docs/usage.md**:
   - New section "Templates disponibles" with:
     - Overview table with use cases
     - 14-feature comparison matrix (REST, GraphQL, JWT, Swagger, etc.)
     - Structural differences (minimal: 6 endpoints, full: 8 endpoints with auth, graphql: 2 endpoints)
     - Generated schema documentation for GraphQL template
     - Decision guide: "Comment choisir le bon template?"
   - Updated "Options disponibles" section with --template examples

**Verification**:
- ✅ CLI help output matches documentation (AC3)
- ✅ All 4 Acceptance Criteria met
- ✅ All 9 tasks/subtasks completed
- ✅ Terminology consistent across all files

**Story Status**: `done` (ready for Epic 6 completion)
