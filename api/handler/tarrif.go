package handler

import (
	"net/http"
	"service/api/middleware"
	"service/api/pkg/interfaces/tarrif"
	"service/api/pkg/responsecode"
	"service/pkg/response"
)

//GetTarrifs ...
func GetTarrifs(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSON(w, responsecode.Unauthorized())
		return
	}

	resp := responsecode.Ok()
	resp["tarrifs"] = tarrif.Tarrif.GetList(r.Context(), id)
	response.JSON(w, resp)
	return
}
