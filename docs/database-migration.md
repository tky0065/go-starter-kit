# Guide de migration de base de données

**Navigation:**

- <i class="material-icons">menu_book</i> [Guide de sélection de base de données](./databases.md) - Quelle base de données choisir
- <i class="material-icons">sync</i> [Guide de migration de base de données](./database-migration.md) ← Vous êtes ici
- <i class="material-icons">arrow_back</i> [Retour au README](../README.md)

---

Ce guide vous aide à migrer entre différents systèmes de bases de données dans les projets go-starter-kit.

## Prérequis

!!! warning "Important"
    La migration de base de données n'est pas automatisée. Vous devrez:
1. Exporter les données de la base source
2. Régénérer le projet avec la base cible
3. Importer les données vers la base cible
4. Tester minutieusement

---

## Chemins de migration

### PostgreSQL ↔ MySQL

**Difficulté:** <i class="material-icons success small">circle</i> Facile (tous deux SQL, compatibles GORM)

**Étapes:**

1. **Exporter les données de la base source**
   ```bash
   # Export PostgreSQL
   pg_dump -U postgres -h localhost -d monapp -F c -b -v -f backup.dump

   # Export MySQL
   mysqldump -u root -p monapp > backup.sql
   ```

2. **Régénérer le projet avec la base cible**
   ```bash
   # Dans votre répertoire de projet
   create-go-starter monapp --database=mysql  # ou postgres
   ```

3. **Convertir la syntaxe SQL si nécessaire**
   - La plupart des différences sont gérées par GORM
   - Changements courants:
     - Serial → AUTO_INCREMENT (géré par GORM)
     - BOOLEAN → TINYINT(1) (géré par GORM)
     - TEXT vs LONGTEXT (généralement compatible)

4. **Importer les données**
   ```bash
   # Import PostgreSQL
   pg_restore -U postgres -h localhost -d monapp backup.dump

   # Import MySQL
   mysql -u root -p monapp < backup.sql
   ```

5. **Tester les migrations et requêtes**
   ```bash
   go run ./cmd/main.go
   # Vérifier que tous les endpoints fonctionnent correctement
   ```

**Problèmes courants:**
- Serial vs AUTO_INCREMENT (géré par GORM)
- Certains types de données diffèrent (TEXT vs LONGTEXT)
- La syntaxe des fonctions peut différer (DATE_ADD vs INTERVAL)

---

### SQL → SQLite (Dégradation)

**Difficulté:** <i class="material-icons warning small">circle</i> Moyen (réduction de fonctionnalités)

**Quand faire ceci:**
- Passage de production au développement local
- Création d'une version démo portable
- Simplification de l'infrastructure pour petits projets

**Étapes:**

1. **Exporter les données en SQL ou CSV**
   ```bash
   # Export PostgreSQL en CSV
   psql -U postgres -d monapp -c "COPY users TO '/tmp/users.csv' WITH CSV HEADER;"

   # Export MySQL en CSV
   mysql -u root -p monapp -e "SELECT * FROM users INTO OUTFILE '/tmp/users.csv'
     FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
   ```

2. **Régénérer le projet**
   ```bash
   create-go-starter monapp --database=sqlite
   ```

3. **Importer un jeu de données limité**
   ```bash
   # SQLite ne gère pas bien les grands jeux de données
   # Importer uniquement les données essentielles
   sqlite3 monapp.db < import.sql
   ```

4. **Retirer les fonctionnalités SQL avancées**
   - Pas de procédures stockées
   - Support des transactions limité
   - Pas d'écritures concurrentes

**Limitations:**

- <i class="material-icons warning">warning</i> Pas d'écritures concurrentes (verrouillage au niveau base)
- <i class="material-icons warning">warning</i> Types de données limités
- <i class="material-icons warning">warning</i> Pas de procédures stockées
- <i class="material-icons warning">warning</i> Pas de gestion utilisateurs/permissions
- <i class="material-icons warning">warning</i> Non adapté pour la production à grande échelle

**Cas d'usage recommandés:**
- Environnements de développement local
- Démos portables
- Pipelines de tests et CI/CD
- Applications à petite échelle (<100 utilisateurs)

---

### SQLite → SQL (Mise à niveau)

**Difficulté:** <i class="material-icons success small">circle</i> Facile (ajout de fonctionnalités)

**Quand faire ceci:**
- Passage du prototype à la production
- Besoin d'accès concurrent en écriture
- Nécessité de fonctionnalités SQL avancées

**Étapes:**

1. **Exporter les données SQLite**
   ```bash
   # Dump en SQL
   sqlite3 monapp.db .dump > backup.sql

   # Ou exporter les tables individuellement
   sqlite3 monapp.db -header -csv "SELECT * FROM users;" > users.csv
   ```

2. **Régénérer le projet**
   ```bash
   create-go-starter monapp --database=postgres  # ou mysql
   ```

3. **Convertir le schéma si nécessaire**
   - Les migrations GORM devraient gérer la plupart des différences
   - Réviser et ajuster tout SQL personnalisé

4. **Importer les données**
   ```bash
   # PostgreSQL
   psql -U postgres -d monapp < converted_backup.sql

   # MySQL
   mysql -u root -p monapp < converted_backup.sql
   ```

5. **Tester minutieusement**
   - Vérifier l'intégrité des données
   - Tester tous les endpoints API
   - Exécuter la suite de tests complète

---

## Exemples d'export/import de données

### Export depuis PostgreSQL

```bash
# Dump complet de la base (compressé)
pg_dump -U postgres -h localhost -d monapp -F c -b -v -f monapp_backup.dump

# Format SQL (lisible par l'humain)
pg_dump -U postgres -h localhost -d monapp -f monapp_backup.sql

# Tables spécifiques uniquement
pg_dump -U postgres -h localhost -d monapp -t users -t posts -f partial_backup.sql

# Export CSV pour une table spécifique
psql -U postgres -d monapp -c "COPY users TO '/tmp/users.csv' WITH CSV HEADER;"
```

### Export depuis MySQL

```bash
# Dump complet de la base
mysqldump -u root -p monapp > monapp_backup.sql

# Tables spécifiques uniquement
mysqldump -u root -p monapp users posts > partial_backup.sql

# Export CSV
mysql -u root -p monapp -e "SELECT * FROM users INTO OUTFILE '/tmp/users.csv'
  FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
```

### Export depuis SQLite

```bash
# Dump complet de la base
sqlite3 monapp.db .dump > monapp_backup.sql

# Table spécifique
sqlite3 monapp.db "SELECT * FROM users;" > users.txt

# Export CSV
sqlite3 monapp.db -header -csv "SELECT * FROM users;" > users.csv
```

### Import vers PostgreSQL

```bash
# Depuis un dump PostgreSQL
pg_restore -U postgres -h localhost -d monapp monapp_backup.dump

# Depuis un fichier SQL
psql -U postgres -d monapp < monapp_backup.sql

# Depuis CSV
psql -U postgres -d monapp -c "\COPY users FROM '/tmp/users.csv' WITH CSV HEADER;"
```

### Import vers MySQL

```bash
# Depuis un fichier SQL
mysql -u root -p monapp < monapp_backup.sql

# Depuis CSV
mysql -u root -p monapp -e "LOAD DATA LOCAL INFILE '/tmp/users.csv'
  INTO TABLE users
  FIELDS TERMINATED BY ','
  ENCLOSED BY '\"'
  LINES TERMINATED BY '\n'
  IGNORE 1 ROWS;"
```

### Import vers SQLite

```bash
# Depuis un fichier SQL
sqlite3 monapp.db < monapp_backup.sql

# Depuis CSV
sqlite3 monapp.db <<EOF
.mode csv
.import users.csv users
EOF
```

---

## Tests après migration

**Liste de vérification:**
- [ ] Toutes les migrations s'exécutent avec succès
- [ ] Intégrité des données vérifiée (comptages, relations)
- [ ] Tous les endpoints API fonctionnent correctement
- [ ] L'authentification/autorisation fonctionne
- [ ] Tous les tests passent
- [ ] Les performances sont acceptables
- [ ] Aucune perte de données détectée

**Commandes de validation:**
```bash
# Exécuter tous les tests
go test ./...

# Vérifier l'intégrité des données
go run ./cmd/main.go
# Tester tous les endpoints API manuellement ou avec des tests automatisés

# Vérifier les comptages d'enregistrements
# Comparer les comptages des bases source et cible
```

---

## Plan de rollback

Ayez toujours un plan de rollback avant de migrer:

1. **Sauvegarder la base de données d'origine**
   ```bash
   # Garder la sauvegarde de la base d'origine en sécurité
   cp monapp_backup.dump monapp_backup_SAFE.dump
   ```

2. **Tester la migration en staging d'abord**
   - Ne jamais tester les migrations directement en production
   - Utiliser un environnement de staging identique à la production

3. **Planifier une fenêtre de maintenance**
   - Planifier la migration pendant les périodes de faible trafic
   - Communiquer la maintenance aux utilisateurs

4. **Documenter les étapes de rollback**
   ```bash
   # Exemple de script de rollback
   #!/bin/bash
   echo "Retour à la base de données d'origine..."
   pg_restore -U postgres -h localhost -d monapp monapp_backup_SAFE.dump
   echo "Rollback terminé"
   ```

---

## Liste de vérification de migration

**Avant la migration:**
- [ ] Sauvegarde complète de la base source
- [ ] Sauvegarde vérifiée (test de restauration)
- [ ] Environnement de staging prêt
- [ ] Plan de migration documenté
- [ ] Plan de rollback préparé
- [ ] Maintenance planifiée
- [ ] Équipe notifiée

**Pendant la migration:**
- [ ] Arrêter l'application (empêcher les changements de données)
- [ ] Exporter les données de la source
- [ ] Régénérer le projet avec la nouvelle base
- [ ] Importer les données vers la cible
- [ ] Exécuter les migrations GORM
- [ ] Vérifier l'intégrité des données

**Après la migration:**
- [ ] Tous les tests passent
- [ ] Endpoints API vérifiés
- [ ] Performances acceptables
- [ ] Monitoring en place
- [ ] Ancienne base sauvegardée
- [ ] Documentation mise à jour

---

## Scénarios de migration courants

### Scénario 1: Prototype vers Production
**Chemin:** SQLite → PostgreSQL

Cas d'usage: Vous avez construit un MVP avec SQLite, maintenant passage à la production.

**Étapes:**
1. Exporter les données SQLite
2. Régénérer avec PostgreSQL
3. Configurer PostgreSQL dans Docker/cloud
4. Importer les données
5. Tester minutieusement

---

### Scénario 2: Hébergement partagé vers Cloud
**Chemin:** MySQL → PostgreSQL

Cas d'usage: Passage d'un hébergement partagé vers infrastructure cloud.

**Étapes:**
1. Exporter les données MySQL
2. Régénérer avec PostgreSQL
3. Configurer PostgreSQL managé (AWS RDS, GCP Cloud SQL)
4. Importer les données
5. Mettre à jour les chaînes de connexion

---

### Scénario 3: Simplifier le développement
**Chemin:** PostgreSQL → SQLite

Cas d'usage: Développement local plus rapide sans Docker.

**Étapes:**
1. Exporter un jeu de données de test minimal
2. Régénérer avec SQLite
3. Importer les données de test
4. Garder la production sur PostgreSQL

---

## Besoin d'aide?

- Consulter le [Guide des bases de données](./databases.md) pour des informations détaillées
- Ouvrir une issue sur [GitHub](https://github.com/tky0065/go-starter-kit/issues)
- Voir les exemples dans le répertoire `/examples` (si disponible)

---

## Ressources additionnelles

- [Outils de migration PostgreSQL](https://www.postgresql.org/docs/current/migration.html)
- [Guide de migration MySQL](https://dev.mysql.com/doc/refman/8.0/en/migration.html)
- [Import/Export SQLite](https://www.sqlite.org/cli.html#csv_import)
- [Migration GORM](https://gorm.io/docs/migration.html)

---

**Dernière mise à jour:** 2026-02-09
**Stories liées:** Epic 7 - Support multi-bases de données (Story 7.5)
