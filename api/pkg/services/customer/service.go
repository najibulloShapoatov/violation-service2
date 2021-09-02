package customer

import (
	"context"
	"service/api/pkg/interfaces/customer"
	"service/api/pkg/interfaces/patternphone"
	"service/api/pkg/interfaces/service"
	"service/api/pkg/interfaces/subscription"
	"service/api/pkg/interfaces/tarrif"
	"service/api/pkg/responsecode"
	"service/pkg/repo"
	"service/pkg/validator"
)

//CustomerService ...
type CustomerService struct {
	repo.Customer
	Service       *responseService        `json:"service"`
	Subscriptions []*responseSubscription `json:"subscriptions"`
}

type responseService struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo,omitempty"`
}

type responseSubscription struct {
	repo.Subscription
	Tarrif *responseTarrif `json:"tarrif"`
}
type responseTarrif struct {
	ID        int     `json:"id"`
	ServiceID int     `json:"service_id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	Days      int     `json:"days"`
}

//GetCustomerForService ...
func GetCustomerForService(ctx context.Context,phone string, serviceID int) (*CustomerService, map[string]interface{}) {

	if !validator.ValidatePhone(phone, patternphone.PatternPhone.Get()) {
		return nil, responsecode.NoValidPhone()
	}

	var customerServcie = new(CustomerService)

	customer, err := customer.New(phone, serviceID).GetByPhoneAndService(ctx)
	if err != nil {
		return nil, responsecode.NotFound()
	}

	customerServcie.Customer = *customer

	service, err := service.Service.GetByID(ctx, serviceID)
	if err != nil {
		return nil, responsecode.NotFound()
	}
	customerServcie.Service = &responseService{
		ID:   service.ID,
		Name: service.Name,
		Logo: service.Logo.String,
	}

	responseSubscriptions := make([]*responseSubscription, 0)
	for _, v := range subscription.New(customerServcie.ServiceID, customerServcie.PhoneNo, "", 0).GetAllActiveByService(ctx) {
		responseSubscription := &responseSubscription{
			Subscription: *v,
		}
		tr, _ := tarrif.Tarrif.GetByID(ctx, responseSubscription.TarrifID, responseSubscription.ServiceID)
		responseSubscription.Tarrif = &responseTarrif{
			ID:        tr.ID,
			ServiceID: tr.ServiceID,
			Title:     tr.Title,
			Price:     tr.Price,
			Days:      tr.Days,
		}

		responseSubscriptions = append(responseSubscriptions, responseSubscription)
	}

	customerServcie.Subscriptions = responseSubscriptions

	return customerServcie, nil

}
