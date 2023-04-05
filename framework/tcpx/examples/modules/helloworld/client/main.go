package main

import (
	"fmt"
	"gitlab.musadisca-games.com/wangxw/musae/framework/tcpx"
	"net"

	//"tcpx"
)

func main() {
	conn, e := net.Dial("tcp", "localhost:8080")

	if e != nil {
		panic(e)
	}
	var message = []byte("hello world")
	buf, e := tcpx.PackWithMarshaller(tcpx.Message{
		MessageID: 1,
		Header:    nil,
		Body:      message,
	}, nil)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	_, e = conn.Write(buf)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
}
