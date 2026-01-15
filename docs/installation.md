# Guide d'installation

Ce guide détaille toutes les méthodes pour installer `create-go-starter` sur votre système.

## Prérequis

### Système requis

- **Go 1.25 ou supérieur** - Vérifiez votre version avec `go version`
- **Git** - Pour cloner le repository
- **Espace disque** - ~50MB pour l'outil et ses dépendances

### Vérifier l'installation de Go

```bash
go version
# Devrait afficher: go version go1.25.x ...
```

Si Go n'est pas installé, téléchargez-le depuis [golang.org/dl](https://golang.org/dl/).

### Configurer GOPATH et PATH

Assurez-vous que `$GOPATH/bin` est dans votre `PATH`:

```bash
# Vérifier GOPATH
go env GOPATH
# Généralement: /Users/<username>/go sur macOS
#              $HOME/go sur Linux
#              C:\Users\<username>\go sur Windows

# Ajouter à PATH si nécessaire (dans ~/.zshrc, ~/.bashrc, ou ~/.bash_profile)
export PATH=$PATH:$(go env GOPATH)/bin
```

Après modification, rechargez votre shell:

```bash
source ~/.zshrc  # ou ~/.bashrc selon votre shell
```

### Outils optionnels (recommandés)

- **golangci-lint** - Pour le développement et contribution
  ```bash
  # macOS avec Homebrew
  brew install golangci-lint

  # Linux
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

  # Ou avec go install
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- **Docker** - Pour tester les projets générés
  ```bash
  # macOS avec Homebrew
  brew install --cask docker

  # Linux: Suivez les instructions officielles
  # https://docs.docker.com/engine/install/
  ```

## Méthode 1: Installation directe (Recommandée)

C'est la méthode la plus simple et recommandée. Installation globale en une seule commande, **sans cloner le repository**.

### Installation en une commande

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Cette commande:
- Télécharge automatiquement le code depuis GitHub
- Compile le binaire
- L'installe dans `$GOPATH/bin` (généralement `~/go/bin`)
- Le rend disponible globalement

### Vérification

```bash
create-go-starter --help
```

Vous devriez voir l'aide de l'outil s'afficher.

**Note**: Assurez-vous que `$GOPATH/bin` est dans votre `PATH`. Sinon:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Localisation du binaire

Le binaire est installé dans:
- **macOS/Linux**: `~/go/bin/create-go-starter`
- **Windows**: `C:\Users\<username>\go\bin\create-go-starter.exe`

### Avantages de cette méthode

- **Ultra simple** - Une seule commande
- **Pas besoin de cloner** - Go gère tout automatiquement
- **Toujours à jour** - Utilisez `@latest` pour la dernière version
- Binaire disponible globalement
- Méthode standard de l'écosystème Go

## Méthode 2: Installation depuis les sources

Cette méthode est recommandée pour les contributeurs ou si vous voulez personnaliser l'outil.

### Étapes

1. **Cloner le repository**

   ```bash
   git clone https://github.com/tky0065/go-starter-kit.git
   cd go-starter-kit
   ```

2. **Build avec go build**

   ```bash
   go build -o create-go-starter ./cmd/create-go-starter
   ```

   Le binaire `create-go-starter` sera créé dans le répertoire courant.

3. **Ou build avec Makefile**

   ```bash
   make build
   ```

   Le binaire sera créé dans le répertoire courant.

4. **Installation manuelle (optionnel)**

   Pour rendre le binaire disponible globalement:

   ```bash
   # Option A: Copier vers $GOPATH/bin
   cp create-go-starter $(go env GOPATH)/bin/

   # Option B: Copier vers /usr/local/bin (nécessite sudo sur macOS/Linux)
   sudo cp create-go-starter /usr/local/bin/

   # Option C: Ajouter le répertoire actuel à PATH
   export PATH=$PATH:$(pwd)  # Ajouter ceci dans ~/.zshrc pour le rendre permanent
   ```

5. **Vérifier l'installation**

   ```bash
   # Si installé globalement
   create-go-starter --help

   # Ou utiliser le chemin relatif
   ./create-go-starter --help
   ```

### Avantages de cette méthode

- Contrôle total sur le build
- Facile de modifier le code source
- Idéal pour le développement et les tests
- Permet de créer des builds personnalisés

### Build avec options avancées

```bash
# Build avec optimisations pour production
go build -ldflags="-s -w" -o create-go-starter ./cmd/create-go-starter

# Build pour un OS/architecture spécifique
GOOS=linux GOARCH=amd64 go build -o create-go-starter-linux ./cmd/create-go-starter
GOOS=windows GOARCH=amd64 go build -o create-go-starter.exe ./cmd/create-go-starter
GOOS=darwin GOARCH=arm64 go build -o create-go-starter-macos-arm ./cmd/create-go-starter
```

## Méthode 3: Binaire pré-compilé (À venir)

> **Note**: Cette méthode sera disponible une fois que le projet publiera des releases avec des binaires pré-compilés.

Lorsque disponible, vous pourrez télécharger les binaires depuis la page [Releases](https://github.com/tky0065/go-starter-kit/releases).

### Installation sur macOS/Linux

```bash
# Télécharger le binaire (remplacez VERSION par la version souhaitée)
curl -L https://github.com/tky0065/go-starter-kit/releases/download/vVERSION/create-go-starter-macos -o create-go-starter

# Donner les permissions d'exécution
chmod +x create-go-starter

# Déplacer vers un répertoire dans PATH
sudo mv create-go-starter /usr/local/bin/
```

### Installation sur Windows

1. Téléchargez `create-go-starter.exe` depuis la page Releases
2. Placez le fichier dans un répertoire de votre choix
3. Ajoutez ce répertoire à votre variable PATH système

## Vérification de l'installation

Après installation, vérifiez que l'outil fonctionne correctement:

```bash
# Afficher l'aide
create-go-starter --help

# Devrait afficher quelque chose comme:
# Usage: create-go-starter <project-name>
#   -h, --help    Show this help message
```

### Test de création de projet

Testez en créant un projet simple:

```bash
# Créer un projet test
create-go-starter test-project

# Vérifier que le projet a été créé
ls -la test-project/

# Nettoyer
rm -rf test-project
```

Si vous voyez la structure du projet créée avec succès (en vert dans le terminal), l'installation fonctionne parfaitement!

## Mise à jour

### Mise à jour via go install (Méthode 1)

Si vous avez installé avec la Méthode 1 (installation directe), réexécutez simplement la commande:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Go téléchargera et installera automatiquement la dernière version disponible. C'est la méthode la plus simple!

### Mise à jour depuis les sources (Méthode 2)

Si vous avez installé depuis les sources:

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

### Vérifier la version

> **Note**: La commande `--version` n'est pas encore implémentée dans la version actuelle.

Une fois implémentée:

```bash
create-go-starter --version
```

Pour l'instant, vérifiez le commit Git:

```bash
cd /path/to/go-starter-kit
git log -1 --oneline
```

## Désinstallation

### Si installé via go install ou copie manuelle

```bash
# Supprimer le binaire
rm $(go env GOPATH)/bin/create-go-starter

# Ou si installé dans /usr/local/bin
sudo rm /usr/local/bin/create-go-starter
```

### Si installé depuis les sources (dans le repository)

```bash
# Supprimer le binaire local
cd /path/to/go-starter-kit
rm create-go-starter

# Optionnel: Supprimer le repository complet
cd ..
rm -rf go-starter-kit
```

### Nettoyer le cache Go

Pour libérer de l'espace:

```bash
# Nettoyer le cache de modules
go clean -modcache

# Nettoyer le cache de build
go clean -cache
```

## Résolution de problèmes

### Problème: "command not found: create-go-starter"

**Causes possibles**:
1. `$GOPATH/bin` n'est pas dans votre `PATH`
2. Le binaire n'a pas été installé correctement
3. **Le cache du shell n'a pas été rechargé** (cause la plus fréquente)

**Solutions**:

**Solution 1: Recharger le cache du shell (:material-flash: Rapide - Essayez ça d'abord!)**

Après `go install`, votre shell (zsh/bash) peut avoir une version en cache de la liste des commandes disponibles. Rechargez-la :

```bash
# Pour zsh et bash
hash -r

# Puis vérifier
which create-go-starter
create-go-starter --help
```

**Solution 2: Redémarrer le terminal**

Fermez et rouvrez votre terminal. C'est souvent la solution la plus simple !

**Solution 3: Vérifier et configurer le PATH**

Si les solutions précédentes ne fonctionnent pas :

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

**Solution 4: Utiliser le chemin complet temporairement**

En attendant de résoudre le PATH :

```bash
$(go env GOPATH)/bin/create-go-starter mon-projet
```

### Problème: "permission denied" sur macOS

**Cause**: macOS Gatekeeper bloque les binaires non signés.

**Solution**:

```bash
# Donner les permissions d'exécution
chmod +x create-go-starter

# Autoriser l'exécution (macOS)
xattr -d com.apple.quarantine create-go-starter
```

Ou autorisez l'application dans:
`Préférences Système > Sécurité et confidentialité > Général > Autoriser quand même`

### Problème: Erreurs de compilation

**Cause**: Version de Go trop ancienne ou dépendances manquantes.

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

### Problème: "go: cannot find main module"

**Cause**: Vous n'êtes pas dans le bon répertoire.

**Solution**:

```bash
# Assurez-vous d'être dans le repository go-starter-kit
cd /path/to/go-starter-kit

# Vérifier que go.mod existe
ls go.mod

# Puis réessayer l'installation
go install ./cmd/create-go-starter
```

### Problème: Conflits de versions Go

**Cause**: Plusieurs versions de Go installées.

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

### Problème: Build lent

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

## Installation pour le développement

Si vous voulez contribuer au projet, suivez ces étapes supplémentaires:

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

## Prochaines étapes

Maintenant que `create-go-starter` est installé, consultez:

- **[Guide d'utilisation](./usage.md)** - Apprendre à utiliser l'outil
- **[Guide des projets générés](./generated-project-guide.md)** - Développer avec les projets créés
- **[Guide de contribution](./contributing.md)** - Contribuer au projet

Ou commencez immédiatement:

```bash
create-go-starter mon-premier-projet
cd mon-premier-projet
# Suivez les instructions affichées!
```

Bon coding! :material-rocket-launch:
