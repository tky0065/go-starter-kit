# Foire Aux Questions (FAQ)

Questions fréquemment posées sur `create-go-starter` et les projets générés.

---

## <i class="material-icons">download</i> Installation & Configuration

### Comment installer `create-go-starter`?

La méthode recommandée est l'installation via `go install`:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Cette commande télécharge et installe automatiquement le binaire dans `$GOPATH/bin`.

**Vérifier l'installation:**

```bash
create-go-starter --help
```

<i class="material-icons info">arrow_forward</i> **Détails:** Voir le [guide d'installation complet](../installation.md)

---

### <i class="material-icons warning">warning</i> Erreur "command not found: create-go-starter"

**Cause:** Le répertoire `$GOPATH/bin` n'est pas dans votre `PATH`, ou le cache du shell n'a pas été rechargé.

**Solutions:**

**1. Recharger le cache du shell** (essayez d'abord ceci):

```bash
hash -r
which create-go-starter
```

**2. Redémarrer le terminal** (souvent la solution la plus simple)

Fermez et rouvrez votre terminal.

**3. Ajouter `$GOPATH/bin` au `PATH`:**

```bash
# Vérifier GOPATH
go env GOPATH

# Ajouter au PATH (dans ~/.zshrc ou ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

# Recharger
source ~/.zshrc  # ou ~/.bashrc
```

**4. Utiliser le chemin complet temporairement:**

```bash
$(go env GOPATH)/bin/create-go-starter mon-projet
```

<i class="material-icons info">arrow_forward</i> **Détails:** [Guide d'installation - Résolution de problèmes](../installation.md#résolution-de-problèmes)

---

### Quelle version de Go est requise?

**Go 1.25 ou supérieur** est requis.

**Vérifier votre version:**

```bash
go version
# Devrait afficher: go version go1.25.x ...
```

**Mettre à jour Go:** Téléchargez depuis [golang.org/dl](https://golang.org/dl/)

---

### Comment mettre à jour `create-go-starter`?

Réexécutez simplement la commande d'installation:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Go téléchargera automatiquement la dernière version.

---

## <i class="material-icons">storage</i> Base de données & Configuration

### Pourquoi PostgreSQL est-il la base de données par défaut?

**PostgreSQL** est choisi par défaut pour plusieurs raisons:

- <i class="material-icons success">check</i> **Fiabilité production** - Conforme ACID, forte intégrité des données
- <i class="material-icons success">check</i> **Fonctionnalités avancées** - JSON, arrays, recherche full-text
- <i class="material-icons success">check</i> **Performances** - Excellent pour requêtes complexes
- <i class="material-icons success">check</i> **Communauté** - Écosystème riche, documentation complète

PostgreSQL est le choix standard pour les applications backend modernes.

<i class="material-icons info">arrow_forward</i> **Comparaison:** [Guide de sélection des bases de données](../databases.md)

---

### Comment changer de base de données?

**Lors de la création d'un nouveau projet:**

```bash
# MySQL
create-go-starter mon-app --database=mysql

# SQLite
create-go-starter mon-app --database=sqlite

# PostgreSQL (défaut)
create-go-starter mon-app --database=postgres
```

**Pour un projet existant:**

Vous devez régénérer le projet avec le nouveau flag `--database` et migrer les données.

<i class="material-icons warning">warning</i> **Important:** Sauvegardez vos données avant de migrer.

<i class="material-icons info">arrow_forward</i> **Guide complet:** [Guide de migration de base de données](../database-migration.md)

---

### Quelle base de données choisir pour mon projet?

**PostgreSQL** (défaut):
- Applications de production
- Données relationnelles complexes
- Requêtes avancées (JSON, analytics)
- Fiabilité et scalabilité

**MySQL**:
- Hébergement partagé (shared hosting)
- Compatibilité large
- Charges de lecture intensives
- Équipes familières avec MySQL

**SQLite**:
- Prototypage rapide / MVP
- Petites applications (<100 utilisateurs concurrents)
- Développement et tests
- Applications desktop/embarquées

<i class="material-icons info">arrow_forward</i> **Comparaison détaillée:** [databases.md - Matrice de décision](../databases.md#matrice-de-décision)

---

### <i class="material-icons error">error</i> Erreur "JWT_SECRET not set" - que faire?

**Cause:** Le secret JWT n'est pas configuré dans le fichier `.env`.

**Solution:**

**1. Générer un secret JWT:**

```bash
openssl rand -base64 32
# Exemple de sortie: xK7vZ9mN2pQ8rT4wL6jH3sB5cA1dF0eG=
```

**2. Ajouter au fichier `.env`:**

```bash
# Éditer .env
JWT_SECRET=xK7vZ9mN2pQ8rT4wL6jH3sB5cA1dF0eG=
```

**Ou utiliser le script automatique:**

```bash
./setup.sh
# Le script génère automatiquement JWT_SECRET
```

!!! warning "Sécurité"
    - Ne jamais commiter le fichier `.env` dans Git
    - Utiliser un secret différent pour chaque environnement (dev, staging, prod)
    - Le secret doit être au minimum 32 caractères

<i class="material-icons info">arrow_forward</i> **Documentation:** [Guide des projets générés - Configuration](../generated-project-guide.md#configuration)

---

### Ai-je besoin de Docker?

Cela dépend de la base de données choisie:

| Base de données | Docker requis? | Alternative |
|-----------------|----------------|-------------|
| **PostgreSQL** | <i class="material-icons success">check</i> Oui (recommandé) | Installation locale possible |
| **MySQL** | <i class="material-icons success">check</i> Oui (recommandé) | Installation locale possible |
| **SQLite** | <i class="material-icons error">close</i> Non | Base de données embarquée (fichier) |

**Lancer PostgreSQL avec Docker:**

```bash
docker run -d --name postgres \
  -e POSTGRES_DB=mon-app \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

**Ou utiliser docker-compose** (inclus dans les projets générés):

```bash
docker-compose up -d
```

---

### Comment résoudre les erreurs de connexion à la base de données?

**Vérifications à effectuer:**

**1. Vérifier que la base de données est démarrée:**

```bash
# Pour Docker
docker ps
# Devrait montrer le conteneur postgres en "Up"

# Pour PostgreSQL local
pg_isready -h localhost -p 5432
```

**2. Vérifier la configuration `.env`:**

```bash
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mon-app
DB_SSLMODE=disable

# MySQL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=mon-app

# SQLite
DB_NAME=mon-app.db
```

**3. Tester la connexion manuellement:**

```bash
# PostgreSQL
psql -h localhost -U postgres -d mon-app

# MySQL
mysql -h localhost -u root -p mon-app
```

**4. Vérifier les logs de l'application:**

```bash
make run
# Les logs indiquent souvent la cause exacte
```

<i class="material-icons warning">warning</i> **Erreur fréquente:** Oublier de démarrer Docker ou PostgreSQL avant de lancer l'app.

---

### Comment résoudre les erreurs de migration de base de données?

**Symptômes:**
- Tables non créées
- Colonnes manquantes
- Erreurs "relation does not exist"

**Solutions:**

**1. Forcer la recréation des tables (développement uniquement):**

```go
// Dans infrastructure/database/database.go
db.AutoMigrate(&models.User{}, &models.RefreshToken{})
// Ajouter temporairement:
// db.Migrator().DropTable(&models.User{}, &models.RefreshToken{})
// db.AutoMigrate(&models.User{}, &models.RefreshToken{})
```

**2. Vérifier les logs GORM:**

Les logs indiquent les requêtes SQL exécutées et les erreurs.

**3. Migration manuelle:**

```bash
# Se connecter à la base de données
psql -h localhost -U postgres -d mon-app

# Vérifier les tables
\dt

# Recréer si nécessaire
DROP TABLE users, refresh_tokens CASCADE;
```

**4. Redémarrer l'application:**

```bash
make run
# GORM recrée automatiquement les tables au démarrage
```

!!! danger "Production"
    Ne jamais utiliser `DropTable` ou `AutoMigrate` en production. Utilisez des migrations versionnées (ex: golang-migrate).

---

## <i class="material-icons">architecture</i> Architecture & Design

### Pourquoi l'architecture hexagonale?

L'**architecture hexagonale** (Ports & Adapters) offre plusieurs avantages:

**<i class="material-icons success">check</i> Séparation des responsabilités**
- Domaine métier isolé
- Infrastructure interchangeable
- Code facile à tester

**<i class="material-icons success">check</i> Testabilité**
- Mocks faciles à créer (interfaces)
- Tests unitaires sans base de données
- Tests d'intégration isolés

**<i class="material-icons success">check</i> Maintenabilité**
- Changements localisés
- Ajout de fonctionnalités facile
- Refactoring sûr

**<i class="material-icons success">check</i> Scalabilité**
- Ajout de nouvelles couches simple
- Migration vers microservices facilitée

<i class="material-icons info">arrow_forward</i> **Détails:** [Guide des projets générés - Architecture](../generated-project-guide.md#architecture)

---

### Quelle est la différence entre `models/`, `domain/` et `interfaces/`?

**`internal/models/`** - **Entités partagées**
- Structures de données GORM (User, RefreshToken)
- DTOs (AuthResponse)
- Partagées par tous les modules
- Évite les dépendances circulaires

**`internal/domain/`** - **Logique métier**
- Services (UserService)
- Règles métier
- Cas d'utilisation
- Ne dépend de rien (sauf models et interfaces)

**`internal/interfaces/`** - **Ports (contrats)**
- Interfaces définissant les contrats
- UserRepository (interface)
- Permet l'inversion de dépendances

**Exemple de flux:**

```go
Handler → Service (domain/) → Repository (interfaces/) ← Implementation (adapters/repository/)
```

<i class="material-icons info">arrow_forward</i> **Diagrammes:** [Guide des projets générés - Architecture hexagonale](../generated-project-guide.md#architecture-hexagonale-ports--adapters)

---

### Quand utiliser le template `minimal` vs `full`?

**Template `minimal`** - API REST basique:
- <i class="material-icons small">circle</i> Prototypage rapide
- <i class="material-icons small">circle</i> APIs publiques simples
- <i class="material-icons small">circle</i> Pas besoin d'authentification
- <i class="material-icons small">circle</i> Learning / démonstration

**Template `full`** (défaut) - API complète avec JWT:
- <i class="material-icons small">circle</i> Applications backend complètes
- <i class="material-icons small">circle</i> Authentification nécessaire
- <i class="material-icons small">circle</i> Gestion des utilisateurs
- <i class="material-icons small">circle</i> Production-ready

**Template `graphql`** - API GraphQL:
- <i class="material-icons small">circle</i> Applications nécessitant GraphQL
- <i class="material-icons small">circle</i> Requêtes flexibles
- <i class="material-icons small">circle</i> Frontend React/Vue/Angular

**Exemple:**

```bash
# Prototype rapide sans auth
create-go-starter demo-api --template=minimal --database=sqlite

# Application production avec auth
create-go-starter mon-app --template=full --database=postgres

# API GraphQL moderne
create-go-starter graphql-api --template=graphql --database=mysql
```

<i class="material-icons info">arrow_forward</i> **Comparaison:** [Guide d'utilisation - Templates disponibles](../usage.md#templates-disponibles)

---

### Comment ajouter un nouveau modèle (entité)?

**Commande `add-model` (v1.2.0+):**

```bash
# Modèle simple
create-go-starter add-model Product --fields "name:string,price:float64"

# Avec relation BelongsTo
create-go-starter add-model Comment --fields "content:string" --belongs-to Post

# Avec relation HasMany
create-go-starter add-model Category --fields "name:string:unique" --has-many Product
```

**Ce qui est généré automatiquement:**
- <i class="material-icons success">check</i> Model avec tags GORM (`internal/models/`)
- <i class="material-icons success">check</i> Repository interface (`internal/interfaces/`)
- <i class="material-icons success">check</i> Repository implémentation (`internal/adapters/repository/`)
- <i class="material-icons success">check</i> Service avec logique métier (`internal/domain/`)
- <i class="material-icons success">check</i> Handlers HTTP CRUD (`internal/adapters/handlers/`)
- <i class="material-icons success">check</i> Tests unitaires
- <i class="material-icons success">check</i> Routes ajoutées automatiquement

<i class="material-icons info">arrow_forward</i> **Détails:** [README - Ajouter des modèles avec relations](../index.md#ajouter-des-modèles-avec-relations-nouveau-v120)

---

## <i class="material-icons">bug_report</i> Erreurs courantes

### <i class="material-icons error">error</i> Erreur de compilation "cannot find package"

**Cause:** Dépendances Go non installées.

**Solution:**

```bash
# Installer les dépendances
go mod tidy

# Vérifier go.mod et go.sum
cat go.mod

# Nettoyer et réinstaller
go clean -modcache
go mod tidy
```

---

### <i class="material-icons error">error</i> Port 8080 déjà utilisé

**Cause:** Un autre processus utilise le port 8080.

**Solution 1 - Changer le port dans `.env`:**

```bash
# .env
PORT=3000
```

**Solution 2 - Arrêter le processus existant:**

```bash
# Trouver le processus
lsof -i :8080

# Tuer le processus
kill -9 <PID>
```

**Solution 3 - Utiliser un port différent temporairement:**

```bash
PORT=3000 make run
```

---

### <i class="material-icons error">error</i> Tests échouent avec "database connection failed"

**Cause:** Les tests tentent de se connecter à une base de données réelle.

**Solutions:**

**1. Utiliser des tests unitaires avec mocks:**

Les tests générés utilisent des mocks pour éviter les connexions DB:

```go
// Example de mock repository
type mockUserRepository struct{}

func (m *mockUserRepository) Create(user *models.User) error {
    return nil
}
```

**2. Pour tests d'intégration, configurer une DB de test:**

```bash
# .env.test
DB_NAME=myapp_test
```

```go
// Dans le test
os.Setenv("DB_NAME", "myapp_test")
```

**3. Utiliser SQLite en mémoire pour les tests:**

```go
dsn := ":memory:"
db, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
```

<i class="material-icons info">arrow_forward</i> **Exemples:** [Guide des projets générés - Tests](../generated-project-guide.md#tests)

---

### <i class="material-icons error">error</i> Docker build échoue

**Erreurs fréquentes:**

**1. "go.sum not found":**

```bash
# Générer go.sum
go mod tidy

# Puis rebuild Docker
docker build -t mon-app .
```

**2. "cannot find package":**

Vérifier que `go.mod` et `go.sum` sont à jour:

```bash
go mod download
go mod verify
docker build -t mon-app .
```

**3. Build trop lent:**

Utiliser le cache Docker:

```dockerfile
# Le Dockerfile généré utilise déjà le cache
COPY go.mod go.sum ./
RUN go mod download
# Puis copier le code source
COPY . .
```

---

## <i class="material-icons">verified_user</i> Bonnes pratiques

### Comment organiser mon code?

**Suivre l'architecture hexagonale:**

```
internal/
├── models/          # Entités partagées (User, Product, etc.)
├── domain/          # Logique métier par domaine
│   ├── user/        # Domaine utilisateur
│   │   ├── service.go
│   │   └── module.go
│   └── product/     # Nouveau domaine
│       ├── service.go
│       └── module.go
├── interfaces/      # Ports (interfaces)
│   ├── user_repository.go
│   └── product_repository.go
├── adapters/        # Implémentations
│   ├── handlers/    # HTTP handlers
│   ├── repository/  # DB repositories
│   └── middleware/  # Middleware
└── infrastructure/  # Infrastructure
    ├── database/
    └── server/
```

**Règles:**
- <i class="material-icons success">check</i> Domain ne dépend de rien (sauf models et interfaces)
- <i class="material-icons success">check</i> Un package par domaine métier
- <i class="material-icons success">check</i> Pas de logique métier dans les handlers
- <i class="material-icons success">check</i> Interfaces pour tous les contrats

<i class="material-icons info">arrow_forward</i> **Détails:** [Guide des projets générés - Bonnes pratiques](../generated-project-guide.md#bonnes-pratiques)

---

### Comment tester efficacement mon code?

**Tests unitaires** - Logique métier isolée:

```go
// Tester le service avec mock repository
func TestUserService_Register(t *testing.T) {
    mockRepo := &mockUserRepository{}
    service := NewUserService(mockRepo, logger)

    user, err := service.Register("test@example.com", "password123")

    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

**Tests d'intégration** - Avec vraie DB:

```go
// Utiliser SQLite en mémoire
dsn := ":memory:"
db, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
repo := NewUserRepository(db)

user, err := repo.Create(&models.User{Email: "test@example.com"})
assert.NoError(t, err)
```

**Tests E2E** - Requêtes HTTP complètes:

```go
app := fiber.New()
app.Post("/api/v1/auth/register", handler.Register)

req := httptest.NewRequest("POST", "/api/v1/auth/register", body)
resp, _ := app.Test(req)

assert.Equal(t, 201, resp.StatusCode)
```

**Lancer les tests:**

```bash
make test              # Tous les tests
make test-coverage     # Avec rapport coverage
go test ./... -v       # Mode verbose
go test -run TestName  # Test spécifique
```

---

### Quelles variables d'environnement sont nécessaires?

**Variables essentielles** (`.env`):

```bash
# Application
PORT=8080
ENV=development

# Base de données (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mon-app
DB_SSLMODE=disable

# Sécurité
JWT_SECRET=<généré avec: openssl rand -base64 32>
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Logging
LOG_LEVEL=info
```

**Variables optionnelles:**

```bash
# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://monapp.com

# Rate limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX=100
```

!!! warning "Sécurité"
    - Toujours utiliser `.env.example` comme template
    - Ne jamais commiter `.env` dans Git
    - Utiliser des secrets différents par environnement

<i class="material-icons info">arrow_forward</i> **Configuration complète:** [Guide des projets générés - Configuration](../generated-project-guide.md#configuration)

---

## <i class="material-icons">cloud_upload</i> Déploiement

### Comment déployer en production?

**Option 1 - Docker (Recommandé):**

```bash
# Build l'image
docker build -t mon-app:latest .

# Lancer le conteneur
docker run -d \
  -p 8080:8080 \
  --env-file .env.production \
  mon-app:latest
```

**Option 2 - Binaire natif:**

```bash
# Build pour production
make build

# Copier sur serveur
scp mon-app user@server:/opt/mon-app/

# Lancer avec systemd
sudo systemctl start mon-app
```

**Option 3 - Docker Compose:**

```bash
# Sur le serveur
docker-compose -f docker-compose.prod.yml up -d
```

**Checklist de déploiement:**
- <i class="material-icons small">circle</i> Variables d'environnement configurées (`.env.production`)
- <i class="material-icons small">circle</i> JWT secret généré (différent de dev)
- <i class="material-icons small">circle</i> Base de données configurée et accessible
- <i class="material-icons small">circle</i> Logs configurés (niveau `info` ou `warn`)
- <i class="material-icons small">circle</i> CORS configuré pour domaine production
- <i class="material-icons small">circle</i> Reverse proxy (nginx) configuré
- <i class="material-icons small">circle</i> SSL/TLS activé (Let's Encrypt)
- <i class="material-icons small">circle</i> Monitoring configuré

<i class="material-icons info">arrow_forward</i> **Guide complet:** [Guide des projets générés - Déploiement](../generated-project-guide.md#déploiement)

---

### Comment gérer les migrations en production?

**Approche recommandée - Migrations versionnées:**

Installer `golang-migrate`:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Créer des migrations:**

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

**Appliquer les migrations:**

```bash
migrate -path migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up
```

**Alternative - Script SQL manuel:**

```sql
-- migrations/001_initial.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```bash
psql -h localhost -U postgres -d mon-app < migrations/001_initial.sql
```

!!! danger "Production"
    - Toujours tester les migrations en staging d'abord
    - Sauvegarder la base de données avant migration
    - Ne jamais utiliser `AutoMigrate` en production

---

### Comment monitorer mon application?

**Logs structurés** (zerolog inclus):

```go
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Msg("User logged in successfully")
```

**Health check endpoint:**

```bash
# Déjà inclus dans les projets générés
curl http://localhost:8080/health
# {"status":"ok"}
```

**Monitoring externe:**

Intégrer avec des services comme:
- Prometheus + Grafana (métriques)
- Sentry (error tracking)
- DataDog (monitoring complet)
- New Relic (APM)

**Exemple avec Prometheus:**

```go
// Ajouter des métriques
import "github.com/prometheus/client_golang/prometheus"

var requestCount = prometheus.NewCounterVec(
    prometheus.CounterOpts{Name: "http_requests_total"},
    []string{"method", "endpoint"},
)
```

<i class="material-icons info">arrow_forward</i> **Détails:** [Guide des projets générés - Monitoring](../generated-project-guide.md#monitoring--logging)

---

## <i class="material-icons">terminal</i> CLI - Fonctionnalités v1.5.2

### Comment utiliser le mode interactif?

Le mode interactif lance un assistant pas-à-pas qui vous guide dans la configuration de votre projet:

```bash
create-go-starter --interactive
# ou avec l'alias court
create-go-starter -i
```

L'assistant vous pose des questions pour chaque option (nom du projet, template, base de données, observabilité) et valide vos réponses en temps réel.

<i class="material-icons info">arrow_forward</i> **Détails:** [Guide d'utilisation - Mode Interactif](../usage.md#mode-interactif---interactive)

---

### Comment prévisualiser les fichiers sans les créer?

Utilisez le mode **dry-run** pour voir la liste complète des fichiers qui seraient générés, sans écrire quoi que ce soit sur le disque:

```bash
create-go-starter mon-app --dry-run
# ou avec l'alias court
create-go-starter mon-app -n

# Combiner avec d'autres options
create-go-starter mon-app -t minimal -d sqlite -n
```

Le dry-run affiche le nombre de fichiers et répertoires qui seraient créés, avec un résumé de la configuration.

<i class="material-icons info">arrow_forward</i> **Détails:** [Guide d'utilisation - Dry-Run](../usage.md#prévisualisation-dry-run---dry-run)

---

### Comment diagnostiquer mon environnement avec `doctor`?

La commande `doctor` vérifie que votre environnement dispose de tous les outils nécessaires:

```bash
create-go-starter doctor
```

Elle vérifie:
- <i class="material-icons success">check</i> **Go** (version >= 1.21 requise)
- <i class="material-icons success">check</i> **Git** (détection de présence)
- <i class="material-icons success">check</i> **Docker** (optionnel, pour PostgreSQL/MySQL)

Chaque vérification affiche un statut clair (OK / WARNING / FAIL) avec des recommandations si un problème est détecté.

<i class="material-icons info">arrow_forward</i> **Détails:** [Guide d'utilisation - Commande Doctor](../usage.md#commande-doctor)

---

### Quels sont les alias courts disponibles?

Depuis la v1.5.2, toutes les options CLI ont des alias courts d'une lettre:

| Alias | Option longue | Description |
|-------|---------------|-------------|
| `-h` | `--help` | Afficher l'aide |
| `-i` | `--interactive` | Mode interactif |
| `-n` | `--dry-run` | Prévisualisation sans écriture |
| `-t` | `--template` | Template (minimal, full, graphql) |
| `-d` | `--database` | Base de données (postgres, mysql, sqlite) |
| `-o` | `--observability` | Observabilité (none, basic, advanced) |

**Exemples de syntaxe flexible:**

```bash
# Toutes ces syntaxes sont valides
create-go-starter mon-app -t minimal -d sqlite
create-go-starter mon-app -t=minimal -d=sqlite
create-go-starter mon-app --template minimal --database sqlite
```

---

### La barre de progression ne s'affiche pas, est-ce normal?

La barre de progression visuelle est automatiquement désactivée dans certains contextes:

- **Terminaux non-interactifs** (CI/CD, pipes) - détection automatique via `os.Stdout.Fd()`
- **Variable `NO_COLOR` définie** - respecte la convention [no-color.org](https://no-color.org)
- **Redirection de sortie** (`> fichier`, `| grep`)

Dans ces cas, la génération se fait silencieusement et les statistiques finales sont toujours affichées sous forme texte.

---

## <i class="material-icons">help</i> Questions diverses

### Puis-je utiliser `create-go-starter` pour des projets commerciaux?

**Oui, absolument!** Le projet est sous licence MIT, ce qui permet:
- <i class="material-icons success">check</i> Utilisation commerciale
- <i class="material-icons success">check</i> Modification
- <i class="material-icons success">check</i> Distribution
- <i class="material-icons success">check</i> Usage privé

Aucune attribution requise (mais appréciée!).

---

### Comment contribuer au projet?

**Processus:**

1. **Fork** le repository
2. **Clone** votre fork
3. **Créer** une branche feature
4. **Développer** avec tests
5. **Commit** selon conventions
6. **Push** et créer une Pull Request

```bash
git clone https://github.com/VOTRE-USERNAME/go-starter-kit.git
cd go-starter-kit
git checkout -b feature/ma-fonctionnalite

# Développer, tester
make test
make lint

# Commit
git commit -m "feat: ajouter fonctionnalité X"

# Push
git push origin feature/ma-fonctionnalite
```

<i class="material-icons info">arrow_forward</i> **Guide complet:** [Guide de contribution](../contributing.md)

---

### Où obtenir de l'aide?

**Ressources:**

- <i class="material-icons">menu_book</i> **Documentation:** [docs/](../)
- <i class="material-icons">bug_report</i> **Issues:** [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)
- <i class="material-icons">forum</i> **Discussions:** [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)
- <i class="material-icons">chat</i> **Questions:** Ouvrir une discussion

**Avant de poser une question:**
1. Vérifier cette FAQ
2. Lire la documentation pertinente
3. Chercher dans les Issues/Discussions existantes

---

### La documentation est-elle à jour?

Oui! La documentation est **systématiquement mise à jour** avec chaque changement de code.

**Processus:**
- Code modifié → Documentation mise à jour **immédiatement**
- Tests documentés avec exemples
- Changements testés localement (`mkdocs serve`)

**Si vous trouvez une erreur:**
Ouvrir une [Issue](https://github.com/tky0065/go-starter-kit/issues) ou [Pull Request](https://github.com/tky0065/go-starter-kit/pulls).

---

## <i class="material-icons">lightbulb</i> Besoin d'aide supplémentaire?

**Cette FAQ n'a pas répondu à votre question?**

- <i class="material-icons">arrow_forward</i> **Documentation complète:** [Index](../index.md)
- <i class="material-icons">arrow_forward</i> **Guide d'installation:** [installation.md](../installation.md)
- <i class="material-icons">arrow_forward</i> **Guide d'utilisation:** [usage.md](../usage.md)
- <i class="material-icons">arrow_forward</i> **Guide des projets générés:** [generated-project-guide.md](../generated-project-guide.md)
- <i class="material-icons">arrow_forward</i> **Ouvrir une Discussion:** [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)

---

**Dernière mise à jour:** 2026-02-18
