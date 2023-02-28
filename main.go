package main

import (
	"net/http"

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

	go func() {
		mux := router.Init()
		logger.Info("Service httpd started: <http://", global.CFG.ListenAddr, ">")
		if err := http.ListenAndServe(global.CFG.ListenAddr, logger.Req(logger.Recovery(mux))); err != nil {
			logger.Fatal(err)
		}
	}()

	osutil.WaitForInterrupt()
}
