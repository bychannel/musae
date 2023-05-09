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
	Cache2Redis(mongoDbType service.MongoDbType, key string, value proto.Message) error
	SaveMongoDB(mongoDbName service.MongoDbType, key string, value proto.Message) error
}

type BaseActor struct {
	actor.ServerImplBase
	base.BaseService
	ActorLogger

	MsgFunc    map[int32]base.FProtoMsgHandler
	RpcMethods map[string]*RpcMethod
}

func (s *BaseActor) RegisterProtoHandler(messageId int32, handler base.FProtoMsgHandler) {
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
