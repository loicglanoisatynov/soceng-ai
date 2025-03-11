pour build l'image : `docker build -t <nom_de_l'image> .`
<!-- docker build -t soceng-ai-server . -->
pour lancer l'image : `docker run -it -d <nom_de_l'image>`
<!-- docker run -it -d soceng-ai-server -->
pour rentrer dans l'image en cours de fonctionnement : `docker exec -it 916d12bdb3cd /bin/bash`

lançable aussi comme objet golang indépendant : `go run bin/main.go start`