package database

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/utils"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	return db
}

func InitDB() (*gorm.DB, error) {

	config := utils.LoadEnv()

	host := config.DBHost
	port := config.DBPort
	user := config.DB_User
	password := config.DBPassword
	dbname := config.DBName

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Enable the uuid-ossp extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.User{}, &model.OTP{}, &model.UserSession{},&model.Plan{})

	fmt.Println("Connected to the database")
	return db, nil
}
