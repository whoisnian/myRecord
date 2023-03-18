package global

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whoisnian/glb/logger"
)

var (
	Pool       *pgxpool.Pool
	PoolTracer tracer
)

type ctxKey int

const (
	_ ctxKey = iota
	tracerQueryCtxKey
	tracerCopyFromCtxKey
)

type tracer struct{}

type queryData struct {
	pgx.TraceQueryStartData
	startTime time.Time
}

func (tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, startData pgx.TraceQueryStartData) context.Context {
	return context.WithValue(ctx, tracerQueryCtxKey, &queryData{
		TraceQueryStartData: startData,
		startTime:           time.Now(),
	})
}
func (tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, endData pgx.TraceQueryEndData) {
	data := ctx.Value(tracerQueryCtxKey).(*queryData)
	interval := time.Since(data.startTime)

	result := "OK"
	if endData.Err != nil {
		result = endData.Err.Error()
	}

	logger.Debug("PG Query(", interval, "|", result, "): `", data.SQL, "` with ", data.Args)
}

type copyFromData struct {
	pgx.TraceCopyFromStartData
	startTime time.Time
}

func (tracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, startData pgx.TraceCopyFromStartData) context.Context {
	return context.WithValue(ctx, tracerCopyFromCtxKey, &copyFromData{
		TraceCopyFromStartData: startData,
		startTime:              time.Now(),
	})
}
func (tracer) TraceCopyFromEnd(ctx context.Context, conn *pgx.Conn, endData pgx.TraceCopyFromEndData) {
	data := ctx.Value(tracerCopyFromCtxKey).(*copyFromData)
	interval := time.Since(data.startTime)

	result := "OK"
	if endData.Err != nil {
		result = endData.Err.Error()
	}

	logger.Debug("PG CopyFrom(", interval, "|", result, "): ", data.TableName, " with ", data.ColumnNames)
}
