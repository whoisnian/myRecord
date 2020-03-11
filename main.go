package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var CONFIG Config

var firstDayRecord int = 0
var firstWeekRecord int = 0
var firstMonthRecord int = 0

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

type flag struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

type flags struct {
	Num   int    `json:"num"`
	Flags []flag `json:"flags"`
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
		//w.Header().Set("Access-Control-Allow-Origin", "*")
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
	if (t == DayRecordType && to < firstDayRecord) ||
		(t == WeekRecordType && to < firstWeekRecord) ||
		(t == MonthRecordType && to < firstMonthRecord) {
		w.Write([]byte(`{"num":-1,"records":null}`))
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
	//w.Header().Set("Access-Control-Allow-Origin", "*")
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
		isoYear, isoWeek := nowTime.ISOWeek()
		recordTime = isoYear*100 + isoWeek
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

// 获取单个flag
func getSingleFlag(w http.ResponseWriter, r *http.Request) {
	// 获取请求flag的id
	id, err := strconv.Atoi(r.URL.Path[len("/flag/"):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 根据id查询flag
	row, err := db.Query("SELECT id, content, status FROM flag WHERE id=?", id)
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
		var resStatus int
		row.Scan(&resId, &resContent, &resStatus)
		var res = flag{
			resId,
			resContent,
			resStatus}
		resp, err := json.Marshal(res)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// 获取指定status的flag
func getFlags(w http.ResponseWriter, r *http.Request) {
	// 获取请求flag的status
	var status int
	var err error
	if r.FormValue("status") == "" {
		status = 0
	} else {
		status, err = strconv.Atoi(r.FormValue("status"))
		if err != nil {
			status = 0
		}
	}

	// 根据范围查询flag
	var row *sql.Rows
	if status == 0 {
		row, err = db.Query("SELECT id, content, status FROM flag")
	} else {
		row, err = db.Query("SELECT id, content, status FROM flag WHERE status=?", status)
	}
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回json
	var res flags
	for row.Next() {
		var resId int
		var resContent string
		var resStatus int
		row.Scan(&resId, &resContent, &resStatus)
		res.Flags = append(res.Flags, flag{
			resId,
			resContent,
			resStatus})
	}
	res.Num = len(res.Flags)
	resp, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// 新建flag
func newFlag(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}

	// 查询是否已有flag
	row, err := db.Query("SELECT 1 FROM flag WHERE content=? limit 1", r.FormValue("content"))
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

	// 添加flag
	_, err = db.Exec("INSERT flag SET content=?,status=?", r.FormValue("content"), 1)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// 更新flag
func updateFlag(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}

	// 获取请求flag的id
	id, err := strconv.Atoi(r.URL.Path[len("/flag/"):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 获取要修改的属性
	var content string
	var status int
	if r.FormValue("status") == "" {
		status = 0
	} else {
		status, err = strconv.Atoi(r.FormValue("status"))
		if err != nil {
			status = 0
		}
	}
	if r.FormValue("content") != "" {
		content = r.FormValue("content")
	} else {
		content = ""
	}

	// 根据id查询flag
	row, err := db.Query("SELECT 1 FROM flag WHERE id=?", id)
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if row.Next() {
		// 更新flag
		if content != "" {
			_, err = db.Exec("UPDATE flag SET content=? WHERE id=?", content, id)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if status != 0 {
			_, err = db.Exec("UPDATE flag SET status=? WHERE id=?", status, id)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// 删除flag
func deleteFlag(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}

	// 获取请求flag的id
	id, err := strconv.Atoi(r.URL.Path[len("/flag/"):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 根据id查询flag
	row, err := db.Query("SELECT 1 FROM flag WHERE id=?", id)
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if row.Next() {
		// 删除flag
		_, err = db.Exec("DELETE FROM flag WHERE id=?", id)
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

	row, err := db.Query("SELECT time FROM record_day order by time limit 1")
	if err != nil {
		log.Panicln(err.Error())
	}
	if row.Next() {
		row.Scan(&firstDayRecord)
	}
	row.Close()

	row, err = db.Query("SELECT time FROM record_week order by time limit 1")
	if err != nil {
		log.Panicln(err.Error())
	}
	if row.Next() {
		row.Scan(&firstWeekRecord)
	}
	row.Close()

	row, err = db.Query("SELECT time FROM record_month order by time limit 1")
	if err != nil {
		log.Panicln(err.Error())
	}
	if row.Next() {
		row.Scan(&firstMonthRecord)
	}
	row.Close()

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

	// GET		/flag/{id}						getSingleFlag
	// GET		/flag/?status={status}			getFlags
	// POST		/flag/							newFlag
	// PUT		/flag/{id}						updateFlag
	// DELETE	/flag/{id}						deleteFlag
	http.HandleFunc("/flag/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		if r.Method == http.MethodGet {
			if r.FormValue("status") == "" {
				getSingleFlag(w, r)
			} else {
				getFlags(w, r)
			}
		} else if r.Method == http.MethodPost {
			if r.FormValue("_method") == http.MethodPut {
				updateFlag(w, r)
			} else if r.FormValue("_method") == http.MethodDelete {
				deleteFlag(w, r)
			} else {
				newFlag(w, r)
			}
		} else if r.Method == http.MethodPut {
			updateFlag(w, r)
		} else if r.Method == http.MethodDelete {
			deleteFlag(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("html"))))

	// 启动服务
	log.Printf("Server started: <http://127.0.0.1:%v>\n", CONFIG.PORT)
	err = http.ListenAndServe(":"+CONFIG.PORT, nil)
	if err != nil {
		log.Panicln(err.Error())
	}
}
