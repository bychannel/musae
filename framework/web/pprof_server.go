package web

import (
	"github.com/arl/statsviz"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func PProfServerStart(addr string) {
	go func() {
		runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪，block
		runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪，mutex
		statsviz.RegisterDefault()
		logger.Info(http.ListenAndServe(addr, nil))
	}()
	logger.Infof("pprof server start on %s", addr)
}
