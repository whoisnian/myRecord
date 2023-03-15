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
}

func B() *Batch {
	return &Batch{}
}

func (b *Batch) Where(sql string, args ...any) *Batch {
	b.conditions = append(b.conditions, sql)
	b.arguments = append(b.arguments, args...)
	return b
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
	sample, ok := reflect.Zero(objType).Interface().(descriptor)
	if !ok {
		return errors.New("model: Batch.Create() want slice of descriptor, but got slice of " + objType.String())
	}
	_, err := global.Pool.CopyFrom(
		context.Background(),
		pgx.Identifier{sample.tableName()},
		sample.fieldsNameActive(),
		pgx.CopyFromSlice(sVal.Len(), func(i int) ([]any, error) {
			return sVal.Index(i).Interface().(descriptor).fieldsPtrActive(), nil
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
	sample, ok := reflect.Zero(objType).Interface().(descriptor)
	if !ok {
		return errors.New("model: Batch.Find() want pointer to slice of descriptor, but got pointer to slice of " + objType.String())
	}

	sb := strings.Builder{}
	sb.WriteString(
		fmt.Sprintf("SELECT %s FROM %s",
			strings.Join(sample.fieldsName(), ","),
			sample.tableName(),
		),
	)

	pos := 1
	if len(b.conditions) > 0 {
		sb.WriteString(" WHERE ")
		for i, condition := range b.conditions {
			if i > 0 {
				sb.WriteString(" and ")
			}
			var last rune
			for _, ch := range condition {
				if ch == rune('?') && last == rune('$') {
					sb.WriteString(strconv.Itoa(pos))
					pos++
				} else {
					sb.WriteRune(ch)
				}
				last = ch
			}
		}
	}

	rows, err := global.Pool.Query(context.Background(), sb.String(), b.arguments...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		v := reflect.New(objType.Elem()).Interface().(descriptor)
		err = rows.Scan(v.fieldsPtr()...)
		if err != nil {
			return err
		}
		sVal.Set(reflect.Append(sVal, reflect.ValueOf(v)))
	}
	return rows.Err()
}

// TODO
func (b *Batch) Update(objs any, to M) error {
	return nil
}

// TODO
func (b *Batch) Remove(objs any) error {
	return nil
}
