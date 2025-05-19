package db_sessions

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	database "soceng-ai/database"
	"soceng-ai/database/tables/db_challenges"
	"soceng-ai/database/tables/db_challenges/db_challenges_structs"
	"soceng-ai/database/tables/db_sessions/db_sessions_structs"
	"soceng-ai/database/tables/db_users"
	"soceng-ai/internals/utils/prompts"
	"time"
)

func Create_game_session(username string, challenge_title string) (string, int) {
	var error_status string
	var challenge_id int
	id := get_next_game_session_available_id()
	user_id := db_users.Get_user_id_by_username_or_email(username)
	challenge_id, error_status = db_challenges.Get_challenge_id_by_title(challenge_title)
	if error_status != "OK" {
		return error_status, http.StatusNoContent
	}
	session_key := generate_session_key()
	start_time := get_current_time()
	status := "in_progress"

	error_status = db_challenges.Is_challenge_validated(challenge_id)
	if error_status != "OK" {
		return error_status, http.StatusUnauthorized
	}

	error_status = delete_previous_game_session_data(challenge_id, user_id)
	if error_status != "OK" {
		return "Aborting operation due to error : " + error_status, http.StatusInternalServerError
	}

	// Crée une nouvelle session de jeu dans la base de données
	return_value := create_game_session(id, user_id, challenge_id, session_key, start_time, status)
	if return_value != "OK" {
		return return_value, http.StatusInternalServerError
	}

	characters := db_challenges.Get_characters_by_challenge_id(challenge_id)
	for _, character := range characters {
		create_session_character(id, character.ID, character.Initial_suspicion, character.Is_accessible)
	}
	hints, error_status := db_challenges.Get_hints_by_challenge_id(challenge_id)
	if hints == nil || error_status != "OK" {
		return error_status, http.StatusNoContent
	}
	for _, hint := range hints {
		create_session_hint(id, hint.ID, hint.Is_available_from_start)
	}

	return "OK", http.StatusOK
}

func Get_sessions_by_username(username string) (string, []db_sessions_structs.Session) {
	user_id := db_users.Get_user_id_by_username_or_email(username)
	if user_id == 0 {
		return "Error getting user id", nil
	}
	db := database.Get_DB()
	rows, err := db.Query("SELECT id, user_id, challenge_id, session_key, start_time, status FROM game_sessions WHERE user_id = $1", user_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_sessions_by_username():Error getting sessions by username: " + err.Error())
		return "Error getting sessions by username: " + err.Error(), nil
	}
	defer rows.Close()
	var sessions []db_sessions_structs.Session
	for rows.Next() {
		var session db_sessions_structs.Session
		err := rows.Scan(&session.ID, &session.UserID, &session.ChallengeID, &session.SessionKey, &session.StartTime, &session.Status)
		if err != nil {
			fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_sessions_by_username():Error scanning session: " + err.Error())
			return "Error scanning session: " + err.Error(), nil
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_sessions_by_username():Error iterating over rows: " + err.Error())
		return "Error iterating over rows: " + err.Error(), nil
	}
	if len(sessions) < 0 {
		return "Error : negative number of sessions", nil
	}
	return "OK", sessions
}

func delete_previous_game_session_data(challenge_id int, user_id int) string {
	var err error
	db := database.Get_DB()
	_, err = db.Exec("DELETE FROM session_characters WHERE session_id IN (SELECT id FROM game_sessions WHERE user_id = $1 AND challenge_id = $2)", user_id, challenge_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:delete_previous_game_session_data():Error deleting previous session characters: " + err.Error())
		return "Error deleting previous game session data : " + err.Error()
	}
	_, err = db.Exec("DELETE FROM session_hints WHERE session_id IN (SELECT id FROM game_sessions WHERE user_id = $1 AND challenge_id = $2)", user_id, challenge_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:delete_previous_game_session_data():Error deleting previous session hints: " + err.Error())
		return "Error deleting previous game session data : " + err.Error()
	}
	_, err = db.Exec("DELETE FROM game_sessions WHERE user_id = $1 AND challenge_id = $2", user_id, challenge_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:delete_previous_game_session_data():Error deleting previous game session data: " + err.Error())
		return "Error deleting previous game session data : " + err.Error()
	}
	return "OK"
}

func create_session_hint(session_id int, hint_id int, is_accessible bool) string {
	db := database.Get_DB()
	_, err := db.Exec("INSERT INTO session_hints (id, session_id, hint_id, is_accessible) VALUES ($1, $2, $3, $4)",
		get_next_session_hint_available_id(), session_id, hint_id, is_accessible)
	if err != nil {
		return "soceng-ai/database/tables/db_sessions/db_session.go:create_session_hint():Error creating session hint: " + err.Error()
	}
	return "OK"
}

func get_next_session_hint_available_id() int {
	var next_id int
	db := database.Get_DB()
	db.QueryRow("SELECT COALESCE(MAX(id) + 1, 1) FROM session_hints").Scan(&next_id)
	return next_id
}

func create_session_character(session_id int, character_id int, suspicion_level int, is_accessible bool) string {
	db := database.Get_DB()
	_, err := db.Exec("INSERT INTO session_characters (id, session_id, character_id, suspicion_level, is_accessible) VALUES ($1, $2, $3, $4, $5)",
		get_next_session_character_available_id(), session_id, character_id, suspicion_level, is_accessible)
	if err != nil {
		return "soceng-ai/database/tables/db_sessions/db_session.go:create_session_character():Error creating session character: " + err.Error()
	}

	return "OK"
}

func get_next_session_character_available_id() int {
	var next_id int
	db := database.Get_DB()
	db.QueryRow("SELECT COALESCE(MAX(id) + 1, 1) FROM session_characters").Scan(&next_id)
	return next_id
}

func create_game_session(id int, user_id int, challenge_id int, session_key string, start_time string, status string) string {
	db := database.Get_DB()
	_, err := db.Exec("INSERT INTO game_sessions (id, user_id, challenge_id, session_key, start_time, status) VALUES ($1, $2, $3, $4, $5, $6)",
		id, user_id, challenge_id, session_key, start_time, status)
	if err != nil {
		return "soceng-ai/database/tables/db_sessions/db_session.go:create_game_session():Error creating game session: " + err.Error()
	}
	return "OK"
}

func get_current_time() string {
	return time.Now().Format(time.Now().Format("2006-01-02 15:04:05"))
}

func generate_session_key() string {
	// Génère une clé de session unique à 6 caractères
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	for {
		key := make([]byte, keyLength)
		for i := range key {
			key[i] = letters[rand.Intn(len(letters))]
		}

		// Vérifie si la clé existe déjà
		db := database.Get_DB()
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM game_sessions WHERE session_key = $1)", string(key)).Scan(&exists)
		if err != nil || !exists {
			return string(key)
		}
		// Si la clé existe déjà, on génère une nouvelle clé
	}
}

func get_next_game_session_available_id() int {
	var next_id int
	db := database.Get_DB()
	db.QueryRow("SELECT COALESCE(MAX(id) + 1, 1) FROM game_sessions").Scan(&next_id)
	return next_id
}

func Check_session_id(session_id string) string {
	db := database.Get_DB()
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM game_sessions WHERE session_key = $1)", session_id).Scan(&exists)
	if err != nil {
		return "Error checking session ID: " + err.Error()
	}
	if !exists {
		return "Session ID does not exist"
	}
	return "OK"
}

func Get_challenge_data(session_id string) (string, int, []byte) {
	var error_status string
	var status_code int
	var payload []byte
	var challenge_id int
	var challenge_data db_challenges_structs.Challenge
	db := database.Get_DB()
	err := db.QueryRow("SELECT challenge_id FROM game_sessions WHERE session_key = $1", session_id).Scan(&challenge_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_challenge_data():Error getting challenge data: " + err.Error())
		return "Error getting challenge data: " + err.Error(), http.StatusNoContent, nil
	}
	if challenge_id == 0 {
		return "Error: no challenge ID found", http.StatusNoContent, nil
	}
	challenge_data, error_status = db_challenges.Get_challenge_data(challenge_id)
	if error_status != "OK" {
		return error_status, http.StatusNoContent, nil
	}
	if payload == nil {
		return "Error: no challenge data found", http.StatusNoContent, nil
	}
	status_code = http.StatusOK
	payload, err = json.Marshal(challenge_data)
	return "OK", status_code, payload

}
