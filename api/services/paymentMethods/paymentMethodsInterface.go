package paymentmethods

type PaymentInterface interface {
    Pay(amount float64)
    Charge(amount float64)
    Refund(amount float64)
    Status() string
}