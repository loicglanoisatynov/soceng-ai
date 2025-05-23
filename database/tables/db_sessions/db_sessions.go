package db_sessions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	database "soceng-ai/database"
	db_challenges "soceng-ai/database/tables/db_challenges"
	"soceng-ai/database/tables/db_sessions/db_sessions_structs"
	"soceng-ai/database/tables/db_users"
	"soceng-ai/internals/utils/prompts"
	"time"
)

func Create_game_session(username string, challenge_title string) (string, int, string) {
	var error_status string
	var challenge_id int
	id := get_next_game_session_available_id()
	user_id := db_users.Get_user_id_by_username_or_email(username)
	challenge_id, error_status = db_challenges.Get_challenge_id_by_title(challenge_title)
	if error_status != "OK" {
		return error_status, http.StatusNoContent, ""
	}
	session_key := generate_session_key()
	start_time := get_current_time()
	status := "in_progress"

	error_status = db_challenges.Is_challenge_validated(challenge_id)
	if error_status != "OK" {
		return error_status, http.StatusUnauthorized, ""
	}

	error_status = delete_previous_game_session_data(challenge_id, user_id)
	if error_status != "OK" {
		return "Aborting operation due to error : " + error_status, http.StatusInternalServerError, ""
	}

	// Crée une nouvelle session de jeu dans la base de données
	return_value := create_game_session(id, user_id, challenge_id, session_key, start_time, status)
	if return_value != "OK" {
		return return_value, http.StatusInternalServerError, ""
	}

	characters, error_status := db_challenges.Get_characters_by_challenge_id(challenge_id)
	if characters == nil || error_status != "OK" {
		delete_previous_game_session_data(challenge_id, user_id)
		return "Error getting characters by challenge ID", http.StatusNoContent, ""
	}
	for _, character := range characters {
		create_session_character(id, character.ID, character.Initial_suspicion, character.Is_accessible)
		if error_status != "OK" {
			delete_previous_game_session_data(challenge_id, user_id)
			return error_status, http.StatusNoContent, ""
		}
	}
	hints, error_status := db_challenges.Get_hints_by_challenge_id(challenge_id)
	if hints == nil || error_status != "OK" {
		delete_previous_game_session_data(challenge_id, user_id)
		return error_status, http.StatusNoContent, ""
	}
	for _, hint := range hints {
		error_status = create_session_hint(id, hint.ID, hint.Is_available_from_start)
		if error_status != "OK" {
			delete_previous_game_session_data(challenge_id, user_id)
			return error_status, http.StatusNoContent, ""
		}
	}

	return "OK", http.StatusOK, session_key
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
	if len(sessions) == 0 {
		return "No sessions found", nil
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

func Get_session_data(session_id string) (string, int, []byte) {
	var error_status string
	var status_code int
	var payload []byte
	var challenge_id int
	var session_data db_sessions_structs.Session
	db := database.Get_DB()
	err := db.QueryRow("SELECT challenge_id FROM game_sessions WHERE session_key = $1", session_id).Scan(&challenge_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_sessions.go:Get_challenge_data():Error getting challenge data: " + err.Error())
		return "Error getting challenge data: " + err.Error(), http.StatusNoContent, nil
	}
	if challenge_id == 0 {
		return "Error: no challenge ID found", http.StatusNoContent, nil
	}
	session_data, error_status = Get_session_data_by_session_id(session_id)
	if error_status != "OK" {
		return error_status, http.StatusNoContent, nil
	}
	if session_data.ID == 0 {
		return "Error: no session data found", http.StatusNoContent, nil
	}
	payload, err = json.Marshal(session_data)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_sessions.go:Get_challenge_data():Error marshalling session data: " + err.Error())
		return "Error marshalling session data: " + err.Error(), http.StatusNoContent, nil
	}
	status_code = http.StatusOK
	return "OK", status_code, payload
}

func Get_session_data_by_session_id(session_id string) (db_sessions_structs.Session, string) {
	db := database.Get_DB()
	var session_data db_sessions_structs.Session
	err := db.QueryRow("SELECT id, user_id, challenge_id, session_key, start_time, status FROM game_sessions WHERE session_key = $1", session_id).Scan(&session_data.ID, &session_data.UserID, &session_data.ChallengeID, &session_data.SessionKey, &session_data.StartTime, &session_data.Status)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_sessions.go:Get_session_data_by_session_id():Error getting session data by session ID: " + err.Error())
		return db_sessions_structs.Session{}, "Error getting session data by session ID: " + err.Error()
	}
	session_characters, error_status := Get_session_characters_by_session_id(session_data.ID)
	if error_status != "OK" {
		return db_sessions_structs.Session{}, error_status
	}
	session_hints, error_status := Get_session_hints_by_session_id(session_data.ID)
	if error_status != "OK" {
		return db_sessions_structs.Session{}, error_status
	}
	session_data.Characters = session_characters
	session_data.Hints = session_hints
	return session_data, "OK"
}

func Get_session_characters_by_session_id(session_id int) ([]db_sessions_structs.Session_character, string) {
	db := database.Get_DB()

	query := `
	SELECT 
		sc.id, sc.session_id, sc.character_id, sc.suspicion_level, sc.is_accessible,
		c.character_name, c.title, c.advice_to_user, c.communication_type, c.osint_data
	FROM session_characters sc
	JOIN characters c ON sc.character_id = c.id
	WHERE sc.session_id = $1
	`

	rows, err := db.Query(query, session_id)
	if err != nil {
		return nil, "Error getting session characters by session ID: " + err.Error()
	}
	defer rows.Close()

	var results []db_sessions_structs.Session_character

	for rows.Next() {
		var sc db_sessions_structs.Session_character
		err := rows.Scan(
			&sc.ID, &sc.SessionID, &sc.CharacterID, &sc.Suspicion, &sc.IsAccessible,
			&sc.Name, &sc.Title, &sc.Advice_to_user, &sc.CommunicationType, &sc.OsintData,
		)
		if err != nil {
			return nil, "Error scanning session character: " + err.Error()
		}
		results = append(results, sc)
	}
	if err := rows.Err(); err != nil {
		return nil, "Error iterating over rows: " + err.Error()
	}

	return results, "OK"
}

func Get_session_hints_by_session_id(session_id int) ([]db_sessions_structs.Session_hint, string) {
	db := database.Get_DB()

	query := `
	SELECT 
		sh.id, sh.session_id, sh.hint_id, sh.is_accessible,
		h.hint_title, h.hint_text, h.illustration_type,
		COALESCE(h.mentions, 0), h.is_capital
	FROM session_hints sh
	JOIN hints h ON sh.hint_id = h.id
	WHERE sh.session_id = $1
	`

	rows, err := db.Query(query, session_id)
	if err != nil {
		return nil, "error getting session hints by session ID: " + err.Error()
	}
	defer rows.Close()

	var results []db_sessions_structs.Session_hint
	for rows.Next() {
		var sh db_sessions_structs.Session_hint
		err := rows.Scan(
			&sh.ID, &sh.SessionID, &sh.HintID, &sh.IsAvailable,
			&sh.Title, &sh.Text, &sh.IllustrationType,
			&sh.Mentions, &sh.IsCapital,
		)
		if err != nil {
			return nil, "error scanning session hint: " + err.Error()
		}
		results = append(results, sh)
	}

	if err := rows.Err(); err != nil {
		return nil, "error iterating over rows: " + err.Error()
	}

	return results, "OK"
}

func Check_session_key(r *http.Request, session_key string) (string, int) {
	username_cookie, _ := r.Cookie("socengai-username")
	// Regarde dans la base de données si le cookie username correspond à un utilisateur
	user_id := db_users.Get_user_id_by_username_or_email(username_cookie.Value)
	if user_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_session_key():Error: user ID is 0")
		return "Error: user ID is 0", http.StatusUnauthorized
	}
	// Regarde dans la base de données si une session existe avec le cookie username et la clé de session
	db := database.Get_DB()
	var session_id int
	err := db.QueryRow("SELECT id FROM game_sessions WHERE user_id = $1 AND session_key = $2", user_id, session_key).Scan(&session_id)
	if err == sql.ErrNoRows {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_session_key():Error: " + err.Error())
		return "Error: " + err.Error(), http.StatusInternalServerError
	}
	if session_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_session_key():Error: session ID is 0")
		return "Error: session ID is 0", http.StatusUnauthorized
	}

	return "OK", http.StatusOK
}

// 3. Vérifier que le character_name pointe vers un personnage existant dans le challenge à partir duquel on a créé la session dont la clé est session_key
// error_status, code_status = db_sessions.Check_character_existence(session_key, post_session_data_request.Character_name)

func Check_character_existence(session_key string, character_name string) (string, int) {
	db := database.Get_DB()
	var character_id int
	err := db.QueryRow("SELECT id FROM characters WHERE character_name = $1", character_name).Scan(&character_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_character_existence():Error getting character ID: " + err.Error())
		return "Error getting character ID: " + err.Error(), http.StatusNoContent
	}
	if character_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_character_existence():Error: character ID is 0")
		return "Error: character ID is 0", http.StatusUnauthorized
	}

	// Récupère l'ID du challenge auquel appartient le personnage
	var challenge_id int
	err = db.QueryRow("SELECT challenge_id FROM characters WHERE id = $1", character_id).Scan(&challenge_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_character_existence():Error getting challenge ID: " + err.Error())
		return "Error getting challenge ID: " + err.Error(), http.StatusNoContent
	}
	if challenge_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_character_existence():Error: challenge ID is 0")
		return "Error: challenge ID is 0", http.StatusUnauthorized
	}
	// Vérifie que la session existe pour le challenge_id
	var session_id int
	err = db.QueryRow("SELECT id FROM game_sessions WHERE session_key = $1 AND challenge_id = $2", session_key, challenge_id).Scan(&session_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_character_existence():Error getting session ID: " + err.Error())
		return "Error getting session ID: " + err.Error(), http.StatusNoContent
	}
	if session_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Check_character_existence():Error: session ID is 0")
		return "Error: session ID is 0", http.StatusUnauthorized
	}

	return "OK", http.StatusOK
}
