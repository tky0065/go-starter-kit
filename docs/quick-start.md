# Quick Start (5 minutes)

<i class="material-icons success">bolt</i> Générez et lancez un projet Go production-ready en 5 minutes.

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

<i class="material-icons success small">check</i> **Généré**: ~45 fichiers avec architecture hexagonale complète

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
