package services

import (
	"email-marketing-service/api/v1/dto"
	paymentmethodFactory "email-marketing-service/api/v1/factory/paymentFactory"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type BillingService struct {
	BillingRepo      *repository.BillingRepository
	SubscriptionSVC  *SubscriptionService
	UserRepo         *repository.UserRepository
	SubscriptionRepo *repository.SubscriptionRepository
	PlanRepo         *repository.PlanRepository
}

func NewBillingService(
	billingRepository *repository.BillingRepository,
	subscriptionSVC *SubscriptionService,
	userRepo *repository.UserRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	planRepo *repository.PlanRepository) *BillingService {
	return &BillingService{
		BillingRepo:      billingRepository,
		SubscriptionSVC:  subscriptionSVC,
		UserRepo:         userRepo,
		SubscriptionRepo: subscriptionRepo,
		PlanRepo:         planRepo,
	}
}

func (s *BillingService) ConfirmPayment(paymentmethod string, reference string) (map[string]interface{}, error) {
	tx := s.BillingRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	paymenservice, err := paymentmethodFactory.PaymentFactory(paymentmethod)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error instantiating factory: %s", err)
	}

	params := &dto.BaseProcessPaymentModel{
		PaymentMethod: paymentmethod,
		Reference:     reference,
	}

	data, err := paymenservice.ProcessDeposit(params)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	checkIfRefExists, err := s.BillingRepo.CheckIfRefExists(reference)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if checkIfRefExists {
		return nil, fmt.Errorf("this payment has been confirmed")
	}

	//get the planId
	planR, err := s.PlanRepo.GetSinglePlan(data.PlanID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//get the userId
	userUUID := &model.User{UUID: data.UserID}

	userId, err := s.UserRepo.FindUserById(userUUID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.cancelFreePlan(userId.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	transactionId := utils.GenerateOTP(10)

	billingServiceData := &model.Billing{
		UUID:          uuid.New().String(),
		UserId:        userId.ID,
		AmountPaid:    float32(data.Amount),
		PlanId:        planR.ID,
		Email:         data.Email,
		Duration:      data.Duration,
		ExpiryDate:    calculateExpiryDate(data.Duration),
		Reference:     reference,
		TransactionId: transactionId,
		PaymentMethod: paymentmethod,
		Status:        data.Status,
	}

	billingRepo, err := s.BillingRepo.CreateBilling(billingServiceData)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	subscription := &model.Subscription{
		UUID:          uuid.New().String(),
		UserId:        userId.ID,
		PlanId:        planR.ID,
		PaymentId:     int(billingRepo.ID),
		StartDate:     time.Now(),
		EndDate:       calculateExpiryDate(data.Duration),
		Expired:       false,
		TransactionId: transactionId,
	}

	_, err = s.SubscriptionSVC.CreateSubscription(subscription)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return map[string]interface{}{
		"data": "payment  successful",
	}, nil
}

func calculateExpiryDate(duration string) time.Time {
	parts := strings.Split(duration, " ")
	num, _ := strconv.Atoi(parts[0])
	unit := parts[1]

	switch unit {
	case "week":
		return time.Now().AddDate(0, 0, num*7)
	case "day":
		return time.Now().AddDate(0, 0, num)
	case "year":
		return time.Now().AddDate(num, 0, 0)
	case "month":
		return time.Now().AddDate(0, num, 0)
	default:
		return time.Now()
	}
}

func (s *BillingService) cancelFreePlan(userId uint) error {
	// Retrieve the user's current subscription
	getUserCurrentSub, err := s.SubscriptionSVC.GetUsersCurrentSubscription(userId)
	if err != nil {
		// If there's an error (e.g., no active subscription), just return nil (no action)
		return nil
	}

	// Get the name of the current plan
	planName := getUserCurrentSub.Plan.PlanName

	// Check if the current plan is the free plan (case insensitive)
	if strings.ToLower(planName) == "free" {
		// If it's a free plan, update its status to expired
		err := s.SubscriptionRepo.UpdateExpiredSubscription(getUserCurrentSub.Id)
		if err != nil {
			// Return the error if the update fails
			return err
		}
	}

	// Continue without any errors if no action was needed or the action succeeded
	return nil
}

func (s *BillingService) GetSingleBillingRecord(biilingId string, userId int) (*model.BillingResponse, error) {

	billing, err := s.BillingRepo.GetSingleBillingRecord(biilingId, userId)
	if err != nil {
		return nil, err
	}

	if billing == nil {
		return nil, fmt.Errorf("no record found: %w", err)
	}

	return billing, nil
}

func (s *BillingService) GetAllBillingForAUser(userId string, page int, pageSize int) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}

	usermodel := &model.User{UUID: userId}

	user, err := s.UserRepo.FindUserById(usermodel)
	if err != nil {
		return repository.PaginatedResult{}, err
	}

	billing, err := s.BillingRepo.GetAllPayments(int(user.ID), paginationParams)
	if err != nil {
		return repository.PaginatedResult{}, err
	}

	// Filter out the billing records where the plan name is "free"
	var filteredBilling []model.BillingResponse
	originalLength := len(billing.Data.([]model.BillingResponse))
	for _, b := range billing.Data.([]model.BillingResponse) {
		if strings.ToLower(b.Plan.PlanName) != "free" {
			filteredBilling = append(filteredBilling, b)
		}
	}

	billing.Data = filteredBilling

	if len(filteredBilling) == 0 {
		return repository.PaginatedResult{}, nil
	}

	removedItems := originalLength - len(filteredBilling)

	billing.TotalCount = billing.TotalCount - int64(removedItems)

	return billing, nil
}
