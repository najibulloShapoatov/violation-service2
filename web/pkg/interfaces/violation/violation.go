package violation

import (
	"context"
	"service/pkg/repo"
)

//IViolation ...
type IViolation interface {
	GetViolations(context.Context, int, int, int, int, int) []*repo.Violation
	GetAll(context.Context) []*repo.Violation
	Get(context.Context, string) (*repo.Violation, error)
}

//Violation ....
var Violation IViolation = &repo.Violation{}

//New ...
func New(plateNo string, BID string) IViolation {
	return &repo.Violation{
		VehiclePlate: plateNo,
		BId:          BID,
	}
}
