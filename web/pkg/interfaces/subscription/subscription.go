package subscription

import (
	"context"
	"service/pkg/repo"
)

//ISubscription ....
type ISubscription interface {
	GetAll(context.Context) []*repo.Subscription
	GetAllActive(context.Context) []*repo.Subscription
	GetAllActiveByService(context.Context) []*repo.Subscription
	GetWithoutService(context.Context) (*repo.Subscription, error)
	GetByService(ctx context.Context) (*repo.Subscription, error)
}

//Subscription ...
var Subscription ISubscription = &repo.Subscription{}

//New ...
func New(phone string, plateNo string, serviceID int) ISubscription {
	return &repo.Subscription{
		ServiceID:    serviceID,
		PhoneNo:      phone,
		VehiclePlate: plateNo,
	}
}
