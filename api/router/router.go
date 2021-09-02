package router

import (
	"encoding/json"
	"net/http"
	"service/api/handler"
	"service/api/middleware"

	"github.com/gorilla/mux"
)

//Init ....
func Init() *mux.Router {

	var router = mux.NewRouter()
	router.Use(middleware.LoggingHTTP)
	router.Use(middleware.Recovery)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("API v2")
	})

	service := router.PathPrefix("/api").Subrouter()
	service.Use(middleware.AuthService)

	service.HandleFunc("/tarrifs", handler.GetTarrifs)
	service.HandleFunc("/subscribe", handler.Subscribe)
	service.HandleFunc("/send-sms", handler.SendSMSFromService)
	service.HandleFunc("/customer", handler.GetCustomerService)
	service.HandleFunc("/customer/html", handler.GetHTMLCustomerService)

	//router.HandleFunc("/api/refresh-key", handler.RefreshAPIKey)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("404 | Not found")
	})
	return router
}
