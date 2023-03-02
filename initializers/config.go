package initializers

import (
	"github.com/whoisnian/glb/config"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/myRecord/global"
)

func SetupConfig() *global.Config {
	cfg := &global.Config{}
	if err := config.FromCommandLine(cfg); err != nil {
		logger.Fatal(err)
	}

	logger.SetDebug(cfg.Debug)
	return cfg
}
