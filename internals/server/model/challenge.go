package model

type Challenge struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Flag        string `db:"flag" json:"flag"`
	Difficulty  string `db:"difficulty" json:"difficulty"`
}
