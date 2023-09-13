package services

import (
	"bytes"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type PaymentService struct {
	PaymentRepo *repository.PaymentRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		PaymentRepo: paymentRepo,
	}
}

var key = os.Getenv("PAYSTACK_KEY")
var api_base = os.Getenv("PAYSTACK_BASE_API")

func (s *PaymentService) InitializePayment(d *model.PaymentModel) (map[string]interface{}, error) {
	utils.LoadEnv()

	url := api_base + "transaction/initialize"

	data := map[string]interface{}{
		"amount": d.AmountPaid * 100,
		"email":  d.Email,
		"metadata": map[string]interface{}{
			"user_id": d.UserId,
			"plan_id": d.PlanId,
		},
	}

	//marshal to json

	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	//create a new client
	client := &http.Client{}

	//create a new POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	//set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal response JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *PaymentService) ConfirmPayment(d *model.PaymentModel) error {
	return nil
}

func (s *PaymentService) GetAllPaymentsForAUser(userId int) ([]model.PaymentResponse, error) {
	return nil, nil
}

func (s *PaymentService) GetSinglePaymentForAUser(userId int, paymentId int) (*model.PaymentResponse, error) {
	return nil, nil
}
