package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello go web")
	r := gin.Default()
	r.GET("/")
	r.Run()
}
