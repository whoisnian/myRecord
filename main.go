package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whoisnian/glb/config"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/myRecord/global"
	"github.com/whoisnian/myRecord/router"
)

func main() {
	err := config.FromCommandLine(&global.CFG)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetDebug(global.CFG.Debug)
	logger.Info("Config: ", global.CFG)

	global.Pool, err = pgxpool.New(context.Background(), global.CFG.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}
	defer global.Pool.Close()

	err = global.Pool.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Connect to database successfully")

	go func() {
		mux := router.Init()
		logger.Info("Service httpd started: <http://", global.CFG.ListenAddr, ">")
		if err := http.ListenAndServe(global.CFG.ListenAddr, logger.Req(logger.Recovery(mux))); err != nil {
			logger.Fatal(err)
		}
	}()

	osutil.WaitForInterrupt()
}
