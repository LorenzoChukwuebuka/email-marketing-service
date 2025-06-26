package service

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/admin/plans/dto"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PlanService struct {
	store db.Store
}

func NewPlanService(store db.Store) *PlanService {
	return &PlanService{
		store: store,
	}
}

func (s *PlanService) CreatePlan(ctx context.Context, req *dto.Plan) (*dto.PlanResponse, error) {
	// Check if plan with the same name already exists
	existingPlan, err := s.store.PlanExists(ctx, req.PlanName)

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	if existingPlan {
		return nil, fmt.Errorf("plan with name '%s' already exists", req.PlanName)
	}

	// Generate new UUID for the plan
	planID := uuid.New()
	now := time.Now()

	// Create the plan
	createdPlan, err := s.store.CreatePlan(ctx, db.CreatePlanParams{
		ID:           planID,
		Name:         req.PlanName,
		Description:  sql.NullString{String: req.Description, Valid: req.Description != ""},
		Price:        decimal.NewFromFloat(req.Price),
		BillingCycle: sql.NullString{String: req.BillingCycle, Valid: req.BillingCycle != ""},
		Status:       sql.NullString{String: string(req.Status), Valid: true},
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating plan: %w", err)
	}

	// Create features for the plan
	var createdFeatures []dto.PlanFeature
	for _, feature := range req.Features {
		_, err = s.store.CreatePlanFeature(ctx, db.CreatePlanFeatureParams{
			ID:          uuid.New(),
			PlanID:      planID,
			Name:        sql.NullString{String: feature.Name, Valid: feature.Name != ""},
			Description: sql.NullString{String: feature.Description, Valid: feature.Description != ""},
			Value:       sql.NullString{String: feature.Value, Valid: feature.Value != ""},
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating feature '%s': %w", feature.Name, err)
		}
		createdFeatures = append(createdFeatures, feature)
	}

	// Create mailing limits for the plan
	_, err = s.store.CreateMailingLimit(ctx, db.CreateMailingLimitParams{
		ID:                   uuid.New(),
		PlanID:               planID,
		DailyLimit:           sql.NullInt32{Int32: req.MailingLimits.DailyLimit, Valid: true},
		MonthlyLimit:         sql.NullInt32{Int32: req.MailingLimits.MonthlyLimit, Valid: true},
		MaxRecipientsPerMail: sql.NullInt32{Int32: req.MailingLimits.MaxRecipientsPerMail, Valid: true},
		CreatedAt:            now,
		UpdatedAt:            now,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating mailing limits: %w", err)
	}

	// Prepare response
	response := &dto.PlanResponse{
		ID:            createdPlan.ID,
		Name:          createdPlan.Name,
		Description:   createdPlan.Description.String,
		Price:         createdPlan.Price.InexactFloat64(),
		BillingCycle:  createdPlan.BillingCycle.String,
		Status:        createdPlan.Status.String,
		Features:      createdFeatures,
		MailingLimits: req.MailingLimits,
		CreatedAt:     createdPlan.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     createdPlan.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, req *dto.EditPlan) (*dto.PlanResponse, error) {
	// Parse UUID from string
	planID, err := uuid.Parse(req.UUID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan UUID: %w", err)
	}

	// Check if plan exists
	existingPlan, err := s.store.GetPlanByID(ctx, planID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plan with ID '%s' not found", req.UUID)
		}
		return nil, fmt.Errorf("error fetching plan: %w", err)
	}

	// Check if another plan with the same name exists (excluding current plan)
	if req.PlanName != "" && req.PlanName != existingPlan.Name {
		nameExists, err := s.store.PlanExists(ctx, req.PlanName)
		if err != nil {
			return nil, common.ErrFetchingRecord
		}
		if nameExists {
			return nil, fmt.Errorf("plan with name '%s' already exists", req.PlanName)
		}
	}

	// Prepare update parameters with COALESCE logic
	var name, description, billingCycle, status interface{}
	var price interface{}

	print(price)

	if req.PlanName != "" {
		name = req.PlanName
	} else {
		name = nil
	}

	if req.Description != "" {
		description = req.Description
	} else {
		description = nil
	}

	if req.Price != 0 {
		price = decimal.NewFromFloat(req.Price)
	} else {
		price = nil
	}

	if req.BillingCycle != "" {
		billingCycle = req.BillingCycle
	} else {
		billingCycle = nil
	}

	if string(req.Status) != "" {
		status = string(req.Status)
	} else {
		status = nil
	}

	// Update plan basic information
	_, err = s.store.UpdatePlan(ctx, db.UpdatePlanParams{
		Name:         name.(string),
		Description:  sql.NullString{String: description.(string), Valid: description != nil},
		Price:        decimal.Decimal{},
		BillingCycle: sql.NullString{String: billingCycle.(string), Valid: billingCycle != nil},
		Status:       sql.NullString{String: status.(string), Valid: status != nil},
		ID:           planID,
	})
	if err != nil {
		return nil, fmt.Errorf("error updating plan: %w", err)
	}

	// Update features if provided
	if req.Features != nil {
		// Get existing features and soft delete them
		existingFeatures, err := s.store.GetPlanFeaturesByPlanID(ctx, planID)
		if err != nil {
			return nil, fmt.Errorf("error getting existing features: %w", err)
		}

		// Soft delete existing features
		for _, feature := range existingFeatures {
			err = s.store.DeletePlanFeature(ctx, feature.ID)
			if err != nil {
				return nil, fmt.Errorf("error soft deleting feature '%s': %w", feature.Name.String, err)
			}
		}

		// Create new features
		now := time.Now()
		for _, feature := range req.Features {
			_, err = s.store.CreatePlanFeature(ctx, db.CreatePlanFeatureParams{
				ID:          uuid.New(),
				PlanID:      planID,
				Name:        sql.NullString{String: feature.Name, Valid: feature.Name != ""},
				Description: sql.NullString{String: feature.Description, Valid: feature.Description != ""},
				Value:       sql.NullString{String: feature.Value, Valid: feature.Value != ""},
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err != nil {
				return nil, fmt.Errorf("error creating feature '%s': %w", feature.Name, err)
			}
		}
	}

	// Update mailing limits if provided
	if req.MailingLimits.DailyLimit != 0 || req.MailingLimits.MonthlyLimit != 0 || req.MailingLimits.MaxRecipientsPerMail != 0 {
		var dailyLimit, monthlyLimit, maxRecipients interface{}

		if req.MailingLimits.DailyLimit != 0 {
			dailyLimit = req.MailingLimits.DailyLimit
		} else {
			dailyLimit = nil
		}

		if req.MailingLimits.MonthlyLimit != 0 {
			monthlyLimit = req.MailingLimits.MonthlyLimit
		} else {
			monthlyLimit = nil
		}

		if req.MailingLimits.MaxRecipientsPerMail != 0 {
			maxRecipients = req.MailingLimits.MaxRecipientsPerMail
		} else {
			maxRecipients = nil
		}

		_, err = s.store.UpdateMailingLimit(ctx, db.UpdateMailingLimitParams{
			DailyLimit:           sql.NullInt32{Int32: dailyLimit.(int32), Valid: dailyLimit != nil},
			MonthlyLimit:         sql.NullInt32{Int32: monthlyLimit.(int32), Valid: monthlyLimit != nil},
			MaxRecipientsPerMail: sql.NullInt32{Int32: maxRecipients.(int32), Valid: maxRecipients != nil},
			PlanID:               planID,
		})
		if err != nil {
			return nil, fmt.Errorf("error updating mailing limits: %w", err)
		}
	}

	// Get updated plan details
	response, err := s.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated plan: %w", err)
	}

	return response, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, planID uuid.UUID) error {
	// Check if plan exists
	_, err := s.store.GetPlanByID(ctx, planID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("plan with ID '%s' not found", planID.String())
		}
		return fmt.Errorf("error fetching plan: %w", err)
	}

	// Check if plan is being used by any subscriptions using the ListPlansWithDetails query
	// This gives us the active_subscriptions_count
	plans, err := s.store.ListPlansWithDetails(ctx)
	if err != nil {
		return fmt.Errorf("error checking plan usage: %w", err)
	}

	// Find the specific plan and check if it has active subscriptions
	for _, plan := range plans {
		if plan.ID == planID {
			// Assuming the ListPlansWithDetails returns a field for active subscription count
			// You'll need to check the exact field name generated by SQLC
			// if plan.ActiveSubscriptionsCount > 0 {
			//     return fmt.Errorf("cannot delete plan: it has %d active subscription(s)", plan.ActiveSubscriptionsCount)
			// }
			break
		}
	}

	// Get and soft delete plan features
	features, err := s.store.GetPlanFeaturesByPlanID(ctx, planID)
	if err != nil {
		return fmt.Errorf("error getting plan features: %w", err)
	}

	for _, feature := range features {
		err = s.store.DeletePlanFeature(ctx, feature.ID)
		if err != nil {
			return fmt.Errorf("error soft deleting feature '%s': %w", feature.Name.String, err)
		}
	}

	// Soft delete the plan itself
	err = s.store.SoftDeletePlan(ctx, planID)
	if err != nil {
		return fmt.Errorf("error soft deleting plan: %w", err)
	}

	return nil
}

func (s *PlanService) ArchivePlan(ctx context.Context, planID uuid.UUID) (*dto.PlanResponse, error) {
	// Check if plan exists
	_, err := s.store.GetPlanByID(ctx, planID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plan with ID '%s' not found", planID.String())
		}
		return nil, fmt.Errorf("error fetching plan: %w", err)
	}

	// Archive the plan (sets status to 'archived')
	_, err = s.store.ArchivePlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error archiving plan: %w", err)
	}

	// Get the archived plan details
	response, err := s.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error fetching archived plan: %w", err)
	}

	return response, nil
}

// Helper method to get plan by ID with all related data
func (s *PlanService) GetPlanByID(ctx context.Context, planID uuid.UUID) (*dto.PlanResponse, error) {
	// Get plan
	plan, err := s.store.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting plan: %w", err)
	}

	// Get features
	features, err := s.store.GetPlanFeaturesByPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting plan features: %w", err)
	}

	// Get mailing limits
	mailingLimit, err := s.store.GetMailingLimitByPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting mailing limits: %w", err)
	}

	// Convert features
	var planFeatures []dto.PlanFeature
	for _, feature := range features {
		planFeatures = append(planFeatures, dto.PlanFeature{
			Name:        feature.Name.String,
			Description: feature.Description.String,
			Value:       feature.Value.String,
		})
	}

	// Prepare response
	response := &dto.PlanResponse{
		ID:           plan.ID,
		Name:         plan.Name,
		Description:  plan.Description.String,
		Price:        plan.Price.InexactFloat64(),
		BillingCycle: plan.BillingCycle.String,
		Status:       plan.Status.String,
		Features:     planFeatures,
		MailingLimits: dto.MailingLimits{
			DailyLimit:           mailingLimit.DailyLimit.Int32,
			MonthlyLimit:         mailingLimit.MonthlyLimit.Int32,
			MaxRecipientsPerMail: mailingLimit.MaxRecipientsPerMail.Int32,
		},
		CreatedAt: plan.CreatedAt.Format(time.RFC3339),
		UpdatedAt: plan.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *PlanService) GetPlanWithDetails(ctx context.Context, planID uuid.UUID) (*dto.PlanResponse, error) {
	// Use the optimized query that gets plan with features and mailing limits in one go
	planDetails, err := s.store.GetPlanWithDetails(ctx, planID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plan with ID '%s' not found", planID.String())
		}
		return nil, fmt.Errorf("error getting plan details: %w", err)
	}

	// Parse features from JSON (you'll need to handle this based on your SQLC generated types)
	var planFeatures []dto.PlanFeature
	// Note: You'll need to unmarshal the JSON features field here
	// This depends on how SQLC generates the struct for the json_agg result

	// For now, fall back to separate queries if JSON parsing is complex
	features, err := s.store.GetPlanFeaturesByPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting plan features: %w", err)
	}

	for _, feature := range features {
		planFeatures = append(planFeatures, dto.PlanFeature{
			Name:        feature.Name.String,
			Description: feature.Description.String,
			Value:       feature.Value.String,
		})
	}

	// Prepare response
	response := &dto.PlanResponse{
		ID:           planDetails.ID,
		Name:         planDetails.Name,
		Description:  planDetails.Description.String,
		Price:        planDetails.Price.InexactFloat64(),
		BillingCycle: planDetails.BillingCycle.String,
		Status:       planDetails.Status.String,
		Features:     planFeatures,
		MailingLimits: dto.MailingLimits{
			DailyLimit:           planDetails.DailyLimit.Int32,
			MonthlyLimit:         planDetails.MonthlyLimit.Int32,
			MaxRecipientsPerMail: planDetails.MaxRecipientsPerMail.Int32,
		},
		CreatedAt: planDetails.CreatedAt.Format(time.RFC3339),
		UpdatedAt: planDetails.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *PlanService) GetAllPlansWithDetails(ctx context.Context) ([]dto.PlanResponse, error) {
	// Use the optimized query that gets all plans with details
	plans, err := s.store.ListPlansWithDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting plans with details: %w", err)
	}

	var responses []dto.PlanResponse
	for _, plan := range plans {
		// Parse features from JSON or use separate query
		var planFeatures []dto.PlanFeature

		// For now, use separate query for features
		features, err := s.store.GetPlanFeaturesByPlanID(ctx, plan.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting plan features for plan %s: %w", plan.ID, err)
		}

		for _, feature := range features {
			planFeatures = append(planFeatures, dto.PlanFeature{
				Name:        feature.Name.String,
				Description: feature.Description.String,
				Value:       feature.Value.String,
			})
		}

		response := dto.PlanResponse{
			ID:           plan.ID,
			Name:         plan.Name,
			Description:  plan.Description.String,
			Price:        plan.Price.InexactFloat64(),
			BillingCycle: plan.BillingCycle.String,
			Status:       plan.Status.String,
			Features:     planFeatures,
			MailingLimits: dto.MailingLimits{
				DailyLimit:           plan.DailyLimit.Int32,
				MonthlyLimit:         plan.MonthlyLimit.Int32,
				MaxRecipientsPerMail: plan.MaxRecipientsPerMail.Int32,
			},
			CreatedAt: plan.CreatedAt.Format(time.RFC3339),
			UpdatedAt: plan.UpdatedAt.Format(time.RFC3339),
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *PlanService) GetSinglePlan(ctx context.Context, planID uuid.UUID) (*dto.PlanResponse, error) {
	// Get plan by ID
	plan, err := s.store.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting plan by ID: %w", err)
	}

	// Get features
	features, err := s.store.GetPlanFeaturesByPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting plan features: %w", err)
	}

	// Get mailing limits
	mailingLimit, err := s.store.GetMailingLimitByPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error getting mailing limits: %w", err)
	}

	// Convert features
	var planFeatures []dto.PlanFeature
	for _, feature := range features {
		planFeatures = append(planFeatures, dto.PlanFeature{
			Name:        feature.Name.String,
			Description: feature.Description.String,
			Value:       feature.Value.String,
		})
	}

	// Prepare response
	response := &dto.PlanResponse{
		ID:           plan.ID,
		Name:         plan.Name,
		Description:  plan.Description.String,
		Price:        plan.Price.InexactFloat64(),
		BillingCycle: plan.BillingCycle.String,
		Status:       plan.Status.String,
		Features:     planFeatures,
		MailingLimits: dto.MailingLimits{
			DailyLimit:           mailingLimit.DailyLimit.Int32,
			MonthlyLimit:         mailingLimit.MonthlyLimit.Int32,
			MaxRecipientsPerMail: mailingLimit.MaxRecipientsPerMail.Int32,
		},
		CreatedAt: plan.CreatedAt.Format(time.RFC3339),
		UpdatedAt: plan.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

