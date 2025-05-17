package db_challenges_structs

type Db_character struct {
	ID                      int    `json:"id"`
	Challenge_id            int    `json:"challenge_id"`
	Advice_to_user          string `json:"advice_to_user"`
	Symbolic_name           string `json:"symbolic_name"`
	Title                   string `json:"title"`
	Initial_suspicion       int    `json:"initial_suspicion"`
	Communication_type      string `json:"communication_type"`
	Symbolic_osint_data     string `json:"symbolic_osint_data"`
	Knows_contact_of        string `json:"knows_contact_of"`
	Holds_hint              string `json:"holds_hint"`
	Is_available_from_start bool   `json:"is_available_from_start"`
	Is_accessible           bool   `json:"is_accessible"`
}

type Db_hint struct {
	ID                      int    `json:"id"`
	Challenge_id            int    `json:"challenge_id"`
	Title                   string `json:"title"`
	Text                    string `json:"text"`
	Keywords                string `json:"keywords"`
	Illustration_type       string `json:"hint_illustration_type"`
	Mentions                string `json:"mentions"`
	Is_available_from_start bool   `json:"is_available_from_start"`
	Is_capital              bool   `json:"is_capital"`
}

type Challenge struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Lore_for_player string `json:"lore_for_player"`
	Lore_for_ai     string `json:"lore_for_ai"`
	Difficulty      string `json:"difficulty"`
	Illustration    string `json:"illustration"`
	Created_at      string `json:"created_at"`
	Updated_at      string `json:"updated_at"`
	Validated       bool   `json:"validated"`
	Osint_data      string `json:"osint_data"`
}
