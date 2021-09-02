package customer

import (
	"context"
	"service/pkg/repo"
)

//ICustomer ....
type ICustomer interface {
	Get(context.Context) (*repo.Customer, error)
	GetByPhoneAndService(context.Context) (*repo.Customer, error)
}

//Customer ...
var Customer ICustomer = &repo.Customer{}

//New ...
func New(Phone string, SmsCode string, serviceID int) ICustomer {
	return &repo.Customer{
		PhoneNo:   Phone,
		SmsCode:   SmsCode,
		ServiceID: serviceID,
	}
}
