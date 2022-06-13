package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var dsn string

// DSN設定
func SetDSNTCP(user string, password string, host string, port int, db string) string {
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, db)
	return fmt.Sprintf("%s:********@tcp(%s:%d)/%s?parseTime=true", user, host, port, db)
}

// DB接続
func Open() (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn does not set")
	}
	return sql.Open("mysql", dsn)
}

// トランザクション開始
func Begin() (tx *sql.Tx, err error) {
	d, err := Open()
	if err != nil {
		return
	}
	tx, err = d.Begin()
	return
}
