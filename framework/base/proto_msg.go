package base

import (
	"context"
	"fmt"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils"
	"google.golang.org/protobuf/proto"
)

type IProtobufHandler interface {
	GetID() int32
	Handler(ctx context.Context, req interface{}) (interface{}, error)
}

/*type ProtoMsg struct {
	AppId   string `json:"appId"`   //message src appid
	MsgId   int32  `json:"msgId"`   //message id
	UserId  string `json:"userId"`  //user id
	RoleId  uint64 `json:"roleId"`  //user id
	UAID    string `json:"uaid"`    //UserActor id
	Data    []byte `json:"data"`    // proto Marshal data
	ErrCode int32  `json:"errCode"` // err code
	GUID    uint32 `json:"guid"`    // guid
	Topic   string `json:"-"`       //message pubsub Topic
}

func (p *ProtoMsg) String() string {
	return fmt.Sprintf("ProtoMsg:{AppId:%v,MsgId:%v,UserId:%v,RoleId:%v,UAID:%v,Data:%v,ErrCode:%v}", p.AppId, p.MsgId, p.UserId, p.RoleId, p.UAID, len(p.Data), p.ErrCode)
}

func (p *ProtoMsg) Unmarshal(m interface{}) error {
	return proto.Unmarshal(p.Data, m.(proto.Message))
}*/

func (p *ProtoMsg) Str() string {
	return fmt.Sprintf("ProtoMsg:{AppId:%v,MsgId:%v,Topic:%v,UserId:%v,RoleId:%v,UAID:%v,Data:%v,ErrCode:%v,ReqIdx:%v,Uids:%+v}",
		p.AppId, p.MsgId, p.Topic, p.UserId, p.RoleId, p.UAID, len(p.Data), p.ErrCode, p.ReqIdx, p.Uids)
}

func (p *ProtoMsg) UnmarshalData(m proto.Message) error {
	err := proto.Unmarshal(p.Data, m)
	if err != nil {
		return fmt.Errorf("in:%s, req:%T{%+v}, err:%+v", p, m, m, err)
	}
	logger.Infof("\n===>>>MSG-UP msg:[%T], {%s}, {%s}\n", m, utils.PrettyJsonLimit(m), p.Str())
	return err
}

func (p *ProtoMsg) Marshal() ([]byte, error) {
	return proto.Marshal(p)
}

type RpcError struct {
	Err  error
	Code int32
}

func (e *RpcError) Error() string {
	return fmt.Sprintf("RpcError:{code:%v,error:%v}", e.Code, e.Err)
}

type FProtoMsgHandler = func(ctx context.Context, in *ProtoMsg) (proto.Message, error, int32)

func PackProtoMsg(cmd int32, uid string, roleId uint64, uaid string, data []byte, appId string, uids []string) ([]byte, error) {
	msg := &ProtoMsg{MsgId: cmd, UserId: uid, RoleId: roleId, UAID: uaid, Data: data, AppId: appId, Uids: uids}
	return proto.Marshal(msg)
}

func UnPackProtoMsg(data []byte) (*ProtoMsg, error) {
	msg := &ProtoMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// ParserReqBytes 解析请求数据接口
func UnmarshalData(bytes []byte, m proto.Message) error {
	var (
		err error
	)

	//msgId, uid, reqData = in.MsgId, in.UserId, in.Data
	err = proto.Unmarshal(bytes, m)
	if err != nil {
		return fmt.Errorf("解析请求参数出错: req:%T{%v}, err:%+v", m, m, err)
	}
	logger.Infof("\n===>>>MSG-UP, msg:[%T], {%s}\n", m, utils.PrettyJsonLimit(m))

	return nil
}
