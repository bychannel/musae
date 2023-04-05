package utils

import (
	"fmt"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils/snowflake"
	"math/rand"
	"time"
)

var guidNode *snowflake.Node

func init() {
	rand.Seed(time.Now().Unix())
	var err error
	nodeId := int64(rand.Uint32() % 1024)

	guidNode, err = snowflake.NewNode(nodeId)
	if err != nil {
		fmt.Println("guid init err:", nodeId, err)
		panic(any(fmt.Sprintf("guid init err:%s", err.Error())))
	}
}

func GenStrGUID() string {
	return GenStrUUID()
}

func GenIntGUID() int64 {
	return int64(GenIntUUID())
}
