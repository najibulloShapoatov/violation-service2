package payment

import (
	"context"
	"service/pkg/repo"
	"time"
)

//IPayment ....
type IPayment interface {
	Create(context.Context) (*repo.Payment, error)
}

//Payment ....
var Payment IPayment = &repo.Payment{}

//New ....
func New(Phone string, Amount float64, serviceID int, operdate time.Time) IPayment {
	return &repo.Payment{
		PhoneNo:   Phone,
		Amount:    Amount,
		OperDate:  time.Now(),
		Status:    1,
		ServiceID: serviceID,
	}
}
