package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"soceng-ai/internals/utils/colors"
	"soceng-ai/internals/utils/prompts"
	"time"
)

var (
	port = "80"        // Port par défaut pour le serveur
	host = "localhost" // Host par défaut pour le serveur
)

// Traitement des arguments de la ligne de commande pour configurer le serveur backend.
func parse_args() {
	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-p", "--port":
			if i+1 < len(os.Args) {
				port = os.Args[i+1]
			}
		case "-h", "--host":
			if i+1 < len(os.Args) {
				host = os.Args[i+1]
			}
		}
	}

	prompts.Prompts_tests(time.Now(), prompts.Success+"Test du serveur backend sur l'hôte "+colors.Cyan_ify(host)+" et le port "+colors.Cyan_ify(port)+".")
}

// Envoie une requête HTTP GET à l'URL de test pour vérifier que le serveur répond correctement.
func test_hello_world() bool {
	resp, err := http.Get("http://" + host + ":" + port + "/api/hello-world")
	if err != nil {
		fmt.Println(prompts.Error+"Erreur lors de la requête HTTP GET :", err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la lecture du corps de la réponse : "+err.Error())
		return false
	}

	if string(body) != "Hello, World !" {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le corps de la réponse n'est pas 'Hello, World !'. Réponse reçue : "+string(body))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test hello-world")
	return true
}

func test_not_hello_world() bool {
	resp, err := http.Get("http://" + host + ":" + port + "/api/not-hello-world")
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP GET : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	// Si la réponse n'est pas "Hello, World!", on considère que le test a échoué.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la lecture du corps de la réponse : "+err.Error())
		return false
	}

	if string(body) == "Hello, World !" {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le corps de la réponse n'est pas 'Hello, World !'. Réponse reçue : "+string(body))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test not-hello-world")
	return true
}

// Vérifie que le serveur renvoie une erreur 404 pour une URL non trouvée.
func test_not_found() bool {
	resp, err := http.Get("http://" + host + ":" + port + "/api/not-found")
	if err != nil {
		fmt.Println(prompts.Error+"Erreur lors de la requête HTTP GET :", err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		// fmt.Println(prompts.Prompt_tests+prompts.Error+"Erreur : le serveur a répondu avec le code de statut", resp.StatusCode)
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test not-found")
	return true
}

// Fonction principale du processus de test du serveur backend.
func main() {
	parse_args()
	if !test_hello_world() {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Échec du test hello-world. Le serveur backend ne fonctionne pas correctement.")
		os.Exit(1)
	}

	if !test_not_hello_world() {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Échec du test not-hello-world. Le serveur backend ne fonctionne pas correctement.")
		os.Exit(1)
	}

	if !test_not_found() {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Échec du test not-found. Le serveur backend ne fonctionne pas correctement.")
		os.Exit(1)
	}

	prompts.Prompts_tests(time.Now(), prompts.Success+"Tests du serveur backend terminés avec succès !")
}
