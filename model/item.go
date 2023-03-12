package model

import (
	"time"
)

type Item struct {
	Id        int64     `db:"id"`
	Type      int32     `db:"type"`
	State     int32     `db:"state"`
	Content   string    `db:"content"`
	Date      time.Time `db:"date"`
	CreatedAt time.Time `db:"created_at"`
}

func (*Item) tableName() string {
	return "items"
}
func (*Item) pkeyName() string {
	return "id"
}
func (*Item) fieldsName() []string {
	return []string{"id", "type", "state", "content", "date", "created_at"}
}
func (*Item) fieldsNameActive() []string { // Passive: "id", "created_at"
	return []string{"type", "state", "content", "date"}
}
func (item *Item) pkeyPtr() any {
	return &item.Id
}
func (item *Item) fieldsPtr() []any {
	return []any{&item.Id, &item.Type, &item.State, &item.Content, &item.Date, &item.CreatedAt}
}
func (item *Item) fieldsPtrActive() []any {
	return []any{&item.Type, &item.State, &item.Content, &item.Date}
}
