# Story 6.1: Flag CLI pour sélection de template

Status: done

## Story

**En tant que** développeur,
**Je veux** pouvoir spécifier le type de template via un flag CLI,
**Afin de** générer le type de projet adapté à mes besoins.

## Acceptance Criteria

1. **AC1**: Given l'utilisateur exécute `create-go-starter mon-projet --template=minimal`, When le CLI parse les arguments, Then le template minimal est sélectionné
2. **AC2**: Given l'utilisateur exécute `create-go-starter mon-projet` (sans flag), When le CLI parse les arguments, Then le template "full" est utilisé par défaut
3. **AC3**: Given l'utilisateur exécute `create-go-starter mon-projet --template=invalid`, When le CLI parse les arguments, Then une erreur est affichée avec les valeurs valides
4. **AC4**: Given l'utilisateur exécute `create-go-starter --help`, When l'aide s'affiche, Then le flag --template est documenté avec les options disponibles

## Tasks / Subtasks

- [x] Task 1: Ajouter le flag `--template` dans main.go (AC: 1, 2)
  - [x] 1.1 Définir le flag avec `flag.StringVar` et valeur par défaut "full"
  - [x] 1.2 Documenter le flag dans `flag.Usage`
  - [x] 1.3 Parser le flag après `flag.Parse()`

- [x] Task 2: Créer la validation du template (AC: 3)
  - [x] 2.1 Créer une constante ou slice avec les valeurs valides: `minimal`, `full`, `graphql`
  - [x] 2.2 Créer fonction `validateTemplate(template string) error`
  - [x] 2.3 Retourner erreur descriptive si invalide avec liste des options

- [x] Task 3: Passer le template au générateur (AC: 1, 2)
  - [x] 3.1 Modifier signature de `run(projectName string)` → `run(projectName, template string)`
  - [x] 3.2 Modifier signature de `generateProjectFiles(projectPath, projectName string)` → ajouter `template string`
  - [x] 3.3 Afficher le template sélectionné dans les messages de progression

- [x] Task 4: Mettre à jour l'aide CLI (AC: 4)
  - [x] 4.1 Ajouter description du flag --template dans `flag.Usage`
  - [x] 4.2 Lister les 3 options avec descriptions courtes

- [x] Task 5: Tests unitaires
  - [x] 5.1 Test `validateTemplate` avec valeurs valides et invalides
  - [x] 5.2 Test flag parsing avec différentes combinaisons
  - [x] 5.3 Test erreur quand template invalide

## Dev Notes

### Architecture actuelle (main.go)

Le CLI utilise le package `flag` standard de Go:
- Ligne 113-128: Parsing des flags actuels (`--help`, `-h`)
- Ligne 140: Appel de `run(projectName)`
- Ligne 175: Appel de `generateProjectFiles(projectPath, projectName)`

### Pattern de flag à suivre

```go
// Définition du flag
var template string
flag.StringVar(&template, "template", "full", "Template type: minimal, full, graphql")

// Validation après parse
if err := validateTemplate(template); err != nil {
    fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("Error: %v", err)))
    os.Exit(1)
}

// Passage au générateur
if err := run(projectName, template); err != nil {
    // ...
}
```

### Templates valides

| Template | Description |
|----------|-------------|
| `minimal` | API REST basique avec Swagger, sans auth |
| `full` | Structure complète: JWT, CRUD users, Swagger (actuel) |
| `graphql` | API GraphQL avec gqlgen et Playground |

### Fichiers à modifier

1. `cmd/create-go-starter/main.go`
   - Ajouter flag `--template`
   - Ajouter fonction `validateTemplate()`
   - Modifier signatures de `run()` et `generateProjectFiles()`
   - Mettre à jour `flag.Usage()`

2. `cmd/create-go-starter/generator.go`
   - Modifier signature `generateProjectFiles(projectPath, projectName, template string)`
   - Pour l'instant, ignorer le paramètre template (sera utilisé dans story 6.2+)

3. `cmd/create-go-starter/main_test.go`
   - Ajouter tests pour `validateTemplate()`
   - Mettre à jour tests existants si signatures changent

### Project Structure Notes

- Le flag doit être cohérent avec le style existant (ex: `--help`)
- Les messages d'erreur doivent utiliser `Red()` pour la couleur
- Les messages de succès doivent utiliser `Green()`

### Contraintes

- Ne PAS implémenter la logique de génération par template (stories 6.2, 6.3, 6.4)
- Le template est passé mais ignoré pour l'instant (backward compatible)
- Les tests existants doivent continuer à passer

### References

- [Source: cmd/create-go-starter/main.go] - Structure actuelle du CLI
- [Source: _bmad-output/planning-artifacts/epic-6.md#Story 6.1] - Spécifications

## Dev Agent Record

### Agent Model Used

Claude Sonnet 4 (Anthropic)

### Completion Notes List

- **2026-01-15**: Story 6.1 implémentée avec succès
  - Ajouté le flag `--template` avec valeur par défaut "full" (AC1, AC2)
  - Créé `ValidTemplates` slice et `DefaultTemplate` constante
  - Implémenté `validateTemplate()` avec message d'erreur descriptif (AC3)
  - Modifié les signatures de `run()` et `generateProjectFiles()` pour accepter le paramètre template
  - Mis à jour `flag.Usage()` avec section Templates documentée (AC4)
  - Ajouté 8 nouveaux tests couvrant tous les ACs:
    - `TestValidateTemplateValid` - templates valides
    - `TestValidateTemplateInvalid` - templates invalides avec vérification du message
    - `TestTemplateDefaultValue` - valeur par défaut
    - `TestValidTemplatesContains` - liste des templates
    - `TestTemplateFlagParsing` - parsing avec --template=X
    - `TestTemplateFlagDefault` - sans flag = "full"
    - `TestInvalidTemplateFlagError` - erreur sur template invalide
    - `TestHelpShowsTemplateFlag` - documentation dans --help
  - Tous les tests existants mis à jour pour utiliser la nouvelle signature
  - Le paramètre template est passé mais ignoré (backward compatible) - sera utilisé dans stories 6.2+

### File List

- cmd/create-go-starter/main.go (modified)
- cmd/create-go-starter/generator.go (modified)
- cmd/create-go-starter/main_test.go (modified)
- cmd/create-go-starter/generator_test.go (modified)
- cmd/create-go-starter/scaffold_test.go
- cmd/create-go-starter/templates_test.go
- cmd/create-go-starter/git_test.go (modified)
- cmd/create-go-starter/smoke_test.go (modified)

## Change Log

- 2026-01-15: Implémentation complète de la story 6.1 - Flag CLI pour sélection de template
