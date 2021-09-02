package service

import (
	"context"
	"service/pkg/repo"
)

//IService ...
type IService interface {
	GetByID(context.Context, int) (*repo.Service, error)
	GetAll(context.Context) []*repo.Service
	GetByAPIKey(context.Context, string) (*repo.Service, error)
}

//Service service
var Service IService = &repo.Service{}

//New ...
func New(serviceID int) IService {
	return &repo.Service{
		ID: serviceID,
	}
}
