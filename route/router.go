package route

import "github.com/whoisnian/glb/httpd"

type jsonMap map[string]any

func SetupRouter() *httpd.Mux {
	mux := httpd.NewMux()
	mux.Handle("/api/items", "GET", listItemHandler)
	mux.Handle("/api/items", "POST", createItemHandler)
	mux.Handle("/api/items/:id", "PUT", updateItemHandler)
	mux.Handle("/api/items/:id", "DELETE", deleteItemHandler)
	mux.Handle("/status", "GET", statusHandler)
	return mux
}
