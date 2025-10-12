package paymentmethodFactory

import (
	"crypto/hmac"
	"crypto/sha512"
	"email-marketing-service/internal/config"
	"email-marketing-service/internal/domain"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
)

type PaystackPaymentProcessor struct {
	paid bool
}

var (
	cfg          = config.LoadEnv()
	key          = cfg.PaystackKey
	api_base     = cfg.PaystackBaseURL
	callback_url = cfg.PAYSTACK_CALLBACK_URL
)

func (c *PaystackPaymentProcessor) OpenDeposit(d *domain.BasePaymentModelData) (map[string]interface{}, error) {
	url := api_base + "transaction/initialize"

	data := map[string]interface{}{
		"amount":       d.AmountToPay * 100,
		"email":        d.Email,
		"callback_url": callback_url,
		"metadata": map[string]interface{}{
			"user_id":           d.UserId,
			"plan_id":           d.PlanId,
			"duration":          d.Duration,
			"company_id":        d.CompanyId,
			"payment_intent_id": d.PaymentIntentID,
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

func (c *PaystackPaymentProcessor) ProcessDeposit(d *domain.BaseProcessPaymentModel) (*domain.BasePaymentResponse, error) {
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

	fmt.Printf("Paystack verification data: %+v\n", data)

	amount := data["amount"].(float64) / 100
	planIDStr, _ := data["metadata"].(map[string]any)["plan_id"].(string)
	userID, _ := data["metadata"].(map[string]any)["user_id"].(string)
	duration, _ := data["metadata"].(map[string]any)["duration"].(string)
	email, _ := data["customer"].(map[string]any)["email"].(string)
	paymentIntentID, _ := data["metadata"].(map[string]any)["payment_intent_id"].(string)
	companyIDD, _ := data["metadata"].(map[string]any)["company_id"].(string)
	status, _ := data["status"].(string)

	paymentData := &domain.BasePaymentResponse{
		Amount:          amount,
		PlanID:          planIDStr,
		UserID:          userID,
		Duration:        duration,
		Email:           email,
		Status:          status,
		PaymentIntentID: paymentIntentID,
		CompanyID:       companyIDD,
	}

	return paymentData, nil
}

func (c *PaystackPaymentProcessor) WebhookHandler(r *http.Request) (*domain.BasePaymentResponse, error) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body: %w", err)
	}
	defer r.Body.Close()

	// Verify Paystack signature
	signature := r.Header.Get("x-paystack-signature")
	if !c.verifySignature(body, signature) {
		return nil, fmt.Errorf("invalid webhook signature")
	}

	// Decode JSON body
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	// Extract event + data
	event, _ := payload["event"].(string)
	if event != "charge.success" {
		// Ignore non-successful charge events
		return nil, fmt.Errorf("event %s ignored", event)
	}

	data, ok := payload["data"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid webhook format")
	}

	// Map webhook payload into BasePaymentResponse
	amount := data["amount"].(float64) / 100
	planIDStr, _ := data["metadata"].(map[string]any)["plan_id"].(string)
	userID, _ := data["metadata"].(map[string]any)["user_id"].(string)
	duration, _ := data["metadata"].(map[string]any)["duration"].(string)
	email, _ := data["customer"].(map[string]any)["email"].(string)
	paymentIntentID, _ := data["metadata"].(map[string]any)["payment_intent_id"].(string)
	companyIDD, _ := data["metadata"].(map[string]any)["company_id"].(string)
	status, _ := data["status"].(string)

	return &domain.BasePaymentResponse{
		Amount:          amount,
		PlanID:          planIDStr,
		UserID:          userID,
		Duration:        duration,
		Email:           email,
		Status:          status,
		PaymentIntentID: paymentIntentID,
		CompanyID:       companyIDD,
	}, nil
}

func (c *PaystackPaymentProcessor) verifySignature(body []byte, signature string) bool {
	h := hmac.New(sha512.New, []byte(key))
	h.Write(body)
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
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
