package response

import (
	"encoding/json"
	"net/http"
)

//Response ...
type Response struct {
	Status int
	Data   map[string]interface{}
}

//New ....
func New(code int, msg string) *Response {
	return &Response{
		Status: code,
		Data:   map[string]interface{}{"code": code, "msg": msg},
	}
}

//ChangeMsg ...
func (r *Response) ChangeMsg(message string) *Response {
	r.Data["msg"] = message
	return r
}

//JSONWithStatus .....
func JSONWithStatus(w http.ResponseWriter, response *Response) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response.Data)
	return
}
