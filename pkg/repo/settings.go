package repo

import (
	"context"
	"service/pkg/db"
	"service/pkg/log"
	"strconv"
	"strings"
)

//Settings ....
type Settings struct {
	Key   string
	Value string
}

//Get ...
func (s *Settings) Get(ctx context.Context, key string) (string, error) {
	key = strings.ToUpper(key)

	db := db.GetDB()
	
	if err := db.QueryRowContext(ctx, `select value from settings where key = $1`, key).Scan(&s.Value); err != nil {
		log.Warn(`select value from settings where key = $1;`, key, err)
		return "", err
	}
	return s.Value, nil
}

//GetInt ...
func (s *Settings) GetInt(ctx context.Context, key string) int {

	str, err := s.Get(ctx, key)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Error("func (s *Settings) GetInt(key string)", err)
		return 0
	}
	return i

}
