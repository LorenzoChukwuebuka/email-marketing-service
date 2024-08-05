package paymentmethodFactory

import (
	"email-marketing-service/api/v1/dto"
	 
	"fmt"
)

type FlutterwavePaymentProcessor struct {
	paid bool
}

func (c *FlutterwavePaymentProcessor) OpenDeposit(d *dto.BasePaymentModelData)(map[string]interface{},error) {
	return nil,nil
}

func (c *FlutterwavePaymentProcessor) ProcessDeposit(amount float64) {
	fmt.Printf("Paid $%.2f using FlutterwavePaymentProcessor Card\n", amount)
	c.paid = true
}

func (c *FlutterwavePaymentProcessor) OpenRefund() {
	fmt.Printf("Charged $ to Credit Card\n")
}

func (c *FlutterwavePaymentProcessor) ProcessRefund() {
	fmt.Printf("Refunded to Credit Card")
}

func (c *FlutterwavePaymentProcessor) ChargeCard(amount float64) {}

func (c *FlutterwavePaymentProcessor) Status() string {
	if c.paid {
		return "Paid"
	}
	return "Unpaid"
}
