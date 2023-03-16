package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/whoisnian/myRecord/global"
)

type descriptor interface {
	tableName() string

	pkeyName() string
	fieldsName() []string
	fieldsNameActive() []string

	pkeyPtr() any
	fieldsPtr() []any
	fieldsPtrActive() []any
}

func Create(obj descriptor) error {
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		obj.tableName(),
		strings.Join(obj.fieldsNameActive(), ","),
		posMark(1, len(obj.fieldsNameActive())),
		strings.Join(obj.fieldsName(), ","),
	)

	row := global.Pool.QueryRow(context.Background(), sql, obj.fieldsPtrActive()...)
	return row.Scan(obj.fieldsPtr()...)
}

func Find(obj descriptor) error {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		strings.Join(obj.fieldsName(), ","),
		obj.tableName(),
		obj.pkeyName(),
	)

	row := global.Pool.QueryRow(context.Background(), sql, obj.pkeyPtr())
	return row.Scan(obj.fieldsPtr()...)
}

func Update(obj descriptor) error {
	cnt := len(obj.fieldsName())
	sql := fmt.Sprintf("UPDATE %s SET (%s) = (%s) WHERE %s = %s RETURNING %s",
		obj.tableName(),
		strings.Join(obj.fieldsName(), ","),
		posMark(1, cnt),
		obj.pkeyName(),
		posMark(cnt+1, cnt+1),
		strings.Join(obj.fieldsName(), ","),
	)

	row := global.Pool.QueryRow(context.Background(), sql, append(obj.fieldsPtr(), obj.pkeyPtr())...)
	return row.Scan(obj.fieldsPtr()...)
}

func Remove(obj descriptor) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = $1",
		obj.tableName(),
		obj.pkeyName(),
	)
	_, err := global.Pool.Exec(context.Background(), sql, obj.pkeyPtr())
	return err
}

const nSmalls = 100

const smallsString = " $0, $1, $2, $3, $4, $5, $6, $7, $8, $9," +
	"$10,$11,$12,$13,$14,$15,$16,$17,$18,$19," +
	"$20,$21,$22,$23,$24,$25,$26,$27,$28,$29," +
	"$30,$31,$32,$33,$34,$35,$36,$37,$38,$39," +
	"$40,$41,$42,$43,$44,$45,$46,$47,$48,$49," +
	"$50,$51,$52,$53,$54,$55,$56,$57,$58,$59," +
	"$60,$61,$62,$63,$64,$65,$66,$67,$68,$69," +
	"$70,$71,$72,$73,$74,$75,$76,$77,$78,$79," +
	"$80,$81,$82,$83,$84,$85,$86,$87,$88,$89," +
	"$90,$91,$92,$93,$94,$95,$96,$97,$98,$99,"

// posMark(1, 3) = "$1, $2, $3"
func posMark(from, to int) string {
	if from <= 0 || from > to {
		return ""
	}
	if to < nSmalls {
		return smallsString[from*4 : to*4+3]
	}

	sb := strings.Builder{}
	sb.WriteByte('$')
	sb.WriteString(strconv.Itoa(from))
	for from += 1; from <= to; from++ {
		sb.WriteString(",$")
		sb.WriteString(strconv.Itoa(from))
	}
	return sb.String()
}
