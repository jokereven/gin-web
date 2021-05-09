package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func func1(c *gin.Context) {
	fmt.Println("星期一")
	c.Set("name", "jokerxue")
}
func func2(c *gin.Context) {
	fmt.Println("星期二")
	v, ok := c.Get("name")
	if ok {
		str := v.(string)
		fmt.Println(str)
	}
}
func func3(c *gin.Context) {
	c.Next() //先执行吓一条链路
	fmt.Println("星期三")
}
func func4(c *gin.Context) {
	fmt.Println("星期四")
	c.Abort() //后面的直接不执行
}
func func5(c *gin.Context) {
	c.Next() //先执行下一条
	fmt.Println("星期五")
}
func func6(c *gin.Context) {
	fmt.Println("星期六")
	fmt.Println("星期天")
}

func Gin_source() {
	r := gin.Default()

	r.GET("/joker", func(c *gin.Context) {})

	showGroup := r.Group("/group", func1, func2)

	showGroup.Use(func3, func4)

	{
		showGroup.GET("/index", func5, func6)
	}
	r.Run()
}
