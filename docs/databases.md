# Guide de sélection de base de données

**Navigation:**

- <i class="material-icons">menu_book</i> [Guide de sélection de base de données](./databases.md) ← Vous êtes ici
- <i class="material-icons">sync</i> [Guide de migration de base de données](./database-migration.md) - Comment changer de base de données
- <i class="material-icons">arrow_back</i> [Retour au README](../README.md)

---

go-starter-kit supporte **3 options de base de données** pour s'adapter aux besoins de votre projet.

## Comparaison rapide

| Base de données | Idéal pour | Complexité | Prêt production | Temps installation |
|-----------------|------------|------------|-----------------|-------------------|
| **PostgreSQL** | Apps production, requêtes complexes | Moyen | <i class="material-icons success">check_circle</i> Oui | 2 min (Docker) |
| **MySQL** | Compatibilité large, hébergement partagé | Moyen | <i class="material-icons success">check_circle</i> Oui | 2 min (Docker) |
| **SQLite** | Prototypage, petites apps, embarqué | Faible | <i class="material-icons warning">error</i> Limité | 0 min |

## Comparaison détaillée

### PostgreSQL (Par défaut)

**Commande:**
```bash
create-go-starter mon-app
# OU explicitement:
create-go-starter mon-app --database=postgres
```

**Points forts:**

- <i class="material-icons success">check</i> Fonctionnalités SQL avancées (JSON, arrays, recherche full-text)
- <i class="material-icons success">check</i> Excellentes performances et fiabilité
- <i class="material-icons success">check</i> Conforme ACID, forte intégrité des données
- <i class="material-icons success">check</i> Idéal pour les requêtes complexes et l'analytique
- <i class="material-icons success">check</i> Communauté active et écosystème riche

**Limitations:**

- <i class="material-icons warning">warning</i> Nécessite Docker pour le développement local
- <i class="material-icons warning">warning</i> Légèrement plus gourmand en ressources que MySQL

**Quand l'utiliser:**
- Applications de production avec données complexes
- Applications nécessitant des fonctionnalités SQL avancées
- Projets nécessitant une forte intégrité des données

**Configuration Docker:**
```yaml
# Automatiquement inclus dans docker-compose.yml
docker-compose up -d
```

**Format DSN:**
```
user:password@tcp(host:5432)/dbname?sslmode=disable
```

---

### MySQL/MariaDB

**Commande:**
```bash
create-go-starter mon-app --database=mysql
```

**Points forts:**

- <i class="material-icons success">check</i> Large compatibilité et support d'hébergement
- <i class="material-icons success">check</i> Excellent pour les charges de lecture intensives
- <i class="material-icons success">check</i> Écosystème et outils matures
- <i class="material-icons success">check</i> Facile de trouver des fournisseurs d'hébergement

**Limitations:**

- <i class="material-icons warning">warning</i> Moins de fonctionnalités avancées que PostgreSQL
- <i class="material-icons warning">warning</i> Quelques variations entre MySQL et MariaDB

**Quand l'utiliser:**
- Environnements d'hébergement partagé
- Applications à forte charge de lecture
- Équipes familières avec MySQL
- Besoin de large compatibilité d'hébergement

**Configuration Docker:**
```yaml
# Automatiquement inclus dans docker-compose.yml
docker-compose up -d
```

**Format DSN:**
```
user:password@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

---

### SQLite

**Commande:**
```bash
create-go-starter mon-app --database=sqlite
```

**Points forts:**

- <i class="material-icons success">check</i> Zéro configuration (pas de serveur nécessaire)
- <i class="material-icons success">check</i> Parfait pour le prototypage rapide
- <i class="material-icons success">check</i> Base de données en fichier unique (sauvegarde/partage facile)
- <i class="material-icons success">check</i> Idéal pour les tests et le développement
- <i class="material-icons success">check</i> Très rapide pour les petits jeux de données

**Limitations:**

- <i class="material-icons warning">warning</i> Écritures concurrentes limitées (verrouille toute la DB)
- <i class="material-icons warning">warning</i> Pas de gestion utilisateurs/permissions
- <i class="material-icons warning">warning</i> Non adapté pour la production à fort trafic
- <i class="material-icons warning">warning</i> Scalabilité limitée

**Quand l'utiliser:**
- Prototypage rapide et MVPs
- Applications desktop
- Systèmes embarqués
- Développement et tests
- Production à petite échelle (<100 utilisateurs concurrents)

**Pas besoin de Docker:**
```bash
# Lancez simplement votre app, le fichier SQLite est créé automatiquement
go run ./cmd/main.go
# Crée: ./my_database.db
```

**Format DSN:**
```
./database.db
```

---

## Matrice de décision

**Choisir PostgreSQL si:**

- <i class="material-icons">center_focus_strong</i> Vous hésitez (c'est le choix par défaut pour une bonne raison)
- <i class="material-icons">center_focus_strong</i> Vous avez besoin de fiabilité niveau production
- <i class="material-icons">center_focus_strong</i> Vous avez des données relationnelles complexes

**Choisir MySQL si:**

- <i class="material-icons">center_focus_strong</i> Vous utilisez un hébergement partagé
- <i class="material-icons">center_focus_strong</i> Votre équipe connaît bien MySQL
- <i class="material-icons">center_focus_strong</i> Vous avez des charges de lecture intensives

**Choisir SQLite si:**

- <i class="material-icons">center_focus_strong</i> Vous prototypez ou construisez un MVP
- <i class="material-icons">center_focus_strong</i> Vous voulez zéro infrastructure
- <i class="material-icons">center_focus_strong</i> Vous avez une petite base d'utilisateurs (<100 concurrents)

---

## Exemples de configuration

### PostgreSQL

**.env.example:**
```bash
# Configuration base de données (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=monapp
DB_SSLMODE=disable
```

### MySQL

**.env.example:**
```bash
# Configuration base de données (MySQL)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=monapp
```

### SQLite

**.env.example:**
```bash
# Configuration base de données (SQLite - embarqué)
DB_NAME=monapp.db
```

---

## Guide de migration

Voir [database-migration.md](./database-migration.md) pour les instructions détaillées de migration entre bases de données.

---

## Considérations de performance

### PostgreSQL
- Idéal pour: Requêtes complexes, transactions ACID, écritures concurrentes
- Performance écriture: Excellente (architecture MVCC)
- Performance lecture: Excellente avec indexation appropriée
- Connection pooling: Recommandé pour la production

### MySQL
- Idéal pour: Charges de lecture intensives, requêtes simples
- Performance écriture: Bonne (verrouillage au niveau ligne)
- Performance lecture: Excellente (cache de requêtes)
- Connection pooling: Recommandé pour la production

### SQLite
- Idéal pour: Scénarios mono-utilisateur, faible concurrence
- Performance écriture: Limitée (verrouillage au niveau base de données)
- Performance lecture: Excellente pour petits jeux de données
- Connection pooling: Non applicable (fichier)

---

## Questions fréquentes

### Puis-je changer de base de données plus tard?
Oui, mais cela nécessite de régénérer le projet avec le nouveau flag de base de données et de migrer les données. Voir [database-migration.md](./database-migration.md) pour les détails.

### Quelle base de données dois-je utiliser pour mon SaaS?
Pour une application SaaS de production, nous recommandons **PostgreSQL** pour sa fiabilité, conformité ACID et fonctionnalités avancées. MySQL est aussi un bon choix si vous le connaissez mieux.

### Puis-je utiliser SQLite en production?
SQLite peut être utilisé pour de la production à petite échelle (<100 utilisateurs concurrents), mais nous recommandons PostgreSQL ou MySQL pour les applications qui prévoient de croître.

### Ai-je besoin de Docker?
- **PostgreSQL**: Oui (pour le développement local)
- **MySQL**: Oui (pour le développement local)
- **SQLite**: Non (base de données embarquée)

### Qu'en est-il de MongoDB ou NoSQL?
Le support NoSQL (MongoDB) a été considéré mais reporté aux versions futures. L'accent actuel est mis sur les bases de données SQL avec support GORM.

---

## Ressources additionnelles

- [Documentation PostgreSQL](https://www.postgresql.org/docs/)
- [Documentation MySQL](https://dev.mysql.com/doc/)
- [Documentation SQLite](https://www.sqlite.org/docs.html)
- [Documentation GORM](https://gorm.io/docs/)
- [Guide de migration de base de données](./database-migration.md)

---

**Dernière mise à jour:** 2026-02-09
