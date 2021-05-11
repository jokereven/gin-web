package main

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql" //init()匿名导入
)

var dbsqlxin *sqlx.DB

func initDBsqlxin() (err error) {
	dsn := "root:3144588210XZQxzq@tcp(127.0.0.1:3306)/go_dome?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	dbsqlxin, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	dbsqlxin.SetMaxOpenConns(20)
	dbsqlxin.SetMaxIdleConns(10)
	return
}

type User struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

/*
	使用sqlx.In实现批量插入
前提是需要我们的结构体实现driver.Valuer接口：
*/

func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO user (name, age) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := dbsqlxin.Exec(query, args...)
	return err
}

/*
sqlx.In的查询示例
关于sqlx.In这里再补充一个用法，在sqlx查询语句中实现In查询和FIND_IN_SET函数。即实现SELECT * FROM user WHERE id in (3, 2, 1);和SELECT * FROM user WHERE id in (3, 2, 1) ORDER BY FIND_IN_SET(id, '3,2,1');。

in查询
查询id在给定id集合中的数据。
*/

// QueryByIDs 根据给定ID查询
func QueryByIDs(ids []int) (users []User, err error) {
	// 动态填充id
	query, args, err := sqlx.In("SELECT name, age FROM user WHERE id IN (?)", ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = dbsqlxin.Rebind(query)

	err = dbsqlxin.Select(&users, query, args...)
	return
}

/*
	in查询和FIND_IN_SET函数
查询id在给定id集合的数据并维持给定id集合的顺序。
*/

// QueryAndOrderByIDs 按照指定id查询并维护顺序
func QueryAndOrderByIDs(ids []int) (users []User, err error) {
	// 动态填充id
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In("SELECT name, age FROM user WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)", ids, strings.Join(strIDs, ","))
	if err != nil {
		return
	}

	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = dbsqlxin.Rebind(query)

	err = dbsqlxin.Select(&users, query, args...)
	return
}

func main() {
	if err := initDBsqlxin(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("连接成功...")

	// u1 := User{Name: "小黄鸭", Age: 18}
	// u2 := User{Name: "大白", Age: 18}
	// u3 := User{Name: "咖啡猫🐱", Age: 18}
	// users := []interface{}{u1, u2, u3}
	// BatchInsertUsers2(users)

	user, err := QueryByIDs([]int{23, 24, 25})
	if err != nil {
		fmt.Println(err)
	}
	for _, con := range user {
		fmt.Println(con)
	}

	users, err := QueryAndOrderByIDs([]int{25, 24, 23})
	if err != nil {
		fmt.Println(err)
	}
	for _, con := range users {
		fmt.Println(con)
	}
}
