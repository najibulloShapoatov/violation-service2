package settings

import (
	"context"
	"service/pkg/repo"
)

//ISettings ....
type ISettings interface {
	GetInt(context.Context, string) int
	Get(context.Context,string) (string, error)
}

//Settings ...
var Settings ISettings = &repo.Settings{}
