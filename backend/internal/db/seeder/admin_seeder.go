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

// SeedAdmins populates the database with initial admin data
func SeedAdmins(ctx context.Context, queries db.Store) error {
	// Get admin credentials from environment variables
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// Set default values if environment variables are not set
	if adminEmail == "" {
		adminEmail = "admin@example.com" // Default email
		log.Println("Warning: ADMIN_EMAIL not set, using default email")
	}

	if adminPassword == "" {
		adminPassword = "defaultpassword123" // Default password
		log.Println("Warning: ADMIN_PASSWORD not set, using default password")
	}

	// Check if any admin exists
	admins, err := queries.GetAllAdmins(ctx) // Assuming you have this method
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for existing admins: %w", err)
	}

	// If admins exist, skip seeding
	if len(admins) > 0 {
		log.Println("Admins already exist, skipping admin seeding")
		return nil
	}

	// Alternatively, check by email if GetAllAdmins doesn't exist
	_, err = queries.GetAdminByEmail(ctx, adminEmail)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for existing admin: %w", err)
	}

	if err != sql.ErrNoRows {
		log.Printf("Admin with email %s already exists, skipping seeding", adminEmail)
		return nil
	}

	// Hash the password
	hashedPassword, err := common.HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	// Create the admin
	firstName := "Lawrence"
	middleName := "Chukwuebuka"
	lastName := "Obi"
	adminType := "admin"

	_, err = queries.CreateAdmin(ctx, db.CreateAdminParams{
		Firstname:  sql.NullString{String: firstName, Valid: true},
		Middlename: sql.NullString{String: middleName, Valid: true},
		Lastname:   sql.NullString{String: lastName, Valid: true},
		Email:      adminEmail,
		Password:   hashedPassword,
		Type:       adminType,
	})
	if err != nil {
		return fmt.Errorf("error creating admin: %w", err)
	}

	log.Printf("Created admin with email: %s", adminEmail)
	log.Println("Admin seeding completed successfully")
	return nil
}

// Alternative implementation if you want to create multiple default admins
func SeedMultipleAdmins(ctx context.Context, queries db.Store) error {
	// Define default admins
	defaultAdmins := []struct {
		firstName  string
		middleName string
		lastName   string
		email      string
		password   string
		adminType  string
	}{
		{
			firstName:  "hello",
			middleName: "wedon't really know",
			lastName:   "hello",
			email:      getEnvOrDefault("ADMIN_EMAIL", "admin@example.com"),
			password:   getEnvOrDefault("ADMIN_PASSWORD", "defaultpassword123"),
			adminType:  "admin",
		},
		{
			firstName:  "Super",
			middleName: "Admin",
			lastName:   "User",
			email:      getEnvOrDefault("SUPER_ADMIN_EMAIL", "superadmin@example.com"),
			password:   getEnvOrDefault("SUPER_ADMIN_PASSWORD", "superadminpass123"),
			adminType:  "super_admin",
		},
	}

	for _, admin := range defaultAdmins {
		// Check if admin exists
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

	log.Println("Admin seeding completed successfully")
	return nil
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Usage example in your main seeder file
func RunAllSeeders(ctx context.Context, queries db.Store) error {
	// Seed admins first
	if err := SeedAdmins(ctx, queries); err != nil {
		return fmt.Errorf("failed to seed admins: %w", err)
	}

	// Seed plans
	if err := SeedPlans(ctx, queries); err != nil {
		return fmt.Errorf("failed to seed plans: %w", err)
	}

	// Add other seeders here...

	log.Println("All seeders completed successfully")
	return nil
}
