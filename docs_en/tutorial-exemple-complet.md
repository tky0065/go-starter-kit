# Tutorial: Complete Example

A step-by-step tutorial to create and develop a complete project.

!!! note "Translation in progress"
    This page is being translated from French. For the complete documentation, please refer to the [French version](../tutorial-exemple-complet/).

## Overview

In this tutorial, we will:

1. Create a new project
2. Configure the environment
3. Test the API
4. Add a new feature
5. Deploy with Docker

## Step 1: Create the Project

```bash
# Install the tool
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest

# Create project
create-go-starter blog-api
cd blog-api
```

## Step 2: Configure Environment

```bash
# Automatic setup
./setup.sh

# Or manual setup
go mod tidy
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env

# Start PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_DB=blog-api \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

## Step 3: Run and Test

```bash
# Start the application
make run

# In another terminal, test the API
curl http://localhost:8080/health
# {"status":"ok"}

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@blog.com","password":"admin123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@blog.com","password":"admin123"}'
```

## Step 4: Add a New Feature (Posts)

### 4.1 Create Entity

```go
// internal/models/post.go
package models

import "time"

type Post struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"not null" json:"title"`
    Content   string    `gorm:"type:text" json:"content"`
    AuthorID  uint      `json:"author_id"`
    Author    User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 4.2 Add Migration

Update `internal/infrastructure/database/database.go`:

```go
db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Post{})
```

### 4.3 Create Interface

```go
// internal/interfaces/post_repository.go
package interfaces

type PostRepository interface {
    Create(ctx context.Context, post *models.Post) error
    FindByID(ctx context.Context, id uint) (*models.Post, error)
    FindAll(ctx context.Context) ([]models.Post, error)
    Update(ctx context.Context, post *models.Post) error
    Delete(ctx context.Context, id uint) error
}
```

## Step 5: Deploy with Docker

```bash
# Build the image
make docker-build

# Run with docker-compose
docker-compose up -d

# Check logs
docker-compose logs -f app
```

## Conclusion

You now have a complete blog API with:

- User authentication (JWT)
- CRUD operations
- PostgreSQL database
- Docker deployment
- CI/CD pipeline

## Next Steps

- Add more features (comments, categories, tags)
- Implement pagination
- Add Swagger documentation
- Set up monitoring
