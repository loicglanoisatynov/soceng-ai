package db_challenges

import (
	"fmt"
	"net/http"
	"soceng-ai/database"
	db_challenges_structs "soceng-ai/database/tables/db_challenges/db_challenges_structs"
	"soceng-ai/internals/server/handlers/api/challenge/challenge_structs"
	"soceng-ai/internals/server/handlers/api/dashboard/dashboard_structs"
	"soceng-ai/internals/utils/prompts"
)

func Create_challenge(challenge challenge_structs.Challenge, r *http.Request, w http.ResponseWriter) string {
	db := database.Get_DB()
	var count int
	var err_str string
	query := "SELECT COUNT(*) FROM challenges WHERE title = ?"
	err := db.QueryRow(query, challenge.Title).Scan(&count)
	if err != nil {
		return "Error checking challenge data: " + err.Error()
	}
	if count > 0 {
		fmt.Println("Challenge with this title already exists.")
		return "Challenge with this title already exists."
	}

	challenge_id := get_next_available_id()
	// On teste la validité des blocs de données hints et characters
	err_str = hints_treatment(challenge.Hints, challenge_id)
	if err_str != "" {
		return err_str
	}
	err_str = characters_treatment(challenge.Characters, challenge_id)
	if err_str != "" {
		return err_str
	}

	present_time := get_current_time()
	validated := false
	difficulty := 1 // Default difficulty, can be changed later
	query = "INSERT INTO challenges (id, title, lore_for_player, lore_for_ai, difficulty, illustration, created_at, updated_at, validated) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, challenge_id, challenge.Title, challenge.Description, challenge.Lore_for_ai, difficulty, challenge.Illustration, present_time, present_time, validated)
	if err != nil {
		return "Error inserting challenge: " + err.Error()
	}

	err_str = check_challenge_coherence(challenge)
	if err_str != "" {
		// TODO : delete the challenge and its elements from the database
		purify(challenge_id)
		return "Error checking challenge validity : " + err_str
	}

	return ""
}

func purify(challenge_id int) string {
	db := database.Get_DB()
	query := "DELETE FROM challenges WHERE id = ?"
	_, err := db.Exec(query, challenge_id)
	if err != nil {
		return "Error deleting challenge : " + err.Error()
	}
	query = "DELETE FROM characters WHERE challenge_id = ?"
	_, err = db.Exec(query, challenge_id)
	if err != nil {
		return "Error deleting characters: " + err.Error()
	}
	query = "DELETE FROM hints WHERE challenge_id = ?"
	_, err = db.Exec(query, challenge_id)
	if err != nil {
		return "Error deleting hints: " + err.Error()
	}
	return ""
}

func Get_hints_by_challenge_id(challenge_id int) ([]db_challenges_structs.Db_hint, string) {
	db := database.Get_DB()
	query := "SELECT * FROM hints WHERE challenge_id = ?"
	rows, err := db.Query(query, challenge_id)
	if err != nil {
		fmt.Println("Error getting hints:", err)
		return nil, err.Error()
	}
	defer rows.Close()

	var hints []db_challenges_structs.Db_hint
	// Parcours des résultats de la requête
	// On crée un tableau de hints
	for rows.Next() {
		var hint db_challenges_structs.Db_hint
		err := rows.Scan(&hint.ID, &hint.Challenge_id, &hint.Title, &hint.Text, &hint.Keywords, &hint.Illustration_type, &hint.Mentions, &hint.Is_available_from_start, &hint.Is_capital)
		if err != nil && err.Error() == "sql: Scan error on column index 6, name \"mentions\": converting NULL to string is unsupported" {
			// fmt.Println("Error scanning hint:", err)
			hint.Mentions = ""
			// continue
		} else if err != nil {
			fmt.Println("Error scanning hint:", err)
			return nil, err.Error()
		}
		hints = append(hints, hint)
	}
	return hints, "OK"
}

func Get_characters_by_challenge_id(challenge_id int) []db_challenges_structs.Db_character {
	db := database.Get_DB()
	query := "SELECT * FROM characters WHERE challenge_id = ?"
	rows, err := db.Query(query, challenge_id)
	if err != nil {
		fmt.Println("Error getting characters:", err)
		return nil
	}
	defer rows.Close()

	var characters []db_challenges_structs.Db_character
	// Parcours des résultats de la requête
	// On crée un tableau de personnages
	for rows.Next() {
		var character db_challenges_structs.Db_character
		// On récupère les données de chaque personnage
		// On les stocke dans la structure Characters
		err := rows.Scan(&character.ID, &character.Challenge_id, &character.Advice_to_user, &character.Symbolic_name, &character.Title, &character.Initial_suspicion, &character.Communication_type, &character.Symbolic_osint_data, &character.Knows_contact_of, &character.Holds_hint, &character.Is_available_from_start)
		if err != nil {
			fmt.Println("Error scanning character:", err)
			continue
		}
		characters = append(characters, character)
	}
	return characters
}

func Is_challenge_validated(challenge_id int) string {
	db := database.Get_DB()
	var validated bool
	query := "SELECT validated FROM challenges WHERE id = ?"
	err := db.QueryRow(query, challenge_id).Scan(&validated)
	if err != nil {
		return "Error getting challenge validation status: " + err.Error()
	}
	if !validated {
		return "Challenge not validated"
	}
	return "OK"
}

func Get_challenge_id_by_title(challenge_title string) (int, string) {
	db := database.Get_DB()
	var challenge_id int
	query := "SELECT id FROM challenges WHERE title = ?"
	err := db.QueryRow(query, challenge_title).Scan(&challenge_id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return -1, "Challenge not found"
	} else if err != nil {
		return -1, "Error getting challenge ID: " + err.Error()
	}
	if challenge_id == 0 {
		return -1, "Challenge not found"
	}
	return challenge_id, "OK"
}

func hints_treatment(hints []challenge_structs.Hint, challenge_id int) string {
	// Etape 1 : Vérifier que les attributs de chaque hint sont valides
	// Etape 2 : Injecter les hints dans la base de données

	db := database.Get_DB()
	var next_id int

	for i := 0; i < len(hints); i++ {
		if hints[i].Title == "" {
			return "Hint has no title"
		}
		if hints[i].Text == "" {
			return hints[i].Title + " text is empty"
		}
		if hints[i].Keywords == "" {
			return hints[i].Title + " keywords are empty"
		}
		if hints[i].Illustration_type == "" && hints[i].Illustration_type != "none" {
			return hints[i].Title + " illustration type is empty"
		}
		if hints[i].Is_available_from_start && !hints[i].Is_available_from_start {
			return hints[i].Title + " is_available_from_start is not a boolean"
		}

		next_id = get_next_available_hint_id()

		query := "INSERT INTO hints (id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err := db.Exec(query, next_id, challenge_id, hints[i].Title, hints[i].Text, hints[i].Keywords, hints[i].Illustration_type, hints[i].Mentions, hints[i].Is_available_from_start, hints[i].Is_capital)
		if err != nil {
			return "Error inserting hint: " + err.Error()
		}
	}
	return ""
}

func get_next_available_hint_id() int {
	db := database.Get_DB()
	id := 0
	err := db.QueryRow("SELECT MAX(id) FROM hints").Scan(&id)
	if err != nil && err.Error() != "sql: Scan error on column index 0, name \"MAX(id)\": converting NULL to int is unsupported" {
		fmt.Println("Error getting next available ID:", err)
		return -1
	}
	return id + 1
}

func characters_treatment(characters []challenge_structs.Characters, challenge_id int) string {
	db := database.Get_DB()
	var next_id int

	for i := 0; i < len(characters); i++ {
		if characters[i].Character_name == "" {
			return characters[i].Character_name + " symbolic_name is empty"
		}
		if characters[i].Title == "" {
			return characters[i].Character_name + " title is empty"
		}
		if characters[i].Advice_to_user == "" {
			return characters[i].Character_name + " advice_to_user is empty"
		}
		if characters[i].Initial_suspicion < 0 || characters[i].Initial_suspicion > 10 {
			return characters[i].Character_name + " initial_suspicion is not between 0 and 10"
		}
		if characters[i].Communication_type == "" {
			return characters[i].Character_name + " communication_type is empty (has to be 'email', 'phone', 'in-person', 'social_media')"
		}
		if characters[i].Osint_data == "" {
			return characters[i].Character_name + " symbolic_osint_data is empty"
		}
		if characters[i].Holds_hint == "" && characters[i].Holds_hint != "null" {
			return characters[i].Character_name + " holds_hint is empty"
		}
		if characters[i].Is_available_from_start && !characters[i].Is_available_from_start {
			return characters[i].Character_name + " is_available_from_start is not a boolean"
		}
		if characters[i].Knows_contact_of != "" {
			// On vérifie que le contact est valide
			contact_id := 0
			query := "SELECT id FROM characters WHERE symbolic_name = ? AND challenge_id = ?"
			err := db.QueryRow(query, characters[i].Knows_contact_of, challenge_id).Scan(&contact_id)
			if err != nil {
				return characters[i].Title + " knows_contact_of is not valid"
			}
			if contact_id == 0 {
				return characters[i].Title + " knows_contact_of is not valid"
			}
		}

		if characters[i].Holds_hint != "" {
			// On vérifie que le hint est valide
			hint_id := 0
			query := "SELECT id FROM hints WHERE hint_title = ? AND challenge_id = ?"
			err := db.QueryRow(query, characters[i].Holds_hint, challenge_id).Scan(&hint_id)
			if err != nil {
				return characters[i].Title + " holds_hint is not valid"
			}
			if hint_id == 0 {
				return characters[i].Title + " holds_hint is not valid"
			}
		}
		next_id = get_next_available_character_id()

		// Si tout est valide, on peut injecter le character dans la base de données
		query := "INSERT INTO characters (id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion, communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err := db.Exec(query, next_id, challenge_id, characters[i].Advice_to_user, characters[i].Character_name, characters[i].Title, characters[i].Initial_suspicion, characters[i].Communication_type, characters[i].Osint_data, characters[i].Knows_contact_of, characters[i].Holds_hint, characters[i].Is_available_from_start)
		if err != nil {
			return "Error inserting character: " + err.Error()
		}
	}

	return ""
}

func get_next_available_character_id() int {
	db := database.Get_DB()
	id := 0
	err := db.QueryRow("SELECT MAX(id) FROM characters").Scan(&id)
	if err != nil && err.Error() != "sql: Scan error on column index 0, name \"MAX(id)\": converting NULL to int is unsupported" {
		fmt.Println("Error getting next available ID:", err)
		return -1
	}
	return id + 1
}

// check_challenge_coherence vérifie que le challenge est fonctionnel et que ses éléments sont reliés entre eux
// et qu'ils sont accessibles. Si ce n'est pas le cas, il renvoie une chaîne de caractères contenant l'erreur.
func check_challenge_coherence(challenge challenge_structs.Challenge) string {

	// 1. Cartographie
	charByName := make(map[string]*challenge_structs.Characters)
	for i := range challenge.Characters {
		charByName[challenge.Characters[i].Character_name] = &challenge.Characters[i]
	}
	hintByTitle := make(map[string]*challenge_structs.Hint)
	for i := range challenge.Hints {
		hintByTitle[challenge.Hints[i].Title] = &challenge.Hints[i]
	}

	// 2. Accessible characters / hints = propagation BFS
	accessibleChars := make(map[string]bool)
	accessibleHints := make(map[string]bool)

	// Init: tous ceux qui sont available from start
	queueChars := []string{} // Symbolic_name
	queueHints := []string{} // Title

	for _, c := range challenge.Characters {
		if c.Is_available_from_start {
			accessibleChars[c.Character_name] = true
			queueChars = append(queueChars, c.Character_name)
		}
	}
	for _, h := range challenge.Hints {
		if h.Is_available_from_start {
			accessibleHints[h.Title] = true
			queueHints = append(queueHints, h.Title)
		}
	}

	// Propagation: hints accessibles via character, character accessibles via contact or hint
	for len(queueChars) > 0 || len(queueHints) > 0 {
		// Characters
		for len(queueChars) > 0 {
			symb := queueChars[0]
			queueChars = queueChars[1:]
			c := charByName[symb]
			// Peut posséder un hint à rendre accessible
			if c.Holds_hint != "" && !accessibleHints[c.Holds_hint] {
				if _, ok := hintByTitle[c.Holds_hint]; ok {
					accessibleHints[c.Holds_hint] = true
					queueHints = append(queueHints, c.Holds_hint)
				}
			}
			// Peut relier un autre personnage
			if c.Knows_contact_of != "" && !accessibleChars[c.Knows_contact_of] {
				if _, ok := charByName[c.Knows_contact_of]; ok {
					accessibleChars[c.Knows_contact_of] = true
					queueChars = append(queueChars, c.Knows_contact_of)
				}
			}
		}
		// Hints peuvent rendre "mentionné" accessible (si applicable dans ton modèle)
		for len(queueHints) > 0 {
			title := queueHints[0]
			queueHints = queueHints[1:]
			h := hintByTitle[title]
			// S'il cite un personnage: le rendre accessible
			if h.Mentions != "" { // Mentions = Nom symbolique dans tes structs ?
				if _, ok := charByName[h.Mentions]; ok {
					if !accessibleChars[h.Mentions] {
						accessibleChars[h.Mentions] = true
						queueChars = append(queueChars, h.Mentions)
					}
				}
			}
		}
	}

	// 3. Contrôle : tout est accessible ?
	msg := ""
	for _, c := range challenge.Characters {
		if !accessibleChars[c.Character_name] {
			msg += fmt.Sprintf("Character '%s' inaccessible; ", c.Title)
		}
	}
	for _, h := range challenge.Hints {
		if !accessibleHints[h.Title] {
			msg += fmt.Sprintf("Hint '%s' inaccessible; ", h.Title)
		}
	}
	return msg
}

func get_next_available_id() int {
	db := database.Get_DB()
	id := 0
	err := db.QueryRow("SELECT MAX(id) FROM challenges").Scan(&id)
	if err != nil && err.Error() != "sql: Scan error on column index 0, name \"MAX(id)\": converting NULL to int is unsupported" {
		fmt.Println("Error getting next available ID:", err)
		return -1
	}
	return id + 1
}

func get_current_time() string {
	db := database.Get_DB()
	var current_time string
	err := db.QueryRow("SELECT CURRENT_TIMESTAMP").Scan(&current_time)
	if err != nil {
		fmt.Println("Error getting current time:", err)
		return ""
	}
	return current_time
}

func Validate_challenge(title string) {
	db := database.Get_DB()
	query := "UPDATE challenges SET validated = TRUE WHERE title = ?"
	_, err := db.Exec(query, title)
	if err != nil {
		fmt.Println("Error validating challenge:", err)
	}
}

// Get_available_challenges récupère la liste des challenges disponibles en base de données puis transmet celles-ci dans un array de challenges formatté pour le dashboard
func Get_available_challenges(username string) []dashboard_structs.Challenge {
	db := database.Get_DB()
	var challenges []db_challenges_structs.Challenge
	query := "SELECT * FROM challenges WHERE validated = TRUE"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error getting available challenges:", err)
		return nil
	}
	defer rows.Close()

	// Parcours des résultats de la requête
	// On crée un tableau de challenges
	for rows.Next() {
		var challenge db_challenges_structs.Challenge
		err := rows.Scan(&challenge.ID, &challenge.Title, &challenge.Lore_for_player, &challenge.Lore_for_ai, &challenge.Difficulty, &challenge.Illustration, &challenge.Created_at, &challenge.Updated_at, &challenge.Validated, &challenge.Osint_data)
		if err != nil {
			fmt.Println("Error scanning challenge:", err)
			continue
		}
		challenges = append(challenges, challenge)
	}
	// On crée un tableau de challenges formaté pour le dashboard
	dashboard_challenges := []dashboard_structs.Challenge{}
	for i := 0; i < len(challenges); i++ {
		var dashboard_challenge dashboard_structs.Challenge
		dashboard_challenge.Name = challenges[i].Title
		dashboard_challenge.Description = challenges[i].Lore_for_player
		dashboard_challenge.Illustration_filename = challenges[i].Illustration
		dashboard_challenge.Status = "available"
		dashboard_challenges = append(dashboard_challenges, dashboard_challenge)
	}

	return dashboard_challenges
}

func Get_challenge_data(challenge_id int) (db_challenges_structs.Challenge, string) {
	db := database.Get_DB()
	var challenge db_challenges_structs.Challenge
	query := "SELECT * FROM challenges WHERE id = ?"
	err := db.QueryRow(query, challenge_id).Scan(&challenge.ID, &challenge.Title, &challenge.Lore_for_player, &challenge.Lore_for_ai, &challenge.Difficulty, &challenge.Illustration, &challenge.Created_at, &challenge.Updated_at, &challenge.Validated, &challenge.Osint_data)
	if err != nil {
		fmt.Println(prompts.Error + "soceng-ai/database/tables/db_challenges/db_challenges.go:Get_challenge_data():Error getting challenge data: " + err.Error())
		return db_challenges_structs.Challenge{}, "Error getting challenge data: " + err.Error()
	}
	return challenge, "OK"
}

/*
type Challenge struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Illustration_filename string `json:"illustration_filename"`
	Status                string `json:"status"`
}

CREATE TABLE challenges (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    lore_for_player TEXT NOT NULL,
    lore_for_ai TEXT NOT NULL,
    difficulty INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    illustration VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    validated BOOLEAN DEFAULT FALSE,
    osint_data TEXT
);
*/
