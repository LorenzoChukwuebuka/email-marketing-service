package database

import (
	"database/sql"
	"email-marketing-service/api/utils"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)

//var db *sql.DB

func InitDB() (*sql.DB, error) {

	utils.LoadEnv()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to the database")
	return db, nil
}
