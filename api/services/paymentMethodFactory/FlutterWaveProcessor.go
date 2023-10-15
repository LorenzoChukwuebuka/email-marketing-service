package paymentmethodFactory

import (
	"email-marketing-service/api/model"
	"fmt"
)

type FlutterwavePaymentProcessor struct {
	paid bool
}

func (c *FlutterwavePaymentProcessor) Initialize(d *model.PaymentModel) {

}

func (c *FlutterwavePaymentProcessor) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Credit FlutterwavePaymentProcessor\n", amount)
	c.paid = true
}

func (c *FlutterwavePaymentProcessor) Charge(amount float64) {
	fmt.Printf("Charged $%.2f to Credit Card\n", amount)
}

func (c *FlutterwavePaymentProcessor) Refund(amount float64) {
	fmt.Printf("Refunded $%.2f to Credit Card\n", amount)
}

func (c *FlutterwavePaymentProcessor) Status() string {
	if c.paid {
		return "Paid"
	}
	return "Unpaid"
}
