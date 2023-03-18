package item

import (
	"time"

	"github.com/whoisnian/myRecord/model"
)

type Item struct {
	Id        int64     `db:"id"`
	Type      int32     `db:"type"`
	State     int32     `db:"state"`
	Content   string    `db:"content"`
	Date      time.Time `db:"date"`
	CreatedAt time.Time `db:"created_at"`
}

func (*Item) TableName() string {
	return "items"
}
func (*Item) PkeyName() string {
	return "id"
}
func (*Item) FieldsName() []string {
	return []string{"id", "type", "state", "content", "date", "created_at"}
}
func (*Item) FieldsNameActive() []string { // Passive: "id", "created_at"
	return []string{"type", "state", "content", "date"}
}
func (it *Item) PkeyPtr() any {
	return &it.Id
}
func (it *Item) FieldsPtr() []any {
	return []any{&it.Id, &it.Type, &it.State, &it.Content, &it.Date, &it.CreatedAt}
}
func (it *Item) FieldsPtrActive() []any {
	return []any{&it.Type, &it.State, &it.Content, &it.Date}
}
func (it *Item) New() model.Descriptor {
	return &Item{}
}

const (
	TypeTodo int32 = iota
	TypeHistoryDay
	TypeHistoryWeek
	TypeHistoryMonth
)

const (
	StateDeleted int32 = iota
	StateToBeDone
	StateFinished
	StateAbandoned
)
