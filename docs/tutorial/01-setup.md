# Partie 1: Installation et Configuration

<i class="material-icons info small">circle</i> **Partie 1/4** - Temps estimé: 15 minutes

[<i class="material-icons">arrow_back</i> Retour à l'index](index.md)

---

## Objectif

Créer une API REST complète pour un blog avec:

- **Articles (Posts)** avec auteur, titre, contenu, tags
- **Commentaires** sur les articles
- **Authentification JWT** (déjà incluse dans create-go-starter)
- **Tests complets**
- **Déploiement Docker**

À la fin de ce tutorial, vous aurez une API Blog production-ready avec toutes les bonnes pratiques.

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

---

## Étape 1: Installation du CLI

### Installation globale (recommandée)

La méthode la plus simple pour installer `create-go-starter`:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Cette commande télécharge, compile et installe le CLI globalement.

### Vérification

```bash
create-go-starter --help
```

Vous devriez voir l'aide s'afficher.

**Note**: Si la commande n'est pas trouvée, ajoutez `$GOPATH/bin` à votre PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Étape 2: Génération du projet

### Créer le projet

```bash
create-go-starter blog-api
```

Cette commande génère **~45 fichiers** avec toute l'architecture nécessaire.

### Structure générée

```bash
cd blog-api
tree -L 3
```

**Résultat**:
```
blog-api/
├── cmd/
│   └── main.go                       # Point d'entrée avec fx DI
├── internal/
│   ├── models/
│   │   └── user.go                   # Entités: User, RefreshToken, AuthResponse
│   ├── domain/
│   │   ├── user/                     # Domaine User (pré-généré)
│   │   │   ├── service.go
│   │   │   └── module.go
│   │   └── errors.go
│   ├── adapters/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go
│   │   │   └── error_handler.go
│   │   └── repository/
│   │       └── user_repository.go
│   ├── infrastructure/
│   │   ├── database/
│   │   └── server/
│   └── interfaces/                   # Ports (interfaces)
│       └── user_repository.go
├── pkg/
│   ├── auth/                         # JWT utilities
│   ├── config/                       # Configuration
│   └── logger/                       # Zerolog logger
├── docs/
│   ├── README.md
│   └── quick-start.md
├── .env                              # Configuration (auto-copié)
├── .env.example
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

<i class="material-icons success">check_circle</i> **Checkpoint 1**: Le projet est généré avec succès.

---

## Étape 3: Configuration initiale

### 3.1 Installer les dépendances

```bash
cd blog-api
go mod tidy
```

Cette commande télécharge toutes les dépendances (Fiber, GORM, fx, etc.).

### 3.2 Configurer PostgreSQL

Vous avez 2 options:

#### Option A: Docker (recommandé)

```bash
docker run -d \
  --name blog-postgres \
  -e POSTGRES_DB=blog_api \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

#### Option B: PostgreSQL local

Si PostgreSQL est installé localement:

```bash
createdb blog_api
```

### 3.3 Configurer les variables d'environnement

Générer un secret JWT sécurisé:

```bash
JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET généré: $JWT_SECRET"
```

Éditer le fichier `.env`:

```bash
nano .env
```

Contenu du `.env`:

```env
# Application
APP_NAME=blog-api
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=blog_api
DB_SSLMODE=disable

# JWT
JWT_SECRET=<coller_le_secret_généré_ici>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

**Important**: Remplacez `<coller_le_secret_généré_ici>` par le JWT_SECRET généré.

---

## Étape 4: Tester le projet de base

### 4.1 Lancer l'application

```bash
make run
```

Vous devriez voir:

```
2024/01/10 10:00:00 INF Starting blog-api server on :8080
```

### 4.2 Tester le health check

Dans un autre terminal:

```bash
curl http://localhost:8080/health
```

**Réponse attendue**:
```json
{"status":"ok"}
```

### 4.3 Tester l'authentification par défaut

#### Créer un utilisateur

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Réponse**:
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "user": {
    "id": 1,
    "email": "admin@blog.com",
    "created_at": "2024-01-10T10:05:00Z"
  }
}
```

#### Se connecter

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Même réponse** avec access_token et refresh_token.

#### Tester une route protégée

```bash
# Remplacez <ACCESS_TOKEN> par le token reçu
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Réponse**:
```json
[
  {
    "id": 1,
    "email": "admin@blog.com",
    "created_at": "2024-01-10T10:05:00Z"
  }
]
```

<i class="material-icons success">check_circle</i> **Checkpoint 2**: Le projet de base fonctionne parfaitement avec User et Auth.

---

## Résumé de la Partie 1

<i class="material-icons success">check</i> Installation du CLI `create-go-starter`
<i class="material-icons success">check</i> Génération d'un projet complet
<i class="material-icons success">check</i> Configuration PostgreSQL et JWT
<i class="material-icons success">check</i> Test de l'authentification

Vous avez maintenant un projet fonctionnel avec authentification JWT et gestion des utilisateurs.

---

## Navigation

[<i class="material-icons">arrow_back</i> Retour à l'index](index.md) | [Partie 2: Créer votre premier domaine <i class="material-icons">arrow_forward</i>](02-first-domain.md)
