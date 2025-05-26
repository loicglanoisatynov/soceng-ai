Lançable sur WSL ?
Lançable sur Windows ?
Lançable sur Linux ?
Lançable sur Docker ?

## Sommaire
- [Sommaire](#sommaire)
- [Opérations basiques de serveur](#opérations-basiques-de-serveur)
- [Authentification](#authentification)
  - [Créer un user (et récupérer le cookie de session) :](#créer-un-user-et-récupérer-le-cookie-de-session-)
    - [Commandes valides](#commandes-valides)
    - [Commandes invalides](#commandes-invalides)
  - [Se logger (et récupérer le cookie de session) :](#se-logger-et-récupérer-le-cookie-de-session-)
    - [Commandes valides](#commandes-valides-1)
- [Customization](#customization)
- [API Challenges](#api-challenges)
  - [Création de challenge](#création-de-challenge)
    - [Commandes valides](#commandes-valides-2)
    - [Commandes invalides](#commandes-invalides-1)
  - [Validation de challenge](#validation-de-challenge)
    - [Commandes valides](#commandes-valides-3)
  - [Enumérer les challenges dispo (dashboard) :](#enumérer-les-challenges-dispo-dashboard-)
    - [Commande valide](#commande-valide)
  - [Commencer un challenge](#commencer-un-challenge)
    - [Commande valide](#commande-valide-1)
    - [Commandes invalides](#commandes-invalides-2)
      - [Pas de cookie de session](#pas-de-cookie-de-session)
      - [Pas de payload](#pas-de-payload)
      - [Payload vide](#payload-vide)
      - [Payload non-pertinent (clés hors-sujet)](#payload-non-pertinent-clés-hors-sujet)
      - [Payload mal formé (erreur de syntaxe)](#payload-mal-formé-erreur-de-syntaxe)
      - [Challenge inexistant](#challenge-inexistant)
      - [Challenge non validé](#challenge-non-validé)
  - [Récupérer le challenge en cours](#récupérer-le-challenge-en-cours)
    - [Commande valide](#commande-valide-2)
  - [Envoyer un message dans le chat du challenge](#envoyer-un-message-dans-le-chat-du-challenge)

## Opérations basiques de serveur

Lancer le serveur :
```bash
(sudo) go run internals/main.go
```
 
Vérifier que le serveur est bien lancé sur le port 80 :
```bash
curl -X GET http://localhost:80
```

## Authentification

### Créer un user (et récupérer le cookie de session) :
#### Commandes valides
```bash
curl -X POST http://localhost:80/create-user \
-H "Content-Type: application/json" \
-d @tests/authentification/create/ok/payload.json \
-c cookie.txt
```

#### Commandes invalides
Créer un user dont le nom d'utilisateur est déjà pris :
```bash
curl -X POST http://localhost:80/create-user \
-H "Content-Type: application/json" \
-d @tests/authentification/create/usernameexists/payload.json \
-c cookie.txt
```

### Se logger (et récupérer le cookie de session) :

#### Commandes valides

Se logger (et récupérer le cookie de session) :
```bash
curl -X POST http://localhost:80/login \
-H "Content-Type: application/json" \
-d '{"username": "lglanois", "password": "very_solid_password"}' \
-c cookie.txt
```

Se logout :

Afficher les cookies de session récupérés :
```bash
grep "socengai" cookie.txt -A 1
```

## Customization

Editer ses données utilisateur (implique d'être logged in) :
```bash
curl -X PUT http://localhost:80/edit-user \
-H "Content-Type: application/json" \
-d '{"email": "newemail@gmail.com", "password": "very_solid_password", "newpassword":"even_better_password"}' \
-b cookie.txt -v
```

Editer son profil utilisateur (implique d'être logged in) :
```bash
curl -X PUT http://localhost:80/edit-profile \
-H "Content-Type: application/json" \
-d '{"username": "lglanois", "avatar": "https://example.com/avatar.png", "biography": "Ceci est ma bio"}' \
-b cookie.txt -v
```

## API Challenges

### Création de challenge

#### Commandes valides

Créer un challenge avec deux hints et un personnage : 
```bash
curl -X POST http://localhost:80/api/challenge \
  -H "Content-Type: application/json" \
  -d @tests/challenge/create-challenge-ok.json \
  -b cookie.txt -v
```

#### Commandes invalides

Créer un challenge (sans cookie de session) :
```bash
curl -X POST http://localhost:80/api/challenge \
  -H "Content-Type: application/json" \
  -d '{"challenge": {"title": "Titre du challenge", "description": "Ceci est une description test"}}' \
  -v
```

Créer un challenge avec une erreur de json (erreur de syntaxe) :
Créer un challenge sans titre :
Créer un challenge sans description :
Créer un challenge sans illustration :
Créer un challenge sans lore de joueur :
Créer un challenge sans lore d'ia : 
Créer un challenge sans hint capital (flag) :
Créer un challenge sans personnage :
Créer un challenge sans hint ou personnage accessible en début de jeu :

### Validation de challenge

#### Commandes valides

Valider un challenge (implique d'être admin) : 
```bash
curl -X PUT http://localhost:80/api/challenge -H "Content-Type: application/json" -d '{"operation":"validate", "title": "Welcome to the Game", "description": "Un petit challenge introductif", "illustration": "illustration.png"}' -b cookie.txt -v
```

### Enumérer les challenges dispo (dashboard) :
#### Commande valide
```bash
curl -X GET http://localhost:80/api/dashboard \
-H "Content-Type: application/json" \
-b cookie.txt -v \
&& echo
```

### Commencer un challenge

#### Commande valide

Commencer un challenge (implique d'être logged in, voir ##) : 
```bash
curl -X POST http://localhost:80/api/sessions/start-challenge \
  -H "Content-Type: application/json" \
  -d @tests/session/create/ok/payload.json \
  -b cookie.txt -v \
  && echo
```

#### Commandes invalides

##### Pas de cookie de session

Commande : 
```bash
curl -X POST http://localhost:80/api/sessions/start-challenge \
  -H "Content-Type: application/json" \
  -d @tests/session/create/ok/payload.json -v \
  && echo
```

Output attendu :
TODO

##### Pas de payload
TODO

##### Payload vide
TODO

##### Payload non-pertinent (clés hors-sujet)
TODO

##### Payload mal formé (erreur de syntaxe)
TODO

##### Challenge inexistant

Commande :
```bash
curl -X POST http://localhost:80/api/sessions/start-challenge \
  -H "Content-Type: application/json" \
  -d @tests/session/create/nonexistant/payload.json \
  -b cookie.txt -v \
  && echo
```

Output attendu :
TODO

##### Challenge non validé

Commande :
```bash
curl -X POST http://localhost:80/api/sessions/start-challenge \
  -H "Content-Type: application/json" \
  -d @tests/session/create/notvalidated/payload.json \
  -b cookie.txt -v \
  && echo
```

Output attendu :

### Récupérer le challenge en cours

#### Commande valide
```bash
curl -X GET http://localhost:80/api/sessions/TEST01 \
-H "Content-Type: application/json" \
-b cookie.txt -v \
&& echo
```

### Envoyer un message dans le chat du challenge
```bash
curl -X POST http://localhost:80/api/sessions/TEST01 \
-H "Content-Type: application/json" \
-d '{"character_name":"julie_recpt", "message": "Bonjour"}' \
-b cookie.txt -v \   
&& echo
```

