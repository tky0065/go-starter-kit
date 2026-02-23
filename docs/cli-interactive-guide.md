# Guide du Mode Interactif

[<i class="material-icons">arrow_back</i> Retour à la documentation](index.md)

---

## Introduction

Le mode interactif de `create-go-starter` offre une expérience moderne et intuitive pour créer vos projets Go. Au lieu de mémoriser des drapeaux de ligne de commande, vous naviguez à travers une interface utilisateur terminale (TUI) élégante avec menus, formulaires multi-étapes, et feedback visuel riche.

<i class="material-icons success">check_circle</i> **Interface moderne** avec navigation au clavier
<i class="material-icons success">check_circle</i> **Feedback visuel** en temps réel pendant la génération
<i class="material-icons success">check_circle</i> **Preview avant création** avec le mode dry-run
<i class="material-icons success">check_circle</i> **Aide contextuelle** accessible à tout moment avec `?`

---

## Démarrage Rapide

### Lancer le Mode Interactif

```bash
./create-go-starter --interactive
```

Ou utilisez l'alias court :

```bash
./create-go-starter -i
```

### Prérequis

- **Terminal compatible TTY** : Le mode interactif nécessite un terminal interactif
- **Largeur minimale** : Au moins 60 colonnes recommandées (80+ pour une expérience optimale)
- **Support des couleurs** : Les couleurs sont activées par défaut (désactivables avec `NO_COLOR=1`)

### Désactiver le Mode Interactif

Si vous préférez le mode non-interactif classique, utilisez `--no-interactive` :

```bash
./create-go-starter my-project --no-interactive
```

Ou définissez la variable d'environnement :

```bash
NO_COLOR=1 ./create-go-starter my-project
```

---

## Navigation et Raccourcis Clavier

### Raccourcis Globaux

| Raccourci | Action | Contexte |
|-----------|--------|----------|
| `?` | Afficher l'aide contextuelle | Tous les écrans |
| `Ctrl+C` | Quitter l'application | Tous les écrans |
| `Esc` | Retour / Annuler | Formulaires et écrans secondaires |
| `q` | Quitter (selon contexte) | Écrans d'aide et preview |

### Navigation dans les Menus et Listes

| Raccourci | Action |
|-----------|--------|
| `↑` ou `k` | Déplacer vers le haut |
| `↓` ou `j` | Déplacer vers le bas |
| `Enter` | Sélectionner / Confirmer |
| `Space` | Basculer la sélection (multi-select) |

### Navigation dans les Formulaires

| Raccourci | Action |
|-----------|--------|
| `→` ou `n` | Étape suivante (Next) |
| `←` ou `b` | Étape précédente (Back) |
| `Enter` | Confirmer la saisie et continuer |
| `Tab` | Navigation entre champs (si applicable) |

### Défilement (Preview, Aide)

| Raccourci | Action |
|-----------|--------|
| `↑` / `↓` | Défiler ligne par ligne |
| `PgUp` | Remonter d'une page |
| `PgDn` | Descendre d'une page |
| `Home` | Aller au début |
| `End` | Aller à la fin |

---

## Parcours Complet : Création d'un Projet

### 1. Écran d'Accueil

Lorsque vous lancez le mode interactif, vous arrivez sur l'écran d'accueil avec le logo `go-starter-kit` et un menu principal.

```
╔═══════════════════════════════════╗
║                                   ║
║     🚀  go-starter-kit  🚀       ║
║                                   ║
║   Production-Ready Go API         ║
║   in < 5 minutes                  ║
║                                   ║
╚═══════════════════════════════════╝

  > Create New Project
    Help
    Exit
```

**Actions disponibles :**

- **Create New Project** : Démarre le processus de création
- **Help** : Affiche l'aide contextuelle
- **Exit** : Quitte l'application

**Navigation :** Utilisez `↑` / `↓` pour naviguer, `Enter` pour sélectionner.

---

### 2. Formulaires Multi-Étapes

Une fois "Create New Project" sélectionné, vous entrez dans un processus guidé en plusieurs étapes.

#### Étape 1/4 : Nom du Projet

```
╔═══════════════════════════════════════════════╗
║  [1/4] Nom du Projet                          ║
╚═══════════════════════════════════════════════╝

  ● ○ ○ ○

  Entrez le nom de votre projet Go :

  ┌─────────────────────────────────────────────┐
  │ my-awesome-api                              │
  └─────────────────────────────────────────────┘

  Règles :
  - Lettres minuscules, chiffres, tirets
  - Commence par une lettre
  - Entre 3 et 50 caractères

  [→/n: Suivant] [Esc: Annuler] [?: Aide]
```

**Navigation :** `Enter` ou `→` pour continuer, `Esc` pour annuler.

#### Étape 2/4 : Sélection du Template

```
╔═══════════════════════════════════════════════╗
║  [2/4] Template                               ║
╚═══════════════════════════════════════════════╝

  ○ ● ○ ○

  Choisissez un template de base :

  > Standard (REST API + Auth)
    Microservice (Clean Architecture)
    GraphQL (avec gqlgen)
    gRPC (Protocol Buffers)

  [↑↓: Naviguer] [Enter: Sélectionner]
  [←/b: Retour] [→/n: Suivant]
```

**Navigation :** `↑` / `↓` pour choisir, `Enter` ou `→` pour continuer, `←` ou `b` pour revenir en arrière.

#### Étape 3/4 : Base de Données

```
╔═══════════════════════════════════════════════╗
║  [3/4] Base de Données                        ║
╚═══════════════════════════════════════════════╝

  ○ ○ ● ○

  Sélectionnez votre base de données :

  > PostgreSQL (Recommandé)
    MySQL
    SQLite
    MongoDB
    Aucune (Stateless API)

  [↑↓: Naviguer] [Enter: Sélectionner]
  [←/b: Retour] [→/n: Suivant]
```

#### Étape 4/4 : Niveau d'Observabilité

```
╔═══════════════════════════════════════════════╗
║  [4/4] Observabilité                          ║
╚═══════════════════════════════════════════════╝

  ○ ○ ○ ●

  Niveau de monitoring et tracing :

  > Basique (Logs + Métriques)
    Standard (+ Prometheus)
    Avancé (+ OpenTelemetry + Jaeger)
    Aucun

  [↑↓: Naviguer] [Enter: Sélectionner]
  [←/b: Retour] [→/n: Confirmer]
```

---

### 3. Écran de Résumé

Après avoir complété les 4 étapes, un écran de résumé affiche votre configuration :

```
╔═══════════════════════════════════════════════╗
║  Résumé de la Configuration                   ║
╚═══════════════════════════════════════════════╝

  ✓ Nom du projet    : my-awesome-api
  ✓ Template         : Standard (REST API + Auth)
  ✓ Base de données  : PostgreSQL
  ✓ Observabilité    : Basique (Logs + Métriques)

  Fichiers à créer   : ~34 fichiers
  Taille estimée     : ~850 KB

  Continuer ?

  [ Oui ]  Non

  [←→: Naviguer] [Enter: Confirmer] [Esc: Annuler]
```

**Actions :**

- `←` / `→` pour basculer entre Oui/Non
- `Enter` pour confirmer
- `Esc` ou `b` pour revenir à l'étape précédente

---

### 4. Génération du Projet

Une fois confirmé, le processus de génération démarre avec une barre de progression animée :

```
╔═══════════════════════════════════════════════════════════╗
║  Génération: my-awesome-api                               ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  [████████████████████░░░░░░░░░░] 67% (24/36 fichiers)   ║
║                                                           ║
║  📊 Statistiques :                                        ║
║    ✓ Fichiers créés     : 24                             ║
║    📦 Taille totale      : ~620 KB                        ║
║    ⏱  Temps écoulé       : 2.3s                          ║
║    ⏳ ETA                : ~1.2s                          ║
║                                                           ║
║  🔄 Création de internal/domain/user/service.go...       ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
```

**Informations affichées :**

- Barre de progression visuelle avec pourcentage
- Nombre de fichiers créés / total
- Taille totale du projet
- Temps écoulé et estimation du temps restant (ETA)
- Fichier actuellement en cours de création

---

### 5. Écran de Succès

À la fin de la génération, un écran de succès affiche les prochaines étapes :

```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║     ✓  Projet créé avec succès !                         ║
║                                                           ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  Votre projet my-awesome-api est prêt.                   ║
║                                                           ║
║  Prochaines étapes :                                      ║
║                                                           ║
║    1. cd my-awesome-api                                  ║
║    2. ./setup.sh                                         ║
║    3. make run                                           ║
║                                                           ║
║  Documentation :                                          ║
║    - README.md                                           ║
║    - docs/GETTING_STARTED.md                             ║
║                                                           ║
║  [Appuyez sur Enter pour quitter]                       ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
```

---

## Mode Dry-Run : Preview Avant Création

Le mode dry-run permet de visualiser exactement ce qui sera créé **avant** de générer les fichiers.

### Lancer le Dry-Run

```bash
./create-go-starter --interactive --dry-run
```

Ou avec le mode non-interactif :

```bash
./create-go-starter my-project --dry-run
```

### Écran de Preview

Après avoir configuré votre projet dans les formulaires, au lieu de générer les fichiers, l'écran de preview s'affiche :

```
╔═══════════════════════════════════════════════════════════╗
║  Preview : my-awesome-api (34 fichiers, ~850 KB)         ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  📁 Structure du projet :                                 ║
║                                                           ║
║  + cmd/                                                   ║
║    + main.go                                    (120 B)  ║
║  + internal/                                              ║
║    + adapters/                                            ║
║      + handlers/                                          ║
║        + auth_handler.go                       (3.2 KB)  ║
║        + user_handler.go                       (2.8 KB)  ║
║      + repository/                                        ║
║        + user_repository.go                    (2.1 KB)  ║
║    + domain/                                              ║
║      + user/                                              ║
║        + service.go                            (4.5 KB)  ║
║    + models/                                              ║
║      + user.go                                 (1.8 KB)  ║
║  + pkg/                                                   ║
║    + config/                                              ║
║      + config.go                               (2.5 KB)  ║
║  ...                                                      ║
║                                                           ║
║  ↓ Plus... (↑↓ pour défiler, Enter pour confirmer)       ║
╚═══════════════════════════════════════════════════════════╝
```

### Coloration Syntaxique

Pour les fichiers importants (comme `main.go`, `config.go`, etc.), vous pouvez voir un aperçu du contenu avec coloration syntaxique basique :

```go
package main  // (bleu)

import (      // (bleu)
    "log"     // (vert)
)

func main() {  // (bleu)
    // Start server  // (gris)
    log.Println("Server started")  // (vert)
}
```

**Navigation dans la preview :**

- `↑` / `↓` : Défiler ligne par ligne
- `PgUp` / `PgDn` : Défiler par page
- `Home` / `End` : Aller au début / fin
- `Enter` : Confirmer et créer le projet
- `Esc` ou `q` : Annuler

---

## Écran d'Aide Contextuel

À tout moment, appuyez sur `?` pour afficher l'aide contextuelle adaptée à l'écran actuel.

### Exemple : Aide sur l'Écran d'Accueil

```
╔═══════════════════════════════════════════════════════════╗
║  Raccourcis Clavier                                       ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  Global :                                                 ║
║    ?        Afficher cette aide                          ║
║    Ctrl+C   Quitter l'application                        ║
║    Esc      Retour / Annuler                             ║
║                                                           ║
║  Navigation :                                             ║
║    ↑ ↓      Déplacer dans les listes                     ║
║    ← →      Naviguer dans les formulaires               ║
║    b / n    Retour / Suivant (alternatif)               ║
║    Enter    Sélectionner / Confirmer                     ║
║    Space    Basculer la sélection (multi-select)        ║
║                                                           ║
║  Défilement (Preview, Aide) :                            ║
║    PgUp     Remonter d'une page                          ║
║    PgDn     Descendre d'une page                         ║
║    Home     Aller en haut                                ║
║    End      Aller en bas                                 ║
║                                                           ║
║  [Appuyez sur une touche pour fermer]                   ║
╚═══════════════════════════════════════════════════════════╝
```

---

## Thème Visuel et Couleurs

Le mode interactif utilise une palette de couleurs cohérente avec le branding **go-starter-kit** :

| Élément | Couleur | Hex | Usage |
|---------|---------|-----|-------|
| <i class="material-icons success">circle</i> **Primary** | Vert | `#00c853` | Succès, confirmations, éléments actifs |
| <i class="material-icons info">circle</i> **Secondary** | Bleu | `#00b0ff` | Titres, informations, highlights |
| <i class="material-icons warning">circle</i> **Warning** | Orange | `#ff6d00` | Avertissements, attention |
| <i class="material-icons error">circle</i> **Error** | Rouge | `#ff1744` | Erreurs, actions destructives |
| **Text** | Blanc | `#ffffff` | Texte principal |
| **Muted** | Gris | `#666666` | Texte secondaire, labels |
| **Border** | Gris | `#888888` | Bordures, séparateurs |

### Icônes et Symboles

| Symbole | Signification |
|---------|---------------|
| <i class="material-icons success">check</i> `✓` | Succès / Validé |
| <i class="material-icons error">close</i> `✗` | Erreur / Échec |
| <i class="material-icons warning">warning</i> `⚠` | Avertissement |
| <i class="material-icons info">info</i> `ℹ` | Information |
| `🔄` | En cours de traitement |
| `📊` | Statistiques |
| `📦` | Taille / Package |
| `⏱` | Temps écoulé |
| `⏳` | Estimation |

---

## Animations et Feedback Visuel

### Logo Fade-In

Au démarrage du mode interactif, le logo `go-starter-kit` apparaît progressivement avec une animation de fade-in sur 1 seconde.

### Progress Bar Animée

Pendant la génération, la barre de progression utilise un gradient animé qui "pulse" pour indiquer l'activité :

```
[████████████████████░░░░░░░░░░] 67%
```

Le gradient passe de vert (`#00c853`) à bleu (`#00b0ff`).

### Spinner Élégant

Pour les opérations longues (comme l'installation des dépendances), un spinner tourne avec un texte contextuel :

```
🔄 Initialisation du dépôt Git...
```

---

## Gestion des Erreurs

### Erreur avec Suggestions

Si une erreur survient, elle est affichée dans une box rouge avec des suggestions de résolution :

```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║  ✗  Erreur : Nom de projet invalide                      ║
║                                                           ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  Le nom "MyProject" contient des caractères invalides.   ║
║                                                           ║
║  Suggestions :                                            ║
║    • Utilisez uniquement des minuscules et tirets        ║
║    • Exemple valide : my-project                         ║
║    • Longueur : 3-50 caractères                          ║
║                                                           ║
║  [Appuyez sur Enter pour recommencer]                   ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
```

### Erreurs Courantes

| Erreur | Cause | Solution |
|--------|-------|----------|
| Nom invalide | Majuscules, caractères spéciaux | Utilisez minuscules, chiffres, tirets |
| Projet existe déjà | Dossier déjà présent | Choisissez un autre nom ou supprimez le dossier |
| Pas de TTY | Mode interactif sans terminal interactif | Utilisez `--no-interactive` |
| Terminal trop petit | Largeur < 60 colonnes | Agrandissez votre terminal |

---

## Troubleshooting

### Le Mode Interactif Ne Se Lance Pas

**Symptôme :** L'application démarre en mode non-interactif même avec `--interactive`.

**Causes possibles :**

1. **Pas de TTY détecté**
   ```bash
   # Vérifier si stdin est un TTY
   [ -t 0 ] && echo "TTY" || echo "Pas de TTY"
   ```

2. **Variable NO_COLOR définie**
   ```bash
   # Vérifier
   echo $NO_COLOR
   # Solution : désinifier
   unset NO_COLOR
   ```

3. **Redirection de stdin**
   ```bash
   # ❌ Ceci désactive le TTY
   echo "" | ./create-go-starter --interactive

   # ✅ Lancez directement
   ./create-go-starter --interactive
   ```

---

### Pas de Couleurs Affichées

**Symptôme :** L'interface est monochrome (noir et blanc).

**Solution :**

```bash
# Vérifier NO_COLOR
echo $NO_COLOR

# Si défini, le désactiver
unset NO_COLOR

# Vérifier le support des couleurs de votre terminal
echo $TERM
# Devrait afficher quelque chose comme "xterm-256color"
```

**Terminals supportés :**

- ✅ iTerm2 (macOS)
- ✅ Terminal.app (macOS)
- ✅ GNOME Terminal (Linux)
- ✅ Konsole (Linux)
- ✅ Windows Terminal (Windows 10+)
- ✅ Alacritty
- ✅ Kitty
- ⚠️ CMD.exe (Windows) - Support limité
- ❌ Anciens terminaux sans support ANSI

---

### Terminal Trop Petit

**Symptôme :** L'interface est coupée ou illisible.

**Solution :** Agrandissez votre terminal à au moins 80 colonnes × 24 lignes.

**Vérifier la taille actuelle :**

```bash
# Afficher colonnes × lignes
echo $COLUMNS × $LINES
```

**Largeurs recommandées :**

- **Minimum** : 60 colonnes
- **Recommandé** : 80 colonnes
- **Optimal** : 100+ colonnes

---

### Navigation Ne Fonctionne Pas

**Symptôme :** Les touches fléchées n'ont aucun effet.

**Causes possibles :**

1. **Mode Vi/Emacs du shell interférant**
   - Solution : Le mode interactif utilise son propre système de navigation

2. **Terminal non standard**
   - Vérifiez que votre terminal envoie les séquences ANSI standard

3. **Redirection de stdin**
   - Assurez-vous de lancer directement sans pipe ni redirection

---

### Animations Saccadées

**Symptôme :** Les animations (logo fade-in, progress bar) sont lentes ou saccadées.

**Causes :**

- Machine lente ou surchargée
- Terminal distant via SSH avec latence réseau

**Solutions :**

- Les animations sont légères (30-60 FPS max), elles devraient fonctionner sur la plupart des machines
- Pour SSH : Utilisez `--no-interactive` si la latence est trop élevée

---

## Comparaison : Mode Interactif vs Non-Interactif

| Aspect | Mode Interactif (`-i`) | Mode Non-Interactif |
|--------|------------------------|---------------------|
| **Lancement** | `./create-go-starter -i` | `./create-go-starter my-project` |
| **Configuration** | Guidé via formulaires | Drapeaux CLI (`--database`, `--template`) |
| **Feedback** | Visuel riche (progress bar, spinners) | Texte simple |
| **Preview** | Écran de preview intégré avec `--dry-run` | Liste textuelle |
| **Aide** | Aide contextuelle avec `?` | `--help` pour aide globale |
| **Navigation** | Clavier (↑↓←→) | Ligne de commande |
| **Utilisation** | Débutants, exploration | Scripts, CI/CD, experts |
| **TTY requis** | Oui | Non |
| **Couleurs** | Oui (désactivables) | Limitées |

---

## Exemples d'Utilisation

### Création Standard avec Mode Interactif

```bash
# Lancer le mode interactif
./create-go-starter --interactive

# Suivre le processus guidé :
# 1. Écran d'accueil → Sélectionner "Create New Project"
# 2. Nom du projet → "my-api"
# 3. Template → "Standard (REST API + Auth)"
# 4. Base de données → "PostgreSQL"
# 5. Observabilité → "Basique"
# 6. Résumé → Confirmer
# 7. Génération → Patienter
# 8. Succès → Suivre les instructions
```

---

### Preview Sans Créer (Dry-Run)

```bash
# Mode interactif avec preview
./create-go-starter -i --dry-run

# Configurer le projet normalement, puis :
# - Visualiser l'arborescence complète
# - Voir les aperçus de code
# - Annuler ou confirmer la création
```

---

### Création Rapide pour Experts

Si vous connaissez déjà les options, utilisez le mode non-interactif :

```bash
./create-go-starter my-api \
  --template=standard \
  --database=postgresql \
  --observability=basic
```

---

## Bonnes Pratiques

### <i class="material-icons success">check</i> Recommandations

1. **Utilisez le mode interactif si :**
   - Vous créez un projet pour la première fois
   - Vous explorez les options disponibles
   - Vous voulez une confirmation visuelle avant création
   - Vous travaillez localement sur votre machine

2. **Utilisez le mode non-interactif si :**
   - Vous automatisez la création (scripts, CI/CD)
   - Vous créez plusieurs projets similaires
   - Vous travaillez sur un serveur distant
   - Vous êtes expert et connaissez les options

3. **Dry-Run avant production :**
   - Utilisez `--dry-run` pour vérifier la structure avant de créer
   - Inspectez les fichiers générés pour comprendre l'architecture
   - Validez que le template correspond à vos besoins

4. **Apprenez les raccourcis :**
   - `?` pour l'aide contextuelle
   - `b` / `n` pour naviguer rapidement
   - `Esc` pour annuler proprement
   - `Ctrl+C` pour quitter immédiatement

---

## FAQ

### Puis-je revenir en arrière dans les formulaires ?

**Oui !** Utilisez `←` ou `b` pour revenir à l'étape précédente. Vos valeurs saisies sont conservées.

---

### Les valeurs sont-elles sauvegardées si je quitte ?

**Non.** Si vous quittez (`Ctrl+C` ou `Esc`), la configuration est perdue. Recommencez depuis le début.

---

### Puis-je utiliser le mode interactif dans un script ?

**Non recommandé.** Le mode interactif nécessite un TTY. Pour les scripts, utilisez le mode non-interactif avec drapeaux.

---

### Le mode interactif fonctionne-t-il sur Windows ?

**Oui**, avec Windows Terminal ou un terminal moderne. CMD.exe a un support limité. Préférez PowerShell ou Windows Terminal.

---

### Comment désactiver les couleurs ?

```bash
NO_COLOR=1 ./create-go-starter my-project
```

Cela désactive également le mode interactif.

---

### Puis-je personnaliser les couleurs du thème ?

Actuellement, le thème est fixe (branding go-starter-kit). Une personnalisation future pourrait être ajoutée.

---

## Ressources Complémentaires

### Documentation Connexe

- <i class="material-icons">menu_book</i> [Guide d'utilisation général](usage.md)
- <i class="material-icons">menu_book</i> [Architecture CLI](cli-architecture.md)
- <i class="material-icons">menu_book</i> [Guide du projet généré](generated-project-guide.md)
- <i class="material-icons">menu_book</i> [Tutoriel complet](tutorial-exemple-complet.md)

### Liens Externes

- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea) - Framework TUI utilisé
- [Bubbles Components](https://github.com/charmbracelet/bubbles) - Composants UI
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss) - Système de style

---

## Support et Contributions

### Signaler un Bug

Si vous rencontrez un problème avec le mode interactif :

1. Vérifiez d'abord cette documentation
2. Consultez les [issues GitHub](https://github.com/tky0065/go-starter-kit/issues)
3. Ouvrez une nouvelle issue avec :
   - Version de `create-go-starter`
   - Système d'exploitation et terminal utilisé
   - Commande exacte exécutée
   - Comportement attendu vs observé
   - Captures d'écran si applicable

### Contribuer

Les contributions sont les bienvenues ! Consultez [CONTRIBUTING.md](contributing.md) pour les guidelines.

---

[<i class="material-icons">arrow_back</i> Retour à la documentation](index.md)
