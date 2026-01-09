# Story 4.2: Gestion centralisée des erreurs

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want un mécanisme uniforme pour formater les erreurs en JSON,
so that les clients de mon API reçoivent des réponses cohérentes en cas de problème et que je n'expose pas d'informations sensibles en production.

## Acceptance Criteria

1.  **Structure JSON Unifiée**
    *   **Given** Une erreur survient dans l'application (ex: 400 Bad Request, 404 Not Found, 500 Internal Error).
    *   **When** L'API renvoie la réponse.
    *   **Then** Le corps de la réponse suit STRICTEMENT ce format :
        ```json
        {
          "status": "error",
          "message": "Description lisible de l'erreur",
          "code": "ERROR_CODE_READABLE",
          "details": null // ou objet/array si validation errors
        }
        ```

2.  **Gestion des Erreurs Fiber & 404**
    *   **Given** Je demande une route qui n'existe pas (`/api/v1/ghost`).
    *   **Then** Je reçois une 404 avec le JSON standard.
    *   **Given** Une méthode non autorisée (POST sur un endpoint GET).
    *   **Then** Je reçois une 405 avec le JSON standard.

3.  **Masquage des Erreurs Internes (Sécurité)**
    *   **Given** Une panique ou une erreur 500 inattendue survient.
    *   **When** Je suis en **Production**.
    *   **Then** Le message retourné est générique ("Internal Server Error") et `details` est vide.
    *   **And** Aucune stack trace n'est visible.
    *   **When** Je suis en **Développement**.
    *   **Then** Le message réel de l'erreur peut être affiché pour le débogage.

4.  **Support des Erreurs Métier (Domain Errors)**
    *   **Given** Le code métier retourne une erreur typée (ex: `ErrEntityNotFound`).
    *   **Then** Le middleware la mappe automatiquement vers le bon statut HTTP (404) et le bon code d'erreur (`RESOURCE_NOT_FOUND`).

## Tasks / Subtasks

- [ ] **Core Error Module (`manual-test-project`)**
    - [ ] Créer `internal/domain/errors.go` : Définir la struct `AppError` et les fonctions helper (`NewNotFoundError`, `NewBadRequestError`, etc.).
    - [ ] Définir les interfaces/types pour les erreurs de validation.

- [ ] **Fiber Error Handler (`manual-test-project`)**
    - [ ] Créer `internal/adapters/middleware/error_handler.go` : Implémenter la fonction qui correspond à `fiber.ErrorHandler`.
    - [ ] Logique de mapping : `fiber.Error` -> JSON, `AppError` -> JSON, `unknown error` -> 500 JSON.
    - [ ] Intégration dans `internal/infrastructure/server/server.go` : Passer ce handler dans `fiber.Config`.

- [ ] **CLI Generator Update**
    - [ ] Mettre à jour les templates pour inclure ces nouveaux fichiers (`templates_domain.go`, `templates_middleware.go`).
    - [ ] Mettre à jour le template du serveur (`templates_server.go`) pour utiliser le ErrorHandler.

- [ ] **Validation & Testing**
    - [ ] Ajouter un test d'intégration qui provoque une 404 et vérifie le JSON.
    - [ ] Ajouter un test qui provoque une erreur métier et vérifie le mapping.

## Dev Notes

### Implementation Guide

**1. Domain Error Structure (`internal/domain/errors.go`)**
```go
package domain

type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Status  int    `json:"-"` // HTTP Status, not in JSON
    Details any    `json:"details,omitempty"`
}

func (e *AppError) Error() string { return e.Message }

// Helpers
func NewNotFoundError(msg string, code string) *AppError { ... }
```

**2. Fiber Error Handler (`internal/adapters/middleware/error_handler.go`)**
```go
func ErrorHandler(c *fiber.Ctx, err error) error {
    // Default 500
    code := fiber.StatusInternalServerError
    resp := fiber.Map{
        "status":  "error",
        "code":    "INTERNAL_SERVER_ERROR",
        "message": "Something went wrong",
        "details": nil,
    }

    // Handle Fiber Errors
    var e *fiber.Error
    if errors.As(err, &e) {
        code = e.Code
        resp["message"] = e.Message
        resp["code"] = "HTTP_ERROR" // Or map code to string
    }

    // Handle AppErrors
    var appErr *domain.AppError
    if errors.As(err, &appErr) {
        code = appErr.Status
        resp["message"] = appErr.Message
        resp["code"] = appErr.Code
        resp["details"] = appErr.Details
    }

    // TODO: Add logging here (Zerolog)

    return c.Status(code).JSON(resp)
}
```

### Architecture Compliance
- **FR15:** Fully addressed.
- **NFR7:** Protects against information leakage (Stack traces).
- **Structure:** Keeps domain pure (errors definition) and infrastructure specific (Fiber handler) separated in adapters.

### Integration Points
- `internal/infrastructure/server/server.go`: `app := fiber.New(fiber.Config{ ErrorHandler: middleware.ErrorHandler })`

## Dev Agent Record

### Agent Model Used
Gemini 2.0 Flash

### Debug Log References
- Verified Architecture.md requires specific JSON format.
- Verified Project Context requires use of `internal/domain/errors.go`.

### Completion Notes List
- [ ] Error types defined.
- [ ] Handler implemented and registered.
- [ ] CLI templates updated.
- [ ] Tests passed.

### File List
