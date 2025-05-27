package basics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"soceng-ai/internals/utils/prompts"
	"time"
)

func Tests(host, port string) {
	if !test_hello_world(host, port) || !test_not_hello_world(host, port) || !test_not_found(host, port) {
		os.Exit(1)
	}
}

// Envoie une requête HTTP GET à l'URL de test pour vérifier que le serveur répond correctement.
func test_hello_world(host, port string) bool {
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

func test_not_hello_world(host, port string) bool {
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
func test_not_found(host, port string) bool {
	resp, err := http.Get("http://" + host + ":" + port + "/api/not-found")
	if err != nil {
		fmt.Println(prompts.Error+"Erreur lors de la requête HTTP GET :", err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test not found")
	return true
}
