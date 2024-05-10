package database

import (
	"database/sql"
)

var DBConn *sql.DB

func ConnectDB() {
	connStr := "user=admin dbname=postgres password=admin123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DBConn = db

	// log.Println(DBConn)

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
