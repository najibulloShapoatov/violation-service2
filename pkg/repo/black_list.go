package repo

import (
	"context"
	"service/pkg/db"
	"time"
)

//BlackList modal
type BlackList struct {
	ID           int        `json:"id"`
	VehiclePlate string     `json:"vehicle_plate"`
	Note         NullString `json:"note"`
	CreatedAt    time.Time  `json:"created_at"`
}

//CheckVehicleBlackList func
func (b *BlackList) CheckVehicleBlackList(ctx context.Context, vehiclePlate string) bool {

	db := db.GetDB()

	if err := db.QueryRowContext(ctx, `select * from black_lists where vehicle_plate = $1;`, vehiclePlate).Err(); err != nil {
		return false
	}
	return true

}
