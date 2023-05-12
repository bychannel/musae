package base

import (
	"context"
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/service/common"
	"gitlab.musadisca-games.com/wangxw/musae/framework/tcpx"
)

type FPreInit = func() error
type FServerInit = func() error
type FNetConnect = func(c *tcpx.Context)
type FNetMessage = func(c *tcpx.Context)
type FNetClose = func(c *tcpx.Context)
type FNetHeartBeat = func(c *tcpx.Context)
type FEventHandler = func(context.Context, *common.TopicEvent) (retry bool, err error)
type FInvokeHandler = func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
type FBindingHandler = func(ctx context.Context, in *common.BindingEvent) (out []byte, err error)
type FActorFactory = func() actor.Server
type FProcessOption = func(server IServer) error

type OnCronEveryHour func(context.Context, *common.BindingEvent) (out []byte, err error)

type IServer interface {
	Init() error
	Start() error
	Main()
	Exit()
	GracefulStop()
	Reload() error
	State() PState
	SetState(state PState)
	GetAppID() string
	GetActors() []string
	RegisterActorFactory(fn ...FActorFactory)
}

type IProcess interface {
	Start(opts ...FProcessOption)
	Exit()
	Status() PState
}
