package global

import (
	"encoding/json"

	"github.com/whoisnian/glb/logger"
)

var CFG *Config

type Config struct {
	Debug bool `flag:"d,false,Enable debug output"`

	Version     bool `flag:"v,false,Show version and quit"`
	CreateTable bool `flag:"ct,false,Create tables and quit"`

	ListenAddr  string `flag:"l,127.0.0.1:9000,Server listen addr"`
	DatabaseURI string `flag:"db,postgresql://postgres@127.0.0.1/record,PostgreSQL database connection URI"`
}

func (cfg Config) Json() string {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		logger.Error(err)
		return "{}"
	}
	return string(data)
}
