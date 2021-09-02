package repo

import (
	"service/pkg/db"
	"service/pkg/log"
	"context"
)

//Tarrif model
type Tarrif struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	ServiceID int        `json:"service_id"`
	Price     float64    `json:"price"`
	Comment   NullString `json:"-"`
	Days      int        `json:"days"`
	IsActive  int        `json:"is_active"`
	WithImage int        `json:"with_image"`
	SortOrder int        `json:"sort_order"`
}

//GetList ...
func (t *Tarrif) GetList(ctx context.Context, serviceID int) []*Tarrif {

	var ts = []*Tarrif{}

	db := db.GetDB()

	rows, err := db.QueryContext(ctx, `select * from tarrifs where service_id = $1;`, serviceID)
	if err != nil {
		log.Warn("select * from tarrifs where service_id = $1;", err)
		return ts
	}
	defer rows.Close()

	for rows.Next() {
		tr := Tarrif{}
		if err := rows.Scan(
			&tr.ID,
			&tr.Title,
			&tr.ServiceID,
			&tr.Price,
			&tr.Comment,
			&tr.Days,
			&tr.IsActive,
			&tr.WithImage,
			&tr.SortOrder,
		); err != nil {
			log.Warn("Row Scan Error", "select * from tarrifs where service_id = $1;", err)
		} else {
			ts = append(ts, &tr)
		}
	}

	return ts
}

//GetByID ....
func (t *Tarrif) GetByID(ctx context.Context, id int, serviceID int) (*Tarrif, error) {

	db := db.GetDB()

	err := db.QueryRowContext(ctx, `select * from tarrifs where id=$1 and service_id=$2;`, id, serviceID).Scan(
		&t.ID,
		&t.Title,
		&t.ServiceID,
		&t.Price,
		&t.Comment,
		&t.Days,
		&t.IsActive,
		&t.WithImage,
		&t.SortOrder,
	)
	if err != nil {
		log.Error("select * from tarrifs where id=$1 and service_id=$2;", id, serviceID, "err ", err)
	}

	return t, err
}
