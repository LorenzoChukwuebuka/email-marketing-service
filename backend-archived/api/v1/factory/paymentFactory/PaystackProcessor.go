package paymentmethodFactory

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/utils"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type PaystackPaymentProcessor struct {
	paid bool
}

var (
	config       = utils.LoadEnv()
	key          = config.PaystackKey
	api_base     = config.PaystackBaseURL
	callback_url = config.PAYSTACK_CALLBACK_URL
)

func (c *PaystackPaymentProcessor) OpenDeposit(d *dto.BasePaymentModelData) (map[string]interface{}, error) {
	url := api_base + "transaction/initialize"

	data := map[string]interface{}{
		"amount":       d.AmountToPay * 100,
		"email":        d.Email,
		"callback_url": callback_url,
		"metadata": map[string]interface{}{
			"user_id":  d.UserId,
			"plan_id":  d.PlanId,
			"duration": d.Duration,
		},
	}

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

func (c *PaystackPaymentProcessor) ProcessDeposit(d *dto.BaseProcessPaymentModel) (*dto.BasePaymentResponse, error) {
	url := fmt.Sprintf(api_base+"transaction/verify/%s", d.Reference)
	//ctx := context.TODO()
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+key).
		SetHeader("Accept", "application/json").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("error decoding response body: %s", err)
	}

	if status, ok := result["status"].(bool); ok && !status {
		return nil, fmt.Errorf("transaction failed")
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	amount := data["amount"].(float64) / 100
	planIDStr, _ := data["metadata"].(map[string]any)["plan_id"].(string)
	userID, _ := data["metadata"].(map[string]any)["user_id"].(string)
	duration, _ := data["metadata"].(map[string]any)["duration"].(string)
	email, _ := data["customer"].(map[string]any)["email"].(string)
	status, _ := data["status"].(string)

	paymentData := &dto.BasePaymentResponse{
		Amount:   amount,
		PlanID:   planIDStr,
		UserID:   userID,
		Duration: duration,
		Email:    email,
		Status:   status,
	}

	return paymentData, nil
}

func (c *PaystackPaymentProcessor) OpenRefund() {
	fmt.Printf("Charged $ to Credit Card\n")
}

func (c *PaystackPaymentProcessor) ProcessRefund() {
	fmt.Printf("Refunded to Credit Card")
}

func (c *PaystackPaymentProcessor) ChargeCard(amount float64) {}

func (c *PaystackPaymentProcessor) Status() string {
	if c.paid {
		return "Paid"
	}
	return "Unpaid"
}
