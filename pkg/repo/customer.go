package repo

import (
	"context"
	"service/pkg/db"
	"service/pkg/log"
	"time"
)

//Customer ...
type Customer struct {
	ID         uint64    `json:"id"`
	PhoneNo    string    `json:"phone_no"`
	SmsCode    string    `json:"-"`
	Status     int       `json:"status"` // 0-non-active, 1-active, 2 - waiting activation, 3 - waiting deactivation
	DateCreate time.Time `json:"date_create"`
	ServiceID  int       `json:"service_id"`
	Content    string    `json:"content"`
}

//Save ...
func (c *Customer) Save(ctx context.Context) (*Customer, error) {

	stmt := `insert into customers (phone_no, sms_code, status, service_id, content) 
			values($1, $2, $3, $4, $5) 
			on conflict (phone_no, service_id) 
			do 
			update set phone_no=excluded.phone_no
			returning *;`

	return c.get(ctx, stmt, c.PhoneNo, c.SmsCode, c.Status, c.ServiceID, c.Content)

}

//Get ...
func (c *Customer) Get(ctx context.Context) (*Customer, error) {

	stmt := `select * from customers where phone_no = $1 and sms_code = $2`
	return c.get(ctx, stmt, c.PhoneNo, c.SmsCode)
}

//GetByPhoneAndService ...
func (c *Customer) GetByPhoneAndService(ctx context.Context) (*Customer, error) {

	stmt := `select * from customers where phone_no = $1 and service_id = $2;`
	return c.get(ctx, stmt, c.PhoneNo, c.ServiceID)

}

/*
################
################
################
################
################
*/

func (c *Customer) get(ctx context.Context, query string, args ...interface{}) (*Customer, error) {
	db := db.GetDB()

	err := db.QueryRowContext(ctx, query, args...).Scan(
		&c.ID,
		&c.PhoneNo,
		&c.SmsCode,
		&c.Status,
		&c.DateCreate,
		&c.ServiceID,
		&c.Content,
	)
	if err != nil {
		log.Warn(query, err)
	}
	return c, err

}
