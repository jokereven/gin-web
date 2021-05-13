package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //init()匿名导入
)

// 定义一个全局对象dbmysql
var dbmysqldell *sql.DB

// 定义一个初始化数据库的函数

func initDBMYSQLDELL() (err error) {
	// DSN:Data Source Name
	dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome"
	//dsn := "uroot:proot@tcp(127.0.0.1:3306)/mysql_demo(数据库名)"

	//去初始化一个全局的db对象不是创建一个对象

	dbmysqldell, err = sql.Open("mysql", dsn)
	//db ,err := sql.Open("mysql",dsn) myslq为连接数据库的类型
	if err != nil {
		panic(err) //退出
	}

	// 尝试与数据库建立连接（校验dsn是否正确）
	err = dbmysqldell.Ping()
	if err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
		return
	}
	dbmysqldell.SetConnMaxIdleTime(2000)             //最大连接数
	dbmysqldell.SetConnMaxLifetime(time.Second * 10) //连接可能被重用的最大时间
	dbmysqldell.SetMaxIdleConns(100)                 //最大空闲连接数
	return
}

type userdell struct {
	id   int
	age  int
	name string
}

// 预处理查询示例
func prepareQueryDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	stmt, err := dbmysqldell.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next() {
		var userdell userdell
		err := rows.Scan(&userdell.id, &userdell.name, &userdell.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", userdell.id, userdell.name, userdell.age)
	}
}

// 预处理插入示例
func prepareInsertDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	stmt, err := dbmysqldell.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("小王子", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	_, err = stmt.Exec("沙河娜扎", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	fmt.Println("insert success.")
}

// sql注入示例
func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select id, name, age from user where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)
	var u userdell
	err := dbmysqldell.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}

func main_SQLdell() {

	if err := initDBMYSQLDELL(); err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
	}

	defer dbmysqldell.Close() // 注意这行代码要写在上面err判断的下面 关闭🥁
	// Close()释放掉myslq连接的资源

	// 调用数据库函数
	fmt.Println("connect to the db sucess")
	// db.**操作数据库

	// prepareQueryDemo() //查数据

	// prepareInsertDemo() //预处理插入

	// sqlInjectDemo("鞠婧祎") //sql注入

	//此时以下输入字符串都可以引发SQL注入问题：
	sqlInjectDemo("xxx' or 1=1#")
	sqlInjectDemo("xxx' union select * from user #")
	sqlInjectDemo("xxx' and (select count(*) from user) <10 #")

	//mysql占位符 ?
}

/*
什么是预处理？

1. 普通SQL语句执行过程：

客户端对SQL语句进行占位符替换得到完整的SQL语句。
客户端发送完整SQL语句到MySQL服务端
MySQL服务端执行完整的SQL语句并将结果返回给客户端。

2. 预处理执行过程：

把SQL语句分成两部分，命令部分与数据部分。
先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
然后把数据部分发送给MySQL服务端，MySQL服务端对SQL语句进行占位符替换。
MySQL服务端执行完整的SQL语句并将结果返回给客户端。

为什么要预处理？
优化MySQL服务器重复执行SQL的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
避免SQL注入问题。

Go实现MySQL预处理
database/sql中使用下面的Prepare方法来实现预处理操作。

func (db *DB) Prepare(query string) (*Stmt, error)

Prepare方法会先将sql语句发送给MySQL服务端，返回一个准备好的状态用于之后的查询和命令。
返回值可以同时执行多个查询和命令。
*/
