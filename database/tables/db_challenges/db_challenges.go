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

	// On teste la validité des blocs de données hints et characters
	// err_str := hints_treatment(challenge.Hints)
	// if err_str != "" {
	// 	return err_str
	// }
	// err_str = characters_treatment(challenge.Characters)
	// if err_str != "" {
	// 	return err_str
	// }

	id := get_next_available_id()
	present_time := get_current_time()
	validated := false
	difficulty := 1 // Default difficulty, can be changed later
	query = "INSERT INTO challenges (id, title, lore_for_player, difficulty, illustration, created_at, updated_at, validated) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, id, challenge.Title, challenge.Description, difficulty, challenge.Illustration, present_time, present_time, validated)
	if err != nil {
		fmt.Println("Error creating challenge:", err)
	}

	// TODO : créer la fonction qui vérifie la validité de l'ensemble du challenge et qui supprime le challenge et les entités associées si le challenge est invalide
	err_str = check_challenge_validity(challenge)
	if err_str != "" {
		fmt.Println("Error checking challenge validity:", err_str)
	}

	return ""
}

func hints_treatment(hints []challenge_structs.Hint) string {
	// Etape 1 : séparer
	return ""
}

func characters_treatment(characters []challenge_structs.Characters) string {
	return ""
}

// check_challenge_validity vérifie que le challenge est fonctionnel et que ses éléments sont reliés entre eux
// et qu'ils sont accessibles. Si ce n'est pas le cas, il renvoie une chaîne de caractères contenant l'erreur.
func check_challenge_validity(challenge challenge_structs.Challenge) string {
	// TODO : check the structural validity of the challenge
	// TODO Check if all characters are accessible (known from at least one other character or mentionned in the hints)
	// TODO Check if all hints are accessible (known from at least one character each)
	// TODO Check if capital hint is indirectly linked to at least one starter hint or character

	return ""
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

// type Hint struct {
// 	Title                   string `json:"title"`
// 	Text                    string `json:"text"`
// 	Keywords                string `json:"keywords"`
// 	Illustration_type       string `json:"hint_illustration_type"`
// 	Mentions                int    `json:"mentions"`
// 	Is_available_from_start bool   `json:"is_available_from_start"`
// 	Is_capital              bool   `json:"is_capital"`
// }

// type Characters struct {
// 	Advice_to_user          string `json:"advice_to_user"`
// 	Symbolic_name           string `json:"symbolic_name"`
// 	Title                   string `json:"title"`
// 	Initial_suspicion       int    `json:"initial_suspicion"`
// 	Communication_type      string `json:"communication_type"`
// 	Symbolic_osint_data     string `json:"symbolic_osint_data"`
// 	Knows_contact_of        string `json:"knows_contact_of"`
// 	Holds_hint              string `json:"holds_hint"`
// 	Is_available_from_start bool   `json:"is_available_from_start"`
// }

// type Challenge struct {
// 	Title           string       `json:"title"`
// 	Description     string       `json:"description"`
// 	Illustration    string       `json:"illustration"`
// 	Lore_for_player string       `json:"lore_for_player"`
// 	Lore_for_ai     string       `json:"lore_for_ai"`
// 	Osint_data      string       `json:"osint_data"`
// 	Hints           []Hint       `json:"hints"`
// 	Characters      []Characters `json:"characters"`
// }
