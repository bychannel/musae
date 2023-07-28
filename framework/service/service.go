package service

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/actor/config"
	"github.com/dapr/go-sdk/actor/runtime"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gitlab.musadisca-games.com/wangxw/musae/framework/base"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/metrics"
	"gitlab.musadisca-games.com/wangxw/musae/framework/tcpx"
	"gitlab.musadisca-games.com/wangxw/musae/framework/web"
	"sync"
	"time"
)

const (
	LOCK_TTL_SEC   = 5
	PUBSUB_TTL_SEC = 5
)

type HTTP_METHOD string
type EVENT_TYPE string

const (
	HTTP_POST HTTP_METHOD = "POST"
	HTTP_GET  HTTP_METHOD = "GET"
)

const (
	//pubsub event type
	EVENT_PRIVATE EVENT_TYPE = "topic-private"
	EVENT_APPID   EVENT_TYPE = "topic-appid"
	EVENT_GLOBAL  EVENT_TYPE = "topic-global"

	//global topic
	GLOBAL_TOPIC = "global"
)

type Service struct {
	base.BaseService
	AppId         string    //服务类型ID, 类型唯一,同一类型服务可以有多个实例
	InAddr        string    //服务端口
	OutAddr       string    //用户端口
	WebAddr       string    //web端口
	GRPCPort      string    //grpc端口
	Actors        []string  //Actor 类型, 在Actor服务器进程上启动
	ConfFile      string    //配置文件
	HasPriTopic   bool      //是否订阅私有主题
	PrivateTopic  string    //私有主题
	OnlinePlayers *sync.Map //在线用户
	svc           common.Service
	tcp           *tcpx.TcpX
	http          *web.HttpServer
	Redis         *redis.Client
	RedisCluster  *redis.ClusterClient
	RedisCenter   *redis.Client
	ES            *elasticsearch.TypedClient
	state         base.PState

	OnPreInit        base.FPreInit        //服务框架初始化之前调用,rpc,http 注册,dapr client 可用
	OnServerInit     base.FServerInit     //服务框架初始化之后调用,dapr client 可用
	OnConnect        base.FNetConnect     //网络连接事件
	OnMessage        base.FNetMessage     //消息事件
	OnClose          base.FNetClose       //网络关闭事件
	OnHeartBeat      base.FNetHeartBeat   //网络关闭事件
	OnEventHandler   base.FEventHandler   //订阅事件
	OnInvokeHandler  base.FInvokeHandler  //服务调用
	OnBindHandler    base.FBindingHandler //输入流事件
	ActorFactory     []base.FActorFactory //Actor微服务创建工厂
	OnRegisterMetric metrics.RegisterMetricFunc
	OnCfgCenterCB    dapr.ConfigurationHandleFunction
	OnCronEveryHour  base.OnCronEveryHour
}

func (s *Service) RegisterRpcHandler(name string, fn common.ServiceInvocationHandler) error {
	if err := s.svc.AddServiceInvocationHandler(name, fn); err != nil {
		logger.Errorf("RegisterRpcHandler [%s] error: %v", name, err)
		return err
	}
	return nil
}

func (s *Service) RegisterBindingInvocationHandler(name string, fn common.BindingInvocationHandler) error {
	if err := s.svc.AddBindingInvocationHandler(name, fn); err != nil {
		logger.Errorf("RegisterBindingInvocationHandler [%s] error: %v", name, err)
		return err
	}
	return nil
}

/*
method for GET, POST
*/
func (s *Service) RegisterHttpHandler(method, relativePath string, fn gin.HandlerFunc) error {
	s.http.RegisterHandler(method, relativePath, fn)
	return nil
}

func (s *Service) onCronEveryHour(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	logger.Debugf("Service binding onCronEveryHour - Data:%s, Meta:%v", in.Data, in.Metadata)
	return nil, nil
}

func (s *Service) ImpActorStub(actorStub actor.Client, opt ...config.Option) {
	s.Daprc.ImplActorClientStub(actorStub, opt...)
}

func (s *Service) GetActors() []string {
	return s.Actors
}

func (s *Service) RegisterActorFactory(fn ...base.FActorFactory) {
	for _, factory := range fn {
		s.ActorFactory = append(s.ActorFactory, factory)
	}
}

func (s *Service) GracefulStop() {
	if len(s.ActorFactory) > 0 {
		for _, actor := range s.Actors {
			runtime.GetActorRuntimeInstance().KillAllActors(actor)
		}
	}
	if err := s.svc.GracefulStop(); err != nil {
		logger.Error("[service] GracefulStop error:", err)
	}
	go func() {
		time.Sleep(6 * time.Second)
		szLog := fmt.Sprintf("ServerStop success appid:%s version:%s rolling:%s", global.AppID, global.APP_VERSION, global.ROLLING_VERSION)
		logger.Info(szLog)
		if baseconf.GetBaseConf().IsDebug && !global.IsDev &&
			len(baseconf.GetBaseConf().FeishuNotifyRobot) > 0 {
			logger.PushLog2Chat(baseconf.GetBaseConf().FeishuNotifyRobot, "ServerStop", szLog)
		}
		s.ExitCh <- struct{}{}
	}()
}

func (s *Service) State() base.PState {
	return s.state
}

func (s *Service) SetState(state base.PState) {
	s.state = state
}

func (s *Service) GetAppID() string {
	return s.AppId
}

// add an input binding invocation handler
func (s *Service) bindingCronHandler() error {
	// every 10s
	var (
		handler common.BindingInvocationHandler
	)

	// every hour
	if s.OnCronEveryHour == nil {
		handler = s.onCronEveryHour
	} else {
		handler = common.BindingInvocationHandler(s.OnCronEveryHour)
	}
	if err := s.svc.AddBindingInvocationHandler("/cron-hour", handler); err != nil {
		logger.Fatalf("error adding binding handler: %v", err)
	}

	return nil
}
