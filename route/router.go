package route

import (
	"net/http"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/myRecord/static"
)

type jsonMap map[string]any

func createStaticHandler(handler http.Handler) httpd.HandlerFunc {
	return func(store *httpd.Store) {
		path := store.RouteParamAny()
		if path == "" {
			store.Redirect("/view/", http.StatusFound)
			return
		}
		handler.ServeHTTP(store.W, store.R)
	}
}

func SetupRouter() *httpd.Mux {
	root, err := static.Root()
	if err != nil {
		logger.Fatal(err)
	}
	fileHandler := http.StripPrefix("/view/", http.FileServer(http.FS(root)))

	mux := httpd.NewMux()
	mux.Handle("/api/items", "GET", listItemHandler)
	mux.Handle("/api/items", "POST", createItemHandler)
	mux.Handle("/api/items/:id", "PUT", updateItemHandler)
	mux.Handle("/api/items/:id", "DELETE", deleteItemHandler)
	mux.Handle("/status", "GET", statusHandler)
	mux.Handle("/*", "GET", createStaticHandler(fileHandler))
	return mux
}
