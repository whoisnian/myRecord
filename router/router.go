package router

import (
	"time"

	"github.com/whoisnian/glb/httpd"
)

type jsonMap map[string]interface{}

func timeHandler(store *httpd.Store) {
	store.RespondJson(jsonMap{"time": time.Now()})
}

func Init() *httpd.Mux {
	mux := httpd.NewMux()
	mux.Handle("/time", "GET", timeHandler)
	return mux
}
