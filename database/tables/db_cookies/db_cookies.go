package db_cookies

import (
	"fmt"
	"soceng-ai/database"
	"time"
)

func Register_cookie(user_id int, cookie string) error {
	db := database.Get_DB()

	err := delete_previous_cookies(user_id)

	cookie_id := get_next_available_id()
	date_timestamp := time.Now()
	date := date_timestamp.Format("2006-01-02 15:04:05.000000")

	query := "INSERT INTO cookies (id, user_id, cookie_value, created_at, last_access) VALUES (?, ?, ?, ?, ?)"

	_, err = db.Exec(query, cookie_id, user_id, cookie, date, date)
	if err != nil {
		return err
	}

	return err
}

func delete_previous_cookies(user_id int) error {
	db := database.Get_DB()
	query := "DELETE FROM cookies WHERE user_id = ?"
	_, err := db.Exec(query, user_id)
	if err != nil {
		return err
	}
	return nil
}

func get_next_available_id() int {
	db := database.Get_DB()
	id := 0
	err := db.QueryRow("SELECT MAX(id) FROM cookies").Scan(&id)
	if err != nil && err.Error() != "sql: Scan error on column index 0, name \"MAX(id)\": converting NULL to int is unsupported" {
		fmt.Println("Error getting next available ID:", err)
		return -1
	}
	return id + 1
}

func Get_user_id_by_cookie(username string, cookie string) int {
	db := database.Get_DB()
	var user_id int
	query := "SELECT user_id FROM cookies WHERE cookie_value = ? AND user_id = (SELECT id FROM users WHERE username = ?)"
	err := db.QueryRow(query, cookie, username).Scan(&user_id)
	if err != nil {
		fmt.Println("Error getting user ID by cookie:", err)
		return -1
	}
	return user_id
}

func Delete_cookie(user_id int, cookie string) error {
	db := database.Get_DB()
	query := "DELETE FROM cookies WHERE user_id = ? AND cookie_value = ?"
	_, err := db.Exec(query, user_id, cookie)
	if err != nil {
		return err
	}
	return nil
}

func Is_cookie_valid(username string, cookie string) bool {
	db := database.Get_DB()
	var user_id int
	query := "SELECT id FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&user_id)
	if err != nil {
		fmt.Println("Error getting user ID by username:", err)
		return false
	}
	query = "SELECT COUNT(*) FROM cookies WHERE cookie_value = ? AND user_id = ?"
	var count int
	err = db.QueryRow(query, cookie, user_id).Scan(&count)
	if err != nil {
		fmt.Println("Error checking cookie validity:", err)
		return false
	}
	return count > 0
}
