# Story 10.5: Raccourcis (Aliases) pour les Options CLI

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want utiliser des raccourcis courts pour toutes les options CLI (ex: `-i` pour `--interactive`, `-t` pour `--template`),
so that je puisse saisir les commandes plus rapidement et de façon moins verbeuse.

## Acceptance Criteria

1. **Given** j'exécute `create-go-starter -i`, **When** le CLI parse les arguments, **Then** c'est équivalent à `create-go-starter --interactive`.
2. **Given** j'exécute `create-go-starter -t minimal my-project`, **When** le CLI parse les arguments, **Then** c'est équivalent à `--template=minimal`.
3. **Given** j'exécute `create-go-starter -d sqlite my-project`, **When** le CLI parse les arguments, **Then** c'est équivalent à `--database=sqlite`.
4. **Given** j'exécute `create-go-starter -o advanced my-project`, **When** le CLI parse les arguments, **Then** c'est équivalent à `--observability=advanced`.
5. **Given** j'exécute `create-go-starter -n my-project`, **When** le CLI parse les arguments, **Then** c'est équivalent à `--dry-run`.
6. **Given** j'exécute `create-go-starter -h`, **When** le CLI parse les arguments, **Then** c'est équivalent à `--help` (déjà implémenté — vérifier que ça marche toujours).
7. **Given** les raccourcis sont ajoutés, **When** j'affiche l'aide (`-h` ou `--help`), **Then** chaque option affiche son raccourci à côté du nom long (ex: `-t, --template`).
8. **Given** des formes courtes et longues sont combinées (ex: `-d sqlite -t minimal my-project`), **When** le CLI parse les arguments, **Then** toutes les formes sont correctement résolues.
9. **Given** un raccourci inconnu est utilisé (ex: `-x`), **When** le CLI parse les arguments, **Then** un message d'erreur clair est affiché: `unknown flag: -x`.

## Table des raccourcis

| Option longue | Raccourci court | Valeur par défaut |
|---------------|-----------------|-------------------|
| `--help` | `-h` | — |
| `--interactive` | `-i` | false |
| `--template` | `-t` | `full` |
| `--database` | `-d` | `postgres` |
| `--observability` | `-o` | `none` |
| `--dry-run` | `-n` | false |

> **Note `-h`** : Déjà partiellement implémenté (`-h`, `-help`, `--help`). S'assurer de la cohérence avec le nouveau système.

## Tasks / Subtasks

- [x] Task 1 : Mettre à jour la boucle de parsing des arguments dans `main()` pour supporter les raccourcis (AC: #1-#6, #8, #9)
  - [x] Ajouter détection `-i` (alias `--interactive`) dans la boucle
  - [x] Ajouter détection `-t`/`-t=value` (alias `--template`)
  - [x] Ajouter détection `-d`/`-d=value` (alias `--database`)
  - [x] Ajouter détection `-o`/`-o=value` (alias `--observability`)
  - [x] Ajouter détection `-n` (alias `--dry-run`)
  - [x] Vérifier que `-h` fonctionne encore (déjà présent en `main.go:218`)
  - [x] Ajouter gestion des flags inconnus avec message d'erreur (AC: #9)
- [x] Task 2 : Mettre à jour la fonction `usage()` pour afficher les raccourcis (AC: #7)
  - [x] Modifier chaque ligne d'option pour afficher le format `-x, --long-option`
  - [x] Aligner les colonnes pour la lisibilité
- [x] Task 3 : Mettre à jour les exemples dans `usage()` (AC: #7)
  - [x] Ajouter des exemples avec les formes courtes
  - [x] Exemples: `create-go-starter -d sqlite -t minimal my-project`
- [x] Task 4 : Tests (AC: #1-#9)
  - [x] Créer ou mettre à jour `cmd/create-go-starter/main_test.go`
  - [x] Tester chaque raccourci court individuellement
  - [x] Tester la combinaison de raccourcis courts et longs
  - [x] Tester le message d'erreur pour flag inconnu
  - [x] Tester la rétrocompatibilité : les formes longues fonctionnent toujours

## Dev Notes

### Architecture et Patterns

**Fichier cible principal :** `cmd/create-go-starter/main.go` uniquement (boucle de parsing + `usage()`)

**Boucle de parsing actuelle (main.go:215-251) :** La boucle itère sur `os.Args[1:]` avec des conditions `strings.HasPrefix`. Le pattern est simple et facile à étendre.

**Pattern d'extension — ajouter les raccourcis dans la boucle existante :**

```go
// Avant (exemple pour --template) :
} else if strings.HasPrefix(arg, "-template=") || strings.HasPrefix(arg, "--template=") {
    parts := strings.SplitN(arg, "=", 2)
    if len(parts) == 2 { template = parts[1] }
} else if (arg == "-template" || arg == "--template") && i+1 < len(args) {
    template = args[i+1]; i++
}

// Après (ajouter les alias courts) :
} else if strings.HasPrefix(arg, "-template=") || strings.HasPrefix(arg, "--template=") ||
          strings.HasPrefix(arg, "-t=") {
    parts := strings.SplitN(arg, "=", 2)
    if len(parts) == 2 { template = parts[1] }
} else if (arg == "-template" || arg == "--template" || arg == "-t") && i+1 < len(args) {
    template = args[i+1]; i++
}
```

**Appliquer le même pattern pour chaque flag :**

| Flag | Court | Forme `=` à détecter | Forme espace à détecter |
|------|-------|----------------------|-------------------------|
| interactive | `-i` | (boolean, pas de valeur) | — |
| template | `-t` | `-t=value` | `-t value` |
| database | `-d` | `-d=value` | `-d value` |
| observability | `-o` | `-o=value` | `-o value` |
| dry-run | `-n` | (boolean, pas de valeur) | — |
| help | `-h` | (boolean, déjà implémenté) | — |

**Gestion des flags inconnus (AC: #9) :**

Actuellement, les flags inconnus sont silencieusement ignorés (la boucle passe au suivant sans erreur). Ajouter une détection en fin de boucle :

```go
// À la fin du bloc else-if, ajouter :
} else if strings.HasPrefix(arg, "-") {
    // Unknown flag
    fmt.Fprintln(os.Stderr, Red(fmt.Sprintf("unknown flag: %s", arg)))
    usage()
    os.Exit(1)
} else if !strings.HasPrefix(arg, "-") && projectName == "" {
    projectName = arg
}
```

**ATTENTION — Ordre des conditions dans la boucle :** Les conditions `strings.HasPrefix(arg, "-")` plus génériques doivent être en **dernier** dans le else-if chain, sinon elles absorberaient les flags connus. Vérifier l'ordre après modification.

**ATTENTION — Flags avec valeur vs flags booléens :** Les flags `-i` et `-n` sont booléens (pas de valeur suivante). Les flags `-t`, `-d`, `-o` nécessitent une valeur (consomment `args[i+1]`). Ne pas consommer `args[i+1]` pour les flags booléens.

### Mise à jour de `usage()` (AC: #7)

**Avant :**
```
Options:
  -database string
        Database type to use (default "postgres")
  -template string
        Template type to generate (default "full")
  -observability string
        Observability level: none|basic|advanced (default "none")
  -h, -help
        Show help message
```

**Après :**
```
Options:
  -d, --database string
        Database type to use (default "postgres")
  -t, --template string
        Template type to generate (default "full")
  -o, --observability string
        Observability level: none|basic|advanced (default "none")
  -i, --interactive
        Launch interactive mode (guided configuration)
  -n, --dry-run
        Preview files without creating them
  -h, --help
        Show help message
```

### Mise à jour des exemples dans `usage()`

**Ajouter ces exemples :**
```
Examples:
  create-go-starter my-project
  create-go-starter -d sqlite my-project
  create-go-starter my-project -t minimal
  create-go-starter -d mysql -t minimal my-project
  create-go-starter my-project --observability=advanced
  create-go-starter -d sqlite -t minimal -o none my-project    # Formes courtes
  create-go-starter -i                                          # Mode interactif
  create-go-starter -n my-project                               # Dry-run preview
  create-go-starter doctor                                      # Diagnostics
```

### Rétrocompatibilité

**CRITIQUE :** Tous les flags existants (formes longues) doivent continuer à fonctionner exactement comme avant. Les raccourcis sont des **aliases additionnels**, pas des remplacements.

Les tests existants dans `main_test.go` (s'ils existent) doivent tous passer sans modification.

### Cas limites à gérer

1. **`-t=minimal` vs `-t minimal`** : Les deux formes doivent fonctionner.
2. **`-tminimal`** (sans espace ni `=`) : **Ne pas supporter** — trop ambigu. Uniquement `-t minimal` et `-t=minimal`.
3. **Flags booléens avec `=`** : `-i=true` ou `-n=false` — **Ne pas supporter**, uniquement `-i` et `-n` (boolean flags sans valeur).
4. **Flag court inconnu suivi d'une valeur** : `-x sqlite` → erreur sur `-x`, ne pas consommer `sqlite` comme projectName.

### Project Structure Notes

- **Fichier modifié uniquement :** `cmd/create-go-starter/main.go` (boucle parsing + `usage()`)
- **Fichier test modifié/créé :** `cmd/create-go-starter/main_test.go`
- **Aucun nouveau fichier** requis
- **Aucune nouvelle dépendance** — modifications purement dans la boucle de parsing stdlib

### Testing Approach

```go
// Test chaque raccourci court
func TestShortFlagTemplate(t *testing.T) {
    // Simuler: create-go-starter -t minimal my-project
    // Vérifier que template = "minimal"
}

func TestShortFlagDatabase(t *testing.T) {
    // -d sqlite → database = "sqlite"
}

func TestShortFlagObservability(t *testing.T) {
    // -o advanced → observability = "advanced"
}

func TestShortFlagInteractive(t *testing.T) {
    // -i → interactive = true
}

func TestShortFlagDryRun(t *testing.T) {
    // -n → dryRun = true
}

func TestCombinedShortFlags(t *testing.T) {
    // -d sqlite -t minimal → database = "sqlite", template = "minimal"
}

func TestMixedShortAndLongFlags(t *testing.T) {
    // -d sqlite --template=full → database = "sqlite", template = "full"
}

func TestUnknownFlagError(t *testing.T) {
    // -x → exit 1 avec message d'erreur "unknown flag: -x"
}

func TestLongFlagsStillWork(t *testing.T) {
    // --database=postgres --template=full → rétrocompatibilité OK
}
```

**Note sur le testing de la boucle de parsing :** La boucle est dans `main()`. Pour la tester unitairement, extraire la logique de parsing dans une fonction `parseArgs(args []string) (Config, error)` qui retourne une struct `Config{ ProjectName, Template, Database, Observability string; Interactive, DryRun, Help bool }`. Cela rend les tests simples sans avoir à forker des processus.

**Refactoring optionnel mais recommandé :** Extraire le parsing dans `parseArgs()` améliore la testabilité de toutes les stories 10.x. Si ce refactoring est fait ici, les stories 10.1 et 10.2 pourront profiter de la même struct `Config`.

### References

- Boucle de parsing actuelle : [Source: cmd/create-go-starter/main.go:215-251]
- `usage()` à mettre à jour : [Source: cmd/create-go-starter/main.go:254-285]
- Flag `-h` déjà implémenté : [Source: cmd/create-go-starter/main.go:218]
- Drapeaux interactif (Story 10.1) : [Source: _bmad-output/implementation-artifacts/10-1-interactive-mode.md]
- Drapeau dry-run (Story 10.2) : [Source: _bmad-output/implementation-artifacts/10-2-dry-run-preview.md]
- Epic 10 description : [Source: _bmad-output/planning-artifacts/epics.md#Epic-10]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

Aucun blocage rencontré. Implémentation directe conforme aux spécifications des Dev Notes.

### Completion Notes List

✅ **Story 10.5 complète** (2026-02-18)

- **Task 1** : Boucle de parsing étendue dans `main.go` (lignes ~237-285). Ajout de `-i`, `-t`/`-t=`, `-d`/`-d=`, `-o`/`-o=`, `-n` comme aliases. Detection des flags inconnus via `else if strings.HasPrefix(arg, "-")` en fin de chaîne.
- **Task 2** : `usage()` mise à jour avec format `-x, --long-option` pour toutes les options.
- **Task 3** : Exemples mis à jour avec formes courtes (`-d sqlite -t minimal`, `-i`, `-n`).
- **Task 4** : 15 nouveaux tests dans `main_test.go` couvrant les ACs #1–#9. Tous PASS. Zéro régression.

**Décision technique** : Les flags inconnus déclenchent `os.Exit(1)` avec message `"unknown flag: -x"` en rouge. Le `else if strings.HasPrefix(arg, "-")` doit être en **dernière position** du else-if chain pour ne pas absorber les flags connus.

**Compatibilité** : Tous les flags longs existants (-database, --database, -template, etc.) continuent de fonctionner sans modification.

### File List

- `cmd/create-go-starter/main.go` (modifié — boucle parsing + usage())
- `cmd/create-go-starter/main_test.go` (modifié — 15 nouveaux tests Story 10.5)

> **Note** : Les fichiers `generator.go`, `docs/usage.md` et `docs/cli-architecture.md` sont aussi modifiés dans le working tree mais relèvent des stories 10.2 (dry-run) et 10.4 (progress bar), pas de la story 10.5.

## Change Log

- 2026-02-18: Story 10.5 — Ajout des aliases courts pour les options CLI (-i, -t, -d, -o, -n). Mise à jour de usage() avec format `-x, --long-option`. Gestion des flags inconnus avec message d'erreur clair. 14 nouveaux tests ajoutés. Rétrocompatibilité totale des flags longs.
- 2026-02-18: Code Review (AI) — 5 corrections appliquées :
  - [H1] Remplacé `flag.Usage()` par `usage()` locale et supprimé import `flag` orphelin (main.go:362)
  - [H2] Corrigé décompte tests dans File List (14→15) et ajouté note sur fichiers multi-stories
  - [M1] Amélioré test `TestShortFlagInteractive` : ajout stdin EOF pipe + vérification anti-usage
  - [M2] Ajouté test `TestUnknownFlagWithValueDoesNotConsume` pour flag inconnu suivi d'une valeur
  - [M3] Restauré exemple `--database=mysql --template=full` dans usage() pour rétrocompatibilité doc
