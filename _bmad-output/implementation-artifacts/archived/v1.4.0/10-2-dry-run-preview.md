# Story 10.2: Dry-Run Preview

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want utiliser `--dry-run` pour prévisualiser sans générer de fichiers,
so that je puisse voir exactement ce qui sera créé avant de lancer la génération réelle.

## Acceptance Criteria

1. **Given** j'exécute `create-go-starter mon-projet --dry-run`, **When** la preview s'affiche, **Then** la liste complète des fichiers qui auraient été générés est affichée (path relatif depuis la racine du projet).
2. **Given** `--dry-run` est actif, **When** la liste de fichiers s'affiche, **Then** aucun fichier ni répertoire n'est créé sur le filesystem.
3. **Given** `--dry-run` est actif, **When** la liste de fichiers s'affiche, **Then** le résumé inclut: nombre total de fichiers, nombre de répertoires, template utilisé, database, observability level.
4. **Given** `--dry-run` est combiné avec `--template=minimal` ou `--database=sqlite`, **When** la preview s'affiche, **Then** seuls les fichiers du template/database concerné sont listés (pas les fichiers du template full).
5. **Given** `--dry-run` est actif, **When** la commande se termine, **Then** le code de sortie est 0 (succès) et un message clair indique "Dry-run completed. No files created."
6. **Given** `--dry-run` est actif et le projet existe déjà, **When** la preview s'affiche, **Then** un avertissement indique que le répertoire existe déjà, mais la preview continue (ne bloque pas comme en mode réel).

## Tasks / Subtasks

- [x] Task 1 : Ajouter le flag `--dry-run` dans le parsing d'arguments de `main()` (AC: #1, #2, #5)
  - [x] Détecter `-dry-run` et `--dry-run` dans la boucle de parsing `main.go:213-251`
  - [x] Mettre à jour `usage()` pour documenter `--dry-run`
  - [x] Mettre à jour les exemples CLI
- [x] Task 2 : Extraire la liste de fichiers de `generateFullTemplateFiles`, `generateMinimalTemplateFiles`, `generateGraphQLTemplateFiles` (AC: #1, #4)
  - [x] Créer `cmd/create-go-starter/dryrun.go`
  - [x] Implémenter `getFilesForTemplate(projectPath, projectName, template, database, observabilityLevel string) []FileGenerator` qui retourne la liste sans écrire
  - [x] Cette fonction est un refactoring : extraire la construction de `[]FileGenerator` de chaque `generate*Files()` dans une fonction `build*FileList()` séparée
  - [x] Les `generate*Files()` appellent `build*FileList()` puis écrivent les fichiers
- [x] Task 3 : Implémenter `runDryRun(projectName, template, database, observabilityLevel string) error` (AC: #1, #2, #3, #5, #6)
  - [x] Appeler `getFilesForTemplate()` pour obtenir la liste
  - [x] Afficher l'en-tête dry-run avec template/database/observability
  - [x] Afficher chaque fichier préfixé par `  📄 ` (chemin relatif)
  - [x] Afficher le résumé final: nb fichiers, nb répertoires uniques
  - [x] Afficher "Dry-run completed. No files created."
  - [x] Vérifier si le répertoire projet existe et afficher un warning (AC: #6)
- [x] Task 4 : Intégrer dans `main()` et `run()` (AC: #2, #5)
  - [x] Si `dryRun` flag détecté dans `main()`, appeler `runDryRun()` au lieu de `run()`
  - [x] Ne pas appeler `createProjectStructure()` ni `generateProjectFiles()` en dry-run
- [x] Task 5 : Tests (AC: #1-#6)
  - [x] Créer `cmd/create-go-starter/dryrun_test.go`
  - [x] Tester `getFilesForTemplate` pour chaque combinaison template/database
  - [x] Vérifier qu'aucun fichier n'est créé sur le filesystem (t.TempDir)
  - [x] Tester le comptage correct de fichiers et répertoires
  - [x] Tester le warning quand le répertoire existe

## Dev Notes

### Architecture et Patterns

**Fichier cible principal :** `cmd/create-go-starter/generator.go` (refactoring) + nouveau `cmd/create-go-starter/dryrun.go`

**Refactoring clé :** Les fonctions `generateFullTemplateFiles`, `generateMinimalTemplateFiles`, `generateGraphQLTemplateFiles` dans `generator.go` mélangent la **construction de la liste** et l'**écriture des fichiers**. Le refactoring les sépare :

**Avant (generator.go:114-319) :**
```go
func generateFullTemplateFiles(projectPath, projectName, database, observabilityLevel string) error {
    files := []FileGenerator{ ... } // construction
    for _, file := range files { os.WriteFile(...) } // écriture
}
```

**Après (refactoring) :**
```go
// generator.go — Construction uniquement
func buildFullFileList(projectPath, projectName, database, observabilityLevel string) []FileGenerator {
    // ... retourne la liste sans écrire
}

// generator.go — Écriture (appelle build + écrit)
func generateFullTemplateFiles(projectPath, projectName, database, observabilityLevel string) error {
    files := buildFullFileList(projectPath, projectName, database, observabilityLevel)
    return writeFiles(files) // helper commun
}

// dryrun.go — Dry-run (appelle build sans écrire)
func getFilesForTemplate(projectPath, projectName, template, database, obs string) []FileGenerator {
    switch template {
    case "full": return buildFullFileList(...)
    case "minimal": return buildMinimalFileList(...)
    case "graphql": return buildGraphQLFileList(...)
    }
}
```

**Helper `writeFiles()` :** Extraire la boucle d'écriture en un helper partagé :
```go
func writeFiles(files []FileGenerator) error {
    for _, file := range files {
        if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
            return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
        }
        if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
            return fmt.Errorf("failed to write file %s: %w", file.Path, err)
        }
    }
    return nil
}
```

**ATTENTION - Ne pas casser les tests existants :** Le refactoring doit maintenir le comportement exact de `generateFullTemplateFiles`, `generateMinimalTemplateFiles`, `generateGraphQLTemplateFiles`. Tous les tests existants doivent passer.

**Fichiers conditionnels (observability=advanced) :** `buildFullFileList` doit gérer la logique conditionnelle existante pour les fichiers observability (generator.go:131-140, 312-316).

### Affichage Dry-Run

```
══════════════════════════════════════════════════════
  DRY-RUN PREVIEW — No files will be created
══════════════════════════════════════════════════════
  Project:       mon-projet
  Template:      full
  Database:      postgres
  Observability: none

Files that would be created (34 files):

  📄 mon-projet/go.mod
  📄 mon-projet/cmd/main.go
  📄 mon-projet/pkg/config/env.go
  📄 mon-projet/pkg/logger/logger.go
  ... (liste complète)

Directories that would be created (12 unique):
  📁 mon-projet/cmd/
  📁 mon-projet/pkg/config/
  ...

══════════════════════════════════════════════════════
  ✅ Dry-run completed. No files were created.
  Run without --dry-run to generate the project.
══════════════════════════════════════════════════════
```

**Warning si répertoire existe :**
```
⚠️  Warning: Directory 'mon-projet' already exists.
   In real mode, this would cause an error.
   Continuing dry-run preview...
```

### Comptage des répertoires uniques

Pour calculer les répertoires uniques depuis `files []FileGenerator` :
```go
dirs := make(map[string]bool)
for _, f := range files {
    dirs[filepath.Dir(f.Path)] = true
}
// len(dirs) = nombre de répertoires uniques
```

**Chemins relatifs :** Utiliser `filepath.Rel(projectPath, file.Path)` ou simplement afficher le path tel quel (il commence déjà par le projectName).

### Project Structure Notes

- **Nouveau fichier :** `cmd/create-go-starter/dryrun.go`
- **Nouveau fichier test :** `cmd/create-go-starter/dryrun_test.go`
- **Fichier modifié :** `cmd/create-go-starter/generator.go` (refactoring `generate*Files` → `build*FileList` + `writeFiles`)
- **Fichier modifié :** `cmd/create-go-starter/main.go` (ajout parsing `--dry-run`, appel `runDryRun()`)
- **Pas de nouveaux répertoires** requis

### Setup.sh et chmod

La fonction `generateFullTemplateFiles` fait un `chmod 0755` sur `setup.sh` après écriture (generator.go:306-309). En dry-run, ce chmod ne s'applique pas. S'assurer que `buildFullFileList` inclut `setup.sh` dans la liste des fichiers comme un fichier normal (sans le chmod — c'est une opération post-écriture).

### Testing Approach

```go
func TestDryRunCreatesNoFiles(t *testing.T) {
    tmpDir := t.TempDir()
    projectName := "test-dry-run"
    err := runDryRun(projectName, "full", "postgres", "none")
    // Vérifier qu'aucun répertoire n'a été créé dans le tmpDir
    entries, _ := os.ReadDir(tmpDir)
    assert(t, len(entries) == 0, "expected no files created")
}

func TestGetFilesForTemplateFullCount(t *testing.T) {
    files := getFilesForTemplate("/tmp/p", "p", "full", "postgres", "none")
    assert(t, len(files) > 30, "full template should have many files")
}
```

### References

- `FileGenerator` struct : [Source: cmd/create-go-starter/generator.go:76-79]
- `generateFullTemplateFiles` (à refactorer) : [Source: cmd/create-go-starter/generator.go:114-319]
- `generateMinimalTemplateFiles` (à refactorer) : [Source: cmd/create-go-starter/generator.go:402-510]
- `generateGraphQLTemplateFiles` (à refactorer) : [Source: cmd/create-go-starter/generator.go:514-662]
- Parsing flags `main()` : [Source: cmd/create-go-starter/main.go:213-251]
- `createProjectStructure()` (à NE PAS appeler en dry-run) : [Source: cmd/create-go-starter/main.go:135-161]
- Epic 10 description : [Source: _bmad-output/planning-artifacts/epics.md#Story-10.2]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

Aucun problème notable - le refactoring s'est déroulé sans difficulté.

### Completion Notes List

- ✅ Task 1 : Flag `--dry-run` ajouté dans `main.go` avec parsing `-dry-run`/`--dry-run` et documentation `usage()`
- ✅ Task 2 : Refactoring complet de `generator.go` - extraction de `buildFullFileList`, `buildMinimalFileList`, `buildGraphQLFileList`, `buildObservabilityFileList` et helper `writeFiles`. Les 3 fonctions `generate*TemplateFiles` appellent désormais `build*FileList` + `writeFiles`.
- ✅ Task 3 : `runDryRun` implémenté dans `dryrun.go` avec affichage formaté (header, liste fichiers 📄, répertoires 📁, footer), warning AC#6, retour nil (exit 0)
- ✅ Task 4 : Intégration dans `main()` - après les validations, si `dryRun=true`, appel de `runDryRun()` avec `return` immédiat (pas de `createProjectStructure` ni `generateProjectFiles`)
- ✅ Task 5 : 17 tests dans `dryrun_test.go` couvrant tous les ACs (13 originaux + 4 tests stdout ajoutés en review) - 100% PASS, aucune régression suite complète
- Résultat dry-run : full=36 fichiers/18 dirs, minimal=20 fichiers/9 dirs, graphql=28 fichiers/13 dirs

### File List

- `cmd/create-go-starter/dryrun.go` (nouveau)
- `cmd/create-go-starter/dryrun_test.go` (nouveau, 17 tests)
- `cmd/create-go-starter/generator.go` (modifié - refactoring complet)
- `cmd/create-go-starter/main.go` (modifié - flag `--dry-run`, conflit `--interactive`, exemples usage)
- `docs/usage.md` (modifié - ajout flag `--dry-run` et exemples)
- `docs/cli-architecture.md` (modifié - checkbox dry-run mode)

## Senior Developer Review (AI)

### Review Model Used

claude-opus-4.6 (github-copilot/claude-opus-4.6)

### Review Date

2026-02-18

### Review Type

Adversarial code review (BMAD workflow: code-review)

### Findings Summary

| ID | Severity | Description | Resolution |
|----|----------|-------------|------------|
| H1 | HIGH | Aucun test ne capture stdout pour vérifier le contenu affiché (AC#1, AC#3, AC#5) | FIXED: Ajout de `captureStdout` helper + 4 tests stdout (`TestDryRunOutputContainsFileList`, `TestDryRunOutputContainsSummary`, `TestDryRunOutputContainsCompletionMessage`, `TestDryRunOutputContainsWarningForExistingDir`) |
| H2 | HIGH | Conflit `--dry-run` + `--interactive` non géré | FIXED: Ajout de la vérification de conflit dans `main.go` avec message d'erreur clair |
| M1 | MEDIUM | Pas d'exemple `--dry-run` dans la section Examples de `usage()` | FIXED: Ajout de 2 exemples dry-run dans `usage()` |
| M2 | MEDIUM | Documentation non mise à jour (README, docs/usage.md, docs/cli-architecture.md) | FIXED: Mise à jour de `docs/usage.md` (flag + exemples) et `docs/cli-architecture.md` (checkbox) |
| M3 | MEDIUM | Message AC#5 dit "No files were created" au lieu de "No files created" | FIXED: Aligné le message dans `dryrun.go` avec l'AC exacte |
| M4 | MEDIUM | Story indique "11 tests" mais il y en a 17 (13 originaux + 4 ajoutés en review) | FIXED: Mise à jour du Dev Agent Record |
| L1 | LOW | `runDryRun` utilise un chemin relatif pour `os.Stat` (pattern pré-existant) | ACCEPTED: Pattern cohérent avec le reste du CLI |
| L2 | LOW | Le comptage de répertoires inclut le répertoire racine du projet | ACCEPTED: Comportement raisonnable et documenté |

### Issues Fixed

6 issues corrigées (2 HIGH, 4 MEDIUM). 2 issues LOW acceptées.

## Change Log

- 2026-02-18: Story 10.2 implémentée - Ajout du flag `--dry-run` avec preview des fichiers sans génération. Refactoring de `generator.go` pour extraire `build*FileList` + `writeFiles`. Nouveau `dryrun.go` avec `getFilesForTemplate` + `runDryRun`. 17 tests ajoutés, 100% PASS.
- 2026-02-18: Code review adversarial (claude-opus-4.6) - 8 issues trouvées (2 HIGH, 4 MEDIUM, 2 LOW). 6 issues corrigées automatiquement : tests stdout ajoutés (H1), conflit --interactive/--dry-run (H2), exemples usage (M1), documentation mise à jour (M2), message AC#5 aligné (M3), compteur tests corrigé (M4). Story passe de `review` à `done`.
