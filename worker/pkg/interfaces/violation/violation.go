package violation

import (
	"context"
	"service/pkg/repo"
)

//IViolation ...
type IViolation interface {
	GetDaily(context.Context, string) []*repo.Violation
}

//Violation ....
var Violation IViolation = &repo.Violation{}