# Base de données

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Migrations

Le projet utilise GORM AutoMigrate pour simplifier les migrations en développement:

```go
// internal/infrastructure/database/database.go
db.AutoMigrate(
    &models.User{},
    &models.RefreshToken{},
)
```

**AutoMigrate**:
- Crée les tables si elles n'existent pas
- Ajoute les colonnes manquantes
- Crée les indexes
- **NE supprime PAS** de colonnes ou tables

**Pour la production**, considérer une solution de migrations versionnées:

- **golang-migrate/migrate**: Migrations SQL ou Go
- **pressly/goose**: Migrations up/down
- **GORM Migrator avancé**: API programmatique

**Exemple avec golang-migrate**:

```bash
# Installer migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Créer migration
migrate create -ext sql -dir migrations -seq create_users_table

# Fichiers créés:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql

# Run migrations
migrate -path migrations -database "postgresql://user:pass@localhost/dbname?sslmode=disable" up
```

### Modèles GORM

Conventions et patterns:

```go
type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    Email    string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
    Password string `gorm:"not null" json:"-"`
}
```

**Tags GORM importants**:

| Tag | Description |
|-----|-------------|
| `primarykey` | Clé primaire |
| `uniqueIndex` | Index unique |
| `index` | Index simple |
| `not null` | Colonne NOT NULL |
| `default:value` | Valeur par défaut |
| `size:255` | Taille de colonne |
| `type:varchar(100)` | Type SQL custom |
| `foreignKey:UserID` | Clé étrangère |
| `references:ID` | Référence FK |

**Conventions**:

- **Soft deletes**: `DeletedAt gorm.DeletedAt`
- **Timestamps auto**: `CreatedAt`, `UpdatedAt`
- **JSON hiding**: `json:"-"` pour password
- **Index sur FK**: Toujours indexer les clés étrangères

### Queries avancées

#### Pagination

```go
var users []models.User
limit := 10
offset := 20

db.Limit(limit).Offset(offset).Find(&users)
```

#### Filtering

```go
// Where simple
db.Where("email = ?", "user@example.com").First(&user)

// Where avec multiple conditions
db.Where("created_at > ? AND email LIKE ?", time.Now().Add(-24*time.Hour), "%@example.com").Find(&users)

// Or
db.Where("email = ?", email1).Or("email = ?", email2).Find(&users)
```

#### Sorting

```go
// Order ASC
db.Order("created_at asc").Find(&users)

// Order DESC
db.Order("created_at desc").Find(&users)

// Multiple sorts
db.Order("created_at desc, email asc").Find(&users)
```

#### Joins

```go
// Inner join
db.Joins("LEFT JOIN refresh_tokens ON refresh_tokens.user_id = users.id").
   Where("refresh_tokens.expires_at > ?", time.Now()).
   Find(&users)

// Preload associations
db.Preload("RefreshTokens").Find(&users)
```

#### Aggregations

```go
// Count
var count int64
db.Model(&models.User{}).Count(&count)

// With where
db.Model(&models.User{}).Where("created_at > ?", yesterday).Count(&count)
```

#### Transactions

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Create user
    if err := tx.Create(&user).Error; err != nil {
        return err  // Rollback
    }

    // Create profile
    if err := tx.Create(&profile).Error; err != nil {
        return err  // Rollback
    }

    return nil  // Commit
})
```

#### Raw SQL

```go
// Raw query
var users []models.User
db.Raw("SELECT * FROM users WHERE email LIKE ?", "%@example.com").Scan(&users)

// Exec
db.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, userID)
```

### Performance tips

1. **Index les colonnes fréquemment requêtées**:
   ```go
   Email string `gorm:"uniqueIndex"`
   ```

2. **Éviter N+1 queries avec Preload**:
   ```go
   // ❌ N+1
   for _, user := range users {
       db.Model(&user).Association("RefreshTokens").Find(&tokens)
   }

   // :material-check-circle: Single query
   db.Preload("RefreshTokens").Find(&users)
   ```

3. **Select seulement les colonnes nécessaires**:
   ```go
   db.Select("id, email").Find(&users)
   ```

4. **Utiliser les connections pools**:
   ```go
   sqlDB, _ := db.DB()
   sqlDB.SetMaxIdleConns(10)
   sqlDB.SetMaxOpenConns(100)
   sqlDB.SetConnMaxLifetime(time.Hour)
   ```

---


---

## Navigation

**Previous**: [Tests](testing.md)  
**Next**: [Sécurité](security.md)  
**Index**: [Guide Index](index.md)
