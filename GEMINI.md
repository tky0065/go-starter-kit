# Gemini Project Context: go-starter-kit

## Project Overview

`go-starter-kit` is a CLI productivity tool designed to scaffold production-ready Go API projects in minutes. It aims to eliminate the friction of starting new projects by providing an opinionated yet pragmatic foundation that balances development speed with industrial standards.

### Core Philosophy
- **Zero-to-Hero Experience:** A CLI tool (`create-go-starter`) that handles scaffolding, dependency management, and initial configuration.
- **Best Practices by Default:** Includes JWT authentication, validation, Swagger documentation, Dockerization, and CI/CD out of the box.
- **Lite Hexagonal Architecture:** A simplified hexagonal structure that ensures maintainability without excessive complexity.

### Tech Stack
- **Web Framework:** [Fiber](https://gofiber.io/) (Fast, Express-inspired)
- **Dependency Injection:** [uber-go/fx](https://github.com/uber-go/fx)
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL support.
- **Logging:** [zerolog](https://github.com/rs/zerolog)
- **Documentation:** Swagger (auto-generated)
- **Containerization:** Docker (multi-stage builds)

---

## Directory Structure

### Main Tool (`/cmd/create-go-starter`)
- `main.go`: CLI entry point, handles project name validation and directory creation (refactored with `run()` function for testability).
- `generator.go`: Orchestrates the file generation process using templates.
- `templates.go`: Contains all boilerplate code templates for the generated project.
- `templates_user.go`: User domain specific templates.
- `git.go`: Handles Git repository initialization and initial commit creation.
- `smoke_test.go`: End-to-end smoke tests for project generation validation.
- `scripts/smoke_test.sh`: Bash script for comprehensive E2E validation.

### Generated Project Structure (Boilerplate)
- `cmd/main.go`: Application entry point using `fx` for dependency injection.
- `internal/`:
    - `models/`: Shared domain entities (User, RefreshToken, AuthResponse). Prevents circular dependencies.
    - `domain/`: Core business logic (services, not entities).
    - `adapters/`: External interfaces (HTTP handlers, Middleware).
    - `infrastructure/`: Infrastructure concerns (Database connection, Server configuration).
    - `interfaces/`: Port definitions (Ports in Hexagonal Architecture).
- `pkg/`: Public/Shared libraries (Config loader, Logger).
- `deployments/`: Deployment configurations (Docker, K8s).

---

## Building and Running

### The Scaffolding Tool (Main Repository)

- **Build the CLI:**
  ```bash
  make build
  ```
- **Run the CLI (Generates a new project):**
  ```bash
  make run PROJECT_NAME=my-new-api
  ```
- **Run Tests:**
  ```bash
  make test
  ```
- **Install CLI Globally:**
  ```bash
  make install
  ```

### The Generated Project

**Option 1: Automated Setup (Recommended)**
```bash
cd <project-name>
./setup.sh
make run
```

The `setup.sh` script automatically handles:
- Go dependency installation (`go mod tidy`)
- JWT secret generation and configuration
- PostgreSQL setup (Docker or local)
- Installation verification

**Option 2: Manual Setup**
1.  **Initialize:**
    ```bash
    cd <project-name>
    go mod tidy
    ```
2.  **Configure Environment:**
    ```bash
    # Generate JWT secret
    openssl rand -base64 32
    # Edit .env and add: JWT_SECRET=<generated-secret>
    ```
3.  **Start PostgreSQL:**
    ```bash
    docker run -d --name postgres \
      -e POSTGRES_DB=<project-name> \
      -e POSTGRES_PASSWORD=postgres \
      -p 5432:5432 \
      postgres:16-alpine
    ```
4.  **Run in Development:**
    ```bash
    make run
    # OR with hot-reload (requires air)
    make dev
    ```
5.  **Run Tests:**
    ```bash
    make test
    ```

---

## Development Conventions

- **Hexagonal Architecture:** Maintain a strict separation between the domain logic and external adapters. Domain logic should not depend on infrastructure or adapters.
- **Dependency Injection:** Always use `fx.Module` and `fx.Provide` for defining components. Avoid global variables.
- **Error Handling:** Use centralized error handling (middleware) and define domain-specific errors in `internal/domain/errors.go`.
- **Environment Variables:** All configuration must be driven by environment variables using the `pkg/config` utility.
- **Linting:** Ensure code passes `golangci-lint` (config provided in templates).
- **Graceful Shutdown:** The server and database connections must handle OS signals for clean termination.

## Key Files for Reference
- `cmd/create-go-starter/main.go`: Logic for the CLI tool.
- `cmd/create-go-starter/templates.go`: Source of truth for the generated boilerplate.
- `go.mod`: Project dependencies.
- `_bmad-output/planning-artifacts/prd.md`: High-level product requirements and vision.

---

## Documentation Maintenance Policy

**⚠️ MANDATORY REQUIREMENT**: Documentation MUST be updated whenever code changes are made.

### Update Documentation When:

1. **Adding Features**:
   - New templates → Update `docs/cli-architecture.md`
   - New generated files → Update `docs/usage.md` and `docs/generated-project-guide.md`
   - New entities/models → Update all architecture diagrams and code examples

2. **Modifying Architecture**:
   - Package restructuring → Update ALL documentation files
   - Dependency graph changes → Update architecture diagrams
   - Design pattern changes → Update examples and explanations

3. **Template Changes**:
   - Any template modification → Update corresponding documentation
   - Import changes → Update all code examples throughout docs

4. **Bug Fixes**:
   - If fix changes generated project behavior → Documentation update required

### Critical Documentation Files

Always review these after code changes:

| File | Purpose | Update When |
|------|---------|-------------|
| `README.md` | Project overview | Structure or features change |
| `docs/usage.md` | Usage guide | Generated structure changes |
| `docs/generated-project-guide.md` | Complete guide | Any template or architecture change |
| `docs/cli-architecture.md` | CLI internals | Generator or template logic changes |
| `CLAUDE.md` | AI context (this repo) | Project structure or conventions change |
| `GEMINI.md` | AI context (this file) | Project structure or conventions change |

### Standard Workflow

```bash
# Step 1: Code changes
vim cmd/create-go-starter/templates_user.go

# Step 2: Test thoroughly
go build -o create-go-starter ./cmd/create-go-starter
./create-go-starter test-validation-project
cd test-validation-project && go mod tidy && go build ./...

# Step 3: Update ALL affected docs
vim docs/cli-architecture.md
vim docs/generated-project-guide.md
vim README.md

# Step 4: Verify documentation accuracy
grep -r "old_pattern" docs/  # Find outdated references
# Update all found instances

# Step 5: Commit code + docs together
git add cmd/ docs/ README.md CLAUDE.md GEMINI.md
git commit -m "feat: implement feature X

- Add feature implementation in templates
- Update all documentation to reflect changes
- Add comprehensive examples in docs"
```

### Why This Is Critical

- **User Trust**: Outdated docs erode confidence in the project
- **Development Velocity**: Accurate docs enable faster onboarding and contributions
- **AI Assistance**: LLMs rely on current documentation for context
- **Maintenance**: Prevents technical debt accumulation

### Enforcement

- Pull requests with code changes but no doc updates will be questioned
- Documentation updates are NOT optional—they are part of the feature
- When in doubt, over-document rather than under-document

**Remember**: Undocumented code changes are considered incomplete work.

---

## MkDocs Documentation Site

### Overview

The project uses **Material for MkDocs** to generate and publish documentation to GitHub Pages. All documentation is written in Markdown (French language).

### Setup and Commands

```bash
# Activate Python virtual environment (required)
source venv/bin/activate

# Local development server (with live reload)
mkdocs serve
# Access at: http://127.0.0.1:8000/go-starter-kit/

# Build static site
mkdocs build --clean

# Deploy to GitHub Pages
mkdocs gh-deploy
```

### Directory Structure

```
docs/                           # Documentation source (Markdown)
├── index.md                    # Home page
├── installation.md
├── usage.md
├── databases.md                # Database selection guide
├── database-migration.md       # Migration guide
├── generated-project-guide.md
├── tutorial-exemple-complet.md
├── cli-architecture.md
├── contributing.md
├── stylesheets/
│   └── extra.css              # Material Icons + custom styles
└── overrides/                 # Theme customizations

mkdocs.yml                     # MkDocs configuration file
site/                          # Generated static site (git-ignored)
venv/                          # Python dependencies
```

### Material Icons in Markdown

**⚠️ CRITICAL RULE**: Use HTML with Material Icons web font, NOT emoji syntax.

#### Why Not Emoji Syntax?

The syntax `:material-icon-name:` is **NOT supported** by Material for MkDocs in markdown content. It renders as TWemoji (Twitter emoji), not Material Design icons.

#### Correct Usage

Material Icons are configured via `docs/stylesheets/extra.css` which imports the Google Fonts Material Icons web font and defines helper classes.

**HTML Syntax (Correct):**

```html
<!-- Success icon (green) -->
<i class="material-icons success">check</i>

<!-- Warning icon (orange) -->
<i class="material-icons warning">warning</i>

<!-- Small indicator -->
<i class="material-icons small">circle</i>

<!-- Default size, no color -->
<i class="material-icons">menu_book</i>
```

**Emoji Syntax (WRONG - renders as emoji):**

```markdown
:material-check:        ❌ Renders as emoji
:material-warning:      ❌ Renders as emoji
:material-book-open:    ❌ Renders as emoji
```

#### Available Classes

**Color Classes:**
- `.success` - Green (#00c853) - Use for positive states, checkmarks, completed items
- `.warning` - Orange (#ff6d00) - Use for warnings, cautions, requires attention
- `.error` - Red (#ff1744) - Use for errors, critical issues, blocked items
- `.info` - Blue (#00b0ff) - Use for informational items, tips

**Size Classes:**
- `.small` - 1em (inline with text, good for indicators)
- (no class) - 1.2em (default, slightly larger than text)
- `.large` - 1.5em (prominent icons)

**Combining Classes:**

```html
<i class="material-icons success small">circle</i>     <!-- Small green dot -->
<i class="material-icons warning large">error</i>      <!-- Large orange error -->
```

#### Common Icons Reference

| Use Case | Icon Name | Example Code |
|----------|-----------|--------------|
| Success/Checkmark | `check` | `<i class="material-icons success">check</i>` |
| Success Badge | `check_circle` | `<i class="material-icons success">check_circle</i>` |
| Warning | `warning` | `<i class="material-icons warning">warning</i>` |
| Error | `error` | `<i class="material-icons warning">error</i>` |
| Info | `info` | `<i class="material-icons info">info</i>` |
| Indicator Dot | `circle` | `<i class="material-icons small">circle</i>` |
| Target/Focus | `center_focus_strong` | `<i class="material-icons">center_focus_strong</i>` |
| Book/Documentation | `menu_book` | `<i class="material-icons">menu_book</i>` |
| Sync/Refresh | `sync` | `<i class="material-icons">sync</i>` |
| Back Arrow | `arrow_back` | `<i class="material-icons">arrow_back</i>` |

**Full Icon Library:** [Google Material Icons](https://fonts.google.com/icons?icon.set=Material+Icons)

### Configuration

Material Icons are enabled through two files:

#### 1. `docs/stylesheets/extra.css`

```css
/* Import Material Icons font from Google Fonts */
@import url('https://fonts.googleapis.com/icon?family=Material+Icons');

/* Base class */
.material-icons {
  font-family: 'Material Icons';
  font-size: 1.2em;
  display: inline-block;
  vertical-align: middle;
  /* ... additional styles ... */
}

/* Color and size variants */
.material-icons.success { color: var(--md-success-fg-color, #00c853); }
.material-icons.warning { color: var(--md-warning-fg-color, #ff6d00); }
/* ... etc ... */
```

#### 2. `mkdocs.yml`

```yaml
extra_css:
  - stylesheets/extra.css
```

### Troubleshooting

| Problem | Solution |
|---------|----------|
| Icons render as emoji | Replace `:material-icon:` with `<i class="material-icons">icon</i>` |
| Icons have no color | Add color class: `.success`, `.warning`, `.error`, or `.info` |
| Icons don't appear | Check `extra_css` in `mkdocs.yml` and verify `extra.css` exists |
| Icons missing after deploy | Ensure Google Fonts CDN is accessible (check browser console) |
| Wrong icon size | Add `.small` or `.large` class, or remove class for default |

### Documentation Workflow

1. **Edit** - Modify markdown files in `docs/`
2. **Preview** - Run `mkdocs serve` and check at http://127.0.0.1:8000
3. **Use icons correctly** - HTML syntax with Material Icons font
4. **Build** - Run `mkdocs build --clean` to verify
5. **Deploy** - Run `mkdocs gh-deploy` or push to trigger CI/CD

### Best Practices

- ✅ Use Material Icons HTML syntax consistently throughout documentation
- ✅ Apply color classes to convey meaning (success=green, warning=orange, etc.)
- ✅ Test documentation locally before deploying
- ✅ Keep line length reasonable (80-120 chars) for readability
- ✅ Use French language for all content
- ❌ Don't use emoji syntax (`:material-icon:`) - it won't render correctly
- ❌ Don't forget to activate venv before running mkdocs commands
