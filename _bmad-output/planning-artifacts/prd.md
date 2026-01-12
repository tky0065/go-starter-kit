---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
inputDocuments:
  - '_bmad-output/planning-artifacts/product-brief-go-starter-kit-2026-01-07.md'
workflowType: 'prd'
lastStep: 10
project_name: 'go-starter-kit'
user_name: 'Yacoubakone'
date: '2026-01-07'
briefCount: 1
researchCount: 0
brainstormingCount: 0
projectDocsCount: 0
---

# Product Requirements Document - go-starter-kit

**Author:** Yacoubakone
**Date:** 2026-01-07

## Executive Summary

**go-starter-kit** est un outil de productivité pour développeurs conçu pour éliminer les frictions liées au démarrage de projets API en Go. Il permet de passer d'une idée à une API "production-ready" en moins de 5 minutes grâce à une approche opinionated et pragmatique. En combinant une interface CLI moderne avec une architecture hexagonale simplifiée et une stack technique pré-configurée (Fiber, GORM, fx), il garantit un équilibre entre rapidité de développement et respect des standards industriels.

### Ce qui rend ce projet spécial

*   **Expérience "Zero-to-Hero" :** Un CLI de scaffolding inspiré des meilleurs standards (type `create-next-app`) pour une prise en main instantanée.
*   **Best Practices par défaut :** L'authentification JWT, la validation, la documentation Swagger, Docker et la CI/CD ne sont plus des options à configurer mais des acquis dès le premier jour.
*   **Architecture Inclusive :** Une structure hexagonale "Lite" qui offre les avantages de la maintenabilité sans la complexité dogmatique, rendant le code accessible aux juniors tout en étant robuste pour les seniors.

## Classification du Projet

**Type Technique :** Outil de développement (CLI & Boilerplate API)
**Domaine :** Outils de développement logiciel
**Complexité :** Moyenne
**Contexte du Projet :** Greenfield - Nouveau projet

## Success Criteria

### User Success

*   **Vitesse d'Onboarding (Time-to-First-Value) :** L'utilisateur doit avoir une API fonctionnelle sur `localhost` en **moins de 5 minutes** après avoir lancé la commande d'installation.
*   **Simplicité d'Extension :** Un développeur (même junior) doit être capable de comprendre la structure et d'ajouter son premier endpoint métier personnalisé en **moins de 30 minutes**.
*   **Confiance en Production :** L'utilisateur doit pouvoir déployer son projet en production (avec CI/CD et conteneurisation) en **moins d'une semaine**, avec la certitude que les bases (sécurité, tests) sont solides.
*   **Réduction de la Charge Mentale :** L'utilisateur ne doit pas avoir à prendre de décisions d'infrastructure (quel ORM ? quel routeur ?) lors de l'initialisation.

### Business Success

*   **Adoption (Lancement - 3 mois) :** Atteindre 1 000+ installations du CLI et 500+ étoiles GitHub.
*   **Rétention & Loyauté (12 mois) :** 60 0es utilisateurs créent plus d'un projet avec le starter kit.
*   **Réputation (12 mois) :** Devenir un standard reconnu (Top 5 des starters Go sur GitHub avec 5 000+ étoiles) et atteindre un Net Promoter Score (NPS) > 50.
*   **Communauté :** Bâtir une communauté active avec des contributeurs réguliers et un temps de réponse aux issues < 24h.

### Technical Success

*   **Fiabilité du Scaffolding :** Taux de réussite de la génération de projet de 100ur les principaux OS (Mac, Linux, Windows).
*   **Qualité du Code Généré :** Le code généré doit passer par défaut tous les linters (golangci-lint), avoir une couverture de tests > 70%, et être libre de vulnérabilités connues.
*   **Performance de Développement :** L'environnement local (`make dev`) doit démarrer (API + DB) en moins de 30 secondes.
*   **Sécurité par Défaut :** Les projets générés doivent inclure par défaut les meilleures pratiques de sécurité (JWT robuste, CORS configuré, Docker non-root).

### Measurable Outcomes

*   **95%+** des utilisateurs atteignent le "Hello World" en < 5 min.
*   **90%+** de taux de complétion de l'onboarding.
*   **< 10** bugs critiques ouverts sur le repo principal à tout moment.

## Product Scope

### MVP - Minimum Viable Product

L'objectif est une expérience "zéro friction" pour une stack opinionated spécifique.
*   **CLI :** Commande unique `create-go-starter` (scaffolding, git init, install dependencies).
*   **Stack Core :** Fiber (Web), GORM (ORM), PostgreSQL (DB), fx (DI).
*   **Architecture :** Hexagonale simplifiée ("Lite").
*   **Fonctionnalités "Batteries Included" :** Auth JWT complète, Validation automatique, Gestion d'erreurs centralisée, Swagger UI auto-généré.
*   **DevOps Ready :** Dockerfile optimisé, docker-compose, CI/CD GitHub Actions de base.

### Growth Features (Post-MVP)

Fonctionnalités pour élargir la base d'utilisateurs une fois le core validé.
*   Support Multi-Base de données (MySQL, SQLite).
*   Observabilité avancée (Prometheus/Grafana).
*   Génération de code (CRUD scaffolding via CLI).

### Vision (Future)

Transformation en plateforme écosystémique.
*   Architecture de Plugins et Marketplace communautaire.
*   Support Multi-Framework (Gin, Echo).
*   Interface Web/Dashboard pour la gestion de projets.

## User Journeys

### Journey 1: Alex - De la Confusion à la Clarté (Apprenant/Junior)
Alex est un développeur JavaScript en reconversion vers Go. Il est enthousiaste mais frustré : ses trois dernières tentatives de créer une API REST propre ont fini en "code spaghetti". Il se perd entre les tutoriels contradictoires sur la structure des dossiers et la configuration de JWT. Il a peur de "mal faire".

Un soir, il découvre `go-starter-kit`. Sceptique, il lance `go install ...` puis `create-go-starter learning-api`. En 15 secondes, le terminal lui affiche "Success". Il suit l'instruction `cd learning-api && make dev`.
**Le moment clé :** Il voit les logs défiler, la base de données Docker démarrer, et l'API se lancer sur le port 8080. Il ouvre `localhost:8080/swagger` et voit une documentation interactive complète. Il teste l'authentification : ça marche.
Il ouvre le code. Au lieu d'être effrayé, il trouve des commentaires clairs. Il voit exactement où ajouter son modèle "Todo". En copiant le pattern existant (Service -> Repository -> Handler), il crée son premier endpoint en 30 minutes. Il lance `make test` et voit que les tests passent. Pour la première fois, Alex se sent comme un "vrai" développeur Go professionnel.

### Journey 2: Sarah - Le Sprint du MVP (Tech Founder)
Sarah a une idée de SaaS B2B révolutionnaire et un rendez-vous investisseur dans 10 jours. Elle doit prouver que son concept fonctionne, pas juste montrer des slides. Elle ne peut pas se permettre de perdre 3 jours à configurer l'infrastructure, mais elle sait qu'un code "sale" tuera sa startup plus tard si elle réussit.

Elle utilise `go-starter-kit` un lundi matin. À 10h00, l'échafaudage est prêt. À 11h00, elle a déjà modifié le modèle `User` pour ajouter des champs métier. L'après-midi, elle se concentre uniquement sur sa "Killer Feature" (un algorithme de matching). Elle n'a pas à se soucier de savoir si ses tokens JWT sont sécurisés ou si son Dockerfile est optimisé pour la prod : c'est déjà fait.
**Le dénouement :** Le vendredi, elle déploie sur un VPS bon marché via le GitHub Action fourni. Lors de la démo investisseur, l'API répond en millisecondes et ne crashe pas. Elle a gagné 2 semaines de dev et sécurisé son pré-seed.

### Journey 3: Marc - L'Harmonisation de l'Équipe (Tech Lead)
L'équipe de Marc a grandi trop vite. Ils ont 5 microservices Go, chacun avec une structure différente. L'onboarding des nouveaux est un cauchemar ("Ah, sur ce projet, les handlers sont dans ce dossier..."). Marc passe son temps en Code Review à corriger des erreurs triviales de configuration.

Il décide d'adopter `go-starter-kit` comme standard interne. Il génère un projet, le customize avec les libs spécifiques de l'entreprise, et en fait un "template interne".
**L'impact :** Deux semaines plus tard, un nouveau junior arrive. Marc lui dit : "Lance le générateur". Le junior est opérationnel le jour même car la structure est standardisée. Les Code Reviews se concentrent désormais sur la logique métier complexe, plus sur la syntaxe ou l'architecture de base. La "charge cognitive" de l'équipe a chuté, la vélocité a augmenté.

### Journey 4: Sam - Le Gardien du Temple (DevOps/SRE)
Sam est fatigué de recevoir des Dockerfiles de développeurs qui pèsent 1GB ou qui tournent en root. Quand Sarah lui envoie le lien de son repo basé sur `go-starter-kit`, il s'attend au pire.

Il ouvre le `Dockerfile` : c'est un multi-stage build optimisé sur Alpine. Il vérifie le `docker-compose.yml` : les healthchecks sont là. Il regarde la config : tout est piloté par variables d'environnement.
**La réaction :** "Enfin !". Il n'a presque rien à faire pour intégrer ce projet dans leur pipeline Kubernetes existant. Le projet est "Ops-friendly" par défaut.

### Journey Requirements Summary

Ces parcours mettent en lumière des exigences critiques :

*   **Onboarding (Alex) :** Documentation Swagger automatique, commandes Makefile simples (`dev`, `test`), commentaires de code pédagogiques.
*   **Vitesse (Sarah) :** Scaffolding instantané, Auth & DB pré-configurés, CI/CD de base prêt à l'emploi.
*   **Standardisation (Marc) :** Architecture hexagonale stricte mais simple, gestion d'erreurs centralisée et uniforme.
*   **Opérabilité (Sam) :** Configuration par ENV, Docker optimisé, Health checks, Graceful shutdown.

## Developer Tool & API Backend Specific Requirements

### Project-Type Overview
Le produit est un **générateur de code (CLI)** qui produit un **backend API** structuré. L'objectif est de fournir une fondation technique "opinionated" qui respecte les standards de production sans intervention manuelle lourde de l'utilisateur.

### Technical Architecture Considerations
*   **Scaffolding Engine :** Le CLI doit être capable de cloner/générer une structure de dossiers complexe, d'initialiser un module Go et de gérer le remplacement dynamique du nom du projet dans les fichiers source.
*   **Dependency Injection (fx) :** L'architecture doit utiliser `uber-go/fx` pour gérer le cycle de vie des composants (DB, Server, Handlers) de manière modulaire.
*   **Hexagonal Architecture (Lite) :** Séparation claire entre les Entités (Domain), les Interfaces (Ports) et les Implémentations (Adapters) pour garantir la testabilité.

### Command Structure (CLI)
*   **Utilisation :** `create-go-starter <project-name>`
*   **Mode MVP :** Non-interactif. Génération d'une configuration standard complète pour éviter la paralysie décisionnelle.

### API & Authentication Model
*   **Organisation des Routes (Grouping) :** Utilisation systématique des groupes de routes (Fiber Groups). Toutes les routes API doivent être préfixées par `/api` (ex: `/api/v1/auth`, `/api/v1/users`) pour permettre une versioning facile et une séparation claire des responsabilités.
*   **Modèle d'Authentification :** JWT (JSON Web Token) avec support des Refresh Tokens et rotation sécurisée. 
*   **Endpoints de base (MVP) :**
    *   `POST /api/v1/auth/register` : Création de compte.
    *   `POST /api/v1/auth/login` : Obtention des tokens.
    *   `POST /api/v1/auth/refresh` : Renouvellement du token d'accès.
    *   `GET /api/v1/users/me` : Profil de l'utilisateur connecté.
*   **Documentation :** Swagger UI auto-généré via des annotations de code, accessible sur `/swagger`.

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach:** Experience & Problem-Solving. La priorité absolue est de tenir la promesse du "< 5 minutes pour une API pro". Le focus est mis sur la fluidité du CLI et l'aspect "clef en main" de la stack Auth+DB+Swagger.
**Resource Requirements:** 1 à 2 développeurs Go seniors (le code généré étant un modèle pour d'autres, il doit être d'une qualité technique irréprochable).

### MVP Feature Set (Phase 1)

**Core User Journeys Supported:**
*   **Alex (Junior) :** Apprentissage par l'exemple avec une structure propre et documentée.
*   **Sarah (Founder) :** Déploiement ultra-rapide d'un backend sécurisé.

**Must-Have Capabilities:**
*   CLI `create-go-starter` (scaffolding, injection du nom de projet).
*   Stack pré-configurée : Fiber, GORM, PostgreSQL, fx.
*   Authentification JWT complète (Login, Register, Refresh).
*   Groupage des routes sous `/api/v1`.
*   Un (1) exemple de module complet (User) suivant l'architecture hexagonale.
*   Swagger UI auto-généré.
*   Configuration via `.env` et Dockerfile de base.

### Post-MVP Features

**Phase 2 (Growth):**
*   Support multi-bases de données (MySQL, SQLite).
*   Génération de CRUD supplémentaire via CLI (`gsk generate`).
*   Templates de CI/CD avancés (multi-cloud).
*   Monitoring de base (Prometheus metrics).

**Phase 3 (Expansion):**
*   Support multi-frameworks (Gin, Echo).
*   Système de plugins pour ajouter des modules (Paiement, S3, etc.).
*   Dashboard web pour visualiser et gérer les projets générés.

### Risk Mitigation Strategy

**Technical Risks:** La complexité de l'injection de dépendances (`fx`) peut effrayer les juniors. 
*   **Mitigation :** Commentaires ultra-pédagogiques dans le code et documentation "Deep Dive" dans le README.
**Market Risks:** Les développeurs Go sont souvent très attachés à leur propre structure ou à des bibliothèques spécifiques (ex: SQL natif vs ORM). 
*   **Mitigation :** Assumer l'aspect "Opinionated" comme une force de gain de temps, tout en gardant l'architecture assez découplée pour permettre des remplacements.

## Functional Requirements

### 1. Project Scaffolding & Generation (CLI)
*   **FR1 :** L'utilisateur peut installer l'outil via une commande `go install`.
*   **FR2 :** L'utilisateur peut générer un nouveau projet en fournissant un nom de projet via le CLI.
*   **FR3 :** Le système peut créer automatiquement une structure de dossiers respectant l'architecture hexagonale lite.
*   **FR4 :** Le système peut injecter dynamiquement le nom du projet dans les fichiers générés (fichiers Go, `go.mod`, Dockerfile).
*   **FR5 :** Le système peut initialiser automatiquement un module Go et télécharger les dépendances nécessaires.
*   **FR6 :** Le système peut créer un fichier `.env` par défaut à partir d'un template `.env.example`.

### 2. User Authentication & Security
*   **FR7 :** Un visiteur peut créer un compte utilisateur (Register).
*   **FR8 :** Un utilisateur peut s'authentifier via des identifiants sécurisés (Login).
*   **FR9 :** Le système peut générer des jetons d'accès (Access Tokens) et de renouvellement (Refresh Tokens) sécurisés.
*   **FR10 :** Un utilisateur peut renouveler son jeton d'accès sans se reconnecter manuellement.
*   **FR11 :** Le système peut hacher les mots de passe de manière sécurisée avant stockage.
*   **FR12 :** Le système peut protéger des routes spécifiques pour n'autoriser que les utilisateurs authentifiés.

### 3. API Infrastructure & Routing
*   **FR13 :** Le système peut regrouper les routes API par domaine fonctionnel.
*   **FR14 :** Le système peut appliquer un préfixe global `/api/v1` à toutes les routes métier.
*   **FR15 :** Le système peut gérer automatiquement les erreurs et renvoyer des réponses JSON standardisées.
*   **FR16 :** Le système peut valider automatiquement les données d'entrée des requêtes HTTP selon des règles prédéfinies.
*   **FR17 :** Le système peut exposer une documentation interactive (Swagger UI) mise à jour automatiquement.

### 4. Database & Persistence
*   **FR18 :** Le système peut se connecter à une base de données PostgreSQL de manière résiliente (pool de connexions).
*   **FR19 :** Le système peut exécuter des migrations de base de données au démarrage ou via commande.
*   **FR20 :** Un développeur peut effectuer des opérations CRUD (Créer, Lire, Mettre à jour, Supprimer) sur l'entité User.

### 5. Dependency Injection & Lifecycle
*   **FR21 :** Le système peut gérer l'injection de dépendances pour tous les composants majeurs (DB, Serveur, Handlers).
*   **FR22 :** Le système peut assurer un démarrage et un arrêt "propre" (graceful shutdown) de l'application.

### 6. Developer Experience & DevOps
*   **FR23 :** Le développeur peut lancer l'application en mode développement avec rechargement automatique (hot-reload).
*   **FR24 :** Le développeur peut exécuter l'ensemble des tests (unitaires et intégration) via une commande unique.
*   **FR25 :** Le système peut être exécuté dans un environnement conteneurisé (Docker).
*   **FR26 :** Le système peut être déployé via un pipeline de CI/CD pré-configuré (GitHub Actions).

## Non-Functional Requirements

### Performance
*   **Temps de réponse de l'API :** Les endpoints de base (Auth/Health) doivent répondre en moins de **100ms** (hors latence réseau).
*   **Démarrage à froid (Cold Start) :** L'application conteneurisée doit être opérationnelle en moins de **2 secondes**.
*   **Légèreté :** L'image Docker finale doit être optimisée (utilisation de multi-stage builds) pour peser moins de **50 Mo** (hors assets volumineux).

### Security
*   **Chiffrement :** Toutes les données sensibles (mots de passe) sont hachées avec `bcrypt` (coût par défaut >= 10).
*   **JWT :** Utilisation de l'algorithme `HS256` ou `RS256` avec une gestion stricte de l'expiration.
*   **Secrets :** Aucune clé ou secret ne doit être en dur dans le code (utilisation obligatoire de `.env` ou variables d'environnement).
*   **OWASP :** Protection native contre les vulnérabilités courantes (CORS, CSRF, Injection SQL via l'ORM).

### Maintenabilité & Qualité du Code
*   **Standardisation :** 100% du code généré doit respecter les standards `golangci-lint`.
*   **Documentation :** Chaque fonction publique doit être documentée. Le README doit permettre un démarrage en 5 minutes.
*   **Testabilité :** L'architecture hexagonale doit permettre de mocker 100% des dépendances externes pour les tests unitaires.

### Opérabilité & Observabilité
*   **Graceful Shutdown :** Le système doit intercepter les signaux d'arrêt et fermer les connexions proprement en moins de **5 secondes**.
*   **Logs :** Utilisation de logs structurés (JSON) en production pour faciliter l'indexation.
*   **Health Checks :** Présence d'un endpoint `/health` pour les orchestrateurs (Kubernetes/Docker).
