package api

import (
	"net/http"
	"soceng-ai/database/tables/db_cookies"
	challenge "soceng-ai/internals/server/handlers/api/challenge"
	sessions "soceng-ai/internals/server/handlers/api/sessions"
	authentification "soceng-ai/internals/server/handlers/authentification"
)

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

// Récupère les requêtes commençant par /api/sessions
func Sessions_handler(w http.ResponseWriter, r *http.Request) {
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
	}

	// On renvoie la réponse
	response.Write(w)
}

/*🔹 1. Création et gestion des sessions

POST /api/sessions/start-challenge
➤ Démarre une nouvelle session de jeu pour un challenge donné.
Payload : { challenge_id }
Opération backend : créer un objet game_session pointant vers l'objet challenge associé. Création des substituts (noms uniques générés).
Retourne : session_id, personnages initiaux, documents initiaux

GET /api/sessions/{session_id}
➤ Récupère les métadonnées et l’état courant d’une session existante.
Inclut : état (en cours, terminé), timestamp, progression éventuelle.

🔹 2. Personnages (agents simulés)

GET /api/sessions/{session_id}/characters
➤ Liste des personnages disponibles dans la session, avec leur nom unique généré.

GET /api/sessions/{session_id}/characters/{character_id}/chat
➤ Récupère l’historique des messages échangés avec ce personnage.

POST /api/sessions/{session_id}/characters/{character_id}/chat
➤ Envoie un message à un personnage, reçoit la réponse IA.
Payload : { message }

🔹 3. Documents et indices

GET /api/sessions/{session_id}/documents
➤ Récupère la liste des documents découverts par le joueur jusqu’à présent.

GET /api/sessions/{session_id}/documents/{doc_id}
➤ Récupère le contenu d’un document spécifique.

POST /api/sessions/{session_id}/documents/unlock
➤ Débloque un document par une action utilisateur (indice résolu, interaction IA, etc.)
Payload : { trigger }

🔹 4. Progression et état final

POST /api/sessions/{session_id}/submit-flag
➤ L’utilisateur pense avoir trouvé le flag, on vérifie et marque la session comme terminée.
Payload : { flag }

GET /api/sessions/{session_id}/progress
➤ (Optionnel) Donne un feedback sur la progression (nombre de documents trouvés, interactions faites…)*/

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
