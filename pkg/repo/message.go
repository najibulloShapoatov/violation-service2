package repo

import (
	"service/pkg/db"
	"service/pkg/log"
)

//Message ...
type Message struct {
	ID        int
	Tag       string
	ServiceID int
	Content   string
}

//GetMessageEndSubscription ...
func (m *Message) GetMessageEndSubscription(serviceID int) string {

	m, err := m.GetByTagAndService("endSubscription", serviceID)
	if err != nil {
		m, err = m.GetByTagAndService("endSubscription", 0)
	}

	return m.Content
}

//GetByTagAndService ...
func (m *Message) GetByTagAndService(tagName string, serviceID int) (*Message, error) {

	query := `select * from messages where tag = $1 and service_id = $2;`

	return m.get(query, tagName, serviceID)
}

/*
------
------
------
------
------
------
*/

func (m *Message) get(query string, args ...interface{}) (*Message, error) {
	db := db.GetDB()
	

	err := db.QueryRow(query, args...).Scan(
		&m.ID,
		&m.Tag,
		&m.ServiceID,
		&m.Content,
	)
	if err != nil {
		log.ErrorDepth(query, 1, err)
		return nil, err
	}
	return m, nil
}
