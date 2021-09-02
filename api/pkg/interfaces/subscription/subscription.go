package subscription

import (
	"context"
	"service/pkg/repo"
)

//ISubscription ....
type ISubscription interface {
	CheckSubscribtionLimit(context.Context,int) bool
	Check(context.Context) bool
	Save(context.Context,int) (*repo.Subscription, error)
	GetAllActive(context.Context) []*repo.Subscription
	GetAllActiveByService(context.Context) []*repo.Subscription
	GetWithoutService(context.Context) (*repo.Subscription, error)
}

//Subscription ...
var Subscription ISubscription = &repo.Subscription{}

//New ....
func New(serviceID int, Phone string, VehiclePlate string, TarrifID int) ISubscription {
	return &repo.Subscription{
		ServiceID:    serviceID,
		PhoneNo:      Phone,
		VehiclePlate: VehiclePlate,
		TarrifID:     TarrifID,
		Status:       1,
	}
}
