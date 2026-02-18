# Best Practices

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Architecture

**1. Domain isolation**

The domain must **never** import other packages:

```go
// ❌ BAD - Domain importing adapter
package user

import "mon-projet/internal/adapters/repository"  // NO!

// :material-check-circle: GOOD - Domain only imports interfaces
package user

import "mon-projet/internal/interfaces"
```

**2. Single Responsibility Principle**

Each component has a single responsibility:

- **Handlers**: Parse + validate + call service
- **Services**: Business logic only
- **Repositories**: Data access only

**3. Dependency Injection**

Always via fx.Provide, no global variables:

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

Always format:

```bash
go fmt ./...
```

Or configure the IDE to format on save.

**2. golangci-lint**

Follow the rules:

```bash
make lint
```

**3. GoDoc Documentation**

For public exports:

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

**4. Explicit error handling**

Always handle errors, do not use `panic`:

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
- Suffix `-er` or `-Service`
- Examples: `UserRepository`, `AuthService`, `Logger`

**Repositories**:
- Suffix `-Repository`
- Examples: `UserRepository`, `ProductRepository`

**Handlers**:
- Suffix `-Handler`
- Examples: `AuthHandler`, `UserHandler`

**Constructors**:
- Prefix `New`
- Examples: `NewUserService`, `NewAuthHandler`

**Private methods**:
- lowerCamelCase
- Examples: `hashPassword`, `validateEmail`

### Error handling patterns

**Wrap errors with context**:

```go
// :material-check-circle: GOOD
if err != nil {
    return fmt.Errorf("failed to create user %s: %w", email, err)
}
```

**Domain errors for business logic**:

```go
if user == nil {
    return domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

**Do not handle HTTP status codes in the service**:

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

**2. Table-driven tests**

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

**3. Descriptive names**

```go
func TestUserService_Register_WhenEmailAlreadyExists_ReturnsConflictError(t *testing.T)
```

**4. Setup/teardown with t.Cleanup()**

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

**1. GORM - Avoid N+1 queries**

```go
// ❌ N+1 problem
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&posts)
}

// :material-check-circle: Single query with Preload
db.Preload("Posts").Find(&users)
```

**2. Context - Always pass context.Context**

```go
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
    return s.repo.FindByID(ctx, id)
}
```

**3. Database indexes**

```go
Email string `gorm:"uniqueIndex"`  // Index on frequently queried columns
```

**4. Connection pooling**

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Security recap

- [ ] Validate all user inputs
- [ ] Never log passwords or tokens
- [ ] Rate limiting on public endpoints
- [ ] HTTPS in production
- [ ] Strong JWT secret (32+ characters)
- [ ] Bcrypt for passwords
- [ ] Update dependencies regularly

```bash
# Check for vulnerabilities
go list -json -m all | nancy sleuth
```

---

## Conclusion

This guide covers all aspects of development with projects generated by `create-go-starter`. To go further:

- **Code examples**: All patterns are in the generated code
- **Tests**: Look at `*_test.go` files for examples
- **Official documentation**:
  - [Fiber](https://docs.gofiber.io/)
  - [GORM](https://gorm.io/docs/)
  - [fx](https://uber-go.github.io/fx/)
  - [zerolog](https://github.com/rs/zerolog)

**Happy coding!** <i class="material-icons info">rocket_launch</i>

If you encounter problems or have questions, check:
- [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)
- [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)


---

## Navigation

**Previous**: [Monitoring](monitoring.md)  
**Index**: [Guide Index](index.md)
