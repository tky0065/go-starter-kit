# Story 10.4: Visual Feedback & Progress

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want voir une progress bar et des statistiques pendant et après la génération du projet,
so que je sache exactement ce qui se passe et puisse évaluer la performance de la génération.

## Acceptance Criteria

1. **Given** la génération démarre, **When** les fichiers sont créés un par un, **Then** une progress bar est affichée dans le terminal indiquant le nombre de fichiers traités sur le total (ex: `[██████░░░░] 18/34 files`).
2. **Given** la génération est terminée, **When** le message final s'affiche, **Then** les statistiques incluent: nombre total de fichiers générés, taille totale estimée (en KB/MB), durée de la génération.
3. **Given** la génération est terminée, **When** les statistiques s'affichent, **Then** chaque étape principale (structure créée, fichiers générés, git initialisé) affiche son temps d'exécution partiel.
4. **Given** le terminal ne supporte pas les ANSI codes (NO_COLOR=1 ou non-TTY), **When** la progress bar est désactivée, **Then** uniquement le texte est affiché sans artefacts visuels.
5. **Given** la génération avec `--dry-run`, **When** la preview s'affiche, **Then** la progress bar n'est pas affichée (dry-run est instantané, pas de vraie génération).
6. **Given** la génération avec `--interactive` suivie de la génération réelle, **When** les fichiers sont créés, **Then** la progress bar apparaît normalement.

## Tasks / Subtasks

- [x] Task 1 : Implémenter la progress bar en stdlib (AC: #1, #4)
  - [x] Créer `cmd/create-go-starter/progress.go`
  - [x] Implémenter `ProgressBar` struct avec `total`, `current`, `width` (défaut: 30 chars)
  - [x] Implémenter `progressBar.Update(current int)` qui affiche et rafraîchit la barre
  - [x] Implémenter `progressBar.Complete()` qui affiche la barre à 100% et passe à la ligne
  - [x] Détecter si TTY est disponible via `isatty(os.Stdout.Fd())` ou vérifier `NO_COLOR` env var
  - [x] Implémenter `isTerminal() bool` helper
- [x] Task 2 : Implémenter le timing et statistiques (AC: #2, #3)
  - [x] Créer `cmd/create-go-starter/stats.go`
  - [x] Implémenter `GenerationStats` struct: `StartTime time.Time`, `StepTimes map[string]time.Duration`, `FilesGenerated int`, `TotalSize int64`
  - [x] Implémenter `stats.StartStep(name string)` et `stats.EndStep(name string)`
  - [x] Implémenter `stats.AddFile(path string)` pour accumuler taille via `os.Stat`
  - [x] Implémenter `stats.Display()` pour afficher le rapport final
- [x] Task 3 : Intégrer la progress bar dans `generateFullTemplateFiles` et consorts (AC: #1, #2)
  - [x] Modifier `writeFiles(files []FileGenerator)` (introduit en story 10.2) pour accepter un callback de progression ou un `ProgressBar`
  - [x] Signature: `writeFiles(files []FileGenerator, onProgress func(current, total int)) error`
  - [x] Appeler `onProgress(i+1, len(files))` dans la boucle
  - [x] En mode non-TTY ou si `onProgress == nil`: pas de progress bar
- [x] Task 4 : Intégrer les stats dans la fonction `run()` (AC: #2, #3)
  - [x] Modifier `run()` pour instancier `GenerationStats`
  - [x] Enregistrer le timing des étapes: "Creating directories", "Generating files", "Git initialization"
  - [x] Afficher les statistiques après `printSuccessMessage`
- [x] Task 5 : Affichage amélioré du message de succès (AC: #2, #3)
  - [x] Modifier `printSuccessMessage()` ou créer une version qui inclut les stats
  - [x] Afficher les stats en dessous du message de succès (ne pas remplacer le message existant)
- [x] Task 6 : Tests (AC: #1-#6)
  - [x] Créer `cmd/create-go-starter/progress_test.go`
  - [x] Tester `ProgressBar.Update()` avec différentes valeurs
  - [x] Tester `isTerminal()`
  - [x] Tester `GenerationStats` timing et affichage
  - [x] Tester que les stats sont correctes après une génération (avec t.TempDir)

## Dev Notes

### Architecture et Patterns

**Fichier cible principal :** Nouveaux `cmd/create-go-starter/progress.go` et `stats.go`, modification `cmd/create-go-starter/generator.go` et `main.go`

**IMPORTANT — Dépendance sur Story 10.2 :** Cette story réutilise le refactoring de `generator.go` introduit en 10.2 (fonction `writeFiles()`). Si 10.2 n'est pas encore complète, implémenter le refactoring `writeFiles` ici aussi.

**Implémentation Progress Bar (stdlib uniquement, pas de dépendances externes) :**

```go
// cmd/create-go-starter/progress.go
package main

import (
    "fmt"
    "os"
    "strings"
)

// ProgressBar is a simple terminal progress bar
type ProgressBar struct {
    total   int
    current int
    width   int
    enabled bool
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total, width int) *ProgressBar {
    return &ProgressBar{
        total:   total,
        width:   width,
        enabled: isTerminal() && os.Getenv("NO_COLOR") == "",
    }
}

// Update refreshes the progress bar display
func (pb *ProgressBar) Update(current int) {
    if !pb.enabled {
        return
    }
    pb.current = current
    percent := float64(current) / float64(pb.total)
    filled := int(percent * float64(pb.width))
    empty := pb.width - filled

    bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
    fmt.Printf("\r  [%s] %d/%d files", bar, current, pb.total)
}

// Complete marks the progress bar as done
func (pb *ProgressBar) Complete() {
    if !pb.enabled {
        return
    }
    pb.Update(pb.total)
    fmt.Println() // newline after progress bar
}
```

**Détection TTY — Approche simple :**
```go
// isTerminal returns true if stdout is a terminal (not piped/redirected)
func isTerminal() bool {
    fileInfo, err := os.Stdout.Stat()
    if err != nil {
        return false
    }
    return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
```

**Stats struct :**
```go
// cmd/create-go-starter/stats.go
package main

import (
    "fmt"
    "time"
    "os"
)

type GenerationStats struct {
    StartTime    time.Time
    stepStarts   map[string]time.Time
    StepDurations map[string]time.Duration
    FilesCount   int
    TotalBytes   int64
}

func NewGenerationStats() *GenerationStats {
    return &GenerationStats{
        StartTime:     time.Now(),
        stepStarts:    make(map[string]time.Time),
        StepDurations: make(map[string]time.Duration),
    }
}

func (s *GenerationStats) StartStep(name string) {
    s.stepStarts[name] = time.Now()
}

func (s *GenerationStats) EndStep(name string) {
    if start, ok := s.stepStarts[name]; ok {
        s.StepDurations[name] = time.Since(start)
    }
}

func (s *GenerationStats) AddFile(path string) {
    s.FilesCount++
    if info, err := os.Stat(path); err == nil {
        s.TotalBytes += info.Size()
    }
}

func (s *GenerationStats) Display() {
    total := time.Since(s.StartTime)
    sizeStr := formatBytes(s.TotalBytes)

    fmt.Printf("\n%s\n", Green("📊 Generation Statistics:"))
    fmt.Printf("   Files generated : %d\n", s.FilesCount)
    fmt.Printf("   Total size      : %s\n", sizeStr)
    fmt.Printf("   Total time      : %s\n", total.Round(time.Millisecond))

    if len(s.StepDurations) > 0 {
        fmt.Println("   Step breakdown  :")
        for step, dur := range s.StepDurations {
            fmt.Printf("     • %-25s %s\n", step+":", dur.Round(time.Millisecond))
        }
    }
}

func formatBytes(bytes int64) string {
    switch {
    case bytes >= 1024*1024:
        return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1024*1024))
    case bytes >= 1024:
        return fmt.Sprintf("%.1f KB", float64(bytes)/float64(1024))
    default:
        return fmt.Sprintf("%d bytes", bytes)
    }
}
```

**Intégration dans `run()` (main.go:343-395) :**
```go
func run(projectName, template, database, observabilityLevel string) error {
    stats := NewGenerationStats()

    // ... existing code ...

    stats.StartStep("Creating directories")
    if err := createProjectStructure(projectPath, template); err != nil { return err }
    stats.EndStep("Creating directories")

    // Calculer le nombre total de fichiers pour la progress bar
    fileList := getFilesForTemplate(projectPath, projectName, template, database, observabilityLevel)
    pb := NewProgressBar(len(fileList), 30)

    stats.StartStep("Generating files")
    if err := generateProjectFilesWithProgress(projectPath, projectName, template, database, observabilityLevel, pb, stats); err != nil {
        return err
    }
    stats.EndStep("Generating files")

    // ... git init ...
    stats.StartStep("Git initialization")
    // ... existing git code ...
    stats.EndStep("Git initialization")

    printSuccessMessage(projectName, database)
    stats.Display()

    return nil
}
```

**Signature `writeFiles` avec progress callback (réutilisé de story 10.2) :**
```go
func writeFiles(files []FileGenerator, onProgress func(current, total int)) error {
    for i, file := range files {
        if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
            return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
        }
        if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
            return fmt.Errorf("failed to write file %s: %w", file.Path, err)
        }
        if onProgress != nil {
            onProgress(i+1, len(files))
        }
    }
    return nil
}
```

### Affichage Expected

```
Creating project: my-project (template: full, database: postgres, observability: none)
📁 Creating directories...
✅ Structure created
📝 Generating core files...
  [████████████████████████░░░░░░] 34/34 files
✅ Files generated successfully
🔑 Configuring environment...
🔧 Initializing Git repository...
✅ Git repository initialized with initial commit

════════════════════════════════════════════════════════════════
🎉 Project 'my-project' created successfully!
════════════════════════════════════════════════════════════════

📊 Generation Statistics:
   Files generated : 34
   Total size      : 87.3 KB
   Total time      : 1.247s
   Step breakdown  :
     • Creating directories:         2ms
     • Generating files:             1.183s
     • Git initialization:           62ms

📋 Next steps - Initial setup:
...
```

### Comportement Non-TTY

Quand `isTerminal()` retourne false (ex: pipeline CI, `| tee output.txt`) :
- La progress bar n'est pas affichée (`\r` causerait des artefacts dans les logs)
- Les stats sont toujours affichées (texte simple)
- `NO_COLOR=1` désactive aussi la progress bar

### Notes sur `run()` actuel

La fonction `run()` (main.go:343-395) affiche déjà des messages de progression avec `fmt.Println`. Ces messages **restent** — la progress bar s'insère **entre** "📝 Generating core files..." et "✅ Files generated successfully". La progress bar est sur la ligne entre ces deux messages, overwritten par `\r`.

**ATTENTION `\r` et les messages existants :** Le `\r` de la progress bar ramène le curseur au début de la ligne courante. S'assurer que la progress bar affiche sur sa propre ligne (sans interférer avec les messages `fmt.Println` avant/après).

### Project Structure Notes

- **Nouveau fichier :** `cmd/create-go-starter/progress.go`
- **Nouveau fichier :** `cmd/create-go-starter/stats.go`
- **Nouveau fichier test :** `cmd/create-go-starter/progress_test.go`
- **Fichier modifié :** `cmd/create-go-starter/main.go` (intégration stats dans `run()`)
- **Fichier modifié :** `cmd/create-go-starter/generator.go` (callback progress dans `writeFiles`)
- **Pas de nouvelles dépendances** — `time`, `os`, `fmt`, `strings` sont stdlib
- **Note sur story 10.2 :** Si `writeFiles()` n'existe pas encore (story 10.2 non implémentée), créer ici aussi

### Testing Approach

```go
func TestProgressBarUpdate(t *testing.T) {
    pb := &ProgressBar{total: 10, width: 10, enabled: false} // disabled for test
    pb.Update(5) // Should not panic
}

func TestProgressBarComplete(t *testing.T) {
    pb := &ProgressBar{total: 10, width: 10, enabled: false}
    pb.Complete() // Should not panic
}

func TestGenerationStatsDisplay(t *testing.T) {
    stats := NewGenerationStats()
    stats.StartStep("test")
    time.Sleep(1 * time.Millisecond)
    stats.EndStep("test")
    stats.FilesCount = 5
    stats.TotalBytes = 1024
    // Capture stdout and verify output
    stats.Display()
}

func TestIsTerminal(t *testing.T) {
    // In CI, isTerminal() should return false
    // This is just a smoke test to ensure no panic
    _ = isTerminal()
}

func TestFormatBytes(t *testing.T) {
    tests := []struct{ bytes int64; expected string }{
        {500, "500 bytes"},
        {1500, "1.5 KB"},
        {1500000, "1.4 MB"},
    }
    for _, tt := range tests {
        result := formatBytes(tt.bytes)
        if result != tt.expected {
            t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, result, tt.expected)
        }
    }
}
```

### References

- Fonction `run()` à modifier : [Source: cmd/create-go-starter/main.go:343-395]
- `writeFiles()` (introduit Story 10.2) : [Source: cmd/create-go-starter/generator.go]
- `getFilesForTemplate()` (introduit Story 10.2) : [Source: cmd/create-go-starter/dryrun.go]
- `printSuccessMessage()` : [Source: cmd/create-go-starter/main.go:471-534]
- Epic 10 description : [Source: _bmad-output/planning-artifacts/epics.md#Story-10.4]
- Architecture CLI : [Source: _bmad-output/planning-artifacts/architecture.md]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

- ✅ `progress.go` créé : `ProgressBar` struct avec `Update()`, `Complete()`, `isTerminal()` — stdlib uniquement
- ✅ `stats.go` créé : `GenerationStats` avec timing par étape, comptage de fichiers, affichage final
- ✅ `generator.go` modifié : `writeFiles` accepte maintenant `onProgress func(current, total int)` — tous les appelants existants passent `nil`
- ✅ `main.go` modifié : `run()` instancie `GenerationStats`, chronomètre chaque étape, affiche les stats après `printSuccessMessage`
- ✅ `progress_test.go` créé : 14 tests couvrant ProgressBar, isTerminal, formatBytes, GenerationStats, writeFiles avec callback
- ✅ Tests : tous passent (`go test -short ./cmd/create-go-starter/` : OK)
- ✅ `go vet` : aucun avertissement
- ✅ E2E validé : stats affichées correctement (ex: `Files generated: 20, Total size: 30.1 KB, Total time: 61ms`)
- ✅ Non-TTY : progress bar désactivée automatiquement (pas d'artefacts `\r` dans les logs)
- ✅ Dry-run : aucune progress bar (dry-run n'appelle pas `writeFiles`)
- Note : `formatBytes` utilise des seuils stricts (>= 1024 pour KB, >= 1024*1024 pour MB)

### Code Review Fixes Applied

- ✅ Fix: Guard `ProgressBar.Update()` against `total <= 0` to prevent division by zero (NaN)
- ✅ Fix: `stats.go` `Display()` now uses dynamic `stepOrder` instead of hardcoded step list
- ✅ Fix: `isTerminalFn` package variable allows test mocking of TTY detection
- ✅ Fix: Added `TestNewProgressBarEnabledWhenTTYAndNoColorAbsent` and `TestNewProgressBarDisabledWhenNotTTY` tests
- ✅ Fix: `TestProgressBarTotalZeroNoPanic` now tests both `enabled: false` and `enabled: true`
- ✅ Fix: Added documentation note on `generateProjectFiles()` being kept for API usage
- ✅ Fix: Updated File List to include all modified files

### File List

- `cmd/create-go-starter/progress.go` (nouveau)
- `cmd/create-go-starter/stats.go` (nouveau)
- `cmd/create-go-starter/progress_test.go` (nouveau)
- `cmd/create-go-starter/generator.go` (modifié : signature `writeFiles`, `buildFullFileList`, `buildMinimalFileList`, `buildGraphQLFileList`)
- `cmd/create-go-starter/main.go` (modifié : `run()` avec stats et progress bar, doctor/interactive/dry-run intégration)
- `docs/cli-architecture.md` (modifié : documentation mise à jour)
- `docs/usage.md` (modifié : documentation mise à jour)
