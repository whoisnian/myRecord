package routes

import "github.com/whoisnian/glb/httpd"

type jsonMap map[string]any

func SetupRouter() *httpd.Mux {
	mux := httpd.NewMux()
	mux.Handle("/status", "GET", statusHandler)
	return mux
}
