package repo

import (
	"service/pkg/db"
	"service/pkg/log"
)

//PatternPhone ...
type PatternPhone struct {
	ID      int
	Pattern string
}

//Get ...
func (p *PatternPhone) Get() []string {

	ps := make([]string, 0)

	db := db.GetDB()
	

	rows, err := db.Query(`select pattern from pattern_phones`)
	if err != nil {
		log.Warn("select pattern from pattern_phones;", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			log.Warn("Row Scan Error", "select pattern from pattern_phones;", err)
		}
		ps = append(ps, s)
	}
	return ps
}
