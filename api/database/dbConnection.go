package database

import (
	"email-marketing-service/api/model"
	adminmodel "email-marketing-service/api/model/admin"
	"email-marketing-service/api/utils"
	"fmt"
	"log"
	"sync"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	return db
}

func InitDB() (*gorm.DB, error) {
	once.Do(func() {
		initializeDatabase()
	})

	return db, nil
}

func initializeDatabase() {
	config := utils.LoadEnv()

	host, port, user, password, dbname := config.DBHost, config.DBPort, config.DB_User, config.DBPassword, config.DBName

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal(err)
	}

	// Enable the uuid-ossp extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatal(err)
	}

	autoMigrateModels()
	fmt.Println("Connected to the database")
}

func autoMigrateModels() {
	db.AutoMigrate(
		&model.User{},
		&model.OTP{},
		&model.UserSession{},
		&model.Plan{},
		&model.APIKey{},
		&model.DailyMailCalc{},
		&model.Subscription{},
		&model.Billing{},
		&model.Logger{},
		&adminmodel.Admin{},
		&model.Ticket{},
		&model.TicketMessage{},
		&model.TicketFiles{},
		&model.KnowledgeBaseArticle{},
		&model.KnowledgeBaseCategory{},
	)
}
