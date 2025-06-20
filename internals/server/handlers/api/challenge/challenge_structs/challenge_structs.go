package challenge_structs

type Hint struct {
	Title                   string `json:"title"`
	Text                    string `json:"text"`
	Keywords                string `json:"keywords"`
	Illustration_type       string `json:"hint_illustration_type"`
	Mentions                string `json:"mentions"`
	Is_available_from_start bool   `json:"is_available_from_start"`
	Is_capital              bool   `json:"is_capital"`
}

type Characters struct {
	Advice_to_user          string `json:"advice_to_user"`
	Character_name          string `json:"character_name"`
	Title                   string `json:"title"`
	Initial_suspicion       int    `json:"initial_suspicion"`
	Communication_type      string `json:"communication_type"`
	Osint_data              string `json:"osint_data"`
	Knows_contact_of        string `json:"knows_contact_of"`
	Holds_hint              string `json:"holds_hint"`
	Is_available_from_start bool   `json:"is_available_from_start"`
}

type Challenge struct {
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Illustration    string       `json:"illustration"`
	Lore_for_player string       `json:"lore_for_player"`
	Lore_for_ai     string       `json:"lore_for_ai"`
	Osint_data      string       `json:"osint_data"`
	Hints           []Hint       `json:"hints"`
	Characters      []Characters `json:"characters"`
}
