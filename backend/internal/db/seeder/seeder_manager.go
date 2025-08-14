package seeders

import (
	"context"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"log"
	"sort"
)

// Manager handles registration and execution of all seeders
type Manager struct {
	seeders []Seeder
}

// NewManager creates a new seeder manager
func NewManager() *Manager {
	return &Manager{
		seeders: make([]Seeder, 0),
	}
}

// Register adds a seeder to the manager
func (m *Manager) Register(seeder Seeder) {
	m.seeders = append(m.seeders, seeder)
}

// RegisterAll adds multiple seeders to the manager
func (m *Manager) RegisterAll(seeders ...Seeder) {
	m.seeders = append(m.seeders, seeders...)
}

// SeedAll runs all registered seeders in priority order
func (m *Manager) SeedAll(ctx context.Context, store db.Store) error {
	if len(m.seeders) == 0 {
		log.Println("No seeders registered")
		return nil
	}

	// Sort seeders by priority
	sort.Slice(m.seeders, func(i, j int) bool {
		return m.seeders[i].Priority() < m.seeders[j].Priority()
	})

	log.Printf("Starting to seed %d seeders...", len(m.seeders))

	for _, seeder := range m.seeders {
		log.Printf("Running seeder: %s (priority: %d)", seeder.Name(), seeder.Priority())

		if err := seeder.Seed(ctx, store); err != nil {
			return fmt.Errorf("seeder %s failed: %w", seeder.Name(), err)
		}

		log.Printf("Seeder %s completed successfully", seeder.Name())
	}

	log.Println("All seeders completed successfully")
	return nil
}

// GetRegisteredSeeders returns a list of all registered seeder names
func (m *Manager) GetRegisteredSeeders() []string {
	names := make([]string, len(m.seeders))
	for i, seeder := range m.seeders {
		names[i] = seeder.Name()
	}
	return names
}

func GetAllSeeders() []Seeder {
	return []Seeder{
		&PlanSeeder{},
		&SMTPSeeder{},
		// Add more seeders here as you create them
		// &UserSeeder{},
		// &CompanySeeder{},
		// &TemplateSeeder{},
	}
}
