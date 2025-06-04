package db_sessions

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	database "soceng-ai/database"
	db_challenges "soceng-ai/database/tables/db_challenges"
	"soceng-ai/database/tables/db_sessions/db_sessions_structs"
	"soceng-ai/database/tables/db_users"
	"soceng-ai/internals/server/handlers/api/sessions/sessions_structs"
	"soceng-ai/internals/utils/prompts"
	"strings"
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

func Get_session_data(session_id string) (string, int, db_sessions_structs.Session) {
	var error_status string
	var status_code int
	// var payload []byte
	var challenge_id int
	var session_data db_sessions_structs.Session
	db := database.Get_DB()
	err := db.QueryRow("SELECT challenge_id FROM game_sessions WHERE session_key = $1", session_id).Scan(&challenge_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_sessions.go:Get_challenge_data():Error getting challenge data: " + err.Error())
		return "Error getting challenge data: " + err.Error(), http.StatusNoContent, db_sessions_structs.Session{}
	}
	if challenge_id == 0 {
		return "Error: no challenge ID found", http.StatusNoContent, db_sessions_structs.Session{}
	}
	session_data, error_status = Get_session_data_by_session_id(session_id)
	if error_status != "OK" {
		return error_status, http.StatusNoContent, db_sessions_structs.Session{}
	}
	if session_data.ID == 0 {
		return "Error: no session data found", http.StatusNoContent, db_sessions_structs.Session{}
	}
	// payload, err = json.Marshal(session_data)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_sessions.go:Get_challenge_data():Error marshalling session data: " + err.Error())
		return "Error marshalling session data: " + err.Error(), http.StatusNoContent, db_sessions_structs.Session{}
	}
	status_code = http.StatusOK
	return "OK", status_code, session_data
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
	session_messages, error_status := get_session_messages_by_session_id(session_data.ID)
	if error_status != "OK" {
		return db_sessions_structs.Session{}, error_status
	}
	session_data.Messages = session_messages
	session_data.Characters = session_characters
	session_data.Hints = session_hints
	return session_data, "OK"
}

/*
CREATE TABLE session_messages (
    id SERIAL PRIMARY KEY,
    session_character_id INT NOT NULL REFERENCES session_characters(id) ON DELETE CASCADE,
    sender VARCHAR(50) NOT NULL CHECK (sender IN ('user', 'character')),
    message TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    hint_given BOOLEAN DEFAULT FALSE,
    contact_given BOOLEAN DEFAULT FALSE
);
*/

func get_session_messages_by_session_id(session_id int) ([]db_sessions_structs.Session_message, string) {
	db := database.Get_DB()

	query := `
	SELECT 
		sm.id, sm.session_character_id, sm.sender, sm.message, sm.timestamp, sm.hint_given, sm.contact_given,
		sc.character_id
	FROM session_messages sm
	JOIN session_characters sc ON sm.session_character_id = sc.id
	JOIN characters c ON sc.character_id = c.id
	WHERE sc.session_id = $1
	ORDER BY sm.timestamp ASC
	`

	rows, err := db.Query(query, session_id)
	if err != nil {
		return nil, "Error getting session messages by session ID: " + err.Error()
	}
	defer rows.Close()

	var results []db_sessions_structs.Session_message

	for rows.Next() {
		var sm db_sessions_structs.Session_message
		err := rows.Scan(
			&sm.ID, &sm.SessionCharacterID, &sm.Sender, &sm.Message, &sm.Timestamp,
			&sm.HintGiven, &sm.ContactGiven, &sm.SessionCharacterID,
		)
		if err != nil {
			return nil, "Error scanning session message: " + err.Error()
		}
		results = append(results, sm)
	}
	if err := rows.Err(); err != nil {
		return nil, "Error iterating over rows: " + err.Error()
	}

	return results, "OK"
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

func Get_session_character_by_session_id(session_key string) (db_sessions_structs.Session_character, string) {
	db := database.Get_DB()
	var session_character db_sessions_structs.Session_character
	err := db.QueryRow("SELECT id, session_id, character_id, suspicion_level, is_accessible FROM session_characters WHERE session_id = (SELECT id FROM game_sessions WHERE session_key = $1)", session_key).Scan(&session_character.ID, &session_character.SessionID, &session_character.CharacterID, &session_character.Suspicion, &session_character.IsAccessible)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_session_character_by_session_id():Error getting session character by session ID: " + err.Error())
		return db_sessions_structs.Session_character{}, "Error getting session character by session ID: " + err.Error()
	}
	return session_character, "OK"
}

func Get_previous_character_message(session_key string, character_name string) string {
	db := database.Get_DB()
	var previous_message string

	err := db.QueryRow("SELECT message FROM session_messages WHERE session_character_id = (SELECT id FROM session_characters WHERE session_id = (SELECT id FROM game_sessions WHERE session_key = $1) AND character_id = (SELECT id FROM characters WHERE character_name = $2)) ORDER BY timestamp DESC LIMIT 1", session_key, character_name).Scan(&previous_message)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_previous_character_message():Error getting previous character message: " + err.Error())
		return "Error getting previous character message: " + err.Error()
	}

	return previous_message
}

func Get_session_id_from_session_key(session_key string) int {
	db := database.Get_DB()
	var session_id int
	err := db.QueryRow("SELECT id FROM game_sessions WHERE session_key = $1", session_key).Scan(&session_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_session_id_from_session_key():Error getting session ID from session key: " + err.Error())
		return 0
	}
	if session_id == 0 {
		return 0
	}
	return session_id
}

/*
CREATE TABLE session_characters (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    suspicion_level INT NOT NULL CHECK (suspicion_level BETWEEN 0 AND 100),
    is_accessible BOOLEAN DEFAULT FALSE
);

CREATE TABLE game_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id INT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    session_key VARCHAR(50) NOT NULL UNIQUE,
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL CHECK (status IN ('in_progress', 'completed'))
);

CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    challenge_id INT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    advice_to_user TEXT, -- Passable à l'API de l'IA.
    character_name VARCHAR(50) NOT NULL, -- Passable à l'IA
    title VARCHAR(50) NOT NULL, -- Passable à l'API de l'IA.
    initial_suspicion INT NOT NULL CHECK (initial_suspicion BETWEEN 1 AND 10), -- Non-passable à l'API de l'IA (sert à générer la suspicion initiale du personnage, dynamique pendant la partie). Entre 1 et 10
    communication_type VARCHAR(50) NOT NULL CHECK (communication_type IN ('email', 'phone', 'in-person', 'social_media')), -- Passable à l'API de l'IA (type de communication : email, phone, in-person, etc.)
    osint_data TEXT, -- Non-passable à l'API de l'IA (sert à générer les données osint du personnage, change pour chaque partie/session)
    knows_contact_of INT REFERENCES characters(id) ON DELETE CASCADE, -- passable à API de l'IA (passe le contact_string de la personne)
    holds_hint INT REFERENCES hints(id) ON DELETE CASCADE, -- Non-passable à l'API de l'IA (sert à générer le hint du personnage, change pour chaque partie/session)
    is_available_from_start BOOLEAN DEFAULT FALSE -- Non-passable à l'API de l'IA (sert à générer la disponibilité du personnage, change pour chaque partie/session)
);

CREATE TABLE challenges (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    lore_for_player TEXT NOT NULL,
    lore_for_ai TEXT NOT NULL,
    organisation VARCHAR(100),
    difficulty INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    illustration VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    validated BOOLEAN DEFAULT FALSE,
    osint_data TEXT
);
*/

func Get_session_character_id_by_session_id(session_id int, character_name string) int {
	db := database.Get_DB()
	var character_id int

	// session_id est l'ID de la session, session ayant pour attribut un challenge_id
	// character_name est le nom du personnage également lié à un challenge par un challenge_id
	// On récupère l'ID du personnage de session dont le nom correspond à character_name

	err := db.QueryRow("SELECT character_id FROM session_characters WHERE session_id = $1 AND character_id = (SELECT id FROM characters WHERE character_name = $2)", session_id, character_name).Scan(&character_id)
	// Si la requête échoue, on gère l'erreur

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_session_character_id_by_session_id():Error: no character found for session ID " + fmt.Sprint(session_id) + " and character name " + character_name)
			return 0
		}
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_session_character_id_by_session_id():Error getting character ID: " + err.Error())
		return 0
	}
	if character_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Get_session_character_id_by_session_id():Error: character ID is 0 for session ID " + fmt.Sprint(session_id) + " and character name " + character_name)
		return 0
	}

	return character_id
}

func get_session_id_from_session_character_id(session_character_id int) int {
	db := database.Get_DB()
	var session_id int
	err := db.QueryRow("SELECT session_id FROM session_characters WHERE id = $1", session_character_id).Scan(&session_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:get_session_id_from_session_character_id():Error getting session ID from session character ID: " + err.Error())
		return 0
	}
	if session_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:get_session_id_from_session_character_id():Error: session ID is 0 for session character ID " + fmt.Sprint(session_character_id))
		return 0
	}
	return session_id
}

func Register_messages(session_character_id int, user_message string, ai_response sessions_structs.Chall_message) (string, int, []byte, string, string, int) {
	var hint_given, contact_given string
	var status_code int
	texte_ia := extract_text_from_ai_response(ai_response.Message)
	holds_hint, holds_contact := extract_concessions_from_ai_response(ai_response.Message, get_session_id_from_session_character_id(session_character_id))
	if holds_hint {
		hint_given = get_hint_name_by_id(session_character_id)
	}
	if holds_contact {
		contact_given = get_contact_name_by_id(session_character_id)
	}
	suspicion := update_suspicion(session_character_id, ai_response.Message)

	if session_character_id == 0 {
		return "Error: character ID is 0", http.StatusNoContent, nil, "", "", 0
	}

	// On insère le message de l'utilisateur, puis la réponse de l'IA dans la base de données
	db := database.Get_DB()
	// Message de l'utilisateur
	_, err := db.Exec("INSERT INTO session_messages (id, session_character_id, sender, message, timestamp) VALUES ($1, $2, $3, $4, $5)",
		get_next_session_message_available_id(), session_character_id, "user", user_message, time.Now())
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Register_messages():Error registering messages: " + err.Error())
		return "Error registering messages: " + err.Error(), http.StatusNoContent, nil, "", "", 0
	}

	// Réponse de l'IA
	_, err = db.Exec("INSERT INTO session_messages (id, session_character_id, sender, message, timestamp, hint_given, contact_given) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		get_next_session_message_available_id(), session_character_id, "character", texte_ia, time.Now(), holds_hint, holds_contact)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:Register_messages():Error registering AI response: " + err.Error())
		return "Error registering AI response: " + err.Error(), http.StatusNoContent, nil, "", "", 0
	}

	status_code = http.StatusOK
	return "OK", status_code, []byte(texte_ia), hint_given, contact_given, suspicion
}

func get_next_session_message_available_id() int {
	var next_id int
	db := database.Get_DB()
	db.QueryRow("SELECT COALESCE(MAX(id) + 1, 1) FROM session_messages").Scan(&next_id)
	return next_id
}

func get_hint_name_by_id(session_character_id int) string {
	db := database.Get_DB()
	var hint_name string
	err := db.QueryRow("SELECT h.hint_title FROM session_hints sh JOIN hints h ON sh.hint_id = h.id WHERE sh.session_id = (SELECT session_id FROM session_characters WHERE id = $1)", session_character_id).Scan(&hint_name)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:get_hint_name_by_id():Error getting hint name: " + err.Error())
		return ""
	}
	return hint_name
}

func get_contact_name_by_id(session_character_id int) string {
	db := database.Get_DB()
	var contact_name string
	err := db.QueryRow("SELECT c.character_name FROM session_characters sc JOIN characters c ON sc.character_id = c.id WHERE sc.id = $1", session_character_id).Scan(&contact_name)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:get_contact_name_by_id():Error getting contact name: " + err.Error())
		return ""
	}
	return contact_name
}

func extract_text_from_ai_response(ai_response string) string {
	key := "\"Réponse (dialogue libre)\":"
	start := strings.Index(ai_response, key)
	if start == -1 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:extract_text_from_ai_response():Error: key not found in AI response")
		return ""
	}
	// Avancer après la clé
	start += len(key)
	// Sauter les espaces/blancs
	rest := ai_response[start:]
	rest = strings.TrimLeft(rest, " \t\n\r")
	// Le champ doit maintenant commencer par un guillemet
	if !strings.HasPrefix(rest, "\"") {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:extract_text_from_ai_response():Error: expected a quote at the start of the field")
		return ""
	}
	// Récupérer juste après le premier guillemet
	rest = rest[1:]
	end := strings.Index(rest, "\"")
	if end == -1 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:extract_text_from_ai_response():Error: closing quote not found in AI response")
		return ""
	}
	ai_response = rest[:end]

	return ai_response
}

func update_suspicion(session_character_id int, ai_response string) int {
	key := "\"Suspicion (entre 1 et 10)\":"
	start := strings.Index(ai_response, key)
	if start == -1 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_suspicion():Error: key not found in AI response")
		return 0
	}
	// Avancer après la clé
	start += len(key)
	// Sauter les espaces/blancs
	rest := ai_response[start:]
	rest = strings.TrimLeft(rest, " \t\n\r")
	// Le champ doit maintenant commencer par un chiffre
	if !strings.HasPrefix(rest, "\"") {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_suspicion():Error: expected a quote at the start of the field")
		return 0
	}
	// Récupérer juste après le premier guillemet
	rest = rest[1:]
	end := strings.Index(rest, "\"")
	if end == -1 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_suspicion():Error: closing quote not found in AI response")
		return 0
	}
	suspicion_level_str := rest[:end]
	suspicion_level := 0
	fmt.Sscanf(suspicion_level_str, "%d", &suspicion_level)

	db := database.Get_DB()
	_, err := db.Exec("UPDATE session_characters SET suspicion_level = $1 WHERE id = $2", suspicion_level, session_character_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_suspicion():Error updating suspicion level: " + err.Error())
		return 0
	}
	// fmt.Println(prompts.Info + "soceng-ai/database/tables/db_sessions/db_session.go:update_suspicion():Updated suspicion level to " + fmt.Sprint(suspicion_level) + " for session character ID " + fmt.Sprint(session_character_id))
	return suspicion_level
}

func extract_concessions_from_ai_response(ai_response string, session_id int) (bool, bool) {
	key_hint := "\"Si_convaincu_donne_document (oui ou non)\":"
	key_contact := "\"Si_convaincu_donne_contact (oui ou non)\":"
	start_hint := strings.Index(ai_response, key_hint)
	start_contact := strings.Index(ai_response, key_contact)
	// Avancer après la clé
	start_hint += len(key_hint)
	start_contact += len(key_contact)
	// Sauter les espaces/blancs
	rest_hint := ai_response[start_hint:]
	rest_contact := ai_response[start_contact:]
	rest_hint = strings.TrimLeft(rest_hint, " \t\n\r")
	rest_contact = strings.TrimLeft(rest_contact, " \t\n\r")
	// Le champ doit maintenant commencer par un guillemet
	if !strings.HasPrefix(rest_hint, "\"") || !strings.HasPrefix(rest_contact, "\"") {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:extract_concessions_from_ai_response():Error: expected a quote at the start of the field")
		return false, false
	}
	// Récupérer juste après le premier guillemet
	rest_hint = rest_hint[1:]
	rest_contact = rest_contact[1:]
	end_hint := strings.Index(rest_hint, "\"")
	end_contact := strings.Index(rest_contact, "\"")
	if end_hint == -1 || end_contact == -1 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:extract_concessions_from_ai_response():Error: closing quote not found in AI response")
		return false, false
	}
	ai_response_hint := rest_hint[:end_hint]
	ai_response_contact := rest_contact[:end_contact]
	// Convertir les réponses en booléens
	gave_hint := strings.EqualFold(ai_response_hint, "oui") || strings.EqualFold(ai_response_hint, "yes") || strings.EqualFold(ai_response_hint, "true")
	gave_contact := strings.EqualFold(ai_response_contact, "oui") || strings.EqualFold(ai_response_contact, "yes") || strings.EqualFold(ai_response_contact, "true")
	if gave_hint {
		update_hint_availability(gave_hint, ai_response_hint, session_id)
	} else {
	}
	if gave_contact {
		update_contact_availability(gave_contact, ai_response_contact, session_id)
	} else {

	}
	return gave_hint, gave_contact
}

func update_hint_availability(gave_hint bool, ai_response_hint string, session_id int) {
	db := database.Get_DB()
	var hint_id int
	err := db.QueryRow("SELECT id FROM hints WHERE hint_title = $1", ai_response_hint).Scan(&hint_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_hint_availability():Error getting hint ID: " + err.Error())
		return
	}
	if hint_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_hint_availability():Error: hint ID is 0")
		return
	}
	if gave_hint {
		create_session_hint(session_id, hint_id, true)
	} else {
		create_session_hint(session_id, hint_id, false)
	}
}

func update_contact_availability(gave_contact bool, ai_response_contact string, session_id int) {
	db := database.Get_DB()
	var character_id int
	err := db.QueryRow("SELECT id FROM characters WHERE character_name = $1", ai_response_contact).Scan(&character_id)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_contact_availability():Error getting character ID: " + err.Error())
		return
	}
	if character_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_sessions/db_session.go:update_contact_availability():Error: character ID is 0")
		return
	}
	if gave_contact {
		create_session_character(session_id, character_id, 0, true) // Suspicion level is set to 0 for contacts
	} else {
		create_session_character(session_id, character_id, 0, false)
	}
}
