package main

import (
	"net/http"

	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/myRecord/global"
	"github.com/whoisnian/myRecord/initializer"
	"github.com/whoisnian/myRecord/route"
)

func main() {
	global.CFG = initializer.SetupConfig()
	logger.Info("Config: ", global.CFG.Json())

	global.Pool = initializer.SetupPostgres()
	defer global.Pool.Close()
	logger.Info("Connect to postgresql successfully")

	if global.CFG.CreateTable {
		initializer.ApplySchema()
		logger.Info("Apply schema successfully")
		return
	}

	go func() {
		mux := route.SetupRouter()
		logger.Info("Service httpd started: <http://", global.CFG.ListenAddr, ">")
		if err := http.ListenAndServe(global.CFG.ListenAddr, logger.Req(logger.Recovery(mux))); err != nil {
			logger.Fatal(err)
		}
	}()

	osutil.WaitForInterrupt()
}
