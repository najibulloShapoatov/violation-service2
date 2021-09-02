package tools

import (
	"context"
	"service/pkg/log"
	"service/pkg/repo"
	"service/pkg/smssender"
	"strconv"
	"time"
)

var newLine = "\n"

/*---------------------
"1302": "Красный цвет",
"1625": "Красный цвет",
"1230": "Сплошная линия",
"1301": "Против движения",
"1345": "Стоп-линия"
*/

//SentSMSStatusFromServiceWithContent ...
func SentSMSStatusFromServiceWithContent(ctx context.Context, phone string, serviceID int, content string) (int, error) {
	customer, err := (&repo.Customer{PhoneNo: phone, ServiceID: serviceID}).GetByPhoneAndService(ctx)
	if err != nil {
		return 0, err
	}
	subscription, err := (&repo.Subscription{PhoneNo: customer.PhoneNo, ServiceID: customer.ServiceID}).GetByService(ctx)
	if err != nil {
		return 0, err
	}
	detailsLink, err := (&repo.Settings{}).Get(ctx, "DEFAULT_LINK")
	if err != nil {
		log.Error("DEFAULT_LINK ", err)
	}
	sentsms := &repo.SentSms{
		ServiceID: customer.ServiceID,
		PhoneNo:   customer.PhoneNo,
		Status:    0,
		SmsText:   getStatusSmsText(ctx, subscription, customer.SmsCode, detailsLink, 1),
		BID:       repo.NullString{},
		Content:   repo.NullString{String: content, Valid: true},
	}
	sentsms.CreatedAt.Scan(time.Now())
	log.Info("sms text > \n ", sentsms.SmsText)

	sentsms.Status = smssender.SentSms(sentsms.PhoneNo, sentsms.SmsText)

	sentsms.Create()
	return sentsms.Status, nil
}

//SentSMSFromServcie ...
func SentSMSFromServcie(ctx context.Context, phone string, serviceID int) (int, error) {

	customer, err := (&repo.Customer{PhoneNo: phone, ServiceID: serviceID}).GetByPhoneAndService(ctx)
	if err != nil {
		return 0, err
	}
	log.Info(" (&repo.Customer{PhoneNo: phone, ServiceID: serviceID}).GetByPhoneAndService()", err, customer)

	subscription, err := (&repo.Subscription{PhoneNo: customer.PhoneNo, ServiceID: customer.ServiceID}).GetByService(ctx)
	if err != nil {
		return 0, err
	}
	log.Info(" (&repo.Subscription{PhoneNo: customer.PhoneNo, ServiceID: customer.ServiceID}).GetByService()", err, customer)

	return SentSMSStatus(ctx, customer, subscription, 1), nil

}

//SentSMSStatus ...
func SentSMSStatus(ctx context.Context, customer *repo.Customer, subscription *repo.Subscription, withImage int) int {

	detailsLink, err := (&repo.Settings{}).Get(ctx, "DEFAULT_LINK")
	if err != nil {
		log.Error("DEFAULT_LINK ", err)
	}

	sentsms := &repo.SentSms{
		ServiceID: customer.ServiceID,
		PhoneNo:   customer.PhoneNo,
		Status:    0,
		SmsText:   getStatusSmsText(ctx, subscription, customer.SmsCode, detailsLink, withImage),
		BID:       repo.NullString{},
		Content:   repo.NullString{String: " ", Valid: true},
	}
	sentsms.CreatedAt.Scan(time.Now())
	log.Info("sms text > \n ", sentsms.SmsText)

	sentsms.Status = smssender.SentSms(sentsms.PhoneNo, sentsms.SmsText)

	sentsms.Create()
	return sentsms.Status
}

//SentSMSTextWithServiceID ...
func SentSMSTextWithServiceID(phone string, serviceID int, text string, content string) int {

	sentsms := &repo.SentSms{
		ServiceID: serviceID,
		PhoneNo:   phone,
		Status:    0,
		SmsText:   text,
		BID:       repo.NullString{},
		Content:   repo.NullString{String: content, Valid: true},
		CreatedAt: repo.NullTime{Time: time.Now(), Valid: true},
	}

	log.Info("sms text > \n ", sentsms.SmsText)

	sentsms.Status = smssender.SentSms(sentsms.PhoneNo, sentsms.SmsText)

	sentsms.Create()
	return sentsms.Status
}

//GetStatusSmsText returns string
func getStatusSmsText(ctx context.Context, subscribption *repo.Subscription, code, detailslink string, WithImage int) string {
	smsText := ""
	subs := subscribption.GetAllActiveByService(ctx)
	if len(subs) == 1 {
		smsText += getForOneSubscribtion(ctx, subscribption)
		if WithImage == 1 {
			smsText += "Подробно:" + detailslink + subscribption.PhoneNo[4:] + "/" + code + "/" + subscribption.VehiclePlate
		}
		return smsText
	} else if len(subs) == 0 {
		return "У вас нет активных подписок !\nПодробно: " + detailslink + subscribption.PhoneNo[4:] + "/" + code
	}
	titles := []string{"подписка", "подписки", "подписок"}
	smsText += "У вас " + strconv.Itoa(len(subs)) + " " + skloneniya(len(subs), titles) + "." + newLine
	for _, elem := range subs {
		smsText += elem.VehiclePlate + newLine
	}
	smsText += "Подробно: " + detailslink + subscribption.PhoneNo[4:] + "/" + code
	return smsText
}

//GetForOneSubscribtion string
func getForOneSubscribtion(ctx context.Context, s *repo.Subscription) string {
	smsText := ""
	//detail subscribtion
	smsText += s.VehiclePlate + newLine
	violations := (&repo.Violation{VehiclePlate: s.VehiclePlate}).GetAll(ctx)
	if len(violations) == 0 {
		smsText += "На Вашей автомашине штрафов не обнаружено."
		return smsText
	}
	//smsText += "Всего-" + strconv.Itoa(count) + newLine
	processed := 0
	unprocessed := 0
	//paid := 0
	total := 0
	var kr int //Красный цвет
	var cl int //Сплошная линия
	var pd int //Против движения
	var sl int //Стоп-линия
	for _, i := range violations {
		if i.IsPaid == 0 {
			total++
			if i.ProcessStatus == 1 {
				processed++
				switch i.VId {
				case "1302":
					kr++
				case "1625":
					kr++
				case "1230":
					cl++
				case "1301":
					pd++
				case "1345":
					sl++
				}
			}
			if i.ProcessStatus != 1 {
				unprocessed++
			}
		}
	}
	smsText += "Всего-" + strconv.Itoa(total) + newLine

	if processed > 0 {
		smsText += "Подтвержденные-" + strconv.Itoa(processed) + newLine
	}
	if unprocessed > 0 {
		smsText += "В процессе проверки-" + strconv.Itoa(unprocessed) + newLine
	}

	if kr > 0 || cl > 0 || pd > 0 || sl > 0 {
		smsText += "Из них:" + newLine
	}
	if kr > 0 {
		smsText += "Красный цвет-" + strconv.Itoa(kr) + newLine
	}
	if cl > 0 {
		smsText += "Сплошная линия-" + strconv.Itoa(cl) + newLine
	}
	if pd > 0 {
		smsText += "Против движения-" + strconv.Itoa(pd) + newLine
	}
	if sl > 0 {
		smsText += "Стоп-линия-" + strconv.Itoa(sl) + newLine
	}
	return smsText
}

func skloneniya(number int, titles []string) string {
	cases := []int{2, 0, 1, 1, 1, 2}
	var currentCase int
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	return titles[currentCase]
}
