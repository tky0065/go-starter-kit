# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `go-starter-kit`, a Go CLI tool generator that scaffolds new Go projects. The tool is called `create-go-starter` and helps bootstrap Go project structures.

## Commands

### Build and Install
```bash
# Build the binary
go build -o create-go-starter ./cmd/create-go-starter

# Install the binary to GOBIN (typically ~/go/bin/)
go install ./cmd/create-go-starter

# Run directly without installing
go run ./cmd/create-go-starter <project-name>
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./cmd/create-go-starter

# Run tests with verbose output
go test -v ./...

# Run a specific test function
go test -run TestColors ./cmd/create-go-starter

# Run tests with coverage report
go test -cover ./...

# Using Makefile targets
make test              # Run all tests
make test-short        # Quick unit tests (skip long-running tests)
make smoke-test        # Full E2E validation with runtime tests
make smoke-test-quick  # E2E validation without runtime (no Docker needed)
```

### Development
```bash
# Format code
go fmt ./...

# Vet code for common mistakes
go vet ./...

# Run the tool
./create-go-starter <project-name>

# Or with go run
go run ./cmd/create-go-starter <project-name>
```

### After Creating a Project

Once you've created a project with `create-go-starter`, you have two setup options:

**Option 1: Automated Setup (Recommended)**
```bash
cd <project-name>
./setup.sh
make run
```

The `setup.sh` script automates:
- Go dependency installation (`go mod tidy`)
- JWT secret generation (`openssl rand -base64 32`)
- PostgreSQL configuration (Docker or local)
- Installation verification

**Option 2: Manual Setup**
```bash
cd <project-name>

# Install dependencies
go mod tidy

# Generate JWT secret and add to .env
openssl rand -base64 32
# Edit .env and add: JWT_SECRET=<generated-secret>

# Start PostgreSQL (Docker)
docker run -d --name postgres \
  -e POSTGRES_DB=<project-name> \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Run the application
make run
```

## Project Structure

### CLI Tool Structure (this repository)
- `cmd/create-go-starter/` - Main CLI application entry point
  - `main.go` - CLI orchestration, validation, and entry point (refactored with `run()` function for testability)
  - `generator.go` - File generation orchestrator
  - `templates.go` - Core templates (config, server, domain, setup.sh)
  - `templates_user.go` - User domain specific templates
  - `tui/` - Interactive terminal UI package (welcome, forms, progress, help, styles)
  - `git.go` - Git repository initialization and initial commit
  - `git_test.go` - Tests for Git functionality
  - `smoke_test.go` - End-to-end smoke tests for project generation validation
  - `*_test.go` - Unit tests for each module (colors, env, generator, main, scaffold, templates)
- `scripts/smoke_test.sh` - Bash script for comprehensive E2E validation
- `go.mod` - Go module definition (requires Go 1.25)
- `_bmad/` - BMAD workflow automation system (not part of core application)

### Generated Project Structure (projects created by the CLI)
When you run `create-go-starter my-project`, the following structure is generated:
- `internal/models/` - Shared domain entities (User, RefreshToken, AuthResponse). Prevents circular dependencies.
- `internal/domain/` - Business logic (services, not entities)
- `internal/interfaces/` - Port definitions (Hexagonal Architecture)
- `internal/adapters/` - HTTP handlers, middleware, repository implementations
- `internal/infrastructure/` - Database, server configuration
- `pkg/` - Shared libraries (config, logger, auth)

## Architecture Notes

### CLI Tool Design
The tool uses a standard flag-driven CLI plus an interactive terminal UI:
- Flag parsing using the standard `flag` package
- ANSI color support via helper functions `Green()` and `Red()`
- Interactive flows implemented in `cmd/create-go-starter/tui/`

### Code Organization
- CLI logic is organized across multiple files for maintainability:
  - `main.go` - Entry point and orchestration
  - `generator.go` - File generation logic
  - `templates.go` / `templates_user.go` - Template definitions
  - `git.go` - Git integration (automatic repository initialization)
- Color utilities are embedded in the main package and tested separately
- The tool is designed to be distributed as a single binary via `go install`
- Test coverage: 83%+ with comprehensive unit and smoke tests

### Testing Approach
- Unit tests are co-located with source code (e.g., `colors_test.go` alongside `main.go`)
- Tests use standard Go testing patterns with table-driven tests where appropriate

## Documentation Maintenance

**CRITICAL REQUIREMENT**: Documentation must ALWAYS be updated when making changes to the codebase.

### When to Update Documentation

Update documentation immediately after ANY of the following changes:

1. **Adding new features**:
   - New templates added → Update `docs/cli-architecture.md`
   - New generated files → Update `docs/usage.md` and `docs/generated-project-guide.md`
   - New domain entities → Update all architecture diagrams and examples

2. **Modifying architecture**:
   - Package structure changes → Update ALL documentation files
   - Dependency changes → Update architecture diagrams
   - Pattern changes → Update examples in `docs/generated-project-guide.md`

3. **Changing templates**:
   - Template modifications → Update corresponding documentation sections
   - New imports or dependencies → Update code examples throughout docs

4. **Bug fixes that affect structure**:
   - If the fix changes how generated projects work → Update docs immediately

### Documentation Files to Check

After ANY code change, review and update these files as needed:

- `README.md` - Main project overview and quick start
- `docs/usage.md` - Usage guide and generated structure
- `docs/generated-project-guide.md` - Complete guide for generated projects
- `docs/cli-architecture.md` - CLI tool architecture documentation
- `CLAUDE.md` - This file (AI context)
- `GEMINI.md` - Gemini AI context
- `AGENTS.md` - Agents AI context
- `site/**`   - Officiel docs site 

### Documentation Update Process

1. **Make code changes** (templates, generator, etc.)
2. **Test the changes** (generate a test project, verify it works)
3. **Update ALL affected documentation** (use search to find all references)
4. **Verify documentation** (read through to ensure consistency)
5. **Commit code AND documentation together** in the same commit or sequential commits

### Example Workflow

```bash
# 1. Make code changes
vim cmd/create-go-starter/templates_user.go

# 2. Test changes
go build -o create-go-starter ./cmd/create-go-starter
./create-go-starter test-project
cd test-project && go build ./...

# 3. Update documentation
vim docs/cli-architecture.md
vim docs/generated-project-guide.md
vim docs/usage.md

# 4. Commit together
git add cmd/ docs/ README.md CLAUDE.md
git commit -m "feat: add new feature

- Implement feature X in templates
- Update all documentation to reflect changes
- Add examples in generated-project-guide.md"

git push origin main
```

### Why This Matters

- **Users rely on accurate docs**: Out-of-date documentation causes confusion and wastes time
- **Future contributors**: Need current docs to understand the system
- **AI assistants**: Use these docs as context for helping with the project
- **Consistency**: Prevents drift between code and documentation

**Remember**: Code without updated documentation is incomplete work. Always update docs as part of your changes.

## MkDocs Documentation Site

### Overview

The project uses **Material for MkDocs** to generate the official documentation site. The documentation is written in Markdown and deployed to GitHub Pages.

### Commands

```bash
# Activate Python virtual environment
source venv/bin/activate

# Build the documentation site
mkdocs build

# Build with clean (removes old files)
mkdocs build --clean

# Serve documentation locally
mkdocs serve
# View at: http://127.0.0.1:8000/go-starter-kit/

# Deploy to GitHub Pages
mkdocs gh-deploy
```

### Project Structure

```
docs/                           # Documentation source files
├── index.md                    # Home page
├── installation.md             # Installation guide
├── usage.md                    # Usage guide
├── databases.md                # Database selection guide
├── database-migration.md       # Database migration guide
├── generated-project-guide.md  # Complete guide for generated projects
├── tutorial-exemple-complet.md # Complete tutorial
├── cli-architecture.md         # CLI architecture documentation
├── contributing.md             # Contributing guidelines
├── stylesheets/               # Custom CSS
│   └── extra.css              # Material Icons and custom styles
└── overrides/                 # Theme overrides (if any)

mkdocs.yml                     # MkDocs configuration
site/                          # Generated site (git-ignored)
venv/                          # Python virtual environment
```

### Material Icons Usage

The documentation uses **Material Icons web font** for consistent iconography.

#### Configuration

Material Icons are enabled via custom CSS (`docs/stylesheets/extra.css`):

```css
/* Import Material Icons font */
@import url('https://fonts.googleapis.com/icon?family=Material+Icons');

/* Base Material Icons styles */
.material-icons {
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  font-size: 1.2em;
  display: inline-block;
  line-height: 1;
  text-transform: none;
  letter-spacing: normal;
  word-wrap: normal;
  white-space: nowrap;
  direction: ltr;
  vertical-align: middle;
}

/* Color variants */
.material-icons.success { color: var(--md-success-fg-color, #00c853); }
.material-icons.warning { color: var(--md-warning-fg-color, #ff6d00); }
.material-icons.error { color: var(--md-error-fg-color, #ff1744); }
.material-icons.info { color: var(--md-info-fg-color, #00b0ff); }

/* Size variants */
.material-icons.small { font-size: 1em; }
.material-icons.large { font-size: 1.5em; }
```

The CSS file is loaded in `mkdocs.yml`:

```yaml
extra_css:
  - stylesheets/extra.css
```



#### Using Material Icons in Markdown

**IMPORTANT**: Do NOT use the `:material-icon-name:` syntax as it renders as emoji. Instead, use HTML with the Material Icons web font:

**❌ WRONG (renders as emoji):**
```markdown
:material-check: Feature enabled
:material-alert: Warning message
```

**✅ CORRECT (renders as Material Icons):**
```markdown
<i class="material-icons success">check</i> Feature enabled
<i class="material-icons warning">warning</i> Warning message
```

#### Common Icons Reference

| Purpose | Icon Name | HTML Code |
|---------|-----------|-----------|
| Success/Check | `check` | `<i class="material-icons success">check</i>` |
| Success Circle | `check_circle` | `<i class="material-icons success">check_circle</i>` |
| Warning | `warning` | `<i class="material-icons warning">warning</i>` |
| Error | `error` | `<i class="material-icons warning">error</i>` |
| Info | `info` | `<i class="material-icons info">info</i>` |
| Target/Focus | `center_focus_strong` | `<i class="material-icons">center_focus_strong</i>` |
| Book | `menu_book` | `<i class="material-icons">menu_book</i>` |
| Sync | `sync` | `<i class="material-icons">sync</i>` |
| Back Arrow | `arrow_back` | `<i class="material-icons">arrow_back</i>` |
| Circle Indicator | `circle` | `<i class="material-icons small">circle</i>` |

**Find more icons:** [Google Material Icons](https://fonts.google.com/icons?icon.set=Material+Icons)

#### Icon Style Classes

**Color Classes:**
- `.success` - Green (for positive/confirmed items)
- `.warning` - Orange (for warnings/cautions)
- `.error` - Red (for errors/critical items)
- `.info` - Blue (for informational items)

**Size Classes:**
- `.small` - 1em (inline with text)
- (default) - 1.2em (slightly larger)
- `.large` - 1.5em (prominent)

**Examples:**
```markdown
<!-- Small success indicator -->
<i class="material-icons success small">circle</i> Easy

<!-- Default warning icon -->
<i class="material-icons warning">warning</i> Requires attention

<!-- Large error icon -->
<i class="material-icons error large">error</i> Critical issue
```

### Common Issues and Solutions

#### Issue: Icons render as emoji instead of Material Icons
**Solution:** Replace `:material-icon-name:` syntax with HTML `<i class="material-icons">icon_name</i>` tags.

#### Issue: Icons don't have color
**Solution:** Add color class: `<i class="material-icons success">check</i>`

#### Issue: Custom CSS not loading
**Solution:**
1. Verify `docs/stylesheets/extra.css` exists
2. Check `mkdocs.yml` has `extra_css: - stylesheets/extra.css`
3. Run `mkdocs build --clean` to rebuild

#### Issue: Icons not showing after deployment
**Solution:**
1. Ensure Google Fonts can be accessed (not blocked by CSP)
2. Check browser console for font loading errors
3. Verify `@import url('https://fonts.googleapis.com/icon?family=Material+Icons');` in CSS

### Updating Documentation

When adding or modifying documentation:

1. **Edit markdown files** in `docs/` directory
2. **Use Material Icons HTML syntax** (not emoji syntax)
3. **Test locally** with `mkdocs serve`
4. **Build and verify** with `mkdocs build --clean`
5. **Deploy** with `mkdocs gh-deploy` (or push to trigger CI/CD)

### Documentation Standards

- Use French for all content (this is a French-language project)
- Include code examples for all features
- Add navigation links at the top of each page
- Use admonitions for important notes (warning, info, danger)
- Keep line length reasonable (80-120 characters)
- Use Material Icons consistently throughout docs

## Commit Messages
- Use present tense ("Add feature" not "Added feature")
- NOT Co-Authored-By lines in commit messages .
