package snowflake

import (
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

// Init ...
// 初始化ID生成器（基于雪花算法实现）
func Init(StartTime string, MachineID int64) (err error) {
	var st time.Time
	if st, err = time.Parse("2006-01-02", StartTime); err != nil {
		fmt.Printf("init userID startTime failed, err: %v", err)
		return err
	}

	sf.Epoch = st.UnixNano() / 1000000

	node, err = sf.NewNode(MachineID)
	return err
}

// 获取ID
func GenID() int64 {
	return node.Generate().Int64()
}
