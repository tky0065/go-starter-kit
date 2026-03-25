# Story 11.1: Framework Selection Flag

Status: done

## Story

As a développeur,
I want pouvoir spécifier le type de framework via un flag `--framework`,
So that je peux choisir mon framework web préféré (Fiber, Gin, Echo).

## Acceptance Criteria

**Given** l'utilisateur exécute `create-go-starter mon-projet --framework=gin`
**When** le CLI parse les arguments
**Then** Gin est sélectionné comme framework
**And** Fiber reste le défaut si non spécifié

## Tasks / Subtasks

- [x] Task 1: Define framework constants and validation (AC: #1)
  - [x] Add framework constants (Fiber, Gin, Echo) to main.go
  - [x] Add ValidFrameworks slice and DefaultFramework constant
  - [x] Implement validateFramework() function with proper error messages
  - [x] Add framework descriptions for help text

- [x] Task 2: Add framework flag to CLI (AC: #1)
  - [x] Add --framework/-f flag in run() function
  - [x] Add framework to InteractiveDefaults struct
  - [x] Update validateArgs() to include framework validation
  - [x] Add framework to help text and usage examples

- [x] Task 3: Add framework to TUI interactive mode (AC: #1)
  - [x] Add framework field to tui.InteractiveDefaults
  - [x] Add framework selection step in TUI form (tui/form.go)
  - [x] Update TUI generator function signature to accept framework parameter
  - [x] Add framework option descriptions in TUI

- [x] Task 4: Pass framework to generator functions (AC: #1)
  - [x] Update generateProjectFiles() signature to include framework parameter
  - [x] Update runWithCallback() to pass framework to generator
  - [x] Update all template builder functions to accept framework parameter
  - [x] Add framework parameter to FileGenerator construction

- [x] Task 5: Update tests and documentation (AC: #1)
  - [x] Add unit tests for validateFramework()
  - [x] Update CLI tests to verify framework flag parsing
  - [x] Update TUI E2E tests to verify framework selection
  - [x] Update README.md with framework flag examples
  - [x] Update docs/usage.md with framework options

## Dev Notes

### Current CLI Flag Pattern Analysis

**Existing Flag Implementations:**
```go
// Current flags in main.go (lines 46-75):
- --template/-t: TemplateMinimal, TemplateFull, TemplateGraphQL (default: full)
- --database/-d: DatabasePostgres, DatabaseMySQL, DatabaseSQLite (default: postgres)
- --observability/-o: ObservabilityNone, ObservabilityBasic, ObservabilityAdvanced (default: none)
```

**Framework Flag should follow same pattern:**
```go
// Framework type constants
const (
	FrameworkFiber = "fiber"
	FrameworkGin   = "gin"
	FrameworkEcho  = "echo"
)

// ValidFrameworks contains the list of valid framework types
var ValidFrameworks = []string{FrameworkFiber, FrameworkGin, FrameworkEcho}

// DefaultFramework is the default framework when not specified
const DefaultFramework = FrameworkFiber

// Framework descriptions (in English for consistency)
const (
	FrameworkFiberDesc = "Fiber v2.52.10 - Fast HTTP framework inspired by Express (default)"
	FrameworkGinDesc   = "Gin - High-performance HTTP web framework (planned)"
	FrameworkEchoDesc  = "Echo - Minimalist high-performance HTTP framework (planned)"
)
```

### Architecture Compliance

**From architecture.md:**
- Stack Core: Fiber (Web), GORM (ORM), PostgreSQL (DB), fx (DI)
- Only Fiber is currently implemented
- Multi-framework support is Phase 3 (Epic 11)

**Critical Requirements:**
1. **Backward Compatibility:** Fiber MUST remain the default framework
2. **Validation:** Framework validation MUST happen before any file generation
3. **Error Messages:** Clear messages for unsupported frameworks (Gin/Echo not yet implemented)
4. **TUI Integration:** Framework selection MUST be available in interactive mode

### Implementation Guidelines

**Step 1: Add Constants (main.go, after line 75)**
```go
// Framework type constants define the available web frameworks.
// - FrameworkFiber: Fiber v2.52.10 (default, fully implemented)
// - FrameworkGin: Gin framework (planned for v2.0.0)
// - FrameworkEcho: Echo framework (planned for v2.0.0)
const (
	FrameworkFiber = "fiber"
	FrameworkGin   = "gin"
	FrameworkEcho  = "echo"
)

// ValidFrameworks contains the list of valid framework types
// Note: Only Fiber is currently implemented (Stories 11.2-11.3 add Gin/Echo)
var ValidFrameworks = []string{FrameworkFiber, FrameworkGin, FrameworkEcho}

// DefaultFramework is the default framework when not specified
const DefaultFramework = FrameworkFiber
```

**Step 2: Add Validation Function (main.go, after validateDatabase)**
```go
// validateFramework checks if the framework type is valid and supported.
// Valid frameworks are: fiber, gin, echo
// Note: gin and echo are in Epic 11 roadmap (Stories 11.2-11.3)
func validateFramework(framework string) error {
	// Check if it's a valid framework
	for _, valid := range ValidFrameworks {
		if framework == valid {
			// Only Fiber is currently supported
			if framework != FrameworkFiber {
				return fmt.Errorf("framework '%s' is not yet supported (planned for v2.0.0).\nCurrently supported: fiber", framework)
			}
			return nil
		}
	}

	// Generic error for invalid frameworks
	return fmt.Errorf("invalid framework '%s'.\nSupported frameworks: %s", framework, strings.Join(ValidFrameworks, ", "))
}
```

**Step 3: Add Flag to run() Function**
```go
// In run() function, add flag definition:
frameworkFlag := flags.String("framework", DefaultFramework, "Web framework to use (fiber|gin|echo)")
flags.StringVar(frameworkFlag, "f", DefaultFramework, "Web framework (shorthand)")

// In validation section:
if err := validateFramework(*frameworkFlag); err != nil {
	return err
}

framework := *frameworkFlag
```

**Step 4: Update Generator Signatures**
```go
// Update function signatures to accept framework parameter:
func generateProjectFiles(projectPath, projectName, template, database, observabilityLevel, framework string) error

func runWithCallback(projectName, template, database, observability, framework string, progressCallback func(current, total int), quiet bool) error

// Pass framework to template constructors:
templates := NewProjectTemplatesWithDatabaseAndFramework(projectName, database, framework)
```

**Step 5: TUI Integration**
```go
// In tui/model.go, add to InteractiveDefaults:
type InteractiveDefaults struct {
	ProjectName   string
	Template      string
	Database      string
	Observability string
	Framework     string  // NEW
}

// In tui/form.go, add framework selection step:
// Add after database selection, before observability
```

### File Structure Requirements

**Files to modify:**
1. `cmd/create-go-starter/main.go` - Constants, validation, flag parsing
2. `cmd/create-go-starter/generator.go` - Function signatures, framework parameter passing
3. `cmd/create-go-starter/tui/model.go` - InteractiveDefaults struct
4. `cmd/create-go-starter/tui/form.go` - Framework selection step
5. `cmd/create-go-starter/tui/integration.go` - Generator function signature

**Files NOT to modify in this story:**
- Template files (templates.go, templates_user.go) - Stories 11.2-11.3
- Generator logic for Gin/Echo - Stories 11.2-11.3
- Domain/business logic - Already framework-agnostic

### Testing Requirements

**Unit Tests (main_test.go):**
```go
func TestValidateFramework(t *testing.T) {
	tests := []struct {
		name      string
		framework string
		wantErr   bool
		errMsg    string
	}{
		{"valid fiber", "fiber", false, ""},
		{"valid gin (not implemented)", "gin", true, "not yet supported"},
		{"valid echo (not implemented)", "echo", true, "not yet supported"},
		{"invalid framework", "express", true, "invalid framework"},
	}
	// ...
}
```

**E2E Tests (smoke_test.go):**
```bash
# Verify framework flag acceptance
./create-go-starter test-fiber-project --framework=fiber
./create-go-starter test-gin-project --framework=gin # Should fail with helpful message
```

**TUI E2E Tests (tui/e2e_test.go):**
```go
func TestFrameworkSelectionInTUI(t *testing.T) {
	// Verify framework appears in TUI form
	// Verify default is Fiber
	// Verify selection persists to generation
}
```

### Project Context Notes

**From project-context.md:**
- **Web Framework:** Fiber v2.52.10 (current default)
- **Routing:** MUST use route groups with `/api/v1` prefix
- **Dependency Injection:** ALL components MUST be registered with `fx`

**Critical Rules:**
1. **NO BREAKING CHANGES:** Existing projects using Fiber must continue to work
2. **Clear Error Messages:** Users selecting Gin/Echo must get helpful "not yet implemented" messages
3. **Documentation:** All examples should show Fiber as default with framework flag optional
4. **Conventional Commits:** Use `feat(epic-11):` prefix for all commits

### References

**Source Documentation:**
- [Source: _bmad-output/planning-artifacts/epics.md#Epic-11-Story-11.1]
- [Source: _bmad-output/planning-artifacts/architecture.md#Core-Architectural-Decisions]
- [Source: _bmad-output/project-context.md#Framework-Specific-Rules]
- [Source: cmd/create-go-starter/main.go#L46-L75] - Existing flag patterns

**Related Stories:**
- Story 11.2: Gin Framework Templates (depends on 11.1)
- Story 11.3: Echo Framework Templates (depends on 11.1)
- Story 11.4: Framework-Agnostic Abstraction (depends on 11.2, 11.3)

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-6

### Debug Log References

N/A - No significant debugging required. Standard compilation errors from threading framework parameter through function signatures.

### Completion Notes List

- [x] Framework flag works in CLI mode (`--framework`/`-f`)
- [x] Framework flag works in interactive TUI mode (new `StateFrameworkSelect` step)
- [x] Framework flag works in interactive non-TUI mode (step 5 in `runInteractiveModeWithReader`)
- [x] Validation prevents unsupported frameworks with clear messages (gin/echo return "not yet supported")
- [x] Tests pass (unit + E2E + TUI) — all tests green
- [x] Documentation updated (README, docs/usage.md)

### File List

- `cmd/create-go-starter/main.go` — Framework constants, validateFramework(), flag parsing, help text
- `cmd/create-go-starter/generator.go` — generateProjectFiles() signature updated, framework param passed to generate*TemplateFiles()
- `cmd/create-go-starter/dryrun.go` — getFilesForTemplate() and runDryRun() signatures updated
- `cmd/create-go-starter/interactive.go` — InteractiveDefaults.Framework, step 5 framework selection
- `cmd/create-go-starter/tui/model.go` — GeneratorFunc type, StateFrameworkSelect, framework/frameworkList fields
- `cmd/create-go-starter/tui/interactive.go` — InteractiveDefaults.Framework, initializeFrameworkList(), frameworkItem
- `cmd/create-go-starter/tui/messages.go` — FrameworkSelectedMsg
- `cmd/create-go-starter/tui/update.go` — FrameworkSelectedMsg handler, navigateBack, handleKeyMsg, updateComponents, generateProjectCmd
- `cmd/create-go-starter/tui/view.go` — viewFrameworkSelect(), View() switch, viewSummary(), viewObservabilitySelect() breadcrumb
- `cmd/create-go-starter/tui/form.go` — totalSteps 4→5, Framework step name, framework field, renderFrameworkStep(), View() case 4/5
- `cmd/create-go-starter/tui/confirmation.go` — Updated for framework parameter
- `cmd/create-go-starter/main_test.go` — New framework tests (TestValidateFramework*, TestFrameworkFlag*, etc.)
- `cmd/create-go-starter/generator_test.go` — Updated generateProjectFiles() calls
- `cmd/create-go-starter/dryrun_test.go` — Updated function call signatures
- `cmd/create-go-starter/interactive_test.go` — Updated test inputs for new framework step
- `cmd/create-go-starter/smoke_test.go` — Updated run() calls with DefaultFramework
- `cmd/create-go-starter/database_integration_test.go` — Updated run() calls with DefaultFramework
- `cmd/create-go-starter/e2e_mysql_test.go` — Updated run() calls with DefaultFramework
- `cmd/create-go-starter/e2e_sqlite_test.go` — Updated run() calls with DefaultFramework
- `cmd/create-go-starter/template_minimal_test.go` — Updated run() calls with DefaultFramework
- `cmd/create-go-starter/templates_observability_test.go` — Updated run() calls with DefaultFramework
- `cmd/create-go-starter/tui/interactive_test.go` — dummyGenerator signature updated
- `cmd/create-go-starter/tui/e2e_test.go` — mockGenerator signatures, flow updated for StateFrameworkSelect
- `cmd/create-go-starter/tui/form_test.go` — Updated for 5 steps
- `README.md` — Framework section added (v1.6.0)
- `docs/usage.md` — Framework section added with table
- `_bmad-output/implementation-artifacts/sprint-status.yaml` — Status in-progress→review
