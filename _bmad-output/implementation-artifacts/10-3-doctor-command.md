# Story 10.3: Doctor Command

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want exécuter `create-go-starter doctor`,
so que le CLI vérifie automatiquement mon environnement de développement et m'indique les problèmes détectés avec des solutions.

## Acceptance Criteria

1. **Given** j'exécute `create-go-starter doctor`, **When** le diagnostic s'exécute, **Then** Go version est vérifiée (minimum Go 1.21 requis) avec le résultat affiché (✅ ou ❌).
2. **Given** j'exécute `create-go-starter doctor`, **When** le diagnostic s'exécute, **Then** Git est vérifié (présence de `git` dans le PATH) avec le résultat et la version affichée.
3. **Given** j'exécute `create-go-starter doctor`, **When** le diagnostic s'exécute, **Then** Docker est vérifié (présence de `docker` dans le PATH et si le daemon est actif) avec le résultat et la version affichée.
4. **Given** qu'un problème est détecté, **When** le rapport s'affiche, **Then** une solution concrète est suggérée pour chaque problème (lien de téléchargement ou commande d'installation).
5. **Given** le diagnostic est terminé, **When** le rapport final s'affiche, **Then** un score de santé est affiché: "All checks passed ✅" ou "X issues found ⚠️".
6. **Given** j'exécute `create-go-starter doctor`, **When** le diagnostic s'exécute, **Then** la version de `create-go-starter` elle-même est affichée dans le rapport.
7. **Given** que tous les checks passent, **When** le rapport est affiché, **Then** le code de sortie est 0. **Given** qu'au moins un check échoue (non-fatal), **Then** le code de sortie est 1.

## Tasks / Subtasks

- [x] Task 1 : Détecter la sous-commande `doctor` dans `main()` (AC: #1-#7)
  - [x] Ajouter la détection de `doctor` dans le bloc `if len(os.Args) > 1 && os.Args[1] == "add-model"` (main.go:195-201) — pattern identique à `add-model`
  - [x] Ajouter l'appel `runDoctor()` et exit approprié
  - [x] Mettre à jour `usage()` pour mentionner la commande `doctor`
- [x] Task 2 : Créer `cmd/create-go-starter/doctor.go` avec la logique principale (AC: #1-#7)
  - [x] Définir la struct `CheckResult { Name string; Status bool; Version string; Message string; Fix string }`
  - [x] Implémenter `runDoctor() int` (retourne exit code)
  - [x] Implémenter `checkGoVersion() CheckResult`
  - [x] Implémenter `checkGit() CheckResult`
  - [x] Implémenter `checkDocker() CheckResult`
  - [x] Implémenter `displayDoctorReport(results []CheckResult)`
  - [x] Inclure la version de `create-go-starter` dans le rapport (AC: #6)
- [x] Task 3 : Implémenter chaque check (AC: #1, #2, #3, #4)
  - [x] `checkGoVersion()`: utiliser `exec.Command("go", "version")` + parser la version + vérifier >= 1.21
  - [x] `checkGit()`: utiliser `exec.Command("git", "--version")` + vérifier succès
  - [x] `checkDocker()`: utiliser `exec.Command("docker", "--version")` pour le binaire + `exec.Command("docker", "info")` pour le daemon
  - [x] Pour chaque check KO: fournir un message `Fix` avec solution concrète
- [x] Task 4 : Affichage et rapport (AC: #4, #5, #6, #7)
  - [x] Afficher chaque check avec `✅ OK` ou `❌ FAIL` en couleur
  - [x] Afficher version quand disponible
  - [x] Afficher message de fix en jaune quand nécessaire
  - [x] Afficher résumé final: score + message global
  - [x] Retourner exit code 0 si tout OK, 1 si problèmes
- [x] Task 5 : Versionning de l'outil (AC: #6)
  - [x] Définir une constante `Version = "1.4.0-dev"` dans un fichier `version.go` ou dans `main.go`
  - [x] Afficher dans le rapport doctor
- [x] Task 6 : Tests (AC: #1-#7)
  - [x] Créer `cmd/create-go-starter/doctor_test.go`
  - [x] Tester `checkGoVersion()` (mock exec ou integration test)
  - [x] Tester `checkGit()` avec git disponible
  - [x] Tester l'affichage du rapport avec des CheckResults connus
  - [x] Tester le calcul du exit code

## Dev Notes

### Architecture et Patterns

**Fichier cible principal :** Nouveau `cmd/create-go-starter/doctor.go` + modification `cmd/create-go-starter/main.go`

**Pattern de sous-commande existant (main.go:195-201) :**
```go
if len(os.Args) > 1 && os.Args[1] == "add-model" {
    if err := runAddModel(os.Args[2:]); err != nil {
        fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("%v", err)))
        os.Exit(1)
    }
    return
}
```

**Ajouter juste après (main.go:201-202) :**
```go
if len(os.Args) > 1 && os.Args[1] == "doctor" {
    exitCode := runDoctor()
    os.Exit(exitCode)
}
```

**Structure `CheckResult` :**
```go
type CheckResult struct {
    Name    string // "Go", "Git", "Docker"
    Status  bool   // true = OK, false = FAIL
    Version string // e.g., "go1.25.5 linux/amd64"
    Message string // Message d'erreur si Status=false
    Fix     string // Solution suggérée
}
```

**Implémentation des checks avec `os/exec` :**
```go
import (
    "os/exec"
    "strings"
)

func checkGoVersion() CheckResult {
    cmd := exec.Command("go", "version")
    output, err := cmd.Output()
    if err != nil {
        return CheckResult{
            Name:    "Go",
            Status:  false,
            Message: "Go not found in PATH",
            Fix:     "Download Go from https://go.dev/dl/ and add to PATH",
        }
    }
    version := strings.TrimSpace(string(output))
    // Parse "go version go1.25.5 linux/amd64" → "go1.25.5"
    // Vérifier >= 1.21
    return CheckResult{Name: "Go", Status: true, Version: version}
}
```

**Check Docker daemon vs binaire :**
- `docker --version` → vérifie que le binaire est installé
- `docker info` → vérifie que le daemon est actif (retourne non-zero si daemon down)
- Deux checks séparés ou un seul avec état "daemon non actif" comme warning ?
  → **Recommandation :** Un seul check Docker avec sous-états distincts affichés

**Version parsing Go :** Go version format: `go version go1.25.5 darwin/arm64`
```go
// Extract "go1.25.5" → parse major.minor → compare to 1.21
parts := strings.Fields(string(output)) // ["go", "version", "go1.25.5", "darwin/arm64"]
versionStr := strings.TrimPrefix(parts[2], "go") // "1.25.5"
// Parse semver: major, minor, patch
```

**Version constante :** Créer `cmd/create-go-starter/version.go` :
```go
package main

// Version is the current version of create-go-starter CLI
const Version = "1.4.0-dev"
```

**Message `Fix` pour chaque problème :**
- Go manquant : `"Download Go from https://go.dev/dl/ and follow installation instructions"`
- Go trop vieux : `"Update Go to 1.21+ from https://go.dev/dl/"`
- Git manquant : `"Install Git: macOS: brew install git | Ubuntu: sudo apt install git | Windows: https://git-scm.com/download/win"`
- Docker binaire manquant : `"Install Docker Desktop: https://www.docker.com/products/docker-desktop/"`
- Docker daemon inactif : `"Start Docker: macOS/Windows: open Docker Desktop app | Linux: sudo systemctl start docker"`

### Affichage du rapport

```
══════════════════════════════════════════════════════
  create-go-starter doctor — Environment Check
══════════════════════════════════════════════════════
  Tool version: v1.4.0-dev

Checking your development environment...

  ✅ Go         go1.25.5 darwin/arm64
  ✅ Git        git version 2.43.0
  ❌ Docker     Binary not found
     Fix: Install Docker Desktop: https://www.docker.com/products/docker-desktop/

══════════════════════════════════════════════════════
  ⚠️  1 issue found. Please resolve before using create-go-starter.
══════════════════════════════════════════════════════
```

**Si tout est OK :**
```
══════════════════════════════════════════════════════
  ✅ All checks passed! Your environment is ready.
══════════════════════════════════════════════════════
```

**Couleurs :**
- ✅ `Green()` pour les checks OK
- ❌ `Red()` pour les checks FAIL
- Fix message : `Yellow()` pour la solution suggérée
- Résumé OK : `Green()`, résumé avec problèmes : `Yellow()`

### Version Bump

Avec cette story, introduire un mécanisme de versioning. La version `1.4.0-dev` sera l'identifiant de la version en cours de développement pour l'Epic 10. À la fin de l'epic, elle deviendra `1.4.0`.

**Note sur l'ordre :** La story 10.4 (visual-feedback) complète l'epic. Utiliser `1.4.0-dev` pour toutes les stories 10.x et finaliser en `1.4.0` lors de la release.

### Testing Approach

Les tests pour `doctor` sont principalement des **tests d'intégration légers** car ils dépendent des outils système :
```go
// test de l'affichage avec des données mockées
func TestDisplayDoctorReport(t *testing.T) {
    results := []CheckResult{
        {Name: "Go", Status: true, Version: "go1.25.5"},
        {Name: "Git", Status: false, Message: "not found", Fix: "install git"},
    }
    // Capturer stdout et vérifier la sortie
}

// test intégration : Go est toujours disponible dans l'env de test
func TestCheckGoVersionIntegration(t *testing.T) {
    result := checkGoVersion()
    if !result.Status {
        t.Errorf("Go should be available in test environment: %s", result.Message)
    }
}
```

**Skip Docker test si Docker non disponible :**
```go
func TestCheckDockerIntegration(t *testing.T) {
    result := checkDocker()
    // Just verify the function doesn't panic and returns a valid result
    assert(t, result.Name == "Docker", "name should be Docker")
}
```

### Project Structure Notes

- **Nouveau fichier :** `cmd/create-go-starter/doctor.go`
- **Nouveau fichier :** `cmd/create-go-starter/version.go`
- **Nouveau fichier test :** `cmd/create-go-starter/doctor_test.go`
- **Fichier modifié :** `cmd/create-go-starter/main.go` (ajout détection sous-commande `doctor`, mise à jour `usage()`)
- **Pas de nouveaux répertoires** requis
- **Pas de nouvelles dépendances** — `os/exec` est stdlib

### References

- Pattern sous-commande `add-model` : [Source: cmd/create-go-starter/main.go:195-201]
- Fonctions couleur : [Source: cmd/create-go-starter/main.go:78-91]
- `usage()` à mettre à jour : [Source: cmd/create-go-starter/main.go:254-285]
- Epic 10 description : [Source: _bmad-output/planning-artifacts/epics.md#Story-10.3]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

Aucun blocage rencontré. Implémentation directe suivant les Dev Notes de la story.

### Completion Notes List

- Créé `cmd/create-go-starter/version.go` avec constante `Version = "1.4.0-dev"`
- Créé `cmd/create-go-starter/doctor.go` avec: struct `CheckResult`, `runDoctor()`, `checkGoVersion()` (avec parsing semver et validation >= 1.21), `checkGit()`, `checkDocker()` (binaire + daemon), `displayDoctorReport()`, `computeExitCode()`
- Modifié `cmd/create-go-starter/main.go`: ajout détection sous-commande `doctor` + mise à jour `usage()` avec section "Subcommands"
- Créé `cmd/create-go-starter/doctor_test.go` avec 10 tests: intégration Go/Git/Docker, affichage rapport, exit code, version constant
- Tous les 10 tests doctor passent; suite complète sans régression
- Exit code: 0 si tout OK, 1 si problème(s) détecté(s)
- Couleurs: ✅ vert pour OK, ❌ rouge pour FAIL, Fix en jaune
- Version v1.4.0-dev affichée dans le rapport

### File List

- `cmd/create-go-starter/doctor.go` (nouveau)
- `cmd/create-go-starter/version.go` (nouveau)
- `cmd/create-go-starter/doctor_test.go` (nouveau)
- `cmd/create-go-starter/main.go` (modifié)

## Change Log

- **2026-02-18**: Implémentation complète de la story 10.3 — Doctor Command
  - Nouveau fichier `doctor.go`: commande `create-go-starter doctor` avec checks Go, Git, Docker
  - Nouveau fichier `version.go`: constante `Version = "1.4.0-dev"`
  - Nouveau fichier `doctor_test.go`: 10 tests (intégration + unitaires)
  - `main.go` modifié: détection sous-commande `doctor`, `usage()` mis à jour
- **2026-02-18**: Code Review — Corrections appliquées (6 issues HIGH/MEDIUM)
  - [H1/H2] Fix `isGoVersionSufficient()` : gère les versions pre-release (rc, beta, alpha) via `stripNonNumericSuffix()`
  - [H1] Ajout de 16 tests table-driven pour `isGoVersionSufficient()` et 7 tests pour `stripNonNumericSuffix()`
  - [M2/M3] Refactoring `displayDoctorReport()` → délègue à `writeDoctorReport(io.Writer)` pour testabilité thread-safe
  - [M3] Suppression de la capture `os.Stdout` via `os.Pipe` dans les tests, remplacée par `bytes.Buffer`
  - [M4] Ajout du support `create-go-starter doctor --help` avec description des checks et exit codes
  - [M1] File List mise à jour pour refléter les fichiers réellement modifiés
  - Tests : 16 tests doctor passent, suite complète sans régression
