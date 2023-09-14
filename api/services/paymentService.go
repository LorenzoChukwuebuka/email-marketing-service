package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type PaymentService struct {
	PaymentRepo *repository.PaymentRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		PaymentRepo: paymentRepo,
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

func (s *PaymentService) ConfirmPayment(reference string) (map[string]interface{}, error) {
	utils.LoadEnv()

	var (
		key      = os.Getenv("PAYSTACK_KEY")
		api_base = os.Getenv("PAYSTACK_BASE_URL")
	)

	url := fmt.Sprintf(api_base+"transaction/verify/%s", reference)

	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+key).
		SetHeader("Accept", "application/json").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("error decoding response body: %s", err)
	}

	if status, ok := result["status"].(bool); ok && !status {
		return nil, fmt.Errorf("transaction failed")
	}

	data := result["data"].(map[string]interface{})

	amount := data["amount"].(float64) / 100
	planIDStr := data["metadata"].(map[string]interface{})["plan_id"].(string)
	userIDStr := data["metadata"].(map[string]interface{})["user_id"].(string)
	duration := data["metadata"].(map[string]interface{})["duration"].(string)
	email := data["customer"].(map[string]interface{})["email"].(string)

	//convert planId and userId to int

	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		return nil, fmt.Errorf("error converting planID to int: %s", err)
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("error converting userID to int: %s", err)
	}

	//split the duration parts
	//perform a logic with it to get the exipry date for the payment made
	parts := strings.Split(duration, " ")

	num := parts[0]
	unit := parts[1]

	var expiryDate time.Time

	switch unit {
	case "week":
		numWeeks, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("error converting num to int: %s", err)
		}
		expiryDate = time.Now().AddDate(0, 0, numWeeks*7)
	case "day":
		numDays, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("error converting num to int: %s", err)
		}
		expiryDate = time.Now().AddDate(0, 0, numDays)
	case "year":
		numYears, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("error converting num to int: %s", err)
		}
		expiryDate = time.Now().AddDate(numYears, 0, 0)
	case "month":
		numMonths, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("error converting num to int: %s", err)
		}
		expiryDate = time.Now().AddDate(0, numMonths, 0)
	default:
		return nil, fmt.Errorf("unsupported unit: %s", unit)
	}

	tx, err := s.PaymentRepo.DB.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	payment := &model.PaymentModel{
		AmountPaid: float32(amount),
		PlanId:     planID,
		UserId:     userID,
		Duration:   duration,
		Reference:  reference,
		ExpiryDate: expiryDate,
		Email:      email,
		CreatedAt:  time.Now(),
	}

	paymentRepo, err := s.PaymentRepo.CreatePayment(payment)

	if err != nil {
		return nil, err
	}

	subscription := &model.SubscriptionModel{
		UserId:    userID,
		PlanId:    planID,
		PaymentId: paymentRepo.Id,
		StartDate: time.Now(),
		EndDate:   expiryDate,
		Expired:   false,
		CreatedAt: time.Now(),
	}

	fmt.Print(subscription)

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %s", err)
	}

	return result, nil
}

func (s *PaymentService) GetAllPaymentsForAUser(userId int) ([]model.PaymentResponse, error) {
	return nil, nil
}

func (s *PaymentService) GetSinglePaymentForAUser(userId int, paymentId int) (*model.PaymentResponse, error) {
	return nil, nil
}
