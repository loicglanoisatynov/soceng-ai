package server

import (
	"fmt"
	"net/http"
	"regexp"
	handlers "soceng-ai/internals/server/handlers"

	api "soceng-ai/internals/server/handlers/api"
	authentification "soceng-ai/internals/server/handlers/authentification"
	profiles_handling "soceng-ai/internals/server/handlers/profiles_handling"
	registering "soceng-ai/internals/server/handlers/registering"
)

var routes []Route

func routes_index_handler(w http.ResponseWriter, r *http.Request) {
	for _, route := range routes {
		fmt.Fprintf(w, "%s\n", route.Get_route_regex())
	}
}

func init() {
	routes = []Route{
		newRoute("/routes", routes_index_handler),
		newRoute("/", handlers.Home),
		newRoute("/create-user", registering.Register_user),
		newRoute("/check-register", registering.Check_register),
		newRoute("/delete-user", registering.Delete_user),
		newRoute("/login", authentification.Login),
		newRoute("/logout", authentification.Logout),
		newRoute("/edit-profile", profiles_handling.Edit_profile),
		newRoute("/edit-user", profiles_handling.Edit_user),

		// newRoute("GET", "/api/get-challenges", handlers.Get_challenges), // Récupère la liste des défis (notamment pour le front-end)
		newRoute("/api/challenge", api.Challenge_handler),
		// newRoute("GET", "/api/get-challenge", handlers.Get_challenge),
		// newRoute("PUT", "/api/edit-challenge", handlers.Edit_challenge),
		// newRoute("DELETE", "/api/delete-challenge", handlers.Delete_challenge),

		/* Gestion des sessions de parties */
		newRoute("/api/sessions/([^/]+)", api.Sessions_handler), // Créer une session de jeu (challenge_id)
		newRoute("/api/dashboard", api.Dashboard_handler),       // Récupérer les informations de session de challenge
		//get dashboard-data

		// Récupérer les informations de session de challenge
		// Récupérer les informations de conversation (affichant à la fois les données du personnage et des messages échangés)
		// Envoyer une réponse à personnage (sous-objet du challenge)

		// newRoute("GET", "/contact", contact),
		// newRoute("GET", "/([^/]+)/admin", widgetAdmin),
		// newRoute("POST", "/([^/]+)/image", widgetImage),
	}
}

// newRoute("GET", "/([^/]+)/admin", widgetAdmin),
// newRoute("POST", "/([^/]+)/image", widgetImage),

func Set_routes(new_routes []Route) {
	routes = new_routes
}

func Get_routes() []Route {
	return routes
}

func newRoute(pattern string, handler http.HandlerFunc) Route {
	return Route{regexp.MustCompile("^" + pattern + "$"), handler}
}

type Route struct {
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

// func (r Route) Get_route_method() string {
// 	return r.method
// }

func (r Route) Get_route_regex() *regexp.Regexp {
	return r.regex
}

func (r Route) Get_route_handler() http.HandlerFunc {
	return r.handler
}

/*
🔹 1. Création et gestion des sessions

POST /api/sessions/start-challenge
➤ Démarre une nouvelle session de jeu pour un challenge donné.
Payload : { challenge_id }
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
➤ (Optionnel) Donne un feedback sur la progression (nombre de documents trouvés, interactions faites…)
*/
