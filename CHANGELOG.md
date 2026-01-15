# Changelog

Toutes les modifications notables de ce projet seront documentées dans ce fichier.

Le format est basé sur [Keep a Changelog](https://keepachangelog.com/fr/1.0.0/),
et ce projet adhère au [Semantic Versioning](https://semver.org/lang/fr/).

## [1.0.0] - 2026-01-15

### 🎉 MVP Complete - Production Ready

Premier release officiel de `create-go-starter`, un générateur CLI pour créer des projets Go prêts pour la production avec architecture hexagonale.

### ✨ Fonctionnalités Ajoutées

#### Templates de Projet (Epic 6)
- **3 templates au choix** via le flag `--template`:
  - `minimal` - API REST basique avec Swagger (sans authentification) - ~20 fichiers
  - `full` - API complète avec JWT auth et gestion utilisateurs (défaut) - ~35 fichiers
  - `graphql` - API GraphQL avec gqlgen et GraphQL Playground - ~23 fichiers

#### Architecture & Stack Technique (Epics 1-5)
- **Architecture hexagonale** (Ports & Adapters) pour séparation claire des responsabilités
- **JWT Authentication** (Epic 2):
  - Access tokens + Refresh tokens avec rotation sécurisée
  - Middleware de sécurisation des routes
  - Gestion de session avec renouvellement automatique
- **User CRUD** (Epic 3):
  - Opérations complètes (Create, Read, Update, Delete)
  - Gestion du profil utilisateur
  - Hachage sécurisé des mots de passe (bcrypt)
- **API REST** avec Fiber v2 - Framework web haute performance
- **Base de données** PostgreSQL avec GORM et migrations automatiques
- **Injection de dépendances** avec uber-go/fx
- **Logging structuré** avec rs/zerolog
- **Validation** avec go-playground/validator

#### Documentation & API (Epic 4)
- **Swagger/OpenAPI** - Documentation auto-générée avec swaggo/swag
- **Standardisation des API** - Format de réponse uniforme
- **Gestion centralisée des erreurs** - Codes d'erreur standardisés

#### DevOps & Qualité (Epics 4-5)
- **Docker**:
  - Build multi-stage optimisé
  - docker-compose pré-configuré pour dev
  - Image de production légère basée sur Alpine
- **CI/CD**:
  - Pipeline GitHub Actions pré-configuré
  - Lint automatique avec golangci-lint
  - Tests automatisés
- **Tests**:
  - Tests unitaires pour handlers, services, repositories
  - Tests d'intégration
  - Couverture de tests du CLI
  - Smoke tests pour validation finale
  - 8 tests de résolveurs GraphQL (template graphql)
- **Makefile** avec commandes utiles (dev, test, build, docker)

#### Automatisation (Epic 5)
- **Initialisation Git automatique** avec commit initial
- **Installation automatique des dépendances** Go (`go mod tidy`)
- **Script setup.sh** pour configuration automatique du projet
- **Documentation inline** avec GoDoc pour toutes les fonctions publiques

### 📊 Métriques de Qualité

- ✅ **26/26** exigences fonctionnelles satisfaites (100%)
- ✅ **13/13** exigences non-fonctionnelles validées (100%)
- ✅ **6/6** epics complétées
- ✅ **26** user stories implémentées
- ✅ **100%** de couverture des acceptance criteria
- ✅ Validation end-to-end réussie pour tous les templates

### 📚 Documentation

- Guide d'installation complet
- Guide d'utilisation avec exemples
- Documentation de l'architecture CLI
- Guide du projet généré
- Quick start guide
- GitHub Pages: https://tky0065.github.io/go-starter-kit/

### 🔧 Configuration Requise

- Go 1.23+ (recommandé: 1.25.5)
- PostgreSQL 12+ (ou Docker pour exécution via conteneur)
- Git (optionnel, pour initialisation automatique)

### 📦 Installation

```bash
# Installation globale depuis GitHub
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@v1.0.0

# Ou version latest
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

### 🚀 Utilisation

```bash
# Template par défaut (full)
create-go-starter mon-projet

# Template minimal
create-go-starter --template=minimal mon-projet

# Template GraphQL
create-go-starter --template=graphql mon-projet
```

### 🐛 Corrections de Bugs

- Fix: Ajout des imports manquants (`fmt`, `time`) dans le template de tests GraphQL
- Fix: Gestion correcte du flag `--template` (nécessite `--template=value` ou position avant le nom)

### 🔒 Sécurité

- Tokens JWT sécurisés avec expiration configurable
- Refresh tokens avec rotation automatique
- Hachage bcrypt pour les mots de passe
- Validation stricte des entrées utilisateur
- Configuration des secrets via variables d'environnement

### 📁 Epics Complétées

1. **Epic 1** - CLI Generator Base
   - Installation de l'outil CLI
   - Génération de la structure de base
   - Injection dynamique du contexte projet
   - Initialisation du serveur Fiber + DI fx + DB
   - Environnement de développement (.env, Makefile, Docker)

2. **Epic 2** - JWT Authentication
   - Inscription des utilisateurs
   - Authentification (login/logout)
   - Renouvellement de session
   - Sécurisation des routes

3. **Epic 3** - User CRUD
   - Gestion du profil utilisateur
   - Opérations CRUD utilisateur

4. **Epic 4** - API, Errors, Swagger, CI/CD
   - Standardisation des API
   - Gestion centralisée des erreurs
   - Documentation interactive Swagger
   - Automatisation de la qualité
   - Intégration continue

5. **Epic 5** - Git auto, Tests, Docker, Smoke tests
   - Initialisation Git automatique
   - Installation automatique des dépendances Go
   - Amélioration de la couverture de tests du CLI
   - Optimisation de l'image Docker générée
   - Documentation des fonctions publiques
   - Validation finale et smoke tests

6. **Epic 6** - Templates Multiples
   - Flag CLI pour sélection de template
   - Template minimal (API REST basique)
   - Refactoring du template full (API complète)
   - Template GraphQL avec gqlgen
   - Documentation et aide CLI

### 🙏 Remerciements

Merci à tous les contributeurs et aux projets open-source utilisés :
- [Fiber](https://github.com/gofiber/fiber) - Framework web
- [fx](https://github.com/uber-go/fx) - Injection de dépendances
- [GORM](https://gorm.io/) - ORM
- [zerolog](https://github.com/rs/zerolog) - Logging
- [swaggo](https://github.com/swaggo/swag) - Swagger
- [gqlgen](https://github.com/99designs/gqlgen) - GraphQL

---

## Format du Versioning

- **MAJOR** (X.0.0): Changements incompatibles avec les versions précédentes
- **MINOR** (1.X.0): Ajout de fonctionnalités rétro-compatibles
- **PATCH** (1.0.X): Corrections de bugs rétro-compatibles

[1.0.0]: https://github.com/tky0065/go-starter-kit/releases/tag/v1.0.0
