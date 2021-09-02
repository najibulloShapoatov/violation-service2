package responsecode

/*
* 100 - success
* 101 - not found
* 102 - connection error
* 103 - service not found
* 104 - invalid request
* 105 - customer already subscribed
* 106 - invalid phone number
* 107 - customer not found
* 108 - invalid vehicle number
* 109 - subscription is expired
* 110 - subscription not found
 */

var (
	//Ok ...
	Ok = func()map[string]interface{}{return New(100, "ok")}
	//NotFound ....
	NotFound = func()map[string]interface{}{return New(101, "not found")}
	//ConnectionError ....
	ConnectionError = func()map[string]interface{}{return New(102, "connection error")}
	//Unauthorized ....
	Unauthorized = func()map[string]interface{}{return New(103, "service not found")}
	//BadRequest ....
	BadRequest = func()map[string]interface{}{return New(104, "invalid request")}
	//CustomerAlreadySubscribed ....
	CustomerAlreadySubscribed = func()map[string]interface{}{return New(105, "customer already subscribed")}
	//NoValidPhone ....
	NoValidPhone = func()map[string]interface{}{return New(106, "invalid phone number")}
	//CustomerNotFound ....
	CustomerNotFound = func()map[string]interface{}{return New(107, "customer not found")}
	//NoValidVehicle ....
	NoValidVehicle = func()map[string]interface{}{return New(108, "invalid vehicle number")}
	//SubscriptionExpired ....
	SubscriptionExpired = func()map[string]interface{}{return New(109, "subscription limit is expired")}
	//SubscriptionNotFound ....
	SubscriptionNotFound = func()map[string]interface{}{return New(110, "subscription not found")}
	//AccessDenied ....
	AccessDenied = func()map[string]interface{}{return New(112, "access denied for this vehicle")}
)

//New ....
func New(code int, msg string) map[string]interface{} {
	return map[string]interface{}{"code": code, "msg": msg}
}
