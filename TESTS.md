Lançable sur WSL ?
Lançable sur Windows ?
Lançable sur Linux ?
Lançable sur Docker ?

Créer un user (et récupérer le cookie de session) :
```bash
curl -X POST http://localhost:80/create-user -H "Content-Type: application/json" -d '{"username": "lglanois", "password": "password0!", "email":"loic.glanois@ynov.com"}' -c cookie.txt
```

Se logger (et récupérer le cookie de session) :
```bash
curl -X POST http://localhost:80/login -H "Content-Type: application/json" -d '{"username": "lglanois", "password": "password0!"}' -c cookie.txt
```

Se logout :

Afficher les cookies de session récupérés :
```bash
grep "socengai" cookie.txt -A 1
```

Créer un challenge (commande valide) :
```bash
curl -X POST http://localhost:80/api/challenge -H "Content-Type: application/json" -d '{"title": "Welcome to the Game", "description": "Un petit challenge introductif", "illustration": "illustration.png"}' -b cookie.txt -v
```

Valider un challenge (implique d'être admin) : 
```bash
curl -X PUT http://localhost:80/api/challenge -H "Content-Type: application/json" -d '{"operation":"validate", "title": "Welcome to the Game", "description": "Un petit challenge introductif", "illustration": "illustration.png"}' -b cookie.txt -v
```