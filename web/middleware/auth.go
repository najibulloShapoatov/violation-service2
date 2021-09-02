package middleware

import (
	"context"
	"errors"
	"net/http"
	"service/pkg/log"
	"service/web/pkg/response"
	"service/pkg/utils"
	"service/web/pkg/interfaces/service"
	"service/web/pkg/responsecode"
	"strings"
)

//ErrNoAuthentication ...
var ErrNoAuthentication = errors.New("No authentication")
var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

/* func (c *contextKey) String() string {
	return c.name
} */

//AuthService ... service
func AuthService(nextHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		APIKey := r.URL.Query().Get("api_key")
		if len(strings.TrimSpace(APIKey)) != 128 {
			log.Warn(APIKey, len(strings.TrimSpace(APIKey)))
			response.JSONWithStatus(w, responsecode.Unauthorized())
			return
		}

		serviceID, err := CheckServiceAuth(r.Context(), APIKey, r)
		if err != nil {
			log.Warn("CheckServiceAuth(APIKey, r)", err)
			response.JSONWithStatus(w, responsecode.Unauthorized())
			return
		}

		ctx := context.WithValue(r.Context(), authenticationContextKey, serviceID)
		r = r.WithContext(ctx)

		nextHandler.ServeHTTP(w, r)
		return
	})
}

//CheckServiceAuth Function
func CheckServiceAuth(ctx context.Context, APIKey string, r *http.Request) (int, error) {
	service, err := service.Service.GetByAPIKey(ctx, APIKey)
	if err != nil {
		log.Info("error in CheckServiceAuth", err)
		return 0, errors.New("service not found")
	}
	if !utils.CheckIP(utils.Split(service.AllowedIP, ";"), r) {
		return 0, errors.New("ip not allowed")
	}

	return service.ID, nil
}

//Authentication returned authenticated serviceID
func Authentication(ctx context.Context) (int, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}
