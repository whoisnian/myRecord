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

type Batch[T Descriptor] struct {
	conditions []string
	arguments  []any

	pos int
}

// Usage: model.B[*item.Item]()
func B[T Descriptor]() *Batch[T] {
	return &Batch[T]{pos: 1}
}

// Usage: model.B[*item.Item]().Where("id = $?", 1)
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

// Usage: model.B[*item.Item]().Create([]*item.Item{})
func (b *Batch[T]) Create(objs []T) error {
	if len(objs) < 1 {
		return nil
	}

	var sample T
	_, err := global.Pool.CopyFrom(
		context.Background(),
		pgx.Identifier{sample.TableName()},
		sample.FieldsNameActive(),
		pgx.CopyFromSlice(len(objs), func(i int) ([]any, error) {
			return objs[i].FieldsPtrActive(), nil
		}),
	)
	return err
}

// Usage: model.B[*item.Item]().Find(&[]*item.Item{})
func (b *Batch[T]) Find(objsp *[]T) error {
	var sample T

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(sample.FieldsName(), ","),
		sample.TableName(),
	))
	sb.WriteString(b.where())

	rows, err := global.Pool.Query(context.Background(), sb.String(), b.arguments...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var factory T
	for rows.Next() {
		obj := factory.New().(T)
		err = rows.Scan(obj.FieldsPtr()...)
		if err != nil {
			return err
		}
		*objsp = append(*objsp, obj)
	}
	return rows.Err()
}

// Usage: model.B[*item.Item]().Update(model.M{})
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
			sample.TableName(),
			keys[0],
		))
	} else {
		sb.WriteString(fmt.Sprintf("UPDATE %s SET (%s) = (%s)",
			sample.TableName(),
			strings.Join(keys, ","),
			posMark(1, len(to)),
		))
	}
	b.pos = len(to) + 1
	sb.WriteString(b.where())

	_, err := global.Pool.Exec(context.Background(), sb.String(), append(values, b.arguments...)...)
	return err
}

// Usage: model.B[*item.Item]().Remove()
func (b *Batch[T]) Remove() error {
	var sample T
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("DELETE FROM %s", sample.TableName()))
	sb.WriteString(b.where())

	_, err := global.Pool.Exec(context.Background(), sb.String(), b.arguments...)
	return err
}
