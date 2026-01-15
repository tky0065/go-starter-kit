# Régénération du site GitHub Pages

Ce guide explique comment régénérer le site GitHub Pages après avoir modifié la documentation dans le dossier `docs/`.

## Prérequis

- Python 3.x installé
- pip installé

## Installation de MkDocs (première fois)

```bash
# Créer un environnement virtuel Python
python3 -m venv venv

# Activer l'environnement virtuel
source venv/bin/activate  # Sur macOS/Linux
# .\venv\Scripts\activate  # Sur Windows

# Installer MkDocs et le thème Material
pip install mkdocs mkdocs-material
```

## Régénération du site

Chaque fois que tu modifies les fichiers dans `docs/`, tu dois régénérer le site HTML:

```bash
# Activer l'environnement virtuel
source venv/bin/activate  # Sur macOS/Linux
# .\venv\Scripts\activate  # Sur Windows

# Régénérer le site (supprime l'ancien et crée le nouveau)
mkdocs build --clean

# Désactiver l'environnement virtuel (optionnel)
deactivate
```

Le site sera généré dans le dossier `site/`.

## Prévisualisation locale

Pour prévisualiser le site localement avant de commit:

```bash
source venv/bin/activate
mkdocs serve
```

Ouvre ensuite http://127.0.0.1:8000/ dans ton navigateur.

## Workflow complet

1. Modifier les fichiers Markdown dans `docs/` (par exemple `docs/usage.md`)
2. Régénérer le site:
   ```bash
   source venv/bin/activate
   mkdocs build --clean
   ```
3. Vérifier les changements:
   ```bash
   mkdocs serve  # Optionnel
   ```
4. Commit et push:
   ```bash
   git add docs/ site/
   git commit -m "docs: Update documentation"
   git push origin develop
   ```

## Structure des fichiers

```
.
├── docs/                          # Fichiers Markdown sources
│   ├── index.md                   # Page d'accueil (copie du README.md)
│   ├── usage.md                   # Guide d'utilisation
│   ├── installation.md            # Guide d'installation
│   ├── generated-project-guide.md # Guide projet généré
│   ├── tutorial-exemple-complet.md # Tutoriel
│   ├── cli-architecture.md        # Architecture du CLI
│   └── contributing.md            # Guide de contribution
│
├── site/                          # Site HTML généré (ne pas modifier manuellement)
│   ├── index.html
│   ├── usage/
│   ├── installation/
│   └── ...
│
├── mkdocs.yml                     # Configuration MkDocs
└── venv/                          # Environnement virtuel Python (ignoré par git)
```

## Configuration MkDocs

Le fichier `mkdocs.yml` contient la configuration du site:
- Thème: Material for MkDocs
- Langue: Français
- Navigation: Structure du menu
- Extensions Markdown: Syntax highlighting, admonitions, etc.

## Notes importantes

- **Ne pas modifier** manuellement les fichiers dans `site/` - ils sont régénérés automatiquement
- Toujours **régénérer le site** après avoir modifié `docs/`
- Le dossier `venv/` est ignoré par git (ne pas le commit)
- Le site utilise MkDocs 1.6.1 et Material theme 9.7.1

## Synchronisation avec README.md

Si tu modifies `README.md` à la racine, n'oublie pas de mettre à jour `docs/index.md`:

```bash
cp README.md docs/index.md
mkdocs build --clean
```

## Dépannage

**Problème**: `mkdocs: command not found`
- **Solution**: Active l'environnement virtuel avec `source venv/bin/activate`

**Problème**: Erreurs d'installation de pip
- **Solution**: Sur macOS avec Python système, utilise `--break-system-packages` ou crée un venv

**Problème**: Le site n'affiche pas les dernières modifications
- **Solution**: Utilise `mkdocs build --clean` pour forcer la régénération complète

## Liens utiles

- [Documentation MkDocs](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
- [Markdown Guide](https://www.markdownguide.org/)
