package router

import (
	"encoding/json"
	"net/http"
	"service/web/handler"
	"service/web/middleware"

	"github.com/gorilla/mux"
)

//Init ....
func Init() *mux.Router {

	var router = mux.NewRouter()
	router.Use(middleware.LoggingHTTP)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "msg": "Wellcome to API"})
	})

	router.HandleFunc("/web/subscriptions/{phone}/{code}", handler.GetListSubscriptions)
	router.HandleFunc("/web/violations/{phone}/{code}/{plateno}", handler.GetListViolations)

	//for external agents
	service := router.PathPrefix("/web/customer").Subrouter()
	service.Use(middleware.AuthService)
	service.HandleFunc("/{phone}", handler.GetListSubscriptionsExternal)
	service.HandleFunc("/{phone}/{plateno}", handler.GetListViolationsExternal)
	service.HandleFunc("/{phone}/{plateno}/{bId}", handler.GetViolationExternal)
	//end

	/* staticDir := "/web/public/"
	router.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir)))) */

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "msg": http.StatusText(http.StatusNotFound)})
	})
	return router
}
