package auth

import (
	"fmt"
	"net/http"
	"os"
	"soceng-ai/internals/utils/prompts"
	"strings"
	"time"
)

/* TODO :
Test : tenter de s'inscrire avec un formulaire contenant une exploitation (ex : script XSS, injection SQL, etc.)
Test : tenter de s'inscrire avec un formulaire json contenant des champs non nécessaires
Test : tenter de s'inscrire avec un formulaire json ne contenant pas les champs nécessaires
Test : tenter de s'inscrire sans username
Test : tenter de s'inscrire sans email
Test : tenter de s'inscrire sans mot de passe
Test : tenter de s'inscrire avec un email ne respectant pas le format d'un email
Test : tenter de s'inscrire avec un email comportant des caractères spéciaux non autorisés
Test : tenter de s'inscrire avec un email déjà utilisé
Test : tenter de s'inscrire un mot de passe trop court
Test : tenter de s'inscrire un mot de passe ne comportant pas de chiffres
Test : tenter de s'inscrire un mot de passe ne comportant pas de lettres
Test : tenter de s'inscrire un mot de passe ne comportant pas de caractères spéciaux
Test : tenter de s'inscrire un mot de passe comportant des caractères spéciaux non autorisés
Test : tenter de s'inscrire avec un username comportant des caractères spéciaux non autorisés
Test : tenter de s'inscrire avec un username déjà utilisé
Test : tenter de s'inscrire avec un username trop court
Test : tenter de s'inscrire avec un username trop long
Test : inscription réussie (username, mot de passe, email)

// Fonctionnalités encore non implémentées
Test : vérification de la création de l'utilisateur dans la base de données TODO entrée d'API pour demander l'existence d'un utilisateur
Test : vérification par mail TODO création d'une fonctionnalité d'envoi de mail
*/

var register_route = "/create-user"
var good_username = "loic.glanois"
var good_password = "Password0!"
var good_email = "loic.glanois@ynov.com"

func Tests(host, port string) {
	if !test_bad_method(host, port) || !test_form_no_json(host, port) || !test_no_json_form(host, port) || !test_no_data_sent(host, port) || !test_json_fucked_up1(host, port) || !test_json_fucked_up2(host, port) {
		os.Exit(1)
	}
	if !test_exploit_xss(host, port) || !test_exploit_sql(host, port) {
		os.Exit(1)
	}

	if !test_ok_registering(host, port, "celine", "louis-ferdinand@destouches.fr", "Password0!") || !test_ok_registering(host, port, "sartre", "jean-paul@sartre.fr", "Password2!") {
		os.Exit(1)
	}
}

func test_ok_registering(host, port, username, email, password string) bool {
	var payload = strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email))
	resp, err := http.Post("http://"+host+":"+port+register_route, "application/json", payload)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_ok_registering")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_ok_registering")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête d'inscription réussie")
	prompts.Prompts_tests(time.Now(), prompts.Success+"L'utilisateur a été créé dans la base de données")
	return true
}

func test_exploit_sql(host, port string) bool {
	var payload = strings.NewReader(`{"username": "select * from person where name='' or 1=1;", "password": "Password0!", "email": "loic.glanois@ynov.com"}`)
	resp, err := http.Post("http://"+host+":"+port+register_route, "application/json", payload)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_exploit_sql")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_exploit_sql")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête avec une exploitation SQL pour l'inscription")
	return true
}

func test_exploit_xss(host, port string) bool {
	var payload = strings.NewReader(`{"username": "<script>alert('XSS')</script>", "password": "Password0!", "email": "loic.glanois@ynov.com"}`)
	resp, err := http.Post("http://"+host+":"+port+register_route, "application/json", payload)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_exploit_xss")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_exploit_xss")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête avec une exploitation XSS pour l'inscription")
	return true
}

// func test_exploit_sql(host, port string) bool {
// 	var payload = strings.NewReader(`{"username": "loic.glanois", "password": "Password0!", "email": "

func test_json_fucked_up2(host, port string) bool {
	var payload = strings.NewReader(`{username": "loic.glanois", "password": "Password0!", "email": "loic.glanois@ynov.com"}`)
	resp, err := http.Post("http://"+host+":"+port+register_route, "application/json", payload)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_json_fucked_up")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_json_fucked_up")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête avec un JSON mal formé pour l'inscription")
	return true
}

func test_json_fucked_up1(host, port string) bool {
	var payload = strings.NewReader(`{ "username": "loic.glanois", "password": "Password0!", "email": "loic.glanois@ynov.com"`)
	resp, err := http.Post("http://"+host+":"+port+register_route, "application/json", payload)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_json_fucked_up")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_json_fucked_up")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête avec un JSON mal formé pour l'inscription")
	return true
}

func test_no_data_sent(host, port string) bool {
	resp, err := http.Post("http://"+host+":"+port+register_route, "text/plain", nil)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_no_data_sent")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_no_data_sent")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête sans données pour l'inscription")
	return true
}

func test_no_json_form(host, port string) bool {
	resp, err := http.Post("http://"+host+":"+port+register_route, "application/json", nil)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_no_json_form")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_no_json_form")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}
	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête sans formulaire JSON pour l'inscription")
	return true
}

func test_form_no_json(host, port string) bool {
	resp, err := http.Post("http://"+host+":"+port+register_route, "application/x-www-form-urlencoded", nil)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_form_no_json")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP POST : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_form_no_json")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}

	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la requête sans JSON pour l'inscription")
	return true
}

func test_bad_method(host, port string) bool {
	resp, err := http.Get("http://" + host + ":" + port + register_route)
	if err != nil {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_bad_method")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur lors de la requête HTTP GET : "+err.Error())
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		prompts.Prompts_tests(time.Now(), prompts.Error+"test_bad_method")
		prompts.Prompts_tests(time.Now(), prompts.Error+"Erreur : le serveur a répondu avec le code de statut "+fmt.Sprint(resp.StatusCode))
		return false
	}

	prompts.Prompts_tests(time.Now(), prompts.Success+"Test de la méthode HTTP incorrecte pour l'inscription")
	return true
}

/* DONE
Test : tenter de s'inscrire avec une méthode HTTP différente de POST
Test : tenter de s'inscrire avec un formulaire qui n'est pas au format JSON
Test : tenter de s'inscrire sans formulaire json
Test : pas de données envoyées du tout
Test : tenter de s'inscrire avec un formulaire json mal formé
*/
