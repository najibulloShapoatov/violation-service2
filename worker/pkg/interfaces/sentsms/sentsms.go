package sentsms

import (
	"context"
	"service/pkg/repo"
)

//ISentSms ...
type ISentSms interface{
	CheckEndingSentSms(context.Context, *repo.Subscription)bool
}

// SentSms ....
var SentSms ISentSms = &repo.SentSms{}