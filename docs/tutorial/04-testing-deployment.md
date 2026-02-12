# Partie 4: Tests et Déploiement

<i class="material-icons small">circle</i> **Partie 4/4** - Temps estimé: 20 minutes

[<i class="material-icons">arrow_back</i> Retour à l'index](index.md)

---

## Objectif

Dans cette dernière partie, vous allez ajouter des tests unitaires et déployer votre API Blog avec Docker.

**Ce que vous allez faire**:
- Écrire des tests unitaires pour le service Post
- Lancer les tests
- Créer une image Docker
- Déployer avec docker-compose

---

## Étape 12: Tests unitaires

### 12.1 Tester le service Post

Créer `internal/domain/post/service_test.go`:

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

### 12.2 Lancer les tests

```bash
make test
```

---

## Étape 13: Déploiement Docker

### 13.1 Build l'image Docker

```bash
make docker-build
```

### 13.2 Lancer avec docker-compose

Le fichier `docker-compose.yml` est déjà généré:

```bash
docker-compose up -d
```

Cela lance:
- L'application sur le port 8080
- PostgreSQL sur le port 5432

### 13.3 Vérifier le déploiement

```bash
curl http://localhost:8080/health
```

---

## Conclusion

Félicitations! Vous avez créé une API Blog complète avec:

<i class="material-icons success">check_circle</i> **Authentification JWT** (User, Login, Register)
<i class="material-icons success">check_circle</i> **Articles** (CRUD complet avec slug, tags, publish/unpublish)
<i class="material-icons success">check_circle</i> **Commentaires** (Create, List, Delete)
<i class="material-icons success">check_circle</i> **Relations** (Post → Author, Comment → Post + Author)
<i class="material-icons success">check_circle</i> **Pagination** (Limit/Offset)
<i class="material-icons success">check_circle</i> **Tests unitaires**
<i class="material-icons success">check_circle</i> **Déploiement Docker**
<i class="material-icons success">check_circle</i> **Architecture hexagonale**
<i class="material-icons success">check_circle</i> **Logging structuré**
<i class="material-icons success">check_circle</i> **Error handling centralisé**

### Résumé de ce que vous avez appris

1. **Installation** de create-go-starter
2. **Génération** d'un projet complet
3. **Configuration** (.env, PostgreSQL, JWT)
4. **Architecture hexagonale**:
   - Domain (entities, services)
   - Adapters (handlers, repositories)
   - Interfaces (ports)
5. **Dependency Injection** avec uber-go/fx
6. **GORM** (migrations, queries, relations)
7. **Fiber** (routes, middleware, handlers)
8. **Tests** avec testify et mocks
9. **Docker** et docker-compose

### Prochaines étapes

Pour aller plus loin:

- **Upload d'images** pour les articles
- **Recherche full-text** dans les posts
- **Likes/Votes** sur les articles
- **Catégories** pour organiser les posts
- **Swagger** pour documenter l'API
- **CI/CD** avec GitHub Actions
- **Kubernetes** pour déploiement en production

### Ressources

- [Guide des projets générés](../generated-project-guide.md) - Documentation complète
- [Repository exemple](https://github.com/tky0065/go-starter-kit/tree/main/examples/blog-api) - Code complet
- [Fiber documentation](https://docs.gofiber.io/)
- [GORM documentation](https://gorm.io/docs/)

**Bon coding!** <i class="material-icons">rocket_launch</i>

---

## Navigation

[<i class="material-icons">arrow_back</i> Partie 3: Exposer l'API HTTP](03-api-integration.md) | [<i class="material-icons">arrow_back</i> Retour à l'index](index.md)
