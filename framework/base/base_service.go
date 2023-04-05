package base

import (
	dapr "github.com/dapr/go-sdk/client"
	"gitlab.musadisca-games.com/wangxw/musae/framework/guid"
	"sync"
)

const (
	Actor2GateOnRpc = 1
	Actor2GateOnCh  = 2
	MAX_TIMER_SIZE  = 1024 * 256
)

type TimerEventCB func() error

type BaseService struct {
	Daprc     dapr.Client
	GuidPool  *guid.GUIDPool
	TimerId   uint64
	TimerMap  *sync.Map
	TimeCh    chan TimerEventCB
	ExitCh    chan struct{}
	IsMetric  bool
	DataDir   string    //excel配置路径
	LogDir    string    //日志目录
	PProfAddr string    //pprof server
	Gateway   string    //服务网关地址
	CfgKeys   *sync.Map //配置中心
}
