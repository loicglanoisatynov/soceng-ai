Lançable sur WSL ?
Lançable sur Windows ?
Lançable sur Linux ?
Lançable sur Docker ?

Vérifier que le serveur est bien lancé sur le port 80 :
```bash
curl -X GET http://localhost:80
```

Créer un user (et récupérer le cookie de session) :
```bash
curl -X POST http://localhost:80/create-user -H "Content-Type: application/json" -d '{"username": "lglanois", "password": "password0!", "email":"loic.glanois@ynov.com"}' -c cookie.txt
```

Se logger (et récupérer le cookie de session) :
```bash
curl -X POST http://localhost:80/login -H "Content-Type: application/json" -d '{"username": "lglanois", "password": "very_solid_password"}' -c cookie.txt
```

Se logout :

Afficher les cookies de session récupérés :
```bash
grep "socengai" cookie.txt -A 1
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