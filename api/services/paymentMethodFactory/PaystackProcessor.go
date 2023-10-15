package paymentmethodFactory

import (
	"email-marketing-service/api/model"
	"fmt"
)

type PaystackPaymentProcessor struct {
	paid bool
}

func (c *PaystackPaymentProcessor) Initialize(d *model.PaymentModel) {

}

func (c *PaystackPaymentProcessor) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using PaystackPaymentProcessor Card\n", amount)
	c.paid = true
}

func (c *PaystackPaymentProcessor) Charge(amount float64) {
	fmt.Printf("Charged $%.2f to Credit Card\n", amount)
}

func (c *PaystackPaymentProcessor) Refund(amount float64) {
	fmt.Printf("Refunded $%.2f to Credit Card\n", amount)
}

func (c *PaystackPaymentProcessor) Status() string {
	if c.paid {
		return "Paid"
	}
	return "Unpaid"
}
