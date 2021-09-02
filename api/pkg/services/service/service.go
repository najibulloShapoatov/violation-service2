package service

import (
	"context"
	"service/api/pkg/interfaces/patternphone"
	"service/api/pkg/responsecode"
	"service/pkg/log"
	"service/pkg/validator"
	"service/tools"
)

//SendSMSStatus ....
func SendSMSStatus(ctx context.Context,phone string, serviceID int) map[string]interface{} {

	if !validator.ValidatePhone(phone, patternphone.PatternPhone.Get()) {
		return responsecode.NoValidPhone()
	}

	status, err := tools.SentSMSFromServcie(ctx, phone, serviceID)
	if err != nil {
		log.Warn("customer or subscription not found tools.SentSMSFromServcie(phone, serviceID)", err)
		return responsecode.NotFound()
	}

	if status != 1 {
		return responsecode.ConnectionError()
	}
	res := responsecode.Ok()

	res["msg"] = "sms status sent to " + phone

	return res
}

/*
//RefreshAPIKey ...
func RefreshAPIKey(service *repo.Service) map[string]interface{} {

	service.APIKey = utils.RandSeq(128)
	service.RefreshAPIKey = utils.RandSeq(128)

	Service = service
	service, err := Service.Update()
	if err != nil {
		log.Error("Service.Update()", err)
		return responsecode.BadRequest
	}
	res := responsecode.Ok

	res["api_key"] = service.APIKey
	res["refresh_api_key"] = service.RefreshAPIKey

	return res
} */
