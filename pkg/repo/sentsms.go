package repo

import (
	"context"
	"service/pkg/db"
	"service/pkg/log"

	"github.com/jackc/pgx/v4"
)

//SentSms ...
type SentSms struct {
	ID        int        `json:"id"`
	ServiceID int        `json:"service_id"`
	PhoneNo   string     `json:"phone_no"`
	Status    int        `json:"status"`
	SmsText   string     `json:"sms_text"`
	CreatedAt NullTime   `json:"created_at"`
	BID       NullString `json:"b_id"`
	Content   NullString `json:"content"`
}

//Create ....
func (s *SentSms) Create() (*SentSms, error) {

	db := db.GetDB()
	
	stmt := `insert into sent_sms (service_id, phone_no, status, sms_text, created_at, b_id, content) 
	values ($1, $2, $3, $4, $5, $6, $7)
	 returning *;`

	err := db.QueryRow(stmt, s.ServiceID, s.PhoneNo, s.Status, s.SmsText, s.CreatedAt.Time, s.BID.String, s.Content.String).Scan(
		&s.ID,
		&s.ServiceID,
		&s.PhoneNo,
		&s.Status,
		&s.SmsText,
		&s.CreatedAt,
		&s.BID,
		&s.Content,
	)

	if err != nil {
		log.Error(stmt, err)
	}

	return s, err
}

//CheckEndingSentSms ..
func (s *SentSms) CheckEndingSentSms(ctx context.Context, subscription *Subscription) bool {

	query := `select * from sent_sms where 
	service_id = $1 and status = 1 and phone_no = $2 and
	content = concat($3,'endSubscriptionSent') and 
	created_at BETWEEN CURRENT_TIMESTAMP - interval '30' day AND CURRENT_TIMESTAMP + interval '1' day;`

	s, err := s.get(ctx, query, subscription.ServiceID, subscription.PhoneNo, subscription.VehiclePlate)
	if err == pgx.ErrNoRows {
		return false
	}
	return true

}

/*
-
-
-
-
*/
func (s *SentSms) get(ctx context.Context, query string, args ...interface{}) (*SentSms, error) {

	db := db.GetDB()
	

	//sentsms := SentSms{}

	err := db.QueryRowContext(ctx,query, args...).Scan(
		&s.ID,
		&s.ServiceID,
		&s.PhoneNo,
		&s.Status,
		&s.SmsText,
		&s.CreatedAt,
		&s.BID,
		&s.Content,
	)
	if err != nil {
		log.ErrorDepth(query, 1, err)
		return nil, err
	}
	return s, nil
}
