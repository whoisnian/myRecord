package item

import (
	"context"
	"time"

	"github.com/whoisnian/myRecord/global"
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
	TypeFlag int32 = iota
	TypeHistoryDay
	TypeHistoryWeek
	TypeHistoryMonth
)

const (
	StateDeleted int32 = iota
	StatePending
	StateFinished
)

func (it *Item) Exists() bool {
	sql := "SELECT 1 FROM items WHERE id != $1 AND type = $2 AND state != $3 AND date >= $4 AND date <= $5 LIMIT 1"

	var result int64
	var st, ed time.Time
	if it.Type == TypeHistoryDay {
		st = time.Date(it.Date.Year(), it.Date.Month(), it.Date.Day(), 0, 0, 0, 0, it.Date.Location())
		ed = st.Add(time.Hour*24 - 1)
	} else if it.Type == TypeHistoryWeek {
		weekStart := time.Date(1970, 1, 4, 0, 0, 0, 0, it.Date.Location()) // Sunday
		st = weekStart.Add(it.Date.Sub(weekStart).Truncate(time.Hour * 24 * 7))
		ed = st.Add(time.Hour*24*7 - 1)
	} else if it.Type == TypeHistoryMonth {
		st = time.Date(it.Date.Year(), it.Date.Month(), 1, 0, 0, 0, 0, it.Date.Location())
		ed = st.AddDate(0, 1, 0).Add(-1)
	} else {
		return false
	}
	err := global.Pool.QueryRow(context.Background(), sql, it.Id, it.Type, StateDeleted, st, ed).Scan(&result)
	return err == nil && result == 1
}
