# Roadmap - go-starter-kit

Ce document présente la vision et les prochaines étapes pour l'évolution de `create-go-starter`.

## :material-rocket-launch: Version Actuelle: v1.3.0

**Date de release**: 17 février 2026  
**Statut**: :material-check: Production Ready

### Fonctionnalités Disponibles

- :material-check: 3 templates de projet (minimal, full, graphql)
- :material-check: Architecture hexagonale
- :material-check: JWT Authentication (access + refresh tokens)
- :material-check: User CRUD complet
- :material-check: API REST avec Fiber v2
- :material-check: API GraphQL avec gqlgen
- :material-check: PostgreSQL + GORM
- :material-check: Swagger/OpenAPI docs
- :material-check: Docker multi-stage optimisé
- :material-check: GitHub Actions CI/CD
- :material-check: Tests complets (unitaires + intégration)
- :material-check: Initialisation Git automatique
- :material-check: Installation dépendances automatique
- :material-check: Support Multi-Base de Données (PostgreSQL, MySQL, SQLite)
- :material-check: CRUD Scaffolding Generator (add-model command)
- :material-check: Observabilité Avancée (Prometheus, Jaeger, Grafana, Health Checks K8s)

---

## :material-trending-up: Prochaines Fonctionnalités

Fonctionnalités planifiées pour les prochaines versions.

### ~~v1.1.0 - Support Multi-Base de Données~~ :material-check: Complété

**Description**: Permettre aux utilisateurs de choisir leur base de données préférée.

**Priorité**: Haute  
**Temps estimé**: 3-4 semaines  
**Version cible**: v1.1.0

#### Ce qui sera ajouté

1. **Support MySQL/MariaDB**
   - Ajouter flag `--database=mysql`
   - Adapter templates avec driver MySQL
   - Documentation spécifique MySQL
   - Tests E2E avec MySQL

2. **Support SQLite**
   - Ajouter flag `--database=sqlite`
   - Configuration pour environnement de développement/test
   - Idéal pour prototypes rapides
   - Tests E2E avec SQLite

3. **Support MongoDB** (optionnel)
   - Ajouter flag `--database=mongodb`
   - Adapter architecture pour NoSQL
   - Driver mongo-go-driver
   - Documentation patterns NoSQL

**Objectifs**:
- [x] Utilisateur peut spécifier `--database=postgres|mysql|sqlite`
- [x] Tous les templates fonctionnent avec chaque DB
- [x] Documentation complète pour chaque DB
- [x] Tests E2E passent pour toutes les DB

---

### ~~CRUD Scaffolding Generator~~ :material-check: Complété (v1.2.0)

**Description**: Générer automatiquement du code CRUD pour de nouveaux modèles.

**Priorité**: Haute  
**Temps estimé**: 4-5 semaines  
**Version cible**: v1.2.0

#### Ce qui sera ajouté

1. **Commande `add-model`**
   - Sous-commande CLI `create-go-starter add-model <name>`
   - Parsing de définition de modèle (YAML ou interactif)
   - Génération fichiers: model, repository, service, handler, tests
   - Mise à jour routes automatique

2. **Templates de Modèles**
   - Support types courants (string, int, time, relations)
   - Validation automatique (required, min, max, email, etc.)
   - Génération tests unitaires pour nouveau modèle

3. **Documentation Auto-générée**
   - Mise à jour automatique Swagger annotations
   - Génération exemple requests/responses
   - Update README du projet

**Example Usage**:
```bash
cd mon-projet
create-go-starter add-model Todo --fields "title:string,completed:bool,dueDate:time"
# Génère: model, repository, service, handler, tests, swagger docs
```

**Objectifs**:
- [x] Commande `add-model` fonctionne dans projet existant
- [x] Code généré compile et tests passent
- [x] Swagger mis à jour automatiquement
- [x] Support relations (one-to-many, many-to-many)

---

### ~~Observabilité Avancée~~ :material-check: Complété (v1.3.0)

**Description**: Ajouter monitoring et observabilité pour projets production.

**Priorité**: Moyenne  
**Temps estimé**: 3-4 semaines  
**Version cible**: v1.3.0

#### Ce qui sera ajouté

1. **Prometheus Metrics**
   - Endpoint `/metrics` avec prometheus
   - Métriques HTTP (latence, status codes, throughput)
   - Métriques DB (connections, query time)
   - Métriques custom pour business logic

2. **Distributed Tracing**
   - OpenTelemetry integration
   - Trace propagation entre services
   - Export vers Jaeger/Zipkin
   - Correlation IDs dans logs

3. **Health Checks Avancés**
   - `/health/liveness` et `/health/readiness`
   - Checks pour DB, external services
   - Graceful shutdown
   - Kubernetes-ready

4. **Grafana Dashboard Template**
   - Dashboard pré-configuré pour projets générés
   - Visualisations clés (traffic, errors, latency)
   - Alerting rules

**Objectifs**:
- [x] Flag `--observability=basic|advanced`
- [x] Metrics Prometheus exposés
- [x] Distributed tracing fonctionnel
- [x] Dashboard Grafana importable
- [x] Documentation complète

---

###  Support Multi-Framework 🎭

**Description**: Supporter d'autres frameworks Go populaires (Gin, Echo).

**Priorité**: Basse  
**Temps estimé**: 5-6 semaines  
**Version cible**: v2.0.0

#### Ce qui sera ajouté

1. **Support Gin Framework**
   - Flag `--framework=gin`
   - Templates adaptés pour Gin
   - Middleware Gin
   - Documentation spécifique

2. **Support Echo Framework**
   - Flag `--framework=echo`
   - Templates adaptés pour Echo
   - Middleware Echo
   - Documentation spécifique

3. **Abstraction Framework-Agnostic**
   - Core business logic indépendant du framework
   - Adapters par framework
   - Migration guide entre frameworks

**Note**: Nécessite refactoring majeur, considéré pour v2.0.0

---

## :material-crystal-ball: Vision Long-Terme (Future)

Transformation en plateforme écosystémique communautaire.

###  Plugin System & Marketplace

**Description**: Permettre à la communauté de créer et partager des plugins.

**Features**:
- Architecture de plugins modulaire
- Registry de plugins communautaires
- CLI pour installer/gérer plugins: `create-go-starter plugin install <name>`
- Exemples de plugins:
  - Authentication providers (OAuth2, SAML, LDAP)
  - Payment integrations (Stripe, PayPal)
  - Cloud services (AWS, GCP, Azure)
  - Notification services (Email, SMS, Push)

**Target**: v2.x

---

###  Interface Web/Dashboard

**Description**: Interface graphique pour créer et gérer projets.

**Features**:
- Web UI pour configuration projet (alternative au CLI)
- Visualisation de l'architecture générée
- Code preview avant génération
- Dashboard pour gérer projets multiples
- Metrics et monitoring intégrés

**Target**: v3.x

---

###  Cloud Deployment Automation

**Description**: Déploiement one-click vers cloud providers.

**Features**:
- Commande `create-go-starter deploy --provider=aws|gcp|azure`
- Terraform/Pulumi templates auto-générés
- Kubernetes manifests optimisés
- Helm charts
- Integration avec cloud databases (RDS, Cloud SQL, etc.)
- Auto-scaling configuration
- CDN setup pour assets statiques

**Target**: v2.x

---

## :material-clipboard-list: Backlog d'Améliorations Mineures

Améliorations continues pour versions patch (v1.0.x, v1.1.x, etc.).

### Templates & Code Generation

- [ ] Template avec API versioning (v1, v2)
- [ ] Template microservices (avec gRPC)
- [ ] Template avec WebSockets support
- [ ] Template avec Server-Sent Events (SSE)
- [ ] Template avec file upload/download
- [ ] Template avec email service (SMTP)
- [ ] Template avec cache layer (Redis)
- [ ] Template avec message queue (RabbitMQ, Kafka)

### CLI Improvements

- [ ] Mode interactif pour sélection template (`create-go-starter --interactive`)
- [ ] Flag `--dry-run` pour preview sans génération
- [ ] Flag `--update` pour mettre à jour projet existant
- [ ] Commande `create-go-starter doctor` pour diagnostics
- [ ] Colored diff pour `--dry-run`
- [ ] Progress bar pendant génération
- [ ] Statistiques post-génération (fichiers créés, taille, etc.)

### Documentation

- [ ] Tutoriels vidéo (YouTube)
- [ ] Blog posts techniques
- [ ] Exemples de projets réels
- [ ] Best practices guide
- [ ] Migration guides (from scratch, from other starters)
- [ ] Troubleshooting guide
- [ ] FAQ section
- [ ] Architecture Decision Records (ADRs)

### Testing & Quality

- [ ] Fuzzing tests pour CLI
- [ ] Performance benchmarks
- [ ] Security scanning automatisé (Snyk, Dependabot)
- [ ] SAST/DAST pour code généré
- [ ] Compatibilité testing (Go versions, OS)
- [ ] Load testing templates
- [ ] Chaos engineering templates

### DevEx (Developer Experience)

- [ ] IDE extensions (VSCode, GoLand)
- [ ] GitHub Copilot integration
- [ ] Snippets pour patterns courants
- [ ] Makefile amélioré avec plus de commandes
- [ ] Pre-commit hooks configurés
- [ ] Git hooks pour tests automatiques
- [ ] Dev container configuration (.devcontainer)

### Community & Ecosystem

- [ ] Discord community server
- [ ] Monthly community calls
- [ ] Contributor recognition program
- [ ] Showcase page (projets utilisant le starter)
- [ ] Templates gallery
- [ ] Blog avec success stories
- [ ] Newsletter mensuelle

---

## :material-vote: Community Feedback

Nous écoutons activement la communauté! Si vous avez des idées ou suggestions:

1. **GitHub Discussions**: https://github.com/tky0065/go-starter-kit/discussions
2. **GitHub Issues**: https://github.com/tky0065/go-starter-kit/issues/new
3. **Feature Requests**: Utilisez le label `enhancement`

---

## :material-chart-bar: Métriques de Succès

### Objectifs 3 Mois (Avril 2026)

- [ ] 1,000+ installations du CLI
- [ ] 500+ étoiles GitHub
- [ ] 10+ contributors
- [ ] 50+ projets créés en production

### Objectifs 12 Mois (Janvier 2027)

- [ ] 5,000+ étoiles GitHub (Top 5 Go starters)
- [ ] 60% des utilisateurs créent 2+ projets
- [ ] NPS > 50
- [ ] Temps de réponse issues < 24h
- [ ] Communauté active (Discord, discussions)

---

## :material-handshake: Comment Contribuer

Vous souhaitez contribuer à ces fonctionnalités? Consultez:

- [Guide de contribution](./docs/contributing.md)
- [Architecture du CLI](./docs/cli-architecture.md)
- [Issues "good first issue"](https://github.com/tky0065/go-starter-kit/labels/good%20first%20issue)

---

**Dernière mise à jour**: 17 février 2026  
**Version du document**: 1.3
