# Story 7.4: MongoDB Support (Optional)

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** développeur,
**Je veux** générer un projet avec support MongoDB,
**Afin de** pouvoir utiliser une base de données NoSQL document-oriented.

## Acceptance Criteria

1. **AC1**: Given l'utilisateur exécute `create-go-starter mon-projet --database=mongodb`, When le projet est généré, Then le driver mongo-go-driver est configuré dans go.mod (PAS GORM)
2. **AC2**: Given le flag `--database=mongodb` est utilisé, When les fichiers sont générés, Then l'architecture est adaptée pour NoSQL (repositories utilisent mongo-driver, pas GORM)
3. **AC3**: Given le projet est généré avec `--database=mongodb`, When docker-compose.yml est créé, Then il contient un service MongoDB 7.0+
4. **AC4**: Given le projet généré avec MongoDB, When `make dev` est exécuté, Then le projet se connecte à MongoDB et les opérations CRUD fonctionnent correctement

## Tasks / Subtasks

- [ ] Task 1: Ajouter mongo-driver aux templates (AC: 1)
  - [ ] 1.1 Dans `GoModDependencies()`, cas "mongodb" retourne `go.mongodb.org/mongo-driver v1.13.1`
  - [ ] 1.2 IMPORTANT: Ne PAS inclure GORM pour MongoDB
  - [ ] 1.3 Ajouter dépendance `go.mongodb.org/mongo-driver/bson`

- [ ] Task 2: Créer templates MongoDB-specific (AC: 2)
  - [ ] 2.1 Créer `templates_mongodb.go` pour templates NoSQL spécifiques
  - [ ] 2.2 Template repository sans GORM, utilisant mongo Collection
  - [ ] 2.3 Template models avec tags `bson` au lieu de `gorm`
  - [ ] 2.4 Template database connection avec `mongo.Connect()`

- [ ] Task 3: Configuration MongoDB (AC: 2)
  - [ ] 3.1 DSN MongoDB: `mongodb://user:pass@host:port/dbname`
  - [ ] 3.2 Variables .env: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME
  - [ ] 3.3 Code de connexion: client mongo, pas GORM DB

- [ ] Task 4: Docker-compose pour MongoDB (AC: 3)
  - [ ] 4.1 Service MongoDB 7.0 avec authentification
  - [ ] 4.2 Variables: MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD, MONGO_INITDB_DATABASE
  - [ ] 4.3 Port 27017
  - [ ] 4.4 Volume pour persistance
  - [ ] 4.5 Healthcheck MongoDB

- [ ] Task 5: Adapter l'architecture pour NoSQL (AC: 2)
  - [ ] 5.1 Models: utiliser `bson` tags, types ObjectID
  - [ ] 5.2 Repositories: méthodes avec *mongo.Collection
  - [ ] 5.3 Pas de migrations (MongoDB est schema-less)
  - [ ] 5.4 Services: adapter pour types MongoDB

- [ ] Task 6: Tests MongoDB (AC: 4)
  - [ ] 6.1 Tests unitaires pour templates MongoDB
  - [ ] 6.2 Test E2E génération projet MongoDB
  - [ ] 6.3 Vérifier absence de GORM dans go.mod
  - [ ] 6.4 Vérifier présence de mongo-driver
  - [ ] 6.5 Test optionnel: connexion effective à MongoDB

## Dev Notes

### 🎯 Objectif Principal

MongoDB est **FONDAMENTALEMENT DIFFÉRENT** des autres databases car c'est une base NoSQL. Cette story nécessite une **architecture adaptée** sans GORM, avec des modèles BSON et des repositories utilisant le driver MongoDB natif.

### ⚠️ CRITICAL: Architecture NoSQL

**DIFFÉRENCE MAJEURE**:
- PostgreSQL/MySQL/SQLite → **SQL** → GORM → Structs avec tags `gorm`
- MongoDB → **NoSQL** → mongo-driver → Documents BSON avec tags `bson`

**Implications**:
- ❌ **PAS de GORM** pour MongoDB
- ✅ Utiliser `go.mongodb.org/mongo-driver` directement
- ✅ Repositories utilisent `*mongo.Collection` au lieu de `*gorm.DB`
- ✅ Models utilisent `bson` tags et types (`primitive.ObjectID`)
- ❌ **PAS de migrations** (MongoDB est schema-less)

### 🏗️ Architecture MongoDB vs SQL

#### Comparaison Structurelle

| Aspect | SQL (GORM) | MongoDB (mongo-driver) |
|--------|------------|------------------------|
| **ORM/Driver** | GORM | mongo-driver natif |
| **Connection** | `*gorm.DB` | `*mongo.Client` |
| **Repository** | `db.Find(&users)` | `collection.Find(ctx, filter)` |
| **Model Tags** | `gorm:"primaryKey"` | `bson:"_id"` |
| **ID Type** | `uint` | `primitive.ObjectID` |
| **Migrations** | AutoMigrate | N/A (schema-less) |
| **Transactions** | GORM transactions | mongo Sessions |
| **Queries** | SQL-like | BSON filters |

### 📁 Templates MongoDB Spécifiques

#### go.mod (Sans GORM)

```go
require (
    go.mongodb.org/mongo-driver v1.13.1
    github.com/gofiber/fiber/v2 v2.52.10
    go.uber.org/fx v1.24.1
    // ... PAS de gorm.io/gorm
)
```

#### Model (avec BSON tags)

```go
package domain

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User représente un utilisateur dans MongoDB
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email     string             `bson:"email" json:"email" validate:"required,email"`
    Password  string             `bson:"password" json:"-"`
    Name      string             `bson:"name" json:"name"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

**Différences clés**:
- `primitive.ObjectID` au lieu de `uint`
- Tags `bson` au lieu de `gorm`
- Pas de `gorm.Model` embedded
- Pas de `DeletedAt` (soft delete géré différemment)

#### Repository (mongo Collection)

```go
package repository

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
    return &UserRepository{
        collection: db.Collection("users"),
    }
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    
    var user domain.User
    err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

#### Database Connection

```go
package database

import (
    "context"
    "fmt"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(config *config.Config) (*mongo.Database, error) {
    // DSN MongoDB
    dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
        config.DBUser,
        config.DBPassword,
        config.DBHost,
        config.DBPort,
        config.DBName,
    )
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Options de connexion
    clientOptions := options.Client().ApplyURI(dsn)
    
    // Connexion
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }
    
    // Ping pour vérifier la connexion
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }
    
    return client.Database(config.DBName), nil
}
```

#### docker-compose.yml (Service MongoDB)

```yaml
services:
  db:
    image: mongo:7.0
    container_name: ${PROJECT_NAME}-mongodb
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${DB_USER}
      MONGO_INITDB_ROOT_PASSWORD: ${DB_PASSWORD}
      MONGO_INITDB_DATABASE: ${DB_NAME}
    ports:
      - "${DB_PORT}:27017"
    volumes:
      - mongodb_data:/data/db
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  mongodb_data:
```

### 🔧 Implémentation Technique

#### Structure de templates_mongodb.go (NOUVEAU fichier)

```go
package main

// MongoDBModelsTemplate génère les models avec tags bson
func MongoDBModelsTemplate(projectName string) string {
    return `package domain

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID ` + "`bson:\"_id,omitempty\" json:\"id\"`" + `
    Email     string             ` + "`bson:\"email\" json:\"email\" validate:\"required,email\"`" + `
    Password  string             ` + "`bson:\"password\" json:\"-\"`" + `
    Name      string             ` + "`bson:\"name\" json:\"name\"`" + `
    CreatedAt time.Time          ` + "`bson:\"created_at\" json:\"created_at\"`" + `
    UpdatedAt time.Time          ` + "`bson:\"updated_at\" json:\"updated_at\"`" + `
}
`
}

// MongoDBRepositoryTemplate génère un repository avec mongo.Collection
func MongoDBRepositoryTemplate(projectName string) string {
    // Template du repository MongoDB (voir exemple ci-dessus)
}

// MongoDBDatabaseTemplate génère le code de connexion MongoDB
func MongoDBDatabaseTemplate(projectName string) string {
    // Template de connexion MongoDB (voir exemple ci-dessus)
}
```

#### Modification de generator.go

```go
func generateProjectFiles(projectPath, projectName, template, database string) error {
    t := NewProjectTemplates(projectName, template, database)
    
    // Switch sur database pour générer les bons fichiers
    switch database {
    case "mongodb":
        // Utiliser templates MongoDB spécifiques
        if err := writeFile(filepath.Join(projectPath, "internal/domain/user.go"),
            MongoDBModelsTemplate(projectName)); err != nil {
            return err
        }
        
        if err := writeFile(filepath.Join(projectPath, "internal/adapters/repository/user_repository.go"),
            MongoDBRepositoryTemplate(projectName)); err != nil {
            return err
        }
        
        if err := writeFile(filepath.Join(projectPath, "internal/infrastructure/database/database.go"),
            MongoDBDatabaseTemplate(projectName)); err != nil {
            return err
        }
        
    default: // SQL databases (postgres, mysql, sqlite)
        // Utiliser templates GORM classiques
        if err := writeFile(filepath.Join(projectPath, "internal/domain/user.go"),
            t.UserModelTemplate()); err != nil {
            return err
        }
        
        // ... templates GORM
    }
    
    // ... reste de la génération
}
```

### 🧪 Stratégie de Tests

#### Tests Unitaires

```go
func TestGoModDependenciesMongoDB(t *testing.T) {
    deps := GoModDependencies("mongodb")
    assert.Contains(t, deps, "go.mongodb.org/mongo-driver")
    assert.NotContains(t, deps, "gorm.io/gorm", "MongoDB should not use GORM")
}

func TestMongoDBModelsTemplate(t *testing.T) {
    template := MongoDBModelsTemplate("testproject")
    assert.Contains(t, template, "primitive.ObjectID")
    assert.Contains(t, template, `bson:"_id"`)
    assert.NotContains(t, template, "gorm.Model")
}

func TestMongoDBRepositoryTemplate(t *testing.T) {
    template := MongoDBRepositoryTemplate("testproject")
    assert.Contains(t, template, "*mongo.Collection")
    assert.Contains(t, template, "InsertOne")
    assert.Contains(t, template, "FindOne")
    assert.NotContains(t, template, "*gorm.DB")
}
```

#### Tests E2E

```go
func TestE2EMongoDBProjectGeneration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping E2E test in short mode")
    }
    
    tmpDir := t.TempDir()
    projectName := "test-mongodb-project"
    
    // Générer avec --database=mongodb
    err := generateProjectFiles(tmpDir, projectName, "full", "mongodb")
    require.NoError(t, err)
    
    // Vérifier go.mod contient mongo-driver, PAS GORM
    goMod := readFile(t, filepath.Join(tmpDir, "go.mod"))
    assert.Contains(t, goMod, "go.mongodb.org/mongo-driver")
    assert.NotContains(t, goMod, "gorm.io/gorm")
    
    // Vérifier docker-compose contient MongoDB
    dockerCompose := readFile(t, filepath.Join(tmpDir, "docker-compose.yml"))
    assert.Contains(t, dockerCompose, "mongo:7.0")
    assert.Contains(t, dockerCompose, "MONGO_INITDB_ROOT_USERNAME")
    
    // Vérifier models utilisent bson tags
    userModel := readFile(t, filepath.Join(tmpDir, "internal/domain/user.go"))
    assert.Contains(t, userModel, "primitive.ObjectID")
    assert.Contains(t, userModel, `bson:"_id"`)
    
    // Vérifier repository utilise mongo.Collection
    userRepo := readFile(t, filepath.Join(tmpDir, "internal/adapters/repository/user_repository.go"))
    assert.Contains(t, userRepo, "*mongo.Collection")
    assert.Contains(t, userRepo, "InsertOne")
    
    // Vérifier compilation
    cmd := exec.Command("go", "build", "./...")
    cmd.Dir = tmpDir
    output, err := cmd.CombinedOutput()
    require.NoError(t, err, "MongoDB project should compile: %s", output)
}
```

### 🛡️ Developer Guardrails

#### Erreurs CRITIQUES à Éviter

**ERREUR #1: Utiliser GORM avec MongoDB**
- ❌ Inclure `gorm.io/gorm` dans go.mod
- ❌ Utiliser `*gorm.DB` dans repositories
- ✅ Utiliser `mongo-driver` UNIQUEMENT

**ERREUR #2: Tags gorm dans les models MongoDB**
- ❌ `gorm:"primaryKey"`
- ✅ `bson:"_id,omitempty"`

**ERREUR #3: Essayer de faire des migrations**
- ❌ `AutoMigrate` n'existe pas pour MongoDB
- ✅ MongoDB est schema-less, pas besoin de migrations

**ERREUR #4: Mauvais type pour ID**
- ❌ `ID uint`
- ✅ `ID primitive.ObjectID`

**ERREUR #5: Ne pas gérer le contexte**
- ❌ Méthodes repository sans `context.Context`
- ✅ Toutes les opérations MongoDB nécessitent un contexte

### 📋 Checklist de Validation

**Avant de marquer la story comme complète**:

- [ ] Fichier `templates_mongodb.go` créé avec tous les templates NoSQL
- [ ] `GoModDependencies("mongodb")` retourne mongo-driver, PAS GORM
- [ ] Template models MongoDB utilise `primitive.ObjectID` et tags `bson`
- [ ] Template repository MongoDB utilise `*mongo.Collection`
- [ ] Template database connection utilise `mongo.Connect()`
- [ ] Template docker-compose contient service MongoDB 7.0
- [ ] Pas de références à GORM dans les fichiers générés pour MongoDB
- [ ] Pas de migrations pour MongoDB (schema-less)
- [ ] Tests unitaires (5+) pour templates MongoDB passent
- [ ] Test E2E `TestE2EMongoDBProjectGeneration` passe
- [ ] Projet MongoDB généré compile sans erreurs
- [ ] README généré documente les différences NoSQL
- [ ] `make lint` passe
- [ ] Tous les tests SQL (postgres/mysql/sqlite) continuent de passer

### 🎓 Complexité et Considérations

#### Pourquoi cette story est OPTIONNELLE

**Raisons**:
1. **Architecture complètement différente** (NoSQL vs SQL)
2. **Templates entièrement nouveaux** (pas de réutilisation GORM)
3. **Complexité élevée** (doubles templates pour toutes les couches)
4. **Cas d'usage niche** (la plupart des users veulent SQL)

**Impact**:
- Estimation: **8-12 heures** (vs 4-6h pour MySQL/SQLite)
- Complexité: **Haute** (architecture différente)
- Priorité: **Basse** (optionnel, peut être dans future epic)

#### Approche Alternative Recommandée

**Si trop complexe**:
1. Marquer cette story comme "postponed"
2. Compléter Epic 7 avec postgres/mysql/sqlite uniquement
3. Créer Epic séparé "MongoDB Support" dans future sprint
4. Permettrait de se concentrer sur SQL d'abord (80% des cas)

### 💡 Conseils d'Implémentation

**Si vous décidez d'implémenter**:

1. Créer `templates_mongodb.go` avec TOUS les templates NoSQL
2. Tests unitaires pour chaque template MongoDB
3. Modifier `generateProjectFiles()` avec switch database
4. Pour "mongodb": utiliser templates MongoDB
5. Pour autres: utiliser templates GORM existants
6. Tests E2E exhaustifs
7. Documentation claire sur différences NoSQL/SQL

**Ordre recommandé**:
1. Models MongoDB (bson tags, ObjectID)
2. Repository MongoDB (Collection, BSON filters)
3. Database connection (mongo.Connect)
4. Docker-compose (service MongoDB)
5. Tests unitaires
6. Test E2E
7. Documentation

### 📖 Références

**Documentation MongoDB**:
- mongo-go-driver: https://www.mongodb.com/docs/drivers/go/current/
- BSON package: https://pkg.go.dev/go.mongodb.org/mongo-driver/bson
- MongoDB Docker: https://hub.docker.com/_/mongo

**Références projet**:
- [Source: _bmad-output/planning-artifacts/epics.md#Story 7.4] - Spécifications
- [Source: cmd/create-go-starter/templates_database.go] - Pattern établi pour SQL
- [Source: cmd/create-go-starter/templates.go] - Templates GORM actuels

## Dev Agent Record

### Agent Model Used

_À remplir par le dev agent_

### Completion Notes List

_À remplir par le dev agent_

### File List

_À remplir par le dev agent_

---

**Date de création**: 2026-02-05  
**Epic**: 7 - Multi-Database Support (v1.1.0)  
**Priorité**: Basse (OPTIONNEL)  
**Estimation**: 8-12 heures  
**Complexité**: Haute (architecture NoSQL fondamentalement différente)  
**Recommandation**: Considérer reporter à epic séparé si ressources limitées
