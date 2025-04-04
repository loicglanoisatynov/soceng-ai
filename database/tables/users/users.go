package users

import "database/sql"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Create_user(db *sql.DB, user User) error {
	query := "INSERT INTO users (id, username, email, passwd) VALUES ($1, $2, $3, $4)"
	id, _ := Get_next_id(db)
	_, err := db.Exec(query, id, user.Username, user.Email, user.Password)
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
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
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
