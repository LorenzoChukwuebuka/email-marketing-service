package seeders

import (
	"context"
	"database/sql"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"log"
	"os"
)

// AdminSeeder implements the Seeder interface for admin users
type AdminSeeder struct{}

// Name returns the seeder name
func (a *AdminSeeder) Name() string {
	return "Admins"
}

// Priority returns the execution priority
func (a *AdminSeeder) Priority() int {
	return 0 // Highest priority - admins should be created first
}

// Seed populates the database with initial admin data
func (a *AdminSeeder) Seed(ctx context.Context, queries db.Store) error {
	// Define default admins
	defaultAdmins := []struct {
		firstName  string
		middleName string
		lastName   string
		email      string
		password   string
		adminType  string
		envEmail   string
		envPass    string
	}{
		{
			firstName:  "Lawrence",
			middleName: "Chukwuebuka",
			lastName:   "Obi",
			email:      a.getEnvOrDefault("ADMIN_EMAIL", "admin@example.com"),
			password:   a.getEnvOrDefault("ADMIN_PASSWORD", "defaultpassword123"),
			adminType:  "admin",
			envEmail:   "ADMIN_EMAIL",
			envPass:    "ADMIN_PASSWORD",
		},
		{
			firstName:  "Super",
			middleName: "Admin",
			lastName:   "User",
			email:      a.getEnvOrDefault("SUPER_ADMIN_EMAIL", "superadmin@example.com"),
			password:   a.getEnvOrDefault("SUPER_ADMIN_PASSWORD", "superadminpass123"),
			adminType:  "super_admin",
			envEmail:   "SUPER_ADMIN_EMAIL",
			envPass:    "SUPER_ADMIN_PASSWORD",
		},
	}

	// Check if any admin exists first (optimization)
	admins, err := queries.GetAllAdmins(ctx)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for existing admins: %w", err)
	}

	// If admins exist, check each one individually
	if len(admins) > 0 {
		log.Printf("Found %d existing admins, checking individual admin emails", len(admins))
	}

	for _, admin := range defaultAdmins {
		// Warn about default credentials
		if admin.email == "admin@example.com" || admin.email == "superadmin@example.com" {
			log.Printf("Warning: %s not set, using default email: %s", admin.envEmail, admin.email)
		}
		if admin.password == "defaultpassword123" || admin.password == "superadminpass123" {
			log.Printf("Warning: %s not set, using default password", admin.envPass)
		}

		// Check if this specific admin exists
		_, err := queries.GetAdminByEmail(ctx, admin.email)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error checking for existing admin %s: %w", admin.email, err)
		}

		if err != sql.ErrNoRows {
			log.Printf("Admin with email %s already exists, skipping", admin.email)
			continue
		}

		// Hash the password
		hashedPassword, err := common.HashPassword(admin.password)
		if err != nil {
			return fmt.Errorf("failed to hash password for admin %s: %w", admin.email, err)
		}

		// Create the admin
		_, err = queries.CreateAdmin(ctx, db.CreateAdminParams{
			Firstname:  sql.NullString{String: admin.firstName, Valid: true},
			Middlename: sql.NullString{String: admin.middleName, Valid: true},
			Lastname:   sql.NullString{String: admin.lastName, Valid: true},
			Email:      admin.email,
			Password:   hashedPassword,
			Type:       admin.adminType,
		})
		if err != nil {
			return fmt.Errorf("error creating admin %s: %w", admin.email, err)
		}

		log.Printf("Created admin: %s (%s)", admin.email, admin.adminType)
	}

	return nil
}


func (a *AdminSeeder) getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}