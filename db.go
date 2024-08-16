package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error

	errENV := godotenv.Load(".env")
	if errENV != nil {
		log.Fatal("Error with .env loading", errENV)
	}
	db_connection_string := os.Getenv("DB_CONNECTION_STRING")

	db, err = sql.Open("postgres", db_connection_string)

	if err != nil {
		log.Fatal("Error Connecting to DB start", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error Connecting to DB end", err)
	}

	fmt.Println("db connected succes")
}

func getDB() *sql.DB {
	return db
}
