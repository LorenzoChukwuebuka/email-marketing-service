package database

import (
	"email-marketing-service/api/v1/model"
	adminmodel "email-marketing-service/api/v1/model/admin"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	db     *gorm.DB
	once   sync.Once
	config = utils.LoadEnv()
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
	log.Println("Connected to the database")

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
		&model.MailUsage{},
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
		&model.MailingLimit{},
		&model.Campaign{},
		&model.CampaignGroup{},
		&model.EmailCampaignResult{},
		&model.UserTempEmail{},
		&model.Log{},
		&model.Domains{},
		&model.Sender{},
		&model.WebAuthnCredential{},
		&model.EmailBox{},
		&model.UserNotification{},
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
		// Create the Plan first
		plan := model.Plan{
			UUID:     uuid.New().String(),
			PlanName: "Free",
			Duration: "day",
			Price:    0,
			Details:  "Our best plan for power users",
			Status:   model.PlanStatus(model.StatusActive),
			IsPaid:   false,
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
				{
					UUID:        uuid.New().String(),
					Name:        "Contacts",
					Identifier:  "contacts",
					CountLimit:  50,
					SizeLimit:   500,
					IsActive:    true,
					Description: "Create and save your own email templates",
				},
			},
		}

		// Create the Plan
		if err := db.Create(&plan).Error; err != nil {
			log.Printf("Failed to seed Plan data: %v", err)
			return
		}

		// Now create the MailingLimit with the correct PlanID
		mailingLimit := model.MailingLimit{
			PlanID:      plan.ID,
			LimitAmount: 200,
			LimitPeriod: "day",
		}

		if err := db.Create(&mailingLimit).Error; err != nil {
			log.Printf("Failed to seed MailingLimit data: %v", err)
			return
		}

		log.Println("Plan and MailingLimit data seeded successfully")
	} else {
		log.Println("Plan data already exists, skipping seed")
	}

	// Check if Admin data already exists
	var adminCount int64
	db.Model(&adminmodel.Admin{}).Count(&adminCount)
	if adminCount == 0 {

		firstName := "hello"
		middleName := "wedon't really know"
		lastName := "hello"
		password := config.ADMIN_PASSWORD

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
			Email:      config.ADMIN_EMAIL,
			Password:   string(hashedPassword),
			Type:       "admin",
		}

		result := db.Create(&admin)
		if result.Error != nil {
			log.Printf("Failed to seed Admin data: %v", result.Error)
		} else {
			log.Println("Admin data seeded successfully")
		}
	} else {
		log.Println("Admin data already exists, skipping seed")
	}

	// Check if SMTPMasterKey data already exists
	var smtpMasterKeyCount int64
	db.Model(&model.SMTPMasterKey{}).Count(&smtpMasterKeyCount)

	if smtpMasterKeyCount == 0 {
		// Fetch an existing user to associate with the SMTP key
		var user adminmodel.Admin
		if err := db.First(&user).Error; err != nil {
			log.Printf("Failed to fetch a user for SMTPMasterKey: %v", err)
			return
		}

		// Now create the SMTPMasterKey with the correct UserId
		smtpKey := model.SMTPMasterKey{
			UUID:      uuid.New().String(),
			UserId:    uuid.New().String(), // Use the existing user's ID
			SMTPLogin: "adminuser",
			KeyName:   "adminuser",
			Password:  uuid.New().String(),
			Status:    model.KeyActive,
		}

		if err := db.Create(&smtpKey).Error; err != nil {
			log.Printf("Failed to seed SMTPMasterKey data: %v", err)
		} else {
			log.Println("SMTPMasterKey data seeded successfully")
		}
	} else {
		log.Println("SMTPMasterKey data already exists, skipping seed")
	}
}
