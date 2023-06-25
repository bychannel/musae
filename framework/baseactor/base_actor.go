package baseactor

import (
	"github.com/dapr/go-sdk/actor"
	"gitlab.musadisca-games.com/wangxw/musae/framework/base"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/service"
	"gitlab.musadisca-games.com/wangxw/musae/framework/state"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type RpcMethod struct {
	Handler interface{}
	Method  grpc.MethodDesc
}

type IBaseActor interface {
	actor.Server
	GetCache(mongoDbName service.MongoDbType, key string, msg proto.Message) (*state.KvTable, error)
	Cache2Redis(mongoDbType service.MongoDbType, uaid string, key string, value proto.Message) error
	SaveMongoDB(mongoDbName service.MongoDbType, key string, value proto.Message) error
}

type BaseActor struct {
	actor.ServerImplBase
	base.BaseService
	ActorLogger

	// 延迟落库
	HandlersMap map[service.MongoDbType][]IBaseHandler

	MsgFunc    map[int32]base.FProtoMsgHandler
	RpcMethods map[string]*RpcMethod
	ActorType  string
}

func (s *BaseActor) RegisterProtoHandler(messageId int32, handler base.FProtoMsgHandler) {
	if s.MsgFunc == nil {
		s.MsgFunc = make(map[int32]base.FProtoMsgHandler)
	}

	if _, ok := s.MsgFunc[messageId]; !ok {
		s.MsgFunc[messageId] = handler
		logger.Debugf("register messageId: %d", messageId)
	} else if ok {
		logger.Errorf("Duplicate messageId are registered: %d", messageId)
	}
}

func (s *BaseActor) RegisterRpcMethod(h interface{}, service *grpc.ServiceDesc) {
	service.HandlerType = h
	if service.HandlerType == nil {
		logger.Errorf("register message,invalid handler: %s ", service.ServiceName)
	}
	for _, v := range service.Methods {
		if _, ok := s.RpcMethods[v.MethodName]; !ok {
			s.RpcMethods[v.MethodName] = &RpcMethod{Handler: h, Method: v}
			logger.Debugf("register  message, %s-%s", service.ServiceName, v.MethodName)
		} else if ok {
			logger.Errorf("duplicate message, %s-%s, metadata %s", service.ServiceName, v.MethodName, service.Metadata)
		}
	}

}

func (s *BaseActor) KeepHandler(ib IBaseHandler) {
	dbType, _, _ := ib.DBTable()
	if _, ok := s.HandlersMap[dbType]; !ok {
		s.HandlersMap[dbType] = make([]IBaseHandler, 0)
	}

	s.HandlersMap[dbType] = append(s.HandlersMap[dbType], ib)
}

func (s *BaseActor) Type() string {
	return s.ActorType
}
