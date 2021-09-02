package subscription

import (
	"context"
	"fmt"
	"service/api/pkg/interfaces/blacklist"
	"service/api/pkg/interfaces/customer"
	"service/api/pkg/interfaces/patternphone"
	"service/api/pkg/interfaces/payment"
	"service/api/pkg/interfaces/service"
	"service/api/pkg/interfaces/settings"
	"service/api/pkg/interfaces/subscription"
	"service/api/pkg/interfaces/tarrif"
	"service/api/pkg/responsecode"
	"service/pkg/log"
	"service/pkg/validator"
	"service/tools"
	"time"
)

//RequestSubscription ....
type RequestSubscription struct {
	Phone        string `json:"phone"`
	Tarrif       int    `json:"tarrif"`
	VehiclePlate string `json:"vehicle_plate"`
}

//Subscribe ....
func Subscribe(ctx context.Context, requestSubscription *RequestSubscription, serviceID int) map[string]interface{} {

	patterns := patternphone.PatternPhone.Get()

	if !validator.ValidatePhone(requestSubscription.Phone, patterns) {
		return responsecode.NoValidPhone()
	}
	if !validator.ValidVehiclePlate(requestSubscription.VehiclePlate) {
		return responsecode.NoValidVehicle()
	}
	if !blacklist.BlackList.CheckVehicleBlackList(ctx, requestSubscription.VehiclePlate) {
		return responsecode.AccessDenied()
	}

	tarrif, err := tarrif.Tarrif.GetByID(ctx, requestSubscription.Tarrif, serviceID)
	if err != nil {
		log.Error("invalid request", 104, err)
		return responsecode.BadRequest()
	}

	limit := settings.Settings.GetInt(ctx, "DEFAULT_LIMIT")

	Subscription := subscription.New(serviceID, requestSubscription.Phone, requestSubscription.VehiclePlate, requestSubscription.Tarrif)

	if !Subscription.CheckSubscribtionLimit(ctx, limit) {
		log.Warn("#subscription limit expired ->", serviceID, requestSubscription)
		return responsecode.SubscriptionExpired()
	}

	service, err := service.Service.GetByID(ctx, serviceID)
	if err != nil {
		log.Error("service.Service.GetByID(serviceID)", err)
		return responsecode.BadRequest()
	}

	if Subscription.Check(ctx) && service.Prolonged == 0 {
		return responsecode.CustomerAlreadySubscribed()
	}

	customer, err := customer.New(requestSubscription.Phone, service.ID).Save(ctx)
	if err != nil {
		log.Error("customer.New(requestSubscription.Phone, service.ID).Save()", err)
		return responsecode.BadRequest()
	}

	subscription, err := Subscription.Save(ctx, tarrif.Days)
	if err != nil {
		log.Error("Subscription.Save(tarrif.Days)", err)
		return responsecode.BadRequest()
	}

	payment, err := payment.New(requestSubscription.Phone, tarrif.Price, serviceID, time.Now()).Create(ctx)
	if err != nil {
		log.Error("not created payment", err)
		return responsecode.BadRequest()
	}

	log.Info("Successfully subscribed : ", " customer:"+fmt.Sprintf("%+v", customer), "\t", "suscribtion: "+fmt.Sprintf("%+v", subscription))

	go func() {
		gorutineCtx := context.Background()
		tools.SentSMSStatus(gorutineCtx, customer, subscription, tarrif.WithImage)
		gorutineCtx.Done()
	}() 

	res := responsecode.Ok()

	//res["subscription"] = subscription
	res["TRANSACTION_ID"] = payment.ID
	return res
}
