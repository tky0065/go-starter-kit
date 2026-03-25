# Roadmap - go-starter-kit

Ce document présente la vision et les prochaines étapes pour l'évolution de `create-go-starter`.

## <i class="material-icons">rocket_launch</i> Version Actuelle: v1.5.2

**Date de release**: 23 février 2026  
**Statut**: <i class="material-icons success">check</i> Production Ready

### Fonctionnalités Disponibles

- <i class="material-icons success">check</i> 3 templates de projet (minimal, full, graphql)
- <i class="material-icons success">check</i> Architecture hexagonale
- <i class="material-icons success">check</i> JWT Authentication (access + refresh tokens)
- <i class="material-icons success">check</i> User CRUD complet
- <i class="material-icons success">check</i> API REST avec Fiber v2
- <i class="material-icons success">check</i> API GraphQL avec gqlgen
- <i class="material-icons success">check</i> PostgreSQL + GORM
- <i class="material-icons success">check</i> Swagger/OpenAPI docs
- <i class="material-icons success">check</i> Docker multi-stage optimisé
- <i class="material-icons success">check</i> GitHub Actions CI/CD
- <i class="material-icons success">check</i> Tests complets (unitaires + intégration)
- <i class="material-icons success">check</i> Initialisation Git automatique
- <i class="material-icons success">check</i> Installation dépendances automatique
- <i class="material-icons success">check</i> Support Multi-Base de Données (PostgreSQL, MySQL, SQLite)
- <i class="material-icons success">check</i> CRUD Scaffolding Generator (add-model command)
- <i class="material-icons success">check</i> Observabilité Avancée (Prometheus, Jaeger, Grafana, Health Checks K8s)
- <i class="material-icons success">check</i> Mode Interactif guidé (`--interactive` / `-i`)
- <i class="material-icons success">check</i> Prévisualisation Dry-Run (`--dry-run` / `-n`)
- <i class="material-icons success">check</i> Commande Doctor (diagnostics environnement)
- <i class="material-icons success">check</i> Barre de progression et statistiques de génération
- <i class="material-icons success">check</i> Alias courts pour tous les flags (`-t`, `-d`, `-o`, `-i`, `-n`, `-h`)

---

## <i class="material-icons">trending_up</i> Prochaines Fonctionnalités

Fonctionnalités planifiées pour les prochaines versions.

### ~~v1.1.0 - Support Multi-Base de Données~~ <i class="material-icons success">check</i> Complété

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

### ~~CRUD Scaffolding Generator~~ <i class="material-icons success">check</i> Complété (v1.2.0)

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

### ~~Observabilité Avancée~~ <i class="material-icons success">check</i> Complété (v1.3.0)

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

### ~~Améliorations CLI et Expérience Développeur~~ <i class="material-icons success">check</i> Complété (v1.4.0)

**Description**: Refonte complète de l'expérience utilisateur du CLI.

**Priorité**: Haute  
**Temps estimé**: 2-3 semaines  
**Version cible**: v1.4.0

#### Ce qui a été ajouté

1. **Mode Interactif (`--interactive` / `-i`)**
   - Assistant guidé étape par étape
   - Sélection interactive du template, database, observabilité
   - Résumé de configuration avec confirmation
   - Zéro dépendance externe (stdlib uniquement)

2. **Prévisualisation Dry-Run (`--dry-run` / `-n`)**
   - Preview des fichiers sans écriture sur disque
   - Compteur fichiers/répertoires
   - Compatible avec tous les flags

3. **Commande Doctor**
   - Diagnostics Go (>= 1.21), Git, Docker
   - Rapport clair avec statut de chaque outil
   - Code de sortie pour scripts CI

4. **Feedback Visuel**
   - Barre de progression pendant la génération
   - Statistiques post-génération (fichiers, taille, temps)
   - Désactivation automatique sur non-TTY / NO_COLOR

5. **Alias Courts**
   - `-t`, `-d`, `-o`, `-i`, `-n`, `-h`
   - Syntaxe flexible: `-t=minimal`, `-t minimal`
   - Détection des flags inconnus

**Objectifs**:
- [x] Mode interactif fonctionnel
- [x] Dry-run avec affichage structuré
- [x] Doctor vérifie Go, Git, Docker
- [x] Barre de progression et statistiques
- [x] Alias courts pour tous les flags
- [x] Tests complets pour les 5 fonctionnalités

---

### Support Multi-Framework

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

## <i class="material-icons">auto_awesome</i> Vision Long-Terme (Future)

Transformation en plateforme écosystémique communautaire.

### Plugin System & Marketplace

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

### Interface Web/Dashboard

**Description**: Interface graphique pour créer et gérer projets.

**Features**:
- Web UI pour configuration projet (alternative au CLI)
- Visualisation de l'architecture générée
- Code preview avant génération
- Dashboard pour gérer projets multiples
- Metrics et monitoring intégrés

**Target**: v3.x

---

### Cloud Deployment Automation

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

## <i class="material-icons">assignment</i> Backlog d'Améliorations Mineures

Améliorations continues pour versions patch.

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

- [x] Mode interactif pour sélection template (`create-go-starter --interactive`)
- [x] Flag `--dry-run` pour preview sans génération
- [ ] Flag `--update` pour mettre à jour projet existant
- [x] Commande `create-go-starter doctor` pour diagnostics
- [ ] Colored diff pour `--dry-run`
- [x] Progress bar pendant génération
- [x] Statistiques post-génération (fichiers créés, taille, etc.)

### Documentation

- [ ] Tutoriels vidéo (YouTube)
- [ ] Blog posts techniques
- [ ] Exemples de projets réels
- [ ] Best practices guide
- [ ] Migration guides (from scratch, from other starters)
- [ ] Troubleshooting guide
- [x] FAQ section
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

## <i class="material-icons">how_to_vote</i> Community Feedback

Nous écoutons activement la communauté! Si vous avez des idées ou suggestions:

1. **GitHub Discussions**: https://github.com/tky0065/go-starter-kit/discussions
2. **GitHub Issues**: https://github.com/tky0065/go-starter-kit/issues/new
3. **Feature Requests**: Utilisez le label `enhancement`

---

## <i class="material-icons">bar_chart</i> Métriques de Succès

### Objectifs 3 Mois (Mai 2026)

- [ ] 1,000+ installations du CLI
- [ ] 500+ étoiles GitHub
- [ ] 10+ contributors
- [ ] 50+ projets créés en production

### Objectifs 12 Mois (Février 2027)

- [ ] 5,000+ étoiles GitHub (Top 5 Go starters)
- [ ] 60% des utilisateurs créent 2+ projets
- [ ] NPS > 50
- [ ] Temps de réponse issues < 24h
- [ ] Communauté active (Discord, discussions)

---

## <i class="material-icons">handshake</i> Comment Contribuer

Vous souhaitez contribuer à ces fonctionnalités? Consultez:

- [Guide de contribution](./docs/contributing.md)
- [Architecture du CLI](./docs/cli-architecture.md)
- [Issues "good first issue"](https://github.com/tky0065/go-starter-kit/labels/good%20first%20issue)

---

**Dernière mise à jour**: 23 février 2026  
**Version du document**: 1.5.1
