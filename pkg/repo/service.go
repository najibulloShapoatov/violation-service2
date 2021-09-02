package repo

import (
	"context"
	"service/pkg/db"
	"service/pkg/log"
	"time"
)

//Service model
type Service struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Login         string     `json:"-"`
	Password      string     `json:"-"`
	DateCreate    *time.Time `json:"-"`
	IsActive      int        `json:"-"`
	Prolonged     int        `json:"-"`
	Logo          NullString `json:"logo"`
	APIKey        string     `json:"-"`
	RefreshAPIKey string     `json:"-"`
	AllowedIP     string     `json:"-"`
}

//GetByID ....
func (s *Service) GetByID(ctx context.Context, id int) (*Service, error) {

	query := `select * from services where id = $1 and is_active = 1;`
	return s.get(ctx, query, id)

}

//GetByAPIKey ....
func (s *Service) GetByAPIKey(ctx context.Context, APIKey string) (*Service, error) {

	query := `select * from services where api_key = $1 and is_active = 1;`
	return s.get(ctx, query, APIKey)
}

//GetAll ....
func (s *Service) GetAll(ctx context.Context) []*Service {
	stmt := `select * from services where is_active = 1;`
	return s.getAll(ctx, stmt)
}

func (s *Service) get(ctx context.Context, query string, args ...interface{}) (*Service, error) {
	var db = db.GetDB()

	err := db.QueryRowContext(ctx, query, args...).Scan(
		&s.ID,
		&s.Name,
		&s.Login,
		&s.Password,
		&s.DateCreate,
		&s.IsActive,
		&s.Prolonged,
		&s.Logo,
		&s.APIKey,
		&s.RefreshAPIKey,
		&s.AllowedIP,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) getAll(ctx context.Context, query string, args ...interface{}) []*Service {
	var ss = make([]*Service, 0)
	db := db.GetDB()
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Warn(query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var sr = &Service{}
		if err := rows.Scan(
			&sr.ID,
			&sr.Name,
			&sr.Login,
			&sr.Password,
			&sr.DateCreate,
			&sr.IsActive,
			&sr.Prolonged,
			&sr.Logo,
			&sr.APIKey,
			&sr.RefreshAPIKey,
			&s.AllowedIP,
		); err != nil {
			log.Warn("Rows Scan", query, err)
		} else {
			ss = append(ss, sr)
		}
	}

	return ss
}
