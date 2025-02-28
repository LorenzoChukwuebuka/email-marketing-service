// migrations/migrations.go
package migrations

import (
	"gorm.io/gorm"
	"time"
	"email-marketing-service/api/v1/model"
)

// Migration represents a single database migration
type Migration struct {
	Name string
	Run  func(*gorm.DB) error
}

// Add all your migrations here
var migrations = []Migration{
	{
		Name: "modify_user_google_id",
		Run: func(db *gorm.DB) error {
			// Drop the unique index if it exists
			if err := db.Exec(`DROP INDEX IF EXISTS idx_users_google_id`).Error; err != nil {
				return err
			}

			// Modify the google_id column
			return db.Exec(`ALTER TABLE users ALTER COLUMN google_id DROP NOT NULL`).Error
		},
	},
	// Add more migrations as needed
}

// RunMigrations executes pending migrations
func RunMigrations(db *gorm.DB) error {
	// Create migration history table
	err := db.AutoMigrate(&model.MigrationHistory{})
	if err != nil {
		return err
	}

	// Run each migration
	for _, migration := range migrations {
		var history model.MigrationHistory

		// Check if migration was already applied
		result := db.Where("name = ?", migration.Name).First(&history)
		if result.Error == gorm.ErrRecordNotFound {
			// Run the migration in a transaction
			err := db.Transaction(func(tx *gorm.DB) error {
				if err := migration.Run(tx); err != nil {
					return err
				}

				// Record the migration
				return tx.Create(&model.MigrationHistory{
					Name:      migration.Name,
					AppliedAt: time.Now(),
				}).Error
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
