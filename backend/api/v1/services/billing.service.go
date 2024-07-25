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
}

func NewBillingService(billingRepository *repository.BillingRepository,
	subscriptionSVC *SubscriptionService,
	userRepo *repository.UserRepository, subscriptionRepo *repository.SubscriptionRepository) *BillingService {
	return &BillingService{
		BillingRepo:      billingRepository,
		SubscriptionSVC:  subscriptionSVC,
		UserRepo:         userRepo,
		SubscriptionRepo: subscriptionRepo,
	}
}

func (s *BillingService) ConfirmPayment(paymentmethod string, reference string) (map[string]interface{}, error) {
	paymenservice, err := paymentmethodFactory.PaymentFactory(paymentmethod)

	if err != nil {
		return nil, fmt.Errorf("error instantiating factory: %s", err)
	}

	params := &dto.BaseProcessPaymentModel{
		PaymentMethod: paymentmethod,
		Reference:     reference,
	}

	data, err := paymenservice.ProcessDeposit(params)

	if err != nil {
		return nil, err
	}

	//get the userId
	userUUID := &model.User{UUID: data.UserID}

	userId, err := s.UserRepo.FindUserById(userUUID)

	if err != nil {
		return nil, err
	}

	err = s.cancelFreePlan(userId.ID)

	if err != nil {
		return nil, err
	}

	transactionId := utils.GenerateOTP(10)

	billingServiceData := &model.Billing{
		UUID:          uuid.New().String(),
		UserId:        userId.ID,
		AmountPaid:    float32(data.Amount),
		PlanId:        uint(data.PlanID),
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
		return nil, err
	}

	subscription := &model.Subscription{
		UUID:          uuid.New().String(),
		UserId:        userId.ID,
		PlanId:        uint(data.PlanID),
		PaymentId:     int(billingRepo.ID),
		StartDate:     time.Now(),
		EndDate:       calculateExpiryDate(data.Duration),
		Expired:       false,
		TransactionId: transactionId,
	}

	_, err = s.SubscriptionSVC.CreateSubscription(subscription)

	if err != nil {
		return nil, err
	}

	fmt.Print(billingRepo)
	return nil, nil
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
	getUserCurrentSub, err := s.SubscriptionSVC.GetUsersCurrentSubscription(userId)

	if err != nil {
		return nil
	}

	fmt.Printf("%+v\n", getUserCurrentSub)

	if getUserCurrentSub.Plan.PlanName == "free" {
		err := s.SubscriptionRepo.UpdateExpiredSubscription(getUserCurrentSub.Id)

		if err != nil {
			return nil
		}

	}

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

func (s *BillingService) GetAllBillingForAUser(userId int, page int) ([]model.Billing, error) {
	billing, err := s.BillingRepo.GetAllPayments(userId, page)
	if err != nil {
		return nil, err
	}
	if billing == nil {
		return nil, fmt.Errorf("no record found: %w", err)
	}
	return billing, nil
}
