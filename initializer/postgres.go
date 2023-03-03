package initializer

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/myRecord/global"
)

func SetupPostgres() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), global.CFG.DatabaseURI)
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
