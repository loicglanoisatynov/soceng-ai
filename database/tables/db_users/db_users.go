package db_users

import (
	"database/sql"
	"soceng-ai/database"
	"time"
)

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
}

func Create_user(db *sql.DB, user User) error {
	query := "INSERT INTO users (id, username, email, passwd, created_at) VALUES ($1, $2, $3, $4, $5)"
	time := time.Now().Format("2006-01-02 15:04:05")
	id, _ := Get_next_id(db)
	_, err := db.Exec(query, id, user.Username, user.Email, user.Password, time)
	if err != nil {
		return err
	}

	return nil
}

func Delete_user(db *sql.DB, id int) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// db = la base de données où fouiller. "by" = la colonne à fouiller (username, email, id). "value" = la valeur à chercher
func Get_user(db *sql.DB, by string, value string) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE " + by + " = $1"
	row := db.QueryRow(query, value)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Created_at, &user.Created_at)
	if err != nil {
		return user, err
	}

	return user, nil
}

func Update_user(db *sql.DB, user User) error {
	query := "UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4"

	_, err := db.Exec(query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func Get_next_id(db *sql.DB) (int, error) {
	var id int
	query := "SELECT MAX(id) FROM users"

	row := db.QueryRow(query)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id + 1, nil
}

func Get_user_id_by_username_or_email(identifier string) int {
	var id int
	db := database.Get_DB()
	query := "SELECT id FROM users WHERE username = $1 OR email = $2"

	row := db.QueryRow(query, identifier, identifier)
	err := row.Scan(&id)
	if err != nil {
		return 0
	}

	return id
}

func Is_password_valid(user_id int, password string) bool {
	var db_password string
	db := database.Get_DB()
	query := "SELECT passwd FROM users WHERE id = $1"

	row := db.QueryRow(query, user_id)
	err := row.Scan(&db_password)
	if err != nil {
		return false
	}

	if db_password == password {
		return true
	} else {
		return false
	}
}

func Update_email(user_id int, email string) bool {
	db := database.Get_DB()
	query := "UPDATE users SET email = $1 WHERE id = $2"

	_, err := db.Exec(query, email, user_id)
	if err != nil {
		return false
	}

	return true
}

func Update_password(user_id int, new_password string) {
	db := database.Get_DB()
	query := "UPDATE users SET passwd = $1 WHERE id = $2"

	_, err := db.Exec(query, new_password, user_id)
	if err != nil {
		return
	}
}

func Does_username_exist(username string) bool {
	db := database.Get_DB()
	query := "SELECT COUNT(*) FROM users WHERE username = $1"

	var count int
	row := db.QueryRow(query, username)
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}

func Does_email_exists(email string) bool {
	db := database.Get_DB()
	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	row := db.QueryRow(query, email)
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}

func Is_admin(user_id int) bool {
	db := database.Get_DB()
	query := "SELECT COUNT(*) FROM users WHERE id = $1 AND is_admin = true"

	var count int
	row := db.QueryRow(query, user_id)
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}
