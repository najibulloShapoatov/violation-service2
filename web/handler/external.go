package handler

import (
	"net/http"
	"service/web/middleware"
	"service/web/pkg/response"
	"service/web/pkg/responsecode"
	"service/web/pkg/services/subscription"
	"service/web/pkg/services/violation"
	"strconv"

	"github.com/gorilla/mux"
)

//GetListSubscriptionsExternal ...
func GetListSubscriptionsExternal(w http.ResponseWriter, r *http.Request) {
	serviceID, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSONWithStatus(w, responsecode.Unauthorized())
		return
	}

	var vars = mux.Vars(r)
	response.JSONWithStatus(w, subscription.GetListExternal(r.Context(), serviceID, "+992"+vars["phone"]))
	return
}

func GetViolationExternal(w http.ResponseWriter, r *http.Request) {
	serviceID, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSONWithStatus(w, responsecode.Unauthorized())
		return
	}
	var vars = mux.Vars(r)
	response.JSONWithStatus(w, violation.GetExternal(r.Context(), serviceID, "+992"+vars["phone"], vars["plateno"], vars["bId"]))
}

//GetListViolationsExternal ...
func GetListViolationsExternal(w http.ResponseWriter, r *http.Request) {

	serviceID, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSONWithStatus(w, responsecode.Unauthorized())
		return
	}
	var vars = mux.Vars(r)

	q := r.URL.Query()

	page := 0
	pageSize := 15
	if s, err := strconv.Atoi(q.Get("page")); err == nil {
		page = s
	}
	if s, err := strconv.Atoi(q.Get("pagesize")); err == nil {
		pageSize = s
	}
	if page > 0 {
		page = page - 1
	}

	paid := -1
	sts := -1
	viol := -1

	if s, err := strconv.Atoi(q.Get("paid")); err == nil {
		paid = s
	}
	if s, err := strconv.Atoi(q.Get("sts")); err == nil {
		sts = s
	}
	if s, err := strconv.Atoi(q.Get("viol")); err == nil {
		viol = s
	}

	response.JSONWithStatus(w, violation.GetListExternal(r.Context(), serviceID, "+992"+vars["phone"], vars["plateno"], page, pageSize, sts, viol, paid))
	return
}
