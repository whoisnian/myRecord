package initializer

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/myRecord/global"
)

func SetupPostgres() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(global.CFG.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}

	if logger.IsDebug() {
		config.ConnConfig.Tracer = global.PoolTracer
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Test postgresql ping...")
	err = pool.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	return pool
}
