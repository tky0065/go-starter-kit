# Story 4.5: Intégration Continue (GitHub Actions)

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want que mes tests et mon linter soient exécutés automatiquement sur GitHub à chaque push,
so that je puisse éviter les régressions et maintenir la qualité du code lors du travail en équipe.

## Acceptance Criteria

1.  **Workflow Déclenché au Push/PR**
    *   **Given** Je pousse du code sur la branche `main` OU j'ouvre une Pull Request vers `main`.
    *   **When** L'événement Git se produit.
    *   **Then** Un workflow GitHub Actions nommé "CI" démarre automatiquement.

2.  **Linting Automatisé**
    *   **Given** Le workflow a démarré.
    *   **When** L'étape de linting s'exécute.
    *   **Then** Il utilise `golangci-lint` (via l'action officielle ou `make lint`).
    *   **And** Si le linter trouve des erreurs, le workflow échoue (rouge).

3.  **Tests Automatisés**
    *   **Given** L'étape de linting a réussi.
    *   **When** L'étape de test s'exécute.
    *   **Then** Il lance tous les tests unitaires et d'intégration (via `make test`).
    *   **And** Si un test échoue, le workflow échoue.

4.  **Build Check (Optionnel mais recommandé)**
    *   **When** Les tests passent.
    *   **Then** Le workflow tente de compiler le projet (`go build` ou `docker build`) pour vérifier que l'artefact est constructible.

## Tasks / Subtasks

- [x] **Workflow Setup (`manual-test-project`)**
    - [x] Créer le répertoire `.github/workflows/`.
    - [x] Créer le fichier `.github/workflows/ci.yml`.
    - [x] Configurer les triggers: `push` sur main, `pull_request` sur main.

- [x] **Job Definition**
    - [x] **Job `lint`:**
        -   Utiliser `actions/checkout@v4`.
        -   Utiliser `actions/setup-go@v5` (Go 1.25.x).
        -   Utiliser `golangci/golangci-lint-action@v6` (ou exécuter `make lint`).
    - [x] **Job `test`:**
        -   Utiliser `actions/checkout@v4`.
        -   Utiliser `actions/setup-go@v5`.
        -   Installer les dépendances (si non géré par `go test`).
        -   Lancer `make test`.
    - [x] **Job `build` (Optional):**
        -   Lancer `go build ./...` ou `docker build .`.

- [x] **CLI Generator Update**
    - [x] Mettre à jour `templates.go` pour inclure `.github/workflows/ci.yml` dans les projets générés.
    - [x] S'assurer que le contenu du fichier YAML est correct (indentation, variables).

## Dev Notes

### GitHub Actions Workflow Example (`.github/workflows/ci.yml`)

```yaml
name: CI

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:
  quality:
    name: Quality & Security
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          cache: false # golangci-lint-action handles its own caching usually
      
      - name: Run Linter
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.60
          args: --timeout=5m

  test:
    name: Test & Build
    runs-on: ubuntu-latest
    needs: quality # Run tests only if lint passes (optional, can run parallel)
    services:
      # Service container for DB if integration tests need it
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: user
          POSTGRES_PASSWORD: password
          POSTGRES_DB: dbname
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      
      - name: Run Tests
        run: make test
        env:
          DB_HOST: localhost
          DB_PORT: 5432
          DB_USER: user
          DB_PASSWORD: password
          DB_NAME: dbname

      - name: Build Check
        run: go build -v ./...
```

### Critical Implementation Details
- **Services:** If integration tests run against a real DB (which they likely do given `testcontainers` isn't strictly mandated yet, but `manual-test-project` might use a local DB), the CI MUST provide a Postgres service container.
- **Environment Variables:** Ensure the test step has the necessary ENV vars to connect to the CI database service.

### Architecture Compliance
- **NFR26:** Explicitly fulfills the requirement for CI/CD via GitHub Actions.
- **Standardization:** Uses `make test` from Story 4.4.

## Dev Agent Record

### Agent Model Used
Gemini 2.0 Flash

### Debug Log References
- Checked Architecture for CI requirements.
- Validated standard GitHub Actions versions (Checkout v4, Setup-Go v5).

### Completion Notes List
- [x] Workflow file created at .github/workflows/ci.yml in manual-test-project.
- [x] Triggers configured for push and pull_request on main branch.
- [x] Linter job configured using golangci-lint-action@v6 with Go 1.25.
- [x] Test job configured with PostgreSQL 16-alpine service container, health checks, and environment variables.
- [x] Build check step added to verify compilation with go build -v ./...
- [x] CLI generator updated: Added GitHubActionsWorkflowTemplate() in templates.go.
- [x] Generator.go updated to include .github/workflows/ci.yml in generated projects.

### File List
- manual-test-project/.github/workflows/ci.yml
- cmd/create-go-starter/templates.go
- cmd/create-go-starter/generator.go

## Change Log
- **2026-01-09**: Implemented GitHub Actions CI/CD workflow with automated linting, testing, and build checks. Configured quality job with golangci-lint and test job with PostgreSQL service container. Updated CLI generator to include CI workflow in new projects.
