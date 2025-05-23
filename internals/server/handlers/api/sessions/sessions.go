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
	var error_status string
	var post_session_data_request sessions_structs.Post_session_data_request
	status_code := http.StatusNotImplemented
	payload := []byte(`{"message": "default message"}`)

	// 1. Récupérer le message du corps de la requête
	err := json.NewDecoder(r.Body).Decode(&post_session_data_request)
	// 2. Vérifier que l'objet json passé en requête comporte un sous-objet "message" et un sous-objet "character_name"
	if err != nil {
		payload = []byte(`{"message": "Error decoding JSON : ` + err.Error() + `"}`)
		status_code = http.StatusBadRequest
		fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Post_session_data():Error decoding JSON: " + err.Error())
	} else if post_session_data_request.Character_name == "" || post_session_data_request.Message == "" {
		payload = []byte(`{"message": "Character name and message cannot be empty"}`)
		status_code = http.StatusBadRequest
	} else {
		// 3. Vérifier que le character_name pointe vers un personnage existant dans le challenge à partir duquel on a créé la session dont la clé est session_key
		error_status, status_code = db_sessions.Check_character_existence(session_key, post_session_data_request.Character_name)
		if error_status != "OK" {
			payload = []byte(`{"message": "Error checking character existence : ` + error_status + `"}`)
			status_code = http.StatusNoContent
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Post_session_data():Error checking character existence: " + error_status)
		} else {
			// 4. Si le personnage existe dans le challenge d'origine de la session, on envoie le message à l'API d'IA et on récupère la réponse
		}
	}

	// 5. Si la réponse de l'IA est valide, on ajoute le message de l'utilisateur dans la DB (table session_messages) et on ajoute la réponse de l'IA dans la DB (table session_messages)
	// 6. On envoie un message de validation

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

	return http.Response{
		StatusCode: status_code,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer(payload)),
	}
}
