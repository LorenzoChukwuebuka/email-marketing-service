package paymentmethodFactory

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/utils"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type PaystackPaymentProcessor struct {
	paid bool
}

var (
	config   = utils.LoadEnv()
	key      = config.PaystackKey
	api_base = config.PaystackBaseURL
)

func (c *PaystackPaymentProcessor) OpenDeposit(d *model.BasePaymentModelData) (map[string]interface{}, error) {
	url := api_base + "transaction/initialize"

	data := map[string]interface{}{
		"amount": d.AmountToPay * 100,
		"email":  d.Email,
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

func (c *PaystackPaymentProcessor) ProcessDeposit(amount float64) {
	fmt.Printf("Paid $%.2f using PaystackPaymentProcessor Card\n", amount)
	c.paid = true
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
