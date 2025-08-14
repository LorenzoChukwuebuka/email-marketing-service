package seeders

import (
	"context"
	db "email-marketing-service/internal/db/sqlc"
	"log"

	"github.com/google/uuid"
)

// SMTPSeeder implements the Seeder interface for SMTP keys
type SMTPSeeder struct{}

// Name returns the seeder name
func (s *SMTPSeeder) Name() string {
	return "SMTP Keys"
}

// Priority returns the execution priority
func (s *SMTPSeeder) Priority() int {
	return 2 // Lower priority than plans
}

// Seed populates the database with SMTP key data
func (s *SMTPSeeder) Seed(ctx context.Context, q db.Store) error {
	// Start a transaction to maintain consistency
	tx, err := q.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)

	// Check if SMTP master key already exists
	exists, err := qtx.CheckSMTPMasterKeyExists(ctx, db.CheckSMTPMasterKeyExistsParams{
		SmtpLogin: "adminuser",
		Password:  "abe6f7f6-f407-4ac2-9763-5de47a0ffee4",
	})
	if err != nil {
		return err
	}

	// Skip if key already exists
	if exists {
		log.Println("SMTP master key already exists, skipping seeding")
		return tx.Commit()
	}

	// Disable foreign key constraints temporarily
	_, err = tx.ExecContext(ctx, "SET session_replication_role = 'replica'")
	if err != nil {
		return err
	}

	// Create default SMTP master key with UUID for user_id and company_id
	_, err = qtx.CreateSMTPMasterKey(ctx, db.CreateSMTPMasterKeyParams{
		UserID:    uuid.New(),
		CompanyID: uuid.New(),
		SmtpLogin: "adminuser",
		KeyName:   "adminuser",
		Password:  "abe6f7f6-f407-4ac2-9763-5de47a0ffee4",
		Status:    "active",
	})
	if err != nil {
		return err
	}

	// Re-enable foreign key constraints
	_, err = tx.ExecContext(ctx, "SET session_replication_role = 'origin'")
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("Default SMTP master key seeded successfully")
	return nil
}
