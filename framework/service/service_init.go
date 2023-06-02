package service

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/actor/config"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/dapr/go-sdk/service/http"
	"github.com/go-redis/redis/v8"
	"gitlab.musadisca-games.com/wangxw/musae/framework/base"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/dlog"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"gitlab.musadisca-games.com/wangxw/musae/framework/guid"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/metrics"
	"gitlab.musadisca-games.com/wangxw/musae/framework/tcpx"
	"gitlab.musadisca-games.com/wangxw/musae/framework/web"
	"math/rand"

	"sync"
	"time"
)

func (s *Service) String() string {
	return s.AppId + "," + "," + s.InAddr + "," + s.OutAddr + "," + s.GRPCPort
}

func (s *Service) PrivateTopicID() string {
	return fmt.Sprintf("%s:%d", s.AppId, global.SID)
}

func (s *Service) InitLog() error {
	fmt.Println("InitLog dir:", s.LogDir)
	var fileName string
	if global.IsCloud {
		//fileName = fmt.Sprintf("%s-%v-%v", s.AppId, global.SID, time.Now().Format("2006-01-02-15-04-05"))
		fileName = fmt.Sprintf("%s-%v-%v", s.AppId, global.SID, time.Now().Format("2006-01-02"))
	} else {
		fileName = fmt.Sprintf("%s-%s", s.AppId, time.Now().Format("2006-01-02"))
	}
	if err := logger.Init(s.LogDir+"/plog", fileName); err != nil {
		return err
	}
	if err := dlog.Init(s.LogDir+"/dlog", fileName, 1); err != nil {
		return err
	}
	if s.IsMetric {
		if err := metrics.InitMetric(s.LogDir+"/mlog", fileName, []string{"ns|" + metrics.NameSpace, "appId|" + s.AppId, "rollingVersion|" + global.ROLLING_VERSION}, s.OnRegisterMetric, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) InitBase() error {
	logger.Info("[service], initBase begin")
	s.TimerMap = &sync.Map{}
	s.CfgKeys = &sync.Map{}
	s.Players = &sync.Map{}
	s.TimeCh = make(chan base.TimerEventCB, base.MAX_TIMER_SIZE)
	s.ExitCh = make(chan struct{}, 1)
	s.GuidPool = guid.NewGUIDPool(s.DBNext)

	//设置随机种子
	rand.Seed(time.Now().UnixNano())

	//pprof server
	if s.PProfAddr != "" {
		web.PProfServerStart(s.PProfAddr)
	}

	if err := s.initSvc(); err != nil {
		logger.Fatal("[service], initHttpSvc error: %v", err)
		return err
	}

	logger.Debug("[service], initSvc success")
	// 开放外网端口 initNet
	if s.OutAddr != "" {
		if err := s.initNet(); err != nil {
			logger.Fatal("[service], initNet error: %v", err)
			return err
		}
	}

	// web server
	if s.WebAddr != "" {
		if err := s.initWeb(); err != nil {
			logger.Fatal("[service], initWeb error: %v", err)
			return err
		}
		logger.Debug("[service] initWeb success")
	}

	if err := s.initRedis(); err != nil {
		logger.Fatal("[service], initRedis error: %v", err)
		return err
	}

	if err := s.initES(); err != nil {
		logger.Fatal("[service], initES error: %v", err)
		return err
	}
	logger.Debug("[service] initRedis success")

	logger.Info("[service] initNet success")
	return nil
}

func (s *Service) initNet() error {
	s.tcp = tcpx.NewTcpX(tcpx.ProtobufMarshaller{})
	s.tcp.WithBuiltInPool(true)

	s.tcp.HeartBeatMode(false, 100*time.Second)
	s.tcp.SetMaxBytePerMessage(int32(baseconf.GetBaseConf().SrvMsgMaxSize))
	//s.net.HeartBeatModeDetail(true, 5*time.Second, false, tcpx.DEFAULT_HEARTBEAT_MESSAGEID)

	s.tcp.OnClose = s.OnClose
	s.tcp.OnMessage = s.OnMessage
	s.tcp.OnConnect = s.OnConnect

	tcpx.SetLogMode(tcpx.DEBUG)
	//s.net.RewriteHeartBeatHandler(15, s.OnHeartBeat)
	//fmt.Println("rewrite heartbeat and receive from client")
	//})

	return nil
}

func (s *Service) initWeb() error {
	s.http = web.NewHttpServer()
	return s.http.Init(s.WebAddr)
}

func (s *Service) initRedis() error {
	var opts *redis.Options
	var clusterOpts *redis.ClusterOptions
	if global.IsCloud {
		clusterOpts = &redis.ClusterOptions{
			Addrs: []string{
				baseconf.GetBaseConf().RedisConf.Addr,
			},
			Password:        baseconf.GetBaseConf().RedisConf.Password,
			DialTimeout:     100 * time.Millisecond,
			ReadTimeout:     100 * time.Millisecond,
			WriteTimeout:    100 * time.Millisecond,
			MaxRetries:      baseconf.GetBaseConf().RedisConf.MaxRetries,
			MinRetryBackoff: time.Duration(baseconf.GetBaseConf().RedisConf.MinRetryBackoff) * time.Millisecond,
			MaxRetryBackoff: time.Duration(baseconf.GetBaseConf().RedisConf.MaxRetryBackoff) * time.Millisecond,
		}
		s.RedisCluster = redis.NewClusterClient(clusterOpts)
	} else {
		opts = &redis.Options{
			Addr:            baseconf.GetBaseConf().RedisConf.AddrDev,
			Password:        baseconf.GetBaseConf().RedisConf.Password,
			DialTimeout:     100 * time.Millisecond,
			ReadTimeout:     100 * time.Millisecond,
			WriteTimeout:    100 * time.Millisecond,
			MaxRetries:      baseconf.GetBaseConf().RedisConf.MaxRetries,
			MinRetryBackoff: time.Duration(baseconf.GetBaseConf().RedisConf.MinRetryBackoff) * time.Millisecond,
			MaxRetryBackoff: time.Duration(baseconf.GetBaseConf().RedisConf.MaxRetryBackoff) * time.Millisecond,
		}
		s.Redis = redis.NewClient(opts)
		logger.RedisCli = s.Redis
	}

	logger.Infof("redis Init success,options:%+v", opts)
	return nil
}

func (s *Service) initSvc() error {
	protocol := baseconf.GetBaseConf().Protocol
	if s.AppId == "actor" {
		protocol = "http"
	}
	if protocol == "http" {
		s.svc = http.NewService(s.InAddr)
	} else {
		var err error
		s.svc, err = grpc.NewService(s.InAddr)
		if err != nil {
			logger.Fatal("error NewService, err: ", err)
		}
	}

	if err := s.initsubpub(); err != nil {
		logger.Fatal("error initsubpub, err: ", err)
	}

	// add a service to service invocation handler
	if err := s.svc.AddServiceInvocationHandler("RpcCall", s.OnInvokeHandler); err != nil {
		logger.Fatal("error adding invocation handler: %v", err)
		return err
	}

	// add a binding invocation handler
	if err := s.svc.AddBindingInvocationHandler("Binding", s.OnBindHandler); err != nil {
		logger.Fatal("error adding binding handler: %v", err)
		return err
	}

	// add an input binding invocation handler
	if err := s.bindingCronHandler(); err != nil {
		return err
	}

	// add ready handler
	if err := s.svc.AddHealthCheckHandler("/api/ready", func(ctx context.Context) (err error) {
		logger.Debugf("/api/ready, server state: %v", s.state)
		if s.state == base.PState_Running {
			return nil
		}
		return fmt.Errorf("app is unready")
	}); err != nil {
		return err
	}
	// add healthz handler
	if err := s.svc.AddHealthCheckHandler("/api/healthz", func(ctx context.Context) (err error) {
		logger.Debugf("/api/healthz, server state: %v", s.state)
		return nil
	}); err != nil {
		return err
	}
	logger.Debugf("s.ActorFactory, %v", s.ActorFactory)

	for _, factory := range s.ActorFactory {
		s.svc.RegisterActorImplFactory(factory,
			config.WithActorIdleTimeout(fmt.Sprintf("%ss", baseconf.GetBaseConf().UserActorGCTime)), //600s
			config.WithActorScanInterval(baseconf.GetBaseConf().UserActorGCScan),                    //10s
			config.WithDrainOngingCallTimeout(baseconf.GetBaseConf().UserActorGCScan),               //10s
			config.WithDrainBalancedActors(true))
		logger.Info("[service], init ActorFactory success")
	}

	return nil
}

/*func (s *Service) initGrpcSvc() error {

	var err error
	// create a Dapr service server
	s.svc, err = grpc.NewService(s.InAddr)

	if err != nil {
		logger.Fatal("error NewService, err: ", err)
	}

	if err := s.initsubpub(); err != nil {
		logger.Fatal("error initsubpub, err: ", err)
	}

	// add a service to service invocation handler
	if err := s.svc.AddServiceInvocationHandler("RpcCall", s.OnInvokeHandler); err != nil {
		logger.Fatal("error adding invocation handler: %v", err)
		return err
	}

	// add a binding invocation handler
	if err := s.svc.AddBindingInvocationHandler("Binding", s.OnBindHandler); err != nil {
		logger.Fatal("error adding binding handler: %v", err)
		return err
	}

	// add an input binding invocation handler
	if err := s.bindingCronHandler(); err != nil {
		return err
	}

	// add ready handler
	if err := s.svc.AddHealthCheckHandler("/ready", func(ctx context.Context) (err error) {
		logger.Debugf("/ready, server state: %v", s.state)
		if s.state == base.PState_Running {
			return nil
		}
		return fmt.Errorf("app is unready")
	}); err != nil {
		return err
	}
	// add healthz handler
	if err := s.svc.AddHealthCheckHandler("/healthz", func(ctx context.Context) (err error) {
		logger.Debugf("/healthz, server state: %v", s.state)
		return nil
	}); err != nil {
		return err
	}

	if s.ActorFactory != nil {
		s.svc.RegisterActorImplFactory(s.ActorFactory,
			config.WithActorIdleTimeout(fmt.Sprintf("%ss", baseconf.GetBaseConf().UserActorGCTime)), //600s
			config.WithActorScanInterval(baseconf.GetBaseConf().UserActorGCScan),                    //10s
			config.WithDrainOngingCallTimeout(baseconf.GetBaseConf().UserActorGCScan),               //10s
			config.WithDrainBalancedActors(true),
		)
		logger.Info("[service], init ActorFactory success")
	}
	return nil
}*/

func (s *Service) initsubpub() error {
	if s.HasPriTopic { //临时用 rand替换GUID
		s.PrivateTopic = s.PrivateTopicID()
		priSub := &common.Subscription{
			PubsubName: "topic-private",
			Topic:      s.PrivateTopic,
			Route:      "/event",
		}
		if err := s.svc.AddTopicEventHandler(priSub, s.OnEventHandler); err != nil {
			logger.Fatal("error adding topic subscription: %v, %v", err, priSub)
			return err
		}
		logger.Infof("Subscript private Topic Event: %+v", priSub)
	}

	{
		// add appid topic subscriptions
		appidSub := &common.Subscription{
			PubsubName: "topic-appid",
			Topic:      s.AppId,
			Route:      "/event",
		}

		if err := s.svc.AddTopicEventHandler(appidSub, s.OnEventHandler); err != nil {
			logger.Fatal("error adding topic subscription: %v, %v", err, appidSub)
			return err
		}
		logger.Infof("Subscript Appid Topic Event: %+v", appidSub)

	}

	{
		// add golbal topic subscriptions
		global := &common.Subscription{
			PubsubName: "topic-global",
			Topic:      GLOBAL_TOPIC,
			Route:      "/event",
		}

		if err := s.svc.AddTopicEventHandler(global, s.OnEventHandler); err != nil {
			logger.Fatal("error adding topic subscription: %v, %v", err, global)
			return err
		}
		logger.Infof("Subscript Global Topic Event: %+v", global)
	}

	{
		// add deadletter topic subscriptions
		deadLetter := &common.Subscription{
			PubsubName: "pubsub",
			Topic:      "orders",
			Route:      "/checkout",
		}

		if err := s.svc.AddTopicEventHandler(deadLetter, s.DeadLetterTopic); err != nil {
			logger.Fatal("error adding topic subscription: %v, %v", err, deadLetter)
			return err
		}
		logger.Infof("Subscript deadLetter Topic Event: %+v", deadLetter)
	}
	return nil
}

func (s *Service) DeadLetterTopic(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
	logger.Warnf("DeadLetterTopic Event: %v, %v", event.Topic, event.ID)
	return false, err
}
