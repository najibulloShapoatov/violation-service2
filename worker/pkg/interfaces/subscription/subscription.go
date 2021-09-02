package subscription

import "service/pkg/repo"

//ISubscription ....
type ISubscription interface {
	UpdateStatus()
	InsertToDailyMessagingQueue()
	InsertToMothlyMessagingQueue()
	InsertToEndSubscriptionMessagingQueue() 
	
}

//Subscription ...
var Subscription ISubscription = &repo.Subscription{}
