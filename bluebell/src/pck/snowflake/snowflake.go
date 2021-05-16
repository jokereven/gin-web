package snowflake

import (
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

// Init ...
// 初始化ID生成器（基于雪花算法实现）
func Init(startTime string, machineID int64) error {
	var st time.Time
	var err error

	if st, err = time.Parse("2006-01-02", startTime); err != nil {
		fmt.Printf("init userID startTime failed, err: %v", err)
		return err
	}

	sf.Epoch = st.UnixNano() / 1000000

	node, err = sf.NewNode(machineID)
	return err
}

// GenID ...
// 获取ID
func GenID() int64 {
	return node.Generate().Int64()
}

// func main() {
// 	if err := Init("2020-12-15", 1); err != nil {
// 		panic(err)
// 	}
// 	id := GenID()
// 	fmt.Println(id)
// }
