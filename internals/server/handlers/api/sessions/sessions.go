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
		returned_status, status_code = db_sessions.Create_game_session(username, create_session_request.Challenge_name)
		if returned_status != "OK" {
			payload = []byte(`{"message": "Error creating game session : ` + returned_status + `"}`)
			// status_code = http.StatusNoContent
		} else {
			payload = []byte(`{"message": "Game session created successfully"}`)
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
		payload = []byte(`{"message": "Session ID cannot be empty"}`)
		status_code = http.StatusBadRequest
	} else {
		returned_status, status_code, payload = db_sessions.Get_challenge_data(session_id)
		if returned_status != "OK" {
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
