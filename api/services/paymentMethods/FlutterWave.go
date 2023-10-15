package paymentmethods

import "fmt"


type Flutterwave struct {
    paid bool
}


func (c *Flutterwave) Pay(amount float64) {
    fmt.Printf("Paid $%.2f using Credit flutterwave\n", amount)
    c.paid = true
}

func (c *Flutterwave) Charge(amount float64) {
    fmt.Printf("Charged $%.2f to Credit Card\n", amount)
}

func (c *Flutterwave) Refund(amount float64) {
    fmt.Printf("Refunded $%.2f to Credit Card\n", amount)
}

func (c *Flutterwave) Status() string {
    if c.paid {
        return "Paid"
    }
    return "Unpaid"
}
