package sessions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	db_sessions "soceng-ai/database/tables/db_sessions"
	sessions_structs "soceng-ai/internals/server/handlers/api/sessions/sessions_structs"
	authentification "soceng-ai/internals/server/handlers/authentification"
	"soceng-ai/internals/utils/prompts"
)

// Fonction créant un objet game_session. Informations nécessaires : nom du challenge
func Start_challenge(r *http.Request) http.Response {
	var returned_status string
	var session_key string
	status_code := http.StatusBadRequest
	payload := []byte(`{"message": "default message"}`)

	var create_session_request sessions_structs.Create_session_request
	err := json.NewDecoder(r.Body).Decode(&create_session_request)
	if err != nil {
		payload = []byte(`{"message": "Error decoding JSON : ` + err.Error() + `"}`)
		status_code = http.StatusBadRequest
	} else if create_session_request.Challenge_name == "" {
		payload = []byte(`{"message": "Challenge name cannot be empty"}`)
		status_code = http.StatusBadRequest
	} else {
		username := authentification.Get_cookie_value(r, "socengai-username")
		if username == "" {
			fmt.Println("soceng-ai/internals/server/handlers/api/sessions/sessions.go:Start_challenge():Error: username not found in header")
		}
		returned_status, status_code, session_key = db_sessions.Create_game_session(username, create_session_request.Challenge_name)
		if returned_status != "OK" {
			payload = []byte(`{"message": "Error creating game session : ` + returned_status + `"}`)
			// status_code = http.StatusNoContent TODO
		} else {
			payload = []byte(`{"message": "Game session created successfully", "session_key": "` + session_key + `"}`)
			status_code = http.StatusOK
		}
	}

	return http.Response{
		StatusCode: status_code,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer(payload)),
	}
}

// Récupère les données de la partie en cours à partir de l'ID de session (chaine de 6 caractères aléatoires) et renvoie les données du challenge en JSON

func Get_session_data(r *http.Request, session_id string) http.Response {
	var returned_status string
	var error_status string
	status_code := http.StatusBadRequest
	payload := []byte(`{"message": "default message"}`)

	error_status = db_sessions.Check_session_id(session_id)
	if error_status != "OK" {
		fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Get_session_data():Error checking session ID: " + error_status)
		return http.Response{
			StatusCode: http.StatusNoContent,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"message": "Error checking session ID : ` + error_status + `"}`))),
		}
	}
	if session_id == "" {
		fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Get_session_data():Error: session ID is empty")
		payload = []byte(`{"message": "Session ID cannot be empty"}`)
		status_code = http.StatusBadRequest
	} else {
		returned_status, status_code, payload = db_sessions.Get_session_data(session_id)

		if returned_status != "OK" {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Get_session_data():Error getting challenge data: " + returned_status + " for session ID: " + session_id + " queried by user: " + authentification.Get_cookie_value(r, "socengai-username"))
			payload = []byte(`{"message": "Error getting challenge data : ` + returned_status + `"}`)
			status_code = http.StatusNoContent
		}
	}

	return http.Response{
		StatusCode: status_code,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer(payload)),
	}
}

// Fonction qui traite de la réception des messages sur la partie en cours
func Post_session_data(r *http.Request, session_key string) http.Response {
	// Pseudo-code :
	// 1. Récupérer le message du corps de la requête
	// 2. Vérifier que la requête comporte un attribut "message" et "character_name"
	// 3. Vérifier que le character_name pointe vers un personnage existant dans la session
	// Requête : "est-ce que l'objet session_character peut être relié à une session (session_id) dont la clé est session_key et est-ce que un personnage du challenge a pour nom character_name ?"

	// 4. Vérifier que le personnage existe (requête à db), récupérer l'identifiant du personnage

	// 6. Récupérer l'id du personnage de session à partir de l'id de session et du nom du personnage

	// 7. Vérifier que le personnage est disponible (is chall_character available ?)

	// 8. Vérifier qu'un message se trouve dans le corps de la requête
	//

	// var returned_status string
	// status_code := http.StatusBadRequest
	// payload := []byte(`{"message": "default message"}`)

	// var post_session_data_request sessions_structs.Post_session_data_request
	// err := json.NewDecoder(r.Body).Decode(&post_session_data_request)
	// if err != nil {
	// 	payload = []byte(`{"message": "Error decoding JSON : ` + err.Error() + `"}`)
	// 	status_code = http.StatusBadRequest
	// } else if post_session_data_request.Session_id == "" || post_session_data_request.Character_name == "" || post_session_data_request.Message == "" {
	// 	payload = []byte(`{"message": "Session ID, character name and message cannot be empty"}`)
	// 	status_code = http.StatusBadRequest
	// } else {
	// 	returned_status, status_code, payload = db_sessions.Post_session_data(session_id, post_session_data_request.Character_name, post_session_data_request.Message)
	// 	if returned_status != "OK" {
	// 		payload = []byte(`{"message": "Error posting session data : ` + returned_status + `"}`)
	// 		status_code = http.StatusNoContent
	// 	}
	// }

	// var returned_status string
	// var payload []byte
	status_code := http.StatusBadRequest
	payload := []byte(`{"message": "default message"}`)
	return http.Response{
		StatusCode: status_code,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer(payload)),
	}
}
