package route

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/myRecord/model"
	"github.com/whoisnian/myRecord/model/item"
)

func listItemHandler(store *httpd.Store) {
	Q := store.R.URL.Query()
	qType := Q.Get("type")
	qState := Q.Get("state")
	qFrom := Q.Get("from")
	qTo := Q.Get("to")

	batch := model.B[*item.Item]()
	if qType != "" {
		if argType, err := strconv.ParseUint(qType, 10, 32); err == nil {
			batch = batch.Where("type = $?", argType)
		} else {
			store.W.WriteHeader(http.StatusBadRequest)
			store.W.Write([]byte(err.Error()))
			return
		}
	}
	if qState != "" {
		if argState, err := strconv.ParseUint(qState, 10, 32); err == nil {
			batch = batch.Where("state = $?", argState)
		} else {
			store.W.WriteHeader(http.StatusBadRequest)
			store.W.Write([]byte(err.Error()))
			return
		}
	}
	if qFrom != "" {
		if argFrom, err := strconv.ParseInt(qFrom, 10, 64); err == nil {
			batch = batch.Where("date >= $?", time.UnixMilli(argFrom))
		} else {
			store.W.WriteHeader(http.StatusBadRequest)
			store.W.Write([]byte(err.Error()))
			return
		}
	}
	if qTo != "" {
		if argTo, err := strconv.ParseInt(qTo, 10, 64); err == nil {
			batch = batch.Where("date < $?", time.UnixMilli(argTo))
		} else {
			store.W.WriteHeader(http.StatusBadRequest)
			store.W.Write([]byte(err.Error()))
			return
		}
	}

	result := []*item.Item{}
	if err := batch.Find(&result); err != nil {
		logger.Panic(err)
	}
	store.RespondJson(result)
}

func createItemHandler(store *httpd.Store) {
	result := item.Item{}
	if err := json.NewDecoder(store.R.Body).Decode(&result); err != nil {
		store.W.WriteHeader(http.StatusBadRequest)
		store.W.Write([]byte(err.Error()))
		return
	}
	if result.Content == "" || result.Date.Unix() <= 0 {
		store.W.WriteHeader(http.StatusBadRequest)
		store.W.Write([]byte("Invalid content or date"))
		return
	}
	if result.Exists() {
		store.W.WriteHeader(http.StatusConflict)
		store.W.Write([]byte("Item already exists"))
		return
	}
	if err := model.Create(&result); err != nil {
		logger.Panic(err)
	}
	store.RespondJson(result)
}

func updateItemHandler(store *httpd.Store) {
	id, err := strconv.ParseInt(store.RouteParam("id"), 10, 64)
	if err != nil {
		store.W.WriteHeader(http.StatusBadRequest)
		store.W.Write([]byte("Invalid ID"))
		return
	}

	result := item.Item{}
	if err := json.NewDecoder(store.R.Body).Decode(&result); err != nil {
		store.W.WriteHeader(http.StatusBadRequest)
		store.W.Write([]byte(err.Error()))
		return
	}
	result.Id = id
	if err := model.Update(&result); err != nil {
		logger.Panic(err)
	}
	store.RespondJson(result)
}

func deleteItemHandler(store *httpd.Store) {
	id, err := strconv.ParseInt(store.RouteParam("id"), 10, 64)
	if err != nil {
		store.W.WriteHeader(http.StatusBadRequest)
		store.W.Write([]byte("Invalid ID"))
		return
	}

	result := item.Item{Id: id}
	if err := model.Find(&result); err != nil {
		logger.Panic(err)
	}
	result.State = item.StateDeleted
	if err := model.Update(&result); err != nil {
		logger.Panic(err)
	}
	store.RespondJson(result)
}
