# Story 5.2: Installation automatique des dépendances Go

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want que le CLI génère un script `setup.sh` qui installe les dépendances et configure le projet,
so that l'installation soit centralisée et automatisée pour un projet immédiatement fonctionnel.

## Acceptance Criteria

1.  Un script `setup.sh` est généré à la racine du projet.
2.  Le script `setup.sh` est exécutable (`chmod +x`).
3.  Le script contient une commande pour installer les dépendances Go (`go mod tidy`).
4.  Le CLI n'installe plus les dépendances directement, mais guide l'utilisateur pour qu'il exécute `setup.sh`.
5.  Un message de progression informatif s'affiche dans `setup.sh` pendant l'opération : "📦 Installation des dépendances...".
6.  En cas d'échec dans le script (ex: absence de Go), le script s'arrête avec un message d'erreur.

## Tasks / Subtasks

- [x] Créer le template du script `setup.sh` (AC: 1, 3, 5, 6)
    - [x] Créer `cmd/create-go-starter/templates.go` avec la fonction `SetupScriptTemplate`.
    - [x] Le script doit vérifier les prérequis (go, docker, etc.).
    - [x] Le script doit installer les dépendances avec `go mod tidy`.
    - [x] Le script doit générer un secret JWT.
    - [x] Le script doit aider à la configuration de PostgreSQL.
- [x] Intégrer la génération du script dans `generator.go` (AC: 1, 2)
    - [x] Ajouter `setup.sh` à la liste des fichiers à générer.
    - [x] Rendre le script `setup.sh` exécutable après sa création.
- [x] Supprimer l'ancienne logique d'installation de `main.go` (AC: 4)
    - [x] Supprimer l'appel à `installGoDependencies`.
    - [x] Supprimer les fichiers `deps.go` et `deps_test.go` devenus inutiles.
- [x] Mettre à jour les messages à l'utilisateur dans `main.go` (AC: 4)
    - [x] Les instructions de fin doivent clairement indiquer d'exécuter `./setup.sh`.

## Dev Notes

- **Architecture :** La responsabilité de l'installation des dépendances et de la configuration du projet est déléguée au script `setup.sh`. Le CLI se concentre sur la génération des fichiers.
- **Conventions :** Le script `setup.sh` utilise des couleurs pour une meilleure lisibilité.
- **Handoff :** L'utilisateur est clairement guidé vers la prochaine étape (`cd <projet> && ./setup.sh`).

### Project Structure Notes

- **Composants modifiés :**
    - `cmd/create-go-starter/templates.go` (ajout de `SetupScriptTemplate`)
    - `cmd/create-go-starter/generator.go` (génération de `setup.sh`)
    - `cmd/create-go-starter/main.go` (suppression de l'appel direct, mise à jour des messages)
- **Composants supprimés :**
    - `cmd/create-go-starter/deps.go`
    - `cmd/create-go-starter/deps_test.go`
- **Impact :** Simplifie le code du CLI et centralise la logique de setup dans un script unique et réutilisable, améliorant la clarté et la maintenance.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 5.2: Installation automatique des dépendances Go]
- [Source: cmd/create-go-starter/main.go] flow de génération.

## Dev Agent Record

### Agent Model Used

- Gemini 2.0 Flash (BMad SM Mode) - Story creation
- Gemini Pro - Implementation

### Completion Notes List

- ✅ Créé un template pour `setup.sh` dans `templates.go`.
- ✅ Le script `setup.sh` gère les dépendances, JWT, et PostgreSQL.
- ✅ Intégré la génération et les permissions de `setup.sh` dans `generator.go`.
- ✅ Supprimé la fonction `installGoDependencies` et les fichiers `deps.go` et `deps_test.go`.
- ✅ Mis à jour `main.go` pour supprimer l'appel direct et guider l'utilisateur vers `setup.sh`.
- ✅ Test E2E validé : la génération du projet et l'exécution de `setup.sh` fonctionnent comme prévu.

### File List

- cmd/create-go-starter/main.go (modifié)
- cmd/create-go-starter/generator.go (modifié)
- cmd/create-go-starter/templates.go (modifié)
- cmd/create-go-starter/deps.go (supprimé)
- cmd/create-go-starter/deps_test.go (supprimé)

## Change Log

- 2026-01-14: Story implémentée - La CLI génère maintenant un script `setup.sh` complet pour l'installation et la configuration, et ne tente plus d'installer les dépendances directement.
