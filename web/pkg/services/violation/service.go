package violation

import (
	"context"
	"fmt"
	"service/pkg/log"
	"service/pkg/validator"
	"service/web/pkg/interfaces/customer"
	"service/web/pkg/interfaces/patternphone"
	"service/web/pkg/interfaces/subscription"
	"service/web/pkg/interfaces/violation"
	"service/web/pkg/response"
	"service/web/pkg/responsecode"
	"service/web/thirdparty"
	"strconv"
	"time"
)

type VViolation struct {
	BId           string
	VehiclePlate  string
	VTime         *time.Time
	VLocation     string
	VId           string
	VDescription  string
	ProcessStatus int
	PunishStatus  int
	dateCreate    *time.Time
	dateUpdate    *time.Time
	IsPaid        int
	IsPublished   int
	Images        []string
}

//GetList .....
func GetList(ctx context.Context, phone string, code string, plateNo string, page, pagesize int, sts, viol, paid int) *response.Response {

	if !validator.ValidatePhone(phone, patternphone.PatternPhone.Get()) {
		return responsecode.ValidationFailed()
	}
	customer, err := customer.New(phone, code, 0).Get(ctx)
	if err != nil {
		log.Warn("customer not found", err)
		return responsecode.ContentNotFound().ChangeMsg("customer not found")
	}

	if !validator.ValidVehiclePlate(plateNo) {
		return responsecode.ValidationFailed()
	}

	subscription, err := subscription.New(customer.PhoneNo, plateNo, 0).GetWithoutService(ctx)

	if err != nil {
		log.Warn("subscribtion not found", err)
		return responsecode.ContentNotFound().ChangeMsg("subscribtion not found")
	}

	if time.Now().Unix() > subscription.DateEnd.Unix() {
		return responsecode.Expired().ChangeMsg(fmt.Sprintf("subscription expired in %v ", subscription.DateEnd.Format(" 15:04:05  02.01.2006")))
	}

	var res = responsecode.Ok()
	if page <= 1 {
		res.Data["subscription"] = subscription
		res.Data["violationQnt"] = GetViolationQnt(ctx, plateNo)
	}
	res.Data["violations"] = violation.New(subscription.VehiclePlate, "").GetViolations(ctx, page, pagesize, paid, sts, viol)

	return res
}



//GetListExternal .....
func GetListExternal(ctx context.Context, serviceID int, phone string, plateNo string, page, pagesize int, sts, viol, paid int) *response.Response {

	if !validator.ValidatePhone(phone, patternphone.PatternPhone.Get()) {
		return responsecode.ValidationFailed()
	}
	customer, err := customer.New(phone, "", serviceID).GetByPhoneAndService(ctx)
	if err != nil {
		log.Warn("customer not found", err)
		return responsecode.ContentNotFound().ChangeMsg("customer not found")
	}

	if !validator.ValidVehiclePlate(plateNo) {
		return responsecode.ValidationFailed()
	}

	subscription, err := subscription.New(customer.PhoneNo, plateNo, serviceID).GetByService(ctx)

	if err != nil {
		log.Warn("subscribtion not found", err)
		return responsecode.ContentNotFound().ChangeMsg("subscribtion not found")
	}

	if time.Now().Unix() > subscription.DateEnd.Unix() {
		return responsecode.Expired().ChangeMsg(fmt.Sprintf("subscription expired in %v ", subscription.DateEnd.Format(" 15:04:05  02.01.2006")))
	}

	var res = responsecode.Ok()
	//if page <= 1 {
	res.Data["subscription"] = subscription
	res.Data["violationQnt"] = GetViolationQnt(ctx, plateNo)
	//}
	res.Data["violations"] = violation.New(subscription.VehiclePlate, "").GetViolations(ctx, page, pagesize, paid, sts, viol)

	return res
}

//GetViolationQnt ...
func GetViolationQnt(ctx context.Context, plateNo string) map[string]string {
	violations := violation.New(plateNo, "").GetAll(ctx)
	totalQnt := len(violations)
	approvedQnt := 0
	paidQnt := 0
	processQnt := 0
	for _, i := range violations {
		// is paid
		if i.IsPaid == 1 {
			paidQnt++
		}
		// approved
		if i.IsPaid == 0 && i.ProcessStatus == 1 {
			approvedQnt++
		}
		// process status
		if i.IsPaid == 0 && i.ProcessStatus != 1 {
			processQnt++
		}
	}
	res := make(map[string]string)

	res["total_qnt"] = strconv.Itoa(totalQnt)
	res["approved_qnt"] = strconv.Itoa(approvedQnt)
	res["paid_qnt"] = strconv.Itoa(paidQnt)
	res["process_qnt"] = strconv.Itoa(processQnt)
	return res
}


func GetExternal(ctx context.Context, serviceID int, phone string, plateNo string, bId string) *response.Response {
	if !validator.ValidatePhone(phone, patternphone.PatternPhone.Get()) {
		return responsecode.ValidationFailed()
	}
	customer, err := customer.New(phone, "", serviceID).GetByPhoneAndService(ctx)
	if err != nil {
		log.Warn("customer not found", err)
		return responsecode.ContentNotFound().ChangeMsg("customer not found")
	}

	if !validator.ValidVehiclePlate(plateNo) {
		return responsecode.ValidationFailed().ChangeMsg("invalid vehicle plate")
	}

	subscription, err := subscription.New(customer.PhoneNo, plateNo, serviceID).GetByService(ctx)

	if err != nil {
		log.Warn("subscribtion not found", err)
		return responsecode.ContentNotFound().ChangeMsg("subscribtion not found")
	}

	if time.Now().Unix() > subscription.DateEnd.Unix() {
		return responsecode.Expired().ChangeMsg(fmt.Sprintf("subscription expired in %v ", subscription.DateEnd.Format(" 15:04:05  02.01.2006")))
	}

	violation, err := violation.Violation.Get(ctx, bId)
	if err != nil {
		log.Warn("Violation Get Error", err)
		return responsecode.NotFound().ChangeMsg("violation not found")
	}
	var vviolation = &VViolation{
		BId:           violation.BId,
		VehiclePlate:  violation.VehiclePlate,
		VTime:         violation.VTime,
		VLocation:     violation.VLocation,
		VId:           violation.VId,
		VDescription:  violation.VDescription.String,
		ProcessStatus: violation.ProcessStatus,
		PunishStatus:  violation.PunishStatus,
		IsPaid:        violation.IsPaid,
		IsPublished:   violation.IsPublished,
	}
	vviolation.Images = thirdparty.GetPhotoLinks(vviolation.BId)

	var res = responsecode.Ok()

	res.Data["violation"] = vviolation

	return res
}
