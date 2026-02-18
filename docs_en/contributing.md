# Contributing Guide

Thank you for contributing to `create-go-starter`! This guide will help you get started.

## Code of Conduct

By participating in this project, you agree to:

- Be respectful and inclusive
- Accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other contributors

## How to Contribute

There are several ways to contribute:

1. **Report bugs** - Help us improve by reporting issues
2. **Propose features** - Share your ideas to improve the tool
3. **Improve documentation** - Fix or add documentation
4. **Contribute code** - Solve issues or implement new features
5. **Test** - Try the tool and give your feedback

## Reporting a Bug

If you find a bug, create a GitHub issue with:

### Bug Report Template

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

**Where to report**: [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)

## Proposing a Feature

Before proposing a feature, check that it hasn't already been proposed.

### Feature Request Template

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

**Where to propose**: [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)

## Contributing Code

### Prerequisites

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **Git** - For cloning and contributing
- **golangci-lint** (optional) - For linting
- **Make** - For build commands

### Development Setup

#### 1. Fork and clone

```bash
# Fork le repository sur GitHub (bouton Fork)

# Clone votre fork
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit

# Ajouter le repository upstream
git remote add upstream https://github.com/tky0065/go-starter-kit.git
```

#### 2. Install development dependencies

```bash
# Installer golangci-lint (optionnel)
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Ou avec go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### 3. Build and test

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

### Contribution Workflow

#### 1. Create a branch

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

**Branch naming conventions**:
- `feature/nom-de-la-feature` - New feature
- `fix/description-du-fix` - Bug fix
- `docs/sujet-de-la-doc` - Documentation
- `refactor/description` - Refactoring
- `test/description` - Adding/improving tests

#### 2. Make your changes

```bash
# Éditez les fichiers
# Suivez les standards de code (voir section Standards)

# Ajoutez des tests pour votre code
# internal/adapters/handlers/handler_test.go
```

**Best practices**:
- **One feature/fix per branch**: Don't mix multiple changes
- **Tests required**: Add tests for all new code
- **Documentation**: Update the docs if necessary
- **Atomic commits**: Small, focused commits

#### 3. Test your changes

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

**Test checklist**:
- [ ] All tests pass (`make test`)
- [ ] Lint passes (`make lint`)
- [ ] Coverage > 80% for new code
- [ ] Manual tests performed (project generation)
- [ ] Generated project compiles (`go build ./...`)
- [ ] Generated project works (`make run`)

#### 4. Commit your changes

Use **Conventional Commits**:

```bash
git add .
git commit -m "feat: ajouter support pour MySQL"
# ou
git commit -m "fix: corriger validation du nom de projet"
# ou
git commit -m "docs: améliorer guide d'installation"
```

**Commit format**:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting (no code changes)
- `refactor`: Refactoring (neither feature nor fix)
- `test`: Adding/modifying tests
- `chore`: Maintenance (build, deps, etc.)
- `perf`: Performance improvement

**Scope** (optional):
- `cli`: CLI interface (main.go)
- `generator`: File generation (generator.go)
- `templates`: Templates (templates.go)
- `tests`: Tests
- `docs`: Documentation

**Examples**:

```bash
feat(templates): add MySQL database template
fix(cli): validate project name with correct regex
docs(readme): update installation instructions
refactor(generator): simplify file creation logic
test(templates): add tests for ReadmeTemplate
chore(deps): update golangci-lint to v1.55
```

**Message body** (optional but recommended):

```bash
git commit -m "feat(templates): add MySQL support

- Add MySQLDatabaseTemplate method
- Update go.mod template to include mysql driver
- Add flag --database to choose DB type

Closes #42"
```

#### 5. Push to your fork

```bash
git push origin feature/ma-fonctionnalite
```

#### 6. Create a Pull Request

1. Go to GitHub: https://github.com/tky0065/go-starter-kit
2. Click on "Compare & pull request"
3. Fill in the PR template (see below)
4. Submit the PR

### Pull Request Template

When you create a PR, use this template:

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

## Code Standards

To maintain code quality and consistency:

### 1. Formatting

**Always** format with `gofmt`:

```bash
go fmt ./...
```

Or configure your IDE/editor to format on save.

### 2. Linting

Use `golangci-lint` to check quality:

```bash
make lint
# ou
golangci-lint run ./...
```

**Configuration**: `.golangci.yml` at the root of the project.

**Main rules**:
- errcheck, gosimple, govet, ineffassign, staticcheck
- gofmt, goimports
- misspell, revive
- No unused vars, unused params

### 3. Naming conventions

Follow Go conventions:

**Variables and functions**:
- `camelCase` for private
- `PascalCase` for public
- Descriptive names (no abbreviations except idiomatic ones)

**Examples**:

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

**Constants**:

```go
const (
    DefaultPort     = 8080
    MaxNameLength   = 100
    ValidationRegex = `^[a-zA-Z0-9][a-zA-Z0-9_-]*$`
)
```

### 4. Error handling

**Always** handle errors explicitly:

```go
// <i class="material-icons success">check_circle</i> Bon
if err != nil {
    return fmt.Errorf("failed to create file %s: %w", path, err)
}

// ❌ Mauvais
os.WriteFile(path, content, 0644)  // Ignore error
```

**Wrapping errors** with context:

```go
if err != nil {
    return fmt.Errorf("generateProjectFiles: %w", err)
}
```

### 5. Documentation

Document **all** public exports with GoDoc:

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
- Complete sentence starting with the name
- No period at the end if a single sentence
- Period at the end if multiple sentences

### 6. Tests

**Minimum coverage**: 80%

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

2. **Descriptive test names**:

```go
func TestGenerateProjectFiles_WithValidName_CreatesAllFiles(t *testing.T)
func TestGenerateProjectFiles_WithInvalidName_ReturnsError(t *testing.T)
```

3. **Setup and teardown**:

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

Group imports by category:

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

## Review Process

### When you submit a PR

1. **Self-review**: Review your code before submitting
2. **CI checks**: Wait for GitHub Actions to pass
3. **Reviews**: At least 1 review required
4. **Changes**: Incorporate the feedback
5. **Merge**: Once approved, a maintainer will merge

### When you review a PR

**Review checklist**:

- [ ] The code follows the project standards
- [ ] Tests are present and pass
- [ ] Documentation is up to date
- [ ] No duplicated code
- [ ] No obvious bugs
- [ ] Acceptable performance
- [ ] No secrets/credentials in the code

**How to give feedback**:
- Be constructive and respectful
- Explain the "why", not just the "what"
- Suggest alternatives
- Use GitHub code suggestions

## Roadmap

**Completed features** (v1.0 -- v1.4):

- [x] **Multiple templates** - Three templates (minimal, full, graphql)
- [x] **Multi-database** - PostgreSQL, MySQL, SQLite
- [x] **CRUD Scaffolding** - `add-model` command with BelongsTo/HasMany relations
- [x] **Observability** - Prometheus, Jaeger, Grafana, K8s Health Checks
- [x] **Interactive CLI** - Guided mode (`--interactive`), dry-run, doctor, short aliases

**Planned features** (v1.5+):

- [ ] Choice of web framework (Gin, Echo, Chi)
- [ ] Microservices generation
- [ ] E2E test templates
- [ ] Plugin system for custom templates

**How to contribute to the roadmap**:
- Vote for features in Discussions
- Propose new ideas
- Implement a feature from the roadmap

## Questions?

If you have questions:

- **GitHub Discussions**: For general questions, ideas
- **GitHub Issues**: For specific bugs
- **Pull Requests**: To propose code

## Acknowledgments

Thank you to all contributors who help improve `create-go-starter`!

Every contribution, large or small, makes a difference. <i class="material-icons small" style="color:#e25555">favorite</i>

---

**Happy coding and thank you for contributing!** <i class="material-icons">rocket_launch</i>
