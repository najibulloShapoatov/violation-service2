package repo

import (
	"context"
	"database/sql"
	"fmt"
	"service/pkg/consts"
	"service/pkg/db"
	"service/pkg/log"
	"time"
)

//Subscription ...
type Subscription struct {
	ID           int       `json:"id"`
	PhoneNo      string    `json:"phone_no"`
	Content      string    `json:"content"`
	VehiclePlate string    `json:"vehicle_plate"`
	TarrifID     int       `json:"tarrif_id"`
	DateStart    time.Time `json:"-"`
	DateEnd      time.Time `json:"date_end"`
	Status       int       `json:"status"` // 0-non-active, 1-active, 2-waiting
	ServiceID    int       `json:"service_id"`
}

//CheckSubscribtionLimit ....
func (s *Subscription) CheckSubscribtionLimit(ctx context.Context, limit int) bool {

	var count = 0
	db := db.GetDB()

	stmt := `select count(1) from subscriptions where phone_no = $1 and service_id = $2 and vehicle_plate != $3;`
	if err := db.QueryRowContext(ctx, stmt, s.PhoneNo, s.ServiceID, s.VehiclePlate).Scan(&count); err != nil {
		log.Warn(stmt, s, err)
	}

	if count < limit {
		return true
	}
	return false
}

//Get ...
func (s *Subscription) Get(ctx context.Context) (*Subscription, error) {

	stmt := `select * from subscriptions where phone_no = $1 AND vehicle_plate = $2 AND service_id = $3 AND status = 1;`

	return s.get(ctx, stmt, s.PhoneNo, s.VehiclePlate, s.ServiceID)

}

//GetWithoutService ...
func (s *Subscription) GetWithoutService(ctx context.Context) (*Subscription, error) {

	stmt := `select * from subscriptions where phone_no = $1 AND vehicle_plate = $2 AND status = 1;`

	return s.get(ctx, stmt, s.PhoneNo, s.VehiclePlate)

}

//GetByService ...
func (s *Subscription) GetByService(ctx context.Context) (*Subscription, error) {

	stmt := `select * from subscriptions where phone_no = $1 AND service_id = $2 AND status = 1;`

	return s.get(ctx, stmt, s.PhoneNo, s.ServiceID)
}

//Check ...
func (s *Subscription) Check(ctx context.Context) bool {

	db := db.GetDB()

	stmt := `select * from subscriptions where phone_no = $1 AND vehicle_plate = $2 AND service_id = $3 AND status = 1;`

	if err := db.QueryRowContext(ctx, stmt, s.PhoneNo, s.VehiclePlate, s.ServiceID).Err(); err != nil {
		return false
	}
	return true
}

//Save ...
func (s *Subscription) Save(ctx context.Context, days int) (*Subscription, error) {

	stmt := `INSERT INTO subscriptions ( phone_no, content, vehicle_plate, tarrif_id, date_start, date_end, status, service_id )
			VALUES
				( $1, $2, $3, $4, CURRENT_TIMESTAMP, $5, $6, $7 ) 
			ON CONFLICT ( phone_no, vehicle_plate, service_id) DO
			UPDATE SET tarrif_id = excluded.tarrif_id,
						status = excluded.status,
					date_end = case when subscriptions.date_end > CURRENT_TIMESTAMP  Then
								excluded.date_end + (subscriptions.date_end - CURRENT_TIMESTAMP)
								else 
								excluded.date_end
								end
			RETURNING *;`

	return s.get(ctx, stmt, s.PhoneNo, s.Content, s.VehiclePlate, fmt.Sprint(s.TarrifID), time.Now().AddDate(0, 0, days), s.Status, s.ServiceID)

}

//GetAllActiveByService ...
func (s *Subscription) GetAllActiveByService(ctx context.Context) []*Subscription {

	stmt := `select * from subscriptions where phone_no = $1 and service_id = $2 and status = 1 and date_end > current_timestamp order by date_start DESC;`

	return s.getMore(ctx, stmt, s.PhoneNo, s.ServiceID)
}

//GetAllActive ...
func (s *Subscription) GetAllActive(ctx context.Context) []*Subscription {

	stmt := `select * from subscriptions where phone_no = $1 and status = 1 and date_end > current_timestamp order by date_start DESC;`

	return s.getMore(ctx, stmt, s.PhoneNo)
}

//GetAllActive ...
func (s *Subscription) GetAll(ctx context.Context) []*Subscription {

	stmt := `select * from subscriptions where phone_no = $1 order by date_start DESC;`

	return s.getMore(ctx, stmt, s.PhoneNo)
}

/*
--worker functions
*/

//UpdateStatus ...
func (s *Subscription) UpdateStatus() {

	db := db.GetDB()

	query := `update subscriptions 
	set status = case when subscriptions.date_end < CURRENT_TIMESTAMP Then 0 else 1 end;`

	if _, err := db.Exec(query); err != nil {
		log.Error("Update status error", query, err)
	}
}

//InsertToDailyMessagingQueue ...
func (s *Subscription) InsertToDailyMessagingQueue() {

	db := db.GetDB()

	query := `insert into queues (action, phone_no, vehicle_plate, service_id)
	select 1, phone_no, vehicle_plate, service_id from subscriptions s where s.status = 1 and not exists(
	select * from sent_sms where 
	service_id = s.service_id and status = 1 and phone_no = s.phone_no and
	content = $1 and 
	created_at BETWEEN CURRENT_TIMESTAMP - interval '1' day AND CURRENT_TIMESTAMP + interval '1' day);`

	if _, err := db.Exec(query, consts.SentDaily); err != nil {
		log.Error("Update status error", query, err)
	}
}

//InsertToMothlyMessagingQueue ...
func (s *Subscription) InsertToMothlyMessagingQueue() {

	db := db.GetDB()

	query := `insert into queues (phone_no, vehicle_plate, service_id, action)
	select DISTINCT on (phone_no) phone_no, vehicle_plate, service_id, 2  from subscriptions where status = 1;`

	if _, err := db.Exec(query); err != nil {
		log.Error("Update status error", query, err)
	}
}

//InsertToEndSubscriptionMessagingQueue ...
func (s *Subscription) InsertToEndSubscriptionMessagingQueue() {

	db := db.GetDB()

	query := `insert into queues ( phone_no, vehicle_plate, service_id, action, text)
	select  s.phone_no, s.vehicle_plate, s.service_id, 3 as action,
	COALESCE((select m."content" from messages m where m.service_id=s.service_id and tag='endSubscription'), 
						(select m."content" from messages m where m.service_id=0 and tag='endSubscription')) as text
	from subscriptions  s
	where s.status = 1 and  s.date_end - interval '1' day <= CURRENT_TIMESTAMP and not exists(
	select * from sent_sms where 
	sent_sms.service_id = s.service_id and sent_sms.status = 1 and sent_sms.phone_no = s.phone_no and
	content = concat(s.vehicle_plate,'endSubscriptionSent') and 
	sent_sms.created_at BETWEEN CURRENT_TIMESTAMP - interval '30' day AND CURRENT_TIMESTAMP + interval '1' day)
	order by s.date_end asc;`

	if _, err := db.Exec(query); err != nil {
		log.Error("Update status error", query, err)
	}
}

/*
------------
------------
------------
------------
*/

func (s *Subscription) get(ctx context.Context, query string, args ...interface{}) (*Subscription, error) {

	db := db.GetDB()

	row := db.QueryRowContext(ctx, query, args...)
	err := row.Scan(
		&s.ID,
		&s.PhoneNo,
		&s.Content,
		&s.VehiclePlate,
		&s.TarrifID,
		&s.DateStart,
		&s.DateEnd,
		&s.Status,
		&s.ServiceID,
	)

	if err != nil && err != sql.ErrNoRows {
		log.ErrorDepth(query, 1, err, args)
	}
	//s=subs;
	return s, err

}

func (s *Subscription) getMore(ctx context.Context, query string, args ...interface{}) []*Subscription {

	sbs := make([]*Subscription, 0)

	db := db.GetDB()

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.ErrorDepth(query, 1, err, args)
		return sbs
	}
	defer rows.Close()

	for rows.Next() {
		var sb = &Subscription{}
		if err := rows.Scan(
			&sb.ID,
			&sb.PhoneNo,
			&sb.Content,
			&sb.VehiclePlate,
			&sb.TarrifID,
			&sb.DateStart,
			&sb.DateEnd,
			&sb.Status,
			&sb.ServiceID,
		); err != nil {
			log.WarnDepth("Row Scan", 1, query, err)
		} else {
			sbs = append(sbs, sb)
		}
	}

	return sbs
}
