package services

import "fmt"

import paymentmethods "email-marketing-service/api/services/paymentMethodFactory"

type Transaction struct {
	selectedPaymentMethod paymentmethods.PaymentInterface
}


func (c *Transaction) ChoosePaymentMethod(paymentType paymentmethods.PaymentMethodType) {
    switch paymentType {
    case paymentmethods.CreditCard:
        c.selectedPaymentMethod = &paymentmethods.CreditCardPaymentProcessor{}
    case paymentmethods.FlutterWave:
        c.selectedPaymentMethod = &paymentmethods.FlutterwavePaymentProcessor{}
    case paymentmethods.Paystack:
        c.selectedPaymentMethod = &paymentmethods.PaystackPaymentProcessor{}
    default:
        fmt.Println("Invalid payment type")
    }
}

