package paymentmethodFactory

import (
	"fmt"
)

func PaymentFactory(paymentMethod string) (PaymentInterface, error) {
	var pI PaymentInterface
	switch paymentMethod {
	case "FlutterWave":
		// Return a concrete implementation of Flutterwave
		return nil, nil
	case "Paystack":
		pI = &PaystackPaymentProcessor{}
		return pI, nil
	default:
		return nil, fmt.Errorf("invalid payment type: %s", paymentMethod)
	}

}
