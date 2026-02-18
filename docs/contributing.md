# Guide de contribution

Merci de contribuer à `create-go-starter`! Ce guide vous aidera à démarrer.

## Code de conduite

En participant à ce projet, vous acceptez de:

- Être respectueux et inclusif
- Accepter les critiques constructives
- Vous concentrer sur ce qui est le mieux pour la communauté
- Faire preuve d'empathie envers les autres contributeurs

## Comment contribuer

Il y a plusieurs façons de contribuer:

1. **Signaler des bugs** - Aidez-nous à améliorer en signalant les problèmes
2. **Proposer des fonctionnalités** - Partagez vos idées pour améliorer l'outil
3. **Améliorer la documentation** - Corrigez ou ajoutez de la documentation
4. **Contribuer du code** - Résolvez des issues ou implémentez de nouvelles fonctionnalités
5. **Tester** - Essayez l'outil et donnez votre feedback

## Signaler un bug

Si vous trouvez un bug, créez une issue GitHub avec:

### Template de bug report

```markdown
**Description du bug**
Description claire et concise du bug.

**Steps pour reproduire**
1. Étape 1
2. Étape 2
3. Voir l'erreur

**Comportement attendu**
Ce qui devrait se passer normalement.

**Comportement actuel**
Ce qui se passe réellement.

**Environment**
- OS: [ex: macOS 13.0, Ubuntu 22.04, Windows 11]
- Version Go: [ex: 1.25.0]
- Version create-go-starter: [ex: v1.0.0]
- Commande exécutée: [ex: `create-go-starter my-project`]

**Logs/Screenshots**
Ajoutez des logs ou screenshots si pertinent.

**Contexte additionnel**
Tout autre contexte utile.
```

**Où signaler**: [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)

## Proposer une fonctionnalité

Avant de proposer une feature, vérifiez qu'elle n'a pas déjà été proposée.

### Template de feature request

```markdown
**Problème à résoudre**
Quel problème cette fonctionnalité résout-elle?

**Solution proposée**
Comment pensez-vous que cela devrait fonctionner?

**Alternatives considérées**
Avez-vous pensé à d'autres approches?

**Use cases**
Exemples d'utilisation de cette fonctionnalité.

**Impact**
Qui bénéficierait de cette fonctionnalité?
```

**Où proposer**: [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)

## Contribuer du code

### Prérequis

- **Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **Git** - Pour cloner et contribuer
- **golangci-lint** (optionnel) - Pour le linting
- **Make** - Pour les commandes de build

### Setup de développement

#### 1. Fork et clone

```bash
# Fork le repository sur GitHub (bouton Fork)

# Clone votre fork
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit

# Ajouter le repository upstream
git remote add upstream https://github.com/tky0065/go-starter-kit.git
```

#### 2. Installer les dépendances de développement

```bash
# Installer golangci-lint (optionnel)
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Ou avec go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### 3. Build et test

```bash
# Build le CLI
make build
# ou
go build -o create-go-starter ./cmd/create-go-starter

# Tests
make test

# Lint
make lint
```

### Workflow de contribution

#### 1. Créer une branche

```bash
# Sync avec upstream
git fetch upstream
git checkout main
git merge upstream/main

# Créer une branche pour votre feature/fix
git checkout -b feature/ma-fonctionnalite
# ou
git checkout -b fix/mon-bug-fix
```

**Conventions de nommage de branches**:
- `feature/nom-de-la-feature` - Nouvelle fonctionnalité
- `fix/description-du-fix` - Bug fix
- `docs/sujet-de-la-doc` - Documentation
- `refactor/description` - Refactoring
- `test/description` - Ajout/amélioration de tests

#### 2. Faire vos changements

```bash
# Éditez les fichiers
# Suivez les standards de code (voir section Standards)

# Ajoutez des tests pour votre code
# internal/adapters/handlers/handler_test.go
```

**Best practices**:
- **Une feature/fix par branche**: Ne mélangez pas plusieurs changements
- **Tests obligatoires**: Ajoutez des tests pour tout nouveau code
- **Documentation**: Mettez à jour la doc si nécessaire
- **Commits atomiques**: Commits petits et focalisés

#### 3. Tester vos changements

```bash
# Run tests
make test

# Vérifier coverage
make test-coverage

# Lint
make lint

# Test manuel
./create-go-starter test-project
cd test-project
go mod tidy
make run
```

**Checklist de tests**:
- [ ] Tous les tests passent (`make test`)
- [ ] Lint passe (`make lint`)
- [ ] Coverage > 80% pour le nouveau code
- [ ] Tests manuels effectués (génération de projet)
- [ ] Projet généré compile (`go build ./...`)
- [ ] Projet généré fonctionne (`make run`)

#### 4. Commit vos changements

Utilisez **Conventional Commits**:

```bash
git add .
git commit -m "feat: ajouter support pour MySQL"
# ou
git commit -m "fix: corriger validation du nom de projet"
# ou
git commit -m "docs: améliorer guide d'installation"
```

**Format des commits**:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types**:
- `feat`: Nouvelle fonctionnalité
- `fix`: Bug fix
- `docs`: Documentation seulement
- `style`: Formatage (pas de changement de code)
- `refactor`: Refactoring (ni feature ni fix)
- `test`: Ajout/modification de tests
- `chore`: Maintenance (build, deps, etc.)
- `perf`: Amélioration de performance

**Scope** (optionnel):
- `cli`: CLI interface (main.go)
- `generator`: File generation (generator.go)
- `templates`: Templates (templates.go)
- `tests`: Tests
- `docs`: Documentation

**Exemples**:

```bash
feat(templates): add MySQL database template
fix(cli): validate project name with correct regex
docs(readme): update installation instructions
refactor(generator): simplify file creation logic
test(templates): add tests for ReadmeTemplate
chore(deps): update golangci-lint to v1.55
```

**Message body** (optionnel mais recommandé):

```bash
git commit -m "feat(templates): add MySQL support

- Add MySQLDatabaseTemplate method
- Update go.mod template to include mysql driver
- Add flag --database to choose DB type

Closes #42"
```

#### 5. Push vers votre fork

```bash
git push origin feature/ma-fonctionnalite
```

#### 6. Créer une Pull Request

1. Allez sur GitHub: https://github.com/tky0065/go-starter-kit
2. Cliquez sur "Compare & pull request"
3. Remplissez le template de PR (voir ci-dessous)
4. Soumettez la PR

### Template de Pull Request

Quand vous créez une PR, utilisez ce template:

```markdown
## Description

Décrivez vos changements en quelques lignes.

## Type de changement

- [ ] Bug fix (changement qui corrige un issue)
- [ ] Nouvelle fonctionnalité (changement qui ajoute une fonctionnalité)
- [ ] Breaking change (fix ou feature qui changerait le comportement existant)
- [ ] Documentation seulement

## Motivation et contexte

Pourquoi ce changement est-il nécessaire? Quel problème résout-il?

Fixes #(issue)

## Comment a été testé?

Décrivez les tests que vous avez effectués.

- [ ] Tests unitaires ajoutés/mis à jour
- [ ] Tests manuels effectués
- [ ] Projet généré compile et fonctionne

## Checklist

- [ ] Mon code suit les standards du projet
- [ ] J'ai effectué une auto-review de mon code
- [ ] J'ai commenté mon code dans les zones complexes
- [ ] J'ai mis à jour la documentation si nécessaire
- [ ] Mes changements ne génèrent pas de warnings
- [ ] J'ai ajouté des tests qui prouvent que mon fix/feature fonctionne
- [ ] Tous les tests passent (nouveau et existants)
- [ ] Lint passe sans erreurs

## Screenshots (si applicable)

Ajoutez des screenshots pour aider à expliquer vos changements.
```

## Standards de code

Pour maintenir la qualité et la cohérence du code:

### 1. Formatage

**Toujours** formater avec `gofmt`:

```bash
go fmt ./...
```

Ou configurez votre IDE/éditeur pour formater à la sauvegarde.

### 2. Linting

Utilisez `golangci-lint` pour vérifier la qualité:

```bash
make lint
# ou
golangci-lint run ./...
```

**Configuration**: `.golangci.yml` à la racine du projet.

**Règles principales**:
- errcheck, gosimple, govet, ineffassign, staticcheck
- gofmt, goimports
- misspell, revive
- Pas d'unused vars, unused params

### 3. Naming conventions

Suivez les conventions Go:

**Variables et fonctions**:
- `camelCase` pour privé
- `PascalCase` pour public
- Noms descriptifs (pas d'abréviations sauf idiomatiques)

**Exemples**:

```go
// <i class="material-icons success">check_circle</i> Bon
projectName string
validateProjectName() error
GenerateFiles() error

// ❌ Mauvais
pn string
valProjName() error
genFiles() error
```

**Constantes**:

```go
const (
    DefaultPort     = 8080
    MaxNameLength   = 100
    ValidationRegex = `^[a-zA-Z0-9][a-zA-Z0-9_-]*$`
)
```

### 4. Error handling

**Toujours** gérer les erreurs explicitement:

```go
// <i class="material-icons success">check_circle</i> Bon
if err != nil {
    return fmt.Errorf("failed to create file %s: %w", path, err)
}

// ❌ Mauvais
os.WriteFile(path, content, 0644)  // Ignore error
```

**Wrapping errors** avec contexte:

```go
if err != nil {
    return fmt.Errorf("generateProjectFiles: %w", err)
}
```

### 5. Documentation

Documenter **tous** les exports publics avec GoDoc:

```go
// ProjectTemplates holds all template generation methods for creating
// project files. It uses the project name to inject dynamic content
// into the generated templates.
type ProjectTemplates struct {
    projectName string
}

// NewProjectTemplates creates a new ProjectTemplates instance with
// the given project name. The project name will be used throughout
// all template generation.
func NewProjectTemplates(projectName string) *ProjectTemplates {
    return &ProjectTemplates{projectName: projectName}
}
```

**Format**:
- Phrase complète commençant par le nom
- Pas de point à la fin si une seule phrase
- Point à la fin si plusieurs phrases

### 6. Tests

**Coverage minimum**: 80%

```bash
go test -cover ./...
```

**Patterns**:

1. **Table-driven tests**:

```go
func TestValidateProjectName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid simple", "project", false},
        {"valid with hyphen", "my-project", false},
        {"invalid space", "my project", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateProjectName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, want error = %v", err, tt.wantErr)
            }
        })
    }
}
```

2. **Noms de tests descriptifs**:

```go
func TestGenerateProjectFiles_WithValidName_CreatesAllFiles(t *testing.T)
func TestGenerateProjectFiles_WithInvalidName_ReturnsError(t *testing.T)
```

3. **Setup et teardown**:

```go
func TestSomething(t *testing.T) {
    // Setup
    tmpDir := t.TempDir()  // Auto-cleaned up

    // Test
    // ...

    // Teardown automatique avec t.TempDir()
}
```

### 7. Imports

Grouper les imports par catégorie:

```go
import (
    // Standard library
    "fmt"
    "os"
    "path/filepath"

    // Third-party (si on en avait)
    "github.com/some/package"

    // Internal (si applicable)
    "go-starter-kit/internal/something"
)
```

## Processus de review

### Quand vous soumettez une PR

1. **Auto-review**: Relisez votre code avant de soumettre
2. **CI checks**: Attendez que GitHub Actions passe
3. **Reviews**: Au moins 1 review requis
4. **Changements**: Intégrez le feedback
5. **Merge**: Une fois approuvée, un maintainer mergera

### Quand vous reviewez une PR

**Checklist de review**:

- [ ] Le code suit les standards du projet
- [ ] Les tests sont présents et passent
- [ ] La documentation est à jour
- [ ] Pas de code dupliqué
- [ ] Pas de bug évidents
- [ ] Performance acceptable
- [ ] Pas de secrets/credentials dans le code

**Comment donner du feedback**:
- Soyez constructif et respectueux
- Expliquez le "pourquoi", pas juste le "quoi"
- Proposez des alternatives
- Utilisez des suggestions de code GitHub

## Roadmap

**Fonctionnalités complétées** (v1.0 -- v1.4):

- [x] **Templates multiples** - Trois templates (minimal, full, graphql)
- [x] **Multi-database** - PostgreSQL, MySQL, SQLite
- [x] **CRUD Scaffolding** - Commande `add-model` avec relations BelongsTo/HasMany
- [x] **Observabilité** - Prometheus, Jaeger, Grafana, Health Checks K8s
- [x] **CLI interactif** - Mode guidé (`--interactive`), dry-run, doctor, alias courts

**Fonctionnalités prévues** (v1.5+):

- [ ] Choix du framework web (Gin, Echo, Chi)
- [ ] Génération de microservices
- [ ] Templates de tests E2E
- [ ] Plugin system pour templates custom

**Comment contribuer au roadmap**:
- Votez pour les features dans Discussions
- Proposez de nouvelles idées
- Implémentez une feature du roadmap

## Questions?

Si vous avez des questions:

- **GitHub Discussions**: Pour questions générales, idées
- **GitHub Issues**: Pour bugs spécifiques
- **Pull Requests**: Pour proposer du code

## Remerciements

Merci à tous les contributeurs qui aident à améliorer `create-go-starter`!

Chaque contribution, grande ou petite, fait une différence. <i class="material-icons small" style="color:#e25555">favorite</i>

---

**Bon coding et merci de contribuer!** <i class="material-icons">rocket_launch</i>
