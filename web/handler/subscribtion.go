package handler

import (
	"net/http"
	"service/web/pkg/response"
	"service/web/pkg/services/subscription"

	"github.com/gorilla/mux"
)

//GetListSubscriptions ...
func GetListSubscriptions(w http.ResponseWriter, r *http.Request) {

	var vars = mux.Vars(r)

	response.JSONWithStatus(w, subscription.GetList(r.Context(),"+992"+vars["phone"], vars["code"]))
	return
}
