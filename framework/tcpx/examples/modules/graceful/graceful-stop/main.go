package main

import (
	"fmt"
	"gitlab.musadisca-games.com/wangxw/musae/framework/errorx"
	"gitlab.musadisca-games.com/wangxw/musae/framework/tcpx"
	//"tcpx"
	"time"
)

func main() {
	srv := tcpx.NewTcpX(nil)

	// start server
	go func() {
		fmt.Println("tcp listen on :8080")
		srv.ListenAndServe("tcp", ":8080")
	}()

	// after 10 seconds and stop it
	go func() {
		time.Sleep(10 * time.Second)
		if e := srv.Stop(false); e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		//
		//if e:=srv.Stop(true); e!=nil {
		//	fmt.Println(errorx.Wrap(e).Error())
		//	return
		//}
	}()

	select {}
}
