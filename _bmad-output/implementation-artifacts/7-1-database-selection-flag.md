# Story 7.1: Database Selection Flag

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** développeur,
**Je veux** pouvoir spécifier le type de base de données via un flag `--database`,
**Afin de** générer un projet avec la base de données de mon choix.

## Acceptance Criteria

1. **AC1**: Given l'utilisateur exécute `create-go-starter mon-projet --database=mysql`, When le CLI parse les arguments, Then MySQL est sélectionné comme base de données
2. **AC2**: Given l'utilisateur exécute `create-go-starter mon-projet` (sans flag), When le CLI parse les arguments, Then PostgreSQL est utilisé par défaut
3. **AC3**: Given l'utilisateur exécute `create-go-starter mon-projet --database=invalid`, When le CLI parse les arguments, Then une erreur est affichée avec les valeurs valides
4. **AC4**: Given l'utilisateur exécute `create-go-starter --help`, When l'aide s'affiche, Then le flag --database est documenté avec les options disponibles (postgres, mysql, sqlite, mongodb)

## Tasks / Subtasks

- [x] Task 1: Ajouter le flag `--database` dans main.go (AC: 1, 2)
  - [x] 1.1 Définir le flag avec `flag.StringVar` et valeur par défaut "postgres"
  - [x] 1.2 Documenter le flag dans `flag.Usage`
  - [x] 1.3 Parser le flag après `flag.Parse()`

- [x] Task 2: Créer la validation de la base de données (AC: 3)
  - [x] 2.1 Créer une constante ou slice avec les valeurs valides: `postgres`, `mysql`, `sqlite`, `mongodb`
  - [x] 2.2 Créer fonction `validateDatabase(database string) error`
  - [x] 2.3 Retourner erreur descriptive si invalide avec liste des options

- [x] Task 3: Passer le paramètre database au générateur (AC: 1, 2)
  - [x] 3.1 Modifier signature de `run(projectName, template string)` → `run(projectName, template, database string)`
  - [x] 3.2 Modifier signature de `generateProjectFiles(projectPath, projectName, template string)` → ajouter `database string`
  - [x] 3.3 Afficher la base de données sélectionnée dans les messages de progression

- [x] Task 4: Mettre à jour l'aide CLI (AC: 4)
  - [x] 4.1 Ajouter description du flag --database dans `flag.Usage`
  - [x] 4.2 Lister les 4 options avec descriptions courtes

- [x] Task 5: Tests unitaires
  - [x] 5.1 Test `validateDatabase` avec valeurs valides et invalides
  - [x] 5.2 Test flag parsing avec différentes combinaisons
  - [x] 5.3 Test erreur quand database invalide
  - [x] 5.4 Test valeur par défaut (postgres)

## Dev Notes

### 🎯 Objectif Principal

Cette story est la **FONDATION** de l'Epic 7 - Multi-Database Support. Elle établit l'interface CLI qui permettra aux utilisateurs de choisir leur base de données. Comme pour le flag `--template` (Story 6.1), nous ajoutons le flag mais la logique de génération différenciée sera implémentée dans les stories suivantes (7.2, 7.3, 7.4).

### 🏗️ Architecture Actuelle du CLI

#### Structure des flags existants (main.go)

**Ligne 113-128**: Parsing des flags actuels
```go
var template string
flag.StringVar(&template, "template", "full", "Template type: minimal, full, graphql")

flag.Usage = func() {
    // Custom usage avec sections Templates, etc.
}
```

**Ligne 140**: Appel de `run(projectName, template)`
**Ligne 175**: Appel de `generateProjectFiles(projectPath, projectName, template)`

#### Pattern établi dans Story 6.1 (Référence)

La Story 6.1 a établi le pattern suivant pour ajouter un flag de sélection:
1. Définir le flag avec `flag.StringVar` et une valeur par défaut
2. Créer une liste de valeurs valides (`ValidTemplates`)
3. Créer une fonction de validation (`validateTemplate()`)
4. Modifier les signatures des fonctions `run()` et `generateProjectFiles()`
5. Documenter dans `flag.Usage()`
6. Tests complets (valeurs valides, invalides, défaut, help)

**CRITICAL**: Nous devons suivre EXACTEMENT le même pattern pour maintenir la cohérence du code.

### 🗄️ Bases de Données Supportées

| Database | Description | Driver | Use Case |
|----------|-------------|--------|----------|
| `postgres` | PostgreSQL (défaut actuel) | gorm.io/driver/postgres | Production, features avancées |
| `mysql` | MySQL/MariaDB | gorm.io/driver/mysql | Compatibilité, hosting partagé |
| `sqlite` | SQLite | gorm.io/driver/sqlite | Prototypage, apps embarquées |
| `mongodb` | MongoDB (NoSQL) | go.mongodb.org/mongo-driver | NoSQL, documents JSON |

**Valeur par défaut**: `postgres` (continuité avec l'implémentation actuelle)

### 📁 Fichiers à Modifier

#### 1. `cmd/create-go-starter/main.go`

**Modifications requises**:
```go
// Ajouter après le flag template (ligne ~115)
var database string
flag.StringVar(&database, "database", "postgres", "Database type: postgres, mysql, sqlite, mongodb")

// Ajouter constantes (après ValidTemplates)
var ValidDatabases = []string{"postgres", "mysql", "sqlite", "mongodb"}
const DefaultDatabase = "postgres"

// Ajouter fonction de validation (après validateTemplate)
func validateDatabase(database string) error {
    for _, valid := range ValidDatabases {
        if database == valid {
            return nil
        }
    }
    return fmt.Errorf("invalid database '%s'. Valid options: %s", 
        database, strings.Join(ValidDatabases, ", "))
}

// Modifier signature de run (ligne ~140)
func run(projectName, template, database string) error {
    // ...
    fmt.Printf("  Template: %s\n", Green(template))
    fmt.Printf("  Database: %s\n", Green(database))
    // ...
}

// Appel dans main (ligne ~170)
if err := run(projectName, template, database); err != nil {
    // ...
}
```

**Mise à jour de flag.Usage()** (ajouter section Databases):
```go
flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "%s\n\n", Bold("Usage:"))
    // ... existing content ...
    
    fmt.Fprintf(os.Stderr, "\n%s\n", Bold("Databases:"))
    fmt.Fprintf(os.Stderr, "  postgres    PostgreSQL (default) - Production-ready, advanced features\n")
    fmt.Fprintf(os.Stderr, "  mysql       MySQL/MariaDB - Wide compatibility, shared hosting\n")
    fmt.Fprintf(os.Stderr, "  sqlite      SQLite - Quick prototyping, embedded apps\n")
    fmt.Fprintf(os.Stderr, "  mongodb     MongoDB - NoSQL, document-oriented\n")
}
```

#### 2. `cmd/create-go-starter/generator.go`

**Modification de signature**:
```go
// Ligne ~30
func generateProjectFiles(projectPath, projectName, template, database string) error {
    // Pour l'instant, le paramètre database est ignoré
    // Il sera utilisé dans les stories 7.2, 7.3, 7.4 pour la génération différenciée
    
    // ... existing code ...
}
```

#### 3. `cmd/create-go-starter/main_test.go`

**Tests à ajouter** (suivre le pattern de Story 6.1):
```go
func TestValidateDatabaseValid(t *testing.T) {
    validDBs := []string{"postgres", "mysql", "sqlite", "mongodb"}
    for _, db := range validDBs {
        if err := validateDatabase(db); err != nil {
            t.Errorf("validateDatabase(%s) should be valid, got error: %v", db, err)
        }
    }
}

func TestValidateDatabaseInvalid(t *testing.T) {
    invalidDBs := []string{"oracle", "mssql", "", "POSTGRES"}
    for _, db := range invalidDBs {
        if err := validateDatabase(db); err == nil {
            t.Errorf("validateDatabase(%s) should return error", db)
        }
    }
}

func TestDatabaseDefaultValue(t *testing.T) {
    if DefaultDatabase != "postgres" {
        t.Errorf("Expected default database to be 'postgres', got '%s'", DefaultDatabase)
    }
}

func TestValidDatabasesContains(t *testing.T) {
    expected := []string{"postgres", "mysql", "sqlite", "mongodb"}
    if len(ValidDatabases) != len(expected) {
        t.Errorf("Expected %d valid databases, got %d", len(expected), len(ValidDatabases))
    }
}

// Tests similaires aux templates pour flag parsing, help, etc.
```

#### 4. Autres fichiers de tests

**Tous les tests existants** qui appellent `run()` ou `generateProjectFiles()` doivent être mis à jour avec le nouveau paramètre `database`:
- `cmd/create-go-starter/generator_test.go`
- `cmd/create-go-starter/scaffold_test.go`
- `cmd/create-go-starter/git_test.go`
- `cmd/create-go-starter/smoke_test.go`

**Pattern de mise à jour**:
```go
// Avant (Story 6.1)
err := generateProjectFiles(tmpDir, "test-project", "full")

// Après (Story 7.1)
err := generateProjectFiles(tmpDir, "test-project", "full", "postgres")
```

### 🔍 Analyse du Code Existant

#### Localisation des flags actuels

**Fichier**: `cmd/create-go-starter/main.go`

**Flags existants**:
- `--template` (depuis Story 6.1): Ligne ~115
  - Valeurs: minimal, full, graphql
  - Défaut: full
  - Pattern: `flag.StringVar(&template, "template", "full", "...")`

**Pattern de validation**:
- `ValidTemplates []string` (slice de valeurs valides)
- `DefaultTemplate string` (constante pour la valeur par défaut)
- `validateTemplate(template string) error` (fonction de validation)

**Pattern d'affichage**:
```go
fmt.Printf("  Template: %s\n", Green(template))
```

### ⚠️ Contraintes Critiques

#### 1. **Backward Compatibility**
- Le paramètre `database` est ajouté mais **ignoré** dans la génération
- La logique de génération différenciée sera dans les stories 7.2-7.4
- Tous les tests existants doivent continuer à passer

#### 2. **Cohérence avec Story 6.1**
- Utiliser EXACTEMENT le même pattern que pour `--template`
- Même structure de validation
- Même style de documentation dans `--help`
- Même approche de tests

#### 3. **Pas de Logique de Génération**
- Ne PAS modifier les templates existants
- Ne PAS ajouter de drivers de bases de données
- Ne PAS créer de nouveaux fichiers template
- Focus uniquement sur le **parsing et validation du flag**

#### 4. **Standards Go**
- Respecter golangci-lint
- Utiliser les fonctions de couleur existantes (`Green()`, `Red()`, `Bold()`)
- Commentaires pour fonctions exportées
- Gestion d'erreurs explicite

### 🧪 Stratégie de Tests

#### Tests de validation (suivre Story 6.1)
1. `TestValidateDatabaseValid` - toutes les valeurs valides
2. `TestValidateDatabaseInvalid` - valeurs invalides
3. `TestDatabaseDefaultValue` - vérifier "postgres"
4. `TestValidDatabasesContains` - liste complète
5. `TestDatabaseFlagParsing` - avec --database=X
6. `TestDatabaseFlagDefault` - sans flag
7. `TestInvalidDatabaseFlagError` - erreur sur invalide
8. `TestHelpShowsDatabaseFlag` - documentation

#### Tests de régression
- Mettre à jour TOUS les tests existants pour passer le paramètre `database`
- Vérifier que les tests E2E passent toujours
- Pattern: Ajouter `"postgres"` comme dernier paramètre partout

### 🎨 Messages Utilisateur

#### Message de progression (dans run())
```
Creating project: my-project
  Template: full
  Database: postgres  ← NOUVEAU
```

#### Message d'erreur (validation)
```
Error: invalid database 'oracle'. Valid options: postgres, mysql, sqlite, mongodb
```

#### Message d'aide (--help)
```
Databases:
  postgres    PostgreSQL (default) - Production-ready, advanced features
  mysql       MySQL/MariaDB - Wide compatibility, shared hosting
  sqlite      SQLite - Quick prototyping, embedded apps
  mongodb     MongoDB - NoSQL, document-oriented
```

### 📚 Références Techniques

#### Drivers GORM (pour contexte - pas utilisés dans cette story)
- **PostgreSQL**: `gorm.io/driver/postgres`
- **MySQL**: `gorm.io/driver/mysql`
- **SQLite**: `gorm.io/driver/sqlite`
- **MongoDB**: Nécessite `go.mongodb.org/mongo-driver` (différent de GORM)

#### Documentation Go
- `flag` package: https://pkg.go.dev/flag
- GORM Drivers: https://gorm.io/docs/connecting_to_the_database.html

### 🚀 Ordre d'Implémentation Recommandé

1. **Ajouter les constantes** (`ValidDatabases`, `DefaultDatabase`)
2. **Créer `validateDatabase()`** avec tests
3. **Ajouter le flag** dans `main()`
4. **Modifier signatures** de `run()` et `generateProjectFiles()`
5. **Mettre à jour `flag.Usage()`**
6. **Mettre à jour tous les tests existants**
7. **Ajouter nouveaux tests** pour database
8. **Vérifier compilation** et tous les tests

### 🔗 Context des Stories Suivantes

#### Story 7.2 - MySQL/MariaDB Support
Utilisera `database == "mysql"` pour:
- Générer templates avec driver MySQL
- Configurer docker-compose.yml avec MySQL
- Adapter go.mod avec bonnes dépendances

#### Story 7.3 - SQLite Support
Utilisera `database == "sqlite"` pour:
- Driver SQLite (fichier .db local)
- Pas de docker-compose pour DB
- Configuration simplifiée

#### Story 7.4 - MongoDB Support
Utilisera `database == "mongodb"` pour:
- Architecture sans GORM (mongo-driver natif)
- Templates NoSQL adaptés
- Docker avec MongoDB

### Project Structure Notes

- **Alignement avec structure unifiée**: Le flag suit le pattern existant de `--template`
- **Pas de conflits détectés**: L'ajout du flag `--database` est orthogonal au flag `--template`
- **Naming cohérent**: `validateDatabase`, `ValidDatabases`, `DefaultDatabase` (même pattern que Template)

### References

- [Source: cmd/create-go-starter/main.go#L113-128] - Flags existants et pattern
- [Source: _bmad-output/implementation-artifacts/6-1-flag-cli-pour-selection-de-template.md] - Story 6.1 comme référence
- [Source: _bmad-output/planning-artifacts/epics.md#Epic 7, Story 7.1] - Spécifications
- [Source: _bmad-output/planning-artifacts/architecture.md#Data Architecture] - Contexte bases de données
- [Source: cmd/create-go-starter/generator.go#L30] - Fonction generateProjectFiles actuelle

## 🛡️ Developer Context - Critical Implementation Guardrails

### 🔥 Prévention des Erreurs Communes LLM

**ERREUR #1: Casser les tests existants**
- ❌ **NE PAS** oublier de mettre à jour TOUS les appels existants à `run()` et `generateProjectFiles()`
- ✅ **FAIRE**: Rechercher globalement tous les appels et ajouter le paramètre `"postgres"`
- 🔍 **Vérifier**: `grep -r "generateProjectFiles\|run(" cmd/create-go-starter/*_test.go`

**ERREUR #2: Implémenter trop tôt**
- ❌ **NE PAS** commencer à modifier les templates de génération
- ❌ **NE PAS** ajouter les drivers de bases de données dans go.mod
- ✅ **FAIRE**: Uniquement parsing + validation du flag, RIEN d'autre
- 📌 **Rappel**: La logique de génération différenciée = Stories 7.2-7.4

**ERREUR #3: Inconsistance avec le pattern existant**
- ❌ **NE PAS** inventer un nouveau pattern de validation
- ✅ **FAIRE**: Copier EXACTEMENT le pattern de `validateTemplate()` de Story 6.1
- 📋 **Pattern**: `ValidDatabases []string` + `DefaultDatabase const` + `validateDatabase() error`

**ERREUR #4: Oublier la documentation**
- ❌ **NE PAS** oublier de mettre à jour `flag.Usage()`
- ✅ **FAIRE**: Ajouter section "Databases:" avec descriptions pour chaque option
- 🎨 **Style**: Suivre le format exact de la section "Templates:"

**ERREUR #5: Tests incomplets**
- ❌ **NE PAS** se limiter à 1-2 tests
- ✅ **FAIRE**: Minimum 8 tests comme dans Story 6.1
- 📊 **Coverage**: Valides, invalides, défaut, help, parsing, erreurs

### 🎯 Exigences Techniques Critiques

#### Respect des Standards Go

**golangci-lint** (MUST PASS):
```bash
# Avant de commit, TOUJOURS exécuter
make lint
# OU
golangci-lint run ./cmd/create-go-starter/...
```

**Conventions de nommage**:
- Variables/Fonctions: `camelCase` → `validateDatabase`, `database`
- Constantes: `PascalCase` → `DefaultDatabase`, `ValidDatabases`
- Acronymes: Majuscules → `DBConfig` (pas `DbConfig`)

**Documentation des fonctions exportées**:
```go
// validateDatabase checks if the provided database type is supported.
// Returns an error if the database is not in the ValidDatabases list.
func validateDatabase(database string) error {
    // ...
}
```

#### Cohérence avec l'Architecture Hexagonale

**Emplacement du code**:
- ✅ Flags CLI → `cmd/create-go-starter/main.go`
- ✅ Validation → `cmd/create-go-starter/main.go` (fonctions utilitaires)
- ✅ Génération → `cmd/create-go-starter/generator.go`
- ❌ NE PAS toucher à `/internal/` dans cette story

**Principe de séparation**:
- Le CLI (`main.go`) = Interface utilisateur + validation
- Le générateur (`generator.go`) = Logique de génération
- Les templates (`templates*.go`) = Contenu à générer

#### Gestion des Erreurs

**Pattern obligatoire**:
```go
// GOOD: Erreurs descriptives avec contexte
if err := validateDatabase(database); err != nil {
    fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("Error: %v", err)))
    os.Exit(1)
}

// BAD: Ignorer les erreurs
validateDatabase(database)  // ❌ JAMAIS faire ça
```

**Messages d'erreur utilisateur**:
- ✅ Incluent la valeur invalide
- ✅ Listent les options valides
- ✅ Utilisent la fonction `Red()` pour la couleur
- Format: `"invalid database 'X'. Valid options: postgres, mysql, sqlite, mongodb"`

### 📋 Exigences de Tests

#### Tests Unitaires Obligatoires

**Minimum requis** (suivre Story 6.1):
1. ✅ `TestValidateDatabaseValid` - Tous les cas valides
2. ✅ `TestValidateDatabaseInvalid` - Cas invalides + edge cases
3. ✅ `TestDatabaseDefaultValue` - Vérifier constante
4. ✅ `TestValidDatabasesContains` - Vérifier liste complète
5. ✅ `TestDatabaseFlagParsing` - Parsing de différentes valeurs
6. ✅ `TestDatabaseFlagDefault` - Comportement sans flag
7. ✅ `TestInvalidDatabaseFlagError` - Erreur sur valeur invalide
8. ✅ `TestHelpShowsDatabaseFlag` - Documentation dans --help

**Edge cases à tester**:
```go
// Cas limites à couvrir
invalidCases := []string{
    "",              // Chaîne vide
    "POSTGRES",      // Mauvaise casse
    "postgre",       // Typo
    "oracle",        // DB non supportée
    "pg",            // Abréviation
    "mysql5.7",      // Avec version
}
```

#### Tests de Régression

**Fichiers à mettre à jour** (TOUS):
- `cmd/create-go-starter/generator_test.go`
- `cmd/create-go-starter/scaffold_test.go`
- `cmd/create-go-starter/templates_test.go`
- `cmd/create-go-starter/git_test.go`
- `cmd/create-go-starter/smoke_test.go`

**Pattern de mise à jour**:
```go
// Rechercher tous les appels avec 3 paramètres
generateProjectFiles(path, name, template)

// Remplacer par 4 paramètres
generateProjectFiles(path, name, template, "postgres")
```

#### Commandes de Test

```bash
# Tests unitaires uniquement
go test ./cmd/create-go-starter -v

# Tests avec coverage
go test ./cmd/create-go-starter -cover

# Tests courts (skip E2E)
go test -short ./cmd/create-go-starter

# Test spécifique
go test -run TestValidateDatabase ./cmd/create-go-starter
```

### 🏗️ Exigences d'Architecture et Fichiers

#### Structure de Code Requise

**Dans `main.go`** (après les flags existants):
```go
// Database configuration
var database string
var ValidDatabases = []string{"postgres", "mysql", "sqlite", "mongodb"}
const DefaultDatabase = "postgres"

// validateDatabase vérifie que la base de données est supportée
func validateDatabase(database string) error {
    for _, valid := range ValidDatabases {
        if database == valid {
            return nil
        }
    }
    return fmt.Errorf("invalid database '%s'. Valid options: %s", 
        database, strings.Join(ValidDatabases, ", "))
}
```

**Signatures de fonctions modifiées**:
```go
// Avant
func run(projectName, template string) error

// Après
func run(projectName, template, database string) error

// Avant
func generateProjectFiles(projectPath, projectName, template string) error

// Après
func generateProjectFiles(projectPath, projectName, template, database string) error
```

#### Modification de flag.Usage()

**Emplacement**: Dans `main()`, section `flag.Usage = func()`

**Ajouter après la section "Templates:"**:
```go
fmt.Fprintf(os.Stderr, "\n%s\n", Bold("Databases:"))
fmt.Fprintf(os.Stderr, "  postgres    PostgreSQL (default) - Production-ready, advanced features\n")
fmt.Fprintf(os.Stderr, "  mysql       MySQL/MariaDB - Wide compatibility, shared hosting\n")
fmt.Fprintf(os.Stderr, "  sqlite      SQLite - Quick prototyping, embedded apps\n")
fmt.Fprintf(os.Stderr, "  mongodb     MongoDB - NoSQL, document-oriented\n")
```

### 🔍 Checklist de Validation Pré-Commit

**Avant de marquer la story comme complète**:

- [ ] Le flag `--database` est défini avec valeur par défaut "postgres"
- [ ] La fonction `validateDatabase()` existe et fonctionne correctement
- [ ] Les constantes `ValidDatabases` et `DefaultDatabase` sont définies
- [ ] Les signatures de `run()` et `generateProjectFiles()` incluent le paramètre database
- [ ] La section "Databases:" est présente dans `--help`
- [ ] TOUS les tests existants ont été mis à jour (aucun échec)
- [ ] Au minimum 8 nouveaux tests pour database sont ajoutés
- [ ] `make lint` passe sans erreurs ni warnings
- [ ] `go test ./cmd/create-go-starter -v` passe à 100%
- [ ] `go test -short ./...` passe (tests E2E)
- [ ] Le message de progression affiche "Database: postgres"
- [ ] Une erreur descriptive s'affiche pour database invalide
- [ ] La commande `create-go-starter mon-projet --database=mysql` parse correctement
- [ ] La commande `create-go-starter mon-projet` utilise "postgres" par défaut

### 🚨 Ce qu'il ne faut ABSOLUMENT PAS faire

**❌ INTERDICTIONS STRICTES**:

1. **NE PAS modifier les fichiers de templates**
   - `templates.go`, `templates_user.go`, etc. → INTOUCHABLES
   
2. **NE PAS ajouter de dépendances**
   - Pas de `go get gorm.io/driver/mysql` dans cette story
   - Pas de modification de `go.mod` pour drivers
   
3. **NE PAS créer de logique de génération conditionnelle**
   - Pas de `if database == "mysql" { ... }` dans les templates
   - Le paramètre est passé mais ignoré pour l'instant
   
4. **NE PAS toucher au docker-compose ou Dockerfile**
   - Ces modifications sont pour Stories 7.2-7.4
   
5. **NE PAS modifier la structure `/internal/`**
   - Cette story concerne uniquement le CLI

### 📊 Définition de "Done"

**La story est complète quand**:

✅ **Fonctionnel**:
- Le flag `--database` accepte: postgres, mysql, sqlite, mongodb
- Valeur par défaut "postgres" fonctionne sans flag
- Erreur claire et descriptive pour valeurs invalides
- L'aide (`--help`) documente le flag correctement

✅ **Technique**:
- Toutes les signatures de fonctions sont mises à jour
- Tous les tests existants passent
- Au moins 8 nouveaux tests ajoutés et passent
- golangci-lint passe sans erreurs
- Code documenté selon standards Go

✅ **Qualité**:
- Pattern identique à Story 6.1 (--template)
- Aucune régression introduite
- Messages utilisateur clairs et colorés
- Code lisible et maintenable

### 🎓 Apprentissage des Stories Précédentes

#### Leçons de Story 6.1 (Template Flag)

**Ce qui a bien fonctionné**:
- ✅ Pattern clair: constantes → validation → usage
- ✅ Tests exhaustifs (8 tests couvrant tous les cas)
- ✅ Signature modifiée progressivement (run puis generateProjectFiles)
- ✅ Documentation dans --help bien formatée

**Ce qu'il faut répliquer**:
- Exactement le même pattern pour `ValidDatabases`
- Même structure de tests unitaires
- Même approche de mise à jour progressive des signatures
- Même style de messages d'erreur

**Fichiers modifiés dans 6.1** (référence):
- `main.go` - Flag + validation + usage
- `generator.go` - Signature de generateProjectFiles
- `main_test.go` - Nouveaux tests
- `generator_test.go` - Mise à jour appels
- `scaffold_test.go` - Mise à jour appels
- `git_test.go` - Mise à jour appels
- `smoke_test.go` - Mise à jour appels

#### Pattern Git des Commits Récents

**Analyse des 20 derniers commits**:
- Convention: `feat:`, `fix:`, `docs:`, `chore:`
- Messages en anglais, concis
- Focus sur la valeur ajoutée

**Pour cette story**:
```bash
# Bon exemple de message de commit
git commit -m "feat: add --database flag for database selection

- Add ValidDatabases and DefaultDatabase constants
- Implement validateDatabase() function
- Update run() and generateProjectFiles() signatures
- Add 8 unit tests covering all scenarios
- Update flag.Usage() with Databases section"
```

### 🔗 Préparation pour Stories Suivantes

#### Story 7.2 - MySQL/MariaDB Support

**Ce que cette story prépare**:
- Le flag `--database=mysql` sera prêt à être utilisé
- Le paramètre sera passé à `generateProjectFiles()`
- Story 7.2 ajoutera la logique conditionnelle dans generator.go

**Interface établie**:
```go
// Story 7.1 (cette story) passe le paramètre
generateProjectFiles(projectPath, projectName, template, database)

// Story 7.2 utilisera database pour générer différemment
if database == "mysql" {
    // Générer templates MySQL
    // Configurer docker-compose avec MySQL
}
```

#### Story 7.3 - SQLite Support

**Différence avec MySQL**:
- Pas de docker-compose (fichier .db local)
- Driver différent (gorm.io/driver/sqlite)
- Configuration simplifiée

#### Story 7.4 - MongoDB Support

**Changement architectural**:
- Pas de GORM (mongo-driver natif)
- Templates NoSQL différents
- Architecture adaptée

### 💡 Conseils pour le Dev Agent

**Approche recommandée**:

1. **Commencer petit**: Ajouter constantes et validateDatabase()
2. **Tester tôt**: Écrire les tests de validation immédiatement
3. **Modifier progressivement**: D'abord run(), puis generateProjectFiles()
4. **Vérifier continuellement**: `go test` après chaque modification
5. **Mettre à jour systématiquement**: Chercher tous les appels existants
6. **Documenter proprement**: Comments + flag.Usage()
7. **Valider avant commit**: make lint + tous les tests

**Ordre d'implémentation détaillé**:

1. Ajouter `ValidDatabases` et `DefaultDatabase` dans main.go
2. Créer `validateDatabase()` avec tests
3. Ajouter le flag avec `flag.StringVar`
4. Modifier signature de `run()` + appels
5. Modifier signature de `generateProjectFiles()` + appels
6. Mettre à jour `flag.Usage()`
7. Mettre à jour tous les fichiers de tests
8. Ajouter les 8 tests pour database
9. Exécuter `make lint` et corriger
10. Exécuter tous les tests et vérifier

**Si les tests échouent**:
- Vérifier que TOUS les appels à `run()` ont 3 paramètres
- Vérifier que TOUS les appels à `generateProjectFiles()` ont 4 paramètres
- Chercher avec `grep` pour s'assurer de ne rien manquer

### 📖 Ressources et Documentation

**Documentation Go**:
- Flag package: https://pkg.go.dev/flag
- Testing: https://pkg.go.dev/testing
- Strings: https://pkg.go.dev/strings

**Standards du projet**:
- [AGENTS.md](../../AGENTS.md) - Guidelines pour agents
- [README.md](../../README.md) - Instructions de build/test
- [Project Context](_bmad-output/project-context.md) - Règles critiques

**Fichiers de référence**:
- Story 6.1: `_bmad-output/implementation-artifacts/6-1-flag-cli-pour-selection-de-template.md`
- Architecture: `_bmad-output/planning-artifacts/architecture.md`
- PRD: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

claude-sonnet-4.5 (via OpenCode)

### Debug Log References

N/A - Implémentation directe sans blocage

### Completion Notes List

✅ **Task 1 - Flag CLI ajouté** (2026-02-05)
- Ajout du flag `--database` dans main.go:139
- Valeur par défaut: "postgres"
- Constantes définies: DatabasePostgres, DatabaseMySQL, DatabaseSQLite, DatabaseMongoDB
- ValidDatabases slice créée avec les 4 options

✅ **Task 2 - Validation créée** (2026-02-05)
- Fonction `validateDatabase()` implémentée dans main.go:75-83
- Validation contre ValidDatabases slice
- Messages d'erreur descriptifs avec liste des options valides
- Pattern identique à validateTemplate()

✅ **Task 3 - Paramètre passé au générateur** (2026-02-05)
- Signature `run()` modifiée: ajout paramètre `database string` (main.go:224)
- Signature `generateProjectFiles()` modifiée: ajout paramètre `database string` (generator.go:81)
- Message de progression mis à jour pour afficher le database sélectionné
- Note: Le paramètre est accepté mais pas encore utilisé dans la génération (comme prévu pour Story 7.1)

✅ **Task 4 - Aide CLI mise à jour** (2026-02-05)
- Section "Databases:" ajoutée dans flag.Usage() (main.go:145-148)
- 4 options documentées avec descriptions:
  - postgres: PostgreSQL (default) - Production-ready, advanced features
  - mysql: MySQL/MariaDB - Wide compatibility, shared hosting
  - sqlite: SQLite - Quick prototyping, embedded apps
  - mongodb: MongoDB - NoSQL, document-oriented

✅ **Task 5 - Tests unitaires complets** (2026-02-05)
- 8 nouveaux tests ajoutés dans main_test.go:
  1. TestValidateDatabaseValid - 4 cas valides (postgres, mysql, sqlite, mongodb)
  2. TestValidateDatabaseInvalid - 6 cas invalides (vide, oracle, POSTGRES, etc.)
  3. TestDatabaseDefaultValue - Vérification constante "postgres"
  4. TestValidDatabasesContains - Vérification liste complète
  5. TestDatabaseFlagParsing - Parsing avec 4 combinaisons différentes
  6. TestDatabaseFlagDefault - Comportement sans flag
  7. TestInvalidDatabaseFlagError - Erreur sur valeur invalide
  8. TestHelpShowsDatabaseFlag - Documentation dans --help

✅ **Tests de régression passés** (2026-02-05)
- Tous les tests existants mis à jour avec nouveau paramètre `database`
- Fichiers modifiés:
  - main_test.go: 4 appels à run() mis à jour
  - generator_test.go: 6 appels à generateProjectFiles() mis à jour
  - template_minimal_test.go: 3 appels à generateProjectFiles() mis à jour
  - smoke_test.go: 1 appel à generateProjectFiles() mis à jour
- `go test -short ./cmd/create-go-starter` → PASS (1.725s)
- `go test ./cmd/create-go-starter` → PASS sauf Docker E2E (attendu)

✅ **Qualité code validée** (2026-02-05)
- `go fmt ./cmd/create-go-starter/...` → OK
- `go vet ./cmd/create-go-starter/...` → OK
- Pas d'erreurs de compilation
- Pattern cohérent avec Story 6.1 (--template)

### File List

**Fichiers modifiés** (relatifs au repo root):
- `cmd/create-go-starter/main.go` - Flag database, validation, signatures modifiées
- `cmd/create-go-starter/generator.go` - Signature generateProjectFiles() modifiée
- `cmd/create-go-starter/main_test.go` - 8 nouveaux tests + 4 appels mis à jour
- `cmd/create-go-starter/generator_test.go` - 6 appels mis à jour
- `cmd/create-go-starter/template_minimal_test.go` - 3 appels mis à jour
- `cmd/create-go-starter/smoke_test.go` - 1 appel mis à jour

**Total**: 6 fichiers modifiés, ~200 lignes de code ajoutées (code + tests)

---

## 📊 Story Summary

### Résumé Exécutif

Cette story établit les **fondations de l'Epic 7** en ajoutant le flag `--database` au CLI, permettant aux utilisateurs de spécifier leur base de données préférée lors de la génération d'un projet. Elle suit exactement le même pattern que la Story 6.1 (--template) pour maintenir la cohérence du code.

### Valeur Ajoutée

**Pour l'utilisateur**:
- Interface claire pour choisir la base de données
- Documentation complète via `--help`
- Messages d'erreur descriptifs en cas de valeur invalide
- Valeur par défaut intelligente (postgres)

**Pour le projet**:
- Préparation pour le support multi-database (Epic 7)
- Pattern cohérent avec le flag --template existant
- Tests exhaustifs garantissant la qualité
- Base solide pour les stories 7.2-7.4

### Scope Précis

**✅ Inclus dans cette story**:
- Flag CLI `--database` avec parsing et validation
- Support pour: postgres, mysql, sqlite, mongodb
- Valeur par défaut: postgres
- Documentation dans `--help`
- 8+ tests unitaires complets
- Mise à jour de tous les tests existants

**❌ Explicitement exclus** (Stories suivantes):
- Logique de génération différenciée par database
- Ajout de drivers de bases de données
- Modification des templates de génération
- Configuration docker-compose par database
- Tests d'intégration avec différentes DB

### Impact sur le Code

**Fichiers modifiés** (~8 fichiers):
1. `cmd/create-go-starter/main.go` - Flag + validation + usage
2. `cmd/create-go-starter/generator.go` - Signature
3. `cmd/create-go-starter/main_test.go` - Nouveaux tests
4. `cmd/create-go-starter/generator_test.go` - Mise à jour
5. `cmd/create-go-starter/scaffold_test.go` - Mise à jour
6. `cmd/create-go-starter/templates_test.go` - Mise à jour
7. `cmd/create-go-starter/git_test.go` - Mise à jour
8. `cmd/create-go-starter/smoke_test.go` - Mise à jour

**Lignes de code estimées**: ~150-200 lignes (code + tests)

### Risques et Mitigation

**Risque 1**: Casser les tests existants
- **Mitigation**: Recherche globale de tous les appels + tests de régression
- **Probabilité**: Faible (pattern bien établi)

**Risque 2**: Inconsistance avec pattern existant
- **Mitigation**: Suivre strictement Story 6.1 comme référence
- **Probabilité**: Très faible (documentation claire)

**Risque 3**: Oublier de documenter
- **Mitigation**: Checklist de validation pré-commit
- **Probabilité**: Faible (requirement explicite)

### Prochaines Étapes

**Après cette story**:
1. **Story 7.2**: Implémenter le support MySQL/MariaDB
2. **Story 7.3**: Implémenter le support SQLite
3. **Story 7.4**: Implémenter le support MongoDB (optionnel)
4. **Story 7.5**: Tests E2E et documentation finale

**Dépendances**:
- Aucune story ne dépend de celle-ci pour démarrer
- Mais 7.2-7.5 **utilisent** le flag ajouté ici

---

## 🎯 Success Criteria Checklist

### Acceptance Criteria (AC)

- [ ] **AC1**: `--database=mysql` sélectionne MySQL correctement
- [ ] **AC2**: Sans flag, postgres est utilisé par défaut
- [ ] **AC3**: `--database=invalid` affiche une erreur avec les valeurs valides
- [ ] **AC4**: `--help` documente le flag --database avec les 4 options

### Quality Gates

- [ ] Tous les tests existants passent (0 régression)
- [ ] Au moins 8 nouveaux tests ajoutés et passent
- [ ] `make lint` passe sans erreurs ni warnings
- [ ] Documentation claire dans le code (comments)
- [ ] Messages utilisateur colorés et descriptifs
- [ ] Pattern cohérent avec Story 6.1 (--template)

### Definition of Ready for Next Story (7.2)

- [ ] Le flag `--database` est fonctionnel et testé
- [ ] Le paramètre `database` est passé à `generateProjectFiles()`
- [ ] La documentation indique que postgres/mysql/sqlite/mongodb sont supportés
- [ ] Aucune régression introduite dans le code existant

---

**Date de création**: 2026-02-05  
**Epic**: 7 - Multi-Database Support (v1.1.0)  
**Priorité**: Haute  
**Estimation**: 2-3 heures pour un développeur expérimenté  
**Complexité**: Faible (pattern établi à répliquer)
