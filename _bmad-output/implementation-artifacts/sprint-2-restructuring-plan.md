# Sprint 2: Documentation Restructuring Plan

**Date**: 2026-02-12
**Objective**: Restructure documentation for better organization, navigation, and maintainability
**Scope**: Split large monolithic documents into modular, focused guides

---

## Executive Summary

### Current State
- `generated-project-guide.md`: **4477 lines** (monolithic)
- `tutorial-exemple-complet.md`: **1585 lines** (comprehensive tutorial)
- Flat navigation structure in `mkdocs.yml`
- All content in single-page documents

### Target State
- Modular guide structure with 9 focused documents
- Tutorial split into 4 logical parts
- Hierarchical navigation in MkDocs
- Quick-start guide (5-minute read)
- Improved discoverability and maintenance

### Benefits
- **Faster navigation**: Users find specific topics quickly
- **Better SEO**: More pages with focused keywords
- **Easier maintenance**: Update specific sections without affecting others
- **Progressive learning**: From quick start to deep dives
- **Reduced cognitive load**: Focused content per page

---

## Part A: Dividing `generated-project-guide.md`

### Document Analysis

**Total Lines**: 4477
**Current Structure**: 11 major sections
**Proposed Split**: 9 focused modules + 1 index

### Section Breakdown with Line Numbers

```
Line    1-19   : Header + Table of Contents
Line   20-925  : Architecture (905 lines)
Line  926-1080 : Configuration (155 lines)
Line 1081-2463 : Développement (1383 lines)
Line 2465-2780 : API Reference (316 lines)
Line 2782-3011 : Tests (230 lines)
Line 3013-3214 : Base de données (202 lines)
Line 3216-3535 : Sécurité (320 lines)
Line 3537-3984 : Déploiement (448 lines)
Line 3986-4196 : Monitoring & Logging (211 lines)
Line 4198-4459 : Bonnes pratiques (262 lines)
Line 4461-4477 : Conclusion (17 lines)
```

---

### Proposed File Structure

#### 1. `docs/guide/index.md`
**Purpose**: Overview and orientation
**Content**: Lines 1-19 + new introductory content
**Size**: ~150 lines

**New Content Outline**:
```markdown
# Guide des projets générés

> Guide complet pour développer, tester et déployer les projets
  créés avec create-go-starter

## Vue d'ensemble

[Brief introduction to generated projects]

## Navigation du guide

### Pour commencer
- **Quick Start** - Démarrage rapide en 5 minutes
- **Configuration** - Variables d'environnement et setup
- **Structure du projet** - Organisation des fichiers

### Architecture et design
- **Architecture hexagonale** - Principes et diagrammes
- **Couche domaine** - Business logic et services
- **Adapters** - Handlers HTTP et repositories
- **Infrastructure** - Database et server

### Développement
- **API Reference** - Endpoints disponibles
- **Tests** - Organisation et best practices
- **Base de données** - Migrations et GORM

### Production
- **Sécurité** - JWT, validation, bonnes pratiques
- **Déploiement** - Docker, Kubernetes, CI/CD
- **Monitoring** - Logging et health checks

## Parcours d'apprentissage recommandé

[Suggested learning paths for different user types]
```

---

#### 2. `docs/guide/architecture.md`
**Purpose**: Hexagonal architecture deep dive
**Source**: Lines 20-289
**Size**: ~270 lines

**Content Sections**:
- Architecture hexagonale (Ports & Adapters)
- Diagramme d'architecture complète (Mermaid)
- Flux d'une requête HTTP (Sequence Diagram)
- Principe de l'Inversion de Dépendances
- Avantages et cas d'usage

**Key Breakpoint**: Line 289 (before "Structure des fichiers et responsabilités")

---

#### 3. `docs/guide/structure.md`
**Purpose**: File structure and responsibilities
**Source**: Lines 290-925
**Size**: ~635 lines

**Content Sections**:
- Structure des fichiers et responsabilités
- Stack technique
- Structure des répertoires détaillée
- Code examples for each layer:
  - `/internal/models`
  - `/internal/domain`
  - `/internal/interfaces`
  - `/internal/adapters`
  - `/internal/infrastructure`
  - `/pkg`

**Key Breakpoint**: Line 925 (before "Configuration")

---

#### 4. `docs/guide/configuration.md`
**Purpose**: Environment setup and configuration
**Source**: Lines 926-1080
**Size**: ~155 lines

**Content Sections**:
- Variables d'environnement
- Générer un JWT_SECRET sécurisé
- Configuration par environnement (dev/staging/prod)
- Configuration PostgreSQL
- Best practices configuration

**Key Breakpoint**: Line 1080 (before "Développement")

---

#### 5. `docs/guide/development.md`
**Purpose**: Daily development workflow and tools
**Source**: Lines 1081-1606
**Size**: ~525 lines

**Content Sections**:
- Workflow quotidien
- Commandes Makefile
- Gestion des Modèles avec `add-model`
- Ajouter une nouvelle fonctionnalité (méthode manuelle)

**Key Breakpoint**: Line 1606 (before "Exemple complet : Entité Product")

---

#### 6. `docs/guide/examples.md`
**Purpose**: Complete practical examples
**Source**: Lines 1607-2463
**Size**: ~856 lines

**Content Sections**:
- Exemple complet : Entité Product
  - Step 1: Créer le modèle
  - Step 2: Créer l'interface repository
  - Step 3: Implémenter le repository
  - Step 4: Créer le service
  - Step 5: Créer l'interface service
  - Step 6: Créer le handler
  - Step 7: Créer le module fx
  - Step 8: Enregistrer les routes
  - Step 9: Mettre à jour main.go
- Résumé : Fichiers créés/modifiés
- Patterns à suivre

**Key Breakpoint**: Line 2463 (before "API Reference")

---

#### 7. `docs/guide/api-reference.md`
**Purpose**: API endpoints documentation
**Source**: Lines 2465-2780
**Size**: ~315 lines

**Content Sections**:
- Endpoints disponibles
- Authentication endpoints
  - POST /api/v1/auth/register
  - POST /api/v1/auth/login
  - POST /api/v1/auth/refresh
- Users endpoints (Protected)
  - GET /api/v1/users/me
  - PUT /api/v1/users/me
  - DELETE /api/v1/users/me
- Workflow complet avec l'API

**Key Breakpoint**: Line 2780 (before "Tests")

---

#### 8. `docs/guide/testing.md`
**Purpose**: Testing strategy and practices
**Source**: Lines 2782-3011
**Size**: ~229 lines

**Content Sections**:
- Organisation des tests
- Exécuter les tests
- Types de tests
  - Tests unitaires
  - Tests d'intégration
  - Tests E2E
- Best practices pour les tests
- Coverage

**Key Breakpoint**: Line 3011 (before "Base de données")

---

#### 9. `docs/guide/database.md`
**Purpose**: Database management and GORM
**Source**: Lines 3013-3214
**Size**: ~201 lines

**Content Sections**:
- Migrations
- Modèles GORM
- Queries avancées
- Performance tips

**Key Breakpoint**: Line 3214 (before "Sécurité")

---

#### 10. `docs/guide/security.md`
**Purpose**: Security best practices
**Source**: Lines 3216-3535
**Size**: ~319 lines

**Content Sections**:
- Authentification JWT
- Protection des routes
- Validation des entrées
- Hashage des mots de passe
- Checklist sécurité

**Key Breakpoint**: Line 3535 (before "Déploiement")

---

#### 11. `docs/guide/deployment.md`
**Purpose**: Deployment and production
**Source**: Lines 3537-3984
**Size**: ~447 lines

**Content Sections**:
- Docker
  - Dockerfile
  - docker-compose.yml
  - Build et run
- Kubernetes
  - deployment.yaml
  - service.yaml
  - Déployer sur K8s
- CI/CD avec GitHub Actions
- Déploiement en production

**Key Breakpoint**: Line 3984 (before "Monitoring & Logging")

---

#### 12. `docs/guide/monitoring.md`
**Purpose**: Monitoring, logging, observability
**Source**: Lines 3986-4196
**Size**: ~210 lines

**Content Sections**:
- Logging avec zerolog
- Monitoring (recommandations)
- Health checks

**Key Breakpoint**: Line 4196 (before "Bonnes pratiques")

---

#### 13. `docs/guide/best-practices.md`
**Purpose**: Best practices and patterns
**Source**: Lines 4198-4477
**Size**: ~279 lines

**Content Sections**:
- Architecture
- Code style
- Naming conventions
- Error handling patterns
- Testing best practices
- Performance
- Sécurité recap
- Conclusion

---

### Implementation Strategy for Part A

**Phase 1: Create New Files (Do NOT delete original yet)**

```bash
# Create guide directory
mkdir -p docs/guide

# Create all new files with content extracted from original
# Lines 1-19 + new intro → docs/guide/index.md
# Lines 20-289 → docs/guide/architecture.md
# Lines 290-925 → docs/guide/structure.md
# Lines 926-1080 → docs/guide/configuration.md
# Lines 1081-1606 → docs/guide/development.md
# Lines 1607-2463 → docs/guide/examples.md
# Lines 2465-2780 → docs/guide/api-reference.md
# Lines 2782-3011 → docs/guide/testing.md
# Lines 3013-3214 → docs/guide/database.md
# Lines 3216-3535 → docs/guide/security.md
# Lines 3537-3984 → docs/guide/deployment.md
# Lines 3986-4196 → docs/guide/monitoring.md
# Lines 4198-4477 → docs/guide/best-practices.md
```

**Phase 2: Add Navigation Links**

Each new file should have:
- **Top navigation**: Links to index and adjacent sections
- **Bottom navigation**: Previous/Next links
- **Sidebar**: Auto-generated from mkdocs.yml

Example header template:
```markdown
# Section Title

<div class="navigation-top">
  <i class="material-icons">arrow_back</i>
  <a href="index.md">Guide Index</a> |
  <a href="previous.md">Previous</a> |
  <a href="next.md">Next</a>
</div>

---

[Content here]

---

## Navigation

**Previous**: [Previous Section](previous.md)
**Next**: [Next Section](next.md)
**Index**: [Guide Index](index.md)
```

**Phase 3: Update Links**

- Search for all internal references in extracted content
- Update relative links to new file locations
- Update anchor links to point to correct files

**Phase 4: Validation**

- Build with `mkdocs build --clean`
- Verify all links work
- Check navigation flow
- Verify no content lost

**Phase 5: Deprecate Original**

- Move `generated-project-guide.md` to `docs/_archive/`
- Add redirect notice in original location
- Update all external references

---

## Part B: Splitting `tutorial-exemple-complet.md`

### Document Analysis

**Total Lines**: 1585
**Current Structure**: 16 main steps (Étapes 1-13 + sections)
**Proposed Split**: 4 logical parts

### Section Breakdown with Line Numbers

```
Line    1-25   : Header + Table of Contents
Line   26-37   : Objectif
Line   38-55   : Prérequis
Line   57-84   : Étape 1 - Installation du CLI
Line   85-147  : Étape 2 - Génération du projet
Line  148-223  : Étape 3 - Configuration initiale
Line  224-311  : Étape 4 - Tester le projet de base
Line  312-401  : Étape 5 - Ajouter le domaine Post
Line  402-618  : Étape 6 - Implémenter le service Post
Line  619-730  : Étape 7 - Créer le repository Post
Line  731-1004 : Étape 8 - Créer le handler HTTP
Line 1005-1199 : Étape 9 - Enregistrer les routes
Line 1200-1317 : Étape 10 - Tester l'API Posts
Line 1318-1458 : Étape 11 - Ajouter le domaine Comment
Line 1459-1507 : Étape 12 - Tests unitaires
Line 1508-1535 : Étape 13 - Déploiement Docker
Line 1536-1585 : Conclusion
```

---

### Proposed Split Strategy

#### Option 1: Split by Tutorial Phase (RECOMMENDED)

Split into 4 progressive parts representing natural workflow stages.

##### `docs/tutorial/01-setup.md`
**Title**: "Partie 1: Installation et Configuration"
**Source**: Lines 1-311
**Size**: ~311 lines
**Time to complete**: ~15 minutes

**Content**:
- Objectif
- Prérequis
- Étape 1: Installation du CLI
- Étape 2: Génération du projet
- Étape 3: Configuration initiale
- Étape 4: Tester le projet de base

**Checkpoint**: Project running, auth working

---

##### `docs/tutorial/02-first-domain.md`
**Title**: "Partie 2: Créer votre premier domaine (Posts)"
**Source**: Lines 312-730
**Size**: ~418 lines
**Time to complete**: ~30 minutes

**Content**:
- Étape 5: Ajouter le domaine Post (Article)
- Étape 6: Implémenter le service Post
- Étape 7: Créer le repository Post

**Checkpoint**: Post domain implemented (service + repo)

---

##### `docs/tutorial/03-api-integration.md`
**Title**: "Partie 3: Exposer l'API HTTP"
**Source**: Lines 731-1317
**Size**: ~586 lines
**Time to complete**: ~30 minutes

**Content**:
- Étape 8: Créer le handler HTTP
- Étape 9: Enregistrer les routes et le module
- Étape 10: Tester l'API Posts
- Étape 11: Ajouter le domaine Comment (as exercise)

**Checkpoint**: Full CRUD API working, tested with curl

---

##### `docs/tutorial/04-testing-deployment.md`
**Title**: "Partie 4: Tests et Déploiement"
**Source**: Lines 1459-1585
**Size**: ~126 lines
**Time to complete**: ~20 minutes

**Content**:
- Étape 12: Tests unitaires
- Étape 13: Déploiement Docker
- Conclusion
- Prochaines étapes

**Checkpoint**: Tests passing, Docker deployed

---

##### `docs/tutorial/index.md`
**Title**: "Tutorial: Créer une API Blog complète"
**Source**: New content + summary
**Size**: ~100 lines

**Content**:
```markdown
# Tutorial: Créer une API Blog complète avec create-go-starter

Guide pas-à-pas pour créer une API Blog avec `create-go-starter`,
de l'installation au déploiement.

## Structure du tutorial

Ce tutorial est divisé en 4 parties progressives:

### Partie 1: Installation et Configuration (15 min)
<i class="material-icons info">developer_board</i>
**Objectif**: Installer le CLI, générer le projet, le configurer et le tester

**Vous allez apprendre**:
- Installer `create-go-starter`
- Générer un nouveau projet
- Configurer PostgreSQL et JWT
- Tester l'API d'authentification générée

[Commencer la Partie 1](01-setup.md)

---

### Partie 2: Créer votre premier domaine (30 min)
<i class="material-icons success">architecture</i>
**Objectif**: Implémenter le domaine Posts (articles de blog)

**Vous allez apprendre**:
- Créer une entité GORM
- Implémenter un service métier
- Créer un repository avec GORM
- Appliquer l'architecture hexagonale

[Commencer la Partie 2](02-first-domain.md)

---

### Partie 3: Exposer l'API HTTP (30 min)
<i class="material-icons warning">http</i>
**Objectif**: Créer les endpoints CRUD pour Posts

**Vous allez apprendre**:
- Créer un handler HTTP avec Fiber
- Enregistrer les routes
- Intégrer avec fx (DI)
- Tester l'API avec curl

[Commencer la Partie 3](03-api-integration.md)

---

### Partie 4: Tests et Déploiement (20 min)
<i class="material-icons">rocket_launch</i>
**Objectif**: Ajouter des tests et déployer avec Docker

**Vous allez apprendre**:
- Écrire des tests unitaires
- Tester avec des mocks
- Créer une image Docker
- Déployer avec docker-compose

[Commencer la Partie 4](04-testing-deployment.md)

---

## Prérequis

[Same as current]

## Temps total estimé

**95 minutes** (~1h30) pour compléter l'ensemble du tutorial.

Vous pouvez faire des pauses entre les parties - chaque partie
a un checkpoint clair.

## Navigation

- [Partie 1: Installation et Configuration](01-setup.md)
- [Partie 2: Créer votre premier domaine](02-first-domain.md)
- [Partie 3: Exposer l'API HTTP](03-api-integration.md)
- [Partie 4: Tests et Déploiement](04-testing-deployment.md)
```

---

#### Option 2: Split by Domain (Alternative)

Keep single tutorial, but with better anchors and TOC.

**NOT RECOMMENDED** - less modular, harder to navigate.

---

### Implementation Strategy for Part B

**Phase 1: Create Tutorial Directory**

```bash
mkdir -p docs/tutorial
```

**Phase 2: Create Index**

Create `docs/tutorial/index.md` with overview and navigation.

**Phase 3: Extract Content**

```bash
# Lines 1-311 → docs/tutorial/01-setup.md
# Lines 312-730 → docs/tutorial/02-first-domain.md
# Lines 731-1317 → docs/tutorial/03-api-integration.md
# Lines 1459-1585 → docs/tutorial/04-testing-deployment.md
```

**Phase 4: Add Navigation**

Each part should have:
- Progress indicator (e.g., "Partie 2/4")
- Link to tutorial index
- Previous/Next links
- Checkpoint summary at the end

**Phase 5: Update Links**

- Fix internal cross-references
- Update code block titles to show current file context

**Phase 6: Deprecate Original**

- Move to `docs/_archive/tutorial-exemple-complet.md`
- Add redirect in old location

---

## Part C: `mkdocs.yml` Navigation Structure

### Proposed Hierarchical Navigation

```yaml
nav:
  # Getting Started
  - Accueil: index.md
  - Installation: installation.md
  - Quick Start: quick-start.md  # NEW
  - Utilisation: usage.md

  # Databases
  - Bases de données:
      - Guide de sélection: databases.md
      - Guide de migration: database-migration.md

  # Tutorial
  - Tutorial:
      - tutorial/index.md
      - Partie 1 - Installation: tutorial/01-setup.md
      - Partie 2 - Premier domaine: tutorial/02-first-domain.md
      - Partie 3 - API HTTP: tutorial/03-api-integration.md
      - Partie 4 - Tests & Deploy: tutorial/04-testing-deployment.md

  # Complete Guide
  - Guide complet:
      - guide/index.md
      - Architecture: guide/architecture.md
      - Structure du projet: guide/structure.md
      - Configuration: guide/configuration.md
      - Développement: guide/development.md
      - Exemples pratiques: guide/examples.md
      - API Reference: guide/api-reference.md
      - Tests: guide/testing.md
      - Base de données: guide/database.md
      - Sécurité: guide/security.md
      - Déploiement: guide/deployment.md
      - Monitoring: guide/monitoring.md
      - Bonnes pratiques: guide/best-practices.md

  # Reference
  - Référence:
      - Glossaire: reference/glossary.md      # NEW
      - FAQ: reference/faq.md                  # NEW

  # Development
  - Développement:
      - Architecture CLI: cli-architecture.md
      - Contribuer: contributing.md
```

### Navigation Benefits

**Logical Grouping**:
- Getting Started: Quick paths for new users
- Tutorial: Progressive learning path
- Guide complet: Deep reference material
- Référence: Quick lookups
- Développement: For contributors

**Progressive Disclosure**:
- Beginners → Quick Start, Tutorial
- Intermediate → Guide sections
- Advanced → Complete guide, CLI architecture
- Contributors → Contributing, CLI architecture

**Improved Hierarchy**:
- Top-level categories clear
- Subsections well-organized
- Logical progression

---

### Navigation Features to Enable

Update `mkdocs.yml` theme features:

```yaml
theme:
  features:
    - navigation.sections       # Already enabled
    - navigation.top           # Already enabled
    - navigation.indexes       # Already enabled
    - toc.follow              # Already enabled
    - search.suggest          # Already enabled
    - search.highlight        # Already enabled
    - content.code.copy       # Already enabled
    - navigation.tabs         # NEW - Top-level tabs
    - navigation.tabs.sticky  # NEW - Sticky tabs on scroll
    - navigation.expand       # NEW - Expand all sections by default
    - navigation.path         # NEW - Show breadcrumb path
```

---

## Part D: Quick Start Guide Content Outline

### File: `docs/quick-start.md`

**Objective**: Get users from zero to running project in 5 minutes
**Target audience**: Developers who want to evaluate the tool quickly
**Format**: Ultra-concise, copy-paste friendly

---

### Proposed Structure

```markdown
# Quick Start (5 minutes)

<i class="material-icons success">bolt</i> Générez et lancez
un projet Go production-ready en 5 minutes.

---

## 1. Installer le CLI (30 secondes)

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

**Vérification**:
```bash
create-go-starter --help
```

---

## 2. Créer le projet (10 secondes)

```bash
create-go-starter my-api
cd my-api
```

<i class="material-icons success small">check</i> **Généré**:
~45 fichiers avec architecture hexagonale complète

---

## 3. Setup automatique (2 minutes)

```bash
./setup.sh
```

**Ce script fait**:
- <i class="material-icons success small">check</i> Installe les dépendances Go
- <i class="material-icons success small">check</i> Génère un JWT secret
- <i class="material-icons success small">check</i> Configure PostgreSQL (Docker)
- <i class="material-icons success small">check</i> Lance les migrations

---

## 4. Lancer l'application (5 secondes)

```bash
make run
```

**Console output**:
```
INFO  Server starting on :8080
INFO  Database connected successfully
INFO  Migrations applied: 2
```

---

## 5. Tester l'API (1 minute)

### Health check

```bash
curl http://localhost:8080/health
```

**Response**:
```json
{"status":"healthy","timestamp":"2026-02-12T10:30:00Z"}
```

### Créer un utilisateur

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "SecurePass123!"
  }'
```

**Response**:
```json
{
  "user": {
    "id": 1,
    "email": "demo@example.com"
  },
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc..."
}
```

### Accéder à votre profil (protégé)

```bash
TOKEN="<your_access_token>"

curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

**Response**:
```json
{
  "id": 1,
  "email": "demo@example.com",
  "created_at": "2026-02-12T10:32:00Z"
}
```

---

## <i class="material-icons success">celebration</i> Félicitations!

Vous avez maintenant une API REST complète avec:

- <i class="material-icons success small">check</i> **Architecture hexagonale** (Ports & Adapters)
- <i class="material-icons success small">check</i> **Authentification JWT** (access + refresh tokens)
- <i class="material-icons success small">check</i> **Base de données** (PostgreSQL + GORM)
- <i class="material-icons success small">check</i> **Validation** (go-playground/validator)
- <i class="material-icons success small">check</i> **Logging** (zerolog)
- <i class="material-icons success small">check</i> **Dependency Injection** (uber-go/fx)
- <i class="material-icons success small">check</i> **Tests** (structure prête)
- <i class="material-icons success small">check</i> **Docker** (Dockerfile + compose)

---

## Prochaines étapes

### <i class="material-icons">menu_book</i> Apprendre

- **[Tutorial complet](tutorial/index.md)** - Créer une API Blog (1h30)
- **[Guide complet](guide/index.md)** - Architecture et patterns
- **[Architecture](guide/architecture.md)** - Comprendre l'hexagonale

### <i class="material-icons">code</i> Développer

- **[Ajouter un domaine](guide/development.md#add-model)** - Utiliser `add-model`
- **[API Reference](guide/api-reference.md)** - Tous les endpoints
- **[Tests](guide/testing.md)** - Stratégies de tests

### <i class="material-icons">rocket_launch</i> Déployer

- **[Docker](guide/deployment.md#docker)** - Containerisation
- **[Kubernetes](guide/deployment.md#kubernetes)** - Orchestration
- **[CI/CD](guide/deployment.md#cicd)** - Automatisation

---

## Besoin d'aide?

- <i class="material-icons info">help</i> **[FAQ](reference/faq.md)** - Questions fréquentes
- <i class="material-icons">bug_report</i> **[Issues](https://github.com/tky0065/go-starter-kit/issues)** - Reporter un bug
- <i class="material-icons">forum</i> **[Discussions](https://github.com/tky0065/go-starter-kit/discussions)** - Poser une question
```

---

### Quick Start Features

**Ultra-Concise**:
- No theory, only actions
- Copy-paste friendly commands
- Expected outputs shown
- Quick verification steps

**Progressive**:
- 5 clear steps
- Time estimates for each
- Success indicators
- Clear completion message

**Action-Oriented**:
- Every section has executable commands
- No reading, just doing
- Immediate feedback
- Working result in 5 minutes

**Gateway to Deeper Content**:
- Links to tutorial for learning
- Links to guide for reference
- Links to deployment for production

---

## Implementation Timeline

### Week 1: Guide Split (Part A)
- **Day 1-2**: Create all 13 guide files
- **Day 3**: Add navigation links
- **Day 4**: Update internal links
- **Day 5**: Validation and testing

### Week 2: Tutorial Split (Part B) + Quick Start (Part D)
- **Day 1-2**: Split tutorial into 4 parts
- **Day 3**: Create quick-start.md
- **Day 4**: Add navigation and links
- **Day 5**: Validation and testing

### Week 3: Navigation & Finalization (Part C)
- **Day 1-2**: Update mkdocs.yml navigation
- **Day 3**: Add glossary and FAQ (skeleton)
- **Day 4**: Final validation
- **Day 5**: Deploy and documentation

---

## Risk Mitigation

### Risks

1. **Content Loss**: Accidentally delete content during split
2. **Broken Links**: Internal references break
3. **User Confusion**: Old URLs stop working
4. **SEO Impact**: URL changes affect search rankings

### Mitigation Strategies

**Content Loss**:
- Keep originals in `_archive/` until fully validated
- Use version control (git) for every change
- Create validation checklist
- Word count comparison (before/after)

**Broken Links**:
- Use search to find all internal references
- Create link mapping document
- Test all links with `mkdocs build --strict`
- Use mkdocs validation plugins

**User Confusion**:
- Add redirect notices in old locations
- Keep old URLs working with redirect pages
- Add migration notice in release notes
- Update external references (README, etc.)

**SEO Impact**:
- Implement proper redirects
- Update sitemap
- Submit new URLs to search engines
- Monitor analytics for 404s

---

## Success Metrics

### Quantitative

- **Page Load Time**: < 2s for any page
- **Navigation Depth**: Max 3 clicks to any content
- **Search Effectiveness**: Top 3 results relevant
- **Mobile Usability**: 100/100 on PageSpeed Insights
- **Broken Links**: 0

### Qualitative

- **User Feedback**: Easier to find specific topics
- **Contributor Experience**: Easier to update specific sections
- **Learning Path**: Clear progression from quick start to advanced
- **Maintenance**: Faster to update individual sections

---

## Validation Checklist

### Pre-Split

- [ ] Backup current documentation
- [ ] Count total lines/words
- [ ] List all internal links
- [ ] Document current structure

### Post-Split

- [ ] All content accounted for (word count match)
- [ ] All internal links working
- [ ] Navigation flows logically
- [ ] Mobile rendering correct
- [ ] Search returns relevant results
- [ ] Build completes without errors (`mkdocs build --strict`)
- [ ] All pages have navigation links
- [ ] All code blocks render correctly
- [ ] All Material Icons render correctly
- [ ] No 404s in local testing

### Post-Deployment

- [ ] All URLs accessible
- [ ] Redirects working
- [ ] Search engine indexed
- [ ] Analytics tracking
- [ ] User feedback collected

---

## Appendix A: Line Extraction Commands

### For generated-project-guide.md

```bash
# Extract specific line ranges
sed -n '1,19p' docs/generated-project-guide.md > docs/guide/index-content.md
sed -n '20,289p' docs/generated-project-guide.md > docs/guide/architecture.md
sed -n '290,925p' docs/generated-project-guide.md > docs/guide/structure.md
sed -n '926,1080p' docs/generated-project-guide.md > docs/guide/configuration.md
sed -n '1081,1606p' docs/generated-project-guide.md > docs/guide/development.md
sed -n '1607,2463p' docs/generated-project-guide.md > docs/guide/examples.md
sed -n '2465,2780p' docs/generated-project-guide.md > docs/guide/api-reference.md
sed -n '2782,3011p' docs/generated-project-guide.md > docs/guide/testing.md
sed -n '3013,3214p' docs/generated-project-guide.md > docs/guide/database.md
sed -n '3216,3535p' docs/generated-project-guide.md > docs/guide/security.md
sed -n '3537,3984p' docs/generated-project-guide.md > docs/guide/deployment.md
sed -n '3986,4196p' docs/generated-project-guide.md > docs/guide/monitoring.md
sed -n '4198,4477p' docs/generated-project-guide.md > docs/guide/best-practices.md
```

### For tutorial-exemple-complet.md

```bash
# Extract tutorial parts
sed -n '1,311p' docs/tutorial-exemple-complet.md > docs/tutorial/01-setup.md
sed -n '312,730p' docs/tutorial-exemple-complet.md > docs/tutorial/02-first-domain.md
sed -n '731,1317p' docs/tutorial-exemple-complet.md > docs/tutorial/03-api-integration.md
sed -n '1459,1585p' docs/tutorial-exemple-complet.md > docs/tutorial/04-testing-deployment.md
```

---

## Appendix B: Navigation Template

### Standard Page Header Template

```markdown
# Page Title

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

[Page content]

---

## Navigation

<div class="page-navigation">
  <div class="nav-previous">
    <i class="material-icons">navigate_before</i>
    <a href="previous.md">Previous: Section Name</a>
  </div>
  <div class="nav-next">
    <a href="next.md">Next: Section Name</a>
    <i class="material-icons">navigate_next</i>
  </div>
</div>

<div class="nav-index">
  <a href="index.md"><i class="material-icons">home</i> Back to Guide Index</a>
</div>
```

### CSS for Navigation (add to extra.css)

```css
/* Page Navigation */
.navigation {
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--md-default-fg-color--lightest);
}

.page-navigation {
  display: flex;
  justify-content: space-between;
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 2px solid var(--md-default-fg-color--lightest);
}

.nav-previous,
.nav-next {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-index {
  text-align: center;
  margin-top: 1rem;
}
```

---

## Appendix C: Redirect Page Template

For old URLs that need to redirect:

```markdown
# Page Moved

<i class="material-icons warning">info</i> **This page has been reorganized.**

The content previously on this page has been split into multiple focused guides:

## New Locations

- **Architecture** → [guide/architecture.md](guide/architecture.md)
- **Structure** → [guide/structure.md](guide/structure.md)
- **Configuration** → [guide/configuration.md](guide/configuration.md)
- **Development** → [guide/development.md](guide/development.md)
- **Examples** → [guide/examples.md](guide/examples.md)
- **API Reference** → [guide/api-reference.md](guide/api-reference.md)
- **Tests** → [guide/testing.md](guide/testing.md)
- **Database** → [guide/database.md](guide/database.md)
- **Security** → [guide/security.md](guide/security.md)
- **Deployment** → [guide/deployment.md](guide/deployment.md)
- **Monitoring** → [guide/monitoring.md](guide/monitoring.md)
- **Best Practices** → [guide/best-practices.md](guide/best-practices.md)

## Why the Change?

The complete guide has been restructured for:
- Faster navigation to specific topics
- Better organization and discoverability
- Easier maintenance and updates
- Improved learning experience

**Start here**: [Guide Index](guide/index.md)
```

---

## Conclusion

This plan provides a comprehensive roadmap for restructuring the documentation in Sprint 2. The modular approach will:

1. **Improve user experience** with faster navigation and focused content
2. **Reduce maintenance burden** by isolating updates to specific files
3. **Enable progressive learning** from quick start to deep dives
4. **Enhance discoverability** with clear hierarchical navigation
5. **Future-proof the docs** for easier expansion and localization

**Next Steps**:
1. Review and approve this plan
2. Begin implementation with Part A (guide split)
3. Proceed to Part B (tutorial split) and Part D (quick start)
4. Finalize with Part C (navigation update)
5. Validate, test, and deploy

---

**Plan Status**: READY FOR REVIEW
**Estimated Implementation Time**: 3 weeks
**Risk Level**: LOW (with mitigation strategies)
**Impact**: HIGH (significant UX improvement)
