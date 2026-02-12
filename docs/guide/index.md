# Guide des projets générés

> Guide complet pour développer, tester et déployer les projets créés avec create-go-starter

---

## Vue d'ensemble

Ce guide complet vous accompagne dans toutes les phases de développement d'un projet créé avec `create-go-starter`. Que vous débutiez ou que vous cherchiez des références approfondies, vous trouverez ici toute la documentation nécessaire.

Les projets générés suivent l'**architecture hexagonale** (Ports & Adapters) et intègrent les meilleures pratiques de développement Go.

## Navigation du guide

### Pour commencer

<i class="material-icons success">rocket_launch</i> **[Configuration](configuration.md)** - Variables d'environnement et setup initial

<i class="material-icons">folder</i> **[Structure du projet](structure.md)** - Organisation des fichiers et responsabilités

### Architecture et design

<i class="material-icons info">architecture</i> **[Architecture hexagonale](architecture.md)** - Principes, diagrammes et flux de requêtes

<i class="material-icons">code</i> **[Développement](development.md)** - Workflow quotidien et outils

<i class="material-icons">lightbulb</i> **[Exemples pratiques](examples.md)** - Créer une entité complète de A à Z

### Développement

<i class="material-icons">api</i> **[API Reference](api-reference.md)** - Documentation complète des endpoints

<i class="material-icons success">check_circle</i> **[Tests](testing.md)** - Organisation, types de tests et best practices

<i class="material-icons">storage</i> **[Base de données](database.md)** - Migrations, GORM et requêtes avancées

### Production

<i class="material-icons warning">shield</i> **[Sécurité](security.md)** - JWT, validation, hashage et checklist

<i class="material-icons">cloud_upload</i> **[Déploiement](deployment.md)** - Docker, Kubernetes et CI/CD

<i class="material-icons">monitor_heart</i> **[Monitoring](monitoring.md)** - Logging, health checks et observabilité

<i class="material-icons success">verified</i> **[Bonnes pratiques](best-practices.md)** - Patterns, conventions et recommandations

---

## Parcours d'apprentissage recommandé

### Débutant

Si vous découvrez `create-go-starter`, suivez ce parcours:

1. <i class="material-icons success small">circle</i> [Configuration](configuration.md) - Setup initial
2. <i class="material-icons success small">circle</i> [Structure du projet](structure.md) - Comprendre l'organisation
3. <i class="material-icons success small">circle</i> [Architecture](architecture.md) - Principes hexagonaux
4. <i class="material-icons success small">circle</i> [Exemples pratiques](examples.md) - Créer votre première entité
5. <i class="material-icons success small">circle</i> [API Reference](api-reference.md) - Tester l'API

### Intermédiaire

Pour approfondir vos connaissances:

1. <i class="material-icons info small">circle</i> [Développement](development.md) - Workflow avancé
2. <i class="material-icons info small">circle</i> [Tests](testing.md) - Stratégies de tests
3. <i class="material-icons info small">circle</i> [Base de données](database.md) - Maîtriser GORM
4. <i class="material-icons info small">circle</i> [Sécurité](security.md) - Bonnes pratiques

### Avancé

Pour la mise en production:

1. <i class="material-icons warning small">circle</i> [Déploiement](deployment.md) - Docker, K8s, CI/CD
2. <i class="material-icons warning small">circle</i> [Monitoring](monitoring.md) - Observabilité
3. <i class="material-icons warning small">circle</i> [Bonnes pratiques](best-practices.md) - Patterns avancés

---

## Navigation rapide

| Section | Contenu principal |
|---------|-------------------|
| [Architecture](architecture.md) | Hexagonal architecture, diagrammes, flux de requêtes |
| [Structure](structure.md) | Organisation des fichiers, stack technique, exemples de code |
| [Configuration](configuration.md) | Variables d'environnement, JWT, PostgreSQL |
| [Développement](development.md) | Workflow quotidien, Makefile, add-model |
| [Exemples](examples.md) | Créer l'entité Product complète (9 étapes) |
| [API Reference](api-reference.md) | Endpoints auth et users, workflow API |
| [Tests](testing.md) | Organisation, types de tests, coverage |
| [Base de données](database.md) | Migrations, GORM, queries avancées |
| [Sécurité](security.md) | JWT, validation, hashage, checklist |
| [Déploiement](deployment.md) | Docker, Kubernetes, GitHub Actions |
| [Monitoring](monitoring.md) | Logging, health checks, observabilité |
| [Bonnes pratiques](best-practices.md) | Architecture, code style, patterns |

---

## Prochaines étapes

<i class="material-icons">arrow_forward</i> **[Commencer par l'architecture](architecture.md)** - Comprendre les fondations hexagonales

<i class="material-icons">menu_book</i> **[Voir le tutorial complet](../tutorial/index.md)** - Guide pas-à-pas pour créer une API Blog

<i class="material-icons">home</i> **[Retour à l'accueil](../index.md)** - Documentation principale
