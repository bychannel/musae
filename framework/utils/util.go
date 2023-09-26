package utils

import (
	"encoding/binary"
	"encoding/json"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"os"
	"path/filepath"
	"runtime"
)

func PackRpcMsg(cmd uint32, data []byte) []byte {
	var messageIDBuf = make([]byte, 4)
	binary.BigEndian.PutUint32(messageIDBuf, cmd)
	var content = make([]byte, 0, len(data)+4)
	content = append(content, messageIDBuf...)
	content = append(content, data...)
	return content
}

func UnPackRpcMsg(data []byte) (uint32, []byte) {
	cmd := binary.BigEndian.Uint32(data[0:4])
	return cmd, data[4:]
}

func PathExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// PrettyJson pretty json for log
func PrettyJson(message interface{}) string {
	if !baseconf.GetBaseConf().IsDebug {
		return ""
	}
	var b []byte
	var err error
	if runtime.GOOS == "windows" {
		b, err = json.MarshalIndent(message, "", "\t")
		if err != nil {
			return err.Error()
		}
	} else {
		b, err = json.Marshal(message)
		if err != nil {
			return err.Error()
		}
	}
	return string(b)
}

// PrettyJson pretty json for log
func PrettyJsonLimit(message interface{}) string {
	var b []byte
	var err error
	if runtime.GOOS == "windows" {
		b, err = json.MarshalIndent(message, "", "\t")
		if err != nil {
			return err.Error()
		}
	} else {
		b, err = json.Marshal(message)
		if err != nil {
			return err.Error()
		}
	}
	str := string(b)
	if !baseconf.GetBaseConf().IsDebug && len(str) > 1024 {
		return string([]rune(str)[:1024]) + "...LogTruncation"
	}
	return str
}

func GetProcName() string {
	_, fileName := filepath.Split(os.Args[0])
	return fileName
}

// IsExist 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
