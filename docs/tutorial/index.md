# Tutorial: Créer une API Blog complète avec create-go-starter

Guide pas-à-pas pour créer une API Blog avec `create-go-starter`, de l'installation au déploiement.

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

### Logiciels requis

- **Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **PostgreSQL** ou **Docker** - Pour la base de données
- **curl** ou **Postman** - Pour tester l'API
- Éditeur de code (VS Code, GoLand, etc.)

### Connaissances recommandées

- Bases de Go (structs, interfaces, error handling)
- Concepts REST API
- Familiarité avec SQL/PostgreSQL (basique)

Pas besoin d'être expert! Ce tutorial explique chaque étape en détail.

## Temps total estimé

**95 minutes** (~1h30) pour compléter l'ensemble du tutorial.

Vous pouvez faire des pauses entre les parties - chaque partie a un checkpoint clair.

## Navigation

- [Partie 1: Installation et Configuration](01-setup.md)
- [Partie 2: Créer votre premier domaine](02-first-domain.md)
- [Partie 3: Exposer l'API HTTP](03-api-integration.md)
- [Partie 4: Tests et Déploiement](04-testing-deployment.md)
