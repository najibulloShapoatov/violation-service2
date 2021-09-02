package handler

import (
	"net/http"
	"service/api/middleware"
	"service/api/pkg/responsecode"
	"service/api/pkg/services/customer"
	"service/api/pkg/services/service"
	"service/pkg/log"
	"service/pkg/response"
	"text/template"
)

//SendSMSFromService ...
func SendSMSFromService(w http.ResponseWriter, r *http.Request) {

	serviceID, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSON(w, responsecode.Unauthorized())
		return
	}
	phone := "+" + r.URL.Query().Get("phone")

	response.JSON(w, service.SendSMSStatus(r.Context(), phone, serviceID))

}

//GetCustomerService ...
func GetCustomerService(w http.ResponseWriter, r *http.Request) {
	serviceID, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSON(w, responsecode.Unauthorized())
		return
	}
	phone := "+" + r.URL.Query().Get("phone")

	res := responsecode.Ok()
	customer, result := customer.GetCustomerForService(r.Context(), phone, serviceID)
	if result != nil {
		response.JSON(w, result)
	}
	res["customer"] = customer

	response.JSON(w, res)

}
//GetHTMLCustomerService ...
func GetHTMLCustomerService(w http.ResponseWriter, r *http.Request) {
	serviceID, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSON(w, responsecode.Unauthorized())
		return
	}
	
	phone := "+" + r.URL.Query().Get("phone")
	
	customer, result := customer.GetCustomerForService(r.Context(), phone, serviceID)
	if result != nil {
		response.JSON(w, result)
	}
	
	tmpl, err := template.ParseFiles("./api/templates/customer.html")
	if err != nil {
		log.Error("<<<-", err)
	}
	tmpl.Execute(w, *customer)
	return

} 

