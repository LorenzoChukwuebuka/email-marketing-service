package seeders

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SeedPlans populates the database with initial plan data
func SeedPlans(ctx context.Context, queries db.Store) error {
	// Define plan UUIDs
	freePlanID := uuid.New()
	basicPlanID := uuid.New()
	proPlanID := uuid.New()
	enterprisePlanID := uuid.New()

	now := time.Now()

	// Create plans
	plans := []struct {
		id          uuid.UUID
		name        string
		description string
		price       float64
		cycle       string
	}{
		{freePlanID, "Free", "Basic email marketing capabilities for small businesses.", 0.00, "monthly"},
		{basicPlanID, "Basic", "Essential email marketing tools to grow your audience.", 19.99, "monthly"},
		{proPlanID, "Professional", "Advanced email marketing solutions for growing businesses.", 49.99, "monthly"},
		{enterprisePlanID, "Enterprise", "Complete email marketing platform for large organizations.", 99.99, "monthly"},
	}

	// Insert plans
	for _, plan := range plans {
		// Check if plan exists
		existingPlan, err := queries.GetPlanByName(ctx, plan.name)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error checking for existing plan: %w", err)
		}

		if err == sql.ErrNoRows {
			// Plan doesn't exist, create it
			_, err = queries.CreatePlan(ctx, db.CreatePlanParams{
				ID:           plan.id,
				Name:         plan.name,
				Description:  sql.NullString{String: plan.description, Valid: true},
				Price:        decimal.NewFromFloat(plan.price),
				BillingCycle: sql.NullString{String: plan.cycle, Valid: true},
				Status:       sql.NullString{String: "active", Valid: true},
				CreatedAt:    now,
				UpdatedAt:    now,
			})
			if err != nil {
				return fmt.Errorf("error creating plan %s: %w", plan.name, err)
			}
			log.Printf("Created plan: %s", plan.name)
		} else {
			// Plan exists, use its ID
			log.Printf("Plan already exists: %s", plan.name)
			plan.id = existingPlan.ID
		}
	}

	// Create features for Free plan
	freeFeatures := []struct {
		name        string
		description string
		value       string
	}{
		{"Subscribers", "Maximum number of subscribers", "500"},
		{"Monthly Emails", "Emails per month", "5000"},
		{"Email Templates", "Available email templates", "5"},
		{"Automation", "Automated email sequences", "false"},
		{"Support", "Customer support level", "Email only"},
	}

	for _, feature := range freeFeatures {
		err := createFeatureIfNotExists(ctx, queries, freePlanID, feature.name, feature.description, feature.value)
		if err != nil {
			return err
		}
	}

	// Create features for Basic plan
	basicFeatures := []struct {
		name        string
		description string
		value       string
	}{
		{"Subscribers", "Maximum number of subscribers", "2500"},
		{"Monthly Emails", "Emails per month", "25000"},
		{"Email Templates", "Available email templates", "15"},
		{"Automation", "Automated email sequences", "true"},
		{"Support", "Customer support level", "Email and chat"},
		{"A/B Testing", "Test different email variations", "false"},
	}

	for _, feature := range basicFeatures {
		err := createFeatureIfNotExists(ctx, queries, basicPlanID, feature.name, feature.description, feature.value)
		if err != nil {
			return err
		}
	}

	// Create features for Pro plan
	proFeatures := []struct {
		name        string
		description string
		value       string
	}{
		{"Subscribers", "Maximum number of subscribers", "10000"},
		{"Monthly Emails", "Emails per month", "100000"},
		{"Email Templates", "Available email templates", "50"},
		{"Automation", "Automated email sequences", "true"},
		{"Support", "Customer support level", "Priority email and chat"},
		{"A/B Testing", "Test different email variations", "true"},
		{"Advanced Analytics", "Detailed email performance metrics", "true"},
	}

	for _, feature := range proFeatures {
		err := createFeatureIfNotExists(ctx, queries, proPlanID, feature.name, feature.description, feature.value)
		if err != nil {
			return err
		}
	}

	// Create features for Enterprise plan
	enterpriseFeatures := []struct {
		name        string
		description string
		value       string
	}{
		{"Subscribers", "Maximum number of subscribers", "Unlimited"},
		{"Monthly Emails", "Emails per month", "Unlimited"},
		{"Email Templates", "Available email templates", "Unlimited"},
		{"Automation", "Automated email sequences", "true"},
		{"Support", "Customer support level", "Dedicated account manager"},
		{"A/B Testing", "Test different email variations", "true"},
		{"Advanced Analytics", "Detailed email performance metrics", "true"},
		{"Custom API Integration", "Custom API integrations with your systems", "true"},
		{"Send Time Optimization", "AI-powered send time optimization", "true"},
	}

	for _, feature := range enterpriseFeatures {
		err := createFeatureIfNotExists(ctx, queries, enterprisePlanID, feature.name, feature.description, feature.value)
		if err != nil {
			return err
		}
	}

	// Create mailing limits
	mailingLimits := []struct {
		planID               uuid.UUID
		dailyLimit           int32
		monthlyLimit         int32
		maxRecipientsPerMail int32
	}{
		{freePlanID, 500, 5000, 500},
		{basicPlanID, 2500, 25000, 2500},
		{proPlanID, 10000, 100000, 10000},
		{enterprisePlanID, 0, 0, 0}, // 0 indicates unlimited
	}

	for _, limit := range mailingLimits {
		err := createMailingLimitIfNotExists(ctx, queries, limit.planID, limit.dailyLimit, limit.monthlyLimit, limit.maxRecipientsPerMail)
		if err != nil {
			return err
		}
	}

	log.Println("Plans seeding completed successfully")
	return nil
}

// Helper function to create a feature if it doesn't exist
func createFeatureIfNotExists(ctx context.Context, queries db.Store, planID uuid.UUID, name, description, value string) error {
	// Check if feature exists for this plan
	features, err := queries.GetPlanFeaturesByPlanID(ctx, planID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for existing features: %w", err)
	}

	// Check if feature exists
	featureExists := false
	for _, feature := range features {
		if feature.Name.String == name {
			featureExists = true
			break
		}
	}

	if !featureExists {
		_, err = queries.CreatePlanFeature(ctx, db.CreatePlanFeatureParams{
			ID:          uuid.New(),
			PlanID:      planID,
			Name:        sql.NullString{String: name, Valid: true},
			Description: sql.NullString{String: description, Valid: true},
			Value:       sql.NullString{String: value, Valid: true},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			return fmt.Errorf("error creating feature %s for plan %s: %w", name, planID, err)
		}
		log.Printf("Created feature: %s for plan %s", name, planID)
	} else {
		log.Printf("Feature already exists: %s for plan %s", name, planID)
	}
	return nil
}

// Helper function to create a mailing limit if it doesn't exist
func createMailingLimitIfNotExists(ctx context.Context, queries db.Store, planID uuid.UUID, dailyLimit, monthlyLimit, maxRecipientsPerMail int32) error {
	// Check if mailing limit exists for this plan
	_, err := queries.GetMailingLimitByPlanID(ctx, planID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for existing mailing limit: %w", err)
	}

	if err == sql.ErrNoRows {
		_, err = queries.CreateMailingLimit(ctx, db.CreateMailingLimitParams{
			ID:                   uuid.New(),
			PlanID:               planID,
			DailyLimit:           sql.NullInt32{Int32: dailyLimit, Valid: true},
			MonthlyLimit:         sql.NullInt32{Int32: monthlyLimit, Valid: true},
			MaxRecipientsPerMail: sql.NullInt32{Int32: maxRecipientsPerMail, Valid: true},
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		})
		if err != nil {
			return fmt.Errorf("error creating mailing limit for plan %s: %w", planID, err)
		}
		log.Printf("Created mailing limit for plan %s", planID)
	} else {
		log.Printf("Mailing limit already exists for plan %s", planID)
	}
	return nil
}
