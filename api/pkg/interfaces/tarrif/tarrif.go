package tarrif

import (
	"context"
	"service/pkg/repo"
)

//ITarrif ....
type ITarrif interface {
	GetByID(context.Context,int, int) (*repo.Tarrif, error)
	GetList(context.Context,int) []*repo.Tarrif
}

//Tarrif ...
var Tarrif ITarrif = &repo.Tarrif{}