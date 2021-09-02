package response

import (
	"encoding/json"
	"net/http"
)

//JSON func
func JSON(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	return
}


