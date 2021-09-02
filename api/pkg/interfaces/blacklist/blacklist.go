package blacklist

import (
	"context"
	"service/pkg/repo"
)

//IBlackList ....
type IBlackList interface {
	CheckVehicleBlackList(context.Context, string) bool
}

//BlackList ...
var BlackList IBlackList = &repo.BlackList{}
