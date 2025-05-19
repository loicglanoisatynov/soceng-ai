package dashboard_structs

type Response struct {
	Dashboard Dashboard `json:"dashboard"`
}

type Dashboard struct {
	Challenges []Challenge `json:"challenges"`
	Score      int         `json:"score"`
}

type Challenge struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Illustration_filename string `json:"illustration_filename"`
	Status                string `json:"status"`
	Lore_for_player       string `json:"lore_for_player"`
	Osint_data            string `json:"osint_data"`
	Difficulty            string `json:"difficulty"`
}
