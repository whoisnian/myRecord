package route

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/myRecord/model"
	"github.com/whoisnian/myRecord/model/item"
)

func listItemHandler(store *httpd.Store) {
	result := []*item.Item{}
	if err := model.B[*item.Item]().Where("type = $?", item.TypeTodo).Where("state != $?", item.StateDeleted).Find(&result); err != nil {
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
