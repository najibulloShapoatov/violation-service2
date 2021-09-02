package components

import (
	"service/worker/pkg/interfaces/subscription"
)

//UpdateStatusSubscription ...
func UpdateStatusSubscription() {
	subscription.Subscription.UpdateStatus()
}

//InsertDailySubscription ...
func InsertDailySubscription() {
	subscription.Subscription.InsertToDailyMessagingQueue()
}

//InsertMonthlySubscription ...
func InsertMonthlySubscription() {
	subscription.Subscription.InsertToMothlyMessagingQueue()
}

//InsertEndingSubscription ...
func InsertEndingSubscription(daysAgo int) {

	subscription.Subscription.InsertToEndSubscriptionMessagingQueue()

}
