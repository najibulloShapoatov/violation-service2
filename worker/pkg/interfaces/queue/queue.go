package queue

import (
	"context"
	"service/pkg/repo"
)

//IQueue ...
type IQueue interface {
	Insert(context.Context) (*repo.Queue, error)
	GetAll(context.Context) []*repo.Queue
}

//Queue ...
var Queue IQueue = &repo.Queue{}

//New ....
func New(action int, phone string, vehiclePlate string, serviceID int, text string)IQueue{

	return &repo.Queue{
		Action:       action,
		PhoneNo:      phone,
		VehiclePlate: vehiclePlate,
		ServiceID:    serviceID,
		Text:         repo.NullString{String: text, Valid:true},
	}
}
