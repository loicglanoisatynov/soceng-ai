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
}
