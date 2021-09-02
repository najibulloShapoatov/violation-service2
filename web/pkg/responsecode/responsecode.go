package responsecode

import (
	"net/http"
	"service/web/pkg/response"
)


var (
	//Ok ...
	Ok = func()*response.Response{return response.New(http.StatusOK, http.StatusText(http.StatusOK))}
	//NotFound ....
	NotFound = func()*response.Response{return response.New(http.StatusNotFound, http.StatusText(http.StatusNotFound))}
	//ConnectionError ....
	ConnectionError = func()*response.Response{return response.New(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))}
	//Unauthorized ....
	Unauthorized = func()*response.Response{return response.New(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))}
	//BadRequest ....
	BadRequest = func()*response.Response{return response.New(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))}

	//ValidationFailed ....
	ValidationFailed = func()*response.Response{return response.New(http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))}

	//ContentNotFound ....
	ContentNotFound = func()*response.Response{return response.New(http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable))}

	//Expired ....
	Expired = func()*response.Response{return response.New(419, "Expired")}

	//Forbidden ....
	Forbidden = func()*response.Response{return response.New(http.StatusForbidden, http.StatusText(http.StatusForbidden))}
)

