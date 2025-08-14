package seeders

import (
	"context"
	db "email-marketing-service/internal/db/sqlc"
)

// Seeder defines the interface that all seeders must implement
type Seeder interface {
	// Name returns the name of the seeder for logging purposes
	Name() string

	// Seed executes the seeding logic
	Seed(ctx context.Context, store db.Store) error

	// Priority returns the execution priority (lower number = higher priority)
	// This allows you to control the order of seeder execution
	Priority() int
}
