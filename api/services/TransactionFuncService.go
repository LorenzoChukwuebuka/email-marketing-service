package services

import (
	"fmt"
	"email-marketing-service/api/model"
	paymentmethods "email-marketing-service/api/services/paymentMethodFactory"
)

type Transaction struct {
	selectedPaymentMethod paymentmethods.PaymentInterface
}



func (c *Transaction) ChoosePaymentMethod(paymentType string) error {
	switch paymentType {
	case "FlutterWave":
		c.selectedPaymentMethod = &paymentmethods.FlutterwavePaymentProcessor{}
	case "Paystack":
		c.selectedPaymentMethod = &paymentmethods.PaystackPaymentProcessor{}
	default:
		return fmt.Errorf("invalid payment type: %s", paymentType)
	}

	return nil
}

func (c *Transaction) OpenProcessPayment(d *model.InitPaymentModelData) (map[string]interface{}, error) {
	return c.selectedPaymentMethod.InitializePaymentProcess(d)

}
