package main

import (
	"net/http"

	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/myRecord/global"
	"github.com/whoisnian/myRecord/initializers"
	"github.com/whoisnian/myRecord/routes"
)

func main() {
	global.CFG = initializers.SetupConfig()
	logger.Info("Config: ", global.CFG.Json())

	global.Pool = initializers.SetupPostgres()
	defer global.Pool.Close()
	logger.Info("Connect to postgresql successfully")

	go func() {
		mux := routes.SetupRouter()
		logger.Info("Service httpd started: <http://", global.CFG.ListenAddr, ">")
		if err := http.ListenAndServe(global.CFG.ListenAddr, logger.Req(logger.Recovery(mux))); err != nil {
			logger.Fatal(err)
		}
	}()

	osutil.WaitForInterrupt()
}
