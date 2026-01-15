# Story 6.4: Template GraphQL

**Status:** ready-for-dev
**Epic:** 6 (Templates Multiples)
**Story:** 6.4

## 1. User Story

**En tant que** développeur,
**Je veux** générer un projet Go avec une API GraphQL pré-configurée,
**Afin de** démarrer rapidement le développement d'APIs modernes sans configurer manuellement gqlgen et Fiber.

## 2. Acceptance Criteria

- **AC1:** **Given** l'utilisateur exécute `create-go-starter my-graphql-api --template=graphql`
- **When** la commande se termine avec succès
- **Then** la structure du projet contient les dossiers `graph/`, `graph/model/` et le fichier `gqlgen.yml`
- **And** le fichier `cmd/main.go` est configuré pour servir le Playground GraphQL et les requêtes `/query`

- **AC2:** **Given** le projet généré est lancé avec `make dev`
- **When** j'accède à `http://localhost:8080/` (ou endpoint configuré)
- **Then** le Playground GraphQL s'affiche et permet d'exécuter des requêtes

- **AC3:** **Given** le projet généré
- **When** j'exécute une requête GraphQL `query { user(id: "1") { name } }` (exemple)
- **Then** je reçois une réponse JSON valide conforme au schéma
- **And** la base de données est correctement interrogée (simulé ou réel via GORM)

- **AC4:** **Structure Scaffolding**
- **Then** le projet inclut les dépendances `github.com/99designs/gqlgen` et `github.com/gofiber/adaptor/v2` (pour l'intégration http.Handler -> Fiber)
- **And** l'architecture "Full" (Auth, etc.) est allégée : pas de JWT complexe par défaut si non spécifié, mais focus sur GraphQL (Note: Le PRD Epic 6 dit "API GraphQL avec gqlgen... GORM, Logger". Il n'explicite pas l'auth, mais le modèle Full a l'auth. Pour simplifier et suivre le "Template Minimal" vs "Full", le "GraphQL" semble être une variante standalone. Assumons une structure fonctionnelle de base proche du Minimal + GraphQL + GORM).

## 3. Technical Requirements

### 3.1 Dependencies
- **Generator**: Le CLI `create-go-starter` doit être mis à jour.
- **Library**: `github.com/99designs/gqlgen` (v0.17.84 ou récent).
- **Adapter**: `github.com/gofiber/adaptor/v2` pour wrapper le handler GraphQL standard dans Fiber.

### 3.2 Implementation Details (Generator Side)
- **New Template Source**: Créer un nouveau fichier `cmd/create-go-starter/templates_graphql.go` pour stocker les templates spécifiques :
    - `GqlGenYmlTemplate()`
    - `GraphSchemaTemplate()`
    - `GraphResolverTemplate()`
    - `GraphModelTemplate()` (Note: Généré manuellement pour le starter pour éviter d'exiger `go run` à l'init, ou minimaliste).
    - `MainGoGraphQLTemplate()` : Version modifiée de `main.go` qui initialise le handler GraphQL au lieu des routes REST classiques.

- **Refactoring Integration (Dependency on Story 6.3)**:
    - Implémenter `getGraphQLFiles(t *ProjectTemplates, projectPath string) []FileGenerator` dans `generator.go`.
    - Cette fonction doit retourner la liste des fichiers :
        - `gqlgen.yml`
        - `graph/schema.graphqls`
        - `graph/schema.resolvers.go`
        - `graph/model/models_gen.go` (Optionnel: Peut être généré via `generate` post-install, mais mieux vaut le fournir pour "Zero-to-Hero").
        - `cmd/main.go` (Spécifique GraphQL)
        - `internal/infrastructure/server/server.go` (Adapté pour GraphQL si nécessaire, ou configurer les routes dans `main`).

### 3.3 Generated Project Architecture
- **Adapter Layer**: Le handler GraphQL agit comme un Adapter.
- **Server**: Fiber reste le serveur HTTP.
    ```go
    // Exemple d'intégration simplifiée dans le main généré ou server.go
    import (
        "github.com/99designs/gqlgen/graphql/handler"
        "github.com/99designs/gqlgen/graphql/playground"
        "github.com/gofiber/adaptor/v2"
        "github.com/gofiber/fiber/v2"
    )

    // ...
    srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
    app.All("/query", adaptor.HTTPHandler(srv))
    app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))
    ```
- **Domain/Infrastructure**: Réutilisation de GORM et `pkg/logger` comme dans les autres templates.

## 4. Developer Context & Guardrails

### 4.1 Development Strategy
1.  **Create Templates First**: Définir les chaînes de caractères pour `gqlgen.yml`, `schema.graphqls` (ex: type User), et les résolveurs de base.
2.  **Mock Generation**: Exécuter `gqlgen init` localement dans un dossier temporaire pour obtenir, nettoyer et copier les fichiers générés dans `templates_graphql.go`.
    - *Pourquoi ?* Pour garantir que le projet généré est valide syntaxiquement sans forcer l'utilisateur à debugger des fichiers générés.
3.  **Integrate Generator**: Brancher `getGraphQLFiles` dans le switch du `generator.go` (créé en Story 6.3).

### 4.2 Architecture Compliance
- **Fx Injection**: Même en mode GraphQL, utiliser `fx` pour l'injection si possible, ou rester simple si le template se veut minimaliste. *Décision*: Le PRD "Full" mentionne `fx`. Le template GraphQL doit conserver l'ADN du starter kit (Fx + Fiber + GORM).
    - Le `Resolver` GraphQL doit probablement recevoir les Services/Repositories via injection.
    - `graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}})`

### 4.3 Anti-Patterns to Avoid
- **Mixing REST & GraphQL loosely**: Ne pas mélanger les structure de dossiers inutilement. Mettre tout ce qui est GraphQL dans `graph/` (convention gqlgen).
- **Manual wiring hell**: Ne pas essayer de réinventer la roue `adaptor`. Utiliser `gofiber/adaptor/v2`.
- **Empty Resolver**: Fournir un resolver exemple (ex: `Users` query) qui fonctionne pour que l'utilisateur voit immédiatement un résultat.

## 5. Testing Requirements

- **Unit Test (Generator)**: Ajouter un test case dans `generator_test.go` pour `--template=graphql`. Vérifier que les fichiers clés (`gqlgen.yml`, `schema.graphqls`) sont créés.
- **Smoke Test**: Mettre à jour `scripts/smoke_test.sh` ou créer un step spécifique pour générer un projet graphql, lancer `go mod tidy` et vérifier que ça compile (`go build`).

## 6. Previous Story Intelligence (from 6.3)
- La story 6.3 a introduit l'architecture de sélection de template. Vous DEVEZ suivre la signature `func getGraphQLFiles(...)` définie dans le plan de refactoring.
- Ne modifiez pas la logique principale de validation du nom de projet, réutilisez l'existant.

## 7. Web Research Notes
- **Gqlgen Integration**: Utilise `net/http`. Fiber utilise `fasthttp`. L'adaptateur est OBLIGATOIRE.
- **Playground**: `playground.Handler` est un standard `http.Handler`, nécessite aussi `adaptor.HTTPHandler`.

## 8. Definition of Done
- [ ] Code du CLI mis à jour avec `templates_graphql.go`
- [ ] Logique de génération implémentée pour le cas "graphql"
- [ ] Projet généré compile et lance un serveur GraphQL fonctionnel
- [ ] Playground accessible à la racine ou `/playground`
- [ ] Tests unitaires et smoke tests passent
