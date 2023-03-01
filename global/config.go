package global

type Config struct {
	Debug bool `flag:"d,false,Enable debug output"`

	Version     bool `flag:"v,false,Show version and quit"`
	CreateTable bool `flag:"ct,false,Create tables and quit"`

	ListenAddr  string `flag:"l,127.0.0.1:9000,Server listen addr"`
	DatabaseURI string `flag:"db,postgresql://postgres@127.0.0.1/record,PostgreSQL database connection URI"`
}

var CFG Config
