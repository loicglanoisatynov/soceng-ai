package db_challenges

import (
	"fmt"
	"net/http"
	"soceng-ai/database"
	challenge_structs "soceng-ai/internals/server/handlers/api/challenge/challenge_structs"
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
		if characters[i].Symbolic_name == "" {
			return characters[i].Symbolic_name + " symbolic_name is empty"
		}
		if characters[i].Title == "" {
			return characters[i].Symbolic_name + " title is empty"
		}
		if characters[i].Advice_to_user == "" {
			return characters[i].Symbolic_name + " advice_to_user is empty"
		}
		if characters[i].Initial_suspicion < 0 || characters[i].Initial_suspicion > 10 {
			return characters[i].Symbolic_name + " initial_suspicion is not between 0 and 10"
		}
		if characters[i].Communication_type == "" {
			return characters[i].Symbolic_name + " communication_type is empty (has to be 'email', 'phone', 'in-person', 'social_media')"
		}
		if characters[i].Symbolic_osint_data == "" {
			return characters[i].Symbolic_name + " symbolic_osint_data is empty"
		}
		if characters[i].Holds_hint == "" && characters[i].Holds_hint != "null" {
			return characters[i].Symbolic_name + " holds_hint is empty"
		}
		if characters[i].Is_available_from_start && !characters[i].Is_available_from_start {
			return characters[i].Symbolic_name + " is_available_from_start is not a boolean"
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
		_, err := db.Exec(query, next_id, challenge_id, characters[i].Advice_to_user, characters[i].Symbolic_name, characters[i].Title, characters[i].Initial_suspicion, characters[i].Communication_type, characters[i].Symbolic_osint_data, characters[i].Knows_contact_of, characters[i].Holds_hint, characters[i].Is_available_from_start)
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
		charByName[challenge.Characters[i].Symbolic_name] = &challenge.Characters[i]
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
			accessibleChars[c.Symbolic_name] = true
			queueChars = append(queueChars, c.Symbolic_name)
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
		if !accessibleChars[c.Symbolic_name] {
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
