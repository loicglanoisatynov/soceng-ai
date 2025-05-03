package db_challenges

import (
	"fmt"
	"soceng-ai/database"
)

// CREATE TABLE challenges (
//     id SERIAL PRIMARY KEY,
//     title VARCHAR(100) NOT NULL,
//     description TEXT NOT NULL,
//     difficulty INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
//     illustration VARCHAR(255),
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     validated BOOLEAN DEFAULT FALSE,
// );

func Create_challenge(title string, description string, illustration string) {
	db := database.Get_DB()
	id := get_next_available_id()
	present_time := get_current_time()
	validated := false
	difficulty := 1 // Default difficulty, can be changed later
	query := "INSERT INTO challenges (id, title, description, difficulty, illustration, created_at, updated_at, validated) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, id, title, description, difficulty, illustration, present_time, present_time, validated)
	if err != nil {
		fmt.Println("Error creating challenge:", err)
	}
}

func Is_data_valid(title string, description string, illustration string) (string, bool) {
	db := database.Get_DB()
	var count int
	query := "SELECT COUNT(*) FROM challenges WHERE title = ?"
	err := db.QueryRow(query, title).Scan(&count)
	if err != nil {
		return "Error checking challenge data: " + err.Error(), false
	}
	if count > 0 {
		fmt.Println("Challenge with this title already exists.")
		return "Challenge with this title already exists.", false
	}
	return "", true
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
