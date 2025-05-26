package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"soceng-ai/database/tables/db_cookies"
	"soceng-ai/database/tables/db_sessions"
	challenge "soceng-ai/internals/server/handlers/api/challenge"
	dashboard "soceng-ai/internals/server/handlers/api/dashboard"
	sessions "soceng-ai/internals/server/handlers/api/sessions"
	authentification "soceng-ai/internals/server/handlers/authentification"
	"soceng-ai/internals/utils/prompts"
)

var re = regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)

func Challenge_handler(w http.ResponseWriter, r *http.Request) {

	cookies_status := process_cookies(r)
	if cookies_status != "OK" {
		http.Error(w, "Error processing cookies", http.StatusUnauthorized)
		return
	}

	switch r.Method {

	// Créer le challenge
	case "POST":
		challenge.Create(w, r)

	// Récupérer le challenge
	// case "GET":
	// 	challenge.Read(w, r)

	// Valider le challenge
	case "PUT":
		challenge.Update(w, r)

	// Supprimer le challenge
	// case "DELETE":
	// 	challenge.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// Récupère les requêtes commençant par /api/sessions.
//
// Résumé :
//
// - POST : /api/sessions/start-challenge `'{"challenge_name": "challenge_name"}'` -> Crée une nouvelle session de jeu
//
// - GET : /api/sessions/{session_id} -> Récupère les données de session en JSON
//
// - POST : /api/sessions/{session_id} `'{"character_name": "character_name", "message": "message"}'` -> Envoie un message au personnage
//
// Concerne la collecte des données de jeu et l'envoi de
// messages aux personnages. Contrôle par cookies. Si la méthode est GET, on envoie les données de session. Si la
// méthode est POST, elle doit contenir sa clé de session dans l'URL. On vérifie que la clé de session est valide.
// Si c'est le cas, on récupère les données de session contenant le nom du personnage adressé et le message envoyé.
// Gère également la création de session de jeu pour un challenge donné.
func Sessions_handler(w http.ResponseWriter, r *http.Request) {
	var error_status string
	var response = http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "Bad Request",
	}

	cookies_status := process_cookies(r)
	if cookies_status != "OK" {
		http.Error(w, "Error processing cookies: "+cookies_status, http.StatusUnauthorized)
		return
	}

	if r.URL.Path == "/api/sessions/start-challenge" {
		switch r.Method {
		case "POST":
			response = sessions.Start_challenge(r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		if re.MatchString(r.URL.Path[len("/api/sessions/"):]) {
			session_id := r.URL.Path[len("/api/sessions/"):]
			// On vérifie que la session existe
			error_status, response.StatusCode = db_sessions.Check_session_key(r, session_id)
			if error_status != "OK" {
				fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Sessions_handler():Error checking session ID: " + error_status)
				http.Error(w, "Error checking session ID: "+error_status, http.StatusBadRequest)
				return
			}
			switch r.Method {
			case "GET":
				response = sessions.Get_session_data(r, session_id)
			case "POST":
				response = sessions.Post_session_data(r, session_id)
			default:

				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Invalid session ID", http.StatusBadRequest)
			return
		}
	}

	response.Write(w)
}

func process_cookies(r *http.Request) string {
	cookies := r.Cookies()
	if len(r.Cookies()) < 2 {
		return "Missing cookie"
	}
	cookies_status := authentification.Cookies_relevant(cookies)
	if cookies_status != "OK" {
		return cookies_status
	}

	username_cookie, err := r.Cookie("socengai-username")
	if err != nil {
		return "Error getting username cookie : " + err.Error()
	}
	auth_cookie, err := r.Cookie("socengai-auth")
	if err != nil {
		return "Error getting auth cookie : " + err.Error()
	}

	if username_cookie.Value == "" || auth_cookie.Value == "" {
		return "Cookie empty"
	}

	if !db_cookies.Is_cookie_valid(username_cookie.Value, auth_cookie.Value) {
		return "Invalid cookie"
	}
	return "OK"
}

// Renvoie les données à afficher sur le tableau de bord
func Dashboard_handler(w http.ResponseWriter, r *http.Request) {

	cookies_status := process_cookies(r)
	if cookies_status != "OK" {
		http.Error(w, "Error processing cookies: "+cookies_status, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		error_status, data := dashboard.Get_dashboard(r)
		if error_status != "OK" {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/dashboard/dashboard.go:Dashboard_handler():Error getting dashboard data: " + error_status)
			http.Error(w, "Error getting dashboard data: "+error_status, http.StatusInternalServerError)
			return
		}
		responseData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/dashboard/dashboard.go:Dashboard_handler():Error marshalling dashboard data: " + err.Error())
			http.Error(w, "Error marshalling dashboard data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseData)
		if err != nil {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/dashboard/dashboard.go:Dashboard_handler():Error writing response: " + err.Error())
			http.Error(w, "Error writing response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func HelloWorld_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, World !")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func NotHelloWorld_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Not Hello, World !")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
