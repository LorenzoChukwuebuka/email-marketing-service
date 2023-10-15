package paymentmethods

import "fmt"

type Paystack struct {
	paid bool
}

func (c *Paystack) Pay(amount float64)  {
	fmt.Printf("Paid $%.2f using paystack Card\n", amount)
	c.paid = true
}

func (c *Paystack) Charge(amount float64) {
	fmt.Printf("Charged $%.2f to Credit Card\n", amount)
}

func (c *Paystack) Refund(amount float64) {
	fmt.Printf("Refunded $%.2f to Credit Card\n", amount)
}

func (c *Paystack) Status() string {
	if c.paid {
		return "Paid"
	}
	return "Unpaid"
}
