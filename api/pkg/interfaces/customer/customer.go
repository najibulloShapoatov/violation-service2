package customer

import (
	"context"
	"service/pkg/repo"
	"service/pkg/utils"
)

//ICustomer ....
type ICustomer interface {
	Save(context.Context) (*repo.Customer, error)
	GetByPhoneAndService(context.Context) (*repo.Customer, error)
}

//Customer ...
var Customer ICustomer = &repo.Customer{}

//New ...
func New(Phone string, serviceID int) ICustomer {
	return &repo.Customer{
		PhoneNo:   Phone,
		ServiceID: serviceID,
		SmsCode:   utils.RandSeq(6),
	}
}
