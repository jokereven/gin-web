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

type usersqlx struct {
	ID   int    `db:"id"` //首字母大写对外 映射到数据库里面的东西 `db:"xxx"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// 查询单条数据示例
func queryRowDemosqlxuse() {
	sqlStr := "select id, name, age from user where id=?"
	var u usersqlx
	err := dbsqlxuse.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.ID, u.Name, u.Age)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	var users []usersqlx
	err := dbsqlxuse.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

// 插入数据
func insertRowDemosqlx() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := dbsqlxuse.Exec(sqlStr, "沙河小王子", 19)
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
func updateRowDemosqlx() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := dbsqlxuse.Exec(sqlStr, 39, 6)
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
func deleteRowDemosqlx() {
	sqlStr := "delete from user where id = ?"
	ret, err := dbsqlxuse.Exec(sqlStr, 6)
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

/*
	NamedExec
DB.NamedExec方法用来绑定SQL语句与结构体或map中的同名字段。
*/

func insertUserDemo() (err error) {
	sqlStr := "INSERT INTO user (name,age) VALUES (:name,:age)"
	_, err = dbsqlxuse.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "我爱罗",
			"age":  28,
		})
	return
}

/*
	NamedQuery
与DB.NamedExec同理，这里是支持查询。
*/

func namedQuery() {
	sqlStr := "SELECT * FROM user WHERE name=:name"
	// 使用map做命名查询
	rows, err := dbsqlxuse.NamedQuery(sqlStr, map[string]interface{}{"name": "七米"})
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u usersqlx
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}

	u := usersqlx{
		Name: "七米",
	}
	// 使用结构体命名查询，根据结构体字段的 db tag进行映射
	rows, err = dbsqlxuse.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u usersqlx
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

func main_sql_use() {
	if err := initDBsqlxuse(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("连接成功...")

	queryRowDemosqlxuse() //查询单条根据指定的类容

	queryMultiRowDemo() //查询多条 切片类型

	// insertRowDemosqlx() //插入

	// updateRowDemosqlx() //更新

	// deleteRowDemosqlx() //删除

	insertUserDemo() //
}
