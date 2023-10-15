package paymentmethodFactory

type PaymentMethodType int

const (
    Paystack PaymentMethodType = iota
    FlutterWave
    CreditCard
)