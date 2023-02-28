package global

type Config struct {
	Debug      bool   `flag:"d,false,Enable debug output"`
	ListenAddr string `flag:"l,127.0.0.1:9000,Server listen addr"`
}

var CFG Config
