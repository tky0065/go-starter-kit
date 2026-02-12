# Bonnes pratiques

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Architecture

**1. Domain isolation**

Le domaine ne doit **jamais** importer d'autres packages:

```go
// ❌ BAD - Domain importing adapter
package user

import "mon-projet/internal/adapters/repository"  // NO!

// :material-check-circle: GOOD - Domain only imports interfaces
package user

import "mon-projet/internal/interfaces"
```

**2. Single Responsibility Principle**

Chaque composant a une seule responsabilité:

- **Handlers**: Parse + validate + call service
- **Services**: Business logic uniquement
- **Repositories**: Data access uniquement

**3. Dependency Injection**

Toujours via fx.Provide, pas de variables globales:

```go
// ❌ BAD - Global variable
var db *gorm.DB

// :material-check-circle: GOOD - Injection
type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}
```

### Code style

**1. gofmt**

Toujours formater:

```bash
go fmt ./...
```

Ou configurer l'IDE pour formater à la sauvegarde.

**2. golangci-lint**

Respecter les règles:

```bash
make lint
```

**3. Documentation GoDoc**

Pour les exports publics:

```go
// UserService handles user-related business logic.
// It provides methods for user registration, authentication, and CRUD operations.
type UserService struct {
    repo   interfaces.UserRepository
    logger zerolog.Logger
}

// Register creates a new user with the provided email and password.
// The password is automatically hashed before storage.
// Returns an error if the email already exists or if validation fails.
func (s *UserService) Register(ctx context.Context, email, password string) (*User, error) {
    // ...
}
```

**4. Error handling explicite**

Toujours gérer les erreurs, ne pas utiliser `panic`:

```go
// ❌ BAD
user := getUserByID(id)  // What if error?

// :material-check-circle: GOOD
user, err := getUserByID(id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

### Naming conventions

**Interfaces**:
- Suffixe `-er` ou `-Service`
- Exemples: `UserRepository`, `AuthService`, `Logger`

**Repositories**:
- Suffixe `-Repository`
- Exemples: `UserRepository`, `ProductRepository`

**Handlers**:
- Suffixe `-Handler`
- Exemples: `AuthHandler`, `UserHandler`

**Constructeurs**:
- Préfixe `New`
- Exemples: `NewUserService`, `NewAuthHandler`

**Méthodes privées**:
- lowerCamelCase
- Exemples: `hashPassword`, `validateEmail`

### Error handling patterns

**Wrap errors avec contexte**:

```go
// :material-check-circle: GOOD
if err != nil {
    return fmt.Errorf("failed to create user %s: %w", email, err)
}
```

**Domain errors pour logique métier**:

```go
if user == nil {
    return domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

**Ne pas gérer les HTTP status dans le service**:

```go
// ❌ BAD - Service returning HTTP status
func (s *UserService) GetByID(id uint) (int, *User, error) {
    return 404, nil, errors.New("not found")
}

// :material-check-circle: GOOD - Service returning domain error
func (s *UserService) GetByID(id uint) (*User, error) {
    return nil, domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

### Testing best practices

**1. Coverage > 80%**

```bash
go test -cover ./...
```

**2. Tests table-driven**

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid", "test", "TEST", false},
    {"empty", "", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := ToUpper(tt.input)
        if tt.wantErr {
            assert.Error(t, err)
        } else {
            assert.Equal(t, tt.want, got)
        }
    })
}
```

**3. Noms descriptifs**

```go
func TestUserService_Register_WhenEmailAlreadyExists_ReturnsConflictError(t *testing.T)
```

**4. Setup/teardown avec t.Cleanup()**

```go
func TestSomething(t *testing.T) {
    db := setupTestDB(t)
    t.Cleanup(func() {
        db.Exec("DELETE FROM users")
        db.Close()
    })

    // Test code
}
```

### Performance

**1. GORM - Éviter N+1 queries**

```go
// ❌ N+1 problem
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&posts)
}

// :material-check-circle: Single query with Preload
db.Preload("Posts").Find(&users)
```

**2. Context - Toujours passer context.Context**

```go
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
    return s.repo.FindByID(ctx, id)
}
```

**3. Database indexes**

```go
Email string `gorm:"uniqueIndex"`  // Index sur colonnes fréquemment requêtées
```

**4. Connection pooling**

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Sécurité recap

- [ ] Valider toutes les entrées utilisateur
- [ ] Jamais logger de passwords ou tokens
- [ ] Rate limiting sur endpoints publics
- [ ] HTTPS en production
- [ ] JWT secret fort (32+ caractères)
- [ ] Bcrypt pour passwords
- [ ] Mettre à jour les dépendances régulièrement

```bash
# Vérifier vulnérabilités
go list -json -m all | nancy sleuth
```

---

## Conclusion

Ce guide couvre tous les aspects du développement avec les projets générés par `create-go-starter`. Pour aller plus loin:

- **Exemples de code**: Tous les patterns sont dans le code généré
- **Tests**: Regardez les fichiers `*_test.go` pour des exemples
- **Documentation officielle**:
  - [Fiber](https://docs.gofiber.io/)
  - [GORM](https://gorm.io/docs/)
  - [fx](https://uber-go.github.io/fx/)
  - [zerolog](https://github.com/rs/zerolog)

**Bon développement!** :material-rocket-launch:

Si vous rencontrez des problèmes ou avez des questions, consultez:
- [Issues GitHub](https://github.com/tky0065/go-starter-kit/issues)
- [Discussions GitHub](https://github.com/tky0065/go-starter-kit/discussions)


---

## Navigation

**Previous**: [Monitoring](monitoring.md)  
**Index**: [Guide Index](index.md)
