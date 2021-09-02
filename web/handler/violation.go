package handler

import (
	"net/http"
	"service/web/pkg/response"
	"service/web/pkg/services/violation"
	"strconv"

	"github.com/gorilla/mux"
)

//GetListViolations ...
func GetListViolations(w http.ResponseWriter, r *http.Request) {

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

	response.JSONWithStatus(w, violation.GetList(r.Context(), "+992"+vars["phone"], vars["code"], vars["plateno"], page, pageSize, sts, viol, paid))
	return

}




