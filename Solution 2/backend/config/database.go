package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB(username, password, hostname, dbname string) {
	// Buat string koneksi dari informasi yang diberikan
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)

	// Buka koneksi ke database
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}

	DB = db

	log.Println("Database connected")
}
