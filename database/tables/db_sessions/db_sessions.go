package db_sessions

import (
	"fmt"
	"math/rand"
	"net/http"
	database "soceng-ai/database"
	"soceng-ai/database/tables/db_challenges"
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

/*
CREATE TABLE game_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id INT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    session_key VARCHAR(50) NOT NULL UNIQUE,
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL CHECK (status IN ('in_progress', 'completed'))
);

CREATE TABLE session_characters (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    suspicion_level INT NOT NULL CHECK (suspicion_level BETWEEN 0 AND 100),
    is_accessible BOOLEAN DEFAULT FALSE
);

CREATE TABLE session_hints (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    hint_id INT NOT NULL REFERENCES hints(id) ON DELETE CASCADE,
    is_accessible BOOLEAN DEFAULT FALSE
)
*/
