package paymentmethodFactory

import (
	"email-marketing-service/api/model"
	"fmt"
)

type CreditCardPaymentProcessor struct {
    paid bool
}


func (c *CreditCardPaymentProcessor) Initialize(d *model.PaymentModel){

}

func (c *CreditCardPaymentProcessor) Pay(amount float64) {
    fmt.Printf("Paid $%.2f using Credit Card\n", amount)
    c.paid = true
}

func (c *CreditCardPaymentProcessor) Charge(amount float64) {
    fmt.Printf("Charged $%.2f to Credit Card\n", amount)
}

func (c *CreditCardPaymentProcessor) Refund(amount float64) {
    fmt.Printf("Refunded $%.2f to Credit Card\n", amount)
}

func (c *CreditCardPaymentProcessor) Status() string {
    if c.paid {
        return "Paid"
    }
    return "Unpaid"
}
