package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //init()匿名导入
)

// 定义一个全局对象dbmysql
var dbmysqlacid *sql.DB

// 定义一个初始化数据库的函数

func initDBMYSQLACID() (err error) {
	// DSN:Data Source Name
	dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome"
	//dsn := "uroot:proot@tcp(127.0.0.1:3306)/mysql_demo(数据库名)"

	//去初始化一个全局的db对象不是创建一个对象

	dbmysqlacid, err = sql.Open("mysql", dsn)
	//db ,err := sql.Open("mysql",dsn) myslq为连接数据库的类型
	if err != nil {
		panic(err) //退出
	}

	// 尝试与数据库建立连接（校验dsn是否正确）
	err = dbmysqlacid.Ping()
	if err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
		return
	}
	dbmysqlacid.SetConnMaxIdleTime(2000)             //最大连接数
	dbmysqlacid.SetConnMaxLifetime(time.Second * 10) //连接可能被重用的最大时间
	dbmysqlacid.SetMaxIdleConns(100)                 //最大空闲连接数
	return
}

type useracid struct {
	id   int
	age  int
	name string
}

// 事务操作示例
func transactionDemo() {
	tx, err := dbmysqlacid.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("exec trans success!")
}

func main() {
	if err := initDBMYSQLACID(); err != nil {
		fmt.Printf("connect to the db failed, err:%v", err)
	}

	defer dbmysqlacid.Close() // 注意这行代码要写在上面err判断的下面 关闭🥁
	// Close()释放掉myslq连接的资源

	// 调用数据库函数
	fmt.Println("connect to the db sucess")
	// db.**操作数据库

	transactionDemo() //事务
}

/*
	Go实现MySQL事务
什么是事务？
事务：一个最小的不可再分的工作单元；通常一个事务对应一个完整的业务(例如银行账户转账业务，该业务就是一个最小的工作单元)，同时这个完整的业务需要执行多次的DML(insert、update、delete)语句共同联合完成。A转账给B，这里面就需要执行两次update操作。

在MySQL中只有使用了Innodb数据库引擎的数据库或表才支持事务。事务处理可以用来维护数据库的完整性，保证成批的SQL语句要么全部执行，要么全部不执行。

事务的ACID
通常事务必须满足4个条件（ACID）：原子性（Atomicity，或称不可分割性）、一致性（Consistency）、隔离性（Isolation，又称独立性）、持久性（Durability）。

条件	解释
原子性	一个事务（transaction）中的所有操作，要么全部完成，要么全部不完成，不会结束在中间某个环节。事务在执行过程中发生错误，会被回滚（Rollback）到事务开始前的状态，就像这个事务从来没有执行过一样。
一致性	在事务开始之前和事务结束以后，数据库的完整性没有被破坏。这表示写入的资料必须完全符合所有的预设规则，这包含资料的精确度、串联性以及后续数据库可以自发性地完成预定的工作。
隔离性	数据库允许多个并发事务同时对其数据进行读写和修改的能力，隔离性可以防止多个事务并发执行时由于交叉执行而导致数据的不一致。事务隔离分为不同级别，包括读未提交（Read uncommitted）、读提交（read committed）、可重复读（repeatable read）和串行化（Serializable）。
持久性	事务处理结束后，对数据的修改就是永久的，即便系统故障也不会丢失。
事务相关方法
Go语言中使用以下三个方法实现MySQL中的事务操作。 开始事务

func (db *DB) Begin() (*Tx, error)
提交事务

func (tx *Tx) Commit() error
回滚事务

func (tx *Tx) Rollback() error
*/
