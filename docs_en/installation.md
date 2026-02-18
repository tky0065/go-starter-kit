# Installation Guide

This guide details all the methods to install `create-go-starter` on your system.

## Prerequisites

### System Requirements

- **Go 1.25 or higher** - Check your version with `go version`
- **Git** - To clone the repository
- **Disk space** - ~50MB for the tool and its dependencies

### Verify Go Installation

```bash
go version
# Devrait afficher: go version go1.25.x ...
```

If Go is not installed, download it from [golang.org/dl](https://golang.org/dl/).

### Configure GOPATH and PATH

Make sure that `$GOPATH/bin` is in your `PATH`:

```bash
# Vérifier GOPATH
go env GOPATH
# Généralement: /Users/<username>/go sur macOS
#              $HOME/go sur Linux
#              C:\Users\<username>\go sur Windows

# Ajouter à PATH si nécessaire (dans ~/.zshrc, ~/.bashrc, ou ~/.bash_profile)
export PATH=$PATH:$(go env GOPATH)/bin
```

After modification, reload your shell:

```bash
source ~/.zshrc  # ou ~/.bashrc selon votre shell
```

### Optional Tools (Recommended)

- **golangci-lint** - For development and contribution
  ```bash
  # macOS avec Homebrew
  brew install golangci-lint

  # Linux
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

  # Ou avec go install
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- **Docker** - To test generated projects
  ```bash
  # macOS avec Homebrew
  brew install --cask docker

  # Linux: Suivez les instructions officielles
  # https://docs.docker.com/engine/install/
  ```

## Method 1: Direct Installation (Recommended)

This is the simplest and recommended method. Global installation in a single command, **without cloning the repository**.

### One-command Installation

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

This command:
- Automatically downloads the code from GitHub
- Compiles the binary
- Installs it in `$GOPATH/bin` (usually `~/go/bin`)
- Makes it available globally

### Verification

```bash
create-go-starter --help
```

You should see the tool's help output displayed.

**Note**: Make sure that `$GOPATH/bin` is in your `PATH`. Otherwise:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Binary Location

The binary is installed in:
- **macOS/Linux**: `~/go/bin/create-go-starter`
- **Windows**: `C:\Users\<username>\go\bin\create-go-starter.exe`

### Advantages of This Method

- **Ultra simple** - A single command
- **No need to clone** - Go handles everything automatically
- **Always up to date** - Use `@latest` for the latest version
- Binary available globally
- Standard method in the Go ecosystem

## Method 2: Installation from Source

This method is recommended for contributors or if you want to customize the tool.

### Steps

1. **Clone the repository**

   ```bash
   git clone https://github.com/tky0065/go-starter-kit.git
   cd go-starter-kit
   ```

2. **Build with go build**

   ```bash
   go build -o create-go-starter ./cmd/create-go-starter
   ```

   The `create-go-starter` binary will be created in the current directory.

3. **Or build with Makefile**

   ```bash
   make build
   ```

   The binary will be created in the current directory.

4. **Manual installation (optional)**

   To make the binary available globally:

   ```bash
   # Option A: Copier vers $GOPATH/bin
   cp create-go-starter $(go env GOPATH)/bin/

   # Option B: Copier vers /usr/local/bin (nécessite sudo sur macOS/Linux)
   sudo cp create-go-starter /usr/local/bin/

   # Option C: Ajouter le répertoire actuel à PATH
   export PATH=$PATH:$(pwd)  # Ajouter ceci dans ~/.zshrc pour le rendre permanent
   ```

5. **Verify the installation**

   ```bash
   # Si installé globalement
   create-go-starter --help

   # Ou utiliser le chemin relatif
   ./create-go-starter --help
   ```

### Advantages of This Method

- Full control over the build
- Easy to modify the source code
- Ideal for development and testing
- Allows creating custom builds

### Build with Advanced Options

```bash
# Build avec optimisations pour production
go build -ldflags="-s -w" -o create-go-starter ./cmd/create-go-starter

# Build pour un OS/architecture spécifique
GOOS=linux GOARCH=amd64 go build -o create-go-starter-linux ./cmd/create-go-starter
GOOS=windows GOARCH=amd64 go build -o create-go-starter.exe ./cmd/create-go-starter
GOOS=darwin GOARCH=arm64 go build -o create-go-starter-macos-arm ./cmd/create-go-starter
```

## Method 3: Pre-compiled Binary (Coming Soon)

> **Note**: This method will be available once the project publishes releases with pre-compiled binaries.

When available, you will be able to download binaries from the [Releases](https://github.com/tky0065/go-starter-kit/releases) page.

### Installation on macOS/Linux

```bash
# Télécharger le binaire (remplacez VERSION par la version souhaitée)
curl -L https://github.com/tky0065/go-starter-kit/releases/download/vVERSION/create-go-starter-macos -o create-go-starter

# Donner les permissions d'exécution
chmod +x create-go-starter

# Déplacer vers un répertoire dans PATH
sudo mv create-go-starter /usr/local/bin/
```

### Installation on Windows

1. Download `create-go-starter.exe` from the Releases page
2. Place the file in a directory of your choice
3. Add that directory to your system PATH variable

## Verifying the Installation

After installation, verify that the tool works correctly:

```bash
# Afficher l'aide
create-go-starter --help

# Devrait afficher quelque chose comme:
# Usage: create-go-starter <project-name>
#   -h, --help    Show this help message
```

### Project Creation Test

Test by creating a simple project:

```bash
# Créer un projet test
create-go-starter test-project

# Vérifier que le projet a été créé
ls -la test-project/

# Nettoyer
rm -rf test-project
```

If you see the project structure created successfully (in green in the terminal), the installation is working perfectly!

## Updating

### Update via go install (Method 1)

If you installed with Method 1 (direct installation), simply re-run the command:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Go will automatically download and install the latest available version. This is the simplest method!

### Update from Source (Method 2)

If you installed from source:

```bash
# Se placer dans le repository
cd /path/to/go-starter-kit

# Récupérer les dernières modifications
git pull origin main

# Rebuild
make build  # ou go build -o create-go-starter ./cmd/create-go-starter

# Réinstaller si nécessaire
cp create-go-starter $(go env GOPATH)/bin/
```

### Check the Version

```bash
create-go-starter doctor
```

The `doctor` command displays the CLI version and checks your environment (Go, Git, Docker).

To check the Git commit of the repository:

```bash
cd /path/to/go-starter-kit
git log -1 --oneline
```

## Uninstallation

### If installed via go install or manual copy

```bash
# Supprimer le binaire
rm $(go env GOPATH)/bin/create-go-starter

# Ou si installé dans /usr/local/bin
sudo rm /usr/local/bin/create-go-starter
```

### If installed from source (in the repository)

```bash
# Supprimer le binaire local
cd /path/to/go-starter-kit
rm create-go-starter

# Optionnel: Supprimer le repository complet
cd ..
rm -rf go-starter-kit
```

### Clean Go Cache

To free up space:

```bash
# Nettoyer le cache de modules
go clean -modcache

# Nettoyer le cache de build
go clean -cache
```

## Troubleshooting

### Problem: "command not found: create-go-starter"

**Possible causes**:
1. `$GOPATH/bin` is not in your `PATH`
2. The binary was not installed correctly
3. **The shell cache has not been reloaded** (most common cause)

**Solutions**:

**Solution 1: Reload the shell cache (<i class="material-icons info">bolt</i> Quick - Try this first!)**

After `go install`, your shell (zsh/bash) may have a cached version of the available commands list. Reload it:

```bash
# Pour zsh et bash
hash -r

# Puis vérifier
which create-go-starter
create-go-starter --help
```

**Solution 2: Restart the terminal**

Close and reopen your terminal. This is often the simplest solution!

**Solution 3: Verify and configure the PATH**

If the previous solutions don't work:

```bash
# Vérifier GOPATH
go env GOPATH

# Vérifier si le binaire existe
ls -l $(go env GOPATH)/bin/create-go-starter

# Vérifier si GOPATH/bin est dans PATH
echo $PATH | grep "$(go env GOPATH)/bin"

# Si absent, ajouter GOPATH/bin au PATH (dans ~/.zshrc ou ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

# Recharger le shell
source ~/.zshrc  # ou ~/.bashrc
```

**Solution 4: Use the full path temporarily**

While waiting to resolve the PATH:

```bash
$(go env GOPATH)/bin/create-go-starter mon-projet
```

### Problem: "permission denied" on macOS

**Cause**: macOS Gatekeeper blocks unsigned binaries.

**Solution**:

```bash
# Donner les permissions d'exécution
chmod +x create-go-starter

# Autoriser l'exécution (macOS)
xattr -d com.apple.quarantine create-go-starter
```

Or allow the application in:
`System Preferences > Security & Privacy > General > Allow Anyway`

### Problem: Compilation errors

**Cause**: Go version too old or missing dependencies.

**Solutions**:

```bash
# Vérifier la version de Go (doit être >= 1.25)
go version

# Mettre à jour les dépendances
go mod tidy
go mod tidy

# Nettoyer et rebuild
go clean
go build -o create-go-starter ./cmd/create-go-starter
```

### Problem: "go: cannot find main module"

**Cause**: You are not in the correct directory.

**Solution**:

```bash
# Assurez-vous d'être dans le repository go-starter-kit
cd /path/to/go-starter-kit

# Vérifier que go.mod existe
ls go.mod

# Puis réessayer l'installation
go install ./cmd/create-go-starter
```

### Problem: Go version conflicts

**Cause**: Multiple Go versions installed.

**Solutions**:

```bash
# Vérifier quelle version de Go est utilisée
which go
go version

# Sur macOS avec Homebrew
brew list go
brew upgrade go

# Définir la version de Go à utiliser (avec go.mod)
cat go.mod | grep "^go "
```

### Problem: Slow build

**Solutions**:

```bash
# Activer le cache de build (devrait être activé par défaut)
go env GOCACHE

# Build avec cache
go build -o create-go-starter ./cmd/create-go-starter

# Si toujours lent, nettoyer puis rebuild
go clean -cache
go build -o create-go-starter ./cmd/create-go-starter
```

## Installation for Development

If you want to contribute to the project, follow these additional steps:

```bash
# 1. Fork le repository sur GitHub
# 2. Cloner votre fork
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit

# 3. Ajouter le remote upstream
git remote add upstream https://github.com/tky0065/go-starter-kit.git

# 4. Installer les dépendances de développement
make install-dev  # ou installer golangci-lint manuellement

# 5. Installer l'outil en mode dev
go install ./cmd/create-go-starter

# 6. Créer une branche pour vos changements
git checkout -b feature/ma-fonctionnalite

# 7. Développer, tester, commiter
make test
make lint
git commit -m "feat: description de la fonctionnalité"

# 8. Pousser et créer une PR
git push origin feature/ma-fonctionnalite
```

## Next Steps

Now that `create-go-starter` is installed, check out:

- **[Usage Guide](./usage.md)** - Learn how to use the tool
- **[Generated Projects Guide](./generated-project-guide.md)** - Develop with the created projects
- **[Contributing Guide](./contributing.md)** - Contribute to the project

Or get started right away:

```bash
create-go-starter mon-premier-projet
cd mon-premier-projet
# Suivez les instructions affichées!
```

Happy coding! <i class="material-icons success">rocket_launch</i>
