package repo

import (
	"context"
	"service/pkg/db"
	"service/pkg/log"
)

//Queue model for processing steps
type Queue struct {
	//gorm.Model
	ID           int
	Action       int // 1-send daily sms 2-send Status 3- send end subscribe 4- send rassilka 5- send text
	PhoneNo      string
	VehiclePlate string
	ServiceID    int
	Text         NullString
}

//Insert ....
func (q *Queue) Insert(ctx context.Context) (*Queue, error) {

	query := `INSERT INTO queues( "action", "phone_no", "vehicle_plate", "service_id", "text") 
	VALUES ( $1, $2, $3, $4, $5);`

	return q.get(ctx, query, q.Action, q.PhoneNo, q.VehiclePlate, q.ServiceID, q.Text.String)
}

//Delete ...
func (q *Queue) Delete(ctx context.Context){
	
	db := db.GetDB()
	
	query :=`delete from  queues where id = $1;`

	if _, err := db.ExecContext(ctx, query, q.ID); err != nil{
		log.Error(query, err)
	}
}


//GetAll ...
func (q *Queue) GetAll(ctx context.Context) []*Queue {
	query := `select * from queues order by action asc;`
	return q.getMore(ctx, query)
}

/*
-
-
-
-
*/
func (q *Queue) get(ctx context.Context, query string, args ...interface{}) (*Queue, error) {
	db := db.GetDB()
	

	err := db.QueryRowContext(ctx, query, args...).Scan(
		&q.ID,
		&q.Action,
		&q.PhoneNo,
		&q.VehiclePlate,
		&q.ServiceID,
		&q.Text,
	)
	if err != nil {
		log.ErrorDepth(query, 1, err)
		return nil, err
	}
	return q, nil
}

func (q *Queue) getMore(ctx context.Context, query string, args ...interface{}) []*Queue {

	queues := make([]*Queue, 0)

	db := db.GetDB()
	

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.ErrorDepth(query, 1, err, args)
		return queues
	}
	defer rows.Close()

	for rows.Next() {
		var queue = &Queue{}

		if err := rows.Scan(
			&queue.ID,
			&queue.Action,
			&queue.PhoneNo,
			&queue.VehiclePlate,
			&queue.ServiceID,
			&queue.Text,
		); err != nil {
			log.WarnDepth("Row Scan", 1, query, err)
		} else {
			queues = append(queues, queue)
		}
	}

	return queues
}
