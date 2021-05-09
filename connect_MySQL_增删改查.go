package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //init()匿名导入
)

// 定义一个全局对象dbmysql
var dbmysql *sql.DB

// 定义一个初始化数据库的函数

func initDBMYSQL() (err error) {
	// DSN:Data Source Name
	dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome"
	//dsn := "uroot:proot@tcp(127.0.0.1:3306)/mysql_demo(数据库名)"

	//去初始化一个全局的db对象不是创建一个对象

	dbmysql, err = sql.Open("mysql", dsn)
	//db ,err := sql.Open("mysql",dsn) myslq为连接数据库的类型
	if err != nil {
		panic(err) //退出
	}

	// 尝试与数据库建立连接（校验dsn是否正确）
	err = dbmysql.Ping()
	if err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
		return
	}
	dbmysql.SetConnMaxIdleTime(2000)             //最大连接数
	dbmysql.SetConnMaxLifetime(time.Second * 10) //连接可能被重用的最大时间
	dbmysql.SetMaxIdleConns(100)                 //最大空闲连接数
	return
}

type user struct {
	id   int
	age  int
	name string
}

// 查询单条数据示例
func queryRowDemo() {
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	rows, err := dbmysql.Query(`
        SELECT id,name,age FROM user
    `)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		// 对于遍历, 只需要判断每次是否有错误产生即可
		// 参数绑定需要数量和位置一一对应
		if err := rows.Scan(&u.id, &u.name, &u.age); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := dbmysql.Exec(sqlStr, "我爱罗", 22)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := dbmysql.Exec(sqlStr, 39, 3)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := dbmysql.Exec(sqlStr, 8) //8是删除的id
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

func main() {

	if err := initDBMYSQL(); err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
	}

	defer dbmysql.Close() // 注意这行代码要写在上面err判断的下面 关闭🥁
	// Close()释放掉myslq连接的资源

	// 调用数据库函数
	fmt.Println("connect to the db sucess")
	// db.**操作数据库

	// queryRowDemo() //调用查询函数

	// insertRowDemo() //插入

	// updateRowDemo() //更新数据

	deleteRowDemo() //删除数据
}

/*
	Open函数可能只是验证其参数格式是否正确，实际上并不创建与数据库的连接。如果要检查数据源的名称是否真实有效，应该调用Ping方法。

返回的DB对象可以安全地被多个goroutine并发使用，并且维护其自己的空闲连接池。因此，Open函数应该仅被调用一次，很少需要关闭这个DB对象。
*/
