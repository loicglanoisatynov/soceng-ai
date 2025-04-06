package cookies

import (
	"fmt"
	"soceng-ai/database"
	"time"
)

func Register_cookie(user_id int, cookie string) error {
	db := database.Get_DB()

	cookie_id := get_next_available_id()
	date_timestamp := time.Now()
	date := date_timestamp.Format("2006-01-02 15:04:05.000000")

	query := "INSERT INTO cookies (id, user_id, cookie_value, created_at, last_access) VALUES (?, ?, ?, ?, ?)"

	_, err := db.Exec(query, cookie_id, user_id, cookie, date, date)
	if err != nil {
		return err
	}

	return err
}

func get_next_available_id() int {
	db := database.Get_DB()
	var id int
	err := db.QueryRow("SELECT MAX(id) FROM cookies").Scan(&id)
	if err != nil {
		fmt.Println("Error getting next available ID:", err)
		return 1
	}
	return id + 1
}

func Get_user_id_by_cookie(cookie string) int {
	db := database.Get_DB()
	var user_id int
	query := "SELECT user_id FROM cookies WHERE cookie_value = ?"
	err := db.QueryRow(query, cookie).Scan(&user_id)
	if err != nil {
		fmt.Println("Error getting user ID by cookie:", err)
		return -1
	}
	return user_id
}
