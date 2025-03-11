package database

import (
	"database/sql"
	"os"
)

var db_name string = "database.db"
var db *sql.DB = nil

func Init_DB() {
	// Si la base de données n'existe pas, la créer
	if _, err := os.Stat(Get_DB_path()); os.IsNotExist(err) {
		create_DB()
	} else {
		Set_DB(Get_DB_path())
	}

	// Vérifier si la base de données est ouverte
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

	// Créer les tables
	create_Tables()

}

func create_Tables() {
	// Créer les tables
	schema, err := os.ReadFile("schema.sql")

	if err != nil {
		panic(err)
	}

	// Génère la base de données à partir du schéma SQL
	_, err = db.Exec(string(schema))

	if err != nil {
		panic(err)
	}
}
