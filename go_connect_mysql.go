package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //init()匿名导入
)

/*
	## Go 操作 MySQL

- 下载依赖
  go get -u github.com/go-sql-driver/mysql

- 使用 MySQL 驱动
  func Open(driverName, dataSourceName string) (*DB, error)
*/

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome"
	//dsn := "uroot:proot@tcp(127.0.0.1:3306)/mysql_demo(数据库名)"

	//去初始化一个全局的db对象不是创建一个对象

	db, err = sql.Open("mysql", dsn)
	//db ,err := sql.Open("mysql",dsn) myslq为连接数据库的类型
	if err != nil {
		panic(err) //退出
	}

	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
		return
	}
	db.SetConnMaxIdleTime(2000)             //最大连接数
	db.SetConnMaxLifetime(time.Second * 10) //连接可能被重用的最大时间
	db.SetMaxIdleConns(100)                 //最大空闲连接数
	return
}

func main() {

	/*
		// DSN:Data Source Name
		dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome"
		//dsn := "uroot:proot@tcp(127.0.0.1:3306)/mysql_demo(数据库名)"
		db, err := sql.Open("mysql", dsn)
		//db ,err := sql.Open("mysql",dsn) myslq为连接数据库的类型
		if err != nil {
			panic(err) //退出
		}
		defer db.Close() // 注意这行代码要写在上面err判断的下面 关闭🥁
		// Close()释放掉myslq连接的资源

		// 尝试与数据库建立连接（校验dsn是否正确）
		err = db.Ping()
		if err != nil {
			fmt.Printf("connect to the db failed, err:%v", err)
			return
		}
	*/

	if err := initDB(); err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
	}

	defer db.Close() // 注意这行代码要写在上面err判断的下面 关闭🥁
	// Close()释放掉myslq连接的资源

	// 调用数据库函数
	fmt.Println("connect to the db sucess")
	// db.**操作数据库
}

/*
	Open函数可能只是验证其参数格式是否正确，实际上并不创建与数据库的连接。如果要检查数据源的名称是否真实有效，应该调用Ping方法。

返回的DB对象可以安全地被多个goroutine并发使用，并且维护其自己的空闲连接池。因此，Open函数应该仅被调用一次，很少需要关闭这个DB对象。
*/
