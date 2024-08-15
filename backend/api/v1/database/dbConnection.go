package database

import (
	"email-marketing-service/api/v1/model"
	adminmodel "email-marketing-service/api/v1/model/admin"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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


	seedData(db)
}

func autoMigrateModels() {
	err := db.AutoMigrate(
		&model.User{},
		&model.OTP{},
		&model.UserSession{},
		&model.Plan{},
		&model.PlanFeature{},
		&model.APIKey{},
		&model.DailyMailCalc{},
		&model.Subscription{},
		&model.Billing{},
		&adminmodel.Admin{},
		&model.SupportTicket{},
		&model.TicketFiles{},
		&model.SentEmails{},
		&model.KnowledgeBaseArticle{},
		&model.KnowledgeBaseCategory{},
		&model.SMTPKey{},
		&model.SMTPMasterKey{},
		&model.ContactGroup{},
		&model.Contact{},
		&model.UserContactGroup{},
		&model.Template{},
	)

	if err != nil {
		log.Fatalf("Migration Failed: %v", err)
	}

}

func seedData(db *gorm.DB) {
	// Check if Plan data already exists
	var planCount int64
	db.Model(&model.Plan{}).Count(&planCount)
	if planCount == 0 {
		// Seed Plan and PlanFeature data
		plan := model.Plan{
			UUID:                uuid.New().String(),
			PlanName:            "Free",
			Duration:            "infinite",
			Price:               00,
			NumberOfMailsPerDay: "100",
			Details:             "Our best plan for power users",
			Status:              model.PlanStatus("active"),
			Features: []model.PlanFeature{
				{
					UUID:        uuid.New().String(),
					Name:        "Advanced Analytics",
					Identifier:  "adv_analytics",
					CountLimit:  100,
					SizeLimit:   1000,
					IsActive:    true,
					Description: "Get detailed insights into your email campaigns",
				},
				{
					UUID:        uuid.New().String(),
					Name:        "Custom Templates",
					Identifier:  "custom_templates",
					CountLimit:  50,
					SizeLimit:   500,
					IsActive:    true,
					Description: "Create and save your own email templates",
				},
			},
		}

		result := db.Create(&plan)
		if result.Error != nil {
			log.Printf("Failed to seed Plan data: %v", result.Error)
		} else {
			fmt.Println("Plan data seeded successfully")
		}
	} else {
		fmt.Println("Plan data already exists, skipping seed")
	}

	// Check if Admin data already exists
	var adminCount int64
	db.Model(&adminmodel.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		// Seed Admin data
		firstName := "hello"
		middleName := "wedon't really know"
		lastName := "hello"
		password := "hello123"

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			return
		}

		admin := adminmodel.Admin{
			UUID:       uuid.New().String(),
			FirstName:  &firstName,
			MiddleName: &middleName,
			LastName:   &lastName,
			Email:      "admin@admin.com",
			Password:   hashedPassword,
			Type:       "admin",
		}

		result := db.Create(&admin)
		if result.Error != nil {
			log.Printf("Failed to seed Admin data: %v", result.Error)
		} else {
			fmt.Println("Admin data seeded successfully")
		}
	} else {
		fmt.Println("Admin data already exists, skipping seed")
	}
}
