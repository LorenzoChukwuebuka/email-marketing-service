package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type PaymentService struct {
	PaymentRepo     *repository.PaymentRepository
	SubscriptionSvc *SubscriptionService
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, subscriptionSvc *SubscriptionService) *PaymentService {
	return &PaymentService{
		PaymentRepo:     paymentRepo,
		SubscriptionSvc: subscriptionSvc,
	}
}

func (s *PaymentService) InitializePayment(d *model.PaymentModel) (map[string]interface{}, error) {
	utils.LoadEnv()

	var (
		key      = os.Getenv("PAYSTACK_KEY")
		api_base = os.Getenv("PAYSTACK_BASE_URL")
	)

	url := api_base + "transaction/initialize"

	data := map[string]interface{}{
		"amount": d.AmountPaid * 100,
		"email":  d.Email,
		"metadata": map[string]interface{}{
			"user_id":  d.UserId,
			"plan_id":  d.PlanId,
			"duration": d.Duration,
		},
	}

	fmt.Print(data)
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", key)).
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PaymentService) ConfirmPayment(reference string) (string, error) {
	utils.LoadEnv()

	key := os.Getenv("PAYSTACK_KEY")
	api_base := os.Getenv("PAYSTACK_BASE_URL")

	url := fmt.Sprintf(api_base+"transaction/verify/%s", reference)

	//ctx := context.TODO()

	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+key).
		SetHeader("Accept", "application/json").
		Get(url)

	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("error decoding response body: %s", err)
	}

	if status, ok := result["status"].(bool); ok && !status {
		return "", fmt.Errorf("transaction failed")
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	amount := data["amount"].(float64) / 100
	planIDStr, _ := data["metadata"].(map[string]interface{})["plan_id"].(string)
	userIDStr, _ := data["metadata"].(map[string]interface{})["user_id"].(string)
	duration, _ := data["metadata"].(map[string]interface{})["duration"].(string)
	email, _ := data["customer"].(map[string]interface{})["email"].(string)

	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		return "", fmt.Errorf("error converting planID to int: %s", err)
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return "", fmt.Errorf("error converting userID to int: %s", err)
	}

	// Parse duration and calculate expiry date

	tx, err := s.PaymentRepo.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("error starting transaction: %s", err)
	}

	defer func() {
		if r := recover(); r != nil {
			// An error occurred, rollback the transaction
			tx.Rollback()
		}
	}()

	payment := &model.PaymentModel{
		AmountPaid: float32(amount),
		UUID:       uuid.New().String(),
		PlanId:     planID,
		UserId:     userID,
		Duration:   duration,
		Reference:  reference,
		ExpiryDate: calculateExpiryDate(duration),
		Email:      email,
		CreatedAt:  time.Now(),
		Status:     "active",
	}

	paymentRepo, err := s.PaymentRepo.CreatePayment(payment)
	if err != nil {
		return "", fmt.Errorf("error creating payment: %s", err)
	}

	subscription := &model.SubscriptionModel{
		UUID:      uuid.New().String(),
		UserId:    userID,
		PlanId:    planID,
		PaymentId: paymentRepo.Id,
		StartDate: time.Now(),
		EndDate:   calculateExpiryDate(duration),
		Expired:   false,
		CreatedAt: time.Now(),
	}

	_, err = s.SubscriptionSvc.CreateSubscription(subscription)
	if err != nil {
		return "", fmt.Errorf("error creating subscription: %s", err)
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("error committing transaction: %s", err)
	}

	return "payment verified successfully", nil
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

func (s *PaymentService) GetAllPaymentsForAUser(userId int) ([]model.PaymentResponse, error) {
	paymentRepo, err := s.PaymentRepo.GetAllPayments(userId)

	if err != nil {
		return nil, err
	}

	if len(paymentRepo) == 0 {
		return nil, fmt.Errorf("no records found for this user")
	}

	return paymentRepo, nil
}

func (s *PaymentService) GetSinglePaymentForAUser(userId string, paymentId string) (*model.PaymentResponse, error) {
	paymentRepo, err := s.PaymentRepo.GetSinglePayment(paymentId, userId)

	if err != nil {
		return nil, err
	}

	if paymentRepo == nil {
		return nil, fmt.Errorf("no records found for this user")
	}
	return paymentRepo, nil
}
