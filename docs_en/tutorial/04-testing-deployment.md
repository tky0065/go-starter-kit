# Part 4: Testing and Deployment

<i class="material-icons small">circle</i> **Part 4/4** - Estimated time: 20 minutes

[<i class="material-icons">arrow_back</i> Back to index](index.md)

---

## Goal

In this final part, you will add unit tests and deploy your Blog API with Docker.

**What you will do**:
- Write unit tests for the Post service
- Run the tests
- Create a Docker image
- Deploy with docker-compose

---

## Step 12: Unit Tests

### 12.1 Test the Post service

Create `internal/domain/post/service_test.go`:

```go
package post_test

import (
	"context"
	"testing"

	"blog-api/internal/models"
	"blog-api/internal/interfaces/mocks"
	"blog-api/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostService_Create(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.PostRepository)
	log := logger.New(&config.Config{AppEnv: "test"})
	service := post.NewService(mockRepo, log)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Post")).
		Return(nil)

	// Act
	result, err := service.Create(context.Background(), 1, "Test Title", "Test Content", "tag1,tag2")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Title", result.Title)
	assert.Equal(t, "test-title", result.Slug)
	mockRepo.AssertExpectations(t)
}
```

### 12.2 Run the tests

```bash
make test
```

---

## Step 13: Docker Deployment

### 13.1 Build the Docker image

```bash
make docker-build
```

### 13.2 Launch with docker-compose

The `docker-compose.yml` file is already generated:

```bash
docker-compose up -d
```

This launches:
- The application on port 8080
- PostgreSQL on port 5432

### 13.3 Verify the deployment

```bash
curl http://localhost:8080/health
```

---

## Conclusion

Congratulations! You have built a complete Blog API with:

<i class="material-icons success">check_circle</i> **JWT Authentication** (User, Login, Register)
<i class="material-icons success">check_circle</i> **Articles** (Full CRUD with slug, tags, publish/unpublish)
<i class="material-icons success">check_circle</i> **Comments** (Create, List, Delete)
<i class="material-icons success">check_circle</i> **Relations** (Post → Author, Comment → Post + Author)
<i class="material-icons success">check_circle</i> **Pagination** (Limit/Offset)
<i class="material-icons success">check_circle</i> **Unit tests**
<i class="material-icons success">check_circle</i> **Docker deployment**
<i class="material-icons success">check_circle</i> **Hexagonal architecture**
<i class="material-icons success">check_circle</i> **Structured logging**
<i class="material-icons success">check_circle</i> **Centralized error handling**

### Summary of What You Learned

1. **Installation** of create-go-starter
2. **Generation** of a complete project
3. **Configuration** (.env, PostgreSQL, JWT)
4. **Hexagonal architecture**:
   - Domain (entities, services)
   - Adapters (handlers, repositories)
   - Interfaces (ports)
5. **Dependency Injection** with uber-go/fx
6. **GORM** (migrations, queries, relations)
7. **Fiber** (routes, middleware, handlers)
8. **Testing** with testify and mocks
9. **Docker** and docker-compose

### Next Steps

To go further:

- **Image uploads** for articles
- **Full-text search** in posts
- **Likes/Votes** on articles
- **Categories** to organize posts
- **Swagger** to document the API
- **CI/CD** with GitHub Actions
- **Kubernetes** for production deployment

### Resources

- [Generated Project Guide](../generated-project-guide.md) - Complete documentation
- [Example Repository](https://github.com/tky0065/go-starter-kit/tree/main/examples/blog-api) - Complete code
- [Fiber documentation](https://docs.gofiber.io/)
- [GORM documentation](https://gorm.io/docs/)

**Happy coding!** <i class="material-icons">rocket_launch</i>

---

## Navigation

[<i class="material-icons">arrow_back</i> Part 3: Expose the HTTP API](03-api-integration.md) | [<i class="material-icons">arrow_back</i> Back to index](index.md)
