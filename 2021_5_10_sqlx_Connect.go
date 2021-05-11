package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql" //init()匿名导入
)

var dbsqlx *sqlx.DB

func initDBsqlx() (err error) {
	dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	dbsqlx, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	dbsqlx.SetMaxOpenConns(20)
	dbsqlx.SetMaxIdleConns(10)
	return
}

func main_Connect() {
	if err := initDBsqlx(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("连接成功...")
}
