package test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils/ziputil"
	"io"
	"testing"
	"time"
)

type JsonObj struct {
	Id       int               `json:"id"`
	Val      string            `json:"va"`
	MapVal   map[string]string `json:"MapVal"`
	SliceVal []string          `json:"SliceVal"`
}

func TestConcat(t *testing.T) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write([]byte("012345678900000\n"))
	w.Close()

	r, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(r)
	if string(data) != "012345678900000\n" || err != nil {
		t.Fatalf("ReadAll = %q, %v, want %q, nil", data, err, "hello world")
	}
}

func Test_GZIP(t *testing.T) {

	//szBuf := "DEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 7\nDEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 9\nDEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 1011403\nDEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 1011401\nDEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 1011402\nDEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 1011503\nDEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 1011501\nDEBUG\t2022-12-01 14:35:37.6851023 +0800 CST m=+0.088561501 base_actor.go:24.RegisterProtoHandler: register messageId: 1011504\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011601\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011701\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011702\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011703\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011705\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011704\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011801\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1010613\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011505\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 base_actor.go:24.RegisterProtoHandler: register messageId: 1011502\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 user_actor.go:173.NewUserActor: UserActorFactory create UserActor, \nINFO\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 service_init.go:124.initSvc: [Service, init ActorFactory succeed]\nDEBUG\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 service_init.go:44.InitBase: [Service, initSvc succeed]\nINFO\t2022-12-01 14:35:37.6856342 +0800 CST m=+0.089093401 service_init.go:62.InitBase: [Service, initNet succeed]"
	szBuf := `[
	{
		"ambience_num": 0,
		"favor_effect": 3,
		"id": 1,
		"name": "ambience_name_1",
		"produce_effect": 2
	},
	{
		"ambience_num": 200,
		"favor_effect": 3,
		"id": 2,
		"name": "ambience_name_2",
		"produce_effect": 2
	},
	{
		"ambience_num": 400,
		"favor_effect": 3,
		"id": 3,
		"name": "ambience_name_3",
		"produce_effect": 2
	},
	{
		"ambience_num": 600,
		"favor_effect": 5,
		"id": 4,
		"name": "ambience_name_4",
		"produce_effect": 3
	},
	{
		"ambience_num": 800,
		"favor_effect": 5,
		"id": 5,
		"name": "ambience_name_5",
		"produce_effect": 3
	},
	{
		"ambience_num": 1000,
		"favor_effect": 5,
		"id": 6,
		"name": "ambience_name_6",
		"produce_effect": 5
	},
	{
		"ambience_num": 1200,
		"favor_effect": 8,
		"id": 7,
		"name": "ambience_name_7",
		"produce_effect": 5
	},
	{
		"ambience_num": 1600,
		"favor_effect": 8,
		"id": 8,
		"name": "ambience_name_8",
		"produce_effect": 5
	}
]`
	bytes := []byte(szBuf)

	gzip := &ziputil.Gzip{}
	gzipBytes, err := gzip.ZipEncode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ungzipBytes, err := gzip.ZipDecode(gzipBytes)
	fmt.Println(szBuf, string(ungzipBytes), err)

}

func Test_0023(t *testing.T) {
	jsonCount := 1

	objs := buildJsonObjs(jsonCount)

	bytes, err := json.Marshal(objs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonByteLen := len(bytes)

	startMilSec := time.Now().UnixMilli()
	gzip := &ziputil.Gzip{}
	gzipBytes, err := gzip.ZipEncode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	gzipByteLen := len(gzipBytes)

	// json数量10000, 大小:877781 ===>>> gzip压缩耗时:22ms ===> 压缩后json大小:43690, 压缩比:4.98%
	fmt.Print(fmt.Sprintf("json数量%d, 大小:%d ===>>> gzip压缩耗时:%dms ===> 压缩后json大小:%d, 压缩比:%.2f%%",
		jsonCount,
		jsonByteLen,
		time.Now().UnixMilli()-startMilSec,
		gzipByteLen,
		float64(gzipByteLen)/float64(jsonByteLen)*100))

	startMilSec = time.Now().UnixMilli()
	ungzipBytes, err := gzip.ZipDecode(gzipBytes)
	assert.NoError(t, err)

	ungzipByteLen := len(ungzipBytes)
	fmt.Println(fmt.Sprintf("===>>> gzip解压缩耗时:%dms ===> 解压缩后数据大小:%d",
		time.Now().UnixMilli()-startMilSec,
		ungzipByteLen))

	time.Sleep(1 * time.Second)
}

func buildJsonObjs(n int) []*JsonObj {
	objs := make([]*JsonObj, n)

	for i := 0; i < n; i++ {
		obj := &JsonObj{
			Id:       i,
			Val:      fmt.Sprintf("name_%d", i),
			MapVal:   map[string]string{"1": "a", "2": "b", "3": "c"},
			SliceVal: append(make([]string, 0), "abc"),
		}

		objs = append(objs, obj)
	}

	return objs
}
