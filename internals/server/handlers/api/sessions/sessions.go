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
