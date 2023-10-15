package database

import (
	"database/sql"
	"email-marketing-service/api/utils"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	 
)

var db *sql.DB

func GetDb() *sql.DB {
	return db
}

func InitDB() (*sql.DB, error) {

config :=	utils.LoadEnv()

	host := config.DBHost
	port := config.DBPort
	user := config.DB_User
	password := config.DBPassword
	dbname := config.DBName

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db.SetMaxOpenConns(10) // Set maximum number of open connections
	db.SetMaxIdleConns(5)  // Set maximum number of idle connections

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to the database")
	return db, nil
}



