package subscription

import (
	"context"
	"service/pkg/validator"
	"service/web/pkg/interfaces/customer"
	"service/web/pkg/interfaces/patternphone"
	"service/web/pkg/interfaces/service"
	"service/web/pkg/interfaces/subscription"
	"service/web/pkg/response"
	"service/web/pkg/responsecode"
)

//GetList ....
func GetList(ctx context.Context, phone string, smscode string) *response.Response {

	if !validator.ValidatePhone(phone, patternphone.PatternPhone.Get()) {
		return responsecode.ValidationFailed().ChangeMsg("phone not valid")
	}

	if _, err := customer.New(phone, smscode, 0).Get(ctx); err != nil {
		return responsecode.ContentNotFound().ChangeMsg("customer not found")
	}

	res := responsecode.Ok()
	res.Data["subscriptions"] = subscription.New(phone, "", 0).GetAll(ctx)

	res.Data["services"] = service.Service.GetAll(ctx)

	return res
}

//GetList ....
func GetListExternal(ctx context.Context, serviceID int, phone string) *response.Response {

	if !validator.ValidatePhone(phone, patternphone.PatternPhone.Get()) {
		return responsecode.ValidationFailed().ChangeMsg("phone not valid")
	}

	if _, err := customer.New(phone, "", serviceID).GetByPhoneAndService(ctx); err != nil {
		return responsecode.ContentNotFound().ChangeMsg("customer not found")
	}

	res := responsecode.Ok()
	res.Data["subscriptions"] = subscription.New(phone, "", serviceID).GetAllActiveByService(ctx)

	return res
}
