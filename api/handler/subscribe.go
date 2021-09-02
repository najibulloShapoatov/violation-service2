package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"service/api/middleware"
	"service/api/pkg/responsecode"
	"service/api/pkg/services/subscription"
	"service/pkg/log"
	"service/pkg/response"

	"strings"
)

var newLine = "\n"

//Subscribe ... service
func Subscribe(w http.ResponseWriter, r *http.Request) {

	serviceID, err := middleware.Authentication(r.Context())
	if err != nil {
		response.JSON(w, responsecode.Unauthorized())
		return
	}

	buf, _ := ioutil.ReadAll(r.Body)
	log.Info("\n ### Requested to Subscribe -> req Body:\t", string(buf))

	var requestSubscribtion = &subscription.RequestSubscription{}
	err = json.NewDecoder(strings.NewReader(string(buf))).Decode(requestSubscribtion)
	if err != nil {
		response.JSON(w, responsecode.BadRequest())
		log.Error("invalid request", 104, err)
		return
	}

	requestSubscribtion.VehiclePlate = strings.TrimSpace(strings.ToUpper(requestSubscribtion.VehiclePlate))
	requestSubscribtion.Phone = "+" + requestSubscribtion.Phone

	res := subscription.Subscribe(r.Context(), requestSubscribtion, serviceID)

	response.JSON(w, res)
	return
}
