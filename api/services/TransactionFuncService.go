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
		// return FlutterWavePaymentService
		c.selectedPaymentMethod = &paymentmethods.FlutterwavePaymentProcessor{}
	case "Paystack":
		c.selectedPaymentMethod = &paymentmethods.PaystackPaymentProcessor{}
	default:
		return fmt.Errorf("invalid payment type: %s", paymentType)
	}

	return nil
}

func (c *Transaction) OpenProcessPayment(
	d *model.InitPaymentModelData,
) (map[string]interface{}, error) {
	return c.selectedPaymentMethod.InitializePaymentProcess(d)
}

func Factory(paymentMethod string) (paymentmethods.PaymentInterface, error) {
	var pI paymentmethods.PaymentInterface
	switch paymentMethod {
	case "FlutterWave":
		// Return a concrete implementation of Flutterwave
		return nil, nil
	case "Paystack":
		pI = &paymentmethods.PaystackPaymentProcessor{}
		return pI, nil
	}

	return nil, nil
}

func UsePayment() error {
	paymentService, err := Factory("Paystack")
	if err != nil {
		return err
	}

	paymentService.Pay(20)
	return nil
}
