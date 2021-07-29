package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	Username     = "hzj"
	Password     = "Mysql@123"
	Hostname     = "127.0.0.1:3306"
	DatabaseName = "ttbb"
)

func Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", Username, Password, Hostname, DatabaseName)
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", Dsn())
	if err != nil {
		panic("connect database error")
	}
}

func getPlaceHolders(n int) string {
	auxSlice := make([]string, n)
	for i := 0; i < n; i++ {
		auxSlice[i] = "?"
	}
	return strings.Join(auxSlice, ",")
}
