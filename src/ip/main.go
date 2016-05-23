package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var db *sql.DB

type Response struct {
	Ip       string `json:"ip"`
	Isp      string `json:"isp"`
	Country  string `json:"country"`
	Privince string `json:"province"`
	City     string `json:"city"`
}

func init() {
	db, _ = sql.Open("mysql", "monitor:monitor@tcp(127.0.0.1:3306)/dns?charset=utf8")
	db.SetMaxOpenConns(10)
}

func main() {
	http.HandleFunc("/info", getInfo)
	http.ListenAndServe(":8000", nil)
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
	fmt.Println(ip)
	var (
		i string
		c string
		p string
		a string
	)
	e := db.QueryRow("select op, c, p, a from ips where ip = inet_aton(?)", ip).Scan(&i, &c, &p, &a)
	var b []byte
	if e != nil {
		b, _ = json.Marshal(Response{Ip: ip})
		w.Write(b)
	}
	b, _ = json.Marshal(Response{Ip: ip, Isp: i, Country: c, Privince: p, City: a})
	w.Write(b)
}
