# Roadmap - go-starter-kit

Ce document présente la vision et les prochaines étapes pour l'évolution de `create-go-starter`.

## 🎉 Version Actuelle: v1.0.0 (MVP Complete)

**Date de release**: 15 janvier 2026  
**Statut**: ✅ Production Ready

### Fonctionnalités Disponibles

- ✅ 3 templates de projet (minimal, full, graphql)
- ✅ Architecture hexagonale
- ✅ JWT Authentication (access + refresh tokens)
- ✅ User CRUD complet
- ✅ API REST avec Fiber v2
- ✅ API GraphQL avec gqlgen
- ✅ PostgreSQL + GORM
- ✅ Swagger/OpenAPI docs
- ✅ Docker multi-stage optimisé
- ✅ GitHub Actions CI/CD
- ✅ Tests complets (unitaires + intégration)
- ✅ Initialisation Git automatique
- ✅ Installation dépendances automatique

**Métriques**:
- 6/6 Epics complétées
- 26/26 Exigences fonctionnelles satisfaites
- 13/13 Exigences non-fonctionnelles validées

---

## 🚀 Growth Features (Post-MVP)

Fonctionnalités planifiées pour élargir la base d'utilisateurs après validation du MVP.

### Epic 7: Support Multi-Base de Données 🗄️

**Objectif**: Permettre aux utilisateurs de choisir leur base de données préférée.

**Priority**: High  
**Estimated Effort**: 3-4 semaines  
**Target Release**: v1.1.0

#### Stories Potentielles

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

**Acceptance Criteria**:
- [ ] Utilisateur peut spécifier `--database=postgres|mysql|sqlite`
- [ ] Tous les templates fonctionnent avec chaque DB
- [ ] Documentation complète pour chaque DB
- [ ] Tests E2E passent pour toutes les DB

---

### Epic 8: CRUD Scaffolding Generator 🏗️

**Objectif**: Générer automatiquement du code CRUD pour de nouveaux modèles.

**Priority**: High  
**Estimated Effort**: 4-5 semaines  
**Target Release**: v1.2.0

#### Stories Potentielles

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

**Acceptance Criteria**:
- [ ] Commande `add-model` fonctionne dans projet existant
- [ ] Code généré compile et tests passent
- [ ] Swagger mis à jour automatiquement
- [ ] Support relations (one-to-many, many-to-many)

---

### Epic 9: Observabilité Avancée 📊

**Objectif**: Ajouter monitoring et observabilité pour projets production.

**Priority**: Medium  
**Estimated Effort**: 3-4 semaines  
**Target Release**: v1.3.0

#### Stories Potentielles

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

**Acceptance Criteria**:
- [ ] Flag `--observability=basic|advanced`
- [ ] Metrics Prometheus exposés
- [ ] Distributed tracing fonctionnel
- [ ] Dashboard Grafana importable
- [ ] Documentation complète

---

### Epic 10: Support Multi-Framework 🎭

**Objectif**: Supporter d'autres frameworks Go populaires (Gin, Echo).

**Priority**: Low  
**Estimated Effort**: 5-6 semaines  
**Target Release**: v2.0.0

#### Stories Potentielles

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

## 🔮 Vision Long-Terme (Future)

Transformation en plateforme écosystémique communautaire.

### Epic 11: Plugin System & Marketplace

**Objectif**: Permettre à la communauté de créer et partager des plugins.

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

### Epic 12: Interface Web/Dashboard

**Objectif**: Interface graphique pour créer et gérer projets.

**Features**:
- Web UI pour configuration projet (alternative au CLI)
- Visualisation de l'architecture générée
- Code preview avant génération
- Dashboard pour gérer projets multiples
- Metrics et monitoring intégrés

**Target**: v3.x

---

### Epic 13: Cloud Deployment Automation

**Objectif**: Déploiement one-click vers cloud providers.

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

## 📋 Backlog d'Améliorations Mineures

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

## 🗳️ Community Feedback

Nous écoutons activement la communauté! Si vous avez des idées ou suggestions:

1. **GitHub Discussions**: https://github.com/tky0065/go-starter-kit/discussions
2. **GitHub Issues**: https://github.com/tky0065/go-starter-kit/issues/new
3. **Feature Requests**: Utilisez le label `enhancement`

---

## 📊 Métriques de Succès

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

## 🤝 Comment Contribuer

Vous souhaitez contribuer à ces fonctionnalités? Consultez:

- [Guide de contribution](./docs/contributing.md)
- [Architecture du CLI](./docs/cli-architecture.md)
- [Issues "good first issue"](https://github.com/tky0065/go-starter-kit/labels/good%20first%20issue)

---

**Dernière mise à jour**: 15 janvier 2026  
**Version du document**: 1.0
