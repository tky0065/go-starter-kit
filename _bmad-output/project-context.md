---
project_name: 'go-starter-kit'
user_name: 'Yacoubakone'
date: 'mercredi 7 janvier 2026'
sections_completed:
  ['technology_stack', 'language_rules', 'framework_rules', 'testing_rules', 'quality_rules', 'workflow_rules', 'anti_patterns']
status: 'complete'
rule_count: 22
optimized_for_llm: true
---

# Project Context for AI Agents

_This file contains critical rules and patterns that AI agents must follow when implementing code in this project. Focus on unobvious details that agents might otherwise miss._

---

## Technology Stack & Versions

- **Language:** Go 1.25.5
- **Web Framework:** Fiber v2.52.10
- **ORM:** GORM v1.31.1
- **Dependency Injection:** uber-go/fx v1.24.1
- **Database:** PostgreSQL (via Docker)
- **Logging:** rs/zerolog
- **Validation:** go-playground/validator/v10
- **Auth:** golang-jwt/jwt/v5
- **Docs:** swaggo/swag
- **Linting:** golangci-lint v2.7.2

## Critical Implementation Rules

### Language-Specific Rules (Go)
- **Acronyms:** MUST use all uppercase for acronyms in code (e.g., `UserID`, `APIKey`, `HTTPClient`).
- **Interfaces:** MUST be defined in `/internal/interfaces`. Do not define local interfaces for cross-layer contracts.
- **Error Handling:** Use named errors from `internal/domain/errors.go` for consistent error checking.
- **Pointers:** Use pointers for optional fields in DTOs (e.g., `*string`) to handle null values in JSON correctly.

### Framework-Specific Rules (Fiber & GORM)
- **Routing:** MUST use route groups with `/api/v1` prefix.
- **Dependency Injection:** ALL components MUST be registered with `fx`. No manual instantiation.
- **GORM Tags:** Every domain struct field MUST have explicit `json` (snake_case) and `gorm` tags.
- **Context:** Always propagate the context from Fiber handlers to services and repositories.

### Code Quality & Style Rules
- **Linting:** MUST pass `golangci-lint` v2.7.2 without warnings.
- **Documentation:** Exported functions MUST have comments starting with the function name.
- **Swagger:** Handlers MUST include `@Summary`, `@Router`, and other swag annotations for auto-docs.

### Development Workflow Rules
- **Commits:** Use Conventional Commits (e.g., `feat:`, `fix:`, `chore:`).
- **Makefile:** Use `make dev`, `make test`, `make build` for all development tasks.
- **Docker:** Maintain multi-stage Dockerfile consistency.

### Critical Don't-Miss Rules (Anti-Patterns)
- **NO SECRETS:** NEVER hardcode secrets. Use `.env` or environment variables.
- **ERROR HANDLING:** NEVER ignore errors (`_ = func()`). Handle or return them.
- **NO GLOBALS:** DO NOT use global variables for application state. Use DI (fx).
- **GRACEFUL SHUTDOWN:** Ensure all long-running processes respect the lifecycle managed by `fx`.

---

## Usage Guidelines

**For AI Agents:**
- Read this file before implementing any code.
- Follow ALL rules exactly as documented.
- When in doubt, prefer the more restrictive option.
- Update this file if new patterns emerge.

**For Humans:**
- Keep this file lean and focused on agent needs.
- Update when technology stack changes.
- Review quarterly for outdated rules.

Last Updated: mercredi 7 janvier 2026
