package state

import (
	"fmt"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
)

// KvTable
type KvTable struct {
	Key     string `json:"key"` //关联key
	Id      uint64 `json:"id"`  //关联ID：user,battle
	UID     string `json:"uid"` //关联UID：账号 id
	Data    []byte `json:"data"`
	UpSecTS int64  `json:"update_ts"`
	InSecTS int64  `json:"insert_ts"`
	DataSrc string `json:"data_src"` //data 原数据
}

func (d *KvTable) Str() string {
	if baseconf.GetBaseConf().IsDebug {
		return fmt.Sprintf("KvTable:{Key:%v, Id:%v, UID:%v, DataLen:%v, Data:%v, UpSecTS:%v, InSecTS:%v}", d.Key, d.Id, d.UID, len(d.Data), d.DataSrc, d.UpSecTS, d.InSecTS)
	} else {
		return fmt.Sprintf("KvTable:{Key:%v, Id:%v, UID:%v, DataLen:%v, Data:%v, UpSecTS:%v, InSecTS:%v}", d.Key, d.Id, d.UID, len(d.Data), "", d.UpSecTS, d.InSecTS)
	}
}
