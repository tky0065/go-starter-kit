# Database

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Migrations

The project uses GORM AutoMigrate to simplify migrations in development:

```go
// internal/infrastructure/database/database.go
db.AutoMigrate(
    &models.User{},
    &models.RefreshToken{},
)
```

**AutoMigrate**:
- Creates tables if they don't exist
- Adds missing columns
- Creates indexes
- **Does NOT delete** columns or tables

**For production**, consider a versioned migrations solution:

- **golang-migrate/migrate**: SQL or Go migrations
- **pressly/goose**: Up/down migrations
- **Advanced GORM Migrator**: Programmatic API

**Example with golang-migrate**:

```bash
# Install migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq create_users_table

# Created files:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql

# Run migrations
migrate -path migrations -database "postgresql://user:pass@localhost/dbname?sslmode=disable" up
```

### GORM Models

Conventions and patterns:

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

**Important GORM tags**:

| Tag | Description |
|-----|-------------|
| `primarykey` | Primary key |
| `uniqueIndex` | Unique index |
| `index` | Simple index |
| `not null` | NOT NULL column |
| `default:value` | Default value |
| `size:255` | Column size |
| `type:varchar(100)` | Custom SQL type |
| `foreignKey:UserID` | Foreign key |
| `references:ID` | FK reference |

**Conventions**:

- **Soft deletes**: `DeletedAt gorm.DeletedAt`
- **Auto timestamps**: `CreatedAt`, `UpdatedAt`
- **JSON hiding**: `json:"-"` for password
- **Index on FK**: Always index foreign keys

### Advanced Queries

#### Pagination

```go
var users []models.User
limit := 10
offset := 20

db.Limit(limit).Offset(offset).Find(&users)
```

#### Filtering

```go
// Simple where
db.Where("email = ?", "user@example.com").First(&user)

// Where with multiple conditions
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

### Performance Tips

1. **Index frequently queried columns**:
   ```go
   Email string `gorm:"uniqueIndex"`
   ```

2. **Avoid N+1 queries with Preload**:
   ```go
   // ❌ N+1
   for _, user := range users {
       db.Model(&user).Association("RefreshTokens").Find(&tokens)
   }

   // :material-check-circle: Single query
   db.Preload("RefreshTokens").Find(&users)
   ```

3. **Select only the necessary columns**:
   ```go
   db.Select("id, email").Find(&users)
   ```

4. **Use connection pools**:
   ```go
   sqlDB, _ := db.DB()
   sqlDB.SetMaxIdleConns(10)
   sqlDB.SetMaxOpenConns(100)
   sqlDB.SetConnMaxLifetime(time.Hour)
   ```

---


---

## Navigation

**Previous**: [Testing](testing.md)  
**Next**: [Security](security.md)  
**Index**: [Guide Index](index.md)
