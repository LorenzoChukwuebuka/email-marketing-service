package paymentmethods

import "fmt"

type CreditCard struct {
    paid bool
}


func (c *CreditCard) Pay(amount float64) {
    fmt.Printf("Paid $%.2f using Credit Card\n", amount)
    c.paid = true
}

func (c *CreditCard) Charge(amount float64) {
    fmt.Printf("Charged $%.2f to Credit Card\n", amount)
}

func (c *CreditCard) Refund(amount float64) {
    fmt.Printf("Refunded $%.2f to Credit Card\n", amount)
}

func (c *CreditCard) Status() string {
    if c.paid {
        return "Paid"
    }
    return "Unpaid"
}
