package service

import (
	"context"
	"service/pkg/repo"
)

//IService ...
type IService interface {
	GetByID(context.Context, int) (*repo.Service, error)
	GetByAPIKey(context.Context, string) (*repo.Service, error)
	//GetByRefreshAPIKey(string) (*repo.Service, error)
	//Update() (*repo.Service, error)
}

//Service service
var Service IService = &repo.Service{}
