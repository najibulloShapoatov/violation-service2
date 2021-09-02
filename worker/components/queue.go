package components

import (
	"context"
	"service/pkg/consts"
	"service/pkg/log"
	"service/tools"
	"service/worker/pkg/interfaces/queue"
	"service/worker/pkg/interfaces/violation"
	//"sync"
)

//ProceedQueue ....
func ProceedQueue(ctx context.Context) {

	//wg := sync.WaitGroup{}

	queues := queue.Queue.GetAll(ctx)
	//wg.Add(len(queues))

	for _, queue := range queues {

		/* go func(ctx context.Context, queue *repo.Queue) {
		defer wg.Done() */
		switch queue.Action {
		case 1:
			sentSMSDaily(ctx, queue.PhoneNo, queue.ServiceID, queue.VehiclePlate, consts.SentDaily)
			break
		case 2:
			sentSmsStatus(ctx, queue.PhoneNo, queue.ServiceID, "")
			break
		case 3:
			sentSmsText(queue.PhoneNo, queue.ServiceID, queue.VehiclePlate+"\n"+queue.Text.String, queue.VehiclePlate+consts.EndSubscription)
			break
		default:
			sentSmsText(queue.PhoneNo, queue.ServiceID, queue.Text.String, "default")
			break
		}
		queue.Delete(ctx)
		//}(ctx, queue)

	}
	//wg.Wait()
}

func sentSMSDaily(ctx context.Context, phone string, serviceID int, vehiclePlate string, content string) {
	if len(violation.Violation.GetDaily(ctx, vehiclePlate)) > 0 {
		sentSmsStatus(ctx, phone, serviceID, content)
	}

}

func sentSmsStatus(ctx context.Context, phone string, serviceID int, content string) {

	_, err := tools.SentSMSStatusFromServiceWithContent(ctx, phone, serviceID, content)
	if err != nil {
		log.ErrorDepth("sent sms status error", 1, err)
	}
}

func sentSmsText(phone string, serviceID int, text string, content string) {

	status := tools.SentSMSTextWithServiceID(phone, serviceID, text, content)
	if status != 1 {
		log.ErrorDepth("sms not sent to client", 1, phone, serviceID, text)
	}
}
