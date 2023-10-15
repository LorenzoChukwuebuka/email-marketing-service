package services

import "fmt"

import paymentmethods "email-marketing-service/api/services/paymentMethods"

type Customer struct {
	selectedPaymentMethod paymentmethods.PaymentInterface
}

func (c *Customer) ChoosePaymentMethod(paymentType string) {
	switch paymentType {
	case "creditCard":
		c.selectedPaymentMethod = &paymentmethods.CreditCard{}
	case "flutterwave":
		c.selectedPaymentMethod = &paymentmethods.Flutterwave{}
	case "paystack":
		c.selectedPaymentMethod = &paymentmethods.Paystack{}
	default:
		fmt.Println("Invalid payment type")
	}
}
