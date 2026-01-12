# create-go-starter

Un outil CLI puissant pour générer des projets Go prêts pour la production en quelques secondes.

## Aperçu

`create-go-starter` est un générateur de projets Go qui crée une architecture hexagonale complète avec toutes les fonctionnalités essentielles d'une application backend moderne. En une seule commande, obtenez un projet structuré avec authentification JWT, API REST, base de données, tests, et configuration Docker prête pour le déploiement.

### Fonctionnalités incluses

- **Architecture hexagonale** (Ports & Adapters) - Séparation claire des responsabilités
- **Authentification JWT** - Access tokens + Refresh tokens avec rotation sécurisée
- **API REST** avec Fiber v2 - Framework web haute performance
- **Base de données** - GORM avec PostgreSQL et migrations automatiques
- **Injection de dépendances** - uber-go/fx pour une architecture modulaire
- **Tests complets** - Tests unitaires et d'intégration
- **Documentation Swagger** - API documentée automatiquement avec OpenAPI
- **Docker** - Build multi-stage optimisé et docker-compose
- **CI/CD** - Pipeline GitHub Actions pré-configuré
- **Logging structuré** - rs/zerolog pour des logs professionnels
- **Validation** - go-playground/validator pour valider les entrées
- **Makefile** - Commandes utiles pour dev, test, build et déploiement

## Installation rapide

### Méthode 1: Installation directe (Recommandée)

Installation globale en une seule commande, sans cloner le repository:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Le binaire sera installé dans `$GOPATH/bin` (généralement `~/go/bin`). Assurez-vous que ce répertoire est dans votre PATH.

**Note**: Si `create-go-starter` n'est pas reconnu après l'installation, ajoutez `$GOPATH/bin` à votre PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Méthode 2: Build depuis les sources

Pour contributeurs ou personnalisation:

```bash
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit
go build -o create-go-starter ./cmd/create-go-starter
# Le binaire est maintenant disponible: ./create-go-starter
```

### Méthode 3: Build avec Makefile

```bash
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit
make build
# Le binaire est disponible: ./create-go-starter
```

Pour plus de détails, consultez le [guide d'installation complet](./docs/installation.md).

## Utilisation de base

### Créer un nouveau projet

```bash
create-go-starter mon-super-projet
```

Cette commande va:
1. Créer la structure complète du projet
2. Générer ~30+ fichiers (handlers, services, repositories, tests, etc.)
3. Configurer tous les fichiers nécessaires (.env, Dockerfile, Makefile, etc.)
4. Copier le fichier `.env.example` vers `.env`

### Lancer le projet généré

#### Option 1: Configuration automatique (Recommandé) 🚀

```bash
cd mon-super-projet
./setup.sh
make run
```

Le script `setup.sh` automatise:
- Installation des dépendances Go
- Génération du JWT secret
- Configuration de PostgreSQL (Docker ou local)
- Vérification de l'installation

#### Option 2: Configuration manuelle

```bash
cd mon-super-projet

# Installer les dépendances et générer go.sum
go mod tidy

# Configurer le JWT secret dans .env
# JWT_SECRET=<générer avec: openssl rand -base64 32>

# Lancer la base de données (PostgreSQL)
docker run -d --name postgres \
  -e POSTGRES_DB=mon-super-projet \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Lancer l'application
make run
```

L'API sera disponible sur `http://localhost:8080`

```bash
# Tester le health check
curl http://localhost:8080/health
# {"status":"ok"}
```

## Structure générée

Voici ce que `create-go-starter` génère pour vous:

```
mon-super-projet/
├── cmd/
│   └── main.go                    # Point d'entrée avec fx dependency injection
├── internal/
│   ├── models/                    # Entités de domaine partagées
│   │   └── user.go                # User, RefreshToken, AuthResponse
│   ├── domain/                    # Couche domaine (logique métier)
│   │   ├── user/                  # Domaine User
│   │   │   ├── service.go         # Logique métier (Register, Login, etc.)
│   │   │   └── module.go          # Module fx
│   │   └── errors.go              # Erreurs métier personnalisées
│   ├── adapters/                  # Adapters (HTTP, DB)
│   │   ├── handlers/              # HTTP handlers
│   │   │   ├── auth_handler.go    # Auth endpoints (register, login, refresh)
│   │   │   └── user_handler.go    # User CRUD endpoints
│   │   ├── middleware/            # Middleware Fiber
│   │   │   ├── auth_middleware.go # JWT verification
│   │   │   └── error_handler.go   # Gestion centralisée des erreurs
│   │   ├── repository/            # Implémentation des repositories
│   │   │   └── user_repository.go # GORM implementation
│   │   └── http/                  # Routes health check
│   ├── infrastructure/            # Infrastructure
│   │   ├── database/              # Configuration DB (GORM, migrations)
│   │   └── server/                # Configuration Fiber app et routes
│   └── interfaces/                # Ports (interfaces)
│       └── user_repository.go     # Interface UserRepository
├── pkg/                           # Packages réutilisables
│   ├── auth/                      # JWT utilities
│   ├── config/                    # Chargement configuration (.env)
│   └── logger/                    # Configuration zerolog
├── .github/workflows/
│   └── ci.yml                     # Pipeline CI/CD (lint, test, build)
├── .env                           # Configuration (copié depuis .env.example)
├── .env.example                   # Template de configuration
├── .gitignore                     # Exclusions Git
├── .golangci.yml                  # Configuration du linter
├── Dockerfile                     # Build multi-stage pour production
├── Makefile                       # Commandes utiles (run, test, lint, docker, etc.)
├── setup.sh                       # Script de configuration automatique
├── go.mod                         # Module Go avec dépendances
└── README.md                      # Documentation du projet
```

Pour une explication détaillée de chaque composant, consultez le [guide d'utilisation](./docs/usage.md).

## Stack technique

Les projets générés utilisent les meilleures bibliothèques de l'écosystème Go:

| Composant | Bibliothèque | Version | Description |
|-----------|-------------|---------|-------------|
| Web Framework | [Fiber](https://gofiber.io/) | v2 | Framework HTTP rapide, inspiré d'Express |
| ORM | [GORM](https://gorm.io/) | v1 | ORM Go avec support PostgreSQL |
| Dependency Injection | [fx](https://uber-go.github.io/fx/) | latest | DI framework par Uber |
| Logging | [zerolog](https://github.com/rs/zerolog) | latest | Logger structuré haute performance |
| JWT | [golang-jwt](https://github.com/golang-jwt/jwt) | v5 | Tokens JWT pour authentification |
| Validation | [validator](https://github.com/go-playground/validator) | v10 | Validation de structs |
| Swagger | [swaggo](https://github.com/swaggo/swag) | latest | Documentation API OpenAPI |
| Crypto | golang.org/x/crypto | latest | Hashage bcrypt pour mots de passe |

## Documentation

### Guides essentiels

- **[Guide d'installation](./docs/installation.md)** - Installation détaillée avec toutes les méthodes
- **[Guide d'utilisation](./docs/usage.md)** - Utilisation du CLI et structure complète générée
- **[Guide des projets générés](./docs/generated-project-guide.md)** - Guide complet pour développer avec les projets créés (architecture, API, tests, déploiement)

### Documentation avancée

- **[Architecture du CLI](./docs/cli-architecture.md)** - Documentation technique pour contributeurs
- **[Guide de contribution](./docs/contributing.md)** - Comment contribuer au projet

## Démarrage rapide en 30 secondes

```bash
# 1. Installer l'outil
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest

# 2. Créer un projet
create-go-starter mon-projet

# 3. Configuration automatique
cd mon-projet
./setup.sh

# 4. Lancer
make run

# 5. Tester
curl http://localhost:8080/health
```

Ou configuration manuelle:
```bash
cd mon-projet
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
docker run -d --name postgres -e POSTGRES_DB=mon-projet -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16-alpine
make run
```

## Exemples d'utilisation de l'API

Une fois votre projet lancé, vous pouvez tester l'API:

```bash
# Créer un utilisateur
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securePass123"}'

# Se connecter
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securePass123"}'

# Utiliser le token retourné pour accéder aux endpoints protégés
TOKEN="<access_token_from_login_response>"
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

Pour plus d'exemples et la documentation complète de l'API, consultez le [guide des projets générés](./docs/generated-project-guide.md#api-reference).

## Commandes Makefile disponibles

Les projets générés incluent un Makefile avec des commandes utiles:

```bash
make help           # Afficher toutes les commandes disponibles
make run            # Lancer l'application
make build          # Build le binaire
make test           # Exécuter les tests
make test-coverage  # Tests avec rapport de coverage
make lint           # Linter le code (golangci-lint)
make docker-build   # Build l'image Docker
make docker-run     # Lancer le conteneur Docker
make clean          # Nettoyer les artifacts
```

## Prérequis

- **Go 1.25 ou supérieur** - [Télécharger Go](https://golang.org/dl/)
- **PostgreSQL** - Pour les projets générés (peut être lancé via Docker)
- **Git** - Pour cloner et contribuer
- **Docker** (optionnel) - Pour lancer PostgreSQL et containeriser l'application
- **golangci-lint** (optionnel) - Pour le linting

## Pourquoi create-go-starter?

### Gain de temps

Au lieu de passer des heures à configurer:
- L'architecture du projet
- L'authentification JWT
- La connexion à la base de données
- Les tests
- Le Docker et CI/CD
- La documentation Swagger

Obtenez tout cela en **une seule commande** et commencez immédiatement à développer vos fonctionnalités métier.

### Best practices intégrées

- **Architecture hexagonale** - Séparation claire entre domaine, adapters et infrastructure
- **Dependency injection** - Code testable et modulaire
- **Error handling centralisé** - Gestion cohérente des erreurs
- **Security-first** - JWT, bcrypt, validation, CORS
- **Tests** - Exemples de tests unitaires et d'intégration
- **Clean code** - Respect des conventions Go et linting strict

### Production-ready

Les projets générés sont prêts pour la production:
- Build Docker multi-stage optimisé
- CI/CD avec tests automatiques
- Logging structuré pour monitoring
- Configuration par environnement
- Health checks
- Graceful shutdown

## Contribuer

Les contributions sont les bienvenues! Consultez le [guide de contribution](./docs/contributing.md) pour commencer.

### Processus de contribution

1. Fork le projet
2. Créer une branche (`git checkout -b feature/ma-fonctionnalite`)
3. Commit les changements (`git commit -m 'feat: ajouter une fonctionnalité'`)
4. Push vers la branche (`git push origin feature/ma-fonctionnalite`)
5. Ouvrir une Pull Request

## Roadmap

Fonctionnalités prévues:

- [ ] Templates multiples (minimal, full, api-only, graphql)
- [ ] Support pour d'autres bases de données (MySQL, SQLite, MongoDB)
- [ ] Choix du framework web (Gin, Echo, Chi)
- [ ] CLI interactif avec prompts
- [ ] Génération de microservices
- [ ] Support GraphQL avec gqlgen
- [ ] Templates de tests E2E
- [ ] Configuration Kubernetes

## Licence

[MIT License](LICENSE) - Libre d'utilisation pour projets personnels et commerciaux.

## Support

- **Issues**: [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)
- **Documentation**: [docs/](./docs/)

## Remerciements

Construit avec les excellentes bibliothèques de la communauté Go. Merci aux mainteneurs de Fiber, GORM, fx, zerolog et toutes les autres dépendances.

---

**Fait avec ❤️ pour la communauté Go**

Commencez à construire votre prochaine application backend en secondes, pas en jours!
