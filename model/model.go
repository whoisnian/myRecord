package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/whoisnian/myRecord/global"
)

type pModel interface {
	tableName() string
	fieldsNameAll() []string
	fieldsNamePartial() []string
	fieldsPointerAll() []any
	fieldsPointerPartial() []any
}

// TODO: nSmalls
func placeholders(from, to int) string {
	if to <= from {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteByte('$')
	sb.WriteString(strconv.Itoa(from))
	for from += 1; from <= to; from++ {
		sb.Write([]byte{',', '$'})
		sb.WriteString(strconv.Itoa(from))
	}
	return sb.String()
}

func Create[T pModel](result T) error {
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		result.tableName(),
		strings.Join(result.fieldsNamePartial(), ","),
		placeholders(1, len(result.fieldsNamePartial())),
		strings.Join(result.fieldsNameAll(), ","),
	)

	rows, err := global.Pool.Query(context.Background(), sql, result.fieldsPointerPartial()...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return err
		}
		return pgx.ErrNoRows
	}

	return rows.Scan(result.fieldsPointerAll()...)
}

func Find[T pModel](result T, id int64) error {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1",
		strings.Join(result.fieldsNameAll(), ","),
		result.tableName(),
	)

	rows, err := global.Pool.Query(context.Background(), sql, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return err
		}
		return pgx.ErrNoRows
	}

	return rows.Scan(result.fieldsPointerAll()...)
}
