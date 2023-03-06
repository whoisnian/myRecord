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

// TODO: auto generate
func (*Item) tableName() string {
	return "items"
}
func (*Item) fieldsNameAll() []string {
	return []string{"id", "type", "state", "content", "date", "created_at"}
}
func (*Item) fieldsNamePartial() []string {
	return []string{"type", "state", "content", "date"}
}
func (item *Item) fieldsPointerAll() []any {
	return []any{&item.Id, &item.Type, &item.State, &item.Content, &item.Date, &item.CreatedAt}
}
func (item *Item) fieldsPointerPartial() []any {
	return []any{&item.Type, &item.State, &item.Content, &item.Date}
}
