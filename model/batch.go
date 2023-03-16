package model

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/whoisnian/myRecord/global"
)

// model.B().Create([]T) error
// model.B().Where(M).Find(*[]T) error
// model.B().Where(M).Update([]T, M) error
// model.B().Where(M).Remove([]T) error

type M map[string]any

type Batch struct {
	conditions []string
	arguments  []any

	pos int
}

func B() *Batch {
	return &Batch{pos: 1}
}

func (b *Batch) Where(sql string, args ...any) *Batch {
	b.conditions = append(b.conditions, sql)
	b.arguments = append(b.arguments, args...)
	return b
}

func (b *Batch) where() string {
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

func (b *Batch) Create(objs any) error {
	sVal := reflect.ValueOf(objs)
	if sVal.Kind() != reflect.Slice {
		return errors.New("model: Batch.Create() want slice as input argument, but got " + sVal.Kind().String())
	}
	if sVal.Len() < 1 {
		return nil
	}

	objType := sVal.Type().Elem()
	sample, ok := reflect.New(objType).Interface().(descriptor)
	if !ok {
		return errors.New("model: Batch.Create() pointer to " + objType.String() + " does not implement model.descriptor")
	}
	_, err := global.Pool.CopyFrom(
		context.Background(),
		pgx.Identifier{sample.tableName()},
		sample.fieldsNameActive(),
		pgx.CopyFromSlice(sVal.Len(), func(i int) ([]any, error) {
			return sVal.Index(i).Addr().Interface().(descriptor).fieldsPtrActive(), nil
		}),
	)
	return err
}

func (b *Batch) Find(objsp any) error {
	spVal := reflect.ValueOf(objsp)
	if spVal.Kind() != reflect.Pointer {
		return errors.New("model: Batch.Find() want pointer as input argument, but got " + spVal.Kind().String())
	}
	sVal := spVal.Elem()
	if sVal.Kind() != reflect.Slice {
		return errors.New("model: Batch.Find() want pointer to slice, but got pointer to " + sVal.Kind().String())
	}

	objType := sVal.Type().Elem()
	sample, ok := reflect.New(objType).Interface().(descriptor)
	if !ok {
		return errors.New("model: Batch.Find() pointer to " + objType.String() + " does not implement model.descriptor")
	}

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

	for rows.Next() {
		v := reflect.New(objType)
		err = rows.Scan(v.Interface().(descriptor).fieldsPtr()...)
		if err != nil {
			return err
		}
		sVal.Set(reflect.Append(sVal, v.Elem()))
	}
	return rows.Err()
}

func (b *Batch) Update(obj descriptor, to M) error {
	keys := make([]string, len(to))
	values := make([]any, len(to))

	i := 0
	for k, v := range to {
		keys[i] = k
		values[i] = v
		i++
	}

	sb := strings.Builder{}
	if len(to) == 1 {
		sb.WriteString(fmt.Sprintf("UPDATE %s SET %s = $1",
			obj.tableName(),
			keys[0],
		))
	} else {
		sb.WriteString(fmt.Sprintf("UPDATE %s SET (%s) = (%s)",
			obj.tableName(),
			strings.Join(keys, ","),
			posMark(1, len(to)),
		))
	}
	b.pos = len(to) + 1
	sb.WriteString(b.where())

	_, err := global.Pool.Exec(context.Background(), sb.String(), append(values, b.arguments...)...)
	return err
}

func (b *Batch) Remove(obj descriptor) error {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("DELETE FROM %s", obj.tableName()))
	sb.WriteString(b.where())

	_, err := global.Pool.Exec(context.Background(), sb.String(), b.arguments...)
	return err
}
