package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql" //init()匿名导入
)

var dbsqlxuse *sqlx.DB

func initDBsqlxuse() (err error) {
	dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	dbsqlxuse, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	dbsqlxuse.SetMaxOpenConns(20)
	dbsqlxuse.SetMaxIdleConns(10)
	return
}

func main() {
	if err := initDBsqlxuse(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("连接成功...")
}
