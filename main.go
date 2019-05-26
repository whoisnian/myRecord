package main

import (
	"database/sql"
	"encoding/json"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

var db *sql.DB
var CONFIG Config

type Config struct {
	PORT  string `toml:"port"`
	DSN   string `toml:"dsn"`
	TOKEN string `toml:"token"`
}

type RecordType string

const (
	DayRecordType   RecordType = "record_day"
	WeekRecordType  RecordType = "record_week"
	MonthRecordType RecordType = "record_month"
)

type record struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Time    int    `json:"time"`
}

type records struct {
	Num     int      `json:"num"`
	Records []record `json:"records"`
}

// 检查Token
func checkToken(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Authorization") == "" && r.FormValue("token") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	} else if r.Header.Get("Authorization") != "Bearer "+CONFIG.TOKEN && r.FormValue("token") != CONFIG.TOKEN {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}

// 获取单个record
func getSingleRecord(w http.ResponseWriter, r *http.Request, t RecordType) {
	// 获取请求record的id
	var id int
	var err error
	if t == DayRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/day-record/"):])
	} else if t == WeekRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/week-record/"):])
	} else if t == MonthRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/month-record/"):])
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 根据id查询record
	row, err := db.Query("SELECT id, content, time FROM "+string(t)+" WHERE id=?", id)
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if row.Next() {
		// 返回json
		var resId int
		var resContent string
		var resTime int
		row.Scan(&resId, &resContent, &resTime)
		var res = record{
			resId,
			resContent,
			resTime}
		resp, err := json.Marshal(res)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// 获取多个record
func getRecords(w http.ResponseWriter, r *http.Request, t RecordType) {
	// 获取请求record的范围
	from, err := strconv.Atoi(r.FormValue("from"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	to, err := strconv.Atoi(r.FormValue("to"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 根据范围查询record
	row, err := db.Query("SELECT id, content, time FROM "+string(t)+" WHERE time BETWEEN ? AND ?", from, to)
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回json
	var res records
	for row.Next() {
		var resId int
		var resContent string
		var resTime int
		row.Scan(&resId, &resContent, &resTime)
		res.Records = append(res.Records, record{
			resId,
			resContent,
			resTime})
	}
	res.Num = len(res.Records)
	resp, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// 新建record
func newRecord(w http.ResponseWriter, r *http.Request, t RecordType) {
	if !checkToken(w, r) {
		return
	}

	var recordTime int
	nowTime := time.Now()
	if t == DayRecordType {
		recordTime = nowTime.Year()*10000 + int(nowTime.Month())*100 + nowTime.Day()
	} else if t == WeekRecordType {
		recordTime = nowTime.Year()*100 + (nowTime.YearDay()-int(nowTime.Weekday())+7)/7
	} else if t == MonthRecordType {
		recordTime = nowTime.Year()*100 + int(nowTime.Month())
	}

	// 查询是否已有record
	row, err := db.Query("SELECT 1 FROM "+string(t)+" WHERE time=? limit 1", recordTime)
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if row.Next() {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// 添加record
	_, err = db.Exec("INSERT "+string(t)+" SET content=?,time=?", r.FormValue("content"), recordTime)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// 更新record
func updateRecord(w http.ResponseWriter, r *http.Request, t RecordType) {
	if !checkToken(w, r) {
		return
	}

	// 获取请求record的id
	var id int
	var err error
	if t == DayRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/day-record/"):])
	} else if t == WeekRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/week-record/"):])
	} else if t == MonthRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/month-record/"):])
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 根据id查询record
	row, err := db.Query("SELECT 1 FROM "+string(t)+" WHERE id=?", id)
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if row.Next() {
		// 更新record
		_, err = db.Exec("UPDATE "+string(t)+" SET content=? WHERE id=?", r.FormValue("content"), id)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// 删除record
func deleteRecord(w http.ResponseWriter, r *http.Request, t RecordType) {
	if !checkToken(w, r) {
		return
	}

	// 获取请求record的id
	var id int
	var err error
	if t == DayRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/day-record/"):])
	} else if t == WeekRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/week-record/"):])
	} else if t == MonthRecordType {
		id, err = strconv.Atoi(r.URL.Path[len("/month-record/"):])
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 根据id查询record
	row, err := db.Query("SELECT 1 FROM "+string(t)+" WHERE id=?", id)
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if row.Next() {
		// 删除record
		_, err = db.Exec("DELETE FROM "+string(t)+" WHERE id=?", id)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func main() {
	_, err := toml.DecodeFile("config.toml", &CONFIG)
	if err != nil {
		panic(err)
	}

	// 连接数据库
	db, err = sql.Open("mysql", CONFIG.DSN)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer db.Close()

	// 检查是否连接上
	err = db.Ping()
	if err != nil {
		log.Panicln(err.Error())
	}

	// GET		/day-record/{id}						getSingleRecord
	// GET		/day-record/?from={time1}&to={time2}	getRecords
	// POST		/day-record/							newRecord
	// PUT		/day-record/{id}						updateRecord
	// DELETE	/day-record/{id}						deleteRecord
	http.HandleFunc("/day-record/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		if r.Method == http.MethodGet {
			if r.FormValue("from") == "" {
				getSingleRecord(w, r, DayRecordType)
			} else {
				getRecords(w, r, DayRecordType)
			}
		} else if r.Method == http.MethodPost {
			if r.FormValue("_method") == http.MethodPut {
				updateRecord(w, r, DayRecordType)
			} else if r.FormValue("_method") == http.MethodDelete {
				deleteRecord(w, r, DayRecordType)
			} else {
				newRecord(w, r, DayRecordType)
			}
		} else if r.Method == http.MethodPut {
			updateRecord(w, r, DayRecordType)
		} else if r.Method == http.MethodDelete {
			deleteRecord(w, r, DayRecordType)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// GET		/week-record/{id}						getSingleRecord
	// GET		/week-record/?from={time1}&to={time2}	getRecords
	// POST		/week-record/							newRecord
	// PUT		/week-record/{id}						updateRecord
	// DELETE	/week-record/{id}						deleteRecord
	http.HandleFunc("/week-record/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		if r.Method == http.MethodGet {
			if r.FormValue("from") == "" {
				getSingleRecord(w, r, WeekRecordType)
			} else {
				getRecords(w, r, WeekRecordType)
			}
		} else if r.Method == http.MethodPost {
			if r.FormValue("_method") == http.MethodPut {
				updateRecord(w, r, WeekRecordType)
			} else if r.FormValue("_method") == http.MethodDelete {
				deleteRecord(w, r, WeekRecordType)
			} else {
				newRecord(w, r, WeekRecordType)
			}
		} else if r.Method == http.MethodPut {
			updateRecord(w, r, WeekRecordType)
		} else if r.Method == http.MethodDelete {
			deleteRecord(w, r, WeekRecordType)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// GET		/month-record/{id}						getSingleRecord
	// GET		/month-record/?from={time1}&to={time2}	getRecords
	// POST		/month-record/							newRecord
	// PUT		/month-record/{id}						updateRecord
	// DELETE	/month-record/{id}						deleteRecord
	http.HandleFunc("/month-record/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		if r.Method == http.MethodGet {
			if r.FormValue("from") == "" {
				getSingleRecord(w, r, MonthRecordType)
			} else {
				getRecords(w, r, MonthRecordType)
			}
		} else if r.Method == http.MethodPost {
			if r.FormValue("_method") == http.MethodPut {
				updateRecord(w, r, MonthRecordType)
			} else if r.FormValue("_method") == http.MethodDelete {
				deleteRecord(w, r, MonthRecordType)
			} else {
				newRecord(w, r, MonthRecordType)
			}
		} else if r.Method == http.MethodPut {
			updateRecord(w, r, MonthRecordType)
		} else if r.Method == http.MethodDelete {
			deleteRecord(w, r, MonthRecordType)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// 启动服务
	log.Printf("Server started: <http://127.0.0.1:%v>\n", CONFIG.PORT)
	err = http.ListenAndServe(":"+CONFIG.PORT, nil)
	if err != nil {
		log.Panicln(err.Error())
	}
}
