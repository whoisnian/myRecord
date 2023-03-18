package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/whoisnian/myRecord/global"
)

type M map[string]any

type Batch[T descriptor] struct {
	conditions []string
	arguments  []any

	pos int
}

func B[T descriptor]() *Batch[T] {
	return &Batch[T]{pos: 1}
}

func (b *Batch[T]) Where(sql string, args ...any) *Batch[T] {
	b.conditions = append(b.conditions, sql)
	b.arguments = append(b.arguments, args...)
	return b
}

func (b *Batch[T]) where() string {
	sb := strings.Builder{}
	if len(b.conditions) > 0 {
		sb.WriteString(" WHERE ")
		for i, condition := range b.conditions {
			if i > 0 {
				sb.WriteString(" and ")
			}
			var last rune
			for _, ch := range condition {
				if ch == rune('?') && last == rune('$') {
					sb.WriteString(strconv.Itoa(b.pos))
					b.pos++
				} else {
					sb.WriteRune(ch)
				}
				last = ch
			}
		}
	}
	return sb.String()
}

func (b *Batch[T]) Create(objs []T) error {
	if len(objs) < 1 {
		return nil
	}

	var sample T
	_, err := global.Pool.CopyFrom(
		context.Background(),
		pgx.Identifier{sample.tableName()},
		sample.fieldsNameActive(),
		pgx.CopyFromSlice(len(objs), func(i int) ([]any, error) {
			return objs[i].fieldsPtrActive(), nil
		}),
	)
	return err
}

func (b *Batch[T]) Find(objsp *[]T) error {
	var sample T

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(sample.fieldsName(), ","),
		sample.tableName(),
	))
	sb.WriteString(b.where())

	rows, err := global.Pool.Query(context.Background(), sb.String(), b.arguments...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var factory T
	for rows.Next() {
		obj := factory.new().(T)
		err = rows.Scan(obj.fieldsPtr()...)
		if err != nil {
			return err
		}
		*objsp = append(*objsp, obj)
	}
	return rows.Err()
}

func (b *Batch[T]) Update(to M) error {
	keys := make([]string, len(to))
	values := make([]any, len(to))

	i := 0
	for k, v := range to {
		keys[i] = k
		values[i] = v
		i++
	}

	var sample T
	sb := strings.Builder{}
	if len(to) == 1 {
		sb.WriteString(fmt.Sprintf("UPDATE %s SET %s = $1",
			sample.tableName(),
			keys[0],
		))
	} else {
		sb.WriteString(fmt.Sprintf("UPDATE %s SET (%s) = (%s)",
			sample.tableName(),
			strings.Join(keys, ","),
			posMark(1, len(to)),
		))
	}
	b.pos = len(to) + 1
	sb.WriteString(b.where())

	_, err := global.Pool.Exec(context.Background(), sb.String(), append(values, b.arguments...)...)
	return err
}

func (b *Batch[T]) Remove() error {
	var sample T
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("DELETE FROM %s", sample.tableName()))
	sb.WriteString(b.where())

	_, err := global.Pool.Exec(context.Background(), sb.String(), b.arguments...)
	return err
}
