# Regenerating the GitHub Pages Site

This guide explains how to regenerate the GitHub Pages site after modifying the documentation in the `docs/` folder.

## Prerequisites

- Python 3.x installed
- pip installed

## Installing MkDocs (first time)

```bash
# Create a Python virtual environment
python3 -m venv venv

# Activate the virtual environment
source venv/bin/activate  # On macOS/Linux
# .\venv\Scripts\activate  # On Windows

# Install MkDocs and the Material theme
pip install mkdocs mkdocs-material
```

## Regenerating the Site

Every time you modify files in `docs/`, you need to regenerate the HTML site:

```bash
# Activate the virtual environment
source venv/bin/activate  # On macOS/Linux
# .\venv\Scripts\activate  # On Windows

# Regenerate the site (deletes the old one and creates a new one)
mkdocs build --clean

# Deactivate the virtual environment (optional)
deactivate
```

The site will be generated in the `site/` folder.

## Local Preview

To preview the site locally before committing:

```bash
source venv/bin/activate
mkdocs serve
```

Then open http://127.0.0.1:8000/ in your browser.

## Complete Workflow

1. Edit the Markdown files in `docs/` (for example `docs/usage.md`)
2. Regenerate the site:
   ```bash
   source venv/bin/activate
   mkdocs build --clean
   ```
3. Verify the changes:
   ```bash
   mkdocs serve  # Optional
   ```
4. Commit and push:
   ```bash
   git add docs/ site/
   git commit -m "docs: Update documentation"
   git push origin develop
   ```

## File Structure

```
.
├── docs/                          # Source Markdown files
│   ├── index.md                   # Home page (copy of README.md)
│   ├── usage.md                   # Usage guide
│   ├── installation.md            # Installation guide
│   ├── generated-project-guide.md # Generated project guide
│   ├── tutorial-exemple-complet.md # Tutorial
│   ├── cli-architecture.md        # CLI architecture
│   └── contributing.md            # Contributing guide
│
├── site/                          # Generated HTML site (do not modify manually)
│   ├── index.html
│   ├── usage/
│   ├── installation/
│   └── ...
│
├── mkdocs.yml                     # MkDocs configuration
└── venv/                          # Python virtual environment (ignored by git)
```

## MkDocs Configuration

The `mkdocs.yml` file contains the site configuration:
- Theme: Material for MkDocs
- Language: French
- Navigation: Menu structure
- Markdown extensions: Syntax highlighting, admonitions, etc.

## Important Notes

- **Do not manually modify** files in `site/` - they are automatically regenerated
- Always **regenerate the site** after modifying `docs/`
- The `venv/` folder is ignored by git (do not commit it)
- The site uses MkDocs 1.6.1 and Material theme 9.7.1

## Syncing with README.md

If you modify `README.md` at the root, don't forget to update `docs/index.md`:

```bash
cp README.md docs/index.md
mkdocs build --clean
```

## Troubleshooting

**Problem**: `mkdocs: command not found`
- **Solution**: Activate the virtual environment with `source venv/bin/activate`

**Problem**: pip installation errors
- **Solution**: On macOS with system Python, use `--break-system-packages` or create a venv

**Problem**: The site doesn't show the latest changes
- **Solution**: Use `mkdocs build --clean` to force a full regeneration

## Useful Links

- [MkDocs Documentation](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
- [Markdown Guide](https://www.markdownguide.org/)
