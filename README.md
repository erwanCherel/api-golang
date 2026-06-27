# myAPI — API REST de gestion de vidéos

API REST de gestion de vidéos avec système de commentaires et d'authentification par token. Réalisée en **Go** dans le cadre du module TIC-API3 (Bachelor ETNA, 2024).

> Une réécriture en **PHP/Symfony** est prévue — voir la section Roadmap.

## Stack

- **Go 1.23** — langage principal
- **Fiber v2** — framework HTTP
- **GORM** — ORM MySQL
- **JWT** — authentification par token
- **MySQL / MariaDB** — base de données

## Endpoints

### Authentification
| Méthode | URI | Auth | Description |
|---|---|---|---|
| POST | `/auth` | - | Connexion, retourne un token |

### Utilisateurs
| Méthode | URI | Auth | Description |
|---|---|---|---|
| POST | `/user` | - | Créer un compte |
| GET | `/users` | - | Lister les utilisateurs (pagination, filtre pseudo) |
| GET | `/user/:id` | Token | Récupérer un utilisateur |
| PUT | `/user/:id` | Token* | Modifier un utilisateur |
| DELETE | `/user/:id` | Token* | Supprimer un utilisateur |

### Vidéos
| Méthode | URI | Auth | Description |
|---|---|---|---|
| POST | `/user/:id/video` | Token* | Uploader une vidéo |
| GET | `/videos` | - | Lister les vidéos (pagination, filtres) |
| GET | `/user/:id/videos` | - | Vidéos d'un utilisateur |
| PUT | `/video/:id` | Token* | Modifier une vidéo |
| PATCH | `/video/:id` | Token | Ajouter un format encodé |
| DELETE | `/video/:id` | Token* | Supprimer une vidéo (+ fichier) |

### Commentaires
| Méthode | URI | Auth | Description |
|---|---|---|---|
| POST | `/video/:id/comment` | Token* | Poster un commentaire |
| GET | `/video/:id/comments` | Token* | Lister les commentaires |

`Token*` = le token doit appartenir au propriétaire de la ressource.

## Lancement

### Prérequis

- Go 1.23+
- MySQL ou MariaDB

### Installation

```bash
git clone https://github.com/erwanCherel/myApi.git
cd myApi

cp .env.example .env
# Renseigner les variables dans .env

go run main.go
```

L'API écoute sur `http://localhost:5001`.

### Variables d'environnement

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=myapi
JWT_SECRET=your_secret
```

La base de données est créée automatiquement au démarrage si elle n'existe pas.

## Structure

```
.
├── main.go
├── config/         # Chargement des variables d'environnement
├── database/       # Connexion MySQL + AutoMigrate
├── models/         # User, Video, VideoFormat, Token, Comment
├── middleware/     # Authentification JWT
├── handlers/
│   ├── auth/       # POST /auth
│   ├── user/       # Handlers utilisateurs et upload vidéo
│   └── video/      # Handlers vidéos et commentaires
├── routes/         # Déclaration des routes par domaine
├── router/         # Point d'entrée du routing
└── public/videos/  # Fichiers vidéo uploadés
```

## Collection Postman

Une collection Postman est disponible : `myApi.postman_collection.json`.

## Roadmap — réécriture Symfony

Points à couvrir lors de la réécriture en PHP/Symfony :

- [ ] **Entités Doctrine** — User, Video, VideoFormat, Token, Comment avec leurs relations
- [ ] **Authentification** — token généré à la connexion, stocké en BDD, transmis via header `Authorization`
- [ ] **Middleware d'auth** — Event Listener ou `#[IsGranted]` sur les routes protégées, sans préfixe `/private/` dans les URLs
- [ ] **Upload de fichier** — gestion des conflits de nom (timestamp + nom original), stockage dans `public/videos/`
- [ ] **Suppression en cascade** — supprimer le fichier physique à la suppression d'une vidéo
- [ ] **Pagination** — réponse systématique avec `pager.current` et `pager.total` (nombre de pages, pas d'items)
- [ ] **Format de réponse uniforme** — toujours `{"message": "...", "data": ...}` avec les bons codes HTTP (201 sur création, 204 sans body sur suppression)
- [ ] **Validation** — contraintes sur `username` (`[a-zA-Z0-9_-]`), email, etc. via les annotations Symfony Validator
- [ ] **Séparation des handlers** — un controller par domaine (UserController, VideoController, CommentController)

## Contexte

Projet réalisé dans le cadre du module **TIC-API3** à l'ETNA (2024). Premier projet en Go — l'objectif était la maîtrise des concepts REST : routing, authentification par token, upload de fichiers, pagination.
