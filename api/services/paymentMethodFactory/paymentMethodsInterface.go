package paymentmethodFactory

import "email-marketing-service/api/model"

type PaymentInterface interface {
	InitializePaymentProcess(
		d *model.InitPaymentModelData,
	) (map[string]interface{}, error) // Not required.
	Pay(amount float64) // Not required
	Charge(amount float64)
	Refund(amount float64)
	Status() string // You can call this PaymentStatus - could be better to return a bool
}

// type Enum struct {
// 	data []EnumValues
// }
//
// type EnumValues struct {
// 	Key   string
// 	Value string
// }
//
// func NewEnum() *Enum {
// 	return &Enum{}
// }
//
// func (e *Enum) CreateEnums(enums map[string]string) {
// 	for key, value := range enums {
// 		d := EnumValues{Key: key, Value: value}
// 		e.data = append(e.data, d)
// 	}
// }
//
// var p = map[string]string{
// 	"S": "Successful",
// 	"F": "Failed",
// }
// var PaymentStatus = NewEnum().CreateEnums(p)
