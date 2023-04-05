package main

import (
	//"tcpx"
	"gitlab.musadisca-games.com/wangxw/musae/framework/tcpx"
	"time"
)

func main() {
	srv := tcpx.NewTcpX(nil)

	srv.HeartBeatModeDetail(true, 10*time.Second, false, tcpx.DEFAULT_HEARTBEAT_MESSAGEID)

	//srv.RewriteHeartBeatHandler(1300, func(c *tcpx.Context) {
	//	fmt.Println("rewrite heartbeat handler")
	//	c.RecvHeartBeat()
	//})

	tcpx.SetLogMode(tcpx.DEBUG)

	srv.ListenAndServe("tcp", ":8101")
}
