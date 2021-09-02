package repo

import (
	"context"
	"service/pkg/db"
	"service/pkg/log"
	"time"
)

//Payment ....
type Payment struct {
	ID        int       `json:"id"`
	PhoneNo   string    `json:"phone_no"`
	Amount    float64   `json:"amount"`
	OperDate  time.Time `json:"oper_date"`
	Status    int       `json:"status"` // 0-non-active, 1-active, 2-waiting
	ServiceID int       `json:"service_id"`
}

//Create ...
func (p *Payment) Create(ctx context.Context) (*Payment, error) {

	stmt := `insert into payments (phone_no, amount, oper_date, status, service_id) values($1, $2, $3, $4, $5) returning *;`
	return p.get(ctx, stmt, p.PhoneNo, p.Amount, p.OperDate, p.Status, p.ServiceID)

}

func (p *Payment) get(ctx context.Context, query string, args ...interface{}) (*Payment, error) {

	db := db.GetDB()

	err := db.QueryRowContext(ctx, query, args...).Scan(
		&p.ID,
		&p.PhoneNo,
		&p.Amount,
		&p.OperDate,
		&p.Status,
		&p.ServiceID,
	)
	if err != nil {
		log.Warn(query, err, p)
	}

	return p, err
}
