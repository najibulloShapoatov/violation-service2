package repo

import (
	"context"
	"database/sql"
	"fmt"
	"service/pkg/db"
	"service/pkg/log"
	"time"
)

//Violation - violations table structure
type Violation struct {
	ID            int
	BId           string
	VehiclePlate  string
	VTime         *time.Time
	VLocation     string
	VId           string
	VDescription  NullString
	ProcessStatus int
	PunishStatus  int
	dateCreate    *time.Time
	dateUpdate    *time.Time
	IsPaid        int
	IsPublished   int
	CreateBy      int
	UpdateBy      int
	FileName      NullString
}

//GetViolations ...
func (v *Violation) GetViolations(ctx context.Context, page, pageSize int, paid, sts, viol int) (vs []*Violation) {

	db := db.GetDBS()

	sqlStatement := "select * from violations where vehicle_plate=$1 "

	if paid != -1 {
		sqlStatement += " and is_paid=" + fmt.Sprint(paid) + " "
	}
	if sts != -1 {
		sqlStatement += " and process_status=" + fmt.Sprint(sts) + " "
	}
	if viol != -1 {
		var viols = []string{fmt.Sprint(viol)}

		if viol == 1625 {
			viols = append(viols, "1302")
		}
		strq := ""
		for _, v := range viols {
			strq += "'" + fmt.Sprint(v) + "' ,"
		}
		sqlStatement += " and v_id in ( " + strq[:len(strq)-1] + ") "

		log.Info("viols", viols)
	}

	sqlStatement += " limit $2 offset $3;"

	rows, err := db.QueryContext(ctx, sqlStatement, v.VehiclePlate, pageSize, (page * pageSize))
	if err != nil {
		log.Error("Error notGetByUserID =>", err)
	}
	defer rows.Close()

	for rows.Next() {
		vl := &Violation{}
		err := rows.Scan(
			&vl.ID,
			&vl.BId,
			&vl.VehiclePlate,
			&vl.VTime,
			&vl.VLocation,
			&vl.VId,
			&vl.VDescription,
			&vl.ProcessStatus,
			&vl.PunishStatus,
			&vl.dateCreate,
			&vl.dateUpdate,
			&vl.CreateBy,
			&vl.UpdateBy,
			&vl.FileName,
			&vl.IsPaid,
			&vl.IsPublished,
		)
		if err != nil {
			log.Error(fmt.Sprint(err))
		}
		vs = append(vs, vl)
	}

	return vs

}

//Get ...
func (v *Violation) Get(ctx context.Context, bID string) (*Violation, error) {

	db := db.GetDBS()

	stmt := "select * from violations where b_id=$1;"

	row := db.QueryRowContext(ctx, stmt, bID)

	err := row.Scan(
		&v.ID,
		&v.BId,
		&v.VehiclePlate,
		&v.VTime,
		&v.VLocation,
		&v.VId,
		&v.VDescription,
		&v.ProcessStatus,
		&v.PunishStatus,
		&v.dateCreate,
		&v.dateUpdate,
		&v.CreateBy,
		&v.UpdateBy,
		&v.FileName,
		&v.IsPaid,
		&v.IsPublished,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return v, nil
		}
		log.Error(fmt.Sprint(err))
		return nil, err
	}
	return v, nil
}

//GetAll ...
func (v *Violation) GetAll(ctx context.Context) []*Violation {
	return v.GetViolations(ctx, 0, 10000, -1, -1, -1)
}

//GetDaily ...
func (v *Violation) GetDaily(ctx context.Context, vehiclePlate string) []*Violation {

	query := `select * from violations where vehicle_plate = $1 AND 
	v_time BETWEEN CURRENT_DATE+interval '8' hour AND CURRENT_TIMESTAMP;`

	return v.getMore(ctx, query, vehiclePlate)
}

func (v *Violation) getMore(ctx context.Context,query string, args ...interface{}) []*Violation {

	db := db.GetDBS()

	var vs = make([]*Violation, 0)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.ErrorDepth(query, 1, err)
	}
	defer rows.Close()

	for rows.Next() {
		vl := &Violation{}
		err := rows.Scan(
			&vl.ID,
			&vl.BId,
			&vl.VehiclePlate,
			&vl.VTime,
			&vl.VLocation,
			&vl.VId,
			&vl.VDescription,
			&vl.ProcessStatus,
			&vl.PunishStatus,
			&vl.dateCreate,
			&vl.dateUpdate,
			&vl.CreateBy,
			&vl.UpdateBy,
			&vl.FileName,
			&vl.IsPaid,
			&vl.IsPublished,
		)
		if err != nil {
			log.Error(fmt.Sprint(err))
		}
		vs = append(vs, vl)
	}

	return vs

}
