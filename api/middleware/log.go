package middleware

import (
	"fmt"
	"net/http"
	"service/pkg/log"
	"time"
)

//LoggingHTTP ...
func LoggingHTTP(nextHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//logging request
		log.Info("request in "+time.Now().Format("2006/01/02 15:04:05")+" =>", fmt.Sprintf("%+v", r))
		nextHandler.ServeHTTP(w, r)
		return
	})
}
