package db_profiles

import (
	"fmt"
	database "soceng-ai/database"
	db_users "soceng-ai/database/tables/db_users"
)

type Db_profile struct {
	Username  string `json:"username"`
	Biography string `json:"biography"`
	Avatar    string `json:"avatar"`
}

func Does_profile_exist(username string) bool {
	db := database.Get_DB()
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM profiles WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func Create_profile(username string) error {
	db := database.Get_DB()

	id := get_next_available_id()
	user_id := db_users.Get_user_id_by_username_or_email(username)
	if user_id == -1 {
		return fmt.Errorf("user not found")
	}

	_, err := db.Exec("INSERT INTO profiles (id, user_id) VALUES (?, ?)", id, user_id)
	if err != nil {
		return fmt.Errorf("failed to create profile: %v", err)
	}
	return nil
}

func get_next_available_id() int {
	db := database.Get_DB()
	id := 0
	err := db.QueryRow("SELECT MAX(id) FROM profiles").Scan(&id)
	if err != nil && err.Error() != "sql: Scan error on column index 0, name \"MAX(id)\": converting NULL to int is unsupported" {
		fmt.Println("Error getting next available ID:", err)
		return -1
	}
	return id + 1
}

func Update_profile(username string, user Db_profile) error {
	db := database.Get_DB()

	_, err := db.Exec("UPDATE profiles SET biography = ?, avatar = ? WHERE user_id = (SELECT id FROM users WHERE username = ?)", user.Biography, user.Avatar, username)
	if err != nil {
		fmt.Println("Error updating profile:", err)
		return fmt.Errorf("failed to update profile: %v", err)
	}
	_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", user.Username, db_users.Get_user_id_by_username_or_email(username))
	if err != nil {
		fmt.Println("Error updating user:", err)
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}
