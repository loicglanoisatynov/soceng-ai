package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

var db_name string = "database.db"

var db *sql.DB = nil

func Init_DB() {
	// Si la base de données n'existe pas, la créer
	if _, err := os.Stat(Get_DB_path()); os.IsNotExist(err) {
		create_DB()
		fmt.Println("Database created")
	} else {
		Set_DB(Get_DB_path())
	}

}

func Set_DB(name string) {
	var err error
	db, err = sql.Open("sqlite3", name)
	if err != nil {
		panic(err)
	}
}

func Get_DB() *sql.DB {
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func GetDBStatus() string {
	if db == nil {
		return "closed"
	} else {
		return "open"
	}
}

func Get_DB_path() string {
	return "./" + db_name
}

func Get_DB_name() string {
	return db_name
}

func create_DB() {
	var err error
	db, err = sql.Open("sqlite3", Get_DB_path())
	if err != nil {
		panic(err)
	}

	create_Tables()
}

func create_Tables() {
	schema, err := os.ReadFile("./database/schema.sql")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(schema))

	if err != nil {
		panic(err)
	}
}
