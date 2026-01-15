# Story 6.5: Documentation et aide CLI

**Status:** ready-for-dev
**Epic:** 6 (Templates Multiples)
**Story:** 6.5

## 1. User Story

**En tant que** développeur,
**Je veux** trouver une documentation claire sur les options de templates disponibles (minimal, full, graphql),
**Afin de** choisir le point de départ le plus adapté à mon nouveau projet.

## 2. Acceptance Criteria

- **AC1:** **Given** je consulte le `README.md`
- **When** je cherche la section "Utilisation"
- **Then** je trouve l'explication du flag `--template` et la liste des valeurs possibles (`minimal`, `full`, `graphql`).

- **AC2:** **Given** je consulte `docs/usage.md`
- **When** je lis la section "Options disponibles"
- **Then** elle est à jour et liste les templates avec une brève description de ce que chacun contient.
- **And** la note "D'autres options seront ajoutées" est retirée ou mise à jour.

- **AC3:** **Given** j'exécute `create-go-starter --help`
- **When** la commande s'exécute
- **Then** la sortie correspond à la documentation (liste des templates et descriptions).

- **AC4:** **Roadmap**
- **Then** le point "Templates multiples" dans la Roadmap du `README.md` est marqué comme terminé (`[x]`).

## 3. Technical Requirements

### 3.1 Documentation Updates
- **README.md**:
    - Mettre à jour la section "Utilisation de base" pour inclure `--template`.
    - Mettre à jour la Roadmap.
- **docs/usage.md**:
    - Ajouter une section détaillée sur les templates.
    - Expliquer les différences structurelles (ex: Minimal n'a pas d'auth, GraphQL a `graph/`).
- **docs/cli-architecture.md** (Optionnel):
    - Vérifier si le diagramme ou l'explication du "Generator" mentionne la sélection de stratégie de template. Si non, ajouter une phrase explicative.

### 3.2 CLI Help Verification (Code)
- Le code dans `main.go` implémente déjà l'aide.
- **Task**: Vérifier que les descriptions dans `main.go` (`TemplateMinimalDesc`, etc.) sont cohérentes avec celles ajoutées dans le README.
    - Minimal: "Basic REST API with Swagger (no authentication)"
    - Full: "Complete API with JWT auth, user management, and Swagger (default)"
    - GraphQL: "GraphQL API with gqlgen and GraphQL Playground"

## 4. Developer Context & Guardrails

### 4.1 Consistency
- Assurez-vous que les termes utilisés dans la doc ("minimal", "full", "graphql") correspondent EXACTEMENT aux valeurs du flag dans le code.
- Ne pas inventer de features non implémentées (ex: ne pas documenter "api-only" si ce n'est pas fait).

### 4.2 LLM Optimization Strategy
- **Format**: Utiliser des tableaux Markdown pour comparer les templates dans `docs/usage.md` (Features vs Template).
    - Ex: | Feature | Minimal | Full | GraphQL |
          |---------|---------|------|---------|
          | Auth    | ❌      | ✅   | ❌ (TBD)|
          | Swagger | ✅      | ✅   | ❌      |
          | GQL     | ❌      | ❌   | ✅      |

## 5. Previous Story Intelligence
- **Story 6.1** a ajouté le code du flag.
- **Stories 6.2, 6.3, 6.4** ont implémenté (ou sont en train d'implémenter) les templates.
- Cette story 6.5 est la "finition" de l'Epic pour garantir que l'utilisateur sait utiliser ce qui a été construit.

## 6. Definition of Done
- [ ] `README.md` mis à jour (Usage + Roadmap)
- [ ] `docs/usage.md` mis à jour avec section Templates
- [ ] Vérification manuelle que `create-go-starter --help` est synchronisé avec la doc

