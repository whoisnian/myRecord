package route

import (
	"github.com/whoisnian/glb/httpd"
)

func statusHandler(store *httpd.Store) {
	store.RespondJson(jsonMap{"status": "ok"})
}
